package api

import (
	"net/http"
	"strings"

	"github.com/StevenWeathers/thunderdome-planning-poker/model"
	"github.com/gorilla/mux"
)

type teamResponse struct {
	Team     *model.Team `json:"team"`
	TeamRole string      `json:"teamRole"`
}

// handleGetTeamByUser gets an team with user role
// @Summary Get Team
// @Description Get a team with user role
// @Tags team
// @Produce  json
// @Param teamId path string true "the team ID"
// @Success 200 object standardJsonResponse{data=teamResponse}
// @Success 500 object standardJsonResponse{}
// @Router /teams/{teamId} [get]
func (a *api) handleGetTeamByUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		TeamRole := r.Context().Value(contextKeyTeamRole).(string)
		TeamID := vars["teamId"]

		Team, err := a.db.TeamGet(TeamID)
		if err != nil {
			Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		result := &teamResponse{
			Team:     Team,
			TeamRole: TeamRole,
		}

		Success(w, r, http.StatusOK, result, nil)
	}
}

// handleGetTeamsByUser gets a list of teams the user is a part of
// @Summary Get User Teams
// @Description Get a list of teams the user is a part of
// @Tags team
// @Produce  json
// @Param userId path string true "the user ID"
// @Param teamId path string true "the team ID"
// @Success 200 object standardJsonResponse{data=[]model.Team}
// @Success 403 object standardJsonResponse{}
// @Router /users/{userId}/teams/{teamId} [get]
func (a *api) handleGetTeamsByUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		UserID := vars["userId"]

		Limit, Offset := getLimitOffsetFromRequest(r, w)

		Teams := a.db.TeamListByUser(UserID, Limit, Offset)

		Success(w, r, http.StatusOK, Teams, nil)
	}
}

// handleGetTeamUsers gets a list of users associated to the team
// @Summary Get Team users
// @Description Get a list of users associated to the team
// @Tags team
// @Produce  json
// @Param teamId path string true "the team ID"
// @Success 200 object standardJsonResponse{data=[]model.User}
// @Router /teams/{teamId}/users [get]
func (a *api) handleGetTeamUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		TeamID := vars["teamId"]
		Limit, Offset := getLimitOffsetFromRequest(r, w)

		Teams := a.db.TeamUserList(TeamID, Limit, Offset)

		Success(w, r, http.StatusOK, Teams, nil)
	}
}

// handleCreateTeam handles creating an team with current user as admin
// @Summary Create Team
// @Description Creates a team with the current user as the team admin
// @Tags team
// @Produce  json
// @Param userId path string true "the user ID"
// @Param name body string true "the team name"
// @Success 200 object standardJsonResponse{data=model.Team}
// @Success 403 object standardJsonResponse{}
// @Success 500 object standardJsonResponse{}
// @Router /users/{userId}/teams [post]
func (a *api) handleCreateTeam() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		UserID := vars["userId"]

		keyVal := getJSONRequestBody(r, w)

		TeamName := keyVal["name"].(string)
		NewTeam, err := a.db.TeamCreate(UserID, TeamName)
		if err != nil {
			Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		Success(w, r, http.StatusOK, NewTeam, nil)
	}
}

// handleTeamAddUser handles adding user to a team
// @Summary Add Team User
// @Description Adds a user to the team
// @Tags team
// @Produce  json
// @Param teamId path string true "the team ID"
// @Param email body string true "the users email"
// @Param role body string true "the users team role" Enums(MEMBER, ADMIN)
// @Success 200 object standardJsonResponse{}
// @Success 403 object standardJsonResponse{}
// @Success 500 object standardJsonResponse{}
// @Router /{teamId}/users [post]
func (a *api) handleTeamAddUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		keyVal := getJSONRequestBody(r, w)

		vars := mux.Vars(r)
		TeamID := vars["teamId"]
		UserEmail := strings.ToLower(keyVal["email"].(string))
		Role := keyVal["role"].(string)

		User, UserErr := a.db.GetUserByEmail(UserEmail)
		if UserErr != nil {
			Failure(w, r, http.StatusInternalServerError, Errorf(ENOTFOUND, "USER_NOT_FOUND"))
			return
		}

		_, err := a.db.TeamAddUser(TeamID, User.Id, Role)
		if err != nil {
			Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		Success(w, r, http.StatusOK, nil, nil)
	}
}

// handleTeamRemoveUser handles removing user from a team
// @Summary Remove Team User
// @Description Remove a user from the team
// @Tags team
// @Produce  json
// @Param teamId path string true "the team ID"
// @Param userId path string true "the user ID"
// @Success 200 object standardJsonResponse{}
// @Success 403 object standardJsonResponse{}
// @Success 500 object standardJsonResponse{}
// @Router /{teamId}/users/{userId} [delete]
func (a *api) handleTeamRemoveUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		TeamID := vars["teamId"]
		UserID := vars["userId"]

		err := a.db.TeamRemoveUser(TeamID, UserID)
		if err != nil {
			Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		Success(w, r, http.StatusOK, nil, nil)
	}
}

// handleGetTeamBattles gets a list of battles associated to the team
// @Summary Get Team Battles
// @Description Get a list of battles associated to the team
// @Tags team
// @Produce  json
// @Param teamId path string true "the team ID"
// @Success 200 object standardJsonResponse{data=[]model.Battle}
// @Router /teams/{teamId}/battles [get]
func (a *api) handleGetTeamBattles() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		TeamID := vars["teamId"]
		Limit, Offset := getLimitOffsetFromRequest(r, w)

		Battles := a.db.TeamBattleList(TeamID, Limit, Offset)

		Success(w, r, http.StatusOK, Battles, nil)
	}
}

// handleTeamRemoveBattle handles removing battle from a team
// @Summary Remove Team Battle
// @Description Remove a battle from the team
// @Tags team
// @Produce  json
// @Param teamId path string true "the team ID"
// @Param battleId path string true "the battle ID"
// @Success 200 object standardJsonResponse{}
// @Success 403 object standardJsonResponse{}
// @Success 500 object standardJsonResponse{}
// @Router /{teamId}/battles/{battleId} [delete]
func (a *api) handleTeamRemoveBattle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		TeamID := vars["teamId"]
		BattleID := vars["battleId"]

		err := a.db.TeamRemoveBattle(TeamID, BattleID)
		if err != nil {
			Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		Success(w, r, http.StatusOK, nil, nil)
	}
}

// handleDeleteTeam handles deleting a team
// @Summary Delete Team
// @Description Delete a Team
// @Tags team
// @Produce  json
// @Param teamId path string true "the team ID"
// @Success 200 object standardJsonResponse{}
// @Success 403 object standardJsonResponse{}
// @Success 500 object standardJsonResponse{}
// @Router /{teamId} [delete]
func (a *api) handleDeleteTeam() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		TeamID := vars["teamId"]

		err := a.db.TeamDelete(TeamID)
		if err != nil {
			Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		Success(w, r, http.StatusOK, nil, nil)
	}
}
