package http

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
	"github.com/spf13/viper"

	"github.com/StevenWeathers/thunderdome-planning-poker/http/battle"

	"github.com/gorilla/mux"
)

// handleGetUserBattles looks up battles associated with UserID
// @Summary Get Battles
// @Description get list of battles for the user
// @Tags battle
// @Produce  json
// @Param userId path string true "the user ID to get battles for"
// @Param limit query int false "Max number of results to return"
// @Param offset query int false "Starting point to return rows from, should be multiplied by limit or 0"
// @Success 200 object standardJsonResponse{data=[]thunderdome.Battle}
// @Failure 403 object standardJsonResponse{}
// @Failure 404 object standardJsonResponse{}
// @Security ApiKeyAuth
// @Router /users/{userId}/battles [get]
func (s *Service) handleGetUserBattles() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		Limit, Offset := getLimitOffsetFromRequest(r)
		vars := mux.Vars(r)
		UserID := vars["userId"]

		battles, Count, err := s.BattleService.GetBattlesByUser(UserID, Limit, Offset)
		if err != nil {
			s.Failure(w, r, http.StatusNotFound, Errorf(ENOTFOUND, "BATTLE_NOT_FOUND"))
			return
		}

		Meta := &pagination{
			Count:  Count,
			Offset: Offset,
			Limit:  Limit,
		}

		s.Success(w, r, http.StatusOK, battles, Meta)
	}
}

type battleRequestBody struct {
	BattleName           string              `json:"name" validate:"required"`
	PointValuesAllowed   []string            `json:"pointValuesAllowed" validate:"required"`
	AutoFinishVoting     bool                `json:"autoFinishVoting"`
	Plans                []*thunderdome.Plan `json:"plans"`
	PointAverageRounding string              `json:"pointAverageRounding" validate:"required,oneof=ceil round floor"`
	HideVoterIdentity    bool                `json:"hideVoterIdentity"`
	BattleLeaders        []string            `json:"battleLeaders"`
	JoinCode             string              `json:"joinCode"`
	LeaderCode           string              `json:"leaderCode"`
}

// handleBattleCreate handles creating a battle (arena)
// @Summary Create Battle
// @Description Create a battle associated to the user
// @Tags battle
// @Produce  json
// @Param userId path string true "the user ID"
// @Param orgId path string false "the organization ID"
// @Param departmentId path string false "the department ID"
// @Param teamId path string false "the team ID"
// @Param battle body battleRequestBody false "new battle object"
// @Success 200 object standardJsonResponse{data=thunderdome.Battle}
// @Failure 403 object standardJsonResponse{}
// @Failure 500 object standardJsonResponse{}
// @Security ApiKeyAuth
// @Router /users/{userId}/battles [post]
// @Router /teams/{teamId}/users/{userId}/battles [post]
// @Router /{orgId}/teams/{teamId}/users/{userId}/battles [post]
// @Router /{orgId}/departments/{departmentId}/teams/{teamId}/users/{userId}/battles [post]
func (s *Service) handleBattleCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		vars := mux.Vars(r)
		UserID := vars["userId"]
		TeamID, teamIdExists := vars["teamId"]

		if !teamIdExists && viper.GetBool("config.require_teams") {
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

		var newBattle *thunderdome.Battle
		var err error
		// if battle created with team association
		if teamIdExists {
			if isTeamUserOrAnAdmin(r) {
				newBattle, err = s.BattleService.TeamCreateBattle(ctx, TeamID, UserID, b.BattleName, b.PointValuesAllowed, b.Plans, b.AutoFinishVoting, b.PointAverageRounding, b.JoinCode, b.LeaderCode, b.HideVoterIdentity)
				if err != nil {
					s.Failure(w, r, http.StatusInternalServerError, err)
					return
				}
			} else {
				s.Failure(w, r, http.StatusForbidden, Errorf(EUNAUTHORIZED, "REQUIRES_TEAM_USER"))
				return
			}
		} else {
			newBattle, err = s.BattleService.CreateBattle(ctx, UserID, b.BattleName, b.PointValuesAllowed, b.Plans, b.AutoFinishVoting, b.PointAverageRounding, b.JoinCode, b.LeaderCode, b.HideVoterIdentity)
			if err != nil {
				s.Failure(w, r, http.StatusInternalServerError, err)
				return
			}
		}

		// when battleLeaders array is passed add additional leaders to battle
		if len(b.BattleLeaders) > 0 {
			updatedLeaders, err := s.BattleService.AddBattleLeadersByEmail(ctx, newBattle.Id, b.BattleLeaders)
			if err != nil {
				s.Logger.Error("error adding additional battle leaders")
			} else {
				newBattle.Leaders = updatedLeaders
			}
		}

		s.Success(w, r, http.StatusOK, newBattle, nil)
	}
}

