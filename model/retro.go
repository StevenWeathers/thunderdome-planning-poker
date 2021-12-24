package model

// Color is a color legend
type Color struct {
	Color  string `json:"color"`
	Legend string `json:"legend"`
}

// RetroUser aka user
type RetroUser struct {
	UserID       string `json:"id"`
	UserName     string `json:"name"`
	Active       bool   `json:"active"`
	Avatar       string `json:"avatar"`
	GravatarHash string `json:"gravatarHash"`
}

// Retro A story mapping board
type Retro struct {
	Id          string         `json:"id" db:"id"`
	OwnerID     string         `json:"ownerId" db:"owner_id"`
	Name        string         `json:"name" db:"name"`
	Users       []*RetroUser   `json:"users"`
	Groups      []*RetroGroup  `json:"groups"`
	Items       []*RetroItem   `json:"items"`
	ActionItems []*RetroAction `json:"actionItems"`
	Votes       []*RetroVote   `json:"votes"`
	Format      string         `json:"format" db:"format"`
	Phase       string         `json:"phase" db:"phase"`
	JoinCode    string         `json:"joinCode" db:"join_code"`
	CreatedDate string         `json:"createdDate" db:"created_date"`
	UpdatedDate string         `json:"updatedDate" db:"updated_date"`
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
	ID        string `json:"id" db:"id"`
	Content   string `json:"content" db:"content"`
	Completed bool   `json:"completed" db:"completed"`
}

// RetroVote is a users vote toward a retro item group
type RetroVote struct {
	UserID  string `json:"userId" db:"user_id"`
	GroupID string `json:"groupId" db:"group_id"`
}
