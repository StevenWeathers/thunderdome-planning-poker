package http

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"go.uber.org/zap"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
)

type userLoginRequestBody struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=72"`
}

type loginResponse struct {
	User        *thunderdome.User `json:"user"`
	SessionId   string            `json:"sessionId"`
	MFARequired bool              `json:"mfaRequired"`
	Subscribed  bool              `json:"subscribed"`
}

// handleLogin attempts to log in the user
// @Summary      Login
// @Description  attempts to log the user in with provided credentials
// @Description  *Endpoint only available when LDAP and header auth are not enabled
// @Tags         auth
// @Produce      json
// @Param        credentials  body    userLoginRequestBody  false  "user login object"
// @Success      200          object  standardJsonResponse{data=loginResponse}
// @Failure      401          object  standardJsonResponse{}
// @Failure      500          object  standardJsonResponse{}
// @Router       /auth [post]
func (s *Service) handleLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		body, bodyErr := io.ReadAll(r.Body)
		if bodyErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
			return
		}

		var u = userLoginRequestBody{}
		jsonErr := json.Unmarshal(body, &u)
		if jsonErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		inputErr := validate.Struct(u)
		if inputErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, inputErr.Error()))
			return
		}

		authedUser, sessionId, err := s.AuthDataSvc.AuthUser(ctx, u.Email, u.Password)
		if err != nil {
			userErr := err.Error()
			if userErr == "USER_NOT_FOUND" || userErr == "INVALID_PASSWORD" || userErr == "USER_DISABLED" {
				s.Failure(w, r, http.StatusUnauthorized, Errorf(EINVALID, "INVALID_LOGIN"))
			} else {
				s.Logger.Ctx(ctx).Error("handleLogin error", zap.Error(err),
					zap.String("user_email", sanitizeUserInputForLogs(u.Email)))
				s.Failure(w, r, http.StatusInternalServerError, err)
			}
			return
		}

		subscribed := s.SubscriptionDataSvc.CheckActiveSubscriber(ctx, authedUser.Id)

		res := loginResponse{
			User:        authedUser,
			SessionId:   sessionId,
			MFARequired: authedUser.MFAEnabled,
			Subscribed:  subscribed == nil,
		}

		if authedUser.MFAEnabled {
			s.Success(w, r, http.StatusOK, res, nil)
			return
		}

		cookieErr := s.Cookie.CreateSessionCookie(w, sessionId)
		if cookieErr != nil {
			s.Logger.Ctx(ctx).Error("handleLogin error", zap.Error(cookieErr),
				zap.String("session_id", sessionId), zap.String("session_user_id", authedUser.Id))
			s.Failure(w, r, http.StatusInternalServerError, Errorf(EINVALID, "INVALID_COOKIE"))
			return
		}

		s.Success(w, r, http.StatusOK, res, nil)
	}
}

// handleLdapLogin attempts to authenticate the user by looking up and authenticating
// via ldap, and then creates the user if not existing and logs them in
// @Summary      Login LDAP
// @Description  attempts to log the user in with provided credentials
// @Description  *Endpoint only available when LDAP is enabled
// @Tags         auth
// @Produce      json
// @Param        credentials  body    userLoginRequestBody  false  "user login object"
// @Success      200          object  standardJsonResponse{data=loginResponse}
// @Failure      401          object  standardJsonResponse{}
// @Failure      500          object  standardJsonResponse{}
// @Router       /auth/ldap [post]
func (s *Service) handleLdapLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		body, bodyErr := io.ReadAll(r.Body)
		if bodyErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
			return
		}

		var u = userLoginRequestBody{}
		jsonErr := json.Unmarshal(body, &u)
		if jsonErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		inputErr := validate.Struct(u)
		if inputErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, inputErr.Error()))
			return
		}

		authedUser, sessionId, err := s.authAndCreateUserLdap(ctx, u.Email, u.Password)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleLdapLogin error", zap.Error(err),
				zap.String("user_email", sanitizeUserInputForLogs(u.Email)))
			s.Failure(w, r, http.StatusUnauthorized, Errorf(EINVALID, "INVALID_LOGIN"))
			return
		}

		res := loginResponse{
			User:        authedUser,
			SessionId:   sessionId,
			MFARequired: authedUser.MFAEnabled,
		}

		if authedUser.MFAEnabled {
			s.Success(w, r, http.StatusOK, res, nil)
			return
		}

		cookieErr := s.Cookie.CreateSessionCookie(w, sessionId)
		if cookieErr != nil {
			s.Logger.Ctx(ctx).Error("handleLdapLogin error", zap.Error(cookieErr),
				zap.String("session_user_id", authedUser.Id), zap.String("session_id", sessionId))
			s.Failure(w, r, http.StatusInternalServerError, Errorf(EINVALID, "INVALID_COOKIE"))
			return
		}

		s.Success(w, r, http.StatusOK, res, nil)
	}
}

