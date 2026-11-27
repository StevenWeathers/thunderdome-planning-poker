package retro

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
	"github.com/jackc/pgx/v5/pgconn"
)

// retroTagScopeCondition matches preset tags belonging to exactly the scope given by $1 (organization) and $2 (team),
// both empty meaning the global defaults
const retroTagScopeCondition = `organization_id IS NOT DISTINCT FROM NULLIF($1, '')::uuid
	AND team_id IS NOT DISTINCT FROM NULLIF($2, '')::uuid`

// GetRetroTags gets the preset tags for a scope
func (d *Service) GetRetroTags(ctx context.Context, organizationID string, teamID string) ([]*thunderdome.RetroTag, error) {
	tags := make([]*thunderdome.RetroTag, 0)

	rows, err := d.DB.QueryContext(ctx,
		`SELECT id, name, is_public, COALESCE(created_by::text, ''), organization_id, team_id, created_at, updated_at
		FROM thunderdome.retro_tag
		WHERE `+retroTagScopeCondition+`
		ORDER BY lower(name);`,
		organizationID, teamID,
	)
	if err != nil {
		return nil, fmt.Errorf("get retro tags query error: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		t := &thunderdome.RetroTag{}
		if err := rows.Scan(
			&t.ID, &t.Name, &t.IsPublic, &t.CreatedBy, &t.OrganizationID, &t.TeamID, &t.CreatedAt, &t.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("get retro tags scan error: %v", err)
		}
		tags = append(tags, t)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("get retro tags rows error: %v", err)
	}

	return tags, nil
}

// CreateRetroTag creates a preset tag in the scope set on the tag
func (d *Service) CreateRetroTag(ctx context.Context, tag *thunderdome.RetroTag) error {
	tag.IsPublic = tag.OrganizationID == nil && tag.TeamID == nil

	err := d.DB.QueryRowContext(ctx,
		`INSERT INTO thunderdome.retro_tag (name, is_public, created_by, organization_id, team_id)
		VALUES ($1, $2, NULLIF($3, '')::uuid, $4, $5)
		RETURNING id, created_at, updated_at;`,
		tag.Name, tag.IsPublic, tag.CreatedBy, tag.OrganizationID, tag.TeamID,
	).Scan(&tag.ID, &tag.CreatedAt, &tag.UpdatedAt)
	if err != nil {
		return retroTagWriteError("create retro tag", err)
	}

	return nil
}

// UpdateRetroTag renames a preset tag within the scope set on the tag
func (d *Service) UpdateRetroTag(ctx context.Context, tag *thunderdome.RetroTag) error {
	tag.IsPublic = tag.OrganizationID == nil && tag.TeamID == nil

	err := d.DB.QueryRowContext(ctx,
		`UPDATE thunderdome.retro_tag SET name = $4, updated_at = CURRENT_TIMESTAMP
		WHERE `+retroTagScopeCondition+` AND id = $3
		RETURNING COALESCE(created_by::text, ''), created_at, updated_at;`,
		stringValue(tag.OrganizationID), stringValue(tag.TeamID), tag.ID, tag.Name,
	).Scan(&tag.CreatedBy, &tag.CreatedAt, &tag.UpdatedAt)
	if err != nil {
		return retroTagWriteError("update retro tag", err)
	}

	return nil
}

// DeleteRetroTag removes a preset tag from a scope, tags already applied to retro items are kept
func (d *Service) DeleteRetroTag(ctx context.Context, organizationID string, teamID string, tagID string) error {
	result, err := d.DB.ExecContext(ctx,
		`DELETE FROM thunderdome.retro_tag WHERE `+retroTagScopeCondition+` AND id = $3;`,
		organizationID, teamID, tagID,
	)
	if err != nil {
		return fmt.Errorf("delete retro tag query error: %v", err)
	}

	if count, err := result.RowsAffected(); err == nil && count == 0 {
		return thunderdome.ErrRetroTagNotFound
	}

	return nil
}

// GetRetroTagOptions gets the preset tag names offered in a retro: the global defaults plus those of the retro's
// team and organization. When subscriptions are enabled the team and organization tags require an active subscription.
func (d *Service) GetRetroTagOptions(ctx context.Context, retroID string, subscriptionsEnabled bool) ([]string, error) {
	options := make([]string, 0)

	rows, err := d.DB.QueryContext(ctx,
		`WITH retro_team AS (
			SELECT t.id AS team_id, COALESCE(t.organization_id, od.organization_id) AS organization_id
			FROM thunderdome.retro r
			JOIN thunderdome.team t ON t.id = r.team_id
			LEFT JOIN thunderdome.organization_department od ON od.id = t.department_id
			WHERE r.id = $1
		),
		subscribed AS (
			SELECT
				NOT $2::boolean OR EXISTS (
					SELECT 1 FROM thunderdome.subscription s, retro_team rt
					WHERE s.organization_id = rt.organization_id AND s.active AND s.expires > CURRENT_TIMESTAMP
				) AS organization,
				NOT $2::boolean OR EXISTS (
					SELECT 1 FROM thunderdome.subscription s, retro_team rt
					WHERE s.team_id = rt.team_id AND s.active AND s.expires > CURRENT_TIMESTAMP
				) AS team
		),
		available AS (
			SELECT name, 1 AS priority FROM thunderdome.retro_tag
			WHERE is_public = true
			UNION ALL
			SELECT tg.name, 2 FROM thunderdome.retro_tag tg, retro_team rt, subscribed s
			WHERE tg.organization_id = rt.organization_id AND s.organization
			UNION ALL
			SELECT tg.name, 3 FROM thunderdome.retro_tag tg, retro_team rt, subscribed s
			WHERE tg.team_id = rt.team_id AND (s.team OR s.organization)
		)
		SELECT name FROM (
			-- the most specific scope wins when names only differ by case
			SELECT DISTINCT ON (lower(name)) name FROM available ORDER BY lower(name), priority DESC
		) options
		ORDER BY lower(name);`,
		retroID, subscriptionsEnabled,
	)
	if err != nil {
		return nil, fmt.Errorf("get retro tag options query error: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, fmt.Errorf("get retro tag options scan error: %v", err)
		}
		options = append(options, name)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("get retro tag options rows error: %v", err)
	}

	return options, nil
}

func retroTagWriteError(action string, err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return thunderdome.ErrRetroTagNotFound
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return thunderdome.ErrDuplicateRetroTag
	}

	return fmt.Errorf("%s query error: %v", action, err)
}

func stringValue(s *string) string {
	if s == nil {
		return ""
	}

	return *s
}
