package retrotemplate

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"

	"go.uber.org/zap"
)

// Service represents a PostgreSQL implementation of thunderdome.RetroTemplateDataSvc.
type Service struct {
	DB     *sql.DB
	Logger *otelzap.Logger
}

// GetPublicTemplates retrieves all public retro templates
func (d *Service) GetPublicTemplates(ctx context.Context) ([]*thunderdome.RetroTemplate, error) {
	templates := make([]*thunderdome.RetroTemplate, 0)

	rows, err := d.DB.QueryContext(ctx,
		`SELECT id, name, description, format, is_public, default_template, COALESCE(created_by::text, ''), created_at, updated_at
		FROM thunderdome.retro_template
		WHERE is_public = true;`,
	)

	if err != nil {
		return nil, fmt.Errorf("error querying public templates: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var t thunderdome.RetroTemplate
		var format string
		if err := rows.Scan(
			&t.ID,
			&t.Name,
			&t.Description,
			&format,
			&t.IsPublic,
			&t.DefaultTemplate,
			&t.CreatedBy,
			&t.CreatedAt,
			&t.UpdatedAt,
		); err != nil {
			d.Logger.Ctx(ctx).Error("GetPublicTemplates row scan error", zap.Error(err))
		} else {
			formatErr := json.Unmarshal([]byte(format), &t.Format)
			if formatErr != nil {
				d.Logger.Error("retro template json error", zap.Error(formatErr))
				return nil, fmt.Errorf("get template format error: %v", formatErr)
			}
			templates = append(templates, &t)
		}
	}

	return templates, nil
}

// GetTemplatesByOrganization retrieves all templates for a specific organization
func (d *Service) GetTemplatesByOrganization(ctx context.Context, organizationID string) ([]*thunderdome.RetroTemplate, error) {
	templates := make([]*thunderdome.RetroTemplate, 0)

	rows, err := d.DB.QueryContext(ctx,
		`SELECT id, name, description, format, default_template, COALESCE(created_by::text, ''), organization_id, created_at, updated_at
		FROM thunderdome.retro_template
		WHERE organization_id = $1;`,
		organizationID,
	)

	if err != nil {
		d.Logger.Ctx(ctx).Error("GetTemplatesByOrganization query error", zap.Error(err))
		return nil, fmt.Errorf("error querying templates for organization: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var t thunderdome.RetroTemplate
		var format string
		if err := rows.Scan(
			&t.ID,
			&t.Name,
			&t.Description,
			&format,
			&t.DefaultTemplate,
			&t.CreatedBy,
			&t.OrganizationID,
			&t.CreatedAt,
			&t.UpdatedAt,
		); err != nil {
			d.Logger.Ctx(ctx).Error("GetTemplatesByOrganization row scan error", zap.Error(err))
		} else {
			formatErr := json.Unmarshal([]byte(format), &t.Format)
			if formatErr != nil {
				d.Logger.Error("retro template json error", zap.Error(formatErr))
				return nil, fmt.Errorf("get template format error: %v", formatErr)
			}
			templates = append(templates, &t)
		}
	}

	return templates, nil
}

// GetTemplatesByTeam retrieves all templates for a specific team
func (d *Service) GetTemplatesByTeam(ctx context.Context, teamID string) ([]*thunderdome.RetroTemplate, error) {
	templates := make([]*thunderdome.RetroTemplate, 0)

	rows, err := d.DB.QueryContext(ctx,
		`SELECT id, name, description, format, is_public, default_template, COALESCE(created_by::text, ''), team_id, created_at, updated_at
		FROM thunderdome.retro_template
		WHERE team_id = $1;`,
		teamID,
	)

	if err != nil {
		return nil, fmt.Errorf("error querying templates for team: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var t thunderdome.RetroTemplate
		var format string
		if err := rows.Scan(
			&t.ID,
			&t.Name,
			&t.Description,
			&format,
			&t.IsPublic,
			&t.DefaultTemplate,
			&t.CreatedBy,
			&t.TeamID,
			&t.CreatedAt,
			&t.UpdatedAt,
		); err != nil {
			d.Logger.Ctx(ctx).Error("GetTemplatesByTeam row scan error", zap.Error(err))
		} else {
			formatErr := json.Unmarshal([]byte(format), &t.Format)
			if formatErr != nil {
				d.Logger.Error("retro template json error", zap.Error(formatErr))
				return nil, fmt.Errorf("get template format error: %v", formatErr)
			}
			templates = append(templates, &t)
		}
	}

	return templates, nil
}

// GetTemplateByID retrieves a specific template by its ID
func (d *Service) GetTemplateByID(ctx context.Context, templateID string) (*thunderdome.RetroTemplate, error) {
	var t thunderdome.RetroTemplate
	var format string

	err := d.DB.QueryRowContext(ctx,
		`SELECT id, name, description, format, is_public, default_template, COALESCE(created_by::text, ''), COALESCE(organization_id::text, ''), COALESCE(team_id::text, ''), created_at, updated_at
		FROM thunderdome.retro_template
		WHERE id = $1;`,
		templateID,
	).Scan(
		&t.ID,
		&t.Name,
		&t.Description,
		&format,
		&t.IsPublic,
		&t.DefaultTemplate,
		&t.CreatedBy,
		&t.OrganizationID,
		&t.TeamID,
		&t.CreatedAt,
		&t.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error querying template by ID: %v", err)
	}

	formatErr := json.Unmarshal([]byte(format), &t.Format)
	if formatErr != nil {
		d.Logger.Error("retro template json error", zap.Error(formatErr))
		return nil, fmt.Errorf("get template format error: %v", formatErr)
	}

	return &t, nil
}

// CreateTemplate creates a new retro template
func (d *Service) CreateTemplate(ctx context.Context, template *thunderdome.RetroTemplate) error {
	_, err := d.DB.ExecContext(ctx,
		`INSERT INTO thunderdome.retro_template (
			name, description, format, is_public, default_template, created_by, organization_id, team_id)
		VALUES ($1, $2, $3, $4, $5, $6, NULLIF($7, '')::uuid, NULLIF($8, '')::uuid);`,
		template.Name,
		template.Description,
		template.Format,
		template.IsPublic,
		template.DefaultTemplate,
		template.CreatedBy,
		template.OrganizationID,
		template.TeamID,
	)

	if err != nil {
		return fmt.Errorf("error creating new template: %v", err)
	}

	return nil
}

// UpdateTemplate updates an existing retro template
func (d *Service) UpdateTemplate(ctx context.Context, template *thunderdome.RetroTemplate) error {
	_, err := d.DB.ExecContext(ctx,
		`UPDATE thunderdome.retro_template
		SET name = $2, description = $3, format = $4, is_public = $5, default_template = $6, organization_id = $7, team_id = $8, updated_at = NOW()
		WHERE id = $1;`,
		template.ID,
		template.Name,
		template.Description,
		template.Format,
		template.IsPublic,
		template.DefaultTemplate,
		template.OrganizationID,
		template.TeamID,
	)

	if err != nil {
		return fmt.Errorf("error updating template: %v", err)
	}

	return nil
}

// DeleteTemplate deletes a retro template by its ID
func (d *Service) DeleteTemplate(ctx context.Context, templateID string) error {
	_, err := d.DB.ExecContext(ctx,
		`DELETE FROM thunderdome.retro_template WHERE id = $1;`,
		templateID,
	)

	if err != nil {
		return fmt.Errorf("error deleting template: %v", err)
	}

	return nil
}

// ListTemplates retrieves a paginated list of templates
func (d *Service) ListTemplates(ctx context.Context, limit int, offset int) ([]*thunderdome.RetroTemplate, int, error) {
	templates := make([]*thunderdome.RetroTemplate, 0)
	var totalCount int

	err := d.DB.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM thunderdome.retro_template;",
	).Scan(&totalCount)

	if err != nil {
		d.Logger.Ctx(ctx).Error("ListTemplates count query error", zap.Error(err))
		return nil, 0, fmt.Errorf("error counting templates: %v", err)
	}

	rows, err := d.DB.QueryContext(ctx,
		`SELECT id, name, description, format, is_public, default_template, COALESCE(created_by::text, ''), COALESCE(organization_id::text, ''), COALESCE(team_id::text, ''), created_at, updated_at
		FROM thunderdome.retro_template
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2;`,
		limit,
		offset,
	)

	if err != nil {
		return nil, 0, fmt.Errorf("error querying templates: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var t thunderdome.RetroTemplate
		var format string
		if err := rows.Scan(
			&t.ID,
			&t.Name,
			&t.Description,
			&format,
			&t.IsPublic,
			&t.DefaultTemplate,
			&t.CreatedBy,
			&t.OrganizationID,
			&t.TeamID,
			&t.CreatedAt,
			&t.UpdatedAt,
		); err != nil {
			d.Logger.Ctx(ctx).Error("ListTemplates row scan error", zap.Error(err))
		} else {
			formatErr := json.Unmarshal([]byte(format), &t.Format)
			if formatErr != nil {
				d.Logger.Error("retro template json error", zap.Error(formatErr))
				return nil, totalCount, fmt.Errorf("get template format error: %v", formatErr)
			}

			templates = append(templates, &t)
		}
	}

	return templates, totalCount, nil
}

// GetDefaultPublicTemplate retrieves the default public template
func (d *Service) GetDefaultPublicTemplate(ctx context.Context) (*thunderdome.RetroTemplate, error) {
	var t thunderdome.RetroTemplate
	var format string

	err := d.DB.QueryRowContext(ctx,
		`SELECT id, name, description, format, is_public, default_template, COALESCE(created_by::text, ''), created_at, updated_at
        FROM thunderdome.retro_template
        WHERE is_public = true AND default_template = true
        LIMIT 1;`,
	).Scan(
		&t.ID,
		&t.Name,
		&t.Description,
		&format,
		&t.IsPublic,
		&t.DefaultTemplate,
		&t.CreatedBy,
		&t.CreatedAt,
		&t.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error querying default public template: %v", err)
	}

	formatErr := json.Unmarshal([]byte(format), &t.Format)
	if formatErr != nil {
		d.Logger.Error("retro template json error", zap.Error(formatErr))
		return nil, fmt.Errorf("get default public template format error: %v", formatErr)
	}

	return &t, nil
}

// GetDefaultTeamTemplate retrieves the default template for a given team
func (d *Service) GetDefaultTeamTemplate(ctx context.Context, teamID string) (*thunderdome.RetroTemplate, error) {
	var t thunderdome.RetroTemplate
	var format string

	err := d.DB.QueryRowContext(ctx,
		`SELECT id, name, description, format, is_public, default_template, COALESCE(created_by::text, ''), team_id, created_at, updated_at
        FROM thunderdome.retro_template
        WHERE team_id = $1 AND default_template = true
        LIMIT 1;`,
		teamID,
	).Scan(
		&t.ID,
		&t.Name,
		&t.Description,
		&format,
		&t.IsPublic,
		&t.DefaultTemplate,
		&t.CreatedBy,
		&t.TeamID,
		&t.CreatedAt,
		&t.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No default template found for this team
		}
		return nil, fmt.Errorf("error querying default team template: %v", err)
	}

	formatErr := json.Unmarshal([]byte(format), &t.Format)
	if formatErr != nil {
		d.Logger.Error("retro template json error", zap.Error(formatErr))
		return nil, fmt.Errorf("get default team template format error: %v", formatErr)
	}

	return &t, nil
}

// GetDefaultOrganizationTemplate retrieves the default template for a given organization
func (d *Service) GetDefaultOrganizationTemplate(ctx context.Context, organizationID string) (*thunderdome.RetroTemplate, error) {
	var t thunderdome.RetroTemplate
	var format string

	err := d.DB.QueryRowContext(ctx,
		`SELECT id, name, description, format, is_public, default_template, COALESCE(created_by::text, ''), organization_id, created_at, updated_at
        FROM thunderdome.retro_template
        WHERE organization_id = $1 AND default_template = true
        LIMIT 1;`,
		organizationID,
	).Scan(
		&t.ID,
		&t.Name,
		&t.Description,
		&format,
		&t.IsPublic,
		&t.DefaultTemplate,
		&t.CreatedBy,
		&t.OrganizationID,
		&t.CreatedAt,
		&t.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No default template found for this organization
		}
		return nil, fmt.Errorf("error querying default organization template: %v", err)
	}

	formatErr := json.Unmarshal([]byte(format), &t.Format)
	if formatErr != nil {
		d.Logger.Error("retro template json error", zap.Error(formatErr))
		return nil, fmt.Errorf("get default organization template format error: %v", formatErr)
	}

	return &t, nil
}

// UpdateTeamTemplate updates an existing team retro template
func (d *Service) UpdateTeamTemplate(ctx context.Context, template *thunderdome.RetroTemplate) error {
	_, err := d.DB.ExecContext(ctx,
		`UPDATE thunderdome.retro_template
		SET name = $3, description = $4, format = $5, default_template = $6, updated_at = NOW()
		WHERE id = $1 AND team_id = $2;`,
		template.ID,
		template.TeamID,
		template.Name,
		template.Description,
		template.Format,
		template.DefaultTemplate,
	)

	if err != nil {
		return fmt.Errorf("error updating team template: %v", err)
	}

	return nil
}

// UpdateOrganizationTemplate updates an existing organization retro template
func (d *Service) UpdateOrganizationTemplate(ctx context.Context, template *thunderdome.RetroTemplate) error {
	_, err := d.DB.ExecContext(ctx,
		`UPDATE thunderdome.retro_template
		SET name = $3, description = $4, format = $5, default_template = $6, updated_at = NOW()
		WHERE id = $1 AND organization_id = $2;`,
		template.ID,
		template.OrganizationID,
		template.Name,
		template.Description,
		template.Format,
		template.DefaultTemplate,
	)

	if err != nil {
		return fmt.Errorf("error updating organization template: %v", err)
	}

	return nil
}

// DeleteTeamTemplate deletes a team retro template by its ID
func (d *Service) DeleteTeamTemplate(ctx context.Context, teamID string, templateID string) error {
	_, err := d.DB.ExecContext(ctx,
		`DELETE FROM thunderdome.retro_template WHERE id = $1 AND team_id = $2;`,
		templateID, teamID,
	)

	if err != nil {
		return fmt.Errorf("error deleting team template: %v", err)
	}

	return nil
}

// DeleteOrganizationTemplate deletes an organization retro template by its ID
func (d *Service) DeleteOrganizationTemplate(ctx context.Context, orgID string, templateID string) error {
	_, err := d.DB.ExecContext(ctx,
		`DELETE FROM thunderdome.retro_template WHERE id = $1 AND organization_id = $2;`,
		templateID, orgID,
	)

	if err != nil {
		return fmt.Errorf("error deleting organization template: %v", err)
	}

	return nil
}
