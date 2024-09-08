package thunderdome

import (
	"context"
	"time"
)

// PokerUser aka user
type PokerUser struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	Type         string `json:"rank"`
	Avatar       string `json:"avatar"`
	Active       bool   `json:"active"`
	Abandoned    bool   `json:"abandoned"`
	Spectator    bool   `json:"spectator"`
	GravatarHash string `json:"gravatarHash"`
	PictureURL   string `json:"pictureUrl"`
}

// Poker aka arena
type Poker struct {
	Id                   string           `json:"id"`
	Name                 string           `json:"name"`
	Users                []*PokerUser     `json:"users"`
	Stories              []*Story         `json:"plans"`
	VotingLocked         bool             `json:"votingLocked"`
	ActiveStoryID        string           `json:"activePlanId"`
	PointValuesAllowed   []string         `json:"pointValuesAllowed"`
	AutoFinishVoting     bool             `json:"autoFinishVoting"`
	Facilitators         []string         `json:"leaders"`
	PointAverageRounding string           `json:"pointAverageRounding"`
	HideVoterIdentity    bool             `json:"hideVoterIdentity"`
	JoinCode             string           `json:"joinCode"`
	FacilitatorCode      string           `json:"leaderCode,omitempty"`
	TeamID               string           `json:"teamId"`
	TeamName             string           `json:"teamName"`
	EstimationScaleID    string           `json:"estimationScaleId"`
	EstimationScale      *EstimationScale `json:"estimationScale,omitempty"`
	CreatedDate          time.Time        `json:"createdDate"`
	UpdatedDate          time.Time        `json:"updatedDate"`
}

// Vote structure
type Vote struct {
	UserId    string `json:"warriorId"`
	VoteValue string `json:"vote"`
}

// Story aka Story structure
type Story struct {
	Id                 string    `json:"id"`
	Name               string    `json:"name"`
	Type               string    `json:"type"`
	ReferenceId        string    `json:"referenceId"`
	Link               string    `json:"link"`
	Description        string    `json:"description"`
	AcceptanceCriteria string    `json:"acceptanceCriteria"`
	Priority           int32     `json:"priority"`
	Votes              []*Vote   `json:"votes"`
	Points             string    `json:"points"`
	Active             bool      `json:"active"`
	Skipped            bool      `json:"skipped"`
	VoteStartTime      time.Time `json:"voteStartTime"`
	VoteEndTime        time.Time `json:"voteEndTime"`
	Position           int32     `json:"position"`
}

