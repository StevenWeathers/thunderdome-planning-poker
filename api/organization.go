package api

import (
	"net/http"
	"strings"

	"github.com/StevenWeathers/thunderdome-planning-poker/model"
	"github.com/gorilla/mux"
)

type createOrgResponse struct {
	OrganizationID string `json:"id"`
}

type orgTeamResponse struct {
	Organization     *model.Organization `json:"organization"`
	Team             *model.Team         `json:"team"`
	OrganizationRole string              `json:"organizationRole"`
	TeamRole         string              `json:"teamRole"`
}

// handleGetOrganizationsByUser gets a list of organizations the user is apart of
// @Summary Get Users Organizations
// @Description get list of organizations for the authenticated user
// @Tags organization
// @Produce  json
// @Param id path int false "the user ID to get organizations for"
// @Param limit query int true "Max number of results to return"
// @Param offset query int true "Starting point to return rows from, should be multiplied by limit or 0"
// @Success 200 object standardJsonResponse{data=[]model.Organization}
// @Failure 403 object standardJsonResponse{}
// @Router /users/{id}/organizations [get]
func (a *api) handleGetOrganizationsByUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		UserID := vars["id"]
		AuthedUserID := r.Context().Value(contextKeyUserID).(string)

		if UserID != AuthedUserID {
			Error(w, r, http.StatusForbidden, "")
			return
		}

		Limit, Offset := getLimitOffsetFromRequest(r, w)

		Organizations := a.db.OrganizationListByUser(UserID, Limit, Offset)

		Success(w, r, http.StatusOK, Organizations, nil)
	}
}

