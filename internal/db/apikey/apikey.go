package apikey

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/StevenWeathers/thunderdome-planning-poker/internal/db"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"

	"go.uber.org/zap"
)

// Service represents a PostgreSQL implementation of thunderdome.APIKeyDataSvc.
type Service struct {
	DB     *sql.DB
	Logger *otelzap.Logger
}

// GenerateApiKey generates a new API key for a User
func (d *Service) GenerateApiKey(ctx context.Context, userID string, keyName string) (*thunderdome.APIKey, error) {
	apiPrefix, prefixErr := db.RandomString(8)
	if prefixErr != nil {
		return nil, fmt.Errorf("error generating api prefix: %v", prefixErr)
	}

	apiSecret, secretErr := db.RandomString(32)
	if secretErr != nil {
		return nil, fmt.Errorf("error generating api secret: %v", secretErr)
	}

	rawKey := apiPrefix + "." + apiSecret
	hashedKey := db.HashString(rawKey)
	keyID := apiPrefix + "." + hashedKey

	apiKey := &thunderdome.APIKey{
		ID:          keyID,
		Name:        keyName,
		Key:         rawKey,
		UserID:      userID,
		Prefix:      apiPrefix,
		Active:      true,
		CreatedDate: time.Now(),
	}

	err := d.DB.QueryRowContext(ctx,
		`INSERT INTO thunderdome.api_key (id, name, user_id) VALUES ($1, $2, $3) RETURNING created_date;`,
		keyID,
		keyName,
		userID,
	).Scan(&apiKey.CreatedDate)
	if err != nil {
		return nil, fmt.Errorf("error creating api key: %v", err)
	}

	return apiKey, nil
}

// GetUserApiKeys gets a list of api keys for a user
func (d *Service) GetUserApiKeys(ctx context.Context, userID string) ([]*thunderdome.APIKey, error) {
	var keys = make([]*thunderdome.APIKey, 0)
	rows, err := d.DB.QueryContext(ctx,
		"SELECT id, name, user_id, active, created_date, updated_date FROM thunderdome.api_key WHERE user_id = $1 ORDER BY created_date",
		userID,
	)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var ak thunderdome.APIKey
			var key string

			if err := rows.Scan(
				&key,
				&ak.Name,
				&ak.UserID,
				&ak.Active,
				&ak.CreatedDate,
				&ak.UpdatedDate,
			); err != nil {
				d.Logger.Ctx(ctx).Error("GetUserApiKeys scan error", zap.Error(err))
			} else {
				splitKey := strings.Split(key, ".")
				ak.Prefix = splitKey[0]
				ak.ID = key
				keys = append(keys, &ak)
			}
		}
	}

	return keys, err
}

// UpdateUserApiKey updates a user api key (active column only)
func (d *Service) UpdateUserApiKey(ctx context.Context, userID string, keyID string, active bool) ([]*thunderdome.APIKey, error) {
	if _, err := d.DB.ExecContext(ctx,
		`UPDATE thunderdome.api_key SET active = $3, updated_date = NOW() WHERE id = $1 AND user_id = $2;`,
		keyID, userID, active); err != nil {
		return nil, fmt.Errorf("error updating api key: %v", err)
	}

	keys, keysErr := d.GetUserApiKeys(ctx, userID)
	if keysErr != nil {
		return nil, fmt.Errorf("error getting users api keys: %v", keysErr)
	}

	return keys, nil
}

// DeleteUserApiKey removes a users api key
func (d *Service) DeleteUserApiKey(ctx context.Context, userID string, keyID string) ([]*thunderdome.APIKey, error) {
	if _, err := d.DB.ExecContext(ctx,
		`DELETE FROM thunderdome.api_key WHERE id = $1 AND user_id = $2;`,
		keyID, userID); err != nil {
		return nil, fmt.Errorf("error deleting api key: %v", err)
	}

	keys, keysErr := d.GetUserApiKeys(ctx, userID)
	if keysErr != nil {
		return nil, fmt.Errorf("error getting users api keys: %v", keysErr)
	}

	return keys, nil
}

// GetApiKeyUser checks to see if the API key exists and returns the User
func (d *Service) GetApiKeyUser(ctx context.Context, apiKey string) (*thunderdome.User, error) {
	user := &thunderdome.User{}

	splitKey := strings.Split(apiKey, ".")
	hashedKey := db.HashString(apiKey)
	keyID := splitKey[0] + "." + hashedKey

	err := d.DB.QueryRowContext(ctx, `
		SELECT u.id, u.name, u.email, u.type, u.avatar, u.verified, u.notifications_enabled, COALESCE(u.country, ''), COALESCE(u.locale, ''), COALESCE(u.company, ''), COALESCE(u.job_title, ''), u.created_date, u.updated_date, u.last_active
		FROM thunderdome.api_key ak
		LEFT JOIN thunderdome.users u ON u.id = ak.user_id
		WHERE ak.id = $1 AND ak.active = true
`,
		keyID,
	).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Type,
		&user.Avatar,
		&user.Verified,
		&user.NotificationsEnabled,
		&user.Country,
		&user.Locale,
		&user.Company,
		&user.JobTitle,
		&user.CreatedDate,
		&user.UpdatedDate,
		&user.LastActive)
	if err != nil {
		return nil, fmt.Errorf("active API Key match not found: %v", err)
	}

	user.GravatarHash = db.CreateGravatarHash(user.Email)

	return user, nil
}

// GetAPIKeys gets a list of api keys
func (d *Service) GetAPIKeys(ctx context.Context, limit int, offset int) []*thunderdome.UserAPIKey {
	var keys = make([]*thunderdome.UserAPIKey, 0)
	rows, err := d.DB.QueryContext(ctx,
		`SELECT apk.id, apk.name, u.id, u.name, u.email, apk.active, apk.created_date, apk.updated_date
		FROM thunderdome.api_key apk
		LEFT JOIN thunderdome.users u ON apk.user_id = u.id
		ORDER BY apk.created_date
		LIMIT $1
		OFFSET $2;`,
		limit,
		offset,
	)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var ak thunderdome.UserAPIKey
			var key string

			if err := rows.Scan(
				&key,
				&ak.Name,
				&ak.UserID,
				&ak.UserName,
				&ak.UserEmail,
				&ak.Active,
				&ak.CreatedDate,
				&ak.UpdatedDate,
			); err != nil {
				d.Logger.Ctx(ctx).Error("apikeys_list scan error", zap.Error(err))
			} else {
				splitKey := strings.Split(key, ".")
				ak.Prefix = splitKey[0]
				ak.ID = key
				keys = append(keys, &ak)
			}
		}
	} else {
		d.Logger.Ctx(ctx).Error("apikeys_list query error", zap.Error(err))
	}

	return keys
}
