package api

import (
	"net/http"
	"strings"

	"github.com/spf13/viper"
)

// handleLogin attempts to login the user by comparing email/password to whats in DB
// @Summary Login
// @Description attempts to log the user in with provided credentials
// @Description *Endpoint only available when LDAP is not enabled
// @Tags auth
// @Produce  json
// @Success 200 object standardJsonResponse{data=model.User}
// @Failure 401 object standardJsonResponse{}
// @Failure 500 object standardJsonResponse{}
// @Router /auth [post]
func (a *api) handleLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		keyVal := getJSONRequestBody(r, w)
		UserEmail := strings.ToLower(keyVal["email"].(string))
		UserPassword := keyVal["password"].(string)

		authedUser, sessionId, err := a.db.AuthUser(UserEmail, UserPassword)
		if err != nil {
			a.Failure(w, r, http.StatusUnauthorized, Errorf(EINVALID, "INVALID_LOGIN"))
			return
		}

		cookieErr := a.createSessionCookie(w, sessionId)
		if cookieErr != nil {
			a.Failure(w, r, http.StatusInternalServerError, Errorf(EINVALID, "INVALID_COOKIE"))
			return
		}

		a.Success(w, r, http.StatusOK, authedUser, nil)
	}
}

// handleLdapLogin attempts to authenticate the user by looking up and authenticating
// via ldap, and then creates the user if not existing and logs them in
// @Summary Login LDAP
// @Description attempts to log the user in with provided credentials
// @Description *Endpoint only available when LDAP is enabled
// @Tags auth
// @Produce  json
// @Success 200 object standardJsonResponse{data=model.User}
// @Failure 401 object standardJsonResponse{}
// @Failure 500 object standardJsonResponse{}
// @Router /auth/ldap [post]
func (a *api) handleLdapLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		keyVal := getJSONRequestBody(r, w)
		UserEmail := strings.ToLower(keyVal["email"].(string))
		UserPassword := keyVal["password"].(string)

		authedUser, sessionId, err := a.authAndCreateUserLdap(UserEmail, UserPassword)
		if err != nil {
			a.Failure(w, r, http.StatusUnauthorized, Errorf(EINVALID, "INVALID_LOGIN"))
			return
		}

		cookieErr := a.createSessionCookie(w, sessionId)
		if cookieErr != nil {
			a.Failure(w, r, http.StatusInternalServerError, Errorf(EINVALID, "INVALID_COOKIE"))
			return
		}

		a.Success(w, r, http.StatusOK, authedUser, nil)
	}
}

// handleLogout clears the user cookie(s) ending session
// @Summary Logout
// @Description Logs the user out by deleting session cookies
// @Tags auth
// @Success 200
// @Router /auth/logout [delete]
func (a *api) handleLogout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		SessionId, cookieErr := a.validateSessionCookie(w, r)
		if cookieErr != nil {
			a.Failure(w, r, http.StatusUnauthorized, Errorf(EINVALID, "INVALID_USER"))
			return
		}

		err := a.db.DeleteSession(SessionId)
		if err != nil {
			a.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		a.clearUserCookies(w)
		a.Success(w, r, http.StatusOK, nil, nil)
	}
}

// handleCreateGuestUser registers a user as a guest user
// @Summary Create Guest User
// @Description Registers a user as a guest (non-authenticated)
// @Tags auth
// @Success 200 object standardJsonResponse{data=model.User}
// @Failure 400 object standardJsonResponse{}
// @Failure 500 object standardJsonResponse{}
// @Router /auth/guest [post]
func (a *api) handleCreateGuestUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		AllowGuests := viper.GetBool("config.allow_guests")
		if !AllowGuests {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, "GUESTS_USERS_DISABLED"))
			return
		}

		keyVal := getJSONRequestBody(r, w)

		UserName := keyVal["name"].(string)

		newUser, err := a.db.CreateUserGuest(UserName)
		if err != nil {
			a.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		cookieErr := a.createUserCookie(w, newUser.Id)
		if cookieErr != nil {
			a.Failure(w, r, http.StatusInternalServerError, Errorf(EINVALID, "INVALID_COOKIE"))
			return
		}

		a.Success(w, r, http.StatusOK, newUser, nil)
	}
}

