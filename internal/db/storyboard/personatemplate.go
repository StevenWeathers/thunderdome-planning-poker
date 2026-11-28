package storyboard

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
)

// personaTemplateScopeCondition matches templates belonging to exactly the scope given by $1 (organization) and
// $2 (team), both empty meaning the global defaults
const personaTemplateScopeCondition = `organization_id IS NOT DISTINCT FROM NULLIF($1, '')::uuid
	AND team_id IS NOT DISTINCT FROM NULLIF($2, '')::uuid`

const personaTemplateColumns = `id, name, description, personas, is_public, COALESCE(created_by::text, ''),
	organization_id, team_id, created_at, updated_at`

// GetPersonaTemplates gets the persona templates for a scope
func (d *Service) GetPersonaTemplates(ctx context.Context, organizationID string, teamID string) ([]*thunderdome.PersonaTemplate, error) {
	return d.queryPersonaTemplates(ctx,
		`SELECT `+personaTemplateColumns+`
		FROM thunderdome.persona_template
		WHERE `+personaTemplateScopeCondition+`
		ORDER BY lower(name);`,
		organizationID, teamID,
	)
}

// GetStoryboardPersonaTemplates gets the templates that can be loaded into a storyboard: the global defaults plus those
// of the storyboard's team and organization. When subscriptions are enabled the team and organization templates require
// an active subscription.
func (d *Service) GetStoryboardPersonaTemplates(ctx context.Context, storyboardID string, subscriptionsEnabled bool) ([]*thunderdome.PersonaTemplate, error) {
	return d.queryPersonaTemplates(ctx,
		`WITH storyboard_team AS (
			SELECT t.id AS team_id, COALESCE(t.organization_id, od.organization_id) AS organization_id
			FROM thunderdome.storyboard s
			JOIN thunderdome.team t ON t.id = s.team_id
			LEFT JOIN thunderdome.organization_department od ON od.id = t.department_id
			WHERE s.id = $1
		),
		subscribed AS (
			SELECT
				NOT $2::boolean OR EXISTS (
					SELECT 1 FROM thunderdome.subscription s, storyboard_team st
					WHERE s.organization_id = st.organization_id AND s.active AND s.expires > CURRENT_TIMESTAMP
				) AS organization,
				NOT $2::boolean OR EXISTS (
					SELECT 1 FROM thunderdome.subscription s, storyboard_team st
					WHERE s.team_id = st.team_id AND s.active AND s.expires > CURRENT_TIMESTAMP
				) AS team
		)
		SELECT `+personaTemplateColumns+`
		FROM thunderdome.persona_template
		WHERE is_public = true
			OR organization_id = (SELECT st.organization_id FROM storyboard_team st, subscribed s WHERE s.organization)
			OR team_id = (SELECT st.team_id FROM storyboard_team st, subscribed s WHERE s.team OR s.organization)
		-- the most specific templates are listed first
		ORDER BY (team_id IS NULL), is_public, lower(name);`,
		storyboardID, subscriptionsEnabled,
	)
}

// AddStoryboardPersonasFromTemplate copies the personas of an available template into a storyboard,
// personas whose name is already on the storyboard are left as they are
func (d *Service) AddStoryboardPersonasFromTemplate(ctx context.Context, storyboardID string, templateID string, subscriptionsEnabled bool) ([]*thunderdome.StoryboardPersona, error) {
	templates, err := d.GetStoryboardPersonaTemplates(ctx, storyboardID, subscriptionsEnabled)
	if err != nil {
		return nil, err
	}

	var template *thunderdome.PersonaTemplate
	for _, t := range templates {
		if t.ID == templateID {
			template = t
			break
		}
	}
	if template == nil {
		return nil, thunderdome.ErrPersonaTemplateNotFound
	}

	personas, err := json.Marshal(template.Personas)
	if err != nil {
		return nil, fmt.Errorf("add storyboard personas from template marshal error: %v", err)
	}

	if _, err := d.DB.ExecContext(ctx,
		`INSERT INTO thunderdome.storyboard_persona (storyboard_id, name, role, description)
		SELECT $1, p.name, p.role, p.description
		FROM jsonb_to_recordset($2::jsonb) AS p(name text, role text, description text)
		ON CONFLICT (storyboard_id, name) DO NOTHING;`,
		storyboardID, string(personas),
	); err != nil {
		return nil, fmt.Errorf("add storyboard personas from template query error: %v", err)
	}

	return d.GetStoryboardPersonas(storyboardID), nil
}

