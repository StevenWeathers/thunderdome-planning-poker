package api

import (
	"net/http"
	"strings"

	"github.com/StevenWeathers/thunderdome-planning-poker/model"
	"github.com/gorilla/mux"
)

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
			a.respondWithStandardJSON(w, http.StatusForbidden, false, nil, nil, nil)
			return
		}

		Limit, Offset := a.getLimitOffsetFromRequest(r, w)

		Organizations := a.db.OrganizationListByUser(UserID, Limit, Offset)

		a.respondWithStandardJSON(w, http.StatusOK, true, nil, Organizations, nil)
	}
}

// handleGetOrganizationByUser gets an organization with user role
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
			errors := make([]string, 0)
			errors = append(errors, err.Error())
			a.respondWithStandardJSON(w, http.StatusInternalServerError, false, errors, nil, nil)
			return
		}

		result := &OrganizationResponse{
			Organization: Organization,
			Role:         OrgRole,
		}

		a.respondWithStandardJSON(w, http.StatusOK, true, nil, result, nil)
	}
}

// handleCreateOrganization handles creating an organization with current user as admin
func (a *api) handleCreateOrganization() http.HandlerFunc {
	type CreateOrgResponse struct {
		OrganizationID string `json:"id"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		UserID := vars["id"]
		AuthedUserID := r.Context().Value(contextKeyUserID).(string)

		if UserID != AuthedUserID {
			a.respondWithStandardJSON(w, http.StatusForbidden, false, nil, nil, nil)
			return
		}

		keyVal := a.getJSONRequestBody(r, w)

		OrgName := keyVal["name"].(string)
		OrgId, err := a.db.OrganizationCreate(UserID, OrgName)
		if err != nil {
			errors := make([]string, 0)
			errors = append(errors, err.Error())
			a.respondWithStandardJSON(w, http.StatusInternalServerError, false, errors, nil, nil)
			return
		}

		var NewOrg = &CreateOrgResponse{
			OrganizationID: OrgId,
		}

		a.respondWithStandardJSON(w, http.StatusOK, true, nil, NewOrg, nil)
	}
}

// handleGetOrganizationTeams gets a list of teams associated to the organization
func (a *api) handleGetOrganizationTeams() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		OrgID := vars["orgId"]
		Limit, Offset := a.getLimitOffsetFromRequest(r, w)

		Teams := a.db.OrganizationTeamList(OrgID, Limit, Offset)

		a.respondWithStandardJSON(w, http.StatusOK, true, nil, Teams, nil)
	}
}

// handleGetOrganizationUsers gets a list of users associated to the organization
func (a *api) handleGetOrganizationUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		OrgID := vars["orgId"]
		Limit, Offset := a.getLimitOffsetFromRequest(r, w)

		Teams := a.db.OrganizationUserList(OrgID, Limit, Offset)

		a.respondWithStandardJSON(w, http.StatusOK, true, nil, Teams, nil)
	}
}

// handleCreateOrganizationTeam handles creating an organization team
func (a *api) handleCreateOrganizationTeam() http.HandlerFunc {
	type CreateTeamResponse struct {
		TeamID string `json:"id"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		// UserID := r.Context().Value(contextKeyUserID).(string)
		vars := mux.Vars(r)
		keyVal := a.getJSONRequestBody(r, w)

		TeamName := keyVal["name"].(string)
		OrgID := vars["orgId"]
		TeamID, err := a.db.OrganizationTeamCreate(OrgID, TeamName)
		if err != nil {
			errors := make([]string, 0)
			errors = append(errors, err.Error())
			a.respondWithStandardJSON(w, http.StatusInternalServerError, false, errors, nil, nil)
			return
		}

		var NewTeam = &CreateTeamResponse{
			TeamID: TeamID,
		}

		a.respondWithStandardJSON(w, http.StatusOK, true, nil, NewTeam, nil)
	}
}

