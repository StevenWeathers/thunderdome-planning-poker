package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/StevenWeathers/thunderdome-planning-poker/internal/db"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
)

// OauthCreateNonce creates a new oauth nonce
func (d *Service) OauthCreateNonce(ctx context.Context) (string, error) {
	nonceId, err := db.RandomBase64String(32)
	if err != nil {
		return "", err
	}

	if _, nonceErr := d.DB.ExecContext(ctx, `
		INSERT INTO thunderdome.auth_nonce (nonce_id) VALUES ($1);
		`,
		nonceId,
	); nonceErr != nil {
		return "", fmt.Errorf("create oauth nonce query error: %v", nonceErr)
	}

	return nonceId, nil
}

func (d *Service) OauthValidateNonce(ctx context.Context, nonceId string) error {
	var expireDate *time.Time
	if err := d.DB.QueryRowContext(ctx,
		`DELETE FROM thunderdome.auth_nonce WHERE nonce_id = $1 RETURNING expire_date;`,
		nonceId,
	).Scan(&expireDate); err != nil {
		return err
	}

	if expireDate == nil || time.Now().After(*expireDate) {
		return fmt.Errorf("nonce invalid")
	}

	return nil
}

// OauthAuthUser authenticate the oauth user or creates a new user
func (d *Service) OauthAuthUser(ctx context.Context, provider string, sub string, email string, emailVerified bool, name string, pictureUrl string) (*thunderdome.User, string, error) {
	var user thunderdome.User

	err := d.DB.QueryRowContext(ctx,
		`SELECT u.id, u.name, ai.email, u.type, ai.verified, u.notifications_enabled,
 				 COALESCE(u.locale, ''), u.disabled, u.theme, COALESCE(ai.picture, u.picture, '')
 				 FROM thunderdome.auth_identity ai
 				 JOIN thunderdome.users u ON u.id = ai.user_id
 				 WHERE ai.provider = $1 AND ai.sub = $2;`,
		provider, sub,
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
		&user.Picture,
	)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		tx, txErr := d.DB.BeginTx(ctx, nil)
		if txErr != nil {
			return nil, "", txErr
		}
		userInsertErr := tx.QueryRowContext(ctx,
			`INSERT INTO thunderdome.users (name, email, type, verified, picture) 
					VALUES ($1, $2, $3, $4, $5)
					RETURNING id, name, email, type, verified, picture;`,
			name, email, thunderdome.RegisteredUserType, emailVerified, pictureUrl,
		).Scan(
			&user.Id,
			&user.Name,
			&user.Email,
			&user.Type,
			&user.Verified,
			&user.Picture,
		)
		if userInsertErr != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				return nil, "", fmt.Errorf("insert failed: %v, unable to rollback: %v\n", userInsertErr, rollbackErr)
			}
			return nil, "", userInsertErr
		}

		_, identityInsertErr := tx.ExecContext(ctx,
			`INSERT INTO thunderdome.auth_identity (user_id, provider, sub, email, picture, verified) 
					VALUES ($1, $2, $3, $4, $5, $6);`,
			user.Id, provider, sub, email, pictureUrl, emailVerified,
		)
		if identityInsertErr != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				return nil, "", fmt.Errorf("insert failed: %v, unable to rollback: %v\n", identityInsertErr, rollbackErr)
			}
			return nil, "", fmt.Errorf("update failed: %v", identityInsertErr)
		}
		if err := tx.Commit(); err != nil {
			return nil, "", err
		}
	} else if err != nil {
		return nil, "", err
	}

	if user.Disabled {
		return nil, "", errors.New("USER_DISABLED")
	}

	SessionId, sessErr := d.CreateSession(ctx, user.Id, true)
	if sessErr != nil {
		return nil, "", sessErr
	}

	return &user, SessionId, nil
}
