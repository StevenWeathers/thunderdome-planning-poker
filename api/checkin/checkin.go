package checkin

import (
	"context"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"net/http"

	"github.com/StevenWeathers/thunderdome-planning-poker/db"
)

// Service provides retro service
type Service struct {
	db                    *db.Database
	logger                *otelzap.Logger
	validateSessionCookie func(w http.ResponseWriter, r *http.Request) (string, error)
	validateUserCookie    func(w http.ResponseWriter, r *http.Request) (string, error)
	eventHandlers         map[string]func(context.Context, string, string, string) ([]byte, error, bool)
}

// New returns a new retro with websocket hub/client and event handlers
func New(
	db *db.Database,
	logger *otelzap.Logger,
	validateSessionCookie func(w http.ResponseWriter, r *http.Request) (string, error),
	validateUserCookie func(w http.ResponseWriter, r *http.Request) (string, error),
) *Service {
	c := &Service{
		db:                    db,
		logger:                logger,
		validateSessionCookie: validateSessionCookie,
		validateUserCookie:    validateUserCookie,
	}

	c.eventHandlers = map[string]func(context.Context, string, string, string) ([]byte, error, bool){
		"checkin_create": c.CheckinCreate,
		"checkin_update": c.CheckinUpdate,
		"checkin_delete": c.CheckinDelete,
		"comment_create": c.CommentCreate,
		"comment_update": c.CommentUpdate,
		"comment_delete": c.CommentDelete,
	}

	go h.run()

	return c
}
