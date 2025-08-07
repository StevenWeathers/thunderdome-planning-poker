package http

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/StevenWeathers/thunderdome-planning-poker/internal/http/checkin"
)

// handleCheckinsGet gets a list of team checkins
//
//	@Summary		Get Team Checkins
//	@Description	Get a list of team checkins
//	@Tags			team
//	@Produce		json
//	@Param			teamId	path	string	true	"the team ID"
//	@Param			date	query	string	false	"the date in YYYY-MM-DD format"
//	@Param			tz		query	string	false	"the timezone name e.g. America/New_York"
//	@Success		200		object	standardJsonResponse{data=[]thunderdome.TeamCheckin}
//	@Security		ApiKeyAuth
//	@Router			/teams/{teamId}/checkins [get]
func (s *Service) handleCheckinsGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)
		teamID := r.PathValue("teamId")
		idErr := validate.Var(teamID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
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

		checkins, err := s.CheckinDataSvc.CheckinList(ctx, teamID, date, tz)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleCheckinsGet error", zap.Error(err), zap.String("team_id", teamID),
				zap.String("checkins_date", date), zap.String("checkins_timezone", tz),
				zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, checkins, nil)
	}
}

// handleCheckinLastByUser gets the last checkin by user
//
//	@Summary		Get Users last checkin for team
//	@Description	Get Users last checkin for team
//	@Tags			team
//	@Produce		json
//	@Param			teamId	path	string	true	"the team ID"
//	@Param			userId	path	string	false	"the user id to get last checkin for"
//	@Success		200		object	standardJsonResponse{data=thunderdome.TeamCheckin}
//	@Security		ApiKeyAuth
//	@Router			/teams/{teamId}/checkins/users/{userId}/last [get]
func (s *Service) handleCheckinLastByUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)
		teamID := r.PathValue("teamId")
		idErr := validate.Var(teamID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		userID := r.PathValue("userId")
		uidErr := validate.Var(userID, "required,uuid")
		if uidErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, uidErr.Error()))
			return
		}

		checkin, err := s.CheckinDataSvc.CheckinLastByUser(ctx, teamID, userID)
		if err != nil && err.Error() != "NO_LAST_CHECKIN" {
			s.Logger.Ctx(ctx).Error("handleCheckinLastByUser error", zap.Error(err),
				zap.String("team_id", teamID),
				zap.String("user_id", userID),
				zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		} else if err != nil && err.Error() == "NO_LAST_CHECKIN" {
			s.Logger.Ctx(ctx).Warn("handleCheckinLastByUser NO_LAST_CHECKIN",
				zap.String("team_id", teamID),
				zap.String("user_id", userID),
				zap.String("session_user_id", sessionUserID))
			s.Success(w, r, http.StatusNoContent, nil, nil)
			return
		}

		s.Success(w, r, http.StatusOK, checkin, nil)
	}
}

type checkinCreateRequestBody struct {
	UserID    string `json:"userId" validate:"required,uuid"`
	Yesterday string `json:"yesterday"`
	Today     string `json:"today"`
	Blockers  string `json:"blockers"`
	Discuss   string `json:"discuss"`
	GoalsMet  bool   `json:"goalsMet"`
}

// handleCheckinCreate handles creating a team user checkin
//
//	@Summary		Create Team Checkin
//	@Description	Creates a team user checkin
//	@Param			teamId	path	string						true	"the team ID"
//	@Param			checkin	body	checkinCreateRequestBody	true	"new check in object"
//	@Tags			team
//	@Produce		json
//	@Success		200	object	standardJsonResponse{}
//	@Success		403	object	standardJsonResponse{}
//	@Success		500	object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/teams/{teamId}/checkins [post]
func (s *Service) handleCheckinCreate(tc *checkin.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)
		teamID := r.PathValue("teamId")
		idErr := validate.Var(teamID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		var c = checkinCreateRequestBody{}
		body, bodyErr := io.ReadAll(r.Body)
		if bodyErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
			return
		}

		jsonErr := json.Unmarshal(body, &c)
		if jsonErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		inputErr := validate.Struct(c)
		if inputErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, inputErr.Error()))
			return
		}

		_, err := tc.APIEvent(ctx, teamID, c.UserID, "checkin_create", string(body))
		if err != nil {
			if err.Error() == "REQUIRES_TEAM_USER" {
				s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, err.Error()))
				return
			}
			s.Logger.Ctx(ctx).Error("handleCheckinCreate error", zap.Error(err), zap.String("team_id", teamID),
				zap.String("entity_user_id", c.UserID), zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, nil, nil)
	}
}

