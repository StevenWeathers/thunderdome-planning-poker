package main

import (
	"context"

	"github.com/spf13/viper"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
)

// InitConfig initializes the application configuration
func InitConfig(logger *otelzap.Logger) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AddConfigPath("/etc/thunderdome/")
	viper.AddConfigPath("$HOME/.config/thunderdome/")
	viper.AddConfigPath(".")

	viper.SetDefault("http.cookie_hashkey", "strongest-avenger")
	viper.SetDefault("http.port", "8080")
	viper.SetDefault("http.secure_cookie", true)
	viper.SetDefault("http.backend_cookie_name", "warriorId")
	viper.SetDefault("http.session_cookie_name", "sessionId")
	viper.SetDefault("http.frontend_cookie_name", "warrior")
	viper.SetDefault("http.domain", "thunderdome.dev")
	viper.SetDefault("http.path_prefix", "")
	viper.SetDefault("http.write_timeout", 5)
	viper.SetDefault("http.read_timeout", 5)
	viper.SetDefault("http.idle_timeout", 30)
	viper.SetDefault("http.read_header_timeout", 2)

	viper.SetDefault("analytics.enabled", true)
	viper.SetDefault("analytics.id", "UA-140245309-1")

	viper.SetDefault("otel.enabled", false)
	viper.SetDefault("otel.service_name", "thunderdome")
	viper.SetDefault("otel.collector_url", "localhost:4317")
	viper.SetDefault("otel.insecure_mode", false)

	viper.SetDefault("db.host", "db")
	viper.SetDefault("db.port", 5432)
	viper.SetDefault("db.user", "thor")
	viper.SetDefault("db.pass", "odinson")
	viper.SetDefault("db.name", "thunderdome")
	viper.SetDefault("db.sslmode", "disable")
	viper.SetDefault("db.max_open_conns", 25)
	viper.SetDefault("db.max_idle_conns", 25)
	viper.SetDefault("db.conn_max_lifetime", 5)

	viper.SetDefault("smtp.enabled", true)
	viper.SetDefault("smtp.host", "localhost")
	viper.SetDefault("smtp.port", "25")
	viper.SetDefault("smtp.secure", true)
	viper.SetDefault("smtp.sender", "no-reply@thunderdome.dev")

	viper.SetDefault("config.aes_hashkey", "therevengers")
	viper.SetDefault("config.allowedPointValues",
		[]string{"0", "1/2", "1", "2", "3", "5", "8", "13", "20", "21", "34", "40", "55", "100", "?", "☕️"})
	viper.SetDefault("config.defaultPointValues",
		[]string{"1", "2", "3", "5", "8", "13", "?"})
	viper.SetDefault("config.show_warrior_rank", false)
	viper.SetDefault("config.avatar_service", "gravatar")
	viper.SetDefault("config.toast_timeout", 1000)
	viper.SetDefault("config.allow_guests", true)
	viper.SetDefault("config.allow_registration", true)
	viper.SetDefault("config.allow_jira_import", true)
	viper.SetDefault("config.allow_csv_import", true)
	viper.SetDefault("config.default_locale", "en")
	viper.SetDefault("config.friendly_ui_verbs", false)
	viper.SetDefault("config.allow_external_api", true)
	viper.SetDefault("config.external_api_verify_required", true)
	viper.SetDefault("config.user_apikey_limit", 5)
	viper.SetDefault("config.show_active_countries", false)
	viper.SetDefault("config.cleanup_battles_days_old", 180)
	viper.SetDefault("config.cleanup_guests_days_old", 180)
	viper.SetDefault("config.cleanup_retros_days_old", 180)
	viper.SetDefault("config.cleanup_storyboards_days_old", 180)
	viper.SetDefault("config.organizations_enabled", true)
	viper.SetDefault("config.require_teams", false)

	// feature flags
	viper.SetDefault("feature.poker", true)
	viper.SetDefault("feature.retro", true)
	viper.SetDefault("feature.storyboard", true)

	viper.SetDefault("auth.method", "normal")
	viper.SetDefault("auth.ldap.url", "")
	viper.SetDefault("auth.ldap.use_tls", true)
	viper.SetDefault("auth.ldap.bindname", "")
	viper.SetDefault("auth.ldap.bindpass", "")
	viper.SetDefault("auth.ldap.basedn", "")
	viper.SetDefault("auth.ldap.filter", "(&(objectClass=posixAccount)(mail=%s))")
	viper.SetDefault("auth.ldap.mail_attr", "mail")
	viper.SetDefault("auth.ldap.cn_attr", "cn")
	viper.SetDefault("auth.header.usernameHeader", "Remote-User")
	viper.SetDefault("auth.header.emailHeader", "Remote-Email")

	_ = viper.BindEnv("http.cookie_hashkey", "COOKIE_HASHKEY")
	_ = viper.BindEnv("http.port", "PORT")
	_ = viper.BindEnv("http.secure_cookie", "COOKIE_SECURE")
	_ = viper.BindEnv("http.backend_cookie_name", "SECURE_COOKIE_NAME")
	_ = viper.BindEnv("http.session_cookie_name", "SESSION_COOKIE_NAME")
	_ = viper.BindEnv("http.frontend_cookie_name", "FRONTEND_COOKIE_NAME")
	_ = viper.BindEnv("http.domain", "APP_DOMAIN")
	_ = viper.BindEnv("http.path_prefix", "PATH_PREFIX")
	_ = viper.BindEnv("http.write_timeout", "HTTP_WRITE_TIMEOUT")
	_ = viper.BindEnv("http.read_timeout", "HTTP_READ_TIMEOUT")
	_ = viper.BindEnv("http.idle_timeout", "HTTP_IDLE_TIMEOUT")
	_ = viper.BindEnv("http.read_header_timeout", "HTTP_READ_HEADER_TIMEOUT")

	_ = viper.BindEnv("analytics.enabled", "ANALYTICS_ENABLED")
	_ = viper.BindEnv("analytics.id", "ANALYTICS_ID")
	_ = viper.BindEnv("admin.email", "ADMIN_EMAIL")

	_ = viper.BindEnv("otel.enabled", "OTEL_ENABLED")
	_ = viper.BindEnv("otel.service_name", "OTEL_SERVICE_NAME")
	_ = viper.BindEnv("otel.collector_url", "OTEL_COLLECTOR_URL")
	_ = viper.BindEnv("otel.insecure_mode", "OTEL_INSECURE_MODE")

	_ = viper.BindEnv("db.host", "DB_HOST")
	_ = viper.BindEnv("db.port", "DB_PORT")
	_ = viper.BindEnv("db.user", "DB_USER")
	_ = viper.BindEnv("db.pass", "DB_PASS")
	_ = viper.BindEnv("db.name", "DB_NAME")
	_ = viper.BindEnv("db.sslmode", "DB_SSLMODE")
	_ = viper.BindEnv("db.max_open_conns", "DB_MAX_OPEN_CONNS")
	_ = viper.BindEnv("db.max_idle_conns", "DB_MAX_IDLE_CONNS")
	_ = viper.BindEnv("db.conn_max_lifetime", "DB_CONN_MAX_LIFETIME")

	_ = viper.BindEnv("smtp.enabled", "SMTP_ENABLED")
	_ = viper.BindEnv("smtp.host", "SMTP_HOST")
	_ = viper.BindEnv("smtp.port", "SMTP_PORT")
	_ = viper.BindEnv("smtp.secure", "SMTP_SECURE")
	_ = viper.BindEnv("smtp.identity", "SMTP_IDENTITY")
	_ = viper.BindEnv("smtp.user", "SMTP_USER")
	_ = viper.BindEnv("smtp.pass", "SMTP_PASS")
	_ = viper.BindEnv("smtp.sender", "SMTP_SENDER")

	_ = viper.BindEnv("config.aes_hashkey", "CONFIG_AES_HASHKEY")
	_ = viper.BindEnv("config.allowedPointValues", "CONFIG_POINTS_ALLOWED")
	_ = viper.BindEnv("config.defaultPointValues", "CONFIG_POINTS_DEFAULT")
	_ = viper.BindEnv("config.show_warrior_rank", "CONFIG_SHOW_RANK")
	_ = viper.BindEnv("config.avatar_service", "CONFIG_AVATAR_SERVICE")
	_ = viper.BindEnv("config.toast_timeout", "CONFIG_TOAST_TIMEOUT")
	_ = viper.BindEnv("config.allow_guests", "CONFIG_ALLOW_GUESTS")
	_ = viper.BindEnv("config.allow_registration", "CONFIG_ALLOW_REGISTRATION")
	_ = viper.BindEnv("config.allow_jira_import", "CONFIG_ALLOW_JIRA_IMPORT")
	_ = viper.BindEnv("config.allow_CSV_import", "CONFIG_ALLOW_CSV_IMPORT")
	_ = viper.BindEnv("config.default_locale", "CONFIG_DEFAULT_LOCALE")
	_ = viper.BindEnv("config.friendly_ui_verbs", "CONFIG_FRIENDLY_UI_VERBS")
	_ = viper.BindEnv("config.allow_external_api", "CONFIG_ALLOW_EXTERNAL_API")
	_ = viper.BindEnv("config.external_api_verify_required", "CONFIG_EXTERNAL_API_VERIFY_REQUIRED")
	_ = viper.BindEnv("config.user_apikey_limit", "CONFIG_USER_APIKEY_LIMIT")
	_ = viper.BindEnv("config.show_active_countries", "CONFIG_SHOW_ACTIVE_COUNTRIES")
	_ = viper.BindEnv("config.cleanup_battles_days_old", "CONFIG_CLEANUP_BATTLES_DAYS_OLD")
	_ = viper.BindEnv("config.cleanup_guests_days_old", "CONFIG_CLEANUP_GUESTS_DAYS_OLD")
	_ = viper.BindEnv("config.cleanup_retros_days_old", "CONFIG_CLEANUP_RETROS_DAYS_OLD")
	_ = viper.BindEnv("config.cleanup_storyboards_days_old", "CONFIG_CLEANUP_STORYBOARDS_DAYS_OLD")
	_ = viper.BindEnv("config.organizations_enabled", "CONFIG_ORGANIZATIONS_ENABLED")
	_ = viper.BindEnv("config.require_teams", "CONFIG_REQUIRE_TEAMS")

	_ = viper.BindEnv("feature.poker", "FEATURE_POKER")
	_ = viper.BindEnv("feature.retro", "FEATURE_RETRO")
	_ = viper.BindEnv("feature.storyboard", "FEATURE_STORYBOARD")

	_ = viper.BindEnv("auth.method", "AUTH_METHOD")
	_ = viper.BindEnv("auth.ldap.url", "AUTH_LDAP_URL")
	_ = viper.BindEnv("auth.ldap.use_tls", "AUTH_LDAP_USE_TLS")
	_ = viper.BindEnv("auth.ldap.bindname", "AUTH_LDAP_BINDNAME")
	_ = viper.BindEnv("auth.ldap.bindpass", "AUTH_LDAP_BINDPASS")
	_ = viper.BindEnv("auth.ldap.basedn", "AUTH_LDAP_BASEDN")
	_ = viper.BindEnv("auth.ldap.filter", "AUTH_LDAP_FILTER")
	_ = viper.BindEnv("auth.ldap.mail_attr", "AUTH_LDAP_MAIL_ATTR")
	_ = viper.BindEnv("auth.ldap.cn_attr", "AUTH_LDAP_CN_ATTR")
	_ = viper.BindEnv("auth.header.usernameHeader", "AUTH_HEADER_USERNAME_HEADER")
	_ = viper.BindEnv("auth.header.emailHeader", "AUTH_HEADER_EMAIL_HEADER")

	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			logger.Ctx(context.Background()).Fatal(err.Error())
		}
	}
}
