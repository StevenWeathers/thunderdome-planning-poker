package database

import (
	"crypto/aes"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"
)

// EncryptAES encrypts a string using AES
func (d *Database) EncryptAES(plaintext string) string {
	key := []byte(d.config.SecretKey)
	c, _ := aes.NewCipher(key)
	out := make([]byte, len(plaintext))
	c.Encrypt(out, []byte(plaintext))

	return hex.EncodeToString(out)
}

// GenerateAPIKey generates a new API key for a Warrior
func (d *Database) GenerateAPIKey(WarriorID string, KeyName string) (*APIKey, error) {
	prefix := make([]byte, 4)
	if _, prefixErr := rand.Read(prefix); prefixErr != nil {
		err := errors.New("error generating api prefix")
		log.Println(err)
		log.Println(prefixErr)
		return nil, err
	}
	apiPrefix := fmt.Sprintf("%x", prefix)

	secret := make([]byte, 16)
	if _, secretErr := rand.Read(secret); secretErr != nil {
		err := errors.New("error generating api secret")
		log.Println(err)
		log.Println(secretErr)
		return nil, err
	}
	apiSecret := fmt.Sprintf("%x", secret)

	APIKEY := &APIKey{
		Name:        KeyName,
		Key:         apiPrefix + "." + apiSecret,
		WarriorID:   WarriorID,
		Prefix:      apiPrefix,
		Active:      true,
		CreatedDate: time.Now(),
	}
	hashedSecret := d.EncryptAES(apiSecret)
	keyID := apiPrefix + "." + hashedSecret

	e := d.db.QueryRow(
		`INSERT INTO api_keys (id, name, warrior_id ) VALUES ($1, $2, $3) RETURNING created_date`,
		keyID,
		KeyName,
		WarriorID,
	).Scan(&APIKEY.CreatedDate)
	if e != nil {
		log.Println(e)
		return nil, errors.New("unable to create new api key")
	}

	return APIKEY, nil
}

// GetWarriorAPIKeys gets a list of api keys for a warrior
func (d *Database) GetWarriorAPIKeys(WarriorID string) ([]*APIKey, error) {
	var APIKeys = make([]*APIKey, 0)
	rows, err := d.db.Query(
		"SELECT id, name, warrior_id, active, created_date, updated_date FROM api_keys WHERE warrior_id = $1 ORDER BY created_date",
		WarriorID,
	)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var ak APIKey
			var key string

			if err := rows.Scan(
				&key,
				&ak.Name,
				&ak.WarriorID,
				&ak.Active,
				&ak.CreatedDate,
				&ak.UpdatedDate,
			); err != nil {
				log.Println(err)
			} else {
				splitKey := strings.Split(key, ".")
				ak.Prefix = splitKey[0]
				ak.ID = key
				APIKeys = append(APIKeys, &ak)
			}
		}
	}

	return APIKeys, err
}

// UpdateWarriorAPIKey updates a warriors api key (active column only)
func (d *Database) UpdateWarriorAPIKey(WarriorID string, KeyID string, Active bool) ([]*APIKey, error) {
	if _, err := d.db.Exec(
		`UPDATE api_keys SET active = $3, updated_date = NOW() WHERE id = $1 AND warrior_id = $2;`, KeyID, WarriorID, Active); err != nil {
		log.Println(err)
		return nil, err
	}

	keys, keysErr := d.GetWarriorAPIKeys(WarriorID)
	if keysErr != nil {
		log.Println(keysErr)
		return nil, keysErr
	}

	return keys, nil
}

// DeleteWarriorAPIKey removes a warriors api key
func (d *Database) DeleteWarriorAPIKey(WarriorID string, KeyID string) ([]*APIKey, error) {
	if _, err := d.db.Exec(
		`DELETE FROM api_keys WHERE id = $1 AND warrior_id = $2;`, KeyID, WarriorID); err != nil {
		log.Println(err)
		return nil, err
	}

	keys, keysErr := d.GetWarriorAPIKeys(WarriorID)
	if keysErr != nil {
		log.Println(keysErr)
		return nil, keysErr
	}

	return keys, nil
}

// ValidateAPIKey checks to see if the API key exists in the database and if so returns WarriorID
func (d *Database) ValidateAPIKey(APK string) (WarriorID string, ValidatationErr error) {
	var warID string = ""

	splitKey := strings.Split(APK, ".")
	apiSecret := splitKey[1]
	hashedSecret := d.EncryptAES(apiSecret)
	keyID := splitKey[0] + "." + hashedSecret

	e := d.db.QueryRow(
		`SELECT warrior_id FROM api_keys WHERE id = $1 AND active = true`,
		keyID,
	).Scan(&warID)
	if e != nil {
		log.Println(e)
		return "", errors.New("active API Key match not found")
	}

	return warID, nil
}
