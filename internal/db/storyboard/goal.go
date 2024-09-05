package storyboard

import (
	"encoding/json"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"

	"go.uber.org/zap"
)

// CreateStoryboardGoal adds a new goal to a Storyboard
func (d *Service) CreateStoryboardGoal(StoryboardID string, userID string, GoalName string) ([]*thunderdome.StoryboardGoal, error) {
	if _, err := d.DB.Exec(
		`INSERT INTO
        thunderdome.storyboard_goal
        (storyboard_id, sort_order, name)
        VALUES ($1, ((SELECT coalesce(MAX(sort_order), 0) FROM thunderdome.storyboard_goal WHERE storyboard_id = $1) + 1), $2);`,
		StoryboardID, GoalName,
	); err != nil {
		d.Logger.Error("CALL thunderdome.create_storyboard_goal error", zap.Error(err))
	}

	goals := d.GetStoryboardGoals(StoryboardID)

	return goals, nil
}

// ReviseGoalName updates the plan name by ID
func (d *Service) ReviseGoalName(StoryboardID string, userID string, GoalID string, GoalName string) ([]*thunderdome.StoryboardGoal, error) {
	if _, err := d.DB.Exec(
		`UPDATE thunderdome.storyboard_goal SET name = $2, updated_date = NOW() WHERE id = $1;`,
		GoalID,
		GoalName,
	); err != nil {
		d.Logger.Error("CALL thunderdome.update_storyboard_goal error", zap.Error(err))
	}

	goals := d.GetStoryboardGoals(StoryboardID)

	return goals, nil
}

// DeleteStoryboardGoal removes a goal from the current board by ID
func (d *Service) DeleteStoryboardGoal(StoryboardID string, userID string, GoalID string) ([]*thunderdome.StoryboardGoal, error) {
	if _, err := d.DB.Exec(
		`CALL thunderdome.sb_goal_delete($1);`, GoalID); err != nil {
		d.Logger.Error("CALL thunderdome.sb_goal_delete error", zap.Error(err))
	}

	goals := d.GetStoryboardGoals(StoryboardID)

	return goals, nil
}

// GetStoryboardGoals retrieves goals for given storyboard from db
func (d *Service) GetStoryboardGoals(StoryboardID string) []*thunderdome.StoryboardGoal {
	var goals = make([]*thunderdome.StoryboardGoal, 0)

	goalRows, goalsErr := d.DB.Query(
		`SELECT
            sg.id,
            sg.sort_order,
            sg.name,
            COALESCE(json_agg(to_jsonb(t) - 'goal_id' ORDER BY t.sort_order) FILTER (WHERE t.id IS NOT NULL), '[]') AS columns,
            COALESCE(json_agg(to_jsonb(sgp) - 'goal_id') FILTER (WHERE sgp.goal_id IS NOT NULL), '[]') AS personas
        FROM thunderdome.storyboard_goal sg
        LEFT JOIN (
            SELECT
                sc.*,
                COALESCE(
                    json_agg(stss ORDER BY stss.sort_order) FILTER (WHERE stss.id IS NOT NULL), '[]'
                ) AS stories,
                (SELECT COALESCE(
                    json_agg(sp) FILTER (WHERE cp.column_id IS NOT NULL), '[]'
                ) AS personas
                FROM thunderdome.storyboard_column_persona cp
                LEFT JOIN thunderdome.storyboard_persona sp ON sp.id = cp.persona_id
                WHERE cp.column_id = sc.id) AS personas
            FROM thunderdome.storyboard_column sc
            LEFT JOIN (
                SELECT
                    ss.*,
                    COALESCE(
                        json_agg(stcm ORDER BY stcm.created_date) FILTER (WHERE stcm.id IS NOT NULL), '[]'
                    ) AS comments
                FROM thunderdome.storyboard_story ss
                LEFT JOIN thunderdome.storyboard_story_comment stcm ON stcm.story_id = ss.id
                GROUP BY ss.id
            ) stss ON stss.column_id = sc.id
            GROUP BY sc.id
        ) t ON t.goal_id = sg.id
        LEFT JOIN (
            SELECT gp.goal_id, sp.*
            FROM thunderdome.storyboard_goal_persona gp
            LEFT JOIN thunderdome.storyboard_persona sp ON sp.id = gp.persona_id
        ) sgp ON sgp.goal_id = sg.id
        WHERE sg.storyboard_id = $1
        GROUP BY sg.id, sg.sort_order
        ORDER BY sg.sort_order;`,
		StoryboardID,
	)
	if goalsErr == nil {
		defer goalRows.Close()
		for goalRows.Next() {
			var columns string
			var personas string
			var sg = &thunderdome.StoryboardGoal{
				Id:        "",
				Name:      "",
				SortOrder: 0,
				Columns:   make([]*thunderdome.StoryboardColumn, 0),
			}
			if err := goalRows.Scan(&sg.Id, &sg.SortOrder, &sg.Name, &columns, &personas); err != nil {
				d.Logger.Error("get_storyboard_goals query scan error", zap.Error(err))
			} else {
				goalColumns := make([]*thunderdome.StoryboardColumn, 0)
				jsonErr := json.Unmarshal([]byte(columns), &goalColumns)
				if jsonErr != nil {
					d.Logger.Error("storyboard goals json error", zap.Error(jsonErr))
				}
				sg.Columns = goalColumns
				goalPersonas := make([]*thunderdome.StoryboardPersona, 0)
				pJsonErr := json.Unmarshal([]byte(personas), &goalPersonas)
				if jsonErr != nil {
					d.Logger.Error("storyboard goals json error", zap.Error(pJsonErr))
				}
				sg.Personas = goalPersonas
				goals = append(goals, sg)
			}
		}
	}

	return goals
}
