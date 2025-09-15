// Package email provides email sending for Thunderdome
package email

import (
	"context"
	"crypto/tls"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"github.com/wneessen/go-mail"

	"github.com/matcornic/hermes/v2"
)

// Config contains all the mail server values
type Config struct {
	AppURL            string
	RepoURL           string
	SenderName        string
	SmtpHost          string
	SmtpPort          int
	SmtpSecure        bool
	SmtpUser          string
	SmtpPass          string
	SmtpSender        string
	SmtpEnabled       bool
	SmtpSkipTLSVerify bool
	SmtpAuth          string
}

// Service contains all the methods to send application emails
type Service struct {
	Config    *Config
	Logger    *otelzap.Logger
	tlsConfig *tls.Config
	authType  mail.SMTPAuthType
}

// New creates a new instance of Service
func New(config *Config, logger *otelzap.Logger) *Service {
	var s = &Service{
		// read environment variables and sets up mail server configuration values
		Config: config,
		Logger: logger,
	}

	s.authType = mail.SMTPAuthType(s.Config.SmtpAuth)
	s.tlsConfig = &tls.Config{
		InsecureSkipVerify: s.Config.SmtpSkipTLSVerify || !s.Config.SmtpSecure,
		ServerName:         s.Config.SmtpHost,
	}

	return s
}

// Generates an email Body with hermes
func (s *Service) generateBody(body hermes.Body) (emailBody string, generateErr error) {
	currentTime := time.Now()
	year := strconv.Itoa(currentTime.Year())
	hms := hermes.Hermes{
		Product: hermes.Product{
			Name:      "Thunderdome",
			Link:      s.Config.AppURL,
			Logo:      s.Config.AppURL + "img/thunderdome-email-logo.png",
			Copyright: "Copyright © " + year + " Thunderdome. All rights reserved.",
		},
	}

	email := hermes.Email{
		Body: body,
	}

	// Generate an HTML email with the provided contents (for modern clients)
	emailBody, err := hms.GenerateHTML(email)
	if err != nil {
		return "", fmt.Errorf("failed to generate email body html: %v", err)
	}

	return emailBody, nil
}

// send - utility function to send emails
func (s *Service) send(userName string, userEmail string, subject string, body string) error {
	var err error
	var c *mail.Client
	if !s.Config.SmtpEnabled {
		return nil
	}
	cleanUsername, err := removeAccents(userName)
	if err != nil {
		return fmt.Errorf("failed to clean username %s: %v", userName, err)
	}

	m := mail.NewMsg()
	if err = m.From(s.Config.SmtpSender); err != nil {
		return fmt.Errorf("failed to set From address %s error: %v", s.Config.SmtpSender, err)
	}
	if err = m.To(userEmail); err != nil {
		return fmt.Errorf("failed to set To address %s error: %v", userEmail, err)
	}

	m.Subject(subject)
	m.SetBodyString(mail.TypeTextHTML, body)
	if err = m.SetAddrHeader(mail.HeaderFrom, fmt.Sprintf("%s <%s>", s.Config.SenderName, s.Config.SmtpSender)); err != nil {
		return fmt.Errorf("failed to set FROM header: %v", err)
	}
	if err = m.SetAddrHeader(mail.HeaderTo, fmt.Sprintf("%s <%s>", cleanUsername, userEmail)); err != nil {
		return fmt.Errorf("failed to set TO header: %v", err)
	}

	if s.Config.SmtpSecure {
		c, err = mail.NewClient(s.Config.SmtpHost, mail.WithPort(s.Config.SmtpPort), mail.WithSMTPAuth(s.authType),
			mail.WithUsername(s.Config.SmtpUser), mail.WithPassword(s.Config.SmtpPass), mail.WithTLSConfig(s.tlsConfig))
	} else {
		c, err = mail.NewClient(s.Config.SmtpHost, mail.WithPort(s.Config.SmtpPort), mail.WithTLSConfig(s.tlsConfig),
			mail.WithTLSPolicy(mail.TLSOpportunistic))
	}
	if err != nil {
		return fmt.Errorf("failed to create mail client: %v", err)
	}

	if err = c.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send mail: %v", err)
	}

	return err
}

// TestConnection tests the SMTP connection without sending an email
func (s *Service) TestConnection() error {
	if !s.Config.SmtpEnabled {
		return fmt.Errorf("SMTP is not enabled")
	}

	var c *mail.Client
	var err error

	if s.Config.SmtpSecure {
		c, err = mail.NewClient(s.Config.SmtpHost, mail.WithPort(s.Config.SmtpPort), mail.WithSMTPAuth(s.authType),
			mail.WithUsername(s.Config.SmtpUser), mail.WithPassword(s.Config.SmtpPass), mail.WithTLSConfig(s.tlsConfig))
	} else {
		c, err = mail.NewClient(s.Config.SmtpHost, mail.WithPort(s.Config.SmtpPort), mail.WithTLSConfig(s.tlsConfig),
			mail.WithTLSPolicy(mail.TLSOpportunistic))
	}
	if err != nil {
		return fmt.Errorf("failed to create mail client: %v", err)
	}

	// Test connection by dialing
	return c.DialWithContext(context.Background())
}

// SendTestEmail sends a test email to verify SMTP configuration
func (s *Service) SendTestEmail(recipient string, recipientName string) error {
	if !s.Config.SmtpEnabled {
		return fmt.Errorf("SMTP is not enabled")
	}

	testSubject := "Thunderdome SMTP Test Email"
	testBody, err := s.generateBody(hermes.Body{
		Title: "SMTP Configuration Test",
		Intros: []string{
			"This is a test email from your Thunderdome Planning Poker application.",
			"If you received this email, your SMTP configuration is working correctly!",
		},
		Dictionary: []hermes.Entry{
			{
				Key:   "Test Date",
				Value: time.Now().Format("2006-01-02 15:04:05 MST"),
			},
			{
				Key:   "Server",
				Value: s.Config.SmtpHost + ":" + strconv.Itoa(s.Config.SmtpPort),
			},
			{
				Key:   "Sender",
				Value: s.Config.SmtpSender,
			},
		},
		Outros: []string{
			"This email was sent from the admin panel for testing purposes.",
		},
	})

	if err != nil {
		return fmt.Errorf("failed to generate test email body: %v", err)
	}

	return s.send(recipientName, recipient, testSubject, testBody)
}

// GetSanitizedConfig returns SMTP configuration with sensitive data masked
func (s *Service) GetSanitizedConfig() map[string]interface{} {
	config := map[string]interface{}{
		"enabled":       s.Config.SmtpEnabled,
		"host":          s.Config.SmtpHost,
		"port":          s.Config.SmtpPort,
		"secure":        s.Config.SmtpSecure,
		"sender":        s.Config.SmtpSender,
		"senderName":    s.Config.SenderName,
		"skipTLSVerify": s.Config.SmtpSkipTLSVerify,
		"authType":      s.Config.SmtpAuth,
	}

	// Mask sensitive information
	if s.Config.SmtpUser != "" {
		config["username"] = maskString(s.Config.SmtpUser)
	}
	if s.Config.SmtpPass != "" {
		config["hasPassword"] = true
	} else {
		config["hasPassword"] = false
	}

	return config
}

// maskString masks sensitive string information for display
func maskString(input string) string {
	if len(input) <= 4 {
		return strings.Repeat("*", len(input))
	}
	return input[:2] + strings.Repeat("*", len(input)-4) + input[len(input)-2:]
}
