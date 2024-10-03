package thunderdome

import (
	"time"
)

// Organization can be a company
type Organization struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Subscribed  *bool     `json:"subscribed,omitempty"`
	CreatedDate time.Time `json:"createdDate"`
	UpdatedDate time.Time `json:"updatedDate"`
}

type UserOrganization struct {
	Organization
	Role string `json:"role"`
}

type UserDepartment struct {
	Department
	Role string `json:"role"`
}

type OrganizationUser struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	Role         string `json:"role"`
	Avatar       string `json:"avatar"`
	GravatarHash string `json:"gravatarHash"`
	PictureURL   string `json:"pictureUrl"`
}

type OrganizationUserInvite struct {
	InviteID       string    `json:"invite_id"`
	OrganizationID string    `json:"organization_id"`
	Email          string    `json:"email"`
	Role           string    `json:"role"`
	CreatedDate    time.Time `json:"created_date"`
	ExpireDate     time.Time `json:"expire_date"`
}

type Department struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	OrganizationID string    `json:"organization_id"`
	CreatedDate    time.Time `json:"createdDate"`
	UpdatedDate    time.Time `json:"updatedDate"`
}

type DepartmentUser struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	Role         string `json:"role"`
	Avatar       string `json:"avatar"`
	GravatarHash string `json:"gravatarHash"`
	PictureURL   string `json:"pictureUrl"`
}

type DepartmentUserInvite struct {
	InviteID     string    `json:"invite_id"`
	DepartmentID string    `json:"department_id"`
	Email        string    `json:"email"`
	Role         string    `json:"role"`
	CreatedDate  time.Time `json:"created_date"`
	ExpireDate   time.Time `json:"expire_date"`
}

// OrganizationMetrics represents the metrics for a single organization
type OrganizationMetrics struct {
	OrganizationID       string `json:"organization_id"`
	OrganizationName     string `json:"organization_name"`
	DepartmentCount      int    `json:"department_count"`
	TeamCount            int    `json:"team_count"`
	RetroCount           int    `json:"retro_count"`
	PokerCount           int    `json:"poker_count"`
	StoryboardCount      int    `json:"storyboard_count"`
	TeamCheckinCount     int    `json:"team_checkin_count"`
	UserCount            int    `json:"user_count"`
	EstimationScaleCount int    `json:"estimation_scale_count"`
	RetroTemplateCount   int    `json:"retro_template_count"`
}
