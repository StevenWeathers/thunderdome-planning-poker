package main

import (
    "log"
    "encoding/json"
    "net/http"
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

type SocketEvent struct {
    EventType string `json:"type"`
    WarriorId string `json:"id"`
    EventValue string `json:"value"`
}

// readPump pumps messages from the websocket connection to the hub.
func (s subscription) readPump() {
    c := s.conn
    defer func() {
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
        warriorId := keyVal["id"]
        battleId := s.arena
        
        switch keyVal["type"] {
        case "join":
            Battles[battleId].Warriors = append(Battles[battleId].Warriors, Warriors[warriorId])
        case "retreat":
            var warriorIndex int
            for i := range Battles[battleId].Warriors {
                if Battles[battleId].Warriors[i].WarriorId == warriorId {
                    warriorIndex = i
                    break;
                }
            }

            Battles[battleId].Warriors = append(Battles[battleId].Warriors[:warriorIndex], Battles[battleId].Warriors[warriorIndex+1:]...)
        case "vote":
            vote := keyVal["value"]
            voteUpdated := false
            for i := range Battles[battleId].Votes {
                if Battles[battleId].Votes[i].WarriorId == warriorId {
                    Battles[battleId].Votes[i].VoteValue = vote
                    voteUpdated = true
                    break;
                }
            }

            if !voteUpdated {
                newVote := &Vote{
                    WarriorId: warriorId,
                    VoteValue: vote }
    
                Battles[battleId].Votes = append(Battles[battleId].Votes, newVote)
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
    ws, err := upgrader.Upgrade(w, r, nil)
    vars := mux.Vars(r)
    log.Println(vars["id"])
    if err != nil {
        log.Println(err)
        return
    }
    c := &connection{send: make(chan []byte, 256), ws: ws}
    s := subscription{c, vars["id"]}
    h.register <- s
    go s.writePump()
    s.readPump()
}