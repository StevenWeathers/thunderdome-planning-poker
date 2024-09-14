package retro

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/StevenWeathers/thunderdome-planning-poker/internal/db"
	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
	"go.uber.org/zap"
)

// RetroGetUsers retrieves the users for a given retro from db
func (d *Service) RetroGetUsers(RetroID string) []*thunderdome.RetroUser {
	var users = make([]*thunderdome.RetroUser, 0)
	rows, err := d.DB.Query(
		`SELECT
			u.id, u.name, su.active, u.avatar, COALESCE(u.email, ''), COALESCE(u.picture, '')
		FROM thunderdome.retro_user su
		LEFT JOIN thunderdome.users u ON su.user_id = u.id
		WHERE su.retro_id = $1
		ORDER BY u.name;`,
		RetroID,
	)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var w thunderdome.RetroUser
			if err := rows.Scan(&w.ID, &w.Name, &w.Active, &w.Avatar, &w.Email, &w.PictureURL); err != nil {
				d.Logger.Error("get retro users error", zap.Error(err))
			} else {
				if w.Email != "" {
					w.GravatarHash = db.CreateGravatarHash(w.Email)
				} else {
					w.GravatarHash = db.CreateGravatarHash(w.ID)
				}
				users = append(users, &w)
			}
		}
	}

	return users
}

// RetroAddUser adds a user by ID to the retro by ID
func (d *Service) RetroAddUser(RetroID string, UserID string) ([]*thunderdome.RetroUser, error) {
	if _, err := d.DB.Exec(
		`INSERT INTO thunderdome.retro_user (retro_id, user_id, active)
		VALUES ($1, $2, true)
		ON CONFLICT (retro_id, user_id) DO UPDATE SET active = true, abandoned = false`,
		RetroID,
		UserID,
	); err != nil {
		d.Logger.Error("insert retro user error", zap.Error(err))
	}

	users := d.RetroGetUsers(RetroID)

	return users, nil
}

// RetroRetreatUser removes a user from the current retro by ID
func (d *Service) RetroRetreatUser(RetroID string, UserID string) []*thunderdome.RetroUser {
	if _, err := d.DB.Exec(
		`UPDATE thunderdome.retro_user SET active = false WHERE retro_id = $1 AND user_id = $2`, RetroID, UserID); err != nil {
		d.Logger.Error("update retro user active false error", zap.Error(err))
	}

	if _, err := d.DB.Exec(
		`UPDATE thunderdome.users SET last_active = NOW() WHERE id = $1`, UserID); err != nil {
		d.Logger.Error("update user last active timestamp error", zap.Error(err))
	}

	users := d.RetroGetUsers(RetroID)

	return users
}

// RetroAbandon removes a user from the current retro by ID and sets abandoned true
func (d *Service) RetroAbandon(RetroID string, UserID string) ([]*thunderdome.RetroUser, error) {
	if _, err := d.DB.Exec(
		`UPDATE thunderdome.retro_user SET active = false, abandoned = true WHERE retro_id = $1 AND user_id = $2`, RetroID, UserID); err != nil {
		return nil, fmt.Errorf("abandon retro query error: %v", err)
	}

	if _, err := d.DB.Exec(
		`UPDATE thunderdome.users SET last_active = NOW() WHERE id = $1`, UserID); err != nil {
		return nil, fmt.Errorf("abandon retro update user query error: %v", err)
	}

	users := d.RetroGetUsers(RetroID)

	return users, nil
}

// GetRetroUserActiveStatus checks retro active status of User for given retro
func (d *Service) GetRetroUserActiveStatus(RetroID string, UserID string) error {
	var active bool

	err := d.DB.QueryRow(`
		SELECT coalesce(active, FALSE)
		FROM thunderdome.retro_user
		WHERE user_id = $2 AND retro_id = $1;`,
		RetroID,
		UserID,
	).Scan(
		&active,
	)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("get retro user active status query error: %v", err)
	} else if err != nil && errors.Is(err, sql.ErrNoRows) {
		return err
	}

	if active {
		return errors.New("DUPLICATE_RETRO_USER")
	}

	return nil
}

// MarkUserReady marks a user as ready for next phase
func (d *Service) MarkUserReady(RetroID string, userID string) ([]string, error) {
	var rawReadyUsers string
	readyUsers := make([]string, 0)

	err := d.DB.QueryRow(
		`UPDATE thunderdome.retro 
		SET updated_date = NOW(), ready_users = ready_users::jsonb || to_jsonb(array[$2])
		WHERE id = $1 
		RETURNING ready_users;`,
		RetroID, userID,
	).Scan(&rawReadyUsers)
	if err != nil {
		return readyUsers, fmt.Errorf("retro MarkUserReady query error: %v", err)
	}

	err = json.Unmarshal([]byte(rawReadyUsers), &readyUsers)
	if err != nil {
		d.Logger.Error("ready_users json error", zap.Error(err))
	}

	return readyUsers, nil
}

// UnmarkUserReady un-marks a user as ready for next phase
func (d *Service) UnmarkUserReady(RetroID string, userID string) ([]string, error) {
	var rawReadyUsers string
	readyUsers := make([]string, 0)

	err := d.DB.QueryRow(
		`UPDATE thunderdome.retro 
		SET updated_date = NOW(), ready_users = ready_users::jsonb - $2
		WHERE id = $1 
		RETURNING ready_users;`,
		RetroID, userID,
	).Scan(&rawReadyUsers)
	if err != nil {
		return readyUsers, fmt.Errorf("retro UnmarkUserReady query error: %v", err)
	}

	err = json.Unmarshal([]byte(rawReadyUsers), &readyUsers)
	if err != nil {
		d.Logger.Error("ready_users json error", zap.Error(err))
	}

	return readyUsers, nil
}
