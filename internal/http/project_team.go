package http

import (
	"encoding/json"
	"io"
	"net/http"

	"go.uber.org/zap"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
)

// handleGetTeamProjects gets a list of projects for a team (direct team route)
//
//	@Summary		Get Team Projects
//	@Description	get list of projects for a team via direct team route
//	@Tags			project
//	@Produce		json
//	@Param			teamId	path		string	true	"Team ID"
//	@Success		200		{object}	standardJsonResponse{data=[]thunderdome.Project}
//	@Failure		400		{object}	standardJsonResponse{}
//	@Failure		500		{object}	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/teams/{teamId}/projects [get]
func (s *Service) handleGetTeamProjects() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID, _ := ctx.Value(contextKeyUserID).(*string)

		teamID := r.PathValue("teamId")
		teamIDErr := validate.Var(teamID, "required,uuid")
		if teamIDErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, teamIDErr.Error()))
			return
		}

		projects, err := s.ProjectDataSvc.GetProjectsByTeam(ctx, teamID)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleGetTeamProjects error", zap.Error(err),
				zap.String("team_id", teamID),
				zap.Stringp("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, projects, nil)
	}
}

// handleTeamProjectCreate creates a new team project (direct team route)
//
//	@Summary		Create Team Project
//	@Description	Creates a Team project via direct team route
//	@Tags			project
//	@Produce		json
//	@Param			teamId	path	string						true	"the team ID"
//	@Param			project	body	scopedProjectRequestBody	true	"new project object"
//	@Success		200		object	standardJsonResponse{data=thunderdome.Project}
//	@Failure		400		object	standardJsonResponse{}
//	@Failure		500		object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/teams/{teamId}/projects [post]
func (s *Service) handleTeamProjectCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)

		teamID := r.PathValue("teamId")
		teamIDErr := validate.Var(teamID, "required,uuid")
		if teamIDErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, teamIDErr.Error()))
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
			ProjectKey:  project.ProjectKey,
			Name:        project.Name,
			Description: project.Description,
			TeamID:      &teamID,
		}

		err := s.ProjectDataSvc.CreateProject(ctx, newProject)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleTeamProjectCreate error", zap.Error(err),
				zap.String("project_name", project.Name),
				zap.String("team_id", teamID),
				zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, newProject, nil)
	}
}

// handleTeamProjectUpdate updates a team project (direct team route)
//
//	@Summary		Update Team Project
//	@Description	Updates a Team Project via direct team route
//	@Tags			project
//	@Produce		json
//	@Param			teamId		path	string						true	"the team ID"
//	@Param			projectId	path	string						true	"the project ID to update"
//	@Param			project		body	scopedProjectRequestBody	true	"project object to update"
//	@Success		200			object	standardJsonResponse{data=thunderdome.Project}
//	@Failure		400			object	standardJsonResponse{}
//	@Failure		500			object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/teams/{teamId}/projects/{projectId} [put]
func (s *Service) handleTeamProjectUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)

		projectID := r.PathValue("projectId")
		idErr := validate.Var(projectID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		teamID := r.PathValue("teamId")
		teamIDErr := validate.Var(teamID, "required,uuid")
		if teamIDErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, teamIDErr.Error()))
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
			ID:          projectID,
			ProjectKey:  project.ProjectKey,
			Name:        project.Name,
			Description: project.Description,
			TeamID:      &teamID,
		}

		err := s.ProjectDataSvc.UpdateTeamProject(ctx, updatedProject)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleTeamProjectUpdate error", zap.Error(err),
				zap.String("project_id", projectID),
				zap.String("team_id", teamID),
				zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, updatedProject, nil)
	}
}

// handleTeamProjectDelete handles deleting a team project (direct team route)
//
//	@Summary		Delete Team Project
//	@Description	Deletes a Team Project via direct team route
//	@Tags			project
//	@Produce		json
//	@Param			teamId		path	string	true	"the team ID"
//	@Param			projectId	path	string	true	"the project ID to delete"
//	@Success		200			object	standardJsonResponse{}
//	@Failure		400			object	standardJsonResponse{}
//	@Failure		500			object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/teams/{teamId}/projects/{projectId} [delete]
func (s *Service) handleTeamProjectDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)

		projectID := r.PathValue("projectId")
		idErr := validate.Var(projectID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		teamID := r.PathValue("teamId")
		teamIDErr := validate.Var(teamID, "required,uuid")
		if teamIDErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, teamIDErr.Error()))
			return
		}

		err := s.ProjectDataSvc.DeleteTeamProject(ctx, teamID, projectID)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleTeamProjectDelete error", zap.Error(err),
				zap.String("project_id", projectID),
				zap.String("team_id", teamID),
				zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, nil, nil)
	}
}
