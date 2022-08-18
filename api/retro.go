package api

import (
	"encoding/json"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"strconv"

	"github.com/StevenWeathers/thunderdome-planning-poker/api/retro"

	"github.com/StevenWeathers/thunderdome-planning-poker/model"
	"github.com/gorilla/mux"
)

type retroCreateRequestBody struct {
	RetroName            string `json:"retroName" example:"sprint 10 retro" validate:"required"`
	Format               string `json:"format" example:"worked_improve_question" validate:"required,oneof=worked_improve_question"`
	JoinCode             string `json:"joinCode" example:"iammadmax"`
	FacilitatorCode      string `json:"facilitatorCode" example:"likeaboss"`
	MaxVotes             int    `json:"maxVotes" validate:"required,min=1,max=9"`
	BrainstormVisibility string `json:"brainstormVisibility" validate:"required,oneof=visible concealed hidden"`
}

// handleRetroCreate handles creating a retro
// @Summary Create Retro
// @Description Create a retro associated to the user
// @Tags retro
// @Produce  json
// @Param userId path string true "the user ID"
// @Param orgId path string false "the organization ID"
// @Param departmentId path string false "the department ID"
// @Param teamId path string false "the team ID"
// @Param retro body retroCreateRequestBody false "new retro object"
// @Success 200 object standardJsonResponse{data=model.Retro}
// @Failure 403 object standardJsonResponse{}
// @Failure 500 object standardJsonResponse{}
// @Security ApiKeyAuth
// @Router /users/{userId}/retros [post]
// @Router /teams/{teamId}/users/{userId}/retros [post]
// @Router /{orgId}/teams/{teamId}/users/{userId}/retros [post]
// @Router /{orgId}/departments/{departmentId}/teams/{teamId}/users/{userId}/retros [post]
func (a *api) handleRetroCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		vars := mux.Vars(r)
		UserID := vars["userId"]
		TeamID, teamIdExists := vars["teamId"]

		if !teamIdExists && viper.GetBool("config.require_teams") {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, "RETRO_CREATION_REQUIRES_TEAM"))
			return
		}

		body, bodyErr := io.ReadAll(r.Body)
		if bodyErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
			return
		}

		var nr = retroCreateRequestBody{}
		jsonErr := json.Unmarshal(body, &nr)
		if jsonErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		inputErr := validate.Struct(nr)
		if inputErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, inputErr.Error()))
			return
		}

		var newRetro *model.Retro
		var err error
		// if retro created with team association

		if teamIdExists {
			if isTeamUserOrAnAdmin(r) {
				newRetro, err = a.db.TeamRetroCreate(ctx, TeamID, UserID, nr.RetroName, nr.Format, nr.JoinCode, nr.FacilitatorCode, nr.MaxVotes, nr.BrainstormVisibility)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
			} else {
				a.Failure(w, r, http.StatusForbidden, Errorf(EUNAUTHORIZED, "REQUIRES_TEAM_USER"))
				return
			}
		} else {
			newRetro, err = a.db.RetroCreate(UserID, nr.RetroName, nr.Format, nr.JoinCode, nr.FacilitatorCode, nr.MaxVotes, nr.BrainstormVisibility)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		a.Success(w, r, http.StatusOK, newRetro, nil)
	}
}

// handleRetroGet looks up retro or returns notfound status
// @Summary Get Retro
// @Description get retro by ID
// @Tags retro
// @Produce  json
// @Param retroId path string true "the retro ID to get"
// @Success 200 object standardJsonResponse{data=model.Retro}
// @Failure 403 object standardJsonResponse{}
// @Failure 404 object standardJsonResponse{}
// @Security ApiKeyAuth
// @Router /retros/{retroId} [get]
func (a *api) handleRetroGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		RetroID := vars["retroId"]
		UserID := r.Context().Value(contextKeyUserID).(string)

		re, err := a.db.RetroGet(RetroID, UserID)

		if err != nil {
			http.NotFound(w, r)
			return
		}

		a.Success(w, r, http.StatusOK, re, nil)
	}
}

// handleRetrosGetByUser looks up retros associated with userID
// @Summary Get Retros by User
// @Description get list of retros for the user
// @Tags retro
// @Produce  json
// @Param userId path string true "the user ID to get retros for"
// @Param limit query int false "Max number of results to return"
// @Param offset query int false "Starting point to return rows from, should be multiplied by limit or 0"
// @Success 200 object standardJsonResponse{data=[]model.Retro}
// @Failure 403 object standardJsonResponse{}
// @Failure 404 object standardJsonResponse{}
// @Security ApiKeyAuth
// @Router /users/{userId}/retros [get]
func (a *api) handleRetrosGetByUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		UserID := vars["userId"]

		retros, err := a.db.RetroGetByUser(UserID)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		a.Success(w, r, http.StatusOK, retros, nil)
	}
}

