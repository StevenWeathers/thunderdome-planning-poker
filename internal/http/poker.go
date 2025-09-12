package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strconv"

	"go.uber.org/zap"

	"github.com/StevenWeathers/thunderdome-planning-poker/internal/http/poker"
	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
)

// handleGetUserGames looks up poker games associated with UserID
//
//	@Summary		Get PokerGames
//	@Description	get list of poker games for the user
//	@Tags			poker
//	@Produce		json
//	@Param			userId	path	string	true	"the user ID to get poker games for"
//	@Param			limit	query	int		false	"Max number of results to return"
//	@Param			offset	query	int		false	"Starting point to return rows from, should be multiplied by limit or 0"
//	@Success		200		object	standardJsonResponse{data=[]thunderdome.Poker}
//	@Failure		403		object	standardJsonResponse{}
//	@Failure		404		object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/users/{userId}/battles [get]
func (s *Service) handleGetUserGames() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		limit, offset := getLimitOffsetFromRequest(r)

		userID := r.PathValue("userId")
		idErr := validate.Var(userID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		games, count, err := s.PokerDataSvc.GetGamesByUser(userID, limit, offset)
		if err != nil {
			s.Failure(w, r, http.StatusNotFound, Errorf(ENOTFOUND, "BATTLE_NOT_FOUND"))
			return
		}

		meta := &pagination{
			Count:  count,
			Offset: offset,
			Limit:  limit,
		}

		s.Success(w, r, http.StatusOK, games, meta)
	}
}

type battleRequestBody struct {
	Name                 string               `json:"name" validate:"required"`
	EstimationScaleID    string               `json:"estimationScaleId"`
	PointValuesAllowed   []string             `json:"pointValuesAllowed" validate:"required"`
	AutoFinishVoting     bool                 `json:"autoFinishVoting"`
	Stories              []*thunderdome.Story `json:"plans"`
	PointAverageRounding string               `json:"pointAverageRounding" validate:"required,oneof=ceil round floor"`
	HideVoterIdentity    bool                 `json:"hideVoterIdentity"`
	Facilitators         []string             `json:"battleLeaders"`
	JoinCode             string               `json:"joinCode"`
	FacilitatorCode      string               `json:"leaderCode"`
}

// handlePokerCreate handles creating a poker game
//
//	@Summary		Create Poker Game
//	@Description	Create a poker game associated to the user
//	@Tags			poker
//	@Produce		json
//	@Param			userId			path	string				true	"the user ID"
//	@Param			orgId			path	string				false	"the organization ID"
//	@Param			departmentId	path	string				false	"the department ID"
//	@Param			teamId			path	string				false	"the team ID"
//	@Param			battle			body	battleRequestBody	false	"new poker game object"
//	@Success		200				object	standardJsonResponse{data=thunderdome.Poker}
//	@Failure		403				object	standardJsonResponse{}
//	@Failure		500				object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/users/{userId}/battles [post]
//	@Router			/teams/{teamId}/users/{userId}/battles [post]
//	@Router			/{orgId}/teams/{teamId}/users/{userId}/battles [post]
//	@Router			/{orgId}/departments/{departmentId}/teams/{teamId}/users/{userId}/battles [post]
func (s *Service) handlePokerCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)

		userID := r.PathValue("userId")
		idErr := validate.Var(userID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		teamID := r.PathValue("teamId")
		if teamID == "" && s.Config.RequireTeams {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, "BATTLE_CREATION_REQUIRES_TEAM"))
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
		// if battle created with team association
		if teamID != "" {
			if isTeamUserOrAnAdmin(r) {
				newGame, err = s.PokerDataSvc.TeamCreateGame(ctx, teamID, userID, b.Name, b.EstimationScaleID, b.PointValuesAllowed, b.Stories, b.AutoFinishVoting, b.PointAverageRounding, b.JoinCode, b.FacilitatorCode, b.HideVoterIdentity)
				if err != nil {
					s.Logger.Ctx(ctx).Error("handlePokerCreate error", zap.Error(err),
						zap.String("entity_user_id", userID), zap.String("team_id", teamID),
						zap.String("poker_name", b.Name), zap.String("session_user_id", sessionUserID))
					s.Failure(w, r, http.StatusInternalServerError, err)
					return
				}
			} else {
				s.Failure(w, r, http.StatusForbidden, Errorf(EUNAUTHORIZED, "REQUIRES_TEAM_USER"))
				return
			}
		} else {
			newGame, err = s.PokerDataSvc.CreateGame(ctx, userID, b.Name, b.EstimationScaleID, b.PointValuesAllowed, b.Stories, b.AutoFinishVoting, b.PointAverageRounding, b.JoinCode, b.FacilitatorCode, b.HideVoterIdentity)
			if err != nil {
				s.Logger.Ctx(ctx).Error("handlePokerCreate error", zap.Error(err),
					zap.String("entity_user_id", userID), zap.String("poker_name", b.Name),
					zap.String("session_user_id", sessionUserID))
				s.Failure(w, r, http.StatusInternalServerError, err)
				return
			}
		}

		// when battleLeaders array is passed add additional leaders to battle
		if len(b.Facilitators) > 0 {
			updatedFacilitators, err := s.PokerDataSvc.AddFacilitatorsByEmail(ctx, newGame.ID, b.Facilitators)
			if err != nil {
				s.Logger.Error("error adding additional battle leaders", zap.String("poker_id", newGame.ID),
					zap.String("entity_user_id", userID), zap.Any("poker_facilitators", b.Facilitators),
					zap.String("session_user_id", sessionUserID))
			} else {
				newGame.Facilitators = updatedFacilitators
			}
		}

		s.Success(w, r, http.StatusOK, newGame, nil)
	}
}

