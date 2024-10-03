package email

import (
	"fmt"

	"github.com/matcornic/hermes/v2"
	"go.uber.org/zap"
)

// SendWelcome sends the welcome email to new registered user
func (s *Service) SendWelcome(userName string, userEmail string, verifyID string) error {
	emailBody, err := s.generateBody(
		hermes.Body{
			Name: userName,
			Intros: []string{
				"Welcome to the Thunderdome!",
			},
			Actions: []hermes.Action{
				{
					Instructions: "Please validate your email, the following link will expire in 24 hours.",
					Button: hermes.Button{
						Color: "#22BC66",
						Text:  "Verify Account",
						Link:  s.Config.AppURL + "verify-account/" + verifyID,
					},
				},
				{
					Instructions: "Need help, or have questions? Visit our Github page",
					Button: hermes.Button{
						Text: "Github Repo",
						Link: s.Config.RepoURL,
					},
				},
			},
		},
	)
	if err != nil {
		s.Logger.Error("Error Generating Welcome Email HTML", zap.Error(err),
			zap.String("user_email", userEmail))

		return err
	}

	sendErr := s.send(
		userName,
		userEmail,
		"Welcome to the Thunderdome!",
		emailBody,
	)
	if sendErr != nil {
		s.Logger.Error("Error sending Welcome Email", zap.Error(sendErr),
			zap.String("user_email", userEmail))
		return sendErr
	}

	return nil
}

// SendEmailVerification sends the verification email to registered user
func (s *Service) SendEmailVerification(userName string, userEmail string, verifyID string) error {
	emailBody, err := s.generateBody(
		hermes.Body{
			Name: userName,
			Intros: []string{
				"Please verify your Thunderdome account email.",
			},
			Actions: []hermes.Action{
				{
					Instructions: "Please validate your email, the following link will expire in 24 hours.",
					Button: hermes.Button{
						Color: "#22BC66",
						Text:  "Verify Account",
						Link:  s.Config.AppURL + "verify-account/" + verifyID,
					},
				},
				{
					Instructions: "Need help, or have questions? Visit our Github page",
					Button: hermes.Button{
						Text: "Github Repo",
						Link: s.Config.RepoURL,
					},
				},
			},
		},
	)
	if err != nil {
		s.Logger.Error("Error Generating Verification Email HTML", zap.Error(err),
			zap.String("user_email", userEmail))
		return err
	}

	sendErr := s.send(
		userName,
		userEmail,
		"Verify your Thunderdome account email",
		emailBody,
	)
	if sendErr != nil {
		s.Logger.Error("Error sending Verification Email", zap.Error(sendErr),
			zap.String("user_email", userEmail))
		return sendErr
	}

	return nil
}

// SendForgotPassword Sends a Forgot Password reset email to user
func (s *Service) SendForgotPassword(userName string, userEmail string, resetID string) error {
	emailBody, err := s.generateBody(
		hermes.Body{
			Name: userName,
			Intros: []string{
				"It seems you've forgot your Thunderdome password.",
			},
			Actions: []hermes.Action{
				{
					Instructions: "Reset your password now, the following link will expire within an hour of the original request.",
					Button: hermes.Button{
						Text: "Reset Password",
						Link: s.Config.AppURL + "reset-password/" + resetID,
					},
				},
				{
					Instructions: "Need help, or have questions? Visit our Github page",
					Button: hermes.Button{
						Text: "Github Repo",
						Link: s.Config.RepoURL,
					},
				},
			},
		},
	)
	if err != nil {
		s.Logger.Error("Error Generating Forgot Password Email HTML", zap.Error(err),
			zap.String("user_email", userEmail))
		return err
	}

	sendErr := s.send(
		userName,
		userEmail,
		"Forgot your Thunderdome password?",
		emailBody,
	)
	if sendErr != nil {
		s.Logger.Error("Error sending Forgot Password Email", zap.Error(sendErr),
			zap.String("user_email", userEmail))
		return sendErr
	}

	return nil
}

// SendPasswordReset Sends a Reset Password confirmation email to user
func (s *Service) SendPasswordReset(userName string, userEmail string) error {
	emailBody, err := s.generateBody(
		hermes.Body{
			Name: userName,
			Intros: []string{
				"Your Thunderdome password was successfully reset.",
			},
			Actions: []hermes.Action{
				{
					Instructions: "Need help, or have questions? Visit our Github page",
					Button: hermes.Button{
						Text: "Github Repo",
						Link: s.Config.RepoURL,
					},
				},
			},
		},
	)
	if err != nil {
		s.Logger.Error("Error Generating Reset Password Email HTML", zap.Error(err),
			zap.String("user_email", userEmail))
		return err
	}

	sendErr := s.send(
		userName,
		userEmail,
		"Your Thunderdome password was successfully reset.",
		emailBody,
	)
	if sendErr != nil {
		s.Logger.Error("Error sending Reset Password Email", zap.Error(sendErr),
			zap.String("user_email", userEmail))
		return sendErr
	}

	return nil
}

