package http

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"go.uber.org/zap"

	"github.com/StevenWeathers/thunderdome-planning-poker/internal/http/retro"
	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"

	"github.com/gorilla/mux"
)

type retroCreateRequestBody struct {
	RetroName             string  `json:"retroName" example:"sprint 10 retro" validate:"required"`
	JoinCode              string  `json:"joinCode" example:"iammadmax"`
	FacilitatorCode       string  `json:"facilitatorCode" example:"likeaboss"`
	MaxVotes              int     `json:"maxVotes" validate:"required,min=1,max=9"`
	BrainstormVisibility  string  `json:"brainstormVisibility" validate:"required,oneof=visible concealed hidden"`
	PhaseTimeLimitMin     int     `json:"phaseTimeLimitMin" validate:"min=0,max=59" example:"10"`
	PhaseAutoAdvance      bool    `json:"phaseAutoAdvance"`
	AllowCumulativeVoting bool    `json:"allowCumulativeVoting"`
	TemplateID            *string `json:"templateId"`
}

// handleRetroCreate handles creating a retro
// @Summary      Create Retro
// @Description  Create a retro associated to the user
// @Tags         retro
// @Produce      json
// @Param        userId        path    string                  true   "the user ID"
// @Param        orgId         path    string                  false  "the organization ID"
// @Param        departmentId  path    string                  false  "the department ID"
// @Param        teamId        path    string                  false  "the team ID"
// @Param        retro         body    retroCreateRequestBody  false  "new retro object"
// @Success      200           object  standardJsonResponse{data=thunderdome.Retro}
// @Failure      403           object  standardJsonResponse{}
// @Failure      500           object  standardJsonResponse{}
// @Security     ApiKeyAuth
// @Router       /users/{userId}/retros [post]
// @Router       /teams/{teamId}/users/{userId}/retros [post]
// @Router       /{orgId}/teams/{teamId}/users/{userId}/retros [post]
// @Router       /{orgId}/departments/{departmentId}/teams/{teamId}/users/{userId}/retros [post]
func (s *Service) handleRetroCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)
		vars := mux.Vars(r)
		userID := vars["userId"]
		idErr := validate.Var(userID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		teamID, teamIDExists := vars["teamId"]
		if !teamIDExists && s.Config.RequireTeams {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, "RETRO_CREATION_REQUIRES_TEAM"))
			return
		}

		body, bodyErr := io.ReadAll(r.Body)
		if bodyErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
			return
		}

		var nr = retroCreateRequestBody{}
		jsonErr := json.Unmarshal(body, &nr)
		if jsonErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		inputErr := validate.Struct(nr)
		if inputErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, inputErr.Error()))
			return
		}

		if nr.TemplateID == nil {
			// get default template
			template, err := s.RetroTemplateDataSvc.GetDefaultPublicTemplate(ctx)
			if err != nil {
				s.Logger.Ctx(ctx).Error("handleRetroCreate get default template by id error", zap.Error(err),
					zap.String("session_user_id", sessionUserID))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			nr.TemplateID = &template.ID
		}

		var newRetro *thunderdome.Retro
		var err error

		// if retro created with team association and user is not a team user or admin, return 403
		if teamIDExists && !isTeamUserOrAnAdmin(r) {
			s.Failure(w, r, http.StatusForbidden, Errorf(EUNAUTHORIZED, "REQUIRES_TEAM_USER"))
			return
		}

		newRetro, err = s.RetroDataSvc.CreateRetro(ctx, userID, teamID, nr.RetroName, nr.JoinCode, nr.FacilitatorCode, nr.MaxVotes, nr.BrainstormVisibility, nr.PhaseTimeLimitMin, nr.PhaseAutoAdvance, nr.AllowCumulativeVoting, *nr.TemplateID)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleRetroCreate error", zap.Error(err),
				zap.String("entity_user_id", userID),
				zap.String("retro_name", nr.RetroName),
				zap.String("session_user_id", sessionUserID),
				zap.String("team_id", teamID))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		s.Success(w, r, http.StatusOK, newRetro, nil)
	}
}

