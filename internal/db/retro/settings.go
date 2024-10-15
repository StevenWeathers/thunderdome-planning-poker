package retro

import (
	"context"
	"database/sql"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
)

// GetSettingsByOrganization retrieves retro settings for an organization
func (s *Service) GetSettingsByOrganization(ctx context.Context, orgID string) (*thunderdome.RetroSettings, error) {
	var settings thunderdome.RetroSettings
	err := s.DB.QueryRowContext(ctx, `
		SELECT id, organization_id, max_votes, allow_multiple_votes, brainstorm_visibility,
		       phase_time_limit_min, phase_auto_advance, allow_cumulative_voting, template_id,
		       join_code, facilitator_code, created_at, updated_at
		FROM thunderdome.retro_settings
		WHERE organization_id = $1`, orgID).Scan(
		&settings.ID, &settings.OrganizationID, &settings.MaxVotes, &settings.AllowMultipleVotes,
		&settings.BrainstormVisibility, &settings.PhaseTimeLimit, &settings.PhaseAutoAdvance,
		&settings.AllowCumulativeVoting, &settings.TemplateID, &settings.JoinCode,
		&settings.FacilitatorCode, &settings.CreatedAt, &settings.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &settings, nil
}

// GetSettingsByDepartment retrieves retro settings for a department
func (s *Service) GetSettingsByDepartment(ctx context.Context, deptID string) (*thunderdome.RetroSettings, error) {
	var settings thunderdome.RetroSettings
	err := s.DB.QueryRowContext(ctx, `
		SELECT id, department_id, max_votes, allow_multiple_votes, brainstorm_visibility,
		       phase_time_limit_min, phase_auto_advance, allow_cumulative_voting, template_id,
		       join_code, facilitator_code, created_at, updated_at
		FROM thunderdome.retro_settings
		WHERE department_id = $1`, deptID).Scan(
		&settings.ID, &settings.DepartmentID, &settings.MaxVotes, &settings.AllowMultipleVotes,
		&settings.BrainstormVisibility, &settings.PhaseTimeLimit, &settings.PhaseAutoAdvance,
		&settings.AllowCumulativeVoting, &settings.TemplateID, &settings.JoinCode,
		&settings.FacilitatorCode, &settings.CreatedAt, &settings.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &settings, nil
}

// GetSettingsByTeam retrieves retro settings for a team
func (s *Service) GetSettingsByTeam(ctx context.Context, teamID string) (*thunderdome.RetroSettings, error) {
	var settings thunderdome.RetroSettings
	err := s.DB.QueryRowContext(ctx, `
		SELECT id, team_id, max_votes, allow_multiple_votes, brainstorm_visibility,
		       phase_time_limit_min, phase_auto_advance, allow_cumulative_voting, template_id,
		       join_code, facilitator_code, created_at, updated_at
		FROM thunderdome.retro_settings
		WHERE team_id = $1`, teamID).Scan(
		&settings.ID, &settings.TeamID, &settings.MaxVotes, &settings.AllowMultipleVotes,
		&settings.BrainstormVisibility, &settings.PhaseTimeLimit, &settings.PhaseAutoAdvance,
		&settings.AllowCumulativeVoting, &settings.TemplateID, &settings.JoinCode,
		&settings.FacilitatorCode, &settings.CreatedAt, &settings.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &settings, nil
}

// CreateSettings creates new retro settings
func (s *Service) CreateSettings(ctx context.Context, settings *thunderdome.RetroSettings) (*thunderdome.RetroSettings, error) {
	err := s.DB.QueryRowContext(ctx, `
		INSERT INTO thunderdome.retro_settings (
			organization_id, department_id, team_id, max_votes, allow_multiple_votes,
			brainstorm_visibility, phase_time_limit_min, phase_auto_advance,
			allow_cumulative_voting, template_id, join_code, facilitator_code
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING id, created_at, updated_at`,
		settings.OrganizationID, settings.DepartmentID, settings.TeamID, settings.MaxVotes,
		settings.AllowMultipleVotes, settings.BrainstormVisibility, settings.PhaseTimeLimit,
		settings.PhaseAutoAdvance, settings.AllowCumulativeVoting, settings.TemplateID,
		settings.JoinCode, settings.FacilitatorCode).Scan(
		&settings.ID, &settings.CreatedAt, &settings.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return settings, nil
}

// UpdateSettings updates existing retro settings
func (s *Service) UpdateSettings(ctx context.Context, settings *thunderdome.RetroSettings) (*thunderdome.RetroSettings, error) {
	err := s.DB.QueryRowContext(ctx, `
		UPDATE thunderdome.retro_settings
		SET max_votes = $1, allow_multiple_votes = $2, brainstorm_visibility = $3,
		    phase_time_limit_min = $4, phase_auto_advance = $5, allow_cumulative_voting = $6,
		    template_id = $7, join_code = $8, facilitator_code = $9, updated_at = CURRENT_TIMESTAMP
		WHERE id = $10 RETURNING created_at, updated_at, organization_id, department_id, team_id`,
		settings.MaxVotes, settings.AllowMultipleVotes, settings.BrainstormVisibility,
		settings.PhaseTimeLimit, settings.PhaseAutoAdvance, settings.AllowCumulativeVoting,
		settings.TemplateID, settings.JoinCode, settings.FacilitatorCode, settings.ID).Scan(
		&settings.CreatedAt, &settings.UpdatedAt, &settings.OrganizationID, &settings.DepartmentID, &settings.TeamID,
	)
	if err != nil {
		return nil, err
	}

	return settings, nil
}

// UpdateOrganizationSettings updates existing organization retro settings
func (s *Service) UpdateOrganizationSettings(ctx context.Context, settings *thunderdome.RetroSettings) (*thunderdome.RetroSettings, error) {
	err := s.DB.QueryRowContext(ctx, `
		UPDATE thunderdome.retro_settings
		SET max_votes = $1, allow_multiple_votes = $2, brainstorm_visibility = $3,
		    phase_time_limit_min = $4, phase_auto_advance = $5, allow_cumulative_voting = $6,
		    template_id = $7, join_code = $8, facilitator_code = $9, updated_at = CURRENT_TIMESTAMP
		WHERE organization_id = $10 RETURNING id, created_at, updated_at`,
		settings.MaxVotes, settings.AllowMultipleVotes, settings.BrainstormVisibility,
		settings.PhaseTimeLimit, settings.PhaseAutoAdvance, settings.AllowCumulativeVoting,
		settings.TemplateID, settings.JoinCode, settings.FacilitatorCode, settings.OrganizationID).Scan(
		&settings.ID, &settings.CreatedAt, &settings.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return settings, nil
}

// UpdateDepartmentSettings updates existing department retro settings
func (s *Service) UpdateDepartmentSettings(ctx context.Context, settings *thunderdome.RetroSettings) (*thunderdome.RetroSettings, error) {
	err := s.DB.QueryRowContext(ctx, `
		UPDATE thunderdome.retro_settings
		SET max_votes = $1, allow_multiple_votes = $2, brainstorm_visibility = $3,
		    phase_time_limit_min = $4, phase_auto_advance = $5, allow_cumulative_voting = $6,
		    template_id = $7, join_code = $8, facilitator_code = $9, updated_at = CURRENT_TIMESTAMP
		WHERE department_id = $10 RETURNING id, created_at, updated_at`,
		settings.MaxVotes, settings.AllowMultipleVotes, settings.BrainstormVisibility,
		settings.PhaseTimeLimit, settings.PhaseAutoAdvance, settings.AllowCumulativeVoting,
		settings.TemplateID, settings.JoinCode, settings.FacilitatorCode, settings.DepartmentID).Scan(
		&settings.ID, &settings.CreatedAt, &settings.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return settings, nil
}

// UpdateTeamSettings updates existing team retro settings
func (s *Service) UpdateTeamSettings(ctx context.Context, settings *thunderdome.RetroSettings) (*thunderdome.RetroSettings, error) {
	err := s.DB.QueryRowContext(ctx, `
		UPDATE thunderdome.retro_settings
		SET max_votes = $1, allow_multiple_votes = $2, brainstorm_visibility = $3,
		    phase_time_limit_min = $4, phase_auto_advance = $5, allow_cumulative_voting = $6,
		    template_id = $7, join_code = $8, facilitator_code = $9, updated_at = CURRENT_TIMESTAMP
		WHERE team_id = $10 RETURNING id, created_at, updated_at`,
		settings.MaxVotes, settings.AllowMultipleVotes, settings.BrainstormVisibility,
		settings.PhaseTimeLimit, settings.PhaseAutoAdvance, settings.AllowCumulativeVoting,
		settings.TemplateID, settings.JoinCode, settings.FacilitatorCode, settings.TeamID).Scan(
		&settings.ID, &settings.CreatedAt, &settings.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return settings, nil
}

// DeleteSettings deletes retro settings
func (s *Service) DeleteSettings(ctx context.Context, id string) error {
	_, err := s.DB.ExecContext(ctx, "DELETE FROM thunderdome.retro_settings WHERE id = $1", id)
	return err
}

// GetSettingsByID retrieves retro settings by ID
func (s *Service) GetSettingsByID(ctx context.Context, id string) (*thunderdome.RetroSettings, error) {
	var settings thunderdome.RetroSettings
	err := s.DB.QueryRowContext(ctx, `
		SELECT id, organization_id, department_id, team_id, max_votes, allow_multiple_votes,
		       brainstorm_visibility, phase_time_limit_min, phase_auto_advance,
		       allow_cumulative_voting, template_id, join_code, facilitator_code,
		       created_at, updated_at
		FROM thunderdome.retro_settings
		WHERE id = $1`, id).Scan(
		&settings.ID, &settings.OrganizationID, &settings.DepartmentID, &settings.TeamID,
		&settings.MaxVotes, &settings.AllowMultipleVotes, &settings.BrainstormVisibility,
		&settings.PhaseTimeLimit, &settings.PhaseAutoAdvance, &settings.AllowCumulativeVoting,
		&settings.TemplateID, &settings.JoinCode, &settings.FacilitatorCode,
		&settings.CreatedAt, &settings.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &settings, nil
}
