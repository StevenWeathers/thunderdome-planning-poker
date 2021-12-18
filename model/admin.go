package model

import "time"

// ApplicationStats includes counts of different data points of the application
type ApplicationStats struct {
	UnregisteredCount     int `json:"unregisteredUserCount"`
	RegisteredCount       int `json:"registeredUserCount"`
	APIKeyCount           int `json:"apikeyCount"`
	BattleCount           int `json:"battleCount"`
	ActiveBattleCount     int `json:"activeBattleCount"`
	ActiveBattleUserCount int `json:"activeBattleUserCount"`
	PlanCount             int `json:"planCount"`
	OrganizationCount     int `json:"organizationCount"`
	DepartmentCount       int `json:"departmentCount"`
	TeamCount             int `json:"teamCount"`
	TeamCheckinsCount     int `json:"teamCheckinsCount"`
}

type Alert struct {
	Id             string    `json:"id" db:"id"`
	Name           string    `json:"name" db:"name"`
	Type           string    `json:"type" db:"type"`
	Content        string    `json:"content" db:"content"`
	Active         bool      `json:"active" db:"active"`
	AllowDismiss   bool      `json:"allowDismiss" db:"allow_dismiss"`
	RegisteredOnly bool      `json:"registeredOnly" db:"registered_only"`
	CreatedDate    time.Time `json:"createdDate" db:"created_date"`
	UpdatedDate    time.Time `json:"updatedDate" db:"updated_date"`
}
