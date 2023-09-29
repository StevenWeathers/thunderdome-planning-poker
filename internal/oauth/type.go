package oauth

import (
	"github.com/StevenWeathers/thunderdome-planning-poker/internal/cookie"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"golang.org/x/oauth2"
)

type Config struct {
	ProviderURL      string
	ClientID         string
	ClientSecret     string
	RedirectURL      string
	StateCookieName  string
	PathPrefix       string
	AppDomain        string
	SecureCookieFlag bool
}

type AuthProvider struct {
	config       Config
	cookie       *cookie.Cookie
	oauth2Config oauth2.Config
	logger       *otelzap.Logger
	verifier     *oidc.IDTokenVerifier
}
