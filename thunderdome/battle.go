package thunderdome

import (
	"context"
	"time"
)

// BattleUser aka user
type BattleUser struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	Type         string `json:"rank"`
	Avatar       string `json:"avatar"`
	Active       bool   `json:"active"`
	Abandoned    bool   `json:"abandoned"`
	Spectator    bool   `json:"spectator"`
	GravatarHash string `json:"gravatarHash"`
}

// Battle aka arena
type Battle struct {
	Id                   string        `json:"id"`
	Name                 string        `json:"name"`
	Users                []*BattleUser `json:"users"`
	Plans                []*Plan       `json:"plans"`
	VotingLocked         bool          `json:"votingLocked"`
	ActivePlanID         string        `json:"activePlanId"`
	PointValuesAllowed   []string      `json:"pointValuesAllowed"`
	AutoFinishVoting     bool          `json:"autoFinishVoting"`
	Leaders              []string      `json:"leaders"`
	PointAverageRounding string        `json:"pointAverageRounding"`
	HideVoterIdentity    bool          `json:"hideVoterIdentity"`
	JoinCode             string        `json:"joinCode"`
	LeaderCode           string        `json:"leaderCode,omitempty"`
	CreatedDate          time.Time     `json:"createdDate"`
	UpdatedDate          time.Time     `json:"updatedDate"`
}

// Vote structure
type Vote struct {
	UserId    string `json:"warriorId"`
	VoteValue string `json:"vote"`
}

// Plan aka Story structure
type Plan struct {
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

type BattleService interface {
	CreateBattle(ctx context.Context, LeaderID string, BattleName string, PointValuesAllowed []string, Plans []*Plan, AutoFinishVoting bool, PointAverageRounding string, JoinCode string, LeaderCode string, HideVoterIdentity bool) (*Battle, error)
	TeamCreateBattle(ctx context.Context, TeamID string, LeaderID string, BattleName string, PointValuesAllowed []string, Plans []*Plan, AutoFinishVoting bool, PointAverageRounding string, JoinCode string, LeaderCode string, HideVoterIdentity bool) (*Battle, error)
	ReviseBattle(BattleID string, BattleName string, PointValuesAllowed []string, AutoFinishVoting bool, PointAverageRounding string, HideVoterIdentity bool, JoinCode string, LeaderCode string) error
	GetBattleLeaderCode(BattleID string) (string, error)
	GetBattle(BattleID string, UserID string) (*Battle, error)
	GetBattlesByUser(UserID string, Limit int, Offset int) ([]*Battle, int, error)
	ConfirmLeader(BattleID string, UserID string) error
	GetBattleUserActiveStatus(BattleID string, UserID string) error
	GetBattleUsers(BattleID string) []*BattleUser
	GetBattleActiveUsers(BattleID string) []*BattleUser
	AddUserToBattle(BattleID string, UserID string) ([]*BattleUser, error)
	RetreatUser(BattleID string, UserID string) []*BattleUser
	AbandonBattle(BattleID string, UserID string) ([]*BattleUser, error)
	SetBattleLeader(BattleID string, LeaderID string) ([]string, error)
	DemoteBattleLeader(BattleID string, LeaderID string) ([]string, error)
	ToggleSpectator(BattleID string, UserID string, Spectator bool) ([]*BattleUser, error)
	DeleteBattle(BattleID string) error
	AddBattleLeadersByEmail(ctx context.Context, BattleID string, LeaderEmails []string) ([]string, error)
	GetBattles(Limit int, Offset int) ([]*Battle, int, error)
	GetActiveBattles(Limit int, Offset int) ([]*Battle, int, error)
	CleanBattles(ctx context.Context, DaysOld int) error
	GetPlans(BattleID string, UserID string) []*Plan
	CreatePlan(BattleID string, PlanName string, PlanType string, ReferenceID string, Link string, Description string, AcceptanceCriteria string, Priority int32) ([]*Plan, error)
	ActivatePlanVoting(BattleID string, PlanID string) ([]*Plan, error)
	SetVote(BattleID string, UserID string, PlanID string, VoteValue string) (BattlePlans []*Plan, AllUsersVoted bool)
	RetractVote(BattleID string, UserID string, PlanID string) ([]*Plan, error)
	EndPlanVoting(BattleID string, PlanID string) ([]*Plan, error)
	SkipPlan(BattleID string, PlanID string) ([]*Plan, error)
	RevisePlan(BattleID string, PlanID string, PlanName string, PlanType string, ReferenceID string, Link string, Description string, AcceptanceCriteria string, Priority int32) ([]*Plan, error)
	BurnPlan(BattleID string, PlanID string) ([]*Plan, error)
	FinalizePlan(BattleID string, PlanID string, PlanPoints string) ([]*Plan, error)
}
