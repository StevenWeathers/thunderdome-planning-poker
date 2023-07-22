package http

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

// handleUserAPIKeys handles getting user API keys
// @Summary      Get API Keys
// @Description  get list of API keys for the user
// @Tags         apikey
// @Produce      json
// @Param        userId  path    string  true  "the user ID to get API keys for"
// @Success      200     object  standardJsonResponse{data=[]thunderdome.APIKey}
// @Failure      403     object  standardJsonResponse{}
// @Failure      500     object  standardJsonResponse{}
// @Security     ApiKeyAuth
// @Router       /users/{userId}/apikeys [get]
func (s *Service) handleUserAPIKeys() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		UserID := vars["userId"]
		idErr := validate.Var(UserID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		APIKeys, keysErr := s.ApiKeyDataSvc.GetUserApiKeys(r.Context(), UserID)
		if keysErr != nil {
			s.Failure(w, r, http.StatusInternalServerError, keysErr)
			return
		}

		s.Success(w, r, http.StatusOK, APIKeys, nil)
	}
}

type apikeyGenerateRequestBody struct {
	Name string `json:"name" validate:"required"`
}

// handleAPIKeyGenerate handles generating an API key for a user
// @Summary      Generate API Key
// @Description  Generates an API key for the user
// @Tags         apikey
// @Produce      json
// @Param        userId  path    string                     true  "the user ID to generate API key for"
// @Param        key     body    apikeyGenerateRequestBody  true  "new APIKey key object"
// @Success      200     object  standardJsonResponse{data=thunderdome.APIKey}
// @Failure      403     object  standardJsonResponse{}
// @Failure      500     object  standardJsonResponse{}
// @Security     ApiKeyAuth
// @Router       /users/{userId}/apikeys [post]
func (s *Service) handleAPIKeyGenerate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		UserID := vars["userId"]
		ctx := r.Context()
		idErr := validate.Var(UserID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		body, bodyErr := io.ReadAll(r.Body)
		if bodyErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
			return
		}

		var k = apikeyGenerateRequestBody{}
		jsonErr := json.Unmarshal(body, &k)
		if jsonErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		inputErr := validate.Struct(k)
		if inputErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, inputErr.Error()))
			return
		}

		APIKeys, keysErr := s.ApiKeyDataSvc.GetUserApiKeys(ctx, UserID)
		if keysErr != nil {
			s.Failure(w, r, http.StatusInternalServerError, keysErr)
			return
		}

		if len(APIKeys) == s.Config.UserAPIKeyLimit {
			s.Failure(w, r, http.StatusForbidden, Errorf(EINVALID, "USER_APIKEY_LIMIT_REACHED"))
			return
		}

		APIKey, keyErr := s.ApiKeyDataSvc.GenerateApiKey(ctx, UserID, k.Name)
		if keyErr != nil {
			s.Failure(w, r, http.StatusInternalServerError, keyErr)
			return
		}

		s.Success(w, r, http.StatusOK, APIKey, nil)
	}
}

type apikeyUpdateRequestBody struct {
	Active bool `json:"active"`
}

// handleUserAPIKeyUpdate handles updating a users API key
// @Summary      Update API Key
// @Description  Updates the API key of the user
// @Tags         apikey
// @Produce      json
// @Param        userId  path    string                   true  "the user ID"
// @Param        keyID   path    string                   true  "the API Key ID to update"
// @Param        key     body    apikeyUpdateRequestBody  true  "APIKey key object to update"
// @Success      200     object  standardJsonResponse{data=[]thunderdome.APIKey}
// @Failure      403     object  standardJsonResponse{}
// @Failure      500     object  standardJsonResponse{}
// @Security     ApiKeyAuth
// @Router       /users/{userId}/apikeys/{keyID} [put]
func (s *Service) handleUserAPIKeyUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		UserID := vars["userId"]
		idErr := validate.Var(UserID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		APK := vars["keyID"]

		body, bodyErr := io.ReadAll(r.Body)
		if bodyErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
			return
		}

		var k = apikeyUpdateRequestBody{}
		jsonErr := json.Unmarshal(body, &k)
		if jsonErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		inputErr := validate.Struct(k)
		if inputErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, inputErr.Error()))
			return
		}

		APIKeys, keysErr := s.ApiKeyDataSvc.UpdateUserApiKey(r.Context(), UserID, APK, k.Active)
		if keysErr != nil {
			s.Failure(w, r, http.StatusInternalServerError, keysErr)
			return
		}

		s.Success(w, r, http.StatusOK, APIKeys, nil)
	}
}

// handleUserAPIKeyDelete handles deleting a users API key
// @Summary      Delete API Key
// @Description  Deletes the API key
// @Tags         apikey
// @Produce      json
// @Param        userId  path    int  false  "the user ID"
// @Param        keyID   path    int  false  "the API Key ID to update"
// @Success      200     object  standardJsonResponse{data=[]thunderdome.APIKey}
// @Failure      403     object  standardJsonResponse{}
// @Failure      500     object  standardJsonResponse{}
// @Security     ApiKeyAuth
// @Router       /users/{userId}/apikeys/{keyID} [delete]
func (s *Service) handleUserAPIKeyDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		UserID := vars["userId"]
		idErr := validate.Var(UserID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		APK := vars["keyID"]

		APIKeys, keysErr := s.ApiKeyDataSvc.DeleteUserApiKey(r.Context(), UserID, APK)
		if keysErr != nil {
			s.Failure(w, r, http.StatusInternalServerError, keysErr)
			return
		}

		s.Success(w, r, http.StatusOK, APIKeys, nil)
	}
}
