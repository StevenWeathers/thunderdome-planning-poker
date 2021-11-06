package api

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/StevenWeathers/thunderdome-planning-poker/model"
	"github.com/gorilla/mux"
)

// handleBattlesGet looks up battles associated with UserID
// @Summary Get Battles
// @Description get list of battles for the user
// @Tags battle
// @Produce  json
// @Param userId path int false "the user ID to get battles for"
// @Success 200 object standardJsonResponse{data=[]model.Battle}
// @Failure 403 object standardJsonResponse{}
// @Failure 404 object standardJsonResponse{}
// @Router /users/{userId}/battles [get]
func (a *api) handleBattlesGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		UserID := vars["userId"]

		battles, err := a.db.GetBattlesByUser(UserID)
		if err != nil {
			Failure(w, r, http.StatusNotFound, Errorf(ENOTFOUND, "BATTLE_NOT_FOUND"))
			return
		}

		Success(w, r, http.StatusOK, battles, nil)
	}
}

// handleBattleCreate handles creating a battle (arena)
// @Summary Create Battle
// @Description Create a battle associated to the user
// @Tags battle
// @Produce  json
// @Param userId path int false "the user ID"
// @Param orgId path int false "the organization ID"
// @Param departmentId path int false "the department ID"
// @Param teamId path int false "the team ID"
// @Param name body string false "the battle name"
// @Param pointValuesAllowed body []string false "the allowed point values e.g. 1,2,3,5,8"
// @Param autoFinishVoting body string false "whether or not to automatically complete voting when all users have voted"
// @Param plans body []model.Plan false "the battle plans"
// @Param pointAverageRounding body string false "which javascript math method to use for rounding point average"
// @Param battleLeaders body []string true "additional battle leaders beyond the user creating the battle"
// @Success 200 object standardJsonResponse{data=model.Battle}
// @Failure 403 object standardJsonResponse{}
// @Failure 500 object standardJsonResponse{}
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
			updatedLeaders, err := a.db.AddBattleLeadersByEmail(newBattle.BattleID, UserID, keyVal.BattleLeaders)
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
				err := a.db.TeamAddBattle(TeamID, newBattle.BattleID)

				if err != nil {
					Failure(w, r, http.StatusInternalServerError, err)
					return
				}
			}
		}

		Success(w, r, http.StatusOK, newBattle, nil)
	}
}
