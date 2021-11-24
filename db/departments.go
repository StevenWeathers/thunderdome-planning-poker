package db

import (
	"errors"
	"log"

	"github.com/StevenWeathers/thunderdome-planning-poker/model"
)

// DepartmentUserRole gets a users role in department (and organization)
func (d *Database) DepartmentUserRole(UserID string, OrgID string, DepartmentID string) (string, string, error) {
	var orgRole string
	var departmentRole string

	e := d.db.QueryRow(
		`SELECT orgRole, departmentRole FROM department_get_user_role($1, $2, $3)`,
		UserID,
		OrgID,
		DepartmentID,
	).Scan(
		&orgRole,
		&departmentRole,
	)
	if e != nil {
		log.Println(e)
		return "", "", errors.New("error getting department users role")
	}

	return orgRole, departmentRole, nil
}

// DepartmentGet gets a department
func (d *Database) DepartmentGet(DepartmentID string) (*model.Department, error) {
	var org = &model.Department{}

	e := d.db.QueryRow(
		`SELECT id, name, created_date, updated_date FROM department_get_by_id($1)`,
		DepartmentID,
	).Scan(
		&org.Id,
		&org.Name,
		&org.CreatedDate,
		&org.UpdatedDate,
	)
	if e != nil {
		log.Println(e)
		return nil, errors.New("department not found")
	}

	return org, nil
}

// OrganizationDepartmentList gets a list of organization departments
func (d *Database) OrganizationDepartmentList(OrgID string, Limit int, Offset int) []*model.Department {
	var departments = make([]*model.Department, 0)
	rows, err := d.db.Query(
		`SELECT id, name, created_date, updated_date FROM department_list($1, $2, $3);`,
		OrgID,
		Limit,
		Offset,
	)

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var department model.Department

			if err := rows.Scan(
				&department.Id,
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
func (d *Database) DepartmentCreate(OrgID string, OrgName string) (*model.Department, error) {
	od := &model.Department{}

	err := d.db.QueryRow(`
		SELECT id, name, created_date, updated_date FROM department_create($1, $2);`,
		OrgID,
		OrgName,
	).Scan(&od.Id, &od.Name, &od.CreatedDate, &od.UpdatedDate)

	if err != nil {
		log.Println("Unable to create organization department: ", err)
		return nil, err
	}

	return od, nil
}

// DepartmentTeamList gets a list of department teams
func (d *Database) DepartmentTeamList(DepartmentID string, Limit int, Offset int) []*model.Team {
	var teams = make([]*model.Team, 0)
	rows, err := d.db.Query(
		`SELECT id, name, created_date, updated_date FROM department_team_list($1, $2, $3);`,
		DepartmentID,
		Limit,
		Offset,
	)

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var team model.Team

			if err := rows.Scan(
				&team.Id,
				&team.Name,
				&team.CreatedDate,
				&team.UpdatedDate,
			); err != nil {
				log.Println(err)
			} else {
				teams = append(teams, &team)
			}
		}
	} else {
		log.Println(err)
	}

	return teams
}

// DepartmentTeamCreate creates a department team
func (d *Database) DepartmentTeamCreate(DepartmentID string, TeamName string) (*model.Team, error) {
	t := &model.Team{}

	err := d.db.QueryRow(`
		SELECT id, name, created_date, updated_date FROM department_team_create($1, $2);`,
		DepartmentID,
		TeamName,
	).Scan(&t.Id, &t.Name, &t.CreatedDate, &t.UpdatedDate)

	if err != nil {
		log.Println("Unable to create department team: ", err)
		return nil, err
	}

	return t, nil
}

// DepartmentUserList gets a list of department users
func (d *Database) DepartmentUserList(DepartmentID string, Limit int, Offset int) []*model.DepartmentUser {
	var users = make([]*model.DepartmentUser, 0)
	rows, err := d.db.Query(
		`SELECT id, name, email, role FROM department_user_list($1, $2, $3);`,
		DepartmentID,
		Limit,
		Offset,
	)

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var usr model.DepartmentUser

			if err := rows.Scan(
				&usr.Id,
				&usr.Name,
				&usr.Email,
				&usr.Role,
			); err != nil {
				log.Println(err)
			} else {
				users = append(users, &usr)
			}
		}
	} else {
		log.Println(err)
	}

	return users
}

// DepartmentAddUser adds a user to an organization department
func (d *Database) DepartmentAddUser(DepartmentID string, UserID string, Role string) (string, error) {
	_, err := d.db.Exec(
		`SELECT department_user_add($1, $2, $3);`,
		DepartmentID,
		UserID,
		Role,
	)

	if err != nil {
		log.Println("Unable to add user to department: ", err)
		return "", err
	}

	return DepartmentID, nil
}

// DepartmentRemoveUser removes a user from a department (and department teams)
func (d *Database) DepartmentRemoveUser(DepartmentID string, UserID string) error {
	_, err := d.db.Exec(
		`CALL department_user_remove($1, $2);`,
		DepartmentID,
		UserID,
	)

	if err != nil {
		log.Println("Unable to remove user from department: ", err)
		return err
	}

	return nil
}

// DepartmentTeamUserRole gets a users role in organization department team
func (d *Database) DepartmentTeamUserRole(UserID string, OrgID string, DepartmentID string, TeamID string) (string, string, string, error) {
	var orgRole string
	var departmentRole string
	var teamRole string

	e := d.db.QueryRow(
		`SELECT orgRole, departmentRole, teamRole FROM department_team_user_role($1, $2, $3, $4)`,
		UserID,
		OrgID,
		DepartmentID,
		TeamID,
	).Scan(
		&orgRole,
		&departmentRole,
		&teamRole,
	)
	if e != nil {
		log.Println(e)
		return "", "", "", errors.New("error getting department team users role")
	}

	return orgRole, departmentRole, teamRole, nil
}
