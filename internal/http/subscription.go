package http

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
	"github.com/gorilla/mux"

	"go.uber.org/zap"
)

type subscriptionRequestBody struct {
	UserID         string    `json:"user_id"`
	TeamID         string    `json:"team_id"`
	OrganizationID string    `json:"organization_id"`
	CustomerID     string    `json:"customer_id"`
	SubscriptionID string    `json:"subscription_id"`
	Type           string    `json:"type" enums:"user, team, organization" validate:"required,oneof=user organization team"`
	Active         bool      `json:"active"`
	Expires        time.Time `json:"expires"`
}

// handleSubscriptionGet gets a subscription
// @Summary      Get Subscription
// @Description  Get a subscription
// @Tags         subscription
// @Produce      json
// @Param        subscriptionId  path    string  true  "the subscription ID"
// @Success      200     object  standardJsonResponse{data=thunderdome.Subscription}
// @Success      500     object  standardJsonResponse{}
// @Security     ApiKeyAuth
// @Router       /subscriptions/{subscriptionId} [get]
func (s *Service) handleSubscriptionGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)
		vars := mux.Vars(r)
		id := vars["subscriptionId"]
		idErr := validate.Var(id, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		sub, err := s.SubscriptionDataSvc.GetSubscriptionByID(ctx, id)
		if err != nil {
			s.Logger.Ctx(ctx).Error(
				"handleGetTeamByUser error", zap.Error(err), zap.String("subscription_id", id),
				zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, sub, nil)
	}
}

// handleGetSubscriptions gets a list of subscriptions
// @Summary      Get Subscriptions
// @Description  get list of subscriptions
// @Tags         subscription
// @Produce      json
// @Param        limit   query   int  false  "Max number of results to return"
// @Param        offset  query   int  false  "Starting point to return rows from, should be multiplied by limit or 0"
// @Success      200     object  standardJsonResponse{data=[]thunderdome.Subscription}
// @Failure      500     object  standardJsonResponse{}
// @Security     ApiKeyAuth
// @Router       /subscriptions [get]
func (s *Service) handleGetSubscriptions() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID, _ := ctx.Value(contextKeyUserID).(*string)
		limit, offset := getLimitOffsetFromRequest(r)
		subscriptions, count, err := s.SubscriptionDataSvc.GetSubscriptions(ctx, limit, offset)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleGetSubscriptions error", zap.Error(err),
				zap.Int("limit", limit), zap.Int("offset", offset),
				zap.Stringp("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		meta := &pagination{
			Count:  count,
			Offset: offset,
			Limit:  limit,
		}

		s.Success(w, r, http.StatusOK, subscriptions, meta)
	}
}

// handleSubscriptionCreate creates a new subscription
// @Summary      Create Subscription
// @Description  Creates a subscription
// @Tags         subscription
// @Produce      json
// @Param        subscription  body    subscriptionRequestBody  true  "new subscription object"
// @Success      200    object  standardJsonResponse{data=thunderdome.Subscription}
// @Failure      500    object  standardJsonResponse{}
// @Security     ApiKeyAuth
// @Router       /subscriptions [post]
func (s *Service) handleSubscriptionCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)
		var sub = subscriptionRequestBody{}
		body, bodyErr := io.ReadAll(r.Body)
		if bodyErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
			return
		}

		jsonErr := json.Unmarshal(body, &sub)
		if jsonErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		inputErr := validate.Struct(sub)
		if inputErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, inputErr.Error()))
			return
		}

		subscription, err := s.SubscriptionDataSvc.CreateSubscription(ctx, thunderdome.Subscription{
			UserID:         sub.UserID,
			TeamID:         sub.TeamID,
			OrganizationID: sub.OrganizationID,
			CustomerID:     sub.CustomerID,
			SubscriptionID: sub.SubscriptionID,
			Type:           sub.Type,
			Expires:        sub.Expires,
		})
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleSubscriptionCreate error", zap.Error(err),
				zap.String("sub_user_id", sub.UserID),
				zap.String("sub_customer_id", sub.CustomerID),
				zap.String("sub_subscription_id", sub.SubscriptionID),
				zap.String("sub_type", sub.Type),
				zap.Time("sub_expires", sub.Expires),
				zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, subscription, nil)
	}
}

// handleSubscriptionUpdate updates a subscription
// @Summary      Update Subscription
// @Description  Updates a Subscription
// @Tags         subscription
// @Produce      json
// @Param        subscriptionId  path    string    true  "the subscription ID to update"
// @Param        subscription    body    subscriptionRequestBody true  "subscription object to update"
// @Success      200      object  standardJsonResponse{data=thunderdome.Subscription}
// @Failure      500      object  standardJsonResponse{}
// @Security     ApiKeyAuth
// @Router       /subscriptions/{subscriptionId} [put]
func (s *Service) handleSubscriptionUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)
		vars := mux.Vars(r)
		subscriptionID := vars["subscriptionId"]
		idErr := validate.Var(subscriptionID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		var sub = subscriptionRequestBody{}
		body, bodyErr := io.ReadAll(r.Body)
		if bodyErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
			return
		}

		jsonErr := json.Unmarshal(body, &sub)
		if jsonErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		inputErr := validate.Struct(sub)
		if inputErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, inputErr.Error()))
			return
		}

		subscription, err := s.SubscriptionDataSvc.UpdateSubscription(ctx, subscriptionID, thunderdome.Subscription{
			UserID:         sub.UserID,
			TeamID:         sub.TeamID,
			OrganizationID: sub.OrganizationID,
			CustomerID:     sub.CustomerID,
			SubscriptionID: sub.SubscriptionID,
			Active:         sub.Active,
			Type:           sub.Type,
			Expires:        sub.Expires,
		})
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleSubscriptionUpdate error",
				zap.Error(err), zap.String("subscription_id", subscriptionID),
				zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, subscription, nil)
	}
}

