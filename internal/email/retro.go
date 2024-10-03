package email

import (
	"fmt"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
	"github.com/matcornic/hermes/v2"
	"go.uber.org/zap"
)

// SendRetroOverview sends the retro overview (items, action items) email to attendees
func (s *Service) SendRetroOverview(retro *thunderdome.Retro, template *thunderdome.RetroTemplate, userName string, userEmail string) error {
	columnMap := make(map[string]string)
	var columnsList string
	for _, column := range template.Format.Columns {
		columnMap[column.Name] = fmt.Sprintf(`
## %s

`, column.Label)
	}
	var retroActionsList string
	for _, action := range retro.ActionItems {
		retroActionsList += formatRetroActionWithAssignee(action)
	}
	for _, item := range retro.Items {
		columnMap[item.Type] += formatRetroItemForMarkdownList(item.Content)
	}
	for _, column := range template.Format.Columns {
		columnsList += fmt.Sprintf(`
%s

`, columnMap[column.Name])
	}

	subject := fmt.Sprintf("Here is your %s Retro Overview", retro.Name)
	emailBody, err := s.generateBody(
		hermes.Body{
			Name: userName,
			Intros: []string{
				subject,
			},
			FreeMarkdown: `
## Action Items
` + hermes.Markdown(retroActionsList) + `
` + hermes.Markdown(columnsList) + `

`,
		},
	)
	if err != nil {
		s.Logger.Error("Error Generating Retro Overview Email HTML", zap.Error(err),
			zap.String("user_email", userEmail))

		return err
	}

	//s.Logger.Info("Sending Retro Overview Email", zap.String("email_body", emailBody))
	sendErr := s.send(
		userName,
		userEmail,
		subject,
		emailBody,
	)
	if sendErr != nil {
		s.Logger.Error("Error sending Retro Overview Email", zap.Error(sendErr),
			zap.String("user_email", userEmail))
		return sendErr
	}

	return nil
}
