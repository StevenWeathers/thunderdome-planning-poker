package wshub

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

// CreateWebsocketUpgrader creates a websocket.Upgrader with the given AppDomain and WebsocketSubdomain
func (h *Hub) CreateWebsocketUpgrader() websocket.Upgrader {
	return websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return checkOrigin(r, h.config.AppDomain, h.config.WebsocketSubdomain)
		},
	}
}

// HandleSocketClose sets the format close message and closes the websocket
func (h *Hub) HandleSocketClose(ctx context.Context, ws *websocket.Conn, closeCode int, text string) {
	cm := websocket.FormatCloseMessage(closeCode, text)
	if err := ws.WriteMessage(websocket.CloseMessage, cm); err != nil {
		h.logger.Ctx(ctx).Error("unauthorized close error", zap.Error(err))
	}
	if err := ws.Close(); err != nil {
		h.logger.Ctx(ctx).Error("close error", zap.Error(err))
	}
}

// AuthError is a custom error type for handling websocket authentication errors
type AuthError struct {
	Code    int
	Message string
}

// Error returns the error message
func (e *AuthError) Error() string {
	return e.Message
}

// WebSocketHandler creates a http.HandlerFunc for handling WebSocket connections
func (h *Hub) WebSocketHandler(
	roomIDVar string,
	authFunc func(w http.ResponseWriter, r *http.Request, c *Connection, roomID string) *AuthError,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		vars := mux.Vars(r)
		RoomID := vars[roomIDVar]

		// upgrade to WebSocket connection
		var upgrader = h.CreateWebsocketUpgrader()
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			h.logger.Ctx(ctx).Error("websocket upgrade error", zap.Error(err),
				zap.String("room_id", RoomID))
			return
		}
		c := h.NewConnection(ws)

		authErr := authFunc(w, r, &c, RoomID)
		if authErr != nil {
			h.HandleSocketClose(ctx, c.Ws, authErr.Code, authErr.Error())
			return
		}
	}
}

// ProcessAPIEventHandler processes an event from the API through the websocket hub.
func (h *Hub) ProcessAPIEventHandler(ctx context.Context, userID, roomID, eventType string, eventValue string) error {
	// find event handler and execute otherwise invalid event
	if _, ok := h.eventHandlers[eventType]; ok {
		// confirm leader for any operation that requires it
		if h.confirmFacilitator != nil {
			if _, ok := h.facilitatorOnlyOperations[eventType]; ok {
				err := h.confirmFacilitator(roomID, userID)
				if err != nil {
					return err
				}
			}
		}

		msg, eventErr, _ := h.eventHandlers[eventType](ctx, roomID, userID, eventValue)
		if eventErr != nil {
			return eventErr
		}

		if h.RoomExists(roomID) {
			h.Broadcast(Message{Data: msg, Room: roomID})
		}
	}

	return nil
}
