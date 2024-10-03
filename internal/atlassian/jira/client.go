// Package jira provides JIRA cloud integration
package jira

import (
	"net/http"
	"time"

	jira "github.com/ctreminiom/go-atlassian/jira/v3"
)

// New creates a new JIRA client
func New(config Config) (*Client, error) {
	httpClient := http.Client{
		Timeout: time.Second * 10,
	}
	instance, err := jira.New(&httpClient, config.InstanceHost)
	if err != nil {
		return nil, err
	}
	instance.Auth.SetBasicAuth(config.ClientMail, config.AccessToken)

	return &Client{
		instance: instance,
	}, nil
}
