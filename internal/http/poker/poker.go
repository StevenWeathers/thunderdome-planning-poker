// Package poker provides Poker Game http handlers for Thunderdome
package poker

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

	// Websocket Subdomain (for Websocket origin check)
	WebsocketSubdomain string
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

type AuthDataSvc interface {
	GetSessionUser(ctx context.Context, SessionId string) (*thunderdome.User, error)
}

// Service provides battle service
type Service struct {
	config                Config
	logger                *otelzap.Logger
	validateSessionCookie func(w http.ResponseWriter, r *http.Request) (string, error)
	validateUserCookie    func(w http.ResponseWriter, r *http.Request) (string, error)
	eventHandlers         map[string]func(context.Context, string, string, string) ([]byte, error, bool)
	UserService           thunderdome.UserDataSvc
	AuthService           AuthDataSvc
	BattleService         thunderdome.PokerDataSvc
}

// New returns a new battle with websocket hub/client and event handlers
func New(
	config Config, logger *otelzap.Logger,
	validateSessionCookie func(w http.ResponseWriter, r *http.Request) (string, error),
	validateUserCookie func(w http.ResponseWriter, r *http.Request) (string, error),
	userService thunderdome.UserDataSvc, authService AuthDataSvc,
	battleService thunderdome.PokerDataSvc,
) *Service {
	b := &Service{
		config:                config,
		logger:                logger,
		validateSessionCookie: validateSessionCookie,
		validateUserCookie:    validateUserCookie,
		UserService:           userService,
		AuthService:           authService,
		BattleService:         battleService,
	}

	b.eventHandlers = map[string]func(context.Context, string, string, string) ([]byte, error, bool){
		"jab_warrior":      b.UserNudge,
		"vote":             b.UserVote,
		"retract_vote":     b.UserVoteRetract,
		"end_voting":       b.StoryVoteEnd,
		"add_plan":         b.StoryAdd,
		"revise_plan":      b.StoryRevise,
		"burn_plan":        b.StoryDelete,
		"story_arrange":    b.StoryArrange,
		"activate_plan":    b.StoryActivate,
		"skip_plan":        b.StorySkip,
		"finalize_plan":    b.StoryFinalize,
		"promote_leader":   b.UserPromote,
		"demote_leader":    b.UserDemote,
		"become_leader":    b.UserPromoteSelf,
		"spectator_toggle": b.UserSpectatorToggle,
		"revise_battle":    b.Revise,
		"concede_battle":   b.Delete,
		"abandon_battle":   b.Abandon,
	}

	go h.run()

	return b
}
