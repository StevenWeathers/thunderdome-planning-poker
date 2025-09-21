package email

import (
	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
	"github.com/matcornic/hermes/v2"
	"go.uber.org/zap"
)

// SendNewTicketToAdmins sends the new support ticket email to the admins
func (s *Service) SendNewTicketToAdmins(adminUser thunderdome.User, ticketID string) error {
	emailBody, err := s.generateBody(
		hermes.Body{
			Name: adminUser.Name,
			Intros: []string{
				"There is a New Support Ticket requiring your attention.",
			},
			Actions: []hermes.Action{
				{
					Instructions: "Please take a moment to review the ticket details.",
					Button: hermes.Button{
						Color: "#22BC66",
						Text:  "View Ticket",
						Link:  s.Config.AppURL + "/admin/support-tickets/" + ticketID,
					},
				},
			},
		},
	)
	if err != nil {
		s.Logger.Error("Error Generating Admin Notification Support Ticket Email HTML", zap.Error(err),
			zap.String("user_email", adminUser.Email))

		return err
	}

	sendErr := s.send(
		adminUser.Name,
		adminUser.Email,
		"There is a New Support Ticket requiring your attention.",
		emailBody,
	)
	if sendErr != nil {
		s.Logger.Error("Error sending Admin Notification Support Ticket Email", zap.Error(sendErr),
			zap.String("user_email", adminUser.Email))
		return sendErr
	}

	return nil
}
