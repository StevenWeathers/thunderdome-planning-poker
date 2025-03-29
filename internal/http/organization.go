package http

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"go.uber.org/zap"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
)

type addUserRequestBody struct {
	UserID string `json:"user_id" validate:"required,uuid"`
	Role   string `json:"role" enums:"MEMBER,ADMIN" validate:"required,oneof=MEMBER ADMIN"`
}

type organizationResponse struct {
	Organization *thunderdome.Organization `json:"organization"`
	Role         string                    `json:"role"`
}

type orgTeamResponse struct {
	Organization     *thunderdome.Organization `json:"organization"`
	Team             *thunderdome.Team         `json:"team"`
	OrganizationRole string                    `json:"organizationRole"`
	TeamRole         string                    `json:"teamRole"`
}

// handleGetOrganizationsByUser gets a list of organizations the user is a part of
//
//	@Summary		Get Users Organizations
//	@Description	Get list of organizations for the authenticated user
//	@Tags			organization
//	@Produce		json
//	@Param			userId	path	string	true	"the user ID to get organizations for"
//	@Param			limit	query	int		false	"Max number of results to return"
//	@Param			offset	query	int		false	"Starting point to return rows from, should be multiplied by limit or 0"
//	@Success		200		object	standardJsonResponse{data=[]thunderdome.UserOrganization}
//	@Failure		403		object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/users/{userId}/organizations [get]
func (s *Service) handleGetOrganizationsByUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if !s.Config.OrganizationsEnabled {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, "ORGANIZATIONS_DISABLED"))
			return
		}

		userID := r.PathValue("userId")
		idErr := validate.Var(userID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		limit, offset := getLimitOffsetFromRequest(r)

		organizations := s.OrganizationDataSvc.OrganizationListByUser(ctx, userID, limit, offset)

		s.Success(w, r, http.StatusOK, organizations, nil)
	}
}

// handleGetOrganizationByUser gets an organization with user role
//
//	@Summary		Get Organization
//	@Description	Get an organization with user role
//	@Tags			organization
//	@Produce		json
//	@Param			orgId	path	string	true	"organization id"
//	@Success		200		object	standardJsonResponse{data=organizationResponse}
//	@Failure		403		object	standardJsonResponse{}
//	@Failure		500		object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/organizations/{orgId} [get]
func (s *Service) handleGetOrganizationByUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !s.Config.OrganizationsEnabled {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, "ORGANIZATIONS_DISABLED"))
			return
		}
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)
		orgRole := ctx.Value(contextKeyOrgRole).(string)

		orgID := r.PathValue("orgId")
		idErr := validate.Var(orgID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		organization, err := s.OrganizationDataSvc.OrganizationGetByID(ctx, orgID)
		if err != nil {
			s.Logger.Ctx(ctx).Error(
				"handleGetOrganizationByUser error", zap.Error(err), zap.String("organization_id", orgID),
				zap.String("session_user_id", sessionUserID), zap.String("organization_role", orgRole))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		result := &organizationResponse{
			Organization: organization,
			Role:         orgRole,
		}

		s.Success(w, r, http.StatusOK, result, nil)
	}
}

// handleCreateOrganization handles creating an organization with current user as admin
//
//	@Summary		Create Organization
//	@Description	Create organization with current user as admin
//	@Tags			organization
//	@Produce		json
//	@Param			userId			path	string					true	"user id"
//	@Param			organization	body	teamCreateRequestBody	true	"new organization object"
//	@Success		200				object	standardJsonResponse{data=thunderdome.Organization}
//	@Failure		403				object	standardJsonResponse{}
//	@Failure		500				object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/users/{userId}/organizations [post]
func (s *Service) handleCreateOrganization() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !s.Config.OrganizationsEnabled {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, "ORGANIZATIONS_DISABLED"))
			return
		}
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)

		userID := r.PathValue("userId")
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

		organization, err := s.OrganizationDataSvc.OrganizationCreate(ctx, userID, team.Name)
		if err != nil {
			s.Logger.Ctx(ctx).Error(
				"handleCreateOrganization error", zap.Error(err), zap.String("entity_user_id", userID),
				zap.String("session_user_id", sessionUserID), zap.String("organization_name", team.Name))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, organization, nil)
	}
}

