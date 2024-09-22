package http

import (
	"context"
	"fmt"
	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

// MockTeamDataSvc is a mock implementation of the TeamDataSvc
type MockTeamDataSvc struct {
	mock.Mock
}

func (m *MockTeamDataSvc) TeamUserRole(ctx context.Context, UserID string, TeamID string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockTeamDataSvc) TeamGet(ctx context.Context, TeamID string) (*thunderdome.Team, error) {
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

func (m *MockTeamDataSvc) TeamUpdate(ctx context.Context, TeamId string, TeamName string) (*thunderdome.Team, error) {
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

func (m *MockTeamDataSvc) TeamGetUserInvites(ctx context.Context, teamId string) ([]thunderdome.TeamUserInvite, error) {
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

func (m *MockTeamDataSvc) TeamIsSubscribed(ctx context.Context, TeamID string) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockTeamDataSvc) GetTeamMetrics(ctx context.Context, teamID string) (*thunderdome.TeamMetrics, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockTeamDataSvc) TeamUserRoles(ctx context.Context, userID, teamID string) (*thunderdome.UserTeamRoleInfo, error) {
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
			userType: "REGISTERED",
			teamID:   "128ee064-62ca-43b2-9fca-9c1089c89bd2",
			setupMocks: func(mtds *MockTeamDataSvc, ml *MockLogger) {
				mtds.On(
					"TeamUserRoles",
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
					"TeamUserRoles",
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
			userType: "REGISTERED",
			teamID:   "7d4d0a17-cb20-4499-bc7f-ecaf2e77c15a",
			setupMocks: func(mtds *MockTeamDataSvc, ml *MockLogger) {
				deptRole := thunderdome.AdminUserType
				mtds.On(
					"TeamUserRoles",
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
			userType: "REGISTERED",
			teamID:   "0ea230df-b5fe-47ae-a473-5153004eebdd",
			setupMocks: func(mtds *MockTeamDataSvc, ml *MockLogger) {
				orgRole := thunderdome.AdminUserType
				mtds.On(
					"TeamUserRoles",
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
			userType: "REGISTERED",
			teamID:   "invalid-team-id",
			setupMocks: func(mtds *MockTeamDataSvc, ml *MockLogger) {
				ml.On("Ctx", mock.Anything).Return(zap.NewNop())
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:     "Team not found",
			userID:   "1b853ef6-2c28-4c8e-ac29-a9d4827774fe",
			userType: "REGISTERED",
			teamID:   "e52b9251-4722-4b4d-8f97-204ce7e51eec",
			setupMocks: func(mtds *MockTeamDataSvc, ml *MockLogger) {
				mtds.On(
					"TeamUserRoles",
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
			userType: "REGISTERED",
			teamID:   "a805def1-e1fa-42a9-b5f6-ee338799fa77",
			setupMocks: func(mtds *MockTeamDataSvc, ml *MockLogger) {
				mtds.On("TeamUserRoles",
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
