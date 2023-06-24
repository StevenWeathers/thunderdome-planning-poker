package db

import (
	"errors"
	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
	"go.uber.org/zap"
)

// CreateRetroItem adds a feedback item to the retro
func (d *RetroService) CreateRetroItem(RetroID string, UserID string, ItemType string, Content string) ([]*thunderdome.RetroItem, error) {
	var groupId string
	err := d.DB.QueryRow(
		`INSERT INTO retro_group
		(retro_id)
		VALUES ($1) RETURNING id;`,
		RetroID,
	).Scan(&groupId)
	if err != nil {
		d.Logger.Error("insert retro group error", zap.Error(err))
		return nil, err
	}

	if _, err := d.DB.Exec(
		`INSERT INTO retro_item
		(retro_id, group_id, type, content, user_id)
		VALUES ($1, $2, $3, $4, $5);`,
		RetroID, groupId, ItemType, Content, UserID,
	); err != nil {
		d.Logger.Error("insert retro item error", zap.Error(err))
	}

	items := d.GetRetroItems(RetroID)

	return items, nil
}

// GroupRetroItem changes the group_id of retro item
func (d *RetroService) GroupRetroItem(RetroID string, ItemId string, GroupId string) ([]*thunderdome.RetroItem, error) {
	if _, err := d.DB.Exec(
		`UPDATE retro_item SET group_id = $3 WHERE retro_id = $1 AND id = $2;`,
		RetroID, ItemId, GroupId,
	); err != nil {
		d.Logger.Error("update retro item error", zap.Error(err))
	}

	items := d.GetRetroItems(RetroID)

	return items, nil
}

// DeleteRetroItem removes item from the current board by ID
func (d *RetroService) DeleteRetroItem(RetroID string, userID string, Type string, ItemID string) ([]*thunderdome.RetroItem, error) {
	if _, err := d.DB.Exec(
		`DELETE FROM retro_item WHERE id = $1 AND type = $2;`, ItemID, Type); err != nil {
		d.Logger.Error("delete retro item error", zap.Error(err))
	}

	items := d.GetRetroItems(RetroID)

	return items, nil
}

// GetRetroItems retrieves retro items
func (d *RetroService) GetRetroItems(RetroID string) []*thunderdome.RetroItem {
	var items = make([]*thunderdome.RetroItem, 0)

	itemRows, itemsErr := d.DB.Query(
		`SELECT id, user_id, group_id, content, type FROM retro_item WHERE retro_id = $1 ORDER BY created_date ASC;`,
		RetroID,
	)
	if itemsErr == nil {
		defer itemRows.Close()
		for itemRows.Next() {
			var ri = &thunderdome.RetroItem{}
			if err := itemRows.Scan(&ri.ID, &ri.UserID, &ri.GroupID, &ri.Content, &ri.Type); err != nil {
				d.Logger.Error("get retro items query scan error", zap.Error(err))
			} else {
				items = append(items, ri)
			}
		}
	} else {
		d.Logger.Error("get retro items query error", zap.Error(itemsErr))
	}

	return items
}

// GetRetroGroups retrieves retro groups
func (d *RetroService) GetRetroGroups(RetroID string) []*thunderdome.RetroGroup {
	var groups = make([]*thunderdome.RetroGroup, 0)

	itemRows, itemsErr := d.DB.Query(
		`SELECT id, COALESCE(name, '') FROM retro_group WHERE retro_id = $1 ORDER BY created_date ASC;`,
		RetroID,
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
func (d *RetroService) GroupNameChange(RetroID string, GroupId string, Name string) ([]*thunderdome.RetroGroup, error) {
	if _, err := d.DB.Exec(
		`UPDATE retro_group SET name = $3 WHERE retro_id = $1 AND id = $2;`,
		RetroID, GroupId, Name,
	); err != nil {
		d.Logger.Error("update retro group error", zap.Error(err))
	}

	groups := d.GetRetroGroups(RetroID)

	return groups, nil
}

// GetRetroVotes gets retro votes
func (d *RetroService) GetRetroVotes(RetroID string) []*thunderdome.RetroVote {
	var votes = make([]*thunderdome.RetroVote, 0)

	itemRows, itemsErr := d.DB.Query(
		`SELECT group_id, user_id FROM retro_group_vote WHERE retro_id = $1;`,
		RetroID,
	)
	if itemsErr == nil {
		defer itemRows.Close()
		for itemRows.Next() {
			var rv = &thunderdome.RetroVote{}
			if err := itemRows.Scan(&rv.GroupID, &rv.UserID); err != nil {
				d.Logger.Error("get retro votes query scan error", zap.Error(err))
			} else {
				votes = append(votes, rv)
			}
		}
	} else {
		d.Logger.Error("get retro votes query error", zap.Error(itemsErr))
	}

	return votes
}

// GroupUserVote inserts a user vote for the retro item group
func (d *RetroService) GroupUserVote(RetroID string, GroupID string, UserID string) ([]*thunderdome.RetroVote, error) {
	var voteCount int
	var maxVotes int
	err := d.DB.QueryRow(
		`SELECT r.max_votes
				FROM retro r
				WHERE r.id = $1;`,
		RetroID,
	).Scan(&maxVotes)
	if err != nil {
		d.Logger.Error("retro max votes query error", zap.Error(err))
	}

	err = d.DB.QueryRow(
		`SELECT count(rgv.group_id)
				FROM retro_group_vote rgv
				WHERE rgv.retro_id = $1 AND rgv.user_id = $2;`,
		RetroID, UserID,
	).Scan(&voteCount)
	if err != nil {
		d.Logger.Error("retro group vote count query error", zap.Error(err))
	}

	if voteCount == maxVotes {
		return nil, errors.New("VOTE_LIMIT_REACHED")
	}

	if _, err = d.DB.Exec(
		`INSERT INTO retro_group_vote
		(retro_id, group_id, user_id)
		VALUES ($1, $2, $3);`,
		RetroID, GroupID, UserID,
	); err != nil {
		d.Logger.Error("retro group vote query error", zap.Error(err))
	}

	votes := d.GetRetroVotes(RetroID)

	return votes, nil
}

// GroupUserSubtractVote deletes a user vote for the retro item group
func (d *RetroService) GroupUserSubtractVote(RetroID string, GroupID string, UserID string) ([]*thunderdome.RetroVote, error) {
	if _, err := d.DB.Exec(
		`DELETE FROM retro_group_vote
		WHERE retro_id = $1 AND group_id = $2 AND user_id = $3;`,
		RetroID, GroupID, UserID,
	); err != nil {
		d.Logger.Error("retro group subtract vote query error", zap.Error(err))
	}

	votes := d.GetRetroVotes(RetroID)

	return votes, nil
}
