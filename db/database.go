// Package db provides access to Thunderdome application data in postgres
package db

import (
	"context"
	"embed"
	"fmt"
	"time"

	"github.com/uptrace/opentelemetry-go-extra/otelzap"

	"github.com/XSAM/otelsql"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"

	"github.com/golang-migrate/migrate/v4/source/iofs"
	"go.uber.org/zap"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq" // necessary for postgres
	"github.com/microcosm-cc/bluemonday"
)

//go:embed migrations/*.sql
var fs embed.FS

// New runs db migrations, sets up a db connection pool
// and sets previously active users to false during startup
func New(AdminEmail string, config *Config, logger *otelzap.Logger) *Database {
	ctx := context.Background()
	dms, err := iofs.New(fs, "migrations")
	if err != nil {
		logger.Ctx(ctx).Fatal("error loading db migrations", zap.Error(err))
	}

	// Do this once for each unique policy, and use the policy for the life of the program
	// Policy creation/editing is not safe to use in multiple goroutines
	bmp := bluemonday.UGCPolicy()

	var d = &Database{
		// read environment variables and sets up database configuration values
		Config:              config,
		HTMLSanitizerPolicy: bmp,
		Logger:              logger,
	}

	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		d.Config.Host,
		d.Config.Port,
		d.Config.User,
		d.Config.Password,
		d.Config.Name,
		d.Config.SSLMode,
	)

	pdb, err := otelsql.Open("postgres", psqlInfo, otelsql.WithAttributes(
		semconv.DBSystemPostgreSQL,
	))
	if err != nil {
		d.Logger.Ctx(ctx).Fatal("error connecting to the database: ", zap.Error(err))
	}
	d.DB = pdb
	d.DB.SetMaxOpenConns(d.Config.MaxOpenConns)
	d.DB.SetMaxIdleConns(d.Config.MaxIdleConns)
	d.DB.SetConnMaxLifetime(time.Duration(d.Config.ConnMaxLifetime) * time.Minute)

	err = otelsql.RegisterDBStatsMetrics(pdb, otelsql.WithAttributes(
		semconv.DBSystemPostgreSQL,
	))
	if err != nil {
		d.Logger.Ctx(ctx).Error("RegisterDBStatsMetrics error", zap.Error(err))
	}

	driver, err := postgres.WithInstance(pdb, &postgres.Config{})
	if err != nil {
		d.Logger.Ctx(ctx).Error("db driver error", zap.Error(err))
	}
	m, err := migrate.NewWithInstance(
		"iofs",
		dms,
		"postgres",
		driver)
	if err != nil {
		d.Logger.Error("new db migration instance", zap.Error(err))
	}
	if err := m.Up(); err != nil && err.Error() != "no change" {
		d.Logger.Ctx(ctx).Error("db migration up error", zap.Error(err))
	}

	// on server start reset all users to active false for battles
	if _, err := d.DB.Exec(
		`call deactivate_all_users();`); err != nil {
		d.Logger.Ctx(ctx).Error("call deactivate_all_users error", zap.Error(err))
	}

	// on server start if admin email is specified set that user to admin type
	if AdminEmail != "" {
		if _, err := d.DB.Exec(
			`call promote_user_by_email($1);`,
			AdminEmail,
		); err != nil {
			d.Logger.Ctx(ctx).Error("call promote_user_by_email error", zap.Error(err), zap.String("admin_email", AdminEmail))
		}
	}

	return d
}
