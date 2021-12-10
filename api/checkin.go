package api

import (
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

// handleGetTeamCheckins gets a list of team checkins
// @Summary Get Team Checkins
// @Description Get a list of team checkins
// @Tags team
// @Produce  json
// @Param teamId path string true "the team ID"
// @Param date query string false "the date in YYYY-MM-DD format"
// @Param tz query string false "the timezone name e.g. America/New_York"
// @Success 200 object standardJsonResponse{data=[]model.TeamCheckin}
// @Security ApiKeyAuth
// @Router /teams/{teamId}/checkins [get]
func (a *api) handleGetTeamCheckins() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		TeamID := vars["teamId"]
		query := r.URL.Query()
		date := query.Get("date")
		tz := query.Get("tz")

		if date == "" {
			date = time.Now().Format("1988-01-02")
		}

		if tz == "" {
			tz = "America/New_York"
		}

		Checkins, err := a.db.CheckinList(TeamID, date, tz)
		if err != nil {
			Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		Success(w, r, http.StatusOK, Checkins, nil)
	}
}

// handleCreateCheckin handles creating a team user checkin
// @Summary Create Team Checkin
// @Description Creates a team user checkin
// @Tags team
// @Produce  json
// @Success 200 object standardJsonResponse{}
// @Success 403 object standardJsonResponse{}
// @Success 500 object standardJsonResponse{}
// @Security ApiKeyAuth
// @Router /teams/{teamId}/checkins [post]
func (a *api) handleCreateCheckin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		TeamId := vars["teamId"]
		UserId := r.Context().Value(contextKeyUserID).(string)

		keyVal := getJSONRequestBody(r, w)
		Yesterday := keyVal["yesterday"].(string)
		Today := keyVal["today"].(string)
		Blockers := keyVal["blockers"].(string)
		Discuss := keyVal["discuss"].(string)
		GoalsMet := keyVal["goalsMet"].(bool)

		err := a.db.CheckinCreate(TeamId, UserId, Yesterday, Today, Blockers, Discuss, GoalsMet)
		if err != nil {
			Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		Success(w, r, http.StatusOK, nil, nil)
	}
}
