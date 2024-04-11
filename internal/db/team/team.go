package team

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/StevenWeathers/thunderdome-planning-poker/internal/db"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"

	"go.uber.org/zap"
)

// Service represents a PostgreSQL implementation of thunderdome.TeamDataSvc.
type Service struct {
	DB     *sql.DB
	Logger *otelzap.Logger
}

// TeamUserRole gets a user's role in team
func (d *Service) TeamUserRole(ctx context.Context, UserID string, TeamID string) (string, error) {
	var teamRole string

	err := d.DB.QueryRowContext(ctx,
		`SELECT tu.role
        FROM thunderdome.team_user tu
        WHERE tu.team_id = $2 AND tu.user_id = $1;`,
		UserID,
		TeamID,
	).Scan(
		&teamRole,
	)
	if err != nil {
		return "", fmt.Errorf("error getting team users role: %v", err)
	}

	return teamRole, nil
}

// TeamGet gets a team
func (d *Service) TeamGet(ctx context.Context, TeamID string) (*thunderdome.Team, error) {
	var team = &thunderdome.Team{}

	err := d.DB.QueryRowContext(ctx,
		`SELECT o.id, o.name, o.created_date, o.updated_date
        FROM thunderdome.team o
        WHERE o.id = $1;`,
		TeamID,
	).Scan(
		&team.Id,
		&team.Name,
		&team.CreatedDate,
		&team.UpdatedDate,
	)
	if err != nil {
		return nil, fmt.Errorf("get team query error: %v", err)
	}

	return team, nil
}

// TeamListByUser gets a list of teams the user is on
func (d *Service) TeamListByUser(ctx context.Context, UserID string, Limit int, Offset int) []*thunderdome.UserTeam {
	var teams = make([]*thunderdome.UserTeam, 0)
	rows, err := d.DB.QueryContext(ctx,
		`SELECT t.id, t.name, t.created_date, t.updated_date, tu.role
        FROM thunderdome.team_user tu
        LEFT JOIN thunderdome.team t ON tu.team_id = t.id
        WHERE tu.user_id = $1
        ORDER BY t.created_date
		LIMIT $2
		OFFSET $3;`,
		UserID,
		Limit,
		Offset,
	)

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var team thunderdome.UserTeam

			if err := rows.Scan(
				&team.Id,
				&team.Name,
				&team.CreatedDate,
				&team.UpdatedDate,
				&team.Role,
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
func (d *Service) TeamCreate(ctx context.Context, UserID string, TeamName string) (*thunderdome.Team, error) {
	t := &thunderdome.Team{}
	err := d.DB.QueryRowContext(ctx, `
		SELECT id, name, created_date, updated_date FROM thunderdome.team_create($1, $2);`,
		UserID,
		TeamName,
	).Scan(&t.Id, &t.Name, &t.CreatedDate, &t.UpdatedDate)

	if err != nil {
		return nil, fmt.Errorf("create team query error: %v", err)
	}

	return t, nil
}

// TeamUpdate updates a team
func (d *Service) TeamUpdate(ctx context.Context, TeamId string, TeamName string) (*thunderdome.Team, error) {
	t := &thunderdome.Team{}
	err := d.DB.QueryRowContext(ctx, `
		UPDATE thunderdome.team
		SET name = $1, updated_date = NOW()
		WHERE id = $2
		RETURNING id, name, created_date, updated_date;`,
		TeamName,
		TeamId,
	).Scan(&t.Id, &t.Name, &t.CreatedDate, &t.UpdatedDate)

	if err != nil {
		return nil, fmt.Errorf("team update query error: %v", err)
	}

	return t, nil
}

// TeamAddUser adds a user to a team
func (d *Service) TeamAddUser(ctx context.Context, TeamID string, UserID string, Role string) (string, error) {
	_, err := d.DB.ExecContext(ctx,
		`INSERT INTO thunderdome.team_user (team_id, user_id, role) VALUES ($1, $2, $3);`,
		TeamID,
		UserID,
		Role,
	)

	if err != nil {
		return "", fmt.Errorf("team add user query error: %v", err)
	}

	return TeamID, nil
}

// TeamUpdateUser updates a team user
func (d *Service) TeamUpdateUser(ctx context.Context, TeamID string, UserID string, Role string) (string, error) {
	_, err := d.DB.ExecContext(ctx,
		`UPDATE thunderdome.team_user SET role = $3 WHERE team_id = $1 AND user_id = $2;`,
		TeamID,
		UserID,
		Role,
	)

	if err != nil {
		return "", fmt.Errorf("team update user query error: %v", err)
	}

	return TeamID, nil
}

// TeamInviteUser invites a user to a team
func (d *Service) TeamInviteUser(ctx context.Context, TeamID string, Email string, Role string) (string, error) {
	var inviteId string
	err := d.DB.QueryRowContext(ctx,
		`INSERT INTO thunderdome.team_user_invite (team_id, email, role) VALUES ($1, $2, $3) RETURNING invite_id;`,
		TeamID,
		Email,
		Role,
	).Scan(&inviteId)

	if err != nil {
		return "", fmt.Errorf("team invite user query error: %v", err)
	}

	return inviteId, nil
}

// TeamUserGetInviteByID gets a team user invite
func (d *Service) TeamUserGetInviteByID(ctx context.Context, InviteID string) (thunderdome.TeamUserInvite, error) {
	tui := thunderdome.TeamUserInvite{}
	err := d.DB.QueryRowContext(ctx,
		`SELECT invite_id, team_id, email, role, created_date, expire_date
 				FROM thunderdome.team_user_invite WHERE invite_id = $1;`,
		InviteID,
	).Scan(&tui.InviteId, &tui.TeamId, &tui.Email, &tui.Role, &tui.CreatedDate, &tui.ExpireDate)

	if err != nil {
		return tui, fmt.Errorf("team get user invite query error: %v", err)
	}

	return tui, nil
}

// TeamDeleteUserInvite deletes a user team invite
func (d *Service) TeamDeleteUserInvite(ctx context.Context, InviteID string) error {
	_, err := d.DB.ExecContext(ctx,
		`DELETE FROM thunderdome.team_user_invite where invite_id = $1;`,
		InviteID,
	)

	if err != nil {
		return fmt.Errorf("team delete user invite query error: %v", err)
	}

	return nil
}

// TeamGetUserInvites gets teams user invites
func (d *Service) TeamGetUserInvites(ctx context.Context, teamId string) ([]thunderdome.TeamUserInvite, error) {
	var invites = make([]thunderdome.TeamUserInvite, 0)
	rows, err := d.DB.QueryContext(ctx,
		`SELECT invite_id, team_id, email, role, created_date, expire_date
 				FROM thunderdome.team_user_invite WHERE team_id = $1;`,
		teamId,
	)

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var invite thunderdome.TeamUserInvite

			if err := rows.Scan(
				&invite.InviteId,
				&invite.TeamId,
				&invite.Email,
				&invite.Role,
				&invite.CreatedDate,
				&invite.ExpireDate,
			); err != nil {
				d.Logger.Ctx(ctx).Error("TeamGetUserInvites query scan error", zap.Error(err))
			} else {
				invites = append(invites, invite)
			}
		}
	} else {
		if !errors.Is(err, sql.ErrNoRows) {
			d.Logger.Ctx(ctx).Error("TeamGetUserInvites query error", zap.Error(err))
		}
	}

	return invites, nil
}

