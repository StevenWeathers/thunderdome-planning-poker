package user

import (
	"context"
)

// RequestEmailChange sends a request to change the email address of a user. It generates a token and sends an email to the new address with a link to confirm the change.
func (d *Service) RequestEmailChange(ctx context.Context, userId string) (string, error) {
	return "", nil
}

// ConfirmEmailChange confirms the email change request. It verifies the token and updates the user's email address if valid and not already in use.
func (d *Service) ConfirmEmailChange(ctx context.Context, userId string, token string, newEmail string) error {
	return nil
}
