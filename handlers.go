package main

import (
	"net/http"
	"io/ioutil"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Battle struct {
	BattleId string `json:"id"`
	LeaderId string `json:"leaderId"`
	BattleName string `json:"name"`
	Warriors []*Warrior `json:"warriors"`
}

type Warrior struct {
	WarriorId string `json:"id"`
	WarriorName string `json:"name"`
}

var Warriors = make(map[string]*Warrior)
var Battles = make(map[string]*Battle)

func CreateBattleHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body) // check for errors

    keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal) // check for errors
	leaderId := keyVal["leaderId"]
	
	newId, _ := uuid.NewUUID()
	id := newId.String()
	
	Battles[id] = &Battle{
		BattleId: id,
		LeaderId: leaderId,
		BattleName: keyVal["battleName"],
		Warriors: []*Warrior{Warriors[leaderId]} }

	RespondWithJSON(w, http.StatusOK, Battles[id])
}

func RecruitWarriorHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body) // check for errors

    keyVal := make(map[string]string)
    json.Unmarshal(body, &keyVal) // check for errors

	newId, _ := uuid.NewUUID()
	id := newId.String()
	Warriors[id] = &Warrior{ WarriorId: id, WarriorName: keyVal["warriorName"] }

	RespondWithJSON(w, http.StatusOK, Warriors[id])
}

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