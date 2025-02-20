package http

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"go.uber.org/zap"

	jira "github.com/StevenWeathers/thunderdome-planning-poker/internal/atlassian/jira"
	jira_data_center "github.com/StevenWeathers/thunderdome-planning-poker/internal/atlassian/jiraDataCenter"
	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
	"github.com/gorilla/mux"
)

// handleGetUserJiraInstances gets a list of jira instances associated to user
//
//	@Summary		Get User Jira Instances
//	@Description	get list of Jira instances associated to user
//	@Tags			jira
//	@Produce		json
//	@Param			userId	path	string	true	"the user ID to find jira instances for"
//	@Success		200		object	standardJsonResponse{data=[]thunderdome.JiraInstance}
//	@Failure		500		object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/users/{userId}/jira-instances [get]
func (s *Service) handleGetUserJiraInstances() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)
		vars := mux.Vars(r)
		userID := vars["userId"]

		idErr := validate.Var(userID, "required,uuid")
		s.ErrCheck(idErr, w, r)

		instances, err := s.JiraDataSvc.FindInstancesByUserID(ctx, userID)
		if err != nil {
			s.Logger.Ctx(ctx).Error(
				"handleGetUserJiraInstances error", zap.Error(err), zap.String("entity_user_id", userID),
				zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, instances, nil)
	}
}

type jiraInstanceRequestBody struct {
	Host           string `json:"host" validate:"required,http_url"`
	ClientMail     string `json:"client_mail" validate:"required,email"`
	AccessToken    string `json:"access_token" validate:"required"`
	JiraDataCenter bool   `json:"jira_data_center"` // Checkbox for enabling Jira Data Center
}

// handleJiraInstanceCreate creates a new Jira Instance
//
//	@Summary		Create Jira Instance
//	@Description	Creates a Jira Instance associated to user
//	@Tags			jira
//	@Produce		json
//	@Param			userId	path	string												true	"the user ID to associate jira instance to"
//	@Param			jira	body	jiraInstanceRequestBody								true	"new jira_instance object"
//	@Success		200		object	standardJsonResponse{data=thunderdome.JiraInstance}	"returns new jira instance"
//	@Failure		500		object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/users/{userId}/jira-instances [post]
func (s *Service) handleJiraInstanceCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)
		vars := mux.Vars(r)
		userID := vars["userId"]

		idErr := validate.Var(userID, "required,uuid")
		s.ErrCheck(idErr, w, r)

		var req = jiraInstanceRequestBody{}
		body, bodyErr := io.ReadAll(r.Body)
		s.ErrCheck(bodyErr, w, r)

		jsonErr := json.Unmarshal(body, &req)
		s.ErrCheck(jsonErr, w, r)

		inputErr := validate.Struct(req)
		s.ErrCheck(inputErr, w, r)

		instance, err := s.JiraDataSvc.CreateInstance(ctx, userID, req.Host, req.ClientMail, req.AccessToken, req.JiraDataCenter)
		if err != nil {
			s.Logger.Ctx(ctx).Error(
				"handleJiraInstanceCreate error", zap.Error(err), zap.String("entity_user_id", userID),
				zap.String("session_user_id", sessionUserID), zap.Bool("jira_data_center", req.JiraDataCenter), zap.Stack("stacktrace"))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, instance, nil)
	}
}

// handleJiraInstanceUpdate updates a Jira Instance
//
//	@Summary		Update Jira Instance
//	@Description	Updates a Jira Instance associated to user
//	@Tags			jira
//	@Produce		json
//	@Param			userId		path	string												true	"the user ID jira instance associated to"
//	@Param			instanceId	path	string												true	"the jira_instance ID to update"
//	@Param			jira		body	jiraInstanceRequestBody								true	"updated jira_instance object"
//	@Success		200			object	standardJsonResponse{data=thunderdome.JiraInstance}	"returns updated jira instance"
//	@Failure		500			object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/users/{userId}/jira-instances/{instanceId} [put]
func (s *Service) handleJiraInstanceUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)
		vars := mux.Vars(r)
		userID := vars["userId"]
		instanceID := vars["instanceId"]

		iidErr := validate.Var(instanceID, "required,uuid")
		s.ErrCheck(iidErr, w, r)

		idErr := validate.Var(userID, "required,uuid")
		s.ErrCheck(idErr, w, r)

		var req = jiraInstanceRequestBody{}
		body, bodyErr := io.ReadAll(r.Body)
		s.ErrCheck(bodyErr, w, r)

		jsonErr := json.Unmarshal(body, &req)
		s.ErrCheck(jsonErr, w, r)

		inputErr := validate.Struct(req)
		s.ErrCheck(inputErr, w, r)

		instance, err := s.JiraDataSvc.UpdateInstance(ctx, instanceID, req.Host, req.ClientMail, req.AccessToken)
		if err != nil {
			s.Logger.Ctx(ctx).Error(
				"handleJiraInstanceUpdate error", zap.Error(err), zap.String("entity_user_id", userID),
				zap.String("session_user_id", sessionUserID), zap.String("jira_instance_id", instanceID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, instance, nil)
	}
}

