package email

import (
	"github.com/matcornic/hermes/v2"
	"go.uber.org/zap"
)

// SendWelcome sends the welcome email to new registered user
func (m *Email) SendWelcome(UserName string, UserEmail string, VerifyID string) error {
	emailBody, err := m.generateBody(
		hermes.Body{
			Name: UserName,
			Intros: []string{
				"Welcome to the Thunderdome! Bring your own mouthguard.",
			},
			Actions: []hermes.Action{
				{
					Instructions: "Please validate your email, the following link will expire in 24 hours.",
					Button: hermes.Button{
						Color: "#22BC66",
						Text:  "Verify Account",
						Link:  m.config.AppURL + "verify-account/" + VerifyID,
					},
				},
				{
					Instructions: "Need help, or have questions? Visit our Github page",
					Button: hermes.Button{
						Text: "Github Repo",
						Link: "https://github.com/StevenWeathers/thunderdome-planning-poker/",
					},
				},
			},
		},
	)
	if err != nil {
		m.logger.Error("Error Generating Welcome Email HTML", zap.Error(err))

		return err
	}

	sendErr := m.Send(
		UserName,
		UserEmail,
		"Welcome to the Thunderdome!",
		emailBody,
	)
	if sendErr != nil {
		m.logger.Error("Error sending Welcome Email", zap.Error(sendErr))
		return sendErr
	}

	return nil
}

// SendEmailVerification sends the verification email to registered user
func (m *Email) SendEmailVerification(UserName string, UserEmail string, VerifyID string) error {
	emailBody, err := m.generateBody(
		hermes.Body{
			Name: UserName,
			Intros: []string{
				"Please verify your Thunderdome account email.",
			},
			Actions: []hermes.Action{
				{
					Instructions: "Please validate your email, the following link will expire in 24 hours.",
					Button: hermes.Button{
						Color: "#22BC66",
						Text:  "Verify Account",
						Link:  m.config.AppURL + "verify-account/" + VerifyID,
					},
				},
				{
					Instructions: "Need help, or have questions? Visit our Github page",
					Button: hermes.Button{
						Text: "Github Repo",
						Link: "https://github.com/StevenWeathers/thunderdome-planning-poker/",
					},
				},
			},
		},
	)
	if err != nil {
		m.logger.Error("Error Generating Verification Email HTML", zap.Error(err))
		return err
	}

	sendErr := m.Send(
		UserName,
		UserEmail,
		"Verify your Thunderdome account email",
		emailBody,
	)
	if sendErr != nil {
		m.logger.Error("Error sending Verification Email", zap.Error(sendErr))
		return sendErr
	}

	return nil
}

// SendForgotPassword Sends a Forgot Password reset email to user
func (m *Email) SendForgotPassword(UserName string, UserEmail string, ResetID string) error {
	emailBody, err := m.generateBody(
		hermes.Body{
			Name: UserName,
			Intros: []string{
				"It seems you've forgot your Thunderdome password.",
			},
			Actions: []hermes.Action{
				{
					Instructions: "Reset your password now, the following link will expire within an hour of the original request.",
					Button: hermes.Button{
						Text: "Reset Password",
						Link: m.config.AppURL + "reset-password/" + ResetID,
					},
				},
				{
					Instructions: "Need help, or have questions? Visit our Github page",
					Button: hermes.Button{
						Text: "Github Repo",
						Link: "https://github.com/StevenWeathers/thunderdome-planning-poker/",
					},
				},
			},
		},
	)
	if err != nil {
		m.logger.Error("Error Generating Forgot Password Email HTML", zap.Error(err))
		return err
	}

	sendErr := m.Send(
		UserName,
		UserEmail,
		"Forgot your Thunderdome password?",
		emailBody,
	)
	if sendErr != nil {
		m.logger.Error("Error sending Forgot Password Email", zap.Error(sendErr))
		return sendErr
	}

	return nil
}

// SendPasswordReset Sends a Reset Password confirmation email to user
func (m *Email) SendPasswordReset(UserName string, UserEmail string) error {
	emailBody, err := m.generateBody(
		hermes.Body{
			Name: UserName,
			Intros: []string{
				"Your Thunderdome password was successfully reset.",
			},
			Actions: []hermes.Action{
				{
					Instructions: "Need help, or have questions? Visit our Github page",
					Button: hermes.Button{
						Text: "Github Repo",
						Link: "https://github.com/StevenWeathers/thunderdome-planning-poker/",
					},
				},
			},
		},
	)
	if err != nil {
		m.logger.Error("Error Generating Reset Password Email HTML", zap.Error(err))
		return err
	}

	sendErr := m.Send(
		UserName,
		UserEmail,
		"Your Thunderdome password was successfully reset.",
		emailBody,
	)
	if sendErr != nil {
		m.logger.Error("Error sending Reset Password Email", zap.Error(sendErr))
		return sendErr
	}

	return nil
}

