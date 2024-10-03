package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest/observer"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func ptr[T any](v T) *T {
	return &v
}

// MockTeamDataSvc is a mock implementation of the TeamDataSvc
type MockTeamDataSvc struct {
	mock.Mock
}

func (m *MockTeamDataSvc) TeamUserRoleByUserID(ctx context.Context, UserID string, TeamID string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockTeamDataSvc) TeamGetByID(ctx context.Context, TeamID string) (*thunderdome.Team, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockTeamDataSvc) TeamListByUser(ctx context.Context, UserID string, Limit int, Offset int) []*thunderdome.UserTeam {
	//TODO implement me
	panic("implement me")
}

func (m *MockTeamDataSvc) TeamListByUserNonOrg(ctx context.Context, UserID string, Limit int, Offset int) []*thunderdome.UserTeam {
	//TODO implement me
	panic("implement me")
}

func (m *MockTeamDataSvc) TeamCreate(ctx context.Context, UserID string, TeamName string) (*thunderdome.Team, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockTeamDataSvc) TeamUpdate(ctx context.Context, teamID string, TeamName string) (*thunderdome.Team, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockTeamDataSvc) TeamAddUser(ctx context.Context, TeamID string, UserID string, Role string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockTeamDataSvc) TeamUserList(ctx context.Context, TeamID string, Limit int, Offset int) ([]*thunderdome.TeamUser, int, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockTeamDataSvc) TeamUpdateUser(ctx context.Context, TeamID string, UserID string, Role string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockTeamDataSvc) TeamRemoveUser(ctx context.Context, TeamID string, UserID string) error {
	//TODO implement me
	panic("implement me")
}

func (m *MockTeamDataSvc) TeamInviteUser(ctx context.Context, TeamID string, Email string, Role string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockTeamDataSvc) TeamUserGetInviteByID(ctx context.Context, InviteID string) (thunderdome.TeamUserInvite, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockTeamDataSvc) TeamDeleteUserInvite(ctx context.Context, InviteID string) error {
	//TODO implement me
	panic("implement me")
}

func (m *MockTeamDataSvc) TeamGetUserInvites(ctx context.Context, teamID string) ([]thunderdome.TeamUserInvite, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockTeamDataSvc) TeamPokerList(ctx context.Context, TeamID string, Limit int, Offset int) []*thunderdome.Poker {
	//TODO implement me
	panic("implement me")
}

func (m *MockTeamDataSvc) TeamAddPoker(ctx context.Context, TeamID string, PokerID string) error {
	//TODO implement me
	panic("implement me")
}

func (m *MockTeamDataSvc) TeamRemovePoker(ctx context.Context, TeamID string, PokerID string) error {
	//TODO implement me
	panic("implement me")
}

func (m *MockTeamDataSvc) TeamDelete(ctx context.Context, TeamID string) error {
	//TODO implement me
	panic("implement me")
}

func (m *MockTeamDataSvc) TeamRetroList(ctx context.Context, TeamID string, Limit int, Offset int) []*thunderdome.Retro {
	//TODO implement me
	panic("implement me")
}

func (m *MockTeamDataSvc) TeamAddRetro(ctx context.Context, TeamID string, RetroID string) error {
	//TODO implement me
	panic("implement me")
}

func (m *MockTeamDataSvc) TeamRemoveRetro(ctx context.Context, TeamID string, RetroID string) error {
	//TODO implement me
	panic("implement me")
}

func (m *MockTeamDataSvc) TeamStoryboardList(ctx context.Context, TeamID string, Limit int, Offset int) []*thunderdome.Storyboard {
	//TODO implement me
	panic("implement me")
}

func (m *MockTeamDataSvc) TeamAddStoryboard(ctx context.Context, TeamID string, StoryboardID string) error {
	//TODO implement me
	panic("implement me")
}

func (m *MockTeamDataSvc) TeamRemoveStoryboard(ctx context.Context, TeamID string, StoryboardID string) error {
	//TODO implement me
	panic("implement me")
}

func (m *MockTeamDataSvc) TeamList(ctx context.Context, Limit int, Offset int) ([]*thunderdome.Team, int) {
	//TODO implement me
	panic("implement me")
}

func (m *MockTeamDataSvc) TeamIsSubscribed(ctx context.Context, teamID string) (bool, error) {
	args := m.Called(ctx, teamID)
	return args.Bool(0), args.Error(1)
}

func (m *MockTeamDataSvc) GetTeamMetrics(ctx context.Context, teamID string) (*thunderdome.TeamMetrics, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockTeamDataSvc) TeamUserRolesByUserID(ctx context.Context, userID, teamID string) (*thunderdome.UserTeamRoleInfo, error) {
	args := m.Called(ctx, userID, teamID)
	utr := args.Get(0).(thunderdome.UserTeamRoleInfo)
	return &utr, args.Error(1)
}

// MockLogger is a mock implementation of the Logger
type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Ctx(ctx context.Context) *zap.Logger {
	args := m.Called(ctx)
	return args.Get(0).(*zap.Logger)
}