// TeamUserList gets a list of team users
func (d *Service) TeamUserList(ctx context.Context, TeamID string, Limit int, Offset int) ([]*thunderdome.TeamUser, int, error) {
	var users = make([]*thunderdome.TeamUser, 0)
	var userCount int

	err := d.DB.QueryRowContext(ctx,
		`SELECT count(user_id) FROM thunderdome.team_user WHERE team_id = $1;`,
		TeamID,
	).Scan(&userCount)
	if err != nil {
		return nil, 0, fmt.Errorf("team user list user count query error: %v", err)
	}

	if userCount == 0 {
		return users, userCount, nil
	}

	rows, err := d.DB.QueryContext(ctx,
		`SELECT u.id, u.name, COALESCE(u.email, ''), tu.role, u.avatar, COALESCE(u.picture, '')
        FROM thunderdome.team_user tu
        LEFT JOIN thunderdome.users u ON tu.user_id = u.id
        WHERE tu.team_id = $1
        ORDER BY tu.created_date
		LIMIT $2
		OFFSET $3;`,
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
				&usr.PictureURL,
			); err != nil {
				d.Logger.Ctx(ctx).Error("team_user_list query scan error", zap.Error(err))
			} else {
				usr.GravatarHash = db.CreateGravatarHash(usr.Email)
				users = append(users, &usr)
			}
		}
	} else {
		return nil, 0, fmt.Errorf("team user list query error: %v", err)
	}

	return users, userCount, nil
}

// TeamRemoveUser removes a user from a team
func (d *Service) TeamRemoveUser(ctx context.Context, TeamID string, UserID string) error {
	_, err := d.DB.ExecContext(ctx,
		`DELETE FROM thunderdome.team_user WHERE team_id = $1 AND user_id = $2;`,
		TeamID,
		UserID,
	)

	if err != nil {
		return fmt.Errorf("team remove user query error: %v", err)
	}

	return nil
}

// TeamPokerList gets a list of team poker games
func (d *Service) TeamPokerList(ctx context.Context, TeamID string, Limit int, Offset int) []*thunderdome.Poker {
	var pokers = make([]*thunderdome.Poker, 0)
	rows, err := d.DB.QueryContext(ctx,
		`SELECT p.id, p.name
        FROM thunderdome.poker p
        WHERE p.team_id = $1
        ORDER BY p.created_date DESC
		LIMIT $2
		OFFSET $3;`,
		TeamID,
		Limit,
		Offset,
	)

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var tb thunderdome.Poker

			if err := rows.Scan(
				&tb.Id,
				&tb.Name,
			); err != nil {
				d.Logger.Ctx(ctx).Error("team_poker list query scan error", zap.Error(err))
			} else {
				pokers = append(pokers, &tb)
			}
		}
	} else {
		d.Logger.Ctx(ctx).Error("team_poker list query error", zap.Error(err))
	}

	return pokers
}

