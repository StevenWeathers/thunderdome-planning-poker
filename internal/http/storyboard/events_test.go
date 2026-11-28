package storyboard

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/StevenWeathers/thunderdome-planning-poker/internal/wshub"
	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type moveGoalStub struct {
	StoryboardDataSvc
	moveErr error
	moved   []string
}

func (s *moveGoalStub) MoveStoryboardGoal(storyboardID string, userID string, goalID string, placeBeforeID string) error {
	s.moved = []string{storyboardID, userID, goalID, placeBeforeID}
	return s.moveErr
}

func (s *moveGoalStub) GetStoryboardGoals(storyboardID string) []*thunderdome.StoryboardGoal {
	return []*thunderdome.StoryboardGoal{{ID: "goal-2"}, {ID: "goal-1"}}
}

func TestMoveGoal(t *testing.T) {
	input := `{"goalId":"goal-2","placeBefore":"goal-1"}`
	moved := []string{"board", "user", "goal-2", "goal-1"}

	tests := []struct {
		name      string
		input     string
		moveErr   error
		wantMoved []string
		wantErr   bool
	}{
		{name: "moves the goal", input: input, wantMoved: moved},
		{name: "invalid input", input: "nope", wantErr: true},
		{name: "data service error", input: input, moveErr: errors.New("move failed"), wantMoved: moved, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stub := &moveGoalStub{moveErr: tt.moveErr}
			s := &Service{StoryboardService: stub}

			_, msg, err, _ := s.MoveGoal(context.Background(), "board", "user", tt.input)

			assert.Equal(t, tt.wantMoved, stub.moved)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, msg)
				return
			}
			require.NoError(t, err)
			goals, err := json.Marshal(stub.GetStoryboardGoals("board"))
			require.NoError(t, err)
			assert.Equal(t, wshub.CreateSocketEvent("goal_moved", string(goals), ""), msg)
		})
	}
}

type editStoryboardStub struct {
	StoryboardDataSvc
	scalesErr error
	edited    bool
	editedID  *string

	facilitatorCode string
}

func (s *editStoryboardStub) GetStoryboardEstimationScales(_ context.Context, _ string, _ bool) ([]*thunderdome.EstimationScale, error) {
	return []*thunderdome.EstimationScale{{ID: "scale-1", Name: "T-Shirt", Values: []string{"S", "M"}}}, s.scalesErr
}

func (s *editStoryboardStub) EditStoryboard(_ string, _ string, _ string, facilitatorCode string, estimationScaleID *string) error {
	s.facilitatorCode = facilitatorCode
	s.edited = true
	s.editedID = estimationScaleID
	return nil
}

func TestEditStoryboardEstimationScale(t *testing.T) {
	scaleID, empty := "scale-1", ""

	tests := []struct {
		name      string
		input     string
		scalesErr error
		wantID    *string
		wantScale bool
		wantErr   bool
	}{
		{name: "omitted keeps the scale", input: `{"storyboardName":"Board"}`},
		{name: "empty removes the scale", input: `{"storyboardName":"Board","estimationScaleId":""}`, wantID: &empty},
		{name: "assigns an available scale", input: `{"storyboardName":"Board","estimationScaleId":"scale-1"}`, wantID: &scaleID, wantScale: true},
		{name: "ignores a client supplied scale", input: `{"storyboardName":"Board","estimationScale":{"id":"forged"}}`},
		{name: "rejects an unavailable scale", input: `{"storyboardName":"Board","estimationScaleId":"other"}`, wantErr: true},
		{name: "scale lookup failure", input: `{"storyboardName":"Board","estimationScaleId":"scale-1"}`, scalesErr: errors.New("boom"), wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stub := &editStoryboardStub{scalesErr: tt.scalesErr}
			s := &Service{StoryboardService: stub}

			_, msg, err, _ := s.EditStoryboard(context.Background(), "board", "user", tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				assert.False(t, stub.edited)
				assert.Nil(t, msg)
				return
			}
			require.NoError(t, err)
			assert.True(t, stub.edited)
			assert.Equal(t, tt.wantID, stub.editedID)

			var event struct {
				Value string `json:"value"`
			}
			require.NoError(t, json.Unmarshal(msg, &event))
			assert.Equal(t, tt.wantScale, strings.Contains(event.Value, `"values":["S","M"]`))
			assert.NotContains(t, event.Value, "forged")
		})
	}
}

func TestEditStoryboardDoesNotBroadcastFacilitatorCode(t *testing.T) {
	stub := &editStoryboardStub{}
	s := &Service{StoryboardService: stub}

	_, msg, err, _ := s.EditStoryboard(context.Background(), "board", "user",
		`{"storyboardName":"Board","joinCode":"join","facilitatorCode":"top-secret"}`)

	require.NoError(t, err)
	assert.Equal(t, "top-secret", stub.facilitatorCode, "the code is still saved")
	assert.NotContains(t, string(msg), "top-secret")
	assert.Contains(t, string(msg), "join", "participants still receive the join code")
}
