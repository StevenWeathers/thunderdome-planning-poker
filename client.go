package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"time"

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
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// connection is an middleman between the websocket connection and the hub.
type connection struct {
	// The websocket connection.
	ws *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

// SocketEvent is the event structure used for socket messages
type SocketEvent struct {
	EventType  string `json:"type"`
	WarriorID  string `json:"id"`
	EventValue string `json:"value"`
}

// readPump pumps messages from the websocket connection to the hub.
func (s subscription) readPump() {
	c := s.conn
	defer func() {
		log.Println(s.warriorID + " has left the arena")
		var warriorIndex int
		for i := range Battles[s.arena].Warriors {
			if Battles[s.arena].Warriors[i].WarriorID == s.warriorID {
				warriorIndex = i
				break
			}
		}

		Battles[s.arena].Warriors = append(Battles[s.arena].Warriors[:warriorIndex], Battles[s.arena].Warriors[warriorIndex+1:]...)

		joinEvent := &SocketEvent{
			EventType:  "retreat",
			WarriorID:  s.warriorID,
			EventValue: ""}
		event, _ := json.Marshal(joinEvent)
		m := message{event, s.arena}
		h.broadcast <- m

		h.unregister <- s
		c.ws.Close()
	}()
	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, msg, err := c.ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}

		keyVal := make(map[string]string)
		json.Unmarshal(msg, &keyVal) // check for errors
		warriorID := keyVal["id"]
		battleID := s.arena

		switch keyVal["type"] {
		case "vote":
			vote := keyVal["value"]
			voteUpdated := false
			for i := range Battles[battleID].Votes {
				if Battles[battleID].Votes[i].WarriorID == warriorID {
					Battles[battleID].Votes[i].VoteValue = vote
					voteUpdated = true
					break
				}
			}

			if !voteUpdated {
				newVote := &Vote{
					WarriorID: warriorID,
					VoteValue: vote}

				Battles[battleID].Votes = append(Battles[battleID].Votes, newVote)
			}
		default:
		}

		m := message{msg, s.arena}
		h.broadcast <- m
	}
}

// write writes a message with the given message type and payload.
func (c *connection) write(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(mt, payload)
}

// writePump pumps messages from the hub to the websocket connection.
func (s *subscription) writePump() {
	c := s.conn
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
			log.Println(string(message))
			if err := c.write(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

// serveWs handles websocket requests from the peer.
func serveWs(w http.ResponseWriter, r *http.Request) {
	warrior, err := r.Cookie("warrior")
	var warriorID string
	var warriorName string

	if err != nil {
		log.Println("error in reading warrior cookie : " + err.Error() + "\n")
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {
		value, _ := url.PathUnescape(warrior.Value)
		keyVal := make(map[string]string)
		json.Unmarshal([]byte(value), &keyVal) // check for errors
		warriorID = keyVal["id"]
		warriorName = keyVal["name"]

		if Warriors[warriorID] == nil {
			Warriors[warriorID] = &Warrior{
				WarriorID:   warriorID,
				WarriorName: warriorName}
		}
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	vars := mux.Vars(r)
	if err != nil {
		log.Println(err)
		return
	}
	c := &connection{send: make(chan []byte, 256), ws: ws}
	s := subscription{c, vars["id"], warriorID}
	h.register <- s

	Battles[s.arena].Warriors = append(Battles[s.arena].Warriors, Warriors[warriorID])
	joinEvent := &SocketEvent{
		EventType:  "joined",
		WarriorID:  warriorID,
		EventValue: warriorName}
	event, _ := json.Marshal(joinEvent)
	m := message{event, s.arena}
	h.broadcast <- m

	go s.writePump()
	s.readPump()
}
