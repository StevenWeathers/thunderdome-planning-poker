package main

import (
	"encoding/json"
	"fmt"
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
	EventUser  string `json:"warriorId"`
}

// CreateSocketEvent makes a SocketEvent struct and turns it into json []byte
func CreateSocketEvent(EventType string, EventValue string, EventUser string) []byte {
	newEvent := &SocketEvent{
		EventType:  EventType,
		EventValue: EventValue,
		EventUser:  EventUser,
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
		UserID := s.UserID

		Users := srv.database.RetreatUser(BattleID, UserID)
		UpdatedUsers, _ := json.Marshal(Users)

		retreatEvent := CreateSocketEvent("warrior_retreated", string(UpdatedUsers), UserID)
		m := message{retreatEvent, BattleID}
		h.broadcast <- m

		h.unregister <- s
		if forceClosed {
			cm := websocket.FormatCloseMessage(4002, "abandoned")
			if err := c.ws.WriteControl(websocket.CloseMessage, cm, time.Now().Add(writeWait)); err != nil {
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
		UserID := s.UserID
		battleID := s.arena

		switch keyVal["type"] {
		case "vote":
			var wv struct {
				VoteValue        string `json:"voteValue"`
				PlanID           string `json:"planId"`
				AutoFinishVoting bool   `json:"autoFinishVoting"`
			}
			json.Unmarshal([]byte(keyVal["value"]), &wv)

			Plans, AllVoted := srv.database.SetVote(battleID, UserID, wv.PlanID, wv.VoteValue)

			updatedPlans, _ := json.Marshal(Plans)
			msg = CreateSocketEvent("vote_activity", string(updatedPlans), UserID)

			if AllVoted && wv.AutoFinishVoting {
				plans, err := srv.database.EndPlanVoting(battleID, UserID, wv.PlanID, true)
				if err != nil {
					badEvent = true
					break
				}
				updatedPlans, _ := json.Marshal(plans)
				msg = CreateSocketEvent("voting_ended", string(updatedPlans), "")
			}
		case "retract_vote":
			PlanID := keyVal["value"]

			plans := srv.database.RetractVote(battleID, UserID, PlanID)

			updatedPlans, _ := json.Marshal(plans)
			msg = CreateSocketEvent("vote_retracted", string(updatedPlans), UserID)
		case "add_plan":
			planObj := make(map[string]string)
			json.Unmarshal([]byte(keyVal["value"]), &planObj)
			PlanName := planObj["planName"]
			PlanType := planObj["type"]
			ReferenceID := planObj["referenceId"]
			Link := planObj["link"]
			Description := planObj["description"]
			AcceptanceCriteria := planObj["acceptanceCriteria"]

			plans, err := srv.database.CreatePlan(battleID, UserID, PlanName, PlanType, ReferenceID, Link, Description, AcceptanceCriteria)
			if err != nil {
				badEvent = true
				break
			}
			updatedPlans, _ := json.Marshal(plans)
			msg = CreateSocketEvent("plan_added", string(updatedPlans), "")
		case "activate_plan":
			plans, err := srv.database.ActivatePlanVoting(battleID, UserID, keyVal["value"])
			if err != nil {
				badEvent = true
				break
			}
			updatedPlans, _ := json.Marshal(plans)
			msg = CreateSocketEvent("plan_activated", string(updatedPlans), "")
		case "skip_plan":
			plans, err := srv.database.SkipPlan(battleID, UserID, keyVal["value"])
			if err != nil {
				badEvent = true
				break
			}
			updatedPlans, _ := json.Marshal(plans)
			msg = CreateSocketEvent("plan_skipped", string(updatedPlans), "")
		case "end_voting":
			plans, err := srv.database.EndPlanVoting(battleID, UserID, keyVal["value"], false)
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

			plans, err := srv.database.FinalizePlan(battleID, UserID, PlanID, PlanPoints)
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

			plans, err := srv.database.RevisePlan(battleID, UserID, PlanID, PlanName, PlanType, ReferenceID, Link, Description, AcceptanceCriteria)
			if err != nil {
				badEvent = true
				break
			}
			updatedPlans, _ := json.Marshal(plans)
			msg = CreateSocketEvent("plan_revised", string(updatedPlans), "")
		case "burn_plan":
			plans, err := srv.database.BurnPlan(battleID, UserID, keyVal["value"])
			if err != nil {
				badEvent = true
				break
			}
			updatedPlans, _ := json.Marshal(plans)
			msg = CreateSocketEvent("plan_burned", string(updatedPlans), "")
		case "promote_leader":
			leaders, err := srv.database.SetBattleLeader(battleID, UserID, keyVal["value"])
			if err != nil {
				badEvent = true
				break
			}
			leadersJson, _ := json.Marshal(leaders)

			msg = CreateSocketEvent("leaders_updated", string(leadersJson), "")
		case "demote_leader":
			leaders, err := srv.database.DemoteBattleLeader(battleID, UserID, keyVal["value"])
			if err != nil {
				badEvent = true
				break
			}
			leadersJson, _ := json.Marshal(leaders)

			msg = CreateSocketEvent("leaders_updated", string(leadersJson), "")
		case "revise_battle":
			var revisedBattle struct {
				BattleName           string   `json:"battleName"`
				PointValuesAllowed   []string `json:"pointValuesAllowed"`
				AutoFinishVoting     bool     `json:"autoFinishVoting"`
				PointAverageRounding string   `json:"pointAverageRounding"`
			}
			json.Unmarshal([]byte(keyVal["value"]), &revisedBattle)

			err := srv.database.ReviseBattle(battleID, UserID, revisedBattle.BattleName, revisedBattle.PointValuesAllowed, revisedBattle.AutoFinishVoting, revisedBattle.PointAverageRounding)
			if err != nil {
				badEvent = true
				break
			}

			updatedBattle, _ := json.Marshal(revisedBattle)
			msg = CreateSocketEvent("battle_revised", string(updatedBattle), "")
		case "concede_battle":
			err := srv.database.DeleteBattle(battleID, UserID)
			if err != nil {
				badEvent = true
				break
			}
			msg = CreateSocketEvent("battle_conceded", "", "")
		case "jab_warrior":
			err := srv.database.ConfirmLeader(battleID, UserID)
			if err != nil {
				badEvent = true
				break
			}
		case "abandon_battle":
			_, err := srv.database.AbandonBattle(battleID, UserID)
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

		// make sure user cookies are valid
		UserID, cookieErr := s.validateUserCookie(w, r)
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
		b, battleErr := s.database.GetBattle(battleID, UserID)
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

		// make sure user exists
		_, UserErr := s.database.GetBattleUser(battleID, UserID)

		if UserErr != nil {
			log.Println("error finding user : " + UserErr.Error() + "\n")
			cm := websocket.FormatCloseMessage(4003, "duplicate session")

			if fmt.Sprint(UserErr) == "user not found" {
				s.clearUserCookies(w)
				cm = websocket.FormatCloseMessage(4001, "unauthorized")
			}

			if err := ws.WriteMessage(websocket.CloseMessage, cm); err != nil {
				log.Printf("unauthorized close error: %v", err)
			}
			if err := ws.Close(); err != nil {
				log.Printf("close error: %v", err)
			}
			return
		}

		c := &connection{send: make(chan []byte, 256), ws: ws}
		ss := subscription{c, battleID, UserID}
		h.register <- ss

		Users, _ := s.database.AddUserToBattle(ss.arena, UserID)
		UpdatedUsers, _ := json.Marshal(Users)

		initEvent := CreateSocketEvent("init", string(battle), UserID)
		_ = c.write(websocket.TextMessage, initEvent)

		joinedEvent := CreateSocketEvent("warrior_joined", string(UpdatedUsers), UserID)
		m := message{joinedEvent, ss.arena}
		h.broadcast <- m

		go ss.writePump()
		go ss.readPump(s)
	}
}
