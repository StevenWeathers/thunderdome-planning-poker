package http

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"go.uber.org/zap"

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
		ctx := r.Context()
		SessionUserID := ctx.Value(contextKeyUserID).(string)
		vars := mux.Vars(r)
		TeamID := vars["teamId"]
		TeamRole := r.Context().Value(contextKeyTeamRole).(string)

		Team, err := s.TeamDataSvc.TeamGet(ctx, TeamID)
		if err != nil {
			s.Logger.Ctx(ctx).Error(
				"handleGetTeamByUser error", zap.Error(err), zap.String("team_id", TeamID),
				zap.String("session_user_id", SessionUserID), zap.String("team_role", TeamRole))
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
		ctx := r.Context()
		vars := mux.Vars(r)
		UserID := vars["userId"]

		Limit, Offset := getLimitOffsetFromRequest(r)

		Teams := s.TeamDataSvc.TeamListByUser(ctx, UserID, Limit, Offset)

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
		ctx := r.Context()
		SessionUserID := ctx.Value(contextKeyUserID).(string)
		vars := mux.Vars(r)
		TeamID := vars["teamId"]
		Limit, Offset := getLimitOffsetFromRequest(r)

		Users, UserCount, err := s.TeamDataSvc.TeamUserList(ctx, TeamID, Limit, Offset)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleGetTeamUsers error", zap.Error(err), zap.String("team_id", TeamID),
				zap.Int("limit", Limit), zap.Int("offset", Offset), zap.String("session_user_id", SessionUserID))
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
		ctx := r.Context()
		SessionUserID := ctx.Value(contextKeyUserID).(string)
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

		NewTeam, err := s.TeamDataSvc.TeamCreate(ctx, UserID, team.Name)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleCreateTeam error", zap.Error(err), zap.String("entity_user_id", UserID),
				zap.String("session_user_id", SessionUserID))
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
		ctx := r.Context()
		SessionUserID := ctx.Value(contextKeyUserID).(string)
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

		UserEmail := strings.ToLower(u.Email)

		User, UserErr := s.UserDataSvc.GetUserByEmail(ctx, UserEmail)
		if UserErr != nil && errors.Is(UserErr, sql.ErrNoRows) {
			inviteID, inviteErr := s.TeamDataSvc.TeamInviteUser(ctx, TeamID, UserEmail, u.Role)
			if inviteErr != nil {
				s.Logger.Ctx(ctx).Error("handleTeamAddUser error", zap.Error(inviteErr),
					zap.String("team_id", TeamID), zap.String("session_user_id", SessionUserID))
				s.Failure(w, r, http.StatusInternalServerError, UserErr)
				return
			}
			team, teamErr := s.TeamDataSvc.TeamGet(ctx, TeamID)
			if teamErr != nil {
				s.Logger.Ctx(ctx).Error("handleTeamAddUser error", zap.Error(teamErr),
					zap.String("team_id", TeamID), zap.String("session_user_id", SessionUserID))
				s.Failure(w, r, http.StatusInternalServerError, teamErr)
				return
			}
			emailErr := s.Email.SendTeamInvite(team.Name, UserEmail, inviteID)
			if emailErr != nil {
				s.Logger.Ctx(ctx).Error("handleTeamAddUser error", zap.Error(emailErr),
					zap.String("team_id", TeamID), zap.String("session_user_id", SessionUserID))
				s.Failure(w, r, http.StatusInternalServerError, emailErr)
				return
			}
			s.Success(w, r, http.StatusOK, nil, nil)
			return
		} else if UserErr != nil {
			s.Logger.Ctx(ctx).Error("handleTeamAddUser error", zap.Error(UserErr),
				zap.String("team_id", TeamID), zap.String("session_user_id", SessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, UserErr)
			return
		}

		_, err := s.TeamDataSvc.TeamAddUser(ctx, TeamID, User.Id, u.Role)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleTeamAddUser error", zap.Error(err), zap.String("team_id", TeamID),
				zap.String("user_id", User.Id), zap.String("team_role", u.Role),
				zap.String("session_user_id", SessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, nil, nil)
	}
}

type teamUpdateUserRequestBody struct {
	Role string `json:"role" enums:"MEMBER,ADMIN" validate:"required,oneof=MEMBER ADMIN"`
}

