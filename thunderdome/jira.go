package thunderdome

import (
	"context"
	"time"
)

type JiraInstance struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Host        string    `json:"host"`
	ClientMail  string    `json:"client_mail"`
	AccessToken string    `json:"access_token"`
	CreatedDate time.Time `json:"created_date"`
	UpdatedDate time.Time `json:"updated_date"`
}

type JiraDataSvc interface {
	FindInstancesByUserId(ctx context.Context, userId string) ([]JiraInstance, error)
	CreateInstance(ctx context.Context, userId string, host string, clientMail string, accessToken string) (JiraInstance, error)
	UpdateInstance(ctx context.Context, instanceId string, host string, clientMail string, accessToken string) (JiraInstance, error)
	DeleteInstance(ctx context.Context, instanceId string) error
}
