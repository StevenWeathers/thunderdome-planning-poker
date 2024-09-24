package wshub

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type Subscription struct {
	Conn   Connection
	RoomID string
	UserID string
}

// WritePump pumps messages from the Hub to the websocket connection.
func (s *Subscription) WritePump() {
	ticker := time.NewTicker(s.Conn.PingPeriod)
	defer func() {
		ticker.Stop()
		_ = s.Conn.Ws.Close()
	}()
	for {
		select {
		case message, ok := <-s.Conn.send:
			if !ok {
				_ = s.Conn.Write(websocket.CloseMessage, []byte{})
				return
			}
			if err := s.Conn.Write(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			if err := s.Conn.Write(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (s *Subscription) ReadPump(
	ctx context.Context,
	hub *Hub,
) {
	ctx = context.WithoutCancel(ctx)
	forceClosed := false
	defer func() {
		var UpdatedUsers string
		if hub.retreatUser != nil {
			UpdatedUsers = hub.retreatUser(s.RoomID, s.UserID)
		}

		if forceClosed {
			cm := websocket.FormatCloseMessage(4002, "abandoned")
			if err := s.Conn.Ws.WriteControl(websocket.CloseMessage, cm, time.Now().Add(s.Conn.WriteWait)); err != nil {
				hub.logger.Ctx(ctx).Error("abandon error", zap.Error(err),
					zap.String("room_id", s.RoomID), zap.String("session_user_id", s.UserID))
			}
		}
		_ = s.Conn.Ws.Close() // close connection, don't care about error in attempting to close unclosed connection

		hub.Unregister(*s)

		if hub.retreatUser != nil {
			userLeaveEvent := CreateSocketEvent("user_left", UpdatedUsers, s.UserID)
			if hub.RoomExists(s.RoomID) {
				hub.Broadcast(Message{Data: userLeaveEvent, Room: s.RoomID})
			}
		}
	}()

	s.Conn.Ws.SetReadLimit(maxMessageSize)
	_ = s.Conn.Ws.SetReadDeadline(time.Now().Add(s.Conn.PongWait))
	s.Conn.Ws.SetPongHandler(func(string) error {
		_ = s.Conn.Ws.SetReadDeadline(time.Now().Add(s.Conn.PongWait))
		return nil
	})

	// Read messages from the websocket connection.
	for {
		var badEvent bool
		var eventErr error
		_, msg, err := s.Conn.Ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(
				err, websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseAbnormalClosure,
			) {
				hub.logger.Ctx(ctx).Error("unexpected close error", zap.Error(err),
					zap.String("room_id", s.RoomID), zap.String("session_user_id", s.UserID))
			}
			break
		}

		keyVal := make(map[string]string)
		err = json.Unmarshal(msg, &keyVal)
		if err != nil {
			badEvent = true
			hub.logger.Error("unexpected room event json error", zap.Error(err),
				zap.String("room_id", s.RoomID), zap.String("session_user_id", s.UserID))
		}

		eventType := keyVal["type"]
		eventValue := keyVal["value"]

		// confirm leader for any operation that requires it (if the room requires)
		if hub.confirmFacilitator != nil {
			if _, ok := hub.facilitatorOnlyOperations[eventType]; ok && !badEvent {
				err := hub.confirmFacilitator(s.RoomID, s.UserID)
				if err != nil {
					badEvent = true
				}
			}
		}

		// find event handler and execute otherwise invalid event
		if _, ok := hub.eventHandlers[eventType]; ok && !badEvent {
			msg, eventErr, forceClosed = hub.eventHandlers[eventType](ctx, s.RoomID, s.UserID, eventValue)
			if eventErr != nil {
				badEvent = true

				// don't log forceClosed events e.g. Abandon
				if !forceClosed {
					hub.logger.Ctx(ctx).Error("close error", zap.Error(eventErr),
						zap.String("room_id", s.RoomID), zap.String("session_user_id", s.UserID),
						zap.String("room_event_type", eventType))
				}
			}
		}

		if !badEvent && hub.RoomExists(s.RoomID) {
			hub.Broadcast(Message{Data: msg, Room: s.RoomID})
		}

		if forceClosed {
			break
		}
	}
}
