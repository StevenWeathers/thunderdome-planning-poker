package email

import (
	"fmt"
	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
	"github.com/matcornic/hermes/v2"
	"go.uber.org/zap"
)

// SendRetroOverview sends the retro overview (items, action items) email to attendees
func (s *Service) SendRetroOverview(retro *thunderdome.Retro, UserName string, UserEmail string) error {
	var retroActionsList string
	var retroWorksList string
	var retroImproveList string
	var retroQuestionList string
	for _, action := range retro.ActionItems {
		retroActionsList += fmt.Sprintf("- %s\n", action.Content)
	}
	for _, item := range retro.Items {
		switch item.Type {
		case "worked":
			retroWorksList += fmt.Sprintf("- %s\n", item.Content)
		case "improve":
			retroImproveList += fmt.Sprintf("- %s\n", item.Content)
		case "question":
			retroQuestionList += fmt.Sprintf("- %s\n", item.Content)
		}
	}

	subject := fmt.Sprintf("Here is your %s Retro Overview", retro.Name)
	emailBody, err := s.generateBody(
		hermes.Body{
			Name: UserName,
			Intros: []string{
				subject,
			},
			FreeMarkdown: `
## Action Items
` + hermes.Markdown(retroActionsList) + `

## Works
` + hermes.Markdown(retroWorksList) + `

## Needs Improvement
` + hermes.Markdown(retroImproveList) + `

## Questions
` + hermes.Markdown(retroQuestionList) + `

`,
		},
	)
	if err != nil {
		s.Logger.Error("Error Generating Retro Overview Email HTML", zap.Error(err))

		return err
	}

	sendErr := s.send(
		UserName,
		UserEmail,
		subject,
		emailBody,
	)
	if sendErr != nil {
		s.Logger.Error("Error sending Retro Overview Email", zap.Error(sendErr))
		return sendErr
	}

	return nil
}
