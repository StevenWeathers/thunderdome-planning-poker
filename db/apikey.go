package db

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/StevenWeathers/thunderdome-planning-poker/model"
	"go.uber.org/zap"
)

// GenerateApiKey generates a new API key for a User
func (d *Database) GenerateApiKey(ctx context.Context, UserID string, KeyName string) (*model.APIKey, error) {
	apiPrefix, prefixErr := randomString(8)
	if prefixErr != nil {
		err := errors.New("error generating api prefix")
		d.logger.Ctx(ctx).Error("error generating api prefix", zap.Error(prefixErr))
		return nil, err
	}

	apiSecret, secretErr := randomString(32)
	if secretErr != nil {
		err := errors.New("error generating api secret")
		d.logger.Ctx(ctx).Error("error generating api secret", zap.Error(prefixErr))
		return nil, err
	}

	APIKEY := &model.APIKey{
		Name:        KeyName,
		Key:         apiPrefix + "." + apiSecret,
		UserId:      UserID,
		Prefix:      apiPrefix,
		Active:      true,
		CreatedDate: time.Now(),
	}
	hashedKey := hashString(APIKEY.Key)
	keyID := apiPrefix + "." + hashedKey

	e := d.db.QueryRowContext(ctx,
		`SELECT createdDate FROM user_apikey_add($1, $2, $3);`,
		keyID,
		KeyName,
		UserID,
	).Scan(&APIKEY.CreatedDate)
	if e != nil {
		d.logger.Ctx(ctx).Error("user_apikey_add query error", zap.Error(e))
		return nil, errors.New("unable to create new api key")
	}

	return APIKEY, nil
}

// GetUserApiKeys gets a list of api keys for a user
func (d *Database) GetUserApiKeys(ctx context.Context, UserID string) ([]*model.APIKey, error) {
	var APIKeys = make([]*model.APIKey, 0)
	rows, err := d.db.QueryContext(ctx,
		"SELECT id, name, user_id, active, created_date, updated_date FROM api_keys WHERE user_id = $1 ORDER BY created_date",
		UserID,
	)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var ak model.APIKey
			var key string

			if err := rows.Scan(
				&key,
				&ak.Name,
				&ak.UserId,
				&ak.Active,
				&ak.CreatedDate,
				&ak.UpdatedDate,
			); err != nil {
				d.logger.Ctx(ctx).Error("GetUserApiKeys scan error", zap.Error(err))
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
func (d *Database) UpdateUserApiKey(ctx context.Context, UserID string, KeyID string, Active bool) ([]*model.APIKey, error) {
	if _, err := d.db.ExecContext(ctx,
		`CALL user_apikey_update($1, $2, $3);`, KeyID, UserID, Active); err != nil {
		d.logger.Ctx(ctx).Error("UpdateUserApiKey query error", zap.Error(err))
		return nil, err
	}

	keys, keysErr := d.GetUserApiKeys(ctx, UserID)
	if keysErr != nil {
		d.logger.Ctx(ctx).Error("GetUserApiKeys query error", zap.Error(keysErr))
		return nil, keysErr
	}

	return keys, nil
}

// DeleteUserApiKey removes a users api key
func (d *Database) DeleteUserApiKey(ctx context.Context, UserID string, KeyID string) ([]*model.APIKey, error) {
	if _, err := d.db.ExecContext(ctx,
		`CALL user_apikey_delete($1, $2);`, KeyID, UserID); err != nil {
		d.logger.Ctx(ctx).Error("call user_apikey_delete error", zap.Error(err))
		return nil, err
	}

	keys, keysErr := d.GetUserApiKeys(ctx, UserID)
	if keysErr != nil {
		d.logger.Ctx(ctx).Error("GetUserApiKeys query error", zap.Error(keysErr))
		return nil, keysErr
	}

	return keys, nil
}

// GetApiKeyUser checks to see if the API key exists and returns the User
func (d *Database) GetApiKeyUser(ctx context.Context, APK string) (*model.User, error) {
	User := &model.User{}

	splitKey := strings.Split(APK, ".")
	hashedKey := hashString(APK)
	keyID := splitKey[0] + "." + hashedKey

	e := d.db.QueryRowContext(ctx, `
		SELECT u.id, u.name, u.email, u.type, u.avatar, u.verified, u.notifications_enabled, COALESCE(u.country, ''), COALESCE(u.locale, ''), COALESCE(u.company, ''), COALESCE(u.job_title, ''), u.created_date, u.updated_date, u.last_active 
		FROM api_keys ak
		LEFT JOIN users u ON u.id = ak.user_id
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
		d.logger.Ctx(ctx).Error("GetApiKeyUser query error", zap.Error(e))
		return nil, errors.New("active API Key match not found")
	}

	User.GravatarHash = createGravatarHash(User.Email)

	return User, nil
}
