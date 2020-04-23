package database

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	_ "github.com/lib/pq" // necessary for postgres
	"github.com/markbates/pkger"
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

// GetEnv gets environment variable matching key string
// and if it finds none uses fallback string
// returning either the matching or fallback string
func GetEnv(key string, fallback string) string {
	var result = os.Getenv(key)

	if result == "" {
		result = fallback
	}

	return result
}

// GetIntEnv gets an environment variable and converts it to an int
// and if it finds none uses fallback
func GetIntEnv(key string, fallback int) int {
	var intResult = fallback
	var stringResult = os.Getenv(key)

	if stringResult != "" {
		v, _ := strconv.Atoi(stringResult)
		intResult = v
	}

	return intResult
}

// GetBoolEnv gets an environment variable and converts it to a bool
// and if it finds none uses fallback
func GetBoolEnv(key string, fallback bool) bool {
	var boolResult = fallback
	var stringResult = os.Getenv(key)

	if stringResult != "" {
		b, _ := strconv.ParseBool(stringResult)
		boolResult = b
	}

	return boolResult
}

// New runs db migrations, sets up a db connection pool
// and sets previously active warriors to false during startup
func New(AdminEmail string) *Database {
	var d = &Database{
		// read environment variables and sets up database configuration values
		config: &Config{
			host:     GetEnv("DB_HOST", "db"),
			port:     GetIntEnv("DB_PORT", 5432),
			user:     GetEnv("DB_USER", "thor"),
			password: GetEnv("DB_PASS", "odinson"),
			dbname:   GetEnv("DB_NAME", "thunderdome"),
		},
	}

	sqlFile, ioErr := pkger.Open("/schema.sql")
	if ioErr != nil {
		log.Println("Error reading schema.sql file required to migrate db")
		log.Fatal(ioErr)
	}
	sqlContent, ioErr := ioutil.ReadAll(sqlFile)
	if ioErr != nil {
		// this will hopefully only possibly panic during development as the file is already in memory otherwise
		log.Println("Error reading schema.sql file required to migrate db")
		log.Fatal(ioErr)
	}
	migrationSQL := string(sqlContent)

	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		d.config.host,
		d.config.port,
		d.config.user,
		d.config.password,
		d.config.dbname,
	)

	pdb, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("error connecting to the database: ", err)
	}
	d.db = pdb

	if _, err := d.db.Exec(migrationSQL); err != nil {
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
