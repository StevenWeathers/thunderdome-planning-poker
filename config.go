package main

import (
	"log"

	"github.com/spf13/viper"
)

// InitConfig initializes the application configuration
func InitConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AddConfigPath("/etc/thunderdome/")
	viper.AddConfigPath("$HOME/.config/thunderdome/")
	viper.AddConfigPath(".")

	viper.SetDefault("http.cookie_hashkey", "strongest-avenger")
	viper.SetDefault("http.port", "8080")
	viper.SetDefault("http.secure_cookie", true)
	viper.SetDefault("http.backend_cookie_name", "warriorId")
	viper.SetDefault("http.frontend_cookie_name", "warrior")
	viper.SetDefault("http.domain", "thunderdome.dev")
	viper.SetDefault("http.path_prefix", "")

	viper.SetDefault("analytics.enabled", true)
	viper.SetDefault("analytics.id", "UA-140245309-1")

	viper.SetDefault("db.host", "db")
	viper.SetDefault("db.port", 5432)
	viper.SetDefault("db.user", "thor")
	viper.SetDefault("db.pass", "odinson")
	viper.SetDefault("db.name", "thunderdome")
	viper.SetDefault("db.sslmode", "disable")

	viper.SetDefault("smtp.host", "localhost")
	viper.SetDefault("smtp.port", "25")
	viper.SetDefault("smtp.secure", true)
	viper.SetDefault("smtp.sender", "no-reply@thunderdome.dev")

	viper.SetDefault("config.allowedPointValues",
		[]string{"0", "1/2", "1", "2", "3", "5", "8", "13", "20", "40", "100", "?"})
	viper.SetDefault("config.defaultPointValues",
		[]string{"1", "2", "3", "5", "8", "13", "?"})
	viper.SetDefault("config.show_warrior_rank", false)
	viper.SetDefault("config.avatar_service", "goadorable")
	viper.SetDefault("config.toast_timeout", 1000)
	viper.SetDefault("config.allow_guests", true)
	viper.SetDefault("config.allow_registration", true)
	viper.SetDefault("config.allow_jira_import", true)
	viper.SetDefault("config.default_locale", "en")
	viper.SetDefault("config.friendly_ui_verbs", false)
	viper.SetDefault("config.allow_external_api", false)

	viper.SetDefault("auth.method", "normal")
	viper.SetDefault("auth.ldap.url", "")
	viper.SetDefault("auth.ldap.use_tls", true)
	viper.SetDefault("auth.ldap.bindname", "")
	viper.SetDefault("auth.ldap.bindpass", "")
	viper.SetDefault("auth.ldap.basedn", "")
	viper.SetDefault("auth.ldap.filter", "(&(objectClass=posixAccount)(mail=%s))")
	viper.SetDefault("auth.ldap.mail_attr", "mail")
	viper.SetDefault("auth.ldap.cn_attr", "cn")

	viper.BindEnv("http.cookie_hashkey", "COOKIE_HASHKEY")
	viper.BindEnv("http.port", "PORT")
	viper.BindEnv("http.secure_cookie", "COOKIE_SECURE")
	viper.BindEnv("http.backend_cookie_name", "SECURE_COOKIE_NAME")
	viper.BindEnv("http.frontend_cookie_name", "FRONTEND_COOKIE_NAME")
	viper.BindEnv("http.domain", "APP_DOMAIN")
	viper.BindEnv("http.path_prefix", "PATH_PREFIX")

	viper.BindEnv("analytics.enabled", "ANALYTICS_ENABLED")
	viper.BindEnv("analytics.id", "ANALYTICS_ID")
	viper.BindEnv("admin.email", "ADMIN_EMAIL")

	viper.BindEnv("db.host", "DB_HOST")
	viper.BindEnv("db.port", "DB_PORT")
	viper.BindEnv("db.user", "DB_USER")
	viper.BindEnv("db.pass", "DB_PASS")
	viper.BindEnv("db.name", "DB_NAME")
	viper.BindEnv("db.sslmode", "DB_SSLMODE")

	viper.BindEnv("smtp.host", "SMTP_HOST")
	viper.BindEnv("smtp.port", "SMTP_PORT")
	viper.BindEnv("smtp.secure", "SMTP_SECURE")
	viper.BindEnv("smtp.identity", "SMTP_IDENTITY")
	viper.BindEnv("smtp.user", "SMTP_USER")
	viper.BindEnv("smtp.pass", "SMTP_PASS")
	viper.BindEnv("smtp.sender", "SMTP_SENDER")

	viper.BindEnv("config.allowedPointValues", "CONFIG_POINTS_ALLOWED")
	viper.BindEnv("config.defaultPointValues", "CONFIG_POINTS_DEFAULT")
	viper.BindEnv("config.show_warrior_rank", "CONFIG_SHOW_RANK")
	viper.BindEnv("config.avatar_service", "CONFIG_AVATAR_SERVICE")
	viper.BindEnv("config.toast_timeout", "CONFIG_TOAST_TIMEOUT")
	viper.BindEnv("config.allow_guests", "CONFIG_ALLOW_GUESTS")
	viper.BindEnv("config.allow_registration", "CONFIG_ALLOW_REGISTRATION")
	viper.BindEnv("config.allow_jira_import", "CONFIG_ALLOW_JIRA_IMPORT")
	viper.BindEnv("config.default_locale", "CONFIG_DEFAULT_LOCALE")
	viper.BindEnv("config.friendly_ui_verbs", "CONFIG_FRIENDLY_UI_VERBS")
	viper.BindEnv("config.allow_external_api", "CONFIG_ALLOW_EXTERNAL_API")

	viper.BindEnv("auth.method", "AUTH_METHOD")
	viper.BindEnv("auth.ldap.url", "AUTH_LDAP_URL")
	viper.BindEnv("auth.ldap.use_tls", "AUTH_LDAP_USE_TLS")
	viper.BindEnv("auth.ldap.bindname", "AUTH_LDAP_BINDNAME")
	viper.BindEnv("auth.ldap.bindpass", "AUTH_LDAP_BINDPASS")
	viper.BindEnv("auth.ldap.basedn", "AUTH_LDAP_BASEDN")
	viper.BindEnv("auth.ldap.filter", "AUTH_LDAP_FILTER")
	viper.BindEnv("auth.ldap.mail_attr", "AUTH_LDAP_MAIL_ATTR")
	viper.BindEnv("auth.ldap.cn_attr", "AUTH_LDAP_CN_ATTR")

	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Fatal(err)
		}
	}
}
