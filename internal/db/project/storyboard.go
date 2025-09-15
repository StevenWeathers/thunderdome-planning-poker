package project

import (
	"context"
	"fmt"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
	"go.uber.org/zap"
)

// AssociateStoryboard associates a Storyboard with a project
func (s *Service) AssociateStoryboard(ctx context.Context, projectID string, storyboardID string) error {
	_, err := s.DB.ExecContext(ctx,
		`INSERT INTO thunderdome.project_storyboard (storyboard_id, project_id)
		VALUES ($1, $2);`,
		storyboardID,
		projectID,
	)

	if err != nil {
		return fmt.Errorf("error associating storyboard with project: %v", err)
	}

	return nil
}

// ListStoryboards retrieves a list of Storyboards associated with a project
func (d *Service) ListStoryboards(ctx context.Context, projectId string, limit int, offset int) ([]*thunderdome.Storyboard, error) {
	var storyboards = make([]*thunderdome.Storyboard, 0)
	rows, err := d.DB.QueryContext(ctx,
		`SELECT s.id, s.name
        FROM thunderdome.project_storyboard ps
		JOIN thunderdome.storyboard s ON s.id = ps.storyboard_id
        WHERE ps.project_id = $1
        ORDER BY s.created_date DESC
		LIMIT $2
		OFFSET $3;`,
		projectId,
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
				return nil, fmt.Errorf("error scanning storyboard: %v", err)
			} else {
				storyboards = append(storyboards, &tb)
			}
		}
	} else {
		d.Logger.Ctx(ctx).Error("team_storyboard_list query error", zap.Error(err))
		return nil, fmt.Errorf("error listing storyboards: %v", err)
	}

	return storyboards, nil
}

// RemoveStoryboard removes the association of a Storyboard with a project
func (s *Service) RemoveStoryboard(ctx context.Context, projectID string, storyboardID string) error {
	res, err := s.DB.ExecContext(ctx,
		`DELETE FROM thunderdome.project_storyboard
		WHERE project_id = $1
		AND storyboard_id = $2;`,
		projectID,
		storyboardID,
	)

	if err != nil {
		return fmt.Errorf("error removing storyboard from project: %v", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("error removing storyboard from project: no rows affected")
	}

	return nil
}
