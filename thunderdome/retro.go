package thunderdome

import "context"

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
}

// Retro A story mapping board
type Retro struct {
	Id                   string         `json:"id" db:"id"`
	OwnerID              string         `json:"ownerId" db:"owner_id"`
	Name                 string         `json:"name" db:"name"`
	Users                []*RetroUser   `json:"users"`
	Groups               []*RetroGroup  `json:"groups"`
	Items                []*RetroItem   `json:"items"`
	ActionItems          []*RetroAction `json:"actionItems"`
	Votes                []*RetroVote   `json:"votes"`
	Facilitators         []string       `json:"facilitators"`
	Format               string         `json:"format" db:"format"`
	Phase                string         `json:"phase" db:"phase"`
	JoinCode             string         `json:"joinCode" db:"join_code"`
	FacilitatorCode      string         `json:"facilitatorCode" db:"facilitator_code"`
	MaxVotes             int            `json:"maxVotes" db:"max_votes"`
	BrainstormVisibility string         `json:"brainstormVisibility" db:"brainstorm_visibility"`
	TeamName             string         `json:"teamName"`
	CreatedDate          string         `json:"createdDate" db:"created_date"`
	UpdatedDate          string         `json:"updatedDate" db:"updated_date"`
}

// RetroItem can be a pro (went well/worked), con (needs improvement), or a question
type RetroItem struct {
	ID      string `json:"id" db:"id"`
	UserID  string `json:"userId" db:"user_id"`
	GroupID string `json:"groupId" db:"group_id"`
	Content string `json:"content" db:"content"`
	Type    string `json:"type" db:"type"`
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
	RetroID     string `json:"retro_id"`
	UserID      string `json:"user_id"`
	Comment     string `json:"comment"`
	CreateDate  string `json:"created_date"`
	UpdatedDate string `json:"updated_date"`
}

// RetroVote is a users vote toward a retro item group
type RetroVote struct {
	UserID  string `json:"userId" db:"user_id"`
	GroupID string `json:"groupId" db:"group_id"`
}

type RetroDataSvc interface {
	RetroCreate(OwnerID string, RetroName string, Format string, JoinCode string, FacilitatorCode string, MaxVotes int, BrainstormVisibility string) (*Retro, error)
	TeamRetroCreate(ctx context.Context, TeamID string, OwnerID string, RetroName string, Format string, JoinCode string, FacilitatorCode string, MaxVotes int, BrainstormVisibility string) (*Retro, error)
	EditRetro(RetroID string, RetroName string, JoinCode string, FacilitatorCode string, maxVotes int, brainstormVisibility string) error
	RetroGet(RetroID string, UserID string) (*Retro, error)
	RetroGetByUser(UserID string, Limit int, Offset int) ([]*Retro, int, error)
	RetroConfirmFacilitator(RetroID string, userID string) error
	RetroGetUsers(RetroID string) []*RetroUser
	GetRetroFacilitators(RetroID string) []string
	RetroAddUser(RetroID string, UserID string) ([]*RetroUser, error)
	RetroFacilitatorAdd(RetroID string, UserID string) ([]string, error)
	RetroFacilitatorRemove(RetroID string, UserID string) ([]string, error)
	RetroRetreatUser(RetroID string, UserID string) []*RetroUser
	RetroAbandon(RetroID string, UserID string) ([]*RetroUser, error)
	RetroAdvancePhase(RetroID string, Phase string) (*Retro, error)
	RetroDelete(RetroID string) error
	GetRetroUserActiveStatus(RetroID string, UserID string) error
	GetRetros(Limit int, Offset int) ([]*Retro, int, error)
	GetActiveRetros(Limit int, Offset int) ([]*Retro, int, error)
	GetRetroFacilitatorCode(RetroID string) (string, error)
	CleanRetros(ctx context.Context, DaysOld int) error

	CreateRetroAction(RetroID string, UserID string, Content string) ([]*RetroAction, error)
	UpdateRetroAction(RetroID string, ActionID string, Content string, Completed bool) (Actions []*RetroAction, DeleteError error)
	DeleteRetroAction(RetroID string, userID string, ActionID string) ([]*RetroAction, error)
	GetRetroActions(RetroID string) []*RetroAction
	GetTeamRetroActions(TeamID string, Limit int, Offset int, Completed bool) ([]*RetroAction, int, error)
	RetroActionCommentAdd(RetroID string, ActionID string, UserID string, Comment string) ([]*RetroAction, error)
	RetroActionCommentEdit(RetroID string, ActionID string, CommentID string, Comment string) ([]*RetroAction, error)
	RetroActionCommentDelete(RetroID string, ActionID string, CommentID string) ([]*RetroAction, error)
	RetroActionAssigneeAdd(RetroID string, ActionID string, UserID string) ([]*RetroAction, error)
	RetroActionAssigneeDelete(RetroID string, ActionID string, UserID string) ([]*RetroAction, error)

	CreateRetroItem(RetroID string, UserID string, ItemType string, Content string) ([]*RetroItem, error)
	GroupRetroItem(RetroID string, ItemId string, GroupId string) (RetroItem, error)
	DeleteRetroItem(RetroID string, userID string, Type string, ItemID string) ([]*RetroItem, error)
	GetRetroItems(RetroID string) []*RetroItem
	GetRetroGroups(RetroID string) []*RetroGroup
	GroupNameChange(RetroID string, GroupId string, Name string) (RetroGroup, error)
	GetRetroVotes(RetroID string) []*RetroVote
	GroupUserVote(RetroID string, GroupID string, UserID string) ([]*RetroVote, error)
	GroupUserSubtractVote(RetroID string, GroupID string, UserID string) ([]*RetroVote, error)
}