// CreatePersonaTemplate creates a persona template in the scope set on the template
func (d *Service) CreatePersonaTemplate(ctx context.Context, template *thunderdome.PersonaTemplate) error {
	personas, err := json.Marshal(template.Personas)
	if err != nil {
		return fmt.Errorf("create persona template marshal error: %v", err)
	}

	template.IsPublic = template.OrganizationID == nil && template.TeamID == nil

	if err := d.DB.QueryRowContext(ctx,
		`INSERT INTO thunderdome.persona_template (
			name, description, personas, is_public, created_by, organization_id, team_id
		)
		VALUES ($1, $2, $3, $4, NULLIF($5, '')::uuid, $6, $7)
		RETURNING id, created_at, updated_at;`,
		template.Name, template.Description, string(personas), template.IsPublic, template.CreatedBy,
		template.OrganizationID, template.TeamID,
	).Scan(&template.ID, &template.CreatedAt, &template.UpdatedAt); err != nil {
		return fmt.Errorf("create persona template query error: %v", err)
	}

	return nil
}

// UpdatePersonaTemplate updates a persona template within the scope set on the template
func (d *Service) UpdatePersonaTemplate(ctx context.Context, template *thunderdome.PersonaTemplate) error {
	personas, err := json.Marshal(template.Personas)
	if err != nil {
		return fmt.Errorf("update persona template marshal error: %v", err)
	}

	template.IsPublic = template.OrganizationID == nil && template.TeamID == nil

	err = d.DB.QueryRowContext(ctx,
		`UPDATE thunderdome.persona_template
		SET name = $4, description = $5, personas = $6, updated_at = CURRENT_TIMESTAMP
		WHERE `+personaTemplateScopeCondition+` AND id = $3
		RETURNING COALESCE(created_by::text, ''), created_at, updated_at;`,
		stringValue(template.OrganizationID), stringValue(template.TeamID), template.ID,
		template.Name, template.Description, string(personas),
	).Scan(&template.CreatedBy, &template.CreatedAt, &template.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return thunderdome.ErrPersonaTemplateNotFound
	}
	if err != nil {
		return fmt.Errorf("update persona template query error: %v", err)
	}

	return nil
}

// DeletePersonaTemplate removes a persona template from a scope, personas already added to storyboards are kept
func (d *Service) DeletePersonaTemplate(ctx context.Context, organizationID string, teamID string, templateID string) error {
	result, err := d.DB.ExecContext(ctx,
		`DELETE FROM thunderdome.persona_template WHERE `+personaTemplateScopeCondition+` AND id = $3;`,
		organizationID, teamID, templateID,
	)
	if err != nil {
		return fmt.Errorf("delete persona template query error: %v", err)
	}

	if count, err := result.RowsAffected(); err == nil && count == 0 {
		return thunderdome.ErrPersonaTemplateNotFound
	}

	return nil
}

func (d *Service) queryPersonaTemplates(ctx context.Context, query string, args ...any) ([]*thunderdome.PersonaTemplate, error) {
	templates := make([]*thunderdome.PersonaTemplate, 0)

	rows, err := d.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("get persona templates query error: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		t := &thunderdome.PersonaTemplate{}
		var personas string
		if err := rows.Scan(
			&t.ID, &t.Name, &t.Description, &personas, &t.IsPublic, &t.CreatedBy,
			&t.OrganizationID, &t.TeamID, &t.CreatedAt, &t.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("get persona templates scan error: %v", err)
		}
		if err := json.Unmarshal([]byte(personas), &t.Personas); err != nil {
			return nil, fmt.Errorf("get persona templates unmarshal error: %v", err)
		}
		templates = append(templates, t)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("get persona templates rows error: %v", err)
	}

	return templates, nil
}

func stringValue(s *string) string {
	if s == nil {
		return ""
	}

	return *s
}
