package config

import (
	"fmt"
	"strings"
)

const (
	defaultHTTPDomain        = "thunderdome.dev"
	defaultHTTPCookieHashkey = "strongest-avenger"
	defaultConfigAESHashkey  = "therevengers"
	defaultDBUser            = "thor"
	defaultDBPass            = "odinson"
	defaultAuthMethod        = "normal"
)

type ValidationIssue struct {
	Key     string
	Message string
}

func (i ValidationIssue) Error() string {
	return fmt.Sprintf("%s: %s", i.Key, i.Message)
}

// Validate returns config issues that should be resolved before running in a non-development environment.
func (c Config) Validate() []ValidationIssue {
	issues := make([]ValidationIssue, 0)

	issues = appendIfInvalid(issues, "http.domain", c.Http.Domain == "", "must be configured")
	issues = appendIfInvalid(issues, "http.domain", c.Http.Domain == defaultHTTPDomain, "must not use the default value")
	issues = appendIfInvalid(issues, "http.cookie_hashkey", c.Http.CookieHashkey == "", "must be configured")
	issues = appendIfInvalid(issues, "http.cookie_hashkey", c.Http.CookieHashkey == defaultHTTPCookieHashkey, "must not use the default value")
	issues = appendIfInvalid(issues, "config.aes_hashkey", c.Config.AesHashkey == "", "must be configured")
	issues = appendIfInvalid(issues, "config.aes_hashkey", c.Config.AesHashkey == defaultConfigAESHashkey, "must not use the default value")
	issues = appendIfInvalid(issues, "db.user", c.Db.User == "", "must be configured")
	issues = appendIfInvalid(issues, "db.user", c.Db.User == defaultDBUser, "should not use the default value")
	issues = appendIfInvalid(issues, "db.pass", c.Db.Pass == "", "must be configured")
	issues = appendIfInvalid(issues, "db.pass", c.Db.Pass == defaultDBPass, "must not use the default value")

	switch strings.ToLower(strings.TrimSpace(c.Auth.Method)) {
	case "", defaultAuthMethod:
	case "header":
		issues = appendIfInvalid(issues, "auth.header.usernameHeader", strings.TrimSpace(c.Auth.Header.UsernameHeader) == "", "must be configured when auth.method=header")
		issues = appendIfInvalid(issues, "auth.header.emailHeader", strings.TrimSpace(c.Auth.Header.EmailHeader) == "", "must be configured when auth.method=header")
	case "ldap":
		issues = appendIfInvalid(issues, "auth.ldap.url", strings.TrimSpace(c.Auth.Ldap.Url) == "", "must be configured when auth.method=ldap")
		issues = appendIfInvalid(issues, "auth.ldap.bindname", strings.TrimSpace(c.Auth.Ldap.Bindname) == "", "must be configured when auth.method=ldap")
		issues = appendIfInvalid(issues, "auth.ldap.bindpass", strings.TrimSpace(c.Auth.Ldap.Bindpass) == "", "must be configured when auth.method=ldap")
		issues = appendIfInvalid(issues, "auth.ldap.basedn", strings.TrimSpace(c.Auth.Ldap.Basedn) == "", "must be configured when auth.method=ldap")
		issues = appendIfInvalid(issues, "auth.ldap.filter", strings.TrimSpace(c.Auth.Ldap.Filter) == "", "must be configured when auth.method=ldap")
		issues = appendIfInvalid(issues, "auth.ldap.filter", !strings.Contains(c.Auth.Ldap.Filter, "%s"), "must include %s when auth.method=ldap")
		issues = appendIfInvalid(issues, "auth.ldap.mail_attr", strings.TrimSpace(c.Auth.Ldap.MailAttr) == "", "must be configured when auth.method=ldap")
		issues = appendIfInvalid(issues, "auth.ldap.cn_attr", strings.TrimSpace(c.Auth.Ldap.CnAttr) == "", "must be configured when auth.method=ldap")
	case "oidc":
		issues = appendIfInvalid(issues, "auth.oidc.provider_name", strings.TrimSpace(c.Auth.OIDC.ProviderName) == "", "must be configured when auth.method=oidc")
		issues = appendIfInvalid(issues, "auth.oidc.provider_url", strings.TrimSpace(c.Auth.OIDC.ProviderURL) == "", "must be configured when auth.method=oidc")
		issues = appendIfInvalid(issues, "auth.oidc.client_id", strings.TrimSpace(c.Auth.OIDC.ClientID) == "", "must be configured when auth.method=oidc")
		issues = appendIfInvalid(issues, "auth.oidc.client_secret", strings.TrimSpace(c.Auth.OIDC.ClientSecret) == "", "must be configured when auth.method=oidc")
	default:
		issues = append(issues, ValidationIssue{Key: "auth.method", Message: "must be one of normal, header, ldap, oidc"})
	}

	if c.Auth.Google.Enabled {
		issues = appendIfInvalid(issues, "auth.google.client_id", strings.TrimSpace(c.Auth.Google.ClientID) == "", "must be configured when auth.google.enabled=true")
		issues = appendIfInvalid(issues, "auth.google.client_secret", strings.TrimSpace(c.Auth.Google.ClientSecret) == "", "must be configured when auth.google.enabled=true")
	}

	if c.Config.SubscriptionsEnabled {
		issues = appendIfInvalid(issues, "subscription.account_secret", strings.TrimSpace(c.Subscription.AccountSecret) == "", "must be configured when config.subscriptions_enabled=true")
		issues = appendIfInvalid(issues, "subscription.webhook_secret", strings.TrimSpace(c.Subscription.WebhookSecret) == "", "must be configured when config.subscriptions_enabled=true")
	}

	return issues
}

func appendIfInvalid(issues []ValidationIssue, key string, invalid bool, message string) []ValidationIssue {
	if invalid {
		issues = append(issues, ValidationIssue{Key: key, Message: message})
	}

	return issues
}
