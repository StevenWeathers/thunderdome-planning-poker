// Package subscription provides stripe subscription webhook functionality
package subscription

import (
	"fmt"
	"io"
	"net/http"
	"time"

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
		ctx := req.Context()
		logger := s.logger.Ctx(ctx)
		const MaxBodyBytes = int64(65536)
		req.Body = http.MaxBytesReader(w, req.Body, MaxBodyBytes)
		payload, err := io.ReadAll(req.Body)
		if err != nil {
			logger.Error(fmt.Sprintf("Error reading request body: %v", err))
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		// This is your Stripe CLI webhook secret for testing your endpoint locally.
		endpointSecret := s.config.WebhookSecret
		// Pass the request body and Stripe-Signature header to ConstructEvent, along
		// with the webhook signing key.
		event, err := webhook.ConstructEvent(payload, req.Header.Get("Stripe-Signature"),
			endpointSecret)

		if err != nil {
			logger.Error(fmt.Sprintf("Error verifying webhook signature: %v", err))
			w.WriteHeader(http.StatusBadRequest) // Return a 400 error on a bad signature
			return
		}

		// Unmarshal the event data into an appropriate struct depending on its Type
		switch event.Type {
		case "checkout.session.completed":
			// store stripe customer id for future events
			clientReferenceId, ok := event.Data.Object["client_reference_id"]
			if !ok || clientReferenceId == nil {
				logger.Error("Error getting client_reference_id from event")
				w.WriteHeader(http.StatusBadRequest) // Return a 400 error on a bad signature
				return
			}
			customerId, ok := event.Data.Object["customer"]
			if !ok || customerId == nil {
				logger.Error("Error getting customer from event")
				w.WriteHeader(http.StatusBadRequest) // Return a 400 error on a bad signature
				return
			}
			expires := time.Now().Add(time.Hour * 24 * 33) // start with 33 day subscription
			_, err = s.dataSvc.CreateSubscription(ctx, clientReferenceId.(string), customerId.(string), expires)
			if err != nil {
				logger.Error(fmt.Sprintf("Error creating subscription: %v", err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		case "customer.subscription.updated":
			customerId, ok := event.Data.Object["customer"]
			if !ok || customerId == nil {
				logger.Error("Error getting customer from event")
				w.WriteHeader(http.StatusBadRequest) // Return a 400 error on a bad signature
				return
			}
			subscription, err := s.dataSvc.GetSubscriptionByCustomerID(ctx, customerId.(string))
			if err != nil {
				logger.Error(fmt.Sprintf("Error getting customer id %s subscription: %v", customerId.(string), err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			subStatus, ok := event.Data.Object["status"]
			if !ok || subStatus == nil {
				logger.Error("Error getting subscription status from event")
				w.WriteHeader(http.StatusBadRequest) // Return a 400 error on a bad signature
				return
			}
			periodEnd, ok := event.Data.Object["current_period_end"]
			if !ok || periodEnd == nil {
				logger.Error("Error getting subscription current_period_end from event")
				w.WriteHeader(http.StatusBadRequest) // Return a 400 error on a bad signature
				return
			}
			expires := time.Unix(int64(periodEnd.(float64)), 0)
			active := subStatus == "active"
			_, err = s.dataSvc.UpdateSubscription(ctx, subscription.ID, active, expires)
			if err != nil {
				logger.Error(fmt.Sprintf("Error creating subscription: %v", err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		default:
			logger.Error(fmt.Sprintf("Unhandled event type: %s\n", event.Type))
		}

		w.WriteHeader(http.StatusOK)
	}
}