// handleHeaderLogin authenticates the user by looking at configurable headers
// and then creates the user if one does not exist and logs them in
// @Summary      Login Header
// @Description  attempts to log the user in with provided credentials
// @Description  *Endpoint only available when Header auth is enabled
// @Tags         auth
// @Produce      json
// @Success      200  object  standardJsonResponse{data=loginResponse}
// @Failure      401  object  standardJsonResponse{}
// @Failure      500  object  standardJsonResponse{}
// @Router       /auth [get]
func (s *Service) handleHeaderLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		username := r.Header.Get(s.Config.AuthHeaderUsernameHeader)
		useremail := r.Header.Get(s.Config.AuthHeaderEmailHeader)

		if username == "" {
			s.Failure(w, r, http.StatusUnauthorized, Errorf(EUNAUTHORIZED, "MISSING_AUTH_HEADER"))
			return
		}

		authedUser, sessionId, err := s.authAndCreateUserHeader(ctx, username, useremail)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleHeaderLogin error", zap.Error(err),
				zap.String("user_name", sanitizeUserInputForLogs(username)),
				zap.String("user_email", sanitizeUserInputForLogs(username)))
			s.Failure(w, r, http.StatusUnauthorized, Errorf(EINVALID, "INVALID_LOGIN"))
			return
		}

		res := loginResponse{
			User:        authedUser,
			SessionId:   sessionId,
			MFARequired: authedUser.MFAEnabled,
		}

		if authedUser.MFAEnabled {
			s.Success(w, r, http.StatusOK, res, nil)
			return
		}

		cookieErr := s.Cookie.CreateSessionCookie(w, sessionId)
		if cookieErr != nil {
			s.Logger.Ctx(ctx).Error("handleHeaderLogin error", zap.Error(cookieErr),
				zap.String("session_user_id", authedUser.Id), zap.String("session_id", sessionId))
			s.Failure(w, r, http.StatusInternalServerError, Errorf(EINVALID, "INVALID_COOKIE"))
			return
		}

		s.Success(w, r, http.StatusOK, res, nil)
	}
}

type mfaLoginRequestBody struct {
	Passcode  string `json:"passcode" validate:"required"`
	SessionId string `json:"sessionId" validate:"required"`
}

