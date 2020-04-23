package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/markbates/pkger"
)

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

		adminErr := ConfirmAdmin(warriorID)
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
	type UIConfig struct {
		AnalyticsEnabled bool
		AnalyticsID      string
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

	data := UIConfig{
		AnalyticsEnabled: s.config.AnalyticsEnabled,
		AnalyticsID:      s.config.AnalyticsID,
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

		authedWarrior, err := AuthWarrior(WarriorEmail, WarriorPassword)
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

		_, warErr := GetWarrior(warriorID)

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
			BattleName         string   `json:"battleName"`
			PointValuesAllowed []string `json:"pointValuesAllowed"`
			Plans              []*Plan  `json:"plans"`
		}
		json.Unmarshal(body, &keyVal) // check for errors

		newBattle, err := CreateBattle(warriorID, keyVal.BattleName, keyVal.PointValuesAllowed, keyVal.Plans)
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

		newWarrior, err := CreateWarriorPrivate(WarriorName)
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

		newWarrior, VerifyID, err := CreateWarriorCorporal(WarriorName, WarriorEmail, WarriorPassword, ActiveWarriorID)
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

		battle, err := GetBattle(BattleID)

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

		_, warErr := GetWarrior(warriorID)

		if warErr != nil {
			log.Println("error finding warrior : " + warErr.Error() + "\n")
			s.clearWarriorCookies(w)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		battles, err := GetBattlesByWarrior(warriorID)

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

		ResetID, WarriorName, resetErr := WarriorResetRequest(WarriorEmail)
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

		WarriorName, WarriorEmail, resetErr := WarriorResetPassword(ResetID, WarriorPassword)
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

		WarriorName, WarriorEmail, updateErr := WarriorUpdatePassword(warriorID, WarriorPassword)
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

		warrior, warErr := GetWarrior(WarriorID)
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

		WarriorID := vars["id"]
		warriorCookieID, cookieErr := s.validateWarriorCookie(w, r)
		if cookieErr != nil || WarriorID != warriorCookieID {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		updateErr := UpdateWarriorProfile(WarriorID, WarriorName)
		if updateErr != nil {
			log.Println("error attempting to update warrior profile : " + updateErr.Error() + "\n")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		return
	}
}

// handleAccountVerification attempts to verify a warriors account
func (s *server) handleAccountVerification() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body) // check for errors

		keyVal := make(map[string]string)
		json.Unmarshal(body, &keyVal) // check for errors
		VerifyID := keyVal["verifyId"]

		verifyErr := VerifyWarriorAccount(VerifyID)
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
		AppStats, err := GetAppStats()

		if err != nil {
			http.NotFound(w, r)
			return
		}

		RespondWithJSON(w, http.StatusOK, AppStats)
	}
}
