package db

import (
	"database/sql"

	"github.com/uptrace/opentelemetry-go-extra/otelzap"

	"github.com/microcosm-cc/bluemonday"
)

// Config holds all the configuration for the db
type Config struct {
	Host            string
	Port            int
	User            string
	Password        string
	Name            string
	SSLMode         string
	AESHashkey      string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime int
}

// Service contains all the methods to interact with DB
type Service struct {
	Config              *Config
	DB                  *sql.DB
	HTMLSanitizerPolicy *bluemonday.Policy
	Logger              *otelzap.Logger
}
