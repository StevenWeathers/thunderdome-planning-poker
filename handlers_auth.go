package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/spf13/viper"
)

// handleLogin attempts to login the user by comparing email/password to whats in DB
func (s *server) handleLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		keyVal := s.getJSONRequestBody(r, w)
		UserEmail := keyVal["warriorEmail"].(string)
		UserPassword := keyVal["warriorPassword"].(string)

		authedUser, err := s.authUserDatabase(UserEmail, UserPassword)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		cookie := s.createCookie(authedUser.UserID)
		if cookie != nil {
			http.SetCookie(w, cookie)
		} else {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		s.respondWithJSON(w, http.StatusOK, authedUser)
	}
}

// handleLdapLogin attempts to authenticate the user by looking up and authenticating
// via ldap, and then creates the user if not existing and logs them in
func (s *server) handleLdapLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		keyVal := s.getJSONRequestBody(r, w)
		UserEmail := keyVal["warriorEmail"].(string)
		UserPassword := keyVal["warriorPassword"].(string)

		authedUser, err := s.authAndCreateUserLdap(UserEmail, UserPassword)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		cookie := s.createCookie(authedUser.UserID)
		if cookie != nil {
			http.SetCookie(w, cookie)
		} else {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		s.respondWithJSON(w, http.StatusOK, authedUser)
	}
}

// handleLogout clears the user cookie(s) ending session
func (s *server) handleLogout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.clearUserCookies(w)
		return
	}
}

// handleUserRecruit registers a user as a guest user
func (s *server) handleUserRecruit() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		AllowGuests := viper.GetBool("config.allow_guests")
		if !AllowGuests {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		keyVal := s.getJSONRequestBody(r, w)

		UserName := keyVal["warriorName"].(string)

		newUser, err := s.database.CreateUserGuest(UserName)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		s.createUserCookie(w, false, newUser.UserID)

		s.respondWithJSON(w, http.StatusOK, newUser)
	}
}

// handleUserEnlist registers a new authenticated user
func (s *server) handleUserEnlist() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		AllowRegistration := viper.GetBool("config.allow_registration")
		if !AllowRegistration {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		keyVal := s.getJSONRequestBody(r, w)

		ActiveUserID, _ := s.validateUserCookie(w, r)

		UserName, UserEmail, UserPassword, accountErr := ValidateUserAccount(
			keyVal["warriorName"].(string),
			strings.ToLower(keyVal["warriorEmail"].(string)),
			keyVal["warriorPassword1"].(string),
			keyVal["warriorPassword2"].(string),
		)

		if accountErr != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		newUser, VerifyID, err := s.database.CreateUserRegistered(UserName, UserEmail, UserPassword, ActiveUserID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		s.createUserCookie(w, true, newUser.UserID)

		s.email.SendWelcome(UserName, UserEmail, VerifyID)

		s.respondWithJSON(w, http.StatusOK, newUser)
	}
}

// handleForgotPassword attempts to send a password reset email
func (s *server) handleForgotPassword() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		keyVal := s.getJSONRequestBody(r, w)
		UserEmail := keyVal["warriorEmail"].(string)

		ResetID, UserName, resetErr := s.database.UserResetRequest(UserEmail)
		if resetErr == nil {
			s.email.SendForgotPassword(UserName, UserEmail, ResetID)
		}

		w.WriteHeader(http.StatusOK)
		return
	}
}

// handleResetPassword attempts to reset a users password
func (s *server) handleResetPassword() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		keyVal := s.getJSONRequestBody(r, w)
		ResetID := keyVal["resetId"].(string)

		UserPassword, passwordErr := ValidateUserPassword(
			keyVal["warriorPassword1"].(string),
			keyVal["warriorPassword2"].(string),
		)

		if passwordErr != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		UserName, UserEmail, resetErr := s.database.UserResetPassword(ResetID, UserPassword)
		if resetErr != nil {
			log.Println("error attempting to reset user password : " + resetErr.Error() + "\n")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		s.email.SendPasswordReset(UserName, UserEmail)

		return
	}
}
