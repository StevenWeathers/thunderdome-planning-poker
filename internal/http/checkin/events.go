package checkin

import (
	"context"
	"encoding/json"

	"github.com/StevenWeathers/thunderdome-planning-poker/internal/wshub"
)

// CheckinCreate creates a checkin
func (b *Service) CheckinCreate(ctx context.Context, teamID string, userID string, eventValue string) (any, []byte, error, bool) {
	var c struct {
		UserID    string `json:"userId"`
		Yesterday string `json:"yesterday"`
		Today     string `json:"today"`
		Blockers  string `json:"blockers"`
		Discuss   string `json:"discuss"`
		GoalsMet  bool   `json:"goalsMet"`
	}
	err := json.Unmarshal([]byte(eventValue), &c)
	if err != nil {
		return nil, nil, err, false
	}

	err = b.CheckinService.CheckinCreate(context.Background(), teamID, c.UserID, c.Yesterday, c.Today, c.Blockers, c.Discuss, c.GoalsMet)
	if err != nil {
		return nil, nil, err, false
	}

	msg := wshub.CreateSocketEvent("checkin_added", "", "")

	return nil, msg, nil, false
}

// CheckinUpdate updates a checkin
func (b *Service) CheckinUpdate(ctx context.Context, teamID string, userID string, eventValue string) (any, []byte, error, bool) {
	var c struct {
		CheckinID string `json:"checkinId"`
		Yesterday string `json:"yesterday"`
		Today     string `json:"today"`
		Blockers  string `json:"blockers"`
		Discuss   string `json:"discuss"`
		GoalsMet  bool   `json:"goalsMet"`
	}
	err := json.Unmarshal([]byte(eventValue), &c)
	if err != nil {
		return nil, nil, err, false
	}

	err = b.CheckinService.CheckinUpdate(context.Background(), c.CheckinID, c.Yesterday, c.Today, c.Blockers, c.Discuss, c.GoalsMet)
	if err != nil {
		return nil, nil, err, false
	}

	msg := wshub.CreateSocketEvent("checkin_updated", "", "")

	return nil, msg, nil, false
}

// CheckinDelete deletes a checkin
func (b *Service) CheckinDelete(ctx context.Context, teamID string, userID string, eventValue string) (any, []byte, error, bool) {
	var c struct {
		CheckinID string `json:"checkinId"`
	}
	err := json.Unmarshal([]byte(eventValue), &c)
	if err != nil {
		return nil, nil, err, false
	}

	err = b.CheckinService.CheckinDelete(context.Background(), c.CheckinID)
	if err != nil {
		return nil, nil, err, false
	}

	msg := wshub.CreateSocketEvent("checkin_deleted", "", "")

	return nil, msg, nil, false
}

// CommentCreate creates a checkin comment
func (b *Service) CommentCreate(ctx context.Context, teamID string, userID string, eventValue string) (any, []byte, error, bool) {
	var c struct {
		CheckinID string `json:"checkinId"`
		UserID    string `json:"userId"`
		Comment   string `json:"comment"`
	}
	err := json.Unmarshal([]byte(eventValue), &c)
	if err != nil {
		return nil, nil, err, false
	}

	err = b.CheckinService.CheckinComment(ctx, teamID, c.CheckinID, c.UserID, c.Comment)
	if err != nil {
		return nil, nil, err, false
	}

	msg := wshub.CreateSocketEvent("comment_added", "", "")

	return nil, msg, nil, false
}

// CommentUpdate updates a checkin comment
func (b *Service) CommentUpdate(ctx context.Context, teamID string, userID string, eventValue string) (any, []byte, error, bool) {
	var c struct {
		CommentID string `json:"commentId"`
		UserID    string `json:"userId"`
		Comment   string `json:"comment"`
	}
	err := json.Unmarshal([]byte(eventValue), &c)
	if err != nil {
		return nil, nil, err, false
	}

	err = b.CheckinService.CheckinCommentEdit(ctx, teamID, c.UserID, c.CommentID, c.Comment)
	if err != nil {
		return nil, nil, err, false
	}

	msg := wshub.CreateSocketEvent("comment_updated", "", "")

	return nil, msg, nil, false
}

// CommentDelete deletes a checkin comment
func (b *Service) CommentDelete(ctx context.Context, teamID string, userID string, eventValue string) (any, []byte, error, bool) {
	var c struct {
		CommentID string `json:"commentId"`
	}
	err := json.Unmarshal([]byte(eventValue), &c)
	if err != nil {
		return nil, nil, err, false
	}

	err = b.CheckinService.CheckinCommentDelete(ctx, c.CommentID)
	if err != nil {
		return nil, nil, err, false
	}

	msg := wshub.CreateSocketEvent("comment_deleted", "", "")

	return nil, msg, nil, false
}
