package oauth

import (
	"context"
	"net/http"
	"time"

	"github.com/StevenWeathers/thunderdome-planning-poker/internal/cookie"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

func New(config Config, cookie *cookie.Cookie, ctx context.Context) (*AuthProvider, error) {
	ap := AuthProvider{
		config: config,
		cookie: cookie,
	}
	provider, err := oidc.NewProvider(ctx, config.ProviderURL)
	if err != nil {
		return nil, err
	}

	ap.oauth2Config = oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		RedirectURL:  config.RedirectURL,

		// Discovery returns the OAuth2 endpoints.
		Endpoint: provider.Endpoint(),

		// "openid" is a required scope for OpenID Connect flows.
		Scopes: []string{oidc.ScopeOpenID, "profile", "email"},
	}

	ap.verifier = provider.Verifier(&oidc.Config{ClientID: config.ClientID})

	return &ap, nil
}

func (ap *AuthProvider) handleRedirect(w http.ResponseWriter, r *http.Request) {
	// @TODO - create a nonce in DB to send to oauth provider

	state, err := uuid.NewUUID()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = ap.cookie.CreateCookie(w, ap.config.StateCookieName, state.String(), int(time.Minute.Seconds()*10))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, ap.oauth2Config.AuthCodeURL(state.String()), http.StatusFound)
}

func (ap *AuthProvider) handleOAuth2Callback(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := ap.logger.Ctx(ctx)
	rq := r.URL.Query()
	state := rq.Get("state")

	// Verify state
	stateCookie, err := ap.cookie.GetCookie(w, r, ap.config.StateCookieName)
	if err != nil || state != stateCookie {
		logger.Error("invalid oauth state", zap.String("stateParam", state),
			zap.String("stateCookie", stateCookie), zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Exchange code for oauth token
	oauth2Token, err := ap.oauth2Config.Exchange(ctx, rq.Get("code"))
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
	idToken, err := ap.verifier.Verify(ctx, rawIDToken)
	if err != nil {
		logger.Error("error parsing and verifying id_token", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Extract custom claims
	var claims struct {
		Email    string `json:"email"`
		Verified bool   `json:"email_verified"`
		Nonce    string `json:"nonce"`
	}
	if err := idToken.Claims(&claims); err != nil {
		logger.Error("error extracting custom claims from id_token", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// @TODO - verify nonce

	// @TODO - find user and set session cookies then redirect
}
