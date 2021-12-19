package db

import (
	"errors"
	"log"

	"github.com/StevenWeathers/thunderdome-planning-poker/model"
)

// TeamUserRole gets a user's role in team
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
func (d *Database) TeamGet(TeamID string) (*model.Team, error) {
	var team = &model.Team{}

	e := d.db.QueryRow(
		`SELECT id, name, created_date, updated_date FROM team_get_by_id($1)`,
		TeamID,
	).Scan(
		&team.Id,
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
func (d *Database) TeamListByUser(UserID string, Limit int, Offset int) []*model.Team {
	var teams = make([]*model.Team, 0)
	rows, err := d.db.Query(
		`SELECT id, name, created_date, updated_date FROM team_list_by_user($1, $2, $3);`,
		UserID,
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

// TeamCreate creates a team with current user as an ADMIN
func (d *Database) TeamCreate(UserID string, TeamName string) (*model.Team, error) {
	t := &model.Team{}
	err := d.db.QueryRow(`
		SELECT id, name, created_date, updated_date FROM team_create($1, $2);`,
		UserID,
		TeamName,
	).Scan(&t.Id, &t.Name, &t.CreatedDate, &t.UpdatedDate)

	if err != nil {
		log.Println("Unable to create team: ", err)
		return nil, err
	}

	return t, nil
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
func (d *Database) TeamUserList(TeamID string, Limit int, Offset int) ([]*model.TeamUser, int, error) {
	var users = make([]*model.TeamUser, 0)
	var userCount int

	err := d.db.QueryRow(
		`SELECT count(user_id) FROM team_user WHERE team_id = $1;`,
		TeamID,
	).Scan(&userCount)
	if err != nil {
		return nil, 0, err
	}

	if userCount == 0 {
		return users, userCount, nil
	}

	rows, err := d.db.Query(
		`SELECT id, name, email, role, avatar FROM team_user_list($1, $2, $3);`,
		TeamID,
		Limit,
		Offset,
	)

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var usr model.TeamUser

			if err = rows.Scan(
				&usr.Id,
				&usr.Name,
				&usr.Email,
				&usr.Role,
				&usr.Avatar,
			); err != nil {
				log.Println(err)
			} else {
				usr.GravatarHash = createGravatarHash(usr.Email)
				users = append(users, &usr)
			}
		}
	} else {
		return nil, 0, err
	}

	return users, userCount, nil
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
func (d *Database) TeamBattleList(TeamID string, Limit int, Offset int) []*model.Battle {
	var battles = make([]*model.Battle, 0)
	rows, err := d.db.Query(
		`SELECT id, name FROM team_battle_list($1, $2, $3);`,
		TeamID,
		Limit,
		Offset,
	)

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var tb model.Battle

			if err := rows.Scan(
				&tb.Id,
				&tb.Name,
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

// TeamRetroList gets a list of team retros
func (d *Database) TeamRetroList(TeamID string, Limit int, Offset int) []*model.Retro {
	var retros = make([]*model.Retro, 0)
	rows, err := d.db.Query(
		`SELECT id, name, format, phase FROM team_retro_list($1, $2, $3);`,
		TeamID,
		Limit,
		Offset,
	)

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var tb model.Retro

			if err := rows.Scan(
				&tb.Id,
				&tb.Name,
				&tb.Format,
				&tb.Phase,
			); err != nil {
				log.Println(err)
			} else {
				retros = append(retros, &tb)
			}
		}
	} else {
		log.Println(err)
	}

	return retros
}

// TeamAddRetro adds a retro to a team
func (d *Database) TeamAddRetro(TeamID string, RetroID string) error {
	_, err := d.db.Exec(
		`SELECT team_retro_add($1, $2);`,
		TeamID,
		RetroID,
	)

	if err != nil {
		log.Println("Unable to add retro to team: ", err)
		return err
	}

	return nil
}

// TeamRemoveRetro removes a retro from a team
func (d *Database) TeamRemoveRetro(TeamID string, RetroID string) error {
	_, err := d.db.Exec(
		`SELECT team_retro_remove($1, $2);`,
		TeamID,
		RetroID,
	)

	if err != nil {
		log.Println("Unable to remove retro from team: ", err)
		return err
	}

	return nil
}
