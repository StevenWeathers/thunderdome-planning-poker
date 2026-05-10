package thunderdome

import "time"

// ColorLegendTemplate is a reusable storyboard color legend template.
type ColorLegendTemplate struct {
	ID             string    `json:"id" db:"id"`
	Name           string    `json:"name" db:"name"`
	Description    string    `json:"description" db:"description"`
	ColorLegend    []*Color  `json:"colorLegend" db:"color_legend"`
	IsPublic       bool      `json:"isPublic" db:"is_public"`
	CreatedBy      string    `json:"createdBy" db:"created_by"`
	OrganizationID *string   `json:"organizationId" db:"organization_id"`
	TeamID         *string   `json:"teamId" db:"team_id"`
	CreatedAt      time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt      time.Time `json:"updatedAt" db:"updated_at"`
}
