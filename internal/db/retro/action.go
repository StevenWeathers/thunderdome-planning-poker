package retro

import (
	"encoding/json"
	"fmt"

	"github.com/StevenWeathers/thunderdome-planning-poker/internal/db"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"

	"go.uber.org/zap"
)

// CreateRetroAction adds a new action to the retro
func (d *Service) CreateRetroAction(RetroID string, UserID string, Content string) ([]*thunderdome.RetroAction, error) {
	if _, err := d.DB.Exec(
		`INSERT INTO thunderdome.retro_action (retro_id, content) VALUES ($1, $2);`, RetroID, Content,
	); err != nil {
		d.Logger.Error("insert retro_action error", zap.Error(err))
	}

	actions := d.GetRetroActions(RetroID)

	return actions, nil
}

// UpdateRetroAction updates an actions status
func (d *Service) UpdateRetroAction(RetroID string, ActionID string, Content string, Completed bool) (Actions []*thunderdome.RetroAction, DeleteError error) {
	if _, err := d.DB.Exec(
		`UPDATE thunderdome.retro_action SET completed = $2, content = $3, updated_date = NOW() WHERE id = $1;`, ActionID, Completed, Content); err != nil {
		d.Logger.Error("update retro_action error", zap.Error(err))
	}

	actions := d.GetRetroActions(RetroID)

	return actions, nil
}

// DeleteRetroAction removes a goal from the current board by ID
func (d *Service) DeleteRetroAction(RetroID string, userID string, ActionID string) ([]*thunderdome.RetroAction, error) {
	if _, err := d.DB.Exec(
		`DELETE FROM thunderdome.retro_action WHERE id = $1;`, ActionID); err != nil {
		d.Logger.Error("delete retro_action error", zap.Error(err))
	}

	actions := d.GetRetroActions(RetroID)

	return actions, nil
}

// GetRetroActions retrieves retro actions from the DB
func (d *Service) GetRetroActions(RetroID string) []*thunderdome.RetroAction {
	var actions = make([]*thunderdome.RetroAction, 0)

	actionRows, actionsErr := d.DB.Query(
		`SELECT a.id, a.content, a.completed,
 		COALESCE(json_agg(json_build_object('id', u.id, 'name', u.name, 'email', COALESCE(u.email, ''), 'avatar', u.avatar))
 		 FILTER (WHERE u.id IS NOT NULL), '[]') AS assignees
		FROM thunderdome.retro_action a
		LEFT JOIN thunderdome.retro_action_assignee as t ON t.action_id = a.id
		LEFT JOIN thunderdome.users u ON t.user_id = u.id
		WHERE a.retro_id = $1
		GROUP BY a.id
		ORDER BY MAX(a.created_date) ASC;`,
		RetroID,
	)
	if actionsErr == nil {
		defer actionRows.Close()
		for actionRows.Next() {
			var ri = &thunderdome.RetroAction{
				ID:        "",
				Content:   "",
				Completed: false,
				Assignees: make([]*thunderdome.User, 0),
			}
			var assignees string
			if err := actionRows.Scan(&ri.ID, &ri.Content, &ri.Completed, &assignees); err != nil {
				d.Logger.Error("get retro actions error", zap.Error(err))
			} else {
				jsonErr := json.Unmarshal([]byte(assignees), &ri.Assignees)
				if jsonErr != nil {
					d.Logger.Error("retro action assignees json error", zap.Error(jsonErr))
				}
				for i, assignee := range ri.Assignees {
					if assignee.Email != "" {
						ri.Assignees[i].GravatarHash = db.CreateGravatarHash(assignee.Email)
					} else {
						ri.Assignees[i].GravatarHash = db.CreateGravatarHash(assignee.Id)
					}
				}
				actions = append(actions, ri)
			}
		}
	}

	return actions
}

