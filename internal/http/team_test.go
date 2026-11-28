package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
)

type teamUpdateUserStub struct {
	TeamDataSvc
	called bool
	title  *string
	notes  *string
}

func (s *teamUpdateUserStub) TeamUpdateUser(_ context.Context, teamID string, _ string, _ string, title *string, notes *string) (string, error) {
	s.called = true
	s.title = title
	s.notes = notes

	return teamID, nil
}

func TestHandleTeamUpdateUser(t *testing.T) {
	validate = validator.New()

	strPtr := func(v string) *string { return &v }

	tests := []struct {
		name       string
		body       string
		wantStatus int
		wantTitle  *string
		wantNotes  *string
	}{
		{
			name:       "role only leaves title and notes unchanged",
			body:       `{"role":"MEMBER"}`,
			wantStatus: http.StatusOK,
		},
		{
			name:       "title and notes are trimmed",
			body:       `{"role":"MEMBER","title":"  Backend Developer ","notes":" On call Fridays\n"}`,
			wantStatus: http.StatusOK,
			wantTitle:  strPtr("Backend Developer"),
			wantNotes:  strPtr("On call Fridays"),
		},
		{
			name:       "empty values clear title and notes",
			body:       `{"role":"ADMIN","title":"","notes":""}`,
			wantStatus: http.StatusOK,
			wantTitle:  strPtr(""),
			wantNotes:  strPtr(""),
		},
		{
			name:       "max length counts characters not bytes",
			body:       `{"role":"MEMBER","title":"` + strings.Repeat("é", 128) + `"}`,
			wantStatus: http.StatusOK,
			wantTitle:  strPtr(strings.Repeat("é", 128)),
		},
		{
			name:       "title too long",
			body:       `{"role":"MEMBER","title":"` + strings.Repeat("a", 129) + `"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "notes too long",
			body:       `{"role":"MEMBER","notes":"` + strings.Repeat("a", 1001) + `"}`,
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stub := &teamUpdateUserStub{}
			s := &Service{TeamDataSvc: stub, Logger: otelzap.New(zap.NewNop())}

			req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(tt.body))
			req.SetPathValue("teamId", "6f1b0c2e-2d6a-4f0b-9a55-6c2f1a7d9b10")
			req.SetPathValue("userId", "0b8e4c7a-91f3-4a6e-8d2b-5e7c3f1a2b44")
			req = req.WithContext(context.WithValue(req.Context(), contextKeyUserID, "session-user"))
			w := httptest.NewRecorder()

			s.handleTeamUpdateUser()(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
			assert.Equal(t, tt.wantStatus == http.StatusOK, stub.called)
			assert.Equal(t, tt.wantTitle, stub.title)
			assert.Equal(t, tt.wantNotes, stub.notes)
		})
	}
}

type orgUpdateUserStub struct {
	OrganizationDataSvc
	called bool
}

func (s *orgUpdateUserStub) OrganizationUpdateUser(_ context.Context, orgID string, _ string, _ string, _ *string, _ *string) (string, error) {
	s.called = true
	return orgID, nil
}

func (s *orgUpdateUserStub) DepartmentUpdateUser(_ context.Context, departmentID string, _ string, _ string, _ *string, _ *string) (string, error) {
	s.called = true
	return departmentID, nil
}

func TestOrganizationAndDepartmentUpdateUserValidation(t *testing.T) {
	validate = validator.New()

	handlers := map[string]func(*Service) http.HandlerFunc{
		"organization": (*Service).handleOrganizationUpdateUser,
		"department":   (*Service).handleDepartmentUpdateUser,
	}
	bodies := map[string]int{
		`{"role":"MEMBER","title":"Designer"}`:                          http.StatusOK,
		`{"role":"MEMBER","title":"` + strings.Repeat("a", 129) + `"}`:  http.StatusBadRequest,
		`{"role":"MEMBER","notes":"` + strings.Repeat("a", 1001) + `"}`: http.StatusBadRequest,
		`{"role":"OWNER"}`: http.StatusBadRequest,
	}

	for name, handler := range handlers {
		for body, wantStatus := range bodies {
			stub := &orgUpdateUserStub{}
			s := &Service{
				Config:              &Config{OrganizationsEnabled: true},
				OrganizationDataSvc: stub,
				Logger:              otelzap.New(zap.NewNop()),
			}

			req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(body))
			req.SetPathValue("orgId", "6f1b0c2e-2d6a-4f0b-9a55-6c2f1a7d9b10")
			req.SetPathValue("departmentId", "1c9d5f3a-7b2e-4d8c-a6f1-2e4b8c0d9a77")
			req.SetPathValue("userId", "0b8e4c7a-91f3-4a6e-8d2b-5e7c3f1a2b44")
			req = req.WithContext(context.WithValue(req.Context(), contextKeyUserID, "session-user"))
			w := httptest.NewRecorder()

			handler(s)(w, req)

			assert.Equal(t, wantStatus, w.Code, "%s %s", name, body)
			assert.Equal(t, wantStatus == http.StatusOK, stub.called, "%s %s", name, body)
		}
	}
}
