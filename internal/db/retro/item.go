package retro

import (
	"database/sql"
	"encoding/json"
	"errors"
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

// DeleteRetroItem removes the user's own item from the current board
func (d *Service) DeleteRetroItem(retroID string, userID string, itemType string, itemID string) ([]*thunderdome.RetroItem, error) {
	result, err := d.DB.Exec(
		`DELETE FROM thunderdome.retro_item WHERE id = $1 AND type = $2 AND retro_id = $3 AND user_id = $4;`,
		itemID, itemType, retroID, userID)
	if err := requireChange(result, err, "ITEM_CHANGE_NOT_ALLOWED"); err != nil {
		return nil, err
	}

	items := d.GetRetroItems(retroID)

	return items, nil
}

// GetRetroItems retrieves retro items
func (d *Service) GetRetroItems(retroID string) []*thunderdome.RetroItem {
	var items = make([]*thunderdome.RetroItem, 0)

	itemRows, itemsErr := d.DB.Query(
		`SELECT
				ri.id, ri.user_id, ri.group_id, ri.content, ri.type, array_to_json(ri.tags) AS tags,
				COALESCE(
					(
						SELECT json_agg(rc ORDER BY rc.created_date)
						FROM thunderdome.retro_item_comment rc
						WHERE rc.item_id = ri.id
					),
					'[]'
				) AS comments,
				COALESCE(
					(
						SELECT json_agg(rr ORDER BY rr.created_date)
						FROM thunderdome.retro_item_reaction rr
						WHERE rr.item_id = ri.id
					),
					'[]'
				) AS reactions
			FROM thunderdome.retro_item ri
			WHERE ri.retro_id = $1
			ORDER BY ri.created_date ASC;`,
		retroID,
	)
	if itemsErr == nil {
		defer itemRows.Close()
		for itemRows.Next() {
			var tags string
			var comments string
			var reactions string
			var ri = &thunderdome.RetroItem{
				Tags:      make([]string, 0),
				Comments:  make([]*thunderdome.RetroItemComment, 0),
				Reactions: make([]*thunderdome.RetroItemReaction, 0),
			}
			if err := itemRows.Scan(&ri.ID, &ri.UserID, &ri.GroupID, &ri.Content, &ri.Type, &tags, &comments, &reactions); err != nil {
				d.Logger.Error("get retro items query scan error", zap.Error(err))
			} else {
				jsonErr := json.Unmarshal([]byte(tags), &ri.Tags)
				if jsonErr != nil {
					d.Logger.Error("retro item tags json error", zap.Error(jsonErr))
				}
				jsonErr = json.Unmarshal([]byte(comments), &ri.Comments)
				if jsonErr != nil {
					d.Logger.Error("retro item comments json error", zap.Error(jsonErr))
				}
				jsonErr = json.Unmarshal([]byte(reactions), &ri.Reactions)
				if jsonErr != nil {
					d.Logger.Error("retro item reactions json error", zap.Error(jsonErr))
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

// ItemTagAdd adds a tag to a retro item, ignoring tags it already has (regardless of case)
func (d *Service) ItemTagAdd(retroID string, itemID string, tag string) ([]*thunderdome.RetroItem, error) {
	// a single statement so that concurrent changes to the same item don't overwrite each other
	if _, err := d.DB.Exec(
		`UPDATE thunderdome.retro_item SET tags = array_append(tags, $3::text), updated_date = NOW()
		WHERE retro_id = $1 AND id = $2 AND cardinality(tags) < $4
			AND NOT EXISTS (SELECT 1 FROM unnest(tags) AS existing WHERE lower(existing) = lower($3::text));`,
		retroID, itemID, tag, thunderdome.RetroItemMaxTags,
	); err != nil {
		return nil, fmt.Errorf("ItemTagAdd query error: %v", err)
	}

	items := d.GetRetroItems(retroID)

	return items, nil
}

// ItemTagRemove removes a tag from a retro item
func (d *Service) ItemTagRemove(retroID string, itemID string, tag string) ([]*thunderdome.RetroItem, error) {
	if _, err := d.DB.Exec(
		`UPDATE thunderdome.retro_item SET tags = array_remove(tags, $3::text), updated_date = NOW()
		WHERE retro_id = $1 AND id = $2;`,
		retroID, itemID, tag,
	); err != nil {
		return nil, fmt.Errorf("ItemTagRemove query error: %v", err)
	}

	items := d.GetRetroItems(retroID)

	return items, nil
}

// ItemCommentAdd adds a comment to a retro item
func (d *Service) ItemCommentAdd(retroID string, itemID string, userID string, comment string) ([]*thunderdome.RetroItem, error) {
	if _, err := d.DB.Exec(
		`INSERT INTO thunderdome.retro_item_comment (item_id, user_id, comment)
		SELECT id, $2::uuid, $3 FROM thunderdome.retro_item WHERE id = $1 AND retro_id = $4;`,
		itemID,
		userID,
		comment,
		retroID,
	); err != nil {
		d.Logger.Error("ItemCommentAdd error", zap.Error(err))
	}

	items := d.GetRetroItems(retroID)

	return items, nil
}

// ItemCommentEdit edits the user's own retro item comment
func (d *Service) ItemCommentEdit(retroID string, userID string, commentID string, comment string) ([]*thunderdome.RetroItem, error) {
	result, err := d.DB.Exec(
		`UPDATE thunderdome.retro_item_comment c SET comment = $4
		FROM thunderdome.retro_item i
		WHERE c.id = $1 AND i.id = c.item_id AND i.retro_id = $2 AND c.user_id = $3;`,
		commentID,
		retroID,
		userID,
		comment,
	)
	if err := requireChange(result, err, "COMMENT_CHANGE_NOT_ALLOWED"); err != nil {
		return nil, err
	}

	items := d.GetRetroItems(retroID)

	return items, nil
}

// ItemCommentDelete deletes the user's own retro item comment
func (d *Service) ItemCommentDelete(retroID string, userID string, commentID string) ([]*thunderdome.RetroItem, error) {
	result, err := d.DB.Exec(
		`DELETE FROM thunderdome.retro_item_comment c
		USING thunderdome.retro_item i
		WHERE c.id = $1 AND i.id = c.item_id AND i.retro_id = $2 AND c.user_id = $3;`,
		commentID,
		retroID,
		userID,
	)
	if err := requireChange(result, err, "COMMENT_CHANGE_NOT_ALLOWED"); err != nil {
		return nil, err
	}

	items := d.GetRetroItems(retroID)

	return items, nil
}

// ItemReactionAdd adds a reaction to a retro item
func (d *Service) ItemReactionAdd(retroID string, itemID string, userID string, reaction string) ([]*thunderdome.RetroItem, error) {
	if _, err := d.DB.Exec(
		`INSERT INTO thunderdome.retro_item_reaction (item_id, user_id, reaction)
		SELECT id, $2::uuid, $3 FROM thunderdome.retro_item WHERE id = $1 AND retro_id = $4;`,
		itemID,
		userID,
		reaction,
		retroID,
	); err != nil {
		d.Logger.Error("ItemReactionAdd error", zap.Error(err))
	}

	items := d.GetRetroItems(retroID)

	return items, nil
}

// requireChange fails when a write errored or its ownership conditions matched nothing
func requireChange(result sql.Result, execErr error, refused string) error {
	if execErr != nil {
		return fmt.Errorf("retro write query error: %v", execErr)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("retro write rows affected error: %v", err)
	}
	if rows == 0 {
		return errors.New(refused)
	}

	return nil
}

// ItemReactionDelete removes the user's own reaction from a retro item
func (d *Service) ItemReactionDelete(retroID string, userID string, reactionID string) ([]*thunderdome.RetroItem, error) {
	result, err := d.DB.Exec(
		`DELETE FROM thunderdome.retro_item_reaction r
		USING thunderdome.retro_item i
		WHERE r.id = $1 AND i.id = r.item_id AND i.retro_id = $2 AND r.user_id = $3;`,
		reactionID,
		retroID,
		userID,
	)
	if err := requireChange(result, err, "REACTION_CHANGE_NOT_ALLOWED"); err != nil {
		return nil, err
	}

	items := d.GetRetroItems(retroID)

	return items, nil
}
