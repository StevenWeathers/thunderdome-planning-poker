package storyboard

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

// ownerOnlyOperations contains a map of operations that only a storyboard leader can execute
var ownerOnlyOperations = map[string]struct{}{
	"facilitator_add":    {},
	"facilitator_remove": {},
	"edit_storyboard":    {},
	"concede_storyboard": {},
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
	eventHandlers := map[string]func(context.Context, string, string, string) ([]byte, error, bool){
		"add_goal":             b.AddGoal,
		"revise_goal":          b.ReviseGoal,
		"delete_goal":          b.DeleteGoal,
		"add_column":           b.AddColumn,
		"revise_column":        b.ReviseColumn,
		"delete_column":        b.DeleteColumn,
		"add_story":            b.AddStory,
		"update_story_name":    b.UpdateStoryName,
		"update_story_content": b.UpdateStoryContent,
		"update_story_color":   b.UpdateStoryColor,
		"update_story_points":  b.UpdateStoryPoints,
		"update_story_closed":  b.UpdateStoryClosed,
		"update_story_link":    b.UpdateStoryLink,
		"move_story":           b.MoveStory,
		"add_story_comment":    b.AddStoryComment,
		"edit_story_comment":   b.EditStoryComment,
		"delete_story_comment": b.DeleteStoryComment,
		"delete_story":         b.DeleteStory,
		"add_persona":          b.AddPersona,
		"update_persona":       b.UpdatePersona,
		"delete_persona":       b.DeletePersona,
		"facilitator_add":      b.FacilitatorAdd,
		"facilitator_remove":   b.FacilitatorRemove,
		"facilitator_self":     b.FacilitatorSelf,
		"revise_color_legend":  b.ReviseColorLegend,
		"edit_storyboard":      b.EditStoryboard,
		"concede_storyboard":   b.Delete,
		"abandon_storyboard":   b.Abandon,
	}

	var forceClosed bool
	c := sub.conn
	UserID := sub.UserID
	StoryboardID := sub.arena

	defer func() {
		Users := b.db.RetreatStoryboardUser(StoryboardID, UserID)
		UpdatedUsers, _ := json.Marshal(Users)

		retreatEvent := createSocketEvent("user_left", string(UpdatedUsers), UserID)
		m := message{retreatEvent, StoryboardID}
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
			b.logger.Error("unexpected storyboard event json error", zap.Error(err))
		}

		eventType := keyVal["type"]
		eventValue := keyVal["value"]

		// confirm owner for any operation that requires it
		if _, ok := ownerOnlyOperations[eventType]; ok && !badEvent {
			err := b.db.ConfirmStoryboardFacilitator(StoryboardID, UserID)
			if err != nil {
				badEvent = true
			}
		}

		// find event handler and execute otherwise invalid event
		if _, ok := eventHandlers[eventType]; ok && !badEvent {
			msg, eventErr, forceClosed = eventHandlers[eventType](ctx, StoryboardID, UserID, eventValue)
			if eventErr != nil {
				badEvent = true

				// don't log forceClosed events e.g. Abandon
				if !forceClosed {
					b.logger.Ctx(ctx).Error("unexpected close error", zap.Error(eventErr))
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

// ServeWs handles websocket requests from the peer.
func (b *Service) ServeWs() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		storyboardID := vars["storyboardId"]
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

		// make sure storyboard is legit
		storyboard, storyboardErr := b.db.GetStoryboard(storyboardID, User.Id)
		if storyboardErr != nil {
			b.handleSocketClose(ctx, ws, 4004, "storyboard not found")
			return
		}

		// check users storyboard active status
		UserErr := b.db.GetStoryboardUserActiveStatus(storyboardID, User.Id)
		if UserErr != nil && !errors.Is(UserErr, sql.ErrNoRows) {
			usrErrMsg := UserErr.Error()

			if usrErrMsg == "DUPLICATE_STORYBOARD_USER" {
				b.handleSocketClose(ctx, ws, 4003, "duplicate session")
			} else {
				b.logger.Ctx(ctx).Error("error finding user", zap.Error(UserErr))
				b.handleSocketClose(ctx, ws, 4005, "internal error")
			}
			return
		}

		if storyboard.JoinCode != "" && (UserErr != nil && errors.Is(UserErr, sql.ErrNoRows)) {
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
					b.logger.Error("unexpected storyboard message error", zap.Error(err))
				}

				if keyVal["type"] == "auth_storyboard" && keyVal["value"] == storyboard.JoinCode {
					UserAuthed = true
					break
				} else if keyVal["type"] == "auth_storyboard" {
					authIncorrect := createSocketEvent("join_code_incorrect", "", User.Id)
					_ = c.write(websocket.TextMessage, authIncorrect)
				}
			}
		} else {
			UserAuthed = true
		}

		for {
			if UserAuthed {
				ss := subscription{c, storyboardID, User.Id}
				h.register <- ss

				Users, _ := b.db.AddUserToStoryboard(ss.arena, User.Id)
				UpdatedUsers, _ := json.Marshal(Users)

				Storyboard, _ := json.Marshal(storyboard)
				initEvent := createSocketEvent("init", string(Storyboard), User.Id)
				_ = c.write(websocket.TextMessage, initEvent)

				joinedEvent := createSocketEvent("user_joined", string(UpdatedUsers), User.Id)
				m := message{joinedEvent, ss.arena}
				h.broadcast <- m

				go ss.writePump()
				go ss.readPump(b, ctx)

				break
			}
		}
	}
}