func TestTeamUserOnly(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		userType       string
		teamID         string
		setupMocks     func(*MockTeamDataSvc, *MockLogger)
		expectedStatus int
	}{
		{
			name:     "Valid team user",
			userID:   "2d6176c8-50d6-4963-8172-2c20ca5022a3",
			userType: thunderdome.RegisteredUserType,
			teamID:   "128ee064-62ca-43b2-9fca-9c1089c89bd2",
			setupMocks: func(mtds *MockTeamDataSvc, ml *MockLogger) {
				mtds.On(
					"TeamUserRolesByUserID",
					mock.Anything,
					"2d6176c8-50d6-4963-8172-2c20ca5022a3",
					"128ee064-62ca-43b2-9fca-9c1089c89bd2",
				).Return(thunderdome.UserTeamRoleInfo{AssociationLevel: "TEAM"}, nil)
				ml.On("Ctx", mock.Anything).Return(zap.NewNop())
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:     "Global admin bypass",
			userID:   "306cf097-d75d-4ec0-8960-a5c2914b28b9",
			userType: thunderdome.AdminUserType,
			teamID:   "67cdc2f5-b4b0-444b-af0b-ae686e4cf9c8",
			setupMocks: func(mtds *MockTeamDataSvc, ml *MockLogger) {
				mtds.On(
					"TeamUserRolesByUserID",
					mock.Anything,
					"306cf097-d75d-4ec0-8960-a5c2914b28b9",
					"67cdc2f5-b4b0-444b-af0b-ae686e4cf9c8",
				).Return(thunderdome.UserTeamRoleInfo{}, nil)
				ml.On("Ctx", mock.Anything).Return(zap.NewNop())
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:     "Department admin",
			userID:   "6a12ef8c-1140-4faa-a505-22910a7593f9",
			userType: thunderdome.RegisteredUserType,
			teamID:   "7d4d0a17-cb20-4499-bc7f-ecaf2e77c15a",
			setupMocks: func(mtds *MockTeamDataSvc, ml *MockLogger) {
				deptRole := thunderdome.AdminUserType
				mtds.On(
					"TeamUserRolesByUserID",
					mock.Anything,
					"6a12ef8c-1140-4faa-a505-22910a7593f9",
					"7d4d0a17-cb20-4499-bc7f-ecaf2e77c15a",
				).Return(thunderdome.UserTeamRoleInfo{DepartmentRole: &deptRole}, nil)
				ml.On("Ctx", mock.Anything).Return(zap.NewNop())
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:     "Organization admin",
			userID:   "31c8521e-2e68-4898-b3cf-e919cf80dbe2",
			userType: thunderdome.RegisteredUserType,
			teamID:   "0ea230df-b5fe-47ae-a473-5153004eebdd",
			setupMocks: func(mtds *MockTeamDataSvc, ml *MockLogger) {
				orgRole := thunderdome.AdminUserType
				mtds.On(
					"TeamUserRolesByUserID",
					mock.Anything,
					"31c8521e-2e68-4898-b3cf-e919cf80dbe2",
					"0ea230df-b5fe-47ae-a473-5153004eebdd",
				).Return(thunderdome.UserTeamRoleInfo{OrganizationRole: &orgRole}, nil)
				ml.On("Ctx", mock.Anything).Return(zap.NewNop())
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:     "Invalid team ID",
			userID:   "1353a056-d239-41e1-ad1a-b3f0777e6c3a",
			userType: thunderdome.RegisteredUserType,
			teamID:   "invalid-team-id",
			setupMocks: func(mtds *MockTeamDataSvc, ml *MockLogger) {
				ml.On("Ctx", mock.Anything).Return(zap.NewNop())
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:     "Team not found",
			userID:   "1b853ef6-2c28-4c8e-ac29-a9d4827774fe",
			userType: thunderdome.RegisteredUserType,
			teamID:   "e52b9251-4722-4b4d-8f97-204ce7e51eec",
			setupMocks: func(mtds *MockTeamDataSvc, ml *MockLogger) {
				mtds.On(
					"TeamUserRolesByUserID",
					mock.Anything,
					"1b853ef6-2c28-4c8e-ac29-a9d4827774fe",
					"e52b9251-4722-4b4d-8f97-204ce7e51eec",
				).Return(thunderdome.UserTeamRoleInfo{}, fmt.Errorf("TEAM_NOT_FOUND"))
				ml.On("Ctx", mock.Anything).Return(zap.NewNop())
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:     "Unauthorized user",
			userID:   "a805def1-e1fa-42a9-b5f6-ee338799fa77",
			userType: thunderdome.RegisteredUserType,
			teamID:   "a805def1-e1fa-42a9-b5f6-ee338799fa77",
			setupMocks: func(mtds *MockTeamDataSvc, ml *MockLogger) {
				mtds.On("TeamUserRolesByUserID",
					mock.Anything,
					"a805def1-e1fa-42a9-b5f6-ee338799fa77",
					"a805def1-e1fa-42a9-b5f6-ee338799fa77",
				).Return(thunderdome.UserTeamRoleInfo{}, nil)
				ml.On("Ctx", mock.Anything).Return(zap.NewNop())
			},
			expectedStatus: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockTeamDataSvc := new(MockTeamDataSvc)
			mockLogger := new(MockLogger)

			tt.setupMocks(mockTeamDataSvc, mockLogger)

			s := &Service{
				TeamDataSvc: mockTeamDataSvc,
				Logger:      otelzap.New(mockLogger.Ctx(context.Background())),
			}

			handler := s.teamUserOnly(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			req, _ := http.NewRequest("GET", "/teams/"+tt.teamID, nil)
			req = mux.SetURLVars(req, map[string]string{"teamId": tt.teamID})
			req = req.WithContext(context.WithValue(req.Context(), contextKeyUserID, tt.userID))
			req = req.WithContext(context.WithValue(req.Context(), contextKeyUserType, tt.userType))

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			mockTeamDataSvc.AssertExpectations(t)
			mockLogger.AssertExpectations(t)
		})
	}
}

func TestTeamAdminOnly(t *testing.T) {
	tests := []struct {
		name             string
		userType         string
		teamRole         *string
		departmentRole   *string
		organizationRole *string
		setupMocks       func(*MockLogger)
		expectedStatus   int
	}{
		{
			name:             "Global Admin",
			userType:         thunderdome.AdminUserType,
			teamRole:         nil,
			departmentRole:   nil,
			organizationRole: nil,
			expectedStatus:   http.StatusOK,
		},
		{
			name:             "Team Admin",
			userType:         thunderdome.RegisteredUserType,
			teamRole:         ptr(thunderdome.AdminUserType),
			departmentRole:   nil,
			organizationRole: nil,
			expectedStatus:   http.StatusOK,
		},
		{
			name:             "Department Admin",
			userType:         thunderdome.RegisteredUserType,
			teamRole:         nil,
			departmentRole:   ptr(thunderdome.AdminUserType),
			organizationRole: nil,
			expectedStatus:   http.StatusOK,
		},
		{
			name:             "Organization Admin",
			userType:         thunderdome.RegisteredUserType,
			teamRole:         nil,
			departmentRole:   nil,
			organizationRole: ptr(thunderdome.AdminUserType),
			expectedStatus:   http.StatusOK,
		},
		{
			name:             "Non-Admin User",
			userType:         thunderdome.RegisteredUserType,
			teamRole:         ptr(thunderdome.EntityMemberUserType),
			departmentRole:   nil,
			organizationRole: nil,
			expectedStatus:   http.StatusForbidden,
		},
		{
			name:             "No Roles",
			userType:         thunderdome.RegisteredUserType,
			teamRole:         nil,
			departmentRole:   nil,
			organizationRole: nil,
			expectedStatus:   http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLogger := new(MockLogger)

			mockLogger.On("Ctx", mock.Anything).Return(zap.NewNop())

			// Create a test service
			s := &Service{
				Logger: otelzap.New(mockLogger.Ctx(context.Background())),
			}

			// Mock handler
			mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			// Create a new request
			req, err := http.NewRequest("GET", "/test", nil)
			assert.NoError(t, err)

			// Set up the context with user type and roles
			ctx := context.WithValue(req.Context(), contextKeyUserType, tt.userType)
			ctx = context.WithValue(ctx, contextKeyUserTeamRoles, &thunderdome.UserTeamRoleInfo{
				TeamRole:         tt.teamRole,
				DepartmentRole:   tt.departmentRole,
				OrganizationRole: tt.organizationRole,
			})
			req = req.WithContext(ctx)

			// Create a response recorder
			rr := httptest.NewRecorder()

			// Call the middleware
			handler := s.teamAdminOnly(mockHandler)
			handler.ServeHTTP(rr, req)

			// Check the status code
			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}

func TestSubscribedTeamOnly(t *testing.T) {
	tests := []struct {
		name                 string
		userType             string
		teamID               string
		subscriptionsEnabled bool
		expectedStatusCode   int
		setupMocks           func(*MockTeamDataSvc)
	}{
		{
			name:                 "Admin user bypasses subscription check",
			userType:             thunderdome.AdminUserType,
			teamID:               "1353a056-d239-41e1-ad1a-b3f0777e6c3a",
			subscriptionsEnabled: true,
			expectedStatusCode:   http.StatusOK,
		},
		{
			name:                 "Subscribed team allowed",
			userType:             thunderdome.EntityMemberUserType,
			teamID:               "2d6176c8-50d6-4963-8172-2c20ca5022a3",
			subscriptionsEnabled: true,
			expectedStatusCode:   http.StatusOK,
			setupMocks: func(mockTeamDataSvc *MockTeamDataSvc) {
				mockTeamDataSvc.On(
					"TeamIsSubscribed",
					mock.Anything,
					"2d6176c8-50d6-4963-8172-2c20ca5022a3",
				).Return(true, nil).Once()
			},
		},
		{
			name:                 "Unsubscribed team forbidden",
			userType:             thunderdome.EntityMemberUserType,
			teamID:               "128ee064-62ca-43b2-9fca-9c1089c89bd2",
			subscriptionsEnabled: true,
			expectedStatusCode:   http.StatusForbidden,
			setupMocks: func(mockTeamDataSvc *MockTeamDataSvc) {
				mockTeamDataSvc.On(
					"TeamIsSubscribed",
					mock.Anything,
					"128ee064-62ca-43b2-9fca-9c1089c89bd2",
				).Return(false, nil).Once()
			},
		},
		{
			name:               "Invalid teamID",
			userType:           thunderdome.EntityMemberUserType,
			teamID:             "invalid-uuid",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:                 "Subscriptions disabled",
			userType:             thunderdome.EntityMemberUserType,
			teamID:               "31c8521e-2e68-4898-b3cf-e919cf80dbe2",
			subscriptionsEnabled: false,
			expectedStatusCode:   http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockTeamDataSvc := new(MockTeamDataSvc)
			mockConfig := &Config{SubscriptionsEnabled: true}
			service := &Service{
				TeamDataSvc: mockTeamDataSvc,
				Config:      mockConfig,
			}

			mockConfig.SubscriptionsEnabled = tt.subscriptionsEnabled

			if tt.setupMocks != nil {
				tt.setupMocks(mockTeamDataSvc)
			}

			handler := service.subscribedTeamOnly(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			req := httptest.NewRequest("GET", "/teams/"+tt.teamID+"/test", nil)
			req = mux.SetURLVars(req, map[string]string{"teamId": tt.teamID})
			req = req.WithContext(context.WithValue(req.Context(), contextKeyUserType, tt.userType))

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatusCode, rr.Code)

			mockTeamDataSvc.AssertExpectations(t)
		})
	}
}

type MockOrganizationDataService struct {
	mock.Mock
}

func (m *MockOrganizationDataService) OrganizationGetByID(ctx context.Context, OrgID string) (*thunderdome.Organization, error) {
	args := m.Called(ctx, OrgID)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	to := args.Get(0).(thunderdome.Organization)
	return &to, nil
}

func (m *MockOrganizationDataService) OrganizationUserRole(ctx context.Context, userID, orgID string) (string, error) {
	args := m.Called(ctx, userID, orgID)
	return args.String(0), args.Error(1)
}

func (m *MockOrganizationDataService) OrganizationListByUser(ctx context.Context, UserID string, Limit int, Offset int) []*thunderdome.UserOrganization {
	//TODO implement me
	panic("implement me")
}

func (m *MockOrganizationDataService) OrganizationCreate(ctx context.Context, UserID string, OrgName string) (*thunderdome.Organization, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockOrganizationDataService) OrganizationUpdate(ctx context.Context, OrgID string, OrgName string) (*thunderdome.Organization, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockOrganizationDataService) OrganizationUserList(ctx context.Context, OrgID string, Limit int, Offset int) []*thunderdome.OrganizationUser {
	//TODO implement me
	panic("implement me")
}

func (m *MockOrganizationDataService) OrganizationAddUser(ctx context.Context, OrgID string, UserID string, Role string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockOrganizationDataService) OrganizationUpsertUser(ctx context.Context, OrgID string, UserID string, Role string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockOrganizationDataService) OrganizationUpdateUser(ctx context.Context, OrgID string, UserID string, Role string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockOrganizationDataService) OrganizationRemoveUser(ctx context.Context, OrganizationID string, UserID string) error {
	//TODO implement me
	panic("implement me")
}

func (m *MockOrganizationDataService) OrganizationInviteUser(ctx context.Context, OrgID string, Email string, Role string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockOrganizationDataService) OrganizationUserGetInviteByID(ctx context.Context, InviteID string) (thunderdome.OrganizationUserInvite, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockOrganizationDataService) OrganizationDeleteUserInvite(ctx context.Context, InviteID string) error {
	//TODO implement me
	panic("implement me")
}

func (m *MockOrganizationDataService) OrganizationGetUserInvites(ctx context.Context, orgID string) ([]thunderdome.OrganizationUserInvite, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockOrganizationDataService) OrganizationTeamList(ctx context.Context, OrgID string, Limit int, Offset int) []*thunderdome.Team {
	//TODO implement me
	panic("implement me")
}

func (m *MockOrganizationDataService) OrganizationTeamCreate(ctx context.Context, OrgID string, TeamName string) (*thunderdome.Team, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockOrganizationDataService) OrganizationTeamUserRole(ctx context.Context, UserID string, OrgID string, TeamID string) (string, string, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockOrganizationDataService) OrganizationDelete(ctx context.Context, OrgID string) error {
	//TODO implement me
	panic("implement me")
}

func (m *MockOrganizationDataService) OrganizationList(ctx context.Context, Limit int, Offset int) []*thunderdome.Organization {
	//TODO implement me
	panic("implement me")
}

func (m *MockOrganizationDataService) GetOrganizationMetrics(ctx context.Context, organizationID string) (*thunderdome.OrganizationMetrics, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockOrganizationDataService) DepartmentUserRole(ctx context.Context, userID, orgID, departmentID string) (string, string, error) {
	args := m.Called(ctx, userID, orgID, departmentID)
	return args.String(0), args.String(1), args.Error(2)
}

func (m *MockOrganizationDataService) DepartmentGetByID(ctx context.Context, DepartmentID string) (*thunderdome.Department, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockOrganizationDataService) OrganizationDepartmentList(ctx context.Context, OrgID string, Limit int, Offset int) []*thunderdome.Department {
	//TODO implement me
	panic("implement me")
}

func (m *MockOrganizationDataService) DepartmentCreate(ctx context.Context, OrgID string, OrgName string) (*thunderdome.Department, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockOrganizationDataService) DepartmentUpdate(ctx context.Context, deptID string, DeptName string) (*thunderdome.Department, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockOrganizationDataService) DepartmentTeamList(ctx context.Context, DepartmentID string, Limit int, Offset int) []*thunderdome.Team {
	//TODO implement me
	panic("implement me")
}

func (m *MockOrganizationDataService) DepartmentTeamCreate(ctx context.Context, DepartmentID string, TeamName string) (*thunderdome.Team, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockOrganizationDataService) DepartmentUserList(ctx context.Context, DepartmentID string, Limit int, Offset int) []*thunderdome.DepartmentUser {
	//TODO implement me
	panic("implement me")
}

func (m *MockOrganizationDataService) DepartmentAddUser(ctx context.Context, DepartmentID string, UserID string, Role string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockOrganizationDataService) DepartmentUpsertUser(ctx context.Context, DepartmentID string, UserID string, Role string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockOrganizationDataService) DepartmentUpdateUser(ctx context.Context, DepartmentID string, UserID string, Role string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockOrganizationDataService) DepartmentRemoveUser(ctx context.Context, DepartmentID string, UserID string) error {
	//TODO implement me
	panic("implement me")
}

func (m *MockOrganizationDataService) DepartmentTeamUserRole(ctx context.Context, UserID string, OrgID string, DepartmentID string, TeamID string) (string, string, string, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockOrganizationDataService) DepartmentDelete(ctx context.Context, DepartmentID string) error {
	//TODO implement me
	panic("implement me")
}

func (m *MockOrganizationDataService) DepartmentInviteUser(ctx context.Context, DeptID string, Email string, Role string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockOrganizationDataService) DepartmentUserGetInviteByID(ctx context.Context, InviteID string) (thunderdome.DepartmentUserInvite, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockOrganizationDataService) DepartmentDeleteUserInvite(ctx context.Context, InviteID string) error {
	//TODO implement me
	panic("implement me")
}

func (m *MockOrganizationDataService) DepartmentGetUserInvites(ctx context.Context, deptID string) ([]thunderdome.DepartmentUserInvite, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockOrganizationDataService) OrganizationIsSubscribed(ctx context.Context, orgID string) (bool, error) {
	args := m.Called(ctx, orgID)
	return args.Bool(0), args.Error(1)
}

func TestSubscribedOrgOnly(t *testing.T) {
	tests := []struct {
		name                 string
		userType             string
		orgID                string
		subscriptionsEnabled bool
		isSubscribed         bool
		expectedStatusCode   int
		setupMocks           func(*MockOrganizationDataService)
	}{
		{
			name:                 "Admin user bypasses subscription check",
			userType:             thunderdome.AdminUserType,
			orgID:                "1353a056-d239-41e1-ad1a-b3f0777e6c3a",
			subscriptionsEnabled: true,
			isSubscribed:         false,
			expectedStatusCode:   http.StatusOK,
		},
		{
			name:                 "Subscribed organization allowed",
			userType:             thunderdome.EntityMemberUserType,
			orgID:                "2d6176c8-50d6-4963-8172-2c20ca5022a3",
			subscriptionsEnabled: true,
			isSubscribed:         true,
			expectedStatusCode:   http.StatusOK,
			setupMocks: func(mockOrgDataSvc *MockOrganizationDataService) {
				mockOrgDataSvc.On(
					"OrganizationIsSubscribed",
					mock.Anything,
					"2d6176c8-50d6-4963-8172-2c20ca5022a3",
				).Return(true, nil).Once()
			},
		},
		{
			name:                 "Unsubscribed organization forbidden",
			userType:             thunderdome.EntityMemberUserType,
			orgID:                "128ee064-62ca-43b2-9fca-9c1089c89bd2",
			subscriptionsEnabled: true,
			isSubscribed:         false,
			expectedStatusCode:   http.StatusForbidden,
			setupMocks: func(mockOrgDataSvc *MockOrganizationDataService) {
				mockOrgDataSvc.On(
					"OrganizationIsSubscribed",
					mock.Anything,
					"128ee064-62ca-43b2-9fca-9c1089c89bd2",
				).Return(false, nil).Once()
			},
		},
		{
			name:               "Invalid orgID",
			userType:           thunderdome.EntityMemberUserType,
			orgID:              "invalid-uuid",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:                 "Subscriptions disabled",
			userType:             thunderdome.EntityMemberUserType,
			orgID:                "31c8521e-2e68-4898-b3cf-e919cf80dbe2",
			subscriptionsEnabled: false,
			isSubscribed:         false,
			expectedStatusCode:   http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockOrgDataSvc := new(MockOrganizationDataService)
			mockConfig := &Config{SubscriptionsEnabled: true}
			service := &Service{
				OrganizationDataSvc: mockOrgDataSvc,
				Config:              mockConfig,
			}

			mockConfig.SubscriptionsEnabled = tt.subscriptionsEnabled

			if tt.setupMocks != nil {
				tt.setupMocks(mockOrgDataSvc)
			}

			handler := service.subscribedOrgOnly(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			req := httptest.NewRequest("GET", "/organizations/"+tt.orgID+"/test", nil)
			req = mux.SetURLVars(req, map[string]string{"orgId": tt.orgID})
			req = req.WithContext(context.WithValue(req.Context(), contextKeyUserType, tt.userType))

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatusCode, rr.Code)

			mockOrgDataSvc.AssertExpectations(t)
		})
	}
}

func TestDepartmentAdminOnly(t *testing.T) {
	tests := []struct {
		name               string
		userType           string
		userID             string
		orgID              string
		departmentID       string
		expectedStatusCode int
		expectedOrgRole    string
		expectedDeptRole   string
		setupMocks         func(*MockOrganizationDataService)
	}{
		{
			name:               "Admin user bypasses role check",
			userType:           thunderdome.AdminUserType,
			userID:             "002738c2-fcf2-438e-a755-2bf9c4233b74",
			orgID:              "00280040-11e2-4c00-9cd0-325b4efd7de2",
			departmentID:       "002a5613-5d31-430d-859f-0e97fee8c5d2",
			expectedStatusCode: http.StatusOK,
			expectedOrgRole:    thunderdome.AdminUserType,
			expectedDeptRole:   thunderdome.AdminUserType,
		},
		{
			name:               "Department admin allowed",
			userType:           "USER",
			userID:             "0023f0d5-19d0-403f-a30d-d5b7616c72bd",
			orgID:              "00241406-94cb-4fe1-9b4a-e979f5761f10",
			departmentID:       "0024d5f8-42b1-46c7-b5a7-da97d59d6b36",
			expectedStatusCode: http.StatusOK,
			expectedOrgRole:    thunderdome.EntityMemberUserType,
			expectedDeptRole:   thunderdome.AdminUserType,
			setupMocks: func(mockOrgDataSvc *MockOrganizationDataService) {
				mockOrgDataSvc.On(
					"DepartmentUserRole",
					mock.Anything,
					"0023f0d5-19d0-403f-a30d-d5b7616c72bd",
					"00241406-94cb-4fe1-9b4a-e979f5761f10",
					"0024d5f8-42b1-46c7-b5a7-da97d59d6b36",
				).Return(thunderdome.EntityMemberUserType, thunderdome.AdminUserType, nil).Once()
			},
		},
		{
			name:               "Organization admin allowed",
			userType:           "USER",
			userID:             "0014c2dc-3e89-4369-857d-b420f5786eff",
			orgID:              "0019f42b-de8f-41ee-904e-3cb6c9ddc8f4",
			departmentID:       "001a6acb-c174-43f0-9930-0b1242931123",
			expectedStatusCode: http.StatusOK,
			expectedOrgRole:    thunderdome.AdminUserType,
			expectedDeptRole:   "",
			setupMocks: func(mockOrgDataSvc *MockOrganizationDataService) {
				mockOrgDataSvc.On(
					"DepartmentUserRole",
					mock.Anything,
					"0014c2dc-3e89-4369-857d-b420f5786eff",
					"0019f42b-de8f-41ee-904e-3cb6c9ddc8f4",
					"001a6acb-c174-43f0-9930-0b1242931123",
				).Return(thunderdome.AdminUserType, "", nil).Once()
			},
		},
		{
			name:               "Non-admin user forbidden",
			userType:           "USER",
			userID:             "3f3d4ca5-6eae-4372-81ba-de8bbaa2dac2",
			orgID:              "2d6176c8-50d6-4963-8172-2c20ca5022a3",
			departmentID:       "002738c2-fcf2-438e-a755-2bf9c4233b74",
			expectedStatusCode: http.StatusForbidden,
			setupMocks: func(mockOrgDataSvc *MockOrganizationDataService) {
				mockOrgDataSvc.On(
					"DepartmentUserRole",
					mock.Anything,
					"3f3d4ca5-6eae-4372-81ba-de8bbaa2dac2",
					"2d6176c8-50d6-4963-8172-2c20ca5022a3",
					"002738c2-fcf2-438e-a755-2bf9c4233b74",
				).Return(thunderdome.EntityMemberUserType, thunderdome.EntityMemberUserType, nil).Once()
			},
		},
		{
			name:               "Invalid orgID",
			userType:           "USER",
			userID:             "ea840339-2e16-4c10-8744-33ee1b636596",
			orgID:              "invalid-org-id",
			departmentID:       "6d72a7dd-4183-4772-8e58-bda948416974",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "Invalid departmentID",
			userType:           "USER",
			userID:             "ea840339-2e16-4c10-8744-33ee1b636596",
			orgID:              "31c8521e-2e68-4898-b3cf-e919cf80dbe2",
			departmentID:       "invalid-dept-id",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "User not found in department",
			userType:           "USER",
			userID:             "ea840339-2e16-4c10-8744-33ee1b636596",
			orgID:              "128ee064-62ca-43b2-9fca-9c1089c89bd2",
			departmentID:       "31c8521e-2e68-4898-b3cf-e919cf80dbe2",
			expectedStatusCode: http.StatusForbidden,
			setupMocks: func(mockOrgDataSvc *MockOrganizationDataService) {
				mockOrgDataSvc.On(
					"DepartmentUserRole",
					mock.Anything,
					"ea840339-2e16-4c10-8744-33ee1b636596",
					"128ee064-62ca-43b2-9fca-9c1089c89bd2",
					"31c8521e-2e68-4898-b3cf-e919cf80dbe2",
				).Return("", "", fmt.Errorf("User not found")).Once()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockOrgDataSvc := new(MockOrganizationDataService)
			service := &Service{
				OrganizationDataSvc: mockOrgDataSvc,
			}

			if tt.setupMocks != nil {
				tt.setupMocks(mockOrgDataSvc)
			}

			var capturedOrgRole, capturedDeptRole string
			handler := service.departmentAdminOnly(func(w http.ResponseWriter, r *http.Request) {
				capturedOrgRole = r.Context().Value(contextKeyOrgRole).(string)
				capturedDeptRole = r.Context().Value(contextKeyDepartmentRole).(string)
				w.WriteHeader(http.StatusOK)
			})

			req := httptest.NewRequest("GET", "/organizations/"+tt.orgID+"/departments/"+tt.departmentID+"/test", nil)
			req = mux.SetURLVars(req, map[string]string{"orgId": tt.orgID, "departmentId": tt.departmentID})
			req = req.WithContext(context.WithValue(req.Context(), contextKeyUserID, tt.userID))
			req = req.WithContext(context.WithValue(req.Context(), contextKeyUserType, tt.userType))

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatusCode, rr.Code)

			if tt.expectedStatusCode == http.StatusOK {
				assert.Equal(t, tt.expectedOrgRole, capturedOrgRole)
				assert.Equal(t, tt.expectedDeptRole, capturedDeptRole)
			}

			mockOrgDataSvc.AssertExpectations(t)
		})
	}
}

func TestDepartmentUserOnly(t *testing.T) {
	tests := []struct {
		name               string
		userType           string
		userID             string
		orgID              string
		departmentID       string
		expectedStatusCode int
		expectedOrgRole    string
		expectedDeptRole   string
		setupMocks         func(*MockOrganizationDataService, *MockLogger)
	}{
		{
			name:               "Admin user bypasses role check",
			userType:           thunderdome.AdminUserType,
			userID:             "00320d41-1a7c-4e3f-b202-b74ab6ea582c",
			orgID:              "128ee064-62ca-43b2-9fca-9c1089c89bd2",
			departmentID:       "31c8521e-2e68-4898-b3cf-e919cf80dbe2",
			expectedStatusCode: http.StatusOK,
			expectedOrgRole:    thunderdome.AdminUserType,
			expectedDeptRole:   thunderdome.AdminUserType,
		},
		{
			name:               "Department user allowed",
			userType:           "USER",
			userID:             "ea840339-2e16-4c10-8744-33ee1b636596",
			orgID:              "128ee064-62ca-43b2-9fca-9c1089c89bd2",
			departmentID:       "31c8521e-2e68-4898-b3cf-e919cf80dbe2",
			expectedStatusCode: http.StatusOK,
			expectedOrgRole:    thunderdome.EntityMemberUserType,
			expectedDeptRole:   thunderdome.EntityMemberUserType,
			setupMocks: func(mockOrgDataSvc *MockOrganizationDataService, mockLogger *MockLogger) {
				mockOrgDataSvc.On(
					"DepartmentUserRole",
					mock.Anything,
					"ea840339-2e16-4c10-8744-33ee1b636596",
					"128ee064-62ca-43b2-9fca-9c1089c89bd2",
					"31c8521e-2e68-4898-b3cf-e919cf80dbe2",
				).Return(thunderdome.EntityMemberUserType, thunderdome.EntityMemberUserType, nil).Once()
			},
		},
		{
			name:               "Department admin allowed",
			userType:           "USER",
			userID:             "3f3d4ca5-6eae-4372-81ba-de8bbaa2dac2",
			orgID:              "2d6176c8-50d6-4963-8172-2c20ca5022a3",
			departmentID:       "002738c2-fcf2-438e-a755-2bf9c4233b74",
			expectedStatusCode: http.StatusOK,
			expectedOrgRole:    thunderdome.EntityMemberUserType,
			expectedDeptRole:   thunderdome.AdminUserType,
			setupMocks: func(mockOrgDataSvc *MockOrganizationDataService, mockLogger *MockLogger) {
				mockOrgDataSvc.On(
					"DepartmentUserRole",
					mock.Anything,
					"3f3d4ca5-6eae-4372-81ba-de8bbaa2dac2",
					"2d6176c8-50d6-4963-8172-2c20ca5022a3",
					"002738c2-fcf2-438e-a755-2bf9c4233b74",
				).Return(thunderdome.EntityMemberUserType, thunderdome.AdminUserType, nil).Once()
			},
		},
		{
			name:               "Org Member Non-department user forbidden",
			userType:           "USER",
			userID:             "0014c2dc-3e89-4369-857d-b420f5786eff",
			orgID:              "0019f42b-de8f-41ee-904e-3cb6c9ddc8f4",
			departmentID:       "001a6acb-c174-43f0-9930-0b1242931123",
			expectedStatusCode: http.StatusForbidden,
			setupMocks: func(mockOrgDataSvc *MockOrganizationDataService, mockLogger *MockLogger) {
				mockOrgDataSvc.On(
					"DepartmentUserRole",
					mock.Anything,
					"0014c2dc-3e89-4369-857d-b420f5786eff",
					"0019f42b-de8f-41ee-904e-3cb6c9ddc8f4",
					"001a6acb-c174-43f0-9930-0b1242931123",
				).Return(thunderdome.EntityMemberUserType, "", nil).Once()
			},
		},
		{
			name:               "Non-org Non-department user forbidden",
			userType:           "USER",
			userID:             "0014c2dc-3e89-4369-857d-b420f5786efg",
			orgID:              "0019f42b-de8f-41ee-904e-3cb6c9ddc8f5",
			departmentID:       "001a6acb-c174-43f0-9930-0b1242931124",
			expectedStatusCode: http.StatusForbidden,
			setupMocks: func(mockOrgDataSvc *MockOrganizationDataService, mockLogger *MockLogger) {
				mockOrgDataSvc.On(
					"DepartmentUserRole",
					mock.Anything,
					"0014c2dc-3e89-4369-857d-b420f5786efg",
					"0019f42b-de8f-41ee-904e-3cb6c9ddc8f5",
					"001a6acb-c174-43f0-9930-0b1242931124",
				).Return("", "", fmt.Errorf("error getting department users role")).Once()
			},
		},
		{
			name:               "Invalid orgID",
			userType:           "USER",
			userID:             "003a65d1-1f9c-4428-a27f-7a7156481412",
			orgID:              "invalid-org-id",
			departmentID:       "003ac614-10ff-472e-93cb-0fea025bccb8",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "Invalid departmentID",
			userType:           "USER",
			userID:             "003a65d1-1f9c-4428-a27f-7a7156481412",
			orgID:              "0019f42b-de8f-41ee-904e-3cb6c9ddc8f5",
			departmentID:       "invalid-dept-id",
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockOrgDataSvc := new(MockOrganizationDataService)
			mockLogger := new(MockLogger)

			mockLogger.On("Ctx", mock.Anything).Return(zap.NewNop())

			service := &Service{
				OrganizationDataSvc: mockOrgDataSvc,
				Logger:              otelzap.New(mockLogger.Ctx(context.Background())),
			}

			if tt.setupMocks != nil {
				tt.setupMocks(mockOrgDataSvc, mockLogger)
			}

			var capturedOrgRole, capturedDeptRole string
			handler := service.departmentUserOnly(func(w http.ResponseWriter, r *http.Request) {
				capturedOrgRole = r.Context().Value(contextKeyOrgRole).(string)
				capturedDeptRole = r.Context().Value(contextKeyDepartmentRole).(string)
				w.WriteHeader(http.StatusOK)
			})

			req := httptest.NewRequest("GET", "/organizations/"+tt.orgID+"/departments/"+tt.departmentID+"/test", nil)
			req = mux.SetURLVars(req, map[string]string{"orgId": tt.orgID, "departmentId": tt.departmentID})
			req = req.WithContext(context.WithValue(req.Context(), contextKeyUserID, tt.userID))
			req = req.WithContext(context.WithValue(req.Context(), contextKeyUserType, tt.userType))

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatusCode, rr.Code)

			if tt.expectedStatusCode == http.StatusOK {
				assert.Equal(t, tt.expectedOrgRole, capturedOrgRole)
				assert.Equal(t, tt.expectedDeptRole, capturedDeptRole)
			}

			mockOrgDataSvc.AssertExpectations(t)
		})
	}
}

func TestOrgAdminOnly(t *testing.T) {
	tests := []struct {
		name            string
		userID          string
		userType        string
		orgID           string
		expectedOrgRole string
		expectedStatus  int
		mockSetup       func(mockOrgDataSvc *MockOrganizationDataService)
	}{
		{
			name:            "Valid Admin User",
			userID:          "123e4567-e89b-12d3-a456-426614174000",
			userType:        thunderdome.AdminUserType,
			orgID:           "123e4567-e89b-12d3-a456-426614174001",
			expectedOrgRole: thunderdome.AdminUserType,
			expectedStatus:  http.StatusOK,
			mockSetup:       func(mockOrgDataSvc *MockOrganizationDataService) {},
		},
		{
			name:            "Valid Org Admin",
			userID:          "223e4567-e89b-12d3-a456-426614174000",
			userType:        "REGULAR",
			orgID:           "223e4567-e89b-12d3-a456-426614174001",
			expectedOrgRole: thunderdome.AdminUserType,
			expectedStatus:  http.StatusOK,
			mockSetup: func(mockOrgDataSvc *MockOrganizationDataService) {
				mockOrgDataSvc.On(
					"OrganizationUserRole",
					mock.Anything,
					"223e4567-e89b-12d3-a456-426614174000",
					"223e4567-e89b-12d3-a456-426614174001",
				).Return(thunderdome.AdminUserType, nil)
			},
		},
		{
			name:            "Non-Admin User",
			userID:          "323e4567-e89b-12d3-a456-426614174000",
			userType:        thunderdome.RegisteredUserType,
			orgID:           "323e4567-e89b-12d3-a456-426614174001",
			expectedOrgRole: thunderdome.EntityMemberUserType,
			expectedStatus:  http.StatusForbidden,
			mockSetup: func(mockOrgDataSvc *MockOrganizationDataService) {
				mockOrgDataSvc.On("OrganizationUserRole", mock.Anything, "323e4567-e89b-12d3-a456-426614174000", "323e4567-e89b-12d3-a456-426614174001").Return(thunderdome.EntityMemberUserType, nil)
			},
		},
		{
			name:           "User Not in Organization",
			userID:         "423e4567-e89b-12d3-a456-426614174000",
			userType:       thunderdome.RegisteredUserType,
			orgID:          "423e4567-e89b-12d3-a456-426614174001",
			expectedStatus: http.StatusForbidden,
			mockSetup: func(mockOrgDataSvc *MockOrganizationDataService) {
				mockOrgDataSvc.On(
					"OrganizationUserRole",
					mock.Anything,
					"423e4567-e89b-12d3-a456-426614174000",
					"423e4567-e89b-12d3-a456-426614174001",
				).Return("", errors.New("user not found"))
			},
		},
		{
			name:           "Invalid OrgID",
			userID:         "523e4567-e89b-12d3-a456-426614174000",
			userType:       thunderdome.RegisteredUserType,
			orgID:          "invalid-org-id",
			expectedStatus: http.StatusBadRequest,
			mockSetup:      func(mockOrgDataSvc *MockOrganizationDataService) {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock OrganizationDataSvc
			mockOrgDataSvc := new(MockOrganizationDataService)

			// Create a new service with the mock
			s := &Service{
				OrganizationDataSvc: mockOrgDataSvc,
			}

			// Define a dummy handler for testing
			var capturedOrgRole string
			dummyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				capturedOrgRole = r.Context().Value(contextKeyOrgRole).(string)
				w.WriteHeader(http.StatusOK)
			})

			// Setup mock expectations
			tt.mockSetup(mockOrgDataSvc)

			// Create a new request
			req, err := http.NewRequest("GET", "/organizations/"+tt.orgID, nil)
			assert.NoError(t, err)

			// Create a new response recorder
			rr := httptest.NewRecorder()

			// Set up the context with user information
			ctx := context.WithValue(req.Context(), contextKeyUserID, tt.userID)
			ctx = context.WithValue(ctx, contextKeyUserType, tt.userType)
			req = req.WithContext(ctx)

			// Set up router with vars
			router := mux.NewRouter()
			router.HandleFunc("/organizations/{orgId}", s.orgAdminOnly(dummyHandler))

			// Serve the request
			router.ServeHTTP(rr, req)

			// Check the status code
			assert.Equal(t, tt.expectedStatus, rr.Code)

			// If the request was successful, check that the context was updated
			if tt.expectedStatus == http.StatusOK {
				assert.Equal(t, tt.expectedOrgRole, capturedOrgRole)
			}

			// Clear mock expectations for the next test
			mockOrgDataSvc.AssertExpectations(t)
		})
	}
}

