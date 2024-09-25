package storyboard

import (
	"context"
	"errors"
	"fmt"

	"github.com/StevenWeathers/thunderdome-planning-poker/internal/fracindex"
	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
	"go.uber.org/zap"
)

// CreateStoryboardStory adds a new story to a Storyboard
func (d *Service) CreateStoryboardStory(StoryboardID string, GoalID string, ColumnID string, userID string) ([]*thunderdome.StoryboardGoal, error) {
	var betweenAkey *string
	var logger = d.Logger.With(
		zap.String("user_id", userID),
		zap.String("storyboard_id", StoryboardID),
		zap.String("column_id", ColumnID),
		zap.String("goal_id", GoalID),
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
         FROM thunderdome.storyboard_story
         WHERE column_id = $1 AND goal_id = $2 AND storyboard_id = $3),
        'a0'
    ) AS last_display_order;`,
		ColumnID, GoalID, StoryboardID,
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

	if _, err := d.DB.Exec(
		`INSERT INTO thunderdome.storyboard_story (storyboard_id, goal_id, column_id, display_order) 
		VALUES ($1, $2, $3, $4);`,
		StoryboardID, GoalID, ColumnID, displayOrder,
	); err != nil {
		logger.Error(
			"create story error",
			zap.Error(err),
			zap.Stringp("display_order_a", betweenAkey),
		)
		return nil, err
	}

	if commitErr := tx.Commit(); commitErr != nil {
		logger.Error("update drivers: unable to commit", zap.Error(commitErr))
		return nil, fmt.Errorf("failed to update storyboard story display_order: %v", commitErr)
	}

	goals := d.GetStoryboardGoals(StoryboardID)

	return goals, nil
}

// ReviseStoryName updates the story name by ID
func (d *Service) ReviseStoryName(StoryboardID string, userID string, StoryID string, StoryName string) ([]*thunderdome.StoryboardGoal, error) {
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
func (d *Service) ReviseStoryContent(StoryboardID string, userID string, StoryID string, StoryContent string) ([]*thunderdome.StoryboardGoal, error) {
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
func (d *Service) ReviseStoryColor(StoryboardID string, userID string, StoryID string, StoryColor string) ([]*thunderdome.StoryboardGoal, error) {
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
func (d *Service) ReviseStoryPoints(StoryboardID string, userID string, StoryID string, Points int) ([]*thunderdome.StoryboardGoal, error) {
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
func (d *Service) ReviseStoryClosed(StoryboardID string, userID string, StoryID string, Closed bool) ([]*thunderdome.StoryboardGoal, error) {
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
func (d *Service) ReviseStoryLink(StoryboardID string, userID string, StoryID string, Link string) ([]*thunderdome.StoryboardGoal, error) {
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
func (d *Service) MoveStoryboardStory(StoryboardID string, userID string, StoryID string, GoalID string, ColumnID string, PlaceBefore string) ([]*thunderdome.StoryboardGoal, error) {
	var betweenAkey *string
	var betweenBkey *string
	var logger = d.Logger.With(
		zap.String("user_id", userID),
		zap.String("storyboard_id", StoryboardID),
		zap.String("story_id", StoryID),
		zap.String("place_before", PlaceBefore),
		zap.String("column_id", ColumnID),
		zap.String("goal_id", GoalID),
	)

	tx, err := d.DB.BeginTx(context.Background(), nil)
	if err != nil {
		logger.Error("begin transaction error", zap.Error(err))
		return nil, err
	}
	defer tx.Rollback()

	if PlaceBefore == "" {
		if err := tx.QueryRow(
			`
		SELECT 
        (SELECT MAX(display_order)
         FROM thunderdome.storyboard_story
         WHERE column_id = $1 AND goal_id = $2 AND storyboard_id = $3)
          AS last_display_order;`,
			ColumnID, GoalID, StoryboardID,
		).Scan(&betweenAkey); err != nil {
			logger.Error("get display_order between query error",
				zap.Error(err),
			)
			return nil, err
		}
	} else {
		if err := tx.QueryRow(
			`
		WITH current_story AS (
			SELECT id, column_id, display_order
			FROM thunderdome.storyboard_story
			WHERE id = $1 AND column_id = $2 AND goal_id = $3
		),
		preceding_story AS (
			SELECT id, display_order
			FROM thunderdome.storyboard_story
			WHERE column_id = (SELECT column_id FROM current_story)
			AND goal_id = (SELECT goal_id FROM current_story)
			AND display_order < (SELECT display_order FROM current_story)
			ORDER BY display_order DESC
			LIMIT 1
		)
		SELECT 
			cs.display_order AS current_display_order,
			ps.display_order AS preceding_display_order
		FROM current_story cs
		LEFT JOIN preceding_story ps ON true;
		`,
			PlaceBefore, ColumnID, GoalID,
		).Scan(&betweenBkey, &betweenAkey); err != nil {
			logger.Error("get display_order between query error",
				zap.Error(err),
			)
			return nil, err
		}
	}

	displayOrder, err := fracindex.KeyBetween(betweenAkey, betweenBkey)
	if err != nil {
		logger.Error("get display_order between error",
			zap.Error(err),
			zap.Stringp("display_order_a", betweenAkey),
			zap.Stringp("display_order_b", betweenBkey),
		)
		return nil, err
	}

	if displayOrder == nil {
		logger.Error("get display_order returned nil",
			zap.Stringp("display_order_a", betweenAkey),
			zap.Stringp("display_order_b", betweenBkey),
		)
		return nil, errors.New("display order is nil")
	}

	if _, err := tx.Exec(
		`UPDATE thunderdome.storyboard_story 
		SET display_order = $1, column_id = $2, goal_id = $3, updated_date = NOW() WHERE id = $4;`,
		displayOrder, ColumnID, GoalID, StoryID,
	); err != nil {
		logger.Error(
			"update story display_order",
			zap.Error(err),
			zap.Stringp("display_order_a", betweenAkey),
			zap.Stringp("display_order_b", betweenBkey),
			zap.Stringp("display_order", displayOrder),
		)
		return nil, err
	}

	if commitErr := tx.Commit(); commitErr != nil {
		logger.Error("update drivers: unable to commit", zap.Error(commitErr))
		return nil, fmt.Errorf("failed to update storyboard story display_order: %v", commitErr)
	}

	goals := d.GetStoryboardGoals(StoryboardID)

	return goals, nil
}

// DeleteStoryboardStory removes a story from the current board by ID
func (d *Service) DeleteStoryboardStory(StoryboardID string, userID string, StoryID string) ([]*thunderdome.StoryboardGoal, error) {
	if _, err := d.DB.Exec(
		`DELETE FROM thunderdome.storyboard_story WHERE id = $1`, StoryID); err != nil {
		d.Logger.Error("storyboard story delete error", zap.Error(err))
	}

	goals := d.GetStoryboardGoals(StoryboardID)

	return goals, nil
}

// AddStoryComment adds a comment to a story
func (d *Service) AddStoryComment(StoryboardID string, UserID string, StoryID string, Comment string) ([]*thunderdome.StoryboardGoal, error) {
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
func (d *Service) EditStoryComment(StoryboardID string, CommentID string, Comment string) ([]*thunderdome.StoryboardGoal, error) {
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
func (d *Service) DeleteStoryComment(StoryboardID string, CommentID string) ([]*thunderdome.StoryboardGoal, error) {
	if _, err := d.DB.Exec(
		`DELETE FROM thunderdome.storyboard_story_comment WHERE id = $1;`,
		CommentID,
	); err != nil {
		d.Logger.Error("CALL thunderdome.story_comment_delete error", zap.Error(err))
	}

	goals := d.GetStoryboardGoals(StoryboardID)

	return goals, nil
}
