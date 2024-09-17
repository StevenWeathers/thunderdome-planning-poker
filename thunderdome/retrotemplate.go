package thunderdome

import (
	"context"
	"time"
)

// RetroTemplate is a template for a retro board
type RetroTemplate struct {
	Id              string               `json:"id" db:"id"`
	Name            string               `json:"name" db:"name"`
	Description     string               `json:"description" db:"description"`
	Format          *RetroTemplateFormat `json:"format" db:"format"`
	IsPublic        bool                 `json:"isPublic" db:"is_public"`
	DefaultTemplate bool                 `json:"defaultTemplate" db:"default_template"`
	CreatedBy       string               `json:"createdBy" db:"created_by"`
	OrganizationId  *string              `json:"organizationId" db:"organization_id"`
	TeamId          *string              `json:"teamId" db:"team_id"`
	CreatedAt       time.Time            `json:"createdAt" db:"created_at"`
	UpdatedAt       time.Time            `json:"updatedAt" db:"updated_at"`
}

// RetroTemplateFormat is the format of a retro template
type RetroTemplateFormat struct {
	Columns []struct {
		Name  string `json:"name"`
		Label string `json:"label"`
		Color string `json:"color"`
		Icon  string `json:"icon"`
	} `json:"columns"`
}

type RetroTemplateDataSvc interface {
	// GetPublicTemplates retrieves all public retro templates
	GetPublicTemplates(ctx context.Context) ([]*RetroTemplate, error)
	// GetTemplatesByOrganization retrieves all templates for a specific organization
	GetTemplatesByOrganization(ctx context.Context, organizationId string) ([]*RetroTemplate, error)
	// GetTemplatesByTeam retrieves all templates for a specific team
	GetTemplatesByTeam(ctx context.Context, teamId string) ([]*RetroTemplate, error)
	// GetTemplateById retrieves a specific template by its ID
	GetTemplateById(ctx context.Context, templateId string) (*RetroTemplate, error)
	// CreateTemplate creates a new retro template
	CreateTemplate(ctx context.Context, template *RetroTemplate) error
	// UpdateTemplate updates an existing retro template
	UpdateTemplate(ctx context.Context, template *RetroTemplate) error
	// DeleteTemplate deletes a retro template by its ID
	DeleteTemplate(ctx context.Context, templateId string) error
	// ListTemplates retrieves a paginated list of templates
	ListTemplates(ctx context.Context, limit int, offset int) ([]*RetroTemplate, int, error)
	// GetDefaultPublicTemplate retrieves the default public template
	GetDefaultPublicTemplate(ctx context.Context) (*RetroTemplate, error)
	// GetDefaultTeamTemplate retrieves the default template for a given team
	GetDefaultTeamTemplate(ctx context.Context, teamID string) (*RetroTemplate, error)
	// GetDefaultOrganizationTemplate retrieves the default template for a given organization
	GetDefaultOrganizationTemplate(ctx context.Context, organizationID string) (*RetroTemplate, error)
	// UpdateTeamTemplate updates an existing team retro template
	UpdateTeamTemplate(ctx context.Context, template *RetroTemplate) error
	// UpdateOrganizationTemplate updates an existing organization retro template
	UpdateOrganizationTemplate(ctx context.Context, template *RetroTemplate) error
	// DeleteOrganizationTemplate deletes an organization retro template by its ID
	DeleteOrganizationTemplate(ctx context.Context, orgId string, templateId string) error
	// DeleteTeamTemplate deletes a team retro template by its ID
	DeleteTeamTemplate(ctx context.Context, teamId string, templateId string) error
}
