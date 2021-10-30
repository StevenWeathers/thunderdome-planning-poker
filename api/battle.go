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
// @Description get list of battles for authenticated user
// @Tags battle
// @Produce  json
// @Param id path int false "the user ID to get battles for"
// @Success 200 object standardJsonResponse{data=[]model.Battle}
// @Failure 403 object standardJsonResponse{}
// @Failure 404 object standardJsonResponse{}
// @Router /users/{id}/battles [get]
func (a *api) handleBattlesGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		UserID := vars["id"]
		AuthedUserID := r.Context().Value(contextKeyUserID).(string)

		if UserID != AuthedUserID {
			a.respondWithStandardJSON(w, http.StatusForbidden, false, nil, nil, nil)
			return
		}

		battles, err := a.db.GetBattlesByUser(UserID)

		if err != nil {
			a.respondWithStandardJSON(w, http.StatusNotFound, false, nil, nil, nil)
			return
		}

		a.respondWithStandardJSON(w, http.StatusOK, true, nil, battles, nil)
	}
}

/*
	Battle Handlers
*/
// handleBattleCreate handles creating a battle (arena)
// @Summary Create Battle
// @Description Create a battle associated to authenticated user
// @Tags battle
// @Produce  json
// @Param id path int false "the user ID"
// @Success 200 object standardJsonResponse{data=model.Battle}
// @Failure 403 object standardJsonResponse{}
// @Failure 500 object standardJsonResponse{}
// @Router /users/{id}/battles [post]
func (a *api) handleBattleCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		UserID := vars["id"]
		AuthedUserID := r.Context().Value(contextKeyUserID).(string)

		if UserID != AuthedUserID {
			a.respondWithStandardJSON(w, http.StatusForbidden, false, nil, nil, nil)
			return
		}

		body, bodyErr := ioutil.ReadAll(r.Body) // check for errors
		if bodyErr != nil {
			log.Println("error in reading request body: " + bodyErr.Error() + "\n")
			errors := make([]string, 0)
			errors = append(errors, bodyErr.Error())
			a.respondWithStandardJSON(w, http.StatusInternalServerError, false, errors, nil, nil)
			return
		}

		var keyVal struct {
			BattleName           string        `json:"battleName"`
			PointValuesAllowed   []string      `json:"pointValuesAllowed"`
			AutoFinishVoting     bool          `json:"autoFinishVoting"`
			Plans                []*model.Plan `json:"plans"`
			PointAverageRounding string        `json:"pointAverageRounding"`
			BattleLeaders        []string      `json:"battleLeaders"`
		}
		json.Unmarshal(body, &keyVal) // check for errors

		newBattle, err := a.db.CreateBattle(UserID, keyVal.BattleName, keyVal.PointValuesAllowed, keyVal.Plans, keyVal.AutoFinishVoting, keyVal.PointAverageRounding)
		if err != nil {
			errors := make([]string, 0)
			errors = append(errors, err.Error())
			a.respondWithStandardJSON(w, http.StatusInternalServerError, false, errors, nil, nil)
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
			if DepartmentRole != nil && DepartmentRole.(string) == "ADMIN" {
				isAdmin = true
			}
			if OrgRole != nil && OrgRole.(string) == "ADMIN" {
				isAdmin = true
			}

			if isAdmin == true || TeamRole != "" {
				err := a.db.TeamAddBattle(TeamID, newBattle.BattleID)

				if err != nil {
					errors := make([]string, 0)
					errors = append(errors, err.Error())
					a.respondWithStandardJSON(w, http.StatusInternalServerError, false, errors, nil, nil)
					return
				}
			}
		}

		a.respondWithStandardJSON(w, http.StatusOK, true, nil, newBattle, nil)
	}
}
