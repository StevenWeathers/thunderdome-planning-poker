package thunderdome

import "errors"

// StoryboardStoryPointsMaxLength is the number of characters a storyboard story's points can hold
const StoryboardStoryPointsMaxLength = 8

// ErrStoryboardEstimationScaleNotFound is returned when an estimation scale doesn't exist or isn't available to a storyboard
var ErrStoryboardEstimationScaleNotFound = errors.New("ESTIMATION_SCALE_NOT_FOUND")

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
	// EstimationScaleID is empty when story points are free text
	EstimationScaleID string           `json:"estimationScaleId,omitempty" db:"estimation_scale_id"`
	EstimationScale   *EstimationScale `json:"estimationScale,omitempty"`
}

// StoryboardGoal A row in a story mapping board
type StoryboardGoal struct {
	ID                string               `json:"id"`
	Name              string               `json:"name"`
	DefaultStoryColor *string              `json:"default_story_color"`
	Personas          []*StoryboardPersona `json:"personas"`
	Columns           []*StoryboardColumn  `json:"columns"`
	SortOrder         string               `json:"sort_order"`
}

// StoryboardColumn A column in a storyboard goal
type StoryboardColumn struct {
	ID                string               `json:"id"`
	Name              string               `json:"name"`
	DefaultStoryColor *string              `json:"default_story_color"`
	Personas          []*StoryboardPersona `json:"personas"`
	Stories           []*StoryboardStory   `json:"stories"`
	SortOrder         string               `json:"sort_order"`
}

// StoryboardStory A story in a storyboard goal column
type StoryboardStory struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	Content     string          `json:"content"`
	Color       string          `json:"color"`
	Points      string          `json:"points"`
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
