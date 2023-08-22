// Package email provides email sending for Thunderdome
package email

import (
	"crypto/tls"
	"fmt"
	"net/mail"
	"net/smtp"
	"strconv"
	"time"

	"github.com/uptrace/opentelemetry-go-extra/otelzap"

	"github.com/matcornic/hermes/v2"
	"go.uber.org/zap"
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

// Config contains all the mail server values
type Config struct {
	AppURL            string
	SenderName        string
	SmtpHost          string
	SmtpPort          string
	SmtpSecure        bool
	SmtpIdentity      string
	SmtpUser          string
	SmtpPass          string
	SmtpSender        string
	SmtpEnabled       bool
	SmtpSkipTLSVerify bool
}

// Service contains all the methods to send application emails
type Service struct {
	Config *Config
	Logger *otelzap.Logger
}

// New creates a new instance of Service
func New(config *Config, logger *otelzap.Logger) *Service {
	var s = &Service{
		// read environment variables and sets up mailserver configuration values
		Config: config,
		Logger: logger,
	}

	// smtp server configuration.
	smtpServerConfig = smtpServer{host: s.Config.SmtpHost, port: s.Config.SmtpPort}

	// smtp sender info
	smtpFrom = mail.Address{
		Name:    s.Config.SenderName,
		Address: s.Config.SmtpSender,
	}

	// TLS config
	tlsConfig = &tls.Config{
		InsecureSkipVerify: s.Config.SmtpSkipTLSVerify || !s.Config.SmtpSecure,
		ServerName:         s.Config.SmtpHost,
	}

	smtpAuth = smtp.PlainAuth(s.Config.SmtpIdentity, s.Config.SmtpUser, s.Config.SmtpPass, s.Config.SmtpHost)

	return s
}

// Generates an email Body with hermes
func (s *Service) generateBody(Body hermes.Body) (emailBody string, generateErr error) {
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
		Body: Body,
	}

	// Generate an HTML email with the provided contents (for modern clients)
	emailBody, err := hms.GenerateHTML(email)
	if err != nil {
		return "", err
	}

	return emailBody, nil
}

// send - utility function to send emails
func (s *Service) send(UserName string, UserEmail string, Subject string, Body string) error {
	if !s.Config.SmtpEnabled {
		return nil
	}

	to := mail.Address{
		Name:    UserName,
		Address: UserEmail,
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
		s.Logger.Error("Error dialing SMTP", zap.Error(err))
		return err
	}

	tlsErr := c.StartTLS(tlsConfig)
	if tlsErr != nil {
		s.Logger.Error("Error starting TLS", zap.Error(tlsErr))
	}

	// Auth
	if s.Config.SmtpSecure {
		if err = c.Auth(smtpAuth); err != nil {
			s.Logger.Error("Error authenticating SMTP", zap.Error(err))
			return err
		}
	}

	// To && From
	if err = c.Mail(smtpFrom.Address); err != nil {
		s.Logger.Error("Error setting SMTP from", zap.Error(err))
		return err
	}

	if err = c.Rcpt(to.Address); err != nil {
		s.Logger.Error("Error setting SMTP to", zap.Error(err))
		return err
	}

	// Data
	w, err := c.Data()
	if err != nil {
		s.Logger.Error("Error setting SMTP data", zap.Error(err))
		return err
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		s.Logger.Error("Error sending email", zap.Error(err))
		return err
	}

	err = w.Close()
	if err != nil {
		s.Logger.Error("Error closing SMTP", zap.Error(err))
		return err
	}

	quitErr := c.Quit()
	if quitErr != nil {
		s.Logger.Error("Error quitting smtp server connection", zap.Error(quitErr))
	}

	return nil
}
