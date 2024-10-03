package retro

import (
	"context"
	"fmt"
)

// CleanRetros deletes retros older than {daysOld} days
func (d *Service) CleanRetros(ctx context.Context, daysOld int) error {
	if _, err := d.DB.ExecContext(ctx,
		`DELETE FROM thunderdome.retro WHERE updated_date < (NOW() - $1 * interval '1 day');`,
		daysOld,
	); err != nil {
		return fmt.Errorf("clean retros query error: %v", err)
	}

	return nil
}
