package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq" // necessary for postgres
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

// HashAndSalt takes a password byte and salt + hashes it
// returning a hash string to store in db
func HashAndSalt(pwd []byte) (string, error) {
	// Use GenerateFromPassword to hash & salt pwd.
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
		return "", err
	}
	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash), nil
}

// ComparePasswords takes a password hash and compares it to entered password bytes
// returning true if matches false if not
func ComparePasswords(hashedPwd string, plainPwd []byte) bool {
	// Since we'll be getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

// New runs db migrations, sets up a db connection pool
// and sets previously active warriors to false during startup
func New(AdminEmail string, schemaSQL string) *Database {
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

	if _, err := d.db.Exec(schemaSQL); err != nil {
		log.Fatal(err)
	}

	// on server start reset all warriors to active false for battles
	if _, err := d.db.Exec(
		`call deactivate_all_warriors();`); err != nil {
		log.Println(err)
	}

	// on server start if admin email is specified set that warrior to GENERAL rank
	if AdminEmail != "" {
		if _, err := d.db.Exec(
			`call promote_warrior_by_email($1);`,
			AdminEmail,
		); err != nil {
			log.Println(err)
		}
	}

	return d
}
