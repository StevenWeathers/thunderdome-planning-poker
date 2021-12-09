package db

import (
	"database/sql"
)

// Config holds all the configuration for the db
type Config struct {
	Host       string
	Port       int
	User       string
	Password   string
	Name       string
	SSLMode    string
	AESHashkey string
}

// Database contains all the methods to interact with DB
type Database struct {
	config *Config
	db     *sql.DB
}
