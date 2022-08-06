package db

import (
	"context"
	"errors"

	"github.com/StevenWeathers/thunderdome-planning-poker/model"
	"go.uber.org/zap"
)

// TeamUserRole gets a user's role in team
func (d *Database) TeamUserRole(ctx context.Context, UserID string, TeamID string) (string, error) {
	var teamRole string

	err := d.db.QueryRowContext(ctx,
		`SELECT role FROM team_get_user_role($1, $2)`,
		UserID,
		TeamID,
	).Scan(
		&teamRole,
	)
	if err != nil {
		d.logger.Ctx(ctx).Error("team_get_user_role query error", zap.Error(err))
		return "", errors.New("error getting team users role")
	}

	return teamRole, nil
}

// TeamGet gets an team
func (d *Database) TeamGet(ctx context.Context, TeamID string) (*model.Team, error) {
	var team = &model.Team{}

	err := d.db.QueryRowContext(ctx,
		`SELECT id, name, created_date, updated_date FROM team_get_by_id($1)`,
		TeamID,
	).Scan(
		&team.Id,
		&team.Name,
		&team.CreatedDate,
		&team.UpdatedDate,
	)
	if err != nil {
		d.logger.Ctx(ctx).Error("team_get_by_id query error", zap.Error(err))
		return nil, errors.New("team not found")
	}

	return team, nil
}

// TeamListByUser gets a list of teams the user is on
func (d *Database) TeamListByUser(ctx context.Context, UserID string, Limit int, Offset int) []*model.Team {
	var teams = make([]*model.Team, 0)
	rows, err := d.db.QueryContext(ctx,
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
				d.logger.Ctx(ctx).Error("team_list_by_user query scan error", zap.Error(err))
			} else {
				teams = append(teams, &team)
			}
		}
	} else {
		d.logger.Ctx(ctx).Error("team_list_by_user query error", zap.Error(err))
	}

	return teams
}

// TeamCreate creates a team with current user as an ADMIN
func (d *Database) TeamCreate(ctx context.Context, UserID string, TeamName string) (*model.Team, error) {
	t := &model.Team{}
	err := d.db.QueryRowContext(ctx, `
		SELECT id, name, created_date, updated_date FROM team_create($1, $2);`,
		UserID,
		TeamName,
	).Scan(&t.Id, &t.Name, &t.CreatedDate, &t.UpdatedDate)

	if err != nil {
		d.logger.Ctx(ctx).Error("team_create query error", zap.Error(err))
		return nil, err
	}

	return t, nil
}

// TeamAddUser adds a user to a team
func (d *Database) TeamAddUser(ctx context.Context, TeamID string, UserID string, Role string) (string, error) {
	_, err := d.db.ExecContext(ctx,
		`SELECT team_user_add($1, $2, $3);`,
		TeamID,
		UserID,
		Role,
	)

	if err != nil {
		d.logger.Ctx(ctx).Error("team_user_add query error", zap.Error(err))
		return "", err
	}

	return TeamID, nil
}

// TeamUserList gets a list of team users
func (d *Database) TeamUserList(ctx context.Context, TeamID string, Limit int, Offset int) ([]*model.TeamUser, int, error) {
	var users = make([]*model.TeamUser, 0)
	var userCount int

	err := d.db.QueryRowContext(ctx,
		`SELECT count(user_id) FROM team_user WHERE team_id = $1;`,
		TeamID,
	).Scan(&userCount)
	if err != nil {
		return nil, 0, err
	}

	if userCount == 0 {
		return users, userCount, nil
	}

	rows, err := d.db.QueryContext(ctx,
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
				d.logger.Ctx(ctx).Error("team_user_list query scan error", zap.Error(err))
			} else {
				usr.GravatarHash = createGravatarHash(usr.Email)
				users = append(users, &usr)
			}
		}
	} else {
		d.logger.Ctx(ctx).Error("team_user_list query error", zap.Error(err))
		return nil, 0, err
	}

	return users, userCount, nil
}

// TeamRemoveUser removes a user from a team
func (d *Database) TeamRemoveUser(ctx context.Context, TeamID string, UserID string) error {
	_, err := d.db.ExecContext(ctx,
		`CALL team_user_remove($1, $2);`,
		TeamID,
		UserID,
	)

	if err != nil {
		d.logger.Ctx(ctx).Error("team_user_remove query error", zap.Error(err))
		return err
	}

	return nil
}

// TeamBattleList gets a list of team battles
func (d *Database) TeamBattleList(ctx context.Context, TeamID string, Limit int, Offset int) []*model.Battle {
	var battles = make([]*model.Battle, 0)
	rows, err := d.db.QueryContext(ctx,
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
				d.logger.Ctx(ctx).Error("team_battle_list query scan error", zap.Error(err))
			} else {
				battles = append(battles, &tb)
			}
		}
	} else {
		d.logger.Ctx(ctx).Error("team_battle_list query error", zap.Error(err))
	}

	return battles
}

