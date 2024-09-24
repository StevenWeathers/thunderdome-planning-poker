package wshub

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
)

// TestNewHub tests the creation of a new Hub
func TestNewHub(t *testing.T) {
	config := Config{
		WriteWaitSec:       10,
		PongWaitSec:        60,
		PingPeriodSec:      54,
		AppDomain:          "example.com",
		WebsocketSubdomain: "ws",
	}
	eventHandlers := make(map[string]func(context.Context, string, string, string) ([]byte, error, bool))
	facilitatorOnlyOperations := make(map[string]struct{})
	confirmFacilitator := func(roomId string, userId string) error { return nil }
	retreatUser := func(roomId string, userId string) string { return "" }

	hub := NewHub(otelzap.New(zap.NewNop()), config, eventHandlers, facilitatorOnlyOperations, confirmFacilitator, retreatUser)

	assert.NotNil(t, hub)
	assert.Equal(t, &config, hub.config)
	//assert.Equal(t, logger, hub.logger)
	assert.NotNil(t, hub.rooms)
	assert.NotNil(t, hub.broadcast)
	assert.NotNil(t, hub.register)
	assert.NotNil(t, hub.unregister)
}

// TestCreateWebsocketUpgrader tests the creation of a websocket upgrader
func TestCreateWebsocketUpgrader(t *testing.T) {
	hub := NewHub(otelzap.New(zap.NewNop()), Config{AppDomain: "example.com", WebsocketSubdomain: "ws"}, nil, nil, nil, nil)
	upgrader := hub.CreateWebsocketUpgrader()

	assert.NotNil(t, upgrader)
	assert.Equal(t, 1024, upgrader.ReadBufferSize)
	assert.Equal(t, 1024, upgrader.WriteBufferSize)
	assert.NotNil(t, upgrader.CheckOrigin)
}