// handleGetBattles gets a list of battles
// @Summary Get Battles
// @Description get list of battles
// @Tags battle
// @Produce  json
// @Param limit query int false "Max number of results to return"
// @Param offset query int false "Starting point to return rows from, should be multiplied by limit or 0"
// @Param active query boolean false "Only active battles"
// @Success 200 object standardJsonResponse{data=[]thunderdome.Battle}
// @Failure 500 object standardJsonResponse{}
// @Security ApiKeyAuth
// @Router /battles [get]
func (s *Service) handleGetBattles() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		Limit, Offset := getLimitOffsetFromRequest(r)
		query := r.URL.Query()
		var err error
		var Count int
		var Battles []*thunderdome.Battle
		Active, _ := strconv.ParseBool(query.Get("active"))

		if Active {
			Battles, Count, err = s.BattleService.GetActiveBattles(Limit, Offset)
		} else {
			Battles, Count, err = s.BattleService.GetBattles(Limit, Offset)
		}

		if err != nil {
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		Meta := &pagination{
			Count:  Count,
			Offset: Offset,
			Limit:  Limit,
		}

		s.Success(w, r, http.StatusOK, Battles, Meta)
	}
}

// handleGetBattle gets the battle by ID
// @Summary Get Battle
// @Description get battle by ID
// @Tags battle
// @Produce  json
// @Param battleId path string true "the battle ID to get"
// @Success 200 object standardJsonResponse{data=thunderdome.Battle}
// @Failure 403 object standardJsonResponse{}
// @Failure 404 object standardJsonResponse{}
// @Security ApiKeyAuth
// @Router /battles/{battleId} [get]
func (s *Service) handleGetBattle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		BattleId := vars["battleId"]
		idErr := validate.Var(BattleId, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		UserId := r.Context().Value(contextKeyUserID).(string)
		UserType := r.Context().Value(contextKeyUserType).(string)

		b, err := s.BattleService.GetBattle(BattleId, UserId)
		if err != nil {
			s.Failure(w, r, http.StatusNotFound, Errorf(ENOTFOUND, "BATTLE_NOT_FOUND"))
			return
		}

		// don't allow retrieving battle details if battle has JoinCode and user hasn't joined yet
		if b.JoinCode != "" {
			UserErr := s.BattleService.GetBattleUserActiveStatus(BattleId, UserId)
			if UserErr != nil && UserType != adminUserType {
				s.Failure(w, r, http.StatusForbidden, Errorf(EUNAUTHORIZED, "USER_MUST_JOIN_BATTLE"))
				return
			}
		}

		s.Success(w, r, http.StatusOK, b, nil)
	}
}

type planRequestBody struct {
	Name               string `json:"planName"`
	Type               string `json:"type"`
	ReferenceID        string `json:"referenceId"`
	Link               string `json:"link"`
	Description        string `json:"description"`
	AcceptanceCriteria string `json:"acceptanceCriteria"`
}

// handleBattlePlanAdd handles adding a plan to battle
// @Summary Create Battle Plan
// @Description Creates a battle plan
// @Param battleId path string true "the battle ID"
// @Param plan body planRequestBody true "new plan object"
// @Tags battle
// @Produce  json
// @Success 200 object standardJsonResponse{}
// @Success 403 object standardJsonResponse{}
// @Success 500 object standardJsonResponse{}
// @Security ApiKeyAuth
// @Router /battles/{battleId}/plans [post]
func (s *Service) handleBattlePlanAdd(b *battle.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		BattleID := vars["battleId"]
		idErr := validate.Var(BattleID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		UserID := r.Context().Value(contextKeyUserID).(string)

		body, bodyErr := io.ReadAll(r.Body)
		if bodyErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
			return
		}

		var plan = planRequestBody{}
		jsonErr := json.Unmarshal(body, &plan)
		if jsonErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		inputErr := validate.Struct(plan)
		if inputErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, inputErr.Error()))
			return
		}

		err := b.APIEvent(r.Context(), BattleID, UserID, "add_plan", string(body))
		if err != nil {
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, nil, nil)
	}
}

// handleBattleDelete handles deleting a battle
// @Summary Delete Battle
// @Description Deletes a battle
// @Param battleId path string true "the battle ID"
// @Tags battle
// @Produce  json
// @Success 200 object standardJsonResponse{}
// @Success 403 object standardJsonResponse{}
// @Success 500 object standardJsonResponse{}
// @Security ApiKeyAuth
// @Router /battles/{battleId} [delete]
func (s *Service) handleBattleDelete(b *battle.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		BattleID := vars["battleId"]
		idErr := validate.Var(BattleID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		UserID := r.Context().Value(contextKeyUserID).(string)

		err := b.APIEvent(r.Context(), BattleID, UserID, "concede_battle", "")
		if err != nil {
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, nil, nil)
	}
}
