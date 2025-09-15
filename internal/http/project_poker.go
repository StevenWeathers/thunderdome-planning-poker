package http

import (
	"encoding/json"
	"io"
	"net/http"
	"slices"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
	"go.uber.org/zap"
)

// handleGetProjectPokerGames retrieves poker games for a specific project
//
//	@Summary		Get Project Pokers
//	@Description	Retrieve poker games for a specific project
//	@Tags			projects,poker
//	@Produce		json
//	@Param			projectId	path	string	true	"the project ID to get pokergames for"
//	@Param			limit		query	int		false	"Max number of results to return"
//	@Param			offset		query	int		false	"Starting point to return rows from, should be multiplied by limit or 0"
//	@Success		200			object	standardJsonResponse{data=[]thunderdome.Poker}
//	@Failure		403			object	standardJsonResponse{}
//	@Failure		404			object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/projects/{projectId}/poker [get]
func (s *Service) handleGetProjectPokerGames() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		projectID := r.PathValue("projectId")
		idErr := validate.Var(projectID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		limit, offset := getLimitOffsetFromRequest(r)

		games, err := s.ProjectDataSvc.ListPokerGames(ctx, projectID, limit, offset)
		if err != nil {
			s.Failure(w, r, http.StatusInternalServerError, Errorf(EINTERNAL, err.Error()))
			return
		}

		s.Success(w, r, http.StatusOK, games, nil)
	}
}

// handleCreateProjectPokerGame creates a new poker game associated with a specific project
//
//	@Summary		Create Project Poker
//	@Description	Create a new poker game associated with a specific project
//	@Tags			projects,poker
//	@Accept			json
//	@Produce		json
//	@Param			projectId	path	string				true	"the project ID to associate the poker game with"
//	@Param			body		body	battleRequestBody	true	"The poker game request body"
//	@Success		200			object	standardJsonResponse{data=thunderdome.Poker}
//	@Failure		400			object	standardJsonResponse{}
//	@Failure		403			object	standardJsonResponse{}
//	@Failure		500			object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/projects/{projectId}/poker [post]
func (s *Service) handleCreateProjectPokerGame() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)

		projectID := r.PathValue("projectId")
		if err := validate.Var(projectID, "required,uuid"); err != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, err.Error()))
			return
		}

		body, bodyErr := io.ReadAll(r.Body)
		if bodyErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
			return
		}

		var b = battleRequestBody{}
		jsonErr := json.Unmarshal(body, &b)
		if jsonErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		inputErr := validate.Struct(b)
		if inputErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, inputErr.Error()))
			return
		}

		// set a default for backwards compatibility
		scale := &thunderdome.EstimationScale{}
		var scaleErr error
		if b.EstimationScaleID == "" {
			scale, scaleErr = s.PokerDataSvc.GetDefaultPublicEstimationScale(ctx)
			if scaleErr != nil {
				s.Logger.Error("create poker error", zap.Error(scaleErr))
				s.Failure(w, r, http.StatusInternalServerError, scaleErr)
				return
			}
			b.EstimationScaleID = scale.ID
		} else {
			scale, scaleErr = s.PokerDataSvc.GetEstimationScale(ctx, b.EstimationScaleID)
			if scaleErr != nil {
				s.Logger.Error("create poker error", zap.Error(scaleErr))
				s.Failure(w, r, http.StatusInternalServerError, scaleErr)
				return
			}
		}

		// verify that the point values allowed are in the estimation scale
		for _, point := range b.PointValuesAllowed {
			if !slices.Contains(scale.Values, point) {
				s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, "POINT_VALUES_NOT_IN_SCALE"))
				return
			}
		}

		var newGame *thunderdome.Poker
		var err error
		newGame, err = s.PokerDataSvc.CreateGame(ctx, sessionUserID, b.Name, b.EstimationScaleID, b.PointValuesAllowed, b.Stories, b.AutoFinishVoting, b.PointAverageRounding, b.JoinCode, b.FacilitatorCode, b.HideVoterIdentity)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handlePokerCreate error", zap.Error(err),
				zap.String("entity_user_id", sessionUserID), zap.String("poker_name", b.Name),
				zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		err = s.ProjectDataSvc.AssociatePoker(ctx, projectID, newGame.ID)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handlePokerCreate associate poker with project error", zap.Error(err),
				zap.String("poker_id", newGame.ID),
				zap.String("project_id", projectID),
				zap.String("session_user_id", sessionUserID))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// when facilitators array is passed add additional facilitators to poker
		if len(b.Facilitators) > 0 {
			updatedFacilitators, err := s.PokerDataSvc.AddFacilitatorsByEmail(ctx, newGame.ID, b.Facilitators)
			if err != nil {
				s.Logger.Error("error adding additional poker facilitators", zap.String("poker_id", newGame.ID),
					zap.String("entity_user_id", sessionUserID), zap.Any("poker_facilitators", b.Facilitators),
					zap.String("session_user_id", sessionUserID))
			} else {
				newGame.Facilitators = updatedFacilitators
			}
		}

		s.Success(w, r, http.StatusOK, newGame, nil)
	}
}

// handleProjectPokerGameRemove handles removing poker game from a project
//
//	@Summary		Remove Project Poker
//	@Description	Remove a poker game from the project
//	@Tags			projects,poker
//	@Produce		json
//	@Param			projectId	path	string	true	"the project ID"
//	@Param			gameId		path	string	true	"the game ID"
//	@Success		200			object	standardJsonResponse{}
//	@Success		403			object	standardJsonResponse{}
//	@Success		500			object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/projects/{projectId}/poker/{gameId} [delete]
func (s *Service) handleProjectPokerGameRemove() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)

		projectID := r.PathValue("projectId")
		idErr := validate.Var(projectID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		gameID := r.PathValue("gameId")
		idErr = validate.Var(gameID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		err := s.ProjectDataSvc.RemovePokerGame(ctx, projectID, gameID)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleProjectPokerGameRemove error", zap.Error(err), zap.String("project_id", projectID),
				zap.String("game_id", gameID), zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, nil, nil)
	}
}
