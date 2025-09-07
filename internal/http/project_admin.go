package http

import (
	"encoding/json"
	"io"
	"net/http"

	"go.uber.org/zap"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
)

type projectRequestBody struct {
	ProjectKey     string  `json:"projectKey" validate:"required,min=2,max=10"`
	Name           string  `json:"name" validate:"required,min=1,max=255"`
	Description    string  `json:"description"`
	OrganizationID *string `json:"organizationId"`
	DepartmentID   *string `json:"departmentId"`
	TeamID         *string `json:"teamId"`
}

// handleAdminGetProjects gets a list of projects
//
//	@Summary		Get Projects
//	@Description	get list of projects
//	@Tags			project, admin
//	@Produce		json
//	@Param			limit	query	int	false	"Max number of results to return"
//	@Param			offset	query	int	false	"Starting point to return rows from, should be multiplied by limit or 0"
//	@Success		200		object	standardJsonResponse{data=[]thunderdome.Project}
//	@Failure		500		object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/admin/projects [get]
func (s *Service) handleAdminGetProjects() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID, _ := ctx.Value(contextKeyUserID).(*string)
		limit, offset := getLimitOffsetFromRequest(r)
		projects, count, err := s.ProjectDataSvc.ListProjects(ctx, limit, offset)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleGetProjects error", zap.Error(err),
				zap.Int("limit", limit), zap.Int("offset", offset),
				zap.Stringp("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		meta := &pagination{
			Count:  count,
			Offset: offset,
			Limit:  limit,
		}

		s.Success(w, r, http.StatusOK, projects, meta)
	}
}

// handleAdminProjectCreate creates a new project
//
//	@Summary		Create Project
//	@Description	Creates a project
//	@Tags			project, admin
//	@Produce		json
//	@Param			project	body	projectRequestBody	true	"new project object"
//	@Success		200		object	standardJsonResponse{data=thunderdome.Project}
//	@Failure		400		object	standardJsonResponse{}
//	@Failure		500		object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/admin/projects [post]
func (s *Service) handleAdminProjectCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)
		var project = projectRequestBody{}
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

		// Validate that at least one association exists (organization, department, or team)
		if project.OrganizationID == nil && project.DepartmentID == nil && project.TeamID == nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, "project must be associated with at least one organization, department, or team"))
			return
		}

		newProject := &thunderdome.Project{
			ProjectKey:     project.ProjectKey,
			Name:           project.Name,
			Description:    project.Description,
			OrganizationID: project.OrganizationID,
			DepartmentID:   project.DepartmentID,
			TeamID:         project.TeamID,
		}

		err := s.ProjectDataSvc.CreateProject(ctx, newProject)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleProjectCreate error", zap.Error(err),
				zap.String("project_name", project.Name),
				zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, newProject, nil)
	}
}

// handleAdminProjectUpdate updates a project
//
//	@Summary		Update Project
//	@Description	Updates a Project
//	@Tags			project, admin
//	@Produce		json
//	@Param			projectId	path	string				true	"the project ID to update"
//	@Param			project		body	projectRequestBody	true	"project object to update"
//	@Success		200			object	standardJsonResponse{data=thunderdome.Project}
//	@Failure		400			object	standardJsonResponse{}
//	@Failure		500			object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/admin/projects/{projectId} [put]
func (s *Service) handleAdminProjectUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)

		projectID := r.PathValue("projectId")
		idErr := validate.Var(projectID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		var project = projectRequestBody{}
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

		// Validate that at least one association exists (organization, department, or team)
		if project.OrganizationID == nil && project.DepartmentID == nil && project.TeamID == nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, "project must be associated with at least one organization, department, or team"))
			return
		}

		updatedProject := &thunderdome.Project{
			ID:             projectID,
			ProjectKey:     project.ProjectKey,
			Name:           project.Name,
			Description:    project.Description,
			OrganizationID: project.OrganizationID,
			DepartmentID:   project.DepartmentID,
			TeamID:         project.TeamID,
		}

		err := s.ProjectDataSvc.UpdateProject(ctx, updatedProject)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleProjectUpdate error", zap.Error(err), zap.String("project_id", projectID),
				zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, updatedProject, nil)
	}
}

// handleAdminProjectDelete handles deleting a project
//
//	@Summary		Delete Project
//	@Description	Deletes a Project
//	@Tags			project, admin
//	@Produce		json
//	@Param			projectId	path	string	true	"the project ID to delete"
//	@Success		200			object	standardJsonResponse{}
//	@Failure		400			object	standardJsonResponse{}
//	@Failure		500			object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/admin/projects/{projectId} [delete]
func (s *Service) handleAdminProjectDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)

		projectID := r.PathValue("projectId")
		idErr := validate.Var(projectID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		err := s.ProjectDataSvc.DeleteProject(ctx, projectID)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleProjectDelete error", zap.Error(err), zap.String("project_id", projectID),
				zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, nil, nil)
	}
}

// handleAdminGetProjectByID gets a specific project by ID
//
//	@Summary		Get Project by ID
//	@Description	get a specific project by its ID
//	@Tags			project, admin
//	@Produce		json
//	@Param			projectId	path	string	true	"the project ID"
//	@Success		200			object	standardJsonResponse{data=thunderdome.Project}
//	@Failure		404			object	standardJsonResponse{}
//	@Failure		500			object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/admin/projects/{projectId} [get]
func (s *Service) handleAdminGetProjectByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID, _ := ctx.Value(contextKeyUserID).(*string)

		projectID := r.PathValue("projectId")
		idErr := validate.Var(projectID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		project, err := s.ProjectDataSvc.GetProjectByID(ctx, projectID)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleGetProjectByID error", zap.Error(err), zap.String("project_id", projectID),
				zap.Stringp("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		if project == nil {
			s.Failure(w, r, http.StatusNotFound, Errorf(ENOTFOUND, "Project not found"))
			return
		}

		s.Success(w, r, http.StatusOK, project, nil)
	}
}
