package model

import "time"

// User aka user
type User struct {
	UserID               string `json:"id"`
	UserName             string `json:"name"`
	UserEmail            string `json:"email"`
	UserType             string `json:"rank"`
	UserAvatar           string `json:"avatar"`
	Verified             bool   `json:"verified"`
	NotificationsEnabled bool   `json:"notificationsEnabled"`
	Country              string `json:"country"`
	Locale               string `json:"locale"`
	Company              string `json:"company"`
	JobTitle             string `json:"jobTitle"`
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
