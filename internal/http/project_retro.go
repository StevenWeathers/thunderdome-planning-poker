package http

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
	"go.uber.org/zap"
)

// handleGetProjectRetros retrieves retros for a specific project
//
//	@Summary		Get Project Retros
//	@Description	Retrieve retros for a specific project
//	@Tags			projects,retros
//	@Produce		json
//	@Param			projectId	path	string	true	"the project ID to get retros for"
//	@Param			limit		query	int		false	"Max number of results to return"
//	@Param			offset		query	int		false	"Starting point to return rows from, should be multiplied by limit or 0"
//	@Success		200			object	standardJsonResponse{data=[]thunderdome.Retro}
//	@Failure		403			object	standardJsonResponse{}
//	@Failure		404			object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/projects/{projectId}/retros [get]
func (s *Service) handleGetProjectRetros() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		projectID := r.PathValue("projectId")
		idErr := validate.Var(projectID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		limit, offset := getLimitOffsetFromRequest(r)

		retros, err := s.ProjectDataSvc.ListRetros(ctx, projectID, limit, offset)
		if err != nil {
			s.Failure(w, r, http.StatusInternalServerError, Errorf(EINTERNAL, err.Error()))
			return
		}

		s.Success(w, r, http.StatusOK, retros, nil)
	}
}

// handleCreateProjectRetro handles creating a retro associated with a project
//
//	@Summary		Create Retro
//	@Description	Create a retro associated with a project
//	@Tags			project,retro
//	@Produce		json
//	@Param			projectId	path	string					true	"the project ID"
//	@Param			retro		body	retroCreateRequestBody	false	"new retro object"
//	@Success		200			object	standardJsonResponse{data=thunderdome.Retro}
//	@Failure		403			object	standardJsonResponse{}
//	@Failure		500			object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/projects/{projectId}/retros [post]
func (s *Service) handleCreateProjectRetro() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)

		projectID := r.PathValue("projectId")
		idErr := validate.Var(projectID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		body, bodyErr := io.ReadAll(r.Body)
		if bodyErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
			return
		}

		var nr = retroCreateRequestBody{}
		jsonErr := json.Unmarshal(body, &nr)
		if jsonErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		inputErr := validate.Struct(nr)
		if inputErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, inputErr.Error()))
			return
		}

		if nr.TemplateID == nil {
			// get default template
			template, err := s.RetroTemplateDataSvc.GetDefaultPublicTemplate(ctx)
			if err != nil {
				s.Logger.Ctx(ctx).Error("handleRetroCreate get default template by id error", zap.Error(err),
					zap.String("session_user_id", sessionUserID))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			nr.TemplateID = &template.ID
		}

		var newRetro *thunderdome.Retro
		var err error

		newRetro, err = s.RetroDataSvc.CreateRetro(ctx, sessionUserID, "", nr.RetroName, nr.JoinCode, nr.FacilitatorCode, nr.MaxVotes, nr.BrainstormVisibility, nr.PhaseTimeLimitMin, nr.PhaseAutoAdvance, nr.AllowCumulativeVoting, nr.HideVotesDuringVoting, *nr.TemplateID)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleRetroCreate error", zap.Error(err),
				zap.String("entity_user_id", sessionUserID),
				zap.String("retro_name", nr.RetroName),
				zap.String("session_user_id", sessionUserID),
				zap.String("project_id", projectID))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = s.ProjectDataSvc.AssociateRetro(ctx, projectID, newRetro.ID)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleRetroCreate associate retro with project error", zap.Error(err),
				zap.String("retro_id", newRetro.ID),
				zap.String("project_id", projectID),
				zap.String("session_user_id", sessionUserID))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		s.Success(w, r, http.StatusOK, newRetro, nil)
	}
}

// handleProjectRetroRemove handles removing retro from a project
//
//	@Summary		Remove Project Retro
//	@Description	Remove a retro from the project
//	@Tags			projects,retro
//	@Produce		json
//	@Param			projectId	path	string	true	"the project ID"
//	@Param			retroId		path	string	true	"the retro ID"
//	@Success		200			object	standardJsonResponse{}
//	@Success		403			object	standardJsonResponse{}
//	@Success		500			object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/projects/{projectId}/retros/{retroId} [delete]
func (s *Service) handleProjectRetroRemove() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)

		projectID := r.PathValue("projectId")
		idErr := validate.Var(projectID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		retroID := r.PathValue("retroId")
		idErr = validate.Var(retroID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		err := s.ProjectDataSvc.RemoveRetro(ctx, projectID, retroID)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleProjectRetroRemove error", zap.Error(err), zap.String("project_id", projectID),
				zap.String("retro_id", retroID), zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, nil, nil)
	}
}
