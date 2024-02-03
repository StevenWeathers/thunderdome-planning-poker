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

	"github.com/gorilla/mux"
)

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
// @Summary      Get Users Organizations
// @Description  Get list of organizations for the authenticated user
// @Tags         organization
// @Produce      json
// @Param        userId  path    string  true   "the user ID to get organizations for"
// @Param        limit   query   int     false  "Max number of results to return"
// @Param        offset  query   int     false  "Starting point to return rows from, should be multiplied by limit or 0"
// @Success      200     object  standardJsonResponse{data=[]thunderdome.Organization}
// @Failure      403     object  standardJsonResponse{}
// @Security     ApiKeyAuth
// @Router       /users/{userId}/organizations [get]
func (s *Service) handleGetOrganizationsByUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if !s.Config.OrganizationsEnabled {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, "ORGANIZATIONS_DISABLED"))
			return
		}
		vars := mux.Vars(r)
		UserID := vars["userId"]

		Limit, Offset := getLimitOffsetFromRequest(r)

		Organizations := s.OrganizationDataSvc.OrganizationListByUser(ctx, UserID, Limit, Offset)

		s.Success(w, r, http.StatusOK, Organizations, nil)
	}
}

// handleGetOrganizationByUser gets an organization with user role
// @Summary      Get Organization
// @Description  Get an organization with user role
// @Tags         organization
// @Produce      json
// @Param        orgId  path    string  true  "organization id"
// @Success      200    object  standardJsonResponse{data=organizationResponse}
// @Failure      403    object  standardJsonResponse{}
// @Failure      500    object  standardJsonResponse{}
// @Security     ApiKeyAuth
// @Router       /organizations/{orgId} [get]
func (s *Service) handleGetOrganizationByUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !s.Config.OrganizationsEnabled {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, "ORGANIZATIONS_DISABLED"))
			return
		}
		ctx := r.Context()
		SessionUserID := ctx.Value(contextKeyUserID).(string)
		OrgRole := ctx.Value(contextKeyOrgRole).(string)
		vars := mux.Vars(r)
		OrgID := vars["orgId"]

		Organization, err := s.OrganizationDataSvc.OrganizationGet(ctx, OrgID)
		if err != nil {
			s.Logger.Ctx(ctx).Error(
				"handleGetOrganizationByUser error", zap.Error(err), zap.String("organization_id", OrgID),
				zap.String("session_user_id", SessionUserID), zap.String("organization_role", OrgRole))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		result := &organizationResponse{
			Organization: Organization,
			Role:         OrgRole,
		}

		s.Success(w, r, http.StatusOK, result, nil)
	}
}

// handleCreateOrganization handles creating an organization with current user as admin
// @Summary      Create Organization
// @Description  Create organization with current user as admin
// @Tags         organization
// @Produce      json
// @Param        userId        path    string                 true  "user id"
// @Param        organization  body    teamCreateRequestBody  true  "new organization object"
// @Success      200           object  standardJsonResponse{data=thunderdome.Organization}
// @Failure      403           object  standardJsonResponse{}
// @Failure      500           object  standardJsonResponse{}
// @Security     ApiKeyAuth
// @Router       /users/{userId}/organizations [post]
func (s *Service) handleCreateOrganization() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !s.Config.OrganizationsEnabled {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, "ORGANIZATIONS_DISABLED"))
			return
		}
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

		Organization, err := s.OrganizationDataSvc.OrganizationCreate(ctx, UserID, team.Name)
		if err != nil {
			s.Logger.Ctx(ctx).Error(
				"handleCreateOrganization error", zap.Error(err), zap.String("entity_user_id", UserID),
				zap.String("session_user_id", SessionUserID), zap.String("organization_name", team.Name))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, Organization, nil)
	}
}

// handleGetOrganizationTeams gets a list of teams associated to the organization
// @Summary      Get Organization Teams
// @Description  Get a list of organization teams
// @Tags         organization
// @Produce      json
// @Param        orgId  path    string  true  "organization id"
// @Success      200    object  standardJsonResponse{data=[]thunderdome.Team}
// @Failure      403    object  standardJsonResponse{}
// @Security     ApiKeyAuth
// @Router       /organizations/{orgId}/teams [get]
func (s *Service) handleGetOrganizationTeams() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !s.Config.OrganizationsEnabled {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, "ORGANIZATIONS_DISABLED"))
			return
		}
		ctx := r.Context()
		vars := mux.Vars(r)
		OrgID := vars["orgId"]
		Limit, Offset := getLimitOffsetFromRequest(r)

		Teams := s.OrganizationDataSvc.OrganizationTeamList(ctx, OrgID, Limit, Offset)

		s.Success(w, r, http.StatusOK, Teams, nil)
	}
}