// handleRetroGet looks up retro or returns notfound status
// @Summary      Get Retro
// @Description  get retro by ID
// @Tags         retro
// @Produce      json
// @Param        retroId  path    string  true  "the retro ID to get"
// @Success      200      object  standardJsonResponse{data=thunderdome.Retro}
// @Failure      403      object  standardJsonResponse{}
// @Failure      404      object  standardJsonResponse{}
// @Security     ApiKeyAuth
// @Router       /retros/{retroId} [get]
func (s *Service) handleRetroGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		retroID := vars["retroId"]
		idErr := validate.Var(retroID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		sessionUserID := r.Context().Value(contextKeyUserID).(string)

		re, err := s.RetroDataSvc.RetroGet(retroID, sessionUserID)

		if err != nil {
			http.NotFound(w, r)
			return
		}

		s.Success(w, r, http.StatusOK, re, nil)
	}
}

// handleRetrosGetByUser looks up retros associated with userID
// @Summary      Get Retros by User
// @Description  get list of retros for the user
// @Tags         retro
// @Produce      json
// @Param        userId  path    string  true   "the user ID to get retros for"
// @Param        limit   query   int     false  "Max number of results to return"
// @Param        offset  query   int     false  "Starting point to return rows from, should be multiplied by limit or 0"
// @Success      200     object  standardJsonResponse{data=[]thunderdome.Retro}
// @Failure      403     object  standardJsonResponse{}
// @Failure      404     object  standardJsonResponse{}
// @Security     ApiKeyAuth
// @Router       /users/{userId}/retros [get]
func (s *Service) handleRetrosGetByUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		limit, offset := getLimitOffsetFromRequest(r)
		vars := mux.Vars(r)
		userID := vars["userId"]
		idErr := validate.Var(userID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		retros, count, err := s.RetroDataSvc.RetroGetByUser(userID, limit, offset)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		meta := &pagination{
			Count:  count,
			Offset: offset,
			Limit:  limit,
		}

		s.Success(w, r, http.StatusOK, retros, meta)
	}
}

// handleGetRetros gets a list of retros
// @Summary      Get Retros
// @Description  get list of retros
// @Tags         retro
// @Produce      json
// @Param        limit   query   int      false  "Max number of results to return"
// @Param        offset  query   int      false  "Starting point to return rows from, should be multiplied by limit or 0"
// @Param        active  query   boolean  false  "Only active retros"
// @Success      200     object  standardJsonResponse{data=[]thunderdome.Retro}
// @Failure      500     object  standardJsonResponse{}
// @Security     ApiKeyAuth
// @Router       /retros [get]
func (s *Service) handleGetRetros() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)
		limit, offset := getLimitOffsetFromRequest(r)
		query := r.URL.Query()
		var err error
		var count int
		var retros []*thunderdome.Retro
		active, _ := strconv.ParseBool(query.Get("active"))

		if active {
			retros, count, err = s.RetroDataSvc.GetActiveRetros(limit, offset)
		} else {
			retros, count, err = s.RetroDataSvc.GetRetros(limit, offset)
		}

		if err != nil {
			s.Logger.Ctx(ctx).Error("handleGetRetros error", zap.Error(err),
				zap.Int("limit", limit), zap.Int("offset", offset), zap.Bool("retro_active", active),
				zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		meta := &pagination{
			Count:  count,
			Offset: offset,
			Limit:  limit,
		}

		s.Success(w, r, http.StatusOK, retros, meta)
	}
}

type actionUpdateRequestBody struct {
	ActionID  string `json:"id" swaggerignore:"true" validate:"required,uuid"`
	Completed bool   `json:"completed" example:"false"`
	Content   string `json:"content" example:"update documentation" validate:"required"`
}

