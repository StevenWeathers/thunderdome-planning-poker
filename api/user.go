package api

import (
	"encoding/json"
	"io"
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
		ctx := r.Context()
		UserID := ctx.Value(contextKeyUserID).(string)

		User, UserErr := a.db.GetUser(ctx, UserID)
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

		User, UserErr := a.db.GetUser(r.Context(), UserID)
		if UserErr != nil {
			a.Failure(w, r, http.StatusInternalServerError, UserErr)
			return
		}

		a.Success(w, r, http.StatusOK, User, nil)
	}
}

type userprofileUpdateRequestBody struct {
	Name                 string `json:"name" validate:"required,max=64"`
	Avatar               string `json:"avatar" validate:"max=128"`
	NotificationsEnabled bool   `json:"notificationsEnabled"`
	Country              string `json:"country" validate:"len=2"`
	Locale               string `json:"locale" validate:"len=2"`
	Company              string `json:"company" validate:"max=256"`
	JobTitle             string `json:"jobTitle" validate:"max=128"`
	Email                string `json:"email" validate:"email"`
}

// handleUserProfileUpdate attempts to update users profile
// @Summary Update User Profile
// @Description Update a users profile
// @Tags user
// @Produce  json
// @Param userId path string true "the user ID"
// @Param user body userprofileUpdateRequestBody true "the user profile object to update"
// @Success 200 object standardJsonResponse{data=model.User}
// @Failure 403 object standardJsonResponse{}
// @Failure 500 object standardJsonResponse{}
// @Security ApiKeyAuth
// @Router /users/{userId} [put]
func (a *api) handleUserProfileUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		SessionUserType := ctx.Value(contextKeyUserType).(string)
		vars := mux.Vars(r)
		UserID := vars["userId"]

		var profile = userprofileUpdateRequestBody{}
		body, bodyErr := io.ReadAll(r.Body)
		if bodyErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
			return
		}

		jsonErr := json.Unmarshal(body, &profile)
		if jsonErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		inputErr := validate.Struct(profile)
		if inputErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, inputErr.Error()))
			return
		}

		if SessionUserType == adminUserType {
			_, _, vErr := validateUserAccount(profile.Name, profile.Email)
			if vErr != nil {
				a.Failure(w, r, http.StatusBadRequest, vErr)
				return
			}
			updateErr := a.db.UpdateUserAccount(ctx, UserID, profile.Name, profile.Email, profile.Avatar, profile.NotificationsEnabled, profile.Country, profile.Locale, profile.Company, profile.JobTitle)
			if updateErr != nil {
				a.Failure(w, r, http.StatusInternalServerError, updateErr)
				return
			}
		} else {
			var updateErr error
			if !a.config.LdapEnabled {
				if profile.Name == "" {
					a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, "INVALID_USERNAME"))
					return
				}
				updateErr = a.db.UpdateUserProfile(ctx, UserID, profile.Name, profile.Avatar, profile.NotificationsEnabled, profile.Country, profile.Locale, profile.Company, profile.JobTitle)
			} else {
				updateErr = a.db.UpdateUserProfileLdap(ctx, UserID, profile.Avatar, profile.NotificationsEnabled, profile.Country, profile.Locale, profile.Company, profile.JobTitle)
			}
			if updateErr != nil {
				a.Failure(w, r, http.StatusInternalServerError, updateErr)
				return
			}
		}

		user, UserErr := a.db.GetUser(ctx, UserID)
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
		ctx := r.Context()
		UserCookieID := ctx.Value(contextKeyUserID).(string)

		User, UserErr := a.db.GetUser(ctx, UserID)
		if UserErr != nil {
			a.Failure(w, r, http.StatusInternalServerError, UserErr)
			return
		}

		updateErr := a.db.DeleteUser(ctx, UserID)
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

		User, VerifyId, err := a.db.UserVerifyRequest(r.Context(), UserID)
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
		countries, err := a.db.GetActiveCountries(r.Context())

		if err != nil {
			a.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		w.Header().Set("Cache-Control", "max-age=3600") // cache for 1 hour just to decrease load
		a.Success(w, r, http.StatusOK, countries, nil)
	}
}