// handleGetPokerGames gets a list of poker games
//
//	@Summary		Get Poker Games
//	@Description	get list of poker games
//	@Tags			poker
//	@Produce		json
//	@Param			limit	query	int		false	"Max number of results to return"
//	@Param			offset	query	int		false	"Starting point to return rows from, should be multiplied by limit or 0"
//	@Param			active	query	boolean	false	"Only active poker games"
//	@Success		200		object	standardJsonResponse{data=[]thunderdome.Poker}
//	@Failure		500		object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/battles [get]
func (s *Service) handleGetPokerGames() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)
		limit, offset := getLimitOffsetFromRequest(r)
		query := r.URL.Query()
		var err error
		var count int
		var games []*thunderdome.Poker
		Active, _ := strconv.ParseBool(query.Get("active"))

		if Active {
			games, count, err = s.PokerDataSvc.GetActiveGames(limit, offset)
		} else {
			games, count, err = s.PokerDataSvc.GetGames(limit, offset)
		}

		if err != nil {
			s.Logger.Ctx(ctx).Error("handleGetPokerGames error", zap.Error(err),
				zap.Int("limit", limit), zap.Int("offset", offset), zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		meta := &pagination{
			Count:  count,
			Offset: offset,
			Limit:  limit,
		}

		s.Success(w, r, http.StatusOK, games, meta)
	}
}

// handleGetPokerGame gets the poker game by ID
//
//	@Summary		Get Poker Game
//	@Description	get poker game by ID
//	@Tags			poker
//	@Produce		json
//	@Param			battleId	path	string	true	"the poker game ID to get"
//	@Success		200			object	standardJsonResponse{data=thunderdome.Poker}
//	@Failure		403			object	standardJsonResponse{}
//	@Failure		404			object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/battles/{battleId} [get]
func (s *Service) handleGetPokerGame() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		gameID := r.PathValue("battleId")
		idErr := validate.Var(gameID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		sessionUserID := r.Context().Value(contextKeyUserID).(string)
		userType := r.Context().Value(contextKeyUserType).(string)

		game, err := s.PokerDataSvc.GetGameByID(gameID, sessionUserID)
		if err != nil {
			s.Failure(w, r, http.StatusNotFound, Errorf(ENOTFOUND, "BATTLE_NOT_FOUND"))
			return
		}

		// don't allow retrieving battle details if battle has JoinCode and user hasn't joined yet
		if game.JoinCode != "" {
			userErr := s.PokerDataSvc.GetUserActiveStatus(gameID, sessionUserID)
			if userErr != nil && userErr.Error() != "DUPLICATE_BATTLE_USER" && userType != thunderdome.AdminUserType {
				s.Failure(w, r, http.StatusForbidden, Errorf(EUNAUTHORIZED, "USER_MUST_JOIN_BATTLE"))
				return
			}
		}

		s.Success(w, r, http.StatusOK, game, nil)
	}
}

