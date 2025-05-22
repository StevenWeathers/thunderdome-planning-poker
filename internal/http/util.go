package http

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"

	"github.com/go-ldap/ldap/v3"
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
func validateUserAccount(name string, email string) (userName string, userEmail string, validateErr error) {
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
func validateUserAccountWithPasswords(name string, email string, pwd1 string, pwd2 string) (userName string, userEmail string, updatedPassword string, validateErr error) {
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
func validateUserPassword(pwd1 string, pwd2 string) (updatedPassword string, validateErr error) {
	a := userPassword{
		Password1: pwd1,
		Password2: pwd2,
	}
	err := validate.Struct(a)

	return pwd1, err
}

// Success returns the successful response including any data and meta
func (s *Service) Success(w http.ResponseWriter, r *http.Request, code int, data interface{}, meta interface{}) {
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
func (s *Service) Failure(w http.ResponseWriter, r *http.Request, code int, err error) {
	// Extract error message.
	_, errMessage := ErrorCode(err), ErrorMessage(err)

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
	if limitErr != nil || Limit <= 0 {
		Limit = defaultLimit
	}

	Offset, offsetErr := strconv.Atoi(query.Get("offset"))
	if offsetErr != nil || Offset < 0 {
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
func (s *Service) authAndCreateUserLdap(ctx context.Context, userName string, userPassword string) (*thunderdome.User, string, error) {
	var authedUser *thunderdome.User
	var sessionID string
	var sessErr error

	l, err := ldap.DialURL(s.Config.AuthLdapUrl)
	if err != nil {
		s.Logger.Ctx(ctx).Error("Failed connecting to ldap server at " + s.Config.AuthLdapUrl)
		return authedUser, sessionID, err
	}
	defer l.Close()
	if s.Config.AuthLdapUseTls {
		err = l.StartTLS(&tls.Config{InsecureSkipVerify: true})
		if err != nil {
			s.Logger.Ctx(ctx).Error("Failed securing ldap connection", zap.Error(err))
			return authedUser, sessionID, err
		}
	}

	if s.Config.AuthLdapBindname != "" {
		err = l.Bind(s.Config.AuthLdapBindname, s.Config.AuthLdapBindpass)
		if err != nil {
			s.Logger.Ctx(ctx).Error("Failed binding for authentication", zap.Error(err))
			return authedUser, sessionID, err
		}
	}

	searchRequest := ldap.NewSearchRequest(s.Config.AuthLdapBasedn,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf(s.Config.AuthLdapFilter, ldap.EscapeFilter(userName)),
		[]string{"dn", s.Config.AuthLdapMailAttr, s.Config.AuthLdapCnAttr},
		nil,
	)

	sr, err := l.Search(searchRequest)
	if err != nil {
		s.Logger.Ctx(ctx).Error("Failed performing ldap search query", zap.String("username", sanitizeUserInputForLogs(userName)), zap.Error(err))
		return authedUser, sessionID, err
	}

	if len(sr.Entries) != 1 {
		s.Logger.Ctx(ctx).Error("User does not exist or too many entries returned", zap.String("username", sanitizeUserInputForLogs(userName)))
		return authedUser, sessionID, errors.New("user not found")
	}

	userdn := sr.Entries[0].DN
	useremail := sr.Entries[0].GetAttributeValue(s.Config.AuthLdapMailAttr)
	usercn := sr.Entries[0].GetAttributeValue(s.Config.AuthLdapCnAttr)

	err = l.Bind(userdn, userPassword)
	if err != nil {
		s.Logger.Ctx(ctx).Error("Failed authenticating user", zap.String("username", sanitizeUserInputForLogs(userName)))
		return authedUser, sessionID, err
	}

	authedUser, err = s.UserDataSvc.GetUserByEmail(ctx, useremail)

	if authedUser == nil {
		var verifyID string
		s.Logger.Ctx(ctx).Error("User does not exist in database, auto-recruit", zap.String("useremail", sanitizeUserInputForLogs(useremail)))
		authedUser, verifyID, err = s.UserDataSvc.CreateUserRegistered(ctx, usercn, useremail, "", "")
		if err != nil {
			s.Logger.Ctx(ctx).Error("Failed auto-creating new user", zap.Error(err))
			return authedUser, sessionID, err
		}
		err = s.AuthDataSvc.VerifyUserAccount(ctx, verifyID)
		if err != nil {
			s.Logger.Ctx(ctx).Error("Failed verifying new user", zap.Error(err))
			return authedUser, sessionID, err
		}
		sessionID, err = s.AuthDataSvc.CreateSession(ctx, authedUser.ID, true)
		if err != nil {
			s.Logger.Ctx(ctx).Error("Failed creating user session", zap.Error(err))
			return authedUser, sessionID, err
		}
	} else {
		if authedUser.Disabled {
			return nil, "", fmt.Errorf("user is disabled")
		}

		sessionID, sessErr = s.AuthDataSvc.CreateSession(ctx, authedUser.ID, true)
		if sessErr != nil {
			s.Logger.Ctx(ctx).Error("Failed creating user session", zap.Error(err))
			return nil, "", err
		}
	}

	return authedUser, sessionID, nil
}

// Authenticate using HTTP headers and if user does not exist, automatically add user as a verified user
func (s *Service) authAndCreateUserHeader(ctx context.Context, username string, useremail string) (*thunderdome.User, string, error) {
	var authedUser *thunderdome.User
	var sessionId string
	var sessErr error

	authedUser, err := s.UserDataSvc.GetUserByEmail(ctx, useremail)

	if authedUser == nil {
		s.Logger.Ctx(ctx).Error("User does not exist in database, auto-recruit", zap.String("useremail", sanitizeUserInputForLogs(useremail)))
		var verifyID string
		authedUser, verifyID, err = s.UserDataSvc.CreateUserRegistered(ctx, username, useremail, "", "")
		if err != nil {
			s.Logger.Ctx(ctx).Error("Failed auto-creating new user", zap.Error(err))
			return authedUser, sessionId, err
		}
		err = s.AuthDataSvc.VerifyUserAccount(ctx, verifyID)
		if err != nil {
			s.Logger.Ctx(ctx).Error("Failed verifying new user", zap.Error(err))
			return authedUser, sessionId, err
		}
		sessionId, err = s.AuthDataSvc.CreateSession(ctx, authedUser.ID, true)
		if err != nil {
			s.Logger.Ctx(ctx).Error("Failed creating user session", zap.Error(err))
			return authedUser, sessionId, err
		}
	} else {
		if authedUser.Disabled {
			return nil, "", fmt.Errorf("user is disabled")
		}

		sessionId, sessErr = s.AuthDataSvc.CreateSession(ctx, authedUser.ID, true)
		if sessErr != nil {
			s.Logger.Ctx(ctx).Error("Failed creating user session", zap.Error(err))
			return nil, "", err
		}
	}

	return authedUser, sessionId, nil
}

// isTeamUserOrAnAdmin determines if the request user is a team user
// or team admin, or department admin (if applicable), or organization admin (if applicable), or application admin
func isTeamUserOrAnAdmin(r *http.Request) bool {
	ctx := r.Context()
	teamUserRoles := ctx.Value(contextKeyUserTeamRoles).(*thunderdome.UserTeamRoleInfo)
	var emptyRole = ""
	orgRole := teamUserRoles.OrganizationRole
	if orgRole == nil {
		orgRole = &emptyRole
	}
	departmentRole := teamUserRoles.DepartmentRole
	if departmentRole == nil {
		departmentRole = &emptyRole
	}
	teamRole := teamUserRoles.TeamRole
	if teamRole == nil {
		teamRole = &emptyRole
	}

	userType := ctx.Value(contextKeyUserType).(string)
	var isAdmin = userType == thunderdome.AdminUserType
	if departmentRole != nil && *departmentRole == thunderdome.AdminUserType {
		isAdmin = true
	}
	if orgRole != nil && *orgRole == thunderdome.AdminUserType {
		isAdmin = true
	}

	return isAdmin || *teamRole != ""
}

// get the index template from embedded filesystem
func (s *Service) getIndexTemplate(filesystem fs.FS) *template.Template {
	ctx := context.Background()
	// get the html template from dist, have it ready for requests
	tmplContent, ioErr := fs.ReadFile(filesystem, "index.html")
	if ioErr != nil {
		s.Logger.Ctx(ctx).Error("Error opening index template")
		if !s.Config.EmbedUseOS {
			s.Logger.Ctx(ctx).Fatal(ioErr.Error())
		}
	}

	tmplString := string(tmplContent)
	tmpl, tmplErr := template.New("index").Parse(tmplString)
	if tmplErr != nil {
		s.Logger.Ctx(ctx).Error("Error parsing index template")
		if !s.Config.EmbedUseOS {
			s.Logger.Ctx(ctx).Fatal(tmplErr.Error())
		}
	}

	return tmpl
}

func getWebsocketConnectSrc(secureProtocol bool, websocketSubdomain string, appDomain string) string {
	wcs := "wss://"
	if !secureProtocol {
		wcs = "ws://"
	}
	sub := websocketSubdomain
	if sub != "" {
		sub = fmt.Sprintf("%s.", sub)
	}

	return fmt.Sprintf("%s%s%s", wcs, sub, appDomain)
}

func retroTemplateBuildFormatFromRequest(requestFormat retroTemplateFormatRequestBody) *thunderdome.RetroTemplateFormat {
	tf := &thunderdome.RetroTemplateFormat{
		Columns: make([]thunderdome.RetroTemplateFormatColumn, 0),
	}

	for _, col := range requestFormat.Columns {
		tf.Columns = append(tf.Columns, thunderdome.RetroTemplateFormatColumn{
			Name:  col.Name,
			Label: col.Label,
			Color: col.Color,
			Icon:  col.Icon,
		})
	}

	return tf
}

// containsLink checks if the input string contains a link
func containsLink(input string) bool {
	urlPattern := `((http|https):\/\/[a-zA-Z0-9\-._~:/?#[\]@!$&'()*+,;=]+)`

	re := regexp.MustCompile(urlPattern)

	return re.MatchString(input)
}

// ssoEnabled checks if SSO is enabled based on the configuration
// It checks if LDAP, Header Auth, or OIDC Auth is enabled
// and returns true if any of them are enabled, otherwise false
// This is used to determine if SSO is enabled for the application
func (s *Service) ssoEnabled() bool {
	if s.Config.LdapEnabled {
		return true
	} else if s.Config.HeaderAuthEnabled {
		return true
	} else if s.Config.OIDCAuth.Enabled {
		return true
	} else {
		return false
	}
}
