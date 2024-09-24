package wshub

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
)

// TestHandleSocketClose tests the handling of socket closure
func TestHandleSocketClose(t *testing.T) {
	hub := NewHub(otelzap.New(zap.NewNop()), Config{}, nil, nil, nil, nil)

	// Create channels for synchronization
	serverReady := make(chan struct{})
	clientClosed := make(chan struct{})

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		upgrader := websocket.Upgrader{}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Fatalf("Failed to upgrade connection: %v", err)
		}
		defer conn.Close()

		// Signal that the server is ready
		close(serverReady)

		// Wait for the client to close the connection
		<-clientClosed

		hub.HandleSocketClose(context.Background(), conn, websocket.CloseNormalClosure, "test close")
	}))
	defer server.Close()

	url := "ws" + strings.TrimPrefix(server.URL, "http")
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Fatalf("Failed to connect to WebSocket server: %v", err)
	}

	// Wait for the server to be ready
	<-serverReady

	// Close the connection from the client side
	err = conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "client closing"))
	if err != nil {
		t.Fatalf("Failed to send close message: %v", err)
	}

	// Signal that the client has closed the connection
	close(clientClosed)

	// Wait a short time for the server to process the close
	time.Sleep(100 * time.Millisecond)

	// Attempt to read from the closed connection
	_, _, err = conn.ReadMessage()

	// Check for the expected error
	if err == nil {
		t.Fatal("Expected an error, but got nil")
	}

	// The error message might vary, so we'll check for both possible outcomes
	expectedErrors := []string{
		"websocket: close 1000 (normal)",
		"use of closed network connection",
	}

	errorMatched := false
	for _, expectedError := range expectedErrors {
		if strings.Contains(err.Error(), expectedError) {
			errorMatched = true
			break
		}
	}

	assert.True(t, errorMatched, "Error should be one of the expected closure messages")
}

// TestWebSocketHandler tests the WebSocket handler
func TestWebSocketHandler(t *testing.T) {
	hub := NewHub(otelzap.New(zap.NewNop()), Config{}, nil, nil, nil, nil)

	authFunc := func(w http.ResponseWriter, r *http.Request, c *Connection, roomID string) *AuthError {
		return nil
	}

	handler := hub.WebSocketHandler("roomID", authFunc)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)
	}))
	defer server.Close()

	url := "ws" + strings.TrimPrefix(server.URL, "http")
	_, _, err := websocket.DefaultDialer.Dial(url, nil)
	assert.NoError(t, err)
}
