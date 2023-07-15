package db

import (
	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
	"go.uber.org/zap"
)

// CreateStoryboardStory adds a new story to a Storyboard
func (d *StoryboardService) CreateStoryboardStory(StoryboardID string, GoalID string, ColumnID string, userID string) ([]*thunderdome.StoryboardGoal, error) {
	if _, err := d.DB.Exec(
		`INSERT INTO thunderdome.storyboard_story (storyboard_id, goal_id, column_id, sort_order) 
		VALUES ($1, $2, $3, ((SELECT coalesce(MAX(sort_order), 0) FROM thunderdome.storyboard_story WHERE column_id = $3) + 1));`,
		StoryboardID, GoalID, ColumnID,
	); err != nil {
		d.Logger.Error("CALL thunderdome.create_storyboard_story error", zap.Error(err))
	}

	goals := d.GetStoryboardGoals(StoryboardID)

	return goals, nil
}

// ReviseStoryName updates the story name by ID
func (d *StoryboardService) ReviseStoryName(StoryboardID string, userID string, StoryID string, StoryName string) ([]*thunderdome.StoryboardGoal, error) {
	if _, err := d.DB.Exec(
		`UPDATE thunderdome.storyboard_story SET name = $2, updated_date = NOW() WHERE id = $1;`,
		StoryID,
		StoryName,
	); err != nil {
		d.Logger.Error("CALL thunderdome.update_story_name error", zap.Error(err))
	}

	goals := d.GetStoryboardGoals(StoryboardID)

	return goals, nil
}

// ReviseStoryContent updates the story content by ID
func (d *StoryboardService) ReviseStoryContent(StoryboardID string, userID string, StoryID string, StoryContent string) ([]*thunderdome.StoryboardGoal, error) {
	if _, err := d.DB.Exec(
		`UPDATE thunderdome.storyboard_story SET content = $2, updated_date = NOW() WHERE id = $1;`,
		StoryID,
		StoryContent,
	); err != nil {
		d.Logger.Error("CALL thunderdome.update_story_content error", zap.Error(err))
	}

	goals := d.GetStoryboardGoals(StoryboardID)

	return goals, nil
}

// ReviseStoryColor updates the story color by ID
func (d *StoryboardService) ReviseStoryColor(StoryboardID string, userID string, StoryID string, StoryColor string) ([]*thunderdome.StoryboardGoal, error) {
	if _, err := d.DB.Exec(
		`UPDATE thunderdome.storyboard_story SET color = $2, updated_date = NOW() WHERE id = $1;`,
		StoryID,
		StoryColor,
	); err != nil {
		d.Logger.Error("CALL thunderdome.update_story_color error", zap.Error(err))
	}

	goals := d.GetStoryboardGoals(StoryboardID)

	return goals, nil
}

// ReviseStoryPoints updates the story points by ID
func (d *StoryboardService) ReviseStoryPoints(StoryboardID string, userID string, StoryID string, Points int) ([]*thunderdome.StoryboardGoal, error) {
	if _, err := d.DB.Exec(
		`UPDATE thunderdome.storyboard_story SET points = $2, updated_date = NOW() WHERE id = $1;`,
		StoryID,
		Points,
	); err != nil {
		d.Logger.Error("CALL thunderdome.update_story_points error", zap.Error(err))
	}

	goals := d.GetStoryboardGoals(StoryboardID)

	return goals, nil
}

// ReviseStoryClosed updates the story closed status by ID
func (d *StoryboardService) ReviseStoryClosed(StoryboardID string, userID string, StoryID string, Closed bool) ([]*thunderdome.StoryboardGoal, error) {
	if _, err := d.DB.Exec(
		`UPDATE thunderdome.storyboard_story SET closed = $2, updated_date = NOW() WHERE id = $1;`,
		StoryID,
		Closed,
	); err != nil {
		d.Logger.Error("CALL thunderdome.update_story_closed error", zap.Error(err))
	}

	goals := d.GetStoryboardGoals(StoryboardID)

	return goals, nil
}

// ReviseStoryLink updates the story link by ID
func (d *StoryboardService) ReviseStoryLink(StoryboardID string, userID string, StoryID string, Link string) ([]*thunderdome.StoryboardGoal, error) {
	if _, err := d.DB.Exec(
		`UPDATE thunderdome.storyboard_story SET link = $2, updated_date = NOW() WHERE id = $1;`,
		StoryID,
		Link,
	); err != nil {
		d.Logger.Error("CALL thunderdome.sb_story_link_edit error", zap.Error(err))
	}

	goals := d.GetStoryboardGoals(StoryboardID)

	return goals, nil
}

// MoveStoryboardStory moves the story by ID to Goal/Column by ID
func (d *StoryboardService) MoveStoryboardStory(StoryboardID string, userID string, StoryID string, GoalID string, ColumnID string, PlaceBefore string) ([]*thunderdome.StoryboardGoal, error) {
	if _, err := d.DB.Exec(
		`CALL thunderdome.sb_story_move($1, $2, $3, $4);`,
		StoryID,
		GoalID,
		ColumnID,
		PlaceBefore,
	); err != nil {
		d.Logger.Error("CALL thunderdome.sb_story_move error", zap.Error(err))
	}

	goals := d.GetStoryboardGoals(StoryboardID)

	return goals, nil
}

// DeleteStoryboardStory removes a story from the current board by ID
func (d *StoryboardService) DeleteStoryboardStory(StoryboardID string, userID string, StoryID string) ([]*thunderdome.StoryboardGoal, error) {
	if _, err := d.DB.Exec(
		`CALL thunderdome.sb_story_delete($1);`, StoryID); err != nil {
		d.Logger.Error("CALL thunderdome.sb_story_delete error", zap.Error(err))
	}

	goals := d.GetStoryboardGoals(StoryboardID)

	return goals, nil
}

// AddStoryComment adds a comment to a story
func (d *StoryboardService) AddStoryComment(StoryboardID string, UserID string, StoryID string, Comment string) ([]*thunderdome.StoryboardGoal, error) {
	if _, err := d.DB.Exec(
		`INSERT INTO thunderdome.storyboard_story_comment (storyboard_id, story_id, user_id, comment) VALUES ($1, $2, $3, $4);`,
		StoryboardID,
		StoryID,
		UserID,
		Comment,
	); err != nil {
		d.Logger.Error("CALL thunderdome.story_comment_add error", zap.Error(err))
	}

	goals := d.GetStoryboardGoals(StoryboardID)

	return goals, nil
}

// EditStoryComment edits a story comment
func (d *StoryboardService) EditStoryComment(StoryboardID string, CommentID string, Comment string) ([]*thunderdome.StoryboardGoal, error) {
	if _, err := d.DB.Exec(
		`UPDATE thunderdome.storyboard_story_comment SET comment = $2
        WHERE id = $1;`,
		CommentID,
		Comment,
	); err != nil {
		d.Logger.Error("CALL thunderdome.story_comment_edit error", zap.Error(err))
	}

	goals := d.GetStoryboardGoals(StoryboardID)

	return goals, nil
}

// DeleteStoryComment deletes a story comment
func (d *StoryboardService) DeleteStoryComment(StoryboardID string, CommentID string) ([]*thunderdome.StoryboardGoal, error) {
	if _, err := d.DB.Exec(
		`DELETE FROM thunderdome.storyboard_story_comment WHERE id = $1;`,
		CommentID,
	); err != nil {
		d.Logger.Error("CALL thunderdome.story_comment_delete error", zap.Error(err))
	}

	goals := d.GetStoryboardGoals(StoryboardID)

	return goals, nil
}
