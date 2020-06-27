package database

import (
	"database/sql"
	"time"
)

// Config holds all the configuration for the db
type Config struct {
	host     string
	port     int
	user     string
	password string
	dbname   string
}

// Database contains all the methods to interact with DB
type Database struct {
	config *Config
	db     *sql.DB
}

// BattleWarrior aka user
type BattleWarrior struct {
	WarriorID   string `json:"id"`
	WarriorName string `json:"name"`
	WarriorRank string `json:"rank"`
	Active      bool   `json:"active"`
}

// Battle aka arena
type Battle struct {
	BattleID           string           `json:"id"`
	LeaderID           string           `json:"leaderId"`
	BattleName         string           `json:"name"`
	Warriors           []*BattleWarrior `json:"warriors"`
	Plans              []*Plan          `json:"plans"`
	VotingLocked       bool             `json:"votingLocked"`
	ActivePlanID       string           `json:"activePlanId"`
	PointValuesAllowed []string         `json:"pointValuesAllowed"`
}

// Warrior aka user
type Warrior struct {
	WarriorID    string `json:"id"`
	WarriorName  string `json:"name"`
	WarriorEmail string `json:"email"`
	WarriorRank  string `json:"rank"`
	Verified     bool   `json:"verified"`
}

// Vote structure
type Vote struct {
	WarriorID string `json:"warriorId"`
	VoteValue string `json:"vote"`
}

// Plan aka Story structure
type Plan struct {
	PlanID        string    `json:"id"`
	PlanName      string    `json:"name"`
	Votes         []*Vote   `json:"votes"`
	Points        string    `json:"points"`
	PlanActive    bool      `json:"active"`
	PlanSkipped   bool      `json:"skipped"`
	VoteStartTime time.Time `json:"voteStartTime"`
	VoteEndTime   time.Time `json:"voteEndTime"`
}
