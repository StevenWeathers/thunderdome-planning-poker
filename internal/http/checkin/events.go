package checkin

import (
	"context"
	"encoding/json"

	"github.com/StevenWeathers/thunderdome-planning-poker/internal/wshub"
)

// CheckinCreate creates a checkin
func (s *Service) CheckinCreate(ctx context.Context, teamID string, userID string, eventValue string) (any, []byte, error, bool) {
	var c struct {
		CheckinDate string `json:"checkinDate"`
		TimeZone    string `json:"timeZone"`
		Yesterday   string `json:"yesterday"`
		Today       string `json:"today"`
		Blockers    string `json:"blockers"`
		Discuss     string `json:"discuss"`
		GoalsMet    bool   `json:"goalsMet"`
	}
	err := json.Unmarshal([]byte(eventValue), &c)
	if err != nil {
		return nil, nil, err, false
	}

	err = s.CheckinService.CheckinCreate(context.Background(), teamID, userID, c.CheckinDate, c.Yesterday, c.Today, c.Blockers, c.Discuss, c.GoalsMet)
	if err != nil {
		return nil, nil, err, false
	}

	msg := wshub.CreateSocketEvent("checkin_added", "", "")

	return nil, msg, nil, false
}

// CheckinUpdate updates a checkin
func (s *Service) CheckinUpdate(ctx context.Context, teamID string, userID string, eventValue string) (any, []byte, error, bool) {
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

	err = s.CheckinService.CheckinUpdate(context.Background(), teamID, userID, c.CheckinID, c.Yesterday, c.Today, c.Blockers, c.Discuss, c.GoalsMet)
	if err != nil {
		return nil, nil, err, false
	}

	msg := wshub.CreateSocketEvent("checkin_updated", "", "")

	return nil, msg, nil, false
}

// CheckinDelete deletes a checkin
func (s *Service) CheckinDelete(ctx context.Context, teamID string, userID string, eventValue string) (any, []byte, error, bool) {
	var c struct {
		CheckinID string `json:"checkinId"`
	}
	err := json.Unmarshal([]byte(eventValue), &c)
	if err != nil {
		return nil, nil, err, false
	}

	err = s.CheckinService.CheckinDelete(context.Background(), teamID, userID, c.CheckinID)
	if err != nil {
		return nil, nil, err, false
	}

	msg := wshub.CreateSocketEvent("checkin_deleted", "", "")

	return nil, msg, nil, false
}

// KudoCreate creates a team kudo.
func (s *Service) KudoCreate(ctx context.Context, teamID string, userID string, eventValue string) (any, []byte, error, bool) {
	var c struct {
		TargetUserID string `json:"targetUserId"`
		KudosDate    string `json:"kudosDate"`
		Comment      string `json:"comment"`
	}
	err := json.Unmarshal([]byte(eventValue), &c)
	if err != nil {
		return nil, nil, err, false
	}

	kudo, err := s.CheckinService.KudoCreate(ctx, teamID, userID, c.TargetUserID, c.KudosDate, c.Comment)
	if err != nil {
		return nil, nil, err, false
	}

	msg := wshub.CreateSocketEvent("kudo_added", "", "")

	return kudo, msg, nil, false
}

// CommentCreate creates a checkin comment
func (s *Service) CommentCreate(ctx context.Context, teamID string, userID string, eventValue string) (any, []byte, error, bool) {
	// the author is always the acting user, never one named in the payload
	var c struct {
		CheckinID string `json:"checkinId"`
		Comment   string `json:"comment"`
	}
	err := json.Unmarshal([]byte(eventValue), &c)
	if err != nil {
		return nil, nil, err, false
	}

	err = s.CheckinService.CheckinComment(ctx, teamID, c.CheckinID, userID, c.Comment)
	if err != nil {
		return nil, nil, err, false
	}

	msg := wshub.CreateSocketEvent("comment_added", "", "")

	return nil, msg, nil, false
}

// CommentUpdate updates a checkin comment
func (s *Service) CommentUpdate(ctx context.Context, teamID string, userID string, eventValue string) (any, []byte, error, bool) {
	var c struct {
		CommentID string `json:"commentId"`
		Comment   string `json:"comment"`
	}
	err := json.Unmarshal([]byte(eventValue), &c)
	if err != nil {
		return nil, nil, err, false
	}

	err = s.CheckinService.CheckinCommentEdit(ctx, teamID, userID, c.CommentID, c.Comment)
	if err != nil {
		return nil, nil, err, false
	}

	msg := wshub.CreateSocketEvent("comment_updated", "", "")

	return nil, msg, nil, false
}

// CommentDelete deletes a checkin comment
func (s *Service) CommentDelete(ctx context.Context, teamID string, userID string, eventValue string) (any, []byte, error, bool) {
	var c struct {
		CommentID string `json:"commentId"`
	}
	err := json.Unmarshal([]byte(eventValue), &c)
	if err != nil {
		return nil, nil, err, false
	}

	err = s.CheckinService.CheckinCommentDelete(ctx, teamID, userID, c.CommentID)
	if err != nil {
		return nil, nil, err, false
	}

	msg := wshub.CreateSocketEvent("comment_deleted", "", "")

	return nil, msg, nil, false
}

// KudoUpdate updates a team kudo.
func (s *Service) KudoUpdate(ctx context.Context, teamID string, userID string, eventValue string) (any, []byte, error, bool) {
	var c struct {
		KudoID       string `json:"kudoId"`
		TargetUserID string `json:"targetUserId"`
		KudosDate    string `json:"kudosDate"`
		Comment      string `json:"comment"`
	}
	err := json.Unmarshal([]byte(eventValue), &c)
	if err != nil {
		return nil, nil, err, false
	}

	kudo, err := s.CheckinService.KudoUpdate(ctx, teamID, userID, c.KudoID, c.TargetUserID, c.KudosDate, c.Comment)
	if err != nil {
		return nil, nil, err, false
	}

	msg := wshub.CreateSocketEvent("kudo_updated", "", "")

	return kudo, msg, nil, false
}

// KudoDelete deletes a team kudo.
func (s *Service) KudoDelete(ctx context.Context, teamID string, userID string, eventValue string) (any, []byte, error, bool) {
	var c struct {
		KudoID string `json:"kudoId"`
	}
	err := json.Unmarshal([]byte(eventValue), &c)
	if err != nil {
		return nil, nil, err, false
	}

	err = s.CheckinService.KudoDelete(ctx, teamID, userID, c.KudoID)
	if err != nil {
		return nil, nil, err, false
	}

	msg := wshub.CreateSocketEvent("kudo_deleted", "", "")

	return nil, msg, nil, false
}
