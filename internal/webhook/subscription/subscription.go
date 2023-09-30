// Package subscription provides stripe subscription webhook functionality
package subscription

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"

	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/webhook"
)

type Config struct {
	AccountSecret string
	WebhookSecret string
}

type Service struct {
	config  Config
	logger  *otelzap.Logger
	dataSvc thunderdome.SubscriptionDataSvc
}

func New(config Config, logger *otelzap.Logger, dataSvc thunderdome.SubscriptionDataSvc) *Service {
	// The library needs to be configured with your account's secret key.
	// Ensure the key is kept out of any version control system you might be using.
	stripe.Key = config.AccountSecret
	return &Service{
		logger:  logger,
		config:  config,
		dataSvc: dataSvc,
	}
}

func (s *Service) HandleWebhook() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		const MaxBodyBytes = int64(65536)
		ctx := req.Context()
		logger := s.logger.Ctx(ctx)
		logger.Info("Stripe webhook request received")
		req.Body = http.MaxBytesReader(w, req.Body, MaxBodyBytes)
		payload, err := io.ReadAll(req.Body)
		if err != nil {
			logger.Error(fmt.Sprintf("Error reading request body: %v", err), zap.String("payload", string(payload)))
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		// Pass the request body and Stripe-Signature header to ConstructEvent, along with the webhook signing key.
		event, err := webhook.ConstructEvent(payload, req.Header.Get("Stripe-Signature"), s.config.WebhookSecret)
		if err != nil {
			logger.Error(fmt.Sprintf("Error verifying webhook signature: %v", err), zap.String("eventId", event.ID))
			w.WriteHeader(http.StatusBadRequest) // Return a 400 error on a bad signature
			return
		}

		// Unmarshal the event data into an appropriate struct depending on its Type
		logger.Info(fmt.Sprintf("Processing Stripe webhook event type: %s", event.Type), zap.String("eventId", event.ID))
		switch event.Type {
		case "checkout.session.completed":
			// store stripe customer id for future events
			clientReferenceId, ok := event.Data.Object["client_reference_id"]
			if !ok || clientReferenceId == nil {
				logger.Error("Error getting client_reference_id from event", zap.String("eventId", event.ID))
				w.WriteHeader(http.StatusBadRequest) // Return a 400 error on a bad signature
				return
			}
			customerId, ok := event.Data.Object["customer"]
			if !ok || customerId == nil {
				logger.Error("Error getting customer from event", zap.String("eventId", event.ID))
				w.WriteHeader(http.StatusBadRequest) // Return a 400 error on a bad signature
				return
			}
			expires := time.Now().Add(time.Hour * 24 * 33) // start with 33 day subscription
			_, err = s.dataSvc.CreateSubscription(ctx, clientReferenceId.(string), customerId.(string), expires)
			if err != nil {
				logger.Error(fmt.Sprintf("Error creating subscription: %v", err), zap.String("eventId", event.ID))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		case "customer.subscription.updated":
			customerId, ok := event.Data.Object["customer"]
			if !ok || customerId == nil {
				logger.Error("Error getting customer from event", zap.String("eventId", event.ID))
				w.WriteHeader(http.StatusBadRequest) // Return a 400 error on a bad signature
				return
			}
			subscription, err := s.dataSvc.GetSubscriptionByCustomerID(ctx, customerId.(string))
			if err != nil {
				logger.Error(fmt.Sprintf("Error getting customer id %s subscription: %v", customerId.(string), err), zap.String("eventId", event.ID))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			subStatus, ok := event.Data.Object["status"]
			if !ok || subStatus == nil {
				logger.Error("Error getting subscription status from event", zap.String("eventId", event.ID))
				w.WriteHeader(http.StatusBadRequest) // Return a 400 error on a bad signature
				return
			}
			periodEnd, ok := event.Data.Object["current_period_end"]
			if !ok || periodEnd == nil {
				logger.Error("Error getting subscription current_period_end from event", zap.String("eventId", event.ID))
				w.WriteHeader(http.StatusBadRequest) // Return a 400 error on a bad signature
				return
			}
			expires := time.Unix(int64(periodEnd.(float64)), 0)
			active := subStatus == "active"
			_, err = s.dataSvc.UpdateSubscription(ctx, subscription.ID, active, expires)
			if err != nil {
				logger.Error(fmt.Sprintf("Error creating subscription: %v", err), zap.String("eventId", event.ID))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		default:
			logger.Error(fmt.Sprintf("Unhandled Stripe webhook event type: %s", event.Type), zap.String("eventId", event.ID))
		}

		logger.Info(fmt.Sprintf("Successfully processed Stripe webhook event type: %s", event.Type), zap.String("eventId", event.ID))

		w.WriteHeader(http.StatusOK)
	}
}
