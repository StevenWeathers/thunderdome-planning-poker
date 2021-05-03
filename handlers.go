package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"html/template"
	"image"
	"image/png"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/StevenWeathers/thunderdome-planning-poker/pkg/database"
	"github.com/anthonynsimon/bild/transform"
	"github.com/gorilla/mux"
	"github.com/ipsn/go-adorable"
	"github.com/o1egl/govatar"
	"github.com/spf13/viper"
	"gopkg.in/go-playground/validator.v9"
)

type contextKey string

var (
	contextKeyWarriorID contextKey
	apiKeyHeaderName    string = "X-API-Key"
)

type warriorAccount struct {
	Name      string `json:"name" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password1 string `json:"password1" validate:"required,min=6,max=72"`
	Password2 string `json:"password2" validate:"required,min=6,max=72,eqfield=Password1"`
}

type warriorPassword struct {
	Password1 string `json:"password1" validate:"required,min=6,max=72"`
	Password2 string `json:"password2" validate:"required,min=6,max=72,eqfield=Password1"`
}

// ValidateUserAccount makes sure warrior name, email, and password are valid before creating the account
func ValidateUserAccount(name string, email string, pwd1 string, pwd2 string) (WarriorName string, WarriorEmail string, WarriorPassword string, validateErr error) {
	v := validator.New()
	a := warriorAccount{
		Name:      name,
		Email:     email,
		Password1: pwd1,
		Password2: pwd2,
	}
	err := v.Struct(a)

	return name, email, pwd1, err
}

// ValidateUserPassword makes sure warrior password is valid before updating the password
func ValidateUserPassword(pwd1 string, pwd2 string) (WarriorPassword string, validateErr error) {
	v := validator.New()
	a := warriorPassword{
		Password1: pwd1,
		Password2: pwd2,
	}
	err := v.Struct(a)

	return pwd1, err
}

// RespondWithJSON takes a payload and writes the response
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// createUserCookie creates the warriors cookie
func (s *server) createUserCookie(w http.ResponseWriter, isRegistered bool, WarriorID string) {
	var cookiedays = 365 // 356 days
	if isRegistered {
		cookiedays = 30 // 30 days
	}

	encoded, err := s.cookie.Encode(s.config.SecureCookieName, WarriorID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return

	}

	cookie := &http.Cookie{
		Name:     s.config.SecureCookieName,
		Value:    encoded,
		Path:     s.config.PathPrefix + "/",
		HttpOnly: true,
		Domain:   s.config.AppDomain,
		MaxAge:   86400 * cookiedays,
		Secure:   s.config.SecureCookieFlag,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, cookie)
}

// clearUserCookies wipes the frontend and backend cookies
// used in the event of bad cookie reads
func (s *server) clearUserCookies(w http.ResponseWriter) {
	feCookie := &http.Cookie{
		Name:   s.config.FrontendCookieName,
		Value:  "",
		Path:   s.config.PathPrefix + "/",
		MaxAge: -1,
	}
	beCookie := &http.Cookie{
		Name:     s.config.SecureCookieName,
		Value:    "",
		Path:     s.config.PathPrefix + "/",
		Domain:   s.config.AppDomain,
		Secure:   s.config.SecureCookieFlag,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1,
		HttpOnly: true,
	}

	http.SetCookie(w, feCookie)
	http.SetCookie(w, beCookie)
}

// validateUserCookie returns the warriorID from secure cookies or errors if failures getting it
func (s *server) validateUserCookie(w http.ResponseWriter, r *http.Request) (string, error) {
	var warriorID string

	if cookie, err := r.Cookie(s.config.SecureCookieName); err == nil {
		var value string
		if err = s.cookie.Decode(s.config.SecureCookieName, cookie.Value, &value); err == nil {
			warriorID = value
		} else {
			log.Println("error in reading warrior cookie : " + err.Error() + "\n")
			s.clearUserCookies(w)
			return "", errors.New("invalid warrior cookies")
		}
	} else {
		log.Println("error in reading warrior cookie : " + err.Error() + "\n")
		s.clearUserCookies(w)
		return "", errors.New("invalid warrior cookies")
	}

	return warriorID, nil
}

/*
	Middlewares
*/

// adminOnly middleware checks if the user is an admin, otherwise reject their request
func (s *server) adminOnly(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get(apiKeyHeaderName)
		apiKey = strings.TrimSpace(apiKey)
		var warriorID string

		if apiKey != "" {
			var apiKeyErr error
			warriorID, apiKeyErr = s.database.ValidateAPIKey(apiKey)
			if apiKeyErr != nil {
				log.Println("error validating api key : " + apiKeyErr.Error() + "\n")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		} else {
			var cookieErr error
			warriorID, cookieErr = s.validateUserCookie(w, r)
			if cookieErr != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		}

		adminErr := s.database.ConfirmAdmin(warriorID)
		if adminErr != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), contextKeyWarriorID, warriorID)

		h(w, r.WithContext(ctx))
	}
}

// userOnly validates that the request was made by a valid warrior
func (s *server) userOnly(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get(apiKeyHeaderName)
		apiKey = strings.TrimSpace(apiKey)
		var warriorID string

		if apiKey != "" {
			var apiKeyErr error
			warriorID, apiKeyErr = s.database.ValidateAPIKey(apiKey)
			if apiKeyErr != nil {
				log.Println("error validating api key : " + apiKeyErr.Error() + "\n")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		} else {
			var cookieErr error
			warriorID, cookieErr = s.validateUserCookie(w, r)
			if cookieErr != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		}

		_, warErr := s.database.GetUser(warriorID)
		if warErr != nil {
			log.Println("error finding warrior : " + warErr.Error() + "\n")
			s.clearUserCookies(w)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), contextKeyWarriorID, warriorID)

		h(w, r.WithContext(ctx))
	}
}

/*
	Handlers
*/

// handleIndex parses the index html file, injecting any relevant data
func (s *server) handleIndex() http.HandlerFunc {
	type AppConfig struct {
		AllowedPointValues    []string
		DefaultPointValues    []string
		ShowWarriorRank       bool
		AvatarService         string
		ToastTimeout          int
		AllowGuests           bool
		AllowRegistration     bool
		AllowJiraImport       bool
		DefaultLocale         string
		FriendlyUIVerbs       bool
		AuthMethod            string
		AppVersion            string
		CookieName            string
		PathPrefix            string
		APIEnabled            bool
		CleanupGuestsDaysOld  int
		CleanupBattlesDaysOld int
	}
	type UIConfig struct {
		AnalyticsEnabled bool
		AnalyticsID      string
		AppConfig        AppConfig
	}

	// get the html template from dist, have it ready for requests
	tmplContent, ioErr := fs.ReadFile(f, "dist/index.html")
	if ioErr != nil {
		log.Println("Error opening index template")
		log.Fatal(ioErr)
	}

	tmplString := string(tmplContent)
	tmpl, tmplErr := template.New("index").Parse(tmplString)
	if tmplErr != nil {
		log.Println("Error parsing index template")
		log.Fatal(tmplErr)
	}

	appConfig := AppConfig{
		AllowedPointValues:    viper.GetStringSlice("config.allowedPointValues"),
		DefaultPointValues:    viper.GetStringSlice("config.defaultPointValues"),
		ShowWarriorRank:       viper.GetBool("config.show_warrior_rank"),
		AvatarService:         viper.GetString("config.avatar_service"),
		ToastTimeout:          viper.GetInt("config.toast_timeout"),
		AllowGuests:           viper.GetBool("config.allow_guests"),
		AllowRegistration:     viper.GetBool("config.allow_registration") && viper.GetString("auth.method") == "normal",
		AllowJiraImport:       viper.GetBool("config.allow_jira_import"),
		DefaultLocale:         viper.GetString("config.default_locale"),
		FriendlyUIVerbs:       viper.GetBool("config.friendly_ui_verbs"),
		AuthMethod:            viper.GetString("auth.method"),
		APIEnabled:            viper.GetBool("config.allow_external_api"),
		AppVersion:            s.config.Version,
		CookieName:            s.config.FrontendCookieName,
		PathPrefix:            s.config.PathPrefix,
		CleanupGuestsDaysOld:  viper.GetInt("config.cleanup_guests_days_old"),
		CleanupBattlesDaysOld: viper.GetInt("config.cleanup_battles_days_old"),
	}

	data := UIConfig{
		AnalyticsEnabled: s.config.AnalyticsEnabled,
		AnalyticsID:      s.config.AnalyticsID,
		AppConfig:        appConfig,
	}

	return func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, data)
	}
}

/*
	Auth Handlers
*/

// handleLogin attempts to login the warrior by comparing email/password to whats in DB
func (s *server) handleLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body) // check for errors

		keyVal := make(map[string]string)
		json.Unmarshal(body, &keyVal) // check for errors
		WarriorEmail := keyVal["warriorEmail"]
		WarriorPassword := keyVal["warriorPassword"]

		authedWarrior, err := s.authUserDatabase(WarriorEmail, WarriorPassword)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		cookie := s.createCookie(authedWarrior.WarriorID)
		if cookie != nil {
			http.SetCookie(w, cookie)
		} else {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		RespondWithJSON(w, http.StatusOK, authedWarrior)
	}
}

// handleLdapLogin attempts to authenticate the warrior by looking up and authenticating
// via ldap, and then creates the warrior if not existing and logs them in
func (s *server) handleLdapLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		keyVal := make(map[string]string)
		json.Unmarshal(body, &keyVal)
		WarriorEmail := keyVal["warriorEmail"]
		WarriorPassword := keyVal["warriorPassword"]

		authedWarrior, err := s.authAndCreateUserLdap(WarriorEmail, WarriorPassword)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		cookie := s.createCookie(authedWarrior.WarriorID)
		if cookie != nil {
			http.SetCookie(w, cookie)
		} else {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		RespondWithJSON(w, http.StatusOK, authedWarrior)
	}
}

// handleLogout clears the warrior cookie(s) ending session
func (s *server) handleLogout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.clearUserCookies(w)
		return
	}
}

// handleUserRecruit registers a user as a private warrior (guest)
func (s *server) handleUserRecruit() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		AllowGuests := viper.GetBool("config.allow_guests")
		if !AllowGuests {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		body, _ := ioutil.ReadAll(r.Body) // check for errors

		keyVal := make(map[string]string)
		jsonErr := json.Unmarshal(body, &keyVal) // check for errors
		if jsonErr != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		WarriorName := keyVal["warriorName"]

		newWarrior, err := s.database.CreateUserGuest(WarriorName)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		s.createUserCookie(w, false, newWarrior.WarriorID)

		RespondWithJSON(w, http.StatusOK, newWarrior)
	}
}

// handleUserEnlist registers a user as a corporal warrior (authenticated)
func (s *server) handleUserEnlist() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		AllowRegistration := viper.GetBool("config.allow_registration")
		if !AllowRegistration {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		body, _ := ioutil.ReadAll(r.Body) // check for errors
		keyVal := make(map[string]string)
		jsonErr := json.Unmarshal(body, &keyVal) // check for errors
		if jsonErr != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		ActiveWarriorID, _ := s.validateUserCookie(w, r)

		WarriorName, WarriorEmail, WarriorPassword, accountErr := ValidateUserAccount(
			keyVal["warriorName"],
			keyVal["warriorEmail"],
			keyVal["warriorPassword1"],
			keyVal["warriorPassword2"],
		)

		if accountErr != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		newWarrior, VerifyID, err := s.database.CreateUserRegistered(WarriorName, WarriorEmail, WarriorPassword, ActiveWarriorID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		s.createUserCookie(w, true, newWarrior.WarriorID)

		s.email.SendWelcome(WarriorName, WarriorEmail, VerifyID)

		RespondWithJSON(w, http.StatusOK, newWarrior)
	}
}

// handleForgotPassword attempts to send a password reset email
func (s *server) handleForgotPassword() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body) // check for errors

		keyVal := make(map[string]string)
		json.Unmarshal(body, &keyVal) // check for errors
		WarriorEmail := keyVal["warriorEmail"]

		ResetID, WarriorName, resetErr := s.database.UserResetRequest(WarriorEmail)
		if resetErr == nil {
			s.email.SendForgotPassword(WarriorName, WarriorEmail, ResetID)
		}

		w.WriteHeader(http.StatusOK)
		return
	}
}

// handleResetPassword attempts to reset a warriors password
func (s *server) handleResetPassword() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body) // check for errors

		keyVal := make(map[string]string)
		json.Unmarshal(body, &keyVal) // check for errors
		ResetID := keyVal["resetId"]

		WarriorPassword, passwordErr := ValidateUserPassword(
			keyVal["warriorPassword1"],
			keyVal["warriorPassword2"],
		)

		if passwordErr != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		WarriorName, WarriorEmail, resetErr := s.database.UserResetPassword(ResetID, WarriorPassword)
		if resetErr != nil {
			log.Println("error attempting to reset warrior password : " + resetErr.Error() + "\n")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		s.email.SendPasswordReset(WarriorName, WarriorEmail)

		return
	}
}

/*
	Warrior Handlers
*/

// handleUpdatePassword attempts to update a warriors password
func (s *server) handleUpdatePassword() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body) // check for errors
		keyVal := make(map[string]string)
		json.Unmarshal(body, &keyVal) // check for errors
		warriorID := r.Context().Value(contextKeyWarriorID).(string)

		WarriorPassword, passwordErr := ValidateUserPassword(
			keyVal["warriorPassword1"],
			keyVal["warriorPassword2"],
		)

		if passwordErr != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		WarriorName, WarriorEmail, updateErr := s.database.UserUpdatePassword(warriorID, WarriorPassword)
		if updateErr != nil {
			log.Println("error attempting to update warrior password : " + updateErr.Error() + "\n")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		s.email.SendPasswordUpdate(WarriorName, WarriorEmail)

		return
	}
}

// handleUserProfile returns the warriors profile if it matches their session
func (s *server) handleUserProfile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		WarriorID := vars["id"]
		warriorCookieID := r.Context().Value(contextKeyWarriorID).(string)

		if WarriorID != warriorCookieID {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		warrior, warErr := s.database.GetUser(WarriorID)
		if warErr != nil {
			log.Println("error finding warrior : " + warErr.Error() + "\n")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		RespondWithJSON(w, http.StatusOK, warrior)
	}
}

// handleUserProfileUpdate attempts to update warriors profile (currently limited to name)
func (s *server) handleUserProfileUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		body, _ := ioutil.ReadAll(r.Body) // check for errors
		keyVal := make(map[string]interface{})
		json.Unmarshal(body, &keyVal) // check for errors
		WarriorName := keyVal["warriorName"].(string)
		WarriorAvatar := keyVal["warriorAvatar"].(string)
		NotificationsEnabled, _ := keyVal["notificationsEnabled"].(bool)

		WarriorID := vars["id"]
		warriorCookieID := r.Context().Value(contextKeyWarriorID).(string)
		if WarriorID != warriorCookieID {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		updateErr := s.database.UpdateUserProfile(WarriorID, WarriorName, WarriorAvatar, NotificationsEnabled)
		if updateErr != nil {
			log.Println("error attempting to update warrior profile : " + updateErr.Error() + "\n")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		warrior, warErr := s.database.GetUser(WarriorID)
		if warErr != nil {
			log.Println("error reloading warrior after update : " + warErr.Error() + "\n")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		RespondWithJSON(w, http.StatusOK, warrior)
	}
}

// handleAccountVerification attempts to verify a warriors account
func (s *server) handleAccountVerification() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body) // check for errors

		keyVal := make(map[string]string)
		json.Unmarshal(body, &keyVal) // check for errors
		VerifyID := keyVal["verifyId"]

		verifyErr := s.database.VerifyUserAccount(VerifyID)
		if verifyErr != nil {
			log.Println("error attempting to verify warrior account : " + verifyErr.Error() + "\n")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		return
	}
}

// handleUserDelete attempts to delete a warriors account
func (s *server) handleUserDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		WarriorID := vars["id"]
		warriorCookieID := r.Context().Value(contextKeyWarriorID).(string)
		if WarriorID != warriorCookieID {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		updateErr := s.database.DeleteUser(WarriorID)
		if updateErr != nil {
			log.Println("error attempting to delete warrior : " + updateErr.Error() + "\n")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		s.clearUserCookies(w)

		return
	}
}

// handleUserAvatar creates an avatar for the given warrior by ID
func (s *server) handleUserAvatar() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		Width, _ := strconv.Atoi(vars["width"])
		WarriorID := vars["id"]
		AvatarGender := govatar.MALE
		warriorGender, ok := vars["avatar"]
		if ok {
			if warriorGender == "female" {
				AvatarGender = govatar.FEMALE
			}
		}

		var avatar image.Image
		if s.config.AvatarService == "govatar" {
			avatar, _ = govatar.GenerateForUsername(AvatarGender, WarriorID)
		} else { // must be goadorable
			var err error
			avatar, _, err = image.Decode(bytes.NewReader(adorable.PseudoRandom([]byte(WarriorID))))
			if err != nil {
				log.Fatalln(err)
			}
		}

		img := transform.Resize(avatar, Width, Width, transform.Linear)
		buffer := new(bytes.Buffer)

		if err := png.Encode(buffer, img); err != nil {
			log.Println("unable to encode image.")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "image/png")
		w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))

		if _, err := w.Write(buffer.Bytes()); err != nil {
			log.Println("unable to write image.")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

/*
	API Key Handlers
*/

// handleAPIKeyGenerate handles generating an API key for a warrior
func (s *server) handleAPIKeyGenerate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		body, _ := ioutil.ReadAll(r.Body) // check for errors
		keyVal := make(map[string]interface{})
		json.Unmarshal(body, &keyVal) // check for errors
		APIKeyName := keyVal["name"].(string)

		WarriorID := vars["id"]
		warriorCookieID := r.Context().Value(contextKeyWarriorID).(string)
		if WarriorID != warriorCookieID {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		APIKey, keyErr := s.database.GenerateAPIKey(WarriorID, APIKeyName)
		if keyErr != nil {
			log.Println("error attempting to generate api key : " + keyErr.Error() + "\n")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		RespondWithJSON(w, http.StatusOK, APIKey)
	}
}

// handleUserAPIKeys handles getting warrior API keys
func (s *server) handleUserAPIKeys() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		WarriorID := vars["id"]
		warriorCookieID := r.Context().Value(contextKeyWarriorID).(string)
		if WarriorID != warriorCookieID {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		APIKeys, keysErr := s.database.GetUserAPIKeys(WarriorID)
		if keysErr != nil {
			log.Println("error retrieving api keys : " + keysErr.Error() + "\n")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		RespondWithJSON(w, http.StatusOK, APIKeys)
	}
}

// handleUserAPIKeyUpdate handles getting warrior API keys
func (s *server) handleUserAPIKeyUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		WarriorID := vars["id"]
		warriorCookieID := r.Context().Value(contextKeyWarriorID).(string)
		if WarriorID != warriorCookieID {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		APK := vars["keyID"]
		body, _ := ioutil.ReadAll(r.Body) // check for errors
		keyVal := make(map[string]interface{})
		json.Unmarshal(body, &keyVal) // check for errors
		active := keyVal["active"].(bool)

		APIKeys, keysErr := s.database.UpdateUserAPIKey(WarriorID, APK, active)
		if keysErr != nil {
			log.Println("error updating api key : " + keysErr.Error() + "\n")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		RespondWithJSON(w, http.StatusOK, APIKeys)
	}
}

// handleUserAPIKeyDelete handles getting warrior API keys
func (s *server) handleUserAPIKeyDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		WarriorID := vars["id"]
		warriorCookieID := r.Context().Value(contextKeyWarriorID).(string)
		if WarriorID != warriorCookieID {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		APK := vars["keyID"]

		APIKeys, keysErr := s.database.DeleteUserAPIKey(WarriorID, APK)
		if keysErr != nil {
			log.Println("error deleting api key : " + keysErr.Error() + "\n")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		RespondWithJSON(w, http.StatusOK, APIKeys)
	}
}

/*
	Battle Handlers
*/
// handleBattleCreate handles creating a battle (arena)
func (s *server) handleBattleCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		warriorID := r.Context().Value(contextKeyWarriorID).(string)
		body, bodyErr := ioutil.ReadAll(r.Body) // check for errors
		if bodyErr != nil {
			log.Println("error in reading request body: " + bodyErr.Error() + "\n")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var keyVal struct {
			BattleName           string           `json:"battleName"`
			PointValuesAllowed   []string         `json:"pointValuesAllowed"`
			AutoFinishVoting     bool             `json:"autoFinishVoting"`
			Plans                []*database.Plan `json:"plans"`
			PointAverageRounding string           `json:"pointAverageRounding"`
		}
		json.Unmarshal(body, &keyVal) // check for errors

		newBattle, err := s.database.CreateBattle(warriorID, keyVal.BattleName, keyVal.PointValuesAllowed, keyVal.Plans, keyVal.AutoFinishVoting, keyVal.PointAverageRounding)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		RespondWithJSON(w, http.StatusOK, newBattle)
	}
}

// handleBattlesGet looks up battles associated with warriorID
func (s *server) handleBattlesGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		warriorID := r.Context().Value(contextKeyWarriorID).(string)
		battles, err := s.database.GetBattlesByUser(warriorID)

		if err != nil {
			http.NotFound(w, r)
			return
		}

		RespondWithJSON(w, http.StatusOK, battles)
	}
}

/*
	Admin Handlers
*/

// handleAppStats gets the applications stats
func (s *server) handleAppStats() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		AppStats, err := s.database.GetAppStats()

		if err != nil {
			http.NotFound(w, r)
			return
		}

		RespondWithJSON(w, http.StatusOK, AppStats)
	}
}

// handleGetRegisteredUsers gets a list of registered warriors
func (s *server) handleGetRegisteredUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		Limit, _ := strconv.Atoi(vars["limit"])
		Offset, _ := strconv.Atoi(vars["offset"])

		Warriors := s.database.GetRegisteredUsers(Limit, Offset)

		RespondWithJSON(w, http.StatusOK, Warriors)
	}
}

// handleUserCreate registers a user as a corporal warrior (authenticated)
func (s *server) handleUserCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body) // check for errors
		keyVal := make(map[string]string)
		jsonErr := json.Unmarshal(body, &keyVal) // check for errors
		if jsonErr != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		WarriorName, WarriorEmail, WarriorPassword, accountErr := ValidateUserAccount(
			keyVal["warriorName"],
			keyVal["warriorEmail"],
			keyVal["warriorPassword1"],
			keyVal["warriorPassword2"],
		)

		if accountErr != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		newWarrior, VerifyID, err := s.database.CreateUserRegistered(WarriorName, WarriorEmail, WarriorPassword, "")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		s.email.SendWelcome(WarriorName, WarriorEmail, VerifyID)

		RespondWithJSON(w, http.StatusOK, newWarrior)
	}
}

// handleUserPromote handles promoting a warrior to General (ADMIN) by ID
func (s *server) handleUserPromote() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body) // check for errors
		keyVal := make(map[string]string)
		jsonErr := json.Unmarshal(body, &keyVal) // check for errors
		if jsonErr != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err := s.database.PromoteUser(keyVal["warriorId"])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		return
	}
}

// handleUserDemote handles demoting a warrior to Corporal (Registered) by ID
func (s *server) handleUserDemote() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body) // check for errors
		keyVal := make(map[string]string)
		jsonErr := json.Unmarshal(body, &keyVal) // check for errors
		if jsonErr != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err := s.database.DemoteUser(keyVal["warriorId"])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		return
	}
}

// handleCleanBattles handles cleaning up old battles (ADMIN Manaually Triggered)
func (s *server) handleCleanBattles() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		DaysOld := viper.GetInt("config.cleanup_battles_days_old")

		err := s.database.CleanBattles(DaysOld)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		return
	}
}

// handleCleanGuests handles cleaning up old guests (ADMIN Manaually Triggered)
func (s *server) handleCleanGuests() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		DaysOld := viper.GetInt("config.cleanup_guests_days_old")

		err := s.database.CleanGuests(DaysOld)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		return
	}
}
