package auth

import (
	"context"
	"database/sql"
	"errors"

	"github.com/StevenWeathers/thunderdome-planning-poker/internal/db"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"

	"go.uber.org/zap"
)

// CreateSession creates a new user authenticated session
func (d *Service) CreateSession(ctx context.Context, UserId string) (string, error) {
	SessionId, err := db.RandomBase64String(32)
	if err != nil {
		return "", err
	}

	if _, sessionErr := d.DB.ExecContext(ctx, `
		INSERT INTO thunderdome.user_session (session_id, user_id, disabled) VALUES ($1, $2, (SELECT mfa_enabled FROM thunderdome.users WHERE id = $2));
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
func (d *Service) EnableSession(ctx context.Context, SessionId string) error {
	if _, sessionErr := d.DB.ExecContext(ctx, `
		UPDATE thunderdome.user_session SET disabled = false WHERE session_id = $1;
		`,
		SessionId,
	); sessionErr != nil {
		d.Logger.Ctx(ctx).Error("Unable to enable user session", zap.Error(sessionErr))
		return sessionErr
	}

	return nil
}

// GetSessionUser gets a user session by sessionId
func (d *Service) GetSessionUser(ctx context.Context, SessionId string) (*thunderdome.User, error) {
	User := &thunderdome.User{}

	e := d.DB.QueryRowContext(ctx, `
		SELECT
        u.id,
        u.name,
        u.email,
        u.type,
        u.avatar,
        u.verified,
        u.notifications_enabled,
        COALESCE(u.country, ''),
        COALESCE(u.locale, ''),
        COALESCE(u.company, ''),
        COALESCE(u.job_title, ''),
        u.created_date,
        u.updated_date,
        u.last_active,
        u.subscribed
    FROM thunderdome.user_session us
    LEFT JOIN thunderdome.users u ON u.id = us.user_id
    WHERE us.session_id = $1 AND NOW() < us.expire_date`,
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
		&User.LastActive,
		&User.Subscribed,
	)
	if e != nil {
		if !errors.Is(e, sql.ErrNoRows) {
			d.Logger.Ctx(ctx).Error("user_session_get query error", zap.Error(e))
		}
		return nil, errors.New("active session match not found")
	}

	User.GravatarHash = db.CreateGravatarHash(User.Email)

	return User, nil
}

// DeleteSession deletes a user authenticated session
func (d *Service) DeleteSession(ctx context.Context, SessionId string) error {
	if _, sessionErr := d.DB.ExecContext(ctx, `
		DELETE FROM thunderdome.user_session WHERE session_id = $1;
		`,
		SessionId,
	); sessionErr != nil {
		d.Logger.Ctx(ctx).Error("Unable to delete user session", zap.Error(sessionErr))
		return sessionErr
	}

	return nil
}
