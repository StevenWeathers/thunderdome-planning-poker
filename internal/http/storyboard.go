package http

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"go.uber.org/zap"

	"github.com/StevenWeathers/thunderdome-planning-poker/internal/http/storyboard"
	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
)

type storyboardCreateRequestBody struct {
	StoryboardName  string `json:"storyboardName" validate:"required"`
	JoinCode        string `json:"joinCode"`
	FacilitatorCode string `json:"facilitatorCode"`
}

// handleStoryboardCreate handles creating a storyboard (arena)
//
//	@Summary		Create Storyboard
//	@Description	Create a storyboard associated to the user
//	@Tags			storyboard
//	@Produce		json
//	@Param			userId			path	string						true	"the user ID"
//	@Param			orgId			path	string						false	"the organization ID"
//	@Param			departmentId	path	string						false	"the department ID"
//	@Param			teamId			path	string						false	"the team ID"
//	@Param			storyboard		body	storyboardCreateRequestBody	false	"new storyboard object"
//	@Success		200				object	standardJsonResponse{data=thunderdome.Storyboard}
//	@Failure		403				object	standardJsonResponse{}
//	@Failure		500				object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/users/{userId}/storyboards [post]
//	@Router			/teams/{teamId}/users/{userId}/storyboards [post]
//	@Router			/{orgId}/teams/{teamId}/users/{userId}/storyboards [post]
//	@Router			/{orgId}/departments/{departmentId}/teams/{teamId}/users/{userId}/storyboards [post]
func (s *Service) handleStoryboardCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)

		userID := r.PathValue("userId")
		idErr := validate.Var(userID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		teamID := r.PathValue("teamId")
		if teamID == "" && s.Config.RequireTeams {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, "STORYBOARD_CREATION_REQUIRES_TEAM"))
			return
		}

		body, bodyErr := io.ReadAll(r.Body) // check for errors
		if bodyErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
			return
		}

		var sb = storyboardCreateRequestBody{}
		jsonErr := json.Unmarshal(body, &sb)
		if jsonErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		inputErr := validate.Struct(sb)
		if inputErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, inputErr.Error()))
			return
		}

		var newStoryboard *thunderdome.Storyboard
		var err error
		// if storyboard created with team association
		if teamID != "" {
			if isTeamUserOrAnAdmin(r) {
				newStoryboard, err = s.StoryboardDataSvc.TeamCreateStoryboard(ctx, teamID, userID, sb.StoryboardName, sb.JoinCode, sb.FacilitatorCode)

				if err != nil {
					s.Logger.Ctx(ctx).Error("handleStoryboardCreate error", zap.Error(err),
						zap.String("entity_user_id", userID), zap.String("team_id", teamID),
						zap.String("storyboard_name", sb.StoryboardName), zap.String("session_user_id", sessionUserID))
					s.Failure(w, r, http.StatusInternalServerError, err)
					return
				}
			} else {
				s.Failure(w, r, http.StatusForbidden, Errorf(EUNAUTHORIZED, "REQUIRES_TEAM_USER"))
				return
			}
		} else {
			newStoryboard, err = s.StoryboardDataSvc.CreateStoryboard(ctx, userID, sb.StoryboardName, sb.JoinCode, sb.FacilitatorCode)
			if err != nil {
				s.Logger.Ctx(ctx).Error("handleStoryboardCreate error", zap.Error(err),
					zap.String("entity_user_id", userID), zap.String("session_user_id", sessionUserID),
					zap.String("storyboard_name", sb.StoryboardName))
				s.Failure(w, r, http.StatusInternalServerError, err)
				return
			}
		}

		s.Success(w, r, http.StatusOK, newStoryboard, nil)
	}
}

