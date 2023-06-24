package db

import (
	"context"
	"database/sql"
	"errors"
	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"

	"go.uber.org/zap"
)

// TeamService represents a PostgreSQL implementation of thunderdome.TeamService.
type TeamService struct {
	DB     *sql.DB
	Logger *otelzap.Logger
}

// TeamUserRole gets a user's role in team
func (d *TeamService) TeamUserRole(ctx context.Context, UserID string, TeamID string) (string, error) {
	var teamRole string

	err := d.DB.QueryRowContext(ctx,
		`SELECT role FROM team_get_user_role($1, $2)`,
		UserID,
		TeamID,
	).Scan(
		&teamRole,
	)
	if err != nil {
		d.Logger.Ctx(ctx).Error("team_get_user_role query error", zap.Error(err))
		return "", errors.New("error getting team users role")
	}

	return teamRole, nil
}

// TeamGet gets a team
func (d *TeamService) TeamGet(ctx context.Context, TeamID string) (*thunderdome.Team, error) {
	var team = &thunderdome.Team{}

	err := d.DB.QueryRowContext(ctx,
		`SELECT id, name, created_date, updated_date FROM team_get_by_id($1)`,
		TeamID,
	).Scan(
		&team.Id,
		&team.Name,
		&team.CreatedDate,
		&team.UpdatedDate,
	)
	if err != nil {
		d.Logger.Ctx(ctx).Error("team_get_by_id query error", zap.Error(err))
		return nil, errors.New("team not found")
	}

	return team, nil
}

// TeamListByUser gets a list of teams the user is on
func (d *TeamService) TeamListByUser(ctx context.Context, UserID string, Limit int, Offset int) []*thunderdome.Team {
	var teams = make([]*thunderdome.Team, 0)
	rows, err := d.DB.QueryContext(ctx,
		`SELECT id, name, created_date, updated_date FROM team_list_by_user($1, $2, $3);`,
		UserID,
		Limit,
		Offset,
	)

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var team thunderdome.Team

			if err := rows.Scan(
				&team.Id,
				&team.Name,
				&team.CreatedDate,
				&team.UpdatedDate,
			); err != nil {
				d.Logger.Ctx(ctx).Error("team_list_by_user query scan error", zap.Error(err))
			} else {
				teams = append(teams, &team)
			}
		}
	} else {
		d.Logger.Ctx(ctx).Error("team_list_by_user query error", zap.Error(err))
	}

	return teams
}

// TeamCreate creates a team with current user as an ADMIN
func (d *TeamService) TeamCreate(ctx context.Context, UserID string, TeamName string) (*thunderdome.Team, error) {
	t := &thunderdome.Team{}
	err := d.DB.QueryRowContext(ctx, `
		SELECT id, name, created_date, updated_date FROM team_create($1, $2);`,
		UserID,
		TeamName,
	).Scan(&t.Id, &t.Name, &t.CreatedDate, &t.UpdatedDate)

	if err != nil {
		d.Logger.Ctx(ctx).Error("team_create query error", zap.Error(err))
		return nil, err
	}

	return t, nil
}

// TeamAddUser adds a user to a team
func (d *TeamService) TeamAddUser(ctx context.Context, TeamID string, UserID string, Role string) (string, error) {
	_, err := d.DB.ExecContext(ctx,
		`SELECT team_user_add($1, $2, $3);`,
		TeamID,
		UserID,
		Role,
	)

	if err != nil {
		d.Logger.Ctx(ctx).Error("team_user_add query error", zap.Error(err))
		return "", err
	}

	return TeamID, nil
}

