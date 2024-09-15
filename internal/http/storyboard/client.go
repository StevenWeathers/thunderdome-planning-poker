package storyboard

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"
	"unicode/utf8"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"

	"go.uber.org/zap"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

const (
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

// connection is a middleman between the websocket connection and the hub.
type connection struct {
	config *Config
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
	StoryboardID := sub.arena

	defer func() {
		Users := b.StoryboardService.RetreatStoryboardUser(StoryboardID, UserID)
		UpdatedUsers, _ := json.Marshal(Users)

		retreatEvent := createSocketEvent("user_left", string(UpdatedUsers), UserID)
		m := message{retreatEvent, StoryboardID}
		h.broadcast <- m

		h.unregister <- sub
		if forceClosed {
			cm := websocket.FormatCloseMessage(4002, "abandoned")
			if err := c.ws.WriteControl(websocket.CloseMessage, cm, time.Now().Add(sub.config.WriteWait())); err != nil {
				b.Logger.Ctx(ctx).Error("abandon error", zap.Error(err),
					zap.String("session_user_id", UserID), zap.String("storyboard_id", StoryboardID))
			}
		}
		if err := c.ws.Close(); err != nil {
			b.Logger.Ctx(ctx).Error("close error", zap.Error(err),
				zap.String("session_user_id", UserID), zap.String("storyboard_id", StoryboardID))
		}
	}()
	c.ws.SetReadLimit(maxMessageSize)
	_ = c.ws.SetReadDeadline(time.Now().Add(sub.config.PongWait()))
	c.ws.SetPongHandler(func(string) error {
		_ = c.ws.SetReadDeadline(time.Now().Add(sub.config.PongWait()))
		return nil
	})

	for {
		var badEvent bool
		var eventErr error
		_, msg, err := c.ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				b.Logger.Ctx(ctx).Error("unexpected close error", zap.Error(err),
					zap.String("session_user_id", UserID), zap.String("storyboard_id", StoryboardID))
			}
			break
		}

		keyVal := make(map[string]string)
		err = json.Unmarshal(msg, &keyVal)
		if err != nil {
			badEvent = true
			b.Logger.Error("unexpected storyboard event json error", zap.Error(err),
				zap.String("session_user_id", UserID), zap.String("storyboard_id", StoryboardID))
		}

		eventType := keyVal["type"]
		eventValue := keyVal["value"]

		// confirm owner for any operation that requires it
		if _, ok := ownerOnlyOperations[eventType]; ok && !badEvent {
			err := b.StoryboardService.ConfirmStoryboardFacilitator(StoryboardID, UserID)
			if err != nil {
				badEvent = true
			}
		}

		// find event handler and execute otherwise invalid event
		if _, ok := b.EventHandlers[eventType]; ok && !badEvent {
			msg, eventErr, forceClosed = b.EventHandlers[eventType](ctx, StoryboardID, UserID, eventValue)
			if eventErr != nil {
				badEvent = true

				// don't log forceClosed events e.g. Abandon
				if !forceClosed {
					b.Logger.Ctx(ctx).Error("unexpected close error", zap.Error(eventErr),
						zap.String("session_user_id", UserID), zap.String("storyboard_id", StoryboardID),
						zap.String("storyboard_event_type", eventType))
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
	_ = c.ws.SetWriteDeadline(time.Now().Add(c.config.WriteWait()))
	return c.ws.WriteMessage(mt, payload)
}

// writePump pumps messages from the hub to the websocket connection.
func (sub *subscription) writePump() {
	c := sub.conn
	ticker := time.NewTicker(sub.config.PingPeriod())
	defer func() {
		ticker.Stop()
		_ = c.ws.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				_ = c.write(websocket.CloseMessage, []byte{})
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

func (b *Service) createWebsocketUpgrader() websocket.Upgrader {
	return websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return checkOrigin(r, b.config.AppDomain, b.config.WebsocketSubdomain)
		},
	}
}

func checkOrigin(r *http.Request, appDomain string, subDomain string) bool {
	origin := r.Header.Get("Origin")
	if len(origin) == 0 {
		return true
	}
	originUrl, err := url.Parse(origin)
	if err != nil {
		return false
	}
	appDomainCheck := equalASCIIFold(originUrl.Host, appDomain)
	subDomainCheck := equalASCIIFold(originUrl.Host, fmt.Sprintf("%s.%s", subDomain, appDomain))
	hostCheck := equalASCIIFold(originUrl.Host, r.Host)

	return appDomainCheck || subDomainCheck || hostCheck
}

// equalASCIIFold returns true if s is equal to t with ASCII case folding as
// defined in RFC 4790.
// Taken from Gorilla Websocket, https://github.com/gorilla/websocket/blob/main/util.go
func equalASCIIFold(s, t string) bool {
	for s != "" && t != "" {
		sr, size := utf8.DecodeRuneInString(s)
		s = s[size:]
		tr, size := utf8.DecodeRuneInString(t)
		t = t[size:]
		if sr == tr {
			continue
		}
		if 'A' <= sr && sr <= 'Z' {
			sr = sr + 'a' - 'A'
		}
		if 'A' <= tr && tr <= 'Z' {
			tr = tr + 'a' - 'A'
		}
		if sr != tr {
			return false
		}
	}
	return s == t
}

// handleSocketUnauthorized sets the format close message and closes the websocket
func (b *Service) handleSocketClose(ctx context.Context, ws *websocket.Conn, closeCode int, text string) {
	cm := websocket.FormatCloseMessage(closeCode, text)
	if err := ws.WriteMessage(websocket.CloseMessage, cm); err != nil {
		b.Logger.Ctx(ctx).Error("unauthorized close error", zap.Error(err))
	}
	if err := ws.Close(); err != nil {
		b.Logger.Ctx(ctx).Error("close error", zap.Error(err))
	}
}

// ServeWs handles websocket requests from the peer.
func (b *Service) ServeWs() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		storyboardID := vars["storyboardId"]
		ctx := r.Context()
		var User *thunderdome.User
		var UserAuthed bool

		// upgrade to WebSocket connection
		var upgrader = b.createWebsocketUpgrader()
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			b.Logger.Ctx(ctx).Error("websocket upgrade error", zap.Error(err),
				zap.String("storyboard_id", storyboardID))
			return
		}
		c := &connection{config: &b.config, send: make(chan []byte, 256), ws: ws}

		SessionId, cookieErr := b.ValidateSessionCookie(w, r)
		if cookieErr != nil && cookieErr.Error() != "COOKIE_NOT_FOUND" {
			b.handleSocketClose(ctx, ws, 4001, "unauthorized")
			return
		}

		if SessionId != "" {
			var userErr error
			User, userErr = b.AuthService.GetSessionUser(ctx, SessionId)
			if userErr != nil {
				b.handleSocketClose(ctx, ws, 4001, "unauthorized")
				return
			}
		} else {
			UserID, err := b.ValidateUserCookie(w, r)
			if err != nil {
				b.handleSocketClose(ctx, ws, 4001, "unauthorized")
				return
			}

			var userErr error
			User, userErr = b.UserService.GetGuestUser(ctx, UserID)
			if userErr != nil {
				b.handleSocketClose(ctx, ws, 4001, "unauthorized")
				return
			}
		}

		// make sure storyboard is legit
		storyboard, storyboardErr := b.StoryboardService.GetStoryboard(storyboardID, User.Id)
		if storyboardErr != nil {
			b.handleSocketClose(ctx, ws, 4004, "storyboard not found")
			return
		}

		// check users storyboard active status
		UserErr := b.StoryboardService.GetStoryboardUserActiveStatus(storyboardID, User.Id)
		if UserErr != nil && !errors.Is(UserErr, sql.ErrNoRows) {
			usrErrMsg := UserErr.Error()

			if usrErrMsg == "DUPLICATE_STORYBOARD_USER" {
				b.handleSocketClose(ctx, ws, 4003, "duplicate session")
			} else {
				b.Logger.Ctx(ctx).Error("error finding user", zap.Error(UserErr),
					zap.String("storyboard_id", storyboardID), zap.String("session_user_id", User.Id))
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
						b.Logger.Ctx(ctx).Error("unexpected close error", zap.Error(err),
							zap.String("storyboard_id", storyboardID), zap.String("session_user_id", User.Id))
					}
					break
				}

				keyVal := make(map[string]string)
				err = json.Unmarshal(msg, &keyVal)
				if err != nil {
					b.Logger.Error("unexpected storyboard message error", zap.Error(err),
						zap.String("storyboard_id", storyboardID), zap.String("session_user_id", User.Id))
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

		if UserAuthed {
			ss := subscription{&b.config, c, storyboardID, User.Id}
			h.register <- ss

			Users, _ := b.StoryboardService.AddUserToStoryboard(ss.arena, User.Id)
			UpdatedUsers, _ := json.Marshal(Users)

			Storyboard, _ := json.Marshal(storyboard)
			initEvent := createSocketEvent("init", string(Storyboard), User.Id)
			_ = c.write(websocket.TextMessage, initEvent)

			joinedEvent := createSocketEvent("user_joined", string(UpdatedUsers), User.Id)
			m := message{joinedEvent, ss.arena}
			h.broadcast <- m

			go ss.writePump()
			go ss.readPump(b, ctx)
		}
	}
}

// APIEvent handles api driven events into the arena (if active)
func (b *Service) APIEvent(ctx context.Context, arenaID string, UserID, eventType string, eventValue string) error {
	// confirm leader for any operation that requires it
	if _, ok := ownerOnlyOperations[eventType]; ok {
		err := b.StoryboardService.ConfirmStoryboardFacilitator(arenaID, UserID)
		if err != nil {
			return err
		}
	}

	// find event handler and execute otherwise invalid event
	if _, ok := b.EventHandlers[eventType]; ok {
		msg, eventErr, _ := b.EventHandlers[eventType](ctx, arenaID, UserID, eventValue)
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
