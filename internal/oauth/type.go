package oauth

import (
	"context"
	"net/http"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"golang.org/x/oauth2"
)

// Config holds the configuration for the oauth service
type Config struct {
	thunderdome.AuthProviderConfig
	CallbackRedirectURL string
	UIRedirectURL       string
	InternalOnlyOidc    bool
}

// CookieManager is an interface for managing cookies
type CookieManager interface {
	CreateUserCookie(w http.ResponseWriter, userID string) error
	CreateSessionCookie(w http.ResponseWriter, sessionID string) error
	CreateUserUICookie(w http.ResponseWriter, userUiCookie thunderdome.UserUICookie) error
	ClearUserCookies(w http.ResponseWriter)
	ValidateUserCookie(w http.ResponseWriter, r *http.Request) (string, error)
	ValidateSessionCookie(w http.ResponseWriter, r *http.Request) (string, error)
	CreateCookie(w http.ResponseWriter, cookieName string, value string, maxAge int) error
	GetCookie(w http.ResponseWriter, r *http.Request, cookieName string) (string, error)
	DeleteCookie(w http.ResponseWriter, cookieName string)
	CreateAuthStateCookie(w http.ResponseWriter, state string) error
	ValidateAuthStateCookie(w http.ResponseWriter, r *http.Request, state string) error
	DeleteAuthStateCookie(w http.ResponseWriter) error
}

// AuthDataSvc is an interface for the auth data service
type AuthDataSvc interface {
	OauthCreateNonce(ctx context.Context) (string, error)
	OauthValidateNonce(ctx context.Context, nonceId string) error
	OauthAuthUser(ctx context.Context, provider string, sub string, email string, emailVerified bool, name string, pictureUrl string) (*thunderdome.User, string, error)
	OauthUpsertUser(ctx context.Context, provider string, sub string, email string, emailVerified bool, name string, pictureUrl string) (*thunderdome.User, string, error)
}

// SubscriptionDataSvc is an interface for the subscription data service
type SubscriptionDataSvc interface {
	CheckActiveSubscriber(ctx context.Context, userID string) error
}

// Service is the oauth service
type Service struct {
	config              Config
	cookie              CookieManager
	oauth2Config        *oauth2.Config
	logger              *otelzap.Logger
	verifier            *oidc.IDTokenVerifier
	authDataSvc         AuthDataSvc
	subscriptionDataSvc SubscriptionDataSvc
}
