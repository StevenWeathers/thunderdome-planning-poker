package http

import (
	"encoding/json"
	"io"
	"net/http"

	"go.uber.org/zap"

	"github.com/StevenWeathers/thunderdome-planning-poker/internal/http/storyboard"
)

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

		_, err := sb.APIEvent(ctx, storyboardID, sessionUserID, "add_goal", sbm.Name)
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

type storyboardGoalUpdateRequestBody struct {
	Name string `json:"name" validate:"required,min=1"`
}

// handleStoryboardGoalUpdate handles updating a goal in a storyboard
//
//	@Summary		Storyboard Goal Update
//	@Description	Update a goal in a storyboard
//	@Param			storyboardId	path	string							true	"the storyboard ID"
//	@Param			goalId			path	string							true	"the goal ID"
//	@Param			storyboard		body	storyboardGoalUpdateRequestBody	false	"the goal to update"
//	@Tags			storyboard
//	@Produce		json
//	@Success		200	object	standardJsonResponse{}
//	@Success		403	object	standardJsonResponse{}
//	@Success		500	object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/storyboards/{storyboardId}/goals/{goalId} [put]
func (s *Service) handleStoryboardGoalUpdate(sb *storyboard.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		storyboardID := r.PathValue("storyboardId")
		idErr := validate.Var(storyboardID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		goalID := r.PathValue("goalId")
		idErr = validate.Var(goalID, "required,uuid")
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

		var sbm = storyboardGoalUpdateRequestBody{}
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

		type updateEvent struct {
			GoalID string `json:"goalId"`
			Name   string `json:"name"`
		}
		sbue := updateEvent{
			GoalID: goalID,
			Name:   sbm.Name,
		}
		updateEventJSON, updateEventErr := json.Marshal(sbue)
		if updateEventErr != nil {
			s.Failure(w, r, http.StatusInternalServerError, Errorf(EINVALID, updateEventErr.Error()))
			return
		}

		_, err := sb.APIEvent(ctx, storyboardID, sessionUserID, "revise_goal", string(updateEventJSON))
		if err != nil {
			s.Logger.Ctx(ctx).Error("handle storyboard goal update error",
				zap.Error(err),
				zap.String("storyboard_id", storyboardID),
				zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, nil, nil)
	}
}

// handleStoryboardGoalDelete handles the deletion of a goal from a storyboard
//
//	@Summary		Storyboard Goal Delete
//	@Description	Delete a goal from a storyboard
//	@Param			storyboardId	path	string	true	"the storyboard ID"
//	@Param			goalId			path	string	true	"the goal ID"
//	@Tags			storyboard
//	@Produce		json
//	@Success		200	object	standardJsonResponse{}
//	@Success		403	object	standardJsonResponse{}
//	@Success		500	object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/storyboards/{storyboardId}/goals/{goalId} [delete]
func (s *Service) handleStoryboardGoalDelete(sb *storyboard.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		storyboardID := r.PathValue("storyboardId")
		idErr := validate.Var(storyboardID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		goalID := r.PathValue("goalId")
		idErr = validate.Var(goalID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		sessionUserID := r.Context().Value(contextKeyUserID).(string)

		_, err := sb.APIEvent(ctx, storyboardID, sessionUserID, "delete_goal", goalID)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handle storyboard goal delete error",
				zap.Error(err),
				zap.String("storyboard_id", storyboardID),
				zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, nil, nil)
	}
}
