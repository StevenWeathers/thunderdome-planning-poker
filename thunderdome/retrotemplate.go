package thunderdome

import (
	"time"
)

// RetroTemplate is a template for a retro board
type RetroTemplate struct {
	ID              string               `json:"id" db:"id"`
	Name            string               `json:"name" db:"name"`
	Description     string               `json:"description" db:"description"`
	Format          *RetroTemplateFormat `json:"format" db:"format"`
	IsPublic        bool                 `json:"isPublic" db:"is_public"`
	DefaultTemplate bool                 `json:"defaultTemplate" db:"default_template"`
	CreatedBy       string               `json:"createdBy" db:"created_by"`
	OrganizationID  *string              `json:"organizationId" db:"organization_id"`
	TeamID          *string              `json:"teamId" db:"team_id"`
	CreatedAt       time.Time            `json:"createdAt" db:"created_at"`
	UpdatedAt       time.Time            `json:"updatedAt" db:"updated_at"`
}

type RetroTemplateFormatColumn struct {
	Name  string `json:"name"`
	Label string `json:"label"`
	Color string `json:"color"`
	Icon  string `json:"icon"`
}

// RetroTemplateFormat is the format of a retro template
type RetroTemplateFormat struct {
	Columns []RetroTemplateFormatColumn `json:"columns"`
}
