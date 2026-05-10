package storyboard

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
	"go.uber.org/zap"
)

func (d *Service) GetColorLegendTemplatesByOrganization(ctx context.Context, organizationID string) ([]*thunderdome.ColorLegendTemplate, error) {
	templates := make([]*thunderdome.ColorLegendTemplate, 0)

	rows, err := d.DB.QueryContext(ctx,
		`SELECT id, name, description, color_legend, is_public, COALESCE(created_by::text, ''),
			COALESCE(organization_id::text, ''), COALESCE(team_id::text, ''), created_at, updated_at
		FROM thunderdome.color_legend_template
		WHERE organization_id = $1
		ORDER BY created_at DESC;`,
		organizationID,
	)
	if err != nil {
		d.Logger.Ctx(ctx).Error("GetColorLegendTemplatesByOrganization query error", zap.Error(err))
		return nil, fmt.Errorf("error querying color legend templates for organization: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		template, scanErr := scanColorLegendTemplate(rows.Scan)
		if scanErr != nil {
			d.Logger.Ctx(ctx).Error("GetColorLegendTemplatesByOrganization row scan error", zap.Error(scanErr))
			continue
		}
		templates = append(templates, template)
	}

	return templates, nil
}

func (d *Service) GetColorLegendTemplatesByTeam(ctx context.Context, teamID string) ([]*thunderdome.ColorLegendTemplate, error) {
	templates := make([]*thunderdome.ColorLegendTemplate, 0)

	rows, err := d.DB.QueryContext(ctx,
		`SELECT id, name, description, color_legend, is_public, COALESCE(created_by::text, ''),
			COALESCE(organization_id::text, ''), COALESCE(team_id::text, ''), created_at, updated_at
		FROM thunderdome.color_legend_template
		WHERE team_id = $1
		ORDER BY created_at DESC;`,
		teamID,
	)
	if err != nil {
		d.Logger.Ctx(ctx).Error("GetColorLegendTemplatesByTeam query error", zap.Error(err))
		return nil, fmt.Errorf("error querying color legend templates for team: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		template, scanErr := scanColorLegendTemplate(rows.Scan)
		if scanErr != nil {
			d.Logger.Ctx(ctx).Error("GetColorLegendTemplatesByTeam row scan error", zap.Error(scanErr))
			continue
		}
		templates = append(templates, template)
	}

	return templates, nil
}

func (d *Service) CreateColorLegendTemplate(ctx context.Context, template *thunderdome.ColorLegendTemplate) (*thunderdome.ColorLegendTemplate, error) {
	colorLegend, marshalErr := json.Marshal(template.ColorLegend)
	if marshalErr != nil {
		return nil, fmt.Errorf("error marshaling color legend template: %v", marshalErr)
	}

	err := d.DB.QueryRowContext(ctx,
		`INSERT INTO thunderdome.color_legend_template (
			name, description, color_legend, is_public, created_by, organization_id, team_id
		)
		VALUES ($1, $2, $3, $4, $5, NULLIF($6, '')::uuid, NULLIF($7, '')::uuid)
		RETURNING id, created_at, updated_at;`,
		template.Name,
		template.Description,
		colorLegend,
		template.IsPublic,
		template.CreatedBy,
		template.OrganizationID,
		template.TeamID,
	).Scan(&template.ID, &template.CreatedAt, &template.UpdatedAt)
	if err != nil {
		d.Logger.Ctx(ctx).Error("CreateColorLegendTemplate query error", zap.Error(err))
		return nil, fmt.Errorf("error creating color legend template: %v", err)
	}

	return template, nil
}

func (d *Service) UpdateOrganizationColorLegendTemplate(ctx context.Context, template *thunderdome.ColorLegendTemplate) (*thunderdome.ColorLegendTemplate, error) {
	colorLegend, marshalErr := json.Marshal(template.ColorLegend)
	if marshalErr != nil {
		return nil, fmt.Errorf("error marshaling color legend template: %v", marshalErr)
	}

	err := d.DB.QueryRowContext(ctx,
		`UPDATE thunderdome.color_legend_template
		SET name = $3, description = $4, color_legend = $5, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND organization_id = $2
		RETURNING updated_at;`,
		template.ID,
		template.OrganizationID,
		template.Name,
		template.Description,
		colorLegend,
	).Scan(&template.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("error updating organization color legend template: %v", err)
	}

	return template, nil
}

func (d *Service) UpdateTeamColorLegendTemplate(ctx context.Context, template *thunderdome.ColorLegendTemplate) (*thunderdome.ColorLegendTemplate, error) {
	colorLegend, marshalErr := json.Marshal(template.ColorLegend)
	if marshalErr != nil {
		return nil, fmt.Errorf("error marshaling color legend template: %v", marshalErr)
	}

	err := d.DB.QueryRowContext(ctx,
		`UPDATE thunderdome.color_legend_template
		SET name = $3, description = $4, color_legend = $5, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND team_id = $2
		RETURNING updated_at;`,
		template.ID,
		template.TeamID,
		template.Name,
		template.Description,
		colorLegend,
	).Scan(&template.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("error updating team color legend template: %v", err)
	}

	return template, nil
}

func (d *Service) DeleteOrganizationColorLegendTemplate(ctx context.Context, organizationID string, templateID string) error {
	_, err := d.DB.ExecContext(ctx,
		`DELETE FROM thunderdome.color_legend_template WHERE id = $1 AND organization_id = $2;`,
		templateID,
		organizationID,
	)
	if err != nil {
		return fmt.Errorf("error deleting organization color legend template: %v", err)
	}

	return nil
}

func (d *Service) DeleteTeamColorLegendTemplate(ctx context.Context, teamID string, templateID string) error {
	_, err := d.DB.ExecContext(ctx,
		`DELETE FROM thunderdome.color_legend_template WHERE id = $1 AND team_id = $2;`,
		templateID,
		teamID,
	)
	if err != nil {
		return fmt.Errorf("error deleting team color legend template: %v", err)
	}

	return nil
}

func scanColorLegendTemplate(scan func(dest ...any) error) (*thunderdome.ColorLegendTemplate, error) {
	template := &thunderdome.ColorLegendTemplate{}
	var rawColorLegend string
	var rawOrganizationID string
	var rawTeamID string

	err := scan(
		&template.ID,
		&template.Name,
		&template.Description,
		&rawColorLegend,
		&template.IsPublic,
		&template.CreatedBy,
		&rawOrganizationID,
		&rawTeamID,
		&template.CreatedAt,
		&template.UpdatedAt,
	)
	if err == nil {
		if unmarshalErr := json.Unmarshal([]byte(rawColorLegend), &template.ColorLegend); unmarshalErr != nil {
			return nil, fmt.Errorf("error unmarshaling color legend template: %v", unmarshalErr)
		}
		if rawOrganizationID != "" {
			template.OrganizationID = &rawOrganizationID
		}
		if rawTeamID != "" {
			template.TeamID = &rawTeamID
		}
	}

	return template, err
}
