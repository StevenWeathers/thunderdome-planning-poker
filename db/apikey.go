package db

import (
	"errors"
	"log"
	"strings"
	"time"

	"github.com/StevenWeathers/thunderdome-planning-poker/model"
)

// GenerateApiKey generates a new API key for a User
func (d *Database) GenerateApiKey(UserID string, KeyName string) (*model.APIKey, error) {
	apiPrefix, prefixErr := randomString(8)
	if prefixErr != nil {
		err := errors.New("error generating api prefix")
		log.Println(err)
		log.Println(prefixErr)
		return nil, err
	}

	apiSecret, secretErr := randomString(32)
	if secretErr != nil {
		err := errors.New("error generating api secret")
		log.Println(err)
		log.Println(secretErr)
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

	e := d.db.QueryRow(
		`SELECT createdDate FROM user_apikey_add($1, $2, $3);`,
		keyID,
		KeyName,
		UserID,
	).Scan(&APIKEY.CreatedDate)
	if e != nil {
		log.Println(e)
		return nil, errors.New("unable to create new api key")
	}

	return APIKEY, nil
}

// GetUserApiKeys gets a list of api keys for a user
func (d *Database) GetUserApiKeys(UserID string) ([]*model.APIKey, error) {
	var APIKeys = make([]*model.APIKey, 0)
	rows, err := d.db.Query(
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
				log.Println(err)
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
func (d *Database) UpdateUserApiKey(UserID string, KeyID string, Active bool) ([]*model.APIKey, error) {
	if _, err := d.db.Exec(
		`CALL user_apikey_update($1, $2, $3);`, KeyID, UserID, Active); err != nil {
		log.Println(err)
		return nil, err
	}

	keys, keysErr := d.GetUserApiKeys(UserID)
	if keysErr != nil {
		log.Println(keysErr)
		return nil, keysErr
	}

	return keys, nil
}

// DeleteUserApiKey removes a users api key
func (d *Database) DeleteUserApiKey(UserID string, KeyID string) ([]*model.APIKey, error) {
	if _, err := d.db.Exec(
		`CALL user_apikey_delete($1, $2);`, KeyID, UserID); err != nil {
		log.Println(err)
		return nil, err
	}

	keys, keysErr := d.GetUserApiKeys(UserID)
	if keysErr != nil {
		log.Println(keysErr)
		return nil, keysErr
	}

	return keys, nil
}

// ValidateApiKey checks to see if the API key exists in the database and if so returns UserID
func (d *Database) ValidateApiKey(APK string) (UserID string, ValidatationErr error) {
	var usrID string = ""

	splitKey := strings.Split(APK, ".")
	hashedKey := hashString(APK)
	keyID := splitKey[0] + "." + hashedKey

	e := d.db.QueryRow(
		`SELECT user_id FROM api_keys WHERE id = $1 AND active = true`,
		keyID,
	).Scan(&usrID)
	if e != nil {
		log.Println(e)
		return "", errors.New("active API Key match not found")
	}

	return usrID, nil
}
