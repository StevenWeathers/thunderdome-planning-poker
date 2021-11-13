package model

type Team struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	CreatedDate string `json:"createdDate"`
	UpdatedDate string `json:"updatedDate"`
}

type TeamUser struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}
