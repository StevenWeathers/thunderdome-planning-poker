package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"golang.org/x/crypto/bcrypt"
	"gopkg.in/go-playground/validator.v9"
)

type warriorAccount struct {
	Name      string `json:"name" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password1 string `json:"password1" validate:"required,min=6,max=72"`
	Password2 string `json:"password2" validate:"required,min=6,max=72,eqfield=Password1"`
}

type warriorPassword struct {
	Password1 string `json:"password1" validate:"required,min=6,max=72"`
	Password2 string `json:"password2" validate:"required,min=6,max=72,eqfield=Password1"`
}

// ValidateWarriorAccount makes sure warrior name, email, and password are valid before creating the account
func ValidateWarriorAccount(name string, email string, pwd1 string, pwd2 string) (WarriorName string, WarriorEmail string, WarriorPassword string, validateErr error) {
	v := validator.New()
	a := warriorAccount{
		Name:      name,
		Email:     email,
		Password1: pwd1,
		Password2: pwd2,
	}
	err := v.Struct(a)

	return name, email, pwd1, err
}

// ValidateWarriorPassword makes sure warrior password is valid before updating the password
func ValidateWarriorPassword(pwd1 string, pwd2 string) (WarriorPassword string, validateErr error) {
	v := validator.New()
	a := warriorPassword{
		Password1: pwd1,
		Password2: pwd2,
	}
	err := v.Struct(a)

	return pwd1, err
}

// HashAndSalt takes a password byte and salt + hashes it
// returning a hash string to store in db
func HashAndSalt(pwd []byte) (string, error) {
	// Use GenerateFromPassword to hash & salt pwd.
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
		return "", err
	}
	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash), nil
}

// ComparePasswords takes a password hash and compares it to entered password bytes
// returning true if matches false if not
func ComparePasswords(hashedPwd string, plainPwd []byte) bool {
	// Since we'll be getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

// GetEnv gets environment variable matching key string
// and if it finds none uses fallback string
// returning either the matching or fallback string
func GetEnv(key string, fallback string) string {
	var result = os.Getenv(key)

	if result == "" {
		result = fallback
	}

	return result
}

// GetIntEnv gets an environment variable and converts it to an int
// and if it finds none uses fallback
func GetIntEnv(key string, fallback int) int {
	var intResult = fallback
	var stringResult = os.Getenv(key)

	if stringResult != "" {
		v, _ := strconv.Atoi(stringResult)
		intResult = v
	}

	return intResult
}

// GetBoolEnv gets an environment variable and converts it to a bool
// and if it finds none uses fallback
func GetBoolEnv(key string, fallback bool) bool {
	var boolResult = fallback
	var stringResult = os.Getenv(key)

	if stringResult != "" {
		b, _ := strconv.ParseBool(stringResult)
		boolResult = b
	}

	return boolResult
}

// RespondWithJSON takes a payload and writes the response
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
