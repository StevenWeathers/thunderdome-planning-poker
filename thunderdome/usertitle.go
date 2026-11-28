package thunderdome

import (
	"errors"
	"time"
)

var (
	ErrDuplicateUserTitle = errors.New("DUPLICATE_USER_TITLE")
	ErrUserTitleNotFound  = errors.New("USER_TITLE_NOT_FOUND")
)

// UserTitle is a title an organization offers for its members
type UserTitle struct {
	ID             string    `json:"id"`
	OrganizationID string    `json:"organizationId"`
	Name           string    `json:"name"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}
