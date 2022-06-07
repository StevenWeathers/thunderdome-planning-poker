package db

import (
	"errors"
	"github.com/StevenWeathers/thunderdome-planning-poker/model"
	"log"
)

// CreateStoryboardColumn adds a new column to a Storyboard
func (d *Database) CreateStoryboardColumn(StoryboardID string, GoalID string, userID string) ([]*model.StoryboardGoal, error) {
	err := d.ConfirmStoryboardOwner(StoryboardID, userID)
	if err != nil {
		return nil, errors.New("Incorrect permissions")
	}

	if _, err := d.db.Exec(
		`call create_storyboard_column($1, $2);`, StoryboardID, GoalID,
	); err != nil {
		log.Println(err)
	}

	goals := d.GetStoryboardGoals(StoryboardID)

	return goals, nil
}

// ReviseStoryboardColumn revises a storyboard column
func (d *Database) ReviseStoryboardColumn(StoryboardID string, UserID string, ColumnID string, ColumnName string) ([]*model.StoryboardGoal, error) {
	err := d.ConfirmStoryboardOwner(StoryboardID, UserID)
	if err != nil {
		return nil, errors.New("Incorrect permissions")
	}

	if _, err := d.db.Exec(
		`call revise_storyboard_column($1, $2, $3);`,
		StoryboardID,
		ColumnID,
		ColumnName,
	); err != nil {
		log.Println(err)
	}

	goals := d.GetStoryboardGoals(StoryboardID)

	return goals, nil
}

// DeleteStoryboardColumn removes a column from the current board by ID
func (d *Database) DeleteStoryboardColumn(StoryboardID string, userID string, ColumnID string) ([]*model.StoryboardGoal, error) {
	err := d.ConfirmStoryboardOwner(StoryboardID, userID)
	if err != nil {
		return nil, errors.New("Incorrect permissions")
	}

	if _, err := d.db.Exec(
		`call delete_storyboard_column($1);`, ColumnID); err != nil {
		log.Println(err)
	}

	goals := d.GetStoryboardGoals(StoryboardID)

	return goals, nil
}
