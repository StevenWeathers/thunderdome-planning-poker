package database

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"log"
	"strings"
	"time"
)

// HashAPIKey hashes the API key using SHA256 (not reversible)
func (d *Database) HashAPIKey(apikey string) string {
	data := []byte(apikey)
	hash := sha256.Sum256(data)
	result := hex.EncodeToString(hash[:])

	return result
}

// generate random secure string of X length
func random(length int) (string, error) {
	chars := "-_+=!$0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := make([]byte, length)

	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	for i, b := range bytes {
		bytes[i] = chars[b%byte(len(chars))]
	}

	return string(bytes), nil
}

// GenerateAPIKey generates a new API key for a Warrior
func (d *Database) GenerateAPIKey(WarriorID string, KeyName string) (*APIKey, error) {
	apiPrefix, prefixErr := random(8)
	if prefixErr != nil {
		err := errors.New("error generating api prefix")
		log.Println(err)
		log.Println(prefixErr)
		return nil, err
	}

	apiSecret, secretErr := random(32)
	if secretErr != nil {
		err := errors.New("error generating api secret")
		log.Println(err)
		log.Println(secretErr)
		return nil, err
	}

	APIKEY := &APIKey{
		Name:        KeyName,
		Key:         apiPrefix + "." + apiSecret,
		WarriorID:   WarriorID,
		Prefix:      apiPrefix,
		Active:      true,
		CreatedDate: time.Now(),
	}
	hashedKey := d.HashAPIKey(APIKEY.Key)
	keyID := apiPrefix + "." + hashedKey

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
	hashedKey := d.HashAPIKey(APK)
	keyID := splitKey[0] + "." + hashedKey

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
