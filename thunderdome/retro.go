package thunderdome

import (
	"time"
)

// Color is a color legend
type Color struct {
	Color  string `json:"color"`
	Legend string `json:"legend"`
}

// RetroUser aka user
type RetroUser struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	Active       bool   `json:"active"`
	Avatar       string `json:"avatar"`
	GravatarHash string `json:"gravatarHash"`
	PictureURL   string `json:"pictureUrl"`
}

// Retro A story mapping board
type Retro struct {
	ID                    string         `json:"id" db:"id"`
	OwnerID               string         `json:"ownerId" db:"owner_id"`
	Name                  string         `json:"name" db:"name"`
	TemplateID            string         `json:"template_id" db:"template_id"`
	Users                 []*RetroUser   `json:"users"`
	Groups                []*RetroGroup  `json:"groups"`
	Items                 []*RetroItem   `json:"items"`
	ActionItems           []*RetroAction `json:"actionItems"`
	Votes                 []*RetroVote   `json:"votes"`
	ReadyUsers            []string       `json:"readyUsers"`
	Facilitators          []string       `json:"facilitators"`
	Phase                 string         `json:"phase" db:"phase"`
	PhaseTimeLimitMin     int            `json:"phase_time_limit_min" db:"phase_time_limit_min"`
	PhaseTimeStart        time.Time      `json:"phase_time_start" db:"phase_time_start"`
	PhaseAutoAdvance      bool           `json:"phase_auto_advance" db:"phase_auto_advance"`
	JoinCode              string         `json:"joinCode" db:"join_code"`
	FacilitatorCode       string         `json:"facilitatorCode" db:"facilitator_code"`
	MaxVotes              int            `json:"maxVotes" db:"max_votes"`
	BrainstormVisibility  string         `json:"brainstormVisibility" db:"brainstorm_visibility"`
	AllowCumulativeVoting bool           `json:"allowCumulativeVoting" db:"allow_cumulative_voting"`
	HideVotesDuringVoting bool           `json:"hideVotesDuringVoting" db:"hide_votes_during_voting"`
	Template              RetroTemplate  `json:"template"`
	TeamID                string         `json:"teamId" db:"team_id"`
	TeamName              string         `json:"teamName"`
	CreatedDate           string         `json:"createdDate" db:"created_date"`
	UpdatedDate           string         `json:"updatedDate" db:"updated_date"`
}

// RetroItem can be a pro (went well/worked), con (needs improvement), or a question
type RetroItem struct {
	ID       string              `json:"id" db:"id"`
	UserID   string              `json:"userId" db:"user_id"`
	GroupID  string              `json:"groupId" db:"group_id"`
	Content  string              `json:"content" db:"content"`
	Type     string              `json:"type" db:"type"`
	Comments []*RetroItemComment `json:"comments"`
}

// RetroGroup is a grouping of retro items
type RetroGroup struct {
	ID   string `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

// RetroAction is an action the team can take based on retro feedback
type RetroAction struct {
	RetroID   string                `json:"retroId,omitempty"`
	ID        string                `json:"id" db:"id"`
	Content   string                `json:"content" db:"content"`
	Completed bool                  `json:"completed" db:"completed"`
	Comments  []*RetroActionComment `json:"comments"`
	Assignees []*User               `json:"assignees"`
}

// RetroActionComment A retro action comment by a user
type RetroActionComment struct {
	ID          string `json:"id"`
	ActionID    string `json:"action_id"`
	UserID      string `json:"user_id"`
	Comment     string `json:"comment"`
	CreateDate  string `json:"created_date"`
	UpdatedDate string `json:"updated_date"`
}

// RetroItemComment A retro item comment by a user
type RetroItemComment struct {
	ID          string `json:"id"`
	ItemID      string `json:"item_id"`
	UserID      string `json:"user_id"`
	Comment     string `json:"comment"`
	CreateDate  string `json:"created_date"`
	UpdatedDate string `json:"updated_date"`
}

// RetroVote is a users vote toward a retro item group
type RetroVote struct {
	UserID  string `json:"userId" db:"user_id"`
	GroupID string `json:"groupId" db:"group_id"`
	Count   int    `json:"count" db:"vote_count"`
}

// RetroSettings represents the settings for a retro session
type RetroSettings struct {
	ID                    string    `json:"id"`
	OrganizationID        *string   `json:"organizationId"`
	DepartmentID          *string   `json:"departmentId"`
	TeamID                *string   `json:"teamId"`
	MaxVotes              int16     `json:"maxVotes"`
	AllowMultipleVotes    bool      `json:"allowMultipleVotes"`
	BrainstormVisibility  string    `json:"brainstormVisibility"`
	PhaseTimeLimit        int16     `json:"phaseTimeLimit"`
	PhaseAutoAdvance      bool      `json:"phaseAutoAdvance"`
	AllowCumulativeVoting bool      `json:"allowCumulativeVoting"`
	TemplateID            *string   `json:"templateId"`
	JoinCode              string    `json:"joinCode"`
	FacilitatorCode       string    `json:"facilitatorCode"`
	CreatedAt             time.Time `json:"createdAt"`
	UpdatedAt             time.Time `json:"updatedAt"`
}
