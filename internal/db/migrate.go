package db

import "github.com/pressly/goose/v3"

// MigrateUp runs all up migrations using goose
func (s *Service) MigrateUp() error {
	return goose.Up(s.DB, migrationsDir, goose.WithAllowMissing())
}

// MigrateDown rolls back the most recent migration using goose
func (s *Service) MigrateDown() error {
	return goose.Down(s.DB, migrationsDir, goose.WithAllowMissing())
}

// MigrateStatus prints the status of all migrations
func (s *Service) MigrateStatus() error {
	return goose.Status(s.DB, migrationsDir, goose.WithAllowMissing())
}

// MigrateVersion returns the current migration version
func (s *Service) MigrateVersion() (int64, error) {
	return goose.GetDBVersion(s.DB)
}
