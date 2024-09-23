package wshub

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"unicode/utf8"
)

// EqualASCIIFold returns true if s is equal to t with ASCII case folding as
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

// SocketEvent is the event structure used for socket messages
type SocketEvent struct {
	Type   string `json:"type"`
	Value  string `json:"value"`
	UserID string `json:"userId"`
}

func CreateSocketEvent(Type string, Value string, UserID string) []byte {
	newEvent := &SocketEvent{
		Type:   Type,
		Value:  Value,
		UserID: UserID,
	}

	event, _ := json.Marshal(newEvent)

	return event
}
