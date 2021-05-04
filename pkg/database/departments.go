package database

import (
	"log"
)

// OrganizationDepartmentList gets a list of organization departments
func (d *Database) OrganizationDepartmentList(OrgID string, Limit int, Offset int) []*Department {
	var departments = make([]*Department, 0)
	rows, err := d.db.Query(
		`SELECT id, name, created_date, updated_date FROM organization_department_list($1, $2, $3);`,
		OrgID,
		Limit,
		Offset,
	)

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var department Department

			if err := rows.Scan(
				&department.DepartmentID,
				&department.Name,
				&department.CreatedDate,
				&department.UpdatedDate,
			); err != nil {
				log.Println(err)
			} else {
				departments = append(departments, &department)
			}
		}
	} else {
		log.Println(err)
	}

	return departments
}

// DepartmentCreate creates an organization department
func (d *Database) DepartmentCreate(OrgID string, OrgName string) (string, error) {
	var DepartmentID string
	err := d.db.QueryRow(`
		SELECT departmentId FROM organization_department_create($1, $2);`,
		OrgID,
		OrgName,
	).Scan(&DepartmentID)

	if err != nil {
		log.Println("Unable to create organization department: ", err)
		return "", err
	}

	return DepartmentID, nil
}
