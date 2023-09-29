package storyboard

import (
	"context"
	"net/http"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
)

// Service provides storyboard service
type Service struct {
	Logger                *otelzap.Logger
	ValidateSessionCookie func(w http.ResponseWriter, r *http.Request) (string, error)
	ValidateUserCookie    func(w http.ResponseWriter, r *http.Request) (string, error)
	EventHandlers         map[string]func(context.Context, string, string, string) ([]byte, error, bool)
	UserService           thunderdome.UserDataSvc
	AuthService           thunderdome.AuthDataSvc
	StoryboardService     thunderdome.StoryboardDataSvc
}

// New returns a new storyboard with websocket hub/client and event handlers
func New(
	logger *otelzap.Logger,
	validateSessionCookie func(w http.ResponseWriter, r *http.Request) (string, error),
	validateUserCookie func(w http.ResponseWriter, r *http.Request) (string, error),
	userService thunderdome.UserDataSvc, authService thunderdome.AuthDataSvc,
	storyboardService thunderdome.StoryboardDataSvc,
) *Service {
	sb := &Service{
		Logger:                logger,
		ValidateSessionCookie: validateSessionCookie,
		ValidateUserCookie:    validateUserCookie,
		UserService:           userService,
		AuthService:           authService,
		StoryboardService:     storyboardService,
	}

	sb.EventHandlers = map[string]func(context.Context, string, string, string) ([]byte, error, bool){
		"add_goal":             sb.AddGoal,
		"revise_goal":          sb.ReviseGoal,
		"delete_goal":          sb.DeleteGoal,
		"add_column":           sb.AddColumn,
		"revise_column":        sb.ReviseColumn,
		"delete_column":        sb.DeleteColumn,
		"add_story":            sb.AddStory,
		"update_story_name":    sb.UpdateStoryName,
		"update_story_content": sb.UpdateStoryContent,
		"update_story_color":   sb.UpdateStoryColor,
		"update_story_points":  sb.UpdateStoryPoints,
		"update_story_closed":  sb.UpdateStoryClosed,
		"update_story_link":    sb.UpdateStoryLink,
		"move_story":           sb.MoveStory,
		"add_story_comment":    sb.AddStoryComment,
		"edit_story_comment":   sb.EditStoryComment,
		"delete_story_comment": sb.DeleteStoryComment,
		"delete_story":         sb.DeleteStory,
		"add_persona":          sb.AddPersona,
		"update_persona":       sb.UpdatePersona,
		"delete_persona":       sb.DeletePersona,
		"facilitator_add":      sb.FacilitatorAdd,
		"facilitator_remove":   sb.FacilitatorRemove,
		"facilitator_self":     sb.FacilitatorSelf,
		"revise_color_legend":  sb.ReviseColorLegend,
		"edit_storyboard":      sb.EditStoryboard,
		"concede_storyboard":   sb.Delete,
		"abandon_storyboard":   sb.Abandon,
	}

	go h.run()

	return sb
}
