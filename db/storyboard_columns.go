package db

import (
	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
	"go.uber.org/zap"
)

// CreateStoryboardColumn adds a new column to a Storyboard
func (d *StoryboardService) CreateStoryboardColumn(StoryboardID string, GoalID string, userID string) ([]*thunderdome.StoryboardGoal, error) {
	if _, err := d.DB.Exec(
		`INSERT INTO thunderdome.storyboard_column (storyboard_id, goal_id, sort_order) 
		VALUES ($1, $2, ((SELECT coalesce(MAX(sort_order), 0) FROM thunderdome.storyboard_column WHERE goal_id = $2) + 1));`,
		StoryboardID, GoalID,
	); err != nil {
		d.Logger.Error("CALL thunderdome.create_storyboard_column error", zap.Error(err))
	}

	goals := d.GetStoryboardGoals(StoryboardID)

	return goals, nil
}

// ReviseStoryboardColumn revises a storyboard column
func (d *StoryboardService) ReviseStoryboardColumn(StoryboardID string, UserID string, ColumnID string, ColumnName string) ([]*thunderdome.StoryboardGoal, error) {
	if _, err := d.DB.Exec(
		`UPDATE thunderdome.storyboard_column SET name = $2, updated_date = NOW() WHERE id = $1;`,
		ColumnID,
		ColumnName,
	); err != nil {
		d.Logger.Error("CALL thunderdome.revise_storyboard_column error", zap.Error(err))
	}

	goals := d.GetStoryboardGoals(StoryboardID)

	return goals, nil
}

// DeleteStoryboardColumn removes a column from the current board by ID
func (d *StoryboardService) DeleteStoryboardColumn(StoryboardID string, userID string, ColumnID string) ([]*thunderdome.StoryboardGoal, error) {
	if _, err := d.DB.Exec(
		`CALL thunderdome.sb_column_delete($1);`, ColumnID); err != nil {
		d.Logger.Error("CALL thunderdome.sb_column_delete error", zap.Error(err))
	}

	goals := d.GetStoryboardGoals(StoryboardID)

	return goals, nil
}
