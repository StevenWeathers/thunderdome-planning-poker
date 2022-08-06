package checkin

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// Message holds the message to send to connections
type Message struct {
	arena string
	msg   string
}

// Subscription holds the users sse connection
type Subscription struct {
	send  chan Message
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
func New() *Broker {
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
		send:  make(chan Message),
		arena: TeamId,
	}
	b.register <- s

	// send connected message to initiate
	fmt.Fprintf(w, "data: connected\n\n")
	flusher.Flush()

	// send ping every 60 seconds to keep connection alive in browser
	ticker := time.NewTicker(60 * time.Second)

	// unregister connection upon close
	defer func() {
		b.unregister <- s
		ticker.Stop()
	}()

	for {
		select {
		case m, ok := <-s.send:
			if !ok {
				return
			}
			fmt.Fprintf(w, "data: %s\n\n", m.msg)
			flusher.Flush()
		case <-ticker.C:
			fmt.Fprintf(w, "data: ping\n\n")
			flusher.Flush()
		case <-ctx.Done():
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
			b.arenas[a.arena][a.send] = struct{}{}
		case a := <-b.unregister:
			connections := b.arenas[a.arena]
			if connections != nil {
				if _, ok := connections[a.send]; ok {
					delete(connections, a.send)
					close(a.send)
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
func (broker *Broker) BroadcastMessage(arena string, msg string) {
	m := Message{
		arena: arena,
		msg:   msg,
	}

	b.broadcast <- m
}
