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

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
)

const (
	ownerUserID = "e0a1f5c2-3b7d-4a4e-9d0e-6f1a2b3c4d5e"
	otherUserID = "0b8e4c7a-91f3-4a6e-8d2b-5e7c3f1a2b44"
	childID     = "0f8fad5b-d9cb-469f-a165-70867728950e"
	targetID    = "7c9e6679-7425-40de-944b-e07fc1f90ae7"
)

type subscriptionStub struct {
	SubscriptionDataSvc
	owner   string
	updated bool
}

func (s *subscriptionStub) GetSubscriptionByID(context.Context, string) (thunderdome.Subscription, error) {
	return thunderdome.Subscription{ID: childID, UserID: s.owner, Type: "team"}, nil
}

func (s *subscriptionStub) UpdateSubscription(_ context.Context, _ string, sub thunderdome.Subscription) (thunderdome.Subscription, error) {
	s.updated = true
	return sub, nil
}

type teamRolesStub struct {
	TeamDataSvc
	level string
}

func (s *teamRolesStub) TeamUserRolesByUserID(context.Context, string, string) (*thunderdome.UserTeamRoleInfo, error) {
	return &thunderdome.UserTeamRoleInfo{AssociationLevel: s.level}, nil
}

func TestHandleEntityUserUpdateSubscriptionOwnership(t *testing.T) {
	validate = validator.New()

	tests := []struct {
		name       string
		owner      string
		level      string
		wantStatus int
	}{
		{name: "own subscription and team", owner: ownerUserID, level: "TEAM", wantStatus: http.StatusOK},
		{name: "another user's subscription", owner: otherUserID, level: "TEAM", wantStatus: http.StatusNotFound},
		{name: "team the user is not part of", owner: ownerUserID, level: "NONE", wantStatus: http.StatusForbidden},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			subs := &subscriptionStub{owner: tt.owner}
			s := &Service{
				SubscriptionDataSvc: subs,
				TeamDataSvc:         &teamRolesStub{level: tt.level},
				Logger:              otelzap.New(zap.NewNop()),
			}
			req := httptest.NewRequest(http.MethodPatch, "/", strings.NewReader(`{"team_id":"`+targetID+`"}`))
			req.SetPathValue("userId", ownerUserID)
			req.SetPathValue("subscriptionId", childID)
			rec := httptest.NewRecorder()

			s.handleEntityUserUpdateSubscription()(rec, req)

			assert.Equal(t, tt.wantStatus, rec.Code)
			assert.Equal(t, tt.wantStatus == http.StatusOK, subs.updated)
		})
	}
}

type jiraStub struct {
	JiraDataSvc
}

func (s *jiraStub) GetInstanceByID(context.Context, string) (thunderdome.JiraInstance, error) {
	return thunderdome.JiraInstance{ID: childID, UserID: otherUserID, Host: "https://example.com"}, nil
}

func TestHandleJiraStoryJQLSearchRejectsAnotherUsersInstance(t *testing.T) {
	validate = validator.New()
	s := &Service{JiraDataSvc: &jiraStub{}, Logger: otelzap.New(zap.NewNop())}

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"jql":"project = X","startAt":0,"maxResults":10}`))
	req.SetPathValue("userId", ownerUserID)
	req.SetPathValue("instanceId", childID)
	rec := httptest.NewRecorder()

	s.handleJiraStoryJQLSearch()(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.Contains(t, rec.Body.String(), "JIRA_INSTANCE_NOT_FOUND")
}
