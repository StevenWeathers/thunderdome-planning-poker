package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

// handleUserAPIKeys handles getting user API keys
// @Summary Get API Keys
// @Description get list of API keys for authenticated user
// @Tags apikey
// @Produce  json
// @Param id path int false "the user ID to get API keys for"
// @Success 200 object standardJsonResponse{data=[]model.APIKey}
// @Failure 403 object standardJsonResponse{}
// @Failure 500 object standardJsonResponse{}
// @Router /users/{id}/apikeys [get]
func (a *api) handleUserAPIKeys() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		UserID := vars["id"]
		UserCookieID := r.Context().Value(contextKeyUserID).(string)
		if UserID != UserCookieID {
			Error(w, r, http.StatusForbidden, "")
			return
		}

		APIKeys, keysErr := a.db.GetUserAPIKeys(UserID)
		if keysErr != nil {
			Error(w, r, http.StatusInternalServerError, keysErr.Error())
			return
		}

		Success(w, r, http.StatusOK, APIKeys, nil)
	}
}

// handleAPIKeyGenerate handles generating an API key for a user
// @Summary Generate API Key
// @Description Generates an API key for the authenticated user
// @Tags apikey
// @Produce  json
// @Param id path int false "the user ID to generate API key for"
// @Success 200 object standardJsonResponse{data=model.APIKey}
// @Failure 403 object standardJsonResponse{}
// @Failure 500 object standardJsonResponse{}
// @Router /users/{id}/apikeys [post]
func (a *api) handleAPIKeyGenerate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		keyVal := getJSONRequestBody(r, w)
		APIKeyName := keyVal["name"].(string)

		UserID := vars["id"]
		UserCookieID := r.Context().Value(contextKeyUserID).(string)
		if UserID != UserCookieID {
			Error(w, r, http.StatusForbidden, "")
			return
		}

		APIKeys, keysErr := a.db.GetUserAPIKeys(UserID)
		if keysErr != nil {
			Error(w, r, http.StatusInternalServerError, keysErr.Error())
			return
		}

		if len(APIKeys) == a.config.UserAPIKeyLimit {
			Error(w, r, http.StatusForbidden, "USER_APIKEY_LIMIT_REACHED")
			return
		}

		APIKey, keyErr := a.db.GenerateAPIKey(UserID, APIKeyName)
		if keyErr != nil {
			Error(w, r, http.StatusInternalServerError, keyErr.Error())
			return
		}

		Success(w, r, http.StatusOK, APIKey, nil)
	}
}

// handleUserAPIKeyUpdate handles updating a users API key
// @Summary Update API Key
// @Description Updates the API key of the authenticated user
// @Tags apikey
// @Produce  json
// @Param id path int false "the user ID"
// @Param keyID path int false "the API Key ID to update"
// @Success 200 object standardJsonResponse{data=[]model.APIKey}
// @Failure 403 object standardJsonResponse{}
// @Failure 500 object standardJsonResponse{}
// @Router /users/{id}/apikeys/{keyID} [put]
func (a *api) handleUserAPIKeyUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		UserID := vars["id"]
		UserCookieID := r.Context().Value(contextKeyUserID).(string)
		if UserID != UserCookieID {
			Error(w, r, http.StatusForbidden, "")
			return
		}
		APK := vars["keyID"]
		keyVal := getJSONRequestBody(r, w)
		active := keyVal["active"].(bool)

		APIKeys, keysErr := a.db.UpdateUserAPIKey(UserID, APK, active)
		if keysErr != nil {
			Error(w, r, http.StatusInternalServerError, keysErr.Error())
			return
		}

		Success(w, r, http.StatusOK, APIKeys, nil)
	}
}

// handleUserAPIKeyDelete handles deleting a users API key
// @Summary Delete API Key
// @Description Deletes the API key
// @Tags apikey
// @Produce  json
// @Param id path int false "the user ID"
// @Param keyID path int false "the API Key ID to update"
// @Success 200 object standardJsonResponse{data=[]model.APIKey}
// @Failure 403 object standardJsonResponse{}
// @Failure 500 object standardJsonResponse{}
// @Router /users/{id}/apikeys/{keyID} [delete]
func (a *api) handleUserAPIKeyDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		UserID := vars["id"]
		UserCookieID := r.Context().Value(contextKeyUserID).(string)
		if UserID != UserCookieID {
			Error(w, r, http.StatusForbidden, "")
			return
		}
		APK := vars["keyID"]

		APIKeys, keysErr := a.db.DeleteUserAPIKey(UserID, APK)
		if keysErr != nil {
			Error(w, r, http.StatusInternalServerError, keysErr.Error())
			return
		}

		Success(w, r, http.StatusOK, APIKeys, nil)
	}
}
