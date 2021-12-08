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
	eventHandlers := map[string]func(string, string, string) ([]byte, error, bool){
		"jab_warrior":      api.battle.UserNudge,
		"vote":             api.battle.UserVote,
		"retract_vote":     api.battle.UserVoteRetract,
		"end_voting":       api.battle.PlanVoteEnd,
		"add_plan":         api.battle.PlanAdd,
		"revise_plan":      api.battle.PlanRevise,
		"burn_plan":        api.battle.PlanDelete,
		"activate_plan":    api.battle.PlanActivate,
		"skip_plan":        api.battle.PlanSkip,
		"finalize_plan":    api.battle.PlanFinalize,
		"promote_leader":   api.battle.UserPromote,
		"demote_leader":    api.battle.UserDemote,
		"become_leader":    api.battle.UserPromoteSelf,
		"spectator_toggle": api.battle.UserSpectatorToggle,
		"revise_battle":    api.battle.Revise,
		"concede_battle":   api.battle.Delete,
		"abandon_battle":   api.battle.Abandon,
	}

	var forceClosed bool
	c := sub.conn
	UserID := sub.UserID
	BattleID := sub.arena

	defer func() {
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
		var badEvent bool
		var eventErr error
		_, msg, err := c.ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}

		keyVal := make(map[string]string)
		json.Unmarshal(msg, &keyVal) // check for errors

		eventType := keyVal["type"]
		eventValue := keyVal["value"]

		// confirm leader for any operation that requires it
		if _, ok := leaderOnlyOperations[eventType]; ok {
			err := api.db.ConfirmLeader(BattleID, UserID)
			if err != nil {
				badEvent = true
			}
		}

		// find event handler and execute otherwise invalid event
		if _, ok := eventHandlers[eventType]; ok {
			msg, eventErr, forceClosed = eventHandlers[eventType](BattleID, UserID, eventValue)
			if eventErr != nil {
				badEvent = true

				// don't log forceClosed events e.g. Abandon
				if !forceClosed {
					log.Println(eventErr)
				}
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
		b, battleErr := a.db.GetBattle(battleID, User.Id)
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
