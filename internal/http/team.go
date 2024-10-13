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

// handleGetTeamByUser gets a team with user role
//
//	@Summary		Get Team
//	@Description	Get a team with user role
//	@Tags			team
//	@Produce		json
//	@Param			teamId	path	string	true	"the team ID"
//	@Success		200		object	standardJsonResponse{data=teamResponse}
//	@Success		500		object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/teams/{teamId} [get]
func (s *Service) handleGetTeamByUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)
		userTeamRoles := ctx.Value(contextKeyUserTeamRoles).(*thunderdome.UserTeamRoleInfo)
		vars := mux.Vars(r)
		teamID := vars["teamId"]
		idErr := validate.Var(teamID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		// Workaround until v5 refactor to replace meta's with consistency
		teamRole := userTeamRoles.TeamRole
		if teamRole == nil {
			var adminStr string
			teamRole = &adminStr
		}

		team, err := s.TeamDataSvc.TeamGetByID(ctx, teamID)
		if err != nil {
			s.Logger.Ctx(ctx).Error(
				"handleGetTeamByUser error", zap.Error(err),
				zap.String("team_id", teamID),
				zap.String("session_user_id", sessionUserID),
				zap.String("team_role", *teamRole))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		result := &teamResponse{
			Team:     team,
			TeamRole: *teamRole,
		}

		s.Success(w, r, http.StatusOK, result, nil)
	}
}

// handleGetTeamsByUser gets a list of teams the user is a part of
//
//	@Summary		Get User Teams
//	@Description	Get a list of teams the user is a part of
//	@Tags			team
//	@Produce		json
//	@Param			userId	path	string	true	"the user ID"
//	@Success		200		object	standardJsonResponse{data=[]thunderdome.UserTeam}
//	@Success		403		object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/users/{userId}/teams [get]
func (s *Service) handleGetTeamsByUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		vars := mux.Vars(r)
		userID := vars["userId"]
		idErr := validate.Var(userID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		limit, offset := getLimitOffsetFromRequest(r)

		teams := s.TeamDataSvc.TeamListByUser(ctx, userID, limit, offset)

		s.Success(w, r, http.StatusOK, teams, nil)
	}
}

// handleGetTeamsByUser gets a list of teams the user is a part of that are not associated with an organization
//
//	@Summary		Get User Teams Non Org
//	@Description	Get a list of teams the user is a part of that are not associated with an organization
//	@Tags			team
//	@Produce		json
//	@Param			userId	path	string	true	"the user ID"
//	@Success		200		object	standardJsonResponse{data=[]thunderdome.UserTeam}
//	@Success		403		object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/users/{userId}/teams-non-org [get]
func (s *Service) handleGetTeamsByUserNonOrg() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		vars := mux.Vars(r)
		userID := vars["userId"]
		idErr := validate.Var(userID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		limit, offset := getLimitOffsetFromRequest(r)

		teams := s.TeamDataSvc.TeamListByUserNonOrg(ctx, userID, limit, offset)

		s.Success(w, r, http.StatusOK, teams, nil)
	}
}

// handleGetTeamUsers gets a list of users associated to the team
//
//	@Summary		Get Team users
//	@Description	Get a list of users associated to the team
//	@Tags			team
//	@Produce		json
//	@Param			teamId	path	string	true	"the team ID"
//	@Success		200		object	standardJsonResponse{data=[]thunderdome.User}
//	@Security		ApiKeyAuth
//	@Router			/teams/{teamId}/users [get]
func (s *Service) handleGetTeamUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)
		vars := mux.Vars(r)
		teamID := vars["teamId"]
		idErr := validate.Var(teamID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		limit, offset := getLimitOffsetFromRequest(r)

		users, userCount, err := s.TeamDataSvc.TeamUserList(ctx, teamID, limit, offset)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleGetTeamUsers error", zap.Error(err), zap.String("team_id", teamID),
				zap.Int("limit", limit), zap.Int("offset", offset), zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		meta := &pagination{
			Count:  userCount,
			Offset: offset,
			Limit:  limit,
		}

		s.Success(w, r, http.StatusOK, users, meta)
	}
}

type teamCreateRequestBody struct {
	Name string `json:"name" validate:"required"`
}