type planRequestBody struct {
	Name               string `json:"planName"`
	Type               string `json:"type"`
	ReferenceID        string `json:"referenceId"`
	Link               string `json:"link"`
	Description        string `json:"description"`
	AcceptanceCriteria string `json:"acceptanceCriteria"`
	Priority           int32  `json:"priority"`
}

// handlePokerStoryAdd handles adding a story to poker
//
//	@Summary		Create Poker Story
//	@Description	Creates a poker story
//	@Param			battleId	path	string			true	"the poker game ID"
//	@Param			plan		body	planRequestBody	true	"new story object"
//	@Tags			poker
//	@Produce		json
//	@Success		200	object	standardJsonResponse{}
//	@Success		403	object	standardJsonResponse{}
//	@Success		500	object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/battles/{battleId}/plans [post]
func (s *Service) handlePokerStoryAdd(pokerSvc *poker.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		gameID := r.PathValue("battleId")
		idErr := validate.Var(gameID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		sessionUserID := ctx.Value(contextKeyUserID).(string)

		body, bodyErr := io.ReadAll(r.Body)
		if bodyErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
			return
		}

		var story = planRequestBody{}
		jsonErr := json.Unmarshal(body, &story)
		if jsonErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		inputErr := validate.Struct(story)
		if inputErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, inputErr.Error()))
			return
		}

		_, err := pokerSvc.APIEvent(ctx, gameID, sessionUserID, "add_plan", string(body))
		if err != nil {
			s.Logger.Ctx(ctx).Error("handlePokerStoryAdd error", zap.Error(err),
				zap.String("poker_id", gameID), zap.String("session_user_id", sessionUserID),
				zap.String("story_name", story.Name))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, nil, nil)
	}
}

type storyUpdateRequestBody struct {
	ID                 string `json:"planId" swaggerignore:"true"`
	Name               string `json:"planName"`
	Type               string `json:"type"`
	ReferenceID        string `json:"referenceId"`
	Link               string `json:"link"`
	Description        string `json:"description"`
	AcceptanceCriteria string `json:"acceptanceCriteria"`
	Priority           int32  `json:"priority"`
}

// handlePokerStoryUpdate handles updating a poker story
//
//	@Summary		Update Poker Story
//	@Description	Updates a poker story
//	@Param			battleId	path	string					true	"the poker game ID"
//	@Param			planId		path	string					true	"the poker story ID"
//	@Param			story		body	storyUpdateRequestBody	true	"updated story object"
//	@Tags			poker
//	@Produce		json
//	@Success		200	object	standardJsonResponse{}
//	@Success		403	object	standardJsonResponse{}
//	@Success		500	object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/battles/{battleId}/plans/{planId} [put]
func (s *Service) handlePokerStoryUpdate(pokerSvc *poker.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		gameID := r.PathValue("battleId")
		idErr := validate.Var(gameID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(
				EINVALID, fmt.Errorf("invalid battleId: %v", idErr.Error()).Error()),
			)
			return
		}
		storyID := r.PathValue("planId")
		idErr = validate.Var(storyID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(
				EINVALID, fmt.Errorf("invalid storyId: %v", idErr.Error()).Error()),
			)
			return
		}
		sessionUserID := ctx.Value(contextKeyUserID).(string)

		body, bodyErr := io.ReadAll(r.Body)
		if bodyErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
			return
		}

		var story = storyUpdateRequestBody{}
		jsonErr := json.Unmarshal(body, &story)
		if jsonErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		story.ID = storyID
		inputErr := validate.Struct(story)
		if inputErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, inputErr.Error()))
			return
		}

		updatedStory, err := json.Marshal(story)
		if err != nil {
			s.Failure(w, r, http.StatusInternalServerError, Errorf(EINTERNAL, err.Error()))
			return
		}

		_, err = pokerSvc.APIEvent(ctx, gameID, sessionUserID, "revise_plan", string(updatedStory))
		if err != nil {
			s.Logger.Ctx(ctx).Error("handlePokerStoryUpdate error", zap.Error(err),
				zap.String("poker_id", gameID), zap.String("session_user_id", sessionUserID),
				zap.String("story_id", storyID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, nil, nil)
	}
}

