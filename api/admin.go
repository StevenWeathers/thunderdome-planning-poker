package api

import (
	"log"
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
			errors := make([]string, 0)
			errors = append(errors, err.Error())
			a.respondWithStandardJSON(w, http.StatusInternalServerError, false, errors, nil, nil)
		}

		ActiveBattleUserCount := 0
		for _, s := range h.arenas {
			ActiveBattleUserCount = ActiveBattleUserCount + len(s)
		}

		AppStats.ActiveBattleCount = len(h.arenas)
		AppStats.ActiveBattleUserCount = ActiveBattleUserCount

		a.respondWithStandardJSON(w, http.StatusOK, true, nil, AppStats, nil)
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
		Limit, Offset := a.getLimitOffsetFromRequest(r, w)

		Users := a.db.GetRegisteredUsers(Limit, Offset)

		a.respondWithStandardJSON(w, http.StatusOK, true, nil, Users, nil)
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
		keyVal := a.getJSONRequestBody(r, w)

		UserName, UserEmail, UserPassword, accountErr := ValidateUserAccount(
			keyVal["warriorName"].(string),
			strings.ToLower(keyVal["warriorEmail"].(string)),
			keyVal["warriorPassword1"].(string),
			keyVal["warriorPassword2"].(string),
		)

		if accountErr != nil {
			errors := make([]string, 0)
			errors = append(errors, accountErr.Error())
			a.respondWithStandardJSON(w, http.StatusBadRequest, false, errors, nil, nil)
		}

		newUser, VerifyID, err := a.db.CreateUserRegistered(UserName, UserEmail, UserPassword, "")
		if err != nil {
			errors := make([]string, 0)
			errors = append(errors, err.Error())
			a.respondWithStandardJSON(w, http.StatusInternalServerError, false, errors, nil, nil)
		}

		a.email.SendWelcome(UserName, UserEmail, VerifyID)

		a.respondWithStandardJSON(w, http.StatusOK, true, nil, newUser, nil)
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
			log.Println("error finding user : " + UserErr.Error() + "\n")
			errors := make([]string, 0)
			errors = append(errors, UserErr.Error())
			a.respondWithStandardJSON(w, http.StatusInternalServerError, false, errors, nil, nil)
		}

		updateErr := a.db.DeleteUser(UserID)
		if updateErr != nil {
			log.Println("error attempting to delete user : " + updateErr.Error() + "\n")
			errors := make([]string, 0)
			errors = append(errors, updateErr.Error())
			a.respondWithStandardJSON(w, http.StatusInternalServerError, false, errors, nil, nil)
		}

		a.email.SendDeleteConfirmation(User.UserName, User.UserEmail)

		a.respondWithStandardJSON(w, http.StatusOK, true, nil, nil, nil)
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
			errors := make([]string, 0)
			errors = append(errors, err.Error())
			a.respondWithStandardJSON(w, http.StatusInternalServerError, false, errors, nil, nil)
		}

		a.respondWithStandardJSON(w, http.StatusOK, true, nil, nil, nil)
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
			errors := make([]string, 0)
			errors = append(errors, err.Error())
			a.respondWithStandardJSON(w, http.StatusInternalServerError, false, errors, nil, nil)
		}

		a.respondWithStandardJSON(w, http.StatusOK, true, nil, nil, nil)
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
		Limit, Offset := a.getLimitOffsetFromRequest(r, w)

		Organizations := a.db.OrganizationList(Limit, Offset)

		a.respondWithStandardJSON(w, http.StatusOK, true, nil, Organizations, nil)
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
		Limit, Offset := a.getLimitOffsetFromRequest(r, w)

		Teams := a.db.TeamList(Limit, Offset)

		a.respondWithStandardJSON(w, http.StatusOK, true, nil, Teams, nil)
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
		Limit, Offset := a.getLimitOffsetFromRequest(r, w)

		Teams := a.db.GetAPIKeys(Limit, Offset)

		a.respondWithStandardJSON(w, http.StatusOK, true, nil, Teams, nil)
	}
}
