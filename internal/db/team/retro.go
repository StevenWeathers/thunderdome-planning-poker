package team

import (
	"context"
	"fmt"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
	"go.uber.org/zap"
)

// TeamRetroList gets a list of team retros
func (d *Service) TeamRetroList(ctx context.Context, teamID string, limit int, offset int) []*thunderdome.Retro {
	var retros = make([]*thunderdome.Retro, 0)
	rows, err := d.DB.QueryContext(ctx,
		`SELECT r.id, r.name, r.template_id, r.phase
        FROM thunderdome.retro r
        WHERE r.team_id = $1
        ORDER BY r.created_date DESC
		LIMIT $2
		OFFSET $3;`,
		teamID,
		limit,
		offset,
	)

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var tb thunderdome.Retro

			if err := rows.Scan(
				&tb.ID,
				&tb.Name,
				&tb.TemplateID,
				&tb.Phase,
			); err != nil {
				d.Logger.Ctx(ctx).Error("team_retro_list query scan error", zap.Error(err))
			} else {
				retros = append(retros, &tb)
			}
		}
	} else {
		d.Logger.Ctx(ctx).Error("team_retro_list query error", zap.Error(err))
	}

	return retros
}

// TeamAddRetro adds a retro to a team
func (d *Service) TeamAddRetro(ctx context.Context, teamID string, retroID string) error {
	_, err := d.DB.ExecContext(ctx,
		`UPDATE thunderdome.retro SET team_id = $1 WHERE id = $2;`,
		teamID,
		retroID,
	)

	if err != nil {
		return fmt.Errorf("team add retro query error: %v", err)
	}

	return nil
}

// TeamRemoveRetro removes a retro from a team
func (d *Service) TeamRemoveRetro(ctx context.Context, teamID string, retroID string) error {
	_, err := d.DB.ExecContext(ctx,
		`UPDATE thunderdome.retro SET team_id = $1 WHERE id = $2;`,
		teamID,
		retroID,
	)

	if err != nil {
		return fmt.Errorf("team remove retro query error: %v", err)
	}

	return nil
}
