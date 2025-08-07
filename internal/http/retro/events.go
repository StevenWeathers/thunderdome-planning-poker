package retro

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/StevenWeathers/thunderdome-planning-poker/internal/wshub"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
	"go.uber.org/zap"
)

// CreateItem creates a retro item
func (b *Service) CreateItem(ctx context.Context, RetroID string, UserID string, EventValue string) (any, []byte, error, bool) {
	var rs struct {
		Type    string `json:"type"`
		Content string `json:"content"`
		Phase   string `json:"phase"`
	}
	err := json.Unmarshal([]byte(EventValue), &rs)
	if err != nil {
		return nil, nil, err, false
	}

	items, err := b.RetroService.CreateRetroItem(RetroID, UserID, rs.Type, rs.Content)
	if err != nil {
		return nil, nil, err, false
	}

	updatedItems, _ := json.Marshal(items)
	msg := wshub.CreateSocketEvent("items_updated", string(updatedItems), "")

	return nil, msg, nil, false
}

// ItemCommentAdd creates a retro item comment
func (b *Service) ItemCommentAdd(ctx context.Context, RetroID string, UserID string, EventValue string) (any, []byte, error, bool) {
	var rs struct {
		ItemID  string `json:"item_id"`
		Comment string `json:"comment"`
	}
	err := json.Unmarshal([]byte(EventValue), &rs)
	if err != nil {
		return nil, nil, err, false
	}

	items, err := b.RetroService.ItemCommentAdd(RetroID, rs.ItemID, UserID, rs.Comment)
	if err != nil {
		return nil, nil, err, false
	}

	updatedItems, _ := json.Marshal(items)
	msg := wshub.CreateSocketEvent("items_updated", string(updatedItems), "")

	return nil, msg, nil, false
}

// ItemCommentEdit updates a retro item comment
func (b *Service) ItemCommentEdit(ctx context.Context, RetroID string, UserID string, EventValue string) (any, []byte, error, bool) {
	var rs struct {
		CommentID string `json:"comment_id"`
		Comment   string `json:"comment"`
	}
	err := json.Unmarshal([]byte(EventValue), &rs)
	if err != nil {
		return nil, nil, err, false
	}

	items, err := b.RetroService.ItemCommentEdit(RetroID, rs.CommentID, rs.Comment)
	if err != nil {
		return nil, nil, err, false
	}

	updatedItems, _ := json.Marshal(items)
	msg := wshub.CreateSocketEvent("items_updated", string(updatedItems), "")

	return nil, msg, nil, false
}

// ItemCommentDelete deletes a retro item comment
func (b *Service) ItemCommentDelete(ctx context.Context, RetroID string, UserID string, EventValue string) (any, []byte, error, bool) {
	var rs struct {
		CommentID string `json:"comment_id"`
	}
	err := json.Unmarshal([]byte(EventValue), &rs)
	if err != nil {
		return nil, nil, err, false
	}

	items, err := b.RetroService.ItemCommentDelete(RetroID, rs.CommentID)
	if err != nil {
		return nil, nil, err, false
	}

	updatedItems, _ := json.Marshal(items)
	msg := wshub.CreateSocketEvent("items_updated", string(updatedItems), "")

	return nil, msg, nil, false
}

// UserMarkReady marks a user as ready to advance to next phase
func (b *Service) UserMarkReady(ctx context.Context, RetroID string, UserID string, EventValue string) (any, []byte, error, bool) {
	readyUsers, err := b.RetroService.MarkUserReady(RetroID, UserID)
	if err != nil {
		return nil, nil, err, false
	}

	updatedReadyUsers, _ := json.Marshal(readyUsers)
	msg := wshub.CreateSocketEvent("user_marked_ready", string(updatedReadyUsers), UserID)

	return nil, msg, nil, false
}

// UserUnMarkReady unsets a user from ready to advance to next phase
func (b *Service) UserUnMarkReady(ctx context.Context, RetroID string, UserID string, EventValue string) (any, []byte, error, bool) {
	readyUsers, err := b.RetroService.UnmarkUserReady(RetroID, UserID)
	if err != nil {
		return nil, nil, err, false
	}

	updatedReadyUsers, _ := json.Marshal(readyUsers)
	msg := wshub.CreateSocketEvent("user_marked_unready", string(updatedReadyUsers), UserID)

	return nil, msg, nil, false
}

// GroupItem changes a retro item's group_id
func (b *Service) GroupItem(ctx context.Context, RetroID string, UserID string, EventValue string) (any, []byte, error, bool) {
	var rs struct {
		ItemID  string `json:"itemId"`
		GroupID string `json:"groupId"`
	}
	err := json.Unmarshal([]byte(EventValue), &rs)
	if err != nil {
		return nil, nil, err, false
	}

	item, err := b.RetroService.GroupRetroItem(RetroID, rs.ItemID, rs.GroupID)
	if err != nil {
		return nil, nil, err, false
	}

	updatedItem, _ := json.Marshal(item)
	msg := wshub.CreateSocketEvent("item_moved", string(updatedItem), "")

	return nil, msg, nil, false
}

