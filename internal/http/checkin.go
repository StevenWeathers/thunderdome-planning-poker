package http

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
	"go.uber.org/zap"

	"github.com/StevenWeathers/thunderdome-planning-poker/internal/http/checkin"
)

func resolveTeamCheckinDate(date string, timezone string) (string, error) {
	if timezone == "" {
		timezone = "America/New_York"
	}

	location, err := time.LoadLocation(timezone)
	if err != nil {
		return "", err
	}

	if date == "" {
		return time.Now().In(location).Format("2006-01-02"), nil
	}

	if _, err := time.Parse("2006-01-02", date); err != nil {
		return "", err
	}

	return date, nil
}

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

		if tz == "" {
			tz = "America/New_York"
		}

		location, tzErr := time.LoadLocation(tz)
		if tzErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, tzErr.Error()))
			return
		}

		if date == "" {
			date = time.Now().In(location).Format("2006-01-02")
		}

		checkins, err := s.CheckinDataSvc.CheckinList(ctx, teamID, date)
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
		query := r.URL.Query()
		date := query.Get("date")
		tz := query.Get("tz")

		if tz == "" {
			tz = "America/New_York"
		}

		location, tzErr := time.LoadLocation(tz)
		if tzErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, tzErr.Error()))
			return
		}

		if date == "" {
			date = time.Now().In(location).Format("2006-01-02")
		}

		checkin, err := s.CheckinDataSvc.CheckinLastByUser(ctx, teamID, userID, date)
		if err != nil && err.Error() != "NO_LAST_CHECKIN" {
			s.Logger.Ctx(ctx).Error("handleCheckinLastByUser error", zap.Error(err),
				zap.String("team_id", teamID),
				zap.String("user_id", userID),
				zap.String("checkins_date", date),
				zap.String("checkins_timezone", tz),
				zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		} else if err != nil && err.Error() == "NO_LAST_CHECKIN" {
			s.Logger.Ctx(ctx).Warn("handleCheckinLastByUser NO_LAST_CHECKIN",
				zap.String("team_id", teamID),
				zap.String("user_id", userID),
				zap.String("checkins_date", date),
				zap.String("checkins_timezone", tz),
				zap.String("session_user_id", sessionUserID))
			s.Success(w, r, http.StatusNoContent, nil, nil)
			return
		}

		s.Success(w, r, http.StatusOK, checkin, nil)
	}
}

type checkinCreateRequestBody struct {
	UserID      string `json:"userId" validate:"required,uuid"`
	CheckinDate string `json:"checkinDate"`
	TimeZone    string `json:"timeZone"`
	Yesterday   string `json:"yesterday"`
	Today       string `json:"today"`
	Blockers    string `json:"blockers"`
	Discuss     string `json:"discuss"`
	GoalsMet    bool   `json:"goalsMet"`
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

		if c.CheckinDate != "" {
			if _, err := time.Parse("2006-01-02", c.CheckinDate); err != nil {
				s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, err.Error()))
				return
			}
		}

		if c.TimeZone == "" {
			c.TimeZone = "America/New_York"
		}

		if _, err := time.LoadLocation(c.TimeZone); err != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, err.Error()))
			return
		}

		if c.CheckinDate == "" {
			location, _ := time.LoadLocation(c.TimeZone)
			c.CheckinDate = time.Now().In(location).Format("2006-01-02")
		}

		if s.Config.SubscriptionsEnabled {
			userType := ctx.Value(contextKeyUserType).(string)
			if userType != thunderdome.AdminUserType {
				location, _ := time.LoadLocation(c.TimeZone)
				todayInLocation := time.Now().In(location).Format("2006-01-02")
				targetDate := c.CheckinDate

				if targetDate > todayInLocation {
					orgID := r.PathValue("orgId")
					if orgID != "" {
						orgSubscribed, orgErr := s.OrganizationDataSvc.OrganizationIsSubscribed(ctx, orgID)
						teamSubscribed, teamErr := s.TeamDataSvc.TeamIsSubscribed(ctx, teamID)
						if (orgErr != nil || !orgSubscribed) && (teamErr != nil || !teamSubscribed) {
							s.Failure(w, r, http.StatusForbidden, Errorf(EUNAUTHORIZED, "ORGANIZATION_OR_TEAM_SUBSCRIPTION_REQUIRED"))
							return
						}
					} else {
						subscribed, err := s.TeamDataSvc.TeamIsSubscribed(ctx, teamID)
						if err != nil || !subscribed {
							s.Failure(w, r, http.StatusForbidden, Errorf(EUNAUTHORIZED, "TEAM_SUBSCRIPTION_REQUIRED"))
							return
						}
					}
				}
			}
		}

		normalizedBody, marshalErr := json.Marshal(c)
		if marshalErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, marshalErr.Error()))
			return
		}
		body = normalizedBody

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

