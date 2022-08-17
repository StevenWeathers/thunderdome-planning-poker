package api

import (
	"encoding/json"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"strconv"

	"github.com/StevenWeathers/thunderdome-planning-poker/api/battle"

	"github.com/StevenWeathers/thunderdome-planning-poker/model"
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
// @Success 200 object standardJsonResponse{data=[]model.Battle}
// @Failure 403 object standardJsonResponse{}
// @Failure 404 object standardJsonResponse{}
// @Security ApiKeyAuth
// @Router /users/{userId}/battles [get]
func (a *api) handleGetUserBattles() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		Limit, Offset := getLimitOffsetFromRequest(r)
		vars := mux.Vars(r)
		UserID := vars["userId"]

		battles, Count, err := a.db.GetBattlesByUser(UserID, Limit, Offset)
		if err != nil {
			a.Failure(w, r, http.StatusNotFound, Errorf(ENOTFOUND, "BATTLE_NOT_FOUND"))
			return
		}

		Meta := &pagination{
			Count:  Count,
			Offset: Offset,
			Limit:  Limit,
		}

		a.Success(w, r, http.StatusOK, battles, Meta)
	}
}

type battleRequestBody struct {
	BattleName           string        `json:"name" validate:"required"`
	PointValuesAllowed   []string      `json:"pointValuesAllowed" validate:"required"`
	AutoFinishVoting     bool          `json:"autoFinishVoting"`
	Plans                []*model.Plan `json:"plans"`
	PointAverageRounding string        `json:"pointAverageRounding" validate:"required,oneof=ceil round floor"`
	BattleLeaders        []string      `json:"battleLeaders"`
	JoinCode             string        `json:"joinCode"`
	LeaderCode           string        `json:"leaderCode"`
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
// @Success 200 object standardJsonResponse{data=model.Battle}
// @Failure 403 object standardJsonResponse{}
// @Failure 500 object standardJsonResponse{}
// @Security ApiKeyAuth
// @Router /users/{userId}/battles [post]
// @Router /teams/{teamId}/users/{userId}/battles [post]
// @Router /{orgId}/teams/{teamId}/users/{userId}/battles [post]
// @Router /{orgId}/departments/{departmentId}/teams/{teamId}/users/{userId}/battles [post]
func (a *api) handleBattleCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		vars := mux.Vars(r)
		UserID := vars["userId"]
		TeamID, teamIdExists := vars["teamId"]

		if !teamIdExists && viper.GetBool("config.require_teams") {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, "BATTLE_CREATION_REQUIRES_TEAM"))
			return
		}

		body, bodyErr := io.ReadAll(r.Body)
		if bodyErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
			return
		}

		var b = battleRequestBody{}
		jsonErr := json.Unmarshal(body, &b)
		if jsonErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		inputErr := validate.Struct(b)
		if inputErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, inputErr.Error()))
			return
		}

		var newBattle *model.Battle
		var err error
		// if battle created with team association
		if teamIdExists {
			if isTeamUserOrAnAdmin(r) {
				newBattle, err = a.db.TeamCreateBattle(ctx, TeamID, UserID, b.BattleName, b.PointValuesAllowed, b.Plans, b.AutoFinishVoting, b.PointAverageRounding, b.JoinCode, b.LeaderCode)
				if err != nil {
					a.Failure(w, r, http.StatusInternalServerError, err)
					return
				}
			} else {
				a.Failure(w, r, http.StatusForbidden, Errorf(EUNAUTHORIZED, "REQUIRES_TEAM_USER"))
				return
			}
		} else {
			newBattle, err = a.db.CreateBattle(ctx, UserID, b.BattleName, b.PointValuesAllowed, b.Plans, b.AutoFinishVoting, b.PointAverageRounding, b.JoinCode, b.LeaderCode)
			if err != nil {
				a.Failure(w, r, http.StatusInternalServerError, err)
				return
			}
		}

		// when battleLeaders array is passed add additional leaders to battle
		if len(b.BattleLeaders) > 0 {
			updatedLeaders, err := a.db.AddBattleLeadersByEmail(ctx, newBattle.Id, b.BattleLeaders)
			if err != nil {
				a.logger.Error("error adding additional battle leaders")
			} else {
				newBattle.Leaders = updatedLeaders
			}
		}

		a.Success(w, r, http.StatusOK, newBattle, nil)
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
// @Success 200 object standardJsonResponse{data=[]model.Battle}
// @Failure 500 object standardJsonResponse{}
// @Security ApiKeyAuth
// @Router /battles [get]
func (a *api) handleGetBattles() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		Limit, Offset := getLimitOffsetFromRequest(r)
		query := r.URL.Query()
		var err error
		var Count int
		var Battles []*model.Battle
		Active, _ := strconv.ParseBool(query.Get("active"))

		if Active {
			Battles, Count, err = a.db.GetActiveBattles(Limit, Offset)
		} else {
			Battles, Count, err = a.db.GetBattles(Limit, Offset)
		}

		if err != nil {
			a.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		Meta := &pagination{
			Count:  Count,
			Offset: Offset,
			Limit:  Limit,
		}

		a.Success(w, r, http.StatusOK, Battles, Meta)
	}
}

// handleGetBattle gets the battle by ID
// @Summary Get Battle
// @Description get battle by ID
// @Tags battle
// @Produce  json
// @Param battleId path string true "the battle ID to get"
// @Success 200 object standardJsonResponse{data=model.Battle}
// @Failure 403 object standardJsonResponse{}
// @Failure 404 object standardJsonResponse{}
// @Security ApiKeyAuth
// @Router /battles/{battleId} [get]
func (a *api) handleGetBattle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		BattleId := vars["battleId"]
		idErr := validate.Var(BattleId, "required,uuid")
		if idErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		UserId := r.Context().Value(contextKeyUserID).(string)
		UserType := r.Context().Value(contextKeyUserType).(string)

		b, err := a.db.GetBattle(BattleId, UserId)
		if err != nil {
			a.Failure(w, r, http.StatusNotFound, Errorf(ENOTFOUND, "BATTLE_NOT_FOUND"))
			return
		}

		// don't allow retrieving battle details if battle has JoinCode and user hasn't joined yet
		if b.JoinCode != "" {
			UserErr := a.db.GetBattleUserActiveStatus(BattleId, UserId)
			if UserErr != nil && UserType != adminUserType {
				a.Failure(w, r, http.StatusForbidden, Errorf(EUNAUTHORIZED, "USER_MUST_JOIN_BATTLE"))
				return
			}
		}

		a.Success(w, r, http.StatusOK, b, nil)
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
func (a *api) handleBattlePlanAdd(b *battle.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		BattleID := vars["battleId"]
		idErr := validate.Var(BattleID, "required,uuid")
		if idErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		UserID := r.Context().Value(contextKeyUserID).(string)

		body, bodyErr := io.ReadAll(r.Body)
		if bodyErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
			return
		}

		err := b.APIEvent(r.Context(), BattleID, UserID, "add_plan", string(body))
		if err != nil {
			a.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		a.Success(w, r, http.StatusOK, nil, nil)
	}
}
