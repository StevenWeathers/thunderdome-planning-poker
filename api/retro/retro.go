package retro

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
	rs := &Service{
		db:                    db,
		logger:                logger,
		validateSessionCookie: validateSessionCookie,
		validateUserCookie:    validateUserCookie,
	}

	rs.eventHandlers = map[string]func(context.Context, string, string, string) ([]byte, error, bool){
		"create_item":         rs.CreateItem,
		"group_item":          rs.GroupItem,
		"group_name_change":   rs.GroupNameChange,
		"group_vote":          rs.GroupUserVote,
		"group_vote_subtract": rs.GroupUserSubtractVote,
		"delete_item":         rs.DeleteItem,
		"create_action":       rs.CreateAction,
		"update_action":       rs.UpdateAction,
		"delete_action":       rs.DeleteAction,
		"advance_phase":       rs.AdvancePhase,
		"add_facilitator":     rs.FacilitatorAdd,
		"remove_facilitator":  rs.FacilitatorRemove,
		"self_facilitator":    rs.FacilitatorSelf,
		"edit_retro":          rs.EditRetro,
		"concede_retro":       rs.Delete,
		"abandon_retro":       rs.Abandon,
	}

	go h.run()

	return rs
}
