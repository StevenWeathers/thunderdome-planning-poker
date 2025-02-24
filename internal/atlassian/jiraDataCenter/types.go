package jiradatacenter

import (
	jira_data_center "github.com/andygrunwald/go-jira/v2/onpremise"
)

// Config is the configuration for the Jira client
type Config struct {
	InstanceHost   string `json:"instance_host"`
	AccessToken    string `json:"access_token"`
	ClientMail     string `json:"client_mail"`
	JiraDataCenter bool   `json:"jira_data_center"`
}

// Client is the Jira client
type Client struct {
	instance *jira_data_center.Client
}

// IssuesSearchResult is the result of a search for issues
type IssuesSearchResult struct {
	Total  int                      `json:"total"`
	Issues []jira_data_center.Issue `json:"issues"`
}
