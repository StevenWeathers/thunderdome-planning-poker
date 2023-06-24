package api

import (
	"encoding/json"
	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
	"io"
	"net/http"
	"strings"

	"github.com/spf13/viper"
)

type userLoginRequestBody struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=72"`
}

type loginResponse struct {
	User        *thunderdome.User `json:"user"`
	SessionId   string            `json:"sessionId"`
	MFARequired bool              `json:"mfaRequired"`
}

// handleLogin attempts to log in the user
// @Summary Login
// @Description attempts to log the user in with provided credentials
// @Description *Endpoint only available when LDAP and header auth are not enabled
// @Tags auth
// @Produce  json
// @Param credentials body userLoginRequestBody false "user login object"
// @Success 200 object standardJsonResponse{data=loginResponse}
// @Failure 401 object standardJsonResponse{}
// @Failure 500 object standardJsonResponse{}
// @Router /auth [post]
func (a *Service) handleLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, bodyErr := io.ReadAll(r.Body)
		if bodyErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
			return
		}

		var u = userLoginRequestBody{}
		jsonErr := json.Unmarshal(body, &u)
		if jsonErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		inputErr := validate.Struct(u)
		if inputErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, inputErr.Error()))
			return
		}

		authedUser, sessionId, err := a.AuthService.AuthUser(r.Context(), u.Email, u.Password)
		if err != nil {
			userErr := err.Error()
			if userErr == "USER_NOT_FOUND" || userErr == "INVALID_PASSWORD" || userErr == "USER_DISABLED" {
				a.Failure(w, r, http.StatusUnauthorized, Errorf(EINVALID, "INVALID_LOGIN"))
			} else {
				a.Failure(w, r, http.StatusInternalServerError, err)
			}
			return
		}

		res := loginResponse{
			User:        authedUser,
			SessionId:   sessionId,
			MFARequired: authedUser.MFAEnabled,
		}

		if authedUser.MFAEnabled {
			a.Success(w, r, http.StatusOK, res, nil)
			return
		}

		cookieErr := a.createSessionCookie(w, sessionId)
		if cookieErr != nil {
			a.Failure(w, r, http.StatusInternalServerError, Errorf(EINVALID, "INVALID_COOKIE"))
			return
		}

		a.Success(w, r, http.StatusOK, res, nil)
	}
}

// handleLdapLogin attempts to authenticate the user by looking up and authenticating
// via ldap, and then creates the user if not existing and logs them in
// @Summary Login LDAP
// @Description attempts to log the user in with provided credentials
// @Description *Endpoint only available when LDAP is enabled
// @Tags auth
// @Produce json
// @Param credentials body userLoginRequestBody false "user login object"
// @Success 200 object standardJsonResponse{data=loginResponse}
// @Failure 401 object standardJsonResponse{}
// @Failure 500 object standardJsonResponse{}
// @Router /auth/ldap [post]
func (a *Service) handleLdapLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, bodyErr := io.ReadAll(r.Body)
		if bodyErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
			return
		}

		var u = userLoginRequestBody{}
		jsonErr := json.Unmarshal(body, &u)
		if jsonErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		inputErr := validate.Struct(u)
		if inputErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, inputErr.Error()))
			return
		}

		authedUser, sessionId, err := a.authAndCreateUserLdap(r.Context(), u.Email, u.Password)
		if err != nil {
			a.Failure(w, r, http.StatusUnauthorized, Errorf(EINVALID, "INVALID_LOGIN"))
			return
		}

		res := loginResponse{
			User:        authedUser,
			SessionId:   sessionId,
			MFARequired: authedUser.MFAEnabled,
		}

		if authedUser.MFAEnabled {
			a.Success(w, r, http.StatusOK, res, nil)
			return
		}

		cookieErr := a.createSessionCookie(w, sessionId)
		if cookieErr != nil {
			a.Failure(w, r, http.StatusInternalServerError, Errorf(EINVALID, "INVALID_COOKIE"))
			return
		}

		a.Success(w, r, http.StatusOK, res, nil)
	}
}

