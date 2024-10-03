package jira

import (
	"database/sql"

	"github.com/uptrace/opentelemetry-go-extra/otelzap"
)

// Service represents the JIRA database service
type Service struct {
	DB         *sql.DB
	Logger     *otelzap.Logger
	AESHashKey string
}
