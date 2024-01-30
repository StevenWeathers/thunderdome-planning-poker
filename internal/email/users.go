package email

import (
	"fmt"

	"github.com/matcornic/hermes/v2"
	"go.uber.org/zap"
)

// SendWelcome sends the welcome email to new registered user
func (s *Service) SendWelcome(UserName string, UserEmail string, VerifyID string) error {
	emailBody, err := s.generateBody(
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
						Link:  s.Config.AppURL + "verify-account/" + VerifyID,
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
		s.Logger.Error("Error Generating Welcome Email HTML", zap.Error(err),
			zap.String("user_email", UserEmail))

		return err
	}

	sendErr := s.send(
		UserName,
		UserEmail,
		"Welcome to the Thunderdome!",
		emailBody,
	)
	if sendErr != nil {
		s.Logger.Error("Error sending Welcome Email", zap.Error(sendErr),
			zap.String("user_email", UserEmail))
		return sendErr
	}

	return nil
}

// SendEmailVerification sends the verification email to registered user
func (s *Service) SendEmailVerification(UserName string, UserEmail string, VerifyID string) error {
	emailBody, err := s.generateBody(
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
						Link:  s.Config.AppURL + "verify-account/" + VerifyID,
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
		s.Logger.Error("Error Generating Verification Email HTML", zap.Error(err),
			zap.String("user_email", UserEmail))
		return err
	}

	sendErr := s.send(
		UserName,
		UserEmail,
		"Verify your Thunderdome account email",
		emailBody,
	)
	if sendErr != nil {
		s.Logger.Error("Error sending Verification Email", zap.Error(sendErr),
			zap.String("user_email", UserEmail))
		return sendErr
	}

	return nil
}

// SendForgotPassword Sends a Forgot Password reset email to user
func (s *Service) SendForgotPassword(UserName string, UserEmail string, ResetID string) error {
	emailBody, err := s.generateBody(
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
						Link: s.Config.AppURL + "reset-password/" + ResetID,
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
		s.Logger.Error("Error Generating Forgot Password Email HTML", zap.Error(err),
			zap.String("user_email", UserEmail))
		return err
	}

	sendErr := s.send(
		UserName,
		UserEmail,
		"Forgot your Thunderdome password?",
		emailBody,
	)
	if sendErr != nil {
		s.Logger.Error("Error sending Forgot Password Email", zap.Error(sendErr),
			zap.String("user_email", UserEmail))
		return sendErr
	}

	return nil
}

// SendPasswordReset Sends a Reset Password confirmation email to user
func (s *Service) SendPasswordReset(UserName string, UserEmail string) error {
	emailBody, err := s.generateBody(
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
		s.Logger.Error("Error Generating Reset Password Email HTML", zap.Error(err),
			zap.String("user_email", UserEmail))
		return err
	}

	sendErr := s.send(
		UserName,
		UserEmail,
		"Your Thunderdome password was successfully reset.",
		emailBody,
	)
	if sendErr != nil {
		s.Logger.Error("Error sending Reset Password Email", zap.Error(sendErr),
			zap.String("user_email", UserEmail))
		return sendErr
	}

	return nil
}

// SendPasswordUpdate Sends an Update Password confirmation email to user
func (s *Service) SendPasswordUpdate(UserName string, UserEmail string) error {
	emailBody, err := s.generateBody(
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
		s.Logger.Error("Error Generating Update Password Email HTML", zap.Error(err),
			zap.String("user_email", UserEmail))
		return err
	}

	sendErr := s.send(
		UserName,
		UserEmail,
		"Your Thunderdome password was successfully updated.",
		emailBody,
	)
	if sendErr != nil {
		s.Logger.Error("Error sending Update Password Email", zap.Error(sendErr),
			zap.String("user_email", UserEmail))
		return sendErr
	}

	return nil
}

// SendDeleteConfirmation Sends an delete account confirmation email to user
func (s *Service) SendDeleteConfirmation(UserName string, UserEmail string) error {
	emailBody, err := s.generateBody(
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
		s.Logger.Error("Error Generating Delete Account Confirmation Email HTML", zap.Error(err),
			zap.String("user_email", UserEmail))
		return err
	}

	sendErr := s.send(
		UserName,
		UserEmail,
		"Your Thunderdome account was deleted.",
		emailBody,
	)
	if sendErr != nil {
		s.Logger.Error("Error sending Delete Account Confirmation Email", zap.Error(sendErr),
			zap.String("user_email", UserEmail))
		return sendErr
	}

	return nil
}

