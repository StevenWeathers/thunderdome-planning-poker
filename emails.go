package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/mail"
	"net/smtp"
	"strconv"
	"time"

	"github.com/matcornic/hermes/v2"
)

// smtpServer data to smtp server
type smtpServer struct {
	host string
	port string
}

// Address URI to smtp server
func (s *smtpServer) Address() string {
	return s.host + ":" + s.port
}

var smtpHost string
var smtpPort string
var smtpSecure bool
var smtpIdentity string
var smtpUser string
var smtpPass string
var smtpSender string
var smtpServerConfig = smtpServer{}
var tlsConfig = &tls.Config{}
var smtpFrom = mail.Address{}
var smtpAuth smtp.Auth

// GetMailserverConfig reads environment variables and sets up mailserver configuration values
func GetMailserverConfig() {
	smtpHost = GetEnv("SMTP_HOST", "localhost")
	smtpPort = GetEnv("SMTP_PORT", "25")
	smtpSecure = GetBoolEnv("SMTP_SECURE", true)
	smtpIdentity = GetEnv("SMTP_IDENTITY", "")
	smtpUser = GetEnv("SMTP_USER", "")
	smtpPass = GetEnv("SMTP_PASS", "")
	smtpSender = GetEnv("SMTP_SENDER", "no-reply@thunderdome.dev")

	// smtp server configuration.
	smtpServerConfig = smtpServer{host: smtpHost, port: smtpPort}

	// smtp sender info
	smtpFrom = mail.Address{
		Name:    "Thunderdome",
		Address: smtpSender,
	}

	// TLS config
	tlsConfig = &tls.Config{
		InsecureSkipVerify: !smtpSecure,
		ServerName:         smtpHost,
	}

	smtpAuth = smtp.PlainAuth(smtpIdentity, smtpUser, smtpPass, smtpHost)
}

// Generates an Email Body with hermes
func generateEmailBody(Body hermes.Body) (emailBody string, generateErr error) {
	currentTime := time.Now()
	year := strconv.Itoa(currentTime.Year())
	hms := hermes.Hermes{
		Product: hermes.Product{
			Name:      "Thunderdome",
			Link:      "https://thunderdome.dev/",
			Logo:      "https://thunderdome.dev/img/thunderdome-email-logo.png",
			Copyright: "Copyright © " + year + " Thunderdome. All rights reserved.",
		},
	}

	email := hermes.Email{
		Body: Body,
	}

	// Generate an HTML email with the provided contents (for modern clients)
	emailBody, err := hms.GenerateHTML(email)
	if err != nil {
		return "", err
	}

	return emailBody, nil
}

// utility function to send emails
func sendEmail(WarriorName string, WarriorEmail string, Subject string, Body string) error {
	to := mail.Address{
		Name:    WarriorName,
		Address: WarriorEmail,
	}

	// Setup headers
	headers := make(map[string]string)
	headers["From"] = smtpFrom.String()
	headers["To"] = to.String()
	headers["Subject"] = Subject
	headers["MIME-version"] = "1.0"
	headers["Content-Type"] = "text/html"

	// Setup message
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + Body

	c, err := smtp.Dial(smtpServerConfig.Address())
	if err != nil {
		log.Println("Error dialing SMTP: ", err)
		return err
	}

	c.StartTLS(tlsConfig)

	// Auth
	if smtpSecure == true {
		if err = c.Auth(smtpAuth); err != nil {
			log.Println("Error authenticating SMTP: ", err)
			return err
		}
	}

	// To && From
	if err = c.Mail(smtpFrom.Address); err != nil {
		log.Println("Error setting SMTP from: ", err)
		return err
	}

	if err = c.Rcpt(to.Address); err != nil {
		log.Println("Error setting SMTP to: ", err)
		return err
	}

	// Data
	w, err := c.Data()
	if err != nil {
		log.Println("Error setting SMTP data: ", err)
		return err
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		log.Println("Error sending email: ", err)
		return err
	}

	err = w.Close()
	if err != nil {
		log.Println("Error closing SMTP: ", err)
		return err
	}

	c.Quit()

	return nil
}

// SendWelcomeEmail sends the welcome email to new registered user
func SendWelcomeEmail(WarriorName string, WarriorEmail string) error {
	emailBody, err := generateEmailBody(
		hermes.Body{
			Name: WarriorName,
			Intros: []string{
				"Welcome to the Thunderdome! Bring your own mouthguard.",
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
		log.Println("Error Generating Welcome Email HTML: ", err)
		return err
	}

	sendErr := sendEmail(
		WarriorName,
		WarriorEmail,
		"Welcome to the Thunderdome!",
		emailBody,
	)
	if sendErr != nil {
		log.Println("Error sending Welcome Email: ", sendErr)
		return sendErr
	}

	return nil
}

// SendForgotPasswordEmail Sends a Forgot Password reset email to warrior
func SendForgotPasswordEmail(WarriorName string, WarriorEmail string, ResetID string) error {
	emailBody, err := generateEmailBody(
		hermes.Body{
			Name: WarriorName,
			Intros: []string{
				"It seems you've forgot your Thunderdome password.",
			},
			Actions: []hermes.Action{
				{
					Instructions: "Reset your password now, the following link will expire within an hour of the original request.",
					Button: hermes.Button{
						Text: "Reset Password",
						Link: "https://thunderdome.dev/reset-password/" + ResetID,
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

	sendErr := sendEmail(
		WarriorName,
		WarriorEmail,
		"Forgot your Thunderdome password?",
		emailBody,
	)
	if sendErr != nil {
		log.Println("Error sending Forgot Password Email: ", sendErr)
		return sendErr
	}

	return nil
}

// SendPasswordResetEmail Sends a Reset Password confirmation email to warrior
func SendPasswordResetEmail(WarriorName string, WarriorEmail string) error {
	emailBody, err := generateEmailBody(
		hermes.Body{
			Name: WarriorName,
			Intros: []string{
				"Your Thunderdome password was succesfully reset.",
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
		log.Println("Error Generating Forgot Password Email HTML: ", err)
		return err
	}

	sendErr := sendEmail(
		WarriorName,
		WarriorEmail,
		"Your Thunderdome password was succesfully reset.",
		emailBody,
	)
	if sendErr != nil {
		log.Println("Error sending Forgot Password Email: ", sendErr)
		return sendErr
	}

	return nil
}