// SendPasswordUpdate Sends an Update Password confirmation email to user
func (s *Service) SendPasswordUpdate(userName string, userEmail string) error {
	emailBody, err := s.generateBody(
		hermes.Body{
			Name: userName,
			Intros: []string{
				"Your Thunderdome password was successfully been updated.",
			},
			Actions: []hermes.Action{
				{
					Instructions: "Need help, or have questions? Visit our Github page",
					Button: hermes.Button{
						Text: "Github Repo",
						Link: s.Config.RepoURL,
					},
				},
			},
		},
	)
	if err != nil {
		s.Logger.Error("Error Generating Update Password Email HTML", zap.Error(err),
			zap.String("user_email", userEmail))
		return err
	}

	sendErr := s.send(
		userName,
		userEmail,
		"Your Thunderdome password was successfully updated.",
		emailBody,
	)
	if sendErr != nil {
		s.Logger.Error("Error sending Update Password Email", zap.Error(sendErr),
			zap.String("user_email", userEmail))
		return sendErr
	}

	return nil
}

// SendDeleteConfirmation Sends an delete account confirmation email to user
func (s *Service) SendDeleteConfirmation(userName string, userEmail string) error {
	emailBody, err := s.generateBody(
		hermes.Body{
			Name: userName,
			Intros: []string{
				"Your Thunderdome account was successfully been deleted.",
			},
			Actions: []hermes.Action{
				{
					Instructions: "Need help, or have questions? Visit our Github page",
					Button: hermes.Button{
						Text: "Github Repo",
						Link: s.Config.RepoURL,
					},
				},
			},
		},
	)
	if err != nil {
		s.Logger.Error("Error Generating Delete Account Confirmation Email HTML", zap.Error(err),
			zap.String("user_email", userEmail))
		return err
	}

	sendErr := s.send(
		userName,
		userEmail,
		"Your Thunderdome account was deleted.",
		emailBody,
	)
	if sendErr != nil {
		s.Logger.Error("Error sending Delete Account Confirmation Email", zap.Error(sendErr),
			zap.String("user_email", userEmail))
		return sendErr
	}

	return nil
}

// SendEmailUpdate Sends an Update Service confirmation email to user
func (s *Service) SendEmailUpdate(userName string, userEmail string) error {
	emailBody, err := s.generateBody(
		hermes.Body{
			Name: userName,
			Intros: []string{
				"Your Thunderdome account email has been lowercased in order to improve unique constraints.",
			},
			Actions: []hermes.Action{
				{
					Instructions: "Need help, or have questions? Visit our Github page",
					Button: hermes.Button{
						Text: "Github Repo",
						Link: s.Config.RepoURL,
					},
				},
			},
		},
	)
	if err != nil {
		s.Logger.Error("Error Generating Service Update Email HTML", zap.Error(err),
			zap.String("user_email", userEmail))
		return err
	}

	sendErr := s.send(
		userName,
		userEmail,
		"Your Thunderdome account email has been updated.",
		emailBody,
	)
	if sendErr != nil {
		s.Logger.Error("Error sending Service Update Email", zap.Error(sendErr),
			zap.String("user_email", userEmail))
		return sendErr
	}

	return nil
}

// SendMergedUpdate Sends an Update Service confirmation email to user
func (s *Service) SendMergedUpdate(userName string, userEmail string) error {
	emailBody, err := s.generateBody(
		hermes.Body{
			Name: userName,
			Intros: []string{
				"Your duplicate Thunderdome accounts under the same email (lowercased) have been merged in order to improve unique constraints. The last active account password was used, in the event you can't login try resetting your password.",
			},
			Actions: []hermes.Action{
				{
					Instructions: "Need help, or have questions? Visit our Github page",
					Button: hermes.Button{
						Text: "Github Repo",
						Link: s.Config.RepoURL,
					},
				},
			},
		},
	)
	if err != nil {
		s.Logger.Error("Error Generating Update Merged Service Email HTML", zap.Error(err),
			zap.String("user_email", userEmail))
		return err
	}

	sendErr := s.send(
		userName,
		userEmail,
		"Your Thunderdome duplicate accounts have been merged.",
		emailBody,
	)
	if sendErr != nil {
		s.Logger.Error("Error sending Update Merged Email", zap.Error(sendErr),
			zap.String("user_email", userEmail))
		return sendErr
	}

	return nil
}

