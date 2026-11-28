package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
)

func TestImmutableCache(t *testing.T) {
	handler := immutableCache(http.FileServer(http.FS(fstest.MapFS{
		"assets/index-abc123.js": {Data: []byte("export {}")},
	})))

	tests := []struct {
		name         string
		path         string
		status       int
		cacheControl string
	}{
		{"existing asset", "/assets/index-abc123.js", http.StatusOK, "public, max-age=31536000, immutable"},
		{"missing asset", "/assets/index-old456.js", http.StatusNotFound, ""},
		{"directory listing", "/assets/", http.StatusOK, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, tt.path, nil))

			assert.Equal(t, tt.status, rr.Code)
			assert.Equal(t, tt.cacheControl, rr.Header().Get("Cache-Control"))
		})
	}
}

type alertsStub struct {
	AlertDataSvc
}

func (alertsStub) GetActiveAlerts(context.Context) []any {
	return []any{}
}

func TestHandleIndexCacheControl(t *testing.T) {
	s := &Service{Config: &Config{}, AlertDataSvc: alertsStub{}}
	handler := s.handleIndex(fstest.MapFS{"index.html": {Data: []byte("<html></html>")}}, thunderdome.UIConfig{})
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/", nil))

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "no-cache", rr.Header().Get("Cache-Control"))
}
