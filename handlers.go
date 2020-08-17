package main

import (
	"encoding/json"
	"errors"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/StevenWeathers/thunderdome-planning-poker/pkg/database"
	"github.com/gorilla/mux"
	"github.com/markbates/pkger"
	"github.com/spf13/viper"
	"gopkg.in/go-playground/validator.v9"
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

// ValidateWarriorAccount makes sure warrior name, email, and password are valid before creating the account
func ValidateWarriorAccount(name string, email string, pwd1 string, pwd2 string) (WarriorName string, WarriorEmail string, WarriorPassword string, validateErr error) {
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

// ValidateWarriorPassword makes sure warrior password is valid before updating the password
func ValidateWarriorPassword(pwd1 string, pwd2 string) (WarriorPassword string, validateErr error) {
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

// createWarriorCookie creates the warriors cookie
func (s *server) createWarriorCookie(w http.ResponseWriter, isRegistered bool, WarriorID string) {
	var cookiedays = 365 // 356 days
	if isRegistered == true {
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
		Path:     "/",
		HttpOnly: true,
		Domain:   s.config.AppDomain,
		MaxAge:   86400 * cookiedays,
		Secure:   s.config.SecureCookieFlag,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, cookie)
}

// clearWarriorCookies wipes the frontend and backend cookies
// used in the event of bad cookie reads
func (s *server) clearWarriorCookies(w http.ResponseWriter) {
	feCookie := &http.Cookie{
		Name:   s.config.FrontendCookieName,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	beCookie := &http.Cookie{
		Name:     s.config.SecureCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	}

	http.SetCookie(w, feCookie)
	http.SetCookie(w, beCookie)
}

// validateWarriorCookie returns the warriorID from secure cookies or errors if failures getting it
func (s *server) validateWarriorCookie(w http.ResponseWriter, r *http.Request) (string, error) {
	var warriorID string

	if cookie, err := r.Cookie(s.config.SecureCookieName); err == nil {
		var value string
		if err = s.cookie.Decode(s.config.SecureCookieName, cookie.Value, &value); err == nil {
			warriorID = value
		} else {
			log.Println("error in reading warrior cookie : " + err.Error() + "\n")
			s.clearWarriorCookies(w)
			return "", errors.New("invalid warrior cookies")
		}
	} else {
		log.Println("error in reading warrior cookie : " + err.Error() + "\n")
		s.clearWarriorCookies(w)
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
		warriorID, cookieErr := s.validateWarriorCookie(w, r)
		if cookieErr != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		adminErr := s.database.ConfirmAdmin(warriorID)
		if adminErr != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		h(w, r)
	}
}

/*
	Handlers
*/

// handleIndex parses the index html file, injecting any relevant data
func (s *server) handleIndex() http.HandlerFunc {
	type AppConfig struct {
		AllowedPointValues []string
		DefaultPointValues []string
		ShowWarriorRank    bool
		AvatarService      string
		ToastTimeout       int
	}
	type UIConfig struct {
		AnalyticsEnabled bool
		AnalyticsID      string
		AppConfig        AppConfig
	}

	// get the html template from dist, have it ready for requests
	indexFile, ioErr := pkger.Open("/dist/index.html")
	if ioErr != nil {
		log.Println("Error opening index template")
		log.Fatal(ioErr)
	}
	tmplContent, ioReadErr := ioutil.ReadAll(indexFile)
	if ioReadErr != nil {
		// this will hopefully only possibly panic during development as the file is already in memory otherwise
		log.Println("Error reading index template")
		log.Fatal(ioReadErr)
	}

	tmplString := string(tmplContent)
	tmpl, tmplErr := template.New("index").Parse(tmplString)
	if tmplErr != nil {
		log.Println("Error parsing index template")
		log.Fatal(tmplErr)
	}

	appConfig := AppConfig{
		AllowedPointValues: viper.GetStringSlice("config.allowedPointValues"),
		DefaultPointValues: viper.GetStringSlice("config.defaultPointValues"),
		ShowWarriorRank:    viper.GetBool("config.show_warrior_rank"),
		AvatarService:      viper.GetString("config.avatar_service"),
		ToastTimeout:       viper.GetInt("config.toast_timeout"),
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

// handleLogin attempts to login the warrior by comparing email/password to whats in DB
func (s *server) handleLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body) // check for errors

		keyVal := make(map[string]string)
		json.Unmarshal(body, &keyVal) // check for errors
		WarriorEmail := keyVal["warriorEmail"]
		WarriorPassword := keyVal["warriorPassword"]

		authedWarrior, err := s.database.AuthWarrior(WarriorEmail, WarriorPassword)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		encoded, err := s.cookie.Encode(s.config.SecureCookieName, authedWarrior.WarriorID)
		if err == nil {
			cookie := &http.Cookie{
				Name:     s.config.SecureCookieName,
				Value:    encoded,
				Path:     "/",
				HttpOnly: true,
				Domain:   s.config.AppDomain,
				MaxAge:   86400 * 30, // 30 days
				Secure:   s.config.SecureCookieFlag,
				SameSite: http.SameSiteStrictMode,
			}
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
		s.clearWarriorCookies(w)
		return
	}
}

// handleBattleCreate handles creating a battle (arena)
func (s *server) handleBattleCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		warriorID, cookieErr := s.validateWarriorCookie(w, r)
		if cookieErr != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		_, warErr := s.database.GetWarrior(warriorID)

		if warErr != nil {
			log.Println("error finding warrior : " + warErr.Error() + "\n")
			s.clearWarriorCookies(w)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		body, bodyErr := ioutil.ReadAll(r.Body) // check for errors
		if bodyErr != nil {
			log.Println("error in reading warrior cookie : " + bodyErr.Error() + "\n")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var keyVal struct {
			BattleName         string           `json:"battleName"`
			PointValuesAllowed []string         `json:"pointValuesAllowed"`
			Plans              []*database.Plan `json:"plans"`
		}
		json.Unmarshal(body, &keyVal) // check for errors

		newBattle, err := s.database.CreateBattle(warriorID, keyVal.BattleName, keyVal.PointValuesAllowed, keyVal.Plans)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		RespondWithJSON(w, http.StatusOK, newBattle)
	}
}

// handleWarriorRecruit registers a user as a private warrior (guest)
func (s *server) handleWarriorRecruit() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body) // check for errors

		keyVal := make(map[string]string)
		jsonErr := json.Unmarshal(body, &keyVal) // check for errors
		if jsonErr != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		WarriorName := keyVal["warriorName"]

		newWarrior, err := s.database.CreateWarriorPrivate(WarriorName)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		s.createWarriorCookie(w, false, newWarrior.WarriorID)

		RespondWithJSON(w, http.StatusOK, newWarrior)
	}
}

// handleWarriorEnlist registers a user as a corporal warrior (authenticated)
func (s *server) handleWarriorEnlist() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body) // check for errors
		keyVal := make(map[string]string)
		jsonErr := json.Unmarshal(body, &keyVal) // check for errors
		if jsonErr != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		ActiveWarriorID, _ := s.validateWarriorCookie(w, r)

		WarriorName, WarriorEmail, WarriorPassword, accountErr := ValidateWarriorAccount(
			keyVal["warriorName"],
			keyVal["warriorEmail"],
			keyVal["warriorPassword1"],
			keyVal["warriorPassword2"],
		)

		if accountErr != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		newWarrior, VerifyID, err := s.database.CreateWarriorCorporal(WarriorName, WarriorEmail, WarriorPassword, ActiveWarriorID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		s.createWarriorCookie(w, true, newWarrior.WarriorID)

		s.email.SendWelcome(WarriorName, WarriorEmail, VerifyID)

		RespondWithJSON(w, http.StatusOK, newWarrior)
	}
}

