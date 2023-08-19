package retro

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
)

// CreateItem creates a retro item
func (b *Service) CreateItem(ctx context.Context, RetroID string, UserID string, EventValue string) ([]byte, error, bool) {
	var rs struct {
		Type    string `json:"type"`
		Content string `json:"content"`
		Phase   string `json:"phase"`
	}
	err := json.Unmarshal([]byte(EventValue), &rs)
	if err != nil {
		return nil, err, false
	}

	items, err := b.RetroService.CreateRetroItem(RetroID, UserID, rs.Type, rs.Content)
	if err != nil {
		return nil, err, false
	}

	updatedItems, _ := json.Marshal(items)
	msg := createSocketEvent("items_updated", string(updatedItems), "")

	return msg, nil, false
}

// GroupItem changes a retro item's group_id
func (b *Service) GroupItem(ctx context.Context, RetroID string, UserID string, EventValue string) ([]byte, error, bool) {
	var rs struct {
		ItemId  string `json:"itemId"`
		GroupId string `json:"groupId"`
	}
	err := json.Unmarshal([]byte(EventValue), &rs)
	if err != nil {
		return nil, err, false
	}

	items, err := b.RetroService.GroupRetroItem(RetroID, rs.ItemId, rs.GroupId)
	if err != nil {
		return nil, err, false
	}

	updatedItems, _ := json.Marshal(items)
	msg := createSocketEvent("items_updated", string(updatedItems), "")

	return msg, nil, false
}

// DeleteItem deletes a retro item
func (b *Service) DeleteItem(ctx context.Context, RetroID string, UserID string, EventValue string) ([]byte, error, bool) {
	var rs struct {
		ItemID string `json:"id"`
		Phase  string `json:"phase"`
		Type   string `json:"type"`
	}
	err := json.Unmarshal([]byte(EventValue), &rs)
	if err != nil {
		return nil, err, false
	}

	items, err := b.RetroService.DeleteRetroItem(RetroID, UserID, rs.Type, rs.ItemID)
	if err != nil {
		return nil, err, false
	}

	updatedItems, _ := json.Marshal(items)
	msg := createSocketEvent("items_updated", string(updatedItems), "")

	return msg, nil, false
}

// GroupNameChange changes a retro group's name
func (b *Service) GroupNameChange(ctx context.Context, RetroID string, UserID string, EventValue string) ([]byte, error, bool) {
	var rs struct {
		GroupId string `json:"groupId"`
		Name    string `json:"name"`
	}
	err := json.Unmarshal([]byte(EventValue), &rs)
	if err != nil {
		return nil, err, false
	}

	groups, err := b.RetroService.GroupNameChange(RetroID, rs.GroupId, rs.Name)
	if err != nil {
		return nil, err, false
	}

	updatedGroups, _ := json.Marshal(groups)
	msg := createSocketEvent("groups_updated", string(updatedGroups), "")

	return msg, nil, false
}

// GroupUserVote handles a users vote for an item group
func (b *Service) GroupUserVote(ctx context.Context, RetroID string, UserID string, EventValue string) ([]byte, error, bool) {
	var rs struct {
		GroupId string `json:"groupId"`
	}
	err := json.Unmarshal([]byte(EventValue), &rs)
	if err != nil {
		return nil, err, false
	}

	votes, err := b.RetroService.GroupUserVote(RetroID, rs.GroupId, UserID)
	if err != nil {
		return nil, err, false
	}

	updatedVotes, _ := json.Marshal(votes)
	msg := createSocketEvent("votes_updated", string(updatedVotes), "")

	return msg, nil, false
}

// GroupUserSubtractVote handles removing a users vote from an item group
func (b *Service) GroupUserSubtractVote(ctx context.Context, RetroID string, UserID string, EventValue string) ([]byte, error, bool) {
	var rs struct {
		GroupId string `json:"groupId"`
	}
	err := json.Unmarshal([]byte(EventValue), &rs)
	if err != nil {
		return nil, err, false
	}

	votes, err := b.RetroService.GroupUserSubtractVote(RetroID, rs.GroupId, UserID)
	if err != nil {
		return nil, err, false
	}

	updatedVotes, _ := json.Marshal(votes)
	msg := createSocketEvent("votes_updated", string(updatedVotes), "")

	return msg, nil, false
}

// CreateAction creates a retro action
func (b *Service) CreateAction(ctx context.Context, RetroID string, UserID string, EventValue string) ([]byte, error, bool) {
	var rs struct {
		Content string `json:"content"`
	}
	err := json.Unmarshal([]byte(EventValue), &rs)
	if err != nil {
		return nil, err, false
	}

	items, err := b.RetroService.CreateRetroAction(RetroID, UserID, rs.Content)
	if err != nil {
		return nil, err, false
	}

	updatedItems, _ := json.Marshal(items)
	msg := createSocketEvent("action_updated", string(updatedItems), "")

	return msg, nil, false
}

// UpdateAction updates a retro action
func (b *Service) UpdateAction(ctx context.Context, RetroID string, UserID string, EventValue string) ([]byte, error, bool) {
	var rs struct {
		ActionID  string `json:"id"`
		Completed bool   `json:"completed"`
		Content   string `json:"content"`
	}
	err := json.Unmarshal([]byte(EventValue), &rs)
	if err != nil {
		return nil, err, false
	}

	items, err := b.RetroService.UpdateRetroAction(RetroID, rs.ActionID, rs.Content, rs.Completed)
	if err != nil {
		return nil, err, false
	}

	updatedItems, _ := json.Marshal(items)
	msg := createSocketEvent("action_updated", string(updatedItems), "")

	return msg, nil, false
}