type checkinUpdateRequestBody struct {
	CheckinID string `json:"checkinId" swaggerignore:"true"`
	Yesterday string `json:"yesterday"`
	Today     string `json:"today"`
	Blockers  string `json:"blockers"`
	Discuss   string `json:"discuss"`
	GoalsMet  bool   `json:"goalsMet"`
}

// handleCheckinUpdate handles updating a team user checkin
//
//	@Summary		Update Team Checkin
//	@Description	Updates a team user checkin
//	@Param			teamId		path	string						true	"the team ID"
//	@Param			checkinId	path	string						true	"the checkin ID"
//	@Param			checkin		body	checkinUpdateRequestBody	true	"updated check in object"
//	@Tags			team
//	@Produce		json
//	@Success		200	object	standardJsonResponse{}
//	@Success		403	object	standardJsonResponse{}
//	@Success		500	object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/teams/{teamId}/checkins/{checkinId} [put]
func (s *Service) handleCheckinUpdate(tc *checkin.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)
		teamID := r.PathValue("teamId")
		idErr := validate.Var(teamID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		checkinID := r.PathValue("checkinId")
		idErr = validate.Var(checkinID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		var c = checkinUpdateRequestBody{}
		body, bodyErr := io.ReadAll(r.Body)
		if bodyErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
			return
		}

		jsonErr := json.Unmarshal(body, &c)
		if jsonErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		c.CheckinID = checkinID
		cu, jsonErr := json.Marshal(c)
		if jsonErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		_, err := tc.APIEvent(ctx, teamID, sessionUserID, "checkin_update", string(cu))
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleCheckinUpdate error", zap.Error(err), zap.String("team_id", teamID),
				zap.String("session_user_id", sessionUserID), zap.String("checkin_id", checkinID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, nil, nil)
	}
}

// handleCheckinDelete handles deleting a team user checkin
//
//	@Summary		Delete Team Checkin
//	@Description	Deletes a team user checkin
//	@Param			teamId		path	string	true	"the team ID"
//	@Param			checkinId	path	string	true	"the checkin ID"
//	@Tags			team
//	@Produce		json
//	@Success		200	object	standardJsonResponse{}
//	@Success		403	object	standardJsonResponse{}
//	@Success		500	object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/teams/{teamId}/checkins/{checkinId} [delete]
func (s *Service) handleCheckinDelete(tc *checkin.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)
		teamID := r.PathValue("teamId")
		idErr := validate.Var(teamID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		checkinID := r.PathValue("checkinId")
		idErr = validate.Var(checkinID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		type checkinDeleteRequest struct {
			CheckinId string `json:"checkinId"`
		}
		var c = checkinDeleteRequest{
			CheckinId: checkinID,
		}
		cu, jsonErr := json.Marshal(c)
		if jsonErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		_, err := tc.APIEvent(ctx, teamID, sessionUserID, "checkin_delete", string(cu))
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleCheckinDelete error", zap.Error(err), zap.String("team_id", teamID),
				zap.String("session_user_id", sessionUserID), zap.String("checkin_id", checkinID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, nil, nil)
	}
}

type checkinCommentRequestBody struct {
	CheckinID string `json:"checkinId" swaggerignore:"true"`
	CommentID string `json:"commentId" swaggerignore:"true"`
	UserID    string `json:"userId" validate:"required,uuid"`
	Comment   string `json:"comment" validate:"required"`
}

