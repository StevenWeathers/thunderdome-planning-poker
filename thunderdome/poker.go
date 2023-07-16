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
}

type PokerDataSvc interface {
	CreateBattle(ctx context.Context, LeaderID string, BattleName string, PointValuesAllowed []string, Plans []*Story, AutoFinishVoting bool, PointAverageRounding string, JoinCode string, LeaderCode string, HideVoterIdentity bool) (*Poker, error)
	TeamCreateBattle(ctx context.Context, TeamID string, LeaderID string, BattleName string, PointValuesAllowed []string, Plans []*Story, AutoFinishVoting bool, PointAverageRounding string, JoinCode string, LeaderCode string, HideVoterIdentity bool) (*Poker, error)
	ReviseBattle(BattleID string, BattleName string, PointValuesAllowed []string, AutoFinishVoting bool, PointAverageRounding string, HideVoterIdentity bool, JoinCode string, LeaderCode string) error
	GetBattleLeaderCode(BattleID string) (string, error)
	GetBattle(BattleID string, UserID string) (*Poker, error)
	GetBattlesByUser(UserID string, Limit int, Offset int) ([]*Poker, int, error)
	ConfirmLeader(BattleID string, UserID string) error
	GetBattleUserActiveStatus(BattleID string, UserID string) error
	GetBattleUsers(BattleID string) []*PokerUser
	GetBattleActiveUsers(BattleID string) []*PokerUser
	AddUserToBattle(BattleID string, UserID string) ([]*PokerUser, error)
	RetreatUser(BattleID string, UserID string) []*PokerUser
	AbandonBattle(BattleID string, UserID string) ([]*PokerUser, error)
	SetBattleLeader(BattleID string, LeaderID string) ([]string, error)
	DemoteBattleLeader(BattleID string, LeaderID string) ([]string, error)
	ToggleSpectator(BattleID string, UserID string, Spectator bool) ([]*PokerUser, error)
	DeleteBattle(BattleID string) error
	AddBattleLeadersByEmail(ctx context.Context, BattleID string, LeaderEmails []string) ([]string, error)
	GetBattles(Limit int, Offset int) ([]*Poker, int, error)
	GetActiveBattles(Limit int, Offset int) ([]*Poker, int, error)
	CleanBattles(ctx context.Context, DaysOld int) error
	GetPlans(BattleID string, UserID string) []*Story
	CreatePlan(BattleID string, PlanName string, PlanType string, ReferenceID string, Link string, Description string, AcceptanceCriteria string, Priority int32) ([]*Story, error)
	ActivatePlanVoting(BattleID string, PlanID string) ([]*Story, error)
	SetVote(BattleID string, UserID string, PlanID string, VoteValue string) (BattlePlans []*Story, AllUsersVoted bool)
	RetractVote(BattleID string, UserID string, PlanID string) ([]*Story, error)
	EndPlanVoting(BattleID string, PlanID string) ([]*Story, error)
	SkipPlan(BattleID string, PlanID string) ([]*Story, error)
	RevisePlan(BattleID string, PlanID string, PlanName string, PlanType string, ReferenceID string, Link string, Description string, AcceptanceCriteria string, Priority int32) ([]*Story, error)
	BurnPlan(BattleID string, PlanID string) ([]*Story, error)
	FinalizePlan(BattleID string, PlanID string, PlanPoints string) ([]*Story, error)
}
