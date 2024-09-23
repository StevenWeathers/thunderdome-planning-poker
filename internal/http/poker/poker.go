// Package poker provides Poker Game http handlers for Thunderdome
package poker

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

type AuthDataSvc interface {
	GetSessionUser(ctx context.Context, SessionId string) (*thunderdome.User, error)
}

type UserDataSvc interface {
	GetGuestUser(ctx context.Context, UserID string) (*thunderdome.User, error)
}

// Service provides battle service
type Service struct {
	config                Config
	logger                *otelzap.Logger
	validateSessionCookie func(w http.ResponseWriter, r *http.Request) (string, error)
	validateUserCookie    func(w http.ResponseWriter, r *http.Request) (string, error)
	UserService           UserDataSvc
	AuthService           AuthDataSvc
	BattleService         thunderdome.PokerDataSvc
	hub                   *wshub.Hub
}

// New returns a new battle with websocket hub/client and event handlers
func New(
	config Config, logger *otelzap.Logger,
	validateSessionCookie func(w http.ResponseWriter, r *http.Request) (string, error),
	validateUserCookie func(w http.ResponseWriter, r *http.Request) (string, error),
	userService UserDataSvc, authService AuthDataSvc,
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

	b.hub = wshub.NewHub(logger, wshub.Config{
		AppDomain:          config.AppDomain,
		WebsocketSubdomain: config.WebsocketSubdomain,
		WriteWaitSec:       config.WriteWaitSec,
		PongWaitSec:        config.PongWaitSec,
		PingPeriodSec:      config.PingPeriodSec,
	}, map[string]func(context.Context, string, string, string) ([]byte, error, bool){
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
	},
		map[string]struct{}{
			"add_plan":       {},
			"revise_plan":    {},
			"burn_plan":      {},
			"activate_plan":  {},
			"skip_plan":      {},
			"end_voting":     {},
			"finalize_plan":  {},
			"jab_warrior":    {},
			"promote_leader": {},
			"demote_leader":  {},
			"revise_battle":  {},
			"concede_battle": {},
		},
		b.BattleService.ConfirmFacilitator,
		b.RetreatUser,
	)

	go b.hub.Run()

	return b
}