// SendEmailUpdate Sends an Update Service confirmation email to user
func (s *Service) SendEmailUpdate(UserName string, UserEmail string) error {
	emailBody, err := s.generateBody(
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
		s.Logger.Error("Error Generating Service Update Email HTML", zap.Error(err),
			zap.String("user_email", UserEmail))
		return err
	}

	sendErr := s.send(
		UserName,
		UserEmail,
		"Your Thunderdome account email has been updated.",
		emailBody,
	)
	if sendErr != nil {
		s.Logger.Error("Error sending Service Update Email", zap.Error(sendErr),
			zap.String("user_email", UserEmail))
		return sendErr
	}

	return nil
}

// SendMergedUpdate Sends an Update Service confirmation email to user
func (s *Service) SendMergedUpdate(UserName string, UserEmail string) error {
	emailBody, err := s.generateBody(
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
		s.Logger.Error("Error Generating Update Merged Service Email HTML", zap.Error(err),
			zap.String("user_email", UserEmail))
		return err
	}

	sendErr := s.send(
		UserName,
		UserEmail,
		"Your Thunderdome duplicate accounts have been merged.",
		emailBody,
	)
	if sendErr != nil {
		s.Logger.Error("Error sending Update Merged Email", zap.Error(sendErr),
			zap.String("user_email", UserEmail))
		return sendErr
	}

	return nil
}

// SendTeamInvite sends the team invite email to unregistered user
func (s *Service) SendTeamInvite(TeamName string, UserEmail string, InviteID string) error {
	emailBody, err := s.generateBody(
		hermes.Body{
			Name: "",
			Intros: []string{
				"Register to join your team on Thunderdome!",
			},
			Actions: []hermes.Action{
				{
					Instructions: fmt.Sprintf(
						"Please register for Thunderdome using the following link (expires in 24 hours) to join the %s Team.",
						TeamName),
					Button: hermes.Button{
						Color: "#22BC66",
						Text:  "Register Account",
						Link:  s.Config.AppURL + "register/team/" + InviteID,
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
		s.Logger.Error("Error Generating Team Invite Email HTML", zap.Error(err),
			zap.String("user_email", UserEmail))

		return err
	}

	sendErr := s.send(
		"",
		UserEmail,
		fmt.Sprintf("Join team %s on Thunderdome!", TeamName),
		emailBody,
	)
	if sendErr != nil {
		s.Logger.Error("Error sending Team Invite Email", zap.Error(sendErr),
			zap.String("user_email", UserEmail),
			zap.String("invite_id", InviteID))
		return sendErr
	}

	return nil
}

// SendOrganizationInvite sends the organization invite email to unregistered user
func (s *Service) SendOrganizationInvite(OrganizationName string, UserEmail string, InviteID string) error {
	emailBody, err := s.generateBody(
		hermes.Body{
			Name: "",
			Intros: []string{
				"Register to join your organization on Thunderdome!",
			},
			Actions: []hermes.Action{
				{
					Instructions: fmt.Sprintf(
						"Please register for Thunderdome using the following link (expires in 24 hours) to join the %s Organization.",
						OrganizationName),
					Button: hermes.Button{
						Color: "#22BC66",
						Text:  "Register Account",
						Link:  s.Config.AppURL + "register/organization/" + InviteID,
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
		s.Logger.Error("Error Generating Organization Invite Email HTML", zap.Error(err),
			zap.String("user_email", UserEmail))

		return err
	}

	sendErr := s.send(
		"",
		UserEmail,
		fmt.Sprintf("Join %s organization on Thunderdome!", OrganizationName),
		emailBody,
	)
	if sendErr != nil {
		s.Logger.Error("Error sending Organization Invite Email", zap.Error(sendErr),
			zap.String("user_email", UserEmail),
			zap.String("invite_id", InviteID))
		return sendErr
	}

	return nil
}
