package thunderdome

import "time"

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
