package thunderdome

import (
	"context"
	"time"
)

type Support struct {
	Id           string    `json:"id" db:"support_id"`
	UserId       string    `json:"user_id" db:"user_id"`
	UserName     string    `json:"user_name" db:"user_name"`
	UserEmail    string    `json:"user_email" db:"user_email"`
	UserQuestion string    `json:"user_question" db:"user_question"`
	Resolved     bool      `json:"resolved" db:"resolved"`
	ResolvedBy   string    `json:"resolved_by" db:"resolved_by"`
	CreatedDate  time.Time `json:"created_date" db:"created_date"`
	UpdatedDate  time.Time `json:"updated_date" db:"updated_date"`
}

type SupportDataSvc interface {
	CreateSupportTicket(ctx context.Context, userId string, userName string, userEmail string, userQuestion string) (Support, error)
	GetSupportTickets(ctx context.Context, Limit int, Offset int) ([]*Support, int, error)
	GetSupportTicketByID(ctx context.Context, id string) (Support, error)
	UpdateSupportTicket(ctx context.Context, id string, ticket Support) (Support, error)
	DeleteSupportTicket(ctx context.Context, id string) error
}
