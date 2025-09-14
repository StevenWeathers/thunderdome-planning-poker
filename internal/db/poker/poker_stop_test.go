package poker_test

import (
	"testing"
	"time"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
	"github.com/stretchr/testify/assert"
)

func TestPokerStructEndedDate(t *testing.T) {
	// Test that the Poker struct has the EndedDate field
	now := time.Now()
	poker := &thunderdome.Poker{
		ID:        "test-id",
		Name:      "Test Poker",
		EndedDate: &now,
	}

	assert.NotNil(t, poker.EndedDate)
	assert.Equal(t, now, *poker.EndedDate)

	// Test with nil EndedDate (game not stopped)
	pokerActive := &thunderdome.Poker{
		ID:        "active-id",
		Name:      "Active Poker",
		EndedDate: nil,
	}

	assert.Nil(t, pokerActive.EndedDate)
}
