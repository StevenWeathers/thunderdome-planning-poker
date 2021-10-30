package api

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// handleAPIKeyGenerate handles generating an API key for a user
// @Summary Generate API Key
// @Description Generates an API key for the authenticated user
// @Tags apikey
// @Produce  json
// @Param id path int false "the user ID to generate API key for"
// @Success 200
// @Router /warrior/{id}/apikey [post]
func (a *api) handleAPIKeyGenerate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		keyVal := a.getJSONRequestBody(r, w)
		APIKeyName := keyVal["name"].(string)

		UserID := vars["id"]
		UserCookieID := r.Context().Value(contextKeyUserID).(string)
		if UserID != UserCookieID {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		APIKey, keyErr := a.db.GenerateAPIKey(UserID, APIKeyName)
		if keyErr != nil {
			log.Println("error attempting to generate api key : " + keyErr.Error() + "\n")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		a.respondWithJSON(w, http.StatusOK, APIKey)
	}
}

// handleUserAPIKeys handles getting user API keys
// @Summary Get API Keys
// @Description get list of API keys for authenticated user
// @Tags apikey
// @Produce  json
// @Param id path int false "the user ID to get API keys for"
// @Success 200
// @Router /warrior/{id}/apikeys [get]
func (a *api) handleUserAPIKeys() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		UserID := vars["id"]
		UserCookieID := r.Context().Value(contextKeyUserID).(string)
		if UserID != UserCookieID {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		APIKeys, keysErr := a.db.GetUserAPIKeys(UserID)
		if keysErr != nil {
			log.Println("error retrieving api keys : " + keysErr.Error() + "\n")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		a.respondWithJSON(w, http.StatusOK, APIKeys)
	}
}

// handleUserAPIKeyUpdate handles updating a users API key
// @Summary Update API Key
// @Description Updates the API key of the authenticated user
// @Tags apikey
// @Produce  json
// @Param id path int false "the user ID to update API key"
// @Param keyID path int false "the API Key ID to update"
// @Success 200
// @Router /warrior/{id}/apikey/{keyID} [put]
func (a *api) handleUserAPIKeyUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		UserID := vars["id"]
		UserCookieID := r.Context().Value(contextKeyUserID).(string)
		if UserID != UserCookieID {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		APK := vars["keyID"]
		keyVal := a.getJSONRequestBody(r, w)
		active := keyVal["active"].(bool)

		APIKeys, keysErr := a.db.UpdateUserAPIKey(UserID, APK, active)
		if keysErr != nil {
			log.Println("error updating api key : " + keysErr.Error() + "\n")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		a.respondWithJSON(w, http.StatusOK, APIKeys)
	}
}

// handleUserAPIKeyDelete handles deleting a users API key
// @Summary Delete API Key
// @Description Deletes the API key
// @Tags apikey
// @Produce  json
// @Param id path int false "the user ID to update API key"
// @Param keyID path int false "the API Key ID to update"
// @Success 200
// @Router /warrior/{id}/apikey/{keyID} [delete]
func (a *api) handleUserAPIKeyDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		UserID := vars["id"]
		UserCookieID := r.Context().Value(contextKeyUserID).(string)
		if UserID != UserCookieID {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		APK := vars["keyID"]

		APIKeys, keysErr := a.db.DeleteUserAPIKey(UserID, APK)
		if keysErr != nil {
			log.Println("error deleting api key : " + keysErr.Error() + "\n")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		a.respondWithJSON(w, http.StatusOK, APIKeys)
	}
}
