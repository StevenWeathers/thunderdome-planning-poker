package wshub

import (
	"context"

	"github.com/gorilla/websocket"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
)

const (
	// Maximum message size allowed from peer.
	maxMessageSize = 1024 * 1024
)

// Message represents a message sent to the websocket hub.
type Message struct {
	Data []byte `json:"data"`
	Room string `json:"room"`
}

type roomExistsRequest struct {
	room     string
	response chan bool
}

// Hub maintains the set of active connections and broadcasts messages to the connections.
type Hub struct {
	rooms                     map[string]map[Connection]struct{}
	broadcast                 chan Message
	register                  chan Subscription
	unregister                chan Subscription
	roomExists                chan roomExistsRequest
	logger                    *otelzap.Logger
	config                    *Config
	eventHandlers             map[string]func(context.Context, string, string, string) ([]byte, error, bool)
	facilitatorOnlyOperations map[string]struct{}
	confirmFacilitator        func(roomId string, userId string) error
	retreatUser               func(roomId string, userId string) string
}

// NewHub creates a new websocket hub.
func NewHub(
	logger *otelzap.Logger,
	config Config,
	eventHandlers map[string]func(context.Context, string, string, string) ([]byte, error, bool),
	facilitatorOnlyOperations map[string]struct{},
	confirmFacilitator func(roomID string, userID string) error,
	retreatUser func(roomID string, userID string) string,
) *Hub {
	return &Hub{
		broadcast:                 make(chan Message),
		register:                  make(chan Subscription),
		unregister:                make(chan Subscription),
		rooms:                     make(map[string]map[Connection]struct{}),
		roomExists:                make(chan roomExistsRequest),
		logger:                    logger,
		config:                    &config,
		eventHandlers:             eventHandlers,
		facilitatorOnlyOperations: facilitatorOnlyOperations,
		confirmFacilitator:        confirmFacilitator,
		retreatUser:               retreatUser,
	}
}

// Run starts the hub.
func (h *Hub) Run() {
	for {
		select {
		case sub := <-h.register:
			if _, ok := h.rooms[sub.RoomID]; !ok {
				h.rooms[sub.RoomID] = make(map[Connection]struct{})
			}
			h.rooms[sub.RoomID][sub.Conn] = struct{}{}

		case sub := <-h.unregister:
			if _, ok := h.rooms[sub.RoomID]; ok {
				if _, ok := h.rooms[sub.RoomID][sub.Conn]; ok {
					delete(h.rooms[sub.RoomID], sub.Conn)
					sub.Conn.Close()
					if len(h.rooms[sub.RoomID]) == 0 {
						delete(h.rooms, sub.RoomID)
					}
				}
			}

		case m := <-h.broadcast:
			if connections, ok := h.rooms[m.Room]; ok {
				for conn := range connections {
					select {
					case conn.Send() <- m.Data:
					default:
						close(conn.Send())
						delete(connections, conn)
						if len(connections) == 0 {
							delete(h.rooms, m.Room)
						}
					}
				}
			}

		case req := <-h.roomExists:
			_, exists := h.rooms[req.room]
			req.response <- exists
		}
	}
}

// Register adds a subscription to the room.
func (h *Hub) Register(sub Subscription) {
	h.register <- sub
}

// Unregister removes a subscription from the room.
func (h *Hub) Unregister(sub Subscription) {
	h.unregister <- sub
}

// Broadcast sends a message to all connections in the room.
func (h *Hub) Broadcast(msg Message) {
	h.broadcast <- msg
}

// RoomExists checks if a room exists in the hub.
func (h *Hub) RoomExists(room string) bool {
	response := make(chan bool)
	h.roomExists <- roomExistsRequest{room: room, response: response}
	return <-response
}

// NewConnection creates a new websocket connection.
func (h *Hub) NewConnection(ws *websocket.Conn) Connection {
	return Connection{
		send:       make(chan []byte, 256),
		Ws:         ws,
		PingPeriod: h.config.PingPeriod(),
		WriteWait:  h.config.WriteWait(),
		PongWait:   h.config.PongWait(),
	}
}

// NewSubscriber creates a new subscription to the room for the given websocket connection.
func (h *Hub) NewSubscriber(ws *websocket.Conn, userID string, roomID string) Subscription {
	sub := Subscription{
		Conn:   h.NewConnection(ws),
		RoomID: roomID,
		UserID: userID,
	}

	h.Register(sub)

	return sub
}
