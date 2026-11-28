package email

import (
	"fmt"
	"time"

	"github.com/matcornic/hermes/v2"
	"go.uber.org/zap"
)

// SendHealthCheckInvite sends a team member the link to fill out a health check
func (s *Service) SendHealthCheckInvite(userName string, userEmail string, teamName string, healthCheckName string, deadline time.Time, teamID string, healthCheckID string) error {
	subject := fmt.Sprintf("%s health check for team %s", healthCheckName, teamName)
	emailBody, err := s.generateBody(
		hermes.Body{
			Name: userName,
			Intros: []string{
				subject,
			},
			Actions: []hermes.Action{
				{
					// the recipient's timezone isn't known, so the deadline is given in UTC
					Instructions: fmt.Sprintf(
						"Please fill out the health check before %s. Responses are anonymous.",
						deadline.UTC().Format("Monday, January 2, 2006 15:04 MST")),
					Button: hermes.Button{
						Color: "#22BC66",
						Text:  "Fill Out Health Check",
						Link:  fmt.Sprintf("%steam/%s/health-checks/%s", s.Config.AppURL, teamID, healthCheckID),
					},
				},
			},
		},
	)
	if err != nil {
		s.Logger.Error("Error Generating Health Check Invite Email HTML", zap.Error(err),
			zap.String("user_email", userEmail))

		return err
	}

	if err := s.send(userName, userEmail, subject, emailBody); err != nil {
		s.Logger.Error("Error sending Health Check Invite Email", zap.Error(err),
			zap.String("user_email", userEmail))

		return err
	}

	return nil
}
