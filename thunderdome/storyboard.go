package thunderdome

import "context"

// StoryboardUser aka user
type StoryboardUser struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Active       bool   `json:"active"`
	Avatar       string `json:"avatar"`
	Abandoned    bool   `json:"abandoned"`
	GravatarHash string `json:"gravatarHash"`
	PictureURL   string `json:"pictureUrl"`
}

// Storyboard A story mapping board
type Storyboard struct {
	ID              string               `json:"id"`
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
	ID        string               `json:"id"`
	Name      string               `json:"name"`
	Personas  []*StoryboardPersona `json:"personas"`
	Columns   []*StoryboardColumn  `json:"columns"`
	SortOrder string               `json:"sort_order"`
}

// StoryboardColumn A column in a storyboard goal
type StoryboardColumn struct {
	ID        string               `json:"id"`
	Name      string               `json:"name"`
	Personas  []*StoryboardPersona `json:"personas"`
	Stories   []*StoryboardStory   `json:"stories"`
	SortOrder string               `json:"sort_order"`
}

// StoryboardStory A story in a storyboard goal column
type StoryboardStory struct {
	ID          string          `json:"id"`
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
	ID          string `json:"id"`
	StoryID     string `json:"story_id"`
	UserID      string `json:"user_id"`
	Comment     string `json:"comment"`
	CreateDate  string `json:"created_date"`
	UpdatedDate string `json:"updated_date"`
}

// StoryboardPersona A storyboards personas
type StoryboardPersona struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Role        string `json:"role"`
	Description string `json:"description"`
}

type StoryboardDataSvc interface {
	CreateStoryboard(ctx context.Context, ownerID string, storyboardName string, joinCode string, facilitatorCode string) (*Storyboard, error)
	TeamCreateStoryboard(ctx context.Context, TeamID string, ownerID string, storyboardName string, joinCode string, facilitatorCode string) (*Storyboard, error)
	EditStoryboard(storyboardID string, storyboardName string, joinCode string, facilitatorCode string) error
	GetStoryboard(storyboardID string, userID string) (*Storyboard, error)
	GetStoryboardsByUser(userID string, limit int, offset int) ([]*Storyboard, int, error)
	ConfirmStoryboardFacilitator(storyboardID string, userID string) error
	GetStoryboardUsers(storyboardID string) []*StoryboardUser
	GetStoryboardPersonas(storyboardID string) []*StoryboardPersona
	GetStoryboards(limit int, offset int) ([]*Storyboard, int, error)
	GetActiveStoryboards(limit int, offset int) ([]*Storyboard, int, error)
	AddUserToStoryboard(storyboardID string, userID string) ([]*StoryboardUser, error)
	RetreatStoryboardUser(storyboardID string, userID string) []*StoryboardUser
	GetStoryboardUserActiveStatus(storyboardID string, userID string) error
	AbandonStoryboard(storyboardID string, userID string) ([]*StoryboardUser, error)
	StoryboardFacilitatorAdd(StoryboardId string, userID string) (*Storyboard, error)
	StoryboardFacilitatorRemove(StoryboardId string, userID string) (*Storyboard, error)
	GetStoryboardFacilitatorCode(storyboardID string) (string, error)
	StoryboardReviseColorLegend(storyboardID string, userID string, colorLegend string) (*Storyboard, error)
	DeleteStoryboard(storyboardID string, userID string) error
	CleanStoryboards(ctx context.Context, daysOld int) error

	AddStoryboardPersona(storyboardID string, userID string, name string, role string, description string) ([]*StoryboardPersona, error)
	UpdateStoryboardPersona(storyboardID string, userID string, personaID string, name string, role string, description string) ([]*StoryboardPersona, error)
	DeleteStoryboardPersona(storyboardID string, userID string, personaID string) ([]*StoryboardPersona, error)

	CreateStoryboardGoal(storyboardID string, userID string, goalName string) ([]*StoryboardGoal, error)
	ReviseGoalName(storyboardID string, userID string, goalID string, goalName string) ([]*StoryboardGoal, error)
	DeleteStoryboardGoal(storyboardID string, userID string, goalID string) ([]*StoryboardGoal, error)
	GetStoryboardGoals(storyboardID string) []*StoryboardGoal

	CreateStoryboardColumn(storyboardID string, goalID string, userID string) ([]*StoryboardGoal, error)
	ReviseStoryboardColumn(storyboardID string, userID string, columnID string, columnName string) ([]*StoryboardGoal, error)
	DeleteStoryboardColumn(storyboardID string, userID string, columnID string) ([]*StoryboardGoal, error)
	ColumnPersonaAdd(storyboardID string, columnID string, personaID string) ([]*StoryboardGoal, error)
	ColumnPersonaRemove(storyboardID string, columnID string, personaID string) ([]*StoryboardGoal, error)

	CreateStoryboardStory(storyboardID string, goalID string, columnID string, userID string) ([]*StoryboardGoal, error)
	ReviseStoryName(storyboardID string, userID string, storyID string, storyName string) ([]*StoryboardGoal, error)
	ReviseStoryContent(storyboardID string, userID string, storyID string, storyContent string) ([]*StoryboardGoal, error)
	ReviseStoryColor(storyboardID string, userID string, storyID string, storyColor string) ([]*StoryboardGoal, error)
	ReviseStoryPoints(storyboardID string, userID string, storyID string, points int) ([]*StoryboardGoal, error)
	ReviseStoryClosed(storyboardID string, userID string, storyID string, closed bool) ([]*StoryboardGoal, error)
	ReviseStoryLink(storyboardID string, userID string, storyID string, link string) ([]*StoryboardGoal, error)
	MoveStoryboardStory(storyboardID string, userID string, storyID string, goalID string, columnID string, placeBefore string) ([]*StoryboardGoal, error)
	DeleteStoryboardStory(storyboardID string, userID string, storyID string) ([]*StoryboardGoal, error)
	AddStoryComment(storyboardID string, userID string, storyID string, comment string) ([]*StoryboardGoal, error)
	EditStoryComment(storyboardID string, commentID string, comment string) ([]*StoryboardGoal, error)
	DeleteStoryComment(storyboardID string, commentID string) ([]*StoryboardGoal, error)
}