// handleGetOrganizationUsers gets a list of users associated to the organization
// @Summary      Get Organization Users
// @Description  get a list of organization users
// @Tags         organization
// @Produce      json
// @Param        orgId  path    string  true  "organization id"
// @Success      200    object  standardJsonResponse{data=[]thunderdome.User}
// @Failure      403    object  standardJsonResponse{}
// @Security     ApiKeyAuth
// @Router       /organizations/{orgId}/users [get]
func (s *Service) handleGetOrganizationUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !s.Config.OrganizationsEnabled {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, "ORGANIZATIONS_DISABLED"))
			return
		}
		ctx := r.Context()
		vars := mux.Vars(r)
		OrgID := vars["orgId"]
		Limit, Offset := getLimitOffsetFromRequest(r)

		Teams := s.OrganizationDataSvc.OrganizationUserList(ctx, OrgID, Limit, Offset)

		s.Success(w, r, http.StatusOK, Teams, nil)
	}
}

// handleCreateOrganizationTeam handles creating an organization team
// @Summary      Create Organization Team
// @Description  Create organization team with current user as admin
// @Tags         organization
// @Produce      json
// @Param        orgId  path    string                 true  "organization id"
// @Param        team   body    teamCreateRequestBody  true  "new team object"
// @Success      200    object  standardJsonResponse{data=thunderdome.Team}
// @Failure      403    object  standardJsonResponse{}
// @Failure      500    object  standardJsonResponse{}
// @Security     ApiKeyAuth
// @Router       /organizations/{orgId}/teams [post]
func (s *Service) handleCreateOrganizationTeam() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !s.Config.OrganizationsEnabled {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, "ORGANIZATIONS_DISABLED"))
			return
		}
		ctx := r.Context()
		SessionUserID := ctx.Value(contextKeyUserID).(string)
		vars := mux.Vars(r)
		OrgID := vars["orgId"]

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

		NewTeam, err := s.OrganizationDataSvc.OrganizationTeamCreate(ctx, OrgID, team.Name)
		s.Logger.Ctx(ctx).Error(
			"handleCreateOrganizationTeam error", zap.Error(err), zap.String("organization_id", OrgID),
			zap.String("session_user_id", SessionUserID), zap.String("team_name", team.Name))
		if err != nil {
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, NewTeam, nil)
	}
}

// handleOrganizationAddUser handles adding user to an organization
// @Summary      Add Org User
// @Description  Add user to organization
// @Tags         organization
// @Produce      json
// @Param        orgId  path    string                  true  "organization id"
// @Param        user   body    teamAddUserRequestBody  true  "new organization user object"
// @Success      200    object  standardJsonResponse{}
// @Failure      403    object  standardJsonResponse{}
// @Failure      500    object  standardJsonResponse{}
// @Security     ApiKeyAuth
// @Router       /organizations/{orgId}/users [post]
func (s *Service) handleOrganizationAddUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !s.Config.OrganizationsEnabled {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, "ORGANIZATIONS_DISABLED"))
			return
		}
		ctx := r.Context()
		SessionUserID := ctx.Value(contextKeyUserID).(string)
		vars := mux.Vars(r)
		OrgID := vars["orgId"]

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

		UserEmail := strings.ToLower(u.Email)

		User, UserErr := s.UserDataSvc.GetUserByEmail(ctx, UserEmail)
		if UserErr != nil && errors.Is(UserErr, sql.ErrNoRows) {
			inviteID, inviteErr := s.OrganizationDataSvc.OrganizationInviteUser(ctx, OrgID, UserEmail, u.Role)
			if inviteErr != nil {
				s.Logger.Ctx(ctx).Error("handleOrganizationAddUser error", zap.Error(inviteErr),
					zap.String("organization_id", OrgID), zap.String("session_user_id", SessionUserID))
				s.Failure(w, r, http.StatusInternalServerError, UserErr)
				return
			}
			org, orgErr := s.OrganizationDataSvc.OrganizationGet(ctx, OrgID)
			if orgErr != nil {
				s.Logger.Ctx(ctx).Error("handleOrganizationAddUser error", zap.Error(orgErr),
					zap.String("organization_id", OrgID), zap.String("session_user_id", SessionUserID))
				s.Failure(w, r, http.StatusInternalServerError, orgErr)
				return
			}
			emailErr := s.Email.SendOrganizationInvite(org.Name, UserEmail, inviteID)
			if emailErr != nil {
				s.Logger.Ctx(ctx).Error("handleOrganizationAddUser error", zap.Error(emailErr),
					zap.String("organization_id", OrgID), zap.String("session_user_id", SessionUserID))
				s.Failure(w, r, http.StatusInternalServerError, emailErr)
				return
			}
			s.Success(w, r, http.StatusOK, nil, nil)
			return
		} else if UserErr != nil {
			s.Logger.Ctx(ctx).Error(
				"handleOrganizationAddUser error", zap.Error(UserErr), zap.String("user_email", UserEmail),
				zap.String("session_user_id", SessionUserID), zap.String("organization_id", OrgID))
			s.Failure(w, r, http.StatusInternalServerError, Errorf(ENOTFOUND, "USER_NOT_FOUND"))
			return
		}

		_, err := s.OrganizationDataSvc.OrganizationAddUser(ctx, OrgID, User.Id, u.Role)
		if err != nil {
			s.Logger.Ctx(ctx).Error(
				"handleOrganizationAddUser error", zap.Error(err), zap.String("user_id", User.Id),
				zap.String("session_user_id", SessionUserID), zap.String("organization_id", OrgID),
				zap.String("user_role", u.Role))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, nil, nil)
	}
}