// handleMFALogin attempts to log in the user with MFA token
// @Summary      MFA Login
// @Description  attempts to log the user in with provided MFA token
// @Tags         auth
// @Produce      json
// @Param        credentials  body    mfaLoginRequestBody  false  "mfa login object"
// @Success      200          object  standardJsonResponse{}
// @Failure      401          object  standardJsonResponse{}
// @Failure      500          object  standardJsonResponse{}
// @Router       /auth/mfa [post]
func (s *Service) handleMFALogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		body, bodyErr := io.ReadAll(r.Body)
		if bodyErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
			return
		}

		var u = mfaLoginRequestBody{}
		jsonErr := json.Unmarshal(body, &u)
		if jsonErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		inputErr := validate.Struct(u)
		if inputErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, inputErr.Error()))
			return
		}

		err := s.AuthDataSvc.MFATokenValidate(ctx, u.SessionId, u.Passcode)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleMFALogin error", zap.Error(err),
				zap.String("session_id", u.SessionId))
			s.Failure(w, r, http.StatusUnauthorized, Errorf(EINVALID, "INVALID_AUTHENTICATOR_TOKEN"))
			return
		}

		cookieErr := s.Cookie.CreateSessionCookie(w, u.SessionId)
		if cookieErr != nil {
			s.Logger.Ctx(ctx).Error("handleMFALogin error", zap.Error(cookieErr),
				zap.String("session_id", u.SessionId))
			s.Failure(w, r, http.StatusInternalServerError, Errorf(EINVALID, "INVALID_COOKIE"))
			return
		}

		s.Success(w, r, http.StatusOK, nil, nil)
	}
}

// handleLogout clears the user Cookie(s) ending session
// @Summary      Logout
// @Description  Logs the user out by deleting session cookies
// @Tags         auth
// @Success      200
// @Router       /auth/logout [delete]
func (s *Service) handleLogout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		SessionId, cookieErr := s.Cookie.ValidateSessionCookie(w, r)
		if cookieErr != nil {
			s.Failure(w, r, http.StatusUnauthorized, Errorf(EINVALID, "INVALID_USER"))
			return
		}

		err := s.AuthDataSvc.DeleteSession(ctx, SessionId)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleLogout error", zap.Error(err),
				zap.String("session_id", SessionId))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Cookie.ClearUserCookies(w)
		s.Success(w, r, http.StatusOK, nil, nil)
	}
}

type guestUserCreateRequestBody struct {
	Name string `json:"name" validate:"required"`
}

// handleCreateGuestUser registers a user as a guest user
// @Summary      Create Guest User
// @Description  Registers a user as a guest (non-authenticated)
// @Tags         auth
// @Produce      json
// @Param        user  body    guestUserCreateRequestBody  false  "guest user object"
// @Success      200   object  standardJsonResponse{data=thunderdome.User}
// @Failure      400   object  standardJsonResponse{}
// @Failure      500   object  standardJsonResponse{}
// @Router       /auth/guest [post]
func (s *Service) handleCreateGuestUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		AllowGuests := s.Config.AllowGuests
		if !AllowGuests {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, "GUESTS_USERS_DISABLED"))
			return
		}

		body, bodyErr := io.ReadAll(r.Body)
		if bodyErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
			return
		}

		var u = guestUserCreateRequestBody{}
		jsonErr := json.Unmarshal(body, &u)
		if jsonErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		inputErr := validate.Struct(u)
		if inputErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, inputErr.Error()))
			return
		}

		newUser, err := s.UserDataSvc.CreateUserGuest(ctx, u.Name)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleCreateGuestUser error", zap.Error(err),
				zap.String("user_name", u.Name))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		cookieErr := s.Cookie.CreateUserCookie(w, newUser.Id)
		if cookieErr != nil {
			s.Logger.Ctx(ctx).Error("handleCreateGuestUser error", zap.Error(cookieErr),
				zap.String("session_user_id", newUser.Id))
			s.Failure(w, r, http.StatusInternalServerError, Errorf(EINVALID, "INVALID_COOKIE"))
			return
		}

		s.Success(w, r, http.StatusOK, newUser, nil)
	}
}

type userRegisterRequestBody struct {
	Name                 string `json:"name" validate:"required"`
	Email                string `json:"email" validate:"required,email"`
	Password1            string `json:"password1" validate:"required,min=6,max=72"`
	Password2            string `json:"password2" validate:"required,min=6,max=72,eqfield=Password1"`
	TeamInviteId         string `json:"teamInviteId"`
	OrganizationInviteId string `json:"orgInviteId"`
}

