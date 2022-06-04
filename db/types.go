package db

import (
	"database/sql"

	"github.com/microcosm-cc/bluemonday"
	"go.uber.org/zap"
)

// Config holds all the configuration for the db
type Config struct {
	Host       string
	Port       int
	User       string
	Password   string
	Name       string
	SSLMode    string
	AESHashkey string
}

// Database contains all the methods to interact with DB
type Database struct {
	config              *Config
	db                  *sql.DB
	htmlSanitizerPolicy *bluemonday.Policy
	logger              *zap.Logger
}
