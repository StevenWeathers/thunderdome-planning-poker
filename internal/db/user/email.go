package user

import (
	"context"
	"fmt"

	"github.com/StevenWeathers/thunderdome-planning-poker/internal/db"
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

// ConfirmEmailChange confirms the email change request. It verifies the changeId and updates the user's email address if valid and not already in use.
func (d *Service) ConfirmEmailChange(ctx context.Context, userId string, changeId string, newEmail string) error {
	sanitizedEmail := db.SanitizeEmail(newEmail)

	// Start a transaction
	tx, err := d.DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() // Rollback if not committed

	// Delete the email change request
	result, delErr := tx.ExecContext(ctx, `
		DELETE FROM thunderdome.user_email_change WHERE change_id = $1 AND user_id = $2 AND expire_date > CURRENT_TIMESTAMP;
		`,
		changeId, userId,
	)
	if delErr != nil {
		return fmt.Errorf("change user email confirm query error: %v", delErr)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("change user email confirm query error: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no matching email change request found or already expired")
	}

	// Check if the new email is already in use
	var emailCount int
	err = tx.QueryRowContext(ctx, `
		SELECT COUNT(*)
		FROM thunderdome.user
		WHERE email = $1;
		`,
		sanitizedEmail,
	).Scan(&emailCount)
	if err != nil {
		return fmt.Errorf("check new email query error: %v", err)
	}
	if emailCount > 0 {
		return fmt.Errorf("email address already in use")
	}

	// Update the user's email address
	_, err = tx.ExecContext(ctx, `
		UPDATE thunderdome.user
		SET email = $1
		WHERE user_id = $2;
		`,
		sanitizedEmail, userId,
	)
	if err != nil {
		return fmt.Errorf("update user email query error: %v", err)
	}

	// update the user's email in the auth_credential table
	_, err = tx.ExecContext(ctx, `
		UPDATE thunderdome.auth_credential
		SET email = $1
		WHERE user_id = $2;
		`,
		sanitizedEmail, userId,
	)
	if err != nil {
		return fmt.Errorf("update user email in auth_credential query error: %v", err)
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
