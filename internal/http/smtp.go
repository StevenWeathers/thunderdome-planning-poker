package http

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

type smtpTestRequestBody struct {
	Email string `json:"email" validate:"required,email" example:"admin@example.com"`
	Name  string `json:"name" validate:"required" example:"Administrator"`
}

// handleGetSMTPConfig returns the current SMTP configuration (with sensitive data masked)
//
//	@Summary		Get SMTP Configuration
//	@Description	Returns the current SMTP configuration with sensitive data masked for security
//	@Tags			admin
//	@Produce		json
//	@Success		200	object	standardJsonResponse{data=object}
//	@Failure		401	object	standardJsonResponse{}
//	@Failure		403	object	standardJsonResponse{}
//	@Failure		500	object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/admin/smtp/config [get]
func (s *Service) handleGetSMTPConfig() http.HandlerFunc {
	return s.userOnly(s.adminOnly(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)

		s.Logger.Ctx(ctx).Info("Get SMTP config called",
			zap.String("session_user_id", sessionUserID))

		config := s.Email.GetSanitizedConfig()

		s.Success(w, r, http.StatusOK, config, nil)
	}))
}

// handleTestSMTPConnection tests the SMTP connection without sending an email
//
//	@Summary		Test SMTP Connection
//	@Description	Tests the SMTP connection configuration without sending an email
//	@Tags			admin
//	@Produce		json
//	@Success		200	object	standardJsonResponse{}
//	@Failure		401	object	standardJsonResponse{}
//	@Failure		403	object	standardJsonResponse{}
//	@Failure		500	object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/admin/smtp/test-connection [post]
func (s *Service) handleTestSMTPConnection() http.HandlerFunc {
	return s.userOnly(s.adminOnly(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)

		s.Logger.Ctx(ctx).Info("Test SMTP connection called",
			zap.String("session_user_id", sessionUserID))

		err := s.Email.TestConnection()
		if err != nil {
			s.Logger.Ctx(ctx).Error("SMTP connection test failed",
				zap.String("session_user_id", sessionUserID),
				zap.Error(err))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Logger.Ctx(ctx).Info("SMTP connection test successful",
			zap.String("session_user_id", sessionUserID))

		s.Success(w, r, http.StatusOK, "SMTP connection test successful", nil)
	}))
}

// handleSendTestEmail sends a test email to verify SMTP configuration
//
//	@Summary		Send Test Email
//	@Description	Sends a test email to the specified recipient to verify SMTP configuration
//	@Tags			admin
//	@Accept			json
//	@Produce		json
//	@Param			testEmail	body	smtpTestRequestBody	true	"Test email details"
//	@Success		200			object	standardJsonResponse{}
//	@Failure		400			object	standardJsonResponse{}
//	@Failure		401			object	standardJsonResponse{}
//	@Failure		403			object	standardJsonResponse{}
//	@Failure		500			object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/admin/smtp/test-email [post]
func (s *Service) handleSendTestEmail() http.HandlerFunc {
	return s.userOnly(s.adminOnly(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)
		var body smtpTestRequestBody

		jsonErr := json.NewDecoder(r.Body).Decode(&body)
		if jsonErr != nil {
			s.Logger.Ctx(ctx).Error("Failed to decode SMTP test request body",
				zap.String("session_user_id", sessionUserID),
				zap.Error(jsonErr))
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINTERNAL, "Invalid JSON request body: "+jsonErr.Error()))
			return
		}

		s.Logger.Ctx(ctx).Info("SMTP test request received",
			zap.String("session_user_id", sessionUserID),
			zap.String("email", body.Email),
			zap.String("name", body.Name),
			zap.Int("email_len", len(body.Email)),
			zap.Int("name_len", len(body.Name)))

		// Additional validation logging for debugging
		if body.Email == "" {
			s.Logger.Ctx(ctx).Error("SMTP test email field is empty")
		}
		if body.Name == "" {
			s.Logger.Ctx(ctx).Error("SMTP test name field is empty")
		}

		inputErr := validate.Struct(body)
		if inputErr != nil {
			s.Logger.Ctx(ctx).Error("SMTP test request validation failed",
				zap.String("session_user_id", sessionUserID),
				zap.String("email", body.Email),
				zap.String("name", body.Name),
				zap.String("validation_error", inputErr.Error()),
				zap.Error(inputErr))
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, "Validation failed: "+inputErr.Error()))
			return
		}

		s.Logger.Ctx(ctx).Info("Send test email called",
			zap.String("session_user_id", sessionUserID),
			zap.String("recipient_email", body.Email),
			zap.String("recipient_name", body.Name))

		err := s.Email.SendTestEmail(body.Email, body.Name)
		if err != nil {
			s.Logger.Ctx(ctx).Error("Failed to send test email",
				zap.String("session_user_id", sessionUserID),
				zap.String("recipient_email", body.Email),
				zap.Error(err))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Logger.Ctx(ctx).Info("Test email sent successfully",
			zap.String("session_user_id", sessionUserID),
			zap.String("recipient_email", body.Email))

		s.Success(w, r, http.StatusOK, "Test email sent successfully", nil)
	}))
}
