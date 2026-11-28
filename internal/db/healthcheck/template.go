package healthcheck

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
)

// templateScopeCondition matches templates belonging to exactly the scope given by $1 (organization) and
// $2 (team), both empty meaning the global defaults
const templateScopeCondition = `organization_id IS NOT DISTINCT FROM NULLIF($1, '')::uuid
	AND team_id IS NOT DISTINCT FROM NULLIF($2, '')::uuid`

// availableTemplateCondition matches the templates (aliased ht) that the team given by $1 can use
const availableTemplateCondition = `(ht.is_public
	OR ht.team_id = $1
	OR ht.organization_id = (
		SELECT COALESCE(t.organization_id, od.organization_id)
		FROM thunderdome.team t
		LEFT JOIN thunderdome.organization_department od ON od.id = t.department_id
		WHERE t.id = $1
	))`

const templateColumns = `ht.id, ht.name, ht.description, ht.areas, ht.is_public, COALESCE(ht.created_by::text, ''),
	ht.organization_id, ht.team_id, ht.created_at, ht.updated_at`

// HealthCheckTemplateList gets the templates of a scope
func (d *Service) HealthCheckTemplateList(ctx context.Context, organizationID string, teamID string) ([]*thunderdome.HealthCheckTemplate, error) {
	return d.queryTemplates(ctx,
		`SELECT `+templateColumns+`
		FROM thunderdome.health_check_template ht
		WHERE `+templateScopeCondition+`
		ORDER BY lower(ht.name);`,
		organizationID, teamID,
	)
}

// HealthCheckTemplatesAvailable gets the templates a team can create health checks from
func (d *Service) HealthCheckTemplatesAvailable(ctx context.Context, teamID string) ([]*thunderdome.HealthCheckTemplate, error) {
	return d.queryTemplates(ctx,
		`SELECT `+templateColumns+`
		FROM thunderdome.health_check_template ht
		WHERE `+availableTemplateCondition+`
		-- the most specific templates are listed first
		ORDER BY (ht.team_id IS NULL), ht.is_public, lower(ht.name);`,
		teamID,
	)
}

// HealthCheckTemplateCreate creates a template in the scope set on the template
func (d *Service) HealthCheckTemplateCreate(ctx context.Context, template *thunderdome.HealthCheckTemplate) error {
	areas, err := json.Marshal(template.Areas)
	if err != nil {
		return fmt.Errorf("health check template create marshal error: %v", err)
	}

	template.IsPublic = template.OrganizationID == nil && template.TeamID == nil

	if err := d.DB.QueryRowContext(ctx,
		`INSERT INTO thunderdome.health_check_template (
			name, description, areas, is_public, created_by, organization_id, team_id
		)
		VALUES ($1, $2, $3, $4, NULLIF($5, '')::uuid, $6, $7)
		RETURNING id, created_at, updated_at;`,
		template.Name, template.Description, string(areas), template.IsPublic, template.CreatedBy,
		template.OrganizationID, template.TeamID,
	).Scan(&template.ID, &template.CreatedAt, &template.UpdatedAt); err != nil {
		return fmt.Errorf("health check template create query error: %v", err)
	}

	return nil
}

// HealthCheckTemplateUpdate updates a template within the scope set on the template
func (d *Service) HealthCheckTemplateUpdate(ctx context.Context, template *thunderdome.HealthCheckTemplate) error {
	areas, err := json.Marshal(template.Areas)
	if err != nil {
		return fmt.Errorf("health check template update marshal error: %v", err)
	}

	var organizationID, teamID string
	if template.OrganizationID != nil {
		organizationID = *template.OrganizationID
	}
	if template.TeamID != nil {
		teamID = *template.TeamID
	}
	template.IsPublic = organizationID == "" && teamID == ""

	err = d.DB.QueryRowContext(ctx,
		`UPDATE thunderdome.health_check_template
		SET name = $4, description = $5, areas = $6, updated_at = CURRENT_TIMESTAMP
		WHERE `+templateScopeCondition+` AND id = $3
		RETURNING COALESCE(created_by::text, ''), created_at, updated_at;`,
		organizationID, teamID, template.ID, template.Name, template.Description, string(areas),
	).Scan(&template.CreatedBy, &template.CreatedAt, &template.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return thunderdome.ErrHealthCheckTemplateNotFound
	}
	if err != nil {
		return fmt.Errorf("health check template update query error: %v", err)
	}

	return nil
}

// HealthCheckTemplateDelete removes a template from a scope, health checks created from it are not changed
func (d *Service) HealthCheckTemplateDelete(ctx context.Context, organizationID string, teamID string, templateID string) error {
	result, err := d.DB.ExecContext(ctx,
		`DELETE FROM thunderdome.health_check_template WHERE `+templateScopeCondition+` AND id = $3;`,
		organizationID, teamID, templateID,
	)
	if err != nil {
		return fmt.Errorf("health check template delete query error: %v", err)
	}

	if count, err := result.RowsAffected(); err == nil && count == 0 {
		return thunderdome.ErrHealthCheckTemplateNotFound
	}

	return nil
}

func (d *Service) queryTemplates(ctx context.Context, query string, args ...any) ([]*thunderdome.HealthCheckTemplate, error) {
	templates := make([]*thunderdome.HealthCheckTemplate, 0)

	rows, err := d.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("health check templates query error: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		t := &thunderdome.HealthCheckTemplate{}
		var areas string
		if err := rows.Scan(
			&t.ID, &t.Name, &t.Description, &areas, &t.IsPublic, &t.CreatedBy,
			&t.OrganizationID, &t.TeamID, &t.CreatedAt, &t.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("health check templates scan error: %v", err)
		}
		if err := json.Unmarshal([]byte(areas), &t.Areas); err != nil {
			return nil, fmt.Errorf("health check templates unmarshal error: %v", err)
		}
		templates = append(templates, t)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("health check templates rows error: %v", err)
	}

	return templates, nil
}
