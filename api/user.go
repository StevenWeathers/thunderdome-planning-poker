package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

// handleSessionUserProfile returns the users profile by session user ID
// @Summary Get Session User Profile
// @Description Gets a users profile by session user ID
// @Tags auth, user
// @Produce  json
// @Success 200 object standardJsonResponse{data=model.User}
// @Failure 403 object standardJsonResponse{}
// @Failure 500 object standardJsonResponse{}
// @Security ApiKeyAuth
// @Router /auth/user [get]
func (a *api) handleSessionUserProfile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		UserID := r.Context().Value(contextKeyUserID).(string)

		User, UserErr := a.db.GetUser(UserID)
		if UserErr != nil {
			a.Failure(w, r, http.StatusInternalServerError, UserErr)
			return
		}

		a.Success(w, r, http.StatusOK, User, nil)
	}
}

// handleUserProfile returns the users profile if it matches their session
// @Summary Get User Profile
// @Description Gets a users profile
// @Tags user
// @Produce  json
// @Param userId path string true "the user ID"
// @Success 200 object standardJsonResponse{data=model.User}
// @Failure 403 object standardJsonResponse{}
// @Failure 500 object standardJsonResponse{}
// @Security ApiKeyAuth
// @Router /users/{userId} [get]
func (a *api) handleUserProfile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		UserID := vars["userId"]

		User, UserErr := a.db.GetUser(UserID)
		if UserErr != nil {
			a.Failure(w, r, http.StatusInternalServerError, UserErr)
			return
		}

		a.Success(w, r, http.StatusOK, User, nil)
	}
}

// handleUserProfileUpdate attempts to update users profile
// @Summary Update User Profile
// @Description Update a users profile
// @Tags user
// @Produce  json
// @Param userId path string true "the user ID"
// @Param name body string true "the users name"
// @Param avatar body string true "avatar"
// @Param notificationsEnabled body boolean true "whether battle notifications are enabled"
// @Param country body string false "user's country code e.g. US"
// @Param locale body string false "the user's locale e.g. en"
// @Param company body string false "the user's company name"
// @Param jobTitle body string false "the user's job title"
// @Param email body string false "the user's email [admin only param]"
// @Success 200 object standardJsonResponse{data=model.User}
// @Failure 403 object standardJsonResponse{}
// @Failure 500 object standardJsonResponse{}
// @Security ApiKeyAuth
// @Router /users/{userId} [put]
func (a *api) handleUserProfileUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		SessionUserType := r.Context().Value(contextKeyUserType).(string)
		vars := mux.Vars(r)
		keyVal := getJSONRequestBody(r, w)
		UserName := keyVal["name"].(string)
		UserAvatar := keyVal["avatar"].(string)
		NotificationsEnabled, _ := keyVal["notificationsEnabled"].(bool)
		Country := keyVal["country"].(string)
		Locale := keyVal["locale"].(string)
		Company := keyVal["company"].(string)
		JobTitle := keyVal["jobTitle"].(string)
		UserID := vars["userId"]
		var UserEmail string
		if keyVal["email"] != nil {
			UserEmail = keyVal["email"].(string)
		}

		if SessionUserType == adminUserType {
			_, _, vErr := validateUserAccount(UserName, UserEmail)
			if vErr != nil {
				a.Failure(w, r, http.StatusBadRequest, vErr)
				return
			}
			updateErr := a.db.UpdateUserAccount(UserID, UserName, UserEmail, UserAvatar, NotificationsEnabled, Country, Locale, Company, JobTitle)
			if updateErr != nil {
				a.Failure(w, r, http.StatusInternalServerError, updateErr)
				return
			}
		} else {
			var updateErr error
			if a.config.LdapEnabled == false {
				if UserName == "" {
					a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, "INVALID_USERNAME"))
					return
				}
				updateErr = a.db.UpdateUserProfile(UserID, UserName, UserAvatar, NotificationsEnabled, Country, Locale, Company, JobTitle)
			} else {
				updateErr = a.db.UpdateUserProfileLdap(UserID, UserAvatar, NotificationsEnabled, Country, Locale, Company, JobTitle)
			}
			if updateErr != nil {
				a.Failure(w, r, http.StatusInternalServerError, updateErr)
				return
			}
		}

		user, UserErr := a.db.GetUser(UserID)
		if UserErr != nil {
			a.Failure(w, r, http.StatusInternalServerError, UserErr)
			return
		}

		a.Success(w, r, http.StatusOK, user, nil)
	}
}

// handleUserDelete attempts to delete a users account
// @Summary Delete User
// @Description Deletes a user
// @Tags user
// @Produce  json
// @Param userId path string true "the user ID"
// @Success 200 object standardJsonResponse{}
// @Failure 403 object standardJsonResponse{}
// @Failure 500 object standardJsonResponse{}
// @Security ApiKeyAuth
// @Router /users/{userId} [delete]
func (a *api) handleUserDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		UserID := vars["userId"]
		UserCookieID := r.Context().Value(contextKeyUserID).(string)

		User, UserErr := a.db.GetUser(UserID)
		if UserErr != nil {
			a.Failure(w, r, http.StatusInternalServerError, UserErr)
			return
		}

		updateErr := a.db.DeleteUser(UserID)
		if updateErr != nil {
			a.Failure(w, r, http.StatusInternalServerError, updateErr)
			return
		}

		a.email.SendDeleteConfirmation(User.Name, User.Email)

		// don't clear admins user cookies when deleting other users
		if UserID == UserCookieID {
			a.clearUserCookies(w)
		}

		a.Success(w, r, http.StatusOK, nil, nil)
	}
}

// handleVerifyRequest sends verification email
// @Summary Request Verification Email
// @Description Sends verification email
// @Tags user
// @Param userId path string true "the user ID"
// @Success 200 object standardJsonResponse{}
// @Success 400 object standardJsonResponse{}
// @Success 500 object standardJsonResponse{}
// @Router /users/{userId}/request-verify [post]
func (a *api) handleVerifyRequest() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		UserID := vars["userId"]

		User, VerifyId, err := a.db.UserVerifyRequest(UserID)
		if err != nil {
			a.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		a.email.SendEmailVerification(User.Name, User.Email, VerifyId)

		a.Success(w, r, http.StatusOK, nil, nil)
	}
}

// handleGetActiveCountries gets a list of registered users countries
// @Summary Get Active Countries
// @Description Gets a list of users countries
// @Produce  json
// @Success 200 object standardJsonResponse{[]string}
// @Failure 500 object standardJsonResponse{}
// @Router /active-countries [get]
func (a *api) handleGetActiveCountries() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		countries, err := a.db.GetActiveCountries()

		if err != nil {
			a.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		w.Header().Set("Cache-Control", "max-age=3600") // cache for 1 hour just to decrease load
		a.Success(w, r, http.StatusOK, countries, nil)
	}
}
