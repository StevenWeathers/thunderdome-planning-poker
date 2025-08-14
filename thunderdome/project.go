package thunderdome

import (
	"time"
)

// Project represents a project within an organization, department, or team
type Project struct {
	ID             string    `json:"id" db:"id"`
	ProjectKey     string    `json:"projectKey" db:"project_key"`
	Name           string    `json:"name" db:"name"`
	Description    string    `json:"description" db:"description"`
	OrganizationID *string   `json:"organizationId" db:"organization_id"`
	DepartmentID   *string   `json:"departmentId" db:"department_id"`
	TeamID         *string   `json:"teamId" db:"team_id"`
	CreatedAt      time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt      time.Time `json:"updatedAt" db:"updated_at"`
}
