// Package poker provides Poker Game http handlers for Thunderdome
package poker

import (
	"context"
	"net/http"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
)

// Service provides battle service
type Service struct {
	logger                *otelzap.Logger
	validateSessionCookie func(w http.ResponseWriter, r *http.Request) (string, error)
	validateUserCookie    func(w http.ResponseWriter, r *http.Request) (string, error)
	eventHandlers         map[string]func(context.Context, string, string, string) ([]byte, error, bool)
	UserService           thunderdome.UserDataSvc
	AuthService           thunderdome.AuthDataSvc
	BattleService         thunderdome.PokerDataSvc
}

// New returns a new battle with websocket hub/client and event handlers
func New(
	logger *otelzap.Logger,
	validateSessionCookie func(w http.ResponseWriter, r *http.Request) (string, error),
	validateUserCookie func(w http.ResponseWriter, r *http.Request) (string, error),
	userService thunderdome.UserDataSvc, authService thunderdome.AuthDataSvc,
	battleService thunderdome.PokerDataSvc,
) *Service {
	b := &Service{
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
