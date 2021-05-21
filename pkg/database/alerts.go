package database

import (
	"errors"
	"log"
)

// GetActiveAlerts gets alerts from db for UI display
func (d *Database) GetActiveAlerts() []interface{} {
	Alerts := make([]interface{}, 0)

	rows, err := d.db.Query(
		`SELECT id, name, type, content, active, allow_dismiss, registered_only FROM alert WHERE active IS TRUE;`,
	)

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var a Alert

			if err := rows.Scan(
				&a.AlertID,
				&a.Name,
				&a.Type,
				&a.Content,
				&a.Active,
				&a.AllowDismiss,
				&a.RegisteredOnly,
			); err != nil {
				log.Println(err)
			} else {
				Alerts = append(Alerts, &a)
			}
		}
	}

	return Alerts
}

// AlertsList gets alerts from db for admin listing
func (d *Database) AlertsList(Limit int, Offset int) []interface{} {
	Alerts := make([]interface{}, 0)

	rows, err := d.db.Query(
		`SELECT id, name, type, content, active, allow_dismiss, registered_only, created_date, updated_date
		FROM alert WHERE active IS TRUE
		LIMIT $1
		OFFSET $2;
		`,
		Limit,
		Offset,
	)

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var a Alert

			if err := rows.Scan(
				&a.AlertID,
				&a.Name,
				&a.Type,
				&a.Content,
				&a.Active,
				&a.AllowDismiss,
				&a.RegisteredOnly,
				&a.CreatedDate,
				&a.UpdatedDate,
			); err != nil {
				log.Println(err)
			} else {
				Alerts = append(Alerts, &a)
			}
		}
	}

	return Alerts
}

// AlertsCreate creates
func (d *Database) AlertsCreate(Name string, Type string, Content string, Active bool, AllowDismiss bool, RegisteredOnly bool) error {
	if _, err := d.db.Exec(
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
		log.Println(err)
		return errors.New("error attempting to add new alert")
	}

	return nil
}

// AlertsUpdate updates an alert
func (d *Database) AlertsUpdate(ID string, Name string, Type string, Content string, Active bool, AllowDismiss bool, RegisteredOnly bool) error {
	if _, err := d.db.Exec(
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
		log.Println(err)
		return errors.New("error attempting to update alert")
	}

	return nil
}

// AlertDelete deletes an alert
func (d *Database) AlertDelete(AlertID string) error {
	_, err := d.db.Exec(
		`DELETE FROM alert WHERE id = $1;`,
		AlertID,
	)

	if err != nil {
		log.Println("Unable to delete alert: ", err)
		return err
	}

	return nil
}