// handleCreateTeam handles creating a team with current user as admin
//
//	@Summary		Create Team
//	@Description	Creates a team with the current user as the team admin
//	@Tags			team
//	@Produce		json
//	@Param			userId	path	string					true	"the user ID"
//	@Param			team	body	teamCreateRequestBody	true	"new team object"
//	@Success		200		object	standardJsonResponse{data=thunderdome.Team}
//	@Success		403		object	standardJsonResponse{}
//	@Success		500		object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/users/{userId}/teams [post]
func (s *Service) handleCreateTeam() http.HandlerFunc {
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
			return
		}

		newTeam, err := s.TeamDataSvc.TeamCreate(ctx, userID, team.Name)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleCreateTeam error", zap.Error(err), zap.String("entity_user_id", userID),
				zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, newTeam, nil)
	}
}

// handleTeamUpdate handles updating a team
//
//	@Summary		Update Team
//	@Description	Updates a team
//	@Tags			team
//	@Produce		json
//	@Param			teamId	path	string					true	"the team ID"
//	@Param			team	body	teamCreateRequestBody	true	"updated team object"
//	@Success		200		object	standardJsonResponse{data=thunderdome.Team}
//	@Success		403		object	standardJsonResponse{}
//	@Success		500		object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/teams/{teamId} [put]
func (s *Service) handleTeamUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)
		vars := mux.Vars(r)
		teamID := vars["teamId"]
		idErr := validate.Var(teamID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

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
			return
		}

		newTeam, err := s.TeamDataSvc.TeamUpdate(ctx, teamID, team.Name)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleTeamUpdate error", zap.Error(err),
				zap.String("team_id", teamID),
				zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, newTeam, nil)
	}
}

type userAddMeta struct {
	Invited bool `json:"user_invited"`
	Added   bool `json:"user_added"`
}

type teamInviteUserRequestBody struct {
	Email string `json:"email" validate:"required,email"`
	Role  string `json:"role" enums:"MEMBER,ADMIN" validate:"required,oneof=MEMBER ADMIN"`
}

// handleTeamInviteUser handles inviting user to a team
//
//	@Summary		Invite Team User
//	@Description	Invites a user to the team
//	@Tags			team
//	@Produce		json
//	@Param			teamId	path	string						true	"the team ID"
//	@Param			user	body	teamInviteUserRequestBody	true	"new team user object"
//	@Success		200		object	standardJsonResponse{}
//	@Success		403		object	standardJsonResponse{}
//	@Success		500		object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/teams/{teamId}/invites [post]
func (s *Service) handleTeamInviteUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)
		vars := mux.Vars(r)
		teamID := vars["teamId"]
		idErr := validate.Var(teamID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		var u = teamInviteUserRequestBody{}
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
			return
		}

		userEmail := strings.ToLower(u.Email)

		if s.Config.LdapEnabled || s.Config.HeaderAuthEnabled {
			user, userErr := s.UserDataSvc.GetUserByEmail(ctx, userEmail)
			if userErr == nil {
				_, err := s.TeamDataSvc.TeamAddUser(ctx, teamID, user.ID, u.Role)
				if err != nil {
					s.Logger.Ctx(ctx).Error("handleTeamInviteUser error", zap.Error(err), zap.String("team_id", teamID),
						zap.String("user_id", user.ID), zap.String("team_role", u.Role),
						zap.String("session_user_id", sessionUserID))
					s.Failure(w, r, http.StatusInternalServerError, err)
					return
				}
				s.Success(w, r, http.StatusOK, nil, userAddMeta{Invited: false, Added: true})
				return
			} else if userErr != nil && !errors.Is(userErr, sql.ErrNoRows) {
				s.Logger.Ctx(ctx).Error("handleTeamInviteUser error", zap.Error(userErr),
					zap.String("team_id", teamID), zap.String("session_user_id", sessionUserID))
				s.Failure(w, r, http.StatusInternalServerError, userErr)
				return
			}
		}

		inviteID, inviteErr := s.TeamDataSvc.TeamInviteUser(ctx, teamID, userEmail, u.Role)
		if inviteErr != nil {
			s.Logger.Ctx(ctx).Error("handleTeamInviteUser error", zap.Error(inviteErr),
				zap.String("team_id", teamID), zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, inviteErr)
			return
		}
		team, teamErr := s.TeamDataSvc.TeamGetByID(ctx, teamID)
		if teamErr != nil {
			s.Logger.Ctx(ctx).Error("handleTeamInviteUser error", zap.Error(teamErr),
				zap.String("team_id", teamID), zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, teamErr)
			return
		}

		emailErr := s.Email.SendTeamInvite(team.Name, userEmail, inviteID)
		if emailErr != nil {
			s.Logger.Ctx(ctx).Error("handleTeamInviteUser error", zap.Error(emailErr),
				zap.String("team_id", teamID), zap.String("session_user_id", sessionUserID))
			return
		}

		s.Success(w, r, http.StatusOK, nil, userAddMeta{Invited: true, Added: false})
	}
}

