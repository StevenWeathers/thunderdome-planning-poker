package model

import "time"

type Team struct {
	Id          string    `json:"id"`
	Name        string    `json:"name"`
	CreatedDate time.Time `json:"createdDate"`
	UpdatedDate time.Time `json:"updatedDate"`
}

type TeamUser struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	Role         string `json:"role"`
	Avatar       string `json:"avatar"`
	GravatarHash string `json:"gravatarHash"`
}

type TeamCheckin struct {
	Id          string    `json:"id"`
	User        *TeamUser `json:"user"`
	Yesterday   string    `json:"yesterday"`
	Today       string    `json:"today"`
	Blockers    string    `json:"blockers"`
	Discuss     string    `json:"discuss"`
	GoalsMet    bool      `json:"goalsMet"`
	CreatedDate string    `json:"createdDate"`
	UpdatedDate string    `json:"updatedDate"`
}
