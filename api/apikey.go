package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

// handleUserAPIKeys handles getting user API keys
// @Summary Get API Keys
// @Description get list of API keys for the user
// @Tags apikey
// @Produce  json
// @Param userId path string true "the user ID to get API keys for"
// @Success 200 object standardJsonResponse{data=[]model.APIKey}
// @Failure 403 object standardJsonResponse{}
// @Failure 500 object standardJsonResponse{}
// @Router /users/{userId}/apikeys [get]
func (a *api) handleUserAPIKeys() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		UserID := vars["userId"]

		APIKeys, keysErr := a.db.GetUserApiKeys(UserID)
		if keysErr != nil {
			Failure(w, r, http.StatusInternalServerError, keysErr)
			return
		}

		Success(w, r, http.StatusOK, APIKeys, nil)
	}
}

// handleAPIKeyGenerate handles generating an API key for a user
// @Summary Generate API Key
// @Description Generates an API key for the user
// @Tags apikey
// @Produce  json
// @Param userId path string true "the user ID to generate API key for"
// @Param name body string true "Name the API Key to distinguish its use"
// @Success 200 object standardJsonResponse{data=model.APIKey}
// @Failure 403 object standardJsonResponse{}
// @Failure 500 object standardJsonResponse{}
// @Router /users/{userId}/apikeys [post]
func (a *api) handleAPIKeyGenerate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		keyVal := getJSONRequestBody(r, w)
		APIKeyName := keyVal["name"].(string)
		UserID := vars["userId"]

		APIKeys, keysErr := a.db.GetUserApiKeys(UserID)
		if keysErr != nil {
			Failure(w, r, http.StatusInternalServerError, keysErr)
			return
		}

		if len(APIKeys) == a.config.UserAPIKeyLimit {
			Failure(w, r, http.StatusForbidden, Errorf(EINVALID, "USER_APIKEY_LIMIT_REACHED"))
			return
		}

		APIKey, keyErr := a.db.GenerateApiKey(UserID, APIKeyName)
		if keyErr != nil {
			Failure(w, r, http.StatusInternalServerError, keyErr)
			return
		}

		Success(w, r, http.StatusOK, APIKey, nil)
	}
}

// handleUserAPIKeyUpdate handles updating a users API key
// @Summary Update API Key
// @Description Updates the API key of the user
// @Tags apikey
// @Produce  json
// @Param userId path string true "the user ID"
// @Param keyID path string true "the API Key ID to update"
// @Param active body boolean true "Whether the API Key is enabled for use"
// @Success 200 object standardJsonResponse{data=[]model.APIKey}
// @Failure 403 object standardJsonResponse{}
// @Failure 500 object standardJsonResponse{}
// @Router /users/{userId}/apikeys/{keyID} [put]
func (a *api) handleUserAPIKeyUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		UserID := vars["userId"]
		APK := vars["keyID"]
		keyVal := getJSONRequestBody(r, w)
		active := keyVal["active"].(bool)

		APIKeys, keysErr := a.db.UpdateUserApiKey(UserID, APK, active)
		if keysErr != nil {
			Failure(w, r, http.StatusInternalServerError, keysErr)
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
// @Param userId path int false "the user ID"
// @Param keyID path int false "the API Key ID to update"
// @Success 200 object standardJsonResponse{data=[]model.APIKey}
// @Failure 403 object standardJsonResponse{}
// @Failure 500 object standardJsonResponse{}
// @Router /users/{userId}/apikeys/{keyID} [delete]
func (a *api) handleUserAPIKeyDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		UserID := vars["userId"]
		APK := vars["keyID"]

		APIKeys, keysErr := a.db.DeleteUserApiKey(UserID, APK)
		if keysErr != nil {
			Failure(w, r, http.StatusInternalServerError, keysErr)
			return
		}

		Success(w, r, http.StatusOK, APIKeys, nil)
	}
}