type teamUpdateUserRequestBody struct {
	Role string `json:"role" enums:"MEMBER,ADMIN" validate:"required,oneof=MEMBER ADMIN"`
}

// handleTeamUpdateUser handles updating a user on the team
//
//	@Summary		Update Team User
//	@Description	Updates a team user
//	@Tags			team
//	@Produce		json
//	@Param			teamId	path	string						true	"the team ID"
//	@Param			userId	path	string						true	"the user ID"
//	@Param			user	body	teamUpdateUserRequestBody	true	"updated team user object"
//	@Success		200		object	standardJsonResponse{}
//	@Success		403		object	standardJsonResponse{}
//	@Success		500		object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/teams/{teamId}/users/{userId} [put]
func (s *Service) handleTeamUpdateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)
		vars := mux.Vars(r)
		teamID := vars["teamId"]
		tidErr := validate.Var(teamID, "required,uuid")
		if tidErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, tidErr.Error()))
			return
		}
		userID := vars["userId"]
		idErr := validate.Var(userID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

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
			return
		}

		_, err := s.TeamDataSvc.TeamUpdateUser(ctx, teamID, userID, u.Role)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleTeamInviteUser error", zap.Error(err), zap.String("team_id", teamID),
				zap.String("user_id", userID), zap.String("team_role", u.Role),
				zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, nil, nil)
	}
}

// handleTeamRemoveUser handles removing user from a team
//
//	@Summary		Remove Team User
//	@Description	Remove a user from the team
//	@Tags			team
//	@Produce		json
//	@Param			teamId	path	string	true	"the team ID"
//	@Param			userId	path	string	true	"the user ID"
//	@Success		200		object	standardJsonResponse{}
//	@Success		403		object	standardJsonResponse{}
//	@Success		500		object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/teams/{teamId}/users/{userId} [delete]
func (s *Service) handleTeamRemoveUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)
		vars := mux.Vars(r)
		teamID := vars["teamId"]
		tidErr := validate.Var(teamID, "required,uuid")
		if tidErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, tidErr.Error()))
			return
		}
		userID := vars["userId"]
		idErr := validate.Var(userID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		err := s.TeamDataSvc.TeamRemoveUser(ctx, teamID, userID)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleTeamRemoveUser error", zap.Error(err), zap.String("team_id", teamID),
				zap.String("user_id", userID), zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, nil, nil)
	}
}

// handleGetTeamPokerGames gets a list of battles associated to the team
//
//	@Summary		Get Team Battles
//	@Description	Get a list of battles associated to the team
//	@Tags			team
//	@Produce		json
//	@Param			teamId	path	string	true	"the team ID"
//	@Success		200		object	standardJsonResponse{data=[]thunderdome.Poker}
//	@Security		ApiKeyAuth
//	@Router			/teams/{teamId}/battles [get]
func (s *Service) handleGetTeamPokerGames() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		vars := mux.Vars(r)
		teamID := vars["teamId"]
		idErr := validate.Var(teamID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		limit, offset := getLimitOffsetFromRequest(r)

		games := s.TeamDataSvc.TeamPokerList(ctx, teamID, limit, offset)

		s.Success(w, r, http.StatusOK, games, nil)
	}
}

// handleTeamRemovePokerGame handles removing poker game from a team
//
//	@Summary		Remove Team Poker
//	@Description	Remove a poker game from the team
//	@Tags			team
//	@Produce		json
//	@Param			teamId		path	string	true	"the team ID"
//	@Param			battleId	path	string	true	"the game ID"
//	@Success		200			object	standardJsonResponse{}
//	@Success		403			object	standardJsonResponse{}
//	@Success		500			object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/teams/{teamId}/battles/{battleId} [delete]
func (s *Service) handleTeamRemovePokerGame() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)
		vars := mux.Vars(r)
		teamID := vars["teamId"]
		idErr := validate.Var(teamID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		gameID := vars["battleId"]
		idErr = validate.Var(gameID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		err := s.TeamDataSvc.TeamRemovePoker(ctx, teamID, gameID)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleTeamRemoveBattle error", zap.Error(err), zap.String("team_id", teamID),
				zap.String("battle_id", gameID), zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, nil, nil)
	}
}