// handleStoryboardGet gets the storyboard by ID
//
//	@Summary		Get Storyboard
//	@Description	get storyboard by ID
//	@Tags			storyboard
//	@Produce		json
//	@Param			storyboardId	path	string	true	"the storyboard ID to get"
//	@Success		200				object	standardJsonResponse{data=thunderdome.Storyboard}
//	@Failure		403				object	standardJsonResponse{}
//	@Failure		404				object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/storyboards/{storyboardId} [get]
func (s *Service) handleStoryboardGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		storyboardID := r.PathValue("storyboardId")
		idErr := validate.Var(storyboardID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		sessionUserID := r.Context().Value(contextKeyUserID).(string)
		userType := r.Context().Value(contextKeyUserType).(string)

		sb, err := s.StoryboardDataSvc.GetStoryboardByID(storyboardID, sessionUserID)
		if err != nil {
			s.Failure(w, r, http.StatusNotFound, Errorf(ENOTFOUND, "STORYBOARD_NOT_FOUND"))
			return
		}

		// don't allow retrieving storyboard details if storyboard has JoinCode and user hasn't joined yet
		if sb.JoinCode != "" {
			UserErr := s.StoryboardDataSvc.GetStoryboardUserActiveStatus(storyboardID, sessionUserID)
			if UserErr != nil && userType != thunderdome.AdminUserType {
				s.Failure(w, r, http.StatusForbidden, Errorf(EUNAUTHORIZED, "USER_MUST_JOIN_STORYBOARD"))
				return
			}
		}

		s.Success(w, r, http.StatusOK, sb, nil)
	}
}

// handleGetUserStoryboards looks up storyboards associated with UserID
//
//	@Summary		Get Storyboards
//	@Description	get list of storyboards for the user
//	@Tags			storyboard
//	@Produce		json
//	@Param			userId	path	string	true	"the user ID to get storyboards for"
//	@Param			limit	query	int		false	"Max number of results to return"
//	@Param			offset	query	int		false	"Starting point to return rows from, should be multiplied by limit or 0"
//	@Success		200		object	standardJsonResponse{data=[]thunderdome.Storyboard}
//	@Failure		403		object	standardJsonResponse{}
//	@Failure		404		object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/users/{userId}/storyboards [get]
func (s *Service) handleGetUserStoryboards() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)
		limit, offset := getLimitOffsetFromRequest(r)

		userID := r.PathValue("userId")
		idErr := validate.Var(userID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		storyboards, count, err := s.StoryboardDataSvc.GetStoryboardsByUser(userID, limit, offset)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleGetUserStoryboards error", zap.Error(err), zap.Int("limit", limit),
				zap.Int("offset", offset), zap.String("entity_user_id", userID),
				zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusNotFound, Errorf(ENOTFOUND, "STORYBOARDS_NOT_FOUND"))
			return
		}

		meta := &pagination{
			Count:  count,
			Offset: offset,
			Limit:  limit,
		}

		s.Success(w, r, http.StatusOK, storyboards, meta)
	}
}

// handleGetStoryboards gets a list of storyboards
//
//	@Summary		Get Storyboards
//	@Description	get list of storyboards
//	@Tags			storyboard
//	@Produce		json
//	@Param			limit	query	int		false	"Max number of results to return"
//	@Param			offset	query	int		false	"Starting point to return rows from, should be multiplied by limit or 0"
//	@Param			active	query	boolean	false	"Only active storyboards"
//	@Success		200		object	standardJsonResponse{data=[]thunderdome.Storyboard}
//	@Failure		500		object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/storyboards [get]
func (s *Service) handleGetStoryboards() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)
		limit, offset := getLimitOffsetFromRequest(r)
		query := r.URL.Query()
		var err error
		var count int
		var storyboards []*thunderdome.Storyboard
		active, _ := strconv.ParseBool(query.Get("active"))

		if active {
			storyboards, count, err = s.StoryboardDataSvc.GetActiveStoryboards(limit, offset)
		} else {
			storyboards, count, err = s.StoryboardDataSvc.GetStoryboards(limit, offset)
		}

		if err != nil {
			s.Logger.Ctx(ctx).Error("handleGetStoryboards error", zap.Error(err), zap.Int("limit", limit),
				zap.Int("offset", offset), zap.Bool("storyboard_active", active),
				zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		meta := &pagination{
			Count:  count,
			Offset: offset,
			Limit:  limit,
		}

		s.Success(w, r, http.StatusOK, storyboards, meta)
	}
}

