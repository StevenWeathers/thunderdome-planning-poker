package team

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"

	"go.uber.org/zap"
)

// Service represents a PostgreSQL implementation of thunderdome.TeamDataSvc.
type Service struct {
	DB     *sql.DB
	Logger *otelzap.Logger
}

// TeamGet gets a team
func (d *Service) TeamGet(ctx context.Context, teamID string) (*thunderdome.Team, error) {
	var team = &thunderdome.Team{}

	err := d.DB.QueryRowContext(ctx,
		`SELECT o.id, o.name, COALESCE(o.organization_id::TEXT, od.organization_id::TEXT, ''),
 COALESCE(o.department_id::TEXT, ''), o.created_date, o.updated_date,
 CASE WHEN s.id IS NOT NULL AND s.expires > NOW() AND s.active = true THEN true ELSE false END AS is_subscribed
        FROM thunderdome.team o
        LEFT JOIN thunderdome.subscription s ON o.id = s.team_id
        LEFT JOIN thunderdome.organization_department od ON o.department_id = od.id
        WHERE o.id = $1;`,
		teamID,
	).Scan(
		&team.ID,
		&team.Name,
		&team.OrganizationID,
		&team.DepartmentID,
		&team.CreatedDate,
		&team.UpdatedDate,
		&team.Subscribed,
	)
	if err != nil {
		return nil, fmt.Errorf("get team query error: %v", err)
	}

	return team, nil
}

// TeamListByUser gets a list of teams the user is on
func (d *Service) TeamListByUser(ctx context.Context, userID string, limit int, offset int) []*thunderdome.UserTeam {
	var teams = make([]*thunderdome.UserTeam, 0)
	rows, err := d.DB.QueryContext(ctx,
		`SELECT
    t.id,
    t.name,
    COALESCE(t.organization_id::TEXT,
             (SELECT organization_id::TEXT
              FROM thunderdome.organization_department
              WHERE id = t.department_id),
             '') AS organization_id,
    COALESCE(t.department_id::TEXT, ''),
    t.created_date,
    t.updated_date,
    tu.role,
    sub_check.is_subscribed
FROM thunderdome.team_user tu
LEFT JOIN thunderdome.team t ON tu.team_id = t.id
LEFT JOIN LATERAL (
    SELECT
    COALESCE(
        (SELECT TRUE
        FROM (
            -- Check for direct team subscription
            SELECT TRUE
            FROM thunderdome.subscription
            WHERE team_id = t.id
                AND active = TRUE
                AND expires > CURRENT_TIMESTAMP
            UNION ALL
            -- Check for organization subscription (either direct or via department)
            SELECT TRUE
            FROM thunderdome.subscription s
            JOIN thunderdome.team t2 ON
                CASE
                WHEN t2.department_id IS NOT NULL THEN
                    s.organization_id = (SELECT organization_id FROM thunderdome.organization_department WHERE id = t2.department_id)
                ELSE
                    s.organization_id = t2.organization_id
                END
            WHERE t2.id = t.id
                AND s.active = TRUE
                AND s.expires > CURRENT_TIMESTAMP
        ) AS subscriptions
        LIMIT 1),
        FALSE
    ) AS is_subscribed
) AS sub_check ON TRUE
WHERE tu.user_id = $1
ORDER BY t.created_date
LIMIT $2
OFFSET $3;`,
		userID,
		limit,
		offset,
	)

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var team thunderdome.UserTeam

			if err := rows.Scan(
				&team.ID,
				&team.Name,
				&team.OrganizationID,
				&team.DepartmentID,
				&team.CreatedDate,
				&team.UpdatedDate,
				&team.Role,
				&team.Subscribed,
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

// TeamListByUserNonOrg gets a list of teams the user is on that are not part of an organization
func (d *Service) TeamListByUserNonOrg(ctx context.Context, userID string, limit int, offset int) []*thunderdome.UserTeam {
	var teams = make([]*thunderdome.UserTeam, 0)
	rows, err := d.DB.QueryContext(ctx,
		`SELECT
    t.id,
    t.name,
    t.created_date,
    t.updated_date,
    tu.role,
    sub_check.is_subscribed
FROM thunderdome.team_user tu
LEFT JOIN thunderdome.team t ON tu.team_id = t.id
LEFT JOIN LATERAL (
    SELECT
    COALESCE(
        (SELECT TRUE
        FROM (
            -- Check for direct team subscription
            SELECT TRUE
            FROM thunderdome.subscription
            WHERE team_id = t.id
                AND active = TRUE
                AND expires > CURRENT_TIMESTAMP
        ) AS subscriptions
        LIMIT 1),
        FALSE
    ) AS is_subscribed
) AS sub_check ON TRUE
WHERE tu.user_id = $1 AND t.department_id IS NULL AND t.organization_id IS NULL
ORDER BY t.created_date
LIMIT $2
OFFSET $3;`,
		userID,
		limit,
		offset,
	)

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var team thunderdome.UserTeam

			if err := rows.Scan(
				&team.ID,
				&team.Name,
				&team.CreatedDate,
				&team.UpdatedDate,
				&team.Role,
				&team.Subscribed,
			); err != nil {
				d.Logger.Ctx(ctx).Error("team_list_by_user_non_org query scan error", zap.Error(err))
			} else {
				teams = append(teams, &team)
			}
		}
	} else {
		d.Logger.Ctx(ctx).Error("team_list_by_user_non_org query error", zap.Error(err))
	}

	return teams
}