type teamKudoRequestBody struct {
	KudoID       string `json:"kudoId" swaggerignore:"true"`
	TargetUserID string `json:"targetUserId" validate:"required,uuid"`
	KudosDate    string `json:"kudosDate"`
	TimeZone     string `json:"timeZone"`
	Comment      string `json:"comment" validate:"required"`
}

// handleTeamKudosGet gets a list of team kudos.
func (s *Service) handleTeamKudosGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)
		teamID := r.PathValue("teamId")
		if err := validate.Var(teamID, "required,uuid"); err != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, err.Error()))
			return
		}

		query := r.URL.Query()
		date, err := resolveTeamCheckinDate(query.Get("date"), query.Get("tz"))
		if err != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, err.Error()))
			return
		}

		kudos, err := s.CheckinDataSvc.KudoList(ctx, teamID, date)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleTeamKudosGet error", zap.Error(err), zap.String("team_id", teamID),
				zap.String("kudos_date", date), zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, kudos, nil)
	}
}

// handleTeamKudoGet gets a single team kudo.
func (s *Service) handleTeamKudoGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)
		teamID := r.PathValue("teamId")
		kudoID := r.PathValue("kudoId")
		if err := validate.Var(teamID, "required,uuid"); err != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, err.Error()))
			return
		}
		if err := validate.Var(kudoID, "required,uuid"); err != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, err.Error()))
			return
		}

		kudo, err := s.CheckinDataSvc.KudoGet(ctx, teamID, kudoID)
		if err != nil {
			if err.Error() == "NO_TEAM_KUDO" {
				s.Failure(w, r, http.StatusNotFound, Errorf(ENOTFOUND, err.Error()))
				return
			}
			s.Logger.Ctx(ctx).Error("handleTeamKudoGet error", zap.Error(err), zap.String("team_id", teamID),
				zap.String("kudo_id", kudoID), zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, kudo, nil)
	}
}

// handleTeamKudoCreate handles creating a team kudo.
func (s *Service) handleTeamKudoCreate(tc *checkin.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)
		teamID := r.PathValue("teamId")
		if err := validate.Var(teamID, "required,uuid"); err != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, err.Error()))
			return
		}

		var kudo teamKudoRequestBody
		body, bodyErr := io.ReadAll(r.Body)
		if bodyErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
			return
		}
		if err := json.Unmarshal(body, &kudo); err != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, err.Error()))
			return
		}
		if err := validate.Struct(kudo); err != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, err.Error()))
			return
		}

		resolvedDate, err := resolveTeamCheckinDate(kudo.KudosDate, kudo.TimeZone)
		if err != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, err.Error()))
			return
		}
		kudo.KudosDate = resolvedDate

		normalizedBody, err := json.Marshal(kudo)
		if err != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, err.Error()))
			return
		}

		result, err := tc.APIEvent(ctx, teamID, sessionUserID, "kudo_create", string(normalizedBody))
		if err != nil {
			if err.Error() == "TEAM_SUBSCRIPTION_REQUIRED" {
				s.Failure(w, r, http.StatusForbidden, Errorf(EUNAUTHORIZED, err.Error()))
				return
			}
			if err.Error() == "REQUIRES_TEAM_USER" || err.Error() == "KUDO_ALREADY_EXISTS" {
				s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, err.Error()))
				return
			}
			s.Logger.Ctx(ctx).Error("handleTeamKudoCreate error", zap.Error(err), zap.String("team_id", teamID),
				zap.String("target_user_id", kudo.TargetUserID), zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, result, nil)
	}
}

