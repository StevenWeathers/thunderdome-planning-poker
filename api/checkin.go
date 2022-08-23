package api

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/StevenWeathers/thunderdome-planning-poker/api/checkin"

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
		idErr := validate.Var(TeamID, "required,uuid")
		if idErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		query := r.URL.Query()
		date := query.Get("date")
		tz := query.Get("tz")

		if date == "" {
			date = time.Now().Format("1988-01-02")
		}

		if tz == "" {
			tz = "America/New_York"
		}

		Checkins, err := a.db.CheckinList(r.Context(), TeamID, date, tz)
		if err != nil {
			a.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		a.Success(w, r, http.StatusOK, Checkins, nil)
	}
}

type checkinCreateRequestBody struct {
	UserId    string `json:"userId" validate:"required,uuid"`
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
func (a *api) handleCheckinCreate(tc *checkin.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		TeamId := vars["teamId"]
		idErr := validate.Var(TeamId, "required,uuid")
		if idErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		var c = checkinCreateRequestBody{}
		body, bodyErr := io.ReadAll(r.Body)
		if bodyErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
			return
		}

		jsonErr := json.Unmarshal(body, &c)
		if jsonErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		inputErr := validate.Struct(c)
		if inputErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, inputErr.Error()))
			return
		}

		err := tc.APIEvent(r.Context(), TeamId, c.UserId, "checkin_create", string(body))
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
	CheckinId string `json:"checkinId" swaggerignore:"true"`
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
func (a *api) handleCheckinUpdate(tc *checkin.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userId := ctx.Value(contextKeyUserID).(string)
		vars := mux.Vars(r)
		TeamId := vars["teamId"]
		idErr := validate.Var(TeamId, "required,uuid")
		if idErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		CheckinId := vars["checkinId"]
		idErr = validate.Var(CheckinId, "required,uuid")
		if idErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		var c = checkinUpdateRequestBody{}
		body, bodyErr := io.ReadAll(r.Body)
		if bodyErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
			return
		}

		jsonErr := json.Unmarshal(body, &c)
		if jsonErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		c.CheckinId = CheckinId
		cu, jsonErr := json.Marshal(c)
		if jsonErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		err := tc.APIEvent(ctx, TeamId, userId, "checkin_update", string(cu))
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
func (a *api) handleCheckinDelete(tc *checkin.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userId := ctx.Value(contextKeyUserID).(string)
		vars := mux.Vars(r)
		TeamId := vars["teamId"]
		idErr := validate.Var(TeamId, "required,uuid")
		if idErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		CheckinId := vars["checkinId"]
		idErr = validate.Var(CheckinId, "required,uuid")
		if idErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		type checkinDeleteRequest struct {
			CheckinId string `json:"checkinId"`
		}
		var c = checkinDeleteRequest{
			CheckinId: CheckinId,
		}
		cu, jsonErr := json.Marshal(c)
		if jsonErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		err := tc.APIEvent(ctx, TeamId, userId, "checkin_delete", string(cu))
		if err != nil {
			a.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		a.Success(w, r, http.StatusOK, nil, nil)
	}
}

type checkinCommentRequestBody struct {
	CheckinId string `json:"checkinId" swaggerignore:"true"`
	CommentId string `json:"commentId" swaggerignore:"true"`
	UserID    string `json:"userId" validate:"required,uuid"`
	Comment   string `json:"comment" validate:"required"`
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
// @Router /teams/{teamId}/checkins/{checkinId}/comments [post]
func (a *api) handleCheckinComment(tc *checkin.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		vars := mux.Vars(r)
		TeamId := vars["teamId"]
		idErr := validate.Var(TeamId, "required,uuid")
		if idErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		CheckinId := vars["checkinId"]
		idErr = validate.Var(CheckinId, "required,uuid")
		if idErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		var c = checkinCommentRequestBody{}
		body, bodyErr := io.ReadAll(r.Body)
		if bodyErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
			return
		}

		jsonErr := json.Unmarshal(body, &c)
		if jsonErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		inputErr := validate.Struct(c)
		if inputErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, inputErr.Error()))
			return
		}

		c.CheckinId = CheckinId
		cu, jsonErr := json.Marshal(c)
		if jsonErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		err := tc.APIEvent(ctx, TeamId, c.UserID, "comment_create", string(cu))
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

// handleCheckinCommentEdit handles editing a team user checkin comment
// @Summary Edit Team Checkin Comment
// @Description Edits a team user checkin comment
// @Param teamId path string true "the team ID"
// @Param checkinId path string true "the checkin ID"
// @Param comment body checkinCommentRequestBody true "comment object"
// @Tags team
// @Produce  json
// @Success 200 object standardJsonResponse{}
// @Success 403 object standardJsonResponse{}
// @Success 500 object standardJsonResponse{}
// @Security ApiKeyAuth
// @Router /teams/{teamId}/checkins/{checkinId}/comments [put]
func (a *api) handleCheckinCommentEdit(tc *checkin.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		vars := mux.Vars(r)
		TeamId := vars["teamId"]
		idErr := validate.Var(TeamId, "required,uuid")
		if idErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		CommentId := vars["commentId"]
		idErr = validate.Var(CommentId, "required,uuid")
		if idErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		var c = checkinCommentRequestBody{}
		body, bodyErr := io.ReadAll(r.Body)
		if bodyErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
			return
		}

		jsonErr := json.Unmarshal(body, &c)
		if jsonErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		inputErr := validate.Struct(c)
		if inputErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, inputErr.Error()))
			return
		}

		c.CommentId = CommentId
		cu, jsonErr := json.Marshal(c)
		if jsonErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		err := tc.APIEvent(ctx, TeamId, c.UserID, "comment_update", string(cu))
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
func (a *api) handleCheckinCommentDelete(tc *checkin.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userId := ctx.Value(contextKeyUserID).(string)
		vars := mux.Vars(r)
		TeamId := vars["teamId"]
		idErr := validate.Var(TeamId, "required,uuid")
		if idErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		CommentId := vars["commentId"]
		idErr = validate.Var(CommentId, "required,uuid")
		if idErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		c := checkinCommentRequestBody{
			CommentId: CommentId,
		}
		cu, jsonErr := json.Marshal(c)
		if jsonErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		err := tc.APIEvent(ctx, TeamId, userId, "comment_delete", string(cu))
		if err != nil {
			a.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		a.Success(w, r, http.StatusOK, nil, nil)
	}
}