// handleGetRetros gets a list of retros
// @Summary Get Retros
// @Description get list of retros
// @Tags retro
// @Produce  json
// @Param limit query int false "Max number of results to return"
// @Param offset query int false "Starting point to return rows from, should be multiplied by limit or 0"
// @Param active query boolean false "Only active retros"
// @Success 200 object standardJsonResponse{data=[]model.Retro}
// @Failure 500 object standardJsonResponse{}
// @Security ApiKeyAuth
// @Router /retros [get]
func (a *api) handleGetRetros() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		Limit, Offset := getLimitOffsetFromRequest(r)
		query := r.URL.Query()
		var err error
		var Count int
		var Retros []*model.Retro
		Active, _ := strconv.ParseBool(query.Get("active"))

		if Active {
			Retros, Count, err = a.db.GetActiveRetros(Limit, Offset)
		} else {
			Retros, Count, err = a.db.GetRetros(Limit, Offset)
		}

		if err != nil {
			a.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		Meta := &pagination{
			Count:  Count,
			Offset: Offset,
			Limit:  Limit,
		}

		a.Success(w, r, http.StatusOK, Retros, Meta)
	}
}

type actionUpdateRequestBody struct {
	ActionID  string `json:"id" swaggerignore:"true" validate:"required,uuid"`
	Completed bool   `json:"completed" example:"false"`
	Content   string `json:"content" example:"update documentation" validate:"required"`
}

// handleRetroActionUpdate handles updating a retro action item
// @Summary Retro Action Item Update
// @Description Update a retro action item
// @Param retroId path string true "the retro ID"
// @Param actionId path string true "the action ID"
// @Param actionItem body actionUpdateRequestBody true "updated action item"
// @Tags retro
// @Produce  json
// @Success 200 object standardJsonResponse{}
// @Success 403 object standardJsonResponse{}
// @Success 500 object standardJsonResponse{}
// @Security ApiKeyAuth
// @Router /retros/{retroId}/actions/{actionId} [put]
func (a *api) handleRetroActionUpdate(rs *retro.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var ra = actionUpdateRequestBody{}

		vars := mux.Vars(r)
		RetroID := vars["retroId"]
		idErr := validate.Var(RetroID, "required,uuid")
		if idErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		ActionID := vars["actionId"]
		UserID := r.Context().Value(contextKeyUserID).(string)

		body, bodyErr := io.ReadAll(r.Body)
		if bodyErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
			return
		}
		jsonErr := json.Unmarshal(body, &ra)
		if jsonErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		ra.ActionID = ActionID
		inputErr := validate.Struct(ra)
		if inputErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, inputErr.Error()))
			return
		}
		updatedActionJson, _ := json.Marshal(ra)

		err := rs.APIEvent(r.Context(), RetroID, UserID, "update_action", string(updatedActionJson))
		if err != nil {
			a.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		a.Success(w, r, http.StatusOK, nil, nil)
	}
}

// handleRetroActionDelete handles deleting a retro action item
// @Summary Retro Action Item Delete
// @Description Delete a retro action item
// @Param retroId path string true "the retro ID"
// @Param actionId path string true "the action ID"
// @Tags retro
// @Produce  json
// @Success 200 object standardJsonResponse{}
// @Success 403 object standardJsonResponse{}
// @Success 500 object standardJsonResponse{}
// @Security ApiKeyAuth
// @Router /retros/{retroId}/actions/{actionId} [delete]
func (a *api) handleRetroActionDelete(rs *retro.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		RetroID := vars["retroId"]
		idErr := validate.Var(RetroID, "required,uuid")
		if idErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		ActionID := vars["actionId"]
		idErr = validate.Var(ActionID, "required,uuid")
		if idErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		UserID := r.Context().Value(contextKeyUserID).(string)

		type actionItem struct {
			ActionID string `json:"id"`
		}
		deleteItem, _ := json.Marshal(actionItem{ActionID: ActionID})

		err := rs.APIEvent(r.Context(), RetroID, UserID, "delete_action", string(deleteItem))
		if err != nil {
			a.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		a.Success(w, r, http.StatusOK, nil, nil)
	}
}

