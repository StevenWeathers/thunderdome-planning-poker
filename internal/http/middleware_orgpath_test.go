package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
)

func TestOrgPathVerification(t *testing.T) {
	validate = validator.New()

	const (
		userID = "e0a1f5c2-3b7d-4a4e-9d0e-6f1a2b3c4d5e"
		orgID  = "0f8fad5b-d9cb-469f-a165-70867728950e"
		deptID = "7c9e6679-7425-40de-944b-e07fc1f90ae7"
		teamID = "6f1b0c2e-2d6a-4f0b-9a55-6c2f1a7d9b10"
	)
	orgAdmin := (*Service).orgAdminOnly
	deptUser := (*Service).departmentUserOnly
	deptAdmin := (*Service).departmentAdminOnly

	tests := []struct {
		name         string
		middleware   func(*Service, http.HandlerFunc) http.HandlerFunc
		userType     string
		departmentID string
		teamID       string
		exists       bool
		wantStatus   int
	}{
		{name: "org admin, foreign department", middleware: orgAdmin, departmentID: deptID, wantStatus: http.StatusNotFound},
		{name: "org admin, foreign team", middleware: orgAdmin, teamID: teamID, wantStatus: http.StatusNotFound},
		{name: "org admin, own team", middleware: orgAdmin, teamID: teamID, exists: true, wantStatus: http.StatusOK},
		{name: "department user, foreign department", middleware: deptUser, departmentID: deptID, wantStatus: http.StatusNotFound},
		{name: "department admin, foreign department", middleware: deptAdmin, departmentID: deptID, wantStatus: http.StatusNotFound},
		{name: "department admin, team outside department", middleware: deptAdmin, departmentID: deptID, teamID: teamID, wantStatus: http.StatusNotFound},
		{name: "site admin, mismatched path", middleware: deptAdmin, userType: thunderdome.AdminUserType, departmentID: deptID, wantStatus: http.StatusNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orgSvc := new(MockOrganizationDataService)
			if tt.userType != thunderdome.AdminUserType {
				orgSvc.On("OrganizationUserRole", mock.Anything, userID, orgID).Return(thunderdome.AdminUserType, nil).Maybe()
				orgSvc.On("DepartmentUserRole", mock.Anything, userID, orgID, tt.departmentID).Return(thunderdome.AdminUserType, "", nil).Maybe()
			}
			orgSvc.On("OrganizationPathExists", mock.Anything, orgID, tt.departmentID, tt.teamID).Return(tt.exists, nil).Once()
			s := &Service{OrganizationDataSvc: orgSvc}

			called := false
			handler := tt.middleware(s, func(w http.ResponseWriter, r *http.Request) { called = true })

			ctx := context.WithValue(context.Background(), contextKeyUserID, userID)
			ctx = context.WithValue(ctx, contextKeyUserType, tt.userType)
			req := httptest.NewRequestWithContext(ctx, http.MethodGet, "/", nil)
			req.SetPathValue("orgId", orgID)
			req.SetPathValue("departmentId", tt.departmentID)
			req.SetPathValue("teamId", tt.teamID)
			rec := httptest.NewRecorder()

			handler(rec, req)

			assert.Equal(t, tt.wantStatus, rec.Code)
			assert.Equal(t, tt.wantStatus == http.StatusOK, called)
			orgSvc.AssertExpectations(t)
		})
	}
}
