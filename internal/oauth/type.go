package oauth

import (
	"github.com/StevenWeathers/thunderdome-planning-poker/internal/cookie"
	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"golang.org/x/oauth2"
)

type Config struct {
	thunderdome.AuthProviderConfig
	CallbackRedirectURL string
	UIRedirectURL       string
}

type Service struct {
	config              Config
	cookie              *cookie.Cookie
	oauth2Config        *oauth2.Config
	logger              *otelzap.Logger
	verifier            *oidc.IDTokenVerifier
	authDataSvc         thunderdome.AuthDataSvc
	subscriptionDataSvc thunderdome.SubscriptionDataSvc
}
