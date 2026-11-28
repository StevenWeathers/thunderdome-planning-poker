package storyboard

import (
	"context"
	"encoding/json"
	"errors"
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
