package db

import (
	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
	"go.uber.org/zap"
)

// CreateStoryboardColumn adds a new column to a Storyboard
func (d *StoryboardService) CreateStoryboardColumn(StoryboardID string, GoalID string, userID string) ([]*thunderdome.StoryboardGoal, error) {
	if _, err := d.DB.Exec(
		`call create_storyboard_column($1, $2);`, StoryboardID, GoalID,
	); err != nil {
		d.Logger.Error("call create_storyboard_column error", zap.Error(err))
	}

	goals := d.GetStoryboardGoals(StoryboardID)

	return goals, nil
}

// ReviseStoryboardColumn revises a storyboard column
func (d *StoryboardService) ReviseStoryboardColumn(StoryboardID string, UserID string, ColumnID string, ColumnName string) ([]*thunderdome.StoryboardGoal, error) {
	if _, err := d.DB.Exec(
		`call revise_storyboard_column($1, $2, $3);`,
		StoryboardID,
		ColumnID,
		ColumnName,
	); err != nil {
		d.Logger.Error("call revise_storyboard_column error", zap.Error(err))
	}

	goals := d.GetStoryboardGoals(StoryboardID)

	return goals, nil
}

// DeleteStoryboardColumn removes a column from the current board by ID
func (d *StoryboardService) DeleteStoryboardColumn(StoryboardID string, userID string, ColumnID string) ([]*thunderdome.StoryboardGoal, error) {
	if _, err := d.DB.Exec(
		`call delete_storyboard_column($1);`, ColumnID); err != nil {
		d.Logger.Error("call delete_storyboard_column error", zap.Error(err))
	}

	goals := d.GetStoryboardGoals(StoryboardID)

	return goals, nil
}