// handleUserRegistration registers a new authenticated user
// @Summary Create User
// @Description Registers a user (authenticated)
// @Tags auth
// @Success 200 object standardJsonResponse{data=model.User}
// @Failure 400 object standardJsonResponse{}
// @Failure 500 object standardJsonResponse{}
// @Router /auth/register [post]
func (a *api) handleUserRegistration() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		AllowRegistration := viper.GetBool("config.allow_registration")
		if !AllowRegistration {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, "USER_REGISTRATION_DISABLED"))
		}

		keyVal := getJSONRequestBody(r, w)

		ActiveUserID, _ := a.validateUserCookie(w, r)

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

		newUser, VerifyID, SessionID, err := a.db.CreateUserRegistered(UserName, UserEmail, UserPassword, ActiveUserID)
		if err != nil {
			a.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		a.email.SendWelcome(UserName, UserEmail, VerifyID)

		if ActiveUserID != "" {
			a.clearUserCookies(w)
		}

		cookieErr := a.createSessionCookie(w, SessionID)
		if cookieErr != nil {
			a.Failure(w, r, http.StatusInternalServerError, Errorf(EINVALID, "INVALID_COOKIE"))
			return
		}

		a.Success(w, r, http.StatusOK, newUser, nil)
	}
}

// handleForgotPassword attempts to send a password reset email
// @Summary Forgot Password
// @Description Sends a forgot password reset email to user
// @Tags auth
// @Success 200 object standardJsonResponse{}
// @Router /auth/forgot-password [post]
func (a *api) handleForgotPassword() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		keyVal := getJSONRequestBody(r, w)
		UserEmail := strings.ToLower(keyVal["email"].(string))

		ResetID, UserName, resetErr := a.db.UserResetRequest(UserEmail)
		if resetErr == nil {
			a.email.SendForgotPassword(UserName, UserEmail, ResetID)
		}

		a.Success(w, r, http.StatusOK, nil, nil)
	}
}

// handleResetPassword attempts to reset a users password
// @Summary Reset Password
// @Description Resets the users password
// @Tags auth
// @Success 200 object standardJsonResponse{}
// @Success 400 object standardJsonResponse{}
// @Success 500 object standardJsonResponse{}
// @Router /auth/reset-password [patch]
func (a *api) handleResetPassword() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		keyVal := getJSONRequestBody(r, w)
		ResetID := keyVal["resetId"].(string)

		UserPassword, passwordErr := validateUserPassword(
			keyVal["password1"].(string),
			keyVal["password2"].(string),
		)

		if passwordErr != nil {
			a.Failure(w, r, http.StatusBadRequest, passwordErr)
			return
		}

		UserName, UserEmail, resetErr := a.db.UserResetPassword(ResetID, UserPassword)
		if resetErr != nil {
			a.Failure(w, r, http.StatusInternalServerError, resetErr)
			return
		}

		a.email.SendPasswordReset(UserName, UserEmail)

		a.Success(w, r, http.StatusOK, nil, nil)
	}
}

// handleUpdatePassword attempts to update a users password
// @Summary Update Password
// @Description Updates the users password
// @Tags auth
// @Success 200 object standardJsonResponse{}
// @Success 400 object standardJsonResponse{}
// @Success 500 object standardJsonResponse{}
// @Security ApiKeyAuth
// @Router /auth/update-password [patch]
func (a *api) handleUpdatePassword() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		keyVal := getJSONRequestBody(r, w)
		UserID := r.Context().Value(contextKeyUserID).(string)

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

// handleAccountVerification attempts to verify a users account
// @Summary Verify User
// @Description Updates the users verified email status
// @Tags auth
// @Success 200 object standardJsonResponse{}
// @Success 500 object standardJsonResponse{}
// @Router /auth/verify [patch]
func (a *api) handleAccountVerification() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		keyVal := getJSONRequestBody(r, w)
		VerifyID := keyVal["verifyId"].(string)

		verifyErr := a.db.VerifyUserAccount(VerifyID)
		if verifyErr != nil {
			a.Failure(w, r, http.StatusInternalServerError, verifyErr)
			return
		}

		a.Success(w, r, http.StatusOK, nil, nil)
	}
}
