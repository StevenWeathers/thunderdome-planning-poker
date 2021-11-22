package api

import (
	"testing"
)

// TestValidUserAccount calls validateUserAccount with valid user inputs for name, email, password1, and password2
func TestValidUserAccount(t *testing.T) {
	Name := "Thor"
	Email := "thor@thunderdome.dev"
	Password := "lokiIsAJoke"

	name, email, password, err := validateUserAccount(Name, Email, Password, Password)
	if err != nil || (Name != name || Email != email || Password != password) {
		t.Fatalf(`validateUserAccount = %v error`, err)
	}
}

// TestInvalidUserAccount calls validateUserAccount with invalid user input for email
func TestInvalidUserAccount(t *testing.T) {
	_, _, _, err := validateUserAccount("Thor", "thor", "lokiIsAJoke", "lokiIsAJoke")
	if err == nil {
		t.Fatalf(`validateUserAccount = %v, want error`, err)
	}
}

// TestValidUserPassword calls validateUserPassword with valid user password1 and password2
func TestValidUserPassword(t *testing.T) {
	Password := "lokiIsAJoke"

	password, err := validateUserPassword(Password, Password)
	if err != nil || (Password != password) {
		t.Fatalf(`validateUserAccount = %v error`, err)
	}
}

// TestInvalidUserPassword calls validateUserPassword with invalid user password1 or password2
func TestInvalidUserPassword(t *testing.T) {
	Password := "lokiIsAJoke"

	_, err := validateUserPassword(Password, Password+"fail")
	if err == nil {
		t.Fatalf(`validateUserAccount = %v, want error`, err)
	}
}
