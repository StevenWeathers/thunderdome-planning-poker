package db

import (
	"errors"

	"github.com/StevenWeathers/thunderdome-planning-poker/model"
	"go.uber.org/zap"
)

// CreateSession creates a new user authenticated session
func (d *Database) CreateSession(UserId string) (string, error) {
	SessionId, err := randomBase64String(32)
	if err != nil {
		return "", err
	}

	if _, sessionErr := d.db.Exec(`
		INSERT INTO user_session (session_id, user_id, disabled) VALUES ($1, $2, (SELECT mfa_enabled FROM users WHERE id = $2));
		`,
		SessionId,
		UserId,
	); sessionErr != nil {
		d.logger.Error("Unable to create a user session", zap.Error(sessionErr))
		return "", sessionErr
	}

	return SessionId, nil
}

// EnableSession enables a user authenticated session
func (d *Database) EnableSession(SessionId string) error {
	if _, sessionErr := d.db.Exec(`
		UPDATE user_session SET disabled = false WHERE session_id = $1;
		`,
		SessionId,
	); sessionErr != nil {
		d.logger.Error("Unable to enable user session", zap.Error(sessionErr))
		return sessionErr
	}

	return nil
}

// GetSessionUser gets a user session by sessionId
func (d *Database) GetSessionUser(SessionId string) (*model.User, error) {
	User := &model.User{}

	e := d.db.QueryRow(`
		SELECT id, name, email, type, avatar, verified, notifications_enabled, country, locale, company, job_title, created_date, updated_date, last_active 
		FROM user_session_get($1);`,
		SessionId,
	).Scan(
		&User.Id,
		&User.Name,
		&User.Email,
		&User.Type,
		&User.Avatar,
		&User.Verified,
		&User.NotificationsEnabled,
		&User.Country,
		&User.Locale,
		&User.Company,
		&User.JobTitle,
		&User.CreatedDate,
		&User.UpdatedDate,
		&User.LastActive)
	if e != nil {
		d.logger.Error("user_session_get query error", zap.Error(e))
		return nil, errors.New("active session match not found")
	}

	User.GravatarHash = createGravatarHash(User.Email)

	return User, nil
}

// DeleteSession deletes a user authenticated session
func (d *Database) DeleteSession(SessionId string) error {
	if _, sessionErr := d.db.Exec(`
		DELETE FROM user_session WHERE session_id = $1;
		`,
		SessionId,
	); sessionErr != nil {
		d.logger.Error("Unable to delete user session", zap.Error(sessionErr))
		return sessionErr
	}

	return nil
}
