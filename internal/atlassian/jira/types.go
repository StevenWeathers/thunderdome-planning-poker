package jira

import (
	jira "github.com/ctreminiom/go-atlassian/v2/jira/v3"
	"github.com/ctreminiom/go-atlassian/v2/pkg/infra/models"
)

// Config is the configuration for the Jira client
type Config struct {
	InstanceHost   string `json:"instance_host"`
	JiraDataCenter bool   `json:"jira_data_center"`
	AccessToken    string `json:"access_token"`
	ClientMail     string `json:"client_mail"`
}

// Client is the Jira client
type Client struct {
	instance *jira.Client
}

// IssuesSearchResult is the result of a search for issues
type IssuesSearchResult struct {
	Total  int                   `json:"total"`
	Issues []*models.IssueScheme `json:"issues"`
}
