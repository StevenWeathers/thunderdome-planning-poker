package project

import (
	"context"
	"fmt"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
)

// AssociatePoker associates a planning poker game with a project
func (s *Service) AssociatePoker(ctx context.Context, projectID string, pokerID string) error {
	_, err := s.DB.ExecContext(ctx,
		`INSERT INTO thunderdome.project_poker (poker_id, project_id)
		VALUES ($1, $2);`,
		pokerID,
		projectID,
	)

	if err != nil {
		return fmt.Errorf("error associating poker with project: %v", err)
	}

	return nil
}

// ListPokerGames retrieves a list of planning poker games associated with a project
func (s *Service) ListPokerGames(ctx context.Context, projectID string, limit int, offset int) ([]*thunderdome.Poker, error) {
	var games []*thunderdome.Poker

	rows, err := s.DB.QueryContext(ctx,
		`SELECT p.id, p.name
        FROM thunderdome.project_poker pp
		JOIN thunderdome.poker p ON p.id = pp.poker_id
        WHERE pp.project_id = $1
        ORDER BY p.created_date DESC
		LIMIT $2
		OFFSET $3;`,
		projectID,
		limit,
		offset,
	)

	if err != nil {
		return nil, fmt.Errorf("error listing poker games: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var game thunderdome.Poker
		if err := rows.Scan(&game.ID, &game.Name); err != nil {
			return nil, fmt.Errorf("error scanning poker game: %v", err)
		}
		games = append(games, &game)
	}

	return games, nil
}
