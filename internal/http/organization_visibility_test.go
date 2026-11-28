package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
)

type visibilityStub struct {
	OrganizationDataSvc
	called string
	userID string
}

func (s *visibilityStub) OrganizationDepartmentList(context.Context, string, int, int) []*thunderdome.Department {
	s.called = "all"
	return []*thunderdome.Department{}
}

func (s *visibilityStub) OrganizationDepartmentListByUser(_ context.Context, _ string, userID string, _ int, _ int) []*thunderdome.Department {
	s.called, s.userID = "byUser", userID
	return []*thunderdome.Department{}
}

func (s *visibilityStub) OrganizationTeamList(context.Context, string, int, int) []*thunderdome.Team {
	s.called = "all"
	return []*thunderdome.Team{}
}

func (s *visibilityStub) OrganizationTeamListByUser(_ context.Context, _ string, userID string, _ int, _ int) []*thunderdome.Team {
	s.called, s.userID = "byUser", userID
	return []*thunderdome.Team{}
}

func (s *visibilityStub) DepartmentTeamList(context.Context, string, int, int) []*thunderdome.Team {
	s.called = "all"
	return []*thunderdome.Team{}
}

func (s *visibilityStub) DepartmentTeamListByUser(_ context.Context, _ string, userID string, _ int, _ int) []*thunderdome.Team {
	s.called, s.userID = "byUser", userID
	return []*thunderdome.Team{}
}

func TestOrganizationListVisibility(t *testing.T) {
	validate = validator.New()

	const sessionUserID = "e0a1f5c2-3b7d-4a4e-9d0e-6f1a2b3c4d5e"
	departments := (*Service).handleGetOrganizationDepartments
	orgTeams := (*Service).handleGetOrganizationTeams
	deptTeams := (*Service).handleGetDepartmentTeams

	tests := []struct {
		name     string
		handler  func(*Service) http.HandlerFunc
		orgRole  string
		deptRole string
		want     string
	}{
		{name: "org admin gets all departments", handler: departments, orgRole: "ADMIN", want: "all"},
		{name: "org member gets own departments", handler: departments, orgRole: "MEMBER", want: "byUser"},
		{name: "org admin gets all org teams", handler: orgTeams, orgRole: "ADMIN", want: "all"},
		{name: "org member gets own org teams", handler: orgTeams, orgRole: "MEMBER", want: "byUser"},
		{name: "org admin gets all department teams", handler: deptTeams, orgRole: "ADMIN", want: "all"},
		{name: "department admin gets all department teams", handler: deptTeams, orgRole: "MEMBER", deptRole: "ADMIN", want: "all"},
		{name: "department member gets own department teams", handler: deptTeams, orgRole: "MEMBER", deptRole: "MEMBER", want: "byUser"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stub := &visibilityStub{}
			s := &Service{Config: &Config{OrganizationsEnabled: true}, OrganizationDataSvc: stub}

			ctx := context.WithValue(context.Background(), contextKeyUserID, sessionUserID)
			ctx = context.WithValue(ctx, contextKeyOrgRole, tt.orgRole)
			ctx = context.WithValue(ctx, contextKeyDepartmentRole, tt.deptRole)
			req := httptest.NewRequestWithContext(ctx, http.MethodGet, "/", nil)
			req.SetPathValue("orgId", "0f8fad5b-d9cb-469f-a165-70867728950e")
			req.SetPathValue("departmentId", "7c9e6679-7425-40de-944b-e07fc1f90ae7")
			rec := httptest.NewRecorder()

			tt.handler(s)(rec, req)

			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, tt.want, stub.called)
			if tt.want == "byUser" {
				assert.Equal(t, sessionUserID, stub.userID)
			}
		})
	}
}