// TeamUserList gets a list of team users
func (d *TeamService) TeamUserList(ctx context.Context, TeamID string, Limit int, Offset int) ([]*thunderdome.TeamUser, int, error) {
	var users = make([]*thunderdome.TeamUser, 0)
	var userCount int

	err := d.DB.QueryRowContext(ctx,
		`SELECT count(user_id) FROM team_user WHERE team_id = $1;`,
		TeamID,
	).Scan(&userCount)
	if err != nil {
		return nil, 0, err
	}

	if userCount == 0 {
		return users, userCount, nil
	}

	rows, err := d.DB.QueryContext(ctx,
		`SELECT id, name, email, role, avatar FROM team_user_list($1, $2, $3);`,
		TeamID,
		Limit,
		Offset,
	)

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var usr thunderdome.TeamUser

			if err = rows.Scan(
				&usr.Id,
				&usr.Name,
				&usr.Email,
				&usr.Role,
				&usr.Avatar,
			); err != nil {
				d.Logger.Ctx(ctx).Error("team_user_list query scan error", zap.Error(err))
			} else {
				usr.GravatarHash = createGravatarHash(usr.Email)
				users = append(users, &usr)
			}
		}
	} else {
		d.Logger.Ctx(ctx).Error("team_user_list query error", zap.Error(err))
		return nil, 0, err
	}

	return users, userCount, nil
}

// TeamRemoveUser removes a user from a team
func (d *TeamService) TeamRemoveUser(ctx context.Context, TeamID string, UserID string) error {
	_, err := d.DB.ExecContext(ctx,
		`CALL team_user_remove($1, $2);`,
		TeamID,
		UserID,
	)

	if err != nil {
		d.Logger.Ctx(ctx).Error("team_user_remove query error", zap.Error(err))
		return err
	}

	return nil
}

// TeamBattleList gets a list of team battles
func (d *TeamService) TeamBattleList(ctx context.Context, TeamID string, Limit int, Offset int) []*thunderdome.Battle {
	var battles = make([]*thunderdome.Battle, 0)
	rows, err := d.DB.QueryContext(ctx,
		`SELECT id, name FROM team_battle_list($1, $2, $3);`,
		TeamID,
		Limit,
		Offset,
	)

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var tb thunderdome.Battle

			if err := rows.Scan(
				&tb.Id,
				&tb.Name,
			); err != nil {
				d.Logger.Ctx(ctx).Error("team_battle_list query scan error", zap.Error(err))
			} else {
				battles = append(battles, &tb)
			}
		}
	} else {
		d.Logger.Ctx(ctx).Error("team_battle_list query error", zap.Error(err))
	}

	return battles
}

// TeamAddBattle adds a battle to a team
func (d *TeamService) TeamAddBattle(ctx context.Context, TeamID string, BattleID string) error {
	_, err := d.DB.ExecContext(ctx,
		`SELECT team_battle_add($1, $2);`,
		TeamID,
		BattleID,
	)

	if err != nil {
		d.Logger.Ctx(ctx).Error("team_battle_add query error", zap.Error(err))
		return err
	}

	return nil
}

// TeamRemoveBattle removes a battle from a team
func (d *TeamService) TeamRemoveBattle(ctx context.Context, TeamID string, BattleID string) error {
	_, err := d.DB.ExecContext(ctx,
		`SELECT team_battle_remove($1, $2);`,
		TeamID,
		BattleID,
	)

	if err != nil {
		d.Logger.Ctx(ctx).Error("team_battle_remove query error", zap.Error(err))
		return err
	}

	return nil
}

// TeamDelete deletes a team
func (d *TeamService) TeamDelete(ctx context.Context, TeamID string) error {
	_, err := d.DB.ExecContext(ctx,
		`CALL team_delete($1);`,
		TeamID,
	)

	if err != nil {
		d.Logger.Ctx(ctx).Error("team_delete query error", zap.Error(err))
		return err
	}

	return nil
}

// TeamRetroList gets a list of team retros
func (d *TeamService) TeamRetroList(ctx context.Context, TeamID string, Limit int, Offset int) []*thunderdome.Retro {
	var retros = make([]*thunderdome.Retro, 0)
	rows, err := d.DB.QueryContext(ctx,
		`SELECT id, name, format, phase FROM team_retro_list($1, $2, $3);`,
		TeamID,
		Limit,
		Offset,
	)

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var tb thunderdome.Retro

			if err := rows.Scan(
				&tb.Id,
				&tb.Name,
				&tb.Format,
				&tb.Phase,
			); err != nil {
				d.Logger.Ctx(ctx).Error("team_retro_list query scan error", zap.Error(err))
			} else {
				retros = append(retros, &tb)
			}
		}
	} else {
		d.Logger.Ctx(ctx).Error("team_retro_list query error", zap.Error(err))
	}

	return retros
}