// handleOrganizationUpdateUser handles updating an organization user
// @Summary      Update Org User
// @Description  Update organization user
// @Tags         organization
// @Produce      json
// @Param        orgId  path    string                  true  "organization id"
// @Param        userId  path    string                  true  "user id"
// @Param        user   body    teamUpdateUserRequestBody  true  "organization user object"
// @Success      200    object  standardJsonResponse{}
// @Failure      403    object  standardJsonResponse{}
// @Failure      500    object  standardJsonResponse{}
// @Security     ApiKeyAuth
// @Router       /organizations/{orgId}/users/{userId} [put]
func (s *Service) handleOrganizationUpdateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !s.Config.OrganizationsEnabled {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, "ORGANIZATIONS_DISABLED"))
			return
		}
		ctx := r.Context()
		SessionUserID := ctx.Value(contextKeyUserID).(string)
		vars := mux.Vars(r)
		OrgID := vars["orgId"]
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

		_, err := s.OrganizationDataSvc.OrganizationUpdateUser(ctx, OrgID, UserID, u.Role)
		if err != nil {
			s.Logger.Ctx(ctx).Error(
				"handleOrganizationAddUser error", zap.Error(err), zap.String("user_id", UserID),
				zap.String("session_user_id", SessionUserID), zap.String("organization_id", OrgID),
				zap.String("user_role", u.Role))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, nil, nil)
	}
}

