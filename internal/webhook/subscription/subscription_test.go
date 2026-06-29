package subscription

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/webhook"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
)

type stubSubscriptionDataSvc struct {
	getSubscriptionBySubscriptionID func(ctx context.Context, subscriptionID string) (thunderdome.Subscription, error)
	updateSubscription              func(ctx context.Context, subscriptionID string, subscription thunderdome.Subscription) (thunderdome.Subscription, error)
}

func (s stubSubscriptionDataSvc) GetSubscriptionBySubscriptionID(ctx context.Context, subscriptionID string) (thunderdome.Subscription, error) {
	if s.getSubscriptionBySubscriptionID != nil {
		return s.getSubscriptionBySubscriptionID(ctx, subscriptionID)
	}

	return thunderdome.Subscription{}, nil
}

func (s stubSubscriptionDataSvc) CreateSubscription(ctx context.Context, subscription thunderdome.Subscription) (thunderdome.Subscription, error) {
	return thunderdome.Subscription{}, nil
}

func (s stubSubscriptionDataSvc) UpdateSubscription(ctx context.Context, subscriptionID string, subscription thunderdome.Subscription) (thunderdome.Subscription, error) {
	if s.updateSubscription != nil {
		return s.updateSubscription(ctx, subscriptionID, subscription)
	}

	return subscription, nil
}

type stubUserDataSvc struct{}

func (s stubUserDataSvc) GetUserByID(ctx context.Context, userID string) (*thunderdome.User, error) {
	return nil, nil
}

type stubEmailSvc struct{}

func (s stubEmailSvc) SendUserSubscriptionActive(userName string, userEmail string, subscriptionType string) error {
	return nil
}

func (s stubEmailSvc) SendUserSubscriptionDeactivated(userName string, userEmail string, subscriptionType string) error {
	return nil
}

func TestHandleWebhookSubscriptionUpdatedIgnoresUntrackedIncompleteStatuses(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		status string
	}{
		{name: "incomplete", status: "incomplete"},
		{name: "incomplete expired", status: "incomplete_expired"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			dataSvc := stubSubscriptionDataSvc{
				getSubscriptionBySubscriptionID: func(ctx context.Context, subscriptionID string) (thunderdome.Subscription, error) {
					return thunderdome.Subscription{}, errors.New("no subscription sub_123")
				},
				updateSubscription: func(ctx context.Context, subscriptionID string, subscription thunderdome.Subscription) (thunderdome.Subscription, error) {
					t.Fatalf("UpdateSubscription should not be called for %s", tt.status)
					return thunderdome.Subscription{}, nil
				},
			}

			service := New(
				Config{WebhookSecret: "whsec_test"},
				otelzap.New(zap.NewNop()),
				dataSvc,
				stubEmailSvc{},
				stubUserDataSvc{},
			)

			req := newSubscriptionUpdatedRequest(t, "whsec_test", "evt_incomplete", "sub_123", tt.status)
			rr := httptest.NewRecorder()

			service.HandleWebhook().ServeHTTP(rr, req)

			if rr.Code != http.StatusOK {
				t.Fatalf("expected %d, got %d", http.StatusOK, rr.Code)
			}
		})
	}
}

func TestHandleWebhookSubscriptionUpdatedReturnsServerErrorForTrackedStatuses(t *testing.T) {
	t.Parallel()

	service := New(
		Config{WebhookSecret: "whsec_test"},
		otelzap.New(zap.NewNop()),
		stubSubscriptionDataSvc{
			getSubscriptionBySubscriptionID: func(ctx context.Context, subscriptionID string) (thunderdome.Subscription, error) {
				return thunderdome.Subscription{}, errors.New("database unavailable")
			},
		},
		stubEmailSvc{},
		stubUserDataSvc{},
	)

	req := newSubscriptionUpdatedRequest(t, "whsec_test", "evt_active", "sub_456", "active")
	rr := httptest.NewRecorder()

	service.HandleWebhook().ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected %d, got %d", http.StatusInternalServerError, rr.Code)
	}
}

func newSubscriptionUpdatedRequest(t *testing.T, secret string, eventID string, subscriptionID string, status string) *http.Request {
	t.Helper()

	payload, err := json.Marshal(map[string]any{
		"id":               eventID,
		"object":           "event",
		"api_version":      stripe.APIVersion,
		"created":          time.Now().Unix(),
		"livemode":         false,
		"pending_webhooks": 1,
		"type":             "customer.subscription.updated",
		"data": map[string]any{
			"object": map[string]any{
				"id":                 subscriptionID,
				"object":             "subscription",
				"status":             status,
				"current_period_end": time.Now().Unix(),
			},
		},
	})
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}

	signedPayload := webhook.GenerateTestSignedPayload(&webhook.UnsignedPayload{
		Payload:   payload,
		Secret:    secret,
		Timestamp: time.Now(),
	})

	req := httptest.NewRequest(http.MethodPost, "/webhooks/subscriptions", bytes.NewReader(signedPayload.Payload))
	req.Header.Set("Stripe-Signature", signedPayload.Header)

	return req
}