func TestOrgUserOnly(t *testing.T) {
	tests := []struct {
		name            string
		userID          string
		userType        string
		orgID           string
		expectedOrgRole string
		expectedStatus  int
		mockSetup       func(mockOrgDataSvc *MockOrganizationDataService)
	}{
		{
			name:            "Valid Admin User",
			userID:          "123e4567-e89b-12d3-a456-426614174000",
			userType:        thunderdome.AdminUserType,
			orgID:           "123e4567-e89b-12d3-a456-426614174001",
			expectedOrgRole: thunderdome.AdminUserType,
			expectedStatus:  http.StatusOK,
			mockSetup: func(mockOrgDataSvc *MockOrganizationDataService) {
				mockOrgDataSvc.On(
					"OrganizationGetByID",
					mock.Anything,
					"123e4567-e89b-12d3-a456-426614174001",
				).Return(thunderdome.Organization{}, nil)
			},
		},
		{
			name:            "Valid Org User",
			userID:          "223e4567-e89b-12d3-a456-426614174000",
			userType:        thunderdome.RegisteredUserType,
			orgID:           "223e4567-e89b-12d3-a456-426614174001",
			expectedOrgRole: thunderdome.EntityMemberUserType,
			expectedStatus:  http.StatusOK,
			mockSetup: func(mockOrgDataSvc *MockOrganizationDataService) {
				mockOrgDataSvc.On(
					"OrganizationGetByID",
					mock.Anything, "223e4567-e89b-12d3-a456-426614174001",
				).Return(thunderdome.Organization{}, nil)
				mockOrgDataSvc.On(
					"OrganizationUserRole",
					mock.Anything,
					"223e4567-e89b-12d3-a456-426614174000",
					"223e4567-e89b-12d3-a456-426614174001",
				).Return(thunderdome.EntityMemberUserType, nil)
			},
		},
		{
			name:           "User Not in Organization",
			userID:         "323e4567-e89b-12d3-a456-426614174000",
			userType:       thunderdome.RegisteredUserType,
			orgID:          "323e4567-e89b-12d3-a456-426614174001",
			expectedStatus: http.StatusForbidden,
			mockSetup: func(mockOrgDataSvc *MockOrganizationDataService) {
				mockOrgDataSvc.On(
					"OrganizationGetByID",
					mock.Anything,
					"323e4567-e89b-12d3-a456-426614174001",
				).Return(thunderdome.Organization{}, nil)
				mockOrgDataSvc.On(
					"OrganizationUserRole",
					mock.Anything,
					"323e4567-e89b-12d3-a456-426614174000",
					"323e4567-e89b-12d3-a456-426614174001",
				).Return("", errors.New("ORGANIZATION_USER_REQUIRED"))
			},
		},
		{
			name:           "Invalid OrgID",
			userID:         "423e4567-e89b-12d3-a456-426614174000",
			userType:       thunderdome.RegisteredUserType,
			orgID:          "invalid-org-id",
			expectedStatus: http.StatusBadRequest,
			mockSetup:      func(mockOrgDataSvc *MockOrganizationDataService) {},
		},
		{
			name:           "Organization Not Found",
			userID:         "523e4567-e89b-12d3-a456-426614174000",
			userType:       thunderdome.RegisteredUserType,
			orgID:          "523e4567-e89b-12d3-a456-426614174001",
			expectedStatus: http.StatusNotFound,
			mockSetup: func(mockOrgDataSvc *MockOrganizationDataService) {
				mockOrgDataSvc.On(
					"OrganizationGetByID",
					mock.Anything,
					"523e4567-e89b-12d3-a456-426614174001",
				).Return(nil, errors.New("ORGANIZATION_NOT_FOUND"))
			},
		},
		{
			name:           "Internal Server Error",
			userID:         "623e4567-e89b-12d3-a456-426614174000",
			userType:       thunderdome.RegisteredUserType,
			orgID:          "623e4567-e89b-12d3-a456-426614174001",
			expectedStatus: http.StatusInternalServerError,
			mockSetup: func(mockOrgDataSvc *MockOrganizationDataService) {
				mockOrgDataSvc.On(
					"OrganizationGetByID",
					mock.Anything,
					"623e4567-e89b-12d3-a456-426614174001",
				).Return(nil, errors.New("unexpected error"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock OrganizationDataSvc
			mockOrgDataSvc := new(MockOrganizationDataService)

			// Create a new service with the mock
			s := &Service{
				OrganizationDataSvc: mockOrgDataSvc,
			}

			// Define a dummy handler for testing
			var capturedOrgRole string
			dummyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				capturedOrgRole = r.Context().Value(contextKeyOrgRole).(string)
				w.WriteHeader(http.StatusOK)
			})

			// Setup mock expectations
			tt.mockSetup(mockOrgDataSvc)

			// Create a new request
			req, err := http.NewRequest("GET", "/organizations/"+tt.orgID, nil)
			assert.NoError(t, err)

			// Create a new response recorder
			rr := httptest.NewRecorder()

			// Set up the context with user information
			ctx := context.WithValue(req.Context(), contextKeyUserID, tt.userID)
			ctx = context.WithValue(ctx, contextKeyUserType, tt.userType)
			req = req.WithContext(ctx)

			// Set up router with vars
			router := mux.NewRouter()
			router.HandleFunc("/organizations/{orgId}", s.orgUserOnly(dummyHandler))

			// Serve the request
			router.ServeHTTP(rr, req)

			// Check the status code
			assert.Equal(t, tt.expectedStatus, rr.Code)

			// If the request was successful, check that the context was updated
			if tt.expectedStatus == http.StatusOK {
				assert.Equal(t, tt.expectedOrgRole, capturedOrgRole)
			}

			// Clear mock expectations for the next test
			mockOrgDataSvc.AssertExpectations(t)
		})
	}
}

