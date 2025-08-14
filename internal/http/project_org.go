package http

import (
	"encoding/json"
	"io"
	"net/http"

	"go.uber.org/zap"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
)

type scopedProjectRequestBody struct {
	ProjectKey  string `json:"projectKey" validate:"required,min=2,max=10"`
	Name        string `json:"name" validate:"required,min=1,max=255"`
	Description string `json:"description"`
}

// handleGetOrganizationProjects gets a list of projects for an organization
//
//	@Summary		Get Organization Projects
//	@Description	get list of projects for an organization
//	@Tags			project
//	@Produce		json
//	@Param			orgId	path		string	true	"Organization ID"
//	@Success		200		{object}	standardJsonResponse{data=[]thunderdome.Project}
//	@Failure		400		{object}	standardJsonResponse{}
//	@Failure		500		{object}	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/organizations/{orgId}/projects [get]
func (s *Service) handleGetOrganizationProjects() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID, _ := ctx.Value(contextKeyUserID).(*string)

		organizationID := r.PathValue("orgId")
		orgIDErr := validate.Var(organizationID, "required,uuid")
		if orgIDErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, orgIDErr.Error()))
			return
		}

		projects, err := s.ProjectDataSvc.GetProjectsByOrganization(ctx, organizationID)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleGetOrganizationProjects error", zap.Error(err),
				zap.String("organization_id", organizationID),
				zap.Stringp("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, projects, nil)
	}
}

// handleOrganizationProjectCreate creates a new organization project
//
//	@Summary		Create Organization Project
//	@Description	Creates an Organization project
//	@Tags			project
//	@Produce		json
//	@Param			orgId	path	string						true	"the organization ID"
//	@Param			project	body	scopedProjectRequestBody	true	"new project object"
//	@Success		200		object	standardJsonResponse{data=thunderdome.Project}
//	@Failure		400		object	standardJsonResponse{}
//	@Failure		500		object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/organizations/{orgId}/projects [post]
func (s *Service) handleOrganizationProjectCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)

		orgID := r.PathValue("orgId")
		orgIDErr := validate.Var(orgID, "required,uuid")
		if orgIDErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, orgIDErr.Error()))
			return
		}

		var project = scopedProjectRequestBody{}
		body, bodyErr := io.ReadAll(r.Body)
		if bodyErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
			return
		}

		jsonErr := json.Unmarshal(body, &project)
		if jsonErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		inputErr := validate.Struct(project)
		if inputErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, inputErr.Error()))
			return
		}

		newProject := &thunderdome.Project{
			ProjectKey:     project.ProjectKey,
			Name:           project.Name,
			Description:    project.Description,
			OrganizationID: &orgID,
		}

		err := s.ProjectDataSvc.CreateProject(ctx, newProject)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleOrganizationProjectCreate error", zap.Error(err),
				zap.String("project_name", project.Name),
				zap.String("organization_id", orgID),
				zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, newProject, nil)
	}
}

// handleOrganizationProjectUpdate updates an organization project
//
//	@Summary		Update Organization Project
//	@Description	Updates an Organization Project
//	@Tags			project
//	@Produce		json
//	@Param			orgId		path	string						true	"the organization ID"
//	@Param			projectId	path	string						true	"the project ID to update"
//	@Param			project		body	scopedProjectRequestBody	true	"project object to update"
//	@Success		200			object	standardJsonResponse{data=thunderdome.Project}
//	@Failure		400			object	standardJsonResponse{}
//	@Failure		500			object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/organizations/{orgId}/projects/{projectId} [put]
func (s *Service) handleOrganizationProjectUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)

		projectID := r.PathValue("projectId")
		idErr := validate.Var(projectID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		orgID := r.PathValue("orgId")
		orgIDErr := validate.Var(orgID, "required,uuid")
		if orgIDErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, orgIDErr.Error()))
			return
		}

		var project = scopedProjectRequestBody{}
		body, bodyErr := io.ReadAll(r.Body)
		if bodyErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
			return
		}

		jsonErr := json.Unmarshal(body, &project)
		if jsonErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		inputErr := validate.Struct(project)
		if inputErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, inputErr.Error()))
			return
		}

		updatedProject := &thunderdome.Project{
			ID:             projectID,
			ProjectKey:     project.ProjectKey,
			Name:           project.Name,
			Description:    project.Description,
			OrganizationID: &orgID,
		}

		err := s.ProjectDataSvc.UpdateOrganizationProject(ctx, updatedProject)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleOrganizationProjectUpdate error", zap.Error(err),
				zap.String("project_id", projectID),
				zap.String("organization_id", orgID),
				zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, updatedProject, nil)
	}
}

// handleOrganizationProjectDelete handles deleting an organization project
//
//	@Summary		Delete Organization Project
//	@Description	Deletes an Organization Project
//	@Tags			project
//	@Produce		json
//	@Param			orgId		path	string	true	"the Organization ID"
//	@Param			projectId	path	string	true	"the project ID to delete"
//	@Success		200			object	standardJsonResponse{}
//	@Failure		400			object	standardJsonResponse{}
//	@Failure		500			object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/organizations/{orgId}/projects/{projectId} [delete]
func (s *Service) handleOrganizationProjectDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)

		projectID := r.PathValue("projectId")
		idErr := validate.Var(projectID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		organizationID := r.PathValue("orgId")
		oidErr := validate.Var(organizationID, "required,uuid")
		if oidErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, oidErr.Error()))
			return
		}

		err := s.ProjectDataSvc.DeleteOrganizationProject(ctx, organizationID, projectID)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleOrganizationProjectDelete error", zap.Error(err),
				zap.String("project_id", projectID),
				zap.String("organization_id", organizationID),
				zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, nil, nil)
	}
}
