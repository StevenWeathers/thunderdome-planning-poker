package http

import (
	"encoding/json"
	"io"
	"net/http"

	"go.uber.org/zap"
)

// handleGetProjectStoryboards retrieves storyboards for a specific project
//
//	@Summary		Get Project Storyboards
//	@Description	Retrieve storyboards for a specific project
//	@Tags			storyboards
//	@Produce		json
//	@Param			projectId	path	string	true	"the project ID to get storyboards for"
//	@Param			limit		query	int		false	"Max number of results to return"
//	@Param			offset		query	int		false	"Starting point to return rows from, should be multiplied by limit or 0"
//	@Success		200			object	standardJsonResponse{data=[]thunderdome.Storyboard}
//	@Failure		403			object	standardJsonResponse{}
//	@Failure		404			object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/projects/{projectId}/storyboards [get]
func (s *Service) handleGetProjectStoryboards() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		projectID := r.PathValue("projectId")
		idErr := validate.Var(projectID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		limit, offset := getLimitOffsetFromRequest(r)

		storyboards, err := s.ProjectDataSvc.ListStoryboards(ctx, projectID, limit, offset)
		if err != nil {
			s.Failure(w, r, http.StatusInternalServerError, Errorf(EINTERNAL, err.Error()))
			return
		}

		s.Success(w, r, http.StatusOK, storyboards, nil)
	}
}

// handleCreateProjectStoryboard handles creating a storyboard associated to a project
//
//	@Summary		Create Storyboard
//	@Description	Create a storyboard associated to the project
//	@Tags			storyboard
//	@Produce		json
//	@Param			projectId	path	string						true	"the project ID"
//	@Param			storyboard	body	storyboardCreateRequestBody	false	"new storyboard object"
//	@Success		200			object	standardJsonResponse{data=thunderdome.Storyboard}
//	@Failure		403			object	standardJsonResponse{}
//	@Failure		500			object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/projects/{projectId}/storyboards [post]
func (s *Service) handleCreateProjectStoryboard() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)

		projectID := r.PathValue("projectId")
		idErr := validate.Var(projectID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		body, bodyErr := io.ReadAll(r.Body) // check for errors
		if bodyErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
			return
		}

		var sb = storyboardCreateRequestBody{}
		jsonErr := json.Unmarshal(body, &sb)
		if jsonErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		inputErr := validate.Struct(sb)
		if inputErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, inputErr.Error()))
			return
		}

		newStoryboard, err := s.StoryboardDataSvc.CreateStoryboard(ctx, sessionUserID, sb.StoryboardName, sb.JoinCode, sb.FacilitatorCode)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleStoryboardCreate error", zap.Error(err),
				zap.String("entity_user_id", sessionUserID), zap.String("session_user_id", sessionUserID),
				zap.String("storyboard_name", sb.StoryboardName))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		err = s.ProjectDataSvc.AssociateStoryboard(ctx, projectID, newStoryboard.ID)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleStoryboardCreate associate storyboard with project error", zap.Error(err),
				zap.String("storyboard_id", newStoryboard.ID),
				zap.String("project_id", projectID),
				zap.String("session_user_id", sessionUserID))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		s.Success(w, r, http.StatusOK, newStoryboard, nil)
	}
}

// handleProjectStoryboardRemove handles removing storyboard from a project
//
//	@Summary		Remove Project Storyboard
//	@Description	Remove a storyboard from the project
//	@Tags			projects,storyboard
//	@Produce		json
//	@Param			projectId		path	string	true	"the project ID"
//	@Param			storyboardId	path	string	true	"the storyboard ID"
//	@Success		200				object	standardJsonResponse{}
//	@Success		403				object	standardJsonResponse{}
//	@Success		500				object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/projects/{projectId}/storyboards/{storyboardId} [delete]
func (s *Service) handleProjectStoryboardRemove() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)

		projectID := r.PathValue("projectId")
		idErr := validate.Var(projectID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		storyboardID := r.PathValue("storyboardId")
		idErr = validate.Var(storyboardID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		err := s.ProjectDataSvc.RemoveStoryboard(ctx, projectID, storyboardID)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleProjectStoryboardRemove error", zap.Error(err), zap.String("project_id", projectID),
				zap.String("storyboard_id", storyboardID), zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, nil, nil)
	}
}
