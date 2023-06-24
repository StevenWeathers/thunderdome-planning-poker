package thunderdome

import (
	"context"
	"time"
)

// APIKey structure
type APIKey struct {
	Id          string    `json:"id"`
	Prefix      string    `json:"prefix"`
	UserId      string    `json:"userId"`
	Name        string    `json:"name"`
	Key         string    `json:"apiKey"`
	Active      bool      `json:"active"`
	CreatedDate time.Time `json:"createdDate"`
	UpdatedDate time.Time `json:"updatedDate"`
}

// UserAPIKey structure
type UserAPIKey struct {
	Id          string    `json:"id"`
	Prefix      string    `json:"prefix"`
	UserId      string    `json:"userId"`
	UserEmail   string    `json:"userEmail"`
	UserName    string    `json:"userName"`
	Name        string    `json:"name"`
	Key         string    `json:"apiKey"`
	Active      bool      `json:"active"`
	CreatedDate time.Time `json:"createdDate"`
	UpdatedDate time.Time `json:"updatedDate"`
}

type APIKeyService interface {
	GenerateApiKey(ctx context.Context, UserID string, KeyName string) (*APIKey, error)
	GetUserApiKeys(ctx context.Context, UserID string) ([]*APIKey, error)
	GetApiKeyUser(ctx context.Context, APK string) (*User, error)
	UpdateUserApiKey(ctx context.Context, UserID string, KeyID string, Active bool) ([]*APIKey, error)
	DeleteUserApiKey(ctx context.Context, UserID string, KeyID string) ([]*APIKey, error)
}