// DeleteItem deletes a retro item
func (b *Service) DeleteItem(ctx context.Context, RetroID string, UserID string, EventValue string) (any, []byte, error, bool) {
	var rs struct {
		ItemID string `json:"id"`
		Phase  string `json:"phase"`
		Type   string `json:"type"`
	}
	err := json.Unmarshal([]byte(EventValue), &rs)
	if err != nil {
		return nil, nil, err, false
	}

	items, err := b.RetroService.DeleteRetroItem(RetroID, UserID, rs.Type, rs.ItemID)
	if err != nil {
		return nil, nil, err, false
	}

	updatedItems, _ := json.Marshal(items)
	msg := wshub.CreateSocketEvent("items_updated", string(updatedItems), "")

	return nil, msg, nil, false
}

// GroupNameChange changes a retro group's name
func (b *Service) GroupNameChange(ctx context.Context, RetroID string, UserID string, EventValue string) (any, []byte, error, bool) {
	var rs struct {
		GroupID string `json:"groupId"`
		Name    string `json:"name"`
	}
	err := json.Unmarshal([]byte(EventValue), &rs)
	if err != nil {
		return nil, nil, err, false
	}

	group, err := b.RetroService.GroupNameChange(RetroID, rs.GroupID, rs.Name)
	if err != nil {
		return nil, nil, err, false
	}

	updatedGroup, _ := json.Marshal(group)
	msg := wshub.CreateSocketEvent("group_name_updated", string(updatedGroup), "")

	return nil, msg, nil, false
}

// GroupUserVote handles a users vote for an item group
func (b *Service) GroupUserVote(ctx context.Context, RetroID string, UserID string, EventValue string) (any, []byte, error, bool) {
	var rs struct {
		GroupID string `json:"groupId"`
	}
	err := json.Unmarshal([]byte(EventValue), &rs)
	if err != nil {
		return nil, nil, err, false
	}

	votes, err := b.RetroService.GroupUserVote(RetroID, rs.GroupID, UserID)
	if err != nil {
		return nil, nil, err, false
	}

	updatedVotes, _ := json.Marshal(votes)
	msg := wshub.CreateSocketEvent("votes_updated", string(updatedVotes), "")

	return nil, msg, nil, false
}

// GroupUserSubtractVote handles removing a users vote from an item group
func (b *Service) GroupUserSubtractVote(ctx context.Context, RetroID string, UserID string, EventValue string) (any, []byte, error, bool) {
	var rs struct {
		GroupID string `json:"groupId"`
	}
	err := json.Unmarshal([]byte(EventValue), &rs)
	if err != nil {
		return nil, nil, err, false
	}

	votes, err := b.RetroService.GroupUserSubtractVote(RetroID, rs.GroupID, UserID)
	if err != nil {
		return nil, nil, err, false
	}

	updatedVotes, _ := json.Marshal(votes)
	msg := wshub.CreateSocketEvent("votes_updated", string(updatedVotes), "")

	return nil, msg, nil, false
}

// CreateAction creates a retro action
func (b *Service) CreateAction(ctx context.Context, RetroID string, UserID string, EventValue string) (any, []byte, error, bool) {
	var rs struct {
		Content string `json:"content"`
	}
	err := json.Unmarshal([]byte(EventValue), &rs)
	if err != nil {
		return nil, nil, err, false
	}

	items, err := b.RetroService.CreateRetroAction(RetroID, UserID, rs.Content)
	if err != nil {
		return nil, nil, err, false
	}

	updatedItems, _ := json.Marshal(items)
	msg := wshub.CreateSocketEvent("action_updated", string(updatedItems), "")

	return nil, msg, nil, false
}

// UpdateAction updates a retro action
func (b *Service) UpdateAction(ctx context.Context, RetroID string, UserID string, EventValue string) (any, []byte, error, bool) {
	var rs struct {
		ActionID  string `json:"id"`
		Completed bool   `json:"completed"`
		Content   string `json:"content"`
	}
	err := json.Unmarshal([]byte(EventValue), &rs)
	if err != nil {
		return nil, nil, err, false
	}

	items, err := b.RetroService.UpdateRetroAction(RetroID, rs.ActionID, rs.Content, rs.Completed)
	if err != nil {
		return nil, nil, err, false
	}

	updatedItems, _ := json.Marshal(items)
	msg := wshub.CreateSocketEvent("action_updated", string(updatedItems), "")

	return nil, msg, nil, false
}