// handleUserRegistration registers a new authenticated user
// @Summary      Create User
// @Description  Registers a user (authenticated)
// @Tags         auth
// @Produce      json
// @Param        user  body    userRegisterRequestBody  false  "new user object"
// @Success      200   object  standardJsonResponse{data=thunderdome.User}
// @Failure      400   object  standardJsonResponse{}
// @Failure      500   object  standardJsonResponse{}
// @Router       /auth/register [post]
func (s *Service) handleUserRegistration() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		AllowRegistration := s.Config.AllowRegistration
		if !AllowRegistration {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, "USER_REGISTRATION_DISABLED"))
		}

		body, bodyErr := io.ReadAll(r.Body)
		if bodyErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
			return
		}

		var u = userRegisterRequestBody{}
		jsonErr := json.Unmarshal(body, &u)
		if jsonErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		inputErr := validate.Struct(u)
		if inputErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, inputErr.Error()))
			return
		}

		ActiveUserID, _ := s.Cookie.ValidateUserCookie(w, r)

		userEmail := u.Email
		var teamInvite thunderdome.TeamUserInvite
		var orgInvite thunderdome.OrganizationUserInvite

		if u.TeamInviteId != "" {
			var err error
			teamInvite, err = s.TeamDataSvc.TeamUserGetInviteByID(ctx, u.TeamInviteId)
			if err != nil {
				s.Failure(w, r, http.StatusInternalServerError, err)
				return
			}
			if userEmail != teamInvite.Email {
				s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, err.Error()))
				return
			}
		}

		if u.OrganizationInviteId != "" {
			var err error
			orgInvite, err = s.OrganizationDataSvc.OrganizationUserGetInviteByID(ctx, u.OrganizationInviteId)
			if err != nil {
				s.Failure(w, r, http.StatusInternalServerError, err)
				return
			}
			if userEmail != orgInvite.Email {
				s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, err.Error()))
				return
			}
		}

		UserName, UserEmail, UserPassword, accountErr := validateUserAccountWithPasswords(
			u.Name,
			userEmail,
			u.Password1,
			u.Password2,
		)

		if accountErr != nil {
			s.Failure(w, r, http.StatusBadRequest, accountErr)
			return
		}

		newUser, VerifyID, err := s.UserDataSvc.CreateUserRegistered(ctx, UserName, UserEmail, UserPassword, ActiveUserID)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleUserRegistration error", zap.Error(err),
				zap.String("user_email", sanitizeUserInputForLogs(UserEmail)),
				zap.String("active_user_id", ActiveUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		_ = s.Email.SendWelcome(UserName, UserEmail, VerifyID)

		if u.TeamInviteId != "" {
			_, inviteErr := s.TeamDataSvc.TeamAddUser(ctx, teamInvite.TeamId, newUser.Id, teamInvite.Role)
			if inviteErr != nil {
				s.Logger.Ctx(ctx).Error("handleUserRegistration error adding invited user to team", zap.Error(inviteErr),
					zap.String("session_user_id", newUser.Id),
					zap.String("invite_id", teamInvite.InviteId))
			}

			delInviteErr := s.TeamDataSvc.TeamDeleteUserInvite(ctx, teamInvite.InviteId)
			if delInviteErr != nil {
				s.Logger.Ctx(ctx).Error("handleUserRegistration error deleting user invite to team", zap.Error(delInviteErr),
					zap.String("session_user_id", newUser.Id),
					zap.String("invite_id", teamInvite.InviteId))
			}
		}

		if u.OrganizationInviteId != "" {
			_, inviteErr := s.OrganizationDataSvc.OrganizationAddUser(ctx, orgInvite.OrganizationId, newUser.Id, orgInvite.Role)
			if inviteErr != nil {
				s.Logger.Ctx(ctx).Error("handleUserRegistration error adding invited user to organization", zap.Error(inviteErr),
					zap.String("session_user_id", newUser.Id),
					zap.String("invite_id", orgInvite.InviteId))
			}

			delInviteErr := s.OrganizationDataSvc.OrganizationDeleteUserInvite(ctx, orgInvite.InviteId)
			if delInviteErr != nil {
				s.Logger.Ctx(ctx).Error("handleUserRegistration error deleting user invite to organization", zap.Error(delInviteErr),
					zap.String("session_user_id", newUser.Id),
					zap.String("invite_id", orgInvite.InviteId))
			}
		}

		if ActiveUserID != "" {
			s.Cookie.ClearUserCookies(w)
		}

		SessionID, err := s.AuthDataSvc.CreateSession(ctx, newUser.Id)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleUserRegistration error", zap.Error(err),
				zap.String("session_user_id", newUser.Id))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		cookieErr := s.Cookie.CreateSessionCookie(w, SessionID)
		if cookieErr != nil {
			s.Logger.Ctx(ctx).Error("handleUserRegistration error", zap.Error(cookieErr),
				zap.String("session_user_id", newUser.Id),
				zap.String("session_id", SessionID))
			s.Failure(w, r, http.StatusInternalServerError, Errorf(EINVALID, "INVALID_COOKIE"))
			return
		}

		s.Success(w, r, http.StatusOK, newUser, nil)
	}
}