// handleRetroActionUpdate handles updating a retro action item
// @Summary      Retro Action Item Update
// @Description  Update a retro action item
// @Param        retroId     path  string                   true  "the retro ID"
// @Param        actionId    path  string                   true  "the action ID"
// @Param        actionItem  body  actionUpdateRequestBody  true  "updated action item"
// @Tags         retro
// @Produce      json
// @Success      200  object  standardJsonResponse{}
// @Success      403  object  standardJsonResponse{}
// @Success      500  object  standardJsonResponse{}
// @Security     ApiKeyAuth
// @Router       /retros/{retroId}/actions/{actionId} [put]
func (s *Service) handleRetroActionUpdate(retroSvc *retro.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var ra = actionUpdateRequestBody{}

		vars := mux.Vars(r)
		retroID := vars["retroId"]
		idErr := validate.Var(retroID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		actionID := vars["actionId"]
		idErr = validate.Var(actionID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		sessionUserID := r.Context().Value(contextKeyUserID).(string)

		body, bodyErr := io.ReadAll(r.Body)
		if bodyErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
			return
		}
		jsonErr := json.Unmarshal(body, &ra)
		if jsonErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		ra.ActionID = actionID
		inputErr := validate.Struct(ra)
		if inputErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, inputErr.Error()))
			return
		}
		updatedActionJson, _ := json.Marshal(ra)

		err := retroSvc.APIEvent(ctx, retroID, sessionUserID, "update_action", string(updatedActionJson))
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleRetroActionUpdate error", zap.Error(err),
				zap.String("retro_id", retroID), zap.String("session_user_id", sessionUserID),
				zap.String("action_id", actionID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, nil, nil)
	}
}

// handleRetroActionDelete handles deleting a retro action item
// @Summary      Retro Action Item Delete
// @Description  Delete a retro action item
// @Param        retroId   path  string  true  "the retro ID"
// @Param        actionId  path  string  true  "the action ID"
// @Tags         retro
// @Produce      json
// @Success      200  object  standardJsonResponse{}
// @Success      403  object  standardJsonResponse{}
// @Success      500  object  standardJsonResponse{}
// @Security     ApiKeyAuth
// @Router       /retros/{retroId}/actions/{actionId} [delete]
func (s *Service) handleRetroActionDelete(retroSvc *retro.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		vars := mux.Vars(r)
		retroID := vars["retroId"]
		idErr := validate.Var(retroID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		actionID := vars["actionId"]
		idErr = validate.Var(actionID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		sessionUserID := r.Context().Value(contextKeyUserID).(string)

		type actionItem struct {
			ActionID string `json:"id"`
		}
		deleteItem, _ := json.Marshal(actionItem{ActionID: actionID})

		err := retroSvc.APIEvent(ctx, retroID, sessionUserID, "delete_action", string(deleteItem))
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleRetroActionDelete error", zap.Error(err),
				zap.String("retro_id", retroID), zap.String("session_user_id", sessionUserID),
				zap.String("action_id", actionID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, nil, nil)
	}
}

type actionAddAssigneeRequestBody struct {
	ActionID string `json:"id" swaggerignore:"true" validate:"required,uuid"`
	UserID   string `json:"user_id" validate:"required,uuid"`
}

// handleRetroActionAssigneeAdd handles adding a retro action assignee
// @Summary      Retro Action Add Assignee
// @Description  Add a retro action assignee
// @Param        retroId     path  string                   true  "the retro ID"
// @Param        actionId    path  string                   true  "the action ID"
// @Param        actionItem  body  actionAddAssigneeRequestBody  true  "updated action item"
// @Tags         retro
// @Produce      json
// @Success      200  object  standardJsonResponse{}
// @Success      403  object  standardJsonResponse{}
// @Success      500  object  standardJsonResponse{}
// @Security     ApiKeyAuth
// @Router       /retros/{retroId}/actions/{actionId}/assignees [post]
func (s *Service) handleRetroActionAssigneeAdd(retroSvc *retro.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var ra = actionAddAssigneeRequestBody{}
		ctx := r.Context()
		vars := mux.Vars(r)
		retroID := vars["retroId"]
		idErr := validate.Var(retroID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		actionID := vars["actionId"]
		idErr = validate.Var(actionID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		sessionUserID := r.Context().Value(contextKeyUserID).(string)

		body, bodyErr := io.ReadAll(r.Body)
		if bodyErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
			return
		}
		jsonErr := json.Unmarshal(body, &ra)
		if jsonErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		ra.ActionID = actionID
		inputErr := validate.Struct(ra)
		if inputErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, inputErr.Error()))
			return
		}
		updatedActionJson, _ := json.Marshal(ra)

		err := retroSvc.APIEvent(ctx, retroID, sessionUserID, "action_assignee_add", string(updatedActionJson))
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleRetroActionAssigneeAdd error", zap.Error(err),
				zap.String("retro_id", retroID), zap.String("session_user_id", sessionUserID),
				zap.String("action_id", actionID), zap.String("action_user_id", ra.UserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, nil, nil)
	}
}

