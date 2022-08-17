package api

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/StevenWeathers/thunderdome-planning-poker/model"
	"github.com/go-ldap/ldap/v3"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type userAccount struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

type userPassword struct {
	Password1 string `json:"password1" validate:"required,min=6,max=72"`
	Password2 string `json:"password2" validate:"required,min=6,max=72,eqfield=Password1"`
}

// validateUserAccount makes sure user's name, email are valid before creating the account
func validateUserAccount(name string, email string) (UserName string, UserEmail string, validateErr error) {
	a := userAccount{
		Name:  name,
		Email: email,
	}
	aErr := validate.Struct(a)
	if aErr != nil {
		return "", "", aErr
	}

	return name, email, nil
}

// validateUserAccountWithPasswords makes sure user's name, email, and password are valid before creating the account
func validateUserAccountWithPasswords(name string, email string, pwd1 string, pwd2 string) (UserName string, UserEmail string, UpdatedPassword string, validateErr error) {
	a := userAccount{
		Name:  name,
		Email: email,
	}
	p := userPassword{
		Password1: pwd1,
		Password2: pwd2,
	}
	aErr := validate.Struct(a)
	if aErr != nil {
		return "", "", "", aErr
	}
	pErr := validate.Struct(p)
	if pErr != nil {
		return "", "", "", pErr
	}

	return name, email, pwd1, nil
}

// validateUserPassword makes sure user password is valid before updating the password
func validateUserPassword(pwd1 string, pwd2 string) (UpdatedPassword string, validateErr error) {
	a := userPassword{
		Password1: pwd1,
		Password2: pwd2,
	}
	err := validate.Struct(a)

	return pwd1, err
}

// createUserCookie creates the users cookie
func (a *api) createUserCookie(w http.ResponseWriter, UserID string) error {
	encoded, err := a.cookie.Encode(a.config.SecureCookieName, UserID)
	if err != nil {
		return err

	}

	cookie := &http.Cookie{
		Name:     a.config.SecureCookieName,
		Value:    encoded,
		Path:     a.config.PathPrefix + "/",
		HttpOnly: true,
		Domain:   a.config.AppDomain,
		MaxAge:   86400 * 365,
		Secure:   a.config.SecureCookieFlag,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, cookie)

	return nil
}