type actionCommentRequestBody struct {
	Comment string `json:"comment" validate:"required"`
}

// handleRetroActionCommentAdd handles adding a comment to a retro action item
// @Summary Retro Action Item Comment
// @Description Add a comment to a retro action item
// @Param retroId path string true "the retro ID"
// @Param actionId path string true "the action ID"
// @Param actionItem body actionCommentRequestBody true "action comment"
// @Tags retro
// @Produce  json
// @Success 200 object standardJsonResponse{}
// @Success 403 object standardJsonResponse{}
// @Success 500 object standardJsonResponse{}
// @Security ApiKeyAuth
// @Router /retros/{retroId}/actions/{actionId}/comments [post]
func (a *api) handleRetroActionCommentAdd() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var ra = actionCommentRequestBody{}

		vars := mux.Vars(r)
		RetroID := vars["retroId"]
		idErr := validate.Var(RetroID, "required,uuid")
		if idErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		ActionID := vars["actionId"]
		idErr = validate.Var(ActionID, "required,uuid")
		if idErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		UserID := r.Context().Value(contextKeyUserID).(string)

		body, bodyErr := io.ReadAll(r.Body)
		if bodyErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
			return
		}
		jsonErr := json.Unmarshal(body, &ra)
		if jsonErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		inputErr := validate.Struct(ra)
		if inputErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, inputErr.Error()))
			return
		}

		action, err := a.db.RetroActionCommentAdd(RetroID, ActionID, UserID, ra.Comment)
		if err != nil {
			a.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		a.Success(w, r, http.StatusOK, action, nil)
	}
}

// handleRetroActionCommentEdit handles editing a retro action item comment
// @Summary Retro Action Item Comment Edit
// @Description Edit a retro action item comment
// @Param retroId path string true "the retro ID"
// @Param actionId path string true "the action ID"
// @Param commentId path string true "the comment ID"
// @Param actionItem body actionCommentRequestBody true "action comment"
// @Tags retro
// @Produce  json
// @Success 200 object standardJsonResponse{}
// @Success 403 object standardJsonResponse{}
// @Success 500 object standardJsonResponse{}
// @Security ApiKeyAuth
// @Router /retros/{retroId}/actions/{actionId}/comments/{commentId} [put]
func (a *api) handleRetroActionCommentEdit() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var ra = actionCommentRequestBody{}

		vars := mux.Vars(r)
		RetroID := vars["retroId"]
		idErr := validate.Var(RetroID, "required,uuid")
		if idErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		ActionID := vars["actionId"]
		idErr = validate.Var(ActionID, "required,uuid")
		if idErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		CommentID := vars["commentId"]
		idErr = validate.Var(CommentID, "required,uuid")
		if idErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		body, bodyErr := io.ReadAll(r.Body)
		if bodyErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
			return
		}
		jsonErr := json.Unmarshal(body, &ra)
		if jsonErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		inputErr := validate.Struct(ra)
		if inputErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, inputErr.Error()))
			return
		}

		action, err := a.db.RetroActionCommentEdit(RetroID, ActionID, CommentID, ra.Comment)
		if err != nil {
			a.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		a.Success(w, r, http.StatusOK, action, nil)
	}
}

// handleRetroActionCommentDelete handles delete a comment from a retro action item
// @Summary Retro Action Item Comment Delete
// @Description Delete a comment from a retro action item
// @Param retroId path string true "the retro ID"
// @Param actionId path string true "the action ID"
// @Param commentId path string true "the comment ID"
// @Tags retro
// @Produce  json
// @Success 200 object standardJsonResponse{}
// @Success 403 object standardJsonResponse{}
// @Success 500 object standardJsonResponse{}
// @Security ApiKeyAuth
// @Router /retros/{retroId}/actions/{actionId}/comments/{commentId} [post]
func (a *api) handleRetroActionCommentDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		RetroID := vars["retroId"]
		idErr := validate.Var(RetroID, "required,uuid")
		if idErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		ActionID := vars["actionId"]
		idErr = validate.Var(ActionID, "required,uuid")
		if idErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		CommentID := vars["commentId"]
		idErr = validate.Var(CommentID, "required,uuid")
		if idErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		action, err := a.db.RetroActionCommentDelete(RetroID, ActionID, CommentID)
		if err != nil {
			a.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		a.Success(w, r, http.StatusOK, action, nil)
	}
}
