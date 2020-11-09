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
	sslmode  string
}

// Database contains all the methods to interact with DB
type Database struct {
	config *Config
	db     *sql.DB
}

// BattleWarrior aka user
type BattleWarrior struct {
	WarriorID     string `json:"id"`
	WarriorName   string `json:"name"`
	WarriorRank   string `json:"rank"`
	WarriorAvatar string `json:"avatar"`
	Active        bool   `json:"active"`
	Abandoned     bool   `json:"abandoned"`
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
	AutoFinishVoting   bool             `json:"autoFinishVoting"`
}

// Warrior aka user
type Warrior struct {
	WarriorID            string `json:"id"`
	WarriorName          string `json:"name"`
	WarriorEmail         string `json:"email"`
	WarriorRank          string `json:"rank"`
	WarriorAvatar        string `json:"avatar"`
	Verified             bool   `json:"verified"`
	NotificationsEnabled bool   `json:"notificationsEnabled"`
}

// Vote structure
type Vote struct {
	WarriorID string `json:"warriorId"`
	VoteValue string `json:"vote"`
}

// Plan aka Story structure
type Plan struct {
	PlanID             string    `json:"id"`
	PlanName           string    `json:"name"`
	Type               string    `json:"type"`
	ReferenceID        string    `json:"referenceId"`
	Link               string    `json:"link"`
	Description        string    `json:"description"`
	AcceptanceCriteria string    `json:"acceptanceCriteria"`
	Votes              []*Vote   `json:"votes"`
	Points             string    `json:"points"`
	PlanActive         bool      `json:"active"`
	PlanSkipped        bool      `json:"skipped"`
	VoteStartTime      time.Time `json:"voteStartTime"`
	VoteEndTime        time.Time `json:"voteEndTime"`
}

// APIKey structure
type APIKey struct {
	ID          string    `json:"id"`
	Prefix      string    `json:"prefix"`
	WarriorID   string    `json:"warriorId"`
	Name        string    `json:"name"`
	Key         string    `json:"apiKey"`
	Active      bool      `json:"active"`
	CreatedDate time.Time `json:"createdDate"`
	UpdatedDate time.Time `json:"updatedDate"`
}
