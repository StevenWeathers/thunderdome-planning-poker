package db

import (
	"context"
	"errors"

	"github.com/StevenWeathers/thunderdome-planning-poker/model"
	"go.uber.org/zap"
)

// GetActiveAlerts gets a list of active global alerts
func (d *Database) GetActiveAlerts(ctx context.Context) []interface{} {
	Alerts := make([]interface{}, 0)

	rows, err := d.db.QueryContext(ctx,
		`SELECT id, name, type, content, active, allow_dismiss, registered_only FROM alert WHERE active IS TRUE;`,
	)

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var a model.Alert

			if err := rows.Scan(
				&a.Id,
				&a.Name,
				&a.Type,
				&a.Content,
				&a.Active,
				&a.AllowDismiss,
				&a.RegisteredOnly,
			); err != nil {
				d.logger.Ctx(ctx).Error("query scan error", zap.Error(err))
			} else {
				Alerts = append(Alerts, &a)
			}
		}
	}

	return Alerts
}

// AlertsList gets a list of global alerts
func (d *Database) AlertsList(ctx context.Context, Limit int, Offset int) ([]*model.Alert, int, error) {
	Alerts := make([]*model.Alert, 0)
	var AlertCount int

	e := d.db.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM alert;",
	).Scan(
		&AlertCount,
	)
	if e != nil {
		d.logger.Ctx(ctx).Error("query scan error", zap.Error(e))
	}

	rows, err := d.db.QueryContext(ctx,
		`SELECT id, name, type, content, active, allow_dismiss, registered_only, created_date, updated_date
		FROM alert
		LIMIT $1
		OFFSET $2;
		`,
		Limit,
		Offset,
	)

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var a model.Alert

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
				d.logger.Ctx(ctx).Error("query scan error", zap.Error(err))
				return nil, AlertCount, err
			} else {
				Alerts = append(Alerts, &a)
			}
		}
	}

	return Alerts, AlertCount, err
}

// AlertsCreate creates a global alert
func (d *Database) AlertsCreate(ctx context.Context, Name string, Type string, Content string, Active bool, AllowDismiss bool, RegisteredOnly bool) error {
	if _, err := d.db.ExecContext(ctx,
		`INSERT INTO alert (name, type, content, active, allow_dismiss, registered_only)
		VALUES ($1, $2, $3, $4, $5, $6);
		`,
		Name,
		Type,
		Content,
		Active,
		AllowDismiss,
		RegisteredOnly,
	); err != nil {
		d.logger.Ctx(ctx).Error("insert error", zap.Error(err))
		return errors.New("error attempting to add new alert")
	}

	return nil
}

// AlertsUpdate updates a global alert
func (d *Database) AlertsUpdate(ctx context.Context, ID string, Name string, Type string, Content string, Active bool, AllowDismiss bool, RegisteredOnly bool) error {
	if _, err := d.db.ExecContext(ctx,
		`
		UPDATE alert
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
		d.logger.Ctx(ctx).Error("update error", zap.Error(err))
		return errors.New("error attempting to update alert")
	}

	return nil
}

// AlertDelete deletes a global alert
func (d *Database) AlertDelete(ctx context.Context, AlertID string) error {
	_, err := d.db.ExecContext(ctx,
		`DELETE FROM alert WHERE id = $1;`,
		AlertID,
	)

	if err != nil {
		d.logger.Ctx(ctx).Error("Unable to delete alert", zap.Error(err))
		return err
	}

	return nil
}
