package thunderdome

import "context"

// StoryboardUser aka user
type StoryboardUser struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	Active       bool   `json:"active"`
	Avatar       string `json:"avatar"`
	Abandoned    bool   `json:"abandoned"`
	GravatarHash string `json:"gravatarHash"`
	PictureURL   string `json:"pictureUrl"`
}

// Storyboard A story mapping board
type Storyboard struct {
	Id              string               `json:"id"`
	OwnerID         string               `json:"owner_id"`
	Name            string               `json:"name"`
	Users           []*StoryboardUser    `json:"users"`
	Facilitators    []string             `json:"facilitators"`
	Goals           []*StoryboardGoal    `json:"goals"`
	ColorLegend     []*Color             `json:"color_legend"`
	Personas        []*StoryboardPersona `json:"personas"`
	JoinCode        string               `json:"joinCode" db:"join_code"`
	FacilitatorCode string               `json:"facilitatorCode" db:"facilitator_code"`
	TeamID          string               `json:"teamId" db:"team_id"`
	TeamName        string               `json:"teamName"`
	CreatedDate     string               `json:"createdDate" db:"created_date"`
	UpdatedDate     string               `json:"updatedDate" db:"updated_date"`
}

// StoryboardGoal A row in a story mapping board
type StoryboardGoal struct {
	Id        string               `json:"id"`
	Name      string               `json:"name"`
	Personas  []*StoryboardPersona `json:"personas"`
	Columns   []*StoryboardColumn  `json:"columns"`
	SortOrder string               `json:"sort_order"`
}

// StoryboardColumn A column in a storyboard goal
type StoryboardColumn struct {
	Id        string               `json:"id"`
	Name      string               `json:"name"`
	Personas  []*StoryboardPersona `json:"personas"`
	Stories   []*StoryboardStory   `json:"stories"`
	SortOrder string               `json:"sort_order"`
}

// StoryboardStory A story in a storyboard goal column
type StoryboardStory struct {
	Id          string          `json:"id"`
	Name        string          `json:"name"`
	Content     string          `json:"content"`
	Color       string          `json:"color"`
	Points      int             `json:"points"`
	Closed      bool            `json:"closed"`
	Link        string          `json:"link"`
	Annotations []string        `json:"annotations"`
	SortOrder   string          `json:"sort_order"`
	Comments    []*StoryComment `json:"comments"`
}

// StoryComment A story comment by a user
type StoryComment struct {
	Id          string `json:"id"`
	StoryID     string `json:"story_id"`
	UserID      string `json:"user_id"`
	Comment     string `json:"comment"`
	CreateDate  string `json:"created_date"`
	UpdatedDate string `json:"updated_date"`
}

// StoryboardPersona A storyboards personas
type StoryboardPersona struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Role        string `json:"role"`
	Description string `json:"description"`
}

