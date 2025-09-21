package user

import (
	"context"
	"fmt"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
)

func (s *Service) CreateSupportTicket(ctx context.Context, userId, fullName, email, inquiry string) (thunderdome.SupportTicket, error) {
	// make sure user doesn't already have more than 3 pending tickets
	var existingCount int
	err := s.DB.QueryRowContext(ctx, `
		SELECT COUNT(id) FROM thunderdome.support_ticket
		WHERE user_id = $1 AND resolved_at IS NOT NULL
	`, userId).Scan(&existingCount)
	if err != nil {
		return thunderdome.SupportTicket{}, err
	}

	if existingCount >= 3 {
		return thunderdome.SupportTicket{}, fmt.Errorf("EXCEEDED_PENDING_TICKET_LIMIT")
	}

	// Create the support ticket
	ticket := thunderdome.SupportTicket{
		UserID:   &userId,
		FullName: fullName,
		Email:    email,
		Inquiry:  inquiry,
	}
	err = s.DB.QueryRowContext(ctx, `
	INSERT INTO thunderdome.support_ticket (user_id, full_name, email, inquiry)
	 VALUES ($1, $2, $3, $4)
	 RETURNING id, created_at, updated_at
	`, userId, fullName, email, inquiry).Scan(&ticket.ID, &ticket.CreatedAt, &ticket.UpdatedAt)
	if err != nil {
		return thunderdome.SupportTicket{}, err
	}

	return thunderdome.SupportTicket{
		ID:       ticket.ID,
		UserID:   ticket.UserID,
		FullName: fullName,
		Email:    email,
		Inquiry:  inquiry,
	}, nil
}
