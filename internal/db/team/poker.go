package team

import (
	"context"
	"fmt"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
	"go.uber.org/zap"
)

// TeamPokerList gets a list of team poker games
func (d *Service) TeamPokerList(ctx context.Context, teamID string, limit int, offset int) []*thunderdome.Poker {
	var pokers = make([]*thunderdome.Poker, 0)
	rows, err := d.DB.QueryContext(ctx,
		`SELECT p.id, p.name, p.end_time, p.end_reason
        FROM thunderdome.poker p
        WHERE p.team_id = $1
        ORDER BY p.created_date DESC
		LIMIT $2
		OFFSET $3;`,
		teamID,
		limit,
		offset,
	)

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var tb thunderdome.Poker

			if err := rows.Scan(
				&tb.ID,
				&tb.Name,
				&tb.EndTime,
				&tb.EndReason,
			); err != nil {
				d.Logger.Ctx(ctx).Error("team_poker list query scan error", zap.Error(err))
			} else {
				pokers = append(pokers, &tb)
			}
		}
	} else {
		d.Logger.Ctx(ctx).Error("team_poker list query error", zap.Error(err))
	}

	return pokers
}

// TeamAddPoker adds a poker game to a team
func (d *Service) TeamAddPoker(ctx context.Context, teamID string, pokerID string) error {
	_, err := d.DB.ExecContext(ctx,
		`UPDATE thunderdome.poker SET team_id = $1 WHERE id = $2;`,
		teamID,
		pokerID,
	)

	if err != nil {
		return fmt.Errorf("team add poker query error: %v", err)
	}

	return nil
}

// TeamRemovePoker removes a poker game from a team
func (d *Service) TeamRemovePoker(ctx context.Context, teamID string, pokerID string) error {
	_, err := d.DB.ExecContext(ctx,
		`UPDATE thunderdome.poker SET team_id = null WHERE id = $2 AND team_id = $1;`,
		teamID,
		pokerID,
	)

	if err != nil {
		return fmt.Errorf("team remove poker query error: %v", err)
	}

	return nil
}
