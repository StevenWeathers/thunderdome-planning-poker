package thunderdome

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