// GetTeamRetroActions retrieves retro actions for the team
func (d *Service) GetTeamRetroActions(TeamID string, Limit int, Offset int, Completed bool) ([]*thunderdome.RetroAction, int, error) {
	var actions = make([]*thunderdome.RetroAction, 0)

	var Count int

	e := d.DB.QueryRow(
		`SELECT COUNT(ra.*) FROM thunderdome.retro tr
				LEFT JOIN thunderdome.retro_action ra ON ra.retro_id = tr.id
				WHERE tr.team_id = $1 AND ra.completed = $2;`,
		TeamID,
		Completed,
	).Scan(
		&Count,
	)
	if e != nil {
		return nil, Count, fmt.Errorf("get team retro actions count query error: %v", e)
	}

	actionRows, err := d.DB.Query(
		`SELECT ra.id, ra.content, ra.completed, ra.retro_id,
				(SELECT COALESCE(
					json_agg(rac ORDER BY rac.created_date) FILTER (WHERE rac.id IS NOT NULL), '[]'
				) AS comments
				FROM thunderdome.retro_action_comment rac 
				WHERE rac.action_id = ra.id) AS comments,
				COALESCE(json_agg(json_build_object('id', u.id, 'name', u.name, 'email', COALESCE(u.email, ''), 'avatar', u.avatar))
 		 			FILTER (WHERE u.id IS NOT NULL), '[]') AS assignees
				FROM thunderdome.retro_action ra
				LEFT JOIN thunderdome.retro_action_assignee as t ON t.action_id = ra.id
				LEFT JOIN thunderdome.users u ON t.user_id = u.id
				WHERE ra.retro_id IN (SELECT id FROM thunderdome.retro WHERE team_id = $1) AND ra.completed = $2
				GROUP BY ra.id, ra.created_date
				ORDER BY ra.created_date DESC
				LIMIT $3 OFFSET $4;`,
		TeamID,
		Completed,
		Limit,
		Offset,
	)
	if err == nil {
		defer actionRows.Close()
		for actionRows.Next() {
			var ri = &thunderdome.RetroAction{
				Comments:  make([]*thunderdome.RetroActionComment, 0),
				Assignees: make([]*thunderdome.User, 0),
			}
			var comments string
			var assignees string
			if err := actionRows.Scan(&ri.ID, &ri.Content, &ri.Completed, &ri.RetroID, &comments, &assignees); err != nil {
				d.Logger.Error("get retro actions error", zap.Error(err))
			} else {
				jsonErr := json.Unmarshal([]byte(comments), &ri.Comments)
				if jsonErr != nil {
					d.Logger.Error("retro action comments json error", zap.Error(jsonErr))
				}
				jsonErr = json.Unmarshal([]byte(assignees), &ri.Assignees)
				if jsonErr != nil {
					d.Logger.Error("retro action assignees json error", zap.Error(jsonErr))
				}
				for i, assignee := range ri.Assignees {
					if assignee.Email != "" {
						ri.Assignees[i].GravatarHash = db.CreateGravatarHash(assignee.Email)
					} else {
						ri.Assignees[i].GravatarHash = db.CreateGravatarHash(assignee.Id)
					}
				}
				actions = append(actions, ri)
			}
		}
	} else {
		return actions, Count, fmt.Errorf("get retro actions error: %v", err)
	}

	return actions, Count, nil
}

// RetroActionCommentAdd adds a comment to a retro action
func (d *Service) RetroActionCommentAdd(RetroID string, ActionID string, UserID string, Comment string) ([]*thunderdome.RetroAction, error) {
	if _, err := d.DB.Exec(
		`INSERT INTO thunderdome.retro_action_comment (action_id, user_id, comment) VALUES ($1, $2, $3);`,
		ActionID,
		UserID,
		Comment,
	); err != nil {
		d.Logger.Error("RetroActionCommentAdd error", zap.Error(err))
	}

	actions := d.GetRetroActions(RetroID)

	return actions, nil
}

// RetroActionCommentEdit edits a retro action comment
func (d *Service) RetroActionCommentEdit(RetroID string, ActionID string, CommentID string, Comment string) ([]*thunderdome.RetroAction, error) {
	if _, err := d.DB.Exec(
		`UPDATE thunderdome.retro_action_comment SET comment = $2 WHERE id = $1;`,
		CommentID,
		Comment,
	); err != nil {
		d.Logger.Error("RetroActionCommentEdit error", zap.Error(err))
	}

	actions := d.GetRetroActions(RetroID)

	return actions, nil
}

// RetroActionCommentDelete deletes a retro action comment
func (d *Service) RetroActionCommentDelete(RetroID string, ActionID string, CommentID string) ([]*thunderdome.RetroAction, error) {
	if _, err := d.DB.Exec(
		`DELETE FROM thunderdome.retro_action_comment WHERE id = $1;`,
		CommentID,
	); err != nil {
		d.Logger.Error("RetroActionCommentDelete error", zap.Error(err))
	}

	actions := d.GetRetroActions(RetroID)

	return actions, nil
}

// RetroActionAssigneeAdd adds an assignee to a retro action
func (d *Service) RetroActionAssigneeAdd(RetroID string, ActionID string, UserID string) ([]*thunderdome.RetroAction, error) {
	if _, err := d.DB.Exec(
		`INSERT INTO thunderdome.retro_action_assignee (action_id, user_id) VALUES ($1, $2);`,
		ActionID,
		UserID,
	); err != nil {
		d.Logger.Error("RetroActionAssigneeAdd error", zap.Error(err))
	}

	actions := d.GetRetroActions(RetroID)

	return actions, nil
}

// RetroActionAssigneeDelete deletes a retro action assignee
func (d *Service) RetroActionAssigneeDelete(RetroID string, ActionID string, UserID string) ([]*thunderdome.RetroAction, error) {
	if _, err := d.DB.Exec(
		`DELETE FROM thunderdome.retro_action_assignee WHERE action_id = $1 AND user_id = $2;`,
		ActionID,
		UserID,
	); err != nil {
		d.Logger.Error("RetroActionAssigneeDelete error", zap.Error(err))
	}

	actions := d.GetRetroActions(RetroID)

	return actions, nil
}
