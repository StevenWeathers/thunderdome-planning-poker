package db

import (
	"database/sql"
	"encoding/json"
	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"

	"go.uber.org/zap"
)

// CreateRetroAction adds a new action to the retro
func (d *Database) CreateRetroAction(RetroID string, UserID string, Content string) ([]*thunderdome.RetroAction, error) {
	if _, err := d.DB.Exec(
		`INSERT INTO retro_action (retro_id, content) VALUES ($1, $2);`, RetroID, Content,
	); err != nil {
		d.Logger.Error("insert retro_action error", zap.Error(err))
	}

	actions := d.GetRetroActions(RetroID)

	return actions, nil
}

// UpdateRetroAction updates an actions status
func (d *Database) UpdateRetroAction(RetroID string, ActionID string, Content string, Completed bool) (Actions []*thunderdome.RetroAction, DeleteError error) {
	if _, err := d.DB.Exec(
		`UPDATE retro_action SET completed = $2, content = $3, updated_date = NOW() WHERE id = $1;`, ActionID, Completed, Content); err != nil {
		d.Logger.Error("update retro_action error", zap.Error(err))
	}

	actions := d.GetRetroActions(RetroID)

	return actions, nil
}

// DeleteRetroAction removes a goal from the current board by ID
func (d *Database) DeleteRetroAction(RetroID string, userID string, ActionID string) ([]*thunderdome.RetroAction, error) {
	if _, err := d.DB.Exec(
		`DELETE FROM retro_action WHERE id = $1;`, ActionID); err != nil {
		d.Logger.Error("delete retro_action error", zap.Error(err))
	}

	actions := d.GetRetroActions(RetroID)

	return actions, nil
}

// GetRetroActions retrieves retro actions from the DB
func (d *Database) GetRetroActions(RetroID string) []*thunderdome.RetroAction {
	var actions = make([]*thunderdome.RetroAction, 0)

	actionRows, actionsErr := d.DB.Query(
		`SELECT id, content, completed FROM retro_action WHERE retro_id = $1 ORDER BY created_date ASC;`,
		RetroID,
	)
	if actionsErr == nil {
		defer actionRows.Close()
		for actionRows.Next() {
			var ri = &thunderdome.RetroAction{
				ID:        "",
				Content:   "",
				Completed: false,
			}
			if err := actionRows.Scan(&ri.ID, &ri.Content, &ri.Completed); err != nil {
				d.Logger.Error("get retro actions error", zap.Error(err))
			} else {
				actions = append(actions, ri)
			}
		}
	}

	return actions
}

// GetTeamRetroActions retrieves retro actions for the team
func (d *Database) GetTeamRetroActions(TeamID string, Limit int, Offset int, Completed bool) ([]*thunderdome.RetroAction, int, error) {
	var actions = make([]*thunderdome.RetroAction, 0)

	var Count int

	e := d.DB.QueryRow(
		`SELECT COUNT(ra.*) FROM team_retro tr
				LEFT JOIN retro_action ra ON ra.retro_id = tr.retro_id
				WHERE tr.team_id = $1 AND ra.completed = $2;`,
		TeamID,
		Completed,
	).Scan(
		&Count,
	)
	if e != nil {
		return nil, Count, e
	}

	actionRows, err := d.DB.Query(
		`SELECT ra.id, ra.content, ra.completed, ra.retro_id,
				COALESCE(
					json_agg(rac ORDER BY rac.created_date) FILTER (WHERE rac.id IS NOT NULL), '[]'
				) AS comments
				FROM retro_action ra
				LEFT JOIN retro_action_comment rac ON rac.action_id = ra.id
				WHERE ra.retro_id IN (SELECT retro_id FROM team_retro WHERE team_id = $1) AND ra.completed = $2
				GROUP BY ra.id, ra.created_date
				ORDER BY ra.created_date DESC
				LIMIT $3 OFFSET $4;`,
		TeamID,
		Completed,
		Limit,
		Offset,
	)
	if err == nil && err != sql.ErrNoRows {
		defer actionRows.Close()
		for actionRows.Next() {
			var comments string
			var ri = &thunderdome.RetroAction{}
			if err := actionRows.Scan(&ri.ID, &ri.Content, &ri.Completed, &ri.RetroID, &comments); err != nil {
				d.Logger.Error("get retro actions error", zap.Error(err))
			} else {
				Comments := make([]*thunderdome.RetroActionComment, 0)
				jsonErr := json.Unmarshal([]byte(comments), &Comments)
				if jsonErr != nil {
					d.Logger.Error("retro action comments json error", zap.Error(jsonErr))
				}
				ri.Comments = Comments
				actions = append(actions, ri)
			}
		}
	} else {
		d.Logger.Error("get retro actions error", zap.Error(err))
		return actions, Count, err
	}

	return actions, Count, nil
}

// RetroActionCommentAdd adds a comment to a retro action
func (d *Database) RetroActionCommentAdd(RetroID string, ActionID string, UserID string, Comment string) ([]*thunderdome.RetroAction, error) {
	if _, err := d.DB.Exec(
		`call retro_action_comment_add($1, $2, $3, $4);`,
		RetroID,
		ActionID,
		UserID,
		Comment,
	); err != nil {
		d.Logger.Error("call retro_action_comment_add error", zap.Error(err))
	}

	actions := d.GetRetroActions(RetroID)

	return actions, nil
}

// RetroActionCommentEdit edits a retro action comment
func (d *Database) RetroActionCommentEdit(RetroID string, ActionID string, CommentID string, Comment string) ([]*thunderdome.RetroAction, error) {
	if _, err := d.DB.Exec(
		`call retro_action_comment_edit($1, $2, $3, $4);`,
		RetroID,
		ActionID,
		CommentID,
		Comment,
	); err != nil {
		d.Logger.Error("call retro_action_comment_edit error", zap.Error(err))
	}

	actions := d.GetRetroActions(RetroID)

	return actions, nil
}

// RetroActionCommentDelete deletes a retro action comment
func (d *Database) RetroActionCommentDelete(RetroID string, ActionID string, CommentID string) ([]*thunderdome.RetroAction, error) {
	if _, err := d.DB.Exec(
		`call retro_action_comment_delete($1, $2, $3);`,
		RetroID,
		ActionID,
		CommentID,
	); err != nil {
		d.Logger.Error("call retro_action_comment_delete error", zap.Error(err))
	}

	actions := d.GetRetroActions(RetroID)

	return actions, nil
}

// RetroActionAssigneeAdd adds an assignee to a retro action
func (d *Database) RetroActionAssigneeAdd(RetroID string, ActionID string, UserID string) ([]*thunderdome.RetroAction, error) {
	if _, err := d.DB.Exec(
		`call retro_action_assignee_add($1, $2, $3);`,
		RetroID,
		ActionID,
		UserID,
	); err != nil {
		d.Logger.Error("call retro_action_assignee_add error", zap.Error(err))
	}

	actions := d.GetRetroActions(RetroID)

	return actions, nil
}

// RetroActionAssigneeDelete deletes a retro action assignee
func (d *Database) RetroActionAssigneeDelete(RetroID string, ActionID string, UserID string) ([]*thunderdome.RetroAction, error) {
	if _, err := d.DB.Exec(
		`call retro_action_assignee_delete($1, $2, $3);`,
		RetroID,
		ActionID,
		UserID,
	); err != nil {
		d.Logger.Error("call retro_action_assignee_delete error", zap.Error(err))
	}

	actions := d.GetRetroActions(RetroID)

	return actions, nil
}