// ActionAddAssignee adds a retro action assignee
func (b *Service) ActionAddAssignee(ctx context.Context, RetroID string, UserID string, EventValue string) (any, []byte, error, bool) {
	var rs struct {
		ActionID string `json:"id"`
		UserID   string `json:"user_id"`
	}
	err := json.Unmarshal([]byte(EventValue), &rs)
	if err != nil {
		return nil, nil, err, false
	}

	items, err := b.RetroService.RetroActionAssigneeAdd(RetroID, rs.ActionID, rs.UserID)
	if err != nil {
		return nil, nil, err, false
	}

	updatedItems, _ := json.Marshal(items)
	msg := wshub.CreateSocketEvent("action_updated", string(updatedItems), "")

	return nil, msg, nil, false
}

// ActionRemoveAssignee removes a retro action assignee
func (b *Service) ActionRemoveAssignee(ctx context.Context, RetroID string, UserID string, EventValue string) (any, []byte, error, bool) {
	var rs struct {
		ActionID string `json:"id"`
		UserID   string `json:"user_id"`
	}
	err := json.Unmarshal([]byte(EventValue), &rs)
	if err != nil {
		return nil, nil, err, false
	}

	items, err := b.RetroService.RetroActionAssigneeDelete(RetroID, rs.ActionID, rs.UserID)
	if err != nil {
		return nil, nil, err, false
	}

	updatedItems, _ := json.Marshal(items)
	msg := wshub.CreateSocketEvent("action_updated", string(updatedItems), "")

	return nil, msg, nil, false
}

// DeleteAction deletes a retro action
func (b *Service) DeleteAction(ctx context.Context, RetroID string, UserID string, EventValue string) (any, []byte, error, bool) {
	var rs struct {
		ActionID string `json:"id"`
	}
	err := json.Unmarshal([]byte(EventValue), &rs)
	if err != nil {
		return nil, nil, err, false
	}

	items, err := b.RetroService.DeleteRetroAction(RetroID, UserID, rs.ActionID)
	if err != nil {
		return nil, nil, err, false
	}

	updatedItems, _ := json.Marshal(items)
	msg := wshub.CreateSocketEvent("action_updated", string(updatedItems), "")

	return nil, msg, nil, false
}

// AdvancePhase updates a retro phase
func (b *Service) AdvancePhase(ctx context.Context, RetroID string, UserID string, EventValue string) (any, []byte, error, bool) {
	var rs struct {
		Phase string `json:"phase"`
	}
	err := json.Unmarshal([]byte(EventValue), &rs)
	if err != nil {
		return nil, nil, err, false
	}

	retro, err := b.RetroService.RetroAdvancePhase(RetroID, rs.Phase)
	if err != nil {
		return nil, nil, err, false
	}

	updatedItems, _ := json.Marshal(retro)
	msg := wshub.CreateSocketEvent("phase_updated", string(updatedItems), "")

	// if retro is completed send retro email to attendees
	if rs.Phase == "completed" {
		go b.SendCompletedEmails(retro)
	}

	return nil, msg, nil, false
}

// PhaseTimeout advances a retro phase after time countdown
func (b *Service) PhaseTimeout(ctx context.Context, RetroID string, UserID string, EventValue string) (any, []byte, error, bool) {
	var rs struct {
		Phase string `json:"phase"`
	}
	err := json.Unmarshal([]byte(EventValue), &rs)
	if err != nil {
		return nil, nil, err, false
	}

	retro, err := b.RetroService.RetroAdvancePhase(RetroID, rs.Phase)
	if err != nil {
		return nil, nil, err, false
	}

	updatedItems, _ := json.Marshal(retro)
	msg := wshub.CreateSocketEvent("phase_updated", string(updatedItems), "")

	return nil, msg, nil, false
}

// PhaseAllReady advances a retro phase after all users are ready
func (b *Service) PhaseAllReady(ctx context.Context, RetroID string, UserID string, EventValue string) (any, []byte, error, bool) {
	var rs struct {
		Phase string `json:"phase"`
	}
	err := json.Unmarshal([]byte(EventValue), &rs)
	if err != nil {
		return nil, nil, err, false
	}

	retro, err := b.RetroService.RetroAdvancePhase(RetroID, rs.Phase)
	if err != nil {
		return nil, nil, err, false
	}

	updatedItems, _ := json.Marshal(retro)
	msg := wshub.CreateSocketEvent("phase_updated", string(updatedItems), "")

	return nil, msg, nil, false
}

