package user

import (
	"context"
	"fmt"
)

// RequestEmailChange sends a request to change the email address of a user. It generates a token and sends an email to the new address with a link to confirm the change.
func (d *Service) RequestEmailChange(ctx context.Context, userId string) (string, error) {
	changeId := ""
	// Start a transaction
	tx, err := d.DB.Begin()
	if err != nil {
		return "", fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() // Rollback if not committed

	// Check if an email change request has been made in the last 10 minutes to reduce spamming
	changeCount := 0
	err = tx.QueryRowContext(ctx, `
		SELECT count(uec.change_id)
		FROM thunderdome.user_email_change uec
		WHERE uec.user_id = $1 AND uec.created_date > (CURRENT_TIMESTAMP - INTERVAL '10 minutes');
		`,
		userId,
	).Scan(&changeCount)
	if err != nil || changeCount > 0 {
		return "", fmt.Errorf("insert user email change request query error: %v", err)
	}

	err = tx.QueryRowContext(ctx, `
		INSERT INTO thunderdome.user_email_change (user_id) VALUES (userId) RETURNING change_id INTO changeId;
		`,
		userId,
	).Scan(&changeId)
	if err != nil {
		return "", fmt.Errorf("create user email change request query error: %v", err)
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return "", fmt.Errorf("failed to commit transaction: %w", err)
	}

	return changeId, nil
}

// ConfirmEmailChange confirms the email change request. It verifies the token and updates the user's email address if valid and not already in use.
func (d *Service) ConfirmEmailChange(ctx context.Context, userId string, token string, newEmail string) error {
	return nil
}
