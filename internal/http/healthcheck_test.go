package http

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
)

const (
	hcTeamID     = "11111111-1111-1111-1111-111111111111"
	hcID         = "22222222-2222-2222-2222-222222222222"
	hcCreatorID  = "33333333-3333-3333-3333-333333333333"
	hcMemberID   = "44444444-4444-4444-4444-444444444444"
	hcTemplateID = "55555555-5555-5555-5555-555555555555"
)

type healthCheckStub struct {
	HealthCheckDataSvc
	check   thunderdome.HealthCheck
	answers []*thunderdome.HealthCheckAnswer
	calls   []string
	err     error
}

func (s *healthCheckStub) HealthCheckGet(context.Context, string, string) (*thunderdome.HealthCheck, error) {
	hc := s.check
	return &hc, nil
}

func (s *healthCheckStub) HealthCheckUserResponse(context.Context, string, string) ([]*thunderdome.HealthCheckAnswer, error) {
	return []*thunderdome.HealthCheckAnswer{}, nil
}

func (s *healthCheckStub) HealthCheckPendingUsers(context.Context, string, string) ([]*thunderdome.HealthCheckParticipant, error) {
	return []*thunderdome.HealthCheckParticipant{{ID: hcMemberID, Name: "Pending Member"}}, nil
}

func (s *healthCheckStub) HealthCheckResponses(context.Context, string) ([][]*thunderdome.HealthCheckAnswer, error) {
	s.calls = append(s.calls, "responses")
	return [][]*thunderdome.HealthCheckAnswer{
		{{Area: "Fun", Rating: "green", Trend: "stable"}},
	}, nil
}

func (s *healthCheckStub) HealthCheckCreate(context.Context, string, string, string, time.Time, string) (string, error) {
	s.calls = append(s.calls, "create")
	return hcID, s.err
}

func (s *healthCheckStub) HealthCheckClose(context.Context, string, string) error {
	s.calls = append(s.calls, "close")
	return s.err
}

func (s *healthCheckStub) HealthCheckResponseSubmit(_ context.Context, _ string, _ string, _ string, answers []*thunderdome.HealthCheckAnswer) error {
	s.calls = append(s.calls, "submit")
	s.answers = answers
	return s.err
}

func newHealthCheckService(stub *healthCheckStub) *Service {
	validate = validator.New()

	return &Service{
		Config:             &Config{},
		Logger:             otelzap.New(zap.NewNop()),
		HealthCheckDataSvc: stub,
	}
}

func healthCheckRequest(method string, body string, userID string, roles *thunderdome.UserTeamRoleInfo) *http.Request {
	req := httptest.NewRequest(method, "/", strings.NewReader(body))
	req.SetPathValue("teamId", hcTeamID)
	req.SetPathValue("healthCheckId", hcID)

	ctx := context.WithValue(req.Context(), contextKeyUserID, userID)
	ctx = context.WithValue(ctx, contextKeyUserType, thunderdome.RegisteredUserType)
	ctx = context.WithValue(ctx, contextKeyUserTeamRoles, roles)

	return req.WithContext(ctx)
}

func teamRoles(team string, org string) *thunderdome.UserTeamRoleInfo {
	roles := &thunderdome.UserTeamRoleInfo{}
	if team != "" {
		roles.TeamRole = &team
	}
	if org != "" {
		roles.OrganizationRole = &org
	}

	return roles
}

func TestCanManageHealthCheck(t *testing.T) {
	hc := &thunderdome.HealthCheck{CreatedBy: hcCreatorID}

	tests := []struct {
		name     string
		userID   string
		userType string
		roles    *thunderdome.UserTeamRoleInfo
		want     bool
	}{
		{"creator", hcCreatorID, thunderdome.RegisteredUserType, teamRoles("MEMBER", ""), true},
		{"creator who left the team", hcCreatorID, thunderdome.RegisteredUserType, teamRoles("", ""), false},
		{"another member", hcMemberID, thunderdome.RegisteredUserType, teamRoles("MEMBER", ""), false},
		{"team admin", hcMemberID, thunderdome.RegisteredUserType, teamRoles("ADMIN", ""), true},
		{"org admin on the team", hcMemberID, thunderdome.RegisteredUserType, teamRoles("MEMBER", "ADMIN"), true},
		{"org admin not on the team", hcMemberID, thunderdome.RegisteredUserType, teamRoles("", "ADMIN"), false},
		{"application admin", hcMemberID, thunderdome.AdminUserType, teamRoles("", ""), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, canManageHealthCheck(tt.userID, tt.userType, tt.roles, hc))
		})
	}
}

