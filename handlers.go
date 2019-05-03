package main

import (
	// "log"
	"net/http"
)

func CreateBattleHandler(w http.ResponseWriter, r *http.Request) {
	newBattle := &Battle{ BattleId: "1", CreatorId: "1" }

	RespondWithJSON(w, http.StatusOK, newBattle)
}