// SendPasswordUpdate Sends an Update Password confirmation email to user
func (m *Email) SendPasswordUpdate(UserName string, UserEmail string) error {
	emailBody, err := m.generateBody(
		hermes.Body{
			Name: UserName,
			Intros: []string{
				"Your Thunderdome password was successfully been updated.",
			},
			Actions: []hermes.Action{
				{
					Instructions: "Need help, or have questions? Visit our Github page",
					Button: hermes.Button{
						Text: "Github Repo",
						Link: "https://github.com/StevenWeathers/thunderdome-planning-poker/",
					},
				},
			},
		},
	)
	if err != nil {
		m.logger.Error("Error Generating Update Password Email HTML", zap.Error(err))
		return err
	}

	sendErr := m.Send(
		UserName,
		UserEmail,
		"Your Thunderdome password was successfully updated.",
		emailBody,
	)
	if sendErr != nil {
		m.logger.Error("Error sending Update Password Email", zap.Error(sendErr))
		return sendErr
	}

	return nil
}

// SendDeleteConfirmation Sends an delete account confirmation email to user
func (m *Email) SendDeleteConfirmation(UserName string, UserEmail string) error {
	emailBody, err := m.generateBody(
		hermes.Body{
			Name: UserName,
			Intros: []string{
				"Your Thunderdome account was successfully been deleted.",
			},
			Actions: []hermes.Action{
				{
					Instructions: "Need help, or have questions? Visit our Github page",
					Button: hermes.Button{
						Text: "Github Repo",
						Link: "https://github.com/StevenWeathers/thunderdome-planning-poker/",
					},
				},
			},
		},
	)
	if err != nil {
		m.logger.Error("Error Generating Delete Account Confirmation Email HTML", zap.Error(err))
		return err
	}

	sendErr := m.Send(
		UserName,
		UserEmail,
		"Your Thunderdome account was deleted.",
		emailBody,
	)
	if sendErr != nil {
		m.logger.Error("Error sending Delete Account Confirmation Email", zap.Error(sendErr))
		return sendErr
	}

	return nil
}

// SendEmailUpdate Sends an Update Email confirmation email to user
func (m *Email) SendEmailUpdate(UserName string, UserEmail string) error {
	emailBody, err := m.generateBody(
		hermes.Body{
			Name: UserName,
			Intros: []string{
				"Your Thunderdome account email has been lowercased in order to improve unique constraints.",
			},
			Actions: []hermes.Action{
				{
					Instructions: "Need help, or have questions? Visit our Github page",
					Button: hermes.Button{
						Text: "Github Repo",
						Link: "https://github.com/StevenWeathers/thunderdome-planning-poker/",
					},
				},
			},
		},
	)
	if err != nil {
		m.logger.Error("Error Generating Email Update Email HTML", zap.Error(err))
		return err
	}

	sendErr := m.Send(
		UserName,
		UserEmail,
		"Your Thunderdome account email has been updated.",
		emailBody,
	)
	if sendErr != nil {
		m.logger.Error("Error sending Email Update Email", zap.Error(sendErr))
		return sendErr
	}

	return nil
}

// SendMergedUpdate Sends an Update Email confirmation email to user
func (m *Email) SendMergedUpdate(UserName string, UserEmail string) error {
	emailBody, err := m.generateBody(
		hermes.Body{
			Name: UserName,
			Intros: []string{
				"Your duplicate Thunderdome accounts under the same email (lowercased) have been merged in order to improve unique constraints. The last active account password was used, in the event you can't login try resetting your password.",
			},
			Actions: []hermes.Action{
				{
					Instructions: "Need help, or have questions? Visit our Github page",
					Button: hermes.Button{
						Text: "Github Repo",
						Link: "https://github.com/StevenWeathers/thunderdome-planning-poker/",
					},
				},
			},
		},
	)
	if err != nil {
		m.logger.Error("Error Generating Update Merged Email Email HTML", zap.Error(err))
		return err
	}

	sendErr := m.Send(
		UserName,
		UserEmail,
		"Your Thunderdome duplicate accounts have been merged.",
		emailBody,
	)
	if sendErr != nil {
		m.logger.Error("Error sending Update Merged Email", zap.Error(sendErr))
		return sendErr
	}

	return nil
}
