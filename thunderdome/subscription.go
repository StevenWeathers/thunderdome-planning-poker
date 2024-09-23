package thunderdome

import (
	"time"
)

type Subscription struct {
	ID             string    `json:"id"`
	UserID         string    `json:"user_id"`
	TeamID         string    `json:"team_id"`
	OrganizationID string    `json:"organization_id"`
	CustomerID     string    `json:"customer_id"`
	SubscriptionID string    `json:"subscription_id"`
	Active         bool      `json:"active"`
	Expires        time.Time `json:"expires"`
	Type           string    `json:"type"`
	CreatedDate    time.Time `json:"created_date"`
	UpdatedDate    time.Time `json:"updated_date"`
	User           User      `json:"user"`
}
