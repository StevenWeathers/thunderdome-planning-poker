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

// BattleUser aka user
type BattleUser struct {
	UserID     string `json:"id"`
	UserName   string `json:"name"`
	UserType   string `json:"rank"`
	UserAvatar string `json:"avatar"`
	Active     bool   `json:"active"`
	Abandoned  bool   `json:"abandoned"`
}

// Battle aka arena
type Battle struct {
	BattleID             string        `json:"id"`
	BattleName           string        `json:"name"`
	Users                []*BattleUser `json:"warriors"`
	Plans                []*Plan       `json:"plans"`
	VotingLocked         bool          `json:"votingLocked"`
	ActivePlanID         string        `json:"activePlanId"`
	PointValuesAllowed   []string      `json:"pointValuesAllowed"`
	AutoFinishVoting     bool          `json:"autoFinishVoting"`
	Leaders              []string      `json:"leaders"`
	PointAverageRounding string        `json:"pointAverageRounding"`
}

// User aka user
type User struct {
	UserID               string `json:"id"`
	UserName             string `json:"name"`
	UserEmail            string `json:"email"`
	UserType             string `json:"rank"`
	UserAvatar           string `json:"avatar"`
	Verified             bool   `json:"verified"`
	NotificationsEnabled bool   `json:"notificationsEnabled"`
}

// Vote structure
type Vote struct {
	UserID    string `json:"warriorId"`
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
	UserID      string    `json:"warriorId"`
	Name        string    `json:"name"`
	Key         string    `json:"apiKey"`
	Active      bool      `json:"active"`
	CreatedDate time.Time `json:"createdDate"`
	UpdatedDate time.Time `json:"updatedDate"`
}

// ApplicationStats includes user, organization, team, battle, and plan counts
type ApplicationStats struct {
	RegisteredCount   int `json:"registeredUserCount"`
	UnregisteredCount int `json:"unregisteredUserCount"`
	BattleCount       int `json:"battleCount"`
	PlanCount         int `json:"planCount"`
	OrganizationCount int `json:"organizationCount"`
	DepartmentCount   int `json:"departmentCount"`
	TeamCount         int `json:"teamCount"`
}

// Organization can be a company
type Organization struct {
	OrganizationID string `json:"id"`
	Name           string `json:"name"`
	CreatedDate    string `json:"createdDate"`
	UpdatedDate    string `json:"updatedDate"`
}

type OrganizationUser struct {
	UserID string `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Role   string `json:"role"`
}

type Department struct {
	DepartmentID string `json:"id"`
	Name         string `json:"name"`
	CreatedDate  string `json:"createdDate"`
	UpdatedDate  string `json:"updatedDate"`
}

type DepartmentUser struct {
	UserID string `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Role   string `json:"role"`
}

type Team struct {
	TeamID      string `json:"id"`
	Name        string `json:"name"`
	CreatedDate string `json:"createdDate"`
	UpdatedDate string `json:"updatedDate"`
}

type TeamUser struct {
	UserID string `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Role   string `json:"role"`
}
