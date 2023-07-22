package http

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"

	"github.com/gorilla/mux"
)

type teamResponse struct {
	Team     *thunderdome.Team `json:"team"`
	TeamRole string            `json:"teamRole"`
}

// handleGetTeamByUser gets an team with user role
// @Summary      Get Team
// @Description  Get a team with user role
// @Tags         team
// @Produce      json
// @Param        teamId  path    string  true  "the team ID"
// @Success      200     object  standardJsonResponse{data=teamResponse}
// @Success      500     object  standardJsonResponse{}
// @Security     ApiKeyAuth
// @Router       /teams/{teamId} [get]
func (s *Service) handleGetTeamByUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		TeamID := vars["teamId"]
		TeamRole := r.Context().Value(contextKeyTeamRole).(string)

		Team, err := s.TeamDataSvc.TeamGet(r.Context(), TeamID)
		if err != nil {
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		result := &teamResponse{
			Team:     Team,
			TeamRole: TeamRole,
		}

		s.Success(w, r, http.StatusOK, result, nil)
	}
}

// handleGetTeamsByUser gets a list of teams the user is a part of
// @Summary      Get User Teams
// @Description  Get a list of teams the user is a part of
// @Tags         team
// @Produce      json
// @Param        userId  path    string  true  "the user ID"
// @Success      200     object  standardJsonResponse{data=[]thunderdome.Team}
// @Success      403     object  standardJsonResponse{}
// @Security     ApiKeyAuth
// @Router       /users/{userId}/teams [get]
func (s *Service) handleGetTeamsByUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		UserID := vars["userId"]

		Limit, Offset := getLimitOffsetFromRequest(r)

		Teams := s.TeamDataSvc.TeamListByUser(r.Context(), UserID, Limit, Offset)

		s.Success(w, r, http.StatusOK, Teams, nil)
	}
}

// handleGetTeamUsers gets a list of users associated to the team
// @Summary      Get Team users
// @Description  Get a list of users associated to the team
// @Tags         team
// @Produce      json
// @Param        teamId  path    string  true  "the team ID"
// @Success      200     object  standardJsonResponse{data=[]thunderdome.User}
// @Security     ApiKeyAuth
// @Router       /teams/{teamId}/users [get]
func (s *Service) handleGetTeamUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		TeamID := vars["teamId"]
		Limit, Offset := getLimitOffsetFromRequest(r)

		Users, UserCount, err := s.TeamDataSvc.TeamUserList(r.Context(), TeamID, Limit, Offset)
		if err != nil {
			s.Failure(w, r, http.StatusInternalServerError, err)
		}

		Meta := &pagination{
			Count:  UserCount,
			Offset: Offset,
			Limit:  Limit,
		}

		s.Success(w, r, http.StatusOK, Users, Meta)
	}
}

type teamCreateRequestBody struct {
	Name string `json:"name" validate:"required"`
}

// handleCreateTeam handles creating a team with current user as admin
// @Summary      Create Team
// @Description  Creates a team with the current user as the team admin
// @Tags         team
// @Produce      json
// @Param        userId  path    string                 true  "the user ID"
// @Param        team    body    teamCreateRequestBody  true  "new team object"
// @Success      200     object  standardJsonResponse{data=thunderdome.Team}
// @Success      403     object  standardJsonResponse{}
// @Success      500     object  standardJsonResponse{}
// @Security     ApiKeyAuth
// @Router       /users/{userId}/teams [post]
func (s *Service) handleCreateTeam() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		UserID := vars["userId"]

		var team = teamCreateRequestBody{}
		body, bodyErr := io.ReadAll(r.Body)
		if bodyErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
			return
		}

		jsonErr := json.Unmarshal(body, &team)
		if jsonErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		inputErr := validate.Struct(team)
		if inputErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, inputErr.Error()))
		}

		NewTeam, err := s.TeamDataSvc.TeamCreate(r.Context(), UserID, team.Name)
		if err != nil {
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, NewTeam, nil)
	}
}

type teamAddUserRequestBody struct {
	Email string `json:"email" validate:"required,email"`
	Role  string `json:"role" enums:"MEMBER,ADMIN" validate:"required,oneof=MEMBER ADMIN"`
}

