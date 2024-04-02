package oauth

import (
	"context"
	"net/http"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"

	"github.com/uptrace/opentelemetry-go-extra/otelzap"

	"github.com/StevenWeathers/thunderdome-planning-poker/internal/cookie"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

func New(
	config Config,
	cookie *cookie.Cookie,
	logger *otelzap.Logger,
	authDataSvc thunderdome.AuthDataSvc,
	subscriptionDataSvc thunderdome.SubscriptionDataSvc,
	ctx context.Context,
) (*Service, error) {
	s := Service{
		config:              config,
		cookie:              cookie,
		logger:              logger,
		authDataSvc:         authDataSvc,
		subscriptionDataSvc: subscriptionDataSvc,
	}
	provider, err := oidc.NewProvider(ctx, config.ProviderURL)
	if err != nil {
		return nil, err
	}

	s.oauth2Config = &oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		RedirectURL:  config.CallbackRedirectURL,

		// Discovery returns the OAuth2 endpoints.
		Endpoint: provider.Endpoint(),

		// "openid" is a required scope for OpenID Connect flows.
		Scopes: []string{oidc.ScopeOpenID, "profile", "email"},
	}

	s.verifier = provider.VerifierContext(ctx, &oidc.Config{ClientID: config.ClientID})

	return &s, nil
}

func (s *Service) HandleOAuth2Redirect() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		nonce, err := s.authDataSvc.OauthCreateNonce(ctx)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// create state cookie for callback state verification
		state, err := uuid.NewUUID()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		stateString := state.String()
		err = s.cookie.CreateAuthStateCookie(w, stateString)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, s.oauth2Config.AuthCodeURL(stateString, oidc.Nonce(nonce)), http.StatusSeeOther)
	}
}

func (s *Service) HandleOAuth2Callback() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := s.logger.Ctx(ctx)
		rq := r.URL.Query()
		state := rq.Get("state")
		code := rq.Get("code")

		// Verify state
		err := s.cookie.ValidateAuthStateCookie(w, r, state)
		if err != nil {
			logger.Error("invalid oauth state", zap.Error(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Exchange code for oauth token
		oauth2Token, err := s.oauth2Config.Exchange(ctx, code)
		if err != nil {
			logger.Error("error exchanging oidc code for token", zap.Error(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Extract the ID Token from OAuth2 token.
		rawIDToken, ok := oauth2Token.Extra("id_token").(string)
		if !ok {
			logger.Error("missing oauth2 id_token")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Parse and verify ID Token payload.
		idToken, err := s.verifier.Verify(ctx, rawIDToken)
		if err != nil {
			logger.Error("error parsing and verifying id_token", zap.Error(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Extract custom claims
		var claims struct {
			Name          string `json:"name"`
			Email         string `json:"email"`
			EmailVerified bool   `json:"email_verified"`
			Nonce         string `json:"nonce"`
			Picture       string `json:"picture"`
		}
		if err := idToken.Claims(&claims); err != nil {
			logger.Error("error extracting custom claims from id_token", zap.Error(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		nonceErr := s.authDataSvc.OauthValidateNonce(ctx, claims.Nonce)
		if nonceErr != nil {
			logger.Error("nonce validation failed", zap.Error(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		user, sessionId, userErr := s.authDataSvc.OauthAuthUser(ctx, s.config.ProviderName, claims.Email, claims.EmailVerified, claims.Name, claims.Picture)
		if userErr != nil {
			logger.Error("error authenticating oauth user", zap.Error(userErr))
			ue := err.Error()
			if ue == "USER_DISABLED" {
				w.WriteHeader(http.StatusUnauthorized)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}

		if scErr := s.cookie.CreateSessionCookie(w, sessionId); scErr != nil {
			logger.Error("error creating oauth user session cookie", zap.Error(scErr), zap.String("userId", user.Id))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		subscribedErr := s.subscriptionDataSvc.CheckActiveSubscriber(ctx, user.Id)

		if err := s.cookie.CreateUserUICookie(w, thunderdome.UserUICookie{
			Id:                   user.Id,
			Name:                 user.Name,
			Email:                user.Email,
			Rank:                 user.Type,
			Locale:               user.Locale,
			NotificationsEnabled: user.NotificationsEnabled,
			Subscribed:           subscribedErr == nil,
		}); err != nil {
			logger.Error("error creating oauth user ui cookie", zap.Error(err), zap.String("userId", user.Id))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, s.config.UIRedirectURL, http.StatusFound)
	}
}
