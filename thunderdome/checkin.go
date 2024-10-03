package thunderdome

type TeamCheckin struct {
	ID          string            `json:"id"`
	User        *TeamUser         `json:"user"`
	Yesterday   string            `json:"yesterday"`
	Today       string            `json:"today"`
	Blockers    string            `json:"blockers"`
	Discuss     string            `json:"discuss"`
	GoalsMet    bool              `json:"goalsMet"`
	CreatedDate string            `json:"createdDate"`
	UpdatedDate string            `json:"updatedDate"`
	Comments    []*CheckinComment `json:"comments"`
}

// CheckinComment A checkin comment by a user
type CheckinComment struct {
	ID          string `json:"id"`
	CheckinID   string `json:"checkin_id"`
	UserID      string `json:"user_id"`
	Comment     string `json:"comment"`
	CreateDate  string `json:"created_date"`
	UpdatedDate string `json:"updated_date"`
}
