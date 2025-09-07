package http

import (
	"net/http"

	"go.uber.org/zap"
)

// handleGetProjectByID gets a specific project by ID
//
//	@Summary		Get Project by ID
//	@Description	get a specific project by its ID
//	@Tags			project
//	@Produce		json
//	@Param			projectId	path	string	true	"the project ID"
//	@Success		200			object	standardJsonResponse{data=thunderdome.Project}
//	@Failure		404			object	standardJsonResponse{}
//	@Failure		500			object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/projects/{projectId} [get]
func (s *Service) handleGetProjectByID() http.HandlerFunc {
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
