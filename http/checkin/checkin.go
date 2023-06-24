package checkin

import (
	"context"
	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"net/http"
)

// Service provides retro service
type Service struct {
	logger                *otelzap.Logger
	validateSessionCookie func(w http.ResponseWriter, r *http.Request) (string, error)
	validateUserCookie    func(w http.ResponseWriter, r *http.Request) (string, error)
	eventHandlers         map[string]func(context.Context, string, string, string) ([]byte, error, bool)
	UserService           thunderdome.UserService
	AuthService           thunderdome.AuthService
	CheckinService        thunderdome.CheckinService
	TeamService           thunderdome.TeamService
}

// New returns a new retro with websocket hub/client and event handlers
func New(
	logger *otelzap.Logger,
	validateSessionCookie func(w http.ResponseWriter, r *http.Request) (string, error),
	validateUserCookie func(w http.ResponseWriter, r *http.Request) (string, error),
	userService thunderdome.UserService, authService thunderdome.AuthService,
	checkinService thunderdome.CheckinService, teamService thunderdome.TeamService,
) *Service {
	c := &Service{
		logger:                logger,
		validateSessionCookie: validateSessionCookie,
		validateUserCookie:    validateUserCookie,
		UserService:           userService,
		AuthService:           authService,
		CheckinService:        checkinService,
		TeamService:           teamService,
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
