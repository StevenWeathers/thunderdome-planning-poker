package api

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/StevenWeathers/thunderdome-planning-poker/model"
	"github.com/gorilla/mux"
)

type organizationResponse struct {
	Organization *model.Organization `json:"organization"`
	Role         string              `json:"role"`
}

type orgTeamResponse struct {
	Organization     *model.Organization `json:"organization"`
	Team             *model.Team         `json:"team"`
	OrganizationRole string              `json:"organizationRole"`
	TeamRole         string              `json:"teamRole"`
}

// handleGetOrganizationsByUser gets a list of organizations the user is a part of
// @Summary Get Users Organizations
// @Description Get list of organizations for the authenticated user
// @Tags organization
// @Produce  json
// @Param userId path string true "the user ID to get organizations for"
// @Param limit query int false "Max number of results to return"
// @Param offset query int false "Starting point to return rows from, should be multiplied by limit or 0"
// @Success 200 object standardJsonResponse{data=[]model.Organization}
// @Failure 403 object standardJsonResponse{}
// @Security ApiKeyAuth
// @Router /users/{userId}/organizations [get]
func (a *api) handleGetOrganizationsByUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !a.config.OrganizationsEnabled {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, "ORGANIZATIONS_DISABLED"))
			return
		}
		vars := mux.Vars(r)
		UserID := vars["userId"]

		Limit, Offset := getLimitOffsetFromRequest(r)

		Organizations := a.db.OrganizationListByUser(r.Context(), UserID, Limit, Offset)

		a.Success(w, r, http.StatusOK, Organizations, nil)
	}
}

// handleGetOrganizationByUser gets an organization with user role
// @Summary Get Organization
// @Description Get an organization with user role
// @Tags organization
// @Produce  json
// @Param orgId path string true "organization id"
// @Success 200 object standardJsonResponse{data=organizationResponse}
// @Failure 403 object standardJsonResponse{}
// @Failure 500 object standardJsonResponse{}
// @Security ApiKeyAuth
// @Router /organizations/{orgId} [get]
func (a *api) handleGetOrganizationByUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !a.config.OrganizationsEnabled {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, "ORGANIZATIONS_DISABLED"))
			return
		}
		ctx := r.Context()
		OrgRole := ctx.Value(contextKeyOrgRole).(string)
		vars := mux.Vars(r)
		OrgID := vars["orgId"]

		Organization, err := a.db.OrganizationGet(ctx, OrgID)
		if err != nil {
			a.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		result := &organizationResponse{
			Organization: Organization,
			Role:         OrgRole,
		}

		a.Success(w, r, http.StatusOK, result, nil)
	}
}

// handleCreateOrganization handles creating an organization with current user as admin
// @Summary Create Organization
// @Description Create organization with current user as admin
// @Tags organization
// @Produce  json
// @Param userId path string true "user id"
// @Param organization body teamCreateRequestBody true "new organization object"
// @Success 200 object standardJsonResponse{data=model.Organization}
// @Failure 403 object standardJsonResponse{}
// @Failure 500 object standardJsonResponse{}
// @Security ApiKeyAuth
// @Router /users/{userId}/organizations [post]
func (a *api) handleCreateOrganization() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !a.config.OrganizationsEnabled {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, "ORGANIZATIONS_DISABLED"))
			return
		}
		vars := mux.Vars(r)
		UserID := vars["userId"]

		var team = teamCreateRequestBody{}
		body, bodyErr := io.ReadAll(r.Body)
		if bodyErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
			return
		}

		jsonErr := json.Unmarshal(body, &team)
		if jsonErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		Organization, err := a.db.OrganizationCreate(r.Context(), UserID, team.Name)
		if err != nil {
			a.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		a.Success(w, r, http.StatusOK, Organization, nil)
	}
}

// handleGetOrganizationTeams gets a list of teams associated to the organization
// @Summary Get Organization Teams
// @Description Get a list of organization teams
// @Tags organization
// @Produce  json
// @Param orgId path string true "organization id"
// @Success 200 object standardJsonResponse{data=[]model.Team}
// @Failure 403 object standardJsonResponse{}
// @Security ApiKeyAuth
// @Router /organizations/{orgId}/teams [get]
func (a *api) handleGetOrganizationTeams() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !a.config.OrganizationsEnabled {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, "ORGANIZATIONS_DISABLED"))
			return
		}
		vars := mux.Vars(r)
		OrgID := vars["orgId"]
		Limit, Offset := getLimitOffsetFromRequest(r)

		Teams := a.db.OrganizationTeamList(r.Context(), OrgID, Limit, Offset)

		a.Success(w, r, http.StatusOK, Teams, nil)
	}
}

