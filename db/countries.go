package db

import (
	"context"
	"database/sql"
	"errors"
	"go.uber.org/zap"
)

// GetActiveCountries gets a list of user countries
func (d *Database) GetActiveCountries(ctx context.Context) ([]string, error) {
	var countries = make([]string, 0)

	rows, err := d.DB.QueryContext(ctx, `SELECT * FROM countries_active();`)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var country sql.NullString
			if err := rows.Scan(
				&country,
			); err != nil {
				d.logger.Ctx(ctx).Error("countries_active query scan error", zap.Error(err))
			} else {
				if country.String != "" {
					countries = append(countries, country.String)
				}
			}
		}
	} else {
		d.logger.Ctx(ctx).Error("countries_active query error", zap.Error(err))
		return nil, errors.New("error attempting to get active countries")
	}

	return countries, nil
}
