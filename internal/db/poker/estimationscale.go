package poker

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
	"go.uber.org/zap"
)

// GetEstimationScales retrieves a list of estimation scales
func (d *Service) GetEstimationScales(ctx context.Context, limit, offset int) ([]*thunderdome.EstimationScale, int, error) {
	scales := make([]*thunderdome.EstimationScale, 0)
	var scaleCount int

	query := `
		SELECT COUNT(*) FROM thunderdome.estimation_scale;
	`
	err := d.DB.QueryRowContext(ctx, query).Scan(&scaleCount)
	if err != nil {
		return nil, 0, fmt.Errorf("error counting estimation scales: %v", err)
	}

	query = `
		SELECT id, name, description, scale_type, values, COALESCE(created_by::TEXT, ''), created_at, updated_at,
		 is_public, COALESCE(organization_id::TEXT, ''), COALESCE(team_id::TEXT,''), default_scale
		FROM thunderdome.estimation_scale
		ORDER BY name
		LIMIT $1 OFFSET $2;
	`
	rows, err := d.DB.QueryContext(ctx, query, limit, offset)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, 0, fmt.Errorf("error querying estimation scales: %v", err)
	} else if err != nil && errors.Is(err, sql.ErrNoRows) {
		return scales, 0, nil
	}
	defer rows.Close()

	for rows.Next() {
		var vArray pgtype.Array[string]
		m := pgtype.NewMap()
		var s thunderdome.EstimationScale
		err := rows.Scan(
			&s.ID,
			&s.Name,
			&s.Description,
			&s.ScaleType,
			m.SQLScanner(&vArray),
			&s.CreatedBy,
			&s.CreatedAt,
			&s.UpdatedAt,
			&s.IsPublic,
			&s.OrganizationID,
			&s.TeamID,
			&s.DefaultScale,
		)
		if err != nil {
			d.Logger.Ctx(ctx).Error("GetEstimationScales row scan error", zap.Error(err))
		} else {
			s.Values = vArray.Elements
			scales = append(scales, &s)
		}
	}

	return scales, scaleCount, nil
}

// GetPublicEstimationScales retrieves a list of public estimation scales
func (d *Service) GetPublicEstimationScales(ctx context.Context, limit, offset int) ([]*thunderdome.EstimationScale, int, error) {
	scales := make([]*thunderdome.EstimationScale, 0)
	var scaleCount int

	query := `
		SELECT COUNT(*) FROM thunderdome.estimation_scale
		WHERE is_public = true;
	`
	err := d.DB.QueryRowContext(ctx, query).Scan(&scaleCount)
	if err != nil {
		return nil, 0, fmt.Errorf("error counting public estimation scales: %v", err)
	}

	query = `
		SELECT id, name, description, scale_type, values, COALESCE(created_by::TEXT, ''), created_at,
		 updated_at, is_public, COALESCE(organization_id::TEXT, ''), COALESCE(team_id::TEXT, ''), default_scale
		FROM thunderdome.estimation_scale
		WHERE is_public = true
		ORDER BY name
		LIMIT $1 OFFSET $2;
	`
	rows, err := d.DB.QueryContext(ctx, query, limit, offset)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, 0, fmt.Errorf("error querying public estimation scales: %v", err)
	} else if err != nil && errors.Is(err, sql.ErrNoRows) {
		return scales, 0, nil
	}
	defer rows.Close()

	for rows.Next() {
		var s thunderdome.EstimationScale
		var vArray pgtype.Array[string]
		m := pgtype.NewMap()
		err := rows.Scan(
			&s.ID,
			&s.Name,
			&s.Description,
			&s.ScaleType,
			m.SQLScanner(&vArray),
			&s.CreatedBy,
			&s.CreatedAt,
			&s.UpdatedAt,
			&s.IsPublic,
			&s.OrganizationID,
			&s.TeamID,
			&s.DefaultScale,
		)
		if err != nil {
			d.Logger.Ctx(ctx).Error("GetEstimationScales row scan error", zap.Error(err))
		} else {
			s.Values = vArray.Elements
			scales = append(scales, &s)
		}
	}

	return scales, scaleCount, nil
}