// handleJiraInstanceDelete deletes a Jira Instance
//
//	@Summary		Delete Jira Instance
//	@Description	Deletes a Jira Instance associated to user
//	@Tags			jira
//	@Produce		json
//	@Param			userId		path	string	true	"the user ID jira instance associated to"
//	@Param			instanceId	path	string	true	"the jira_instance ID to delete"
//	@Success		200			object	standardJsonResponse{}
//	@Failure		500			object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/users/{userId}/jira-instances/{instanceId} [delete]
func (s *Service) handleJiraInstanceDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)
		vars := mux.Vars(r)
		instanceID := vars["instanceId"]

		iidErr := validate.Var(instanceID, "required,uuid")
		s.ErrCheck(iidErr, w, r)

		userID := vars["userId"]

		idErr := validate.Var(userID, "required,uuid")
		s.ErrCheck(idErr, w, r)

		err := s.JiraDataSvc.DeleteInstance(ctx, instanceID)
		if err != nil {
			s.Logger.Ctx(ctx).Error(
				"handleJiraInstanceDelete error", zap.Error(err), zap.String("entity_user_id", userID),
				zap.String("session_user_id", sessionUserID), zap.String("jira_instance_id", instanceID))
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
//
//	@Summary		Query Jira for Stories by JQL
//	@Description	Queries Jira Instance API for Stories by JQL
//	@Tags			jira
//	@Produce		json
//	@Param			userId		path	string							true	"the user ID associated to jira instance"
//	@Param			instanceId	path	string							true	"the jira_instance ID to query"
//	@Param			jira		body	jiraStoryJQLSearchRequestBody	true	"jql search request"
//	@Success		200			object	standardJsonResponse{}
//	@Failure		500			object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/users/{userId}/jira-instances/{instanceId}/jql-story-search [post]
func (s *Service) handleJiraStoryJQLSearch() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		ctx := r.Context()

		userID := vars["userId"]
		idErr := validate.Var(userID, "required,uuid")
		s.ErrCheck(idErr, w, r)

		instanceID := vars["instanceId"]
		iidErr := validate.Var(instanceID, "required,uuid")
		s.ErrCheck(iidErr, w, r)

		var req = jiraStoryJQLSearchRequestBody{}
		body, bodyErr := io.ReadAll(r.Body)
		s.ErrCheck(bodyErr, w, r)

		jsonErr := json.Unmarshal(body, &req)
		s.ErrCheck(jsonErr, w, r)

		inputErr := validate.Struct(req)
		s.ErrCheck(inputErr, w, r)

		fields := []string{"key", "summary", "priority", "issuetype", "description"}

		instance, err := s.JiraDataSvc.GetInstanceByID(ctx, instanceID)
		errorTitle := "handleJiraStoryJQLSearch error"

		s.logJiraSearchError(err, errorTitle, w, r, ctx, vars, fields, req)

		// check here for DataCenter
		if instance.JiraDataCenter {

			jiraDataCenterClient, err := CreateNewJiraDataCenterInstance(instance)
			s.logJiraSearchError(err, errorTitle, w, r, ctx, vars, fields, req)

			stories, err := jiraDataCenterClient.StoriesJQLSearch(ctx, req.JQL, fields, req.StartAt, req.MaxResults)

			s.logErrorWithJSONResponse(err, errorTitle, w, ctx, vars, fields, req)
			s.Success(w, r, http.StatusOK, stories, nil)

		} else {

			jiraClient, err := CreateNewJiraInstance(instance)
			s.logJiraSearchError(err, errorTitle, w, r, ctx, vars, fields, req)

			stories, err := jiraClient.StoriesJQLSearch(ctx, req.JQL, fields, req.StartAt, req.MaxResults)

			s.logErrorWithJSONResponse(err, errorTitle, w, ctx, vars, fields, req)
			s.Success(w, r, http.StatusOK, stories, nil)
		}

	}

}
func CreateNewJiraDataCenterInstance(instance thunderdome.JiraInstance) (*jira_data_center.Client, error) {

	jiraClient, err := jira_data_center.New(jira_data_center.Config{
		InstanceHost:   instance.Host,
		ClientMail:     instance.ClientMail,
		JiraDataCenter: instance.JiraDataCenter,
		AccessToken:    instance.AccessToken,
	})
	return jiraClient, err
}
func CreateNewJiraInstance(instance thunderdome.JiraInstance) (*jira.Client, error) {

	jiraClient, err := jira.New(jira.Config{
		InstanceHost:   instance.Host,
		ClientMail:     instance.ClientMail,
		JiraDataCenter: instance.JiraDataCenter,
		AccessToken:    instance.AccessToken,
	})
	return jiraClient, err
}

func (s *Service) ErrCheck(err error, w http.ResponseWriter, r *http.Request) {
	if err != nil {
		s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, err.Error()))
		return
	}
}

func (s *Service) logJiraSearchError(err error, errorTitle string, w http.ResponseWriter, r *http.Request, ctx context.Context, vars map[string]string, fields []string, req jiraStoryJQLSearchRequestBody) {
	if err != nil {
		s.createLoggerStructure(err, errorTitle, ctx, vars, fields, req)
		s.Failure(w, r, http.StatusInternalServerError, err)
		return
	}
}

func (s *Service) logErrorWithJSONResponse(err error, errorTitle string, w http.ResponseWriter, ctx context.Context, vars map[string]string, fields []string, req jiraStoryJQLSearchRequestBody) {

	if err != nil {
		s.createLoggerStructure(err, errorTitle, ctx, vars, fields, req)

		result := &standardJsonResponse{
			Success: false,
			Error:   err.Error(),
			Data:    map[string]interface{}{},
			Meta:    map[string]interface{}{},
		}

		response, _ := json.Marshal(result)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(response)
		return
	}

}

func (s *Service) createLoggerStructure(err error, errorTitle string, ctx context.Context, vars map[string]string, fields []string, req jiraStoryJQLSearchRequestBody) {
	s.Logger.Ctx(ctx).Error(
		errorTitle, zap.Error(err), zap.String("entity_user_id", vars["userID"]),
		zap.String("session_user_id", ctx.Value(contextKeyUserID).(string)), zap.String("jira_instance_id", vars["instanceID"]),
		zap.Int("start_at", req.StartAt), zap.Int("max_results", req.MaxResults),
		zap.Any("jira_fields", fields))
}
