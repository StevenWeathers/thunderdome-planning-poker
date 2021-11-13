package model

import "time"

// BattleUser aka user
type BattleUser struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Type      string `json:"rank"`
	Avatar    string `json:"avatar"`
	Active    bool   `json:"active"`
	Abandoned bool   `json:"abandoned"`
	Spectator bool   `json:"spectator"`
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
	CreatedDate          string        `json:"createdDate"`
	UpdatedDate          string        `json:"updatedDate"`
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
	Votes              []*Vote   `json:"votes"`
	Points             string    `json:"points"`
	Active             bool      `json:"active"`
	Skipped            bool      `json:"skipped"`
	VoteStartTime      time.Time `json:"voteStartTime"`
	VoteEndTime        time.Time `json:"voteEndTime"`
}
