package email

import (
	"fmt"

	"github.com/matcornic/hermes/v2"
	"go.uber.org/zap"
)

// SendUserSubscriptionActive sends an email to the user that their subscription is now active
func (s *Service) SendUserSubscriptionActive(userName string, userEmail string, subscriptionType string) error {
	emailBody, err := s.generateBody(
		hermes.Body{
			Name: userName,
			Intros: []string{
				fmt.Sprintf("Your Thunderdome %s subscription is now active!", subscriptionType),
			},
			Actions: []hermes.Action{
				{
					Instructions: "Need help, or have questions? Visit our Github page",
					Button: hermes.Button{
						Text: "Github Repo",
						Link: s.Config.RepoURL,
					},
				},
			},
		},
	)
	if err != nil {
		s.Logger.Error("Error Generating Subscription Active Email HTML", zap.Error(err),
			zap.String("user_email", userEmail))

		return err
	}

	sendErr := s.send(
		userName,
		userEmail,
		"Your Thunderdome subscription is now active",
		emailBody,
	)
	if sendErr != nil {
		s.Logger.Error("Error sending  Subscription Active Email", zap.Error(sendErr),
			zap.String("user_email", userEmail))
		return sendErr
	}

	return nil
}

// SendUserSubscriptionDeactivated sends an email to the user that their subscription is now deactivated
func (s *Service) SendUserSubscriptionDeactivated(userName string, userEmail string, subscriptionType string) error {
	emailBody, err := s.generateBody(
		hermes.Body{
			Name: userName,
			Intros: []string{
				fmt.Sprintf("Your Thunderdome %s subscription is now deactivated, sorry to see you go.", subscriptionType),
			},
			Actions: []hermes.Action{
				{
					Instructions: "Need help, or have questions? Visit our Github page",
					Button: hermes.Button{
						Text: "Github Repo",
						Link: s.Config.RepoURL,
					},
				},
			},
		},
	)
	if err != nil {
		s.Logger.Error("Error Generating Subscription Deactivated Email HTML", zap.Error(err),
			zap.String("user_email", userEmail))

		return err
	}

	sendErr := s.send(
		userName,
		userEmail,
		"Your Thunderdome subscription is now deactivated",
		emailBody,
	)
	if sendErr != nil {
		s.Logger.Error("Error sending  Subscription Deactivated Email", zap.Error(sendErr),
			zap.String("user_email", userEmail))
		return sendErr
	}

	return nil
}
