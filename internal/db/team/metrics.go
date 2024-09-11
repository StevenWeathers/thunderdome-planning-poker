package team

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
)

// GetOrganizationMetrics retrieves metrics for a specific organization
func (d *OrganizationService) GetOrganizationMetrics(ctx context.Context, organizationID string) (*thunderdome.OrganizationMetrics, error) {
	var metrics thunderdome.OrganizationMetrics

	err := d.DB.QueryRowContext(ctx, `
		WITH org_metrics AS (
			SELECT
				o.id AS organization_id,
				o.name AS organization_name,
				COUNT(DISTINCT od.id) AS department_count,
				COUNT(DISTINCT t.id) AS team_count,
				COUNT(DISTINCT r.id) AS retro_count,
				COUNT(DISTINCT p.id) AS poker_count,
				COUNT(DISTINCT s.id) AS storyboard_count,
				COUNT(DISTINCT tc.id) AS team_checkin_count,
				COUNT(DISTINCT ou.user_id) AS user_count
			FROM
				thunderdome.organization o
			LEFT JOIN thunderdome.organization_department od ON o.id = od.organization_id
			LEFT JOIN thunderdome.team t ON (o.id = t.organization_id OR od.id = t.department_id)
			LEFT JOIN thunderdome.retro r ON t.id = r.team_id
			LEFT JOIN thunderdome.poker p ON t.id = p.team_id
			LEFT JOIN thunderdome.storyboard s ON t.id = s.team_id
			LEFT JOIN thunderdome.team_checkin tc ON t.id = tc.team_id
			LEFT JOIN thunderdome.organization_user ou ON o.id = ou.organization_id
			WHERE o.id = $1
			GROUP BY o.id, o.name
		)
		SELECT
			om.organization_id,
			om.organization_name,
			om.department_count,
			om.team_count,
			om.retro_count,
			om.poker_count,
			om.storyboard_count,
			om.team_checkin_count,
			om.user_count,
			COUNT(DISTINCT es.id) AS estimation_scale_count,
			COUNT(DISTINCT rt.id) AS retro_template_count
		FROM
			org_metrics om
		LEFT JOIN thunderdome.estimation_scale es ON om.organization_id = es.organization_id
		LEFT JOIN thunderdome.retro_template rt ON om.organization_id = rt.organization_id
		GROUP BY
			om.organization_id,
			om.organization_name,
			om.department_count,
			om.team_count,
			om.retro_count,
			om.poker_count,
			om.storyboard_count,
			om.team_checkin_count,
			om.user_count
	`, organizationID).Scan(
		&metrics.OrganizationID,
		&metrics.OrganizationName,
		&metrics.DepartmentCount,
		&metrics.TeamCount,
		&metrics.RetroCount,
		&metrics.PokerCount,
		&metrics.StoryboardCount,
		&metrics.TeamCheckinCount,
		&metrics.UserCount,
		&metrics.EstimationScaleCount,
		&metrics.RetroTemplateCount,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no organization found with ID %s", organizationID)
		}
		return nil, fmt.Errorf("unable to get organization metrics: %v", err)
	}

	return &metrics, nil
}

// GetTeamMetrics retrieves metrics for a specific team
func (d *Service) GetTeamMetrics(ctx context.Context, teamID string) (*thunderdome.TeamMetrics, error) {
	var metrics thunderdome.TeamMetrics

	err := d.DB.QueryRowContext(ctx, `
		WITH team_metrics AS (
			SELECT
				t.id AS team_id,
				t.name AS team_name,
				COALESCE(o.id, '') AS organization_id,
				COALESCE(o.name, '') AS organization_name,
				COALESCE(od.id, '') AS department_id,
				COALESCE(od.name, '') AS department_name,
				COUNT(DISTINCT tu.user_id) AS user_count,
				COUNT(DISTINCT p.id) AS poker_count,
				COUNT(DISTINCT r.id) AS retro_count,
				COUNT(DISTINCT s.id) AS storyboard_count,
				COUNT(DISTINCT tc.id) AS team_checkin_count
			FROM
				thunderdome.team t
			LEFT JOIN thunderdome.organization o ON t.organization_id = o.id
			LEFT JOIN thunderdome.organization_department od ON t.department_id = od.id
			LEFT JOIN thunderdome.team_user tu ON t.id = tu.team_id
			LEFT JOIN thunderdome.poker p ON t.id = p.team_id
			LEFT JOIN thunderdome.retro r ON t.id = r.team_id
			LEFT JOIN thunderdome.storyboard s ON t.id = s.team_id
			LEFT JOIN thunderdome.team_checkin tc ON t.id = tc.team_id
			WHERE t.id = $1
			GROUP BY t.id, t.name, o.id, o.name, od.id, od.name
		)
		SELECT
			tm.team_id,
			tm.team_name,
			tm.organization_id,
			tm.organization_name,
			tm.department_id,
			tm.department_name,
			tm.user_count,
			tm.poker_count,
			tm.retro_count,
			tm.storyboard_count,
			tm.team_checkin_count,
			COUNT(DISTINCT es.id) AS estimation_scale_count,
			COUNT(DISTINCT rt.id) AS retro_template_count
		FROM
			team_metrics tm
		LEFT JOIN thunderdome.estimation_scale es ON tm.team_id = es.team_id
		LEFT JOIN thunderdome.retro_template rt ON tm.team_id = rt.team_id
		GROUP BY
			tm.team_id,
			tm.team_name,
			tm.organization_id,
			tm.organization_name,
			tm.department_id,
			tm.department_name,
			tm.user_count,
			tm.poker_count,
			tm.retro_count,
			tm.storyboard_count,
			tm.team_checkin_count
	`, teamID).Scan(
		&metrics.TeamID,
		&metrics.TeamName,
		&metrics.OrganizationID,
		&metrics.OrganizationName,
		&metrics.DepartmentID,
		&metrics.DepartmentName,
		&metrics.UserCount,
		&metrics.PokerCount,
		&metrics.RetroCount,
		&metrics.StoryboardCount,
		&metrics.TeamCheckinCount,
		&metrics.EstimationScaleCount,
		&metrics.RetroTemplateCount,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no team found with ID %s", teamID)
		}
		return nil, fmt.Errorf("unable to get team metrics: %v", err)
	}

	return &metrics, nil
}
