package poker

import (
	"context"
	"fmt"
)

// PurgeOldGames deletes games older than {daysOld} days
func (d *Service) PurgeOldGames(ctx context.Context, daysOld int) error {
	if _, err := d.DB.ExecContext(ctx,
		`DELETE FROM thunderdome.poker WHERE last_active < (NOW() - $1 * interval '1 day');`,
		daysOld,
	); err != nil {
		return fmt.Errorf("clean poker games query error: %v", err)
	}

	return nil
}
