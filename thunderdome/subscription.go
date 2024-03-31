package thunderdome

import (
	"context"
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

type SubscriptionDataSvc interface {
	CheckActiveSubscriber(ctx context.Context, userId string) error
	GetSubscriptionByID(ctx context.Context, id string) (Subscription, error)
	GetSubscriptionBySubscriptionID(ctx context.Context, subscriptionId string) (Subscription, error)
	GetActiveSubscriptionsByUserID(ctx context.Context, userId string) ([]Subscription, error)
	CreateSubscription(ctx context.Context, subscription Subscription) (Subscription, error)
	UpdateSubscription(ctx context.Context, id string, sub Subscription) (Subscription, error)
	GetSubscriptions(ctx context.Context, Limit int, Offset int) ([]Subscription, int, error)
	DeleteSubscription(ctx context.Context, id string) error
}
