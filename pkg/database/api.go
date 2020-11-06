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
func EncryptAES(key []byte, plaintext string) string {
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
	// cipher key
	key := "thisis32bitlongpassphraseimusing"
	hashedSecret := EncryptAES([]byte(key), apiSecret)
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
		"SELECT id, name, warrior_id, active, created_date FROM api_keys WHERE warrior_id = $1 ORDER BY created_date",
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
			); err != nil {
				log.Println(err)
			} else {
				splitKey := strings.Split(key, ".")
				ak.Prefix = splitKey[0]
				APIKeys = append(APIKeys, &ak)
			}
		}
	}

	return APIKeys, err
}
