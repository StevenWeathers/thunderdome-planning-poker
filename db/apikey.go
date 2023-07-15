package db

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"

	"go.uber.org/zap"
)

// APIKeyService represents a PostgreSQL implementation of thunderdome.APIKeyService.
type APIKeyService struct {
	DB     *sql.DB
	Logger *otelzap.Logger
}

// GenerateApiKey generates a new API key for a User
func (d *APIKeyService) GenerateApiKey(ctx context.Context, UserID string, KeyName string) (*thunderdome.APIKey, error) {
	apiPrefix, prefixErr := randomString(8)
	if prefixErr != nil {
		err := errors.New("error generating api prefix")
		d.Logger.Ctx(ctx).Error("error generating api prefix", zap.Error(prefixErr))
		return nil, err
	}

	apiSecret, secretErr := randomString(32)
	if secretErr != nil {
		err := errors.New("error generating api secret")
		d.Logger.Ctx(ctx).Error("error generating api secret", zap.Error(prefixErr))
		return nil, err
	}

	key := apiPrefix + "." + apiSecret
	hashedKey := hashString(key)
	keyID := apiPrefix + "." + hashedKey

	APIKEY := &thunderdome.APIKey{
		Id:          keyID,
		Name:        KeyName,
		Key:         key,
		UserId:      UserID,
		Prefix:      apiPrefix,
		Active:      true,
		CreatedDate: time.Now(),
	}

	e := d.DB.QueryRowContext(ctx,
		`INSERT INTO thunderdome.api_key (id, name, user_id) VALUES ($1, $2, $3) RETURNING created_date;`,
		keyID,
		KeyName,
		UserID,
	).Scan(&APIKEY.CreatedDate)
	if e != nil {
		d.Logger.Ctx(ctx).Error("user_apikey_add query error", zap.Error(e))
		return nil, errors.New("unable to create new api key")
	}

	return APIKEY, nil
}

// GetUserApiKeys gets a list of api keys for a user
func (d *APIKeyService) GetUserApiKeys(ctx context.Context, UserID string) ([]*thunderdome.APIKey, error) {
	var APIKeys = make([]*thunderdome.APIKey, 0)
	rows, err := d.DB.QueryContext(ctx,
		"SELECT id, name, user_id, active, created_date, updated_date FROM thunderdome.api_key WHERE user_id = $1 ORDER BY created_date",
		UserID,
	)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var ak thunderdome.APIKey
			var key string

			if err := rows.Scan(
				&key,
				&ak.Name,
				&ak.UserId,
				&ak.Active,
				&ak.CreatedDate,
				&ak.UpdatedDate,
			); err != nil {
				d.Logger.Ctx(ctx).Error("GetUserApiKeys scan error", zap.Error(err))
			} else {
				splitKey := strings.Split(key, ".")
				ak.Prefix = splitKey[0]
				ak.Id = key
				APIKeys = append(APIKeys, &ak)
			}
		}
	}

	return APIKeys, err
}

// UpdateUserApiKey updates a user api key (active column only)
func (d *APIKeyService) UpdateUserApiKey(ctx context.Context, UserID string, KeyID string, Active bool) ([]*thunderdome.APIKey, error) {
	if _, err := d.DB.ExecContext(ctx,
		`UPDATE thunderdome.api_key SET active = $3, updated_date = NOW() WHERE id = $1 AND user_id = $2;`,
		KeyID, UserID, Active); err != nil {
		d.Logger.Ctx(ctx).Error("UpdateUserApiKey query error", zap.Error(err))
		return nil, err
	}

	keys, keysErr := d.GetUserApiKeys(ctx, UserID)
	if keysErr != nil {
		d.Logger.Ctx(ctx).Error("GetUserApiKeys query error", zap.Error(keysErr))
		return nil, keysErr
	}

	return keys, nil
}

// DeleteUserApiKey removes a users api key
func (d *APIKeyService) DeleteUserApiKey(ctx context.Context, UserID string, KeyID string) ([]*thunderdome.APIKey, error) {
	if _, err := d.DB.ExecContext(ctx,
		`DELETE FROM thunderdome.api_key WHERE id = $1 AND user_id = $2;`,
		KeyID, UserID); err != nil {
		d.Logger.Ctx(ctx).Error("CALL thunderdome.user_apikey_delete error", zap.Error(err))
		return nil, err
	}

	keys, keysErr := d.GetUserApiKeys(ctx, UserID)
	if keysErr != nil {
		d.Logger.Ctx(ctx).Error("GetUserApiKeys query error", zap.Error(keysErr))
		return nil, keysErr
	}

	return keys, nil
}

// GetApiKeyUser checks to see if the API key exists and returns the User
func (d *APIKeyService) GetApiKeyUser(ctx context.Context, APK string) (*thunderdome.User, error) {
	User := &thunderdome.User{}

	splitKey := strings.Split(APK, ".")
	hashedKey := hashString(APK)
	keyID := splitKey[0] + "." + hashedKey

	e := d.DB.QueryRowContext(ctx, `
		SELECT u.id, u.name, u.email, u.type, u.avatar, u.verified, u.notifications_enabled, COALESCE(u.country, ''), COALESCE(u.locale, ''), COALESCE(u.company, ''), COALESCE(u.job_title, ''), u.created_date, u.updated_date, u.last_active 
		FROM thunderdome.api_key ak
		LEFT JOIN thunderdome.users u ON u.id = ak.user_id
		WHERE ak.id = $1 AND ak.active = true
`,
		keyID,
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
		d.Logger.Ctx(ctx).Error("GetApiKeyUser query error", zap.Error(e))
		return nil, errors.New("active API Key match not found")
	}

	User.GravatarHash = createGravatarHash(User.Email)

	return User, nil
}

// GetAPIKeys gets a list of api keys
func (d *APIKeyService) GetAPIKeys(ctx context.Context, Limit int, Offset int) []*thunderdome.UserAPIKey {
	var APIKeys = make([]*thunderdome.UserAPIKey, 0)
	rows, err := d.DB.QueryContext(ctx,
		`SELECT apk.id, apk.name, u.id, u.name, u.email, apk.active, apk.created_date, apk.updated_date
		FROM thunderdome.api_key apk
		LEFT JOIN thunderdome.users u ON apk.user_id = u.id
		ORDER BY apk.created_date
		LIMIT $1
		OFFSET $2;`,
		Limit,
		Offset,
	)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var ak thunderdome.UserAPIKey
			var key string

			if err := rows.Scan(
				&key,
				&ak.Name,
				&ak.UserId,
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
				ak.Id = key
				APIKeys = append(APIKeys, &ak)
			}
		}
	} else {
		d.Logger.Ctx(ctx).Error("apikeys_list query error", zap.Error(err))
	}

	return APIKeys
}