// handleGetOrganizationUsers gets a list of users associated to the organization
// @Summary Get Organization Users
// @Description get a list of organization users
// @Tags organization
// @Produce  json
// @Param orgId path string true "organization id"
// @Success 200 object standardJsonResponse{data=[]model.User}
// @Failure 403 object standardJsonResponse{}
// @Security ApiKeyAuth
// @Router /organizations/{orgId}/users [get]
func (a *api) handleGetOrganizationUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !a.config.OrganizationsEnabled {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, "ORGANIZATIONS_DISABLED"))
			return
		}
		vars := mux.Vars(r)
		OrgID := vars["orgId"]
		Limit, Offset := getLimitOffsetFromRequest(r)

		Teams := a.db.OrganizationUserList(r.Context(), OrgID, Limit, Offset)

		a.Success(w, r, http.StatusOK, Teams, nil)
	}
}

// handleCreateOrganizationTeam handles creating an organization team
// @Summary Create Organization Team
// @Description Create organization team with current user as admin
// @Tags organization
// @Produce  json
// @Param orgId path string true "organization id"
// @Param team body teamCreateRequestBody true "new team object"
// @Success 200 object standardJsonResponse{data=model.Team}
// @Failure 403 object standardJsonResponse{}
// @Failure 500 object standardJsonResponse{}
// @Security ApiKeyAuth
// @Router /organizations/{orgId}/teams [post]
func (a *api) handleCreateOrganizationTeam() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !a.config.OrganizationsEnabled {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, "ORGANIZATIONS_DISABLED"))
			return
		}
		vars := mux.Vars(r)
		OrgID := vars["orgId"]

		var team = teamCreateRequestBody{}
		body, bodyErr := io.ReadAll(r.Body)
		if bodyErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
			return
		}

		jsonErr := json.Unmarshal(body, &team)
		if jsonErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		NewTeam, err := a.db.OrganizationTeamCreate(r.Context(), OrgID, team.Name)
		if err != nil {
			a.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		a.Success(w, r, http.StatusOK, NewTeam, nil)
	}
}

// handleOrganizationAddUser handles adding user to an organization
// @Summary Add Org User
// @Description Add user to organization
// @Tags organization
// @Produce  json
// @Param orgId path string true "organization id"
// @Param user body teamAddUserRequestBody true "new organization user object"
// @Success 200 object standardJsonResponse{}
// @Failure 403 object standardJsonResponse{}
// @Failure 500 object standardJsonResponse{}
// @Security ApiKeyAuth
// @Router /organizations/{orgId}/users [post]
func (a *api) handleOrganizationAddUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !a.config.OrganizationsEnabled {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, "ORGANIZATIONS_DISABLED"))
			return
		}

		vars := mux.Vars(r)
		OrgID := vars["orgId"]

		var u = teamAddUserRequestBody{}
		body, bodyErr := io.ReadAll(r.Body)
		if bodyErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
			return
		}

		jsonErr := json.Unmarshal(body, &u)
		if jsonErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		UserEmail := u.Email

		User, UserErr := a.db.GetUserByEmail(r.Context(), UserEmail)
		if UserErr != nil {
			a.Failure(w, r, http.StatusInternalServerError, Errorf(ENOTFOUND, "USER_NOT_FOUND"))
			return
		}

		_, err := a.db.OrganizationAddUser(r.Context(), OrgID, User.Id, u.Role)
		if err != nil {
			a.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		a.Success(w, r, http.StatusOK, nil, nil)
	}
}

// handleOrganizationRemoveUser handles removing user from an organization (including departments, teams)
// @Summary Remove Org User
// @Description Remove user from organization including departments and teams
// @Tags organization
// @Produce  json
// @Param orgId path string true "organization id"
// @Param userId path string true "user id"
// @Success 200 object standardJsonResponse{}
// @Failure 403 object standardJsonResponse{}
// @Failure 500 object standardJsonResponse{}
// @Security ApiKeyAuth
// @Router /organizations/{orgId}/users/{userId} [delete]
func (a *api) handleOrganizationRemoveUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !a.config.OrganizationsEnabled {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, "ORGANIZATIONS_DISABLED"))
			return
		}
		vars := mux.Vars(r)
		OrgID := vars["orgId"]
		UserID := vars["userId"]
		idErr := validate.Var(UserID, "required,uuid")
		if idErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		err := a.db.OrganizationRemoveUser(r.Context(), OrgID, UserID)
		if err != nil {
			a.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		a.Success(w, r, http.StatusOK, nil, nil)
	}
}

