package poker

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/StevenWeathers/thunderdome-planning-poker/internal/db"
	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
	"go.uber.org/zap"
)

// AddFacilitator makes a user a facilitator of the game
func (d *Service) AddFacilitator(pokerID string, userID string) ([]string, error) {
	facilitators := make([]string, 0)

	if _, err := d.DB.Exec(
		`INSERT INTO thunderdome.poker_facilitator (poker_id, user_id) VALUES ($1, $2);`,
		pokerID, userID); err != nil {
		return nil, fmt.Errorf("poker add facilitator query error: %v", err)
	}

	rows, facilitatorErr := d.DB.Query(`
		SELECT user_id FROM thunderdome.poker_facilitator WHERE poker_id = $1;
	`, pokerID)
	if facilitatorErr != nil {
		return facilitators, nil
	}

	defer rows.Close()
	for rows.Next() {
		var leader string
		if err := rows.Scan(
			&leader,
		); err != nil {
			d.Logger.Error("poker_facilitator query scan error", zap.Error(err))
		} else {
			facilitators = append(facilitators, leader)
		}
	}

	return facilitators, nil
}

// RemoveFacilitator removes a user from game facilitators
func (d *Service) RemoveFacilitator(pokerID string, userID string) ([]string, error) {
	facilitators := make([]string, 0)
	facilitatorCount := 0
	err := d.DB.QueryRow(
		`SELECT count(user_id) FROM thunderdome.poker_facilitator WHERE poker_id = $1;`,
		pokerID,
	).Scan(&facilitatorCount)
	if err != nil {
		return nil, fmt.Errorf("poker remove facilitator query error: %v", err)
	}

	if facilitatorCount == 1 {
		return nil, fmt.Errorf("ONLY_FACILITATOR")
	}

	if _, err := d.DB.Exec(
		`DELETE FROM thunderdome.poker_facilitator WHERE poker_id = $1 AND user_id = $2;`,
		pokerID, userID); err != nil {
		return nil, fmt.Errorf("poker remove facilitator query error: %v", err)
	}

	rows, facilitatorErr := d.DB.Query(`
		SELECT user_id FROM thunderdome.poker_facilitator WHERE poker_id = $1;
	`, pokerID)
	if facilitatorErr != nil {
		return facilitators, nil
	}

	defer rows.Close()
	for rows.Next() {
		var leader string
		if err := rows.Scan(
			&leader,
		); err != nil {
			d.Logger.Error("poker_facilitator query scan error", zap.Error(err))
		} else {
			facilitators = append(facilitators, leader)
		}
	}

	return facilitators, nil
}

// AddFacilitatorsByEmail adds additional game facilitators by email
func (d *Service) AddFacilitatorsByEmail(ctx context.Context, pokerID string, facilitatorEmails []string) ([]string, error) {
	var facilitators string
	var newFacilitators []string

	for i, email := range facilitatorEmails {
		facilitatorEmails[i] = db.SanitizeEmail(email)
	}
	emails := strings.Join(facilitatorEmails[:], ",")

	e := d.DB.QueryRowContext(ctx,
		`SELECT facilitators FROM thunderdome.poker_facilitator_add_by_email($1, $2);`,
		pokerID, emails,
	).Scan(&facilitators)
	if e != nil {
		return nil, fmt.Errorf("error adding poker facilitator by email: %v", e)
	}

	_ = json.Unmarshal([]byte(facilitators), &newFacilitators)

	return newFacilitators, nil
}

// ConfirmFacilitator confirms the user is a facilitator of the game
func (d *Service) ConfirmFacilitator(pokerID string, userID string) error {
	var facilitatorID string
	var role string
	err := d.DB.QueryRow("SELECT type FROM thunderdome.users WHERE id = $1", userID).Scan(&role)
	if err != nil {
		return fmt.Errorf("confirm poker facilitator get user role error: %v", err)
	}

	e := d.DB.QueryRow(
		"SELECT user_id FROM thunderdome.poker_facilitator WHERE poker_id = $1 AND user_id = $2",
		pokerID, userID,
	).Scan(&facilitatorID)
	if e != nil && role != thunderdome.AdminUserType {
		return fmt.Errorf("confirm poker facilitator query error: %v", err)
	}

	return nil
}

// GetFacilitatorCode retrieve the game leader_code
func (d *Service) GetFacilitatorCode(pokerID string) (string, error) {
	var encryptedLeaderCode string

	if err := d.DB.QueryRow(`
		SELECT COALESCE(leader_code, '') FROM thunderdome.poker
		WHERE id = $1`,
		pokerID,
	).Scan(&encryptedLeaderCode); err != nil {
		return "", fmt.Errorf("get poker facilitator code query error: %v", err)
	}

	if encryptedLeaderCode == "" {
		return "", fmt.Errorf("poker facilitator code not set")
	}
	decryptedCode, codeErr := db.Decrypt(encryptedLeaderCode, d.AESHashKey)
	if codeErr != nil {
		return "", fmt.Errorf("get poker facilitator code decrypt error: %v", codeErr)
	}

	return decryptedCode, nil
}
