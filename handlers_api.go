package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// handleAPIKeyGenerate handles generating an API key for a user
func (s *server) handleAPIKeyGenerate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		keyVal := s.getJSONRequestBody(r, w)
		APIKeyName := keyVal["name"].(string)

		UserID := vars["id"]
		UserCookieID := r.Context().Value(contextKeyUserID).(string)
		if UserID != UserCookieID {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		APIKey, keyErr := s.database.GenerateAPIKey(UserID, APIKeyName)
		if keyErr != nil {
			log.Println("error attempting to generate api key : " + keyErr.Error() + "\n")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		s.respondWithJSON(w, http.StatusOK, APIKey)
	}
}

// handleUserAPIKeys handles getting user API keys
func (s *server) handleUserAPIKeys() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		UserID := vars["id"]
		UserCookieID := r.Context().Value(contextKeyUserID).(string)
		if UserID != UserCookieID {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		APIKeys, keysErr := s.database.GetUserAPIKeys(UserID)
		if keysErr != nil {
			log.Println("error retrieving api keys : " + keysErr.Error() + "\n")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		s.respondWithJSON(w, http.StatusOK, APIKeys)
	}
}

// handleUserAPIKeyUpdate handles getting user API keys
func (s *server) handleUserAPIKeyUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		UserID := vars["id"]
		UserCookieID := r.Context().Value(contextKeyUserID).(string)
		if UserID != UserCookieID {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		APK := vars["keyID"]
		keyVal := s.getJSONRequestBody(r, w)
		active := keyVal["active"].(bool)

		APIKeys, keysErr := s.database.UpdateUserAPIKey(UserID, APK, active)
		if keysErr != nil {
			log.Println("error updating api key : " + keysErr.Error() + "\n")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		s.respondWithJSON(w, http.StatusOK, APIKeys)
	}
}

// handleUserAPIKeyDelete handles getting user API keys
func (s *server) handleUserAPIKeyDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		UserID := vars["id"]
		UserCookieID := r.Context().Value(contextKeyUserID).(string)
		if UserID != UserCookieID {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		APK := vars["keyID"]

		APIKeys, keysErr := s.database.DeleteUserAPIKey(UserID, APK)
		if keysErr != nil {
			log.Println("error deleting api key : " + keysErr.Error() + "\n")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		s.respondWithJSON(w, http.StatusOK, APIKeys)
	}
}