// handlePokerStoryDelete handles deleting a story from poker
//
//	@Summary		Delete Poker Story
//	@Description	Deletes a poker story
//	@Param			battleId	path	string	true	"the poker game ID"
//	@Param			planId		path	string	true	"the story ID"
//	@Tags			poker
//	@Produce		json
//	@Success		200	object	standardJsonResponse{}
//	@Success		403	object	standardJsonResponse{}
//	@Success		500	object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/battles/{battleId}/plans/{planId} [delete]
func (s *Service) handlePokerStoryDelete(pokerSvc *poker.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		gameID := r.PathValue("battleId")
		idErr := validate.Var(gameID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		storyID := r.PathValue("planId")
		sidErr := validate.Var(storyID, "required,uuid")
		if sidErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, sidErr.Error()))
			return
		}
		sessionUserID := ctx.Value(contextKeyUserID).(string)

		_, err := pokerSvc.APIEvent(ctx, gameID, sessionUserID, "burn_plan", storyID)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handlePokerStoryDelete error", zap.Error(err),
				zap.String("poker_id", gameID), zap.String("session_user_id", sessionUserID),
				zap.String("story_id", storyID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, nil, nil)
	}
}

// handlePokerDelete handles deleting a poker game
//
//	@Summary		Delete Poker Game
//	@Description	Deletes a poker game
//	@Param			battleId	path	string	true	"the poker game ID"
//	@Tags			poker
//	@Produce		json
//	@Success		200	object	standardJsonResponse{}
//	@Success		403	object	standardJsonResponse{}
//	@Success		500	object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/battles/{battleId} [delete]
func (s *Service) handlePokerDelete(pokerSvc *poker.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		gameID := r.PathValue("battleId")
		idErr := validate.Var(gameID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		sessionUserID := ctx.Value(contextKeyUserID).(string)

		_, err := pokerSvc.APIEvent(ctx, gameID, sessionUserID, "concede_battle", "")
		if err != nil {
			s.Logger.Ctx(ctx).Error("handlePokerDelete error", zap.Error(err),
				zap.String("poker_id", gameID), zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, nil, nil)
	}
}

// handleGameStop handles stopping a poker game via HTTP POST request
//
//	@Summary		Stop Game
//	@Description	Stops a poker game by setting the ended_date
//	@Tags			games
//	@Produce		json
//	@Param			gameId	path		string	true	"the game ID"
//	@Success		200		object		standardJsonResponse{}
//	@Success		403		object		standardJsonResponse{}
//	@Success		500		object		standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/games/{gameId}/stop [post]
func (s *Service) handleGameStop() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		gameID := r.PathValue("gameId")
		idErr := validate.Var(gameID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		sessionUserID := ctx.Value(contextKeyUserID).(string)

		// CRITICAL SECURITY: Verify user is a facilitator of the game
		err := s.PokerDataSvc.ConfirmFacilitator(gameID, sessionUserID)
		if err != nil {
			s.Logger.Ctx(ctx).Warn("SECURITY_VIOLATION: unauthorized game stop attempt via HTTP",
				zap.String("game_id", gameID),
				zap.String("user_id", sessionUserID),
				zap.String("violation_type", "unauthorized_stop_game_http"))
			s.Failure(w, r, http.StatusForbidden, Errorf(EUNAUTHORIZED, "User not authorized to stop this game"))
			return
		}

		// Stop the game
		stopErr := s.PokerDataSvc.StopGame(gameID)
		if stopErr != nil {
			s.Logger.Ctx(ctx).Error("handleGameStop error", zap.Error(stopErr),
				zap.String("game_id", gameID), zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, stopErr)
			return
		}

		// SECURITY AUDIT: Log successful game stop for compliance
		s.Logger.Ctx(ctx).Info("SECURITY_AUDIT: game stopped by authorized facilitator via HTTP",
			zap.String("game_id", gameID),
			zap.String("stopped_by", sessionUserID),
			zap.String("action", "stop_game_http"))

		s.Success(w, r, http.StatusOK, nil, nil)
	}
}
