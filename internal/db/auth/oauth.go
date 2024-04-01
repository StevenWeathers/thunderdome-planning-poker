package auth

import (
	"context"
	"errors"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
)

// OauthAuthUser authenticate the oauth user or creates a new user
func (d *Service) OauthAuthUser(ctx context.Context, provider string, email string, emailVerified bool, name string, pictureUrl string) (*thunderdome.User, string, error) {
	var user thunderdome.User

	err := d.DB.QueryRowContext(ctx,
		`INSERT INTO thunderdome.users (type, provider, email, verified, name, picture_url)
				VALUES ('REGISTERED', $1, $2, $3, $4, $5)
				ON CONFLICT (provider, lower((email)::text)) DO UPDATE
				SET verified = EXCLUDED.verified
 				RETURNING id, name, email, type, verified, notifications_enabled,
 				 COALESCE(locale, ''), disabled, theme, COALESCE(picture_url, '')`,
		provider, email, emailVerified, name, pictureUrl,
	).Scan(
		&user.Id,
		&user.Name,
		&user.Email,
		&user.Type,
		&user.Verified,
		&user.NotificationsEnabled,
		&user.Locale,
		&user.Disabled,
		&user.Theme,
		&user.PictureURL,
	)
	if err != nil {
		return nil, "", err
	}

	if user.Disabled {
		return nil, "", errors.New("USER_DISABLED")
	}

	SessionId, sessErr := d.CreateSession(ctx, user.Id)
	if sessErr != nil {
		return nil, "", sessErr
	}

	return &user, SessionId, nil
}
