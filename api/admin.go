package api

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// handleAppStats gets the applications stats
// @Summary Get Application Stats
// @Description get application stats such as count of registered users
// @Tags admin
// @Produce  json
// @Success 200 object standardJsonResponse{data=[]model.ApplicationStats}
// @Failure 500 object standardJsonResponse{}
// @Security ApiKeyAuth
// @Router /admin/stats [get]
func (a *api) handleAppStats() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		AppStats, err := a.db.GetAppStats()
		if err != nil {
			a.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		a.Success(w, r, http.StatusOK, AppStats, nil)
	}
}

// handleGetRegisteredUsers gets a list of registered users
// @Summary Get Registered Users
// @Description get list of registered users
// @Tags admin
// @Produce  json
// @Param limit query int false "Max number of results to return"
// @Param offset query int false "Starting point to return rows from, should be multiplied by limit or 0"
// @Success 200 object standardJsonResponse{data=[]model.User}
// @Failure 500 object standardJsonResponse{}
// @Security ApiKeyAuth
// @Router /admin/users [get]
func (a *api) handleGetRegisteredUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		Limit, Offset := getLimitOffsetFromRequest(r, w)

		Users, Count, err := a.db.GetRegisteredUsers(Limit, Offset)
		if err != nil {
			a.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		Meta := &pagination{
			Count:  Count,
			Offset: Offset,
			Limit:  Limit,
		}

		a.Success(w, r, http.StatusOK, Users, Meta)
	}
}

// handleUserCreate registers a new authenticated user
// @Summary Create Registered User
// @Description Create a registered user
// @Tags admin
// @Produce  json
// @Param name body string true "the new users name"
// @Param email body string true "the new users email"
// @Param password1 body string true "the new user's password"
// @Param password2 body string true "the new user's password repeated"
// @Success 200 object standardJsonResponse{data=model.User}
// @Failure 400 object standardJsonResponse{}
// @Failure 500 object standardJsonResponse{}
// @Security ApiKeyAuth
// @Router /admin/users [post]
func (a *api) handleUserCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		keyVal := getJSONRequestBody(r, w)

		UserName, UserEmail, UserPassword, accountErr := validateUserAccountWithPasswords(
			keyVal["name"].(string),
			strings.ToLower(keyVal["email"].(string)),
			keyVal["password1"].(string),
			keyVal["password2"].(string),
		)

		if accountErr != nil {
			a.Failure(w, r, http.StatusBadRequest, accountErr)
			return
		}

		newUser, VerifyID, err := a.db.CreateUser(UserName, UserEmail, UserPassword)
		if err != nil {
			a.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		a.email.SendWelcome(UserName, UserEmail, VerifyID)

		a.Success(w, r, http.StatusOK, newUser, nil)
	}
}

// handleUserPromote handles promoting a user to admin
// @Summary Promotes User
// @Description Promotes a user to admin
// @Description Grants read and write access to administrative information
// @Tags admin
// @Produce  json
// @Param userId path string true "the user ID to promote"
// @Success 200 object standardJsonResponse{}
// @Failure 500 object standardJsonResponse{}
// @Security ApiKeyAuth
// @Router /admin/users/{userId}/promote/ [patch]
func (a *api) handleUserPromote() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		UserID := vars["userId"]

		err := a.db.PromoteUser(UserID)
		if err != nil {
			a.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		a.Success(w, r, http.StatusOK, nil, nil)
	}
}

// handleUserDemote handles demoting a user to registered
// @Summary Demote User
// @Description Demotes a user from admin to registered
// @Tags admin
// @Produce  json
// @Param userId path string true "the user ID to demote"
// @Success 200 object standardJsonResponse{}
// @Failure 500 object standardJsonResponse{}
// @Security ApiKeyAuth
// @Router /admin/users/{userId}/demote [patch]
func (a *api) handleUserDemote() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		UserID := vars["userId"]

		err := a.db.DemoteUser(UserID)
		if err != nil {
			a.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		a.Success(w, r, http.StatusOK, nil, nil)
	}
}

