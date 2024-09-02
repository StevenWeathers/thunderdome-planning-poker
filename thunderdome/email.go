package thunderdome

type EmailService interface {
	SendWelcome(UserName string, UserEmail string, VerifyID string) error
	SendEmailVerification(UserName string, UserEmail string, VerifyID string) error
	SendForgotPassword(UserName string, UserEmail string, ResetID string) error
	SendPasswordReset(UserName string, UserEmail string) error
	SendPasswordUpdate(UserName string, UserEmail string) error
	SendDeleteConfirmation(UserName string, UserEmail string) error
	SendEmailUpdate(UserName string, UserEmail string) error
	SendMergedUpdate(UserName string, UserEmail string) error
	SendTeamInvite(TeamName string, UserEmail string, InviteID string) error
	SendOrganizationInvite(OrganizationName string, UserEmail string, InviteID string) error
	SendDepartmentInvite(OrganizationName string, DepartmentName string, UserEmail string, InviteID string) error
	SendUserSubscriptionActive(UserName string, UserEmail string, SubscriptionType string) error
	SendUserSubscriptionDeactivated(UserName string, UserEmail string, SubscriptionType string) error
	// SendRetroOverview sends the retro overview (items, action items) email to attendees
	SendRetroOverview(retro *Retro, UserName string, UserEmail string) error
}
