// Package db provides access to Thunderdome application data in postgres
package db

import (
	"database/sql"
	"embed"
	"fmt"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq" // necessary for postgres
	"github.com/spf13/viper"
)

//go:embed migrations/*.sql
var fs embed.FS

//go:embed procs_functions.sql
var procsFuncsSQL string

// New runs db migrations, sets up a db connection pool
// and sets previously active users to false during startup
func New(AdminEmail string) *Database {
	dms, err := iofs.New(fs, "migrations")
	if err != nil {
		log.Fatal(err)
	}
	var d = &Database{
		// read environment variables and sets up database configuration values
		config: &Config{
			host:     viper.GetString("db.host"),
			port:     viper.GetInt("db.port"),
			user:     viper.GetString("db.user"),
			password: viper.GetString("db.pass"),
			dbname:   viper.GetString("db.name"),
			sslmode:  viper.GetString("db.sslmode"),
		},
	}

	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		d.config.host,
		d.config.port,
		d.config.user,
		d.config.password,
		d.config.dbname,
		d.config.sslmode,
	)

	pdb, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("error connecting to the database: ", err)
	}
	d.db = pdb

	driver, err := postgres.WithInstance(pdb, &postgres.Config{})
	m, err := migrate.NewWithInstance(
		"iofs",
		dms,
		"postgres",
		driver)
	if err := m.Up(); err != nil && err.Error() != "no change" {
		log.Fatal(err)
	}

	if _, err := d.db.Exec(procsFuncsSQL); err != nil {
		log.Fatal(err)
	}

	// on server start reset all users to active false for battles
	if _, err := d.db.Exec(
		`call deactivate_all_users();`); err != nil {
		log.Println(err)
	}

	// on server start if admin email is specified set that user to admin type
	if AdminEmail != "" {
		if _, err := d.db.Exec(
			`call promote_user_by_email($1);`,
			AdminEmail,
		); err != nil {
			log.Println(err)
		}
	}

	return d
}