// DeleteAction deletes a retro action
func (b *Service) DeleteAction(ctx context.Context, RetroID string, UserID string, EventValue string) ([]byte, error, bool) {
	var rs struct {
		ActionID string `json:"id"`
	}
	err := json.Unmarshal([]byte(EventValue), &rs)
	if err != nil {
		return nil, err, false
	}

	items, err := b.RetroService.DeleteRetroAction(RetroID, UserID, rs.ActionID)
	if err != nil {
		return nil, err, false
	}

	updatedItems, _ := json.Marshal(items)
	msg := createSocketEvent("action_updated", string(updatedItems), "")

	return msg, nil, false
}

// AdvancePhase updates a retro phase
func (b *Service) AdvancePhase(ctx context.Context, RetroID string, UserID string, EventValue string) ([]byte, error, bool) {
	var rs struct {
		Phase string `json:"phase"`
	}
	err := json.Unmarshal([]byte(EventValue), &rs)
	if err != nil {
		return nil, err, false
	}

	retro, err := b.RetroService.RetroAdvancePhase(RetroID, rs.Phase)
	if err != nil {
		return nil, err, false
	}

	updatedItems, _ := json.Marshal(retro)
	msg := createSocketEvent("phase_updated", string(updatedItems), "")

	// if retro is completed send retro email to attendees
	if rs.Phase == "completed" {
		go b.SendCompletedEmails(retro)
	}

	return msg, nil, false
}

// FacilitatorAdd adds a user as facilitator of the retro
func (b *Service) FacilitatorAdd(ctx context.Context, RetroID string, UserID string, EventValue string) ([]byte, error, bool) {
	var rs struct {
		UserID string `json:"userId"`
	}
	err := json.Unmarshal([]byte(EventValue), &rs)
	if err != nil {
		return nil, err, false
	}

	facilitators, err := b.RetroService.RetroFacilitatorAdd(RetroID, rs.UserID)
	if err != nil {
		return nil, err, false
	}
	updatedFacilitators, _ := json.Marshal(facilitators)

	msg := createSocketEvent("facilitators_updated", string(updatedFacilitators), "")

	return msg, nil, false
}

// FacilitatorRemove removes a retro facilitator
func (b *Service) FacilitatorRemove(ctx context.Context, RetroID string, UserID string, EventValue string) ([]byte, error, bool) {
	var rs struct {
		UserID string `json:"userId"`
	}
	err := json.Unmarshal([]byte(EventValue), &rs)
	if err != nil {
		return nil, err, false
	}

	facilitators, err := b.RetroService.RetroFacilitatorRemove(RetroID, rs.UserID)
	if err != nil {
		return nil, err, false
	}
	updatedFacilitators, _ := json.Marshal(facilitators)

	msg := createSocketEvent("facilitators_updated", string(updatedFacilitators), "")

	return msg, nil, false
}

// FacilitatorSelf handles self-promoting a user to a facilitator
func (b *Service) FacilitatorSelf(ctx context.Context, RetroID string, UserID string, EventValue string) ([]byte, error, bool) {
	facilitatorCode, err := b.RetroService.GetRetroFacilitatorCode(RetroID)
	if err != nil {
		return nil, err, false
	}

	if EventValue == facilitatorCode {
		facilitators, err := b.RetroService.RetroFacilitatorAdd(RetroID, UserID)
		if err != nil {
			return nil, err, false
		}
		updatedFacilitators, _ := json.Marshal(facilitators)

		msg := createSocketEvent("facilitators_updated", string(updatedFacilitators), "")

		return msg, nil, false
	} else {
		return nil, errors.New("INCORRECT_FACILITATOR_CODE"), false
	}
}

// EditRetro handles editing the retro settings
func (b *Service) EditRetro(ctx context.Context, RetroID string, UserID string, EventValue string) ([]byte, error, bool) {
	var rb struct {
		Name                 string `json:"retroName"`
		JoinCode             string `json:"joinCode"`
		FacilitatorCode      string `json:"facilitatorCode"`
		MaxVotes             int    `json:"maxVotes"`
		BrainstormVisibility string `json:"brainstormVisibility"`
	}
	err := json.Unmarshal([]byte(EventValue), &rb)
	if err != nil {
		return nil, err, false
	}

	err = b.RetroService.EditRetro(
		RetroID,
		rb.Name,
		rb.JoinCode,
		rb.FacilitatorCode,
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
func (b *Service) Delete(ctx context.Context, RetroID string, UserID string, EventValue string) ([]byte, error, bool) {
	err := b.RetroService.RetroDelete(RetroID)
	if err != nil {
		return nil, err, false
	}
	msg := createSocketEvent("conceded", "", "")

	return msg, nil, false
}

// Abandon handles setting abandoned true so retro doesn't show up in users retro list, then leaves retro
func (b *Service) Abandon(ctx context.Context, RetroID string, UserID string, EventValue string) ([]byte, error, bool) {
	_, err := b.RetroService.RetroAbandon(RetroID, UserID)
	if err != nil {
		return nil, err, false
	}

	return nil, errors.New("ABANDONED_RETRO"), true
}

// SendCompletedEmails sends an email to attendees with the retro items and actions
func (b *Service) SendCompletedEmails(retro *thunderdome.Retro) {
	users := b.RetroService.RetroGetUsers(retro.Id)

	for _, user := range users {
		// don't send emails to guest's as they have no email
		if user.Email != "" {
			b.EmailService.SendRetroOverview(retro, user.Name, user.Email)
		}
	}
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
