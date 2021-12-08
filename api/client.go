package api

import (
	"encoding/json"
	"github.com/StevenWeathers/thunderdome-planning-poker/model"
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

// leaderOnlyOperations contains a map of operations that only a battle leader can execute
var leaderOnlyOperations = map[string]bool{
	"add_plan":       true,
	"revise_plan":    true,
	"burn_plan":      true,
	"activate_plan":  true,
	"skip_plan":      true,
	"end_voting":     true,
	"finalize_plan":  true,
	"jab_warrior":    true,
	"promote_leader": true,
	"demote_leader":  true,
	"revise_battle":  true,
	"concede_battle": true,
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// connection is a middleman between the websocket connection and the hub.
type connection struct {
	// The websocket connection.
	ws *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

// socketEvent is the event structure used for socket messages
type socketEvent struct {
	Type  string `json:"type"`
	Value string `json:"value"`
	User  string `json:"warriorId"`
}

func createSocketEvent(Type string, Value string, User string) []byte {
	newEvent := &socketEvent{
		Type:  Type,
		Value: Value,
		User:  User,
	}

	event, _ := json.Marshal(newEvent)

	return event
}

// readPump pumps messages from the websocket connection to the hub.
func (sub subscription) readPump(api *api) {
	var forceClosed bool
	c := sub.conn
	defer func() {
		BattleID := sub.arena
		UserID := sub.UserID

		Users := api.db.RetreatUser(BattleID, UserID)
		UpdatedUsers, _ := json.Marshal(Users)

		retreatEvent := createSocketEvent("warrior_retreated", string(UpdatedUsers), UserID)
		m := message{retreatEvent, BattleID}
		h.broadcast <- m

		h.unregister <- sub
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
		authorizedOperation := true
		keyVal := make(map[string]string)
		json.Unmarshal(msg, &keyVal) // check for errors
		UserID := sub.UserID
		battleID := sub.arena

		eventType := keyVal["type"]

		// confirm leader for any operation that requires it
		if _, ok := leaderOnlyOperations[eventType]; ok {
			err := api.db.ConfirmLeader(battleID, UserID)
			if err != nil {
				badEvent = true
				authorizedOperation = false
			}
		}

		if authorizedOperation {
			switch eventType {
			case "vote":
				var wv struct {
					VoteValue        string `json:"voteValue"`
					PlanID           string `json:"planId"`
					AutoFinishVoting bool   `json:"autoFinishVoting"`
				}
				json.Unmarshal([]byte(keyVal["value"]), &wv)

				Plans, AllVoted := api.db.SetVote(battleID, UserID, wv.PlanID, wv.VoteValue)

				updatedPlans, _ := json.Marshal(Plans)
				msg = createSocketEvent("vote_activity", string(updatedPlans), UserID)

				if AllVoted && wv.AutoFinishVoting {
					plans, err := api.db.EndPlanVoting(battleID, wv.PlanID)
					if err != nil {
						badEvent = true
						break
					}
					updatedPlans, _ := json.Marshal(plans)
					msg = createSocketEvent("voting_ended", string(updatedPlans), "")
				}
			case "retract_vote":
				PlanID := keyVal["value"]

				plans := api.db.RetractVote(battleID, UserID, PlanID)

				updatedPlans, _ := json.Marshal(plans)
				msg = createSocketEvent("vote_retracted", string(updatedPlans), UserID)
			case "add_plan":
				var p struct {
					Name               string `json:"planName"`
					Type               string `json:"type"`
					ReferenceId        string `json:"referenceId"`
					Link               string `json:"link"`
					Description        string `json:"description"`
					AcceptanceCriteria string `json:"acceptanceCriteria"`
				}
				json.Unmarshal([]byte(keyVal["value"]), &p)

				plans, err := api.db.CreatePlan(battleID, p.Name, p.Type, p.ReferenceId, p.Link, p.Description, p.AcceptanceCriteria)
				if err != nil {
					badEvent = true
					break
				}
				updatedPlans, _ := json.Marshal(plans)
				msg = createSocketEvent("plan_added", string(updatedPlans), "")
			case "activate_plan":
				plans, err := api.db.ActivatePlanVoting(battleID, keyVal["value"])
				if err != nil {
					badEvent = true
					break
				}
				updatedPlans, _ := json.Marshal(plans)
				msg = createSocketEvent("plan_activated", string(updatedPlans), "")
			case "skip_plan":
				plans, err := api.db.SkipPlan(battleID, keyVal["value"])
				if err != nil {
					badEvent = true
					break
				}
				updatedPlans, _ := json.Marshal(plans)
				msg = createSocketEvent("plan_skipped", string(updatedPlans), "")
			case "end_voting":
				plans, err := api.db.EndPlanVoting(battleID, keyVal["value"])
				if err != nil {
					badEvent = true
					break
				}
				updatedPlans, _ := json.Marshal(plans)
				msg = createSocketEvent("voting_ended", string(updatedPlans), "")
			case "finalize_plan":
				var p struct {
					Id     string `json:"planId"`
					Points string `json:"planPoints"`
				}
				json.Unmarshal([]byte(keyVal["value"]), &p)

				plans, err := api.db.FinalizePlan(battleID, p.Id, p.Points)
				if err != nil {
					badEvent = true
					break
				}
				updatedPlans, _ := json.Marshal(plans)
				msg = createSocketEvent("plan_finalized", string(updatedPlans), "")
			case "revise_plan":
				var p struct {
					Id                 string `json:"planId"`
					Name               string `json:"planName"`
					Type               string `json:"type"`
					ReferenceId        string `json:"referenceId"`
					Link               string `json:"link"`
					Description        string `json:"description"`
					AcceptanceCriteria string `json:"acceptanceCriteria"`
				}
				json.Unmarshal([]byte(keyVal["value"]), &p)

				plans, err := api.db.RevisePlan(battleID, p.Id, p.Name, p.Type, p.ReferenceId, p.Link, p.Description, p.AcceptanceCriteria)
				if err != nil {
					badEvent = true
					break
				}
				updatedPlans, _ := json.Marshal(plans)
				msg = createSocketEvent("plan_revised", string(updatedPlans), "")
			case "burn_plan":
				plans, err := api.db.BurnPlan(battleID, keyVal["value"])
				if err != nil {
					badEvent = true
					break
				}
				updatedPlans, _ := json.Marshal(plans)
				msg = createSocketEvent("plan_burned", string(updatedPlans), "")
			case "promote_leader":
				leaders, err := api.db.SetBattleLeader(battleID, keyVal["value"])
				if err != nil {
					badEvent = true
					break
				}
				leadersJson, _ := json.Marshal(leaders)

				msg = createSocketEvent("leaders_updated", string(leadersJson), "")
			case "demote_leader":
				leaders, err := api.db.DemoteBattleLeader(battleID, keyVal["value"])
				if err != nil {
					badEvent = true
					break
				}
				leadersJson, _ := json.Marshal(leaders)

				msg = createSocketEvent("leaders_updated", string(leadersJson), "")
			case "become_leader":
				leaderCode, err := api.db.GetBattleLeaderCode(battleID, api.config.AESHashkey)
				if err != nil {
					badEvent = true
					break
				}

				if keyVal["value"] == leaderCode {
					leaders, err := api.db.SetBattleLeader(battleID, UserID)
					if err != nil {
						badEvent = true
						break
					}
					leadersJson, _ := json.Marshal(leaders)

					msg = createSocketEvent("leaders_updated", string(leadersJson), "")
				} else {
					badEvent = true
					break
				}
			case "spectator_toggle":
				var st struct {
					Spectator bool `json:"spectator"`
				}
				json.Unmarshal([]byte(keyVal["value"]), &st)
				users, err := api.db.ToggleSpectator(battleID, UserID, st.Spectator)
				if err != nil {
					badEvent = true
					break
				}
				usersJson, _ := json.Marshal(users)

				msg = createSocketEvent("users_updated", string(usersJson), "")
			case "revise_battle":
				var revisedBattle struct {
					BattleName           string   `json:"battleName"`
					PointValuesAllowed   []string `json:"pointValuesAllowed"`
					AutoFinishVoting     bool     `json:"autoFinishVoting"`
					PointAverageRounding string   `json:"pointAverageRounding"`
					JoinCode             string   `json:"joinCode"`
					LeaderCode           string   `json:"leaderCode"`
				}
				json.Unmarshal([]byte(keyVal["value"]), &revisedBattle)

				err := api.db.ReviseBattle(battleID, revisedBattle.BattleName, revisedBattle.PointValuesAllowed, revisedBattle.AutoFinishVoting, revisedBattle.PointAverageRounding)
				if err != nil {
					badEvent = true
					break
				}

				if revisedBattle.JoinCode != "" {
					err = api.db.ReviseBattleJoinCode(battleID, revisedBattle.JoinCode, api.config.AESHashkey)
					if err != nil {
						badEvent = true
						break
					}
				}

				if revisedBattle.LeaderCode != "" {
					err = api.db.ReviseBattleLeaderCode(battleID, revisedBattle.LeaderCode, api.config.AESHashkey)
					if err != nil {
						badEvent = true
						break
					}

					revisedBattle.LeaderCode = ""
				}

				updatedBattle, _ := json.Marshal(revisedBattle)
				msg = createSocketEvent("battle_revised", string(updatedBattle), "")
			case "concede_battle":
				err := api.db.DeleteBattle(battleID)
				if err != nil {
					badEvent = true
					break
				}
				msg = createSocketEvent("battle_conceded", "", "")
			case "jab_warrior":
			case "abandon_battle":
				_, err := api.db.AbandonBattle(battleID, UserID)
				if err != nil {
					badEvent = true
					break
				}
				badEvent = true // don't want this event to cause write panic
				forceClosed = true
			default:
				badEvent = true
			}
		}

		if !badEvent {
			m := message{msg, sub.arena}
			h.broadcast <- m
		}

		if forceClosed {
			break
		}
	}
}

// write a message with the given message type and payload.
func (c *connection) write(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(mt, payload)
}

// writePump pumps messages from the hub to the websocket connection.
func (sub *subscription) writePump() {
	c := sub.conn
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

// handleSocketUnauthorized sets the format close message and closes the websocket
func handleSocketClose(ws *websocket.Conn, closeCode int, text string) {
	cm := websocket.FormatCloseMessage(closeCode, text)
	if err := ws.WriteMessage(websocket.CloseMessage, cm); err != nil {
		log.Printf("unauthorized close error: %v", err)
	}
	if err := ws.Close(); err != nil {
		log.Printf("close error: %v", err)
	}
}

// serveBattleWs handles websocket requests from the peer.
func (a *api) serveBattleWs() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		battleID := vars["battleId"]
		var User *model.User
		var UserAuthed bool

		// upgrade to WebSocket connection
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		c := &connection{send: make(chan []byte, 256), ws: ws}

		SessionId, cookieErr := a.validateSessionCookie(w, r)
		if cookieErr != nil && cookieErr.Error() != "NO_SESSION_COOKIE" {
			handleSocketClose(ws, 4001, "unauthorized")
			return
		}

		if SessionId != "" {
			var userErr error
			User, userErr = a.db.GetSessionUser(SessionId)
			if userErr != nil {
				handleSocketClose(ws, 4001, "unauthorized")
				return
			}
		} else {
			UserID, err := a.validateUserCookie(w, r)
			if err != nil {
				handleSocketClose(ws, 4001, "unauthorized")
				return
			}

			var userErr error
			User, userErr = a.db.GetGuestUser(UserID)
			if userErr != nil {
				handleSocketClose(ws, 4001, "unauthorized")
				return
			}
		}

		// make sure battle is legit
		b, battleErr := a.db.GetBattle(battleID, User.Id, a.config.AESHashkey)
		if battleErr != nil {
			handleSocketClose(ws, 4004, "battle not found")
			return
		}

		// check users battle active status
		UserErr := a.db.GetBattleUserActiveStatus(battleID, User.Id)
		if UserErr != nil && UserErr.Error() != "sql: no rows in result set" {
			usrErrMsg := UserErr.Error()
			log.Println("error finding user : " + usrErrMsg + "\n")
			if usrErrMsg == "DUPLICATE_BATTLE_USER" {
				handleSocketClose(ws, 4003, "duplicate session")
			} else {
				handleSocketClose(ws, 4005, "internal error")
			}
			return
		}

		if b.JoinCode != "" && (UserErr != nil && UserErr.Error() == "sql: no rows in result set") {
			jcrEvent := createSocketEvent("join_code_required", "", User.Id)
			_ = c.write(websocket.TextMessage, jcrEvent)

			for {
				_, msg, err := c.ws.ReadMessage()
				if err != nil {
					if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
						log.Printf("error: %v", err)
					}
					break
				}

				keyVal := make(map[string]string)
				json.Unmarshal(msg, &keyVal)

				if keyVal["type"] == "auth_battle" && keyVal["value"] == b.JoinCode {
					UserAuthed = true
					break
				} else if keyVal["type"] == "auth_battle" {
					authIncorrect := createSocketEvent("join_code_incorrect", "", User.Id)
					_ = c.write(websocket.TextMessage, authIncorrect)
				}
			}
		} else {
			UserAuthed = true
		}

		for {
			if UserAuthed == true {
				ss := subscription{c, battleID, User.Id}
				h.register <- ss

				Users, _ := a.db.AddUserToBattle(ss.arena, User.Id)
				UpdatedUsers, _ := json.Marshal(Users)

				battle, _ := json.Marshal(b)
				initEvent := createSocketEvent("init", string(battle), User.Id)
				_ = c.write(websocket.TextMessage, initEvent)

				joinedEvent := createSocketEvent("warrior_joined", string(UpdatedUsers), User.Id)
				m := message{joinedEvent, ss.arena}
				h.broadcast <- m

				go ss.writePump()
				go ss.readPump(a)

				break
			}
		}
	}
}
