package api

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/StevenWeathers/thunderdome-planning-poker/pkg/database"
	"github.com/gorilla/mux"
)

// handleGetOrganizationsByUser gets a list of organizations the user is apart of
// @Summary Get Users Organizations
// @Description get list of organizations for the authenticated user
// @Tags organizations
// @Produce  json
// @Param limit path int false "Max number of results to return"
// @Param offset path int false "Starting point to return rows from, should be multiplied by limit or 0"
// @Success 200
// @Router /organizations/{limit}/{offset} [get]
func (a *api) handleGetOrganizationsByUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		UserID := r.Context().Value(contextKeyUserID).(string)
		vars := mux.Vars(r)
		Limit, _ := strconv.Atoi(vars["limit"])
		Offset, _ := strconv.Atoi(vars["offset"])

		Organizations := a.db.OrganizationListByUser(UserID, Limit, Offset)

		a.respondWithJSON(w, http.StatusOK, Organizations)
	}
}

// handleGetOrganizationByUser gets an organization with user role
func (a *api) handleGetOrganizationByUser() http.HandlerFunc {
	type OrganizationResponse struct {
		Organization *database.Organization `json:"organization"`
		Role         string                 `json:"role"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		OrgRole := r.Context().Value(contextKeyOrgRole).(string)
		vars := mux.Vars(r)
		OrgID := vars["orgId"]

		Organization, err := a.db.OrganizationGet(OrgID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		a.respondWithJSON(w, http.StatusOK, &OrganizationResponse{
			Organization: Organization,
			Role:         OrgRole,
		})
	}
}

// handleCreateOrganization handles creating an organization with current user as admin
func (a *api) handleCreateOrganization() http.HandlerFunc {
	type CreateOrgResponse struct {
		OrganizationID string `json:"id"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		UserID := r.Context().Value(contextKeyUserID).(string)
		keyVal := a.getJSONRequestBody(r, w)

		OrgName := keyVal["name"].(string)
		OrgId, err := a.db.OrganizationCreate(UserID, OrgName)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var NewOrg = &CreateOrgResponse{
			OrganizationID: OrgId,
		}

		a.respondWithJSON(w, http.StatusOK, NewOrg)
	}
}

// handleGetOrganizationTeams gets a list of teams associated to the organization
func (a *api) handleGetOrganizationTeams() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		OrgID := vars["orgId"]
		Limit, _ := strconv.Atoi(vars["limit"])
		Offset, _ := strconv.Atoi(vars["offset"])

		Teams := a.db.OrganizationTeamList(OrgID, Limit, Offset)

		a.respondWithJSON(w, http.StatusOK, Teams)
	}
}

// handleGetOrganizationUsers gets a list of users associated to the organization
func (a *api) handleGetOrganizationUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		OrgID := vars["orgId"]
		Limit, _ := strconv.Atoi(vars["limit"])
		Offset, _ := strconv.Atoi(vars["offset"])

		Teams := a.db.OrganizationUserList(OrgID, Limit, Offset)

		a.respondWithJSON(w, http.StatusOK, Teams)
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
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var NewTeam = &CreateTeamResponse{
			TeamID: TeamID,
		}

		a.respondWithJSON(w, http.StatusOK, NewTeam)
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
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, err := a.db.OrganizationAddUser(OrgID, User.UserID, Role)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		return
	}
}

// handleOrganizationRemoveUser handles removing user from an organization (including departments, teams)
func (a *api) handleOrganizationRemoveUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		keyVal := a.getJSONRequestBody(r, w)

		vars := mux.Vars(r)
		OrgID := vars["orgId"]
		UserID := keyVal["id"].(string)

		err := a.db.OrganizationRemoveUser(OrgID, UserID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		return
	}
}

// handleGetOrganizationTeamByUser gets a team with users roles
func (a *api) handleGetOrganizationTeamByUser() http.HandlerFunc {
	type TeamResponse struct {
		Organization     *database.Organization `json:"organization"`
		Team             *database.Team         `json:"team"`
		OrganizationRole string                 `json:"organizationRole"`
		TeamRole         string                 `json:"teamRole"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		OrgRole := r.Context().Value(contextKeyOrgRole).(string)
		TeamRole := r.Context().Value(contextKeyTeamRole).(string)
		vars := mux.Vars(r)
		OrgID := vars["orgId"]
		TeamID := vars["teamId"]

		Organization, err := a.db.OrganizationGet(OrgID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		Team, err := a.db.TeamGet(TeamID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		a.respondWithJSON(w, http.StatusOK, &TeamResponse{
			Organization:     Organization,
			Team:             Team,
			OrganizationRole: OrgRole,
			TeamRole:         TeamRole,
		})
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
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		OrgRole, roleErr := a.db.OrganizationUserRole(User.UserID, OrgID)
		if OrgRole == "" || roleErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, err := a.db.TeamAddUser(TeamID, User.UserID, Role)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		return
	}
}
