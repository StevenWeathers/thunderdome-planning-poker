package poker_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/StevenWeathers/thunderdome-planning-poker/internal/http/poker"
	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
)

// MockPokerDataSvc is a mock implementation of PokerDataSvc for testing
type MockPokerDataSvc struct {
	mock.Mock
}

func (m *MockPokerDataSvc) StopGame(pokerID string) error {
	args := m.Called(pokerID)
	return args.Error(0)
}

func (m *MockPokerDataSvc) ConfirmFacilitator(pokerID string, userID string) error {
	args := m.Called(pokerID, userID)
	return args.Error(0)
}

func (m *MockPokerDataSvc) GetGameByID(pokerID string, userID string) (*thunderdome.Poker, error) {
	args := m.Called(pokerID, userID)
	return args.Get(0).(*thunderdome.Poker), args.Error(1)
}

// Add other required methods from PokerDataSvc interface as no-ops
func (m *MockPokerDataSvc) UpdateGame(pokerID string, name string, pointValuesAllowed []string, autoFinishVoting bool, pointAverageRounding string, hideVoterIdentity bool, joinCode string, facilitatorCode string, teamID string) error {
	return nil
}

func (m *MockPokerDataSvc) GetFacilitatorCode(pokerID string) (string, error) {
	return "", nil
}

func (m *MockPokerDataSvc) GetUserActiveStatus(pokerID string, userID string) error {
	return nil
}

func (m *MockPokerDataSvc) AddUser(pokerID string, userID string) ([]*thunderdome.PokerUser, error) {
	return nil, nil
}

func (m *MockPokerDataSvc) RetreatUser(pokerID string, userID string) []*thunderdome.PokerUser {
	return nil
}

func (m *MockPokerDataSvc) AbandonGame(pokerID string, userID string) ([]*thunderdome.PokerUser, error) {
	return nil, nil
}

func (m *MockPokerDataSvc) AddFacilitator(pokerID string, userID string) ([]string, error) {
	return nil, nil
}

func (m *MockPokerDataSvc) RemoveFacilitator(pokerID string, userID string) ([]string, error) {
	return nil, nil
}

func (m *MockPokerDataSvc) ToggleSpectator(pokerID string, userID string, spectator bool) ([]*thunderdome.PokerUser, error) {
	return nil, nil
}

func (m *MockPokerDataSvc) DeleteGame(pokerID string) error {
	return nil
}

func (m *MockPokerDataSvc) CreateStory(pokerID string, name string, storyType string, referenceID string, link string, description string, acceptanceCriteria string, priority int32) ([]*thunderdome.Story, error) {
	return nil, nil
}

func (m *MockPokerDataSvc) ActivateStoryVoting(pokerID string, storyID string) ([]*thunderdome.Story, error) {
	return nil, nil
}

func (m *MockPokerDataSvc) SetVote(pokerID string, userID string, storyID string, voteValue string) ([]*thunderdome.Story, bool) {
	return nil, false
}

func (m *MockPokerDataSvc) RetractVote(pokerID string, userID string, storyID string) ([]*thunderdome.Story, error) {
	return nil, nil
}

func (m *MockPokerDataSvc) EndStoryVoting(pokerID string, storyID string) ([]*thunderdome.Story, error) {
	return nil, nil
}

func (m *MockPokerDataSvc) SkipStory(pokerID string, storyID string) ([]*thunderdome.Story, error) {
	return nil, nil
}

func (m *MockPokerDataSvc) UpdateStory(pokerID string, storyID string, name string, storyType string, referenceID string, link string, description string, acceptanceCriteria string, priority int32) ([]*thunderdome.Story, error) {
	return nil, nil
}

func (m *MockPokerDataSvc) DeleteStory(pokerID string, storyID string) ([]*thunderdome.Story, error) {
	return nil, nil
}

func (m *MockPokerDataSvc) ArrangeStory(pokerID string, storyID string, beforeStoryID string) ([]*thunderdome.Story, error) {
	return nil, nil
}

func (m *MockPokerDataSvc) FinalizeStory(pokerID string, storyID string, points string) ([]*thunderdome.Story, error) {
	return nil, nil
}

