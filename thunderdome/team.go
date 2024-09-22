package thunderdome

import (
	"time"
)

type Team struct {
	Id             string    `json:"id"`
	Name           string    `json:"name"`
	OrganizationId string    `json:"organization_id"`
	DepartmentId   string    `json:"department_id"`
	Subscribed     *bool     `json:"subscribed,omitempty"`
	CreatedDate    time.Time `json:"createdDate"`
	UpdatedDate    time.Time `json:"updatedDate"`
}

type UserTeam struct {
	Team
	Role string `json:"role"`
}

type TeamUser struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	Role         string `json:"role"`
	Avatar       string `json:"avatar"`
	GravatarHash string `json:"gravatarHash"`
	PictureURL   string `json:"pictureUrl"`
}

type TeamUserInvite struct {
	InviteId    string    `json:"invite_id"`
	TeamId      string    `json:"team_id"`
	Email       string    `json:"email"`
	Role        string    `json:"role"`
	CreatedDate time.Time `json:"created_date"`
	ExpireDate  time.Time `json:"expire_date"`
}

// TeamMetrics represents the metrics for a single team
type TeamMetrics struct {
	TeamID               string `json:"team_id"`
	TeamName             string `json:"team_name"`
	OrganizationID       string `json:"organization_id"`
	OrganizationName     string `json:"organization_name"`
	DepartmentID         string `json:"department_id"`
	DepartmentName       string `json:"department_name"`
	UserCount            int    `json:"user_count"`
	PokerCount           int    `json:"poker_count"`
	RetroCount           int    `json:"retro_count"`
	StoryboardCount      int    `json:"storyboard_count"`
	TeamCheckinCount     int    `json:"team_checkin_count"`
	EstimationScaleCount int    `json:"estimation_scale_count"`
	RetroTemplateCount   int    `json:"retro_template_count"`
}

// UserTeamRoleInfo represents a team's structure and a user's roles (if any) for that team.
type UserTeamRoleInfo struct {
	UserID           string  `db:"user_id" json:"userId"`
	TeamID           string  `db:"team_id" json:"teamId"`
	TeamRole         *string `db:"team_role" json:"teamRole"`
	DepartmentID     *string `db:"department_id" json:"departmentId"`
	DepartmentRole   *string `db:"department_role" json:"departmentRole"`
	OrganizationID   *string `db:"organization_id" json:"organizationId"`
	OrganizationRole *string `db:"organization_role" json:"organizationRole"`
	AssociationLevel string  `db:"association_level" json:"associationLevel"`
}
