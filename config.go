package main

import (
	"log"
	"github.com/spf13/viper"
)

func InitConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AddConfigPath("/etc/thunderdome/")
	viper.AddConfigPath("$HOME/.config/thunderdome/")
	viper.AddConfigPath(".")

	viper.SetDefault("http.cookie_hashkey", "strongest-avenger")
	viper.SetDefault("http.port", "8080")
	viper.SetDefault("http.secure_cookie", true)
	viper.SetDefault("http.domain", "thunderdome.dev")
	viper.SetDefault("analytics.enabled", true)
	viper.SetDefault("analytics.id", "UA-140245309-1")
	viper.SetDefault("db.host", "db")
	viper.SetDefault("db.port", 5432)
	viper.SetDefault("db.user", "thor")
	viper.SetDefault("db.pass", "odinson")
	viper.SetDefault("db.name", "thunderdome")
	viper.SetDefault("smtp.host", "localhost")
	viper.SetDefault("smtp.port", "25")
	viper.SetDefault("smtp.secure", true)
	viper.SetDefault("smtp.sender", "no-reply@thunderdome.dev")

	viper.BindEnv("http.cookie_hashkey", "COOKIE_HASHKEY")
	viper.BindEnv("http.port", "PORT")
	viper.BindEnv("http.secure_cookie", "COOKIE_SECURE")
	viper.BindEnv("http.domain", "APP_DOMAIN")
	viper.BindEnv("analytics.enabled", "ANALYTICS_ENABLED")
	viper.BindEnv("analytics.id", "ANALYTICS_ID")
	viper.BindEnv("admin.email", "ADMIN_EMAIL")
	viper.BindEnv("db.host", "DB_HOST")
	viper.BindEnv("db.port", "DB_PORT")
	viper.BindEnv("db.user", "DB_USER")
	viper.BindEnv("db.pass", "DB_PASS")
	viper.BindEnv("db.name", "DB_NAME")
	viper.BindEnv("smtp.host", "SMTP_HOST")
	viper.BindEnv("smtp.port", "SMTP_PORT")
	viper.BindEnv("smtp.secure", "SMTP_SECURE")
	viper.BindEnv("smtp.identity", "SMTP_IDENTITY")
	viper.BindEnv("smtp.user", "SMTP_USER")
	viper.BindEnv("smtp.pass", "SMTP_PASS")
	viper.BindEnv("smtp.sender", "SMTP_SENDER")

	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Fatal(err)
		}
	}
}
