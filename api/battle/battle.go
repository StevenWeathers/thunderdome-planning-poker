// Package battle provides websocket event handlers for Thunderdome
package battle

import (
	"context"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"net/http"

	"github.com/StevenWeathers/thunderdome-planning-poker/db"
)

// Service provides battle service
type Service struct {
	db                    *db.Database
	logger                *otelzap.Logger
	validateSessionCookie func(w http.ResponseWriter, r *http.Request) (string, error)
	validateUserCookie    func(w http.ResponseWriter, r *http.Request) (string, error)
	eventHandlers         map[string]func(context.Context, string, string, string) ([]byte, error, bool)
}

// New returns a new battle with websocket hub/client and event handlers
func New(
	db *db.Database,
	logger *otelzap.Logger,
	validateSessionCookie func(w http.ResponseWriter, r *http.Request) (string, error),
	validateUserCookie func(w http.ResponseWriter, r *http.Request) (string, error),
) *Service {
	b := &Service{
		db:                    db,
		logger:                logger,
		validateSessionCookie: validateSessionCookie,
		validateUserCookie:    validateUserCookie,
	}

	b.eventHandlers = map[string]func(context.Context, string, string, string) ([]byte, error, bool){
		"jab_warrior":      b.UserNudge,
		"vote":             b.UserVote,
		"retract_vote":     b.UserVoteRetract,
		"end_voting":       b.PlanVoteEnd,
		"add_plan":         b.PlanAdd,
		"revise_plan":      b.PlanRevise,
		"burn_plan":        b.PlanDelete,
		"activate_plan":    b.PlanActivate,
		"skip_plan":        b.PlanSkip,
		"finalize_plan":    b.PlanFinalize,
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
