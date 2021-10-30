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
// @Param limit query int true "Max number of results to return"
// @Param offset query int true "Starting point to return rows from, should be multiplied by limit or 0"
// @Success 200
// @Router /alerts [get]
func (a *api) handleGetAlerts() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		Limit, Offset := a.getLimitOffsetFromRequest(r, w)

		Alerts := a.db.AlertsList(Limit, Offset)

		a.respondWithJSON(w, http.StatusOK, Alerts)
	}
}

// handleAlertCreate creates a new alert
// @Summary Create Alert
// @Description Creates an alert (global notice)
// @Tags alert
// @Produce  json
// @Param name body string false "Name of the alert"
// @Param type body string false "Type of alert" Enums(ERROR, INFO, NEW, SUCCESS, WARNING)
// @Param content body string false "Alert content"
// @Param active body boolean false "Whether alert should be displayed or not"
// @Param allowDismiss body boolean false "Whether or not to allow users to dismiss the alert"
// @Param registeredOnly body boolean false "Whether or not to only show to users with an active session"
// @Success 200
// @Router /alerts [post]
func (a *api) handleAlertCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		keyVal := a.getJSONRequestBody(r, w)

		Name := keyVal["name"].(string)
		Type := keyVal["type"].(string)
		Content := keyVal["content"].(string)
		Active := keyVal["active"].(bool)
		AllowDismiss := keyVal["allowDismiss"].(bool)
		RegisteredOnly := keyVal["registeredOnly"].(bool)

		err := a.db.AlertsCreate(Name, Type, Content, Active, AllowDismiss, RegisteredOnly)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		ActiveAlerts = a.db.GetActiveAlerts()

		a.respondWithJSON(w, http.StatusOK, ActiveAlerts)
	}
}

// handleAlertUpdate updates an alert
// @Summary Update Alert
// @Description Updates an Alert
// @Tags alert
// @Produce  json
// @Success 200
// @Router /alerts/{id} [put]
func (a *api) handleAlertUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		keyVal := a.getJSONRequestBody(r, w)
		vars := mux.Vars(r)

		ID := vars["id"]
		Name := keyVal["name"].(string)
		Type := keyVal["type"].(string)
		Content := keyVal["content"].(string)
		Active := keyVal["active"].(bool)
		AllowDismiss := keyVal["allowDismiss"].(bool)
		RegisteredOnly := keyVal["registeredOnly"].(bool)

		err := a.db.AlertsUpdate(ID, Name, Type, Content, Active, AllowDismiss, RegisteredOnly)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		ActiveAlerts = a.db.GetActiveAlerts()

		a.respondWithJSON(w, http.StatusOK, ActiveAlerts)
	}
}

// handleAlertDelete handles deleting an alert
// @Summary Delete Alert
// @Description Deletes an Alert
// @Tags alert
// @Produce  json
// @Param id path int false "the alert ID to delete"
// @Success 200
// @Router /alerts/{id} [delete]
func (a *api) handleAlertDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		AlertID := vars["id"]

		err := a.db.AlertDelete(AlertID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		ActiveAlerts = a.db.GetActiveAlerts()

		a.respondWithJSON(w, http.StatusOK, ActiveAlerts)
	}
}