// TeamAddRetro adds a retro to a team
func (d *TeamService) TeamAddRetro(ctx context.Context, TeamID string, RetroID string) error {
	_, err := d.DB.ExecContext(ctx,
		`SELECT team_retro_add($1, $2);`,
		TeamID,
		RetroID,
	)

	if err != nil {
		d.Logger.Ctx(ctx).Error("team_retro_add query error", zap.Error(err))
		return err
	}

	return nil
}

// TeamRemoveRetro removes a retro from a team
func (d *TeamService) TeamRemoveRetro(ctx context.Context, TeamID string, RetroID string) error {
	_, err := d.DB.ExecContext(ctx,
		`SELECT team_retro_remove($1, $2);`,
		TeamID,
		RetroID,
	)

	if err != nil {
		d.Logger.Ctx(ctx).Error("team_retro_remove query error", zap.Error(err))
		return err
	}

	return nil
}

// TeamStoryboardList gets a list of team storyboards
func (d *TeamService) TeamStoryboardList(ctx context.Context, TeamID string, Limit int, Offset int) []*thunderdome.Storyboard {
	var storyboards = make([]*thunderdome.Storyboard, 0)
	rows, err := d.DB.QueryContext(ctx,
		`SELECT id, name FROM team_storyboard_list($1, $2, $3);`,
		TeamID,
		Limit,
		Offset,
	)

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var tb thunderdome.Storyboard

			if err := rows.Scan(
				&tb.Id,
				&tb.Name,
			); err != nil {
				d.Logger.Ctx(ctx).Error("team_storyboard_list query scan error", zap.Error(err))
			} else {
				storyboards = append(storyboards, &tb)
			}
		}
	} else {
		d.Logger.Ctx(ctx).Error("team_storyboard_list query error", zap.Error(err))
	}

	return storyboards
}

// TeamAddStoryboard adds a storyboard to a team
func (d *TeamService) TeamAddStoryboard(ctx context.Context, TeamID string, StoryboardID string) error {
	_, err := d.DB.ExecContext(ctx,
		`SELECT team_storyboard_add($1, $2);`,
		TeamID,
		StoryboardID,
	)

	if err != nil {
		d.Logger.Ctx(ctx).Error("team_storyboard_add query error", zap.Error(err))
		return err
	}

	return nil
}

// TeamRemoveStoryboard removes a storyboard from a team
func (d *TeamService) TeamRemoveStoryboard(ctx context.Context, TeamID string, StoryboardID string) error {
	_, err := d.DB.ExecContext(ctx,
		`SELECT team_storyboard_remove($1, $2);`,
		TeamID,
		StoryboardID,
	)

	if err != nil {
		d.Logger.Ctx(ctx).Error("team_storyboard_remove query error", zap.Error(err))
		return err
	}

	return nil
}

// TeamList gets a list of teams
func (d *TeamService) TeamList(ctx context.Context, Limit int, Offset int) ([]*thunderdome.Team, int) {
	var teams = make([]*thunderdome.Team, 0)
	var count = 0

	err := d.DB.QueryRowContext(ctx, `SELECT count FROM team_list_count();`).Scan(&count)
	if err != nil {
		d.Logger.Ctx(ctx).Error("Unable to get application stats", zap.Error(err))
		return teams, count
	}

	rows, err := d.DB.QueryContext(ctx,
		`SELECT id, name, created_date, updated_date FROM team_list($1, $2);`,
		Limit,
		Offset,
	)

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var team thunderdome.Team

			if err := rows.Scan(
				&team.Id,
				&team.Name,
				&team.CreatedDate,
				&team.UpdatedDate,
			); err != nil {
				d.Logger.Ctx(ctx).Error("team_list scan error", zap.Error(err))
			} else {
				teams = append(teams, &team)
			}
		}
	} else {
		d.Logger.Ctx(ctx).Error("team_list query error", zap.Error(err))
	}

	return teams, count
}
