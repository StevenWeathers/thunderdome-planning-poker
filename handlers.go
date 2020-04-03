package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// LoginHandler attempts to login the warrior by comparing email/password to whats in DB
func LoginHandler(w http.ResponseWriter, r *http.Request) {
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

	encoded, err := Sc.Encode(SecureCookieName, authedWarrior.WarriorID)
	if err == nil {
		cookie := &http.Cookie{
			Name:     SecureCookieName,
			Value:    encoded,
			Path:     "/",
			HttpOnly: true,
			Domain:   AppDomain,
			MaxAge:   86400 * 30, // 30 days
			Secure:   SecureCookieFlag,
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

// LogoutHandler clears the warrior cookie(s) ending session
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	ClearWarriorCookies(w)
	return
}

// CreateBattleHandler handles creating a battle (arena)
func CreateBattleHandler(w http.ResponseWriter, r *http.Request) {
	warriorID, cookieErr := ValidateWarriorCookie(w, r)
	if cookieErr != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	_, warErr := GetWarrior(warriorID)

	if warErr != nil {
		log.Println("error finding warrior : " + warErr.Error() + "\n")
		ClearWarriorCookies(w)
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

// RecruitWarriorHandler registers a user as a private warrior (guest)
func RecruitWarriorHandler(w http.ResponseWriter, r *http.Request) {
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

	encoded, err := Sc.Encode(SecureCookieName, newWarrior.WarriorID)
	if err == nil {
		cookie := &http.Cookie{
			Name:     SecureCookieName,
			Value:    encoded,
			Path:     "/",
			HttpOnly: true,
			Domain:   AppDomain,
			MaxAge:   86400 * 365, // 365 days
			Secure:   SecureCookieFlag,
			SameSite: http.SameSiteStrictMode,
		}
		http.SetCookie(w, cookie)
	} else {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	RespondWithJSON(w, http.StatusOK, newWarrior)
}

// EnlistWarriorHandler registers a user as a corporal warrior (authenticated)
func EnlistWarriorHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body) // check for errors
	keyVal := make(map[string]string)
	jsonErr := json.Unmarshal(body, &keyVal) // check for errors
	if jsonErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ActiveWarriorID, _ := ValidateWarriorCookie(w, r)

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

	encoded, err := Sc.Encode(SecureCookieName, newWarrior.WarriorID)
	if err == nil {
		cookie := &http.Cookie{
			Name:     SecureCookieName,
			Value:    encoded,
			Path:     "/",
			HttpOnly: true,
			Domain:   AppDomain,
			MaxAge:   86400 * 30, // 30 days
			Secure:   SecureCookieFlag,
			SameSite: http.SameSiteStrictMode,
		}
		http.SetCookie(w, cookie)
	} else {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	SendWelcomeEmail(WarriorName, WarriorEmail, VerifyID)

	RespondWithJSON(w, http.StatusOK, newWarrior)
}

// GetBattleHandler looks up battle or returns notfound status
func GetBattleHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	BattleID := vars["id"]

	battle, err := GetBattle(BattleID)

	if err != nil {
		http.NotFound(w, r)
		return
	}

	RespondWithJSON(w, http.StatusOK, battle)
}

// GetBattlesHandler looks up battles associated with warriorID
func GetBattlesHandler(w http.ResponseWriter, r *http.Request) {
	warriorID, cookieErr := ValidateWarriorCookie(w, r)
	if cookieErr != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	_, warErr := GetWarrior(warriorID)

	if warErr != nil {
		log.Println("error finding warrior : " + warErr.Error() + "\n")
		ClearWarriorCookies(w)
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

// ForgotPasswordHandler attempts to send a password reset email
func ForgotPasswordHandler(w http.ResponseWriter, r *http.Request) {
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

	SendForgotPasswordEmail(WarriorName, WarriorEmail, ResetID)

	w.WriteHeader(http.StatusOK)
	return
}

// ResetPasswordHandler attempts to reset a warriors password
func ResetPasswordHandler(w http.ResponseWriter, r *http.Request) {
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

	SendPasswordResetEmail(WarriorName, WarriorEmail)

	return
}

// UpdatePasswordHandler attempts to update a warriors password
func UpdatePasswordHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body) // check for errors
	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal) // check for errors

	warriorID, cookieErr := ValidateWarriorCookie(w, r)
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

	SendPasswordUpdateEmail(WarriorName, WarriorEmail)

	return
}

// WarriorProfileHandler returns the warriors profile if it matches their session
func WarriorProfileHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	WarriorID := vars["id"]

	warriorCookieID, cookieErr := ValidateWarriorCookie(w, r)
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

// WarriorUpdateProfileHandler attempts to update warriors profile (currently limited to name)
func WarriorUpdateProfileHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	body, _ := ioutil.ReadAll(r.Body) // check for errors
	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal) // check for errors
	WarriorName := keyVal["warriorName"]

	WarriorID := vars["id"]
	warriorCookieID, cookieErr := ValidateWarriorCookie(w, r)
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

// GetAppStatsHandler gets the applications stats
func GetAppStatsHandler(w http.ResponseWriter, r *http.Request) {
	warriorID, cookieErr := ValidateWarriorCookie(w, r)
	if cookieErr != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	AppStats, err := GetAppStats(warriorID)

	if err != nil {
		http.NotFound(w, r)
		return
	}

	RespondWithJSON(w, http.StatusOK, AppStats)
}

// VerifyAccountHandler attempts to verify a warriors account
func VerifyAccountHandler(w http.ResponseWriter, r *http.Request) {
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
