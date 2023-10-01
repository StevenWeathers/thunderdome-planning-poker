package thunderdome

import (
	"context"
	"time"
)

type Subscription struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	CustomerID  string    `json:"customer_id"`
	Active      bool      `json:"active"`
	Expires     time.Time `json:"expires"`
	CreatedDate time.Time `json:"created_date"`
	UpdatedDate time.Time `json:"updated_date"`
}

type SubscriptionDataSvc interface {
	CheckActiveSubscriber(ctx context.Context, userId string) error
	GetSubscriptionByUserID(ctx context.Context, userId string) (Subscription, error)
	GetSubscriptionByCustomerID(ctx context.Context, customerId string) (Subscription, error)
	CreateSubscription(ctx context.Context, userId string, customerId string, expires time.Time) (Subscription, error)
	UpdateSubscription(ctx context.Context, id string, active bool, expires time.Time) (Subscription, error)
}
