package storyboard

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/StevenWeathers/thunderdome-planning-poker/internal/db"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
	"go.uber.org/zap"
)

// GetStoryboardUsers retrieves the users for a given storyboard from db
func (d *Service) GetStoryboardUsers(storyboardID string) []*thunderdome.StoryboardUser {
	var users = make([]*thunderdome.StoryboardUser, 0)
	rows, err := d.DB.Query(
		`SELECT
			w.id, w.name, su.active, w.avatar, COALESCE(w.email, ''), COALESCE(w.picture, '')
		FROM thunderdome.storyboard_user su
		LEFT JOIN thunderdome.users w ON su.user_id = w.id
		WHERE su.storyboard_id = $1
		ORDER BY w.name;`,
		storyboardID,
	)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var su thunderdome.StoryboardUser
			if err := rows.Scan(&su.ID, &su.Name, &su.Active, &su.Avatar, &su.GravatarHash, &su.PictureURL); err != nil {
				d.Logger.Error("get_storyboard_users query scan error", zap.Error(err))
			} else {
				if su.GravatarHash != "" {
					su.GravatarHash = db.CreateGravatarHash(su.GravatarHash)
				} else {
					su.GravatarHash = db.CreateGravatarHash(su.ID)
				}
				users = append(users, &su)
			}
		}
	}

	return users
}

// AddUserToStoryboard adds a user by ID to the storyboard by ID
func (d *Service) AddUserToStoryboard(storyboardID string, userID string) ([]*thunderdome.StoryboardUser, error) {
	if _, err := d.DB.Exec(
		`INSERT INTO thunderdome.storyboard_user (storyboard_id, user_id, active)
		VALUES ($1, $2, true)
		ON CONFLICT (storyboard_id, user_id) DO UPDATE SET active = true, abandoned = false`,
		storyboardID,
		userID,
	); err != nil {
		d.Logger.Error("insert storybaord user error", zap.Error(err))
	}

	users := d.GetStoryboardUsers(storyboardID)

	return users, nil
}

// RetreatStoryboardUser removes a user from the current storyboard by ID
func (d *Service) RetreatStoryboardUser(storyboardID string, userID string) []*thunderdome.StoryboardUser {
	if _, err := d.DB.Exec(
		`UPDATE thunderdome.storyboard_user SET active = false WHERE storyboard_id = $1 AND user_id = $2`, storyboardID, userID); err != nil {
		d.Logger.Error("set storyboard user active false error", zap.Error(err))
	}

	if _, err := d.DB.Exec(
		`UPDATE thunderdome.users SET last_active = NOW() WHERE id = $1`, userID); err != nil {
		d.Logger.Error("set user last active error", zap.Error(err))
	}

	users := d.GetStoryboardUsers(storyboardID)

	return users
}

// GetStoryboardUserActiveStatus checks storyboard active status of User for given storyboard
func (d *Service) GetStoryboardUserActiveStatus(storyboardID string, userID string) error {
	var active bool

	err := d.DB.QueryRow(`
		SELECT coalesce(active, FALSE)
		FROM thunderdome.storyboard_user
		WHERE user_id = $2 AND storyboard_id = $1;`,
		storyboardID,
		userID,
	).Scan(
		&active,
	)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("get storyboard user active status query error: %v", err)
	} else if err != nil && errors.Is(err, sql.ErrNoRows) {
		return err
	}

	if active {
		return errors.New("DUPLICATE_STORYBOARD_USER")
	}

	return nil
}

// AbandonStoryboard removes a user from the current storyboard by ID and sets abandoned true
func (d *Service) AbandonStoryboard(storyboardID string, userID string) ([]*thunderdome.StoryboardUser, error) {
	if _, err := d.DB.Exec(
		`UPDATE thunderdome.storyboard_user SET active = false, abandoned = true WHERE storyboard_id = $1 AND user_id = $2`, storyboardID, userID); err != nil {
		return nil, fmt.Errorf("set storyboard user active false error: %v", err)
	}

	if _, err := d.DB.Exec(
		`UPDATE thunderdome.users SET last_active = NOW() WHERE id = $1`, userID); err != nil {
		return nil, fmt.Errorf("set user last active error: %v", err)
	}

	users := d.GetStoryboardUsers(storyboardID)

	return users, nil
}