// TeamAddPoker adds a poker game to a team
func (d *Service) TeamAddPoker(ctx context.Context, TeamID string, PokerID string) error {
	_, err := d.DB.ExecContext(ctx,
		`UPDATE thunderdome.poker SET team_id = $1 WHERE id = $2;`,
		TeamID,
		PokerID,
	)

	if err != nil {
		return fmt.Errorf("team add poker query error: %v", err)
	}

	return nil
}

// TeamRemovePoker removes a poker game from a team
func (d *Service) TeamRemovePoker(ctx context.Context, TeamID string, PokerID string) error {
	_, err := d.DB.ExecContext(ctx,
		`UPDATE thunderdome.poker SET team_id = null WHERE id = $2 AND team_id = $1;`,
		TeamID,
		PokerID,
	)

	if err != nil {
		return fmt.Errorf("team remove poker query error: %v", err)
	}

	return nil
}

// TeamDelete deletes a team
func (d *Service) TeamDelete(ctx context.Context, TeamID string) error {
	_, err := d.DB.ExecContext(ctx,
		`DELETE FROM thunderdome.team WHERE id = $1;`,
		TeamID,
	)

	if err != nil {
		return fmt.Errorf("team delete query error: %v", err)
	}

	return nil
}

// TeamRetroList gets a list of team retros
func (d *Service) TeamRetroList(ctx context.Context, TeamID string, Limit int, Offset int) []*thunderdome.Retro {
	var retros = make([]*thunderdome.Retro, 0)
	rows, err := d.DB.QueryContext(ctx,
		`SELECT r.id, r.name, r.format, r.phase
        FROM thunderdome.retro r
        WHERE r.team_id = $1
        ORDER BY r.created_date DESC
		LIMIT $2
		OFFSET $3;`,
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
func (d *Service) TeamAddRetro(ctx context.Context, TeamID string, RetroID string) error {
	_, err := d.DB.ExecContext(ctx,
		`UPDATE thunderdome.retro SET team_id = $1 WHERE id = $2;`,
		TeamID,
		RetroID,
	)

	if err != nil {
		return fmt.Errorf("team add retro query error: %v", err)
	}

	return nil
}

// TeamRemoveRetro removes a retro from a team
func (d *Service) TeamRemoveRetro(ctx context.Context, TeamID string, RetroID string) error {
	_, err := d.DB.ExecContext(ctx,
		`UPDATE thunderdome.retro SET team_id = $1 WHERE id = $2;`,
		TeamID,
		RetroID,
	)

	if err != nil {
		return fmt.Errorf("team remove retro query error: %v", err)
	}

	return nil
}

// TeamStoryboardList gets a list of team storyboards
func (d *Service) TeamStoryboardList(ctx context.Context, TeamID string, Limit int, Offset int) []*thunderdome.Storyboard {
	var storyboards = make([]*thunderdome.Storyboard, 0)
	rows, err := d.DB.QueryContext(ctx,
		`SELECT s.id, s.name
        FROM thunderdome.storyboard s
        WHERE s.team_id = $1
        ORDER BY s.created_date DESC
		LIMIT $2
		OFFSET $3;`,
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
func (d *Service) TeamAddStoryboard(ctx context.Context, TeamID string, StoryboardID string) error {
	_, err := d.DB.ExecContext(ctx,
		`UPDATE thunderdome.storyboard SET team_id = $1 WHERE id = $2;`,
		TeamID,
		StoryboardID,
	)

	if err != nil {
		return fmt.Errorf("team add storyboard query error: %v", err)
	}

	return nil
}

// TeamRemoveStoryboard removes a storyboard from a team
func (d *Service) TeamRemoveStoryboard(ctx context.Context, TeamID string, StoryboardID string) error {
	_, err := d.DB.ExecContext(ctx,
		`UPDATE thunderdome.storyboard SET team_id = $1 WHERE id = $2;`,
		TeamID,
		StoryboardID,
	)

	if err != nil {
		return fmt.Errorf("team remove storyboard query error: %v", err)
	}

	return nil
}

// TeamList gets a list of teams
func (d *Service) TeamList(ctx context.Context, Limit int, Offset int) ([]*thunderdome.Team, int) {
	var teams = make([]*thunderdome.Team, 0)
	var count = 0

	err := d.DB.QueryRowContext(ctx, `SELECT count(t.id)
    FROM thunderdome.team t
    WHERE t.department_id IS NULL AND t.organization_id IS NULL;`).Scan(&count)
	if err != nil {
		d.Logger.Ctx(ctx).Error("Unable to get application stats", zap.Error(err))
		return teams, count
	}

	rows, err := d.DB.QueryContext(ctx,
		`SELECT t.id, t.name, t.created_date, t.updated_date
        FROM thunderdome.team t
        WHERE t.department_id IS NULL AND t.organization_id IS NULL
        ORDER BY t.created_date
		LIMIT $1
		OFFSET $2;`,
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
