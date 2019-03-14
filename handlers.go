package main

import (
	"log"
	"net/http"
)

func CreateBattleHandler(w http.ResponseWriter, r *http.Request) {
	newBattle := &Battle{ BattleId: "1", CreatorId: "1" }

	RespondWithJSON(w, http.StatusOK, newBattle)
}

func BattleHandler(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}