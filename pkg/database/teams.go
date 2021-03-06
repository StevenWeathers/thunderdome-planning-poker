package database

import (
	"errors"
	"log"
)

// TeamUserRole gets a users role in team
func (d *Database) TeamUserRole(UserID string, TeamID string) (string, error) {
	var teamRole string

	e := d.db.QueryRow(
		`SELECT role FROM team_get_user_role($1, $2)`,
		UserID,
		TeamID,
	).Scan(
		&teamRole,
	)
	if e != nil {
		log.Println(e)
		return "", errors.New("error getting team users role")
	}

	return teamRole, nil
}

// TeamGet gets an team
func (d *Database) TeamGet(TeamID string) (*Team, error) {
	var team = &Team{
		TeamID:      "",
		Name:        "",
		CreatedDate: "",
		UpdatedDate: "",
	}

	e := d.db.QueryRow(
		`SELECT id, name, created_date, updated_date FROM team_get_by_id($1)`,
		TeamID,
	).Scan(
		&team.TeamID,
		&team.Name,
		&team.CreatedDate,
		&team.UpdatedDate,
	)
	if e != nil {
		log.Println(e)
		return nil, errors.New("team not found")
	}

	return team, nil
}

// TeamListByUser gets a list of teams the user is on
func (d *Database) TeamListByUser(UserID string, Limit int, Offset int) []*Team {
	var teams = make([]*Team, 0)
	rows, err := d.db.Query(
		`SELECT id, name, created_date, updated_date FROM team_list_by_user($1, $2, $3);`,
		UserID,
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

// TTeamCreate creates a team with current user as an ADMIN
func (d *Database) TeamCreate(UserID string, TeamName string) (string, error) {
	var TeamID string
	err := d.db.QueryRow(`
		SELECT teamId FROM team_create($1, $2);`,
		UserID,
		TeamName,
	).Scan(&TeamID)

	if err != nil {
		log.Println("Unable to create team: ", err)
		return "", err
	}

	return TeamID, nil
}

// TeamAddUser adds a user to a team
func (d *Database) TeamAddUser(TeamID string, UserID string, Role string) (string, error) {
	_, err := d.db.Exec(
		`SELECT team_user_add($1, $2, $3);`,
		TeamID,
		UserID,
		Role,
	)

	if err != nil {
		log.Println("Unable to add user to team: ", err)
		return "", err
	}

	return TeamID, nil
}

// TeamUserList gets a list of team users
func (d *Database) TeamUserList(TeamID string, Limit int, Offset int) []*OrganizationUser {
	var users = make([]*OrganizationUser, 0)
	rows, err := d.db.Query(
		`SELECT id, name, email, role FROM team_user_list($1, $2, $3);`,
		TeamID,
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

// TeamRemoveUser removes a user from a team
func (d *Database) TeamRemoveUser(TeamID string, UserID string) error {
	_, err := d.db.Exec(
		`CALL team_user_remove($1, $2);`,
		TeamID,
		UserID,
	)

	if err != nil {
		log.Println("Unable to remove user from team: ", err)
		return err
	}

	return nil
}

// TeamBattleList gets a list of team battles
func (d *Database) TeamBattleList(TeamID string, Limit int, Offset int) []*Battle {
	var battles = make([]*Battle, 0)
	rows, err := d.db.Query(
		`SELECT id, name FROM team_battle_list($1, $2, $3);`,
		TeamID,
		Limit,
		Offset,
	)

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var tb Battle

			if err := rows.Scan(
				&tb.BattleID,
				&tb.BattleName,
			); err != nil {
				log.Println(err)
			} else {
				battles = append(battles, &tb)
			}
		}
	} else {
		log.Println(err)
	}

	return battles
}

// TeamAddBattle adds a battle to a team
func (d *Database) TeamAddBattle(TeamID string, BattleID string) error {
	_, err := d.db.Exec(
		`SELECT team_battle_add($1, $2);`,
		TeamID,
		BattleID,
	)

	if err != nil {
		log.Println("Unable to add battle to team: ", err)
		return err
	}

	return nil
}

// TeamRemoveBattle removes a battle from a team
func (d *Database) TeamRemoveBattle(TeamID string, BattleID string) error {
	_, err := d.db.Exec(
		`SELECT team_battle_remove($1, $2);`,
		TeamID,
		BattleID,
	)

	if err != nil {
		log.Println("Unable to remove battle from team: ", err)
		return err
	}

	return nil
}

// TeamDelete deletes a team
func (d *Database) TeamDelete(TeamID string) error {
	_, err := d.db.Exec(
		`CALL team_delete($1);`,
		TeamID,
	)

	if err != nil {
		log.Println("Unable to delete team: ", err)
		return err
	}

	return nil
}