// FacilitatorAdd adds a user as facilitator of the retro
func (b *Service) FacilitatorAdd(ctx context.Context, RetroID string, UserID string, EventValue string) (any, []byte, error, bool) {
	var rs struct {
		UserID string `json:"userId"`
	}
	err := json.Unmarshal([]byte(EventValue), &rs)
	if err != nil {
		return nil, nil, err, false
	}

	facilitators, err := b.RetroService.RetroFacilitatorAdd(RetroID, rs.UserID)
	if err != nil {
		return nil, nil, err, false
	}
	updatedFacilitators, _ := json.Marshal(facilitators)

	msg := wshub.CreateSocketEvent("facilitators_updated", string(updatedFacilitators), "")

	return nil, msg, nil, false
}

// FacilitatorRemove removes a retro facilitator
func (b *Service) FacilitatorRemove(ctx context.Context, RetroID string, UserID string, EventValue string) (any, []byte, error, bool) {
	var rs struct {
		UserID string `json:"userId"`
	}
	err := json.Unmarshal([]byte(EventValue), &rs)
	if err != nil {
		return nil, nil, err, false
	}

	facilitators, err := b.RetroService.RetroFacilitatorRemove(RetroID, rs.UserID)
	if err != nil {
		return nil, nil, err, false
	}
	updatedFacilitators, _ := json.Marshal(facilitators)

	msg := wshub.CreateSocketEvent("facilitators_updated", string(updatedFacilitators), "")

	return nil, msg, nil, false
}

// FacilitatorSelf handles self-promoting a user to a facilitator
func (b *Service) FacilitatorSelf(ctx context.Context, RetroID string, UserID string, EventValue string) (any, []byte, error, bool) {
	facilitatorCode, err := b.RetroService.GetRetroFacilitatorCode(RetroID)
	if err != nil {
		return nil, nil, err, false
	}

	if EventValue == facilitatorCode {
		facilitators, err := b.RetroService.RetroFacilitatorAdd(RetroID, UserID)
		if err != nil {
			return nil, nil, err, false
		}
		updatedFacilitators, _ := json.Marshal(facilitators)

		msg := wshub.CreateSocketEvent("facilitators_updated", string(updatedFacilitators), "")

		return nil, msg, nil, false
	} else {
		return nil, nil, errors.New("INCORRECT_FACILITATOR_CODE"), false
	}
}

// EditRetro handles editing the retro settings
func (b *Service) EditRetro(ctx context.Context, RetroID string, UserID string, EventValue string) (any, []byte, error, bool) {
	var rb struct {
		Name                 string `json:"retroName"`
		JoinCode             string `json:"joinCode"`
		FacilitatorCode      string `json:"facilitatorCode"`
		MaxVotes             int    `json:"maxVotes"`
		BrainstormVisibility string `json:"brainstormVisibility"`
		PhaseAutoAdvance     bool   `json:"phase_auto_advance"`
	}
	err := json.Unmarshal([]byte(EventValue), &rb)
	if err != nil {
		return nil, nil, err, false
	}

	err = b.RetroService.EditRetro(
		RetroID,
		rb.Name,
		rb.JoinCode,
		rb.FacilitatorCode,
		rb.MaxVotes,
		rb.BrainstormVisibility,
		rb.PhaseAutoAdvance,
	)
	if err != nil {
		return nil, nil, err, false
	}

	updatedRetro, _ := json.Marshal(rb)
	msg := wshub.CreateSocketEvent("retro_edited", string(updatedRetro), "")

	return nil, msg, nil, false
}

// Delete handles deleting the retro
func (b *Service) Delete(ctx context.Context, RetroID string, UserID string, EventValue string) (any, []byte, error, bool) {
	err := b.RetroService.RetroDelete(RetroID)
	if err != nil {
		return nil, nil, err, false
	}
	msg := wshub.CreateSocketEvent("conceded", "", "")

	return nil, msg, nil, false
}

// Abandon handles setting abandoned true so retro doesn't show up in users retro list, then leaves retro
func (b *Service) Abandon(ctx context.Context, RetroID string, UserID string, EventValue string) (any, []byte, error, bool) {
	_, err := b.RetroService.RetroAbandon(RetroID, UserID)
	if err != nil {
		return nil, nil, err, false
	}

	return nil, nil, errors.New("ABANDONED_RETRO"), true
}

// SendCompletedEmails sends an email to attendees with the retro items and actions
func (b *Service) SendCompletedEmails(retro *thunderdome.Retro) {
	users := b.RetroService.RetroGetUsers(retro.ID)

	for _, user := range users {
		// don't send emails to guest's as they have no email
		if user.Email != "" {
			template, err := b.TemplateService.GetTemplateByID(context.Background(), retro.TemplateID)
			if err != nil {
				b.logger.Error("Error getting template", zap.Error(err))
			} else {
				err := b.EmailService.SendRetroOverview(retro, template, user.Name, user.Email)
				if err != nil {
					b.logger.Error("Error sending retro overview email", zap.Error(err))
				}
			}
		}
	}
}
