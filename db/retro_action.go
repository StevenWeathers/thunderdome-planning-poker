package db

import (
	"github.com/StevenWeathers/thunderdome-planning-poker/model"
	"go.uber.org/zap"
)

// CreateRetroAction adds a new action to the retro
func (d *Database) CreateRetroAction(RetroID string, UserID string, Content string) ([]*model.RetroAction, error) {
	if _, err := d.db.Exec(
		`INSERT INTO retro_action (retro_id, content) VALUES ($1, $2);`, RetroID, Content,
	); err != nil {
		d.logger.Error("insert retro_action error", zap.Error(err))
	}

	actions := d.GetRetroActions(RetroID)

	return actions, nil
}

// UpdateRetroAction updates an actions status
func (d *Database) UpdateRetroAction(RetroID string, ActionID string, Content string, Completed bool) (Actions []*model.RetroAction, DeleteError error) {
	if _, err := d.db.Exec(
		`UPDATE retro_action SET completed = $2, content = $3, updated_date = NOW() WHERE id = $1;`, ActionID, Completed, Content); err != nil {
		d.logger.Error("update retro_action error", zap.Error(err))
	}

	actions := d.GetRetroActions(RetroID)

	return actions, nil
}

// DeleteRetroAction removes a goal from the current board by ID
func (d *Database) DeleteRetroAction(RetroID string, userID string, ActionID string) ([]*model.RetroAction, error) {
	if _, err := d.db.Exec(
		`DELETE FROM retro_action WHERE id = $1;`, ActionID); err != nil {
		d.logger.Error("delete retro_action error", zap.Error(err))
	}

	actions := d.GetRetroActions(RetroID)

	return actions, nil
}

// GetRetroActions retrieves retro actions from the DB
func (d *Database) GetRetroActions(RetroID string) []*model.RetroAction {
	var actions = make([]*model.RetroAction, 0)

	actionRows, actionsErr := d.db.Query(
		`SELECT id, content, completed FROM retro_action WHERE retro_id = $1 ORDER BY created_date ASC;`,
		RetroID,
	)
	if actionsErr == nil {
		defer actionRows.Close()
		for actionRows.Next() {
			var ri = &model.RetroAction{
				ID:        "",
				Content:   "",
				Completed: false,
			}
			if err := actionRows.Scan(&ri.ID, &ri.Content, &ri.Completed); err != nil {
				d.logger.Error("get retro actions error", zap.Error(err))
			} else {
				actions = append(actions, ri)
			}
		}
	}

	return actions
}