// handleDeleteTeam handles deleting a team
//
//	@Summary		Delete Team
//	@Description	Delete a Team
//	@Tags			team
//	@Produce		json
//	@Param			teamId	path	string	true	"the team ID"
//	@Success		200		object	standardJsonResponse{}
//	@Success		403		object	standardJsonResponse{}
//	@Success		500		object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/teams/{teamId} [delete]
func (s *Service) handleDeleteTeam() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)
		vars := mux.Vars(r)
		teamID := vars["teamId"]
		idErr := validate.Var(teamID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		err := s.TeamDataSvc.TeamDelete(ctx, teamID)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleDeleteTeam error", zap.Error(err), zap.String("team_id", teamID),
				zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, nil, nil)
	}
}

// handleGetTeamRetros gets a list of retros associated to the team
//
//	@Summary		Get Team Retros
//	@Description	Get a list of retros associated to the team
//	@Tags			team
//	@Produce		json
//	@Param			teamId	path	string	true	"the team ID"
//	@Success		200		object	standardJsonResponse{data=[]thunderdome.Retro}
//	@Security		ApiKeyAuth
//	@Router			/teams/{teamId}/retros [get]
func (s *Service) handleGetTeamRetros() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		vars := mux.Vars(r)
		teamID := vars["teamId"]
		idErr := validate.Var(teamID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		limit, offset := getLimitOffsetFromRequest(r)

		retrospectives := s.TeamDataSvc.TeamRetroList(ctx, teamID, limit, offset)

		s.Success(w, r, http.StatusOK, retrospectives, nil)
	}
}

// handleTeamRemoveRetro handles removing retro from a team
//
//	@Summary		Remove Team Retro
//	@Description	Remove a retro from the team
//	@Tags			team
//	@Produce		json
//	@Param			teamId	path	string	true	"the team ID"
//	@Param			retroId	path	string	true	"the retro ID"
//	@Success		200		object	standardJsonResponse{}
//	@Success		403		object	standardJsonResponse{}
//	@Success		500		object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/teams/{teamId}/retros/{retroId} [delete]
func (s *Service) handleTeamRemoveRetro() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)
		vars := mux.Vars(r)
		teamID := vars["teamId"]
		idErr := validate.Var(teamID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		retrospectiveID := vars["retroId"]
		idErr = validate.Var(retrospectiveID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		err := s.TeamDataSvc.TeamRemoveRetro(ctx, teamID, retrospectiveID)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleTeamRemoveRetro error", zap.Error(err), zap.String("team_id", teamID),
				zap.String("retro_id", retrospectiveID), zap.String("session_user_id", sessionUserID))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

// handleGetTeamStoryboards gets a list of storyboards associated to the team
//
//	@Summary		Get Team Storyboards
//	@Description	Get a list of storyboards associated to the team
//	@Tags			team
//	@Produce		json
//	@Param			teamId	path	string	true	"the team ID"
//	@Success		200		object	standardJsonResponse{data=[]thunderdome.Storyboard}
//	@Security		ApiKeyAuth
//	@Router			/teams/{teamId}/storyboards [get]
func (s *Service) handleGetTeamStoryboards() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		vars := mux.Vars(r)
		teamID := vars["teamId"]
		idErr := validate.Var(teamID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		limit, offset := getLimitOffsetFromRequest(r)

		storyboards := s.TeamDataSvc.TeamStoryboardList(ctx, teamID, limit, offset)

		s.Success(w, r, http.StatusOK, storyboards, nil)
	}
}

