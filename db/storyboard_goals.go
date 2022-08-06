package db

import (
	"encoding/json"

	"github.com/StevenWeathers/thunderdome-planning-poker/model"
	"go.uber.org/zap"
)

// CreateStoryboardGoal adds a new goal to a Storyboard
func (d *Database) CreateStoryboardGoal(StoryboardID string, userID string, GoalName string) ([]*model.StoryboardGoal, error) {
	if _, err := d.db.Exec(
		`call create_storyboard_goal($1, $2);`, StoryboardID, GoalName,
	); err != nil {
		d.logger.Error("call create_storyboard_goal error", zap.Error(err))
	}

	goals := d.GetStoryboardGoals(StoryboardID)

	return goals, nil
}

// ReviseGoalName updates the plan name by ID
func (d *Database) ReviseGoalName(StoryboardID string, userID string, GoalID string, GoalName string) ([]*model.StoryboardGoal, error) {
	if _, err := d.db.Exec(
		`call update_storyboard_goal($1, $2);`,
		GoalID,
		GoalName,
	); err != nil {
		d.logger.Error("call update_storyboard_goal error", zap.Error(err))
	}

	goals := d.GetStoryboardGoals(StoryboardID)

	return goals, nil
}

// DeleteStoryboardGoal removes a goal from the current board by ID
func (d *Database) DeleteStoryboardGoal(StoryboardID string, userID string, GoalID string) ([]*model.StoryboardGoal, error) {
	if _, err := d.db.Exec(
		`call delete_storyboard_goal($1);`, GoalID); err != nil {
		d.logger.Error("call delete_storyboard_goal error", zap.Error(err))
	}

	goals := d.GetStoryboardGoals(StoryboardID)

	return goals, nil
}

// GetStoryboardGoals retrieves goals for given storyboard from db
func (d *Database) GetStoryboardGoals(StoryboardID string) []*model.StoryboardGoal {
	var goals = make([]*model.StoryboardGoal, 0)

	goalRows, goalsErr := d.db.Query(
		`SELECT * FROM get_storyboard_goals($1);`,
		StoryboardID,
	)
	if goalsErr == nil {
		defer goalRows.Close()
		for goalRows.Next() {
			var columns string
			var personas string
			var sg = &model.StoryboardGoal{
				Id:        "",
				Name:      "",
				SortOrder: 0,
				Columns:   make([]*model.StoryboardColumn, 0),
			}
			if err := goalRows.Scan(&sg.Id, &sg.SortOrder, &sg.Name, &columns, &personas); err != nil {
				d.logger.Error("get_storyboard_goals query scan error", zap.Error(err))
			} else {
				goalColumns := make([]*model.StoryboardColumn, 0)
				jsonErr := json.Unmarshal([]byte(columns), &goalColumns)
				if jsonErr != nil {
					d.logger.Error("storyboard goals json error", zap.Error(jsonErr))
				}
				sg.Columns = goalColumns
				goalPersonas := make([]*model.StoryboardPersona, 0)
				pJsonErr := json.Unmarshal([]byte(personas), &goalPersonas)
				if jsonErr != nil {
					d.logger.Error("storyboard goals json error", zap.Error(pJsonErr))
				}
				sg.Personas = goalPersonas
				goals = append(goals, sg)
			}
		}
	}

	return goals
}
