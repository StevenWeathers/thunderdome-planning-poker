package jira

import (
	jira "github.com/ctreminiom/go-atlassian/jira/v3"
	"github.com/ctreminiom/go-atlassian/pkg/infra/models"
)

type Config struct {
	InstanceHost string `json:"instance_host"`
	AccessToken  string `json:"access_token"`
	ClientMail   string `json:"client_mail"`
}

type Client struct {
	instance *jira.Client
}

type IssuesSearchResult struct {
	Total  int                   `json:"total"`
	Issues []*models.IssueScheme `json:"issues"`
}