type forgotPasswordRequestBody struct {
	Email string `json:"email" validate:"required,email"`
}

// handleForgotPassword attempts to send a password reset Email
// @Summary      Forgot Password
// @Description  Sends a forgot password reset Email to user
// @Tags         auth
// @Produce      json
// @Param        user  body    forgotPasswordRequestBody  false  "forgot password object"
// @Success      200   object  standardJsonResponse{}
// @Router       /auth/forgot-password [post]
func (s *Service) handleForgotPassword() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		body, bodyErr := io.ReadAll(r.Body)
		if bodyErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
			return
		}

		var u = forgotPasswordRequestBody{}
		jsonErr := json.Unmarshal(body, &u)
		if jsonErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		inputErr := validate.Struct(u)
		if inputErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, inputErr.Error()))
			return
		}

		UserEmail := strings.ToLower(u.Email)

		ResetID, UserName, resetErr := s.AuthDataSvc.UserResetRequest(ctx, UserEmail)
		if resetErr == nil {
			_ = s.Email.SendForgotPassword(UserName, UserEmail, ResetID)
		} else {
			s.Logger.Ctx(ctx).Error("handleForgotPassword error", zap.Error(resetErr),
				zap.String("user_email", sanitizeUserInputForLogs(UserEmail)))
		}

		s.Success(w, r, http.StatusOK, nil, nil)
	}
}

type resetPasswordRequestBody struct {
	ResetID   string `json:"resetId" validate:"required"`
	Password1 string `json:"password1" validate:"required,min=6,max=72"`
	Password2 string `json:"password2" validate:"required,min=6,max=72,eqfield=Password1"`
}

// handleResetPassword attempts to reset a user's password
// @Summary      Reset Password
// @Description  Resets the user's password
// @Tags         auth
// @Produce      json
// @Param        reset  body    resetPasswordRequestBody  false  "reset password object"
// @Success      200    object  standardJsonResponse{}
// @Success      400    object  standardJsonResponse{}
// @Success      500    object  standardJsonResponse{}
// @Router       /auth/reset-password [patch]
func (s *Service) handleResetPassword() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		body, bodyErr := io.ReadAll(r.Body)
		if bodyErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
			return
		}

		var u = resetPasswordRequestBody{}
		jsonErr := json.Unmarshal(body, &u)
		if jsonErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		inputErr := validate.Struct(u)
		if inputErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, inputErr.Error()))
			return
		}

		UserName, UserEmail, resetErr := s.AuthDataSvc.UserResetPassword(ctx, u.ResetID, u.Password1)
		if resetErr != nil {
			s.Logger.Ctx(ctx).Error("handleResetPassword error", zap.Error(resetErr),
				zap.String("reset_id", u.ResetID))
			s.Failure(w, r, http.StatusInternalServerError, resetErr)
			return
		}

		_ = s.Email.SendPasswordReset(UserName, UserEmail)

		s.Success(w, r, http.StatusOK, nil, nil)
	}
}