type actionRemoveAssigneeRequestBody struct {
	ActionID string `json:"id" swaggerignore:"true" validate:"required,uuid"`
	UserID   string `json:"user_id" validate:"required,uuid"`
}

// handleRetroActionAssigneeRemove handles removing a retro action assignee
// @Summary      Retro Action Remove Assignee
// @Description  Remove an assignee from a retro action
// @Param        retroId     path  string                   true  "the retro ID"
// @Param        actionId    path  string                   true  "the action ID"
// @Param        actionItem  body  actionRemoveAssigneeRequestBody  true  "updated action item"
// @Tags         retro
// @Produce      json
// @Success      200  object  standardJsonResponse{}
// @Success      403  object  standardJsonResponse{}
// @Success      500  object  standardJsonResponse{}
// @Security     ApiKeyAuth
// @Router       /retros/{retroId}/actions/{actionId}/assignees [delete]
func (s *Service) handleRetroActionAssigneeRemove(retroSvc *retro.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var ra = actionRemoveAssigneeRequestBody{}
		ctx := r.Context()
		vars := mux.Vars(r)
		retroID := vars["retroId"]
		idErr := validate.Var(retroID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		actionID := vars["actionId"]
		idErr = validate.Var(actionID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		sessionUserID := r.Context().Value(contextKeyUserID).(string)

		body, bodyErr := io.ReadAll(r.Body)
		if bodyErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
			return
		}
		jsonErr := json.Unmarshal(body, &ra)
		if jsonErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		ra.ActionID = actionID
		inputErr := validate.Struct(ra)
		if inputErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, inputErr.Error()))
			return
		}
		updatedActionJson, _ := json.Marshal(ra)

		err := retroSvc.APIEvent(ctx, retroID, sessionUserID, "action_assignee_remove", string(updatedActionJson))
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleRetroActionAssigneeRemove error", zap.Error(err),
				zap.String("retro_id", retroID), zap.String("session_user_id", sessionUserID),
				zap.String("action_id", actionID), zap.String("action_user_id", ra.UserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, nil, nil)
	}
}

type actionCommentRequestBody struct {
	Comment string `json:"comment" validate:"required"`
}

// handleRetroActionCommentAdd handles adding a comment to a retro action item
// @Summary      Retro Action Item Comment
// @Description  Add a comment to a retro action item
// @Param        retroId     path  string                    true  "the retro ID"
// @Param        actionId    path  string                    true  "the action ID"
// @Param        actionItem  body  actionCommentRequestBody  true  "action comment"
// @Tags         retro
// @Produce      json
// @Success      200  object  standardJsonResponse{}
// @Success      403  object  standardJsonResponse{}
// @Success      500  object  standardJsonResponse{}
// @Security     ApiKeyAuth
// @Router       /retros/{retroId}/actions/{actionId}/comments [post]
func (s *Service) handleRetroActionCommentAdd() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var ra = actionCommentRequestBody{}
		ctx := r.Context()
		vars := mux.Vars(r)
		retroID := vars["retroId"]
		idErr := validate.Var(retroID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		actionID := vars["actionId"]
		idErr = validate.Var(actionID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		sessionUserID := r.Context().Value(contextKeyUserID).(string)

		body, bodyErr := io.ReadAll(r.Body)
		if bodyErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
			return
		}
		jsonErr := json.Unmarshal(body, &ra)
		if jsonErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		inputErr := validate.Struct(ra)
		if inputErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, inputErr.Error()))
			return
		}

		action, err := s.RetroDataSvc.RetroActionCommentAdd(retroID, actionID, sessionUserID, ra.Comment)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleRetroActionCommentAdd error", zap.Error(err),
				zap.String("retro_id", retroID), zap.String("session_user_id", sessionUserID),
				zap.String("action_id", actionID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, action, nil)
	}
}