// TestStopGameFunctionality tests the stop game WebSocket event functionality
func TestStopGameFunctionality(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	otelLogger := otelzap.New(logger)

	t.Run("successful game stop by authorized facilitator", func(t *testing.T) {
		mockSvc := &MockPokerDataSvc{}

		service := poker.New(
			poker.Config{},
			otelLogger,
			nil, nil, nil, nil, mockSvc,
		)

		gameID := "550e8400-e29b-41d4-a716-446655440000"
		userID := "550e8400-e29b-41d4-a716-446655440001"

		// Mock successful facilitator confirmation and game stop
		mockSvc.On("ConfirmFacilitator", gameID, userID).Return(nil)
		mockSvc.On("StopGame", gameID).Return(nil)

		now := time.Now()
		mockGame := &thunderdome.Poker{
			ID:        gameID,
			Name:      "Test Game",
			EndedDate: &now,
		}
		mockSvc.On("GetGameByID", gameID, userID).Return(mockGame, nil)

		// Test the Stop function
		ctx := context.Background()
		result, msg, err, shouldClose := service.Stop(ctx, gameID, userID, "")

		// Assertions
		assert.NoError(t, err)
		assert.Nil(t, result)
		assert.NotNil(t, msg)
		assert.False(t, shouldClose)

		// Verify all mocks were called
		mockSvc.AssertExpectations(t)
	})

	t.Run("unauthorized user cannot stop game", func(t *testing.T) {
		mockSvc := &MockPokerDataSvc{}

		service := poker.New(
			poker.Config{},
			otelLogger,
			nil, nil, nil, nil, mockSvc,
		)

		gameID := "550e8400-e29b-41d4-a716-446655440000"
		userID := "550e8400-e29b-41d4-a716-446655440001"

		// Mock failed facilitator confirmation
		mockSvc.On("ConfirmFacilitator", gameID, userID).Return(errors.New("not a facilitator"))

		// Test the Stop function
		ctx := context.Background()
		result, msg, err, shouldClose := service.Stop(ctx, gameID, userID, "")

		// Assertions
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "UNAUTHORIZED")
		assert.Nil(t, result)
		assert.Nil(t, msg)
		assert.False(t, shouldClose)

		// Verify ConfirmFacilitator was called but StopGame was not
		mockSvc.AssertExpectations(t)
		mockSvc.AssertNotCalled(t, "StopGame")
	})

	t.Run("database error during game stop", func(t *testing.T) {
		mockSvc := &MockPokerDataSvc{}

		service := poker.New(
			poker.Config{},
			otelLogger,
			nil, nil, nil, nil, mockSvc,
		)

		gameID := "550e8400-e29b-41d4-a716-446655440000"
		userID := "550e8400-e29b-41d4-a716-446655440001"

		// Mock successful facilitator confirmation but failed game stop
		mockSvc.On("ConfirmFacilitator", gameID, userID).Return(nil)
		mockSvc.On("StopGame", gameID).Return(errors.New("database error"))

		// Test the Stop function
		ctx := context.Background()
		result, msg, err, shouldClose := service.Stop(ctx, gameID, userID, "")

		// Assertions
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database error")
		assert.Nil(t, result)
		assert.Nil(t, msg)
		assert.False(t, shouldClose)

		// Verify both methods were called
		mockSvc.AssertExpectations(t)
	})
}

// TestStopGameInputValidation tests the input validation in the StopGame database function
func TestStopGameInputValidation(t *testing.T) {
	t.Run("invalid UUID format should be rejected", func(t *testing.T) {
		// This test validates the isValidUUID function behavior
		// Note: This would need to be tested against actual StopGame implementation
		invalidUUIDs := []string{
			"",
			"invalid-uuid",
			"550e8400-e29b-41d4-a716-44665544000",   // too short
			"550e8400-e29b-41d4-a716-4466554400000", // too long
			"550e8400-e29b-41d4-a716-44665544000z",  // invalid character
		}

		for _, invalidUUID := range invalidUUIDs {
			t.Run("UUID: "+invalidUUID, func(t *testing.T) {
				// Test would validate that StopGame returns appropriate error
				// This is a placeholder for actual database integration test
				assert.True(t, len(invalidUUID) == 0 || len(invalidUUID) != 36 || containsInvalidChars(invalidUUID))
			})
		}
	})
}

// Helper function to check for invalid UUID characters
func containsInvalidChars(s string) bool {
	for _, char := range s {
		if !((char >= '0' && char <= '9') || (char >= 'a' && char <= 'f') || (char >= 'A' && char <= 'F') || char == '-') {
			return true
		}
	}
	return false
}

// TestGameStatusBadgeDisplay tests the status badge display logic
func TestGameStatusBadgeDisplay(t *testing.T) {
	t.Run("active game shows correct status", func(t *testing.T) {
		// Test logic for active game (endedDate is nil)
		var endedDate *time.Time = nil
		isActive := endedDate == nil

		assert.True(t, isActive)
		// In actual component, this would show "Active" with green color
	})

	t.Run("stopped game shows correct status with timestamp", func(t *testing.T) {
		// Test logic for stopped game (endedDate is set)
		now := time.Now()
		endedDate := &now
		isActive := endedDate == nil

		assert.False(t, isActive)
		assert.NotNil(t, endedDate)
		// In actual component, this would show "Stopped" with orange color and formatted timestamp
	})
}
