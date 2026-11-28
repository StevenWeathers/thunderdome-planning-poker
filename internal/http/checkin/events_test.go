package checkin

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
)

type actorStub struct {
	CheckinDataSvc
	actor string
}

func (s *actorStub) CheckinCreate(_ context.Context, _ string, userID string, _, _, _, _, _ string, _ bool) error {
	s.actor = userID
	return nil
}

func (s *actorStub) CheckinUpdate(_ context.Context, _ string, userID string, _, _, _, _, _ string, _ bool) error {
	s.actor = userID
	return nil
}

func (s *actorStub) CheckinDelete(_ context.Context, _ string, userID string, _ string) error {
	s.actor = userID
	return nil
}

func (s *actorStub) CheckinComment(_ context.Context, _ string, _ string, userID string, _ string) error {
	s.actor = userID
	return nil
}

func (s *actorStub) CheckinCommentEdit(_ context.Context, _ string, userID string, _ string, _ string) error {
	s.actor = userID
	return nil
}

func (s *actorStub) CheckinCommentDelete(_ context.Context, _ string, userID string, _ string) error {
	s.actor = userID
	return nil
}

func (s *actorStub) KudoUpdate(_ context.Context, _ string, userID string, _, _, _, _ string) (*thunderdome.TeamKudo, error) {
	s.actor = userID
	return nil, nil
}

func (s *actorStub) KudoDelete(_ context.Context, _ string, userID string, _ string) error {
	s.actor = userID
	return nil
}

func TestEventsActAsTheSender(t *testing.T) {
	stub := &actorStub{}
	s := &Service{CheckinService: stub}

	events := map[string]func(context.Context, string, string, string) (any, []byte, error, bool){
		"checkin_create": s.CheckinCreate,
		"checkin_update": s.CheckinUpdate,
		"checkin_delete": s.CheckinDelete,
		"comment_create": s.CommentCreate,
		"comment_update": s.CommentUpdate,
		"comment_delete": s.CommentDelete,
		"kudo_update":    s.KudoUpdate,
		"kudo_delete":    s.KudoDelete,
	}

	for name, event := range events {
		t.Run(name, func(t *testing.T) {
			stub.actor = ""

			_, _, err, _ := event(context.Background(), "team", "sender", `{"userId":"someone-else"}`)

			assert.NoError(t, err)
			assert.Equal(t, "sender", stub.actor)
		})
	}
}