// CreateEstimationScale creates a new estimation scale
func (d *Service) CreateEstimationScale(ctx context.Context, scale *thunderdome.EstimationScale) (*thunderdome.EstimationScale, error) {
	query := `
		INSERT INTO thunderdome.estimation_scale
		(name, description, scale_type, values, created_by, is_public, organization_id, team_id, default_scale)
		VALUES ($1, $2, $3, $4, $5, $6, NULLIF($7, '')::uuid, NULLIF($8, '')::uuid, $9)
		RETURNING id, created_at, updated_at;
	`
	err := d.DB.QueryRowContext(ctx, query,
		scale.Name,
		scale.Description,
		scale.ScaleType,
		scale.Values,
		scale.CreatedBy,
		scale.IsPublic,
		scale.OrganizationID,
		scale.TeamID,
		scale.DefaultScale,
	).Scan(&scale.ID, &scale.CreatedAt, &scale.UpdatedAt)

	if err != nil {
		d.Logger.Ctx(ctx).Error("CreateEstimationScale query error", zap.Error(err))
		return nil, fmt.Errorf("error creating new estimation scale: %v", err)
	}

	return scale, nil
}

// UpdateEstimationScale updates an existing estimation scale
func (d *Service) UpdateEstimationScale(ctx context.Context, scale *thunderdome.EstimationScale) (*thunderdome.EstimationScale, error) {
	query := `
		UPDATE thunderdome.estimation_scale
		SET name = $2, description = $3, scale_type = $4, values = $5, is_public = $6,
			organization_id = NULLIF($7, '')::uuid, team_id = NULLIF($8, '')::uuid, default_scale = $9, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		RETURNING updated_at;
	`
	err := d.DB.QueryRowContext(ctx, query,
		scale.ID,
		scale.Name,
		scale.Description,
		scale.ScaleType,
		scale.Values,
		scale.IsPublic,
		scale.OrganizationID,
		scale.TeamID,
		scale.DefaultScale,
	).Scan(&scale.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("error updating estimation scale: %v", err)
	}

	return scale, nil
}

// UpdateTeamEstimationScale updates an existing team estimation scale
func (d *Service) UpdateTeamEstimationScale(ctx context.Context, scale *thunderdome.EstimationScale) (*thunderdome.EstimationScale, error) {
	query := `
		UPDATE thunderdome.estimation_scale
		SET name = $2, description = $3, scale_type = $4, values = $5, is_public = $6,
			default_scale = $8, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND team_id = $7
		RETURNING updated_at;
	`
	err := d.DB.QueryRowContext(ctx, query,
		scale.ID,
		scale.Name,
		scale.Description,
		scale.ScaleType,
		scale.Values,
		scale.IsPublic,
		scale.TeamID,
		scale.DefaultScale,
	).Scan(&scale.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("error updating estimation scale: %v", err)
	}

	return scale, nil
}

// UpdateOrganizationEstimationScale updates an existing organization estimation scale
func (d *Service) UpdateOrganizationEstimationScale(ctx context.Context, scale *thunderdome.EstimationScale) (*thunderdome.EstimationScale, error) {
	query := `
		UPDATE thunderdome.estimation_scale
		SET name = $2, description = $3, scale_type = $4, values = $5, is_public = $6,
			default_scale = $8, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND organization_id = $7
		RETURNING updated_at;
	`
	err := d.DB.QueryRowContext(ctx, query,
		scale.ID,
		scale.Name,
		scale.Description,
		scale.ScaleType,
		scale.Values,
		scale.IsPublic,
		scale.OrganizationID,
		scale.DefaultScale,
	).Scan(&scale.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("error updating estimation scale: %v", err)
	}

	return scale, nil
}

// DeleteEstimationScale deletes an estimation scale
func (d *Service) DeleteEstimationScale(ctx context.Context, scaleID string) error {
	query := `DELETE FROM thunderdome.estimation_scale WHERE id = $1;`
	_, err := d.DB.ExecContext(ctx, query, scaleID)

	if err != nil {
		return fmt.Errorf("error deleting estimation scale: %v", err)
	}

	return nil
}

// DeleteTeamEstimationScale deletes a team's estimation scale
func (d *Service) DeleteTeamEstimationScale(ctx context.Context, teamID string, scaleID string) error {
	query := `DELETE FROM thunderdome.estimation_scale WHERE id = $1 AND team_id = $2;`
	_, err := d.DB.ExecContext(ctx, query, scaleID, teamID)

	if err != nil {
		return fmt.Errorf("error deleting team estimation scale: %v", err)
	}

	return nil
}

// DeleteOrganizationEstimationScale deletes an organization's estimation scale
func (d *Service) DeleteOrganizationEstimationScale(ctx context.Context, orgID string, scaleID string) error {
	query := `DELETE FROM thunderdome.estimation_scale WHERE id = $1 AND organization_id = $2;`
	_, err := d.DB.ExecContext(ctx, query, scaleID, orgID)

	if err != nil {
		return fmt.Errorf("error deleting organization estimation scale: %v", err)
	}

	return nil
}

