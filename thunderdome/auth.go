package thunderdome

import (
	"net/http"
	"time"
)

type AuthProviderConfig struct {
	ProviderName string `mapstructure:"provider_name"`
	ProviderURL  string `mapstructure:"provider_url"`
	ClientID     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
}

type Credential struct {
	UserID      string    `json:"user_id"`
	Email       string    `json:"email"`
	Password    string    `json:"-"`
	Verified    bool      `json:"verified"`
	MFAEnabled  bool      `json:"mfa_enabled"`
	CreatedDate time.Time `json:"created_date"`
	UpdatedDate time.Time `json:"updated_date"`
}

type Identity struct {
	UserID      string    `json:"user_id"`
	Provider    string    `json:"provider"`
	Sub         string    `json:"-"`
	Email       string    `json:"email"`
	Verified    bool      `json:"verified"`
	CreatedDate time.Time `json:"created_date"`
	UpdatedDate time.Time `json:"updated_date"`
}

type AuthProviderSvc interface {
	HandleOAuth2Redirect(w http.ResponseWriter, r *http.Request)
	HandleOAuth2Callback(w http.ResponseWriter, r *http.Request)
}
