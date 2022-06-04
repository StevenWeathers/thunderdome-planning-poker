package retro

import (
	"net/http"

	"github.com/StevenWeathers/thunderdome-planning-poker/db"
	"go.uber.org/zap"
)

// Service provides retro service
type Service struct {
	db                    *db.Database
	logger                *zap.Logger
	validateSessionCookie func(w http.ResponseWriter, r *http.Request) (string, error)
	validateUserCookie    func(w http.ResponseWriter, r *http.Request) (string, error)
}

// New returns a new retro with websocket hub/client and event handlers
func New(
	db *db.Database,
	logger *zap.Logger,
	validateSessionCookie func(w http.ResponseWriter, r *http.Request) (string, error),
	validateUserCookie func(w http.ResponseWriter, r *http.Request) (string, error),
) *Service {
	rs := &Service{
		db:                    db,
		logger:                logger,
		validateSessionCookie: validateSessionCookie,
		validateUserCookie:    validateUserCookie,
	}

	go h.run()

	return rs
}
