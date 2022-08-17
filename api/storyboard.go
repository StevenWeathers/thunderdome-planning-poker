package api

import (
	"encoding/json"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"strconv"

	"github.com/StevenWeathers/thunderdome-planning-poker/model"
	"github.com/gorilla/mux"
)

type storyboardCreateRequestBody struct {
	StoryboardName  string `json:"storyboardName" validate:"required"`
	JoinCode        string `json:"joinCode"`
	FacilitatorCode string `json:"facilitatorCode"`
}

// handleStoryboardCreate handles creating a storyboard (arena)
// @Summary Create Storyboard
// @Description Create a storyboard associated to the user
// @Tags storyboard
// @Produce  json
// @Param userId path string true "the user ID"
// @Param orgId path string false "the organization ID"
// @Param departmentId path string false "the department ID"
// @Param teamId path string false "the team ID"
// @Param storyboard body storyboardCreateRequestBody false "new storyboard object"
// @Success 200 object standardJsonResponse{data=model.Storyboard}
// @Failure 403 object standardJsonResponse{}
// @Failure 500 object standardJsonResponse{}
// @Security ApiKeyAuth
// @Router /users/{userId}/storyboards [post]
// @Router /teams/{teamId}/users/{userId}/storyboards [post]
// @Router /{orgId}/teams/{teamId}/users/{userId}/storyboards [post]
// @Router /{orgId}/departments/{departmentId}/teams/{teamId}/users/{userId}/storyboards [post]
func (a *api) handleStoryboardCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		vars := mux.Vars(r)
		UserID := vars["userId"]
		TeamID, teamIdExists := vars["teamId"]

		if !teamIdExists && viper.GetBool("config.require_teams") {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, "STORYBOARD_CREATION_REQUIRES_TEAM"))
			return
		}

		body, bodyErr := io.ReadAll(r.Body) // check for errors
		if bodyErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
			return
		}

		var s = storyboardCreateRequestBody{}
		jsonErr := json.Unmarshal(body, &s)
		if jsonErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		inputErr := validate.Struct(s)
		if inputErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, inputErr.Error()))
			return
		}

		var newStoryboard *model.Storyboard
		var err error
		// if storyboard created with team association
		if teamIdExists {
			if isTeamUserOrAnAdmin(r) {
				newStoryboard, err = a.db.TeamCreateStoryboard(ctx, TeamID, UserID, s.StoryboardName, s.JoinCode, s.FacilitatorCode)

				if err != nil {
					a.Failure(w, r, http.StatusInternalServerError, err)
					return
				}
			} else {
				a.Failure(w, r, http.StatusForbidden, Errorf(EUNAUTHORIZED, "REQUIRES_TEAM_USER"))
				return
			}
		} else {
			newStoryboard, err = a.db.CreateStoryboard(ctx, UserID, s.StoryboardName, s.JoinCode, s.FacilitatorCode)
			if err != nil {
				a.Failure(w, r, http.StatusInternalServerError, err)
				return
			}
		}

		a.Success(w, r, http.StatusOK, newStoryboard, nil)
	}
}

// handleStoryboardGet gets the storyboard by ID
// @Summary Get Storyboard
// @Description get storyboard by ID
// @Tags storyboard
// @Produce  json
// @Param storyboardId path string true "the storyboard ID to get"
// @Success 200 object standardJsonResponse{data=model.Storyboard}
// @Failure 403 object standardJsonResponse{}
// @Failure 404 object standardJsonResponse{}
// @Security ApiKeyAuth
// @Router /storyboards/{storyboardId} [get]
func (a *api) handleStoryboardGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		StoryboardID := vars["storyboardId"]
		idErr := validate.Var(StoryboardID, "required,uuid")
		if idErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		UserId := r.Context().Value(contextKeyUserID).(string)
		UserType := r.Context().Value(contextKeyUserType).(string)

		storyboard, err := a.db.GetStoryboard(StoryboardID, UserId)
		if err != nil {
			a.Failure(w, r, http.StatusNotFound, Errorf(ENOTFOUND, "STORYBOARD_NOT_FOUND"))
			return
		}

		// don't allow retrieving storyboard details if storyboard has JoinCode and user hasn't joined yet
		if storyboard.JoinCode != "" {
			UserErr := a.db.GetStoryboardUserActiveStatus(StoryboardID, UserId)
			if UserErr != nil && UserType != adminUserType {
				a.Failure(w, r, http.StatusForbidden, Errorf(EUNAUTHORIZED, "USER_MUST_JOIN_STORYBOARD"))
				return
			}
		}

		a.Success(w, r, http.StatusOK, storyboard, nil)
	}
}

// handleGetUserStoryboards looks up storyboards associated with UserID
// @Summary Get Storyboards
// @Description get list of storyboards for the user
// @Tags storyboard
// @Produce  json
// @Param userId path string true "the user ID to get storyboards for"
// @Param limit query int false "Max number of results to return"
// @Param offset query int false "Starting point to return rows from, should be multiplied by limit or 0"
// @Success 200 object standardJsonResponse{data=[]model.Storyboard}
// @Failure 403 object standardJsonResponse{}
// @Failure 404 object standardJsonResponse{}
// @Security ApiKeyAuth
// @Router /users/{userId}/storyboards [get]
func (a *api) handleGetUserStoryboards() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		Limit, Offset := getLimitOffsetFromRequest(r)
		vars := mux.Vars(r)
		UserID := vars["userId"]

		storyboards, Count, err := a.db.GetStoryboardsByUser(UserID)
		if err != nil {
			a.Failure(w, r, http.StatusNotFound, Errorf(ENOTFOUND, "STORYBOARDS_NOT_FOUND"))
			return
		}

		Meta := &pagination{
			Count:  Count,
			Offset: Offset,
			Limit:  Limit,
		}

		a.Success(w, r, http.StatusOK, storyboards, Meta)
	}
}

// handleGetStoryboards gets a list of storyboards
// @Summary Get Storyboards
// @Description get list of storyboards
// @Tags storyboard
// @Produce  json
// @Param limit query int false "Max number of results to return"
// @Param offset query int false "Starting point to return rows from, should be multiplied by limit or 0"
// @Param active query boolean false "Only active storyboards"
// @Success 200 object standardJsonResponse{data=[]model.Storyboard}
// @Failure 500 object standardJsonResponse{}
// @Security ApiKeyAuth
// @Router /storyboards [get]
func (a *api) handleGetStoryboards() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		Limit, Offset := getLimitOffsetFromRequest(r)
		query := r.URL.Query()
		var err error
		var Count int
		var storyboards []*model.Storyboard
		Active, _ := strconv.ParseBool(query.Get("active"))

		if Active {
			storyboards, Count, err = a.db.GetActiveStoryboards(Limit, Offset)
		} else {
			storyboards, Count, err = a.db.GetStoryboards(Limit, Offset)
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

		a.Success(w, r, http.StatusOK, storyboards, Meta)
	}
}