// handleSubscriptionDelete handles deleting a subscription
// @Summary      Delete Subscription
// @Description  Deletes a Subscription
// @Tags         subscription
// @Produce      json
// @Param        subscriptionId  path    string  true  "the subscription ID to delete"
// @Success      200      object  standardJsonResponse{}
// @Failure      500      object  standardJsonResponse{}
// @Security     ApiKeyAuth
// @Router       /subscriptions/{subscriptionId} [delete]
func (s *Service) handleSubscriptionDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)
		vars := mux.Vars(r)
		subscriptionID := vars["subscriptionId"]
		idErr := validate.Var(subscriptionID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		err := s.SubscriptionDataSvc.DeleteSubscription(ctx, subscriptionID)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleSubscriptionDelete error", zap.Error(err),
				zap.String("subscription_id", subscriptionID),
				zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, nil, nil)
	}
}

// handleGetEntityUserActiveSubs gets a list of active subscriptions for the entity user
// @Summary      Get Entity User Active Subscriptions
// @Description  get list of active entity user subscriptions
// @Tags         subscription
// @Param        userId  path    string  true  "the entity user ID"
// @Produce      json
// @Success      200     object  standardJsonResponse{data=[]thunderdome.Subscription}
// @Failure      500     object  standardJsonResponse{}
// @Security     ApiKeyAuth
// @Router       /users/{userId}/subscriptions [get]
func (s *Service) handleGetEntityUserActiveSubs() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		vars := mux.Vars(r)
		sessionUserID, _ := ctx.Value(contextKeyUserID).(*string)
		entityUserID := vars["userId"]
		idErr := validate.Var(entityUserID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		subscriptions, err := s.SubscriptionDataSvc.GetActiveSubscriptionsByUserID(ctx, entityUserID)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleGetSubscriptions error", zap.Error(err),
				zap.String("entity_user_id", entityUserID),
				zap.Stringp("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, subscriptions, nil)
	}
}

type subscriptionAssociateRequestBody struct {
	TeamID         string `json:"team_id"`
	OrganizationID string `json:"organization_id"`
}

// handleEntityUserUpdateSubscription updates an entity users subscription association
// @Summary      Update Entity User Subscriptions
// @Description  get list of active entity user subscriptions
// @Tags         subscription
// @Param        userId  path    string  true  "the entity user ID"
// @Param        subscriptionId  path   string  true  "the subscription ID to update"
// @Param        subscription  body    subscriptionAssociateRequestBody  true  "update subscription association object"
// @Produce      json
// @Success      200     object  standardJsonResponse{data=[]thunderdome.Subscription}
// @Failure      500     object  standardJsonResponse{}
// @Security     ApiKeyAuth
// @Router       /users/{userId}/subscriptions/{subscriptionId} [patch]
func (s *Service) handleEntityUserUpdateSubscription() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		vars := mux.Vars(r)
		sessionUserID, _ := ctx.Value(contextKeyUserID).(*string)
		entityUserID := vars["userId"]
		idErr := validate.Var(entityUserID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		subscriptionID := vars["subscriptionId"]
		sidErr := validate.Var(subscriptionID, "required,uuid")
		if sidErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, sidErr.Error()))
			return
		}

		var sar = subscriptionAssociateRequestBody{}
		body, bodyErr := io.ReadAll(r.Body)
		if bodyErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
			return
		}

		jsonErr := json.Unmarshal(body, &sar)
		if jsonErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		sub, err := s.SubscriptionDataSvc.GetSubscriptionByID(ctx, subscriptionID)
		if err != nil {
			s.Logger.Ctx(ctx).Error(
				"handleEntityUserUpdateSubscription GetSubscriptionByID error", zap.Error(err),
				zap.String("subscription_id", subscriptionID),
				zap.String("entity_user_id", entityUserID),
				zap.Stringp("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		if sub.Type == "team" {
			if sar.TeamID == "" {
				s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, "MISSING REQUIRED team_id FIELD"))
				return
			} else {
				sub.TeamID = sar.TeamID
			}
		}

		if sub.Type == "organization" {
			if sar.OrganizationID == "" {
				s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, "MISSING REQUIRED organization_id FIELD"))
				return
			} else {
				sub.OrganizationID = sar.OrganizationID
			}
		}

		subscription, err := s.SubscriptionDataSvc.UpdateSubscription(ctx, sub.ID, sub)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleEntityUserUpdateSubscription error", zap.Error(err),
				zap.String("subscription_id", subscriptionID),
				zap.String("entity_user_id", entityUserID),
				zap.Stringp("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, subscription, nil)
	}
}