// GetDefaultPublicEstimationScale retrieves the default public estimation scale
func (d *Service) GetDefaultPublicEstimationScale(ctx context.Context) (*thunderdome.EstimationScale, error) {
	query := `
		SELECT id, name, description, scale_type, values, COALESCE(created_by::TEXT, ''), created_at, updated_at,
		 is_public, default_scale
		FROM thunderdome.estimation_scale
		WHERE default_scale = true AND is_public = true
		LIMIT 1;
	`
	var s thunderdome.EstimationScale
	var vArray pgtype.Array[string]
	m := pgtype.NewMap()
	err := d.DB.QueryRowContext(ctx, query).Scan(
		&s.ID,
		&s.Name,
		&s.Description,
		&s.ScaleType,
		m.SQLScanner(&vArray),
		&s.CreatedBy,
		&s.CreatedAt,
		&s.UpdatedAt,
		&s.IsPublic,
		&s.DefaultScale,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("error retrieving default public estimation scale: %v", err)
	}

	s.Values = vArray.Elements

	return &s, nil
}

// GetDefaultEstimationScale retrieves the default estimation scale for an organization or team
func (d *Service) GetDefaultEstimationScale(ctx context.Context, organizationID, teamID string) (*thunderdome.EstimationScale, error) {
	query := `
		SELECT id, name, description, scale_type, values, COALESCE(created_by::TEXT, ''), created_at, updated_at,
		 is_public, default_scale, COALESCE(organization_id::TEXT, ''), COALESCE(team_id::TEXT,'')
		FROM thunderdome.estimation_scale
		WHERE default_scale = true
		AND (organization_id = $1 OR team_id = $2)
		LIMIT 1;
	`
	var s thunderdome.EstimationScale
	var vArray pgtype.Array[string]
	m := pgtype.NewMap()
	err := d.DB.QueryRowContext(ctx, query, organizationID, teamID).Scan(
		&s.ID,
		&s.Name,
		&s.Description,
		&s.ScaleType,
		m.SQLScanner(&vArray),
		&s.CreatedBy,
		&s.CreatedAt,
		&s.UpdatedAt,
		&s.IsPublic,
		&s.DefaultScale,
		&s.OrganizationID,
		&s.TeamID,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("error retrieving default estimation scale: %v", err)
	}

	s.Values = vArray.Elements

	return &s, nil
}

// GetEstimationScale retrieves an estimation scale by its ID
func (d *Service) GetEstimationScale(ctx context.Context, scaleID string) (*thunderdome.EstimationScale, error) {
	query := `
		SELECT id, name, description, scale_type, values, COALESCE(created_by::TEXT, ''), created_at,
		 updated_at, is_public, default_scale, COALESCE(organization_id::TEXT, ''), COALESCE(team_id::TEXT,'')
		FROM thunderdome.estimation_scale
		WHERE id = $1
	`

	var scale thunderdome.EstimationScale
	var vArray pgtype.Array[string]
	m := pgtype.NewMap()

	err := d.DB.QueryRowContext(ctx, query, scaleID).Scan(
		&scale.ID,
		&scale.Name,
		&scale.Description,
		&scale.ScaleType,
		m.SQLScanner(&vArray),
		&scale.CreatedBy,
		&scale.CreatedAt,
		&scale.UpdatedAt,
		&scale.IsPublic,
		&scale.DefaultScale,
		&scale.OrganizationID,
		&scale.TeamID,
	)

	if err != nil {
		return nil, err
	}

	scale.Values = vArray.Elements

	return &scale, nil
}

// GetPublicEstimationScale retrieves a public estimation scale by its ID
func (d *Service) GetPublicEstimationScale(ctx context.Context, scaleID string) (*thunderdome.EstimationScale, error) {
	query := `
		SELECT id, name, description, scale_type, values, COALESCE(created_by::TEXT, ''), created_at, updated_at, is_public, default_scale
		FROM thunderdome.estimation_scale
		WHERE id = $1 AND is_public = true
	`

	var scale thunderdome.EstimationScale
	var vArray pgtype.Array[string]
	m := pgtype.NewMap()

	err := d.DB.QueryRowContext(ctx, query, scaleID).Scan(
		&scale.ID,
		&scale.Name,
		&scale.Description,
		&scale.ScaleType,
		m.SQLScanner(&vArray),
		&scale.CreatedBy,
		&scale.CreatedAt,
		&scale.UpdatedAt,
		&scale.IsPublic,
		&scale.DefaultScale,
	)

	if err != nil {
		return nil, err
	}

	scale.Values = vArray.Elements

	return &scale, nil
}

