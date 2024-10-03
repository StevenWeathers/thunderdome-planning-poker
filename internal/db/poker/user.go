package poker

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/StevenWeathers/thunderdome-planning-poker/internal/db"
	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
	"go.uber.org/zap"
)

// GetUserActiveStatus checks game active status of User
func (d *Service) GetUserActiveStatus(pokerID string, userID string) error {
	var active bool

	err := d.DB.QueryRow(`
		SELECT coalesce(active, FALSE)
		FROM thunderdome.poker_user
		WHERE user_id = $2 AND poker_id = $1;`,
		pokerID,
		userID,
	).Scan(
		&active,
	)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("poker get user active status query error: %v", err)
	} else if err != nil && errors.Is(err, sql.ErrNoRows) {
		return err
	}

	if active {
		return errors.New("DUPLICATE_BATTLE_USER")
	}

	return nil
}

// GetUsers retrieves the users for a given game
func (d *Service) GetUsers(pokerID string) []*thunderdome.PokerUser {
	var users = make([]*thunderdome.PokerUser, 0)
	rows, err := d.DB.Query(
		`SELECT
			u.id, u.name, u.type, u.avatar, pu.active, pu.spectator, COALESCE(u.email, ''), COALESCE(u.picture, '')
		FROM thunderdome.poker_user pu
		LEFT JOIN thunderdome.users u ON pu.user_id = u.id
		WHERE pu.poker_id = $1
		ORDER BY u.name`,
		pokerID,
	)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var w thunderdome.PokerUser
			if err := rows.Scan(&w.ID, &w.Name, &w.Type, &w.Avatar, &w.Active, &w.Spectator, &w.GravatarHash, &w.PictureURL); err != nil {
				d.Logger.Error("error getting poker users", zap.Error(err))
			} else {
				if w.GravatarHash != "" {
					w.GravatarHash = db.CreateGravatarHash(w.GravatarHash)
				} else {
					w.GravatarHash = db.CreateGravatarHash(w.ID)
				}
				users = append(users, &w)
			}
		}
	}

	return users
}

// GetActiveUsers retrieves the active users for a given game
func (d *Service) GetActiveUsers(pokerID string) []*thunderdome.PokerUser {
	var users = make([]*thunderdome.PokerUser, 0)
	rows, err := d.DB.Query(
		`SELECT
			w.id, w.name, w.type, w.avatar, bw.active, bw.spectator, COALESCE(w.email, ''), COALESCE(w.picture, '')
		FROM thunderdome.poker_user bw
		LEFT JOIN thunderdome.users w ON bw.user_id = w.id
		WHERE bw.poker_id = $1 AND bw.active = true
		ORDER BY w.name`,
		pokerID,
	)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var w thunderdome.PokerUser
			if err := rows.Scan(&w.ID, &w.Name, &w.Type, &w.Avatar, &w.Active, &w.Spectator, &w.GravatarHash, &w.PictureURL); err != nil {
				d.Logger.Error("error getting active poker users", zap.Error(err))
			} else {
				if w.GravatarHash != "" {
					w.GravatarHash = db.CreateGravatarHash(w.GravatarHash)
				} else {
					w.GravatarHash = db.CreateGravatarHash(w.ID)
				}
				users = append(users, &w)
			}
		}
	}

	return users
}

// AddUser adds a user by ID to the game by ID
func (d *Service) AddUser(pokerID string, userID string) ([]*thunderdome.PokerUser, error) {
	if _, err := d.DB.Exec(
		`INSERT INTO thunderdome.poker_user (poker_id, user_id, active)
		VALUES ($1, $2, true)
		ON CONFLICT (poker_id, user_id) DO UPDATE SET active = true, abandoned = false`,
		pokerID,
		userID,
	); err != nil {
		d.Logger.Error("error adding user to poker", zap.Error(err))
	}

	users := d.GetUsers(pokerID)

	return users, nil
}

// RetreatUser removes a user from the current game by ID
func (d *Service) RetreatUser(pokerID string, userID string) []*thunderdome.PokerUser {
	if _, err := d.DB.Exec(
		`UPDATE thunderdome.poker_user SET active = false WHERE poker_id = $1 AND user_id = $2`, pokerID, userID); err != nil {
		d.Logger.Error("error updating poker user to active false", zap.Error(err))
	}

	if _, err := d.DB.Exec(
		`UPDATE thunderdome.users SET last_active = NOW() WHERE id = $1`, userID); err != nil {
		d.Logger.Error("error updating user last active timestamp", zap.Error(err))
	}

	users := d.GetUsers(pokerID)

	return users
}

// AbandonGame removes a user from the current game by ID and sets abandoned true
func (d *Service) AbandonGame(pokerID string, userID string) ([]*thunderdome.PokerUser, error) {
	if _, err := d.DB.Exec(
		`UPDATE thunderdome.poker_user SET active = false, abandoned = true WHERE poker_id = $1 AND user_id = $2`, pokerID, userID); err != nil {
		return nil, fmt.Errorf("error updating game user to abandoned: %v", err)
	}

	if _, err := d.DB.Exec(
		`UPDATE thunderdome.users SET last_active = NOW() WHERE id = $1`, userID); err != nil {
		return nil, fmt.Errorf("error updating user last active timestamp: %v", err)
	}

	users := d.GetUsers(pokerID)

	return users, nil
}

// ToggleSpectator changes a game users spectator status
func (d *Service) ToggleSpectator(pokerID string, userID string, spectator bool) ([]*thunderdome.PokerUser, error) {
	if _, err := d.DB.Exec(
		`UPDATE thunderdome.poker_user SET spectator = $3 WHERE poker_id = $1 AND user_id = $2`, pokerID, userID, spectator); err != nil {
		return nil, fmt.Errorf("poker toggle spectator query error: %v", err)
	}

	if _, err := d.DB.Exec(
		`UPDATE thunderdome.users SET last_active = NOW() WHERE id = $1`, userID); err != nil {
		d.Logger.Error("error updating user last active timestamp", zap.Error(err))
	}

	users := d.GetUsers(pokerID)

	return users, nil
}
