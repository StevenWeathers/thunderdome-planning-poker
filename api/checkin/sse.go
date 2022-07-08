package checkin

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

// Message holds the message to send to connections
type Message struct {
	arena   string `json:"arena"`
	Message string `json:"msg"`
}

// Subscription holds the users sse connection
type Subscription struct {
	conn  chan Message
	arena string
}

// Broker maintains the set of active connections and broadcasts messages to the
// connections.
type Broker struct {
	// Registered connections.
	arenas map[string]map[chan Message]struct{}
	// Inbound messages from the connections.
	broadcast chan Message
	// Register requests from the connections.
	register chan Subscription
	// Unregister requests from connections.
	unregister chan Subscription
}

var b = &Broker{
	broadcast:  make(chan Message),
	register:   make(chan Subscription),
	unregister: make(chan Subscription),
	arenas:     make(map[string]map[chan Message]struct{}),
}

// New returns a new checkin with sse broker
func New() (broker *Broker) {
	go b.listen()

	return b
}

// Stream handles the server sent event connection with the browser
func (broker *Broker) Stream(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	TeamId := vars["teamId"]
	flusher, ok := w.(http.Flusher)

	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	s := Subscription{
		conn:  make(chan Message),
		arena: TeamId,
	}
	b.register <- s

	// unregister connection upon close
	defer func() {
		b.unregister <- s
	}()

	// send connected message to initiate
	fmt.Fprintf(w, "data: connected\n\n")
	flusher.Flush()

	// send ping every 30 seconds to keep connection alive in browser
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case msg := <-s.conn:
			fmt.Fprintf(w, "data: %s\n\n", msg.Message)
			flusher.Flush()
		case <-ticker.C:
			fmt.Fprintf(w, "data: ping\n\n")
			flusher.Flush()
		case <-ctx.Done():
			b.unregister <- s
			return
		}
	}
}

func (broker *Broker) listen() {
	for {
		select {
		case a := <-b.register:
			connections := b.arenas[a.arena]
			if connections == nil {
				connections = make(map[chan Message]struct{})
				b.arenas[a.arena] = connections
			}
			b.arenas[a.arena][a.conn] = struct{}{}
		case a := <-b.unregister:
			connections := b.arenas[a.arena]
			if connections != nil {
				if _, ok := connections[a.conn]; ok {
					delete(connections, a.conn)
					close(a.conn)
					if len(connections) == 0 {
						delete(b.arenas, a.arena)
					}
				}
			}
		case m := <-b.broadcast:
			connections := b.arenas[m.arena]
			for c := range connections {
				select {
				case c <- m:
				default:
					close(c)
					delete(connections, c)
					if len(connections) == 0 {
						delete(b.arenas, m.arena)
					}
				}
			}
		}
	}

}

// BroadcastMessage sends a message to the broker for broadcasting to the connections
func (broker *Broker) BroadcastMessage(arena string, Msg string) {
	msg := Message{
		arena:   arena,
		Message: Msg,
	}

	b.broadcast <- msg
}
