package retro

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/StevenWeathers/thunderdome-planning-poker/internal/db"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
)

// GetSettingsByOrganization retrieves retro settings for an organization
func (d *Service) GetSettingsByOrganization(ctx context.Context, orgID string) (*thunderdome.RetroSettings, error) {
	var settings thunderdome.RetroSettings
	var joinCode string
	var facilitatorCode string

	err := d.DB.QueryRowContext(ctx, `
		SELECT id, organization_id, max_votes, allow_multiple_votes, brainstorm_visibility,
		       phase_time_limit_min, phase_auto_advance, allow_cumulative_voting, template_id,
		       join_code, facilitator_code, created_at, updated_at
		FROM thunderdome.retro_settings
		WHERE organization_id = $1`, orgID).Scan(
		&settings.ID, &settings.OrganizationID, &settings.MaxVotes, &settings.AllowMultipleVotes,
		&settings.BrainstormVisibility, &settings.PhaseTimeLimit, &settings.PhaseAutoAdvance,
		&settings.AllowCumulativeVoting, &settings.TemplateID, &joinCode,
		&facilitatorCode, &settings.CreatedAt, &settings.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if joinCode != "" {
		decryptedCode, codeErr := db.Decrypt(joinCode, d.AESHashKey)
		if codeErr != nil {
			return nil, fmt.Errorf("poker settings decode join_code error: %v", codeErr)
		}
		settings.JoinCode = decryptedCode
	}

	if facilitatorCode != "" {
		decryptedCode, codeErr := db.Decrypt(facilitatorCode, d.AESHashKey)
		if codeErr != nil {
			return nil, fmt.Errorf("poker settings decode leader_code error: %v", codeErr)
		}
		settings.FacilitatorCode = decryptedCode
	}

	return &settings, nil
}

// GetSettingsByDepartment retrieves retro settings for a department
func (d *Service) GetSettingsByDepartment(ctx context.Context, deptID string) (*thunderdome.RetroSettings, error) {
	var settings thunderdome.RetroSettings
	var joinCode string
	var facilitatorCode string

	err := d.DB.QueryRowContext(ctx, `
		SELECT id, department_id, max_votes, allow_multiple_votes, brainstorm_visibility,
		       phase_time_limit_min, phase_auto_advance, allow_cumulative_voting, template_id,
		       join_code, facilitator_code, created_at, updated_at
		FROM thunderdome.retro_settings
		WHERE department_id = $1`, deptID).Scan(
		&settings.ID, &settings.DepartmentID, &settings.MaxVotes, &settings.AllowMultipleVotes,
		&settings.BrainstormVisibility, &settings.PhaseTimeLimit, &settings.PhaseAutoAdvance,
		&settings.AllowCumulativeVoting, &settings.TemplateID, &joinCode,
		&facilitatorCode, &settings.CreatedAt, &settings.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if joinCode != "" {
		decryptedCode, codeErr := db.Decrypt(joinCode, d.AESHashKey)
		if codeErr != nil {
			return nil, fmt.Errorf("poker settings decode join_code error: %v", codeErr)
		}
		settings.JoinCode = decryptedCode
	}

	if facilitatorCode != "" {
		decryptedCode, codeErr := db.Decrypt(facilitatorCode, d.AESHashKey)
		if codeErr != nil {
			return nil, fmt.Errorf("poker settings decode leader_code error: %v", codeErr)
		}
		settings.FacilitatorCode = decryptedCode
	}

	return &settings, nil
}

// GetSettingsByTeam retrieves retro settings for a team
func (d *Service) GetSettingsByTeam(ctx context.Context, teamID string) (*thunderdome.RetroSettings, error) {
	var settings thunderdome.RetroSettings
	var joinCode string
	var facilitatorCode string

	err := d.DB.QueryRowContext(ctx, `
		SELECT id, team_id, max_votes, allow_multiple_votes, brainstorm_visibility,
		       phase_time_limit_min, phase_auto_advance, allow_cumulative_voting, template_id,
		       join_code, facilitator_code, created_at, updated_at
		FROM thunderdome.retro_settings
		WHERE team_id = $1`, teamID).Scan(
		&settings.ID, &settings.TeamID, &settings.MaxVotes, &settings.AllowMultipleVotes,
		&settings.BrainstormVisibility, &settings.PhaseTimeLimit, &settings.PhaseAutoAdvance,
		&settings.AllowCumulativeVoting, &settings.TemplateID, &joinCode,
		&facilitatorCode, &settings.CreatedAt, &settings.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if joinCode != "" {
		decryptedCode, codeErr := db.Decrypt(joinCode, d.AESHashKey)
		if codeErr != nil {
			return nil, fmt.Errorf("poker settings decode join_code error: %v", codeErr)
		}
		settings.JoinCode = decryptedCode
	}

	if facilitatorCode != "" {
		decryptedCode, codeErr := db.Decrypt(facilitatorCode, d.AESHashKey)
		if codeErr != nil {
			return nil, fmt.Errorf("poker settings decode leader_code error: %v", codeErr)
		}
		settings.FacilitatorCode = decryptedCode
	}

	return &settings, nil
}

// CreateSettings creates new retro settings
func (d *Service) CreateSettings(ctx context.Context, settings *thunderdome.RetroSettings) (*thunderdome.RetroSettings, error) {
	var encryptedJoinCode string
	var encryptedFacilitatorCode string

	if settings.JoinCode != "" {
		EncryptedCode, codeErr := db.Encrypt(settings.JoinCode, d.AESHashKey)
		if codeErr != nil {
			return nil, fmt.Errorf("poker settings encrypt join_code error: %v", codeErr)
		}
		encryptedJoinCode = EncryptedCode
	}

	if settings.FacilitatorCode != "" {
		EncryptedCode, codeErr := db.Encrypt(settings.FacilitatorCode, d.AESHashKey)
		if codeErr != nil {
			return nil, fmt.Errorf("poker settings encrypt facilitator_code error: %v", codeErr)
		}
		encryptedFacilitatorCode = EncryptedCode
	}

	err := d.DB.QueryRowContext(ctx, `
		INSERT INTO thunderdome.retro_settings (
			organization_id, department_id, team_id, max_votes, allow_multiple_votes,
			brainstorm_visibility, phase_time_limit_min, phase_auto_advance,
			allow_cumulative_voting, template_id, join_code, facilitator_code
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING id, created_at, updated_at`,
		settings.OrganizationID, settings.DepartmentID, settings.TeamID, settings.MaxVotes,
		settings.AllowMultipleVotes, settings.BrainstormVisibility, settings.PhaseTimeLimit,
		settings.PhaseAutoAdvance, settings.AllowCumulativeVoting, settings.TemplateID,
		encryptedJoinCode, encryptedFacilitatorCode).Scan(
		&settings.ID, &settings.CreatedAt, &settings.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return settings, nil
}

// UpdateSettings updates existing retro settings
func (d *Service) UpdateSettings(ctx context.Context, settings *thunderdome.RetroSettings) (*thunderdome.RetroSettings, error) {
	var encryptedJoinCode string
	var encryptedFacilitatorCode string

	if settings.JoinCode != "" {
		EncryptedCode, codeErr := db.Encrypt(settings.JoinCode, d.AESHashKey)
		if codeErr != nil {
			return nil, fmt.Errorf("poker settings encrypt join_code error: %v", codeErr)
		}
		encryptedJoinCode = EncryptedCode
	}

	if settings.FacilitatorCode != "" {
		EncryptedCode, codeErr := db.Encrypt(settings.FacilitatorCode, d.AESHashKey)
		if codeErr != nil {
			return nil, fmt.Errorf("poker settings encrypt facilitator_code error: %v", codeErr)
		}
		encryptedFacilitatorCode = EncryptedCode
	}

	err := d.DB.QueryRowContext(ctx, `
		UPDATE thunderdome.retro_settings
		SET max_votes = $1, allow_multiple_votes = $2, brainstorm_visibility = $3,
		    phase_time_limit_min = $4, phase_auto_advance = $5, allow_cumulative_voting = $6,
		    template_id = $7, join_code = $8, facilitator_code = $9, updated_at = CURRENT_TIMESTAMP
		WHERE id = $10 RETURNING created_at, updated_at, organization_id, department_id, team_id`,
		settings.MaxVotes, settings.AllowMultipleVotes, settings.BrainstormVisibility,
		settings.PhaseTimeLimit, settings.PhaseAutoAdvance, settings.AllowCumulativeVoting,
		settings.TemplateID, encryptedJoinCode, encryptedFacilitatorCode, settings.ID).Scan(
		&settings.CreatedAt, &settings.UpdatedAt, &settings.OrganizationID, &settings.DepartmentID, &settings.TeamID,
	)
	if err != nil {
		return nil, err
	}

	return settings, nil
}

// UpdateOrganizationSettings updates existing organization retro settings
func (d *Service) UpdateOrganizationSettings(ctx context.Context, settings *thunderdome.RetroSettings) (*thunderdome.RetroSettings, error) {
	var encryptedJoinCode string
	var encryptedFacilitatorCode string

	if settings.JoinCode != "" {
		EncryptedCode, codeErr := db.Encrypt(settings.JoinCode, d.AESHashKey)
		if codeErr != nil {
			return nil, fmt.Errorf("poker settings encrypt join_code error: %v", codeErr)
		}
		encryptedJoinCode = EncryptedCode
	}

	if settings.FacilitatorCode != "" {
		EncryptedCode, codeErr := db.Encrypt(settings.FacilitatorCode, d.AESHashKey)
		if codeErr != nil {
			return nil, fmt.Errorf("poker settings encrypt facilitator_code error: %v", codeErr)
		}
		encryptedFacilitatorCode = EncryptedCode
	}

	err := d.DB.QueryRowContext(ctx, `
		UPDATE thunderdome.retro_settings
		SET max_votes = $1, allow_multiple_votes = $2, brainstorm_visibility = $3,
		    phase_time_limit_min = $4, phase_auto_advance = $5, allow_cumulative_voting = $6,
		    template_id = $7, join_code = $8, facilitator_code = $9, updated_at = CURRENT_TIMESTAMP
		WHERE organization_id = $10 RETURNING id, created_at, updated_at`,
		settings.MaxVotes, settings.AllowMultipleVotes, settings.BrainstormVisibility,
		settings.PhaseTimeLimit, settings.PhaseAutoAdvance, settings.AllowCumulativeVoting,
		settings.TemplateID, encryptedJoinCode, encryptedFacilitatorCode, settings.OrganizationID).Scan(
		&settings.ID, &settings.CreatedAt, &settings.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return settings, nil
}

// UpdateDepartmentSettings updates existing department retro settings
func (d *Service) UpdateDepartmentSettings(ctx context.Context, settings *thunderdome.RetroSettings) (*thunderdome.RetroSettings, error) {
	var encryptedJoinCode string
	var encryptedFacilitatorCode string

	if settings.JoinCode != "" {
		EncryptedCode, codeErr := db.Encrypt(settings.JoinCode, d.AESHashKey)
		if codeErr != nil {
			return nil, fmt.Errorf("poker settings encrypt join_code error: %v", codeErr)
		}
		encryptedJoinCode = EncryptedCode
	}

	if settings.FacilitatorCode != "" {
		EncryptedCode, codeErr := db.Encrypt(settings.FacilitatorCode, d.AESHashKey)
		if codeErr != nil {
			return nil, fmt.Errorf("poker settings encrypt facilitator_code error: %v", codeErr)
		}
		encryptedFacilitatorCode = EncryptedCode
	}

	err := d.DB.QueryRowContext(ctx, `
		UPDATE thunderdome.retro_settings
		SET max_votes = $1, allow_multiple_votes = $2, brainstorm_visibility = $3,
		    phase_time_limit_min = $4, phase_auto_advance = $5, allow_cumulative_voting = $6,
		    template_id = $7, join_code = $8, facilitator_code = $9, updated_at = CURRENT_TIMESTAMP
		WHERE department_id = $10 RETURNING id, created_at, updated_at`,
		settings.MaxVotes, settings.AllowMultipleVotes, settings.BrainstormVisibility,
		settings.PhaseTimeLimit, settings.PhaseAutoAdvance, settings.AllowCumulativeVoting,
		settings.TemplateID, encryptedJoinCode, encryptedFacilitatorCode, settings.DepartmentID).Scan(
		&settings.ID, &settings.CreatedAt, &settings.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return settings, nil
}

// UpdateTeamSettings updates existing team retro settings
func (d *Service) UpdateTeamSettings(ctx context.Context, settings *thunderdome.RetroSettings) (*thunderdome.RetroSettings, error) {
	var encryptedJoinCode string
	var encryptedFacilitatorCode string

	if settings.JoinCode != "" {
		EncryptedCode, codeErr := db.Encrypt(settings.JoinCode, d.AESHashKey)
		if codeErr != nil {
			return nil, fmt.Errorf("poker settings encrypt join_code error: %v", codeErr)
		}
		encryptedJoinCode = EncryptedCode
	}

	if settings.FacilitatorCode != "" {
		EncryptedCode, codeErr := db.Encrypt(settings.FacilitatorCode, d.AESHashKey)
		if codeErr != nil {
			return nil, fmt.Errorf("poker settings encrypt facilitator_code error: %v", codeErr)
		}
		encryptedFacilitatorCode = EncryptedCode
	}

	err := d.DB.QueryRowContext(ctx, `
		UPDATE thunderdome.retro_settings
		SET max_votes = $1, allow_multiple_votes = $2, brainstorm_visibility = $3,
		    phase_time_limit_min = $4, phase_auto_advance = $5, allow_cumulative_voting = $6,
		    template_id = $7, join_code = $8, facilitator_code = $9, updated_at = CURRENT_TIMESTAMP
		WHERE team_id = $10 RETURNING id, created_at, updated_at`,
		settings.MaxVotes, settings.AllowMultipleVotes, settings.BrainstormVisibility,
		settings.PhaseTimeLimit, settings.PhaseAutoAdvance, settings.AllowCumulativeVoting,
		settings.TemplateID, encryptedJoinCode, encryptedFacilitatorCode, settings.TeamID).Scan(
		&settings.ID, &settings.CreatedAt, &settings.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return settings, nil
}

// DeleteSettings deletes retro settings
func (d *Service) DeleteSettings(ctx context.Context, id string) error {
	_, err := d.DB.ExecContext(ctx, "DELETE FROM thunderdome.retro_settings WHERE id = $1", id)
	return err
}

// GetSettingsByID retrieves retro settings by ID
func (d *Service) GetSettingsByID(ctx context.Context, id string) (*thunderdome.RetroSettings, error) {
	var settings thunderdome.RetroSettings
	var joinCode string
	var facilitatorCode string

	err := d.DB.QueryRowContext(ctx, `
		SELECT id, organization_id, department_id, team_id, max_votes, allow_multiple_votes,
		       brainstorm_visibility, phase_time_limit_min, phase_auto_advance,
		       allow_cumulative_voting, template_id, join_code, facilitator_code,
		       created_at, updated_at
		FROM thunderdome.retro_settings
		WHERE id = $1`, id).Scan(
		&settings.ID, &settings.OrganizationID, &settings.DepartmentID, &settings.TeamID,
		&settings.MaxVotes, &settings.AllowMultipleVotes, &settings.BrainstormVisibility,
		&settings.PhaseTimeLimit, &settings.PhaseAutoAdvance, &settings.AllowCumulativeVoting,
		&settings.TemplateID, &joinCode, &facilitatorCode,
		&settings.CreatedAt, &settings.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if joinCode != "" {
		decryptedCode, codeErr := db.Decrypt(joinCode, d.AESHashKey)
		if codeErr != nil {
			return nil, fmt.Errorf("poker settings decode join_code error: %v", codeErr)
		}
		settings.JoinCode = decryptedCode
	}

	if facilitatorCode != "" {
		decryptedCode, codeErr := db.Decrypt(facilitatorCode, d.AESHashKey)
		if codeErr != nil {
			return nil, fmt.Errorf("poker settings decode leader_code error: %v", codeErr)
		}
		settings.FacilitatorCode = decryptedCode
	}

	return &settings, nil
}
