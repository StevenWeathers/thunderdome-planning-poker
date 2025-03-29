package http

import (
	"encoding/json"
	"io"
	"net/http"

	"go.uber.org/zap"
)

var ActiveAlerts []interface{}

type alertRequestBody struct {
	Name           string `json:"name" validate:"required"`
	Type           string `json:"type" enums:"ERROR, INFO, NEW, SUCCESS, WARNING" validate:"required,oneof=ERROR INFO NEW SUCCESS WARNING"`
	Content        string `json:"content" validate:"required"`
	Active         bool   `json:"active"`
	AllowDismiss   bool   `json:"allowDismiss"`
	RegisteredOnly bool   `json:"registeredOnly"`
}

// handleGetAlerts gets a list of alerts
//
//	@Summary		Get Alerts
//	@Description	get list of alerts (global notices)
//	@Tags			alert
//	@Produce		json
//	@Param			limit	query	int	false	"Max number of results to return"
//	@Param			offset	query	int	false	"Starting point to return rows from, should be multiplied by limit or 0"
//	@Success		200		object	standardJsonResponse{data=[]thunderdome.Alert}
//	@Failure		500		object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/alerts [get]
func (s *Service) handleGetAlerts() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID, _ := ctx.Value(contextKeyUserID).(*string)
		limit, offset := getLimitOffsetFromRequest(r)
		alerts, count, err := s.AlertDataSvc.AlertsList(ctx, limit, offset)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleGetAlerts error", zap.Error(err),
				zap.Int("limit", limit), zap.Int("offset", offset),
				zap.Stringp("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		meta := &pagination{
			Count:  count,
			Offset: offset,
			Limit:  limit,
		}

		s.Success(w, r, http.StatusOK, alerts, meta)
	}
}

// handleAlertCreate creates a new alert
//
//	@Summary		Create Alert
//	@Description	Creates an alert (global notice)
//	@Tags			alert
//	@Produce		json
//	@Param			alert	body	alertRequestBody								true	"new alert object"
//	@Success		200		object	standardJsonResponse{data=[]thunderdome.Alert}	"returns active alerts"
//	@Failure		500		object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/alerts [post]
func (s *Service) handleAlertCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)
		var alert = alertRequestBody{}
		body, bodyErr := io.ReadAll(r.Body)
		if bodyErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
			return
		}

		jsonErr := json.Unmarshal(body, &alert)
		if jsonErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		inputErr := validate.Struct(alert)
		if inputErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, inputErr.Error()))
			return
		}

		err := s.AlertDataSvc.AlertsCreate(ctx, alert.Name, alert.Type, alert.Content, alert.Active, alert.AllowDismiss, alert.RegisteredOnly)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleAlertCreate error", zap.Error(err),
				zap.String("alert_name", alert.Name), zap.String("alert_type", alert.Type),
				zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		ActiveAlerts = s.AlertDataSvc.GetActiveAlerts(ctx)

		s.Success(w, r, http.StatusOK, ActiveAlerts, nil)
	}
}

// handleAlertUpdate updates an alert
//
//	@Summary		Update Alert
//	@Description	Updates an Alert
//	@Tags			alert
//	@Produce		json
//	@Param			alertId	path	string											true	"the alert ID to update"
//	@Param			alert	body	alertRequestBody								true	"alert object to update"
//	@Success		200		object	standardJsonResponse{data=[]thunderdome.Alert}	"returns active alerts"
//	@Failure		500		object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/alerts/{alertId} [put]
func (s *Service) handleAlertUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)
		alertID := r.PathValue("alertId")
		idErr := validate.Var(alertID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		var alert = alertRequestBody{}
		body, bodyErr := io.ReadAll(r.Body)
		if bodyErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
			return
		}

		jsonErr := json.Unmarshal(body, &alert)
		if jsonErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		inputErr := validate.Struct(alert)
		if inputErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, inputErr.Error()))
			return
		}

		err := s.AlertDataSvc.AlertsUpdate(ctx, alertID, alert.Name, alert.Type, alert.Content, alert.Active, alert.AllowDismiss, alert.RegisteredOnly)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleAlertUpdate error", zap.Error(err), zap.String("alert_id", alertID),
				zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		ActiveAlerts = s.AlertDataSvc.GetActiveAlerts(ctx)

		s.Success(w, r, http.StatusOK, ActiveAlerts, nil)
	}
}

// handleAlertDelete handles deleting an alert
//
//	@Summary		Delete Alert
//	@Description	Deletes an Alert
//	@Tags			alert
//	@Produce		json
//	@Param			alertId	path	string											true	"the alert ID to delete"
//	@Success		200		object	standardJsonResponse{data=[]thunderdome.Alert}	"returns active alerts"
//	@Failure		500		object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/alerts/{alertId} [delete]
func (s *Service) handleAlertDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)
		alertID := r.PathValue("alertId")
		idErr := validate.Var(alertID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		err := s.AlertDataSvc.AlertDelete(ctx, alertID)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleAlertDelete error", zap.Error(err), zap.String("alert_id", alertID),
				zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		ActiveAlerts = s.AlertDataSvc.GetActiveAlerts(ctx)

		s.Success(w, r, http.StatusOK, ActiveAlerts, nil)
	}
}
