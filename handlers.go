package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

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

	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal) // check for errors
	BattleName := keyVal["battleName"]

	newBattle, err := CreateBattle(warriorID, BattleName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	RespondWithJSON(w, http.StatusOK, newBattle)
}

// RecruitWarriorHandler registeres a user as a warrior
func RecruitWarriorHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body) // check for errors

	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal) // check for errors
	WarriorName := keyVal["warriorName"]

	newWarrior, err := CreateWarrior(WarriorName)
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
			Expires:  time.Now().Add(365 * (24 * time.Hour)),
			MaxAge:   86400 * 365,
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

	battles, err := GetBattlesByLeader(warriorID)

	if err != nil {
		http.NotFound(w, r)
		return
	}

	RespondWithJSON(w, http.StatusOK, battles)
}
