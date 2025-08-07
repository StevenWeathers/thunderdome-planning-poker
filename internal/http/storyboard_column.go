package http

import (
	"encoding/json"
	"io"
	"net/http"

	"go.uber.org/zap"

	"github.com/StevenWeathers/thunderdome-planning-poker/internal/http/storyboard"
)

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

type storyboardColumnUpdateRequestBody struct {
	Name string `json:"name" validate:"required"`
}

// handleStoryboardColumnUpdate handles updating a column in a storyboard
//
//	@Summary		Storyboard Column Update
//	@Description	Update a column in a storyboard
//	@Param			storyboardId	path	string								true	"the storyboard ID"
//	@Param			columnId		path	string								true	"the column ID"
//	@Param			storyboard		body	storyboardColumnUpdateRequestBody	false	"the column to update"
//	@Tags			storyboard
//	@Produce		json
//	@Success		200	object	standardJsonResponse{}
//	@Success		403	object	standardJsonResponse{}
//	@Success		500	object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/storyboards/{storyboardId}/columns/{columnId} [put]
func (s *Service) handleStoryboardColumnUpdate(sb *storyboard.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		storyboardID := r.PathValue("storyboardId")
		idErr := validate.Var(storyboardID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		columnID := r.PathValue("columnId")
		idErr = validate.Var(columnID, "required,uuid")
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

		var sbm = storyboardColumnUpdateRequestBody{}
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
			ColumnID string `json:"id"`
			Name     string `json:"name"`
		}
		sbue := updateEvent{
			ColumnID: columnID,
			Name:     sbm.Name,
		}
		updateEventJSON, updateEventErr := json.Marshal(sbue)
		if updateEventErr != nil {
			s.Failure(w, r, http.StatusInternalServerError, Errorf(EINVALID, updateEventErr.Error()))
			return
		}

		err := sb.APIEvent(ctx, storyboardID, sessionUserID, "revise_column", string(updateEventJSON))
		if err != nil {
			s.Logger.Ctx(ctx).Error("handle storyboard column update error",
				zap.Error(err),
				zap.String("storyboard_id", storyboardID),
				zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, nil, nil)
	}
}

// handleStoryboardColumnDelete handles the deletion of a column from a storyboard
//
//	@Summary		Storyboard Column Delete
//	@Description	Delete a column from a storyboard
//	@Param			storyboardId	path	string	true	"the storyboard ID"
//	@Param			columnId		path	string	true	"the column ID"
//	@Tags			storyboard
//	@Produce		json
//	@Success		200	object	standardJsonResponse{}
//	@Success		403	object	standardJsonResponse{}
//	@Success		500	object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/storyboards/{storyboardId}/columns/{columnId} [delete]
func (s *Service) handleStoryboardColumnDelete(sb *storyboard.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		storyboardID := r.PathValue("storyboardId")
		idErr := validate.Var(storyboardID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		columnID := r.PathValue("columnId")
		idErr = validate.Var(columnID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		sessionUserID := r.Context().Value(contextKeyUserID).(string)

		err := sb.APIEvent(ctx, storyboardID, sessionUserID, "delete_column", columnID)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handle storyboard column delete error",
				zap.Error(err),
				zap.String("storyboard_id", storyboardID),
				zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, nil, nil)
	}
}
