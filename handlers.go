package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// Battle aka arena
type Battle struct {
	BattleID   string     `json:"id"`
	LeaderID   string     `json:"leaderId"`
	BattleName string     `json:"name"`
	Warriors   []*Warrior `json:"warriors"`
	Votes      []*Vote    `json:"votes"`
}

// Warrior aka user
type Warrior struct {
	WarriorID   string `json:"id"`
	WarriorName string `json:"name"`
}

// Vote structure
type Vote struct {
	WarriorID string `json:"warriorId"`
	VoteValue string `json:"vote"`
}

// Warriors stores all warriors in memory
var Warriors = make(map[string]*Warrior)

// Battles stores all battles in memory
var Battles = make(map[string]*Battle)

// CreateBattleHandler handles creating a battle (arena)
func CreateBattleHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body) // check for errors

	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal) // check for errors
	leaderID := keyVal["leaderId"]

	newID, _ := uuid.NewUUID()
	id := newID.String()

	Battles[id] = &Battle{
		BattleID:   id,
		LeaderID:   leaderID,
		BattleName: keyVal["battleName"],
		Warriors:   make([]*Warrior, 0),
		Votes:      make([]*Vote, 0)}

	RespondWithJSON(w, http.StatusOK, Battles[id])
}

// RecruitWarriorHandler registeres a user as a warrior in memory
func RecruitWarriorHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body) // check for errors

	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal) // check for errors

	newID, _ := uuid.NewUUID()
	id := newID.String()
	Warriors[id] = &Warrior{WarriorID: id, WarriorName: keyVal["warriorName"]}

	RespondWithJSON(w, http.StatusOK, Warriors[id])
}

// GetBattlePlansHandler looks up battle in memory or returns notfound status
func GetBattlePlansHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if battle, ok := Battles[id]; ok {
		RespondWithJSON(w, http.StatusOK, battle)
	} else {
		http.NotFound(w, r)
		return
	}
}
