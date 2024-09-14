package team

import (
	"context"
	"fmt"

	"github.com/StevenWeathers/thunderdome-planning-poker/internal/db"
	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
	"go.uber.org/zap"
)

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
