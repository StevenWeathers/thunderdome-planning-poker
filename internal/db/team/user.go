package team

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/StevenWeathers/thunderdome-planning-poker/internal/db"
	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
	"go.uber.org/zap"
)

// TeamUserRole gets a user's role in team
func (d *Service) TeamUserRole(ctx context.Context, userID string, teamID string) (string, error) {
	var teamRole string

	err := d.DB.QueryRowContext(ctx,
		`SELECT tu.role
        FROM thunderdome.team_user tu
        WHERE tu.team_id = $2 AND tu.user_id = $1;`,
		userID,
		teamID,
	).Scan(
		&teamRole,
	)
	if err != nil {
		return "", fmt.Errorf("error getting team users role: %v", err)
	}

	return teamRole, nil
}

// TeamUserRoles gets a user's set of roles in relation to the team if any, and if application the department and organization
func (d *Service) TeamUserRoles(ctx context.Context, userID string, teamID string) (*thunderdome.UserTeamRoleInfo, error) {
	tr := thunderdome.UserTeamRoleInfo{}

	err := d.DB.QueryRowContext(ctx,
		`WITH team_info AS (
    SELECT
        t.id AS team_id,
        t.department_id,
        COALESCE(t.organization_id, od.organization_id) AS organization_id
    FROM
        thunderdome.team t
    LEFT JOIN
        thunderdome.organization_department od ON t.department_id = od.id
    WHERE
        t.id = $2
),
user_roles AS (
    SELECT
        ti.*,
        $1 AS user_id,
        tu.role AS team_role,
        du.role AS department_role,
        ou.role AS organization_role
    FROM
        team_info ti
    LEFT JOIN
        thunderdome.team_user tu ON ti.team_id = tu.team_id AND tu.user_id = $1
    LEFT JOIN
        thunderdome.department_user du ON ti.department_id = du.department_id AND du.user_id = $1
    LEFT JOIN
        thunderdome.organization_user ou ON ti.organization_id = ou.organization_id AND ou.user_id = $1
)
SELECT
    user_id,
    team_id,
    team_role,
    department_id,
    department_role,
    organization_id,
    organization_role,
    CASE
        WHEN team_role IS NOT NULL THEN 'TEAM'
        WHEN department_role IS NOT NULL THEN 'DEPARTMENT'
        WHEN organization_role IS NOT NULL THEN 'ORGANIZATION'
        ELSE 'NONE'
    END AS association_level
FROM
    user_roles;`,
		userID,
		teamID,
	).Scan(
		&tr.UserID,
		&tr.TeamID,
		&tr.TeamRole,
		&tr.DepartmentID,
		&tr.DepartmentRole,
		&tr.OrganizationID,
		&tr.OrganizationRole,
		&tr.AssociationLevel,
	)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("TEAM_NOT_FOUND")
	} else if err != nil {
		return nil, fmt.Errorf("error getting team users roles: %v", err)
	}

	return &tr, nil
}

// TeamUserList gets a list of team users
func (d *Service) TeamUserList(ctx context.Context, teamID string, limit int, offset int) ([]*thunderdome.TeamUser, int, error) {
	var users = make([]*thunderdome.TeamUser, 0)
	var userCount int

	err := d.DB.QueryRowContext(ctx,
		`SELECT count(user_id) FROM thunderdome.team_user WHERE team_id = $1;`,
		teamID,
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
		teamID,
		limit,
		offset,
	)

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var usr thunderdome.TeamUser

			if err = rows.Scan(
				&usr.ID,
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
func (d *Service) TeamAddUser(ctx context.Context, teamID string, userID string, role string) (string, error) {
	_, err := d.DB.ExecContext(ctx,
		`INSERT INTO thunderdome.team_user (team_id, user_id, role) VALUES ($1, $2, $3);`,
		teamID,
		userID,
		role,
	)

	if err != nil {
		return "", fmt.Errorf("team add user query error: %v", err)
	}

	return teamID, nil
}

// TeamUpdateUser updates a team user
func (d *Service) TeamUpdateUser(ctx context.Context, teamID string, userID string, role string) (string, error) {
	_, err := d.DB.ExecContext(ctx,
		`UPDATE thunderdome.team_user SET role = $3 WHERE team_id = $1 AND user_id = $2;`,
		teamID,
		userID,
		role,
	)

	if err != nil {
		return "", fmt.Errorf("team update user query error: %v", err)
	}

	return teamID, nil
}

// TeamRemoveUser removes a user from a team
func (d *Service) TeamRemoveUser(ctx context.Context, teamID string, userID string) error {
	_, err := d.DB.ExecContext(ctx,
		`DELETE FROM thunderdome.team_user WHERE team_id = $1 AND user_id = $2;`,
		teamID,
		userID,
	)

	if err != nil {
		return fmt.Errorf("team remove user query error: %v", err)
	}

	return nil
}