// handleHeaderLogin authenticates the user by looking at configurable headers
// and then creates the user if one does not exist and logs them in
// @Summary Login Header
// @Description attempts to log the user in with provided credentials
// @Description *Endpoint only available when Header auth is enabled
// @Tags auth
// @Produce json
// @Success 200 object standardJsonResponse{data=loginResponse}
// @Failure 401 object standardJsonResponse{}
// @Failure 500 object standardJsonResponse{}
// @Router /auth [get]
func (a *Service) handleHeaderLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		viper.GetString("auth.ldap.url")

		username := r.Header.Get(viper.GetString("auth.header.usernameHeader"))
		useremail := r.Header.Get(viper.GetString("auth.header.emailHeader"))

		if username == "" {
			a.Failure(w, r, http.StatusUnauthorized, Errorf(EUNAUTHORIZED, "MISSING_AUTH_HEADER"))
			return
		}

		authedUser, sessionId, err := a.authAndCreateUserHeader(r.Context(), username, useremail)
		if err != nil {
			a.Failure(w, r, http.StatusUnauthorized, Errorf(EINVALID, "INVALID_LOGIN"))
			return
		}

		res := loginResponse{
			User:        authedUser,
			SessionId:   sessionId,
			MFARequired: authedUser.MFAEnabled,
		}

		if authedUser.MFAEnabled {
			a.Success(w, r, http.StatusOK, res, nil)
			return
		}

		cookieErr := a.createSessionCookie(w, sessionId)
		if cookieErr != nil {
			a.Failure(w, r, http.StatusInternalServerError, Errorf(EINVALID, "INVALID_COOKIE"))
			return
		}

		a.Success(w, r, http.StatusOK, res, nil)
	}
}

type mfaLoginRequestBody struct {
	Passcode  string `json:"passcode" validate:"required"`
	SessionId string `json:"sessionId" validate:"required"`
}

// handleMFALogin attempts to log in the user with MFA token
// @Summary MFA Login
// @Description attempts to log the user in with provided MFA token
// @Tags auth
// @Produce  json
// @Param credentials body mfaLoginRequestBody false "mfa login object"
// @Success 200 object standardJsonResponse{}
// @Failure 401 object standardJsonResponse{}
// @Failure 500 object standardJsonResponse{}
// @Router /auth/mfa [post]
func (a *Service) handleMFALogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, bodyErr := io.ReadAll(r.Body)
		if bodyErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
			return
		}

		var u = mfaLoginRequestBody{}
		jsonErr := json.Unmarshal(body, &u)
		if jsonErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		inputErr := validate.Struct(u)
		if inputErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, inputErr.Error()))
			return
		}

		err := a.AuthService.MFATokenValidate(r.Context(), u.SessionId, u.Passcode)
		if err != nil {
			a.Failure(w, r, http.StatusUnauthorized, Errorf(EINVALID, "INVALID_AUTHENTICATOR_TOKEN"))
			return
		}

		cookieErr := a.createSessionCookie(w, u.SessionId)
		if cookieErr != nil {
			a.Failure(w, r, http.StatusInternalServerError, Errorf(EINVALID, "INVALID_COOKIE"))
			return
		}

		a.Success(w, r, http.StatusOK, nil, nil)
	}
}

// handleLogout clears the user Cookie(s) ending session
// @Summary Logout
// @Description Logs the user out by deleting session cookies
// @Tags auth
// @Success 200
// @Router /auth/logout [delete]
func (a *Service) handleLogout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		SessionId, cookieErr := a.validateSessionCookie(w, r)
		if cookieErr != nil {
			a.Failure(w, r, http.StatusUnauthorized, Errorf(EINVALID, "INVALID_USER"))
			return
		}

		err := a.AuthService.DeleteSession(r.Context(), SessionId)
		if err != nil {
			a.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		a.clearUserCookies(w)
		a.Success(w, r, http.StatusOK, nil, nil)
	}
}

type guestUserCreateRequestBody struct {
	Name string `json:"name" validate:"required"`
}

// handleCreateGuestUser registers a user as a guest user
// @Summary Create Guest User
// @Description Registers a user as a guest (non-authenticated)
// @Tags auth
// @Produce json
// @Param user body guestUserCreateRequestBody false "guest user object"
// @Success 200 object standardJsonResponse{data=thunderdome.User}
// @Failure 400 object standardJsonResponse{}
// @Failure 500 object standardJsonResponse{}
// @Router /auth/guest [post]
func (a *Service) handleCreateGuestUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		AllowGuests := viper.GetBool("config.allow_guests")
		if !AllowGuests {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, "GUESTS_USERS_DISABLED"))
			return
		}

		body, bodyErr := io.ReadAll(r.Body)
		if bodyErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
			return
		}

		var u = guestUserCreateRequestBody{}
		jsonErr := json.Unmarshal(body, &u)
		if jsonErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		inputErr := validate.Struct(u)
		if inputErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, inputErr.Error()))
			return
		}

		newUser, err := a.UserService.CreateUserGuest(r.Context(), u.Name)
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

