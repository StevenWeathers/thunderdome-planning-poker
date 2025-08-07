package retro

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
	GetSessionUserByID(ctx context.Context, sessionID string) (*thunderdome.User, error)
}

type UserDataSvc interface {
	GetGuestUserByID(ctx context.Context, UserID string) (*thunderdome.User, error)
}

type RetroDataSvc interface {
	EditRetro(retroID string, retroName string, joinCode string, facilitatorCode string, maxVotes int, brainstormVisibility string, phaseAutoAdvance bool) error
	RetroGetByID(retroID string, userID string) (*thunderdome.Retro, error)
	RetroConfirmFacilitator(retroID string, userID string) error
	RetroGetUsers(retroID string) []*thunderdome.RetroUser
	RetroAddUser(retroID string, userID string) ([]*thunderdome.RetroUser, error)
	RetroFacilitatorAdd(retroID string, userID string) ([]string, error)
	RetroFacilitatorRemove(retroID string, userID string) ([]string, error)
	RetroRetreatUser(retroID string, userID string) []*thunderdome.RetroUser
	RetroAbandon(retroID string, userID string) ([]*thunderdome.RetroUser, error)
	RetroAdvancePhase(retroID string, phase string) (*thunderdome.Retro, error)
	RetroDelete(retroID string) error
	GetRetroUserActiveStatus(retroID string, userID string) error
	GetRetroFacilitatorCode(retroID string) (string, error)
	MarkUserReady(retroID string, userID string) ([]string, error)
	UnmarkUserReady(retroID string, userID string) ([]string, error)

	CreateRetroAction(retroID string, userID string, content string) ([]*thunderdome.RetroAction, error)
	UpdateRetroAction(retroID string, actionID string, content string, completed bool) (Actions []*thunderdome.RetroAction, DeleteError error)
	DeleteRetroAction(retroID string, userID string, actionID string) ([]*thunderdome.RetroAction, error)
	RetroActionAssigneeAdd(retroID string, actionID string, userID string) ([]*thunderdome.RetroAction, error)
	RetroActionAssigneeDelete(retroID string, actionID string, userID string) ([]*thunderdome.RetroAction, error)

	CreateRetroItem(retroID string, userID string, itemType string, content string) ([]*thunderdome.RetroItem, error)
	GroupRetroItem(retroID string, itemId string, groupId string) (thunderdome.RetroItem, error)
	DeleteRetroItem(retroID string, userID string, itemType string, itemID string) ([]*thunderdome.RetroItem, error)
	GroupNameChange(retroID string, groupID string, name string) (thunderdome.RetroGroup, error)
	GroupUserVote(retroID string, groupID string, userID string) ([]*thunderdome.RetroVote, error)
	GroupUserSubtractVote(retroID string, groupID string, userID string) ([]*thunderdome.RetroVote, error)
	ItemCommentAdd(retroID string, itemID string, userID string, comment string) ([]*thunderdome.RetroItem, error)
	ItemCommentEdit(retroID string, commentID string, comment string) ([]*thunderdome.RetroItem, error)
	ItemCommentDelete(retroID string, commentID string) ([]*thunderdome.RetroItem, error)
}

type RetroTemplateDataSvc interface {
	// GetTemplateByID retrieves a specific template by its ID
	GetTemplateByID(ctx context.Context, templateID string) (*thunderdome.RetroTemplate, error)
}

type EmailService interface {
	// SendRetroOverview sends the retro overview (items, action items) email to attendees
	SendRetroOverview(retro *thunderdome.Retro, template *thunderdome.RetroTemplate, userName string, userEmail string) error
}

// Service provides retro service
type Service struct {
	config                Config
	logger                *otelzap.Logger
	validateSessionCookie func(w http.ResponseWriter, r *http.Request) (string, error)
	validateUserCookie    func(w http.ResponseWriter, r *http.Request) (string, error)
	UserService           UserDataSvc
	AuthService           AuthDataSvc
	RetroService          RetroDataSvc
	TemplateService       RetroTemplateDataSvc
	EmailService          EmailService
	hub                   *wshub.Hub
}

// New returns a new retro with websocket hub/client and event handlers
func New(
	config Config,
	logger *otelzap.Logger,
	validateSessionCookie func(w http.ResponseWriter, r *http.Request) (string, error),
	validateUserCookie func(w http.ResponseWriter, r *http.Request) (string, error),
	userService UserDataSvc, authService AuthDataSvc,
	retroService RetroDataSvc, templateService RetroTemplateDataSvc,
	emailService EmailService,
) *Service {
	rs := &Service{
		config:                config,
		logger:                logger,
		validateSessionCookie: validateSessionCookie,
		validateUserCookie:    validateUserCookie,
		UserService:           userService,
		AuthService:           authService,
		RetroService:          retroService,
		TemplateService:       templateService,
		EmailService:          emailService,
	}

	rs.hub = wshub.NewHub(logger, wshub.Config{
		AppDomain:          config.AppDomain,
		WebsocketSubdomain: config.WebsocketSubdomain,
		WriteWaitSec:       config.WriteWaitSec,
		PongWaitSec:        config.PongWaitSec,
		PingPeriodSec:      config.PingPeriodSec,
	}, map[string]func(context.Context, string, string, string) (any, []byte, error, bool){
		"create_item":            rs.CreateItem,
		"user_ready":             rs.UserMarkReady,
		"user_unready":           rs.UserUnMarkReady,
		"group_item":             rs.GroupItem,
		"group_name_change":      rs.GroupNameChange,
		"group_vote":             rs.GroupUserVote,
		"group_vote_subtract":    rs.GroupUserSubtractVote,
		"delete_item":            rs.DeleteItem,
		"item_comment_add":       rs.ItemCommentAdd,
		"item_comment_edit":      rs.ItemCommentEdit,
		"item_comment_delete":    rs.ItemCommentDelete,
		"create_action":          rs.CreateAction,
		"update_action":          rs.UpdateAction,
		"delete_action":          rs.DeleteAction,
		"action_assignee_add":    rs.ActionAddAssignee,
		"action_assignee_remove": rs.ActionRemoveAssignee,
		"advance_phase":          rs.AdvancePhase,
		"phase_time_ran_out":     rs.PhaseTimeout,
		"phase_all_ready":        rs.PhaseAllReady,
		"add_facilitator":        rs.FacilitatorAdd,
		"remove_facilitator":     rs.FacilitatorRemove,
		"self_facilitator":       rs.FacilitatorSelf,
		"edit_retro":             rs.EditRetro,
		"concede_retro":          rs.Delete,
		"abandon_retro":          rs.Abandon,
	},
		map[string]struct{}{
			"advance_phase":      {},
			"add_facilitator":    {},
			"remove_facilitator": {},
			"edit_retro":         {},
			"concede_retro":      {},
			"phase_time_ran_out": {},
			"phase_all_ready":    {},
		},
		rs.RetroService.RetroConfirmFacilitator,
		rs.RetreatUser,
	)

	go rs.hub.Run()

	return rs
}
