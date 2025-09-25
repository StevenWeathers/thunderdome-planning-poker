package cmd

import (
	"log/slog"
	"os"

	"github.com/spf13/cobra"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"

	"github.com/StevenWeathers/thunderdome-planning-poker/internal/config"
	"github.com/StevenWeathers/thunderdome-planning-poker/internal/db"
)

// migrateCmd represents the migrate command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migrations",
	Long:  `Run database migrations to set up or update the database schema.`,
	Run:   runMigrate,
}

var migrateUpCmd = &cobra.Command{
	Use:   "up",
	Short: "Run all pending migrations",
	Run:   runMigrateUp,
}

var migrateDownCmd = &cobra.Command{
	Use:   "down",
	Short: "Rollback the last migration",
	Run:   runMigrateDown,
}

var migrateStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show migration status",
	Run:   runMigrateStatus,
}

func init() {
	RootCmd.AddCommand(migrateCmd)
	migrateCmd.AddCommand(migrateUpCmd)
	migrateCmd.AddCommand(migrateDownCmd)
	migrateCmd.AddCommand(migrateStatusCmd)
}

func runMigrate(cmd *cobra.Command, args []string) {
	cmd.Help()
}

func runMigrateUp(cmd *cobra.Command, args []string) {
	version := RootCmd.Version
	zlog, _ := zap.NewProduction(
		zap.Fields(
			zap.String("version", version),
		),
	)
	defer func() {
		_ = zlog.Sync()
	}()
	logger := otelzap.New(zlog)

	c := config.InitConfig(logger)

	d := db.New(c.Admin.Email, &db.Config{
		Host:                   c.Db.Host,
		Port:                   c.Db.Port,
		User:                   c.Db.User,
		Password:               c.Db.Pass,
		Name:                   c.Db.Name,
		SSLMode:                c.Db.Sslmode,
		AESHashkey:             c.Config.AesHashkey,
		MaxIdleConns:           c.Db.MaxIdleConns,
		MaxOpenConns:           c.Db.MaxOpenConns,
		ConnMaxLifetime:        c.Db.ConnMaxLifetime,
		DefaultEstimationScale: c.Config.AllowedPointValues,
	}, logger, true)

	if err := d.MigrateUp(); err != nil {
		slog.Error("Failed to run migrations", slog.Any("error", err))
		os.Exit(1)
	}

	slog.Info("Migrations completed successfully")
}

func runMigrateDown(cmd *cobra.Command, args []string) {
	version := RootCmd.Version
	zlog, _ := zap.NewProduction(
		zap.Fields(
			zap.String("version", version),
		),
	)
	defer func() {
		_ = zlog.Sync()
	}()
	logger := otelzap.New(zlog)

	c := config.InitConfig(logger)

	d := db.New(c.Admin.Email, &db.Config{
		Host:                   c.Db.Host,
		Port:                   c.Db.Port,
		User:                   c.Db.User,
		Password:               c.Db.Pass,
		Name:                   c.Db.Name,
		SSLMode:                c.Db.Sslmode,
		AESHashkey:             c.Config.AesHashkey,
		MaxIdleConns:           c.Db.MaxIdleConns,
		MaxOpenConns:           c.Db.MaxOpenConns,
		ConnMaxLifetime:        c.Db.ConnMaxLifetime,
		DefaultEstimationScale: c.Config.AllowedPointValues,
	}, logger, true)

	if err := d.MigrateDown(); err != nil {
		slog.Error("Failed to rollback migration", slog.Any("error", err))
		os.Exit(1)
	}

	slog.Info("Migration rollback completed successfully")
}

func runMigrateStatus(cmd *cobra.Command, args []string) {
	version := RootCmd.Version
	zlog, _ := zap.NewProduction(
		zap.Fields(
			zap.String("version", version),
		),
	)
	defer func() {
		_ = zlog.Sync()
	}()
	logger := otelzap.New(zlog)

	c := config.InitConfig(logger)

	d := db.New(c.Admin.Email, &db.Config{
		Host:                   c.Db.Host,
		Port:                   c.Db.Port,
		User:                   c.Db.User,
		Password:               c.Db.Pass,
		Name:                   c.Db.Name,
		SSLMode:                c.Db.Sslmode,
		AESHashkey:             c.Config.AesHashkey,
		MaxIdleConns:           c.Db.MaxIdleConns,
		MaxOpenConns:           c.Db.MaxOpenConns,
		ConnMaxLifetime:        c.Db.ConnMaxLifetime,
		DefaultEstimationScale: c.Config.AllowedPointValues,
	}, logger, true)

	if err := d.MigrateStatus(); err != nil {
		slog.Error("Failed to get migration status", slog.Any("error", err))
		os.Exit(1)
	}
}