// TeamAddBattle adds a battle to a team
func (d *Database) TeamAddBattle(ctx context.Context, TeamID string, BattleID string) error {
	_, err := d.db.ExecContext(ctx,
		`SELECT team_battle_add($1, $2);`,
		TeamID,
		BattleID,
	)

	if err != nil {
		d.logger.Ctx(ctx).Error("team_battle_add query error", zap.Error(err))
		return err
	}

	return nil
}

// TeamRemoveBattle removes a battle from a team
func (d *Database) TeamRemoveBattle(ctx context.Context, TeamID string, BattleID string) error {
	_, err := d.db.ExecContext(ctx,
		`SELECT team_battle_remove($1, $2);`,
		TeamID,
		BattleID,
	)

	if err != nil {
		d.logger.Ctx(ctx).Error("team_battle_remove query error", zap.Error(err))
		return err
	}

	return nil
}

// TeamDelete deletes a team
func (d *Database) TeamDelete(ctx context.Context, TeamID string) error {
	_, err := d.db.ExecContext(ctx,
		`CALL team_delete($1);`,
		TeamID,
	)

	if err != nil {
		d.logger.Ctx(ctx).Error("team_delete query error", zap.Error(err))
		return err
	}

	return nil
}

// TeamRetroList gets a list of team retros
func (d *Database) TeamRetroList(ctx context.Context, TeamID string, Limit int, Offset int) []*model.Retro {
	var retros = make([]*model.Retro, 0)
	rows, err := d.db.QueryContext(ctx,
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
				d.logger.Ctx(ctx).Error("team_retro_list query scan error", zap.Error(err))
			} else {
				retros = append(retros, &tb)
			}
		}
	} else {
		d.logger.Ctx(ctx).Error("team_retro_list query error", zap.Error(err))
	}

	return retros
}

// TeamAddRetro adds a retro to a team
func (d *Database) TeamAddRetro(ctx context.Context, TeamID string, RetroID string) error {
	_, err := d.db.ExecContext(ctx,
		`SELECT team_retro_add($1, $2);`,
		TeamID,
		RetroID,
	)

	if err != nil {
		d.logger.Ctx(ctx).Error("team_retro_add query error", zap.Error(err))
		return err
	}

	return nil
}

// TeamRemoveRetro removes a retro from a team
func (d *Database) TeamRemoveRetro(ctx context.Context, TeamID string, RetroID string) error {
	_, err := d.db.ExecContext(ctx,
		`SELECT team_retro_remove($1, $2);`,
		TeamID,
		RetroID,
	)

	if err != nil {
		d.logger.Ctx(ctx).Error("team_retro_remove query error", zap.Error(err))
		return err
	}

	return nil
}

// TeamStoryboardList gets a list of team storyboards
func (d *Database) TeamStoryboardList(ctx context.Context, TeamID string, Limit int, Offset int) []*model.Storyboard {
	var storyboards = make([]*model.Storyboard, 0)
	rows, err := d.db.QueryContext(ctx,
		`SELECT id, name FROM team_storyboard_list($1, $2, $3);`,
		TeamID,
		Limit,
		Offset,
	)

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var tb model.Storyboard

			if err := rows.Scan(
				&tb.Id,
				&tb.Name,
			); err != nil {
				d.logger.Ctx(ctx).Error("team_storyboard_list query scan error", zap.Error(err))
			} else {
				storyboards = append(storyboards, &tb)
			}
		}
	} else {
		d.logger.Ctx(ctx).Error("team_storyboard_list query error", zap.Error(err))
	}

	return storyboards
}

// TeamAddStoryboard adds a storyboard to a team
func (d *Database) TeamAddStoryboard(ctx context.Context, TeamID string, StoryboardID string) error {
	_, err := d.db.ExecContext(ctx,
		`SELECT team_storyboard_add($1, $2);`,
		TeamID,
		StoryboardID,
	)

	if err != nil {
		d.logger.Ctx(ctx).Error("team_storyboard_add query error", zap.Error(err))
		return err
	}

	return nil
}

// TeamRemoveStoryboard removes a storyboard from a team
func (d *Database) TeamRemoveStoryboard(ctx context.Context, TeamID string, StoryboardID string) error {
	_, err := d.db.ExecContext(ctx,
		`SELECT team_storyboard_remove($1, $2);`,
		TeamID,
		StoryboardID,
	)

	if err != nil {
		d.logger.Ctx(ctx).Error("team_storyboard_remove query error", zap.Error(err))
		return err
	}

	return nil
}
