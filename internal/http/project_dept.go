package http

import (
	"encoding/json"
	"io"
	"net/http"

	"go.uber.org/zap"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
)

// handleGetDepartmentProjects gets a list of projects for a department
//
//	@Summary		Get Department Projects
//	@Description	get list of projects for a department
//	@Tags			project
//	@Produce		json
//	@Param			orgId			path		string	true	"Organization ID"
//	@Param			departmentId	path		string	true	"Department ID"
//	@Success		200				{object}	standardJsonResponse{data=[]thunderdome.Project}
//	@Failure		400				{object}	standardJsonResponse{}
//	@Failure		500				{object}	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/organizations/{orgId}/departments/{departmentId}/projects [get]
func (s *Service) handleGetDepartmentProjects() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID, _ := ctx.Value(contextKeyUserID).(*string)

		departmentID := r.PathValue("departmentId")
		deptIDErr := validate.Var(departmentID, "required,uuid")
		if deptIDErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, deptIDErr.Error()))
			return
		}

		projects, err := s.ProjectDataSvc.GetProjectsByDepartment(ctx, departmentID)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleGetDepartmentProjects error", zap.Error(err),
				zap.String("department_id", departmentID),
				zap.Stringp("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, projects, nil)
	}
}

// handleDepartmentProjectCreate creates a new department project
//
//	@Summary		Create Department Project
//	@Description	Creates a Department project
//	@Tags			project
//	@Produce		json
//	@Param			orgId			path	string						true	"the organization ID"
//	@Param			departmentId	path	string						true	"the department ID"
//	@Param			project			body	scopedProjectRequestBody	true	"new project object"
//	@Success		200				object	standardJsonResponse{data=thunderdome.Project}
//	@Failure		400				object	standardJsonResponse{}
//	@Failure		500				object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/organizations/{orgId}/departments/{departmentId}/projects [post]
func (s *Service) handleDepartmentProjectCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)

		deptID := r.PathValue("departmentId")
		deptIDErr := validate.Var(deptID, "required,uuid")
		if deptIDErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, deptIDErr.Error()))
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
			ProjectKey:   project.ProjectKey,
			Name:         project.Name,
			Description:  project.Description,
			DepartmentID: &deptID,
		}

		err := s.ProjectDataSvc.CreateProject(ctx, newProject)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleDepartmentProjectCreate error", zap.Error(err),
				zap.String("project_name", project.Name),
				zap.String("department_id", deptID),
				zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, newProject, nil)
	}
}

// handleDepartmentProjectUpdate updates a department project
//
//	@Summary		Update Department Project
//	@Description	Updates a Department Project
//	@Tags			project
//	@Produce		json
//	@Param			orgId			path	string						true	"the organization ID"
//	@Param			departmentId	path	string						true	"the department ID"
//	@Param			projectId		path	string						true	"the project ID to update"
//	@Param			project			body	scopedProjectRequestBody	true	"project object to update"
//	@Success		200				object	standardJsonResponse{data=thunderdome.Project}
//	@Failure		400				object	standardJsonResponse{}
//	@Failure		500				object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/organizations/{orgId}/departments/{departmentId}/projects/{projectId} [put]
func (s *Service) handleDepartmentProjectUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)

		projectID := r.PathValue("projectId")
		idErr := validate.Var(projectID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		deptID := r.PathValue("departmentId")
		deptIDErr := validate.Var(deptID, "required,uuid")
		if deptIDErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, deptIDErr.Error()))
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
			ID:           projectID,
			ProjectKey:   project.ProjectKey,
			Name:         project.Name,
			Description:  project.Description,
			DepartmentID: &deptID,
		}

		err := s.ProjectDataSvc.UpdateDepartmentProject(ctx, updatedProject)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleDepartmentProjectUpdate error", zap.Error(err),
				zap.String("project_id", projectID),
				zap.String("department_id", deptID),
				zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, updatedProject, nil)
	}
}

// handleDepartmentProjectDelete handles deleting a department project
//
//	@Summary		Delete Department Project
//	@Description	Deletes a Department Project
//	@Tags			project
//	@Produce		json
//	@Param			orgId			path	string	true	"the organization ID"
//	@Param			departmentId	path	string	true	"the department ID"
//	@Param			projectId		path	string	true	"the project ID to delete"
//	@Success		200				object	standardJsonResponse{}
//	@Failure		400				object	standardJsonResponse{}
//	@Failure		500				object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/organizations/{orgId}/departments/{departmentId}/projects/{projectId} [delete]
func (s *Service) handleDepartmentProjectDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)

		projectID := r.PathValue("projectId")
		idErr := validate.Var(projectID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		departmentID := r.PathValue("departmentId")
		deptIDErr := validate.Var(departmentID, "required,uuid")
		if deptIDErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, deptIDErr.Error()))
			return
		}

		err := s.ProjectDataSvc.DeleteDepartmentProject(ctx, departmentID, projectID)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleDepartmentProjectDelete error", zap.Error(err),
				zap.String("project_id", projectID),
				zap.String("department_id", departmentID),
				zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, nil, nil)
	}
}
