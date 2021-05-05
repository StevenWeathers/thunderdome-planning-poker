package main

import (
	"net/http"
	"strconv"

	"github.com/StevenWeathers/thunderdome-planning-poker/pkg/database"
	"github.com/gorilla/mux"
)

// handleGetOrganizationDepartments gets a list of departments associated to the organization
func (s *server) handleGetOrganizationDepartments() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		OrgID := vars["orgId"]
		Limit, _ := strconv.Atoi(vars["limit"])
		Offset, _ := strconv.Atoi(vars["offset"])

		Teams := s.database.OrganizationDepartmentList(OrgID, Limit, Offset)

		s.respondWithJSON(w, http.StatusOK, Teams)
	}
}

// handleGetDepartmentByUser gets an department with user role
func (s *server) handleGetDepartmentByUser() http.HandlerFunc {
	type DepartmentResponse struct {
		Organization     *database.Organization `json:"organization"`
		Department       *database.Department   `json:"department"`
		OrganizationRole string                 `json:"organizationRole"`
		DepartmentRole   string                 `json:"departmentRole"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		OrgRole := r.Context().Value(contextKeyOrgRole).(string)
		DepartmentRole := r.Context().Value(contextKeyDepartmentRole).(string)
		vars := mux.Vars(r)
		OrgID := vars["orgId"]
		DepartmentID := vars["departmentId"]

		Organization, err := s.database.OrganizationGet(OrgID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		Department, err := s.database.DepartmentGet(DepartmentID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		s.respondWithJSON(w, http.StatusOK, &DepartmentResponse{
			Organization:     Organization,
			Department:       Department,
			OrganizationRole: OrgRole,
			DepartmentRole:   DepartmentRole,
		})
	}
}

// handleCreateDepartment handles creating an organization department
func (s *server) handleCreateDepartment() http.HandlerFunc {
	type CreateDepartmentResponse struct {
		DepartmentID string `json:"id"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		// UserID := r.Context().Value(contextKeyUserID).(string)
		vars := mux.Vars(r)
		keyVal := s.getJSONRequestBody(r, w)

		OrgName := keyVal["name"].(string)
		OrgID := vars["orgId"]
		DepartmentID, err := s.database.DepartmentCreate(OrgID, OrgName)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var NewDepartment = &CreateDepartmentResponse{
			DepartmentID: DepartmentID,
		}

		s.respondWithJSON(w, http.StatusOK, NewDepartment)
	}
}

// handleGetDepartmentTeams gets a list of teams associated to the department
func (s *server) handleGetDepartmentTeams() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		DepartmentID := vars["departmentId"]
		Limit, _ := strconv.Atoi(vars["limit"])
		Offset, _ := strconv.Atoi(vars["offset"])

		Teams := s.database.DepartmentTeamList(DepartmentID, Limit, Offset)

		s.respondWithJSON(w, http.StatusOK, Teams)
	}
}

// handleGetDepartmentUsers gets a list of users associated to the department
func (s *server) handleGetDepartmentUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		DepartmentID := vars["departmentId"]
		Limit, _ := strconv.Atoi(vars["limit"])
		Offset, _ := strconv.Atoi(vars["offset"])

		Teams := s.database.DepartmentUserList(DepartmentID, Limit, Offset)

		s.respondWithJSON(w, http.StatusOK, Teams)
	}
}

// handleCreateDepartmentTeam handles creating an department team
func (s *server) handleCreateDepartmentTeam() http.HandlerFunc {
	type CreateTeamResponse struct {
		TeamID string `json:"id"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		// UserID := r.Context().Value(contextKeyUserID).(string)
		vars := mux.Vars(r)
		keyVal := s.getJSONRequestBody(r, w)

		TeamName := keyVal["name"].(string)
		DepartmentID := vars["departmentId"]
		TeamID, err := s.database.DepartmentTeamCreate(DepartmentID, TeamName)
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

// handleDepartmentAddUser handles adding user to an organization department
func (s *server) handleDepartmentAddUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		keyVal := s.getJSONRequestBody(r, w)

		vars := mux.Vars(r)
		DepartmentId := vars["departmentId"]
		UserEmail := keyVal["email"].(string)
		Role := keyVal["role"].(string)

		User, UserErr := s.database.GetUserByEmail(UserEmail)
		if UserErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, err := s.database.DepartmentAddUser(DepartmentId, User.UserID, Role)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		return
	}
}

// handleDepartmentTeamByUser gets a team with users roles
func (s *server) handleDepartmentTeamByUser() http.HandlerFunc {
	type TeamResponse struct {
		Organization     *database.Organization `json:"organization"`
		Department       *database.Department   `json:"department"`
		Team             *database.Team         `json:"team"`
		OrganizationRole string                 `json:"organizationRole"`
		DepartmentRole   string                 `json:"departmentRole"`
		TeamRole         string                 `json:"teamRole"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		OrgRole := r.Context().Value(contextKeyOrgRole).(string)
		DepartmentRole := r.Context().Value(contextKeyDepartmentRole).(string)
		TeamRole := r.Context().Value(contextKeyTeamRole).(string)
		vars := mux.Vars(r)
		OrgID := vars["orgId"]
		DepartmentID := vars["departmentId"]
		TeamID := vars["teamId"]

		Organization, err := s.database.OrganizationGet(OrgID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		Department, err := s.database.DepartmentGet(DepartmentID)
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
			Department:       Department,
			Team:             Team,
			OrganizationRole: OrgRole,
			DepartmentRole:   DepartmentRole,
			TeamRole:         TeamRole,
		})
	}
}
