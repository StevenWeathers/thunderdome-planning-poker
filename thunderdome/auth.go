package thunderdome

import (
	"context"
	"net/http"
)

type AuthProviderConfig struct {
	ProviderName string `mapstructure:"provider_name"`
	ProviderURL  string `mapstructure:"provider_url"`
	ClientID     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
}

type AuthProviderSvc interface {
	HandleOAuth2Redirect(w http.ResponseWriter, r *http.Request)
	HandleOAuth2Callback(w http.ResponseWriter, r *http.Request)
}

type AuthDataSvc interface {
	AuthUser(ctx context.Context, UserEmail string, UserPassword string) (*User, string, error)
	OauthAuthUser(ctx context.Context, provider string, email string, emailVerified bool, name string, pictureUrl string) (*User, string, error)
	UserResetRequest(ctx context.Context, UserEmail string) (resetID string, UserName string, resetErr error)
	UserResetPassword(ctx context.Context, ResetID string, UserPassword string) (UserName string, UserEmail string, resetErr error)
	UserUpdatePassword(ctx context.Context, UserID string, UserPassword string) (Name string, Email string, resetErr error)
	UserVerifyRequest(ctx context.Context, UserId string) (*User, string, error)
	VerifyUserAccount(ctx context.Context, VerifyID string) error
	MFASetupGenerate(email string) (string, string, error)
	MFASetupValidate(ctx context.Context, UserID string, secret string, passcode string) error
	MFARemove(ctx context.Context, UserID string) error
	MFATokenValidate(ctx context.Context, SessionId string, passcode string) error
	CreateSession(ctx context.Context, UserId string) (string, error)
	EnableSession(ctx context.Context, SessionId string) error
	GetSessionUser(ctx context.Context, SessionId string) (*User, error)
	DeleteSession(ctx context.Context, SessionId string) error
}
