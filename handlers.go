package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// CreateBattleHandler handles creating a battle (arena)
func CreateBattleHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body) // check for errors

	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)  // check for errors
	LeaderID := keyVal["leaderId"] // @TODO get this from https cookie when implemented
	BattleName := keyVal["battleName"]

	newBattle, err := CreateBattle(LeaderID, BattleName)
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
