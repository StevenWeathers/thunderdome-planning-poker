package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 1024 * 1024
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// connection is an middleman between the websocket connection and the hub.
type connection struct {
	// The websocket connection.
	ws *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

// SocketEvent is the event structure used for socket messages
type SocketEvent struct {
	EventType  string `json:"type"`
	EventValue string `json:"value"`
}

// CreateSocketEvent makes a SocketEvent struct and turns it into json []byte
func CreateSocketEvent(EventType string, EventValue string) []byte {
	newEvent := &SocketEvent{
		EventType:  EventType,
		EventValue: EventValue}

	event, _ := json.Marshal(newEvent)

	return event
}

// readPump pumps messages from the websocket connection to the hub.
func (s subscription) readPump() {
	c := s.conn
	defer func() {
		BattleID := s.arena
		WarriorID := s.warriorID
		log.Println(s.warriorID + " has left the arena")

		Warriors := RetreatWarrior(BattleID, WarriorID)
		updatedWarriors, _ := json.Marshal(Warriors)

		retreatEvent := CreateSocketEvent("user_activity", string(updatedWarriors))
		m := message{retreatEvent, BattleID}
		h.broadcast <- m

		h.unregister <- s
		c.ws.Close()
	}()
	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, msg, err := c.ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}

		var badEvent bool
		keyVal := make(map[string]string)
		json.Unmarshal(msg, &keyVal) // check for errors
		warriorID := s.warriorID
		battleID := s.arena

		switch keyVal["type"] {
		case "vote":
			voteObj := make(map[string]string)
			json.Unmarshal([]byte(keyVal["value"]), &voteObj)
			VoteValue := voteObj["voteValue"]
			PlanID := voteObj["planId"]

			plans := SetVote(battleID, warriorID, PlanID, VoteValue)
			updatedPlans, _ := json.Marshal(plans)
			msg = CreateSocketEvent("vote_activity", string(updatedPlans))
		case "add_plan":
			plans, err := CreatePlan(battleID, warriorID, keyVal["value"])
			if err != nil {
				badEvent = true
				break
			}
			updatedPlans, _ := json.Marshal(plans)
			msg = CreateSocketEvent("plan_added", string(updatedPlans))
		case "activate_plan":
			plans, err := ActivatePlanVoting(battleID, warriorID, keyVal["value"])
			if err != nil {
				badEvent = true
				break
			}
			updatedPlans, _ := json.Marshal(plans)
			msg = CreateSocketEvent("plan_activated", string(updatedPlans))
		case "end_voting":
			plans, err := EndPlanVoting(battleID, warriorID, keyVal["value"])
			if err != nil {
				badEvent = true
				break
			}
			updatedPlans, _ := json.Marshal(plans)
			msg = CreateSocketEvent("voting_ended", string(updatedPlans))
		case "finalize_plan":
			planObj := make(map[string]string)
			json.Unmarshal([]byte(keyVal["value"]), &planObj)
			PlanID := planObj["planId"]
			PlanPoints := planObj["planPoints"]

			plans, err := FinalizePlan(battleID, warriorID, PlanID, PlanPoints)
			if err != nil {
				badEvent = true
				break
			}
			updatedPlans, _ := json.Marshal(plans)
			msg = CreateSocketEvent("plan_finalized", string(updatedPlans))
		case "revise_plan":
			planObj := make(map[string]string)
			json.Unmarshal([]byte(keyVal["value"]), &planObj)
			PlanID := planObj["planId"]
			PlanName := planObj["planName"]

			plans, err := RevisePlanName(battleID, warriorID, PlanID, PlanName)
			if err != nil {
				badEvent = true
				break
			}
			updatedPlans, _ := json.Marshal(plans)
			msg = CreateSocketEvent("plan_revised", string(updatedPlans))
		case "burn_plan":
			plans, err := BurnPlan(battleID, warriorID, keyVal["value"])
			if err != nil {
				badEvent = true
				break
			}
			updatedPlans, _ := json.Marshal(plans)
			msg = CreateSocketEvent("plan_burned", string(updatedPlans))
		case "promote_leader":
			battle, err := SetBattleLeader(battleID, warriorID, keyVal["value"])
			if err != nil {
				badEvent = true
				break
			}

			updatedBattle, _ := json.Marshal(battle)
			msg = CreateSocketEvent("battle_updated", string(updatedBattle))
		default:
		}

		if badEvent != true {
			m := message{msg, s.arena}
			h.broadcast <- m
		}
	}
}

// write writes a message with the given message type and payload.
func (c *connection) write(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(mt, payload)
}

// writePump pumps messages from the hub to the websocket connection.
func (s *subscription) writePump() {
	c := s.conn
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				return
			}
			log.Println(string(message))
			if err := c.write(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

// serveWs handles websocket requests from the peer.
func serveWs(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	battleID := vars["id"]
	var warriorID string

	if cookie, err := r.Cookie(SecureCookieName); err == nil {
		var value string
		if err = Sc.Decode(SecureCookieName, cookie.Value, &value); err == nil {
			warriorID = value
		} else {
			log.Println("error in reading warrior cookie : " + err.Error() + "\n")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	} else {
		log.Println("error in reading warrior cookie : " + err.Error() + "\n")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, warErr := GetWarrior(battleID, warriorID)

	if warErr != nil {
		log.Println("error finding warrior : " + warErr.Error() + "\n")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	c := &connection{send: make(chan []byte, 256), ws: ws}
	s := subscription{c, battleID, warriorID}
	h.register <- s

	Warriors, _ := AddWarriorToBattle(s.arena, warriorID)
	updatedWarriors, _ := json.Marshal(Warriors)

	joinedEvent := CreateSocketEvent("user_activity", string(updatedWarriors))
	m := message{joinedEvent, s.arena}
	h.broadcast <- m

	go s.writePump()
	s.readPump()
}
