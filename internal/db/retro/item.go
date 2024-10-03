package retro

import (
	"encoding/json"
	"fmt"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
	"go.uber.org/zap"
)

// CreateRetroItem adds a feedback item to the retro
func (d *Service) CreateRetroItem(retroID string, userID string, itemType string, content string) ([]*thunderdome.RetroItem, error) {
	var groupID string
	err := d.DB.QueryRow(
		`INSERT INTO thunderdome.retro_group
		(retro_id)
		VALUES ($1) RETURNING id;`,
		retroID,
	).Scan(&groupID)
	if err != nil {
		return nil, fmt.Errorf("insert retro group error: %v", err)
	}

	if _, err := d.DB.Exec(
		`INSERT INTO thunderdome.retro_item
		(retro_id, group_id, type, content, user_id)
		VALUES ($1, $2, $3, $4, $5);`,
		retroID, groupID, itemType, content, userID,
	); err != nil {
		d.Logger.Error("insert retro item error", zap.Error(err))
	}

	items := d.GetRetroItems(retroID)

	return items, nil
}

// GroupRetroItem changes the group_id of retro item
func (d *Service) GroupRetroItem(retroID string, itemID string, groupID string) (thunderdome.RetroItem, error) {
	ri := thunderdome.RetroItem{}

	err := d.DB.QueryRow(
		`UPDATE thunderdome.retro_item SET group_id = $3
 				WHERE retro_id = $1 AND id = $2
 				RETURNING id, user_id, group_id, content, type;`,
		retroID, itemID, groupID,
	).Scan(&ri.ID, &ri.UserID, &ri.GroupID, &ri.Content, &ri.Type)

	if err != nil {
		d.Logger.Error("move (group) retro item error", zap.Error(err))
		return ri, err
	}

	return ri, nil
}

// DeleteRetroItem removes item from the current board by ID
func (d *Service) DeleteRetroItem(retroID string, userID string, itemType string, itemID string) ([]*thunderdome.RetroItem, error) {
	if _, err := d.DB.Exec(
		`DELETE FROM thunderdome.retro_item WHERE id = $1 AND type = $2;`, itemID, itemType); err != nil {
		d.Logger.Error("delete retro item error", zap.Error(err))
	}

	items := d.GetRetroItems(retroID)

	return items, nil
}

// GetRetroItems retrieves retro items
func (d *Service) GetRetroItems(retroID string) []*thunderdome.RetroItem {
	var items = make([]*thunderdome.RetroItem, 0)

	itemRows, itemsErr := d.DB.Query(
		`SELECT
				ri.id, ri.user_id, ri.group_id, ri.content, ri.type,
				COALESCE(
					json_agg(rc ORDER BY rc.created_date) FILTER (WHERE rc.id IS NOT NULL), '[]'
				) AS comments
			FROM thunderdome.retro_item ri
			LEFT JOIN thunderdome.retro_item_comment rc ON rc.item_id = ri.id
			WHERE ri.retro_id = $1
			GROUP BY ri.id, ri.created_date
			ORDER BY ri.created_date ASC;`,
		retroID,
	)
	if itemsErr == nil {
		defer itemRows.Close()
		for itemRows.Next() {
			var comments string
			var ri = &thunderdome.RetroItem{
				Comments: make([]*thunderdome.RetroItemComment, 0),
			}
			if err := itemRows.Scan(&ri.ID, &ri.UserID, &ri.GroupID, &ri.Content, &ri.Type, &comments); err != nil {
				d.Logger.Error("get retro items query scan error", zap.Error(err))
			} else {
				jsonErr := json.Unmarshal([]byte(comments), &ri.Comments)
				if jsonErr != nil {
					d.Logger.Error("retro item comments json error", zap.Error(jsonErr))
				}
				items = append(items, ri)
			}
		}
	} else {
		d.Logger.Error("get retro items query error", zap.Error(itemsErr))
	}

	return items
}

// GetRetroGroups retrieves retro groups
func (d *Service) GetRetroGroups(retroID string) []*thunderdome.RetroGroup {
	var groups = make([]*thunderdome.RetroGroup, 0)

	itemRows, itemsErr := d.DB.Query(
		`SELECT id, COALESCE(name, '') FROM thunderdome.retro_group WHERE retro_id = $1 ORDER BY created_date ASC;`,
		retroID,
	)
	if itemsErr == nil {
		defer itemRows.Close()
		for itemRows.Next() {
			var ri = &thunderdome.RetroGroup{}
			if err := itemRows.Scan(&ri.ID, &ri.Name); err != nil {
				d.Logger.Error("get retro groups query scan error", zap.Error(err))
			} else {
				groups = append(groups, ri)
			}
		}
	} else {
		d.Logger.Error("get retro groups query error", zap.Error(itemsErr))
	}

	return groups
}

// GroupNameChange changes retro item group name
func (d *Service) GroupNameChange(retroID string, groupID string, name string) (thunderdome.RetroGroup, error) {
	rg := thunderdome.RetroGroup{}

	err := d.DB.QueryRow(
		`UPDATE thunderdome.retro_group SET name = $3
				WHERE retro_id = $1 AND id = $2
				RETURNING id, name;`,
		retroID, groupID, name,
	).Scan(&rg.ID, &rg.Name)

	if err != nil {
		d.Logger.Error("update retro group name error", zap.Error(err))
		return rg, err
	}

	return rg, nil
}

// ItemCommentAdd adds a comment to a retro item
func (d *Service) ItemCommentAdd(retroID string, itemID string, userID string, comment string) ([]*thunderdome.RetroItem, error) {
	if _, err := d.DB.Exec(
		`INSERT INTO thunderdome.retro_item_comment (item_id, user_id, comment) VALUES ($1, $2, $3);`,
		itemID,
		userID,
		comment,
	); err != nil {
		d.Logger.Error("ItemCommentAdd error", zap.Error(err))
	}

	items := d.GetRetroItems(retroID)

	return items, nil
}

// ItemCommentEdit edits a retro item comment
func (d *Service) ItemCommentEdit(retroID string, commentID string, comment string) ([]*thunderdome.RetroItem, error) {
	if _, err := d.DB.Exec(
		`UPDATE thunderdome.retro_item_comment SET comment = $2 WHERE id = $1;`,
		commentID,
		comment,
	); err != nil {
		d.Logger.Error("ItemCommentEdit error", zap.Error(err))
	}

	items := d.GetRetroItems(retroID)

	return items, nil
}

// ItemCommentDelete deletes a retro item comment
func (d *Service) ItemCommentDelete(retroID string, commentID string) ([]*thunderdome.RetroItem, error) {
	if _, err := d.DB.Exec(
		`DELETE FROM thunderdome.retro_item_comment WHERE id = $1;`,
		commentID,
	); err != nil {
		d.Logger.Error("ItemCommentDelete error", zap.Error(err))
	}

	items := d.GetRetroItems(retroID)

	return items, nil
}
