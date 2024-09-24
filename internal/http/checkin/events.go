package checkin

import (
	"context"
	"encoding/json"

	"github.com/StevenWeathers/thunderdome-planning-poker/internal/wshub"
)

// CheckinCreate creates a checkin
func (b *Service) CheckinCreate(ctx context.Context, TeamID string, UserID string, EventValue string) ([]byte, error, bool) {
	var c struct {
		UserId    string `json:"userId"`
		Yesterday string `json:"yesterday"`
		Today     string `json:"today"`
		Blockers  string `json:"blockers"`
		Discuss   string `json:"discuss"`
		GoalsMet  bool   `json:"goalsMet"`
	}
	err := json.Unmarshal([]byte(EventValue), &c)
	if err != nil {
		return nil, err, false
	}

	err = b.CheckinService.CheckinCreate(context.Background(), TeamID, c.UserId, c.Yesterday, c.Today, c.Blockers, c.Discuss, c.GoalsMet)
	if err != nil {
		return nil, err, false
	}

	msg := wshub.CreateSocketEvent("checkin_added", "", "")

	return msg, nil, false
}

// CheckinUpdate updates a checkin
func (b *Service) CheckinUpdate(ctx context.Context, TeamID string, UserID string, EventValue string) ([]byte, error, bool) {
	var c struct {
		CheckinId string `json:"checkinId"`
		Yesterday string `json:"yesterday"`
		Today     string `json:"today"`
		Blockers  string `json:"blockers"`
		Discuss   string `json:"discuss"`
		GoalsMet  bool   `json:"goalsMet"`
	}
	err := json.Unmarshal([]byte(EventValue), &c)
	if err != nil {
		return nil, err, false
	}

	err = b.CheckinService.CheckinUpdate(context.Background(), c.CheckinId, c.Yesterday, c.Today, c.Blockers, c.Discuss, c.GoalsMet)
	if err != nil {
		return nil, err, false
	}

	msg := wshub.CreateSocketEvent("checkin_updated", "", "")

	return msg, nil, false
}

// CheckinDelete deletes a checkin
func (b *Service) CheckinDelete(ctx context.Context, TeamID string, UserID string, EventValue string) ([]byte, error, bool) {
	var c struct {
		CheckinId string `json:"checkinId"`
	}
	err := json.Unmarshal([]byte(EventValue), &c)
	if err != nil {
		return nil, err, false
	}

	err = b.CheckinService.CheckinDelete(context.Background(), c.CheckinId)
	if err != nil {
		return nil, err, false
	}

	msg := wshub.CreateSocketEvent("checkin_deleted", "", "")

	return msg, nil, false
}

// CommentCreate creates a checkin comment
func (b *Service) CommentCreate(ctx context.Context, TeamID string, UserID string, EventValue string) ([]byte, error, bool) {
	var c struct {
		CheckinId string `json:"checkinId"`
		UserID    string `json:"userId"`
		Comment   string `json:"comment"`
	}
	err := json.Unmarshal([]byte(EventValue), &c)
	if err != nil {
		return nil, err, false
	}

	err = b.CheckinService.CheckinComment(ctx, TeamID, c.CheckinId, c.UserID, c.Comment)
	if err != nil {
		return nil, err, false
	}

	msg := wshub.CreateSocketEvent("comment_added", "", "")

	return msg, nil, false
}

// CommentUpdate updates a checkin comment
func (b *Service) CommentUpdate(ctx context.Context, TeamID string, UserID string, EventValue string) ([]byte, error, bool) {
	var c struct {
		CommentId string `json:"commentId"`
		UserID    string `json:"userId"`
		Comment   string `json:"comment"`
	}
	err := json.Unmarshal([]byte(EventValue), &c)
	if err != nil {
		return nil, err, false
	}

	err = b.CheckinService.CheckinCommentEdit(ctx, TeamID, c.UserID, c.CommentId, c.Comment)
	if err != nil {
		return nil, err, false
	}

	msg := wshub.CreateSocketEvent("comment_updated", "", "")

	return msg, nil, false
}

// CommentDelete deletes a checkin comment
func (b *Service) CommentDelete(ctx context.Context, TeamID string, UserID string, EventValue string) ([]byte, error, bool) {
	var c struct {
		CommentId string `json:"commentId"`
	}
	err := json.Unmarshal([]byte(EventValue), &c)
	if err != nil {
		return nil, err, false
	}

	err = b.CheckinService.CheckinCommentDelete(ctx, c.CommentId)
	if err != nil {
		return nil, err, false
	}

	msg := wshub.CreateSocketEvent("comment_deleted", "", "")

	return msg, nil, false
}
