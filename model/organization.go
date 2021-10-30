package model

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