// handleAdminUpdateUserPassword attempts to update a users password
// @Summary Update Password
// @Description Updates the users password
// @Tags admin
// @Param userId path string true "the user ID to update password for"
// @Success 200 object standardJsonResponse{}
// @Success 400 object standardJsonResponse{}
// @Success 500 object standardJsonResponse{}
// @Security ApiKeyAuth
// @Router /admin/users/{userId}/password [patch]
func (a *api) handleAdminUpdateUserPassword() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		keyVal := getJSONRequestBody(r, w)
		vars := mux.Vars(r)
		UserID := vars["userId"]

		UserPassword, passwordErr := validateUserPassword(
			keyVal["password1"].(string),
			keyVal["password2"].(string),
		)

		if passwordErr != nil {
			a.Failure(w, r, http.StatusBadRequest, passwordErr)
			return
		}

		UserName, UserEmail, updateErr := a.db.UserUpdatePassword(UserID, UserPassword)
		if updateErr != nil {
			a.Failure(w, r, http.StatusInternalServerError, updateErr)
			return
		}

		a.email.SendPasswordUpdate(UserName, UserEmail)

		a.Success(w, r, http.StatusOK, nil, nil)
	}
}

// handleGetOrganizations gets a list of organizations
// @Summary Get Organizations
// @Description get a list of organizations
// @Tags admin
// @Produce  json
// @Param limit query int false "Max number of results to return"
// @Param offset query int false "Starting point to return rows from, should be multiplied by limit or 0"
// @Success 200 object standardJsonResponse{data=[]model.Organization}
// @Failure 500 object standardJsonResponse{}
// @Security ApiKeyAuth
// @Router /admin/organizations [get]
func (a *api) handleGetOrganizations() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !a.config.OrganizationsEnabled {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, "ORGANIZATIONS_DISABLED"))
			return
		}
		Limit, Offset := getLimitOffsetFromRequest(r, w)

		Organizations := a.db.OrganizationList(Limit, Offset)

		a.Success(w, r, http.StatusOK, Organizations, nil)
	}
}

// handleGetTeams gets a list of teams
// @Summary Get Teams
// @Description get a list of teams
// @Tags admin
// @Produce  json
// @Param limit query int false "Max number of results to return"
// @Param offset query int false "Starting point to return rows from, should be multiplied by limit or 0"
// @Success 200 object standardJsonResponse{data=[]model.Team}
// @Failure 500 object standardJsonResponse{}
// @Security ApiKeyAuth
// @Router /admin/teams [get]
func (a *api) handleGetTeams() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		Limit, Offset := getLimitOffsetFromRequest(r, w)

		Teams := a.db.TeamList(Limit, Offset)

		a.Success(w, r, http.StatusOK, Teams, nil)
	}
}

// handleGetAPIKeys gets a list of APIKeys
// @Summary Get API Keys
// @Description get a list of users API Keys
// @Tags admin
// @Produce  json
// @Param limit query int false "Max number of results to return"
// @Param offset query int false "Starting point to return rows from, should be multiplied by limit or 0"
// @Success 200 object standardJsonResponse{data=[]model.Team}
// @Failure 500 object standardJsonResponse{}
// @Security ApiKeyAuth
// @Router /admin/apikeys [get]
func (a *api) handleGetAPIKeys() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		Limit, Offset := getLimitOffsetFromRequest(r, w)

		Teams := a.db.GetAPIKeys(Limit, Offset)

		a.Success(w, r, http.StatusOK, Teams, nil)
	}
}

// handleSearchRegisteredUsersByEmail gets a list of registered users filtered by email likeness
// @Summary Search Registered Users by Email
// @Description get list of registered users filtered by email likeness
// @Tags admin
// @Produce  json
// @Param search query string true "The user email to search for"
// @Param limit query int false "Max number of results to return"
// @Param offset query int false "Starting point to return rows from, should be multiplied by limit or 0"
// @Success 200 object standardJsonResponse{data=[]model.User}
// @Failure 400 object standardJsonResponse{}
// @Failure 500 object standardJsonResponse{}
// @Security ApiKeyAuth
// @Router /admin/search/users/email [get]
func (a *api) handleSearchRegisteredUsersByEmail() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		Limit, Offset := getLimitOffsetFromRequest(r, w)
		Search, err := getSearchFromRequest(r, w)
		if err != nil {
			a.Failure(w, r, http.StatusBadRequest, err)
			return
		}

		Users, Count, err := a.db.SearchRegisteredUsersByEmail(Search, Limit, Offset)
		if err != nil {
			a.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		Meta := &pagination{
			Count:  Count,
			Offset: Offset,
			Limit:  Limit,
		}

		a.Success(w, r, http.StatusOK, Users, Meta)
	}
}