type StoryboardDataSvc interface {
	CreateStoryboard(ctx context.Context, OwnerID string, StoryboardName string, JoinCode string, FacilitatorCode string) (*Storyboard, error)
	TeamCreateStoryboard(ctx context.Context, TeamID string, OwnerID string, StoryboardName string, JoinCode string, FacilitatorCode string) (*Storyboard, error)
	EditStoryboard(StoryboardID string, StoryboardName string, JoinCode string, FacilitatorCode string) error
	GetStoryboard(StoryboardID string, UserID string) (*Storyboard, error)
	GetStoryboardsByUser(UserID string, Limit int, Offset int) ([]*Storyboard, int, error)
	ConfirmStoryboardFacilitator(StoryboardID string, UserID string) error
	GetStoryboardUsers(StoryboardID string) []*StoryboardUser
	GetStoryboardPersonas(StoryboardID string) []*StoryboardPersona
	GetStoryboards(Limit int, Offset int) ([]*Storyboard, int, error)
	GetActiveStoryboards(Limit int, Offset int) ([]*Storyboard, int, error)
	AddUserToStoryboard(StoryboardID string, UserID string) ([]*StoryboardUser, error)
	RetreatStoryboardUser(StoryboardID string, UserID string) []*StoryboardUser
	GetStoryboardUserActiveStatus(StoryboardID string, UserID string) error
	AbandonStoryboard(StoryboardID string, UserID string) ([]*StoryboardUser, error)
	StoryboardFacilitatorAdd(StoryboardId string, UserID string) (*Storyboard, error)
	StoryboardFacilitatorRemove(StoryboardId string, UserID string) (*Storyboard, error)
	GetStoryboardFacilitatorCode(StoryboardID string) (string, error)
	StoryboardReviseColorLegend(StoryboardID string, UserID string, ColorLegend string) (*Storyboard, error)
	DeleteStoryboard(StoryboardID string, userID string) error
	CleanStoryboards(ctx context.Context, DaysOld int) error

	AddStoryboardPersona(StoryboardID string, UserID string, Name string, Role string, Description string) ([]*StoryboardPersona, error)
	UpdateStoryboardPersona(StoryboardID string, UserID string, PersonaID string, Name string, Role string, Description string) ([]*StoryboardPersona, error)
	DeleteStoryboardPersona(StoryboardID string, UserID string, PersonaID string) ([]*StoryboardPersona, error)

	CreateStoryboardGoal(StoryboardID string, userID string, GoalName string) ([]*StoryboardGoal, error)
	ReviseGoalName(StoryboardID string, userID string, GoalID string, GoalName string) ([]*StoryboardGoal, error)
	DeleteStoryboardGoal(StoryboardID string, userID string, GoalID string) ([]*StoryboardGoal, error)
	GetStoryboardGoals(StoryboardID string) []*StoryboardGoal

	CreateStoryboardColumn(StoryboardID string, GoalID string, userID string) ([]*StoryboardGoal, error)
	ReviseStoryboardColumn(StoryboardID string, UserID string, ColumnID string, ColumnName string) ([]*StoryboardGoal, error)
	DeleteStoryboardColumn(StoryboardID string, userID string, ColumnID string) ([]*StoryboardGoal, error)
	ColumnPersonaAdd(StoryboardID string, ColumnID string, PersonaID string) ([]*StoryboardGoal, error)
	ColumnPersonaRemove(StoryboardID string, ColumnID string, PersonaID string) ([]*StoryboardGoal, error)

	CreateStoryboardStory(StoryboardID string, GoalID string, ColumnID string, userID string) ([]*StoryboardGoal, error)
	ReviseStoryName(StoryboardID string, userID string, StoryID string, StoryName string) ([]*StoryboardGoal, error)
	ReviseStoryContent(StoryboardID string, userID string, StoryID string, StoryContent string) ([]*StoryboardGoal, error)
	ReviseStoryColor(StoryboardID string, userID string, StoryID string, StoryColor string) ([]*StoryboardGoal, error)
	ReviseStoryPoints(StoryboardID string, userID string, StoryID string, Points int) ([]*StoryboardGoal, error)
	ReviseStoryClosed(StoryboardID string, userID string, StoryID string, Closed bool) ([]*StoryboardGoal, error)
	ReviseStoryLink(StoryboardID string, userID string, StoryID string, Link string) ([]*StoryboardGoal, error)
	MoveStoryboardStory(StoryboardID string, userID string, StoryID string, GoalID string, ColumnID string, PlaceBefore string) ([]*StoryboardGoal, error)
	DeleteStoryboardStory(StoryboardID string, userID string, StoryID string) ([]*StoryboardGoal, error)
	AddStoryComment(StoryboardID string, UserID string, StoryID string, Comment string) ([]*StoryboardGoal, error)
	EditStoryComment(StoryboardID string, CommentID string, Comment string) ([]*StoryboardGoal, error)
	DeleteStoryComment(StoryboardID string, CommentID string) ([]*StoryboardGoal, error)
}
