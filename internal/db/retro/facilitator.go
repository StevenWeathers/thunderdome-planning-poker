package retro

import (
	"fmt"

	"github.com/StevenWeathers/thunderdome-planning-poker/internal/db"
	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
	"go.uber.org/zap"
)

// GetRetroFacilitatorCode retrieve the retro facilitator code
func (d *Service) GetRetroFacilitatorCode(retroID string) (string, error) {
	var encryptedCode string

	if err := d.DB.QueryRow(`
		SELECT COALESCE(facilitator_code, '') FROM thunderdome.retro
		WHERE id = $1`,
		retroID,
	).Scan(&encryptedCode); err != nil {
		return "", fmt.Errorf("get retro facilitator_code query error: %v", err)
	}

	if encryptedCode == "" {
		return "", fmt.Errorf("retro facilitator_code not set")
	}
	decryptedCode, codeErr := db.Decrypt(encryptedCode, d.AESHashKey)
	if codeErr != nil {
		return "", fmt.Errorf("retrieve retro facilitator_code decrypt error: %v", codeErr)
	}

	return decryptedCode, nil
}

// RetroConfirmFacilitator confirms the user is a facilitator of the retro
func (d *Service) RetroConfirmFacilitator(retroID string, userID string) error {
	var facilitatorID string
	var role string
	err := d.DB.QueryRow("SELECT type FROM thunderdome.users WHERE id = $1", userID).Scan(&role)
	if err != nil {
		return fmt.Errorf("retro confirm facilitator get user role error: %v", err)
	}

	err = d.DB.QueryRow(
		"SELECT user_id FROM thunderdome.retro_facilitator WHERE retro_id = $1 AND user_id = $2",
		retroID, userID).Scan(&facilitatorID)
	if err != nil && role != thunderdome.AdminUserType {
		return fmt.Errorf("get retro facilitator error: %v", err)
	}

	return nil
}

// GetRetroFacilitators gets a list of retro facilitator ids
func (d *Service) GetRetroFacilitators(retroID string) []string {
	var facilitators = make([]string, 0)
	rows, err := d.DB.Query(
		`SELECT user_id FROM thunderdome.retro_facilitator WHERE retro_id = $1;`,
		retroID,
	)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var facilitator string
			if err := rows.Scan(&facilitator); err != nil {
				d.Logger.Error("get retro facilitators error", zap.Error(err))
			} else {
				facilitators = append(facilitators, facilitator)
			}
		}
	}

	return facilitators
}

// RetroFacilitatorAdd adds a retro facilitator
func (d *Service) RetroFacilitatorAdd(retroID string, userID string) ([]string, error) {
	if _, err := d.DB.Exec(
		`INSERT INTO thunderdome.retro_facilitator (retro_id, user_id) VALUES ($1, $2);`,
		retroID, userID); err != nil {
		return nil, fmt.Errorf("retro add facilitator query error: %v", err)
	}

	facilitators := d.GetRetroFacilitators(retroID)

	return facilitators, nil
}

// RetroFacilitatorRemove removes a retro facilitator
func (d *Service) RetroFacilitatorRemove(retroID string, userID string) ([]string, error) {
	facilitatorCount := 0
	err := d.DB.QueryRow(
		`SELECT count(user_id) FROM thunderdome.retro_facilitator WHERE retro_id = $1;`,
		retroID,
	).Scan(&facilitatorCount)
	if err != nil {
		return nil, fmt.Errorf("retro remove facilitator query error: %v", err)
	}

	if facilitatorCount == 1 {
		return nil, fmt.Errorf("ONLY_FACILITATOR")
	}

	if _, err := d.DB.Exec(
		`DELETE FROM thunderdome.retro_facilitator WHERE retro_id = $1 AND user_id = $2;`,
		retroID, userID); err != nil {
		return nil, fmt.Errorf("retro remove facilitator query error: %v", err)
	}

	facilitators := d.GetRetroFacilitators(retroID)

	return facilitators, nil
}
