// Package battle provides websocket event handlers for Thunderdome
package battle

import (
	"github.com/StevenWeathers/thunderdome-planning-poker/db"
	"net/http"
)

// Service provides battle service
type Service struct {
	db                    *db.Database
	validateSessionCookie func(w http.ResponseWriter, r *http.Request) (string, error)
	validateUserCookie    func(w http.ResponseWriter, r *http.Request) (string, error)
}

// New returns a new battle with websocket hub/client and event handlers
func New(
	db *db.Database,
	validateSessionCookie func(w http.ResponseWriter, r *http.Request) (string, error),
	validateUserCookie func(w http.ResponseWriter, r *http.Request) (string, error),
) *Service {
	b := &Service{
		db:                    db,
		validateSessionCookie: validateSessionCookie,
		validateUserCookie:    validateUserCookie,
	}

	go h.run()

	return b
}
