package retro

import (
	"fmt"

	"github.com/StevenWeathers/thunderdome-planning-poker/internal/db"
	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
	"go.uber.org/zap"
)

// GetRetroFacilitatorCode retrieve the retro facilitator code
func (d *Service) GetRetroFacilitatorCode(RetroID string) (string, error) {
	var EncryptedCode string

	if err := d.DB.QueryRow(`
		SELECT COALESCE(facilitator_code, '') FROM thunderdome.retro
		WHERE id = $1`,
		RetroID,
	).Scan(&EncryptedCode); err != nil {
		return "", fmt.Errorf("get retro facilitator_code query error: %v", err)
	}

	if EncryptedCode == "" {
		return "", fmt.Errorf("retro facilitator_code not set")
	}
	DecryptedCode, codeErr := db.Decrypt(EncryptedCode, d.AESHashKey)
	if codeErr != nil {
		return "", fmt.Errorf("retrieve retro facilitator_code decrypt error: %v", codeErr)
	}

	return DecryptedCode, nil
}

// RetroConfirmFacilitator confirms the user is a facilitator of the retro
func (d *Service) RetroConfirmFacilitator(RetroID string, userID string) error {
	var facilitatorId string
	var role string
	err := d.DB.QueryRow("SELECT type FROM thunderdome.users WHERE id = $1", userID).Scan(&role)
	if err != nil {
		return fmt.Errorf("retro confirm facilitator get user role error: %v", err)
	}

	err = d.DB.QueryRow(
		"SELECT user_id FROM thunderdome.retro_facilitator WHERE retro_id = $1 AND user_id = $2",
		RetroID, userID).Scan(&facilitatorId)
	if err != nil && role != thunderdome.AdminUserType {
		return fmt.Errorf("get retro facilitator error: %v", err)
	}

	return nil
}

// GetRetroFacilitators gets a list of retro facilitator ids
func (d *Service) GetRetroFacilitators(RetroID string) []string {
	var facilitators = make([]string, 0)
	rows, err := d.DB.Query(
		`SELECT user_id FROM thunderdome.retro_facilitator WHERE retro_id = $1;`,
		RetroID,
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
func (d *Service) RetroFacilitatorAdd(RetroID string, UserID string) ([]string, error) {
	if _, err := d.DB.Exec(
		`INSERT INTO thunderdome.retro_facilitator (retro_id, user_id) VALUES ($1, $2);`,
		RetroID, UserID); err != nil {
		return nil, fmt.Errorf("retro add facilitator query error: %v", err)
	}

	facilitators := d.GetRetroFacilitators(RetroID)

	return facilitators, nil
}

// RetroFacilitatorRemove removes a retro facilitator
func (d *Service) RetroFacilitatorRemove(RetroID string, UserID string) ([]string, error) {
	facilitatorCount := 0
	err := d.DB.QueryRow(
		`SELECT count(user_id) FROM thunderdome.retro_facilitator WHERE retro_id = $1;`,
		RetroID,
	).Scan(&facilitatorCount)
	if err != nil {
		return nil, fmt.Errorf("retro remove facilitator query error: %v", err)
	}

	if facilitatorCount == 1 {
		return nil, fmt.Errorf("ONLY_FACILITATOR")
	}

	if _, err := d.DB.Exec(
		`DELETE FROM thunderdome.retro_facilitator WHERE retro_id = $1 AND user_id = $2;`,
		RetroID, UserID); err != nil {
		return nil, fmt.Errorf("retro remove facilitator query error: %v", err)
	}

	facilitators := d.GetRetroFacilitators(RetroID)

	return facilitators, nil
}
