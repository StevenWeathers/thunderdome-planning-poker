package jira

import (
	"database/sql"

	"github.com/uptrace/opentelemetry-go-extra/otelzap"
)

type Service struct {
	DB         *sql.DB
	Logger     *otelzap.Logger
	AESHashKey string
}
