// Package poker provides Poker Game http handlers for Thunderdome
package poker

import (
	"context"
	"net/http"
	"time"

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

type PokerDataSvc interface {
	// UpdateGame updates an existing poker game
	UpdateGame(pokerID string, name string, pointValuesAllowed []string, autoFinishVoting bool, pointAverageRounding string, hideVoterIdentity bool, joinCode string, facilitatorCode string, teamID string) error
	// GetFacilitatorCode retrieves the facilitator code for a poker game
	GetFacilitatorCode(pokerID string) (string, error)
	// GetGameByID retrieves a poker game by its ID
	GetGameByID(pokerID string, userID string) (*thunderdome.Poker, error)
	// ConfirmFacilitator confirms a user as a facilitator for a poker game
	ConfirmFacilitator(pokerID string, userID string) error
	// GetUserActiveStatus retrieves the active status of a user in a poker game
	GetUserActiveStatus(pokerID string, userID string) error
	// AddUser adds a user to a poker game
	AddUser(pokerID string, userID string) ([]*thunderdome.PokerUser, error)
	// RetreatUser sets a user as inactive in a poker game
	RetreatUser(pokerID string, userID string) []*thunderdome.PokerUser
	// AbandonGame sets a user as abandoned in a poker game
	AbandonGame(pokerID string, userID string) ([]*thunderdome.PokerUser, error)
	// AddFacilitator adds a facilitator to a poker game
	AddFacilitator(pokerID string, userID string) ([]string, error)
	// RemoveFacilitator removes a facilitator from a poker game
	RemoveFacilitator(pokerID string, userID string) ([]string, error)
	// ToggleSpectator toggles a user's spectator status in a poker game
	ToggleSpectator(pokerID string, userID string, spectator bool) ([]*thunderdome.PokerUser, error)
	// DeleteGame deletes a poker game
	DeleteGame(pokerID string) error
	// CreateStory creates a new story in a poker game
	CreateStory(pokerID string, name string, storyType string, referenceID string, link string, description string, acceptanceCriteria string, priority int32) ([]*thunderdome.Story, error)
	// ActivateStoryVoting activates voting for a story in a poker game
	ActivateStoryVoting(pokerID string, storyID string) ([]*thunderdome.Story, error)
	// SetVote sets a user's vote for a story in a poker game
	SetVote(pokerID string, userID string, storyID string, voteValue string) (stories []*thunderdome.Story, allUsersVoted bool)
	// RetractVote retracts a user's vote for a story in a poker game
	RetractVote(pokerID string, userID string, storyID string) ([]*thunderdome.Story, error)
	// EndStoryVoting ends voting for a story in a poker game
	EndStoryVoting(pokerID string, storyID string) ([]*thunderdome.Story, error)
	// SkipStory skips a story in a poker game
	SkipStory(pokerID string, storyID string) ([]*thunderdome.Story, error)
	// UpdateStory updates an existing story in a poker game
	UpdateStory(pokerID string, storyID string, name string, storyType string, referenceID string, link string, description string, acceptanceCriteria string, priority int32) ([]*thunderdome.Story, error)
	// DeleteStory deletes a story from a poker game
	DeleteStory(pokerID string, storyID string) ([]*thunderdome.Story, error)
	// ArrangeStory sets the position of the story relative to the story it's being placed before
	ArrangeStory(pokerID string, storyID string, beforeStoryID string) ([]*thunderdome.Story, error)
	// FinalizeStory finalizes the points for a story in a poker game
	FinalizeStory(pokerID string, storyID string, points string) ([]*thunderdome.Story, error)
	// EndGame ends a poker game with a specified reason
	EndGame(ctx context.Context, pokerID string, endReason string) (string, time.Time, error)
}

type AuthDataSvc interface {
	GetSessionUserByID(ctx context.Context, sessionID string) (*thunderdome.User, error)
}

type UserDataSvc interface {
	GetGuestUserByID(ctx context.Context, userID string) (*thunderdome.User, error)
}

// Service provides battle service
type Service struct {
	config                Config
	logger                *otelzap.Logger
	validateSessionCookie func(w http.ResponseWriter, r *http.Request) (string, error)
	validateUserCookie    func(w http.ResponseWriter, r *http.Request) (string, error)
	UserService           UserDataSvc
	AuthService           AuthDataSvc
	PokerService          PokerDataSvc
	hub                   *wshub.Hub
}

// New returns a new battle with websocket hub/client and event handlers
func New(
	config Config, logger *otelzap.Logger,
	validateSessionCookie func(w http.ResponseWriter, r *http.Request) (string, error),
	validateUserCookie func(w http.ResponseWriter, r *http.Request) (string, error),
	userService UserDataSvc, authService AuthDataSvc,
	pokerDataService PokerDataSvc,
) *Service {
	b := &Service{
		config:                config,
		logger:                logger,
		validateSessionCookie: validateSessionCookie,
		validateUserCookie:    validateUserCookie,
		UserService:           userService,
		AuthService:           authService,
		PokerService:          pokerDataService,
	}

	b.hub = wshub.NewHub(logger, wshub.Config{
		AppDomain:          config.AppDomain,
		WebsocketSubdomain: config.WebsocketSubdomain,
		WriteWaitSec:       config.WriteWaitSec,
		PongWaitSec:        config.PongWaitSec,
		PingPeriodSec:      config.PingPeriodSec,
	}, map[string]func(context.Context, string, string, string) (any, []byte, error, bool){
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
		"end_game":         b.EndGame,
		"concede_battle":   b.Delete,
		"abandon_battle":   b.Abandon,
	},
		map[string]struct{}{
			"add_plan":       {},
			"revise_plan":    {},
			"burn_plan":      {},
			"story_arrange":  {},
			"activate_plan":  {},
			"skip_plan":      {},
			"end_voting":     {},
			"finalize_plan":  {},
			"jab_warrior":    {},
			"promote_leader": {},
			"demote_leader":  {},
			"revise_battle":  {},
			"end_game":       {},
			"concede_battle": {},
		},
		b.PokerService.ConfirmFacilitator,
		b.RetreatUser,
	)

	go b.hub.Run()

	return b
}
