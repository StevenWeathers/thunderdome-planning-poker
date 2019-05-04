package main

import (
	"net/http"

	"github.com/google/uuid"
)

func CreateBattleHandler(w http.ResponseWriter, r *http.Request) {
	id, _ := uuid.NewUUID()	
	newBattle := &Battle{ BattleId: id.String(), CreatorId: "1" }

	RespondWithJSON(w, http.StatusOK, newBattle)
}
