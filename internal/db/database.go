// Package db provides access to Thunderdome application data in postgres
package db

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"time"

	"github.com/uptrace/opentelemetry-go-extra/otelzap"

	"github.com/XSAM/otelsql"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"

	"go.uber.org/zap"

	_ "github.com/jackc/pgx/v5/stdlib" // necessary for postgres
	"github.com/microcosm-cc/bluemonday"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var fs embed.FS

var migrationsDir = "migrations"

// New runs db migrations, sets up a db connection pool
// and sets previously active users to false during startup
func New(adminEmail string, config *Config, logger *otelzap.Logger, instantiate bool) *Service {
	ctx := context.Background()

	// Do this once for each unique policy, and use the policy for the life of the program
	// Policy creation/editing is not safe to use in multiple goroutines
	bmp := bluemonday.UGCPolicy()

	var d = &Service{
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

	pdb := waitForDB(ctx, logger, psqlInfo)

	d.DB = pdb
	d.DB.SetMaxOpenConns(d.Config.MaxOpenConns)
	d.DB.SetMaxIdleConns(d.Config.MaxIdleConns)
	d.DB.SetConnMaxLifetime(time.Duration(d.Config.ConnMaxLifetime) * time.Minute)

	_, err := otelsql.RegisterDBStatsMetrics(pdb, otelsql.WithAttributes(
		semconv.DBSystemPostgreSQL,
	))
	if err != nil {
		d.Logger.Ctx(ctx).Error("RegisterDBStatsMetrics error", zap.Error(err))
	}

	gl := newGooseLogger(logger)
	goose.SetLogger(gl)
	goose.SetBaseFS(fs)

	if err := goose.SetDialect("postgres"); err != nil {
		d.Logger.Ctx(ctx).Error("goose set postgres dialect error", zap.Error(err))
	}

	if instantiate {
		d.Instantiate(ctx, adminEmail)
	}

	return d
}

func (d *Service) Instantiate(ctx context.Context, adminEmail string) {
	if err := goose.Up(d.DB, migrationsDir, goose.WithAllowMissing()); err != nil {
		d.Logger.Ctx(ctx).Error("migrations error", zap.Error(err))
	}

	// on server start reset all users to active false for games
	if _, err := d.DB.Exec(
		`CALL thunderdome.users_deactivate_all();`); err != nil {
		d.Logger.Ctx(ctx).Error("CALL thunderdome.deactivate_all_users error", zap.Error(err))
	}

	// on server start if admin email is specified set that user to admin type
	if adminEmail != "" {
		if _, err := d.DB.Exec(
			`UPDATE thunderdome.users SET type = 'ADMIN', updated_date = NOW() WHERE email = $1;`,
			adminEmail,
		); err != nil {
			d.Logger.Ctx(ctx).Error("CALL thunderdome.promote_user_by_email error", zap.Error(err), zap.String("admin_email", adminEmail))
		}
	}

	// backwards compatibility for self-hosted instances with custom estimation scale configured
	// will be removed in v5 and documented in the release notes
	if len(d.Config.DefaultEstimationScale) > 0 {
		if _, err := d.DB.Exec(
			`UPDATE thunderdome.estimation_scale SET values = $1 WHERE scale_type = 'thunderdome_default' AND
		values <> $1;`,
			d.Config.DefaultEstimationScale,
		); err != nil {
			d.Logger.Ctx(ctx).Error("failed to update thunderdome_default estimation scale", zap.Error(err))
		}
	}
}

// waitForDB attempts to open a database connection and ping it until successful.
// This is to ensure that the database is ready before proceeding with application startup and migrations.
// The function will retry every 2 seconds until a successful connection is made or the context is cancelled.
func waitForDB(ctx context.Context, logger *otelzap.Logger, dsn string) *sql.DB {
	var db *sql.DB
	for {
		var err error
		db, err = otelsql.Open("pgx", dsn, otelsql.WithAttributes(
			semconv.DBSystemPostgreSQL,
		))
		if err != nil {
			logger.Ctx(ctx).Info("Waiting for DB open", zap.Error(err))
			time.Sleep(2 * time.Second)
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		err = db.PingContext(ctx)
		cancel()
		if err != nil {
			logger.Ctx(ctx).Info("Waiting for DB ping", zap.Error(err))
			db.Close()
			time.Sleep(2 * time.Second)
			continue
		}

		break
	}
	logger.Ctx(ctx).Info("Database connection established.")
	return db
}
