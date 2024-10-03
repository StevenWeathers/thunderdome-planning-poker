package storyboard

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
	GetSessionUser(ctx context.Context, sessionID string) (*thunderdome.User, error)
}

type UserDataSvc interface {
	GetGuestUser(ctx context.Context, userID string) (*thunderdome.User, error)
}

// Service provides storyboard service
type Service struct {
	config                Config
	logger                *otelzap.Logger
	validateSessionCookie func(w http.ResponseWriter, r *http.Request) (string, error)
	validateUserCookie    func(w http.ResponseWriter, r *http.Request) (string, error)
	UserService           UserDataSvc
	AuthService           AuthDataSvc
	StoryboardService     thunderdome.StoryboardDataSvc
	hub                   *wshub.Hub
}

// New returns a new storyboard with websocket hub/client and event handlers
func New(
	config Config,
	logger *otelzap.Logger,
	validateSessionCookie func(w http.ResponseWriter, r *http.Request) (string, error),
	validateUserCookie func(w http.ResponseWriter, r *http.Request) (string, error),
	userService UserDataSvc, authService AuthDataSvc,
	storyboardService thunderdome.StoryboardDataSvc,
) *Service {
	sb := &Service{
		config:                config,
		logger:                logger,
		validateSessionCookie: validateSessionCookie,
		validateUserCookie:    validateUserCookie,
		UserService:           userService,
		AuthService:           authService,
		StoryboardService:     storyboardService,
	}

	sb.hub = wshub.NewHub(logger, wshub.Config{
		AppDomain:          config.AppDomain,
		WebsocketSubdomain: config.WebsocketSubdomain,
		WriteWaitSec:       config.WriteWaitSec,
		PongWaitSec:        config.PongWaitSec,
		PingPeriodSec:      config.PingPeriodSec,
	}, map[string]func(context.Context, string, string, string) ([]byte, error, bool){
		"add_goal":              sb.AddGoal,
		"revise_goal":           sb.ReviseGoal,
		"delete_goal":           sb.DeleteGoal,
		"add_column":            sb.AddColumn,
		"revise_column":         sb.ReviseColumn,
		"delete_column":         sb.DeleteColumn,
		"column_persona_add":    sb.ColumnPersonaAdd,
		"column_persona_remove": sb.ColumnPersonaRemove,
		"add_story":             sb.AddStory,
		"update_story_name":     sb.UpdateStoryName,
		"update_story_content":  sb.UpdateStoryContent,
		"update_story_color":    sb.UpdateStoryColor,
		"update_story_points":   sb.UpdateStoryPoints,
		"update_story_closed":   sb.UpdateStoryClosed,
		"update_story_link":     sb.UpdateStoryLink,
		"move_story":            sb.MoveStory,
		"add_story_comment":     sb.AddStoryComment,
		"edit_story_comment":    sb.EditStoryComment,
		"delete_story_comment":  sb.DeleteStoryComment,
		"delete_story":          sb.DeleteStory,
		"add_persona":           sb.AddPersona,
		"update_persona":        sb.UpdatePersona,
		"delete_persona":        sb.DeletePersona,
		"facilitator_add":       sb.FacilitatorAdd,
		"facilitator_remove":    sb.FacilitatorRemove,
		"facilitator_self":      sb.FacilitatorSelf,
		"revise_color_legend":   sb.ReviseColorLegend,
		"edit_storyboard":       sb.EditStoryboard,
		"concede_storyboard":    sb.Delete,
		"abandon_storyboard":    sb.Abandon,
	},
		map[string]struct{}{
			"facilitator_add":    {},
			"facilitator_remove": {},
			"edit_storyboard":    {},
			"concede_storyboard": {},
		},
		sb.StoryboardService.ConfirmStoryboardFacilitator,
		sb.RetreatUser,
	)

	go sb.hub.Run()

	return sb
}
