package retro

import (
	"context"
	"fmt"
)

// CleanRetros deletes retros older than {DaysOld} days
func (d *Service) CleanRetros(ctx context.Context, DaysOld int) error {
	if _, err := d.DB.ExecContext(ctx,
		`DELETE FROM thunderdome.retro WHERE updated_date < (NOW() - $1 * interval '1 day');`,
		DaysOld,
	); err != nil {
		return fmt.Errorf("clean retros query error: %v", err)
	}

	return nil
}
