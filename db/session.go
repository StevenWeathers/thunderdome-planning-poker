package db

import (
	"context"
	"database/sql"
	"errors"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"

	"go.uber.org/zap"
)

// CreateSession creates a new user authenticated session
func (d *AuthService) CreateSession(ctx context.Context, UserId string) (string, error) {
	SessionId, err := randomBase64String(32)
	if err != nil {
		return "", err
	}

	if _, sessionErr := d.DB.ExecContext(ctx, `
		INSERT INTO user_session (session_id, user_id, disabled) VALUES ($1, $2, (SELECT mfa_enabled FROM users WHERE id = $2));
		`,
		SessionId,
		UserId,
	); sessionErr != nil {
		d.Logger.Ctx(ctx).Error("Unable to create a user session", zap.Error(sessionErr))
		return "", sessionErr
	}

	return SessionId, nil
}

// EnableSession enables a user authenticated session
func (d *AuthService) EnableSession(ctx context.Context, SessionId string) error {
	if _, sessionErr := d.DB.ExecContext(ctx, `
		UPDATE user_session SET disabled = false WHERE session_id = $1;
		`,
		SessionId,
	); sessionErr != nil {
		d.Logger.Ctx(ctx).Error("Unable to enable user session", zap.Error(sessionErr))
		return sessionErr
	}

	return nil
}

// GetSessionUser gets a user session by sessionId
func (d *AuthService) GetSessionUser(ctx context.Context, SessionId string) (*thunderdome.User, error) {
	User := &thunderdome.User{}

	e := d.DB.QueryRowContext(ctx, `
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
		if !errors.Is(e, sql.ErrNoRows) {
			d.Logger.Ctx(ctx).Error("user_session_get query error", zap.Error(e))
		}
		return nil, errors.New("active session match not found")
	}

	User.GravatarHash = createGravatarHash(User.Email)

	return User, nil
}

// DeleteSession deletes a user authenticated session
func (d *AuthService) DeleteSession(ctx context.Context, SessionId string) error {
	if _, sessionErr := d.DB.ExecContext(ctx, `
		DELETE FROM user_session WHERE session_id = $1;
		`,
		SessionId,
	); sessionErr != nil {
		d.Logger.Ctx(ctx).Error("Unable to delete user session", zap.Error(sessionErr))
		return sessionErr
	}

	return nil
}
