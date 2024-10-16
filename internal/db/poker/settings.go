package poker

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/StevenWeathers/thunderdome-planning-poker/internal/db"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
)

// GetSettingsByOrganization retrieves poker settings for an organization
func (d *Service) GetSettingsByOrganization(ctx context.Context, orgID string) (*thunderdome.PokerSettings, error) {
	var settings thunderdome.PokerSettings
	var joinCode string
	var facilitatorCode string

	err := d.DB.QueryRowContext(ctx, `
		SELECT id, organization_id, auto_finish_voting, point_average_rounding, hide_voter_identity, 
		       estimation_scale_id, join_code, facilitator_code, created_at, updated_at
		FROM thunderdome.poker_settings
		WHERE organization_id = $1`, orgID).Scan(
		&settings.ID, &settings.OrganizationID, &settings.AutoFinishVoting, &settings.PointAverageRounding,
		&settings.HideVoterIdentity, &settings.EstimationScaleID, &joinCode, &facilitatorCode,
		&settings.CreatedAt, &settings.UpdatedAt)

	if errors.Is(err, sql.ErrNoRows) {
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

// GetSettingsByDepartment retrieves poker settings for a department
func (d *Service) GetSettingsByDepartment(ctx context.Context, deptID string) (*thunderdome.PokerSettings, error) {
	var settings thunderdome.PokerSettings
	var joinCode string
	var facilitatorCode string

	err := d.DB.QueryRowContext(ctx, `
		SELECT id, department_id, auto_finish_voting, point_average_rounding, hide_voter_identity, 
		       estimation_scale_id, join_code, facilitator_code, created_at, updated_at
		FROM thunderdome.poker_settings
		WHERE department_id = $1`, deptID).Scan(
		&settings.ID, &settings.DepartmentID, &settings.AutoFinishVoting, &settings.PointAverageRounding,
		&settings.HideVoterIdentity, &settings.EstimationScaleID, &joinCode, &facilitatorCode,
		&settings.CreatedAt, &settings.UpdatedAt)

	if errors.Is(err, sql.ErrNoRows) {
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

// GetSettingsByTeam retrieves poker settings for a team
func (d *Service) GetSettingsByTeam(ctx context.Context, teamID string) (*thunderdome.PokerSettings, error) {
	var settings thunderdome.PokerSettings
	var joinCode string
	var facilitatorCode string

	err := d.DB.QueryRowContext(ctx, `
		SELECT id, team_id, auto_finish_voting, point_average_rounding, hide_voter_identity, 
		       estimation_scale_id, join_code, facilitator_code, created_at, updated_at
		FROM thunderdome.poker_settings
		WHERE team_id = $1`, teamID).Scan(
		&settings.ID, &settings.TeamID, &settings.AutoFinishVoting, &settings.PointAverageRounding,
		&settings.HideVoterIdentity, &settings.EstimationScaleID, &joinCode, &facilitatorCode,
		&settings.CreatedAt, &settings.UpdatedAt)

	if errors.Is(err, sql.ErrNoRows) {
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

// CreateSettings creates new poker settings
func (d *Service) CreateSettings(ctx context.Context, settings *thunderdome.PokerSettings) (*thunderdome.PokerSettings, error) {
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
		INSERT INTO thunderdome.poker_settings (
			organization_id, department_id, team_id, auto_finish_voting, point_average_rounding,
			hide_voter_identity, estimation_scale_id, join_code, facilitator_code
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id, created_at, updated_at`,
		settings.OrganizationID, settings.DepartmentID, settings.TeamID, settings.AutoFinishVoting,
		settings.PointAverageRounding, settings.HideVoterIdentity, settings.EstimationScaleID,
		encryptedJoinCode, encryptedFacilitatorCode).Scan(
		&settings.ID, &settings.CreatedAt, &settings.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return settings, nil
}

// UpdateSettings updates existing poker settings
func (d *Service) UpdateSettings(ctx context.Context, settings *thunderdome.PokerSettings) (*thunderdome.PokerSettings, error) {
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
		UPDATE thunderdome.poker_settings
		SET auto_finish_voting = $1, point_average_rounding = $2, hide_voter_identity = $3,
		    estimation_scale_id = $4, join_code = $5, facilitator_code = $6, updated_at = CURRENT_TIMESTAMP
		WHERE id = $7 RETURNING created_at, updated_at, organization_id, department_id, team_id`,
		settings.AutoFinishVoting, settings.PointAverageRounding, settings.HideVoterIdentity,
		settings.EstimationScaleID, encryptedJoinCode, encryptedFacilitatorCode, settings.ID).Scan(
		&settings.CreatedAt, &settings.UpdatedAt, &settings.OrganizationID, &settings.DepartmentID, &settings.TeamID,
	)
	if err != nil {
		return nil, err
	}
	return settings, nil
}

// UpdateOrganizationSettings updates existing organization poker settings
func (d *Service) UpdateOrganizationSettings(ctx context.Context, settings *thunderdome.PokerSettings) (*thunderdome.PokerSettings, error) {
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
		UPDATE thunderdome.poker_settings
		SET auto_finish_voting = $1, point_average_rounding = $2, hide_voter_identity = $3,
		    estimation_scale_id = $4, join_code = $5, facilitator_code = $6, updated_at = CURRENT_TIMESTAMP
		WHERE organization_id = $7 RETURNING id, created_at, updated_at`,
		settings.AutoFinishVoting, settings.PointAverageRounding, settings.HideVoterIdentity,
		settings.EstimationScaleID, encryptedJoinCode, encryptedFacilitatorCode, settings.OrganizationID).Scan(
		&settings.ID, &settings.CreatedAt, &settings.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return settings, nil
}

// UpdateDepartmentSettings updates existing department poker settings
func (d *Service) UpdateDepartmentSettings(ctx context.Context, settings *thunderdome.PokerSettings) (*thunderdome.PokerSettings, error) {
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
		UPDATE thunderdome.poker_settings
		SET auto_finish_voting = $1, point_average_rounding = $2, hide_voter_identity = $3,
		    estimation_scale_id = $4, join_code = $5, facilitator_code = $6, updated_at = CURRENT_TIMESTAMP
		WHERE department_id = $7 RETURNING id, created_at, updated_at`,
		settings.AutoFinishVoting, settings.PointAverageRounding, settings.HideVoterIdentity,
		settings.EstimationScaleID, encryptedJoinCode, encryptedFacilitatorCode, settings.DepartmentID).Scan(
		&settings.ID, &settings.CreatedAt, &settings.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return settings, nil
}

// UpdateTeamSettings updates existing team poker settings
func (d *Service) UpdateTeamSettings(ctx context.Context, settings *thunderdome.PokerSettings) (*thunderdome.PokerSettings, error) {
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
		UPDATE thunderdome.poker_settings
		SET auto_finish_voting = $1, point_average_rounding = $2, hide_voter_identity = $3,
		    estimation_scale_id = $4, join_code = $5, facilitator_code = $6, updated_at = CURRENT_TIMESTAMP
		WHERE team_id = $7 RETURNING id, created_at, updated_at`,
		settings.AutoFinishVoting, settings.PointAverageRounding, settings.HideVoterIdentity,
		settings.EstimationScaleID, encryptedJoinCode, encryptedFacilitatorCode, settings.TeamID).Scan(
		&settings.ID, &settings.CreatedAt, &settings.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return settings, nil
}

// DeleteSettings deletes poker settings
func (d *Service) DeleteSettings(ctx context.Context, id string) error {
	_, err := d.DB.ExecContext(ctx, "DELETE FROM thunderdome.poker_settings WHERE id = $1", id)
	return err
}

// GetSettingsByID retrieves poker settings by ID
func (d *Service) GetSettingsByID(ctx context.Context, id string) (*thunderdome.PokerSettings, error) {
	var settings thunderdome.PokerSettings
	var joinCode string
	var facilitatorCode string

	err := d.DB.QueryRowContext(ctx, `
		SELECT id, organization_id, department_id, team_id, auto_finish_voting, point_average_rounding,
		       hide_voter_identity, estimation_scale_id, join_code, facilitator_code, created_at, updated_at
		FROM thunderdome.poker_settings
		WHERE id = $1`, id).Scan(
		&settings.ID, &settings.OrganizationID, &settings.DepartmentID, &settings.TeamID,
		&settings.AutoFinishVoting, &settings.PointAverageRounding, &settings.HideVoterIdentity,
		&settings.EstimationScaleID, &joinCode, &facilitatorCode,
		&settings.CreatedAt, &settings.UpdatedAt)

	if errors.Is(err, sql.ErrNoRows) {
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