// handleOrganizationUpdate handles updating an organization
//
//	@Summary		Update Organization
//	@Description	Update organization
//	@Tags			organization
//	@Produce		json
//	@Param			orgId			path	string					true	"organization id"
//	@Param			organization	body	teamCreateRequestBody	true	"updated organization object"
//	@Success		200				object	standardJsonResponse{data=thunderdome.Organization}
//	@Failure		403				object	standardJsonResponse{}
//	@Failure		500				object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/organizations/{orgId} [put]
func (s *Service) handleOrganizationUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !s.Config.OrganizationsEnabled {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, "ORGANIZATIONS_DISABLED"))
			return
		}
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)

		orgID := r.PathValue("orgId")
		idErr := validate.Var(orgID, "required,uuid")
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

		organization, err := s.OrganizationDataSvc.OrganizationUpdate(ctx, orgID, team.Name)
		if err != nil {
			s.Logger.Ctx(ctx).Error(
				"handleCreateOrganization error", zap.Error(err),
				zap.String("organization_id", orgID),
				zap.String("session_user_id", sessionUserID),
				zap.String("organization_name", team.Name))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, organization, nil)
	}
}

// handleGetOrganizationTeams gets a list of teams associated to the organization
//
//	@Summary		Get Organization Teams
//	@Description	Get a list of organization teams
//	@Tags			organization
//	@Produce		json
//	@Param			orgId	path	string	true	"organization id"
//	@Success		200		object	standardJsonResponse{data=[]thunderdome.Team}
//	@Failure		403		object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/organizations/{orgId}/teams [get]
func (s *Service) handleGetOrganizationTeams() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !s.Config.OrganizationsEnabled {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, "ORGANIZATIONS_DISABLED"))
			return
		}
		ctx := r.Context()

		orgID := r.PathValue("orgId")
		idErr := validate.Var(orgID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		limit, offset := getLimitOffsetFromRequest(r)

		teams := s.OrganizationDataSvc.OrganizationTeamList(ctx, orgID, limit, offset)

		s.Success(w, r, http.StatusOK, teams, nil)
	}
}

// handleGetOrganizationUsers gets a list of users associated to the organization
//
//	@Summary		Get Organization Users
//	@Description	get a list of organization users
//	@Tags			organization
//	@Produce		json
//	@Param			orgId	path	string	true	"organization id"
//	@Success		200		object	standardJsonResponse{data=[]thunderdome.User}
//	@Failure		403		object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/organizations/{orgId}/users [get]
func (s *Service) handleGetOrganizationUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !s.Config.OrganizationsEnabled {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, "ORGANIZATIONS_DISABLED"))
			return
		}
		ctx := r.Context()

		orgID := r.PathValue("orgId")
		idErr := validate.Var(orgID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		limit, offset := getLimitOffsetFromRequest(r)

		teams := s.OrganizationDataSvc.OrganizationUserList(ctx, orgID, limit, offset)

		s.Success(w, r, http.StatusOK, teams, nil)
	}
}

