package main

import (
	"context"
	"log"
	"net/http"
	"strings"
)

// adminOnly middleware checks if the user is an admin, otherwise reject their request
func (s *server) adminOnly(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get(apiKeyHeaderName)
		apiKey = strings.TrimSpace(apiKey)
		var UserID string

		if apiKey != "" {
			var apiKeyErr error
			UserID, apiKeyErr = s.database.ValidateAPIKey(apiKey)
			if apiKeyErr != nil {
				log.Println("error validating api key : " + apiKeyErr.Error() + "\n")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		} else {
			var cookieErr error
			UserID, cookieErr = s.validateUserCookie(w, r)
			if cookieErr != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		}

		adminErr := s.database.ConfirmAdmin(UserID)
		if adminErr != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), contextKeyUserID, UserID)

		h(w, r.WithContext(ctx))
	}
}

// userOnly validates that the request was made by a valid user
func (s *server) userOnly(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get(apiKeyHeaderName)
		apiKey = strings.TrimSpace(apiKey)
		var UserID string

		if apiKey != "" {
			var apiKeyErr error
			UserID, apiKeyErr = s.database.ValidateAPIKey(apiKey)
			if apiKeyErr != nil {
				log.Println("error validating api key : " + apiKeyErr.Error() + "\n")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		} else {
			var cookieErr error
			UserID, cookieErr = s.validateUserCookie(w, r)
			if cookieErr != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		}

		_, UserErr := s.database.GetUser(UserID)
		if UserErr != nil {
			log.Println("error finding user : " + UserErr.Error() + "\n")
			s.clearUserCookies(w)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), contextKeyUserID, UserID)

		h(w, r.WithContext(ctx))
	}
}
