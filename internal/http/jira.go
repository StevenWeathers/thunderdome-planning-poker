package http

import (
	"encoding/json"
	"io"
	"net/http"

	"go.uber.org/zap"

	"github.com/StevenWeathers/thunderdome-planning-poker/internal/atlassian/jira"

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
// @Router       /users/{userId}/jira-instances [get]
func (s *Service) handleGetUserJiraInstances() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		SessionUserID := ctx.Value(contextKeyUserID).(string)
		vars := mux.Vars(r)
		userId := vars["userId"]

		instances, err := s.JiraDataSvc.FindInstancesByUserId(ctx, userId)
		if err != nil {
			s.Logger.Ctx(ctx).Error(
				"handleGetUserJiraInstances error", zap.Error(err), zap.String("entity_user_id", userId),
				zap.String("session_user_id", SessionUserID))
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
// @Router       /users/{userId}/jira-instances [post]
func (s *Service) handleJiraInstanceCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		SessionUserID := ctx.Value(contextKeyUserID).(string)
		vars := mux.Vars(r)
		userId := vars["userId"]

		var req = jiraInstanceRequestBody{}
		body, bodyErr := io.ReadAll(r.Body)
		if bodyErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
			return
		}

		jsonErr := json.Unmarshal(body, &req)
		if jsonErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		inputErr := validate.Struct(req)
		if inputErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, inputErr.Error()))
			return
		}

		instance, err := s.JiraDataSvc.CreateInstance(ctx, userId, req.Host, req.ClientMail, req.AccessToken)
		if err != nil {
			s.Logger.Ctx(ctx).Error(
				"handleJiraInstanceCreate error", zap.Error(err), zap.String("entity_user_id", userId),
				zap.String("session_user_id", SessionUserID))
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
// @Router       /users/{userId}/jira-instances/{instanceId} [put]
func (s *Service) handleJiraInstanceUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		SessionUserID := ctx.Value(contextKeyUserID).(string)
		vars := mux.Vars(r)
		userId := vars["userId"]
		instanceId := vars["instanceId"]

		var req = jiraInstanceRequestBody{}
		body, bodyErr := io.ReadAll(r.Body)
		if bodyErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
			return
		}

		jsonErr := json.Unmarshal(body, &req)
		if jsonErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		inputErr := validate.Struct(req)
		if inputErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, inputErr.Error()))
			return
		}

		instance, err := s.JiraDataSvc.UpdateInstance(ctx, instanceId, req.Host, req.ClientMail, req.AccessToken)
		if err != nil {
			s.Logger.Ctx(ctx).Error(
				"handleJiraInstanceUpdate error", zap.Error(err), zap.String("entity_user_id", userId),
				zap.String("session_user_id", SessionUserID), zap.String("jira_instance_id", instanceId))
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
// @Router       /users/{userId}/jira-instances/{instanceId} [delete]
func (s *Service) handleJiraInstanceDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		SessionUserID := ctx.Value(contextKeyUserID).(string)
		vars := mux.Vars(r)
		instanceId := vars["instanceId"]
		userId := vars["userId"]

		err := s.JiraDataSvc.DeleteInstance(ctx, instanceId)
		if err != nil {
			s.Logger.Ctx(ctx).Error(
				"handleJiraInstanceDelete error", zap.Error(err), zap.String("entity_user_id", userId),
				zap.String("session_user_id", SessionUserID), zap.String("jira_instance_id", instanceId))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, nil, nil)
	}
}

type jiraStoryJQLSearchRequestBody struct {
	JQL        string `json:"jql" validate:"required"`
	StartAt    int    `json:"startAt"`
	MaxResults int    `json:"maxResults"`
}

// handleJiraStoryJQLSearch queries Jira API for Stories by JQL
// @Summary      Query Jira for Stories by JQL
// @Description  Queries Jira Instance API for Stories by JQL
// @Tags         jira
// @Produce      json
// @Param        userId  path    string                                          true  "the user ID associated to jira instance"
// @Param        instanceId  path    string                                          true  "the jira_instance ID to query"
// @Param        jira  body    jiraStoryJQLSearchRequestBody                                true  "jql search request"
// @Success      200    object  standardJsonResponse{}
// @Failure      500    object  standardJsonResponse{}
// @Security     ApiKeyAuth
// @Router       /users/{userId}/jira-instances/{instanceId}/jql-story-search [post]
func (s *Service) handleJiraStoryJQLSearch() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		ctx := r.Context()
		SessionUserID := ctx.Value(contextKeyUserID).(string)
		userId := vars["userId"]
		instanceId := vars["instanceId"]

		var req = jiraStoryJQLSearchRequestBody{}
		body, bodyErr := io.ReadAll(r.Body)
		if bodyErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
			return
		}

		jsonErr := json.Unmarshal(body, &req)
		if jsonErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		inputErr := validate.Struct(req)
		if inputErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, inputErr.Error()))
			return
		}
		fields := []string{"key", "summary", "priority", "issuetype", "description"}

		instance, err := s.JiraDataSvc.GetInstanceById(ctx, instanceId)
		if err != nil {
			s.Logger.Ctx(ctx).Error(
				"handleJiraStoryJQLSearch error", zap.Error(err), zap.String("entity_user_id", userId),
				zap.String("session_user_id", SessionUserID), zap.String("jira_instance_id", instanceId),
				zap.Int("start_at", req.StartAt), zap.Int("max_results", req.MaxResults),
				zap.Any("jira_fields", fields))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		jiraClient, err := jira.New(jira.Config{
			InstanceHost: instance.Host,
			ClientMail:   instance.ClientMail,
			AccessToken:  instance.AccessToken,
		})
		if err != nil {
			s.Logger.Ctx(ctx).Error(
				"handleJiraStoryJQLSearch error", zap.Error(err), zap.String("entity_user_id", userId),
				zap.String("session_user_id", SessionUserID), zap.String("jira_instance_id", instanceId),
				zap.Int("start_at", req.StartAt), zap.Int("max_results", req.MaxResults),
				zap.Any("jira_fields", fields))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		stories, err := jiraClient.StoriesJQLSearch(ctx, req.JQL, fields, req.StartAt, req.MaxResults)
		if err != nil {
			s.Logger.Ctx(ctx).Error(
				"handleJiraStoryJQLSearch error", zap.Error(err), zap.String("entity_user_id", userId),
				zap.String("session_user_id", SessionUserID), zap.String("jira_instance_id", instanceId),
				zap.Int("start_at", req.StartAt), zap.Int("max_results", req.MaxResults),
				zap.Any("jira_fields", fields))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, stories, nil)
	}
}
