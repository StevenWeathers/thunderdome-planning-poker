package retro

import (
	"encoding/json"
	"errors"
)

// CreateItem creates a retro item
func (b *Service) CreateItem(RetroID string, UserID string, EventValue string) ([]byte, error, bool) {
	var rs struct {
		Type    string `json:"type"`
		Content string `json:"content"`
		Phase   string `json:"phase"`
	}
	json.Unmarshal([]byte(EventValue), &rs)

	items, err := b.db.CreateRetroItem(RetroID, UserID, rs.Type, rs.Content)
	if err != nil {
		return nil, err, false
	}

	updatedItems, _ := json.Marshal(items)
	msg := createSocketEvent("items_updated", string(updatedItems), "")

	return msg, nil, false
}

// GroupItem changes a retro item's group_id
func (b *Service) GroupItem(RetroID string, UserID string, EventValue string) ([]byte, error, bool) {
	var rs struct {
		ItemId  string `json:"itemId"`
		GroupId string `json:"groupId"`
	}
	json.Unmarshal([]byte(EventValue), &rs)

	items, err := b.db.GroupRetroItem(RetroID, rs.ItemId, rs.GroupId)
	if err != nil {
		return nil, err, false
	}

	updatedItems, _ := json.Marshal(items)
	msg := createSocketEvent("items_updated", string(updatedItems), "")

	return msg, nil, false
}

// DeleteItem deletes a retro item
func (b *Service) DeleteItem(RetroID string, UserID string, EventValue string) ([]byte, error, bool) {
	var rs struct {
		ItemID string `json:"id"`
		Phase  string `json:"phase"`
		Type   string `json:"type"`
	}
	json.Unmarshal([]byte(EventValue), &rs)

	items, err := b.db.DeleteRetroItem(RetroID, UserID, rs.Type, rs.ItemID)
	if err != nil {
		return nil, err, false
	}

	updatedItems, _ := json.Marshal(items)
	msg := createSocketEvent("items_updated", string(updatedItems), "")

	return msg, nil, false
}

// GroupNameChange changes a retro group's name
func (b *Service) GroupNameChange(RetroID string, UserID string, EventValue string) ([]byte, error, bool) {
	var rs struct {
		GroupId string `json:"groupId"`
		Name    string `json:"name"`
	}
	json.Unmarshal([]byte(EventValue), &rs)

	groups, err := b.db.GroupNameChange(RetroID, rs.GroupId, rs.Name)
	if err != nil {
		return nil, err, false
	}

	updatedGroups, _ := json.Marshal(groups)
	msg := createSocketEvent("groups_updated", string(updatedGroups), "")

	return msg, nil, false
}

// GroupUserVote handles a users vote for an item group
func (b *Service) GroupUserVote(RetroID string, UserID string, EventValue string) ([]byte, error, bool) {
	var rs struct {
		GroupId string `json:"groupId"`
	}
	json.Unmarshal([]byte(EventValue), &rs)

	vc, vcErr := b.db.RetroUserVoteCount(RetroID, UserID)
	if vcErr != nil {
		return nil, vcErr, false
	}
	if vc == 3 {
		return nil, errors.New("VOTE_LIMIT_REACHED"), false
	}

	votes, err := b.db.GroupUserVote(RetroID, rs.GroupId, UserID)
	if err != nil {
		return nil, err, false
	}

	updatedVotes, _ := json.Marshal(votes)
	msg := createSocketEvent("votes_updated", string(updatedVotes), "")

	return msg, nil, false
}

// GroupUserSubtractVote handles removing a users vote from an item group
func (b *Service) GroupUserSubtractVote(RetroID string, UserID string, EventValue string) ([]byte, error, bool) {
	var rs struct {
		GroupId string `json:"groupId"`
	}
	json.Unmarshal([]byte(EventValue), &rs)

	votes, err := b.db.GroupUserSubtractVote(RetroID, rs.GroupId, UserID)
	if err != nil {
		return nil, err, false
	}

	updatedVotes, _ := json.Marshal(votes)
	msg := createSocketEvent("votes_updated", string(updatedVotes), "")

	return msg, nil, false
}

// CreateAction creates a retro action
func (b *Service) CreateAction(RetroID string, UserID string, EventValue string) ([]byte, error, bool) {
	var rs struct {
		Content string `json:"content"`
	}
	json.Unmarshal([]byte(EventValue), &rs)

	items, err := b.db.CreateRetroAction(RetroID, UserID, rs.Content)
	if err != nil {
		return nil, err, false
	}

	updatedItems, _ := json.Marshal(items)
	msg := createSocketEvent("action_updated", string(updatedItems), "")

	return msg, nil, false
}