// createSessionCookie creates the user's session cookie
func (a *api) createSessionCookie(w http.ResponseWriter, SessionID string) error {
	encoded, err := a.cookie.Encode(a.config.SessionCookieName, SessionID)
	if err != nil {
		return err
	}

	cookie := &http.Cookie{
		Name:     a.config.SessionCookieName,
		Value:    encoded,
		Path:     a.config.PathPrefix + "/",
		HttpOnly: true,
		Domain:   a.config.AppDomain,
		MaxAge:   86400 * 30,
		Secure:   a.config.SecureCookieFlag,
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie(w, cookie)

	return nil
}

// clearUserCookies wipes the frontend and backend cookies
// used in the event of bad cookie reads
func (a *api) clearUserCookies(w http.ResponseWriter) {
	feCookie := &http.Cookie{
		Name:   a.config.FrontendCookieName,
		Value:  "",
		Path:   a.config.PathPrefix + "/",
		MaxAge: -1,
	}
	beCookie := &http.Cookie{
		Name:     a.config.SecureCookieName,
		Value:    "",
		Path:     a.config.PathPrefix + "/",
		Domain:   a.config.AppDomain,
		Secure:   a.config.SecureCookieFlag,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1,
		HttpOnly: true,
	}
	seCookie := &http.Cookie{
		Name:     a.config.SessionCookieName,
		Value:    "",
		Path:     a.config.PathPrefix + "/",
		Domain:   a.config.AppDomain,
		Secure:   a.config.SecureCookieFlag,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1,
		HttpOnly: true,
	}

	http.SetCookie(w, feCookie)
	http.SetCookie(w, beCookie)
	http.SetCookie(w, seCookie)
}

// validateUserCookie returns the UserID from secure cookies or errors if failures getting it
func (a *api) validateUserCookie(w http.ResponseWriter, r *http.Request) (string, error) {
	var UserID string

	if cookie, err := r.Cookie(a.config.SecureCookieName); err == nil {
		var value string
		if err = a.cookie.Decode(a.config.SecureCookieName, cookie.Value, &value); err == nil {
			UserID = value
		} else {
			a.clearUserCookies(w)
			return "", errors.New("INVALID_USER_COOKIE")
		}
	} else {
		return "", errors.New("NO_USER_COOKIE")
	}

	return UserID, nil
}

// validateSessionCookie returns the SessionID from secure cookies or errors if failures getting it
func (a *api) validateSessionCookie(w http.ResponseWriter, r *http.Request) (string, error) {
	var SessionID string

	if cookie, err := r.Cookie(a.config.SessionCookieName); err == nil {
		var value string
		if err = a.cookie.Decode(a.config.SessionCookieName, cookie.Value, &value); err == nil {
			SessionID = value
		} else {
			a.clearUserCookies(w)
			return "", errors.New("INVALID_SESSION_COOKIE")
		}
	} else {
		return "", errors.New("NO_SESSION_COOKIE")
	}

	return SessionID, nil
}

// Success returns the successful response including any data and meta
func (a *api) Success(w http.ResponseWriter, r *http.Request, code int, data interface{}, meta interface{}) {
	result := &standardJsonResponse{
		Success: true,
		Error:   "",
		Data:    map[string]interface{}{},
		Meta:    map[string]interface{}{},
	}

	if meta != nil {
		result.Meta = meta
	}

	if data != nil {
		result.Data = data
	}

	response, _ := json.Marshal(result)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// Failure responds with an error and its associated status code header
func (a *api) Failure(w http.ResponseWriter, r *http.Request, code int, err error) {
	ctx := r.Context()
	// Extract error message.
	errCode, errMessage := ErrorCode(err), ErrorMessage(err)

	if errCode == EINTERNAL {
		a.logger.Ctx(ctx).Error(
			"[http] error",
			zap.String("method", r.Method),
			zap.String("url_path", sanitizeUserInputForLogs(r.URL.Path)),
			zap.Error(err),
		)
	}

	result := &standardJsonResponse{
		Success: false,
		Error:   errMessage,
		Data:    map[string]interface{}{},
		Meta:    map[string]interface{}{},
	}

	response, _ := json.Marshal(result)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// getLimitOffsetFromRequest gets the limit and offset query parameters from the request
// defaulting to 20 for limit and 0 for offset
func getLimitOffsetFromRequest(r *http.Request) (limit int, offset int) {
	defaultLimit := 20
	defaultOffset := 0
	query := r.URL.Query()
	Limit, limitErr := strconv.Atoi(query.Get("limit"))
	if limitErr != nil || Limit == 0 {
		Limit = defaultLimit
	}

	Offset, offsetErr := strconv.Atoi(query.Get("offset"))
	if offsetErr != nil {
		Offset = defaultOffset
	}

	return Limit, Offset
}

// getSearchFromRequest gets the search query parameter from the request
func getSearchFromRequest(r *http.Request) (search string, err error) {
	query := r.URL.Query()
	Search := query.Get("search")
	searchErr := validate.Var(Search, "required,min=3")
	if searchErr != nil {
		return "", searchErr
	}

	return Search, nil
}

// for logging purposes sanitize strings by removing new lines
func sanitizeUserInputForLogs(unescapedInput string) string {
	escapedString := strings.Replace(unescapedInput, "\n", "", -1)
	escapedString = strings.Replace(escapedString, "\r", "", -1)
	return escapedString
}

// Authenticate using LDAP and if user does not exist, automatically add user as a verified user
func (a *api) authAndCreateUserLdap(ctx context.Context, UserName string, UserPassword string) (*model.User, string, error) {
	var AuthedUser *model.User
	var SessionId string
	var sessErr error

	l, err := ldap.DialURL(viper.GetString("auth.ldap.url"))
	if err != nil {
		a.logger.Ctx(ctx).Error("Failed connecting to ldap server at " + viper.GetString("auth.ldap.url"))
		return AuthedUser, SessionId, err
	}
	defer l.Close()
	if viper.GetBool("auth.ldap.use_tls") {
		err = l.StartTLS(&tls.Config{InsecureSkipVerify: true})
		if err != nil {
			a.logger.Ctx(ctx).Error("Failed securing ldap connection", zap.Error(err))
			return AuthedUser, SessionId, err
		}
	}

	if viper.GetString("auth.ldap.bindname") != "" {
		err = l.Bind(viper.GetString("auth.ldap.bindname"), viper.GetString("auth.ldap.bindpass"))
		if err != nil {
			a.logger.Ctx(ctx).Error("Failed binding for authentication", zap.Error(err))
			return AuthedUser, SessionId, err
		}
	}

	searchRequest := ldap.NewSearchRequest(viper.GetString("auth.ldap.basedn"),
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf(viper.GetString("auth.ldap.filter"), ldap.EscapeFilter(UserName)),
		[]string{"dn", viper.GetString("auth.ldap.mail_attr"), viper.GetString("auth.ldap.cn_attr")},
		nil,
	)

	sr, err := l.Search(searchRequest)
	if err != nil {
		a.logger.Ctx(ctx).Error("Failed performing ldap search query", zap.String("username", sanitizeUserInputForLogs(UserName)), zap.Error(err))
		return AuthedUser, SessionId, err
	}

	if len(sr.Entries) != 1 {
		a.logger.Ctx(ctx).Error("User does not exist or too many entries returned", zap.String("username", sanitizeUserInputForLogs(UserName)))
		return AuthedUser, SessionId, errors.New("user not found")
	}

	userdn := sr.Entries[0].DN
	useremail := sr.Entries[0].GetAttributeValue(viper.GetString("auth.ldap.mail_attr"))
	usercn := sr.Entries[0].GetAttributeValue(viper.GetString("auth.ldap.cn_attr"))

	err = l.Bind(userdn, UserPassword)
	if err != nil {
		a.logger.Ctx(ctx).Error("Failed authenticating user", zap.String("username", sanitizeUserInputForLogs(UserName)))
		return AuthedUser, SessionId, err
	}

	AuthedUser, err = a.db.GetUserByEmail(ctx, useremail)

	if AuthedUser == nil {
		a.logger.Ctx(ctx).Error("User does not exist in database, auto-recruit", zap.String("useremail", sanitizeUserInputForLogs(useremail)))
		newUser, verifyID, sessionId, err := a.db.CreateUserRegistered(ctx, usercn, useremail, "", "")
		if err != nil {
			a.logger.Ctx(ctx).Error("Failed auto-creating new user", zap.Error(err))
			return AuthedUser, SessionId, err
		}
		err = a.db.VerifyUserAccount(ctx, verifyID)
		if err != nil {
			a.logger.Ctx(ctx).Error("Failed verifying new user", zap.Error(err))
			return AuthedUser, SessionId, err
		}
		AuthedUser = newUser
		SessionId = sessionId
	} else {
		if AuthedUser.Disabled {
			return nil, "", fmt.Errorf("user is disabled")
		}

		SessionId, sessErr = a.db.CreateSession(ctx, AuthedUser.Id)
		if sessErr != nil {
			a.logger.Ctx(ctx).Error("Failed creating user session", zap.Error(err))
			return nil, "", err
		}
	}

	return AuthedUser, SessionId, nil
}

// isTeamUserOrAnAdmin determines if the request user is a team user
// or team admin, or department admin (if applicable), or organization admin (if applicable), or application admin
func isTeamUserOrAnAdmin(r *http.Request) bool {
	UserType := r.Context().Value(contextKeyUserType).(string)
	OrgRole := r.Context().Value(contextKeyOrgRole)
	DepartmentRole := r.Context().Value(contextKeyDepartmentRole)
	TeamRole := r.Context().Value(contextKeyTeamRole).(string)
	var isAdmin = UserType == adminUserType
	if DepartmentRole != nil && DepartmentRole.(string) == adminUserType {
		isAdmin = true
	}
	if OrgRole != nil && OrgRole.(string) == adminUserType {
		isAdmin = true
	}

	return isAdmin || TeamRole != ""
}
