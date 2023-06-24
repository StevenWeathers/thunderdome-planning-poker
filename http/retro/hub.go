package retro

type message struct {
	data  []byte
	arena string
}

type subscription struct {
	conn   *connection
	arena  string
	UserID string
}

// hub maintains the set of active connections and broadcasts messages to the
// connections.
type hub struct {
	// Registered connections.
	arenas map[string]map[*connection]struct{}

	// Inbound messages from the connections.
	broadcast chan message

	// Register requests from the connections.
	register chan subscription

	// Unregister requests from connections.
	unregister chan subscription
}

var h = hub{
	broadcast:  make(chan message),
	register:   make(chan subscription),
	unregister: make(chan subscription),
	arenas:     make(map[string]map[*connection]struct{}),
}

func (h *hub) run() {
	for {
		select {
		case a := <-h.register:
			connections := h.arenas[a.arena]
			if connections == nil {
				connections = make(map[*connection]struct{})
				h.arenas[a.arena] = connections
			}
			h.arenas[a.arena][a.conn] = struct{}{}
		case a := <-h.unregister:
			connections := h.arenas[a.arena]
			if connections != nil {
				if _, ok := connections[a.conn]; ok {
					delete(connections, a.conn)
					close(a.conn.send)
					if len(connections) == 0 {
						delete(h.arenas, a.arena)
					}
				}
			}
		case m := <-h.broadcast:
			connections := h.arenas[m.arena]
			for c := range connections {
				select {
				case c.send <- m.data:
				default:
					close(c.send)
					delete(connections, c)
					if len(connections) == 0 {
						delete(h.arenas, m.arena)
					}
				}
			}
		}
	}
}
