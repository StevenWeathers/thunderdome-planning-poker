package thunderdome

import (
	"errors"
	"strings"
	"time"
	"unicode/utf8"
)

const (
	// RetroTagMaxLength is the maximum number of characters in a retro tag
	RetroTagMaxLength = 64
	// RetroItemMaxTags is the maximum number of tags on a single retro item
	RetroItemMaxTags = 10
)

var (
	// ErrInvalidRetroTag is returned when a tag is empty or too long
	ErrInvalidRetroTag = errors.New("INVALID_RETRO_TAG")
	// ErrDuplicateRetroTag is returned when a preset tag name already exists in its scope
	ErrDuplicateRetroTag = errors.New("DUPLICATE_RETRO_TAG")
	// ErrRetroTagNotFound is returned when a preset tag doesn't exist in the requested scope
	ErrRetroTagNotFound = errors.New("RETRO_TAG_NOT_FOUND")
)

// RetroTag is a preset tag that can be applied to retro items.
// A tag with neither an OrganizationID nor a TeamID is a global default.
type RetroTag struct {
	ID             string    `json:"id" db:"id"`
	Name           string    `json:"name" db:"name"`
	IsPublic       bool      `json:"isPublic" db:"is_public"`
	CreatedBy      string    `json:"createdBy" db:"created_by"`
	OrganizationID *string   `json:"organizationId" db:"organization_id"`
	TeamID         *string   `json:"teamId" db:"team_id"`
	CreatedAt      time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt      time.Time `json:"updatedAt" db:"updated_at"`
}

// NormalizeRetroTag collapses whitespace and validates the length of a tag
func NormalizeRetroTag(tag string) (string, error) {
	normalized := strings.Join(strings.Fields(strings.ToValidUTF8(tag, "")), " ")
	if normalized == "" || utf8.RuneCountInString(normalized) > RetroTagMaxLength {
		return "", ErrInvalidRetroTag
	}

	return normalized, nil
}
