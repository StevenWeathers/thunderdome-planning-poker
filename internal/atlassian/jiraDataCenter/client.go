package jiradatacenter

import (
	jira "github.com/andygrunwald/go-jira/v2/onpremise"
)

// New creates a new JIRA client
func New(config Config) (*Client, error) {
	tp := jira.BearerAuthTransport{
		Token: config.AccessToken,
	}
	instance, err := jira.NewClient(config.InstanceHost, tp.Client())
	if err != nil {
		return nil, err
	}
	return &Client{
		instance: instance,
	}, nil
}
