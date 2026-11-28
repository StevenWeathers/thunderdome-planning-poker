package http

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
)

type stubStoryboardAccessSvc struct {
	StoryboardDataSvc
	allowed bool
	err     error
}

func (s stubStoryboardAccessSvc) CanUserAccessStoryboard(_ context.Context, _ string, _ string) (bool, error) {
	return s.allowed, s.err
}

func TestStoryboardUserOnly(t *testing.T) {
	validate = validator.New()
	storyboardID := "11111111-1111-1111-1111-111111111111"

	tests := []struct {
		name         string
		storyboardID string
		userType     string
		svc          stubStoryboardAccessSvc
		wantStatus   int
	}{
		{"user with access", storyboardID, thunderdome.RegisteredUserType, stubStoryboardAccessSvc{allowed: true}, http.StatusOK},
		{"user without access", storyboardID, thunderdome.RegisteredUserType, stubStoryboardAccessSvc{}, http.StatusForbidden},
		{"admin without access", storyboardID, thunderdome.AdminUserType, stubStoryboardAccessSvc{}, http.StatusOK},
		{"missing storyboard", storyboardID, thunderdome.RegisteredUserType, stubStoryboardAccessSvc{err: fmt.Errorf("query: %w", sql.ErrNoRows)}, http.StatusNotFound},
		{"query failure", storyboardID, thunderdome.RegisteredUserType, stubStoryboardAccessSvc{err: errors.New("boom")}, http.StatusInternalServerError},
		{"invalid storyboard id", "nope", thunderdome.RegisteredUserType, stubStoryboardAccessSvc{allowed: true}, http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &Service{StoryboardDataSvc: tt.svc, Logger: otelzap.New(zap.NewNop())}
			called := false
			handler := service.storyboardUserOnly(func(w http.ResponseWriter, _ *http.Request) {
				called = true
				w.WriteHeader(http.StatusOK)
			})

			ctx := context.WithValue(context.Background(), contextKeyUserID, "user-1")
			ctx = context.WithValue(ctx, contextKeyUserType, tt.userType)
			req := httptest.NewRequest(http.MethodDelete, "/api/storyboards/x/stories/y", nil).WithContext(ctx)
			req.SetPathValue("storyboardId", tt.storyboardID)
			rr := httptest.NewRecorder()

			handler(rr, req)

			assert.Equal(t, tt.wantStatus, rr.Code)
			assert.Equal(t, tt.wantStatus == http.StatusOK, called)
		})
	}
}