// GetOrganizationEstimationScales retrieves estimation scales for a specific organization with pagination
func (d *Service) GetOrganizationEstimationScales(ctx context.Context, orgID string, limit, offset int) ([]*thunderdome.EstimationScale, int, error) {
	scales := make([]*thunderdome.EstimationScale, 0)
	// Query to get the total count of estimation scales for the organization
	countQuery := `
		SELECT COUNT(*)
		FROM thunderdome.estimation_scale
		WHERE organization_id = $1
	`

	var totalCount int
	err := d.DB.QueryRowContext(ctx, countQuery, orgID).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	// Query to get the estimation scales with pagination
	query := `
		SELECT id, name, description, scale_type, values, COALESCE(created_by::TEXT, ''), created_at, updated_at,
		 is_public, organization_id, default_scale
		FROM thunderdome.estimation_scale
		WHERE organization_id = $1
		ORDER BY name
		LIMIT $2 OFFSET $3
	`

	rows, err := d.DB.QueryContext(ctx, query, orgID, limit, offset)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, 0, err
	} else if err != nil && errors.Is(err, sql.ErrNoRows) {
		return scales, 0, nil
	}
	defer rows.Close()

	for rows.Next() {
		var scale thunderdome.EstimationScale
		var vArray pgtype.Array[string]
		m := pgtype.NewMap()
		var orgIDNullable sql.NullString

		err := rows.Scan(
			&scale.ID,
			&scale.Name,
			&scale.Description,
			&scale.ScaleType,
			m.SQLScanner(&vArray),
			&scale.CreatedBy,
			&scale.CreatedAt,
			&scale.UpdatedAt,
			&scale.IsPublic,
			&orgIDNullable,
			&scale.DefaultScale,
		)
		if err != nil {
			return nil, 0, err
		}

		scale.Values = vArray.Elements

		if orgIDNullable.Valid {
			scale.OrganizationID = orgIDNullable.String
		}

		scales = append(scales, &scale)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return scales, totalCount, nil
}

// GetTeamEstimationScales retrieves estimation scales for a specific team with pagination
func (d *Service) GetTeamEstimationScales(ctx context.Context, teamID string, limit, offset int) ([]*thunderdome.EstimationScale, int, error) {
	scales := make([]*thunderdome.EstimationScale, 0)

	// Query to get the total count of estimation scales for the organization
	countQuery := `
		SELECT COUNT(*)
		FROM thunderdome.estimation_scale
		WHERE team_id = $1
	`

	var totalCount int
	err := d.DB.QueryRowContext(ctx, countQuery, teamID).Scan(&totalCount)
	if err != nil {
		d.Logger.Ctx(ctx).Error("GetTeamEstimationScales query error", zap.Error(err))
		return nil, 0, err
	}

	// Query to get the estimation scales with pagination
	query := `
		SELECT id, name, description, scale_type, values, COALESCE(created_by::TEXT, ''), created_at, updated_at,
		 is_public, COALESCE(organization_id::TEXT, ''), default_scale
		FROM thunderdome.estimation_scale
		WHERE team_id = $1
		ORDER BY name
		LIMIT $2 OFFSET $3
	`

	rows, err := d.DB.QueryContext(ctx, query, teamID, limit, offset)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		d.Logger.Ctx(ctx).Error("GetTeamEstimationScales query error", zap.Error(err))
		return nil, 0, err
	} else if err != nil && errors.Is(err, sql.ErrNoRows) {
		return scales, 0, nil
	}
	defer rows.Close()

	for rows.Next() {
		var scale thunderdome.EstimationScale
		var vArray pgtype.Array[string]
		m := pgtype.NewMap()

		err := rows.Scan(
			&scale.ID,
			&scale.Name,
			&scale.Description,
			&scale.ScaleType,
			m.SQLScanner(&vArray),
			&scale.CreatedBy,
			&scale.CreatedAt,
			&scale.UpdatedAt,
			&scale.IsPublic,
			&scale.OrganizationID,
			&scale.DefaultScale,
		)
		if err != nil {
			d.Logger.Ctx(ctx).Error("GetTeamEstimationScales row scan error", zap.Error(err))
			return nil, 0, err
		}

		scale.Values = vArray.Elements

		scales = append(scales, &scale)
	}

	if err = rows.Err(); err != nil {
		d.Logger.Ctx(ctx).Error("GetTeamEstimationScales query error", zap.Error(err))
		return nil, 0, err
	}

	return scales, totalCount, nil
}
