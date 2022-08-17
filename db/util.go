package db

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"io"
	"math/big"
	"regexp"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// sanitizeEmail removes any non-valid email characters and lowercase's email
func sanitizeEmail(email string) string {
	emailRegExp := regexp.MustCompile(`[^a-zA-Z0-9-_.@+]`)

	return string(emailRegExp.ReplaceAll(
		[]byte(strings.ToLower(email)), []byte("")),
	)
}

// contains checks if a string is present in a slice
func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

// random generates a random secure byte of X length
func random(length int) ([]byte, error) {
	chars := "-_+=!$0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := make([]byte, length)

	if _, err := rand.Read(bytes); err != nil {
		return nil, err // out of randomness, should never happen
	}

	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		if err != nil {
			return nil, err
		}
		bytes[i] = chars[num.Int64()]
	}

	return bytes, nil
}

// randomString returns a random secure string of X length
func randomString(l int) (string, error) {
	s, err := random(l)
	if err != nil {
		return "", err
	}
	return string(s), nil
}

// randomBase64String returns a random secure string of X length base64 encoded
func randomBase64String(l int) (string, error) {
	s, err := random(l)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(s), nil
}

// hashString hashes the string using SHA256 (not reversible)
func hashString(s string) string {
	data := []byte(s)
	hash := sha256.Sum256(data)
	result := hex.EncodeToString(hash[:])

	return result
}

// hashSaltPassword takes a password byte then salt + hashes it returning a hash string
func hashSaltPassword(UserPassword string) (string, error) {
	pwd := []byte(UserPassword)
	// Use GenerateFromPassword to hash & salt pwd.
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	// GenerateFromPassword returns a byte slice, so we need to
	// convert the bytes to a string and return it
	return string(hash), nil
}

// comparePasswords takes a password hash and compares it to entered password
func comparePasswords(hashedPwd string, password string) bool {
	// Since we'll be getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice
	byteHash := []byte(hashedPwd)
	SubmittedPassword := []byte(password)
	err := bcrypt.CompareHashAndPassword(byteHash, SubmittedPassword)

	return err == nil
}

// checkPasswordCost checks the passwords stored hash for bcrypt cost
// if it does not match current cost then return true and let auth update the hash
func checkPasswordCost(hashedPwd string) bool {
	byteHash := []byte(hashedPwd)

	hashCost, costErr := bcrypt.Cost(byteHash)
	if costErr == nil && hashCost != bcrypt.DefaultCost {
		return true
	}

	return false
}

// createHash creates a md5 hashed string from string
func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

// encrypt data for storing securely
func encrypt(data string, passphrase string) (string, error) {
	block, _ := aes.NewCipher([]byte(createHash(passphrase)))
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	ciphertext := gcm.Seal(nonce, nonce, []byte(data), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// decrypt data for sending to client
func decrypt(data string, passphrase string) (string, error) {
	dataByte, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}
	key := []byte(createHash(passphrase))
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonceSize := gcm.NonceSize()
	if len(dataByte) < nonceSize {
		return "", errors.New("unable to decrypt data")
	}
	nonce, ciphertext := dataByte[:nonceSize], dataByte[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}

// createGravatarHash md5 hashes email for gravatar use
func createGravatarHash(email string) string {
	gh := md5.Sum([]byte(email))
	return hex.EncodeToString(gh[:])
}
