package config

import (
	"context"
	"strings"

	"github.com/spf13/viper"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
)

// InitConfig initializes the application configuration
func InitConfig(logger *otelzap.Logger) Config {
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
	viper.SetDefault("http.websocket_write_wait_sec", 10)
	viper.SetDefault("http.websocket_pong_wait_sec", 60)
	viper.SetDefault("http.websocket_ping_period_sec", 54)

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
	viper.SetDefault("smtp.skip_tls_verify", false)
	viper.SetDefault("smtp.sender", "no-reply@thunderdome.dev")
	viper.SetDefault("smtp.user", "")
	viper.SetDefault("smtp.pass", "")
	viper.SetDefault("smtp.auth", "PLAIN")

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
	viper.SetDefault("config.subscriptions_enabled", false)

	viper.SetDefault("subscription.account_secret", "")
	viper.SetDefault("subscription.webhook_secret", "")
	viper.SetDefault("subscription.manage_link", "https://billing.stripe.com/p/login/5kA5lKeb7eU9bp6cMM")
	viper.SetDefault("subscription.individual.enabled", true)
	viper.SetDefault("subscription.individual.month_price", "5")
	viper.SetDefault("subscription.individual.year_price", "50")
	viper.SetDefault("subscription.individual.month_checkout_link", "https://buy.stripe.com/7sIcP8gdhc3nc6YeUU")
	viper.SetDefault("subscription.individual.year_checkout_link", "https://buy.stripe.com/14kcP8e590kFb2UdQR")
	viper.SetDefault("subscription.team.enabled", false)
	viper.SetDefault("subscription.team.month_price", "20")
	viper.SetDefault("subscription.team.year_price", "200")
	viper.SetDefault("subscription.team.month_checkout_link", "https://buy.stripe.com/28o6qK5yD4AV3As5ks")
	viper.SetDefault("subscription.team.year_checkout_link", "https://buy.stripe.com/aEUg1kaSX4AV7QI14d")
	viper.SetDefault("subscription.organization.enabled", false)
	viper.SetDefault("subscription.organization.month_price", "50")
	viper.SetDefault("subscription.organization.year_price", "500")
	viper.SetDefault("subscription.organization.month_checkout_link", "https://buy.stripe.com/8wM6qK2mr0kF5IA8wC")
	viper.SetDefault("subscription.organization.year_checkout_link", "https://buy.stripe.com/eVa02m2mr7N74EwcMT")

	viper.SetDefault("admin.email", "")

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

	// automatically load matching envs
	viper.SetEnvKeyReplacer(strings.NewReplacer(`.`, `_`))
	viper.AutomaticEnv()

	// the following envs are not automatic because they didn't match the key structure
	_ = viper.BindEnv("http.cookie_hashkey", "COOKIE_HASHKEY")
	_ = viper.BindEnv("http.port", "PORT")
	_ = viper.BindEnv("http.secure_cookie", "COOKIE_SECURE")
	_ = viper.BindEnv("http.backend_cookie_name", "SECURE_COOKIE_NAME")
	_ = viper.BindEnv("http.session_cookie_name", "SESSION_COOKIE_NAME")
	_ = viper.BindEnv("http.frontend_cookie_name", "FRONTEND_COOKIE_NAME")
	_ = viper.BindEnv("http.domain", "APP_DOMAIN")
	_ = viper.BindEnv("http.path_prefix", "PATH_PREFIX")
	_ = viper.BindEnv("config.allowedPointValues", "CONFIG_POINTS_ALLOWED")
	_ = viper.BindEnv("config.defaultPointValues", "CONFIG_POINTS_DEFAULT")
	_ = viper.BindEnv("config.show_warrior_rank", "CONFIG_SHOW_RANK")
	_ = viper.BindEnv("auth.header.usernameHeader", "AUTH_HEADER_USERNAME_HEADER")
	_ = viper.BindEnv("auth.header.emailHeader", "AUTH_HEADER_EMAIL_HEADER")

	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			logger.Ctx(context.Background()).Fatal(err.Error())
		}
	}

	var c Config
	err = viper.Unmarshal(&c)
	if err != nil {
		logger.Ctx(context.Background()).Fatal(err.Error())
	}

	return c
}