// handleGetOrganizationTeamByUser gets a team with users roles
// @Summary Get Organization Team
// @Description Get an organizations team with users roles
// @Tags organization
// @Produce  json
// @Param orgId path string true "organization id"
// @Param teamId path string true "team id"
// @Success 200 object standardJsonResponse{data=orgTeamResponse}
// @Failure 403 object standardJsonResponse{}
// @Failure 500 object standardJsonResponse{}
// @Security ApiKeyAuth
// @Router /organizations/{orgId}/teams/{teamId} [get]
func (a *api) handleGetOrganizationTeamByUser() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		if !a.config.OrganizationsEnabled {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, "ORGANIZATIONS_DISABLED"))
			return
		}
		ctx := r.Context()
		OrgRole := ctx.Value(contextKeyOrgRole).(string)
		TeamRole := ctx.Value(contextKeyTeamRole).(string)
		vars := mux.Vars(r)
		OrgID := vars["orgId"]
		TeamID := vars["teamId"]

		Organization, err := a.db.OrganizationGet(r.Context(), OrgID)
		if err != nil {
			a.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		Team, err := a.db.TeamGet(ctx, TeamID)
		if err != nil {
			a.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		result := &orgTeamResponse{
			Organization:     Organization,
			Team:             Team,
			OrganizationRole: OrgRole,
			TeamRole:         TeamRole,
		}

		a.Success(w, r, http.StatusOK, result, nil)
	}
}

// handleOrganizationTeamAddUser handles adding user to a team so long as they are in the organization
// @Summary Add Org Team User
// @Description Add user to organization team as long as they are already in the organization
// @Tags organization
// @Produce  json
// @Param orgId path string true "organization id"
// @Param teamId path string true "team id"
// @Param user body teamAddUserRequestBody true "new team user object"
// @Success 200 object standardJsonResponse{}
// @Failure 403 object standardJsonResponse{}
// @Failure 500 object standardJsonResponse{}
// @Security ApiKeyAuth
// @Router /organizations/{orgId}/teams/{teamId}/users [post]
func (a *api) handleOrganizationTeamAddUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !a.config.OrganizationsEnabled {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, "ORGANIZATIONS_DISABLED"))
			return
		}

		ctx := r.Context()
		vars := mux.Vars(r)
		OrgID := vars["orgId"]
		TeamID := vars["teamId"]

		var u = teamAddUserRequestBody{}
		body, bodyErr := io.ReadAll(r.Body)
		if bodyErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
			return
		}

		jsonErr := json.Unmarshal(body, &u)
		if jsonErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		UserEmail := u.Email

		User, UserErr := a.db.GetUserByEmail(ctx, UserEmail)
		if UserErr != nil {
			a.Failure(w, r, http.StatusInternalServerError, Errorf(ENOTFOUND, "USER_NOT_FOUND"))
			return
		}

		OrgRole, roleErr := a.db.OrganizationUserRole(r.Context(), User.Id, OrgID)
		if OrgRole == "" || roleErr != nil {
			a.Failure(w, r, http.StatusInternalServerError, Errorf(EUNAUTHORIZED, "ORGANIZATION_USER_REQUIRED"))
			return
		}

		_, err := a.db.TeamAddUser(ctx, TeamID, User.Id, u.Role)
		if err != nil {
			a.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		a.Success(w, r, http.StatusOK, nil, nil)
	}
}

// handleDeleteOrganization handles deleting an organization
// @Summary Delete Organization
// @Description Delete an Organization
// @Tags organization
// @Produce  json
// @Param orgId path string true "the organization ID"
// @Success 200 object standardJsonResponse{}
// @Success 403 object standardJsonResponse{}
// @Success 500 object standardJsonResponse{}
// @Security ApiKeyAuth
// @Router /organizations/{orgId} [delete]
func (a *api) handleDeleteOrganization() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		OrgID := vars["orgId"]
		idErr := validate.Var(OrgID, "required,uuid")
		if idErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		err := a.db.OrganizationDelete(r.Context(), OrgID)
		if err != nil {
			a.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		a.Success(w, r, http.StatusOK, nil, nil)
	}
}
