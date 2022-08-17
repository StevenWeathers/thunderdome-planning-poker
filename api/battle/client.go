package battle

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/StevenWeathers/thunderdome-planning-poker/model"
	"go.uber.org/zap"

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
var leaderOnlyOperations = map[string]struct{}{
	"add_plan":       {},
	"revise_plan":    {},
	"burn_plan":      {},
	"activate_plan":  {},
	"skip_plan":      {},
	"end_voting":     {},
	"finalize_plan":  {},
	"jab_warrior":    {},
	"promote_leader": {},
	"demote_leader":  {},
	"revise_battle":  {},
	"concede_battle": {},
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

// readPump pumps messages from the websocket connection to the hub.
func (sub subscription) readPump(b *Service, ctx context.Context) {
	var forceClosed bool
	c := sub.conn
	UserID := sub.UserID
	BattleID := sub.arena

	defer func() {
		Users := b.db.RetreatUser(BattleID, UserID)
		UpdatedUsers, _ := json.Marshal(Users)

		retreatEvent := createSocketEvent("warrior_retreated", string(UpdatedUsers), UserID)
		m := message{retreatEvent, BattleID}
		h.broadcast <- m

		h.unregister <- sub
		if forceClosed {
			cm := websocket.FormatCloseMessage(4002, "abandoned")
			if err := c.ws.WriteControl(websocket.CloseMessage, cm, time.Now().Add(writeWait)); err != nil {
				b.logger.Ctx(ctx).Error("abandon error", zap.Error(err))
			}
		}
		if err := c.ws.Close(); err != nil {
			b.logger.Ctx(ctx).Error("close error", zap.Error(err))
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
			if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				b.logger.Ctx(ctx).Error("unexpected close error", zap.Error(err))
			}
			break
		}

		keyVal := make(map[string]string)
		err = json.Unmarshal(msg, &keyVal)
		if err != nil {
			badEvent = true
			b.logger.Error("unexpected battle event json error", zap.Error(err))
		}

		eventType := keyVal["type"]
		eventValue := keyVal["value"]

		// confirm leader for any operation that requires it
		if _, ok := leaderOnlyOperations[eventType]; ok && !badEvent {
			err := b.db.ConfirmLeader(BattleID, UserID)
			if err != nil {
				badEvent = true
			}
		}

		// find event handler and execute otherwise invalid event
		if _, ok := b.eventHandlers[eventType]; ok && !badEvent {
			msg, eventErr, forceClosed = b.eventHandlers[eventType](ctx, BattleID, UserID, eventValue)
			if eventErr != nil {
				badEvent = true

				// don't log forceClosed events e.g. Abandon
				if !forceClosed {
					b.logger.Ctx(ctx).Error("close error", zap.Error(eventErr))
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
			if err := c.write(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// handleSocketUnauthorized sets the format close message and closes the websocket
func (b *Service) handleSocketClose(ctx context.Context, ws *websocket.Conn, closeCode int, text string) {
	cm := websocket.FormatCloseMessage(closeCode, text)
	if err := ws.WriteMessage(websocket.CloseMessage, cm); err != nil {
		b.logger.Ctx(ctx).Error("unauthorized close error", zap.Error(err))
	}
	if err := ws.Close(); err != nil {
		b.logger.Ctx(ctx).Error("close error", zap.Error(err))
	}
}

// ServeBattleWs handles websocket requests from the peer.
func (b *Service) ServeBattleWs() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		battleID := vars["battleId"]
		ctx := r.Context()
		var User *model.User
		var UserAuthed bool

		// upgrade to WebSocket connection
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			b.logger.Ctx(ctx).Error("websocket upgrade error", zap.Error(err))
			return
		}
		c := &connection{send: make(chan []byte, 256), ws: ws}

		SessionId, cookieErr := b.validateSessionCookie(w, r)
		if cookieErr != nil && cookieErr.Error() != "NO_SESSION_COOKIE" {
			b.handleSocketClose(ctx, ws, 4001, "unauthorized")
			return
		}

		if SessionId != "" {
			var userErr error
			User, userErr = b.db.GetSessionUser(ctx, SessionId)
			if userErr != nil {
				b.handleSocketClose(ctx, ws, 4001, "unauthorized")
				return
			}
		} else {
			UserID, err := b.validateUserCookie(w, r)
			if err != nil {
				b.handleSocketClose(ctx, ws, 4001, "unauthorized")
				return
			}

			var userErr error
			User, userErr = b.db.GetGuestUser(ctx, UserID)
			if userErr != nil {
				b.handleSocketClose(ctx, ws, 4001, "unauthorized")
				return
			}
		}

		// make sure battle is legit
		battle, battleErr := b.db.GetBattle(battleID, User.Id)
		if battleErr != nil {
			b.handleSocketClose(ctx, ws, 4004, "battle not found")
			return
		}

		// check users battle active status
		UserErr := b.db.GetBattleUserActiveStatus(battleID, User.Id)
		if UserErr != nil && !errors.Is(UserErr, sql.ErrNoRows) {
			usrErrMsg := UserErr.Error()

			if usrErrMsg == "DUPLICATE_BATTLE_USER" {
				b.handleSocketClose(ctx, ws, 4003, "duplicate session")
			} else {
				b.logger.Ctx(ctx).Error("error finding user", zap.Error(UserErr))
				b.handleSocketClose(ctx, ws, 4005, "internal error")
			}
			return
		}

		if battle.JoinCode != "" && (UserErr != nil && errors.Is(UserErr, sql.ErrNoRows)) {
			jcrEvent := createSocketEvent("join_code_required", "", User.Id)
			_ = c.write(websocket.TextMessage, jcrEvent)

			for {
				_, msg, err := c.ws.ReadMessage()
				if err != nil {
					if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
						b.logger.Ctx(ctx).Error("unexpected close error", zap.Error(err))
					}
					break
				}

				keyVal := make(map[string]string)
				err = json.Unmarshal(msg, &keyVal)
				if err != nil {
					b.logger.Error("unexpected battle message error", zap.Error(err))
				}

				if keyVal["type"] == "auth_battle" && keyVal["value"] == battle.JoinCode {
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
			if UserAuthed {
				ss := subscription{c, battleID, User.Id}
				h.register <- ss

				Users, _ := b.db.AddUserToBattle(ss.arena, User.Id)
				UpdatedUsers, _ := json.Marshal(Users)

				Battle, _ := json.Marshal(battle)
				initEvent := createSocketEvent("init", string(Battle), User.Id)
				_ = c.write(websocket.TextMessage, initEvent)

				joinedEvent := createSocketEvent("warrior_joined", string(UpdatedUsers), User.Id)
				m := message{joinedEvent, ss.arena}
				h.broadcast <- m

				go ss.writePump()
				go ss.readPump(b, ctx)

				break
			}
		}
	}
}

// APIEvent handles api driven events into the arena (if active)
func (b *Service) APIEvent(ctx context.Context, arenaID string, UserID, eventType string, eventValue string) error {

	// confirm leader for any operation that requires it
	if _, ok := leaderOnlyOperations[eventType]; ok {
		err := b.db.ConfirmLeader(arenaID, UserID)
		if err != nil {
			return err
		}
	}

	// find event handler and execute otherwise invalid event
	if _, ok := b.eventHandlers[eventType]; ok {
		msg, eventErr, _ := b.eventHandlers[eventType](ctx, arenaID, UserID, eventValue)
		if eventErr != nil {
			return eventErr
		}

		if _, ok := h.arenas[arenaID]; ok {
			m := message{msg, arenaID}
			h.broadcast <- m
		}
	}

	return nil
}
