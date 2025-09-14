package team

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

// TeamPokerList gets a list of team poker games
func (d *Service) TeamPokerList(ctx context.Context, teamID string, limit int, offset int) []*thunderdome.Poker {
	var pokers = make([]*thunderdome.Poker, 0)
	rows, err := d.DB.QueryContext(ctx,
		`SELECT p.id, p.name, p.voting_locked, COALESCE(p.active_story_id::text, ''), p.point_values_allowed,
		 p.auto_finish_voting, p.point_average_rounding, p.created_date, p.updated_date, p.ended_date,
		 COALESCE(p.team_id::TEXT, ''), p.estimation_scale_id,
		 CASE WHEN COUNT(bl) = 0 THEN '[]'::json ELSE array_to_json(array_agg(bl.user_id)) END AS facilitators
        FROM thunderdome.poker p
        LEFT JOIN thunderdome.poker_facilitator bl ON p.id = bl.poker_id
        WHERE p.team_id = $1
        GROUP BY p.id, p.created_date
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
			var tb = &thunderdome.Poker{
				Users:              make([]*thunderdome.PokerUser, 0),
				Stories:            make([]*thunderdome.Story, 0),
				VotingLocked:       true,
				PointValuesAllowed: make([]string, 0),
				AutoFinishVoting:   true,
				Facilitators:       make([]string, 0),
			}
			var facilitators string
			var vArray pgtype.Array[string]
			m := pgtype.NewMap()

			if err := rows.Scan(
				&tb.ID,
				&tb.Name,
				&tb.VotingLocked,
				&tb.ActiveStoryID,
				m.SQLScanner(&vArray),
				&tb.AutoFinishVoting,
				&tb.PointAverageRounding,
				&tb.CreatedDate,
				&tb.UpdatedDate,
				&tb.EndedDate,
				&tb.TeamID,
				&tb.EstimationScaleID,
				&facilitators,
			); err != nil {
				d.Logger.Ctx(ctx).Error("team_poker list query scan error", zap.Error(err))
			} else {
				// Handle point values array conversion
				tb.PointValuesAllowed = vArray.Elements

				// Handle facilitators JSON conversion
				if facilitators != "" {
					_ = json.Unmarshal([]byte(facilitators), &tb.Facilitators)
				}
				pokers = append(pokers, tb)
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