// handleOrganizationAddUser handles adding user to an organization
func (a *api) handleOrganizationAddUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		keyVal := a.getJSONRequestBody(r, w)

		vars := mux.Vars(r)
		OrgID := vars["orgId"]
		UserEmail := strings.ToLower(keyVal["email"].(string))
		Role := keyVal["role"].(string)

		User, UserErr := a.db.GetUserByEmail(UserEmail)
		if UserErr != nil {
			errors := make([]string, 0)
			errors = append(errors, UserErr.Error())
			a.respondWithStandardJSON(w, http.StatusInternalServerError, false, errors, nil, nil)
			return
		}

		_, err := a.db.OrganizationAddUser(OrgID, User.UserID, Role)
		if err != nil {
			errors := make([]string, 0)
			errors = append(errors, err.Error())
			a.respondWithStandardJSON(w, http.StatusInternalServerError, false, errors, nil, nil)
			return
		}

		a.respondWithStandardJSON(w, http.StatusOK, true, nil, nil, nil)
	}
}

// handleOrganizationRemoveUser handles removing user from an organization (including departments, teams)
func (a *api) handleOrganizationRemoveUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		OrgID := vars["orgId"]
		UserID := vars["userId"]

		err := a.db.OrganizationRemoveUser(OrgID, UserID)
		if err != nil {
			errors := make([]string, 0)
			errors = append(errors, err.Error())
			a.respondWithStandardJSON(w, http.StatusInternalServerError, false, errors, nil, nil)
			return
		}

		a.respondWithStandardJSON(w, http.StatusOK, true, nil, nil, nil)
	}
}

// handleGetOrganizationTeamByUser gets a team with users roles
func (a *api) handleGetOrganizationTeamByUser() http.HandlerFunc {
	type TeamResponse struct {
		Organization     *model.Organization `json:"organization"`
		Team             *model.Team         `json:"team"`
		OrganizationRole string              `json:"organizationRole"`
		TeamRole         string              `json:"teamRole"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		OrgRole := r.Context().Value(contextKeyOrgRole).(string)
		TeamRole := r.Context().Value(contextKeyTeamRole).(string)
		vars := mux.Vars(r)
		OrgID := vars["orgId"]
		TeamID := vars["teamId"]

		Organization, err := a.db.OrganizationGet(OrgID)
		if err != nil {
			errors := make([]string, 0)
			errors = append(errors, err.Error())
			a.respondWithStandardJSON(w, http.StatusInternalServerError, false, errors, nil, nil)
			return
		}

		Team, err := a.db.TeamGet(TeamID)
		if err != nil {
			errors := make([]string, 0)
			errors = append(errors, err.Error())
			a.respondWithStandardJSON(w, http.StatusInternalServerError, false, errors, nil, nil)
			return
		}

		result := &TeamResponse{
			Organization:     Organization,
			Team:             Team,
			OrganizationRole: OrgRole,
			TeamRole:         TeamRole,
		}

		a.respondWithStandardJSON(w, http.StatusOK, true, nil, result, nil)
	}
}

// handleOrganizationTeamAddUser handles adding user to a team so long as they are in the organization
func (a *api) handleOrganizationTeamAddUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		keyVal := a.getJSONRequestBody(r, w)

		vars := mux.Vars(r)
		OrgID := vars["orgId"]
		TeamID := vars["teamId"]
		UserEmail := strings.ToLower(keyVal["email"].(string))
		Role := keyVal["role"].(string)

		User, UserErr := a.db.GetUserByEmail(UserEmail)
		if UserErr != nil {
			errors := make([]string, 0)
			errors = append(errors, UserErr.Error())
			a.respondWithStandardJSON(w, http.StatusInternalServerError, false, errors, nil, nil)
			return
		}

		OrgRole, roleErr := a.db.OrganizationUserRole(User.UserID, OrgID)
		if OrgRole == "" || roleErr != nil {
			errors := make([]string, 0)
			errors = append(errors, roleErr.Error())
			a.respondWithStandardJSON(w, http.StatusInternalServerError, false, errors, nil, nil)
			return
		}

		_, err := a.db.TeamAddUser(TeamID, User.UserID, Role)
		if err != nil {
			errors := make([]string, 0)
			errors = append(errors, err.Error())
			a.respondWithStandardJSON(w, http.StatusInternalServerError, false, errors, nil, nil)
			return
		}

		a.respondWithStandardJSON(w, http.StatusOK, true, nil, nil, nil)
	}
}