// handleStoryboardDelete handles deleting a storyboard
//
//	@Summary		Storyboard Delete
//	@Description	Delete a storyboard
//	@Param			storyboardId	path	string	true	"the storyboard ID"
//	@Tags			storyboard
//	@Produce		json
//	@Success		200	object	standardJsonResponse{}
//	@Success		403	object	standardJsonResponse{}
//	@Success		500	object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/storyboards/{storyboardId} [delete]
func (s *Service) handleStoryboardDelete(sb *storyboard.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		storyboardID := r.PathValue("storyboardId")
		idErr := validate.Var(storyboardID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		sessionUserID := r.Context().Value(contextKeyUserID).(string)

		err := sb.APIEvent(ctx, storyboardID, sessionUserID, "concede_storyboard", "")
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleStoryboardDelete error", zap.Error(err), zap.String("storyboard_id", storyboardID),
				zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, nil, nil)
	}
}

type storyboardGoalAddRequestBody struct {
	Name string `json:"name" validate:"required,min=1"`
}

// handleStoryboardGoalAdd handles adding a goal to a storyboard
//
//	@Summary		Storyboard Goal Add
//	@Description	Add a goal to a storyboard
//	@Param			storyboardId	path	string							true	"the storyboard ID"
//	@Param			storyboard		body	storyboardGoalAddRequestBody	false	"the goal to add"
//	@Tags			storyboard
//	@Produce		json
//	@Success		200	object	standardJsonResponse{}
//	@Success		403	object	standardJsonResponse{}
//	@Success		500	object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/storyboards/{storyboardId}/goals [post]
func (s *Service) handleStoryboardGoalAdd(sb *storyboard.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		storyboardID := r.PathValue("storyboardId")
		idErr := validate.Var(storyboardID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		sessionUserID := r.Context().Value(contextKeyUserID).(string)

		body, bodyErr := io.ReadAll(r.Body) // check for errors
		if bodyErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
			return
		}

		var sbm = storyboardGoalAddRequestBody{}
		jsonErr := json.Unmarshal(body, &sbm)
		if jsonErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		inputErr := validate.Struct(sbm)
		if inputErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, inputErr.Error()))
			return
		}

		err := sb.APIEvent(ctx, storyboardID, sessionUserID, "add_goal", sbm.Name)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handle storyboard goal add error",
				zap.Error(err),
				zap.String("storyboard_id", storyboardID),
				zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, nil, nil)
	}
}

type storyboardColumnAddRequestBody struct {
	GoalID string `json:"goalId" validate:"required,uuid"`
}

// handleStoryboardColumnAdd handles adding a column to a storyboard goal
//
//	@Summary		Storyboard Column Add
//	@Description	Add a column to a storyboard goal
//	@Param			storyboardId	path	string							true	"the storyboard ID"
//	@Param			storyboard		body	storyboardColumnAddRequestBody	false	"request body for adding a column"
//	@Tags			storyboard
//	@Produce		json
//	@Success		200	object	standardJsonResponse{}
//	@Success		403	object	standardJsonResponse{}
//	@Success		500	object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/storyboards/{storyboardId}/columns [post]
func (s *Service) handleStoryboardColumnAdd(sb *storyboard.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		storyboardID := r.PathValue("storyboardId")
		idErr := validate.Var(storyboardID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		sessionUserID := r.Context().Value(contextKeyUserID).(string)

		body, bodyErr := io.ReadAll(r.Body) // check for errors
		if bodyErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
			return
		}

		var sbm = storyboardColumnAddRequestBody{}
		jsonErr := json.Unmarshal(body, &sbm)
		if jsonErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		inputErr := validate.Struct(sbm)
		if inputErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, inputErr.Error()))
			return
		}

		eventValue, err := json.Marshal(sbm)
		if err != nil {
			s.Failure(w, r, http.StatusInternalServerError, Errorf(EINVALID, err.Error()))
			return
		}

		err = sb.APIEvent(ctx, storyboardID, sessionUserID, "add_column", string(eventValue))
		if err != nil {
			s.Logger.Ctx(ctx).Error("handle storyboard column add error",
				zap.Error(err),
				zap.String("storyboard_id", storyboardID),
				zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, nil, nil)
	}
}

type storyboardStoryAddRequestBody struct {
	GoalID   string `json:"goalId" validate:"required,uuid"`
	ColumnID string `json:"columnId" validate:"required,uuid"`
}

