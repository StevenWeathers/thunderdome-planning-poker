package thunderdome

import "time"

// SupportTicket represents a support ticket in the system
type SupportTicket struct {
	ID         string     `json:"id"`
	UserID     *string    `json:"userId"`
	FullName   string     `json:"fullName"`
	Email      string     `json:"email"`
	Inquiry    string     `json:"inquiry"`
	AssignedTo *string    `json:"assignedTo"`
	Notes      *string    `json:"notes"`
	ResolvedAt *time.Time `json:"resolvedAt"`
	ResolvedBy *string    `json:"resolvedBy"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt"`
}
