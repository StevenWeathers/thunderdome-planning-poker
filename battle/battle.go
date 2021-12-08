// Package battle provides websocket event handlers for Thunderdome
package battle

import "github.com/StevenWeathers/thunderdome-planning-poker/db"

// Service provides battle service
type Service struct {
	db *db.Database
}

// New returns a new battle with websocket hub/client and event handlers
func New(db *db.Database) *Service {
	b := &Service{
		db: db,
	}

	return b
}
