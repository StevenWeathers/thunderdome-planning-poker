package team

import (
	"context"
	"fmt"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
	"go.uber.org/zap"
)

// TeamPokerList gets a list of team poker games
func (d *Service) TeamPokerList(ctx context.Context, TeamID string, Limit int, Offset int) []*thunderdome.Poker {
	var pokers = make([]*thunderdome.Poker, 0)
	rows, err := d.DB.QueryContext(ctx,
		`SELECT p.id, p.name
        FROM thunderdome.poker p
        WHERE p.team_id = $1
        ORDER BY p.created_date DESC
		LIMIT $2
		OFFSET $3;`,
		TeamID,
		Limit,
		Offset,
	)

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var tb thunderdome.Poker

			if err := rows.Scan(
				&tb.Id,
				&tb.Name,
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
func (d *Service) TeamAddPoker(ctx context.Context, TeamID string, PokerID string) error {
	_, err := d.DB.ExecContext(ctx,
		`UPDATE thunderdome.poker SET team_id = $1 WHERE id = $2;`,
		TeamID,
		PokerID,
	)

	if err != nil {
		return fmt.Errorf("team add poker query error: %v", err)
	}

	return nil
}

// TeamRemovePoker removes a poker game from a team
func (d *Service) TeamRemovePoker(ctx context.Context, TeamID string, PokerID string) error {
	_, err := d.DB.ExecContext(ctx,
		`UPDATE thunderdome.poker SET team_id = null WHERE id = $2 AND team_id = $1;`,
		TeamID,
		PokerID,
	)

	if err != nil {
		return fmt.Errorf("team remove poker query error: %v", err)
	}

	return nil
}