// UpdateAction updates a retro action
func (b *Service) UpdateAction(RetroID string, UserID string, EventValue string) ([]byte, error, bool) {
	var rs struct {
		ActionID  string `json:"id"`
		Completed bool   `json:"completed"`
		Content   string `json:"content"`
	}
	json.Unmarshal([]byte(EventValue), &rs)

	items, err := b.db.UpdateRetroAction(RetroID, rs.ActionID, rs.Content, rs.Completed)
	if err != nil {
		return nil, err, false
	}

	updatedItems, _ := json.Marshal(items)
	msg := createSocketEvent("action_updated", string(updatedItems), "")

	return msg, nil, false
}

// DeleteAction deletes a retro action
func (b *Service) DeleteAction(RetroID string, UserID string, EventValue string) ([]byte, error, bool) {
	var rs struct {
		ActionID string `json:"id"`
	}
	json.Unmarshal([]byte(EventValue), &rs)

	items, err := b.db.DeleteRetroAction(RetroID, UserID, rs.ActionID)
	if err != nil {
		return nil, err, false
	}

	updatedItems, _ := json.Marshal(items)
	msg := createSocketEvent("action_updated", string(updatedItems), "")

	return msg, nil, false
}

// AdvancePhase updates a retro phase
func (b *Service) AdvancePhase(RetroID string, UserID string, EventValue string) ([]byte, error, bool) {
	var rs struct {
		Phase string `json:"phase"`
	}
	json.Unmarshal([]byte(EventValue), &rs)

	retro, err := b.db.RetroAdvancePhase(RetroID, rs.Phase)
	if err != nil {
		return nil, err, false
	}

	updatedItems, _ := json.Marshal(retro)
	msg := createSocketEvent("retro_updated", string(updatedItems), "")

	return msg, nil, false
}

// FacilitatorAdd adds a user as facilitator of the retro
func (b *Service) FacilitatorAdd(RetroID string, UserID string, EventValue string) ([]byte, error, bool) {
	var rs struct {
		UserID string `json:"userId"`
	}
	json.Unmarshal([]byte(EventValue), &rs)

	retro, err := b.db.RetroFacilitatorAdd(RetroID, rs.UserID)
	if err != nil {
		return nil, err, false
	}

	updatedRetro, _ := json.Marshal(retro)
	msg := createSocketEvent("retro_updated", string(updatedRetro), "")

	return msg, nil, false
}

// FacilitatorRemove removes a retro facilitator
func (b *Service) FacilitatorRemove(RetroID string, UserID string, EventValue string) ([]byte, error, bool) {
	var rs struct {
		UserID string `json:"userId"`
	}
	json.Unmarshal([]byte(EventValue), &rs)

	retro, err := b.db.RetroFacilitatorRemove(RetroID, rs.UserID)
	if err != nil {
		return nil, err, false
	}

	updatedRetro, _ := json.Marshal(retro)
	msg := createSocketEvent("retro_updated", string(updatedRetro), "")

	return msg, nil, false
}

// EditRetro handles editing the retro settings
func (b *Service) EditRetro(RetroID string, UserID string, EventValue string) ([]byte, error, bool) {
	var rb struct {
		Name                 string `json:"retroName"`
		JoinCode             string `json:"joinCode"`
		MaxVotes             int    `json:"maxVotes"`
		BrainstormVisibility string `json:"brainstormVisibility"`
	}
	json.Unmarshal([]byte(EventValue), &rb)

	err := b.db.EditRetro(
		RetroID,
		rb.Name,
		rb.JoinCode,
		rb.MaxVotes,
		rb.BrainstormVisibility,
	)
	if err != nil {
		return nil, err, false
	}

	updatedRetro, _ := json.Marshal(rb)
	msg := createSocketEvent("retro_edited", string(updatedRetro), "")

	return msg, nil, false
}

// Delete handles deleting the retro
func (b *Service) Delete(RetroID string, UserID string, EventValue string) ([]byte, error, bool) {
	err := b.db.RetroDelete(RetroID)
	if err != nil {
		return nil, err, false
	}
	msg := createSocketEvent("conceded", "", "")

	return msg, nil, false
}

// Abandon handles setting abandoned true so retro doesn't show up in users retro list, then leaves retro
func (b *Service) Abandon(RetroID string, UserID string, EventValue string) ([]byte, error, bool) {
	b.db.RetroAbandon(RetroID, UserID)

	return nil, errors.New("ABANDONED_RETRO"), true
}

// socketEvent is the event structure used for socket messages
type socketEvent struct {
	Type  string `json:"type"`
	Value string `json:"value"`
	User  string `json:"userId"`
}

func createSocketEvent(Type string, Value string, User string) []byte {
	newEvent := &socketEvent{
		Type:  Type,
		Value: Value,
		User:  User,
	}

	event, _ := json.Marshal(newEvent)

	return event
}