// handleBattleGet looks up battle or returns notfound status
func (s *server) handleBattleGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		BattleID := vars["id"]

		battle, err := s.database.GetBattle(BattleID)

		if err != nil {
			http.NotFound(w, r)
			return
		}

		RespondWithJSON(w, http.StatusOK, battle)
	}
}

// handleBattlesGet looks up battles associated with warriorID
func (s *server) handleBattlesGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		warriorID, cookieErr := s.validateWarriorCookie(w, r)
		if cookieErr != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		_, warErr := s.database.GetWarrior(warriorID)

		if warErr != nil {
			log.Println("error finding warrior : " + warErr.Error() + "\n")
			s.clearWarriorCookies(w)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		battles, err := s.database.GetBattlesByWarrior(warriorID)

		if err != nil {
			http.NotFound(w, r)
			return
		}

		RespondWithJSON(w, http.StatusOK, battles)
	}
}

// handleForgotPassword attempts to send a password reset email
func (s *server) handleForgotPassword() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body) // check for errors

		keyVal := make(map[string]string)
		json.Unmarshal(body, &keyVal) // check for errors
		WarriorEmail := keyVal["warriorEmail"]

		ResetID, WarriorName, resetErr := s.database.WarriorResetRequest(WarriorEmail)
		if resetErr != nil {
			log.Println("error attempting to send warrior reset : " + resetErr.Error() + "\n")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		s.email.SendForgotPassword(WarriorName, WarriorEmail, ResetID)

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

		WarriorPassword, passwordErr := ValidateWarriorPassword(
			keyVal["warriorPassword1"],
			keyVal["warriorPassword2"],
		)

		if passwordErr != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		WarriorName, WarriorEmail, resetErr := s.database.WarriorResetPassword(ResetID, WarriorPassword)
		if resetErr != nil {
			log.Println("error attempting to reset warrior password : " + resetErr.Error() + "\n")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		s.email.SendPasswordReset(WarriorName, WarriorEmail)

		return
	}
}

