package checkin

import (
	"context"
	"net/http"

	"github.com/StevenWeathers/thunderdome-planning-poker/internal/wshub"
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

	// Websocket Subdomain (for Websocket origin check)
	WebsocketSubdomain string
}

type CheckinDataSvc interface {
	CheckinList(ctx context.Context, teamID string, date string, timeZone string) ([]*thunderdome.TeamCheckin, error)
	CheckinCreate(ctx context.Context, teamID string, userID string, yesterday string, today string, blockers string, discuss string, goalsMet bool) error
	CheckinUpdate(ctx context.Context, checkinID string, yesterday string, today string, blockers string, discuss string, goalsMet bool) error
	CheckinDelete(ctx context.Context, checkinID string) error
	CheckinComment(ctx context.Context, teamID string, checkinID string, userID string, comment string) error
	CheckinCommentEdit(ctx context.Context, teamID string, userID string, commentID string, comment string) error
	CheckinCommentDelete(ctx context.Context, commentID string) error
	CheckinLastByUser(ctx context.Context, teamID string, userID string) (*thunderdome.TeamCheckin, error)
}

type AuthDataSvc interface {
	GetSessionUserByID(ctx context.Context, sessionID string) (*thunderdome.User, error)
}

type TeamDataSvc interface {
	TeamUserRoleByUserID(ctx context.Context, userID string, teamID string) (string, error)
	TeamGetByID(ctx context.Context, teamID string) (*thunderdome.Team, error)
}

type UserDataSvc interface {
	GetGuestUserByID(ctx context.Context, userID string) (*thunderdome.User, error)
}

// Service provides retro service
type Service struct {
	config                Config
	logger                *otelzap.Logger
	validateSessionCookie func(w http.ResponseWriter, r *http.Request) (string, error)
	validateUserCookie    func(w http.ResponseWriter, r *http.Request) (string, error)
	UserService           UserDataSvc
	AuthService           AuthDataSvc
	CheckinService        CheckinDataSvc
	TeamService           TeamDataSvc
	hub                   *wshub.Hub
}

// New returns a new retro with websocket hub/client and event handlers
func New(
	config Config,
	logger *otelzap.Logger,
	validateSessionCookie func(w http.ResponseWriter, r *http.Request) (string, error),
	validateUserCookie func(w http.ResponseWriter, r *http.Request) (string, error),
	userService UserDataSvc, authService AuthDataSvc,
	checkinService CheckinDataSvc, teamService TeamDataSvc,
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

	c.hub = wshub.NewHub(logger, wshub.Config{
		AppDomain:          config.AppDomain,
		WebsocketSubdomain: config.WebsocketSubdomain,
		WriteWaitSec:       config.WriteWaitSec,
		PongWaitSec:        config.PongWaitSec,
		PingPeriodSec:      config.PingPeriodSec,
	}, map[string]func(context.Context, string, string, string) (any, []byte, error, bool){
		"checkin_create": c.CheckinCreate,
		"checkin_update": c.CheckinUpdate,
		"checkin_delete": c.CheckinDelete,
		"comment_create": c.CommentCreate,
		"comment_update": c.CommentUpdate,
		"comment_delete": c.CommentDelete,
	},
		map[string]struct{}{},
		nil,
		nil,
	)

	go c.hub.Run()

	return c
}
