package db

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"golang.org/x/crypto/bcrypt"
	"log"
	"math/big"
)

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

// randomBase64String returns a random secure base64 string of X length
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
		log.Println(err)
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
	if err != nil {
		log.Println(err)
		return false
	}

	return true
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
