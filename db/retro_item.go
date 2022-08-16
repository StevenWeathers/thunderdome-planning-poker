package db

import (
	"github.com/StevenWeathers/thunderdome-planning-poker/model"
	"go.uber.org/zap"
)

// CreateRetroItem adds a feedback item to the retro
func (d *Database) CreateRetroItem(RetroID string, UserID string, ItemType string, Content string) ([]*model.RetroItem, error) {
	var groupId string
	err := d.db.QueryRow(
		`INSERT INTO retro_group
		(retro_id)
		VALUES ($1) RETURNING id;`,
		RetroID,
	).Scan(&groupId)
	if err != nil {
		d.logger.Error("insert retro group error", zap.Error(err))
		return nil, err
	}

	if _, err := d.db.Exec(
		`INSERT INTO retro_item
		(retro_id, group_id, type, content, user_id)
		VALUES ($1, $2, $3, $4, $5);`,
		RetroID, groupId, ItemType, Content, UserID,
	); err != nil {
		d.logger.Error("insert retro item error", zap.Error(err))
	}

	items := d.GetRetroItems(RetroID)

	return items, nil
}

// GroupRetroItem changes the group_id of retro item
func (d *Database) GroupRetroItem(RetroID string, ItemId string, GroupId string) ([]*model.RetroItem, error) {
	if _, err := d.db.Exec(
		`UPDATE retro_item SET group_id = $3 WHERE retro_id = $1 AND id = $2;`,
		RetroID, ItemId, GroupId,
	); err != nil {
		d.logger.Error("update retro item error", zap.Error(err))
	}

	items := d.GetRetroItems(RetroID)

	return items, nil
}

// DeleteRetroItem removes item from the current board by ID
func (d *Database) DeleteRetroItem(RetroID string, userID string, Type string, ItemID string) ([]*model.RetroItem, error) {
	if _, err := d.db.Exec(
		`DELETE FROM retro_item WHERE id = $1 AND type = $2;`, ItemID, Type); err != nil {
		d.logger.Error("delete retro item error", zap.Error(err))
	}

	items := d.GetRetroItems(RetroID)

	return items, nil
}

// GetRetroItems retrieves retro items
func (d *Database) GetRetroItems(RetroID string) []*model.RetroItem {
	var items = make([]*model.RetroItem, 0)

	itemRows, itemsErr := d.db.Query(
		`SELECT id, user_id, group_id, content, type FROM retro_item WHERE retro_id = $1 ORDER BY created_date ASC;`,
		RetroID,
	)
	if itemsErr == nil {
		defer itemRows.Close()
		for itemRows.Next() {
			var ri = &model.RetroItem{}
			if err := itemRows.Scan(&ri.ID, &ri.UserID, &ri.GroupID, &ri.Content, &ri.Type); err != nil {
				d.logger.Error("get retro items query scan error", zap.Error(err))
			} else {
				items = append(items, ri)
			}
		}
	} else {
		d.logger.Error("get retro items query error", zap.Error(itemsErr))
	}

	return items
}

// GetRetroGroups retrieves retro groups
func (d *Database) GetRetroGroups(RetroID string) []*model.RetroGroup {
	var groups = make([]*model.RetroGroup, 0)

	itemRows, itemsErr := d.db.Query(
		`SELECT id, COALESCE(name, '') FROM retro_group WHERE retro_id = $1 ORDER BY created_date ASC;`,
		RetroID,
	)
	if itemsErr == nil {
		defer itemRows.Close()
		for itemRows.Next() {
			var ri = &model.RetroGroup{}
			if err := itemRows.Scan(&ri.ID, &ri.Name); err != nil {
				d.logger.Error("get retro groups query scan error", zap.Error(err))
			} else {
				groups = append(groups, ri)
			}
		}
	} else {
		d.logger.Error("get retro groups query error", zap.Error(itemsErr))
	}

	return groups
}

// GroupNameChange changes retro item group name
func (d *Database) GroupNameChange(RetroID string, GroupId string, Name string) ([]*model.RetroGroup, error) {
	if _, err := d.db.Exec(
		`UPDATE retro_group SET name = $3 WHERE retro_id = $1 AND id = $2;`,
		RetroID, GroupId, Name,
	); err != nil {
		d.logger.Error("update retro group error", zap.Error(err))
	}

	groups := d.GetRetroGroups(RetroID)

	return groups, nil
}

// GetRetroVotes gets retro votes
func (d *Database) GetRetroVotes(RetroID string) []*model.RetroVote {
	var votes = make([]*model.RetroVote, 0)

	itemRows, itemsErr := d.db.Query(
		`SELECT group_id, user_id FROM retro_group_vote WHERE retro_id = $1;`,
		RetroID,
	)
	if itemsErr == nil {
		defer itemRows.Close()
		for itemRows.Next() {
			var rv = &model.RetroVote{}
			if err := itemRows.Scan(&rv.GroupID, &rv.UserID); err != nil {
				d.logger.Error("get retro votes query scan error", zap.Error(err))
			} else {
				votes = append(votes, rv)
			}
		}
	} else {
		d.logger.Error("get retro votes query error", zap.Error(itemsErr))
	}

	return votes
}

// GroupUserVote inserts a user vote for the retro item group
func (d *Database) GroupUserVote(RetroID string, GroupID string, UserID string) ([]*model.RetroVote, error) {
	if _, err := d.db.Exec(
		`INSERT INTO retro_group_vote
		(retro_id, group_id, user_id)
		VALUES ($1, $2, $3);`,
		RetroID, GroupID, UserID,
	); err != nil {
		d.logger.Error("retro group vote query error", zap.Error(err))
	}

	votes := d.GetRetroVotes(RetroID)

	return votes, nil
}

// GroupUserSubtractVote deletes a user vote for the retro item group
func (d *Database) GroupUserSubtractVote(RetroID string, GroupID string, UserID string) ([]*model.RetroVote, error) {
	if _, err := d.db.Exec(
		`DELETE FROM retro_group_vote
		WHERE retro_id = $1 AND group_id = $2 AND user_id = $3;`,
		RetroID, GroupID, UserID,
	); err != nil {
		d.logger.Error("retro group subtract vote query error", zap.Error(err))
	}

	votes := d.GetRetroVotes(RetroID)

	return votes, nil
}

// RetroUserVoteCount gets a count of user's votes for the retro
func (d *Database) RetroUserVoteCount(RetroID string, UserID string) (int, error) {
	var voteCount int

	err := d.db.QueryRow(
		`SELECT count(group_id) FROM retro_group_vote WHERE retro_id = $1 AND user_id = $2;`,
		RetroID,
		UserID,
	).Scan(&voteCount)
	if err != nil {
		d.logger.Error("retro group vote count query error", zap.Error(err))
		return voteCount, err
	}

	return voteCount, nil
}