type updatePasswordRequestBody struct {
	Password1 string `json:"password1" validate:"required,min=6,max=72"`
	Password2 string `json:"password2" validate:"required,min=6,max=72,eqfield=Password1"`
}

// handleUpdatePassword attempts to update a user's password
// @Summary      Update Password
// @Description  Updates the user's password
// @Tags         auth
// @Produce      json
// @Param        passwords  body    updatePasswordRequestBody  false  "update password object"
// @Success      200        object  standardJsonResponse{}
// @Success      400        object  standardJsonResponse{}
// @Success      500        object  standardJsonResponse{}
// @Security     ApiKeyAuth
// @Router       /auth/update-password [patch]
func (s *Service) handleUpdatePassword() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		SessionUserID := r.Context().Value(contextKeyUserID).(string)
		body, bodyErr := io.ReadAll(r.Body)
		if bodyErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
			return
		}

		var u = updatePasswordRequestBody{}
		jsonErr := json.Unmarshal(body, &u)
		if jsonErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		inputErr := validate.Struct(u)
		if inputErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, inputErr.Error()))
			return
		}

		UserName, UserEmail, updateErr := s.AuthDataSvc.UserUpdatePassword(ctx, SessionUserID, u.Password1)
		if updateErr != nil {
			s.Logger.Ctx(ctx).Error("handleResetPassword error", zap.Error(updateErr),
				zap.String("session_user_id", SessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, updateErr)
			return
		}

		_ = s.Email.SendPasswordUpdate(UserName, UserEmail)

		s.Success(w, r, http.StatusOK, nil, nil)
	}
}

type verificationRequestBody struct {
	VerifyID string `json:"verifyId" validate:"required"`
}

// handleAccountVerification attempts to verify a users account
// @Summary      Verify User
// @Description  Updates the users verified Email status
// @Tags         auth
// @Produce      json
// @Param        verify  body    verificationRequestBody  false  "verify object"
// @Success      200     object  standardJsonResponse{}
// @Success      500     object  standardJsonResponse{}
// @Router       /auth/verify [patch]
func (s *Service) handleAccountVerification() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		body, bodyErr := io.ReadAll(r.Body)
		if bodyErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
			return
		}

		var u = verificationRequestBody{}
		jsonErr := json.Unmarshal(body, &u)
		if jsonErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		inputErr := validate.Struct(u)
		if inputErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, inputErr.Error()))
			return
		}

		verifyErr := s.AuthDataSvc.VerifyUserAccount(ctx, u.VerifyID)
		if verifyErr != nil {
			s.Logger.Ctx(ctx).Error("handleAccountVerification error", zap.Error(verifyErr),
				zap.String("verify_id", u.VerifyID))
			s.Failure(w, r, http.StatusInternalServerError, verifyErr)
			return
		}

		s.Success(w, r, http.StatusOK, nil, nil)
	}
}

// handleMFASetupGenerate generates the MFA secret and QR code for setup
// @Summary      MFA Setup Generate secret and QR code
// @Description  Generates MFA secret and QR Code
// @Tags         auth
// @Success      200
// @Router       /auth/mfa/setup/generate [post]
func (s *Service) handleMFASetupGenerate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		SessionUserID := ctx.Value(contextKeyUserID).(string)

		u, err := s.UserDataSvc.GetUser(ctx, SessionUserID)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleMFASetupGenerate error", zap.Error(err),
				zap.String("session_user_id", SessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		secret, png64, err := s.AuthDataSvc.MFASetupGenerate(u.Email)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleMFASetupGenerate error", zap.Error(err),
				zap.String("session_user_id", SessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		type result struct {
			Secret string `json:"secret"`
			QRCode string `json:"qrCode"`
		}

		s.Success(w, r, http.StatusOK, result{Secret: secret, QRCode: png64}, nil)
	}
}

