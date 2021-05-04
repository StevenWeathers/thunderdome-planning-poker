package database

import (
	"errors"
	"log"
)

// OrganizationWithRole gets an organization with current users role
func (d *Database) OrganizationWithRole(UserID string, OrgID string) (*Organization, string, error) {
	var org = &Organization{
		OrganizationID: "",
		Name:           "",
		CreatedDate:    "",
		UpdatedDate:    "",
	}
	var role string

	e := d.db.QueryRow(
		`SELECT id, name, created_date, updated_date, role FROM organization_get_with_role($1, $2)`,
		UserID,
		OrgID,
	).Scan(
		&org.OrganizationID,
		&org.Name,
		&org.CreatedDate,
		&org.CreatedDate,
		&role,
	)
	if e != nil {
		log.Println(e)
		return nil, "", errors.New("organization not found")
	}

	return org, role, nil
}

// OrganizationList gets a list of organizations the user is apart of
func (d *Database) OrganizationListByUser(UserID string, Limit int, Offset int) []*Organization {
	var organizations = make([]*Organization, 0)
	rows, err := d.db.Query(
		`SELECT id, name, created_date, updated_date FROM organization_list_by_user($1, $2, $3);`,
		UserID,
		Limit,
		Offset,
	)

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var org Organization

			if err := rows.Scan(
				&org.OrganizationID,
				&org.Name,
				&org.CreatedDate,
				&org.UpdatedDate,
			); err != nil {
				log.Println(err)
			} else {
				organizations = append(organizations, &org)
			}
		}
	} else {
		log.Println(err)
	}

	return organizations
}

// OrganizationCreate creates an organization
func (d *Database) OrganizationCreate(UserID string, OrgName string) (string, error) {
	var OrgID string
	err := d.db.QueryRow(`
		SELECT organizationId FROM organization_create($1, $2);`,
		UserID,
		OrgName,
	).Scan(&OrgID)

	if err != nil {
		log.Println("Unable to create organization: ", err)
		return "", err
	}

	return OrgID, nil
}

// OrganizationUserList gets a list of organization users
func (d *Database) OrganizationUserList(OrgID string, Limit int, Offset int) []*OrganizationUser {
	var users = make([]*OrganizationUser, 0)
	rows, err := d.db.Query(
		`SELECT id, name, email, role FROM organization_user_list($1, $2, $3);`,
		OrgID,
		Limit,
		Offset,
	)

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var usr OrganizationUser

			if err := rows.Scan(
				&usr.UserID,
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

// OrganizationAddUser adds a user to an organization
func (d *Database) OrganizationAddUser(OrgID string, UserID string, Role string) (string, error) {
	_, err := d.db.Exec(
		`call organization_user_add($1, $2, $3);`,
		OrgID,
		UserID,
		Role,
	)

	if err != nil {
		log.Println("Unable to add user to organization: ", err)
		return "", err
	}

	return OrgID, nil
}

// OrganizationTeamList gets a list of organization teams
func (d *Database) OrganizationTeamList(OrgID string, Limit int, Offset int) []*Team {
	var teams = make([]*Team, 0)
	rows, err := d.db.Query(
		`SELECT id, name, created_date, updated_date FROM organization_team_list($1, $2, $3);`,
		OrgID,
		Limit,
		Offset,
	)

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var team Team

			if err := rows.Scan(
				&team.TeamID,
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