type EstimationScale struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	ScaleType      string    `json:"scaleType"`
	Values         []string  `json:"values"`
	CreatedBy      string    `json:"createdBy"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
	IsPublic       bool      `json:"isPublic"`
	OrganizationID string    `json:"organizationId"`
	TeamID         string    `json:"teamId"`
	DefaultScale   bool      `json:"defaultScale"`
}

type PokerDataSvc interface {
	// CreateGame creates a new poker game
	CreateGame(ctx context.Context, FacilitatorID string, Name string, EstimationScaleID string, PointValuesAllowed []string, Stories []*Story, AutoFinishVoting bool, PointAverageRounding string, JoinCode string, FacilitatorCode string, HideVoterIdentity bool) (*Poker, error)
	// TeamCreateGame creates a new poker game for a team
	TeamCreateGame(ctx context.Context, TeamID string, FacilitatorID string, Name string, EstimationScaleID string, PointValuesAllowed []string, Stories []*Story, AutoFinishVoting bool, PointAverageRounding string, JoinCode string, FacilitatorCode string, HideVoterIdentity bool) (*Poker, error)
	// UpdateGame updates an existing poker game
	UpdateGame(PokerID string, Name string, PointValuesAllowed []string, AutoFinishVoting bool, PointAverageRounding string, HideVoterIdentity bool, JoinCode string, FacilitatorCode string, TeamID string) error
	// GetFacilitatorCode retrieves the facilitator code for a poker game
	GetFacilitatorCode(PokerID string) (string, error)
	// GetGame retrieves a poker game by its ID
	GetGame(PokerID string, UserID string) (*Poker, error)
	// GetGamesByUser retrieves a list of poker games for a user
	GetGamesByUser(UserID string, Limit int, Offset int) ([]*Poker, int, error)
	// ConfirmFacilitator confirms a user as a facilitator for a poker game
	ConfirmFacilitator(PokerID string, UserID string) error
	// GetUserActiveStatus retrieves the active status of a user in a poker game
	GetUserActiveStatus(PokerID string, UserID string) error
	// GetUsers retrieves a list of users in a poker game
	GetUsers(PokerID string) []*PokerUser
	// GetActiveUsers retrieves a list of active users in a poker game
	GetActiveUsers(PokerID string) []*PokerUser
	// AddUser adds a user to a poker game
	AddUser(PokerID string, UserID string) ([]*PokerUser, error)
	// RetreatUser sets a user as inactive in a poker game
	RetreatUser(PokerID string, UserID string) []*PokerUser
	// AbandonGame sets a user as abandoned in a poker game
	AbandonGame(PokerID string, UserID string) ([]*PokerUser, error)
	// AddFacilitator adds a facilitator to a poker game
	AddFacilitator(PokerID string, UserID string) ([]string, error)
	// RemoveFacilitator removes a facilitator from a poker game
	RemoveFacilitator(PokerID string, UserID string) ([]string, error)
	// ToggleSpectator toggles a user's spectator status in a poker game
	ToggleSpectator(PokerID string, UserID string, Spectator bool) ([]*PokerUser, error)
	// DeleteGame deletes a poker game
	DeleteGame(PokerID string) error
	// AddFacilitatorsByEmail adds facilitators to a poker game by email
	AddFacilitatorsByEmail(ctx context.Context, PokerID string, FacilitatorEmails []string) ([]string, error)
	// GetGames retrieves a list of poker games
	GetGames(Limit int, Offset int) ([]*Poker, int, error)
	// GetActiveGames retrieves a list of active poker games
	GetActiveGames(Limit int, Offset int) ([]*Poker, int, error)
	// PurgeOldGames purges poker games older than a specified number of days
	PurgeOldGames(ctx context.Context, DaysOld int) error
	// GetStories retrieves a list of stories in a poker game
	GetStories(PokerID string, UserID string) []*Story
	// CreateStory creates a new story in a poker game
	CreateStory(PokerID string, Name string, Type string, ReferenceID string, Link string, Description string, AcceptanceCriteria string, Priority int32) ([]*Story, error)
	// ActivateStoryVoting activates voting for a story in a poker game
	ActivateStoryVoting(PokerID string, StoryID string) ([]*Story, error)
	// SetVote sets a user's vote for a story in a poker game
	SetVote(PokerID string, UserID string, StoryID string, VoteValue string) (BattlePlans []*Story, AllUsersVoted bool)
	// RetractVote retracts a user's vote for a story in a poker game
	RetractVote(PokerID string, UserID string, StoryID string) ([]*Story, error)
	// EndStoryVoting ends voting for a story in a poker game
	EndStoryVoting(PokerID string, StoryID string) ([]*Story, error)
	// SkipStory skips a story in a poker game
	SkipStory(PokerID string, StoryID string) ([]*Story, error)
	// UpdateStory updates an existing story in a poker game
	UpdateStory(PokerID string, StoryID string, Name string, Type string, ReferenceID string, Link string, Description string, AcceptanceCriteria string, Priority int32) ([]*Story, error)
	// DeleteStory deletes a story from a poker game
	DeleteStory(PokerID string, StoryID string) ([]*Story, error)
	// ArrangeStory sets the position of the story relative to the story it's being placed before
	ArrangeStory(PokerID string, StoryID string, BeforeStoryID string) ([]*Story, error)
	// FinalizeStory finalizes the points for a story in a poker game
	FinalizeStory(PokerID string, StoryID string, Points string) ([]*Story, error)
	// GetEstimationScales retrieves a list of estimation scales
	GetEstimationScales(ctx context.Context, limit, offset int) ([]*EstimationScale, int, error)
	// GetPublicEstimationScales retrieves a list of public estimation scales
	GetPublicEstimationScales(ctx context.Context, limit, offset int) ([]*EstimationScale, int, error)
	// CreateEstimationScale creates a new estimation scale
	CreateEstimationScale(ctx context.Context, scale *EstimationScale) (*EstimationScale, error)
	// UpdateEstimationScale updates an existing estimation scale
	UpdateEstimationScale(ctx context.Context, scale *EstimationScale) (*EstimationScale, error)
	// DeleteEstimationScale deletes an estimation scale by its ID
	DeleteEstimationScale(ctx context.Context, scaleID string) error
	// GetDefaultEstimationScale retrieves the default estimation scale for an organization or team
	GetDefaultEstimationScale(ctx context.Context, organizationID, teamID string) (*EstimationScale, error)
	// GetDefaultPublicEstimationScale retrieves the default public estimation scale
	GetDefaultPublicEstimationScale(ctx context.Context) (*EstimationScale, error)
	// GetPublicEstimationScale retrieves a public estimation scale by its ID
	GetPublicEstimationScale(ctx context.Context, id string) (*EstimationScale, error)
	// GetOrganizationEstimationScales retrieves a list of estimation scales for an organization
	GetOrganizationEstimationScales(ctx context.Context, orgID string, limit, offset int) ([]*EstimationScale, int, error)
	// GetTeamEstimationScales retrieves a list of estimation scales for a team
	GetTeamEstimationScales(ctx context.Context, teamID string, limit, offset int) ([]*EstimationScale, int, error)
	// GetEstimationScale retrieves an estimation scale by its ID
	GetEstimationScale(ctx context.Context, id string) (*EstimationScale, error)
}
