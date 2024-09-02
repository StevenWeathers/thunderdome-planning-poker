package storyboard

import (
	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
	"go.uber.org/zap"
)

// CreateStoryboardColumn adds a new column to a Storyboard
func (d *Service) CreateStoryboardColumn(StoryboardID string, GoalID string, userID string) ([]*thunderdome.StoryboardGoal, error) {
	if _, err := d.DB.Exec(
		`INSERT INTO thunderdome.storyboard_column (storyboard_id, goal_id, sort_order) 
		VALUES ($1, $2, ((SELECT coalesce(MAX(sort_order), 0) FROM thunderdome.storyboard_column WHERE goal_id = $2) + 1));`,
		StoryboardID, GoalID,
	); err != nil {
		d.Logger.Error("CreateStoryboardColumn error", zap.Error(err))
	}

	goals := d.GetStoryboardGoals(StoryboardID)

	return goals, nil
}

// ReviseStoryboardColumn revises a storyboard column
func (d *Service) ReviseStoryboardColumn(StoryboardID string, UserID string, ColumnID string, ColumnName string) ([]*thunderdome.StoryboardGoal, error) {
	if _, err := d.DB.Exec(
		`UPDATE thunderdome.storyboard_column SET name = $2, updated_date = NOW() WHERE id = $1;`,
		ColumnID,
		ColumnName,
	); err != nil {
		d.Logger.Error("ReviseStoryboardColumn error", zap.Error(err))
	}

	goals := d.GetStoryboardGoals(StoryboardID)

	return goals, nil
}

// DeleteStoryboardColumn removes a column from the current board by ID
func (d *Service) DeleteStoryboardColumn(StoryboardID string, userID string, ColumnID string) ([]*thunderdome.StoryboardGoal, error) {
	if _, err := d.DB.Exec(
		`CALL thunderdome.sb_column_delete($1);`, ColumnID); err != nil {
		d.Logger.Error("CALL thunderdome.sb_column_delete error", zap.Error(err))
	}

	goals := d.GetStoryboardGoals(StoryboardID)

	return goals, nil
}

// ColumnPersonaAdd adds a persona column to a Storyboard column
func (d *Service) ColumnPersonaAdd(StoryboardID string, ColumnID string, PersonaID string) ([]*thunderdome.StoryboardGoal, error) {
	if _, err := d.DB.Exec(
		`INSERT INTO thunderdome.storyboard_column_persona (column_id, persona_id, created_date) 
		VALUES ($1, $2, NOW());`,
		ColumnID, PersonaID,
	); err != nil {
		d.Logger.Error("ColumnPersonaAdd error", zap.Error(err))
	}

	goals := d.GetStoryboardGoals(StoryboardID)

	return goals, nil
}

// ColumnPersonaRemove removes a persona column from a Storyboard column
func (d *Service) ColumnPersonaRemove(StoryboardID string, ColumnID string, PersonaID string) ([]*thunderdome.StoryboardGoal, error) {
	if _, err := d.DB.Exec(
		`DELETE FROM thunderdome.storyboard_column_persona WHERE column_id = $1 AND persona_id = $2;`,
		ColumnID, PersonaID,
	); err != nil {
		d.Logger.Error("ColumnPersonaRemove error", zap.Error(err))
	}

	goals := d.GetStoryboardGoals(StoryboardID)

	return goals, nil
}
