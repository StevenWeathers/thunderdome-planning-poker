package thunderdome

type TeamKudo struct {
	ID          string    `json:"id"`
	TeamID      string    `json:"teamId"`
	User        *TeamUser `json:"user"`
	TargetUser  *TeamUser `json:"targetUser"`
	Comment     string    `json:"comment"`
	KudosDate   string    `json:"kudosDate"`
	CreatedDate string    `json:"createdDate"`
	UpdatedDate string    `json:"updatedDate"`
}
