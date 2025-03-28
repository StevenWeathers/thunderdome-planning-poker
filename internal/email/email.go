// Package email provides email sending for Thunderdome
package email

import (
	"crypto/tls"
	"fmt"
	"strconv"
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
