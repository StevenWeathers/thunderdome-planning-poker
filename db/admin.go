package db

import (
	"context"
	"errors"
	"strings"

	"github.com/StevenWeathers/thunderdome-planning-poker/model"
	"go.uber.org/zap"
)

// GetAppStats gets counts of common application metrics such as users and battles
func (d *Database) GetAppStats(ctx context.Context) (*model.ApplicationStats, error) {
	var Appstats model.ApplicationStats

	err := d.db.QueryRowContext(ctx, `
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
			retro_action_count,
			storyboard_count,
			active_storyboard_count,
			active_storyboard_user_count,
			storyboard_goal_count,
			storyboard_column_count,
			storyboard_story_count,
			storyboard_persona_count
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
		&Appstats.StoryboardCount,
		&Appstats.ActiveStoryboardCount,
		&Appstats.ActiveStoryboardUserCount,
		&Appstats.StoryboardGoalCount,
		&Appstats.StoryboardColumnCount,
		&Appstats.StoryboardStoryCount,
		&Appstats.StoryboardPersonaCount,
	)
	if err != nil {
		d.logger.Ctx(ctx).Error("Unable to get application stats", zap.Error(err))
		return nil, err
	}

	return &Appstats, nil
}

// PromoteUser promotes a user to admin type
func (d *Database) PromoteUser(ctx context.Context, UserID string) error {
	if _, err := d.db.ExecContext(ctx,
		`call promote_user($1);`,
		UserID,
	); err != nil {
		d.logger.Ctx(ctx).Error("call promote_user error", zap.Error(err))
		return errors.New("error attempting to promote user to admin")
	}

	return nil
}

// DemoteUser demotes a user to registered type
func (d *Database) DemoteUser(ctx context.Context, UserID string) error {
	if _, err := d.db.ExecContext(ctx,
		`call demote_user($1);`,
		UserID,
	); err != nil {
		d.logger.Ctx(ctx).Error("call demote_user error", zap.Error(err))
		return errors.New("error attempting to demote user to registered")
	}

	return nil
}

// DisableUser disables a user from logging in
func (d *Database) DisableUser(ctx context.Context, UserID string) error {
	if _, err := d.db.ExecContext(ctx,
		`call user_disable($1);`,
		UserID,
	); err != nil {
		d.logger.Ctx(ctx).Error("call user_disable error", zap.Error(err))
		return errors.New("error attempting to disable user")
	}

	return nil
}

// EnableUser enables a user allowing login
func (d *Database) EnableUser(ctx context.Context, UserID string) error {
	if _, err := d.db.ExecContext(ctx,
		`call user_enable($1);`,
		UserID,
	); err != nil {
		d.logger.Ctx(ctx).Error("call user_enable error", zap.Error(err))
		return errors.New("error attempting to enable user")
	}

	return nil
}

// CleanBattles deletes battles older than {DaysOld} days
func (d *Database) CleanBattles(ctx context.Context, DaysOld int) error {
	if _, err := d.db.ExecContext(ctx,
		`call clean_battles($1);`,
		DaysOld,
	); err != nil {
		d.logger.Ctx(ctx).Error("call clean_battles", zap.Error(err))
		return errors.New("error attempting to clean battles")
	}

	return nil
}

// CleanRetros deletes retros older than {DaysOld} days
func (d *Database) CleanRetros(ctx context.Context, DaysOld int) error {
	if _, err := d.db.ExecContext(ctx,
		`call clean_retros($1);`,
		DaysOld,
	); err != nil {
		d.logger.Ctx(ctx).Error("call clean_retros", zap.Error(err))
		return errors.New("error attempting to clean retros")
	}

	return nil
}

// CleanStoryboards deletes storyboards older than {DaysOld} days
func (d *Database) CleanStoryboards(ctx context.Context, DaysOld int) error {
	if _, err := d.db.ExecContext(ctx,
		`call clean_storyboards($1);`,
		DaysOld,
	); err != nil {
		d.logger.Ctx(ctx).Error("call clean_storyboards", zap.Error(err))
		return errors.New("error attempting to clean storyboards")
	}

	return nil
}

// CleanGuests deletes guest users older than {DaysOld} days
func (d *Database) CleanGuests(ctx context.Context, DaysOld int) error {
	if _, err := d.db.ExecContext(ctx,
		`call clean_guest_users($1);`,
		DaysOld,
	); err != nil {
		d.logger.Ctx(ctx).Error("call clean_guest_users", zap.Error(err))
		return errors.New("error attempting to clean Guest Users")
	}

	return nil
}

// LowercaseUserEmails goes through and lower cases any user email that has uppercase letters
// returning the list of updated users
func (d *Database) LowercaseUserEmails(ctx context.Context) ([]*model.User, error) {
	var users = make([]*model.User, 0)
	rows, err := d.db.QueryContext(ctx,
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
				d.logger.Ctx(ctx).Error("lowercase_unique_user_emails scan error", zap.Error(err))
				return nil, err
			} else {
				users = append(users, &usr)
			}
		}
	} else {
		d.logger.Ctx(ctx).Error("lowercase_unique_user_emails query error", zap.Error(err))
		return nil, err
	}

	return users, nil
}

// MergeDuplicateAccounts goes through and merges user accounts with duplicate emails that has uppercase letters
// returning the list of merged users
func (d *Database) MergeDuplicateAccounts(ctx context.Context) ([]*model.User, error) {
	var users = make([]*model.User, 0)
	rows, err := d.db.QueryContext(ctx,
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
				d.logger.Ctx(ctx).Error("merge_nonunique_user_accounts scan error", zap.Error(err))
				return nil, err
			} else {
				users = append(users, &usr)
			}
		}
	} else {
		d.logger.Ctx(ctx).Error("merge_nonunique_user_accounts query error", zap.Error(err))
		return nil, err
	}

	return users, nil
}

// OrganizationList gets a list of organizations
func (d *Database) OrganizationList(ctx context.Context, Limit int, Offset int) []*model.Organization {
	var organizations = make([]*model.Organization, 0)
	rows, err := d.db.QueryContext(ctx,
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
				d.logger.Ctx(ctx).Error("organization_list scan error", zap.Error(err))
			} else {
				organizations = append(organizations, &org)
			}
		}
	} else {
		d.logger.Ctx(ctx).Error("organization_list query error", zap.Error(err))
	}

	return organizations
}

// TeamList gets a list of teams
func (d *Database) TeamList(ctx context.Context, Limit int, Offset int) []*model.Team {
	var teams = make([]*model.Team, 0)
	rows, err := d.db.QueryContext(ctx,
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
				d.logger.Ctx(ctx).Error("team_list scan error", zap.Error(err))
			} else {
				teams = append(teams, &team)
			}
		}
	} else {
		d.logger.Ctx(ctx).Error("team_list query error", zap.Error(err))
	}

	return teams
}

// GetAPIKeys gets a list of api keys
func (d *Database) GetAPIKeys(ctx context.Context, Limit int, Offset int) []*model.APIKey {
	var APIKeys = make([]*model.APIKey, 0)
	rows, err := d.db.QueryContext(ctx,
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
				d.logger.Ctx(ctx).Error("apikeys_list scan error", zap.Error(err))
			} else {
				splitKey := strings.Split(key, ".")
				ak.Prefix = splitKey[0]
				ak.Id = key
				APIKeys = append(APIKeys, &ak)
			}
		}
	} else {
		d.logger.Ctx(ctx).Error("apikeys_list query error", zap.Error(err))
	}

	return APIKeys
}