// handleCreateOrganizationTeam handles creating an organization team
//
//	@Summary		Create Organization Team
//	@Description	Create organization team with current user as admin
//	@Tags			organization
//	@Produce		json
//	@Param			orgId	path	string					true	"organization id"
//	@Param			team	body	teamCreateRequestBody	true	"new team object"
//	@Success		200		object	standardJsonResponse{data=thunderdome.Team}
//	@Failure		403		object	standardJsonResponse{}
//	@Failure		500		object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/organizations/{orgId}/teams [post]
func (s *Service) handleCreateOrganizationTeam() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !s.Config.OrganizationsEnabled {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, "ORGANIZATIONS_DISABLED"))
			return
		}
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)

		orgID := r.PathValue("orgId")
		idErr := validate.Var(orgID, "required,uuid")
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

		newTeam, err := s.OrganizationDataSvc.OrganizationTeamCreate(ctx, orgID, team.Name)
		s.Logger.Ctx(ctx).Error(
			"handleCreateOrganizationTeam error", zap.Error(err), zap.String("organization_id", orgID),
			zap.String("session_user_id", sessionUserID), zap.String("team_name", team.Name))
		if err != nil {
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, newTeam, nil)
	}
}

// handleOrganizationInviteUser handles inviting user to an organization
//
//	@Summary		Invite Org User
//	@Description	Invite user to organization
//	@Tags			organization
//	@Produce		json
//	@Param			orgId	path	string						true	"organization id"
//	@Param			user	body	teamInviteUserRequestBody	true	"new organization user object"
//	@Success		200		object	standardJsonResponse{}
//	@Failure		403		object	standardJsonResponse{}
//	@Failure		500		object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/organizations/{orgId}/invites [post]
func (s *Service) handleOrganizationInviteUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !s.Config.OrganizationsEnabled {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, "ORGANIZATIONS_DISABLED"))
			return
		}
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)

		orgID := r.PathValue("orgId")
		idErr := validate.Var(orgID, "required,uuid")
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

		userEmail := strings.ToLower(u.Email)

		if s.Config.LdapEnabled || s.Config.HeaderAuthEnabled || s.Config.OIDCAuth.Enabled {
			user, userErr := s.UserDataSvc.GetUserByEmail(ctx, userEmail)
			if userErr == nil {
				_, err := s.OrganizationDataSvc.OrganizationAddUser(ctx, orgID, user.ID, u.Role)
				if err != nil {
					s.Logger.Ctx(ctx).Error(
						"handleOrganizationInviteUser error", zap.Error(err), zap.String("user_id", user.ID),
						zap.String("session_user_id", sessionUserID), zap.String("organization_id", orgID),
						zap.String("user_role", u.Role))
					s.Failure(w, r, http.StatusInternalServerError, err)
					return
				}
				s.Success(w, r, http.StatusOK, nil, userAddMeta{Invited: false, Added: true})
				return
			} else if userErr != nil && !errors.Is(userErr, sql.ErrNoRows) {
				s.Logger.Ctx(ctx).Error(
					"handleOrganizationInviteUser error", zap.Error(userErr), zap.String("user_email", userEmail),
					zap.String("session_user_id", sessionUserID), zap.String("organization_id", orgID))
				s.Failure(w, r, http.StatusInternalServerError, Errorf(ENOTFOUND, "USER_NOT_FOUND"))
				return
			}
		}

		inviteID, inviteErr := s.OrganizationDataSvc.OrganizationInviteUser(ctx, orgID, userEmail, u.Role)
		if inviteErr != nil {
			s.Logger.Ctx(ctx).Error("handleOrganizationInviteUser error", zap.Error(inviteErr),
				zap.String("organization_id", orgID), zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, inviteErr)
			return
		}
		org, orgErr := s.OrganizationDataSvc.OrganizationGetByID(ctx, orgID)
		if orgErr != nil {
			s.Logger.Ctx(ctx).Error("handleOrganizationInviteUser error", zap.Error(orgErr),
				zap.String("organization_id", orgID), zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, orgErr)
			return
		}

		emailErr := s.Email.SendOrganizationInvite(org.Name, userEmail, inviteID)
		if emailErr != nil {
			s.Logger.Ctx(ctx).Error("handleOrganizationInviteUser error", zap.Error(emailErr),
				zap.String("organization_id", orgID), zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, emailErr)
			return
		}

		s.Success(w, r, http.StatusOK, nil, userAddMeta{Invited: true, Added: false})
	}
}

