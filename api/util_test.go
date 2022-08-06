package api

import (
	"testing"

	"github.com/go-playground/validator/v10"
)

func TestMain(m *testing.M) {
	validate = validator.New()
}

// TestValidUserAccount calls validateUserAccountWithPasswords with valid user inputs for name, email, password1, and password2
func TestValidUserAccount(t *testing.T) {
	Name := "Thor"
	Email := "thor@thunderdome.dev"
	Password := "lokiIsAJoke"

	name, email, password, err := validateUserAccountWithPasswords(Name, Email, Password, Password)
	if err != nil || (Name != name || Email != email || Password != password) {
		t.Fatalf(`validateUserAccountWithPasswords = %v error`, err)
	}
}

// TestInvalidUserAccount calls validateUserAccountWithPasswords with invalid user input for email
func TestInvalidUserAccount(t *testing.T) {
	_, _, _, err := validateUserAccountWithPasswords("Thor", "thor", "lokiIsAJoke", "lokiIsAJoke")
	if err == nil {
		t.Fatalf(`validateUserAccountWithPasswords = %v, want error`, err)
	}
}

// TestValidUserPassword calls validateUserPassword with valid user password1 and password2
func TestValidUserPassword(t *testing.T) {
	Password := "lokiIsAJoke"

	password, err := validateUserPassword(Password, Password)
	if err != nil || (Password != password) {
		t.Fatalf(`validateUserAccountWithPasswords = %v error`, err)
	}
}

// TestInvalidUserPassword calls validateUserPassword with invalid user password1 or password2
func TestInvalidUserPassword(t *testing.T) {
	Password := "lokiIsAJoke"

	_, err := validateUserPassword(Password, Password+"fail")
	if err == nil {
		t.Fatalf(`validateUserAccountWithPasswords = %v, want error`, err)
	}
}
