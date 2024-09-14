package storyboard

import (
	"fmt"

	"github.com/StevenWeathers/thunderdome-planning-poker/internal/db"
	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
)

// ConfirmStoryboardFacilitator confirms the user is a facilitator of the storyboard
func (d *Service) ConfirmStoryboardFacilitator(StoryboardID string, UserID string) error {
	var facilitatorId string
	var role string
	err := d.DB.QueryRow("SELECT type FROM thunderdome.users WHERE id = $1", UserID).Scan(&role)
	if err != nil {
		return fmt.Errorf("confirm storyboard facilitator user role query error:%v", err)
	}

	err = d.DB.QueryRow(
		`SELECT user_id FROM thunderdome.storyboard_facilitator WHERE storyboard_id = $1 AND user_id = $2;`,
		StoryboardID, UserID).Scan(&facilitatorId)
	if err != nil && role != thunderdome.AdminUserType {
		return fmt.Errorf("confirm storyboard facilitator query error:%v", err)
	}

	return nil
}

// StoryboardFacilitatorAdd adds a storyboard facilitator
func (d *Service) StoryboardFacilitatorAdd(StoryboardId string, UserID string) (*thunderdome.Storyboard, error) {
	if _, err := d.DB.Exec(
		`INSERT INTO thunderdome.storyboard_facilitator (storyboard_id, user_id) VALUES ($1, $2);`,
		StoryboardId, UserID); err != nil {
		return nil, fmt.Errorf("storyboard add faciliator query error: %v", err)
	}

	storyboard, err := d.GetStoryboard(StoryboardId, "")
	if err != nil {
		return nil, fmt.Errorf("storyboard add facilitator get storyboard error: %v", err)
	}

	return storyboard, nil
}

// StoryboardFacilitatorRemove removes a storyboard facilitator
func (d *Service) StoryboardFacilitatorRemove(StoryboardId string, UserID string) (*thunderdome.Storyboard, error) {
	facilitatorCount := 0
	err := d.DB.QueryRow(
		`SELECT count(user_id) FROM thunderdome.storyboard_facilitator WHERE storyboard_id = $1;`,
		StoryboardId,
	).Scan(&facilitatorCount)
	if err != nil {
		return nil, fmt.Errorf("storyboard remove facilitator query error: %v", err)
	}

	if facilitatorCount == 1 {
		return nil, fmt.Errorf("ONLY_FACILITATOR")
	}

	if _, err := d.DB.Exec(
		`DELETE FROM thunderdome.storyboard_facilitator WHERE storyboard_id = $1 AND user_id = $2;`,
		StoryboardId, UserID); err != nil {
		return nil, fmt.Errorf("storyboard remove facilitator query error: %v", err)
	}

	storyboard, err := d.GetStoryboard(StoryboardId, "")
	if err != nil {
		return nil, fmt.Errorf("storyboard remove facilitator get storyboard error: %v", err)
	}

	return storyboard, nil
}

// GetStoryboardFacilitatorCode retrieve the storyboard facilitator code
func (d *Service) GetStoryboardFacilitatorCode(StoryboardID string) (string, error) {
	var EncryptedCode string

	if err := d.DB.QueryRow(`
		SELECT COALESCE(facilitator_code, '') FROM thunderdome.storyboard
		WHERE id = $1`,
		StoryboardID,
	).Scan(&EncryptedCode); err != nil {
		return "", fmt.Errorf("get storyboard facilitator_code query error: %v", err)
	}

	if EncryptedCode == "" {
		return "", fmt.Errorf("storyboard facilitator_code not set")
	}
	DecryptedCode, codeErr := db.Decrypt(EncryptedCode, d.AESHashKey)
	if codeErr != nil {
		return "", fmt.Errorf("get storyboard facilitator_code decrypt error: %v", codeErr)
	}

	return DecryptedCode, nil
}
