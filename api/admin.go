package api

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// handleAppStats gets the applications stats
// @Summary Get Application Stats
// @Description get application stats such as count of registered warriors
// @Tags admin
// @Produce  json
// @Success 200 object standardJsonResponse{data=[]model.ApplicationStats}
// @Failure 500 object standardJsonResponse{}
// @Router /admin/stats [get]
func (a *api) handleAppStats() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		AppStats, err := a.db.GetAppStats()
		if err != nil {
			Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		ActiveBattleUserCount := 0
		for _, s := range h.arenas {
			ActiveBattleUserCount = ActiveBattleUserCount + len(s)
		}

		AppStats.ActiveBattleCount = len(h.arenas)
		AppStats.ActiveBattleUserCount = ActiveBattleUserCount

		Success(w, r, http.StatusOK, AppStats, nil)
	}
}

// handleGetRegisteredUsers gets a list of registered users
// @Summary Get Registered Users
// @Description get list of registered users
// @Tags admin
// @Produce  json
// @Param limit query int true "Max number of results to return"
// @Param offset query int true "Starting point to return rows from, should be multiplied by limit or 0"
// @Success 200 object standardJsonResponse{data=[]model.User}
// @Router /admin/users [get]
func (a *api) handleGetRegisteredUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		Limit, Offset := getLimitOffsetFromRequest(r, w)

		Users := a.db.GetRegisteredUsers(Limit, Offset)

		Success(w, r, http.StatusOK, Users, nil)
	}
}

// handleUserCreate registers a new authenticated user
// @Summary Create Registered User
// @Description Create a registered user
// @Tags admin
// @Produce  json
// @Success 200 object standardJsonResponse{data=model.User}
// @Failure 400 object standardJsonResponse{}
// @Failure 500 object standardJsonResponse{}
// @Router /admin/users [post]
func (a *api) handleUserCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		keyVal := getJSONRequestBody(r, w)

		UserName, UserEmail, UserPassword, accountErr := validateUserAccount(
			keyVal["warriorName"].(string),
			strings.ToLower(keyVal["warriorEmail"].(string)),
			keyVal["warriorPassword1"].(string),
			keyVal["warriorPassword2"].(string),
		)

		if accountErr != nil {
			Failure(w, r, http.StatusBadRequest, accountErr)
			return
		}

		newUser, VerifyID, err := a.db.CreateUserRegistered(UserName, UserEmail, UserPassword, "")
		if err != nil {
			Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		a.email.SendWelcome(UserName, UserEmail, VerifyID)

		Success(w, r, http.StatusOK, newUser, nil)
	}
}

// handleAdminUserDelete attempts to delete a users account
// @Summary Delete User
// @Description Delete a registered user
// @Tags admin
// @Produce  json
// @Param id path int false "the user ID to delete"
// @Success 200 object standardJsonResponse{}
// @Failure 500 object standardJsonResponse{}
// @Router /admin/users/{id} [delete]
func (a *api) handleAdminUserDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		UserID := vars["id"]

		User, UserErr := a.db.GetUser(UserID)
		if UserErr != nil {
			Failure(w, r, http.StatusInternalServerError, UserErr)
			return
		}

		updateErr := a.db.DeleteUser(UserID)
		if updateErr != nil {
			Failure(w, r, http.StatusInternalServerError, updateErr)
			return
		}

		a.email.SendDeleteConfirmation(User.UserName, User.UserEmail)

		Success(w, r, http.StatusOK, nil, nil)
	}
}

// handleUserPromote handles promoting a user to admin
// @Summary Promotes User
// @Description Promotes a user to admin
// @Description Grants read and write access to administrative information
// @Tags admin
// @Produce  json
// @Param id path int false "the user ID to promote"
// @Success 200 object standardJsonResponse{}
// @Failure 500 object standardJsonResponse{}
// @Router /admin/users/{id}/promote/ [patch]
func (a *api) handleUserPromote() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		UserID := vars["id"]

		err := a.db.PromoteUser(UserID)
		if err != nil {
			Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		Success(w, r, http.StatusOK, nil, nil)
	}
}

// handleUserDemote handles demoting a user to registered
// @Summary Demote User
// @Description Demotes a user from admin to registered
// @Tags admin
// @Produce  json
// @Param id path int false "the user ID to demote"
// @Success 200 object standardJsonResponse{}
// @Failure 500 object standardJsonResponse{}
// @Router /admin/users/{id}/demote [patch]
func (a *api) handleUserDemote() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		UserID := vars["id"]

		err := a.db.DemoteUser(UserID)
		if err != nil {
			Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		Success(w, r, http.StatusOK, nil, nil)
	}
}

// handleGetOrganizations gets a list of organizations
// @Summary Get Organizations
// @Description get a list of organizations
// @Tags admin
// @Produce  json
// @Param limit query int true "Max number of results to return"
// @Param offset query int true "Starting point to return rows from, should be multiplied by limit or 0"
// @Success 200 object standardJsonResponse{data=[]model.Organization}
// @Failure 500 object standardJsonResponse{}
// @Router /admin/organizations [get]
func (a *api) handleGetOrganizations() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		Limit, Offset := getLimitOffsetFromRequest(r, w)

		Organizations := a.db.OrganizationList(Limit, Offset)

		Success(w, r, http.StatusOK, Organizations, nil)
	}
}

// handleGetTeams gets a list of teams
// @Summary Get Teams
// @Description get a list of teams
// @Tags admin
// @Produce  json
// @Param limit query int true "Max number of results to return"
// @Param offset query int true "Starting point to return rows from, should be multiplied by limit or 0"
// @Success 200 object standardJsonResponse{data=[]model.Team}
// @Failure 500 object standardJsonResponse{}
// @Router /admin/teams [get]
func (a *api) handleGetTeams() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		Limit, Offset := getLimitOffsetFromRequest(r, w)

		Teams := a.db.TeamList(Limit, Offset)

		Success(w, r, http.StatusOK, Teams, nil)
	}
}

// handleGetAPIKeys gets a list of APIKeys
// @Summary Get API Keys
// @Description get a list of users API Keys
// @Tags admin
// @Produce  json
// @Param limit query int true "Max number of results to return"
// @Param offset query int true "Starting point to return rows from, should be multiplied by limit or 0"
// @Success 200 object standardJsonResponse{data=[]model.Team}
// @Failure 500 object standardJsonResponse{}
// @Router /admin/apikeys [get]
func (a *api) handleGetAPIKeys() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		Limit, Offset := getLimitOffsetFromRequest(r, w)

		Teams := a.db.GetAPIKeys(Limit, Offset)

		Success(w, r, http.StatusOK, Teams, nil)
	}
}