func TestHandleGetTeamHealthCheck(t *testing.T) {
	areas := []*thunderdome.HealthCheckArea{{Name: "Fun"}}

	tests := []struct {
		name        string
		closed      bool
		responses   int
		roles       *thunderdome.UserTeamRoleInfo
		userID      string
		wantResults bool
		wantPending bool
		wantRespond bool
	}{
		{name: "open hides results", responses: 5, roles: teamRoles("MEMBER", ""), userID: hcCreatorID, wantPending: true, wantRespond: true},
		{name: "closed with too few responses hides results", closed: true, responses: thunderdome.HealthCheckMinResponses - 1, roles: teamRoles("MEMBER", ""), userID: hcMemberID, wantRespond: true},
		{name: "closed with enough responses shows results", closed: true, responses: thunderdome.HealthCheckMinResponses, roles: teamRoles("MEMBER", ""), userID: hcMemberID, wantResults: true, wantRespond: true},
		{name: "org admin only views", closed: true, responses: thunderdome.HealthCheckMinResponses, roles: teamRoles("", "ADMIN"), userID: hcMemberID, wantResults: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stub := &healthCheckStub{check: thunderdome.HealthCheck{
				ID: hcID, TeamID: hcTeamID, Areas: areas, CreatedBy: hcCreatorID,
				Closed: tt.closed, ResponseCount: tt.responses,
			}}
			s := newHealthCheckService(stub)
			w := httptest.NewRecorder()

			s.handleGetTeamHealthCheck()(w, healthCheckRequest(http.MethodGet, "", tt.userID, tt.roles))

			require.Equal(t, http.StatusOK, w.Code)
			var body struct {
				Data healthCheckDetail `json:"data"`
			}
			require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
			assert.Equal(t, tt.wantResults, body.Data.Results != nil)
			assert.Equal(t, tt.wantPending, body.Data.Pending != nil)
			assert.Equal(t, tt.wantRespond, body.Data.CanRespond)
			// individual responses aren't even read unless results are shown
			assert.Equal(t, tt.wantResults, len(stub.calls) == 1)
		})
	}
}

func TestHandleTeamHealthCheckRespond(t *testing.T) {
	areas := []*thunderdome.HealthCheckArea{{Name: "Fun"}, {Name: "Speed"}}
	fun := `{"area":"Fun","rating":"green","trend":"stable","comment":"  great  "}`
	speed := `{"area":"Speed","rating":"red","trend":"declining"}`

	tests := []struct {
		name       string
		body       string
		roles      *thunderdome.UserTeamRoleInfo
		svcErr     error
		wantStatus int
	}{
		{name: "every area answered", body: `{"answers":[` + fun + `,` + speed + `]}`, wantStatus: http.StatusOK},
		{name: "missing area", body: `{"answers":[` + fun + `]}`, wantStatus: http.StatusBadRequest},
		{name: "duplicate area", body: `{"answers":[` + fun + `,` + fun + `]}`, wantStatus: http.StatusBadRequest},
		{name: "unknown area", body: `{"answers":[` + fun + `,{"area":"Nope","rating":"green","trend":"stable"}]}`, wantStatus: http.StatusBadRequest},
		{name: "invalid rating", body: `{"answers":[` + fun + `,{"area":"Speed","rating":"blue","trend":"stable"}]}`, wantStatus: http.StatusBadRequest},
		{name: "comment too long", body: `{"answers":[` + speed + `,{"area":"Fun","rating":"green","trend":"stable","comment":"` + strings.Repeat("a", 1001) + `"}]}`, wantStatus: http.StatusBadRequest},
		{name: "org admin not on the team", body: `{"answers":[` + fun + `,` + speed + `]}`, roles: teamRoles("", "ADMIN"), wantStatus: http.StatusForbidden},
		{name: "closed", body: `{"answers":[` + fun + `,` + speed + `]}`, svcErr: thunderdome.ErrHealthCheckClosed, wantStatus: http.StatusConflict},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stub := &healthCheckStub{check: thunderdome.HealthCheck{ID: hcID, TeamID: hcTeamID, Areas: areas}, err: tt.svcErr}
			s := newHealthCheckService(stub)
			w := httptest.NewRecorder()
			roles := tt.roles
			if roles == nil {
				roles = teamRoles("MEMBER", "")
			}

			s.handleTeamHealthCheckRespond()(w, healthCheckRequest(http.MethodPut, tt.body, hcMemberID, roles))

			assert.Equal(t, tt.wantStatus, w.Code, w.Body.String())
			if tt.wantStatus == http.StatusOK {
				require.Len(t, stub.answers, 2)
				assert.Equal(t, "great", stub.answers[0].Comment)
			}
		})
	}
}

