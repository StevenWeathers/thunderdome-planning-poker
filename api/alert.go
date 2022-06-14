package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

var ActiveAlerts []interface{}

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
		Alerts, Count, err := a.db.AlertsList(Limit, Offset)
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
// @Param name body string true "Name of the alert"
// @Param type body string true "Type of alert" Enums(ERROR, INFO, NEW, SUCCESS, WARNING)
// @Param content body string true "Alert content"
// @Param active body boolean true "Whether alert should be displayed or not"
// @Param allowDismiss body boolean true "Whether or not to allow users to dismiss the alert"
// @Param registeredOnly body boolean true "Whether or not to only show to users with an active session"
// @Success 200 object standardJsonResponse{data=[]model.Alert} "returns active alerts"
// @Failure 500 object standardJsonResponse{}
// @Security ApiKeyAuth
// @Router /alerts [post]
func (a *api) handleAlertCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		keyVal := getJSONRequestBody(r, w)

		Name := keyVal["name"].(string)
		Type := keyVal["type"].(string)
		Content := keyVal["content"].(string)
		Active := keyVal["active"].(bool)
		AllowDismiss := keyVal["allowDismiss"].(bool)
		RegisteredOnly := keyVal["registeredOnly"].(bool)

		err := a.db.AlertsCreate(Name, Type, Content, Active, AllowDismiss, RegisteredOnly)
		if err != nil {
			a.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		ActiveAlerts = a.db.GetActiveAlerts()

		a.Success(w, r, http.StatusOK, ActiveAlerts, nil)
	}
}

// handleAlertUpdate updates an alert
// @Summary Update Alert
// @Description Updates an Alert
// @Tags alert
// @Produce  json
// @Param alertId path string true "the alert ID to update"
// @Param name body string true "Name of the alert"
// @Param type body string true "Type of alert" Enums(ERROR, INFO, NEW, SUCCESS, WARNING)
// @Param content body string true "Alert content"
// @Param active body boolean true "Whether alert should be displayed or not"
// @Param allowDismiss body boolean true "Whether or not to allow users to dismiss the alert"
// @Param registeredOnly body boolean true "Whether or not to only show to users with an active session"
// @Success 200 object standardJsonResponse{data=[]model.Alert} "returns active alerts"
// @Failure 500 object standardJsonResponse{}
// @Security ApiKeyAuth
// @Router /alerts/{alertId} [put]
func (a *api) handleAlertUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		keyVal := getJSONRequestBody(r, w)
		vars := mux.Vars(r)

		ID := vars["alertId"]
		Name := keyVal["name"].(string)
		Type := keyVal["type"].(string)
		Content := keyVal["content"].(string)
		Active := keyVal["active"].(bool)
		AllowDismiss := keyVal["allowDismiss"].(bool)
		RegisteredOnly := keyVal["registeredOnly"].(bool)

		err := a.db.AlertsUpdate(ID, Name, Type, Content, Active, AllowDismiss, RegisteredOnly)
		if err != nil {
			a.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		ActiveAlerts = a.db.GetActiveAlerts()

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

		err := a.db.AlertDelete(AlertID)
		if err != nil {
			a.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		ActiveAlerts = a.db.GetActiveAlerts()

		a.Success(w, r, http.StatusOK, ActiveAlerts, nil)
	}
}
