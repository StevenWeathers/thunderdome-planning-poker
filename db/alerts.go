package db

import (
	"context"
	"database/sql"
	"errors"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"

	"go.uber.org/zap"
)

// AlertService represents a PostgreSQL implementation of thunderdome.AlertDataSvc.
type AlertService struct {
	DB     *sql.DB
	Logger *otelzap.Logger
}

// GetActiveAlerts gets a list of active global alerts
func (d *AlertService) GetActiveAlerts(ctx context.Context) []interface{} {
	Alerts := make([]interface{}, 0)

	rows, err := d.DB.QueryContext(ctx,
		`SELECT id, name, type, content, active, allow_dismiss, registered_only FROM thunderdome.alert WHERE active IS TRUE;`,
	)

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var a thunderdome.Alert

			if err := rows.Scan(
				&a.Id,
				&a.Name,
				&a.Type,
				&a.Content,
				&a.Active,
				&a.AllowDismiss,
				&a.RegisteredOnly,
			); err != nil {
				d.Logger.Ctx(ctx).Error("query scan error", zap.Error(err))
			} else {
				Alerts = append(Alerts, &a)
			}
		}
	}

	return Alerts
}

// AlertsList gets a list of global alerts
func (d *AlertService) AlertsList(ctx context.Context, Limit int, Offset int) ([]*thunderdome.Alert, int, error) {
	Alerts := make([]*thunderdome.Alert, 0)
	var AlertCount int

	e := d.DB.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM thunderdome.alert;",
	).Scan(
		&AlertCount,
	)
	if e != nil {
		d.Logger.Ctx(ctx).Error("query scan error", zap.Error(e))
	}

	rows, err := d.DB.QueryContext(ctx,
		`SELECT id, name, type, content, active, allow_dismiss, registered_only, created_date, updated_date
		FROM thunderdome.alert
		LIMIT $1
		OFFSET $2;
		`,
		Limit,
		Offset,
	)

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var a thunderdome.Alert

			if err := rows.Scan(
				&a.Id,
				&a.Name,
				&a.Type,
				&a.Content,
				&a.Active,
				&a.AllowDismiss,
				&a.RegisteredOnly,
				&a.CreatedDate,
				&a.UpdatedDate,
			); err != nil {
				d.Logger.Ctx(ctx).Error("query scan error", zap.Error(err))
				return nil, AlertCount, err
			} else {
				Alerts = append(Alerts, &a)
			}
		}
	}

	return Alerts, AlertCount, err
}

// AlertsCreate creates a global alert
func (d *AlertService) AlertsCreate(ctx context.Context, Name string, Type string, Content string, Active bool, AllowDismiss bool, RegisteredOnly bool) error {
	if _, err := d.DB.ExecContext(ctx,
		`INSERT INTO thunderdome.alert (name, type, content, active, allow_dismiss, registered_only)
		VALUES ($1, $2, $3, $4, $5, $6);
		`,
		Name,
		Type,
		Content,
		Active,
		AllowDismiss,
		RegisteredOnly,
	); err != nil {
		d.Logger.Ctx(ctx).Error("insert error", zap.Error(err))
		return errors.New("error attempting to add new alert")
	}

	return nil
}

// AlertsUpdate updates a global alert
func (d *AlertService) AlertsUpdate(ctx context.Context, ID string, Name string, Type string, Content string, Active bool, AllowDismiss bool, RegisteredOnly bool) error {
	if _, err := d.DB.ExecContext(ctx,
		`
		UPDATE thunderdome.alert
		SET name = $2, type = $3, content = $4, active = $5, allow_dismiss = $6, registered_only = $7
		WHERE id = $1;
		`,
		ID,
		Name,
		Type,
		Content,
		Active,
		AllowDismiss,
		RegisteredOnly,
	); err != nil {
		d.Logger.Ctx(ctx).Error("update error", zap.Error(err))
		return errors.New("error attempting to update alert")
	}

	return nil
}

// AlertDelete deletes a global alert
func (d *AlertService) AlertDelete(ctx context.Context, AlertID string) error {
	_, err := d.DB.ExecContext(ctx,
		`DELETE FROM thunderdome.alert WHERE id = $1;`,
		AlertID,
	)

	if err != nil {
		d.Logger.Ctx(ctx).Error("Unable to delete alert", zap.Error(err))
		return err
	}

	return nil
}
