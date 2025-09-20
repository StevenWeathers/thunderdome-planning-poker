package team

import (
	"context"
	"fmt"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
	"go.uber.org/zap"
)

// TeamStoryboardList gets a list of team storyboards
func (d *Service) TeamStoryboardList(ctx context.Context, teamID string, limit int, offset int) ([]*thunderdome.Storyboard, int) {
	var storyboards = make([]*thunderdome.Storyboard, 0)
	var count int

	err := d.DB.QueryRow(
		"SELECT COUNT(*) FROM thunderdome.storyboard WHERE team_id = $1;",
		teamID,
	).Scan(
		&count,
	)
	if err != nil {
		d.Logger.Ctx(ctx).Error("get storyboards count query error", zap.Error(err))
		return storyboards, count
	}

	rows, err := d.DB.QueryContext(ctx,
		`SELECT s.id, s.name
        FROM thunderdome.storyboard s
        WHERE s.team_id = $1
        ORDER BY s.created_date DESC
		LIMIT $2
		OFFSET $3;`,
		teamID,
		limit,
		offset,
	)

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var tb thunderdome.Storyboard

			if err := rows.Scan(
				&tb.ID,
				&tb.Name,
			); err != nil {
				d.Logger.Ctx(ctx).Error("team_storyboard_list query scan error", zap.Error(err))
			} else {
				storyboards = append(storyboards, &tb)
			}
		}
	} else {
		d.Logger.Ctx(ctx).Error("team_storyboard_list query error", zap.Error(err))
	}

	return storyboards, count
}

// TeamAddStoryboard adds a storyboard to a team
func (d *Service) TeamAddStoryboard(ctx context.Context, teamID string, storyboardID string) error {
	_, err := d.DB.ExecContext(ctx,
		`UPDATE thunderdome.storyboard SET team_id = $1 WHERE id = $2;`,
		teamID,
		storyboardID,
	)

	if err != nil {
		return fmt.Errorf("team add storyboard query error: %v", err)
	}

	return nil
}

// TeamRemoveStoryboard removes a storyboard from a team
func (d *Service) TeamRemoveStoryboard(ctx context.Context, teamID string, storyboardID string) error {
	_, err := d.DB.ExecContext(ctx,
		`UPDATE thunderdome.storyboard SET team_id = null WHERE id = $2 AND team_id = $1;`,
		teamID,
		storyboardID,
	)

	if err != nil {
		return fmt.Errorf("team remove storyboard query error: %v", err)
	}

	return nil
}