func TestSubscribedUserOnly(t *testing.T) {
	tests := []struct {
		name                 string
		userID               string
		userType             string
		subscriptionsEnabled bool
		expectedStatus       int
		mockSetup            func(mockSubDataSvc *MockSubscriptionDataService)
	}{
		{
			name:                 "Subscriptions Disabled",
			userID:               "123e4567-e89b-12d3-a456-426614174000",
			userType:             thunderdome.RegisteredUserType,
			subscriptionsEnabled: false,
			expectedStatus:       http.StatusOK,
			mockSetup:            func(mockSubDataSvc *MockSubscriptionDataService) {},
		},
		{
			name:                 "Admin User",
			userID:               "223e4567-e89b-12d3-a456-426614174000",
			userType:             thunderdome.AdminUserType,
			subscriptionsEnabled: true,
			expectedStatus:       http.StatusOK,
			mockSetup:            func(mockSubDataSvc *MockSubscriptionDataService) {},
		},
		{
			name:                 "Subscribed Regular User",
			userID:               "323e4567-e89b-12d3-a456-426614174000",
			userType:             thunderdome.RegisteredUserType,
			subscriptionsEnabled: true,
			expectedStatus:       http.StatusOK,
			mockSetup: func(mockSubDataSvc *MockSubscriptionDataService) {
				mockSubDataSvc.On("CheckActiveSubscriber", mock.Anything, "323e4567-e89b-12d3-a456-426614174000").Return(nil)
			},
		},
		{
			name:                 "Unsubscribed Regular User",
			userID:               "423e4567-e89b-12d3-a456-426614174000",
			userType:             thunderdome.RegisteredUserType,
			subscriptionsEnabled: true,
			expectedStatus:       http.StatusForbidden,
			mockSetup: func(mockSubDataSvc *MockSubscriptionDataService) {
				mockSubDataSvc.On("CheckActiveSubscriber", mock.Anything, "423e4567-e89b-12d3-a456-426614174000").Return(errors.New("not subscribed"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock SubscriptionDataService
			mockSubDataSvc := new(MockSubscriptionDataService)

			// Create a new service with the mock and config
			s := &Service{
				SubscriptionDataSvc: mockSubDataSvc,
				Config: &Config{
					SubscriptionsEnabled: tt.subscriptionsEnabled,
				},
			}

			// Define a dummy handler for testing
			dummyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			// Setup mock expectations
			tt.mockSetup(mockSubDataSvc)

			// Create a new request
			req, err := http.NewRequest("GET", "/test", nil)
			assert.NoError(t, err)

			// Create a new response recorder
			rr := httptest.NewRecorder()

			// Set up the context with user information
			ctx := context.WithValue(req.Context(), contextKeyUserID, tt.userID)
			ctx = context.WithValue(ctx, contextKeyUserType, tt.userType)
			req = req.WithContext(ctx)

			// Call the middleware
			handler := s.subscribedUserOnly(dummyHandler)
			handler.ServeHTTP(rr, req)

			// Check the status code
			assert.Equal(t, tt.expectedStatus, rr.Code)

			// Clear mock expectations for the next test
			mockSubDataSvc.AssertExpectations(t)
		})
	}
}

func TestSubscribedEntityUserOnly(t *testing.T) {
	tests := []struct {
		name                 string
		userID               string
		userType             string
		entityUserID         string
		subscriptionsEnabled bool
		expectedStatus       int
		mockSetup            func(mockSubDataSvc *MockSubscriptionDataService)
	}{
		{
			name:                 "Admin User",
			userID:               "123e4567-e89b-12d3-a456-426614174000",
			userType:             thunderdome.AdminUserType,
			entityUserID:         "223e4567-e89b-12d3-a456-426614174000",
			subscriptionsEnabled: true,
			expectedStatus:       http.StatusOK,
			mockSetup:            func(mockSubDataSvc *MockSubscriptionDataService) {},
		},
		{
			name:                 "Matching User ID",
			userID:               "323e4567-e89b-12d3-a456-426614174000",
			userType:             thunderdome.RegisteredUserType,
			entityUserID:         "323e4567-e89b-12d3-a456-426614174000",
			subscriptionsEnabled: true,
			expectedStatus:       http.StatusOK,
			mockSetup: func(mockSubDataSvc *MockSubscriptionDataService) {
				mockSubDataSvc.On("CheckActiveSubscriber", mock.Anything, "323e4567-e89b-12d3-a456-426614174000").Return(nil)
			},
		},
		{
			name:                 "Non-Matching User ID",
			userID:               "423e4567-e89b-12d3-a456-426614174000",
			userType:             thunderdome.RegisteredUserType,
			entityUserID:         "523e4567-e89b-12d3-a456-426614174000",
			subscriptionsEnabled: true,
			expectedStatus:       http.StatusForbidden,
			mockSetup:            func(mockSubDataSvc *MockSubscriptionDataService) {},
		},
		{
			name:                 "Subscriptions Disabled",
			userID:               "623e4567-e89b-12d3-a456-426614174000",
			userType:             thunderdome.RegisteredUserType,
			entityUserID:         "623e4567-e89b-12d3-a456-426614174000",
			subscriptionsEnabled: false,
			expectedStatus:       http.StatusOK,
			mockSetup:            func(mockSubDataSvc *MockSubscriptionDataService) {},
		},
		{
			name:                 "Unsubscribed User",
			userID:               "723e4567-e89b-12d3-a456-426614174000",
			userType:             thunderdome.RegisteredUserType,
			entityUserID:         "723e4567-e89b-12d3-a456-426614174000",
			subscriptionsEnabled: true,
			expectedStatus:       http.StatusForbidden,
			mockSetup: func(mockSubDataSvc *MockSubscriptionDataService) {
				mockSubDataSvc.On("CheckActiveSubscriber", mock.Anything, "723e4567-e89b-12d3-a456-426614174000").Return(errors.New("not subscribed"))
			},
		},
		{
			name:                 "Invalid Entity User ID",
			userID:               "823e4567-e89b-12d3-a456-426614174000",
			userType:             thunderdome.RegisteredUserType,
			entityUserID:         "invalid-user-id",
			subscriptionsEnabled: true,
			expectedStatus:       http.StatusBadRequest,
			mockSetup:            func(mockSubDataSvc *MockSubscriptionDataService) {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock SubscriptionDataService
			mockSubDataSvc := new(MockSubscriptionDataService)

			// Create a new service with the mock and config
			s := &Service{
				SubscriptionDataSvc: mockSubDataSvc,
				Config: &Config{
					SubscriptionsEnabled: tt.subscriptionsEnabled,
				},
			}

			// Define a dummy handler for testing
			dummyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			// Setup mock expectations
			tt.mockSetup(mockSubDataSvc)

			// Create a new request
			req, err := http.NewRequest("GET", "/users/"+tt.entityUserID, nil)
			assert.NoError(t, err)

			// Create a new response recorder
			rr := httptest.NewRecorder()

			// Set up the context with user information
			ctx := context.WithValue(req.Context(), contextKeyUserID, tt.userID)
			ctx = context.WithValue(ctx, contextKeyUserType, tt.userType)
			req = req.WithContext(ctx)

			// Set up router with vars
			router := mux.NewRouter()
			router.HandleFunc("/users/{userId}", s.subscribedEntityUserOnly(dummyHandler))

			// Serve the request
			router.ServeHTTP(rr, req)

			// Check the status code
			assert.Equal(t, tt.expectedStatus, rr.Code)

			// Clear mock expectations for the next test
			mockSubDataSvc.AssertExpectations(t)
		})
	}
}

// MockSubscriptionDataService is a mock of SubscriptionDataService
type MockSubscriptionDataService struct {
	mock.Mock
}

func (m *MockSubscriptionDataService) GetSubscriptionByID(ctx context.Context, id string) (thunderdome.Subscription, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockSubscriptionDataService) GetSubscriptionBySubscriptionID(ctx context.Context, subscriptionID string) (thunderdome.Subscription, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockSubscriptionDataService) GetActiveSubscriptionsByUserID(ctx context.Context, userID string) ([]thunderdome.Subscription, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockSubscriptionDataService) CreateSubscription(ctx context.Context, subscription thunderdome.Subscription) (thunderdome.Subscription, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockSubscriptionDataService) UpdateSubscription(ctx context.Context, id string, sub thunderdome.Subscription) (thunderdome.Subscription, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockSubscriptionDataService) GetSubscriptions(ctx context.Context, Limit int, Offset int) ([]thunderdome.Subscription, int, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockSubscriptionDataService) DeleteSubscription(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}

func (m *MockSubscriptionDataService) CheckActiveSubscriber(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func TestVerifiedUserOnly(t *testing.T) {
	tests := []struct {
		name                      string
		sessionUserID             string
		userType                  string
		entityUserID              string
		externalAPIVerifyRequired bool
		userVerified              bool
		expectedStatus            int
		mockSetup                 func(mockUserDataSvc *MockUserDataService)
	}{
		{
			name:                      "Admin User (Entity User is Verified)",
			sessionUserID:             "123e4567-e89b-12d3-a456-426614174000",
			userType:                  thunderdome.AdminUserType,
			entityUserID:              "223e4567-e89b-12d3-a456-426614174000",
			externalAPIVerifyRequired: true,
			userVerified:              true,
			expectedStatus:            http.StatusOK,
			mockSetup: func(mockUserDataSvc *MockUserDataService) {
				mockUserDataSvc.On(
					"GetUserByID",
					mock.Anything,
					"223e4567-e89b-12d3-a456-426614174000",
				).Return(&thunderdome.User{Verified: true}, nil)
			},
		},
		{
			name:                      "Admin User (Entity User is Not Verified)",
			sessionUserID:             "123e4567-e89b-12d3-a456-426614174001",
			userType:                  thunderdome.AdminUserType,
			entityUserID:              "223e4567-e89b-12d3-a456-426614174001",
			externalAPIVerifyRequired: true,
			userVerified:              false,
			expectedStatus:            http.StatusForbidden,
			mockSetup: func(mockUserDataSvc *MockUserDataService) {
				mockUserDataSvc.On(
					"GetUserByID",
					mock.Anything,
					"223e4567-e89b-12d3-a456-426614174001",
				).Return(&thunderdome.User{Verified: false}, nil)
			},
		},
		{
			name:                      "Matching Verified User",
			sessionUserID:             "323e4567-e89b-12d3-a456-426614174000",
			userType:                  thunderdome.RegisteredUserType,
			entityUserID:              "323e4567-e89b-12d3-a456-426614174000",
			externalAPIVerifyRequired: true,
			userVerified:              true,
			expectedStatus:            http.StatusOK,
			mockSetup: func(mockUserDataSvc *MockUserDataService) {
				mockUserDataSvc.On("GetUserByID", mock.Anything, "323e4567-e89b-12d3-a456-426614174000").Return(&thunderdome.User{Verified: true}, nil)
			},
		},
		{
			name:                      "Matching Unverified User",
			sessionUserID:             "423e4567-e89b-12d3-a456-426614174000",
			userType:                  thunderdome.RegisteredUserType,
			entityUserID:              "423e4567-e89b-12d3-a456-426614174000",
			externalAPIVerifyRequired: true,
			userVerified:              false,
			expectedStatus:            http.StatusForbidden,
			mockSetup: func(mockUserDataSvc *MockUserDataService) {
				mockUserDataSvc.On("GetUserByID", mock.Anything, "423e4567-e89b-12d3-a456-426614174000").Return(&thunderdome.User{Verified: false}, nil)
			},
		},
		{
			name:                      "Non-Matching User ID",
			sessionUserID:             "523e4567-e89b-12d3-a456-426614174000",
			userType:                  thunderdome.RegisteredUserType,
			entityUserID:              "623e4567-e89b-12d3-a456-426614174000",
			externalAPIVerifyRequired: true,
			userVerified:              true,
			expectedStatus:            http.StatusForbidden,
			mockSetup:                 func(mockUserDataSvc *MockUserDataService) {},
		},
		{
			name:                      "Verification Not Required",
			sessionUserID:             "723e4567-e89b-12d3-a456-426614174000",
			userType:                  thunderdome.RegisteredUserType,
			entityUserID:              "723e4567-e89b-12d3-a456-426614174000",
			externalAPIVerifyRequired: false,
			userVerified:              false,
			expectedStatus:            http.StatusOK,
			mockSetup: func(mockUserDataSvc *MockUserDataService) {
				mockUserDataSvc.On(
					"GetUserByID",
					mock.Anything, "723e4567-e89b-12d3-a456-426614174000",
				).Return(&thunderdome.User{Verified: false}, nil)
			},
		},
		{
			name:                      "Invalid Entity User ID",
			sessionUserID:             "823e4567-e89b-12d3-a456-426614174000",
			userType:                  thunderdome.RegisteredUserType,
			entityUserID:              "invalid-user-id",
			externalAPIVerifyRequired: true,
			userVerified:              true,
			expectedStatus:            http.StatusBadRequest,
			mockSetup:                 func(mockUserDataSvc *MockUserDataService) {},
		},
		{
			name:                      "User Not Found",
			sessionUserID:             "923e4567-e89b-12d3-a456-426614174000",
			userType:                  thunderdome.RegisteredUserType,
			entityUserID:              "923e4567-e89b-12d3-a456-426614174000",
			externalAPIVerifyRequired: true,
			userVerified:              true,
			expectedStatus:            http.StatusInternalServerError,
			mockSetup: func(mockUserDataSvc *MockUserDataService) {
				mockUserDataSvc.On("GetUserByID", mock.Anything, "923e4567-e89b-12d3-a456-426614174000").Return(nil, errors.New("user not found"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock UserDataService
			mockUserDataSvc := new(MockUserDataService)

			// Create a new service with the mock and config
			s := &Service{
				UserDataSvc: mockUserDataSvc,
				Config: &Config{
					ExternalAPIVerifyRequired: tt.externalAPIVerifyRequired,
				},
				Logger: otelzap.New(zap.NewNop()),
			}

			// Define a dummy handler for testing
			dummyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			// Setup mock expectations
			tt.mockSetup(mockUserDataSvc)

			// Create a new request
			req, err := http.NewRequest("GET", "/users/"+tt.entityUserID, nil)
			assert.NoError(t, err)

			// Create a new response recorder
			rr := httptest.NewRecorder()

			// Set up the context with user information
			ctx := context.WithValue(req.Context(), contextKeyUserID, tt.sessionUserID)
			ctx = context.WithValue(ctx, contextKeyUserType, tt.userType)
			req = req.WithContext(ctx)

			// Set up router with vars
			router := mux.NewRouter()
			router.HandleFunc("/users/{userId}", s.verifiedUserOnly(dummyHandler))

			// Serve the request
			router.ServeHTTP(rr, req)

			// Check the status code
			assert.Equal(t, tt.expectedStatus, rr.Code)

			// Clear mock expectations for the next test
			mockUserDataSvc.AssertExpectations(t)
		})
	}
}

// MockUserDataService is a mock of UserDataService
type MockUserDataService struct {
	mock.Mock
}

func (m *MockUserDataService) GetGuestUserByID(ctx context.Context, UserID string) (*thunderdome.User, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockUserDataService) GetUserByEmail(ctx context.Context, UserEmail string) (*thunderdome.User, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockUserDataService) GetRegisteredUsers(ctx context.Context, Limit int, Offset int) ([]*thunderdome.User, int, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockUserDataService) SearchRegisteredUsersByEmail(ctx context.Context, Email string, Limit int, Offset int) ([]*thunderdome.User, int, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockUserDataService) CreateUser(ctx context.Context, UserName string, UserEmail string, UserPassword string) (NewUser *thunderdome.User, VerifyID string, RegisterErr error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockUserDataService) CreateUserGuest(ctx context.Context, UserName string) (*thunderdome.User, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockUserDataService) CreateUserRegistered(ctx context.Context, UserName string, UserEmail string, UserPassword string, ActiveUserID string) (NewUser *thunderdome.User, VerifyID string, RegisterErr error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockUserDataService) UpdateUserAccount(ctx context.Context, UserID string, UserName string, UserEmail string, UserAvatar string, NotificationsEnabled bool, Country string, Locale string, Company string, JobTitle string, Theme string) error {
	//TODO implement me
	panic("implement me")
}

func (m *MockUserDataService) UpdateUserProfile(ctx context.Context, UserID string, UserName string, UserAvatar string, NotificationsEnabled bool, Country string, Locale string, Company string, JobTitle string, Theme string) error {
	//TODO implement me
	panic("implement me")
}

func (m *MockUserDataService) UpdateUserProfileLdap(ctx context.Context, UserID string, UserAvatar string, NotificationsEnabled bool, Country string, Locale string, Company string, JobTitle string, Theme string) error {
	//TODO implement me
	panic("implement me")
}

func (m *MockUserDataService) PromoteUser(ctx context.Context, UserID string) error {
	//TODO implement me
	panic("implement me")
}

func (m *MockUserDataService) DemoteUser(ctx context.Context, UserID string) error {
	//TODO implement me
	panic("implement me")
}

func (m *MockUserDataService) DisableUser(ctx context.Context, UserID string) error {
	//TODO implement me
	panic("implement me")
}

func (m *MockUserDataService) EnableUser(ctx context.Context, UserID string) error {
	//TODO implement me
	panic("implement me")
}

func (m *MockUserDataService) DeleteUser(ctx context.Context, UserID string) error {
	//TODO implement me
	panic("implement me")
}

func (m *MockUserDataService) CleanGuests(ctx context.Context, DaysOld int) error {
	//TODO implement me
	panic("implement me")
}

func (m *MockUserDataService) GetActiveCountries(ctx context.Context) ([]string, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockUserDataService) GetUserCredentialByUserID(ctx context.Context, UserID string) (*thunderdome.Credential, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockUserDataService) GetUserByID(ctx context.Context, userID string) (*thunderdome.User, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*thunderdome.User), args.Error(1)
}

func TestAdminOnly(t *testing.T) {
	tests := []struct {
		name           string
		userType       string
		expectedStatus int
	}{
		{
			name:           "Admin User",
			userType:       thunderdome.AdminUserType,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Regular User",
			userType:       thunderdome.RegisteredUserType,
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "Guest User",
			userType:       "GUEST",
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "Empty User Type",
			userType:       "",
			expectedStatus: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new service
			s := &Service{
				// Add any necessary service configuration here
			}

			// Define a dummy handler for testing
			dummyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			// Create a new request
			req, err := http.NewRequest("GET", "/admin-only", nil)
			assert.NoError(t, err)

			// Create a new response recorder
			rr := httptest.NewRecorder()

			// Set up the context with user type
			ctx := context.WithValue(req.Context(), contextKeyUserType, tt.userType)
			req = req.WithContext(ctx)

			// Call the middleware
			handler := s.adminOnly(dummyHandler)
			handler.ServeHTTP(rr, req)

			// Check the status code
			assert.Equal(t, tt.expectedStatus, rr.Code)

			// If the status is Forbidden, check for the correct error message
			if tt.expectedStatus == http.StatusForbidden {
				var responseBody map[string]interface{}
				err = json.NewDecoder(rr.Body).Decode(&responseBody)
				assert.NoError(t, err)
				assert.Equal(t, "REQUIRES_ADMIN", responseBody["error"])
			}
		})
	}
}

func TestRegisteredUserOnly(t *testing.T) {
	tests := []struct {
		name           string
		userType       string
		expectedStatus int
	}{
		{
			name:           "Admin User",
			userType:       thunderdome.AdminUserType,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Regular User",
			userType:       thunderdome.RegisteredUserType,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Guest User",
			userType:       thunderdome.GuestUserType,
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "Empty User Type",
			userType:       "",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new service
			s := &Service{
				// Add any necessary service configuration here
			}

			// Define a dummy handler for testing
			dummyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			// Create a new request
			req, err := http.NewRequest("GET", "/registered-only", nil)
			assert.NoError(t, err)

			// Create a new response recorder
			rr := httptest.NewRecorder()

			// Set up the context with user type
			ctx := context.WithValue(req.Context(), contextKeyUserType, tt.userType)
			req = req.WithContext(ctx)

			// Call the middleware
			handler := s.registeredUserOnly(dummyHandler)
			handler.ServeHTTP(rr, req)

			// Check the status code
			assert.Equal(t, tt.expectedStatus, rr.Code)

			// If the status is Forbidden, check for the correct error message
			if tt.expectedStatus == http.StatusForbidden {
				var responseBody map[string]interface{}
				err = json.NewDecoder(rr.Body).Decode(&responseBody)
				assert.NoError(t, err)
				assert.Equal(t, "REGISTERED_USER_ONLY", responseBody["error"])
			}
		})
	}
}

func TestEntityUserOnly(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		userType       string
		entityUserID   string
		expectedStatus int
	}{
		{
			name:           "Matching User ID",
			userID:         "123e4567-e89b-12d3-a456-426614174000",
			userType:       "REGULAR",
			entityUserID:   "123e4567-e89b-12d3-a456-426614174000",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Admin User with Different Entity ID",
			userID:         "223e4567-e89b-12d3-a456-426614174000",
			userType:       thunderdome.AdminUserType,
			entityUserID:   "323e4567-e89b-12d3-a456-426614174000",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Non-Matching User ID",
			userID:         "423e4567-e89b-12d3-a456-426614174000",
			userType:       "REGULAR",
			entityUserID:   "523e4567-e89b-12d3-a456-426614174000",
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "Invalid Entity User ID",
			userID:         "623e4567-e89b-12d3-a456-426614174000",
			userType:       "REGULAR",
			entityUserID:   "invalid-user-id",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new service
			s := &Service{
				// Add any necessary service configuration here
			}

			// Define a dummy handler for testing
			dummyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			// Create a new request
			req, err := http.NewRequest("GET", "/users/"+tt.entityUserID, nil)
			assert.NoError(t, err)

			// Create a new response recorder
			rr := httptest.NewRecorder()

			// Set up the context with user information
			ctx := context.WithValue(req.Context(), contextKeyUserID, tt.userID)
			ctx = context.WithValue(ctx, contextKeyUserType, tt.userType)
			req = req.WithContext(ctx)

			// Set up router with vars
			router := mux.NewRouter()
			router.HandleFunc("/users/{userId}", s.entityUserOnly(dummyHandler))

			// Serve the request
			router.ServeHTTP(rr, req)

			// Check the status code
			assert.Equal(t, tt.expectedStatus, rr.Code)

			// If the status is not OK, check for the correct error message
			if tt.expectedStatus != http.StatusOK {
				var responseBody map[string]interface{}
				err = json.NewDecoder(rr.Body).Decode(&responseBody)
				assert.NoError(t, err)

				if tt.expectedStatus == http.StatusForbidden {
					assert.Equal(t, "INVALID_USER", responseBody["error"])
				}
			}
		})
	}
}

