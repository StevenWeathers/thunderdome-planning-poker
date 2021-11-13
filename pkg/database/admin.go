package database

import (
	"errors"
	"log"
	"strings"

	"github.com/StevenWeathers/thunderdome-planning-poker/model"
)

// ConfirmAdmin confirms whether the users is infact an admin
func (d *Database) ConfirmAdmin(UserID string) error {
	var UserType string
	e := d.db.QueryRow("SELECT coalesce(type, '') FROM users WHERE id = $1;", UserID).Scan(&UserType)
	if e != nil {
		log.Println(e)
		return errors.New("could not find users type")
	}

	if UserType != "GENERAL" {
		return errors.New(("user is not an admin"))
	}

	return nil
}

// GetAppStats gets counts of users (registered and unregistered), battles, and plans
func (d *Database) GetAppStats() (*model.ApplicationStats, error) {
	var Appstats model.ApplicationStats

	err := d.db.QueryRow(`
		SELECT
			unregistered_user_count,
			registered_user_count,
			battle_count,
			plan_count,
			organization_count,
			department_count,
			team_count,
			apikey_count
		FROM get_app_stats();
		`,
	).Scan(
		&Appstats.UnregisteredCount,
		&Appstats.RegisteredCount,
		&Appstats.BattleCount,
		&Appstats.PlanCount,
		&Appstats.OrganizationCount,
		&Appstats.DepartmentCount,
		&Appstats.TeamCount,
		&Appstats.APIKeyCount,
	)
	if err != nil {
		log.Println("Unable to get application stats: ", err)
		return nil, err
	}

	return &Appstats, nil
}

// PromoteUser promotes a user to admin type
func (d *Database) PromoteUser(UserID string) error {
	if _, err := d.db.Exec(
		`call promote_user($1);`,
		UserID,
	); err != nil {
		log.Println(err)
		return errors.New("error attempting to promote user to admin")
	}

	return nil
}

// DemoteUser demotes a user to registered type
func (d *Database) DemoteUser(UserID string) error {
	if _, err := d.db.Exec(
		`call demote_user($1);`,
		UserID,
	); err != nil {
		log.Println(err)
		return errors.New("error attempting to demote user to registered")
	}

	return nil
}

// CleanBattles deletes battles older than X days
func (d *Database) CleanBattles(DaysOld int) error {
	if _, err := d.db.Exec(
		`call clean_battles($1);`,
		DaysOld,
	); err != nil {
		log.Println(err)
		return errors.New("error attempting to clean battles")
	}

	return nil
}

// CleanGuests deletes guest users older than X days
func (d *Database) CleanGuests(DaysOld int) error {
	if _, err := d.db.Exec(
		`call clean_guest_users($1);`,
		DaysOld,
	); err != nil {
		log.Println(err)
		return errors.New("error attempting to clean Guest Users")
	}

	return nil
}

// LowercaseUserEmails goes through and lowercases any user email that has any uppercase letters
// returning the list of updated users
func (d *Database) LowercaseUserEmails() ([]*model.User, error) {
	var users = make([]*model.User, 0)
	rows, err := d.db.Query(
		`SELECT name, email FROM lowercase_unique_user_emails();`,
	)

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var usr model.User

			if err := rows.Scan(
				&usr.Name,
				&usr.Email,
			); err != nil {
				log.Println(err)
				return nil, err
			} else {
				users = append(users, &usr)
			}
		}
	} else {
		log.Println(err)
		return nil, err
	}

	return users, nil
}

// MergeDuplicateAccounts goes through and merges any user accounts with duplicate emails that has any uppercase letters
// returning the list of merged users
func (d *Database) MergeDuplicateAccounts() ([]*model.User, error) {
	var users = make([]*model.User, 0)
	rows, err := d.db.Query(
		`SELECT name, email FROM merge_nonunique_user_accounts();`,
	)

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var usr model.User

			if err := rows.Scan(
				&usr.Name,
				&usr.Email,
			); err != nil {
				log.Println(err)
				return nil, err
			} else {
				users = append(users, &usr)
			}
		}
	} else {
		log.Println(err)
		return nil, err
	}

	return users, nil
}

// OrganizationList gets a list of organizations
func (d *Database) OrganizationList(Limit int, Offset int) []*model.Organization {
	var organizations = make([]*model.Organization, 0)
	rows, err := d.db.Query(
		`SELECT id, name, created_date, updated_date FROM organization_list($1, $2);`,
		Limit,
		Offset,
	)

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var org model.Organization

			if err := rows.Scan(
				&org.Id,
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

// TeamList gets a list of teams
func (d *Database) TeamList(Limit int, Offset int) []*model.Team {
	var teams = make([]*model.Team, 0)
	rows, err := d.db.Query(
		`SELECT id, name, created_date, updated_date FROM team_list($1, $2);`,
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

// GetAPIKeys gets a list of api keys
func (d *Database) GetAPIKeys(Limit int, Offset int) []*model.APIKey {
	var APIKeys = make([]*model.APIKey, 0)
	rows, err := d.db.Query(
		`SELECT id, name, email, active, created_date, updated_date
		FROM apikeys_list($1, $2);`,
		Limit,
		Offset,
	)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var ak model.APIKey
			var key string

			if err := rows.Scan(
				&key,
				&ak.Name,
				&ak.UserID,
				&ak.Active,
				&ak.CreatedDate,
				&ak.UpdatedDate,
			); err != nil {
				log.Println(err)
			} else {
				splitKey := strings.Split(key, ".")
				ak.Prefix = splitKey[0]
				ak.ID = key
				APIKeys = append(APIKeys, &ak)
			}
		}
	}

	return APIKeys
}