// handleTeamKudoUpdate handles updating a team kudo.
func (s *Service) handleTeamKudoUpdate(tc *checkin.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)
		teamID := r.PathValue("teamId")
		kudoID := r.PathValue("kudoId")
		if err := validate.Var(teamID, "required,uuid"); err != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, err.Error()))
			return
		}
		if err := validate.Var(kudoID, "required,uuid"); err != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, err.Error()))
			return
		}

		var kudo teamKudoRequestBody
		body, bodyErr := io.ReadAll(r.Body)
		if bodyErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
			return
		}
		if err := json.Unmarshal(body, &kudo); err != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, err.Error()))
			return
		}
		if err := validate.Struct(kudo); err != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, err.Error()))
			return
		}

		resolvedDate, err := resolveTeamCheckinDate(kudo.KudosDate, kudo.TimeZone)
		if err != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, err.Error()))
			return
		}
		kudo.KudoID = kudoID
		kudo.KudosDate = resolvedDate

		normalizedBody, err := json.Marshal(kudo)
		if err != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, err.Error()))
			return
		}

		result, err := tc.APIEvent(ctx, teamID, sessionUserID, "kudo_update", string(normalizedBody))
		if err != nil {
			if err.Error() == "TEAM_SUBSCRIPTION_REQUIRED" {
				s.Failure(w, r, http.StatusForbidden, Errorf(EUNAUTHORIZED, err.Error()))
				return
			}
			if err.Error() == "NO_TEAM_KUDO" {
				s.Failure(w, r, http.StatusNotFound, Errorf(ENOTFOUND, err.Error()))
				return
			}
			if err.Error() == "REQUIRES_TEAM_USER" || err.Error() == "KUDO_ALREADY_EXISTS" {
				s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, err.Error()))
				return
			}
			s.Logger.Ctx(ctx).Error("handleTeamKudoUpdate error", zap.Error(err), zap.String("team_id", teamID),
				zap.String("kudo_id", kudoID), zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, result, nil)
	}
}

// handleTeamKudoDelete handles deleting a team kudo.
func (s *Service) handleTeamKudoDelete(tc *checkin.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID := ctx.Value(contextKeyUserID).(string)
		teamID := r.PathValue("teamId")
		kudoID := r.PathValue("kudoId")
		if err := validate.Var(teamID, "required,uuid"); err != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, err.Error()))
			return
		}
		if err := validate.Var(kudoID, "required,uuid"); err != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, err.Error()))
			return
		}

		body, err := json.Marshal(struct {
			KudoID string `json:"kudoId"`
		}{KudoID: kudoID})
		if err != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, err.Error()))
			return
		}

		_, err = tc.APIEvent(ctx, teamID, userID, "kudo_delete", string(body))
		if err != nil {
			if err.Error() == "TEAM_SUBSCRIPTION_REQUIRED" {
				s.Failure(w, r, http.StatusForbidden, Errorf(EUNAUTHORIZED, err.Error()))
				return
			}
			if err.Error() == "NO_TEAM_KUDO" {
				s.Failure(w, r, http.StatusNotFound, Errorf(ENOTFOUND, err.Error()))
				return
			}
			s.Logger.Ctx(ctx).Error("handleTeamKudoDelete error", zap.Error(err), zap.String("team_id", teamID),
				zap.String("session_user_id", userID), zap.String("kudo_id", kudoID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, nil, nil)
	}
}
