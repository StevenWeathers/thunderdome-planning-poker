package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/uptrace/opentelemetry-go-extra/otelzap"

	"github.com/microcosm-cc/bluemonday"
)

// Config holds all the configuration for the db
type Config struct {
	Host                   string
	Port                   int
	User                   string
	Password               string
	Name                   string
	SSLMode                string
	AESHashkey             string
	MaxOpenConns           int
	MaxIdleConns           int
	ConnMaxLifetime        int
	DefaultEstimationScale []string
}

// Service contains all the methods to interact with DB
type Service struct {
	Config              *Config
	DB                  *sql.DB
	HTMLSanitizerPolicy *bluemonday.Policy
	Logger              *otelzap.Logger
}

type gooseLogger struct {
	logger *otelzap.Logger
}

func (l *gooseLogger) Fatalf(format string, v ...any) {
	l.logger.Ctx(context.Background()).Fatal(fmt.Sprintf(format, v...))
}
func (l *gooseLogger) Printf(format string, v ...any) {
	l.logger.Ctx(context.Background()).Info(fmt.Sprintf(format, v...))
}

func newGooseLogger(logger *otelzap.Logger) *gooseLogger {
	return &gooseLogger{logger: logger}
}
