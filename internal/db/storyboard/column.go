package storyboard

import (
	"context"
	"errors"
	"fmt"

	"github.com/StevenWeathers/thunderdome-planning-poker/internal/fracindex"
	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
	"go.uber.org/zap"
)

// CreateStoryboardColumn adds a new column to a Storyboard
func (d *Service) CreateStoryboardColumn(storyboardID string, goalID string, userID string) ([]*thunderdome.StoryboardGoal, error) {
	var betweenAkey *string
	var logger = d.Logger.With(
		zap.String("user_id", userID),
		zap.String("storyboard_id", storyboardID),
		zap.String("goal_id", goalID),
	)

	tx, err := d.DB.BeginTx(context.Background(), nil)
	if err != nil {
		logger.Error("begin transaction error", zap.Error(err))
		return nil, err
	}
	defer tx.Rollback()

	if err := tx.QueryRow(
		`
		SELECT
    COALESCE(
        (SELECT MAX(display_order)
         FROM thunderdome.storyboard_column
         WHERE storyboard_id = $1 AND goal_id = $2),
        'a0'
    ) AS last_display_order;`,
		storyboardID, goalID,
	).Scan(&betweenAkey); err != nil {
		logger.Error("get display_order between query error",
			zap.Error(err),
		)
		return nil, err
	}

	displayOrder, err := fracindex.KeyBetween(betweenAkey, nil)
	if err != nil {
		logger.Error("get display_order between error",
			zap.Error(err),
			zap.Stringp("display_order_a", betweenAkey),
		)
		return nil, err
	}

	if displayOrder == nil {
		logger.Error("get display_order returned nil",
			zap.Stringp("display_order_a", betweenAkey),
		)
		return nil, errors.New("display order is nil")
	}

	if _, err := tx.Exec(
		`INSERT INTO thunderdome.storyboard_column (storyboard_id, goal_id, display_order)
		VALUES ($1, $2, $3);`,
		storyboardID, goalID, displayOrder,
	); err != nil {
		logger.Error("CreateStoryboardColumn error",
			zap.Error(err),
			zap.Stringp("display_order", displayOrder),
		)
		return nil, err
	}

	if commitErr := tx.Commit(); commitErr != nil {
		logger.Error("update drivers: unable to commit", zap.Error(commitErr))
		return nil, fmt.Errorf("failed to update storyboard story display_order: %v", commitErr)
	}

	goals := d.GetStoryboardGoals(storyboardID)

	return goals, nil
}

// ReviseStoryboardColumn revises a storyboard column
func (d *Service) ReviseStoryboardColumn(storyboardID string, userID string, columnID string, columnName string) ([]*thunderdome.StoryboardGoal, error) {
	if _, err := d.DB.Exec(
		`UPDATE thunderdome.storyboard_column SET name = $2, updated_date = NOW() WHERE id = $1;`,
		columnID,
		columnName,
	); err != nil {
		d.Logger.Error("revise storyboard column error", zap.Error(err))
	}

	goals := d.GetStoryboardGoals(storyboardID)

	return goals, nil
}

// DeleteStoryboardColumn removes a column from the current board by ID
func (d *Service) DeleteStoryboardColumn(storyboardID string, userID string, columnID string) ([]*thunderdome.StoryboardGoal, error) {
	if _, err := d.DB.Exec(
		`DELETE FROM thunderdome.storyboard_column WHERE id = $1;`, columnID); err != nil {
		d.Logger.Error("delete storyboard column error", zap.Error(err))
	}

	goals := d.GetStoryboardGoals(storyboardID)

	return goals, nil
}

// ColumnPersonaAdd adds a persona column to a Storyboard column
func (d *Service) ColumnPersonaAdd(storyboardID string, columnID string, personaID string) ([]*thunderdome.StoryboardGoal, error) {
	if _, err := d.DB.Exec(
		`INSERT INTO thunderdome.storyboard_column_persona (column_id, persona_id, created_date)
		VALUES ($1, $2, NOW());`,
		columnID, personaID,
	); err != nil {
		d.Logger.Error("ColumnPersonaAdd error", zap.Error(err))
	}

	goals := d.GetStoryboardGoals(storyboardID)

	return goals, nil
}

// ColumnPersonaRemove removes a persona column from a Storyboard column
func (d *Service) ColumnPersonaRemove(storyboardID string, columnID string, personaID string) ([]*thunderdome.StoryboardGoal, error) {
	if _, err := d.DB.Exec(
		`DELETE FROM thunderdome.storyboard_column_persona WHERE column_id = $1 AND persona_id = $2;`,
		columnID, personaID,
	); err != nil {
		d.Logger.Error("ColumnPersonaRemove error", zap.Error(err))
	}

	goals := d.GetStoryboardGoals(storyboardID)

	return goals, nil
}