// handleTeamUpdateUser handles updating a user on the team
// @Summary      Update Team User
// @Description  Updates a team user
// @Tags         team
// @Produce      json
// @Param        teamId  path    string                  true  "the team ID"
// @Param        userId  path    string                  true  "the user ID"
// @Param        user    body    teamUpdateUserRequestBody  true  "updated team user object"
// @Success      200     object  standardJsonResponse{}
// @Success      403     object  standardJsonResponse{}
// @Success      500     object  standardJsonResponse{}
// @Security     ApiKeyAuth
// @Router       /teams/{teamId}/users/{userId} [put]
func (s *Service) handleTeamUpdateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		SessionUserID := ctx.Value(contextKeyUserID).(string)
		vars := mux.Vars(r)
		TeamID := vars["teamId"]
		UserID := vars["userId"]

		var u = teamUpdateUserRequestBody{}
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

		_, err := s.TeamDataSvc.TeamUpdateUser(ctx, TeamID, UserID, u.Role)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleTeamAddUser error", zap.Error(err), zap.String("team_id", TeamID),
				zap.String("user_id", UserID), zap.String("team_role", u.Role),
				zap.String("session_user_id", SessionUserID))
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
		ctx := r.Context()
		SessionUserID := ctx.Value(contextKeyUserID).(string)
		vars := mux.Vars(r)
		TeamID := vars["teamId"]
		UserID := vars["userId"]
		idErr := validate.Var(UserID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		err := s.TeamDataSvc.TeamRemoveUser(ctx, TeamID, UserID)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleTeamRemoveUser error", zap.Error(err), zap.String("team_id", TeamID),
				zap.String("user_id", UserID), zap.String("session_user_id", SessionUserID))
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
		ctx := r.Context()
		vars := mux.Vars(r)
		TeamID := vars["teamId"]

		Limit, Offset := getLimitOffsetFromRequest(r)

		Battles := s.TeamDataSvc.TeamPokerList(ctx, TeamID, Limit, Offset)

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
		ctx := r.Context()
		SessionUserID := ctx.Value(contextKeyUserID).(string)
		vars := mux.Vars(r)
		TeamID := vars["teamId"]
		BattleID := vars["battleId"]
		idErr := validate.Var(BattleID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		err := s.TeamDataSvc.TeamRemovePoker(ctx, TeamID, BattleID)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleTeamRemoveBattle error", zap.Error(err), zap.String("team_id", TeamID),
				zap.String("battle_id", BattleID), zap.String("session_user_id", SessionUserID))
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
		ctx := r.Context()
		SessionUserID := ctx.Value(contextKeyUserID).(string)
		vars := mux.Vars(r)
		TeamID := vars["teamId"]
		idErr := validate.Var(TeamID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		err := s.TeamDataSvc.TeamDelete(ctx, TeamID)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleDeleteTeam error", zap.Error(err), zap.String("team_id", TeamID),
				zap.String("session_user_id", SessionUserID))
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
		ctx := r.Context()
		vars := mux.Vars(r)
		TeamID := vars["teamId"]
		Limit, Offset := getLimitOffsetFromRequest(r)

		Retrospectives := s.TeamDataSvc.TeamRetroList(ctx, TeamID, Limit, Offset)

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
		ctx := r.Context()
		SessionUserID := ctx.Value(contextKeyUserID).(string)
		vars := mux.Vars(r)
		TeamID := vars["teamId"]
		RetrospectiveID := vars["retroId"]
		idErr := validate.Var(RetrospectiveID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		err := s.TeamDataSvc.TeamRemoveRetro(ctx, TeamID, RetrospectiveID)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleTeamRemoveRetro error", zap.Error(err), zap.String("team_id", TeamID),
				zap.String("retro_id", RetrospectiveID), zap.String("session_user_id", SessionUserID))
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
		ctx := r.Context()
		vars := mux.Vars(r)
		TeamID := vars["teamId"]
		Limit, Offset := getLimitOffsetFromRequest(r)

		Storyboards := s.TeamDataSvc.TeamStoryboardList(ctx, TeamID, Limit, Offset)

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
		ctx := r.Context()
		SessionUserID := ctx.Value(contextKeyUserID).(string)
		vars := mux.Vars(r)
		TeamID := vars["teamId"]
		StoryboardID := vars["storyboardId"]
		idErr := validate.Var(StoryboardID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		err := s.TeamDataSvc.TeamRemoveStoryboard(ctx, TeamID, StoryboardID)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleTeamRemoveStoryboard error", zap.Error(err), zap.String("team_id", TeamID),
				zap.String("storyboard_id", StoryboardID), zap.String("session_user_id", SessionUserID))
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
		ctx := r.Context()
		SessionUserID := ctx.Value(contextKeyUserID).(string)
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
			s.Logger.Ctx(ctx).Error("handleGetTeamRetroActions error", zap.Error(err), zap.String("team_id", TeamID),
				zap.Int("limit", Limit), zap.Int("offset", Offset), zap.String("session_user_id", SessionUserID))
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