// handleGetOrganizationByUser gets an organization with user role
// @Summary Get Organization
// @Description get an organization with user role
// @Tags organization
// @Produce  json
// @Param id path int false "organization id"
// @Success 200 object standardJsonResponse{data=model.Organization}
// @Failure 403 object standardJsonResponse{}
// @Failure 500 object standardJsonResponse{}
// @Router /organizations/{id} [get]
func (a *api) handleGetOrganizationByUser() http.HandlerFunc {
	type OrganizationResponse struct {
		Organization *model.Organization `json:"organization"`
		Role         string              `json:"role"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		OrgRole := r.Context().Value(contextKeyOrgRole).(string)
		vars := mux.Vars(r)
		OrgID := vars["orgId"]

		Organization, err := a.db.OrganizationGet(OrgID)
		if err != nil {
			Error(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		result := &OrganizationResponse{
			Organization: Organization,
			Role:         OrgRole,
		}

		Success(w, r, http.StatusOK, result, nil)
	}
}

// handleCreateOrganization handles creating an organization with current user as admin
// @Summary Create Organization
// @Description Create organization with current user as admin
// @Tags organization
// @Produce  json
// @Param id path int false "user id"
// @Success 200 object standardJsonResponse{data=createOrgResponse}
// @Failure 403 object standardJsonResponse{}
// @Failure 500 object standardJsonResponse{}
// @Router /users/{id}/organizations [post]
func (a *api) handleCreateOrganization() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		UserID := vars["id"]
		AuthedUserID := r.Context().Value(contextKeyUserID).(string)

		if UserID != AuthedUserID {
			Error(w, r, http.StatusForbidden, "")
			return
		}

		keyVal := getJSONRequestBody(r, w)

		OrgName := keyVal["name"].(string)
		OrgId, err := a.db.OrganizationCreate(UserID, OrgName)
		if err != nil {
			Error(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		var NewOrg = &createOrgResponse{
			OrganizationID: OrgId,
		}

		Success(w, r, http.StatusOK, NewOrg, nil)
	}
}

// handleGetOrganizationTeams gets a list of teams associated to the organization
// @Summary Get Organization Teams
// @Description get a list of organization teams
// @Tags organization
// @Produce  json
// @Param id path int false "organization id"
// @Success 200 object standardJsonResponse{data=[]model.Team}
// @Failure 403 object standardJsonResponse{}
// @Router /organizations/{id}/teams [get]
func (a *api) handleGetOrganizationTeams() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		OrgID := vars["orgId"]
		Limit, Offset := getLimitOffsetFromRequest(r, w)

		Teams := a.db.OrganizationTeamList(OrgID, Limit, Offset)

		Success(w, r, http.StatusOK, Teams, nil)
	}
}

// handleGetOrganizationUsers gets a list of users associated to the organization
// @Summary Get Organization Users
// @Description get a list of organization users
// @Tags organization
// @Produce  json
// @Param id path int false "organization id"
// @Success 200 object standardJsonResponse{data=[]model.User}
// @Failure 403 object standardJsonResponse{}
// @Router /organizations/{id}/users [get]
func (a *api) handleGetOrganizationUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		OrgID := vars["orgId"]
		Limit, Offset := getLimitOffsetFromRequest(r, w)

		Teams := a.db.OrganizationUserList(OrgID, Limit, Offset)

		Success(w, r, http.StatusOK, Teams, nil)
	}
}

// handleCreateOrganizationTeam handles creating an organization team
// @Summary Create Organization Team
// @Description Create organization team with current user as admin
// @Tags organization
// @Produce  json
// @Param id path int false "organization id"
// @Param name body string false "team name"
// @Success 200 object standardJsonResponse{data=createTeamResponse}
// @Failure 403 object standardJsonResponse{}
// @Failure 500 object standardJsonResponse{}
// @Router /organizations/{id}/teams [post]
func (a *api) handleCreateOrganizationTeam() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// UserID := r.Context().Value(contextKeyUserID).(string)
		vars := mux.Vars(r)
		keyVal := getJSONRequestBody(r, w)

		TeamName := keyVal["name"].(string)
		OrgID := vars["orgId"]
		TeamID, err := a.db.OrganizationTeamCreate(OrgID, TeamName)
		if err != nil {
			Error(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		var NewTeam = &createTeamResponse{
			TeamID: TeamID,
		}

		Success(w, r, http.StatusOK, NewTeam, nil)
	}
}

// handleOrganizationAddUser handles adding user to an organization
// @Summary Add Org User
// @Description Add user to organization
// @Tags organization
// @Produce  json
// @Param id path int false "organization id"
// @Param email body int false "user email"
// @Param role body int false "user team role"
// @Success 200 object standardJsonResponse{data=createTeamResponse}
// @Failure 403 object standardJsonResponse{}
// @Failure 500 object standardJsonResponse{}
// @Router /organizations/{id}/users [post]
func (a *api) handleOrganizationAddUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		keyVal := getJSONRequestBody(r, w)

		vars := mux.Vars(r)
		OrgID := vars["orgId"]
		UserEmail := strings.ToLower(keyVal["email"].(string))
		Role := keyVal["role"].(string)

		User, UserErr := a.db.GetUserByEmail(UserEmail)
		if UserErr != nil {
			Error(w, r, http.StatusInternalServerError, "USER_NOT_FOUND")
			return
		}

		_, err := a.db.OrganizationAddUser(OrgID, User.UserID, Role)
		if err != nil {
			Error(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		Success(w, r, http.StatusOK, nil, nil)
	}
}

// handleOrganizationRemoveUser handles removing user from an organization (including departments, teams)
// @Summary Remove Org User
// @Description Remove user from organization including departments and teams
// @Tags organization
// @Produce  json
// @Param id path int false "organization id"
// @Param userId path int false "user id"
// @Success 200 object standardJsonResponse{}
// @Failure 403 object standardJsonResponse{}
// @Failure 500 object standardJsonResponse{}
// @Router /organizations/{id}/users/{userId} [delete]
func (a *api) handleOrganizationRemoveUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		OrgID := vars["orgId"]
		UserID := vars["userId"]

		err := a.db.OrganizationRemoveUser(OrgID, UserID)
		if err != nil {
			Error(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		Success(w, r, http.StatusOK, nil, nil)
	}
}

// handleGetOrganizationTeamByUser gets a team with users roles
// @Summary Get Organization Team
// @Description Get an organizations team with users roles
// @Tags organization
// @Produce  json
// @Param id path int false "organization id"
// @Param teamId path int false "team id"
// @Success 200 object standardJsonResponse{data=orgTeamResponse}
// @Failure 403 object standardJsonResponse{}
// @Failure 500 object standardJsonResponse{}
// @Router /organizations/{id}/teams/{teamId} [get]
func (a *api) handleGetOrganizationTeamByUser() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		OrgRole := r.Context().Value(contextKeyOrgRole).(string)
		TeamRole := r.Context().Value(contextKeyTeamRole).(string)
		vars := mux.Vars(r)
		OrgID := vars["orgId"]
		TeamID := vars["teamId"]

		Organization, err := a.db.OrganizationGet(OrgID)
		if err != nil {
			Error(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		Team, err := a.db.TeamGet(TeamID)
		if err != nil {
			Error(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		result := &orgTeamResponse{
			Organization:     Organization,
			Team:             Team,
			OrganizationRole: OrgRole,
			TeamRole:         TeamRole,
		}

		Success(w, r, http.StatusOK, result, nil)
	}
}

// handleOrganizationTeamAddUser handles adding user to a team so long as they are in the organization
// @Summary Add Org Team User
// @Description Add user to organization team as long as they are already in the organization
// @Tags organization
// @Produce  json
// @Param id path int false "organization id"
// @Param teamId path int false "team id"
// @Param email body string false "user email"
// @Param role body string false "user team role"
// @Success 200 object standardJsonResponse{}
// @Failure 403 object standardJsonResponse{}
// @Failure 500 object standardJsonResponse{}
// @Router /organizations/{id}/teams/{teamId}/users [post]
func (a *api) handleOrganizationTeamAddUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		keyVal := getJSONRequestBody(r, w)

		vars := mux.Vars(r)
		OrgID := vars["orgId"]
		TeamID := vars["teamId"]
		UserEmail := strings.ToLower(keyVal["email"].(string))
		Role := keyVal["role"].(string)

		User, UserErr := a.db.GetUserByEmail(UserEmail)
		if UserErr != nil {
			Error(w, r, http.StatusInternalServerError, "USER_NOT_FOUND")
			return
		}

		OrgRole, roleErr := a.db.OrganizationUserRole(User.UserID, OrgID)
		if OrgRole == "" || roleErr != nil {
			Error(w, r, http.StatusInternalServerError, "ORGANIZATION_USER_REQUIRED")
			return
		}

		_, err := a.db.TeamAddUser(TeamID, User.UserID, Role)
		if err != nil {
			Error(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		Success(w, r, http.StatusOK, nil, nil)
	}
}
