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
		UserID := r.Context().Value(contextKeyUserID).(string)
		vars := mux.Vars(r)

		Organization, role, err := s.database.OrganizationWithRole(UserID, vars["orgId"])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		s.respondWithJSON(w, http.StatusOK, &OrganizationResponse{
			Organization: Organization,
			Role:         role,
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
			w.WriteHeader(http.StatusBadRequest)
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