// handleCheckinComment handles creating a team user checkin comment
//
//	@Summary		Create Team Checkin Comment
//	@Description	Creates a team user checkin comment
//	@Param			teamId		path	string						true	"the team ID"
//	@Param			checkinId	path	string						true	"the checkin ID"
//	@Param			comment		body	checkinCommentRequestBody	true	"comment object"
//	@Tags			team
//	@Produce		json
//	@Success		200	object	standardJsonResponse{}
//	@Success		403	object	standardJsonResponse{}
//	@Success		500	object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/teams/{teamId}/checkins/{checkinId}/comments [post]
func (s *Service) handleCheckinComment(tc *checkin.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)
		teamID := r.PathValue("teamId")
		idErr := validate.Var(teamID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		checkinID := r.PathValue("checkinId")
		idErr = validate.Var(checkinID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		var c = checkinCommentRequestBody{}
		body, bodyErr := io.ReadAll(r.Body)
		if bodyErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
			return
		}

		jsonErr := json.Unmarshal(body, &c)
		if jsonErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		inputErr := validate.Struct(c)
		if inputErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, inputErr.Error()))
			return
		}

		c.CheckinID = checkinID
		cu, jsonErr := json.Marshal(c)
		if jsonErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		_, err := tc.APIEvent(ctx, teamID, c.UserID, "comment_create", string(cu))
		if err != nil {
			if err.Error() == "REQUIRES_TEAM_USER" {
				s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, err.Error()))
				return
			}
			s.Logger.Ctx(ctx).Error("handleCheckinComment error", zap.Error(err), zap.String("team_id", teamID),
				zap.String("entity_user_id", c.UserID), zap.String("checkin_id", checkinID),
				zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, nil, nil)
	}
}

// handleCheckinCommentEdit handles editing a team user checkin comment
//
//	@Summary		Edit Team Checkin Comment
//	@Description	Edits a team user checkin comment
//	@Param			teamId		path	string						true	"the team ID"
//	@Param			checkinId	path	string						true	"the checkin ID"
//	@Param			commentId	path	string						true	"the comment ID"
//	@Param			comment		body	checkinCommentRequestBody	true	"comment object"
//	@Tags			team
//	@Produce		json
//	@Success		200	object	standardJsonResponse{}
//	@Success		403	object	standardJsonResponse{}
//	@Success		500	object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/teams/{teamId}/checkins/{checkinId}/comments/{commentId} [put]
func (s *Service) handleCheckinCommentEdit(tc *checkin.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)
		teamID := r.PathValue("teamId")
		idErr := validate.Var(teamID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		commentID := r.PathValue("commentId")
		idErr = validate.Var(commentID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		var c = checkinCommentRequestBody{}
		body, bodyErr := io.ReadAll(r.Body)
		if bodyErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
			return
		}

		jsonErr := json.Unmarshal(body, &c)
		if jsonErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		inputErr := validate.Struct(c)
		if inputErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, inputErr.Error()))
			return
		}

		c.CommentID = commentID
		cu, jsonErr := json.Marshal(c)
		if jsonErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		_, err := tc.APIEvent(ctx, teamID, c.UserID, "comment_update", string(cu))
		if err != nil {
			if err.Error() == "REQUIRES_TEAM_USER" {
				s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, err.Error()))
				return
			}
			s.Logger.Ctx(ctx).Error("handleCheckinCommentEdit error", zap.Error(err), zap.String("team_id", teamID),
				zap.String("entity_user_id", c.UserID), zap.String("comment_id", commentID),
				zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, nil, nil)
	}
}

// handleCheckinCommentDelete handles deleting a team user checkin comment
//
//	@Summary		Delete Team Checkin Comment
//	@Description	Deletes a team user checkin comment
//	@Param			teamId		path	string	true	"the team ID"
//	@Param			checkinId	path	string	true	"the checkin ID"
//	@Param			commentId	path	string	true	"the comment ID"
//	@Tags			team
//	@Produce		json
//	@Success		200	object	standardJsonResponse{}
//	@Success		403	object	standardJsonResponse{}
//	@Success		500	object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/teams/{teamId}/checkins/{checkinId}/comments/{commentId} [delete]
func (s *Service) handleCheckinCommentDelete(tc *checkin.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID := ctx.Value(contextKeyUserID).(string)
		teamID := r.PathValue("teamId")
		idErr := validate.Var(teamID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		commentID := r.PathValue("commentId")
		idErr = validate.Var(commentID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		c := checkinCommentRequestBody{
			CommentID: commentID,
		}
		cu, jsonErr := json.Marshal(c)
		if jsonErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		_, err := tc.APIEvent(ctx, teamID, userID, "comment_delete", string(cu))
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleCheckinCommentDelete error", zap.Error(err), zap.String("team_id", teamID),
				zap.String("session_user_id", userID), zap.String("comment_id", commentID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, nil, nil)
	}
}