// TeamCreate creates a team with current user as an ADMIN
func (d *Service) TeamCreate(ctx context.Context, userID string, teamName string) (*thunderdome.Team, error) {
	t := &thunderdome.Team{}
	err := d.DB.QueryRowContext(ctx, `
		SELECT id, name, created_date, updated_date FROM thunderdome.team_create($1, $2);`,
		userID,
		teamName,
	).Scan(&t.ID, &t.Name, &t.CreatedDate, &t.UpdatedDate)

	if err != nil {
		return nil, fmt.Errorf("create team query error: %v", err)
	}

	return t, nil
}

// TeamUpdate updates a team
func (d *Service) TeamUpdate(ctx context.Context, teamID string, teamName string) (*thunderdome.Team, error) {
	t := &thunderdome.Team{}
	err := d.DB.QueryRowContext(ctx, `
		UPDATE thunderdome.team
		SET name = $1, updated_date = NOW()
		WHERE id = $2
		RETURNING id, name, created_date, updated_date;`,
		teamName,
		teamID,
	).Scan(&t.ID, &t.Name, &t.CreatedDate, &t.UpdatedDate)

	if err != nil {
		return nil, fmt.Errorf("team update query error: %v", err)
	}

	return t, nil
}

// TeamDelete deletes a team
func (d *Service) TeamDelete(ctx context.Context, teamID string) error {
	_, err := d.DB.ExecContext(ctx,
		`DELETE FROM thunderdome.team WHERE id = $1;`,
		teamID,
	)

	if err != nil {
		return fmt.Errorf("team delete query error: %v", err)
	}

	return nil
}

// TeamList gets a list of teams
func (d *Service) TeamList(ctx context.Context, limit int, offset int) ([]*thunderdome.Team, int) {
	var teams = make([]*thunderdome.Team, 0)
	var count = 0

	err := d.DB.QueryRowContext(ctx, `SELECT count(t.id)
    FROM thunderdome.team t
    WHERE t.department_id IS NULL AND t.organization_id IS NULL;`).Scan(&count)
	if err != nil {
		d.Logger.Ctx(ctx).Error("Unable to get TeamList", zap.Error(err))
		return teams, count
	}

	rows, err := d.DB.QueryContext(ctx,
		`SELECT t.id, t.name, t.created_date, t.updated_date
        FROM thunderdome.team t
        WHERE t.department_id IS NULL AND t.organization_id IS NULL
        ORDER BY t.created_date
		LIMIT $1
		OFFSET $2;`,
		limit,
		offset,
	)

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var team thunderdome.Team

			if err := rows.Scan(
				&team.ID,
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

// TeamIsSubscribed checks if a team is subscribed
func (d *Service) TeamIsSubscribed(ctx context.Context, teamID string) (bool, error) {
	var subscribed bool

	err := d.DB.QueryRowContext(ctx,
		`SELECT
  COALESCE(
      (SELECT TRUE
       FROM (
           -- Check for direct team subscription
           SELECT TRUE
           FROM thunderdome.subscription
           WHERE team_id = $1
             AND active = TRUE
             AND expires > CURRENT_TIMESTAMP
           UNION ALL
           -- Check for organization subscription (either direct or via department)
           SELECT TRUE
           FROM thunderdome.subscription s
           JOIN thunderdome.team t ON
             CASE
               WHEN t.department_id IS NOT NULL THEN
                 s.organization_id = (SELECT organization_id FROM thunderdome.organization_department WHERE id = t.department_id)
               ELSE
                 s.organization_id = t.organization_id
             END
           WHERE t.id = $1
             AND s.active = TRUE
             AND s.expires > CURRENT_TIMESTAMP
       ) AS subscriptions
       LIMIT 1),
      FALSE
  ) AS is_subscribed;`,
		teamID,
	).Scan(
		&subscribed,
	)
	if err != nil {
		return false, fmt.Errorf("error getting team subscription: %v", err)
	}

	return subscribed, nil
}