// handleStoryboardStoryAdd handles adding a story to a storyboard goal column
//
//	@Summary		Storyboard Story Add
//	@Description	Add a story to a storyboard goal column
//	@Param			storyboardId	path	string							true	"the storyboard ID"
//	@Param			storyboard		body	storyboardStoryAddRequestBody	false	"request body for adding a story"
//	@Tags			storyboard
//	@Produce		json
//	@Success		200	object	standardJsonResponse{}
//	@Success		403	object	standardJsonResponse{}
//	@Success		500	object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/storyboards/{storyboardId}/stories [post]
func (s *Service) handleStoryboardStoryAdd(sb *storyboard.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		storyboardID := r.PathValue("storyboardId")
		idErr := validate.Var(storyboardID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		sessionUserID := r.Context().Value(contextKeyUserID).(string)

		body, bodyErr := io.ReadAll(r.Body) // check for errors
		if bodyErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
			return
		}

		var sbm = storyboardStoryAddRequestBody{}
		jsonErr := json.Unmarshal(body, &sbm)
		if jsonErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		inputErr := validate.Struct(sbm)
		if inputErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, inputErr.Error()))
			return
		}

		eventValue, err := json.Marshal(sbm)
		if err != nil {
			s.Failure(w, r, http.StatusInternalServerError, Errorf(EINVALID, err.Error()))
			return
		}

		err = sb.APIEvent(ctx, storyboardID, sessionUserID, "add_story", string(eventValue))
		if err != nil {
			s.Logger.Ctx(ctx).Error("handle storyboard story add error",
				zap.Error(err),
				zap.String("storyboard_id", storyboardID),
				zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, nil, nil)
	}
}

type storyboardStoryMoveRequestBody struct {
	PlaceBefore string `json:"placeBefore" validate:"omitempty,uuid"`
	GoalID      string `json:"goalId" validate:"required,uuid"`
	ColumnID    string `json:"columnId" validate:"required,uuid"`
}

// handleStoryboardStoryMove handles moving a story in a storyboard
//
//	@Summary		Storyboard Story Move
//	@Description	Move a story in a storyboard
//	@Param			storyboardId	path	string							true	"the storyboard ID"
//	@Param			storyId			path	string							true	"the story ID"
//	@Param			storyboard		body	storyboardStoryMoveRequestBody	false	"target goal column and place before story"
//	@Tags			storyboard
//	@Produce		json
//	@Success		200	object	standardJsonResponse{}
//	@Success		403	object	standardJsonResponse{}
//	@Success		500	object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/storyboards/{storyboardId}/stories/{storyId}/move [put]
func (s *Service) handleStoryboardStoryMove(sb *storyboard.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		storyboardID := r.PathValue("storyboardId")
		idErr := validate.Var(storyboardID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		storyID := r.PathValue("storyId")
		idErr = validate.Var(storyboardID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		sessionUserID := r.Context().Value(contextKeyUserID).(string)

		body, bodyErr := io.ReadAll(r.Body) // check for errors
		if bodyErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
			return
		}

		var sbm = storyboardStoryMoveRequestBody{}
		jsonErr := json.Unmarshal(body, &sbm)
		if jsonErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		inputErr := validate.Struct(sbm)
		if inputErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, inputErr.Error()))
			return
		}

		type moveEvent struct {
			StoryID     string `json:"storyId"`
			PlaceBefore string `json:"placeBefore"`
			GoalID      string `json:"goalId"`
			ColumnID    string `json:"columnId"`
		}
		sbme := moveEvent{
			StoryID:     storyID,
			PlaceBefore: sbm.PlaceBefore,
			GoalID:      sbm.GoalID,
			ColumnID:    sbm.ColumnID,
		}
		moveEventJSON, moveEventErr := json.Marshal(sbme)
		if moveEventErr != nil {
			s.Failure(w, r, http.StatusInternalServerError, Errorf(EINVALID, moveEventErr.Error()))
			return
		}

		err := sb.APIEvent(ctx, storyboardID, sessionUserID, "move_story", string(moveEventJSON))
		if err != nil {
			s.Logger.Ctx(ctx).Error("handle storyboard story move error",
				zap.Error(err),
				zap.String("storyboard_id", storyboardID),
				zap.String("story_id", storyID),
				zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, nil, nil)
	}
}
