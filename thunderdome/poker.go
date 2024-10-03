package thunderdome

import (
	"context"
	"time"
)

// PokerUser aka user
type PokerUser struct {
	ID           string `json:"id"`
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
	ID                   string           `json:"id"`
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
	UserID    string `json:"warriorId"`
	VoteValue string `json:"vote"`
}

// Story aka Story structure
type Story struct {
	ID                 string    `json:"id"`
	Name               string    `json:"name"`
	Type               string    `json:"type"`
	ReferenceID        string    `json:"referenceId"`
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
	CreateGame(ctx context.Context, facilitatorID string, name string, estimationScaleID string, pointValuesAllowed []string, stories []*Story, autoFinishVoting bool, pointAverageRounding string, joinCode string, facilitatorCode string, hideVoterIdentity bool) (*Poker, error)
	// TeamCreateGame creates a new poker game for a team
	TeamCreateGame(ctx context.Context, teamID string, facilitatorID string, name string, estimationScaleID string, pointValuesAllowed []string, stories []*Story, autoFinishVoting bool, pointAverageRounding string, joinCode string, facilitatorCode string, hideVoterIdentity bool) (*Poker, error)
	// UpdateGame updates an existing poker game
	UpdateGame(pokerID string, name string, pointValuesAllowed []string, autoFinishVoting bool, pointAverageRounding string, hideVoterIdentity bool, joinCode string, facilitatorCode string, teamID string) error
	// GetFacilitatorCode retrieves the facilitator code for a poker game
	GetFacilitatorCode(pokerID string) (string, error)
	// GetGame retrieves a poker game by its ID
	GetGame(pokerID string, userID string) (*Poker, error)
	// GetGamesByUser retrieves a list of poker games for a user
	GetGamesByUser(userID string, limit int, offset int) ([]*Poker, int, error)
	// ConfirmFacilitator confirms a user as a facilitator for a poker game
	ConfirmFacilitator(pokerID string, userID string) error
	// GetUserActiveStatus retrieves the active status of a user in a poker game
	GetUserActiveStatus(pokerID string, userID string) error
	// GetUsers retrieves a list of users in a poker game
	GetUsers(pokerID string) []*PokerUser
	// GetActiveUsers retrieves a list of active users in a poker game
	GetActiveUsers(pokerID string) []*PokerUser
	// AddUser adds a user to a poker game
	AddUser(pokerID string, userID string) ([]*PokerUser, error)
	// RetreatUser sets a user as inactive in a poker game
	RetreatUser(pokerID string, userID string) []*PokerUser
	// AbandonGame sets a user as abandoned in a poker game
	AbandonGame(pokerID string, userID string) ([]*PokerUser, error)
	// AddFacilitator adds a facilitator to a poker game
	AddFacilitator(pokerID string, userID string) ([]string, error)
	// RemoveFacilitator removes a facilitator from a poker game
	RemoveFacilitator(pokerID string, userID string) ([]string, error)
	// ToggleSpectator toggles a user's spectator status in a poker game
	ToggleSpectator(pokerID string, userID string, spectator bool) ([]*PokerUser, error)
	// DeleteGame deletes a poker game
	DeleteGame(pokerID string) error
	// AddFacilitatorsByEmail adds facilitators to a poker game by email
	AddFacilitatorsByEmail(ctx context.Context, pokerID string, facilitatorEmails []string) ([]string, error)
	// GetGames retrieves a list of poker games
	GetGames(limit int, offset int) ([]*Poker, int, error)
	// GetActiveGames retrieves a list of active poker games
	GetActiveGames(limit int, offset int) ([]*Poker, int, error)
	// PurgeOldGames purges poker games older than a specified number of days
	PurgeOldGames(ctx context.Context, daysOld int) error
	// GetStories retrieves a list of stories in a poker game
	GetStories(pokerID string, userID string) []*Story
	// CreateStory creates a new story in a poker game
	CreateStory(pokerID string, name string, storyType string, referenceID string, link string, description string, acceptanceCriteria string, priority int32) ([]*Story, error)
	// ActivateStoryVoting activates voting for a story in a poker game
	ActivateStoryVoting(pokerID string, storyID string) ([]*Story, error)
	// SetVote sets a user's vote for a story in a poker game
	SetVote(pokerID string, userID string, storyID string, voteValue string) (stories []*Story, allUsersVoted bool)
	// RetractVote retracts a user's vote for a story in a poker game
	RetractVote(pokerID string, userID string, storyID string) ([]*Story, error)
	// EndStoryVoting ends voting for a story in a poker game
	EndStoryVoting(pokerID string, storyID string) ([]*Story, error)
	// SkipStory skips a story in a poker game
	SkipStory(pokerID string, storyID string) ([]*Story, error)
	// UpdateStory updates an existing story in a poker game
	UpdateStory(pokerID string, storyID string, name string, storyType string, referenceID string, link string, description string, acceptanceCriteria string, priority int32) ([]*Story, error)
	// DeleteStory deletes a story from a poker game
	DeleteStory(pokerID string, storyID string) ([]*Story, error)
	// ArrangeStory sets the position of the story relative to the story it's being placed before
	ArrangeStory(pokerID string, storyID string, beforeStoryID string) ([]*Story, error)
	// FinalizeStory finalizes the points for a story in a poker game
	FinalizeStory(pokerID string, storyID string, points string) ([]*Story, error)
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
	GetEstimationScale(ctx context.Context, scaleID string) (*EstimationScale, error)
	// DeleteOrganizationEstimationScale deletes an organization's estimation scale by its ID
	DeleteOrganizationEstimationScale(ctx context.Context, orgID string, scaleID string) error
	// DeleteTeamEstimationScale deletes a team's estimation scale by its ID
	DeleteTeamEstimationScale(ctx context.Context, teamID string, scaleID string) error
	// UpdateOrganizationEstimationScale updates an existing organization estimation scale
	UpdateOrganizationEstimationScale(ctx context.Context, scale *EstimationScale) (*EstimationScale, error)
	// UpdateTeamEstimationScale updates an existing team estimation scale
	UpdateTeamEstimationScale(ctx context.Context, scale *EstimationScale) (*EstimationScale, error)
}