// handleTeamRemoveStoryboard handles removing storyboard from a team
//
//	@Summary		Remove Team Storyboard
//	@Description	Remove a storyboard from the team
//	@Tags			team
//	@Produce		json
//	@Param			teamId			path	string	true	"the team ID"
//	@Param			storyboardId	path	string	true	"the storyboard ID"
//	@Success		200				object	standardJsonResponse{}
//	@Success		403				object	standardJsonResponse{}
//	@Success		500				object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/teams/{teamId}/storyboards/{storyboardId} [delete]
func (s *Service) handleTeamRemoveStoryboard() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)
		vars := mux.Vars(r)
		teamID := vars["teamId"]
		idErr := validate.Var(teamID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		storyboardID := vars["storyboardId"]
		idErr = validate.Var(storyboardID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		err := s.TeamDataSvc.TeamRemoveStoryboard(ctx, teamID, storyboardID)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleTeamRemoveStoryboard error", zap.Error(err), zap.String("team_id", teamID),
				zap.String("storyboard_id", storyboardID), zap.String("session_user_id", sessionUserID))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

// handleGetTeamRetroActions gets a list of retro actions
//
//	@Summary		Get Retro Actions
//	@Description	get list of retro actions
//	@Tags			team
//	@Produce		json
//	@Param			limit		query	int		false	"Max number of results to return"
//	@Param			offset		query	int		false	"Starting point to return rows from, should be multiplied by limit or 0"
//	@Param			completed	query	boolean	false	"Only completed retro actions"
//	@Success		200			object	standardJsonResponse{data=[]thunderdome.RetroAction}
//	@Failure		500			object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/teams/{teamId}/retro-actions [get]
func (s *Service) handleGetTeamRetroActions() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)
		vars := mux.Vars(r)
		teamID := vars["teamId"]
		idErr := validate.Var(teamID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		limit, offset := getLimitOffsetFromRequest(r)
		var err error
		var count int
		var actions []*thunderdome.RetroAction
		query := r.URL.Query()
		completed, _ := strconv.ParseBool(query.Get("completed"))

		actions, count, err = s.RetroDataSvc.GetTeamRetroActions(teamID, limit, offset, completed)

		if err != nil {
			s.Logger.Ctx(ctx).Error("handleGetTeamRetroActions error", zap.Error(err), zap.String("team_id", teamID),
				zap.Int("limit", limit), zap.Int("offset", offset), zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		meta := &pagination{
			Count:  count,
			Offset: offset,
			Limit:  limit,
		}

		s.Success(w, r, http.StatusOK, actions, meta)
	}
}

// handleGetTeamUserInvites gets a list of user invites associated to the team
//
//	@Summary		Get Team User Invites
//	@Description	Get a list of user invites associated to the team
//	@Tags			team
//	@Produce		json
//	@Param			teamId	path	string	true	"the team ID"
//	@Success		200		object	standardJsonResponse{data=[]thunderdome.TeamUserInvite}
//	@Security		ApiKeyAuth
//	@Router			/teams/{teamId}/invites [get]
func (s *Service) handleGetTeamUserInvites() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)
		vars := mux.Vars(r)
		teamID := vars["teamId"]
		idErr := validate.Var(teamID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		invites, err := s.TeamDataSvc.TeamGetUserInvites(ctx, teamID)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleGetTeamUserInvites error", zap.Error(err), zap.String("team_id", teamID),
				zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, invites, nil)
	}
}

// handleDeleteTeamUserInvite handles deleting user invite from a team
//
//	@Summary		Deletes Team User Invite
//	@Description	Delete a user invite from the team
//	@Tags			team
//	@Produce		json
//	@Param			teamId		path	string	true	"the team ID"
//	@Param			inviteId	path	string	true	"the user invite ID"
//	@Success		200			object	standardJsonResponse{}
//	@Success		403			object	standardJsonResponse{}
//	@Success		500			object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/teams/{teamId}/invites/{inviteId} [delete]
func (s *Service) handleDeleteTeamUserInvite() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)
		vars := mux.Vars(r)
		teamID := vars["teamId"]
		idErr := validate.Var(teamID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		inviteID := vars["inviteId"]
		idErr = validate.Var(inviteID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		err := s.TeamDataSvc.TeamDeleteUserInvite(ctx, inviteID)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleDeleteTeamUserInvite error", zap.Error(err), zap.String("team_id", teamID),
				zap.String("invite_id", inviteID), zap.String("session_user_id", sessionUserID))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

// handleTeamMetrics gets the metrics for a specific team
//
//	@Summary		Get Team Metrics
//	@Description	Get metrics for a specific team such as user count, poker game count, etc.
//	@Tags			admin
//	@Produce		json
//	@Param			teamID	path	string	true	"Team ID"
//	@Success		200		object	standardJsonResponse{data=thunderdome.TeamMetrics}
//	@Failure		400		object	standardJsonResponse{}
//	@Failure		404		object	standardJsonResponse{}
//	@Failure		500		object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/admin/teams/{teamID}/metrics [get]
func (s *Service) handleTeamMetrics() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)
		vars := mux.Vars(r)
		teamID := vars["teamID"]
		idErr := validate.Var(teamID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		if teamID == "" {
			s.Failure(w, r, http.StatusBadRequest, errors.New("team ID is required"))
			return
		}

		metrics, err := s.TeamDataSvc.GetTeamMetrics(ctx, teamID)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleTeamMetrics error", zap.Error(err),
				zap.String("session_user_id", sessionUserID),
				zap.String("team_id", teamID))
			if err.Error() == "no team found with ID "+teamID {
				s.Failure(w, r, http.StatusNotFound, err)
			} else {
				s.Failure(w, r, http.StatusInternalServerError, err)
			}
			return
		}

		s.Success(w, r, http.StatusOK, metrics, nil)
	}
}