func TestHandleTeamHealthCheckCreate(t *testing.T) {
	future := time.Now().Add(time.Hour).UTC().Format(time.RFC3339)
	past := time.Now().Add(-time.Hour).UTC().Format(time.RFC3339)
	body := func(name string, deadline string) string {
		return `{"name":"` + name + `","templateId":"` + hcTemplateID + `","deadline":"` + deadline + `"}`
	}

	tests := []struct {
		name       string
		body       string
		roles      *thunderdome.UserTeamRoleInfo
		svcErr     error
		wantStatus int
		wantCreate bool
	}{
		{name: "member", body: body("Sprint 1", future), roles: teamRoles("MEMBER", ""), wantStatus: http.StatusOK, wantCreate: true},
		{name: "org admin not on the team", body: body("Sprint 1", future), roles: teamRoles("", "ADMIN"), wantStatus: http.StatusForbidden},
		{name: "blank name", body: body("   ", future), roles: teamRoles("MEMBER", ""), wantStatus: http.StatusBadRequest},
		{name: "past deadline", body: body("Sprint 1", past), roles: teamRoles("MEMBER", ""), wantStatus: http.StatusBadRequest},
		{name: "unavailable template", body: body("Sprint 1", future), roles: teamRoles("MEMBER", ""), svcErr: thunderdome.ErrHealthCheckTemplateNotFound, wantStatus: http.StatusNotFound, wantCreate: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stub := &healthCheckStub{err: tt.svcErr}
			s := newHealthCheckService(stub)
			w := httptest.NewRecorder()

			s.handleTeamHealthCheckCreate()(w, healthCheckRequest(http.MethodPost, tt.body, hcMemberID, tt.roles))

			assert.Equal(t, tt.wantStatus, w.Code, w.Body.String())
			assert.Equal(t, tt.wantCreate, len(stub.calls) == 1)
		})
	}
}

func TestHandleTeamHealthCheckClose(t *testing.T) {
	tests := []struct {
		name       string
		userID     string
		wantStatus int
	}{
		{name: "creator", userID: hcCreatorID, wantStatus: http.StatusOK},
		{name: "another member", userID: hcMemberID, wantStatus: http.StatusForbidden},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stub := &healthCheckStub{check: thunderdome.HealthCheck{ID: hcID, TeamID: hcTeamID, CreatedBy: hcCreatorID}}
			s := newHealthCheckService(stub)
			w := httptest.NewRecorder()

			s.handleTeamHealthCheckClose()(w, healthCheckRequest(http.MethodPost, "", tt.userID, teamRoles("MEMBER", "")))

			assert.Equal(t, tt.wantStatus, w.Code)
			assert.Equal(t, tt.wantStatus == http.StatusOK, len(stub.calls) == 1)
		})
	}
}
