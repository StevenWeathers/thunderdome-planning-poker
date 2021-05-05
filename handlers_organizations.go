package main

import (
	"net/http"
	"strconv"

	"github.com/StevenWeathers/thunderdome-planning-poker/pkg/database"
	"github.com/gorilla/mux"
)

// handleGetOrganizationsByUser gets a list of organizations the user is apart of
func (s *server) handleGetOrganizationsByUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		UserID := r.Context().Value(contextKeyUserID).(string)
		vars := mux.Vars(r)
		Limit, _ := strconv.Atoi(vars["limit"])
		Offset, _ := strconv.Atoi(vars["offset"])

		Organizations := s.database.OrganizationListByUser(UserID, Limit, Offset)

		s.respondWithJSON(w, http.StatusOK, Organizations)
	}
}

// handleGetOrganizationByUser gets an organization with user role
func (s *server) handleGetOrganizationByUser() http.HandlerFunc {
	type OrganizationResponse struct {
		Organization *database.Organization `json:"organization"`
		Role         string                 `json:"role"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		OrgRole := r.Context().Value(contextKeyOrgRole).(string)
		vars := mux.Vars(r)
		OrgID := vars["orgId"]

		Organization, err := s.database.OrganizationGet(OrgID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		s.respondWithJSON(w, http.StatusOK, &OrganizationResponse{
			Organization: Organization,
			Role:         OrgRole,
		})
	}
}

// handleCreateOrganization handles creating an organization with current user as admin
func (s *server) handleCreateOrganization() http.HandlerFunc {
	type CreateOrgResponse struct {
		OrganizationID string `json:"id"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		UserID := r.Context().Value(contextKeyUserID).(string)
		keyVal := s.getJSONRequestBody(r, w)

		OrgName := keyVal["name"].(string)
		OrgId, err := s.database.OrganizationCreate(UserID, OrgName)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var NewOrg = &CreateOrgResponse{
			OrganizationID: OrgId,
		}

		s.respondWithJSON(w, http.StatusOK, NewOrg)
	}
}

// handleGetOrganizationTeams gets a list of teams associated to the organization
func (s *server) handleGetOrganizationTeams() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		OrgID := vars["orgId"]
		Limit, _ := strconv.Atoi(vars["limit"])
		Offset, _ := strconv.Atoi(vars["offset"])

		Teams := s.database.OrganizationTeamList(OrgID, Limit, Offset)

		s.respondWithJSON(w, http.StatusOK, Teams)
	}
}

// handleGetOrganizationUsers gets a list of users associated to the organization
func (s *server) handleGetOrganizationUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		OrgID := vars["orgId"]
		Limit, _ := strconv.Atoi(vars["limit"])
		Offset, _ := strconv.Atoi(vars["offset"])

		Teams := s.database.OrganizationUserList(OrgID, Limit, Offset)

		s.respondWithJSON(w, http.StatusOK, Teams)
	}
}

// handleCreateOrganizationTeam handles creating an organization team
func (s *server) handleCreateOrganizationTeam() http.HandlerFunc {
	type CreateTeamResponse struct {
		TeamID string `json:"id"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		// UserID := r.Context().Value(contextKeyUserID).(string)
		vars := mux.Vars(r)
		keyVal := s.getJSONRequestBody(r, w)

		TeamName := keyVal["name"].(string)
		OrgID := vars["orgId"]
		TeamID, err := s.database.OrganizationTeamCreate(OrgID, TeamName)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var NewTeam = &CreateTeamResponse{
			TeamID: TeamID,
		}

		s.respondWithJSON(w, http.StatusOK, NewTeam)
	}
}

// handleOrganizationAddUser handles adding user to an organization
func (s *server) handleOrganizationAddUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		keyVal := s.getJSONRequestBody(r, w)

		vars := mux.Vars(r)
		OrgID := vars["orgId"]
		UserEmail := keyVal["email"].(string)
		Role := keyVal["role"].(string)

		User, UserErr := s.database.GetUserByEmail(UserEmail)
		if UserErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, err := s.database.OrganizationAddUser(OrgID, User.UserID, Role)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		return
	}
}

// handleGetOrganizationTeamByUser gets a team with users roles
func (s *server) handleGetOrganizationTeamByUser() http.HandlerFunc {
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

		Organization, err := s.database.OrganizationGet(OrgID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		Team, err := s.database.TeamGet(TeamID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		s.respondWithJSON(w, http.StatusOK, &TeamResponse{
			Organization:     Organization,
			Team:             Team,
			OrganizationRole: OrgRole,
			TeamRole:         TeamRole,
		})
	}
}
