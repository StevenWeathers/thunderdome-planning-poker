// Package jira provides JIRA cloud integration
package jira

import (
	"log"
	"net/http"
	"time"

	jira "github.com/ctreminiom/go-atlassian/jira/v3"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
)

type Config struct {
	InstanceHost string `json:"instance_host"`
	AccessToken  string `json:"access_token"`
	ClientMail   string `json:"client_mail"`
}

type Client struct {
	instance *jira.Client
	logger   *otelzap.Logger
}

func New(config Config, logger *otelzap.Logger) *Client {
	httpClient := http.Client{
		Timeout: time.Second * 10,
	}
	instance, err := jira.New(&httpClient, config.InstanceHost)
	if err != nil {
		log.Fatal(err)
	}
	instance.Auth.SetBasicAuth(config.ClientMail, config.AccessToken)

	return &Client{
		instance: instance,
		logger:   logger,
	}
}
