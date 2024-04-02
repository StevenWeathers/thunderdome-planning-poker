package checkin

import (
	"context"
	"net/http"
	"time"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
)

type Config struct {
	// Time allowed to write a message to the peer.
	WriteWaitSec int

	// Time allowed to read the next pong message from the peer.
	PongWaitSec int

	// Send pings to peer with this period. Must be less than pongWait.
	PingPeriodSec int

	// App Domain (for Websocket origin check)
  AppDomain string
}

func (c *Config) WriteWait() time.Duration {
	return time.Duration(c.WriteWaitSec) * time.Second
}

func (c *Config) PingPeriod() time.Duration {
	return time.Duration(c.PingPeriodSec) * time.Second
}

func (c *Config) PongWait() time.Duration {
	return time.Duration(c.PongWaitSec) * time.Second
}

// Service provides retro service
type Service struct {
	config                Config
	logger                *otelzap.Logger
	validateSessionCookie func(w http.ResponseWriter, r *http.Request) (string, error)
	validateUserCookie    func(w http.ResponseWriter, r *http.Request) (string, error)
	eventHandlers         map[string]func(context.Context, string, string, string) ([]byte, error, bool)
	UserService           thunderdome.UserDataSvc
	AuthService           thunderdome.AuthDataSvc
	CheckinService        thunderdome.CheckinDataSvc
	TeamService           thunderdome.TeamDataSvc
}

// New returns a new retro with websocket hub/client and event handlers
func New(
	config Config,
	logger *otelzap.Logger,
	validateSessionCookie func(w http.ResponseWriter, r *http.Request) (string, error),
	validateUserCookie func(w http.ResponseWriter, r *http.Request) (string, error),
	userService thunderdome.UserDataSvc, authService thunderdome.AuthDataSvc,
	checkinService thunderdome.CheckinDataSvc, teamService thunderdome.TeamDataSvc,
) *Service {
	c := &Service{
		config:                config,
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
