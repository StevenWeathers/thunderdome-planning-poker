package auth

import (
	"context"
	"fmt"

	"github.com/StevenWeathers/thunderdome-planning-poker/internal/db"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
)

// CreateSession creates a new user authenticated session
func (d *Service) CreateSession(ctx context.Context, UserId string, enabled bool) (string, error) {
	SessionId, err := db.RandomBase64String(32)
	if err != nil {
		return "", err
	}

	if _, sessionErr := d.DB.ExecContext(ctx, `
		INSERT INTO thunderdome.user_session (session_id, user_id, disabled) VALUES ($1, $2, $3);
		`,
		SessionId,
		UserId,
		enabled,
	); sessionErr != nil {
		return "", fmt.Errorf("create user session query error: %v", sessionErr)
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
		return fmt.Errorf("enable user session query error: %v", sessionErr)
	}

	return nil
}

// GetSessionUser gets a user session by sessionId
func (d *Service) GetSessionUser(ctx context.Context, SessionId string) (*thunderdome.User, error) {
	User := &thunderdome.User{}

	err := d.DB.QueryRowContext(ctx, `
		SELECT
        u.id,
        u.name,
        u.email,
        u.type,
        u.avatar,
        u.notifications_enabled,
        COALESCE(u.country, ''),
        COALESCE(u.locale, ''),
        COALESCE(u.company, ''),
        COALESCE(u.job_title, ''),
        COALESCE(u.picture_url, ''),
        u.created_date,
        u.updated_date,
        u.last_active
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
		&User.NotificationsEnabled,
		&User.Country,
		&User.Locale,
		&User.Company,
		&User.JobTitle,
		&User.PictureURL,
		&User.CreatedDate,
		&User.UpdatedDate,
		&User.LastActive,
	)
	if err != nil {
		return nil, fmt.Errorf("get session user query error: %v", err)
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
		return fmt.Errorf("delete user session query error: %v", sessionErr)
	}

	return nil
}
