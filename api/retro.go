package api

import (
	"encoding/json"
	"github.com/StevenWeathers/thunderdome-planning-poker/api/retro"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/StevenWeathers/thunderdome-planning-poker/model"
	"github.com/gorilla/mux"
)

type retroCreateRequestBody struct {
	RetroName            string `json:"retroName" example:"sprint 10 retro"`
	Format               string `json:"format" example:"worked_improve_question"`
	JoinCode             string `json:"joinCode" example:"iammadmax"`
	MaxVotes             int    `json:"maxVotes"`
	BrainstormVisibility string `json:"brainstormVisibility"`
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
		userID := r.Context().Value(contextKeyUserID).(string)
		vars := mux.Vars(r)

		body, bodyErr := ioutil.ReadAll(r.Body)
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

		newRetro, err := a.db.RetroCreate(userID, nr.RetroName, nr.Format, nr.JoinCode, nr.MaxVotes, nr.BrainstormVisibility)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// if retro created with team association
		TeamID, ok := vars["teamId"]
		if ok {
			OrgRole := r.Context().Value(contextKeyOrgRole)
			DepartmentRole := r.Context().Value(contextKeyDepartmentRole)
			TeamRole := r.Context().Value(contextKeyTeamRole).(string)
			var isAdmin bool
			if DepartmentRole != nil && DepartmentRole.(string) == "ADMIN" {
				isAdmin = true
			}
			if OrgRole != nil && OrgRole.(string) == "ADMIN" {
				isAdmin = true
			}

			if isAdmin == true || TeamRole != "" {
				err := a.db.TeamAddRetro(TeamID, newRetro.Id)

				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
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

		re, err := a.db.RetroGet(RetroID)

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
		userID := r.Context().Value(contextKeyUserID).(string)

		retros, err := a.db.RetroGetByUser(userID)
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
	ActionID  string `json:"id" swaggerignore:"true"`
	Completed bool   `json:"completed" example:"false"`
	Content   string `json:"content" example:"update documentation"`
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
		ActionID := vars["actionId"]
		UserID := r.Context().Value(contextKeyUserID).(string)

		body, bodyErr := ioutil.ReadAll(r.Body)
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
		updatedActionJson, _ := json.Marshal(ra)

		err := rs.APIEvent(RetroID, UserID, "update_action", string(updatedActionJson))
		if err != nil {
			a.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		a.Success(w, r, http.StatusOK, nil, nil)
	}
}

type actionCommentRequestBody struct {
	Comment string `json:"comment"`
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
		ActionID := vars["actionId"]
		UserID := r.Context().Value(contextKeyUserID).(string)

		body, bodyErr := ioutil.ReadAll(r.Body)
		if bodyErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
			return
		}
		jsonErr := json.Unmarshal(body, &ra)
		if jsonErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
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
		ActionID := vars["actionId"]
		CommentID := vars["commentId"]

		action, err := a.db.RetroActionCommentDelete(RetroID, ActionID, CommentID)
		if err != nil {
			a.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		a.Success(w, r, http.StatusOK, action, nil)
	}
}