// SendTeamInvite sends the team invite email to user
func (s *Service) SendTeamInvite(teamName string, userEmail string, inviteID string) error {
	subject := fmt.Sprintf("Join team %s on Thunderdome", teamName)
	emailBody, err := s.generateBody(
		hermes.Body{
			Name: "",
			Intros: []string{
				subject,
			},
			Actions: []hermes.Action{
				{
					Instructions: fmt.Sprintf(
						"Please use the following link (expires in 24 hours) to join team %s on Thunderdome today.",
						teamName),
					Button: hermes.Button{
						Color: "#22BC66",
						Text:  "Join Team",
						Link:  s.Config.AppURL + "invite/team/" + inviteID,
					},
				},
				{
					Instructions: "Need help, or have questions? Visit our Github page",
					Button: hermes.Button{
						Text: "Github Repo",
						Link: s.Config.RepoURL,
					},
				},
			},
		},
	)
	if err != nil {
		s.Logger.Error("Error Generating Team Invite Email HTML", zap.Error(err),
			zap.String("user_email", userEmail))

		return err
	}

	sendErr := s.send(
		"",
		userEmail,
		subject,
		emailBody,
	)
	if sendErr != nil {
		s.Logger.Error("Error sending Team Invite Email", zap.Error(sendErr),
			zap.String("user_email", userEmail),
			zap.String("invite_id", inviteID))
		return sendErr
	}

	return nil
}

// SendOrganizationInvite sends the organization invite email to user
func (s *Service) SendOrganizationInvite(organizationName string, userEmail string, inviteID string) error {
	subject := fmt.Sprintf("Join %s organization on Thunderdome", organizationName)
	emailBody, err := s.generateBody(
		hermes.Body{
			Name: "",
			Intros: []string{
				subject,
			},
			Actions: []hermes.Action{
				{
					Instructions: fmt.Sprintf(
						"Please use the following link (expires in 24 hours) to join the %s organization on Thunderdome today.",
						organizationName),
					Button: hermes.Button{
						Color: "#22BC66",
						Text:  "Join Organization",
						Link:  s.Config.AppURL + "invite/organization/" + inviteID,
					},
				},
				{
					Instructions: "Need help, or have questions? Visit our Github page",
					Button: hermes.Button{
						Text: "Github Repo",
						Link: s.Config.RepoURL,
					},
				},
			},
		},
	)
	if err != nil {
		s.Logger.Error("Error Generating Organization Invite Email HTML", zap.Error(err),
			zap.String("user_email", userEmail))

		return err
	}

	sendErr := s.send(
		"",
		userEmail,
		subject,
		emailBody,
	)
	if sendErr != nil {
		s.Logger.Error("Error sending Organization Invite Email", zap.Error(sendErr),
			zap.String("user_email", userEmail),
			zap.String("invite_id", inviteID))
		return sendErr
	}

	return nil
}

// SendDepartmentInvite sends the department invite email to unregistered user
func (s *Service) SendDepartmentInvite(organizationName string, departmentName string, userEmail string, inviteID string) error {
	emailBody, err := s.generateBody(
		hermes.Body{
			Name: "",
			Intros: []string{
				"Register to join your organization's department on Thunderdome!",
			},
			Actions: []hermes.Action{
				{
					Instructions: fmt.Sprintf(
						"Please register for Thunderdome using the following link (expires in 24 hours) to join the %s Organization's %s department.",
						organizationName, departmentName),
					Button: hermes.Button{
						Color: "#22BC66",
						Text:  "Register Account",
						Link:  s.Config.AppURL + "register/department/" + inviteID,
					},
				},
				{
					Instructions: "Need help, or have questions? Visit our Github page",
					Button: hermes.Button{
						Text: "Github Repo",
						Link: s.Config.RepoURL,
					},
				},
			},
		},
	)
	if err != nil {
		s.Logger.Error("Error Generating Department Invite Email HTML", zap.Error(err),
			zap.String("user_email", userEmail))

		return err
	}

	sendErr := s.send(
		"",
		userEmail,
		fmt.Sprintf("Join %s organization's %s department on Thunderdome!", organizationName, departmentName),
		emailBody,
	)
	if sendErr != nil {
		s.Logger.Error("Error sending Department Invite Email", zap.Error(sendErr),
			zap.String("user_email", userEmail),
			zap.String("invite_id", inviteID))
		return sendErr
	}

	return nil
}
