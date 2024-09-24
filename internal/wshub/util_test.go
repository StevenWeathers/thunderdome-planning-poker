package wshub

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCreateSocketEvent tests the creation of socket events
func TestCreateSocketEvent(t *testing.T) {
	event := CreateSocketEvent("test_type", "test_value", "user1")

	var socketEvent SocketEvent
	err := json.Unmarshal(event, &socketEvent)
	assert.NoError(t, err)
	assert.Equal(t, "test_type", socketEvent.Type)
	assert.Equal(t, "test_value", socketEvent.Value)
	assert.Equal(t, "user1", socketEvent.UserID)
}

// TestEqualASCIIFold tests the ASCII folding comparison
func TestEqualASCIIFold(t *testing.T) {
	assert.True(t, equalASCIIFold("test", "TEST"))
	assert.True(t, equalASCIIFold("Test", "tEST"))
	assert.False(t, equalASCIIFold("test", "test1"))
}

// TestCheckOrigin tests the origin checking function
func TestCheckOrigin(t *testing.T) {
	r, _ := http.NewRequest("GET", "http://example.com", nil)
	r.Header.Set("Origin", "http://example.com")
	assert.True(t, checkOrigin(r, "example.com", "ws"))

	r.Header.Set("Origin", "http://ws.example.com")
	assert.True(t, checkOrigin(r, "example.com", "ws"))

	r.Header.Set("Origin", "http://other.com")
	assert.False(t, checkOrigin(r, "example.com", "ws"))
}
