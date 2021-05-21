package main

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// handleGetAlerts gets a list of alerts
func (s *server) handleGetAlerts() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		Limit, _ := strconv.Atoi(vars["limit"])
		Offset, _ := strconv.Atoi(vars["offset"])

		Alerts := s.database.AlertsList(Limit, Offset)

		s.respondWithJSON(w, http.StatusOK, Alerts)
	}
}

// handleAlertCreate creates a new alert
func (s *server) handleAlertCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		keyVal := s.getJSONRequestBody(r, w)

		Name := keyVal["name"].(string)
		Type := keyVal["type"].(string)
		Content := keyVal["content"].(string)
		Active := keyVal["active"].(bool)
		AllowDismiss := keyVal["allowDismiss"].(bool)
		RegisteredOnly := keyVal["registeredOnly"].(bool)

		err := s.database.AlertsCreate(Name, Type, Content, Active, AllowDismiss, RegisteredOnly)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		ActiveAlerts = s.database.GetActiveAlerts()

		s.respondWithJSON(w, http.StatusOK, ActiveAlerts)
	}
}

// handleAlertUpdate updates an alert
func (s *server) handleAlertUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		keyVal := s.getJSONRequestBody(r, w)
		vars := mux.Vars(r)

		ID := vars["id"]
		Name := keyVal["name"].(string)
		Type := keyVal["type"].(string)
		Content := keyVal["content"].(string)
		Active := keyVal["active"].(bool)
		AllowDismiss := keyVal["allowDismiss"].(bool)
		RegisteredOnly := keyVal["registeredOnly"].(bool)

		err := s.database.AlertsUpdate(ID, Name, Type, Content, Active, AllowDismiss, RegisteredOnly)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		ActiveAlerts = s.database.GetActiveAlerts()

		s.respondWithJSON(w, http.StatusOK, ActiveAlerts)
	}
}

// handleAlertDelete handles deleting an alert
func (s *server) handleAlertDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		keyVal := s.getJSONRequestBody(r, w)
		AlertID := keyVal["id"].(string)

		err := s.database.AlertDelete(AlertID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		ActiveAlerts = s.database.GetActiveAlerts()

		s.respondWithJSON(w, http.StatusOK, ActiveAlerts)
	}
}
