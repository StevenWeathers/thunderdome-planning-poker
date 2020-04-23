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

var smtpServerConfig = smtpServer{}
var tlsConfig = &tls.Config{}
var smtpFrom = mail.Address{}
var smtpAuth smtp.Auth

// EmailConfig contains all the mailserver values
type EmailConfig struct {
	AppDomain    string
	SenderName   string
	smtpHost     string
	smtpPort     string
	smtpSecure   bool
	smtpIdentity string
	smtpUser     string
	smtpPass     string
	smtpSender   string
}

// Email contains all the methods to send application emails
type Email struct {
	config *EmailConfig
}

// NewEmail creates a new instance of Email
func NewEmail(AppDomain string) *Email {
	var m = &Email{
		// read environment variables and sets up mailserver configuration values
		config: &EmailConfig{
			AppDomain:    AppDomain,
			SenderName:   "Thunderdome",
			smtpHost:     GetEnv("SMTP_HOST", "localhost"),
			smtpPort:     GetEnv("SMTP_PORT", "25"),
			smtpSecure:   GetBoolEnv("SMTP_SECURE", true),
			smtpIdentity: GetEnv("SMTP_IDENTITY", ""),
			smtpUser:     GetEnv("SMTP_USER", ""),
			smtpPass:     GetEnv("SMTP_PASS", ""),
			smtpSender:   GetEnv("SMTP_SENDER", "no-reply@thunderdome.dev"),
		},
	}

	// smtp server configuration.
	smtpServerConfig = smtpServer{host: m.config.smtpHost, port: m.config.smtpPort}

	// smtp sender info
	smtpFrom = mail.Address{
		Name:    m.config.SenderName,
		Address: m.config.smtpSender,
	}

	// TLS config
	tlsConfig = &tls.Config{
		InsecureSkipVerify: !m.config.smtpSecure,
		ServerName:         m.config.smtpHost,
	}

	smtpAuth = smtp.PlainAuth(m.config.smtpIdentity, m.config.smtpUser, m.config.smtpPass, m.config.smtpHost)

	return m
}

// Generates an Email Body with hermes
func (m *Email) generateBody(Body hermes.Body) (emailBody string, generateErr error) {
	currentTime := time.Now()
	year := strconv.Itoa(currentTime.Year())
	hms := hermes.Hermes{
		Product: hermes.Product{
			Name:      "Thunderdome",
			Link:      "https://" + m.config.AppDomain + "/",
			Logo:      "https://" + m.config.AppDomain + "/img/thunderdome-email-logo.png",
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

// Send - utility function to send emails
func (m *Email) Send(WarriorName string, WarriorEmail string, Subject string, Body string) error {
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
	if m.config.smtpSecure == true {
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

// SendWelcome sends the welcome email to new registered user
func (m *Email) SendWelcome(WarriorName string, WarriorEmail string, VerifyID string) error {
	emailBody, err := m.generateBody(
		hermes.Body{
			Name: WarriorName,
			Intros: []string{
				"Welcome to the Thunderdome! Bring your own mouthguard.",
			},
			Actions: []hermes.Action{
				{
					Instructions: "Please validate your email, the following link will expire in 24 hours.",
					Button: hermes.Button{
						Color: "#22BC66",
						Text:  "Verify Account",
						Link:  "https://" + m.config.AppDomain + "/verify-account/" + VerifyID,
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

// SendForgotPassword Sends a Forgot Password reset email to warrior
func (m *Email) SendForgotPassword(WarriorName string, WarriorEmail string, ResetID string) error {
	emailBody, err := m.generateBody(
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
						Link: "https://" + m.config.AppDomain + "/reset-password/" + ResetID,
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

// SendPasswordReset Sends a Reset Password confirmation email to warrior
func (m *Email) SendPasswordReset(WarriorName string, WarriorEmail string) error {
	emailBody, err := m.generateBody(
		hermes.Body{
			Name: WarriorName,
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
		WarriorName,
		WarriorEmail,
		"Your Thunderdome password was successfully reset.",
		emailBody,
	)
	if sendErr != nil {
		log.Println("Error sending Reset Password Email: ", sendErr)
		return sendErr
	}

	return nil
}

// SendPasswordUpdate Sends an Update Password confirmation email to warrior
func (m *Email) SendPasswordUpdate(WarriorName string, WarriorEmail string) error {
	emailBody, err := m.generateBody(
		hermes.Body{
			Name: WarriorName,
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
		WarriorName,
		WarriorEmail,
		"Your Thunderdome password was successfully updated.",
		emailBody,
	)
	if sendErr != nil {
		log.Println("Error sending Update Password Email: ", sendErr)
		return sendErr
	}

	return nil
}
