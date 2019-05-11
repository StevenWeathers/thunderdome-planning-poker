package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
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
	WarriorID  string `json:"id"`
	EventValue string `json:"value"`
}

// CreateSocketEvent makes a SocketEvent struct and turns it into json []byte
func CreateSocketEvent(EventType string, WarriorID string, EventValue string) []byte {
	newEvent := &SocketEvent{
		EventType:  EventType,
		WarriorID:  WarriorID,
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

		retreatEvent := CreateSocketEvent("user_activity", WarriorID, string(updatedWarriors))
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

		keyVal := make(map[string]string)
		json.Unmarshal(msg, &keyVal) // check for errors
		warriorID := keyVal["id"]
		battleID := s.arena

		switch keyVal["type"] {
		case "vote":
			voteObj := make(map[string]string)
			json.Unmarshal([]byte(keyVal["value"]), &voteObj)
			VoteValue := voteObj["voteValue"]
			PlanID := voteObj["planId"]
			log.Println(VoteValue)
			log.Println(PlanID)

			plans := SetVote(battleID, warriorID, PlanID, VoteValue)
			updatedPlans, _ := json.Marshal(plans)
			msg = CreateSocketEvent("vote_activity", warriorID, string(updatedPlans))
		case "add_plan":
			plans := CreatePlan(battleID, keyVal["value"])
			updatedPlans, _ := json.Marshal(plans)
			msg = CreateSocketEvent("plan_added", warriorID, string(updatedPlans))
		case "activate_plan":
			plans := ActivatePlanVoting(battleID, keyVal["value"])
			updatedPlans, _ := json.Marshal(plans)
			msg = CreateSocketEvent("plan_activated", warriorID, string(updatedPlans))
		case "end_voting":
			plans := EndPlanVoting(battleID, keyVal["value"])
			updatedPlans, _ := json.Marshal(plans)
			msg = CreateSocketEvent("voting_ended", warriorID, string(updatedPlans))
		case "finalize_plan":
			planObj := make(map[string]string)
			json.Unmarshal([]byte(keyVal["value"]), &planObj)
			PlanID := planObj["planId"]
			PlanPoints := planObj["planPoints"]

			plans := FinalizePlan(battleID, PlanID, PlanPoints)
			updatedPlans, _ := json.Marshal(plans)
			msg = CreateSocketEvent("plan_finalized", warriorID, string(updatedPlans))
		case "revise_plan":
			planObj := make(map[string]string)
			json.Unmarshal([]byte(keyVal["value"]), &planObj)
			PlanID := planObj["planId"]
			PlanName := planObj["planName"]

			plans := RevisePlanName(battleID, PlanID, PlanName)
			updatedPlans, _ := json.Marshal(plans)
			msg = CreateSocketEvent("plan_burned", warriorID, string(updatedPlans))
		case "burn_plan":
			plans := BurnPlan(battleID, keyVal["value"])
			updatedPlans, _ := json.Marshal(plans)
			msg = CreateSocketEvent("plan_burned", warriorID, string(updatedPlans))
		default:
		}

		m := message{msg, s.arena}
		h.broadcast <- m
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
	warrior, err := r.Cookie("warrior")
	var warriorID string
	var warriorName string

	if err != nil {
		log.Println("error in reading warrior cookie : " + err.Error() + "\n")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	value, _ := url.PathUnescape(warrior.Value)
	keyVal := make(map[string]string)
	json.Unmarshal([]byte(value), &keyVal) // check for errors

	warriorID = keyVal["id"]
	warriorName = keyVal["name"]

	_, warErr := GetWarrior(warriorID)

	if warErr != nil {
		Warriors[warriorID] = &Warrior{
			WarriorID:   warriorID,
			WarriorName: warriorName}
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	vars := mux.Vars(r)
	c := &connection{send: make(chan []byte, 256), ws: ws}
	s := subscription{c, vars["id"], warriorID}
	h.register <- s

	

	Warriors := AddWarriorToBattle(s.arena, warriorID)
	updatedWarriors, _ := json.Marshal(Warriors)

	joinedEvent := CreateSocketEvent("user_activity", warriorID, string(updatedWarriors))
	m := message{joinedEvent, s.arena}
	h.broadcast <- m

	go s.writePump()
	s.readPump()
}
