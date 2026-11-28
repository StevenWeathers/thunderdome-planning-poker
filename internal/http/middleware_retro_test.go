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

type retroAccessStub struct {
	RetroDataSvc
	allowed bool
	err     error
}

func (s *retroAccessStub) CanUserAccessRetro(context.Context, string, string) (bool, error) {
	return s.allowed, s.err
}

func TestRetroUserOnly(t *testing.T) {
	validate = validator.New()

	tests := []struct {
		name       string
		userType   string
		allowed    bool
		err        error
		wantStatus int
	}{
		{name: "allowed user", allowed: true, wantStatus: http.StatusOK},
		{name: "user without access", wantStatus: http.StatusForbidden},
		{name: "missing retro", err: fmt.Errorf("lookup: %w", sql.ErrNoRows), wantStatus: http.StatusNotFound},
		{name: "lookup failure", err: errors.New("boom"), wantStatus: http.StatusInternalServerError},
		{name: "site admin", userType: thunderdome.AdminUserType, wantStatus: http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				RetroDataSvc: &retroAccessStub{allowed: tt.allowed, err: tt.err},
				Logger:       otelzap.New(zap.NewNop()),
			}
			handler := s.retroUserOnly(func(http.ResponseWriter, *http.Request) {})

			ctx := context.WithValue(context.Background(), contextKeyUserID, "e0a1f5c2-3b7d-4a4e-9d0e-6f1a2b3c4d5e")
			ctx = context.WithValue(ctx, contextKeyUserType, tt.userType)
			req := httptest.NewRequestWithContext(ctx, http.MethodPut, "/", nil)
			req.SetPathValue("retroId", "0f8fad5b-d9cb-469f-a165-70867728950e")
			rec := httptest.NewRecorder()

			handler(rec, req)

			assert.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}