type userRegisterRequestBody struct {
	Name      string `json:"name" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password1 string `json:"password1" validate:"required,min=6,max=72"`
	Password2 string `json:"password2" validate:"required,min=6,max=72,eqfield=Password1"`
}

// handleUserRegistration registers a new authenticated user
// @Summary Create User
// @Description Registers a user (authenticated)
// @Tags auth
// @Produce json
// @Param user body userRegisterRequestBody false "new user object"
// @Success 200 object standardJsonResponse{data=thunderdome.User}
// @Failure 400 object standardJsonResponse{}
// @Failure 500 object standardJsonResponse{}
// @Router /auth/register [post]
func (a *Service) handleUserRegistration() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		AllowRegistration := viper.GetBool("config.allow_registration")
		if !AllowRegistration {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, "USER_REGISTRATION_DISABLED"))
		}

		body, bodyErr := io.ReadAll(r.Body)
		if bodyErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
			return
		}

		var u = userRegisterRequestBody{}
		jsonErr := json.Unmarshal(body, &u)
		if jsonErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		inputErr := validate.Struct(u)
		if inputErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, inputErr.Error()))
			return
		}

		ActiveUserID, _ := a.validateUserCookie(w, r)

		UserName, UserEmail, UserPassword, accountErr := validateUserAccountWithPasswords(
			u.Name,
			u.Email,
			u.Password1,
			u.Password2,
		)

		if accountErr != nil {
			a.Failure(w, r, http.StatusBadRequest, accountErr)
			return
		}

		newUser, VerifyID, err := a.UserService.CreateUserRegistered(r.Context(), UserName, UserEmail, UserPassword, ActiveUserID)
		if err != nil {
			a.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		a.Email.SendWelcome(UserName, UserEmail, VerifyID)

		if ActiveUserID != "" {
			a.clearUserCookies(w)
		}

		SessionID, err := a.AuthService.CreateSession(r.Context(), newUser.Id)
		if err != nil {
			a.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		cookieErr := a.createSessionCookie(w, SessionID)
		if cookieErr != nil {
			a.Failure(w, r, http.StatusInternalServerError, Errorf(EINVALID, "INVALID_COOKIE"))
			return
		}

		a.Success(w, r, http.StatusOK, newUser, nil)
	}
}

type forgotPasswordRequestBody struct {
	Email string `json:"email" validate:"required,email"`
}

// handleForgotPassword attempts to send a password reset Email
// @Summary Forgot Password
// @Description Sends a forgot password reset Email to user
// @Tags auth
// @Produce json
// @Param user body forgotPasswordRequestBody false "forgot password object"
// @Success 200 object standardJsonResponse{}
// @Router /auth/forgot-password [post]
func (a *Service) handleForgotPassword() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, bodyErr := io.ReadAll(r.Body)
		if bodyErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
			return
		}

		var u = forgotPasswordRequestBody{}
		jsonErr := json.Unmarshal(body, &u)
		if jsonErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		inputErr := validate.Struct(u)
		if inputErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, inputErr.Error()))
			return
		}

		UserEmail := strings.ToLower(u.Email)

		ResetID, UserName, resetErr := a.AuthService.UserResetRequest(r.Context(), UserEmail)
		if resetErr == nil {
			a.Email.SendForgotPassword(UserName, UserEmail, ResetID)
		}

		a.Success(w, r, http.StatusOK, nil, nil)
	}
}

type resetPasswordRequestBody struct {
	ResetID   string `json:"resetId" validate:"required"`
	Password1 string `json:"password1" validate:"required,min=6,max=72"`
	Password2 string `json:"password2" validate:"required,min=6,max=72,eqfield=Password1"`
}

// handleResetPassword attempts to reset a user's password
// @Summary Reset Password
// @Description Resets the user's password
// @Tags auth
// @Produce json
// @Param reset body resetPasswordRequestBody false "reset password object"
// @Success 200 object standardJsonResponse{}
// @Success 400 object standardJsonResponse{}
// @Success 500 object standardJsonResponse{}
// @Router /auth/reset-password [patch]
func (a *Service) handleResetPassword() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, bodyErr := io.ReadAll(r.Body)
		if bodyErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
			return
		}

		var u = resetPasswordRequestBody{}
		jsonErr := json.Unmarshal(body, &u)
		if jsonErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		inputErr := validate.Struct(u)
		if inputErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, inputErr.Error()))
			return
		}

		UserName, UserEmail, resetErr := a.AuthService.UserResetPassword(r.Context(), u.ResetID, u.Password1)
		if resetErr != nil {
			a.Failure(w, r, http.StatusInternalServerError, resetErr)
			return
		}

		a.Email.SendPasswordReset(UserName, UserEmail)

		a.Success(w, r, http.StatusOK, nil, nil)
	}
}

type updatePasswordRequestBody struct {
	Password1 string `json:"password1" validate:"required,min=6,max=72"`
	Password2 string `json:"password2" validate:"required,min=6,max=72,eqfield=Password1"`
}

// handleUpdatePassword attempts to update a user's password
// @Summary Update Password
// @Description Updates the user's password
// @Tags auth
// @Produce json
// @Param passwords body updatePasswordRequestBody false "update password object"
// @Success 200 object standardJsonResponse{}
// @Success 400 object standardJsonResponse{}
// @Success 500 object standardJsonResponse{}
// @Security ApiKeyAuth
// @Router /auth/update-password [patch]
func (a *Service) handleUpdatePassword() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		UserID := r.Context().Value(contextKeyUserID).(string)
		body, bodyErr := io.ReadAll(r.Body)
		if bodyErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
			return
		}

		var u = updatePasswordRequestBody{}
		jsonErr := json.Unmarshal(body, &u)
		if jsonErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		inputErr := validate.Struct(u)
		if inputErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, inputErr.Error()))
			return
		}

		UserName, UserEmail, updateErr := a.AuthService.UserUpdatePassword(r.Context(), UserID, u.Password1)
		if updateErr != nil {
			a.Failure(w, r, http.StatusInternalServerError, updateErr)
			return
		}

		a.Email.SendPasswordUpdate(UserName, UserEmail)

		a.Success(w, r, http.StatusOK, nil, nil)
	}
}

type verificationRequestBody struct {
	VerifyID string `json:"verifyId" validate:"required"`
}

// handleAccountVerification attempts to verify a users account
// @Summary Verify User
// @Description Updates the users verified Email status
// @Tags auth
// @Produce json
// @Param verify body verificationRequestBody false "verify object"
// @Success 200 object standardJsonResponse{}
// @Success 500 object standardJsonResponse{}
// @Router /auth/verify [patch]
func (a *Service) handleAccountVerification() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, bodyErr := io.ReadAll(r.Body)
		if bodyErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
			return
		}

		var u = verificationRequestBody{}
		jsonErr := json.Unmarshal(body, &u)
		if jsonErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		inputErr := validate.Struct(u)
		if inputErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, inputErr.Error()))
			return
		}

		verifyErr := a.AuthService.VerifyUserAccount(r.Context(), u.VerifyID)
		if verifyErr != nil {
			a.Failure(w, r, http.StatusInternalServerError, verifyErr)
			return
		}

		a.Success(w, r, http.StatusOK, nil, nil)
	}
}

// handleMFASetupGenerate generates the MFA secret and QR code for setup
// @Summary MFA Setup Generate secret and QR code
// @Description Generates MFA secret and QR Code
// @Tags auth
// @Success 200
// @Router /auth/mfa/setup/generate [post]
func (a *Service) handleMFASetupGenerate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		UserID := ctx.Value(contextKeyUserID).(string)

		u, err := a.UserService.GetUser(ctx, UserID)
		if err != nil {
			a.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		secret, png64, err := a.AuthService.MFASetupGenerate(u.Email)
		if err != nil {
			a.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		type result struct {
			Secret string `json:"secret"`
			QRCode string `json:"qrCode"`
		}

		a.Success(w, r, http.StatusOK, result{Secret: secret, QRCode: png64}, nil)
	}
}

type mfaSetupValidateRequestBody struct {
	Secret   string `json:"secret" validate:"required"`
	Passcode string `json:"passcode" validate:"required"`
}

// handleMFASetupValidate validates the passcode for MFA secret during setup
// @Summary Validate MFA Setup passcode
// @Description Validates the passcode for the MFA secret
// @Param verify body mfaSetupValidateRequestBody false "verify object"
// @Tags auth
// @Success 200
// @Router /auth/mfa/setup/validate [post]
func (a *Service) handleMFASetupValidate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		UserID := ctx.Value(contextKeyUserID).(string)

		body, bodyErr := io.ReadAll(r.Body)
		if bodyErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
			return
		}

		var v = mfaSetupValidateRequestBody{}
		jsonErr := json.Unmarshal(body, &v)
		if jsonErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		inputErr := validate.Struct(v)
		if inputErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, inputErr.Error()))
			return
		}

		type result struct {
			Result string `json:"result"`
		}
		res := result{Result: "SUCCESS"}

		err := a.AuthService.MFASetupValidate(ctx, UserID, v.Secret, v.Passcode)
		if err != nil {
			res.Result = err.Error()
		}

		a.Success(w, r, http.StatusOK, res, nil)
	}
}

// handleMFARemove removes MFA requirement from user auth
// @Summary Remove MFA
// @Description Removes MFA requirement from user auth
// @Tags auth
// @Success 200
// @Router /auth/mfa [delete]
func (a *Service) handleMFARemove() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		UserID := ctx.Value(contextKeyUserID).(string)

		err := a.AuthService.MFARemove(ctx, UserID)
		if err != nil {
			a.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		a.Success(w, r, http.StatusOK, nil, nil)
	}
}
