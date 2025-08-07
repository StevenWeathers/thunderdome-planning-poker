package http

import (
	"encoding/json"
	"io"
	"net/http"

	"go.uber.org/zap"

	"github.com/StevenWeathers/thunderdome-planning-poker/internal/http/storyboard"
)

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
//	@Success		200	object	standardJsonResponse{data=thunderdome.StoryboardStory}
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

		newStory, err := sb.APIEvent(ctx, storyboardID, sessionUserID, "add_story", string(eventValue))
		if err != nil {
			s.Logger.Ctx(ctx).Error("handle storyboard story add error",
				zap.Error(err),
				zap.String("storyboard_id", storyboardID),
				zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, newStory, nil)
	}
}

type storyboardStoryNameRequestBody struct {
	Name string `json:"name" validate:"required"`
}

// handleStoryboardStoryNameUpdate handles updating a story name in a storyboard
//
//	@Summary		Storyboard Story Name Update
//	@Description	Updates a story name in a storyboard
//	@Param			storyboardId	path	string							true	"the storyboard ID"
//	@Param			storyId			path	string							true	"the story ID"
//	@Param			storyboard		body	storyboardStoryMoveRequestBody	false	"new story name"
//	@Tags			storyboard
//	@Produce		json
//	@Success		200	object	standardJsonResponse{}
//	@Success		403	object	standardJsonResponse{}
//	@Success		500	object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/storyboards/{storyboardId}/stories/{storyId}/name [put]
func (s *Service) handleStoryboardStoryNameUpdate(sb *storyboard.Service) http.HandlerFunc {
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

		var ub = storyboardStoryNameRequestBody{}
		jsonErr := json.Unmarshal(body, &ub)
		if jsonErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		inputErr := validate.Struct(ub)
		if inputErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, inputErr.Error()))
			return
		}

		type nameEvent struct {
			StoryID string `json:"storyId"`
			Name    string `json:"name"`
		}
		sbue := nameEvent{
			StoryID: storyID,
			Name:    ub.Name,
		}
		updateEventJSON, updateEventErr := json.Marshal(sbue)
		if updateEventErr != nil {
			s.Failure(w, r, http.StatusInternalServerError, Errorf(EINVALID, updateEventErr.Error()))
			return
		}

		_, err := sb.APIEvent(ctx, storyboardID, sessionUserID, "update_story_name", string(updateEventJSON))
		if err != nil {
			s.Logger.Ctx(ctx).Error("handle storyboard story name update error",
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

type storyboardStoryContentRequestBody struct {
	Content string `json:"content" validate:"required"`
}

// handleStoryboardStoryContentUpdate handles updating a story content in a storyboard
//
//	@Summary		Storyboard Story Content Update
//	@Description	Updates a story content in a storyboard
//	@Param			storyboardId	path	string								true	"the storyboard ID"
//	@Param			storyId			path	string								true	"the story ID"
//	@Param			storyboard		body	storyboardStoryContentRequestBody	false	"new story content"
//	@Tags			storyboard
//	@Produce		json
//	@Success		200	object	standardJsonResponse{}
//	@Success		403	object	standardJsonResponse{}
//	@Success		500	object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/storyboards/{storyboardId}/stories/{storyId}/content [put]
func (s *Service) handleStoryboardStoryContentUpdate(sb *storyboard.Service) http.HandlerFunc {
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

		var ub = storyboardStoryContentRequestBody{}
		jsonErr := json.Unmarshal(body, &ub)
		if jsonErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		inputErr := validate.Struct(ub)
		if inputErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, inputErr.Error()))
			return
		}

		type contentEvent struct {
			StoryID string `json:"storyId"`
			Content string `json:"content"`
		}
		sbue := contentEvent{
			StoryID: storyID,
			Content: ub.Content,
		}
		updateEventJSON, updateEventErr := json.Marshal(sbue)
		if updateEventErr != nil {
			s.Failure(w, r, http.StatusInternalServerError, Errorf(EINVALID, updateEventErr.Error()))
			return
		}

		_, err := sb.APIEvent(ctx, storyboardID, sessionUserID, "update_story_content", string(updateEventJSON))
		if err != nil {
			s.Logger.Ctx(ctx).Error("handle storyboard story content update error",
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

type storyboardStoryColorRequestBody struct {
	Color string `json:"color" validate:"required"`
}

// handleStoryboardStoryColorUpdate handles updating a story color in a storyboard
//
//	@Summary		Storyboard Story Color Update
//	@Description	Updates a story color in a storyboard
//	@Param			storyboardId	path	string							true	"the storyboard ID"
//	@Param			storyId			path	string							true	"the story ID"
//	@Param			storyboard		body	storyboardStoryMoveRequestBody	false	"new story color"
//	@Tags			storyboard
//	@Produce		json
//	@Success		200	object	standardJsonResponse{}
//	@Success		403	object	standardJsonResponse{}
//	@Success		500	object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/storyboards/{storyboardId}/stories/{storyId}/color [put]
func (s *Service) handleStoryboardStoryColorUpdate(sb *storyboard.Service) http.HandlerFunc {
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

		var ub = storyboardStoryColorRequestBody{}
		jsonErr := json.Unmarshal(body, &ub)
		if jsonErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		inputErr := validate.Struct(ub)
		if inputErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, inputErr.Error()))
			return
		}

		type colorEvent struct {
			StoryID string `json:"storyId"`
			Color   string `json:"color"`
		}
		sbue := colorEvent{
			StoryID: storyID,
			Color:   ub.Color,
		}
		updateEventJSON, updateEventErr := json.Marshal(sbue)
		if updateEventErr != nil {
			s.Failure(w, r, http.StatusInternalServerError, Errorf(EINVALID, updateEventErr.Error()))
			return
		}

		_, err := sb.APIEvent(ctx, storyboardID, sessionUserID, "update_story_color", string(updateEventJSON))
		if err != nil {
			s.Logger.Ctx(ctx).Error("handle storyboard story color update error",
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

type storyboardStoryPointsRequestBody struct {
	Points string `json:"points" validate:"required"`
}

// handleStoryboardStoryPointsUpdate handles updating a story points in a storyboard
//
//	@Summary		Storyboard Story Points Update
//	@Description	Updates a story points in a storyboard
//	@Param			storyboardId	path	string							true	"the storyboard ID"
//	@Param			storyId			path	string							true	"the story ID"
//	@Param			storyboard		body	storyboardStoryMoveRequestBody	false	"new story points"
//	@Tags			storyboard
//	@Produce		json
//	@Success		200	object	standardJsonResponse{}
//	@Success		403	object	standardJsonResponse{}
//	@Success		500	object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/storyboards/{storyboardId}/stories/{storyId}/points [put]
func (s *Service) handleStoryboardStoryPointsUpdate(sb *storyboard.Service) http.HandlerFunc {
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

		var ub = storyboardStoryPointsRequestBody{}
		jsonErr := json.Unmarshal(body, &ub)
		if jsonErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		inputErr := validate.Struct(ub)
		if inputErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, inputErr.Error()))
			return
		}

		type pointsEvent struct {
			StoryID string `json:"storyId"`
			Points  string `json:"points"`
		}
		sbue := pointsEvent{
			StoryID: storyID,
			Points:  ub.Points,
		}
		updateEventJSON, updateEventErr := json.Marshal(sbue)
		if updateEventErr != nil {
			s.Failure(w, r, http.StatusInternalServerError, Errorf(EINVALID, updateEventErr.Error()))
			return
		}

		_, err := sb.APIEvent(ctx, storyboardID, sessionUserID, "update_story_points", string(updateEventJSON))
		if err != nil {
			s.Logger.Ctx(ctx).Error("handle storyboard story points update error",
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

type storyboardStoryClosedRequestBody struct {
	Closed string `json:"closed" validate:"required"`
}

// handleStoryboardStoryClosedUpdate handles updating a story closed in a storyboard
//
//	@Summary		Storyboard Story Closed Update
//	@Description	Updates a story closed in a storyboard
//	@Param			storyboardId	path	string							true	"the storyboard ID"
//	@Param			storyId			path	string							true	"the story ID"
//	@Param			storyboard		body	storyboardStoryMoveRequestBody	false	"new story closed"
//	@Tags			storyboard
//	@Produce		json
//	@Success		200	object	standardJsonResponse{}
//	@Success		403	object	standardJsonResponse{}
//	@Success		500	object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/storyboards/{storyboardId}/stories/{storyId}/closed [put]
func (s *Service) handleStoryboardStoryClosedUpdate(sb *storyboard.Service) http.HandlerFunc {
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

		var ub = storyboardStoryClosedRequestBody{}
		jsonErr := json.Unmarshal(body, &ub)
		if jsonErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		inputErr := validate.Struct(ub)
		if inputErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, inputErr.Error()))
			return
		}

		type closedEvent struct {
			StoryID string `json:"storyId"`
			Closed  string `json:"closed"`
		}
		sbue := closedEvent{
			StoryID: storyID,
			Closed:  ub.Closed,
		}
		updateEventJSON, updateEventErr := json.Marshal(sbue)
		if updateEventErr != nil {
			s.Failure(w, r, http.StatusInternalServerError, Errorf(EINVALID, updateEventErr.Error()))
			return
		}

		_, err := sb.APIEvent(ctx, storyboardID, sessionUserID, "update_story_closed", string(updateEventJSON))
		if err != nil {
			s.Logger.Ctx(ctx).Error("handle storyboard story closed update error",
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

type storyboardStoryLinkRequestBody struct {
	Link string `json:"link" validate:"required"`
}

// handleStoryboardStoryLinkUpdate handles updating a story link in a storyboard
//
//	@Summary		Storyboard Story Link Update
//	@Description	Updates a story link in a storyboard
//	@Param			storyboardId	path	string							true	"the storyboard ID"
//	@Param			storyId			path	string							true	"the story ID"
//	@Param			storyboard		body	storyboardStoryMoveRequestBody	false	"new story link"
//	@Tags			storyboard
//	@Produce		json
//	@Success		200	object	standardJsonResponse{}
//	@Success		403	object	standardJsonResponse{}
//	@Success		500	object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/storyboards/{storyboardId}/stories/{storyId}/link [put]
func (s *Service) handleStoryboardStoryLinkUpdate(sb *storyboard.Service) http.HandlerFunc {
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

		var ub = storyboardStoryLinkRequestBody{}
		jsonErr := json.Unmarshal(body, &ub)
		if jsonErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		inputErr := validate.Struct(ub)
		if inputErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, inputErr.Error()))
			return
		}

		type linkEvent struct {
			StoryID string `json:"storyId"`
			Link    string `json:"link"`
		}
		sbue := linkEvent{
			StoryID: storyID,
			Link:    ub.Link,
		}
		updateEventJSON, updateEventErr := json.Marshal(sbue)
		if updateEventErr != nil {
			s.Failure(w, r, http.StatusInternalServerError, Errorf(EINVALID, updateEventErr.Error()))
			return
		}

		_, err := sb.APIEvent(ctx, storyboardID, sessionUserID, "update_story_link", string(updateEventJSON))
		if err != nil {
			s.Logger.Ctx(ctx).Error("handle storyboard story link update error",
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

		_, err := sb.APIEvent(ctx, storyboardID, sessionUserID, "move_story", string(moveEventJSON))
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

// handleStoryboardStoryDelete handles deleting a story in a storyboard
//
//	@Summary		Storyboard Story delete
//	@Description	Deletes a story in a storyboard
//	@Param			storyboardId	path	string	true	"the storyboard ID"
//	@Param			storyId			path	string	true	"the story ID"
//	@Tags			storyboard
//	@Produce		json
//	@Success		200	object	standardJsonResponse{}
//	@Success		403	object	standardJsonResponse{}
//	@Success		500	object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/storyboards/{storyboardId}/stories/{storyId} [delete]
func (s *Service) handleStoryboardStoryDelete(sb *storyboard.Service) http.HandlerFunc {
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

		_, err := sb.APIEvent(ctx, storyboardID, sessionUserID, "delete_story", storyID)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handle storyboard story delete error",
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
