package thunderdome

type EmailService interface {
	SendWelcome(userName string, userEmail string, verifyID string) error
	SendEmailVerification(userName string, userEmail string, verifyID string) error
	SendForgotPassword(userName string, userEmail string, resetID string) error
	SendPasswordReset(userName string, userEmail string) error
	SendPasswordUpdate(userName string, userEmail string) error
	SendDeleteConfirmation(userName string, userEmail string) error
	SendEmailUpdate(userName string, userEmail string) error
	SendMergedUpdate(userName string, userEmail string) error
	SendTeamInvite(TeamName string, userEmail string, inviteID string) error
	SendOrganizationInvite(organizationName string, userEmail string, inviteID string) error
	SendDepartmentInvite(organizationName string, departmentName string, userEmail string, inviteID string) error
	SendUserSubscriptionActive(userName string, userEmail string, subscriptionType string) error
	SendUserSubscriptionDeactivated(userName string, userEmail string, subscriptionType string) error
	// SendRetroOverview sends the retro overview (items, action items) email to attendees
	SendRetroOverview(retro *Retro, template *RetroTemplate, userName string, userEmail string) error
}
