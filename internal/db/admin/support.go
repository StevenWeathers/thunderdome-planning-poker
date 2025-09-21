package admin

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
)

// GetSupportTicketByID retrieves a support ticket by its ID
func (d *Service) GetSupportTicketByID(ctx context.Context, ticketID string) (*thunderdome.SupportTicket, error) {
	ticket := &thunderdome.SupportTicket{}
	query := `SELECT id, user_id, full_name, email, inquiry, assigned_to, notes, resolved_at, resolved_by, created_at, updated_at FROM thunderdome.support_ticket WHERE id = $1`
	err := d.DB.QueryRowContext(ctx, query, ticketID).Scan(
		&ticket.ID,
		&ticket.UserID,
		&ticket.FullName,
		&ticket.Email,
		&ticket.Inquiry,
		&ticket.AssignedTo,
		&ticket.Notes,
		&ticket.ResolvedAt,
		&ticket.ResolvedBy,
		&ticket.CreatedAt,
		&ticket.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get support ticket by id error: %w", err)
	}
	return ticket, nil
}

// UpdateSupportTicket updates an existing support ticket
func (d *Service) UpdateSupportTicket(ctx context.Context, ticket *thunderdome.SupportTicket) error {
	query := `UPDATE thunderdome.support_ticket SET full_name = $1, email = $2, inquiry = $3, assigned_to = $4, notes = $5, resolved_at = $6, resolved_by = $7, updated_at = NOW() WHERE id = $8`
	_, err := d.DB.ExecContext(ctx, query,
		ticket.FullName,
		ticket.Email,
		ticket.Inquiry,
		ticket.AssignedTo,
		ticket.Notes,
		ticket.ResolvedAt,
		ticket.ResolvedBy,
		ticket.ID,
	)
	if err != nil {
		return fmt.Errorf("update support ticket error: %w", err)
	}
	return nil
}

// DeleteSupportTicket deletes a support ticket by its ID
func (d *Service) DeleteSupportTicket(ctx context.Context, ticketID string) error {
	query := `DELETE FROM thunderdome.support_ticket WHERE id = $1`
	_, err := d.DB.ExecContext(ctx, query, ticketID)
	if err != nil {
		return fmt.Errorf("delete support ticket error: %w", err)
	}
	return nil
}

// ListSupportTickets returns a list of support tickets with pagination
func (d *Service) ListSupportTickets(ctx context.Context, limit, offset int) ([]*thunderdome.SupportTicket, int, error) {
	tickets := make([]*thunderdome.SupportTicket, 0)
	count := 0
	countQuery := `SELECT COUNT(*) FROM thunderdome.support_ticket`
	err := d.DB.QueryRowContext(ctx, countQuery).Scan(&count)
	if err != nil {
		return nil, 0, fmt.Errorf("list support tickets count error: %w", err)
	}

	query := `SELECT id, user_id, full_name, email, inquiry, assigned_to, notes, resolved_at, resolved_by, created_at, updated_at FROM thunderdome.support_ticket ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	rows, err := d.DB.QueryContext(ctx, query, limit, offset)
	if err != nil {
		if err == sql.ErrNoRows {
			return tickets, count, nil
		}
		return nil, 0, fmt.Errorf("list support tickets error: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		ticket := &thunderdome.SupportTicket{}
		err := rows.Scan(
			&ticket.ID,
			&ticket.UserID,
			&ticket.FullName,
			&ticket.Email,
			&ticket.Inquiry,
			&ticket.AssignedTo,
			&ticket.Notes,
			&ticket.ResolvedAt,
			&ticket.ResolvedBy,
			&ticket.CreatedAt,
			&ticket.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("list support tickets scan error: %w", err)
		}
		tickets = append(tickets, ticket)
	}
	return tickets, count, nil
}