// handleOrganizationUpdateUser handles updating an organization user
//
//	@Summary		Update Org User
//	@Description	Update organization user
//	@Tags			organization
//	@Produce		json
//	@Param			orgId	path	string						true	"organization id"
//	@Param			userId	path	string						true	"user id"
//	@Param			user	body	teamUpdateUserRequestBody	true	"organization user object"
//	@Success		200		object	standardJsonResponse{}
//	@Failure		403		object	standardJsonResponse{}
//	@Failure		500		object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/organizations/{orgId}/users/{userId} [put]
func (s *Service) handleOrganizationUpdateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !s.Config.OrganizationsEnabled {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, "ORGANIZATIONS_DISABLED"))
			return
		}
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)

		orgID := r.PathValue("orgId")
		idErr := validate.Var(orgID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		userID := r.PathValue("userId")
		idErr = validate.Var(userID, "required,uuid")
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

		_, err := s.OrganizationDataSvc.OrganizationUpdateUser(ctx, orgID, userID, u.Role)
		if err != nil {
			s.Logger.Ctx(ctx).Error(
				"handleOrganizationInviteUser error", zap.Error(err), zap.String("user_id", userID),
				zap.String("session_user_id", sessionUserID), zap.String("organization_id", orgID),
				zap.String("user_role", u.Role))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, nil, nil)
	}
}

