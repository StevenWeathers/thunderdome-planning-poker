package config

import (
	"strings"
	"testing"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
)

func TestConfigValidateAcceptsConfiguredNormalAuth(t *testing.T) {
	c := Config{
		Http: Http{
			Domain:        "planning.example.com",
			CookieHashkey: "cookie-secret",
		},
		Db: Db{
			User: "planner",
			Pass: "db-secret",
		},
		Config: AppConfig{
			AesHashkey: "aes-secret",
		},
		Auth: Auth{
			Method: "normal",
		},
	}

	issues := c.Validate()
	if len(issues) != 0 {
		t.Fatalf("expected no validation issues, got %v", issues)
	}
}

func TestConfigValidateFlagsDefaultSensitiveValues(t *testing.T) {
	c := Config{
		Http: Http{
			Domain:        defaultHTTPDomain,
			CookieHashkey: defaultHTTPCookieHashkey,
		},
		Db: Db{
			User: defaultDBUser,
			Pass: defaultDBPass,
		},
		Config: AppConfig{
			AesHashkey: defaultConfigAESHashkey,
		},
		Auth: Auth{
			Method: "normal",
		},
	}

	issues := c.Validate()
	assertHasIssue(t, issues, "http.domain", "default value")
	assertHasIssue(t, issues, "http.cookie_hashkey", "default value")
	assertHasIssue(t, issues, "config.aes_hashkey", "default value")
	assertHasIssue(t, issues, "db.user", "default value")
	assertHasIssue(t, issues, "db.pass", "default value")
}

func TestConfigValidateFlagsMethodSpecificRequirements(t *testing.T) {
	t.Run("header auth", func(t *testing.T) {
		c := Config{
			Http:   Http{Domain: "planning.example.com", CookieHashkey: "cookie-secret"},
			Db:     Db{User: "planner", Pass: "db-secret"},
			Config: AppConfig{AesHashkey: "aes-secret"},
			Auth:   Auth{Method: "header"},
		}

		issues := c.Validate()
		assertHasIssue(t, issues, "auth.header.usernameHeader", "auth.method=header")
		assertHasIssue(t, issues, "auth.header.emailHeader", "auth.method=header")
	})

	t.Run("ldap auth", func(t *testing.T) {
		c := Config{
			Http:   Http{Domain: "planning.example.com", CookieHashkey: "cookie-secret"},
			Db:     Db{User: "planner", Pass: "db-secret"},
			Config: AppConfig{AesHashkey: "aes-secret"},
			Auth: Auth{
				Method: "ldap",
				Ldap:   AuthLdap{Filter: "(uid=person)"},
			},
		}

		issues := c.Validate()
		assertHasIssue(t, issues, "auth.ldap.url", "auth.method=ldap")
		assertHasIssue(t, issues, "auth.ldap.filter", "must include %s")
	})

	t.Run("oidc auth", func(t *testing.T) {
		c := Config{
			Http:   Http{Domain: "planning.example.com", CookieHashkey: "cookie-secret"},
			Db:     Db{User: "planner", Pass: "db-secret"},
			Config: AppConfig{AesHashkey: "aes-secret"},
			Auth:   Auth{Method: "oidc"},
		}

		issues := c.Validate()
		assertHasIssue(t, issues, "auth.oidc.provider_name", "auth.method=oidc")
		assertHasIssue(t, issues, "auth.oidc.client_secret", "auth.method=oidc")
	})

	t.Run("google auth", func(t *testing.T) {
		c := Config{
			Http:   Http{Domain: "planning.example.com", CookieHashkey: "cookie-secret"},
			Db:     Db{User: "planner", Pass: "db-secret"},
			Config: AppConfig{AesHashkey: "aes-secret"},
			Auth: Auth{
				Method: "normal",
				Google: Google{Enabled: true},
			},
		}

		issues := c.Validate()
		assertHasIssue(t, issues, "auth.google.client_id", "auth.google.enabled=true")
		assertHasIssue(t, issues, "auth.google.client_secret", "auth.google.enabled=true")
	})

	t.Run("subscriptions", func(t *testing.T) {
		c := Config{
			Http: Http{Domain: "planning.example.com", CookieHashkey: "cookie-secret"},
			Db:   Db{User: "planner", Pass: "db-secret"},
			Config: AppConfig{
				AesHashkey:           "aes-secret",
				SubscriptionsEnabled: true,
			},
			Auth:         Auth{Method: "normal"},
			Subscription: thunderdome.SubscriptionConfig{},
		}

		issues := c.Validate()
		assertHasIssue(t, issues, "subscription.account_secret", "subscriptions_enabled=true")
		assertHasIssue(t, issues, "subscription.webhook_secret", "subscriptions_enabled=true")
	})
}

func TestConfigValidateRejectsUnknownAuthMethod(t *testing.T) {
	c := Config{
		Http:   Http{Domain: "planning.example.com", CookieHashkey: "cookie-secret"},
		Db:     Db{User: "planner", Pass: "db-secret"},
		Config: AppConfig{AesHashkey: "aes-secret"},
		Auth:   Auth{Method: "sso"},
	}

	issues := c.Validate()
	assertHasIssue(t, issues, "auth.method", "must be one of")
}

func assertHasIssue(t *testing.T, issues []ValidationIssue, key string, messagePart string) {
	t.Helper()

	for _, issue := range issues {
		if issue.Key == key && strings.Contains(issue.Message, messagePart) {
			return
		}
	}

	t.Fatalf("expected issue for %s containing %q, got %v", key, messagePart, issues)
}
