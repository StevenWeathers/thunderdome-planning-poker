package api

import (
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

// handleCheckinsGet gets a list of team checkins
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
func (a *api) handleCheckinsGet() http.HandlerFunc {
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

// handleCheckinCreate handles creating a team user checkin
// @Summary Create Team Checkin
// @Description Creates a team user checkin
// @Tags team
// @Produce  json
// @Success 200 object standardJsonResponse{}
// @Success 403 object standardJsonResponse{}
// @Success 500 object standardJsonResponse{}
// @Security ApiKeyAuth
// @Router /teams/{teamId}/checkins [post]
func (a *api) handleCheckinCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		TeamId := vars["teamId"]

		keyVal := getJSONRequestBody(r, w)
		UserId := keyVal["userId"].(string)
		Yesterday := keyVal["yesterday"].(string)
		Today := keyVal["today"].(string)
		Blockers := keyVal["blockers"].(string)
		Discuss := keyVal["discuss"].(string)
		GoalsMet := keyVal["goalsMet"].(bool)

		err := a.db.CheckinCreate(TeamId, UserId, Yesterday, Today, Blockers, Discuss, GoalsMet)
		if err != nil {
			if err.Error() == "REQUIRES_TEAM_USER" {
				Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, err.Error()))
				return
			}
			Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		Success(w, r, http.StatusOK, nil, nil)
	}
}

// handleCheckinCreate handles updating a team user checkin
// @Summary Update Team Checkin
// @Description Updates a team user checkin
// @Tags team
// @Produce  json
// @Success 200 object standardJsonResponse{}
// @Success 403 object standardJsonResponse{}
// @Success 500 object standardJsonResponse{}
// @Security ApiKeyAuth
// @Router /teams/{teamId}/checkins/{checkinId} [put]
func (a *api) handleCheckinUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		CheckinId := vars["checkinId"]

		keyVal := getJSONRequestBody(r, w)
		Yesterday := keyVal["yesterday"].(string)
		Today := keyVal["today"].(string)
		Blockers := keyVal["blockers"].(string)
		Discuss := keyVal["discuss"].(string)
		GoalsMet := keyVal["goalsMet"].(bool)

		err := a.db.CheckinUpdate(CheckinId, Yesterday, Today, Blockers, Discuss, GoalsMet)
		if err != nil {
			Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		Success(w, r, http.StatusOK, nil, nil)
	}
}

// handleCheckinDelete handles deleting a team user checkin
// @Summary Delete Team Checkin
// @Description Deletes a team user checkin
// @Tags team
// @Produce  json
// @Success 200 object standardJsonResponse{}
// @Success 403 object standardJsonResponse{}
// @Success 500 object standardJsonResponse{}
// @Security ApiKeyAuth
// @Router /teams/{teamId}/checkins/{checkinId} [delete]
func (a *api) handleCheckinDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		CheckinId := vars["checkinId"]

		err := a.db.CheckinDelete(CheckinId)
		if err != nil {
			Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		Success(w, r, http.StatusOK, nil, nil)
	}
}