// handleOrganizationRemoveUser handles removing user from an organization (including departments, teams)
//
//	@Summary		Remove Org User
//	@Description	Remove user from organization including departments and teams
//	@Tags			organization
//	@Produce		json
//	@Param			orgId	path	string	true	"organization id"
//	@Param			userId	path	string	true	"user id"
//	@Success		200		object	standardJsonResponse{}
//	@Failure		403		object	standardJsonResponse{}
//	@Failure		500		object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/organizations/{orgId}/users/{userId} [delete]
func (s *Service) handleOrganizationRemoveUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !s.Config.OrganizationsEnabled {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, "ORGANIZATIONS_DISABLED"))
			return
		}
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)

		orgID := r.PathValue("orgId")
		idErr := validate.Var(orgID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		userID := r.PathValue("userId")
		idErr = validate.Var(userID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		err := s.OrganizationDataSvc.OrganizationRemoveUser(ctx, orgID, userID)
		if err != nil {
			s.Logger.Ctx(ctx).Error(
				"handleOrganizationRemoveUser error", zap.Error(err), zap.String("user_id", userID),
				zap.String("session_user_id", sessionUserID), zap.String("organization_id", orgID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, nil, nil)
	}
}

// handleGetOrganizationTeamByUser gets a team with users roles
//
//	@Summary		Get Organization Team
//	@Description	Get an organizations team with users roles
//	@Tags			organization
//	@Produce		json
//	@Param			orgId	path	string	true	"organization id"
//	@Param			teamId	path	string	true	"team id"
//	@Success		200		object	standardJsonResponse{data=orgTeamResponse}
//	@Failure		403		object	standardJsonResponse{}
//	@Failure		500		object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/organizations/{orgId}/teams/{teamId} [get]
func (s *Service) handleGetOrganizationTeamByUser() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		if !s.Config.OrganizationsEnabled {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, "ORGANIZATIONS_DISABLED"))
			return
		}
		ctx := r.Context()
		teamUserRoles := ctx.Value(contextKeyUserTeamRoles).(*thunderdome.UserTeamRoleInfo)
		var emptyRole = ""
		orgRole := teamUserRoles.OrganizationRole
		if orgRole == nil {
			orgRole = &emptyRole
		}
		teamRole := teamUserRoles.TeamRole
		if teamRole == nil {
			teamRole = &emptyRole
		}
		sessionUserID := ctx.Value(contextKeyUserID).(string)

		orgID := r.PathValue("orgId")
		idErr := validate.Var(orgID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		teamID := r.PathValue("teamId")
		idErr = validate.Var(teamID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		organization, err := s.OrganizationDataSvc.OrganizationGetByID(ctx, orgID)
		if err != nil {
			s.Logger.Ctx(ctx).Error(
				"handleGetOrganizationTeamByUser error", zap.Error(err), zap.String("organization_id", orgID),
				zap.String("session_user_id", sessionUserID), zap.String("organization_role", *orgRole),
				zap.String("team_role", *teamRole), zap.String("team_id", teamID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		team, err := s.TeamDataSvc.TeamGetByID(ctx, teamID)
		if err != nil {
			s.Logger.Ctx(ctx).Error(
				"handleGetOrganizationTeamByUser error", zap.Error(err), zap.String("organization_id", orgID),
				zap.String("session_user_id", sessionUserID), zap.String("organization_role", *orgRole),
				zap.String("team_role", *teamRole), zap.String("team_id", teamID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		result := &orgTeamResponse{
			Organization:     organization,
			Team:             team,
			OrganizationRole: *orgRole,
			TeamRole:         *teamRole,
		}

		s.Success(w, r, http.StatusOK, result, nil)
	}
}

// handleOrganizationTeamAddUser handles adding user to a team so long as they are in the organization
//
//	@Summary		Add Org Team User
//	@Description	Add user to organization team as long as they are already in the organization
//	@Tags			organization
//	@Produce		json
//	@Param			orgId	path	string				true	"organization id"
//	@Param			teamId	path	string				true	"team id"
//	@Param			user	body	addUserRequestBody	true	"new team user object"
//	@Success		200		object	standardJsonResponse{}
//	@Failure		403		object	standardJsonResponse{}
//	@Failure		500		object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/organizations/{orgId}/teams/{teamId}/users [post]
func (s *Service) handleOrganizationTeamAddUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !s.Config.OrganizationsEnabled {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, "ORGANIZATIONS_DISABLED"))
			return
		}

		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)

		orgID := r.PathValue("orgId")
		idErr := validate.Var(orgID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		teamID := r.PathValue("teamId")
		idErr = validate.Var(teamID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		var u = addUserRequestBody{}
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

		orgRole, roleErr := s.OrganizationDataSvc.OrganizationUserRole(ctx, u.UserID, orgID)
		if orgRole == "" || roleErr != nil {
			s.Logger.Ctx(ctx).Error(
				"handleOrganizationTeamAddUser error", zap.Error(roleErr), zap.String("organization_id", orgID),
				zap.String("session_user_id", sessionUserID), zap.String("organization_role", orgRole),
				zap.String("team_role", u.Role), zap.String("team_id", teamID), zap.String("user_id", u.UserID))
			s.Failure(w, r, http.StatusInternalServerError, Errorf(EUNAUTHORIZED, "ORGANIZATION_USER_REQUIRED"))
			return
		}

		_, err := s.TeamDataSvc.TeamAddUser(ctx, teamID, u.UserID, u.Role)
		if err != nil {
			s.Logger.Ctx(ctx).Error(
				"handleOrganizationTeamAddUser error", zap.Error(err), zap.String("organization_id", orgID),
				zap.String("session_user_id", sessionUserID), zap.String("organization_role", orgRole),
				zap.String("team_role", u.Role), zap.String("team_id", teamID), zap.String("user_id", u.UserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, nil, nil)
	}
}

// handleDeleteOrganization handles deleting an organization
//
//	@Summary		Delete Organization
//	@Description	Delete an Organization
//	@Tags			organization
//	@Produce		json
//	@Param			orgId	path	string	true	"the organization ID"
//	@Success		200		object	standardJsonResponse{}
//	@Success		403		object	standardJsonResponse{}
//	@Success		500		object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/organizations/{orgId} [delete]
func (s *Service) handleDeleteOrganization() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)

		orgID := r.PathValue("orgId")
		idErr := validate.Var(orgID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		err := s.OrganizationDataSvc.OrganizationDelete(ctx, orgID)
		if err != nil {
			s.Logger.Ctx(ctx).Error(
				"handleDeleteOrganization error", zap.Error(err), zap.String("organization_id", orgID),
				zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, nil, nil)
	}
}

// handleGetOrganizationUserInvites gets a list of user invites associated to the organization
//
//	@Summary		Get Organization User Invites
//	@Description	Get a list of user invites associated to the organization
//	@Tags			organization
//	@Produce		json
//	@Param			orgId	path	string	true	"the org ID"
//	@Success		200		object	standardJsonResponse{data=[]thunderdome.OrganizationUserInvite}
//	@Security		ApiKeyAuth
//	@Router			/organizations/{orgId}/invites [get]
func (s *Service) handleGetOrganizationUserInvites() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)

		orgID := r.PathValue("orgId")
		idErr := validate.Var(orgID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		invites, err := s.OrganizationDataSvc.OrganizationGetUserInvites(ctx, orgID)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleGetOrganizationUserInvites error", zap.Error(err), zap.String("organization_id", orgID),
				zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, invites, nil)
	}
}

// handleDeleteOrganizationUserInvite handles deleting user invite from an organization
//
//	@Summary		Delete Org User Invite
//	@Description	Delete user invite from organization
//	@Tags			organization
//	@Produce		json
//	@Param			orgId		path	string	true	"organization id"
//	@Param			inviteId	path	string	true	"invite id"
//	@Success		200			object	standardJsonResponse{}
//	@Failure		403			object	standardJsonResponse{}
//	@Failure		500			object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/organizations/{orgId}/invites/{inviteId} [delete]
func (s *Service) handleDeleteOrganizationUserInvite() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !s.Config.OrganizationsEnabled {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, "ORGANIZATIONS_DISABLED"))
			return
		}
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)

		orgID := r.PathValue("orgId")
		idErr := validate.Var(orgID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		inviteID := r.PathValue("inviteId")
		idErr = validate.Var(inviteID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		err := s.OrganizationDataSvc.OrganizationDeleteUserInvite(ctx, inviteID)
		if err != nil {
			s.Logger.Ctx(ctx).Error(
				"handleDeleteOrganizationUserInvite error", zap.Error(err), zap.String("invite_id", inviteID),
				zap.String("session_user_id", sessionUserID), zap.String("organization_id", orgID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, nil, nil)
	}
}

// handleOrganizationMetrics gets the metrics for a specific organization
//
//	@Summary		Get Organization Metrics
//	@Description	Get metrics for a specific organization such as user count, team count, etc.
//	@Tags			admin
//	@Produce		json
//	@Param			organizationID	path	string	true	"Organization ID"
//	@Success		200				object	standardJsonResponse{data=thunderdome.OrganizationMetrics}
//	@Failure		400				object	standardJsonResponse{}
//	@Failure		404				object	standardJsonResponse{}
//	@Failure		500				object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/organizations/{orgId}/metrics [get]
func (s *Service) handleOrganizationMetrics() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)

		organizationID := r.PathValue("orgId")
		idErr := validate.Var(organizationID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, errors.New("organization ID is required"))
			return
		}

		metrics, err := s.OrganizationDataSvc.GetOrganizationMetrics(ctx, organizationID)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleOrganizationMetrics error", zap.Error(err),
				zap.String("session_user_id", sessionUserID),
				zap.String("organization_id", organizationID))
			if err.Error() == "no organization found with ID "+organizationID {
				s.Failure(w, r, http.StatusNotFound, err)
			} else {
				s.Failure(w, r, http.StatusInternalServerError, err)
			}
			return
		}

		s.Success(w, r, http.StatusOK, metrics, nil)
	}
}