// handleOrganizationRemoveUser handles removing user from an organization (including departments, teams)
// @Summary      Remove Org User
// @Description  Remove user from organization including departments and teams
// @Tags         organization
// @Produce      json
// @Param        orgId   path    string  true  "organization id"
// @Param        userId  path    string  true  "user id"
// @Success      200     object  standardJsonResponse{}
// @Failure      403     object  standardJsonResponse{}
// @Failure      500     object  standardJsonResponse{}
// @Security     ApiKeyAuth
// @Router       /organizations/{orgId}/users/{userId} [delete]
func (s *Service) handleOrganizationRemoveUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !s.Config.OrganizationsEnabled {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, "ORGANIZATIONS_DISABLED"))
			return
		}
		ctx := r.Context()
		SessionUserID := ctx.Value(contextKeyUserID).(string)
		vars := mux.Vars(r)
		OrgID := vars["orgId"]
		UserID := vars["userId"]
		idErr := validate.Var(UserID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		err := s.OrganizationDataSvc.OrganizationRemoveUser(ctx, OrgID, UserID)
		if err != nil {
			s.Logger.Ctx(ctx).Error(
				"handleOrganizationRemoveUser error", zap.Error(err), zap.String("user_id", UserID),
				zap.String("session_user_id", SessionUserID), zap.String("organization_id", OrgID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, nil, nil)
	}
}

// handleGetOrganizationTeamByUser gets a team with users roles
// @Summary      Get Organization Team
// @Description  Get an organizations team with users roles
// @Tags         organization
// @Produce      json
// @Param        orgId   path    string  true  "organization id"
// @Param        teamId  path    string  true  "team id"
// @Success      200     object  standardJsonResponse{data=orgTeamResponse}
// @Failure      403     object  standardJsonResponse{}
// @Failure      500     object  standardJsonResponse{}
// @Security     ApiKeyAuth
// @Router       /organizations/{orgId}/teams/{teamId} [get]
func (s *Service) handleGetOrganizationTeamByUser() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		if !s.Config.OrganizationsEnabled {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, "ORGANIZATIONS_DISABLED"))
			return
		}
		ctx := r.Context()
		OrgRole := ctx.Value(contextKeyOrgRole).(string)
		TeamRole := ctx.Value(contextKeyTeamRole).(string)
		SessionUserID := ctx.Value(contextKeyUserID).(string)
		vars := mux.Vars(r)
		OrgID := vars["orgId"]
		TeamID := vars["teamId"]

		Organization, err := s.OrganizationDataSvc.OrganizationGet(ctx, OrgID)
		if err != nil {
			s.Logger.Ctx(ctx).Error(
				"handleGetOrganizationTeamByUser error", zap.Error(err), zap.String("organization_id", OrgID),
				zap.String("session_user_id", SessionUserID), zap.String("organization_role", OrgRole),
				zap.String("team_role", TeamRole), zap.String("team_id", TeamID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		Team, err := s.TeamDataSvc.TeamGet(ctx, TeamID)
		if err != nil {
			s.Logger.Ctx(ctx).Error(
				"handleGetOrganizationTeamByUser error", zap.Error(err), zap.String("organization_id", OrgID),
				zap.String("session_user_id", SessionUserID), zap.String("organization_role", OrgRole),
				zap.String("team_role", TeamRole), zap.String("team_id", TeamID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		result := &orgTeamResponse{
			Organization:     Organization,
			Team:             Team,
			OrganizationRole: OrgRole,
			TeamRole:         TeamRole,
		}

		s.Success(w, r, http.StatusOK, result, nil)
	}
}

// handleOrganizationTeamAddUser handles adding user to a team so long as they are in the organization
// @Summary      Add Org Team User
// @Description  Add user to organization team as long as they are already in the organization
// @Tags         organization
// @Produce      json
// @Param        orgId   path    string                  true  "organization id"
// @Param        teamId  path    string                  true  "team id"
// @Param        user    body    teamAddUserRequestBody  true  "new team user object"
// @Success      200     object  standardJsonResponse{}
// @Failure      403     object  standardJsonResponse{}
// @Failure      500     object  standardJsonResponse{}
// @Security     ApiKeyAuth
// @Router       /organizations/{orgId}/teams/{teamId}/users [post]
func (s *Service) handleOrganizationTeamAddUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !s.Config.OrganizationsEnabled {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, "ORGANIZATIONS_DISABLED"))
			return
		}

		ctx := r.Context()
		SessionUserID := ctx.Value(contextKeyUserID).(string)
		vars := mux.Vars(r)
		OrgID := vars["orgId"]
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

		UserEmail := u.Email

		User, UserErr := s.UserDataSvc.GetUserByEmail(ctx, UserEmail)
		if UserErr != nil {
			s.Logger.Ctx(ctx).Error(
				"handleOrganizationTeamAddUser error", zap.Error(UserErr), zap.String("organization_id", OrgID),
				zap.String("session_user_id", SessionUserID), zap.String("team_role", u.Role),
				zap.String("team_id", TeamID), zap.String("user_id", User.Id))
			s.Failure(w, r, http.StatusInternalServerError, Errorf(ENOTFOUND, "USER_NOT_FOUND"))
			return
		}

		OrgRole, roleErr := s.OrganizationDataSvc.OrganizationUserRole(ctx, User.Id, OrgID)
		if OrgRole == "" || roleErr != nil {
			s.Logger.Ctx(ctx).Error(
				"handleOrganizationTeamAddUser error", zap.Error(roleErr), zap.String("organization_id", OrgID),
				zap.String("session_user_id", SessionUserID), zap.String("organization_role", OrgRole),
				zap.String("team_role", u.Role), zap.String("team_id", TeamID), zap.String("user_id", User.Id))
			s.Failure(w, r, http.StatusInternalServerError, Errorf(EUNAUTHORIZED, "ORGANIZATION_USER_REQUIRED"))
			return
		}

		_, err := s.TeamDataSvc.TeamAddUser(ctx, TeamID, User.Id, u.Role)
		if err != nil {
			s.Logger.Ctx(ctx).Error(
				"handleOrganizationTeamAddUser error", zap.Error(err), zap.String("organization_id", OrgID),
				zap.String("session_user_id", SessionUserID), zap.String("organization_role", OrgRole),
				zap.String("team_role", u.Role), zap.String("team_id", TeamID), zap.String("user_id", User.Id))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, nil, nil)
	}
}

// handleDeleteOrganization handles deleting an organization
// @Summary      Delete Organization
// @Description  Delete an Organization
// @Tags         organization
// @Produce      json
// @Param        orgId  path    string  true  "the organization ID"
// @Success      200    object  standardJsonResponse{}
// @Success      403    object  standardJsonResponse{}
// @Success      500    object  standardJsonResponse{}
// @Security     ApiKeyAuth
// @Router       /organizations/{orgId} [delete]
func (s *Service) handleDeleteOrganization() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		SessionUserID := ctx.Value(contextKeyUserID).(string)
		vars := mux.Vars(r)
		OrgID := vars["orgId"]
		idErr := validate.Var(OrgID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		err := s.OrganizationDataSvc.OrganizationDelete(ctx, OrgID)
		if err != nil {
			s.Logger.Ctx(ctx).Error(
				"handleDeleteOrganization error", zap.Error(err), zap.String("organization_id", OrgID),
				zap.String("session_user_id", SessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, nil, nil)
	}
}