// handleTeamAddUser handles adding user to a team
// @Summary      Add Team User
// @Description  Adds a user to the team
// @Tags         team
// @Produce      json
// @Param        teamId  path    string                  true  "the team ID"
// @Param        user    body    teamAddUserRequestBody  true  "new team user object"
// @Success      200     object  standardJsonResponse{}
// @Success      403     object  standardJsonResponse{}
// @Success      500     object  standardJsonResponse{}
// @Security     ApiKeyAuth
// @Router       /teams/{teamId}/users [post]
func (s *Service) handleTeamAddUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		TeamID := vars["teamId"]

		var u = teamAddUserRequestBody{}
		body, bodyErr := io.ReadAll(r.Body)
		if bodyErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
			return
		}

		jsonErr := json.Unmarshal(body, &u)
		if jsonErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		inputErr := validate.Struct(u)
		if inputErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, inputErr.Error()))
		}

		UserEmail := u.Email

		User, UserErr := s.UserDataSvc.GetUserByEmail(r.Context(), UserEmail)
		if UserErr != nil {
			s.Failure(w, r, http.StatusInternalServerError, Errorf(ENOTFOUND, "USER_NOT_FOUND"))
			return
		}

		_, err := s.TeamDataSvc.TeamAddUser(r.Context(), TeamID, User.Id, u.Role)
		if err != nil {
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, nil, nil)
	}
}

// handleTeamRemoveUser handles removing user from a team
// @Summary      Remove Team User
// @Description  Remove a user from the team
// @Tags         team
// @Produce      json
// @Param        teamId  path    string  true  "the team ID"
// @Param        userId  path    string  true  "the user ID"
// @Success      200     object  standardJsonResponse{}
// @Success      403     object  standardJsonResponse{}
// @Success      500     object  standardJsonResponse{}
// @Security     ApiKeyAuth
// @Router       /teams/{teamId}/users/{userId} [delete]
func (s *Service) handleTeamRemoveUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		TeamID := vars["teamId"]
		UserID := vars["userId"]
		idErr := validate.Var(UserID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		err := s.TeamDataSvc.TeamRemoveUser(r.Context(), TeamID, UserID)
		if err != nil {
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, nil, nil)
	}
}

// handleGetTeamBattles gets a list of battles associated to the team
// @Summary      Get Team Battles
// @Description  Get a list of battles associated to the team
// @Tags         team
// @Produce      json
// @Param        teamId  path    string  true  "the team ID"
// @Success      200     object  standardJsonResponse{data=[]thunderdome.Poker}
// @Security     ApiKeyAuth
// @Router       /teams/{teamId}/battles [get]
func (s *Service) handleGetTeamBattles() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		TeamID := vars["teamId"]

		Limit, Offset := getLimitOffsetFromRequest(r)

		Battles := s.TeamDataSvc.TeamPokerList(r.Context(), TeamID, Limit, Offset)

		s.Success(w, r, http.StatusOK, Battles, nil)
	}
}

// handleTeamRemoveBattle handles removing battle from a team
// @Summary      Remove Team Poker
// @Description  Remove a battle from the team
// @Tags         team
// @Produce      json
// @Param        teamId    path    string  true  "the team ID"
// @Param        battleId  path    string  true  "the battle ID"
// @Success      200       object  standardJsonResponse{}
// @Success      403       object  standardJsonResponse{}
// @Success      500       object  standardJsonResponse{}
// @Security     ApiKeyAuth
// @Router       /teams/{teamId}/battles/{battleId} [delete]
func (s *Service) handleTeamRemoveBattle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		TeamID := vars["teamId"]
		BattleID := vars["battleId"]
		idErr := validate.Var(BattleID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		err := s.TeamDataSvc.TeamRemovePoker(r.Context(), TeamID, BattleID)
		if err != nil {
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, nil, nil)
	}
}

// handleDeleteTeam handles deleting a team
// @Summary      Delete Team
// @Description  Delete a Team
// @Tags         team
// @Produce      json
// @Param        teamId  path    string  true  "the team ID"
// @Success      200     object  standardJsonResponse{}
// @Success      403     object  standardJsonResponse{}
// @Success      500     object  standardJsonResponse{}
// @Security     ApiKeyAuth
// @Router       /teams/{teamId} [delete]
func (s *Service) handleDeleteTeam() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		TeamID := vars["teamId"]
		idErr := validate.Var(TeamID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		err := s.TeamDataSvc.TeamDelete(r.Context(), TeamID)
		if err != nil {
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, nil, nil)
	}
}

