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
}

// Poker aka arena
type Poker struct {
	Id                   string       `json:"id"`
	Name                 string       `json:"name"`
	Users                []*PokerUser `json:"users"`
	Stories              []*Story     `json:"plans"`
	VotingLocked         bool         `json:"votingLocked"`
	ActiveStoryID        string       `json:"activePlanId"`
	PointValuesAllowed   []string     `json:"pointValuesAllowed"`
	AutoFinishVoting     bool         `json:"autoFinishVoting"`
	Facilitators         []string     `json:"leaders"`
	PointAverageRounding string       `json:"pointAverageRounding"`
	HideVoterIdentity    bool         `json:"hideVoterIdentity"`
	JoinCode             string       `json:"joinCode"`
	FacilitatorCode      string       `json:"leaderCode,omitempty"`
	TeamID               string       `json:"teamId"`
	CreatedDate          time.Time    `json:"createdDate"`
	UpdatedDate          time.Time    `json:"updatedDate"`
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

type PokerDataSvc interface {
	CreateGame(ctx context.Context, FacilitatorID string, Name string, PointValuesAllowed []string, Stories []*Story, AutoFinishVoting bool, PointAverageRounding string, JoinCode string, FacilitatorCode string, HideVoterIdentity bool) (*Poker, error)
	TeamCreateGame(ctx context.Context, TeamID string, FacilitatorID string, Name string, PointValuesAllowed []string, Stories []*Story, AutoFinishVoting bool, PointAverageRounding string, JoinCode string, FacilitatorCode string, HideVoterIdentity bool) (*Poker, error)
	UpdateGame(PokerID string, Name string, PointValuesAllowed []string, AutoFinishVoting bool, PointAverageRounding string, HideVoterIdentity bool, JoinCode string, FacilitatorCode string, TeamID string) error
	GetFacilitatorCode(PokerID string) (string, error)
	GetGame(PokerID string, UserID string) (*Poker, error)
	GetGamesByUser(UserID string, Limit int, Offset int) ([]*Poker, int, error)
	ConfirmFacilitator(PokerID string, UserID string) error
	GetUserActiveStatus(PokerID string, UserID string) error
	GetUsers(PokerID string) []*PokerUser
	GetActiveUsers(PokerID string) []*PokerUser
	AddUser(PokerID string, UserID string) ([]*PokerUser, error)
	RetreatUser(PokerID string, UserID string) []*PokerUser
	AbandonGame(PokerID string, UserID string) ([]*PokerUser, error)
	AddFacilitator(PokerID string, UserID string) ([]string, error)
	RemoveFacilitator(PokerID string, UserID string) ([]string, error)
	ToggleSpectator(PokerID string, UserID string, Spectator bool) ([]*PokerUser, error)
	DeleteGame(PokerID string) error
	AddFacilitatorsByEmail(ctx context.Context, PokerID string, FacilitatorEmails []string) ([]string, error)
	GetGames(Limit int, Offset int) ([]*Poker, int, error)
	GetActiveGames(Limit int, Offset int) ([]*Poker, int, error)
	PurgeOldGames(ctx context.Context, DaysOld int) error
	GetStories(PokerID string, UserID string) []*Story
	CreateStory(PokerID string, Name string, Type string, ReferenceID string, Link string, Description string, AcceptanceCriteria string, Priority int32) ([]*Story, error)
	ActivateStoryVoting(PokerID string, StoryID string) ([]*Story, error)
	SetVote(PokerID string, UserID string, StoryID string, VoteValue string) (BattlePlans []*Story, AllUsersVoted bool)
	RetractVote(PokerID string, UserID string, StoryID string) ([]*Story, error)
	EndStoryVoting(PokerID string, StoryID string) ([]*Story, error)
	SkipStory(PokerID string, StoryID string) ([]*Story, error)
	UpdateStory(PokerID string, StoryID string, Name string, Type string, ReferenceID string, Link string, Description string, AcceptanceCriteria string, Priority int32) ([]*Story, error)
	DeleteStory(PokerID string, StoryID string) ([]*Story, error)
	// ArrangeStory sets the position of the story relative to the story it's being placed before
	ArrangeStory(PokerID string, StoryID string, BeforeStoryID string) ([]*Story, error)
	FinalizeStory(PokerID string, StoryID string, Points string) ([]*Story, error)
}
