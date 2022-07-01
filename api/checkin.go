package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/mux"
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
			a.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		a.Success(w, r, http.StatusOK, Checkins, nil)
	}
}

type checkinCreateRequestBody struct {
	UserId    string `json:"userId"`
	Yesterday string `json:"yesterday"`
	Today     string `json:"today"`
	Blockers  string `json:"blockers"`
	Discuss   string `json:"discuss"`
	GoalsMet  bool   `json:"goalsMet"`
}

// handleCheckinCreate handles creating a team user checkin
// @Summary Create Team Checkin
// @Description Creates a team user checkin
// @Param teamId path string true "the team ID"
// @Param checkin body checkinCreateRequestBody true "new check in object"
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

		var c = checkinCreateRequestBody{}
		body, bodyErr := ioutil.ReadAll(r.Body)
		if bodyErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
			return
		}

		jsonErr := json.Unmarshal(body, &c)
		if jsonErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		err := a.db.CheckinCreate(TeamId, c.UserId, c.Yesterday, c.Today, c.Blockers, c.Discuss, c.GoalsMet)
		if err != nil {
			if err.Error() == "REQUIRES_TEAM_USER" {
				a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, err.Error()))
				return
			}
			a.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		a.Success(w, r, http.StatusOK, nil, nil)
	}
}

type checkinUpdateRequestBody struct {
	Yesterday string `json:"yesterday"`
	Today     string `json:"today"`
	Blockers  string `json:"blockers"`
	Discuss   string `json:"discuss"`
	GoalsMet  bool   `json:"goalsMet"`
}

// handleCheckinUpdate handles updating a team user checkin
// @Summary Update Team Checkin
// @Description Updates a team user checkin
// @Param teamId path string true "the team ID"
// @Param checkinId path string true "the checkin ID"
// @Param checkin body checkinUpdateRequestBody true "updated check in object"
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

		var c = checkinUpdateRequestBody{}
		body, bodyErr := ioutil.ReadAll(r.Body)
		if bodyErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
			return
		}

		jsonErr := json.Unmarshal(body, &c)
		if jsonErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		err := a.db.CheckinUpdate(CheckinId, c.Yesterday, c.Today, c.Blockers, c.Discuss, c.GoalsMet)
		if err != nil {
			a.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		a.Success(w, r, http.StatusOK, nil, nil)
	}
}

// handleCheckinDelete handles deleting a team user checkin
// @Summary Delete Team Checkin
// @Description Deletes a team user checkin
// @Param teamId path string true "the team ID"
// @Param checkinId path string true "the checkin ID"
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
			a.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		a.Success(w, r, http.StatusOK, nil, nil)
	}
}

type checkinCommentRequestBody struct {
	UserID  string `json:"userId"`
	Comment string `json:"comment"`
}

// handleCheckinComment handles creating a team user checkin comment
// @Summary Create Team Checkin Comment
// @Description Creates a team user checkin comment
// @Param teamId path string true "the team ID"
// @Param checkinId path string true "the checkin ID"
// @Param comment body checkinCommentRequestBody true "comment object"
// @Tags team
// @Produce  json
// @Success 200 object standardJsonResponse{}
// @Success 403 object standardJsonResponse{}
// @Success 500 object standardJsonResponse{}
// @Security ApiKeyAuth
// @Router /teams/{teamId}/checkins/{checkinId}/comment [post]
func (a *api) handleCheckinComment() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		TeamId := vars["teamId"]
		CheckinId := vars["checkinId"]

		var c = checkinCommentRequestBody{}
		body, bodyErr := ioutil.ReadAll(r.Body)
		if bodyErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
			return
		}

		jsonErr := json.Unmarshal(body, &c)
		if jsonErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		err := a.db.CheckinComment(TeamId, CheckinId, c.UserID, c.Comment)
		if err != nil {
			if err.Error() == "REQUIRES_TEAM_USER" {
				a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, err.Error()))
				return
			}
			a.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		a.Success(w, r, http.StatusOK, nil, nil)
	}
}

// handleCheckinCommentDelete handles deleting a team user checkin comment
// @Summary Delete Team Checkin Comment
// @Description Deletes a team user checkin comment
// @Param teamId path string true "the team ID"
// @Param checkinId path string true "the checkin ID"
// @Param commentId path string true "the comment ID"
// @Tags team
// @Produce  json
// @Success 200 object standardJsonResponse{}
// @Success 403 object standardJsonResponse{}
// @Success 500 object standardJsonResponse{}
// @Security ApiKeyAuth
// @Router /teams/{teamId}/checkins/{checkinId}/comments/{commentId} [delete]
func (a *api) handleCheckinCommentDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		CommentId := vars["commentId"]

		err := a.db.CheckinCommentDelete(CommentId)
		if err != nil {
			a.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		a.Success(w, r, http.StatusOK, nil, nil)
	}
}
