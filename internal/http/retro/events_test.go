package retro

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type editRetroStub struct {
	RetroDataSvc
	facilitatorCode string
}

func (s *editRetroStub) EditRetro(_ string, _ string, _ string, facilitatorCode string, _ int, _ string, _ bool, _ bool, _ int) error {
	s.facilitatorCode = facilitatorCode
	return nil
}

func TestEditRetroDoesNotBroadcastFacilitatorCode(t *testing.T) {
	stub := &editRetroStub{}
	s := &Service{RetroService: stub}

	_, msg, err, _ := s.EditRetro(context.Background(), "retro", "user",
		`{"retroName":"Retro","joinCode":"join","facilitatorCode":"top-secret"}`)

	require.NoError(t, err)
	assert.Equal(t, "top-secret", stub.facilitatorCode, "the code is still saved")
	assert.NotContains(t, string(msg), "top-secret")
	assert.Contains(t, string(msg), "join", "participants still receive the join code")
}
