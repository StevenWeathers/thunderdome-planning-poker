package api

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gorilla/mux"
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
// @Summary Get Alerts
// @Description get list of alerts (global notices)
// @Tags alert
// @Produce  json
// @Param limit query int false "Max number of results to return"
// @Param offset query int false "Starting point to return rows from, should be multiplied by limit or 0"
// @Success 200 object standardJsonResponse{data=[]model.Alert}
// @Failure 500 object standardJsonResponse{}
// @Security ApiKeyAuth
// @Router /alerts [get]
func (a *api) handleGetAlerts() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		Limit, Offset := getLimitOffsetFromRequest(r)
		Alerts, Count, err := a.db.AlertsList(r.Context(), Limit, Offset)
		if err != nil {
			a.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		Meta := &pagination{
			Count:  Count,
			Offset: Offset,
			Limit:  Limit,
		}

		a.Success(w, r, http.StatusOK, Alerts, Meta)
	}
}

// handleAlertCreate creates a new alert
// @Summary Create Alert
// @Description Creates an alert (global notice)
// @Tags alert
// @Produce  json
// @Param alert body alertRequestBody true "new alert object"
// @Success 200 object standardJsonResponse{data=[]model.Alert} "returns active alerts"
// @Failure 500 object standardJsonResponse{}
// @Security ApiKeyAuth
// @Router /alerts [post]
func (a *api) handleAlertCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var alert = alertRequestBody{}
		body, bodyErr := io.ReadAll(r.Body)
		if bodyErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
			return
		}

		jsonErr := json.Unmarshal(body, &alert)
		if jsonErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		inputErr := validate.Struct(alert)
		if inputErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, inputErr.Error()))
			return
		}

		err := a.db.AlertsCreate(r.Context(), alert.Name, alert.Type, alert.Content, alert.Active, alert.AllowDismiss, alert.RegisteredOnly)
		if err != nil {
			a.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		ActiveAlerts = a.db.GetActiveAlerts(r.Context())

		a.Success(w, r, http.StatusOK, ActiveAlerts, nil)
	}
}

// handleAlertUpdate updates an alert
// @Summary Update Alert
// @Description Updates an Alert
// @Tags alert
// @Produce  json
// @Param alertId path string true "the alert ID to update"
// @Param alert body alertRequestBody true "alert object to update"
// @Success 200 object standardJsonResponse{data=[]model.Alert} "returns active alerts"
// @Failure 500 object standardJsonResponse{}
// @Security ApiKeyAuth
// @Router /alerts/{alertId} [put]
func (a *api) handleAlertUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		ID := vars["alertId"]
		idErr := validate.Var(ID, "required,uuid")
		if idErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		var alert = alertRequestBody{}
		body, bodyErr := io.ReadAll(r.Body)
		if bodyErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
			return
		}

		jsonErr := json.Unmarshal(body, &alert)
		if jsonErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		inputErr := validate.Struct(alert)
		if inputErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, inputErr.Error()))
			return
		}

		err := a.db.AlertsUpdate(r.Context(), ID, alert.Name, alert.Type, alert.Content, alert.Active, alert.AllowDismiss, alert.RegisteredOnly)
		if err != nil {
			a.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		ActiveAlerts = a.db.GetActiveAlerts(r.Context())

		a.Success(w, r, http.StatusOK, ActiveAlerts, nil)
	}
}

// handleAlertDelete handles deleting an alert
// @Summary Delete Alert
// @Description Deletes an Alert
// @Tags alert
// @Produce  json
// @Param alertId path string true "the alert ID to delete"
// @Success 200 object standardJsonResponse{data=[]model.Alert} "returns active alerts"
// @Failure 500 object standardJsonResponse{}
// @Security ApiKeyAuth
// @Router /alerts/{alertId} [delete]
func (a *api) handleAlertDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		AlertID := vars["alertId"]
		idErr := validate.Var(AlertID, "required,uuid")
		if idErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		err := a.db.AlertDelete(r.Context(), AlertID)
		if err != nil {
			a.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		ActiveAlerts = a.db.GetActiveAlerts(r.Context())

		a.Success(w, r, http.StatusOK, ActiveAlerts, nil)
	}
}
