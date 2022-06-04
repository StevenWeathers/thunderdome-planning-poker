package db

import (
	"errors"
	"strings"

	"github.com/StevenWeathers/thunderdome-planning-poker/model"
	"go.uber.org/zap"
)

// GetAppStats gets counts of common application metrics such as users and battles
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
			apikey_count,
			active_battle_count,
			active_battle_user_count,
			team_checkins_count,
			retro_count,
			active_retro_count,
			active_retro_user_count,
			retro_item_count,
			retro_action_count
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
		&Appstats.ActiveBattleCount,
		&Appstats.ActiveBattleUserCount,
		&Appstats.TeamCheckinsCount,
		&Appstats.RetroCount,
		&Appstats.ActiveRetroCount,
		&Appstats.ActiveRetroUserCount,
		&Appstats.RetroItemCount,
		&Appstats.RetroActionCount,
	)
	if err != nil {
		d.logger.Error("Unable to get application stats", zap.Error(err))
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
		d.logger.Error("call promote_user error", zap.Error(err))
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
		d.logger.Error("call demote_user error", zap.Error(err))
		return errors.New("error attempting to demote user to registered")
	}

	return nil
}

// CleanBattles deletes battles older than {DaysOld} days
func (d *Database) CleanBattles(DaysOld int) error {
	if _, err := d.db.Exec(
		`call clean_battles($1);`,
		DaysOld,
	); err != nil {
		d.logger.Error("call clean_battles", zap.Error(err))
		return errors.New("error attempting to clean battles")
	}

	return nil
}

// CleanGuests deletes guest users older than {DaysOld} days
func (d *Database) CleanGuests(DaysOld int) error {
	if _, err := d.db.Exec(
		`call clean_guest_users($1);`,
		DaysOld,
	); err != nil {
		d.logger.Error("call clean_guest_users", zap.Error(err))
		return errors.New("error attempting to clean Guest Users")
	}

	return nil
}

// LowercaseUserEmails goes through and lower cases any user email that has uppercase letters
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
				d.logger.Error("lowercase_unique_user_emails scan error", zap.Error(err))
				return nil, err
			} else {
				users = append(users, &usr)
			}
		}
	} else {
		d.logger.Error("lowercase_unique_user_emails query error", zap.Error(err))
		return nil, err
	}

	return users, nil
}

// MergeDuplicateAccounts goes through and merges user accounts with duplicate emails that has uppercase letters
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
				d.logger.Error("merge_nonunique_user_accounts scan error", zap.Error(err))
				return nil, err
			} else {
				users = append(users, &usr)
			}
		}
	} else {
		d.logger.Error("merge_nonunique_user_accounts query error", zap.Error(err))
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
				d.logger.Error("organization_list scan error", zap.Error(err))
			} else {
				organizations = append(organizations, &org)
			}
		}
	} else {
		d.logger.Error("organization_list query error", zap.Error(err))
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
				d.logger.Error("team_list scan error", zap.Error(err))
			} else {
				teams = append(teams, &team)
			}
		}
	} else {
		d.logger.Error("team_list query error", zap.Error(err))
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
				&ak.UserId,
				&ak.Active,
				&ak.CreatedDate,
				&ak.UpdatedDate,
			); err != nil {
				d.logger.Error("apikeys_list scan error", zap.Error(err))
			} else {
				splitKey := strings.Split(key, ".")
				ak.Prefix = splitKey[0]
				ak.Id = key
				APIKeys = append(APIKeys, &ak)
			}
		}
	} else {
		d.logger.Error("apikeys_list query error", zap.Error(err))
	}

	return APIKeys
}
