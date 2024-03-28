package support

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"go.uber.org/zap"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
)

// Service represents a PostgreSQL implementation of thunderdome.SupportDataSvc.
type Service struct {
	DB     *sql.DB
	Logger *otelzap.Logger
}

// CreateSupportTicket creates a support request (ticket)
func (d *Service) CreateSupportTicket(ctx context.Context, userId string, userName string, userEmail string, userQuestion string) (thunderdome.Support, error) {
	var sup thunderdome.Support

	err := d.DB.QueryRowContext(ctx, `
		INSERT INTO thunderdome.support
		(user_id, user_name, user_email, user_question)
		VALUES 
		(NULLIF($1, ''), $2, $3, $4)
		RETURNING support_id, user_id, user_name, user_email, user_question, resolved, resolved_by,
		 created_date, updated_date;
		`,
		userId, userName, userEmail, userQuestion,
	).Scan(&sup.Id, &sup.UserId, &sup.UserName, &sup.UserEmail, &sup.UserQuestion,
		&sup.Resolved, &sup.ResolvedBy, &sup.CreatedDate, &sup.UpdatedDate,
	)
	if err != nil {
		return sup, fmt.Errorf("unable to create support ticket: %v", err)
	}

	return sup, nil
}

// GetSupportTickets gets a list of support tickets
func (d *Service) GetSupportTickets(ctx context.Context, Limit int, Offset int) ([]*thunderdome.Support, int, error) {
	sup := make([]*thunderdome.Support, 0)
	count := 0

	e := d.DB.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM thunderdome.support;",
	).Scan(
		&count,
	)
	if e != nil {
		d.Logger.Ctx(ctx).Error("GetSupportTickets count query scan error", zap.Error(e))
	}

	rows, err := d.DB.QueryContext(ctx,
		`SELECT support_id, user_id, user_name, user_email, user_question, resolved, resolved_by,
		 created_date, updated_date 
		FROM thunderdome.support LIMIT $1 OFFSET $2;`,
		Limit, Offset,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sup, count, nil
		}
		d.Logger.Ctx(ctx).Error("GetSupportTickets query error", zap.Error(err))
		return sup, count, err
	}

	defer rows.Close()
	for rows.Next() {
		var st thunderdome.Support

		if err := rows.Scan(
			&st.Id,
			&st.UserId,
			&st.UserName,
			&st.UserEmail,
			&st.UserQuestion,
			&st.Resolved,
			&st.ResolvedBy,
			&st.CreatedDate,
			&st.UpdatedDate,
		); err != nil {
			d.Logger.Ctx(ctx).Error("GetSupportTickets row scan error", zap.Error(err))
		} else {
			sup = append(sup, &st)
		}
	}

	return sup, count, nil
}

// GetSupportTicketByID gets a support ticket by ID
func (d *Service) GetSupportTicketByID(ctx context.Context, id string) (thunderdome.Support, error) {
	sup := thunderdome.Support{}

	err := d.DB.QueryRowContext(ctx,
		`SELECT support_id, user_id, user_name, user_email, user_question, resolved, resolved_by,
		 created_date, updated_date 
		FROM thunderdome.support WHERE support_id = $1;`,
		id,
	).Scan(&sup.Id, &sup.UserId, &sup.UserName, &sup.UserEmail, &sup.UserQuestion,
		&sup.Resolved, &sup.ResolvedBy, &sup.CreatedDate, &sup.UpdatedDate,
	)

	if err != nil {
		d.Logger.Ctx(ctx).Error("GetSupportTicketByID query error", zap.Error(err))
		return sup, err
	}

	return sup, nil
}

// UpdateSupportTicket updates a support ticket by ID
func (d *Service) UpdateSupportTicket(ctx context.Context, id string, ticket thunderdome.Support) (thunderdome.Support, error) {
	sup := thunderdome.Support{}

	return sup, nil
}

// DeleteSupportTicket deletes a support ticket by ID
func (d *Service) DeleteSupportTicket(ctx context.Context, id string) error {
	_, err := d.DB.ExecContext(ctx,
		`DELETE FROM thunderdome.support WHERE support_id = $1;`,
		id,
	)

	if err != nil {
		return fmt.Errorf("error deleting support ticket: %v", err)
	}

	return nil
}
