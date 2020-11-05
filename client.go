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
	EventType    string `json:"type"`
	EventValue   string `json:"value"`
	EventWarrior string `json:"warriorId"`
}

// CreateSocketEvent makes a SocketEvent struct and turns it into json []byte
func CreateSocketEvent(EventType string, EventValue string, EventWarrior string) []byte {
	newEvent := &SocketEvent{
		EventType:    EventType,
		EventValue:   EventValue,
		EventWarrior: EventWarrior,
	}

	event, _ := json.Marshal(newEvent)

	return event
}

// readPump pumps messages from the websocket connection to the hub.
func (s subscription) readPump(srv *server) {
	var forceClosed bool
	c := s.conn
	defer func() {
		BattleID := s.arena
		WarriorID := s.warriorID

		Warriors := srv.database.RetreatWarrior(BattleID, WarriorID)
		updatedWarriors, _ := json.Marshal(Warriors)

		retreatEvent := CreateSocketEvent("warrior_retreated", string(updatedWarriors), WarriorID)
		m := message{retreatEvent, BattleID}
		h.broadcast <- m

		h.unregister <- s
		if forceClosed {
			cm := websocket.FormatCloseMessage(4002, "abandoned")
			if err := c.ws.WriteMessage(websocket.CloseMessage, cm); err != nil {
				log.Printf("abandon error: %v", err)
			}
		}
		if err := c.ws.Close(); err != nil {
			log.Printf("close error: %v", err)
		}
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
			var wv struct {
				VoteValue        string `json:"voteValue"`
				PlanID           string `json:"planId"`
				AutoFinishVoting bool   `json:"autoFinishVoting"`
			}
			json.Unmarshal([]byte(keyVal["value"]), &wv)

			Plans, AllVoted := srv.database.SetVote(battleID, warriorID, wv.PlanID, wv.VoteValue)

			updatedPlans, _ := json.Marshal(Plans)
			msg = CreateSocketEvent("vote_activity", string(updatedPlans), warriorID)

			if AllVoted && wv.AutoFinishVoting {
				plans, err := srv.database.EndPlanVoting(battleID, warriorID, wv.PlanID, true)
				if err != nil {
					badEvent = true
					break
				}
				updatedPlans, _ := json.Marshal(plans)
				msg = CreateSocketEvent("voting_ended", string(updatedPlans), "")
			}
		case "retract_vote":
			PlanID := keyVal["value"]

			plans := srv.database.RetractVote(battleID, warriorID, PlanID)

			updatedPlans, _ := json.Marshal(plans)
			msg = CreateSocketEvent("vote_retracted", string(updatedPlans), warriorID)
		case "add_plan":
			planObj := make(map[string]string)
			json.Unmarshal([]byte(keyVal["value"]), &planObj)
			PlanName := planObj["planName"]
			PlanType := planObj["type"]
			ReferenceID := planObj["referenceId"]
			Link := planObj["link"]
			Description := planObj["description"]
			AcceptanceCriteria := planObj["acceptanceCriteria"]

			plans, err := srv.database.CreatePlan(battleID, warriorID, PlanName, PlanType, ReferenceID, Link, Description, AcceptanceCriteria)
			if err != nil {
				badEvent = true
				break
			}
			updatedPlans, _ := json.Marshal(plans)
			msg = CreateSocketEvent("plan_added", string(updatedPlans), "")
		case "activate_plan":
			plans, err := srv.database.ActivatePlanVoting(battleID, warriorID, keyVal["value"])
			if err != nil {
				badEvent = true
				break
			}
			updatedPlans, _ := json.Marshal(plans)
			msg = CreateSocketEvent("plan_activated", string(updatedPlans), "")
		case "skip_plan":
			plans, err := srv.database.SkipPlan(battleID, warriorID, keyVal["value"])
			if err != nil {
				badEvent = true
				break
			}
			updatedPlans, _ := json.Marshal(plans)
			msg = CreateSocketEvent("plan_skipped", string(updatedPlans), "")
		case "end_voting":
			plans, err := srv.database.EndPlanVoting(battleID, warriorID, keyVal["value"], false)
			if err != nil {
				badEvent = true
				break
			}
			updatedPlans, _ := json.Marshal(plans)
			msg = CreateSocketEvent("voting_ended", string(updatedPlans), "")
		case "finalize_plan":
			planObj := make(map[string]string)
			json.Unmarshal([]byte(keyVal["value"]), &planObj)
			PlanID := planObj["planId"]
			PlanPoints := planObj["planPoints"]

			plans, err := srv.database.FinalizePlan(battleID, warriorID, PlanID, PlanPoints)
			if err != nil {
				badEvent = true
				break
			}
			updatedPlans, _ := json.Marshal(plans)
			msg = CreateSocketEvent("plan_finalized", string(updatedPlans), "")
		case "revise_plan":
			planObj := make(map[string]string)
			json.Unmarshal([]byte(keyVal["value"]), &planObj)
			PlanID := planObj["planId"]
			PlanName := planObj["planName"]
			PlanType := planObj["type"]
			ReferenceID := planObj["referenceId"]
			Link := planObj["link"]
			Description := planObj["description"]
			AcceptanceCriteria := planObj["acceptanceCriteria"]

			plans, err := srv.database.RevisePlan(battleID, warriorID, PlanID, PlanName, PlanType, ReferenceID, Link, Description, AcceptanceCriteria)
			if err != nil {
				badEvent = true
				break
			}
			updatedPlans, _ := json.Marshal(plans)
			msg = CreateSocketEvent("plan_revised", string(updatedPlans), "")
		case "burn_plan":
			plans, err := srv.database.BurnPlan(battleID, warriorID, keyVal["value"])
			if err != nil {
				badEvent = true
				break
			}
			updatedPlans, _ := json.Marshal(plans)
			msg = CreateSocketEvent("plan_burned", string(updatedPlans), "")
		case "promote_leader":
			err := srv.database.SetBattleLeader(battleID, warriorID, keyVal["value"])
			if err != nil {
				badEvent = true
				break
			}

			msg = CreateSocketEvent("leader_updated", keyVal["value"], "")
		case "revise_battle":
			var revisedBattle struct {
				BattleName         string   `json:"battleName"`
				PointValuesAllowed []string `json:"pointValuesAllowed"`
				AutoFinishVoting   bool     `json:"autoFinishVoting"`
			}
			json.Unmarshal([]byte(keyVal["value"]), &revisedBattle)

			err := srv.database.ReviseBattle(battleID, warriorID, revisedBattle.BattleName, revisedBattle.PointValuesAllowed, revisedBattle.AutoFinishVoting)
			if err != nil {
				badEvent = true
				break
			}

			updatedBattle, _ := json.Marshal(revisedBattle)
			msg = CreateSocketEvent("battle_revised", string(updatedBattle), "")
		case "concede_battle":
			err := srv.database.DeleteBattle(battleID, warriorID)
			if err != nil {
				badEvent = true
				break
			}
			msg = CreateSocketEvent("battle_conceded", "", "")
		case "jab_warrior":
			err := srv.database.ConfirmLeader(battleID, warriorID)
			if err != nil {
				badEvent = true
				break
			}
		case "abandon_battle":
			_, err := srv.database.AbandonBattle(battleID, warriorID)
			if err != nil {
				badEvent = true
				break
			}
			badEvent = true // don't want this event to cause write panic
			forceClosed = true
		default:
		}

		if !badEvent {
			m := message{msg, s.arena}
			h.broadcast <- m
		}

		if forceClosed {
			break
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
func (s *server) serveWs() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		battleID := vars["id"]

		// upgrade to WebSocket connection
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}

		// make sure warrior cookies are valid
		warriorID, cookieErr := s.validateWarriorCookie(w, r)
		if cookieErr != nil {
			cm := websocket.FormatCloseMessage(4001, "unauthorized")
			if err := ws.WriteMessage(websocket.CloseMessage, cm); err != nil {
				log.Printf("unauthorized close error: %v", err)
			}
			if err := ws.Close(); err != nil {
				log.Printf("close error: %v", err)
			}
			return
		}

		// make sure battle is legit
		b, battleErr := s.database.GetBattle(battleID, warriorID)
		if battleErr != nil {
			cm := websocket.FormatCloseMessage(4004, "battle not found")
			if err := ws.WriteMessage(websocket.CloseMessage, cm); err != nil {
				log.Printf("not found close error: %v", err)
			}
			if err := ws.Close(); err != nil {
				log.Printf("close error: %v", err)
			}
			return
		}
		battle, _ := json.Marshal(b)

		// make sure warrior exists
		_, warErr := s.database.GetBattleWarrior(battleID, warriorID)

		if warErr != nil {
			log.Println("error finding warrior : " + warErr.Error() + "\n")
			s.clearWarriorCookies(w)
			cm := websocket.FormatCloseMessage(4001, "unauthorized")
			if err := ws.WriteMessage(websocket.CloseMessage, cm); err != nil {
				log.Printf("unauthorized close error: %v", err)
			}
			if err := ws.Close(); err != nil {
				log.Printf("close error: %v", err)
			}
			return
		}

		c := &connection{send: make(chan []byte, 256), ws: ws}
		ss := subscription{c, battleID, warriorID}
		h.register <- ss

		Warriors, _ := s.database.AddWarriorToBattle(ss.arena, warriorID)
		updatedWarriors, _ := json.Marshal(Warriors)

		initEvent := CreateSocketEvent("init", string(battle), warriorID)
		_ = c.write(websocket.TextMessage, initEvent)

		joinedEvent := CreateSocketEvent("warrior_joined", string(updatedWarriors), warriorID)
		m := message{joinedEvent, ss.arena}
		h.broadcast <- m

		go ss.writePump()
		go ss.readPump(s)
	}
}