// handleRetroActionCommentEdit handles editing a retro action item comment
// @Summary      Retro Action Item Comment Edit
// @Description  Edit a retro action item comment
// @Param        retroId     path  string                    true  "the retro ID"
// @Param        actionId    path  string                    true  "the action ID"
// @Param        commentId   path  string                    true  "the comment ID"
// @Param        actionItem  body  actionCommentRequestBody  true  "action comment"
// @Tags         retro
// @Produce      json
// @Success      200  object  standardJsonResponse{}
// @Success      403  object  standardJsonResponse{}
// @Success      500  object  standardJsonResponse{}
// @Security     ApiKeyAuth
// @Router       /retros/{retroId}/actions/{actionId}/comments/{commentId} [put]
func (s *Service) handleRetroActionCommentEdit() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var ra = actionCommentRequestBody{}
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)
		vars := mux.Vars(r)
		retroID := vars["retroId"]
		idErr := validate.Var(retroID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		actionID := vars["actionId"]
		idErr = validate.Var(actionID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		commentID := vars["commentId"]
		idErr = validate.Var(commentID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		body, bodyErr := io.ReadAll(r.Body)
		if bodyErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
			return
		}
		jsonErr := json.Unmarshal(body, &ra)
		if jsonErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		inputErr := validate.Struct(ra)
		if inputErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, inputErr.Error()))
			return
		}

		action, err := s.RetroDataSvc.RetroActionCommentEdit(retroID, actionID, commentID, ra.Comment)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleRetroActionCommentEdit error", zap.Error(err),
				zap.String("retro_id", retroID), zap.String("action_id", actionID),
				zap.String("comment_id", commentID), zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, action, nil)
	}
}

// handleRetroActionCommentDelete handles delete a comment from a retro action item
// @Summary      Retro Action Item Comment Delete
// @Description  Delete a comment from a retro action item
// @Param        retroId    path  string  true  "the retro ID"
// @Param        actionId   path  string  true  "the action ID"
// @Param        commentId  path  string  true  "the comment ID"
// @Tags         retro
// @Produce      json
// @Success      200  object  standardJsonResponse{}
// @Success      403  object  standardJsonResponse{}
// @Success      500  object  standardJsonResponse{}
// @Security     ApiKeyAuth
// @Router       /retros/{retroId}/actions/{actionId}/comments/{commentId} [post]
func (s *Service) handleRetroActionCommentDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)
		vars := mux.Vars(r)
		retroID := vars["retroId"]
		idErr := validate.Var(retroID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		actionID := vars["actionId"]
		idErr = validate.Var(actionID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		commentID := vars["commentId"]
		idErr = validate.Var(commentID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		action, err := s.RetroDataSvc.RetroActionCommentDelete(retroID, actionID, commentID)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleRetroActionCommentDelete error", zap.Error(err),
				zap.String("retro_id", retroID), zap.String("action_id", actionID),
				zap.String("comment_id", commentID), zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, action, nil)
	}
}

// handleRetroDelete handles deleting a retro
// @Summary      Retro Delete
// @Description  Delete a retro
// @Param        retroId  path  string  true  "the retro ID"
// @Tags         retro
// @Produce      json
// @Success      200  object  standardJsonResponse{}
// @Success      403  object  standardJsonResponse{}
// @Success      500  object  standardJsonResponse{}
// @Security     ApiKeyAuth
// @Router       /retros/{retroId} [delete]
func (s *Service) handleRetroDelete(retroSvc *retro.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		vars := mux.Vars(r)
		retroID := vars["retroId"]
		idErr := validate.Var(retroID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		sessionUserID := ctx.Value(contextKeyUserID).(string)

		err := retroSvc.APIEvent(ctx, retroID, sessionUserID, "concede_retro", "")
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleRetroDelete error", zap.Error(err),
				zap.String("retro_id", retroID), zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, nil, nil)
	}
}
