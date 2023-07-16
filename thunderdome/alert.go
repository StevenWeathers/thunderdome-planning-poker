package thunderdome

import (
	"context"
	"time"
)

type Alert struct {
	Id             string    `json:"id" db:"id"`
	Name           string    `json:"name" db:"name"`
	Type           string    `json:"type" db:"type"`
	Content        string    `json:"content" db:"content"`
	Active         bool      `json:"active" db:"active"`
	AllowDismiss   bool      `json:"allowDismiss" db:"allow_dismiss"`
	RegisteredOnly bool      `json:"registeredOnly" db:"registered_only"`
	CreatedDate    time.Time `json:"createdDate" db:"created_date"`
	UpdatedDate    time.Time `json:"updatedDate" db:"updated_date"`
}

type AlertDataSvc interface {
	GetActiveAlerts(ctx context.Context) []interface{}
	AlertsList(ctx context.Context, Limit int, Offset int) ([]*Alert, int, error)
	AlertsCreate(ctx context.Context, Name string, Type string, Content string, Active bool, AllowDismiss bool, RegisteredOnly bool) error
	AlertsUpdate(ctx context.Context, ID string, Name string, Type string, Content string, Active bool, AllowDismiss bool, RegisteredOnly bool) error
	AlertDelete(ctx context.Context, AlertID string) error
}