// handleUpdatePassword attempts to update a warriors password
func (s *server) handleUpdatePassword() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body) // check for errors
		keyVal := make(map[string]string)
		json.Unmarshal(body, &keyVal) // check for errors

		warriorID, cookieErr := s.validateWarriorCookie(w, r)
		if cookieErr != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		WarriorPassword, passwordErr := ValidateWarriorPassword(
			keyVal["warriorPassword1"],
			keyVal["warriorPassword2"],
		)

		if passwordErr != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		WarriorName, WarriorEmail, updateErr := s.database.WarriorUpdatePassword(warriorID, WarriorPassword)
		if updateErr != nil {
			log.Println("error attempting to update warrior password : " + updateErr.Error() + "\n")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		s.email.SendPasswordUpdate(WarriorName, WarriorEmail)

		return
	}
}

// handleWarriorProfile returns the warriors profile if it matches their session
func (s *server) handleWarriorProfile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		WarriorID := vars["id"]

		warriorCookieID, cookieErr := s.validateWarriorCookie(w, r)
		if cookieErr != nil || WarriorID != warriorCookieID {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		warrior, warErr := s.database.GetWarrior(WarriorID)
		if warErr != nil {
			log.Println("error finding warrior : " + warErr.Error() + "\n")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		RespondWithJSON(w, http.StatusOK, warrior)
	}
}

// handleWarriorProfileUpdate attempts to update warriors profile (currently limited to name)
func (s *server) handleWarriorProfileUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		body, _ := ioutil.ReadAll(r.Body) // check for errors
		keyVal := make(map[string]string)
		json.Unmarshal(body, &keyVal) // check for errors
		WarriorName := keyVal["warriorName"]
		WarriorAvatar := keyVal["warriorAvatar"]

		WarriorID := vars["id"]
		warriorCookieID, cookieErr := s.validateWarriorCookie(w, r)
		if cookieErr != nil || WarriorID != warriorCookieID {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		updateErr := s.database.UpdateWarriorProfile(WarriorID, WarriorName, WarriorAvatar)
		if updateErr != nil {
			log.Println("error attempting to update warrior profile : " + updateErr.Error() + "\n")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		warrior, warErr := s.database.GetWarrior(WarriorID)
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

		verifyErr := s.database.VerifyWarriorAccount(VerifyID)
		if verifyErr != nil {
			log.Println("error attempting to verify warrior account : " + verifyErr.Error() + "\n")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		return
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
