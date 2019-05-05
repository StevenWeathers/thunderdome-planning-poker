package main

import (
	"net/http"
	"io/ioutil"
	"encoding/json"

	"github.com/google/uuid"
)

type Battle struct {
	BattleId string `json:"id"`
	CreatorId string `json:"creatorId"`
	BattleName string `json:"name"`
}

func CreateBattleHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body) // check for errors

    keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal) // check for errors
	
	id, _ := uuid.NewUUID()	
	newBattle := &Battle{
		BattleId: id.String(),
		CreatorId: keyVal["creatorId"],
		BattleName: keyVal["battleName"] }

	RespondWithJSON(w, http.StatusOK, newBattle)
}

type User struct {
	UserId string `json:"id"`
	UserName string `json:"name"`
}

func RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body) // check for errors

    keyVal := make(map[string]string)
    json.Unmarshal(body, &keyVal) // check for errors

	id, _ := uuid.NewUUID()
	newUser := &User{ UserId: id.String(), UserName: keyVal["userName"] }

	RespondWithJSON(w, http.StatusOK, newUser)
}
