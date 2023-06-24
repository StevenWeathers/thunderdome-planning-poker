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
}
