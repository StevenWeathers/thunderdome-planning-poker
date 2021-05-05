package database

import (
	"errors"
	"log"
)

// OrganizationGet gets an organization
func (d *Database) OrganizationGet(OrgID string) (*Organization, error) {
	var org = &Organization{
		OrganizationID: "",
		Name:           "",
		CreatedDate:    "",
		UpdatedDate:    "",
	}

	e := d.db.QueryRow(
		`SELECT id, name, created_date, updated_date FROM organization_get_by_id($1)`,
		OrgID,
	).Scan(
		&org.OrganizationID,
		&org.Name,
		&org.CreatedDate,
		&org.UpdatedDate,
	)
	if e != nil {
		log.Println(e)
		return nil, errors.New("error getting organization")
	}

	return org, nil
}

// OrganizationUserRole gets a users role in organization
func (d *Database) OrganizationUserRole(UserID string, OrgID string) (string, error) {
	var role string

	e := d.db.QueryRow(
		`SELECT role FROM organization_get_user_role($1, $2)`,
		UserID,
		OrgID,
	).Scan(
		&role,
	)
	if e != nil {
		log.Println(e)
		return "", errors.New("error getting organization users role")
	}

	return role, nil
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
		`SELECT organization_user_add($1, $2, $3);`,
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

// OrganizationRemoveUser removes a user from a organization
func (d *Database) OrganizationRemoveUser(OrganizationID string, UserID string) error {
	_, err := d.db.Exec(
		`CALL organization_user_remove($1, $2);`,
		OrganizationID,
		UserID,
	)

	if err != nil {
		log.Println("Unable to remove user from organization: ", err)
		return err
	}

	return nil
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

// OrganizationTeamCreate creates an organization team
func (d *Database) OrganizationTeamCreate(OrgID string, TeamName string) (string, error) {
	var TeamID string
	err := d.db.QueryRow(`
		SELECT teamId FROM organization_team_create($1, $2);`,
		OrgID,
		TeamName,
	).Scan(&TeamID)

	if err != nil {
		log.Println("Unable to create organization team: ", err)
		return "", err
	}

	return TeamID, nil
}

// OrganizationTeamUserRole gets a users role in organization team
func (d *Database) OrganizationTeamUserRole(UserID string, OrgID string, TeamID string) (string, string, error) {
	var orgRole string
	var teamRole string

	e := d.db.QueryRow(
		`SELECT orgRole, teamRole FROM organization_team_user_role($1, $2, $3)`,
		UserID,
		OrgID,
		TeamID,
	).Scan(
		&orgRole,
		&teamRole,
	)
	if e != nil {
		log.Println(e)
		return "", "", errors.New("error getting organization team users role")
	}

	return orgRole, teamRole, nil
}
