package alert

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"

	"go.uber.org/zap"
)

// Service represents the alert database service
type Service struct {
	DB     *sql.DB
	Logger *otelzap.Logger
}

// GetActiveAlerts gets a list of active global alerts
func (d *Service) GetActiveAlerts(ctx context.Context) []any {
	alerts := make([]any, 0)

	rows, err := d.DB.QueryContext(ctx,
		`SELECT id, name, type, content, active, allow_dismiss, registered_only FROM thunderdome.alert WHERE active IS TRUE;`,
	)

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var a thunderdome.Alert

			if err := rows.Scan(
				&a.ID,
				&a.Name,
				&a.Type,
				&a.Content,
				&a.Active,
				&a.AllowDismiss,
				&a.RegisteredOnly,
			); err != nil {
				d.Logger.Ctx(ctx).Error("GetActiveAlerts row scan error", zap.Error(err))
			} else {
				alerts = append(alerts, &a)
			}
		}
	}

	return alerts
}

// AlertsList gets a list of global alerts
func (d *Service) AlertsList(ctx context.Context, limit int, offset int) ([]*thunderdome.Alert, int, error) {
	alerts := make([]*thunderdome.Alert, 0)
	var alertCount int

	e := d.DB.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM thunderdome.alert;",
	).Scan(
		&alertCount,
	)
	if e != nil {
		d.Logger.Ctx(ctx).Error("AlertsList query scan error", zap.Error(e))
	}

	rows, err := d.DB.QueryContext(ctx,
		`SELECT id, name, type, content, active, allow_dismiss, registered_only, created_date, updated_date
		FROM thunderdome.alert
		LIMIT $1
		OFFSET $2;
		`,
		limit,
		offset,
	)

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var a thunderdome.Alert

			if err := rows.Scan(
				&a.ID,
				&a.Name,
				&a.Type,
				&a.Content,
				&a.Active,
				&a.AllowDismiss,
				&a.RegisteredOnly,
				&a.CreatedDate,
				&a.UpdatedDate,
			); err != nil {
				return nil, alertCount, fmt.Errorf("AlertsList row scan error: %v", err)
			} else {
				alerts = append(alerts, &a)
			}
		}
	}

	return alerts, alertCount, err
}

// AlertsCreate creates a global alert
func (d *Service) AlertsCreate(ctx context.Context, name string, alertType string, content string, active bool, allowDismiss bool, registeredOnly bool) error {
	if _, err := d.DB.ExecContext(ctx,
		`INSERT INTO thunderdome.alert (name, type, content, active, allow_dismiss, registered_only)
		VALUES ($1, $2, $3, $4, $5, $6);
		`,
		name,
		alertType,
		content,
		active,
		allowDismiss,
		registeredOnly,
	); err != nil {
		return fmt.Errorf("error creating new alert: %v", err)
	}

	return nil
}

// AlertsUpdate updates a global alert
func (d *Service) AlertsUpdate(ctx context.Context, alertID string, name string, alertType string, content string, active bool, allowDismiss bool, registeredOnly bool) error {
	if _, err := d.DB.ExecContext(ctx,
		`
		UPDATE thunderdome.alert
		SET name = $2, type = $3, content = $4, active = $5, allow_dismiss = $6, registered_only = $7
		WHERE id = $1;
		`,
		alertID,
		name,
		alertType,
		content,
		active,
		allowDismiss,
		registeredOnly,
	); err != nil {
		return fmt.Errorf("error updating alert: %v", err)
	}

	return nil
}

// AlertDelete deletes a global alert
func (d *Service) AlertDelete(ctx context.Context, alertID string) error {
	_, err := d.DB.ExecContext(ctx,
		`DELETE FROM thunderdome.alert WHERE id = $1;`,
		alertID,
	)

	if err != nil {
		return fmt.Errorf("error deleting alert: %v", err)
	}

	return nil
}
