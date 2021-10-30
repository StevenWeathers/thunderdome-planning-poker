package model

import "time"

// BattleUser aka user
type BattleUser struct {
	UserID     string `json:"id"`
	UserName   string `json:"name"`
	UserType   string `json:"rank"`
	UserAvatar string `json:"avatar"`
	Active     bool   `json:"active"`
	Abandoned  bool   `json:"abandoned"`
	Spectator  bool   `json:"spectator"`
}

// Battle aka arena
type Battle struct {
	BattleID             string        `json:"id"`
	BattleName           string        `json:"name"`
	Users                []*BattleUser `json:"warriors"`
	Plans                []*Plan       `json:"plans"`
	VotingLocked         bool          `json:"votingLocked"`
	ActivePlanID         string        `json:"activePlanId"`
	PointValuesAllowed   []string      `json:"pointValuesAllowed"`
	AutoFinishVoting     bool          `json:"autoFinishVoting"`
	Leaders              []string      `json:"leaders"`
	PointAverageRounding string        `json:"pointAverageRounding"`
}

// Vote structure
type Vote struct {
	UserID    string `json:"warriorId"`
	VoteValue string `json:"vote"`
}

// Plan aka Story structure
type Plan struct {
	PlanID             string    `json:"id"`
	PlanName           string    `json:"name"`
	Type               string    `json:"type"`
	ReferenceID        string    `json:"referenceId"`
	Link               string    `json:"link"`
	Description        string    `json:"description"`
	AcceptanceCriteria string    `json:"acceptanceCriteria"`
	Votes              []*Vote   `json:"votes"`
	Points             string    `json:"points"`
	PlanActive         bool      `json:"active"`
	PlanSkipped        bool      `json:"skipped"`
	VoteStartTime      time.Time `json:"voteStartTime"`
	VoteEndTime        time.Time `json:"voteEndTime"`
}
