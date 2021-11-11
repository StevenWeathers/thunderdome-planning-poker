package api

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/StevenWeathers/thunderdome-planning-poker/model"
	"github.com/go-ldap/ldap/v3"
	"github.com/spf13/viper"
	"gopkg.in/go-playground/validator.v9"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type UserAccount struct {
	Name      string `json:"name" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password1 string `json:"password1" validate:"required,min=6,max=72"`
	Password2 string `json:"password2" validate:"required,min=6,max=72,eqfield=Password1"`
}

type UserPassword struct {
	Password1 string `json:"password1" validate:"required,min=6,max=72"`
	Password2 string `json:"password2" validate:"required,min=6,max=72,eqfield=Password1"`
}

// validateUserAccount makes sure user name, email, and password are valid before creating the account
func validateUserAccount(name string, email string, pwd1 string, pwd2 string) (UserName string, UserEmail string, UpdatedPassword string, validateErr error) {
	v := validator.New()
	a := UserAccount{
		Name:      name,
		Email:     email,
		Password1: pwd1,
		Password2: pwd2,
	}
	err := v.Struct(a)

	return name, email, pwd1, err
}

// validateUserPassword makes sure user password is valid before updating the password
func validateUserPassword(pwd1 string, pwd2 string) (UpdatedPassword string, validateErr error) {
	v := validator.New()
	a := UserPassword{
		Password1: pwd1,
		Password2: pwd2,
	}
	err := v.Struct(a)

	return pwd1, err
}

// createUserCookie creates the users cookie
func (a *api) createUserCookie(w http.ResponseWriter, r *http.Request, isRegistered bool, UserID string) {
	var cookiedays = 365 // 356 days
	if isRegistered {
		cookiedays = 30 // 30 days
	}

	encoded, err := a.cookie.Encode(a.config.SecureCookieName, UserID)
	if err != nil {
		Failure(w, r, http.StatusInternalServerError, Errorf(EINVALID, "INVALID_COOKIE"))
		return

	}

	cookie := &http.Cookie{
		Name:     a.config.SecureCookieName,
		Value:    encoded,
		Path:     a.config.PathPrefix + "/",
		HttpOnly: true,
		Domain:   a.config.AppDomain,
		MaxAge:   86400 * cookiedays,
		Secure:   a.config.SecureCookieFlag,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, cookie)
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

	http.SetCookie(w, feCookie)
	http.SetCookie(w, beCookie)
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
			return "", errors.New("invalid user cookies")
		}
	} else {
		a.clearUserCookies(w)
		return "", errors.New("invalid user cookies")
	}

	return UserID, nil
}

func (a *api) createCookie(UserID string) *http.Cookie {
	encoded, err := a.cookie.Encode(a.config.SecureCookieName, UserID)
	var NewCookie *http.Cookie

	if err == nil {
		NewCookie = &http.Cookie{
			Name:     a.config.SecureCookieName,
			Value:    encoded,
			Path:     a.config.PathPrefix + "/",
			HttpOnly: true,
			Domain:   a.config.AppDomain,
			MaxAge:   86400 * 30, // 30 days
			Secure:   a.config.SecureCookieFlag,
			SameSite: http.SameSiteStrictMode,
		}
	}
	return NewCookie
}

// Success returns the successful response including any data and meta
func Success(w http.ResponseWriter, r *http.Request, code int, data interface{}, meta interface{}) {
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
func Failure(w http.ResponseWriter, r *http.Request, code int, err error) {
	// Extract error message.
	errCode, errMessage := ErrorCode(err), ErrorMessage(err)

	if errCode == EINTERNAL {
		LogError(r, err)
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

// LogError logs an error with the HTTP route information.
func LogError(r *http.Request, err error) {
	log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
}

// getJSONRequestBody gets a JSON request body broken into a key/value map
func getJSONRequestBody(r *http.Request, w http.ResponseWriter) map[string]interface{} {
	body, _ := ioutil.ReadAll(r.Body) // check for errors
	keyVal := make(map[string]interface{})
	jsonErr := json.Unmarshal(body, &keyVal) // check for errors

	if jsonErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		return nil
	}

	return keyVal
}

// getLimitOffsetFromRequest gets the limit and offset query parameters from the request
// defaulting to 20 for limit and 0 for offset
func getLimitOffsetFromRequest(r *http.Request, w http.ResponseWriter) (limit int, offset int) {
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

// Authenticate using LDAP and if user does not exist, automatically add user as a verified user
func (a *api) authAndCreateUserLdap(UserName string, UserPassword string) (*model.User, error) {
	var AuthedUser *model.User
	l, err := ldap.DialURL(viper.GetString("auth.ldap.url"))
	if err != nil {
		log.Println("Failed connecting to ldap server at", viper.GetString("auth.ldap.url"))
		return AuthedUser, err
	}
	defer l.Close()
	if viper.GetBool("auth.ldap.use_tls") {
		err = l.StartTLS(&tls.Config{InsecureSkipVerify: true})
		if err != nil {
			log.Println("Failed securing ldap connection", err)
			return AuthedUser, err
		}
	}

	if viper.GetString("auth.ldap.bindname") != "" {
		err = l.Bind(viper.GetString("auth.ldap.bindname"), viper.GetString("auth.ldap.bindpass"))
		if err != nil {
			log.Println("Failed binding for authentication:", err)
			return AuthedUser, err
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
		log.Println("Failed performing ldap search query for", UserName, ":", err)
		return AuthedUser, err
	}

	if len(sr.Entries) != 1 {
		log.Println("User", UserName, "does not exist or too many entries returned")
		return AuthedUser, errors.New("user not found")
	}

	userdn := sr.Entries[0].DN
	useremail := sr.Entries[0].GetAttributeValue(viper.GetString("auth.ldap.mail_attr"))
	usercn := sr.Entries[0].GetAttributeValue(viper.GetString("auth.ldap.cn_attr"))

	err = l.Bind(userdn, UserPassword)
	if err != nil {
		log.Println("Failed authenticating user ", UserName)
		return AuthedUser, err
	}

	AuthedUser, err = a.db.GetUserByEmail(useremail)
	if AuthedUser == nil {
		log.Println("User", useremail, "does not exist in database, auto-recruit")
		newUser, verifyID, err := a.db.CreateUserRegistered(usercn, useremail, "", "")
		if err != nil {
			log.Println("Failed auto-creating new user", err)
			return AuthedUser, err
		}
		err = a.db.VerifyUserAccount(verifyID)
		if err != nil {
			log.Println("Failed verifying new user", err)
			return AuthedUser, err
		}
		AuthedUser = newUser
	}

	return AuthedUser, nil
}
