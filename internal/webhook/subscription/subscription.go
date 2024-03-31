// Package subscription provides stripe subscription webhook functionality
package subscription

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/stripe/stripe-go/v76/checkout/session"
	"github.com/stripe/stripe-go/v76/product"

	"go.uber.org/zap"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"

	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/webhook"
)

type Config struct {
	AccountSecret string
	WebhookSecret string
}

type Service struct {
	config      Config
	logger      *otelzap.Logger
	dataSvc     thunderdome.SubscriptionDataSvc
	emailSvc    thunderdome.EmailService
	userDataSvc thunderdome.UserDataSvc
}

func New(
	config Config,
	logger *otelzap.Logger,
	dataSvc thunderdome.SubscriptionDataSvc,
	emailSvc thunderdome.EmailService,
	userDataSvc thunderdome.UserDataSvc,
) *Service {
	// The library needs to be configured with your account's secret key.
	// Ensure the key is kept out of any version control system you might be using.
	stripe.Key = config.AccountSecret
	return &Service{
		logger:      logger,
		config:      config,
		dataSvc:     dataSvc,
		emailSvc:    emailSvc,
		userDataSvc: userDataSvc,
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
			var c stripe.CheckoutSession
			err = json.Unmarshal(event.Data.Raw, &c)
			if err != nil {
				logger.Error("Error getting checkout session from event", zap.String("eventId", event.ID))
				w.WriteHeader(http.StatusBadRequest) // Return a 400 error on a bad signature
				return
			}

			sessionParams := stripe.CheckoutSessionParams{}
			sessionParams.AddExpand("subscription")
			cs, err := session.Get(c.ID, &sessionParams)
			if err != nil {
				logger.Error("Error getting session from event", zap.String("eventId", event.ID),
					zap.String("sessionId", cs.ID))
				w.WriteHeader(http.StatusBadRequest) // Return a 400 error on a bad signature
				return
			}
			if cs.Subscription == nil {
				logger.Error("Error getting subscription from event", zap.String("eventId", event.ID),
					zap.String("sessionId", cs.ID))
				w.WriteHeader(http.StatusBadRequest) // Return a 400 error on a bad signature
				return
			}
			if cs.LineItems == nil || len(cs.LineItems.Data) < 1 {
				logger.Error("Error getting subscription product from event", zap.String("eventId", event.ID),
					zap.String("sessionId", cs.ID))
				w.WriteHeader(http.StatusBadRequest) // Return a 400 error on a bad signature
				return
			}
			productId := cs.LineItems.Data[0].ID
			p, err := product.Get(productId, nil)
			if err != nil {
				logger.Error("Error getting product from event", zap.String("eventId", event.ID),
					zap.String("productId", productId))
				w.WriteHeader(http.StatusBadRequest) // Return a 400 error on a bad signature
				return
			}
			subType, ok := p.Metadata["plan_type"]
			if !ok {
				logger.Error("Error getting product type from event", zap.String("eventId", event.ID),
					zap.String("productId", productId))
				w.WriteHeader(http.StatusBadRequest) // Return a 400 error on a bad signature
				return
			}
			expires := time.Unix(cs.Subscription.CurrentPeriodEnd, 0)

			sub := thunderdome.Subscription{
				UserID:         cs.ClientReferenceID,
				CustomerID:     cs.Customer.ID,
				SubscriptionID: cs.Subscription.ID,
				Type:           subType,
				Expires:        expires,
			}

			_, err = s.dataSvc.CreateSubscription(ctx, sub)
			if err != nil {
				logger.Error(fmt.Sprintf("Error creating subscription: %v", err), zap.String("eventId", event.ID))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			go func(ctx context.Context) {
				user, userErr := s.userDataSvc.GetUser(ctx, sub.UserID)
				if userErr != nil {
					logger.Error(fmt.Sprintf("error getting user to send subscription active email: %v", userErr),
						zap.String("eventId", event.ID))
					return
				}
				emailErr := s.emailSvc.SendUserSubscriptionActive(user.Id, user.Email, sub.Type)
				if emailErr != nil {
					logger.Error(fmt.Sprintf("error sending subscription active email: %v", emailErr),
						zap.String("eventId", event.ID))
				}
			}(context.WithoutCancel(ctx))
		case "customer.subscription.updated":
			var sub stripe.Subscription
			err := json.Unmarshal(event.Data.Raw, &sub)
			if err != nil {
				logger.Error("Error getting subscription from event", zap.String("eventId", event.ID))
				w.WriteHeader(http.StatusBadRequest) // Return a 400 error on a bad signature
				return
			}

			subscription, err := s.dataSvc.GetSubscriptionBySubscriptionID(ctx, sub.ID)
			if err != nil {
				logger.Error(fmt.Sprintf("Error getting subscription id %s subscription: %v", sub.ID, err), zap.String("eventId", event.ID))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			subscription.Expires = time.Unix(sub.CurrentPeriodEnd, 0)
			subscriptionStatusChanged := false
			stripeStatusIsActive := sub.Status == "active"
			if subscription.Active != stripeStatusIsActive {
				subscriptionStatusChanged = true
				subscription.Active = stripeStatusIsActive
			}
			//subscription.Type = "user" // @TODO - get subtype from update metadata and update if different

			_, err = s.dataSvc.UpdateSubscription(ctx, subscription.ID, subscription)
			if err != nil {
				logger.Error(fmt.Sprintf("Error creating subscription: %v", err), zap.String("eventId", event.ID))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			if subscriptionStatusChanged {
				if !subscription.Active {
					go func(ctx context.Context) {
						user, userErr := s.userDataSvc.GetUser(ctx, subscription.UserID)
						if userErr != nil {
							logger.Error(fmt.Sprintf("error getting user to send subscription deactivated email: %v", userErr),
								zap.String("eventId", event.ID))
							return
						}
						emailErr := s.emailSvc.SendUserSubscriptionDeactivated(user.Id, user.Email, subscription.Type)
						if emailErr != nil {
							logger.Error(fmt.Sprintf("error sending subscription deactivated email: %v", emailErr),
								zap.String("eventId", event.ID))
						}
					}(context.WithoutCancel(ctx))
				} else {
					go func(ctx context.Context) {
						user, userErr := s.userDataSvc.GetUser(ctx, subscription.UserID)
						if userErr != nil {
							logger.Error(fmt.Sprintf("error getting user to send subscription active email: %v", userErr),
								zap.String("eventId", event.ID))
							return
						}
						emailErr := s.emailSvc.SendUserSubscriptionActive(user.Id, user.Email, subscription.Type)
						if emailErr != nil {
							logger.Error(fmt.Sprintf("error sending subscription activate email: %v", emailErr),
								zap.String("eventId", event.ID))
						}
					}(context.WithoutCancel(ctx))
				}
			}
		default:
			logger.Error(fmt.Sprintf("Unhandled Stripe webhook event type: %s", event.Type), zap.String("eventId", event.ID))
		}

		logger.Info(fmt.Sprintf("Successfully processed Stripe webhook event type: %s", event.Type), zap.String("eventId", event.ID))

		w.WriteHeader(http.StatusOK)
	}
}
