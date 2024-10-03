package wshub

import (
	"time"

	"github.com/gorilla/websocket"
)

// Connection is the websocket connection for a client.
type Connection struct {
	// The websocket connection.
	Ws *websocket.Conn
	// Buffered channel of outbound messages.
	send       chan []byte
	WriteWait  time.Duration
	PingPeriod time.Duration
	PongWait   time.Duration
}

// Send returns the channel to send messages to the client.
func (c *Connection) Send() chan<- []byte { return c.send }

// Close closes the websocket client connection.
func (c *Connection) Close() { c.Ws.Close() }

// Write a message with the given message type and payload.
func (c *Connection) Write(mt int, payload []byte) error {
	_ = c.Ws.SetWriteDeadline(time.Now().Add(c.WriteWait))
	return c.Ws.WriteMessage(mt, payload)
}
