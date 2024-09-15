package checkin

import (
	"context"
	"encoding/json"
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
	maxMessageSize int64 = 1024 * 1024
)

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
	TeamID := sub.arena

	defer func() {
		h.unregister <- sub
		if forceClosed {
			cm := websocket.FormatCloseMessage(4002, "abandoned")
			if err := c.ws.WriteControl(websocket.CloseMessage, cm, time.Now().Add(sub.config.WriteWait())); err != nil {
				b.logger.Ctx(ctx).Error("abandon error", zap.Error(err),
					zap.String("team_id", TeamID), zap.String("session_user_id", UserID))
			}
		}
		if err := c.ws.Close(); err != nil {
			b.logger.Ctx(ctx).Error("close error", zap.Error(err),
				zap.String("team_id", TeamID), zap.String("session_user_id", UserID))
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
				b.logger.Ctx(ctx).Error("unexpected close error", zap.Error(err),
					zap.String("team_id", TeamID), zap.String("session_user_id", UserID))
			}
			break
		}

		keyVal := make(map[string]string)
		err = json.Unmarshal(msg, &keyVal)
		if err != nil {
			badEvent = true
			b.logger.Error("unexpected retro event json error", zap.Error(err),
				zap.String("team_id", TeamID), zap.String("session_user_id", UserID))
		}

		eventType := keyVal["type"]
		eventValue := keyVal["value"]

		// find event handler and execute otherwise invalid event
		if _, ok := b.eventHandlers[eventType]; ok && !badEvent {
			msg, eventErr, forceClosed = b.eventHandlers[eventType](ctx, TeamID, UserID, eventValue)
			if eventErr != nil {
				badEvent = true

				// don't log forceClosed events e.g. Abandon
				if !forceClosed {
					b.logger.Ctx(ctx).Error("unexpected close error", zap.Error(eventErr),
						zap.String("team_id", TeamID), zap.String("session_user_id", UserID),
						zap.String("checkin_event_type", eventType))
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
		teamID := vars["teamId"]
		ctx := r.Context()
		var User *thunderdome.User

		// upgrade to WebSocket connection
		var upgrader = b.createWebsocketUpgrader()
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			b.logger.Ctx(ctx).Error("websocket upgrade error", zap.Error(err),
				zap.String("team_id", teamID))
			return
		}
		c := &connection{config: &b.config, send: make(chan []byte, 256), ws: ws}

		SessionId, cookieErr := b.validateSessionCookie(w, r)
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
			UserID, err := b.validateUserCookie(w, r)
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

		// make sure team is legit
		_, retroErr := b.TeamService.TeamGet(context.Background(), teamID)
		if retroErr != nil {
			b.handleSocketClose(ctx, ws, 4004, "team not found")
			return
		}

		// make sure user is a team user
		_, UserErr := b.TeamService.TeamUserRole(ctx, User.Id, teamID)
		if UserErr != nil {
			b.logger.Ctx(ctx).Error("REQUIRES_TEAM_USER", zap.Error(UserErr),
				zap.String("team_id", teamID), zap.String("session_user_id", User.Id))
			b.handleSocketClose(ctx, ws, 4005, "REQUIRES_TEAM_USER")
			return
		}

		ss := subscription{&b.config, c, teamID, User.Id}
		h.register <- ss

		initEvent := createSocketEvent("init", "", User.Id)
		_ = c.write(websocket.TextMessage, initEvent)

		go ss.writePump()
		go ss.readPump(b, ctx)
	}
}

// APIEvent handles api driven events into the arena (if active)
func (b *Service) APIEvent(ctx context.Context, arenaID string, UserID, eventType string, eventValue string) error {
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
