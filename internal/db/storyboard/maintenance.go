package storyboard

import (
	"context"
	"fmt"
)

// CleanStoryboards deletes storyboards older than {DaysOld} days
func (d *Service) CleanStoryboards(ctx context.Context, DaysOld int) error {
	if _, err := d.DB.ExecContext(ctx,
		`DELETE FROM thunderdome.storyboard WHERE updated_date < (NOW() - $1 * interval '1 day');`,
		DaysOld,
	); err != nil {
		return fmt.Errorf("clean storyboards query error: %v", err)
	}

	return nil
}
