package storyboard

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
	"github.com/jackc/pgx/v5/pgtype"
)

const estimationScaleColumns = `es.id, es.name, COALESCE(es.description, ''), es.scale_type, es.values, es.is_public,
	COALESCE(es.organization_id::text, ''), COALESCE(es.team_id::text, ''), es.default_scale`

// availableEstimationScaleQuery selects the public scales plus those of the team given by $1 and its organization,
// falling back to the organization given by $2. When $3 (subscriptions enabled) is true the team and organization
// scales require an active subscription.
const availableEstimationScaleQuery = `WITH scope_team AS (
		SELECT t.id AS team_id, COALESCE(t.organization_id, od.organization_id) AS organization_id
		FROM thunderdome.team t
		LEFT JOIN thunderdome.organization_department od ON od.id = t.department_id
		WHERE t.id = NULLIF($1, '')::uuid
	),
	scope AS (
		SELECT
			(SELECT team_id FROM scope_team) AS team_id,
			COALESCE((SELECT organization_id FROM scope_team), NULLIF($2, '')::uuid) AS organization_id
	),
	subscribed AS (
		SELECT
			NOT $3::boolean OR EXISTS (
				SELECT 1 FROM thunderdome.subscription s, scope
				WHERE s.organization_id = scope.organization_id AND s.active AND s.expires > CURRENT_TIMESTAMP
			) AS organization,
			NOT $3::boolean OR EXISTS (
				SELECT 1 FROM thunderdome.subscription s, scope
				WHERE s.team_id = scope.team_id AND s.active AND s.expires > CURRENT_TIMESTAMP
			) AS team
	)
	SELECT ` + estimationScaleColumns + `
	FROM thunderdome.estimation_scale es, scope, subscribed
	WHERE es.is_public = true
		OR (es.organization_id = scope.organization_id AND subscribed.organization)
		OR (es.team_id = scope.team_id AND (subscribed.team OR subscribed.organization))`

// the most specific scales are listed first
const estimationScaleOrder = ` ORDER BY (es.team_id IS NULL), es.is_public, lower(es.name);`

// GetStoryboardEstimationScales gets the estimation scales that can be assigned to a storyboard,
// always including the one it currently uses
func (d *Service) GetStoryboardEstimationScales(ctx context.Context, storyboardID string, subscriptionsEnabled bool) ([]*thunderdome.EstimationScale, error) {
	var teamID, currentScaleID string
	if err := d.DB.QueryRowContext(ctx,
		`SELECT COALESCE(team_id::text, ''), COALESCE(estimation_scale_id::text, '')
		FROM thunderdome.storyboard WHERE id = $1;`,
		storyboardID,
	).Scan(&teamID, &currentScaleID); err != nil {
		return nil, fmt.Errorf("get storyboard estimation scales storyboard query error: %w", err)
	}

	scales, err := d.queryEstimationScales(ctx,
		availableEstimationScaleQuery+estimationScaleOrder,
		teamID, "", subscriptionsEnabled,
	)
	if err != nil {
		return nil, err
	}

	// a storyboard can use a scale outside its current scope, such as one chosen when created in a project
	if currentScaleID != "" && !slices.ContainsFunc(scales, func(s *thunderdome.EstimationScale) bool {
		return s.ID == currentScaleID
	}) {
		current, err := d.getEstimationScale(ctx, currentScaleID)
		if err != nil {
			return nil, err
		}
		scales = append(scales, current)
	}

	return scales, nil
}

// ConfirmStoryboardEstimationScale returns an error unless the scale can be assigned to a new storyboard
// in the scope of the team or organization, both being empty for a personal storyboard
func (d *Service) ConfirmStoryboardEstimationScale(ctx context.Context, teamID string, organizationID string, scaleID string, subscriptionsEnabled bool) error {
	scales, err := d.queryEstimationScales(ctx,
		availableEstimationScaleQuery+estimationScaleOrder,
		teamID, organizationID, subscriptionsEnabled,
	)
	if err != nil {
		return err
	}

	if !slices.ContainsFunc(scales, func(s *thunderdome.EstimationScale) bool { return s.ID == scaleID }) {
		return thunderdome.ErrStoryboardEstimationScaleNotFound
	}

	return nil
}

func (d *Service) getEstimationScale(ctx context.Context, scaleID string) (*thunderdome.EstimationScale, error) {
	scales, err := d.queryEstimationScales(ctx,
		`SELECT `+estimationScaleColumns+` FROM thunderdome.estimation_scale es WHERE es.id = $1;`,
		scaleID,
	)
	if err != nil {
		return nil, err
	}
	if len(scales) == 0 {
		return nil, thunderdome.ErrStoryboardEstimationScaleNotFound
	}

	return scales[0], nil
}

func (d *Service) queryEstimationScales(ctx context.Context, query string, args ...any) ([]*thunderdome.EstimationScale, error) {
	scales := make([]*thunderdome.EstimationScale, 0)

	rows, err := d.DB.QueryContext(ctx, query, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return scales, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get estimation scales query error: %w", err)
	}
	defer rows.Close()

	m := pgtype.NewMap()
	for rows.Next() {
		s := &thunderdome.EstimationScale{}
		var values pgtype.Array[string]
		if err := rows.Scan(
			&s.ID, &s.Name, &s.Description, &s.ScaleType, m.SQLScanner(&values), &s.IsPublic,
			&s.OrganizationID, &s.TeamID, &s.DefaultScale,
		); err != nil {
			return nil, fmt.Errorf("get estimation scales scan error: %w", err)
		}
		s.Values = values.Elements
		scales = append(scales, s)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("get estimation scales rows error: %w", err)
	}

	return scales, nil
}
