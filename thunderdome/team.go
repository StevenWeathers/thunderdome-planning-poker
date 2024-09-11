package thunderdome

import (
	"context"
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

type TeamDataSvc interface {
	TeamUserRole(ctx context.Context, UserID string, TeamID string) (string, error)
	TeamGet(ctx context.Context, TeamID string) (*Team, error)
	TeamListByUser(ctx context.Context, UserID string, Limit int, Offset int) []*UserTeam
	TeamCreate(ctx context.Context, UserID string, TeamName string) (*Team, error)
	TeamUpdate(ctx context.Context, TeamId string, TeamName string) (*Team, error)
	TeamAddUser(ctx context.Context, TeamID string, UserID string, Role string) (string, error)
	TeamUserList(ctx context.Context, TeamID string, Limit int, Offset int) ([]*TeamUser, int, error)
	TeamUpdateUser(ctx context.Context, TeamID string, UserID string, Role string) (string, error)
	TeamRemoveUser(ctx context.Context, TeamID string, UserID string) error
	TeamInviteUser(ctx context.Context, TeamID string, Email string, Role string) (string, error)
	TeamUserGetInviteByID(ctx context.Context, InviteID string) (TeamUserInvite, error)
	TeamDeleteUserInvite(ctx context.Context, InviteID string) error
	TeamGetUserInvites(ctx context.Context, teamId string) ([]TeamUserInvite, error)
	TeamPokerList(ctx context.Context, TeamID string, Limit int, Offset int) []*Poker
	TeamAddPoker(ctx context.Context, TeamID string, PokerID string) error
	TeamRemovePoker(ctx context.Context, TeamID string, PokerID string) error
	TeamDelete(ctx context.Context, TeamID string) error
	TeamRetroList(ctx context.Context, TeamID string, Limit int, Offset int) []*Retro
	TeamAddRetro(ctx context.Context, TeamID string, RetroID string) error
	TeamRemoveRetro(ctx context.Context, TeamID string, RetroID string) error
	TeamStoryboardList(ctx context.Context, TeamID string, Limit int, Offset int) []*Storyboard
	TeamAddStoryboard(ctx context.Context, TeamID string, StoryboardID string) error
	TeamRemoveStoryboard(ctx context.Context, TeamID string, StoryboardID string) error
	TeamList(ctx context.Context, Limit int, Offset int) ([]*Team, int)
	TeamIsSubscribed(ctx context.Context, TeamID string) (bool, error)
	GetTeamMetrics(ctx context.Context, teamID string) (*TeamMetrics, error)
}