// handleGetTeamRetros gets a list of retros associated to the team
// @Summary      Get Team Retros
// @Description  Get a list of retros associated to the team
// @Tags         team
// @Produce      json
// @Param        teamId  path    string  true  "the team ID"
// @Success      200     object  standardJsonResponse{data=[]thunderdome.Retro}
// @Security     ApiKeyAuth
// @Router       /teams/{teamId}/retros [get]
func (s *Service) handleGetTeamRetros() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		TeamID := vars["teamId"]
		Limit, Offset := getLimitOffsetFromRequest(r)

		Retrospectives := s.TeamDataSvc.TeamRetroList(r.Context(), TeamID, Limit, Offset)

		s.Success(w, r, http.StatusOK, Retrospectives, nil)
	}
}

// handleTeamRemoveRetro handles removing retro from a team
// @Summary      Remove Team Retro
// @Description  Remove a retro from the team
// @Tags         team
// @Produce      json
// @Param        teamId   path    string  true  "the team ID"
// @Param        retroId  path    string  true  "the retro ID"
// @Success      200      object  standardJsonResponse{}
// @Success      403      object  standardJsonResponse{}
// @Success      500      object  standardJsonResponse{}
// @Security     ApiKeyAuth
// @Router       /teams/{teamId}/retros/{retroId} [delete]
func (s *Service) handleTeamRemoveRetro() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		TeamID := vars["teamId"]
		RetrospectiveID := vars["retroId"]
		idErr := validate.Var(RetrospectiveID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		err := s.TeamDataSvc.TeamRemoveRetro(r.Context(), TeamID, RetrospectiveID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

// handleGetTeamStoryboards gets a list of storyboards associated to the team
// @Summary      Get Team Storyboards
// @Description  Get a list of storyboards associated to the team
// @Tags         team
// @Produce      json
// @Param        teamId  path    string  true  "the team ID"
// @Success      200     object  standardJsonResponse{data=[]thunderdome.Storyboard}
// @Security     ApiKeyAuth
// @Router       /teams/{teamId}/storyboards [get]
func (s *Service) handleGetTeamStoryboards() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		TeamID := vars["teamId"]
		Limit, Offset := getLimitOffsetFromRequest(r)

		Storyboards := s.TeamDataSvc.TeamStoryboardList(r.Context(), TeamID, Limit, Offset)

		s.Success(w, r, http.StatusOK, Storyboards, nil)
	}
}

// handleTeamRemoveStoryboard handles removing storyboard from a team
// @Summary      Remove Team Storyboard
// @Description  Remove a storyboard from the team
// @Tags         team
// @Produce      json
// @Param        teamId        path    string  true  "the team ID"
// @Param        storyboardId  path    string  true  "the storyboard ID"
// @Success      200           object  standardJsonResponse{}
// @Success      403           object  standardJsonResponse{}
// @Success      500           object  standardJsonResponse{}
// @Security     ApiKeyAuth
// @Router       /teams/{teamId}/storyboards/{storyboardId} [delete]
func (s *Service) handleTeamRemoveStoryboard() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		TeamID := vars["teamId"]
		StoryboardID := vars["storyboardId"]
		idErr := validate.Var(StoryboardID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		err := s.TeamDataSvc.TeamRemoveStoryboard(r.Context(), TeamID, StoryboardID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

// handleGetTeamRetroActions gets a list of retro actions
// @Summary      Get Retro Actions
// @Description  get list of retro actions
// @Tags         team
// @Produce      json
// @Param        limit      query   int      false  "Max number of results to return"
// @Param        offset     query   int      false  "Starting point to return rows from, should be multiplied by limit or 0"
// @Param        completed  query   boolean  false  "Only completed retro actions"
// @Success      200        object  standardJsonResponse{data=[]thunderdome.RetroAction}
// @Failure      500        object  standardJsonResponse{}
// @Security     ApiKeyAuth
// @Router       /teams/{teamId}/retro-actions [get]
func (s *Service) handleGetTeamRetroActions() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		TeamID := vars["teamId"]
		Limit, Offset := getLimitOffsetFromRequest(r)
		var err error
		var Count int
		var Actions []*thunderdome.RetroAction
		query := r.URL.Query()
		Completed, _ := strconv.ParseBool(query.Get("completed"))

		Actions, Count, err = s.RetroDataSvc.GetTeamRetroActions(TeamID, Limit, Offset, Completed)

		if err != nil {
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		Meta := &pagination{
			Count:  Count,
			Offset: Offset,
			Limit:  Limit,
		}

		s.Success(w, r, http.StatusOK, Actions, Meta)
	}
}
