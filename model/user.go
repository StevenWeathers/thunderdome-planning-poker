package model

import "time"

// User aka user
type User struct {
	Id                   string    `json:"id"`
	Name                 string    `json:"name"`
	Email                string    `json:"email"`
	Type                 string    `json:"rank"`
	Avatar               string    `json:"avatar"`
	Verified             bool      `json:"verified"`
	NotificationsEnabled bool      `json:"notificationsEnabled"`
	Country              string    `json:"country"`
	Locale               string    `json:"locale"`
	Company              string    `json:"company"`
	JobTitle             string    `json:"jobTitle"`
	GravatarHash         string    `json:"gravatarHash"`
	CreatedDate          time.Time `json:"createdDate"`
	UpdatedDate          time.Time `json:"updatedDate"`
	LastActive           time.Time `json:"lastActive"`
	Disabled             bool      `json:"disabled"`
	MFAEnabled           bool      `json:"mfaEnabled"`
}

// APIKey structure
type APIKey struct {
	Id          string    `json:"id"`
	Prefix      string    `json:"prefix"`
	UserId      string    `json:"userId"`
	Name        string    `json:"name"`
	Key         string    `json:"apiKey"`
	Active      bool      `json:"active"`
	CreatedDate time.Time `json:"createdDate"`
	UpdatedDate time.Time `json:"updatedDate"`
}
