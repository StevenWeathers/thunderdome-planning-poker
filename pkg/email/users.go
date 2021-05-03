package email

import (
	"log"

	"github.com/matcornic/hermes/v2"
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
		log.Println("Error Generating Welcome Email HTML: ", err)
		return err
	}

	sendErr := m.Send(
		UserName,
		UserEmail,
		"Welcome to the Thunderdome!",
		emailBody,
	)
	if sendErr != nil {
		log.Println("Error sending Welcome Email: ", sendErr)
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
		log.Println("Error Generating Forgot Password Email HTML: ", err)
		return err
	}

	sendErr := m.Send(
		UserName,
		UserEmail,
		"Forgot your Thunderdome password?",
		emailBody,
	)
	if sendErr != nil {
		log.Println("Error sending Forgot Password Email: ", sendErr)
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
		log.Println("Error Generating Reset Password Email HTML: ", err)
		return err
	}

	sendErr := m.Send(
		UserName,
		UserEmail,
		"Your Thunderdome password was successfully reset.",
		emailBody,
	)
	if sendErr != nil {
		log.Println("Error sending Reset Password Email: ", sendErr)
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
		log.Println("Error Generating Update Password Email HTML: ", err)
		return err
	}

	sendErr := m.Send(
		UserName,
		UserEmail,
		"Your Thunderdome password was successfully updated.",
		emailBody,
	)
	if sendErr != nil {
		log.Println("Error sending Update Password Email: ", sendErr)
		return sendErr
	}

	return nil
}
