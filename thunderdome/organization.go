package thunderdome

import (
	"context"
	"time"
)

// Organization can be a company
type Organization struct {
	Id          string    `json:"id"`
	Name        string    `json:"name"`
	CreatedDate time.Time `json:"createdDate"`
	UpdatedDate time.Time `json:"updatedDate"`
}

type OrganizationUser struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	Role         string `json:"role"`
	Avatar       string `json:"avatar"`
	GravatarHash string `json:"gravatarHash"`
	PictureURL   string `json:"pictureUrl"`
}

type OrganizationUserInvite struct {
	InviteId       string    `json:"invite_id"`
	OrganizationId string    `json:"organization_id"`
	Email          string    `json:"email"`
	Role           string    `json:"role"`
	CreatedDate    time.Time `json:"created_date"`
	ExpireDate     time.Time `json:"expire_date"`
}

type Department struct {
	Id          string    `json:"id"`
	Name        string    `json:"name"`
	CreatedDate time.Time `json:"createdDate"`
	UpdatedDate time.Time `json:"updatedDate"`
}

type DepartmentUser struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	Role         string `json:"role"`
	Avatar       string `json:"avatar"`
	GravatarHash string `json:"gravatarHash"`
	PictureURL   string `json:"pictureUrl"`
}

type OrganizationDataSvc interface {
	OrganizationGet(ctx context.Context, OrgID string) (*Organization, error)
	OrganizationUserRole(ctx context.Context, UserID string, OrgID string) (string, error)
	OrganizationListByUser(ctx context.Context, UserID string, Limit int, Offset int) []*Organization
	OrganizationCreate(ctx context.Context, UserID string, OrgName string) (*Organization, error)
	OrganizationUserList(ctx context.Context, OrgID string, Limit int, Offset int) []*OrganizationUser
	OrganizationAddUser(ctx context.Context, OrgID string, UserID string, Role string) (string, error)
	OrganizationUpdateUser(ctx context.Context, OrgID string, UserID string, Role string) (string, error)
	OrganizationRemoveUser(ctx context.Context, OrganizationID string, UserID string) error
	OrganizationInviteUser(ctx context.Context, OrgID string, Email string, Role string) (string, error)
	OrganizationUserGetInviteByID(ctx context.Context, InviteID string) (OrganizationUserInvite, error)
	OrganizationDeleteUserInvite(ctx context.Context, InviteID string) error
	OrganizationGetUserInvites(ctx context.Context, orgId string) ([]OrganizationUserInvite, error)
	OrganizationTeamList(ctx context.Context, OrgID string, Limit int, Offset int) []*Team
	OrganizationTeamCreate(ctx context.Context, OrgID string, TeamName string) (*Team, error)
	OrganizationTeamUserRole(ctx context.Context, UserID string, OrgID string, TeamID string) (string, string, error)
	OrganizationDelete(ctx context.Context, OrgID string) error
	OrganizationList(ctx context.Context, Limit int, Offset int) []*Organization

	DepartmentUserRole(ctx context.Context, UserID string, OrgID string, DepartmentID string) (string, string, error)
	DepartmentGet(ctx context.Context, DepartmentID string) (*Department, error)
	OrganizationDepartmentList(ctx context.Context, OrgID string, Limit int, Offset int) []*Department
	DepartmentCreate(ctx context.Context, OrgID string, OrgName string) (*Department, error)
	DepartmentTeamList(ctx context.Context, DepartmentID string, Limit int, Offset int) []*Team
	DepartmentTeamCreate(ctx context.Context, DepartmentID string, TeamName string) (*Team, error)
	DepartmentUserList(ctx context.Context, DepartmentID string, Limit int, Offset int) []*DepartmentUser
	DepartmentAddUser(ctx context.Context, DepartmentID string, UserID string, Role string) (string, error)
	DepartmentUpdateUser(ctx context.Context, DepartmentID string, UserID string, Role string) (string, error)
	DepartmentRemoveUser(ctx context.Context, DepartmentID string, UserID string) error
	DepartmentTeamUserRole(ctx context.Context, UserID string, OrgID string, DepartmentID string, TeamID string) (string, string, string, error)
	DepartmentDelete(ctx context.Context, DepartmentID string) error
}
