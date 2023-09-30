package http

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

// handleGetUserJiraInstances gets a list of jira instances associated to user
// @Summary      Get User Jira Instances
// @Description  get list of Jira instances associated to user
// @Tags         jira
// @Produce      json
// @Param        userId  path    string                                          true  "the user ID to find jira instances for"
// @Success      200     object  standardJsonResponse{data=[]thunderdome.JiraInstance}
// @Failure      500     object  standardJsonResponse{}
// @Security     ApiKeyAuth
// @Router       /{userId}/jira-instances [get]
func (s *Service) handleGetUserJiraInstances() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userId := vars["userId"]

		instances, err := s.JiraDataSvc.FindInstancesByUserId(r.Context(), userId)
		if err != nil {
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, instances, nil)
	}
}

type jiraInstanceRequestBody struct {
	Host        string `json:"host" validate:"required"`
	ClientMail  string `json:"client_mail" validate:"required"`
	AccessToken string `json:"access_token" validate:"required"`
}

// handleJiraInstanceCreate creates a new Jira Instance
// @Summary      Create Jira Instance
// @Description  Creates a Jira Instance associated to user
// @Tags         jira
// @Produce      json
// @Param        userId  path    string                                          true  "the user ID to associate jira instance to"
// @Param        jira  body    jiraInstanceRequestBody                                true  "new jira_instance object"
// @Success      200    object  standardJsonResponse{data=thunderdome.JiraInstance}  "returns new jira instance"
// @Failure      500    object  standardJsonResponse{}
// @Security     ApiKeyAuth
// @Router       /{userId}/jira-instances [post]
func (s *Service) handleJiraInstanceCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userId := vars["userId"]

		var alert = jiraInstanceRequestBody{}
		body, bodyErr := io.ReadAll(r.Body)
		if bodyErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
			return
		}

		jsonErr := json.Unmarshal(body, &alert)
		if jsonErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		inputErr := validate.Struct(alert)
		if inputErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, inputErr.Error()))
			return
		}

		instance, err := s.JiraDataSvc.CreateInstance(r.Context(), userId, alert.Host, alert.ClientMail, alert.AccessToken)
		if err != nil {
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, instance, nil)
	}
}

// handleJiraInstanceUpdate updates a Jira Instance
// @Summary      Update Jira Instance
// @Description  Updates a Jira Instance associated to user
// @Tags         jira
// @Produce      json
// @Param        userId  path    string                                          true  "the user ID jira instance associated to"
// @Param        instanceId  path    string                                          true  "the jira_instance ID to update"
// @Param        jira  body    jiraInstanceRequestBody                                true  "updated jira_instance object"
// @Success      200    object  standardJsonResponse{data=thunderdome.JiraInstance}  "returns updated jira instance"
// @Failure      500    object  standardJsonResponse{}
// @Security     ApiKeyAuth
// @Router       /{userId}/jira-instances/{instanceId} [put]
func (s *Service) handleJiraInstanceUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		instanceId := vars["instanceId"]

		var alert = jiraInstanceRequestBody{}
		body, bodyErr := io.ReadAll(r.Body)
		if bodyErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
			return
		}

		jsonErr := json.Unmarshal(body, &alert)
		if jsonErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		inputErr := validate.Struct(alert)
		if inputErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, inputErr.Error()))
			return
		}

		instance, err := s.JiraDataSvc.UpdateInstance(r.Context(), instanceId, alert.Host, alert.ClientMail, alert.AccessToken)
		if err != nil {
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, instance, nil)
	}
}

// handleJiraInstanceDelete deletes a Jira Instance
// @Summary      Delete Jira Instance
// @Description  Deletes a Jira Instance associated to user
// @Tags         jira
// @Produce      json
// @Param        userId  path    string                                          true  "the user ID jira instance associated to"
// @Param        instanceId  path    string                                          true  "the jira_instance ID to delete"
// @Success      200    object  standardJsonResponse{}
// @Failure      500    object  standardJsonResponse{}
// @Security     ApiKeyAuth
// @Router       /{userId}/jira-instances/{instanceId} [delete]
func (s *Service) handleJiraInstanceDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		instanceId := vars["instanceId"]

		err := s.JiraDataSvc.DeleteInstance(r.Context(), instanceId)
		if err != nil {
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, nil, nil)
	}
}
