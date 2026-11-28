package thunderdome

import (
	"errors"
	"time"
)

// ErrPersonaTemplateNotFound is returned when a persona template doesn't exist or isn't available to the requester
var ErrPersonaTemplateNotFound = errors.New("PERSONA_TEMPLATE_NOT_FOUND")

// PersonaTemplatePersona is a persona that a template adds to a storyboard
type PersonaTemplatePersona struct {
	Name        string `json:"name"`
	Role        string `json:"role"`
	Description string `json:"description"`
}

// PersonaTemplate is a reusable set of storyboard personas.
type PersonaTemplate struct {
	ID             string                    `json:"id"`
	Name           string                    `json:"name"`
	Description    string                    `json:"description"`
	Personas       []*PersonaTemplatePersona `json:"personas"`
	IsPublic       bool                      `json:"isPublic"`
	CreatedBy      string                    `json:"createdBy"`
	OrganizationID *string                   `json:"organizationId"`
	TeamID         *string                   `json:"teamId"`
	CreatedAt      time.Time                 `json:"createdAt"`
	UpdatedAt      time.Time                 `json:"updatedAt"`
}
