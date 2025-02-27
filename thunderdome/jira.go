package thunderdome

import (
	"time"
)

type JiraInstance struct {
	ID             string    `json:"id"`
	UserID         string    `json:"user_id"`
	Host           string    `json:"host"`
	ClientMail     string    `json:"client_mail"`
	AccessToken    string    `json:"access_token"`
	JiraDataCenter bool      `json:"jira_data_center"` // Checkbox for enabling Jira Data Center
	CreatedDate    time.Time `json:"created_date"`
	UpdatedDate    time.Time `json:"updated_date"`
}