func TestPanicRecovery(t *testing.T) {
	tests := []struct {
		name           string
		handler        http.HandlerFunc
		expectedStatus int
		expectPanic    bool
	}{
		{
			name: "No Panic",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}),
			expectedStatus: http.StatusOK,
			expectPanic:    false,
		},
		{
			name: "Panic Occurs",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				panic("test panic")
			}),
			expectedStatus: http.StatusInternalServerError,
			expectPanic:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			observedZapCore, observedLogs := observer.New(zap.InfoLevel)
			observedLogger := zap.New(observedZapCore)

			s := &Service{
				Logger: otelzap.New(observedLogger),
			}

			// Create a new request
			req, err := http.NewRequest("GET", "/test", nil)
			assert.NoError(t, err)

			// Create a new response recorder
			rr := httptest.NewRecorder()

			// Apply the panicRecovery middleware to the test handler
			handler := s.panicRecovery(tt.handler)

			// Serve the request
			handler.ServeHTTP(rr, req)

			// Check the status code
			assert.Equal(t, tt.expectedStatus, rr.Code)

			if tt.expectPanic {
				require.Equal(t, 1, observedLogs.Len())
				allLogs := observedLogs.All()
				assert.Equal(t, zap.ErrorLevel, allLogs[0].Level)
				assert.Equal(t, "http handler recovering from panic error: test panic", allLogs[0].Message)
			} else {
				require.Equal(t, 0, observedLogs.Len())
			}
		})
	}
}
