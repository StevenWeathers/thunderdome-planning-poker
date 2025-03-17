package thunderdome

import (
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

// PokerSettings represents the default settings for a poker session creation
type PokerSettings struct {
	ID                   string    `json:"ID"`
	OrganizationID       *string   `json:"organizationId"`
	DepartmentID         *string   `json:"departmentId"`
	TeamID               *string   `json:"teamId"`
	AutoFinishVoting     bool      `json:"autoFinishVoting"`
	PointAverageRounding string    `json:"pointAverageRounding"`
	HideVoterIdentity    bool      `json:"hideVoterIdentity"`
	EstimationScaleID    *string   `json:"estimationScaleId"`
	JoinCode             string    `json:"joinCode"`
	FacilitatorCode      string    `json:"facilitatorCode"`
	CreatedAt            time.Time `json:"createdAt"`
	UpdatedAt            time.Time `json:"updatedAt"`
}
