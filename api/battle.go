package api

import (
	"encoding/json"
	"github.com/StevenWeathers/thunderdome-planning-poker/model"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
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
		Limit, Offset := getLimitOffsetFromRequest(r, w)
		vars := mux.Vars(r)
		UserID := vars["userId"]

		battles, Count, err := a.db.GetBattlesByUser(UserID, Limit, Offset)
		if err != nil {
			Failure(w, r, http.StatusNotFound, Errorf(ENOTFOUND, "BATTLE_NOT_FOUND"))
			return
		}

		Meta := &pagination{
			Count:  Count,
			Offset: Offset,
			Limit:  Limit,
		}

		Success(w, r, http.StatusOK, battles, Meta)
	}
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
// @Param name body string false "the battle name"
// @Param pointValuesAllowed body []string false "the allowed point values e.g. 1,2,3,5,8"
// @Param autoFinishVoting body string false "whether to automatically complete voting when all users have voted"
// @Param plans body []model.Plan false "the battle plans"
// @Param pointAverageRounding body string false "which javascript math method to use for rounding point average"
// @Param battleLeaders body []string true "additional battle leaders beyond the user creating the battle"
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
		vars := mux.Vars(r)
		UserID := vars["userId"]
		UserType := r.Context().Value(contextKeyUserType).(string)

		body, bodyErr := ioutil.ReadAll(r.Body) // check for errors
		if bodyErr != nil {
			Failure(w, r, http.StatusInternalServerError, bodyErr)
			return
		}

		var keyVal struct {
			BattleName           string        `json:"name"`
			PointValuesAllowed   []string      `json:"pointValuesAllowed"`
			AutoFinishVoting     bool          `json:"autoFinishVoting"`
			Plans                []*model.Plan `json:"plans"`
			PointAverageRounding string        `json:"pointAverageRounding"`
			BattleLeaders        []string      `json:"battleLeaders"`
		}
		json.Unmarshal(body, &keyVal) // check for errors

		newBattle, err := a.db.CreateBattle(UserID, keyVal.BattleName, keyVal.PointValuesAllowed, keyVal.Plans, keyVal.AutoFinishVoting, keyVal.PointAverageRounding)
		if err != nil {
			Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		// when battleLeaders array is passed add additional leaders to battle
		if len(keyVal.BattleLeaders) > 0 {
			updatedLeaders, err := a.db.AddBattleLeadersByEmail(newBattle.Id, keyVal.BattleLeaders)
			if err != nil {
				log.Println("error adding additional battle leaders")
			} else {
				newBattle.Leaders = updatedLeaders
			}
		}

		// if battle created with team association
		TeamID, ok := vars["teamId"]
		if ok {
			OrgRole := r.Context().Value(contextKeyOrgRole)
			DepartmentRole := r.Context().Value(contextKeyDepartmentRole)
			TeamRole := r.Context().Value(contextKeyTeamRole).(string)
			var isAdmin bool
			if UserType != adminUserType || (DepartmentRole != nil && DepartmentRole.(string) == "ADMIN") {
				isAdmin = true
			}
			if UserType != adminUserType || (OrgRole != nil && OrgRole.(string) == "ADMIN") {
				isAdmin = true
			}

			if isAdmin == true || TeamRole != "" {
				err := a.db.TeamAddBattle(TeamID, newBattle.Id)

				if err != nil {
					Failure(w, r, http.StatusInternalServerError, err)
					return
				}
			}
		}

		Success(w, r, http.StatusOK, newBattle, nil)
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
		Limit, Offset := getLimitOffsetFromRequest(r, w)
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
			Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		Meta := &pagination{
			Count:  Count,
			Offset: Offset,
			Limit:  Limit,
		}

		Success(w, r, http.StatusOK, Battles, Meta)
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
		UserId := r.Context().Value(contextKeyUserID).(string)
		UserType := r.Context().Value(contextKeyUserType).(string)

		battle, err := a.db.GetBattle(BattleId, UserId, a.config.AESHashkey)
		if err != nil {
			Failure(w, r, http.StatusNotFound, Errorf(ENOTFOUND, "BATTLE_NOT_FOUND"))
			return
		}

		// don't allow retrieving battle details if battle has JoinCode and user hasn't joined yet
		if battle.JoinCode != "" {
			UserErr := a.db.GetBattleUserActiveStatus(BattleId, UserId)
			if UserErr != nil && UserType != adminUserType {
				Failure(w, r, http.StatusForbidden, Errorf(EUNAUTHORIZED, "USER_MUST_JOIN_BATTLE"))
				return
			}
		}

		Success(w, r, http.StatusOK, battle, nil)
	}
}
