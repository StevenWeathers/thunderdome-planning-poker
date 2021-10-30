package api

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

// handleAppStats gets the applications stats
// @Summary Get Application Stats
// @Description get application stats such as count of registered warriors
// @Tags admin
// @Produce  json
// @Success 200
// @Router /admin/stats [get]
func (a *api) handleAppStats() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		AppStats, err := a.db.GetAppStats()

		if err != nil {
			http.NotFound(w, r)
			return
		}

		ActiveBattleUserCount := 0
		for _, s := range h.arenas {
			ActiveBattleUserCount = ActiveBattleUserCount + len(s)
		}

		AppStats.ActiveBattleCount = len(h.arenas)
		AppStats.ActiveBattleUserCount = ActiveBattleUserCount

		a.respondWithJSON(w, http.StatusOK, AppStats)
	}
}

// handleGetRegisteredUsers gets a list of registered users
// @Summary Get Registered Users
// @Description get list of registered users
// @Tags admin
// @Produce  json
// @Param limit path int false "Max number of results to return"
// @Param offset path int false "Starting point to return rows from, should be multiplied by limit or 0"
// @Success 200
// @Router /admin/warriors/{limit}/{offset} [get]
func (a *api) handleGetRegisteredUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		Limit, _ := strconv.Atoi(vars["limit"])
		Offset, _ := strconv.Atoi(vars["offset"])

		Users := a.db.GetRegisteredUsers(Limit, Offset)

		a.respondWithJSON(w, http.StatusOK, Users)
	}
}

// handleUserCreate registers a new authenticated user
// @Summary Create Registered User
// @Description Create a registered user
// @Tags admin
// @Produce  json
// @Success 200
// @Router /admin/warrior [post]
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
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		newUser, VerifyID, err := a.db.CreateUserRegistered(UserName, UserEmail, UserPassword, "")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		a.email.SendWelcome(UserName, UserEmail, VerifyID)

		a.respondWithJSON(w, http.StatusOK, newUser)
	}
}

// handleAdminUserDelete attempts to delete a users account
func (a *api) handleAdminUserDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		UserID := vars["id"]

		User, UserErr := a.db.GetUser(UserID)
		if UserErr != nil {
			log.Println("error finding user : " + UserErr.Error() + "\n")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		updateErr := a.db.DeleteUser(UserID)
		if updateErr != nil {
			log.Println("error attempting to delete user : " + updateErr.Error() + "\n")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		a.email.SendDeleteConfirmation(User.UserName, User.UserEmail)

		return
	}
}

// handleUserPromote handles promoting a user to admin
// @Summary Promotes User
// @Description Promotes a user to admin
// @Description Grants read and write access to administrative information
// @Tags admin
// @Produce  json
// @Success 200
// @Router /admin/promote [post]
func (a *api) handleUserPromote() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		keyVal := a.getJSONRequestBody(r, w)

		err := a.db.PromoteUser(keyVal["warriorId"].(string))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		return
	}
}

// handleUserDemote handles demoting a user to registered
// @Summary Demote User
// @Description Demotes a user from admin to registered
// @Tags admin
// @Produce  json
// @Success 200
// @Router /admin/demote [post]
func (a *api) handleUserDemote() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		keyVal := a.getJSONRequestBody(r, w)

		err := a.db.DemoteUser(keyVal["warriorId"].(string))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		return
	}
}

// handleGetOrganizations gets a list of organizations
// @Summary Get Organizations
// @Description get a list of organizations
// @Tags admin
// @Produce  json
// @Param limit path int false "Max number of results to return"
// @Param offset path int false "Starting point to return rows from, should be multiplied by limit or 0"
// @Success 200
// @Router /admin/organizations/{limit}/{offset} [get]
func (a *api) handleGetOrganizations() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		Limit, _ := strconv.Atoi(vars["limit"])
		Offset, _ := strconv.Atoi(vars["offset"])

		Organizations := a.db.OrganizationList(Limit, Offset)

		a.respondWithJSON(w, http.StatusOK, Organizations)
	}
}

// handleGetTeams gets a list of teams
// @Summary Get Teams
// @Description get a list of teams
// @Tags admin
// @Produce  json
// @Param limit path int false "Max number of results to return"
// @Param offset path int false "Starting point to return rows from, should be multiplied by limit or 0"
// @Success 200
// @Router /admin/teams/{limit}/{offset} [get]
func (a *api) handleGetTeams() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		Limit, _ := strconv.Atoi(vars["limit"])
		Offset, _ := strconv.Atoi(vars["offset"])

		Teams := a.db.TeamList(Limit, Offset)

		a.respondWithJSON(w, http.StatusOK, Teams)
	}
}

// handleGetAPIKeys gets a list of APIKeys
// @Summary Get API Keys
// @Description get a list of users API Keys
// @Tags admin
// @Produce  json
// @Param limit path int false "Max number of results to return"
// @Param offset path int false "Starting point to return rows from, should be multiplied by limit or 0"
// @Success 200
// @Router /admin/apikeys/{limit}/{offset} [get]
func (a *api) handleGetAPIKeys() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		Limit, _ := strconv.Atoi(vars["limit"])
		Offset, _ := strconv.Atoi(vars["offset"])

		Teams := a.db.GetAPIKeys(Limit, Offset)

		a.respondWithJSON(w, http.StatusOK, Teams)
	}
}
