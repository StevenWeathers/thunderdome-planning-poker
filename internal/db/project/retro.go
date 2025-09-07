package project

import (
	"context"
	"fmt"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
)

// AssociateRetro associates a sprint retrospective with a project
func (s *Service) AssociateRetro(ctx context.Context, projectID string, retroID string) error {
	_, err := s.DB.ExecContext(ctx,
		`INSERT INTO thunderdome.project_retro (retro_id, project_id)
		VALUES ($1, $2);`,
		retroID,
		projectID,
	)

	if err != nil {
		return fmt.Errorf("error associating retro with project: %v", err)
	}

	return nil
}

// ListRetros retrieves a list of retrospectives associated with a project
func (s *Service) ListRetros(ctx context.Context, projectID string, limit int, offset int) ([]*thunderdome.Retro, error) {
	var retros []*thunderdome.Retro

	rows, err := s.DB.QueryContext(ctx,
		`SELECT r.id, r.name
        FROM thunderdome.project_retro pr
		JOIN thunderdome.retro r ON r.id = pr.retro_id
        WHERE pr.project_id = $1
        ORDER BY r.created_date DESC
		LIMIT $2
		OFFSET $3;`,
		projectID,
		limit,
		offset,
	)

	if err != nil {
		return nil, fmt.Errorf("error listing retros: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var retro thunderdome.Retro
		if err := rows.Scan(&retro.ID, &retro.Name); err != nil {
			return nil, fmt.Errorf("error scanning retro: %v", err)
		}
		retros = append(retros, &retro)
	}

	return retros, nil
}
