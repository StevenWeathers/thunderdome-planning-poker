package model

// ApplicationStats includes user, organization, team, battle, and plan counts
type ApplicationStats struct {
	RegisteredCount       int `json:"registeredUserCount"`
	UnregisteredCount     int `json:"unregisteredUserCount"`
	BattleCount           int `json:"battleCount"`
	PlanCount             int `json:"planCount"`
	OrganizationCount     int `json:"organizationCount"`
	DepartmentCount       int `json:"departmentCount"`
	TeamCount             int `json:"teamCount"`
	APIKeyCount           int `json:"apikeyCount"`
	ActiveBattleCount     int `json:"activeBattleCount"`
	ActiveBattleUserCount int `json:"activeBattleUserCount"`
}

type Alert struct {
	AlertID        string `json:"id" db:"id"`
	Name           string `json:"name" db:"name"`
	Type           string `json:"type" db:"type"`
	Content        string `json:"content" db:"content"`
	Active         bool   `json:"active" db:"active"`
	AllowDismiss   bool   `json:"allowDismiss" db:"allow_dismiss"`
	RegisteredOnly bool   `json:"registeredOnly" db:"registered_only"`
	CreatedDate    string `json:"createdDate" db:"created_date"`
	UpdatedDate    string `json:"updatedDate" db:"updated_date"`
}
