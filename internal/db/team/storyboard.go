package team

import (
	"context"
	"fmt"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
	"go.uber.org/zap"
)

// TeamStoryboardList gets a list of team storyboards
func (d *Service) TeamStoryboardList(ctx context.Context, TeamID string, Limit int, Offset int) []*thunderdome.Storyboard {
	var storyboards = make([]*thunderdome.Storyboard, 0)
	rows, err := d.DB.QueryContext(ctx,
		`SELECT s.id, s.name
        FROM thunderdome.storyboard s
        WHERE s.team_id = $1
        ORDER BY s.created_date DESC
		LIMIT $2
		OFFSET $3;`,
		TeamID,
		Limit,
		Offset,
	)

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var tb thunderdome.Storyboard

			if err := rows.Scan(
				&tb.Id,
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

	return storyboards
}

// TeamAddStoryboard adds a storyboard to a team
func (d *Service) TeamAddStoryboard(ctx context.Context, TeamID string, StoryboardID string) error {
	_, err := d.DB.ExecContext(ctx,
		`UPDATE thunderdome.storyboard SET team_id = $1 WHERE id = $2;`,
		TeamID,
		StoryboardID,
	)

	if err != nil {
		return fmt.Errorf("team add storyboard query error: %v", err)
	}

	return nil
}

// TeamRemoveStoryboard removes a storyboard from a team
func (d *Service) TeamRemoveStoryboard(ctx context.Context, TeamID string, StoryboardID string) error {
	_, err := d.DB.ExecContext(ctx,
		`UPDATE thunderdome.storyboard SET team_id = $1 WHERE id = $2;`,
		TeamID,
		StoryboardID,
	)

	if err != nil {
		return fmt.Errorf("team remove storyboard query error: %v", err)
	}

	return nil
}