type mfaSetupValidateRequestBody struct {
	Secret   string `json:"secret" validate:"required"`
	Passcode string `json:"passcode" validate:"required"`
}

// handleMFASetupValidate validates the passcode for MFA secret during setup
// @Summary      Validate MFA Setup passcode
// @Description  Validates the passcode for the MFA secret
// @Param        verify  body  mfaSetupValidateRequestBody  false  "verify object"
// @Tags         auth
// @Success      200
// @Router       /auth/mfa/setup/validate [post]
func (s *Service) handleMFASetupValidate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		SessionUserID := ctx.Value(contextKeyUserID).(string)

		body, bodyErr := io.ReadAll(r.Body)
		if bodyErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
			return
		}

		var v = mfaSetupValidateRequestBody{}
		jsonErr := json.Unmarshal(body, &v)
		if jsonErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		inputErr := validate.Struct(v)
		if inputErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, inputErr.Error()))
			return
		}

		type result struct {
			Result string `json:"result"`
		}
		res := result{Result: "SUCCESS"}

		err := s.AuthDataSvc.MFASetupValidate(ctx, SessionUserID, v.Secret, v.Passcode)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleMFASetupValidate error", zap.Error(err),
				zap.String("session_user_id", SessionUserID))
			res.Result = err.Error()
		}

		s.Success(w, r, http.StatusOK, res, nil)
	}
}

// handleMFARemove removes MFA requirement from user auth
// @Summary      Remove MFA
// @Description  Removes MFA requirement from user auth
// @Tags         auth
// @Success      200
// @Router       /auth/mfa [delete]
func (s *Service) handleMFARemove() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		SessionUserID := ctx.Value(contextKeyUserID).(string)

		err := s.AuthDataSvc.MFARemove(ctx, SessionUserID)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleMFARemove error", zap.Error(err),
				zap.String("session_user_id", SessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, nil, nil)
	}
}

// handleGetTeamInvite gets a team user invite details
// @Summary      Get Team Invite
// @Description  Get a team user invite details
// @Tags         auth
// @Produce      json
// @Param        inviteId  path    string  true  "the invite ID"
// @Success      200     object  standardJsonResponse{data=[]thunderdome.TeamUserInvite}
// @Router       /auth/invite/team/{inviteId} [get]
func (s *Service) handleGetTeamInvite() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		vars := mux.Vars(r)
		InviteID := vars["inviteId"]

		invite, err := s.TeamDataSvc.TeamUserGetInviteByID(ctx, InviteID)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleGetTeamInvite error", zap.Error(err), zap.String("invite_id", InviteID))
			s.Failure(w, r, http.StatusInternalServerError, err)
		}

		s.Success(w, r, http.StatusOK, invite, nil)
	}
}

// handleGetOrganizationInvite gets a organization user invite details
// @Summary      Get Organization Invite
// @Description  Get a organization user invite details
// @Tags         auth
// @Produce      json
// @Param        inviteId  path    string  true  "the invite ID"
// @Success      200     object  standardJsonResponse{data=[]thunderdome.OrganizationUserInvite}
// @Router       /auth/invite/organization/{inviteId} [get]
func (s *Service) handleGetOrganizationInvite() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		vars := mux.Vars(r)
		InviteID := vars["inviteId"]

		invite, err := s.OrganizationDataSvc.OrganizationUserGetInviteByID(ctx, InviteID)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleGetOrganizationInvite error", zap.Error(err), zap.String("invite_id", InviteID))
			s.Failure(w, r, http.StatusInternalServerError, err)
		}

		s.Success(w, r, http.StatusOK, invite, nil)
	}
}
