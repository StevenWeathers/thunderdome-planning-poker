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
		`SELECT id, name, description, format, is_public, COALESCE(created_by::text, ''), created_at, updated_at
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
			&t.Id,
			&t.Name,
			&t.Description,
			&format,
			&t.IsPublic,
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
func (d *Service) GetTemplatesByOrganization(ctx context.Context, organizationId string) ([]*thunderdome.RetroTemplate, error) {
	templates := make([]*thunderdome.RetroTemplate, 0)

	rows, err := d.DB.QueryContext(ctx,
		`SELECT id, name, description, format, is_public, COALESCE(created_by::text, ''), organization_id, created_at, updated_at
		FROM thunderdome.retro_template
		WHERE organization_id = $1;`,
		organizationId,
	)

	if err != nil {
		return nil, fmt.Errorf("error querying templates for organization: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var t thunderdome.RetroTemplate
		var format string
		if err := rows.Scan(
			&t.Id,
			&t.Name,
			&t.Description,
			&format,
			&t.IsPublic,
			&t.CreatedBy,
			&t.OrganizationId,
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
func (d *Service) GetTemplatesByTeam(ctx context.Context, teamId string) ([]*thunderdome.RetroTemplate, error) {
	templates := make([]*thunderdome.RetroTemplate, 0)

	rows, err := d.DB.QueryContext(ctx,
		`SELECT id, name, description, format, is_public, COALESCE(created_by::text, ''), team_id, created_at, updated_at
		FROM thunderdome.retro_template
		WHERE team_id = $1;`,
		teamId,
	)

	if err != nil {
		return nil, fmt.Errorf("error querying templates for team: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var t thunderdome.RetroTemplate
		var format string
		if err := rows.Scan(
			&t.Id,
			&t.Name,
			&t.Description,
			&format,
			&t.IsPublic,
			&t.CreatedBy,
			&t.TeamId,
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

// GetTemplateById retrieves a specific template by its ID
func (d *Service) GetTemplateById(ctx context.Context, templateId string) (*thunderdome.RetroTemplate, error) {
	var t thunderdome.RetroTemplate
	var format string

	err := d.DB.QueryRowContext(ctx,
		`SELECT id, name, description, format, is_public, COALESCE(created_by::text, ''), organization_id, team_id, created_at, updated_at
		FROM thunderdome.retro_template
		WHERE id = $1;`,
		templateId,
	).Scan(
		&t.Id,
		&t.Name,
		&t.Description,
		&format,
		&t.IsPublic,
		&t.CreatedBy,
		&t.OrganizationId,
		&t.TeamId,
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
		`INSERT INTO thunderdome.retro_template (name, description, format, is_public, created_by, organization_id, team_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7);`,
		template.Name,
		template.Description,
		template.Format,
		template.IsPublic,
		template.CreatedBy,
		template.OrganizationId,
		template.TeamId,
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
		SET name = $2, description = $3, format = $4, is_public = $5, organization_id = $6, team_id = $7, updated_at = NOW()
		WHERE id = $1;`,
		template.Id,
		template.Name,
		template.Description,
		template.Format,
		template.IsPublic,
		template.OrganizationId,
		template.TeamId,
	)

	if err != nil {
		return fmt.Errorf("error updating template: %v", err)
	}

	return nil
}

// DeleteTemplate deletes a retro template by its ID
func (d *Service) DeleteTemplate(ctx context.Context, templateId string) error {
	_, err := d.DB.ExecContext(ctx,
		`DELETE FROM thunderdome.retro_template WHERE id = $1;`,
		templateId,
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
		`SELECT id, name, description, format, is_public, COALESCE(created_by::text, ''), organization_id, team_id, created_at, updated_at
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
			&t.Id,
			&t.Name,
			&t.Description,
			&format,
			&t.IsPublic,
			&t.CreatedBy,
			&t.OrganizationId,
			&t.TeamId,
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
