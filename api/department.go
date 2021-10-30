package api

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/StevenWeathers/thunderdome-planning-poker/pkg/database"
	"github.com/gorilla/mux"
)

// handleGetOrganizationDepartments gets a list of departments associated to the organization
func (a *api) handleGetOrganizationDepartments() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		OrgID := vars["orgId"]
		Limit, _ := strconv.Atoi(vars["limit"])
		Offset, _ := strconv.Atoi(vars["offset"])

		Teams := a.db.OrganizationDepartmentList(OrgID, Limit, Offset)

		a.respondWithJSON(w, http.StatusOK, Teams)
	}
}

// handleGetDepartmentByUser gets an department with user role
func (a *api) handleGetDepartmentByUser() http.HandlerFunc {
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

		Organization, err := a.db.OrganizationGet(OrgID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		Department, err := a.db.DepartmentGet(DepartmentID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		a.respondWithJSON(w, http.StatusOK, &DepartmentResponse{
			Organization:     Organization,
			Department:       Department,
			OrganizationRole: OrgRole,
			DepartmentRole:   DepartmentRole,
		})
	}
}

// handleCreateDepartment handles creating an organization department
func (a *api) handleCreateDepartment() http.HandlerFunc {
	type CreateDepartmentResponse struct {
		DepartmentID string `json:"id"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		// UserID := r.Context().Value(contextKeyUserID).(string)
		vars := mux.Vars(r)
		keyVal := a.getJSONRequestBody(r, w)

		OrgName := keyVal["name"].(string)
		OrgID := vars["orgId"]
		DepartmentID, err := a.db.DepartmentCreate(OrgID, OrgName)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var NewDepartment = &CreateDepartmentResponse{
			DepartmentID: DepartmentID,
		}

		a.respondWithJSON(w, http.StatusOK, NewDepartment)
	}
}

// handleGetDepartmentTeams gets a list of teams associated to the department
func (a *api) handleGetDepartmentTeams() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		DepartmentID := vars["departmentId"]
		Limit, _ := strconv.Atoi(vars["limit"])
		Offset, _ := strconv.Atoi(vars["offset"])

		Teams := a.db.DepartmentTeamList(DepartmentID, Limit, Offset)

		a.respondWithJSON(w, http.StatusOK, Teams)
	}
}

// handleGetDepartmentUsers gets a list of users associated to the department
func (a *api) handleGetDepartmentUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		DepartmentID := vars["departmentId"]
		Limit, _ := strconv.Atoi(vars["limit"])
		Offset, _ := strconv.Atoi(vars["offset"])

		Teams := a.db.DepartmentUserList(DepartmentID, Limit, Offset)

		a.respondWithJSON(w, http.StatusOK, Teams)
	}
}

// handleCreateDepartmentTeam handles creating an department team
func (a *api) handleCreateDepartmentTeam() http.HandlerFunc {
	type CreateTeamResponse struct {
		TeamID string `json:"id"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		// UserID := r.Context().Value(contextKeyUserID).(string)
		vars := mux.Vars(r)
		keyVal := a.getJSONRequestBody(r, w)

		TeamName := keyVal["name"].(string)
		DepartmentID := vars["departmentId"]
		TeamID, err := a.db.DepartmentTeamCreate(DepartmentID, TeamName)
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

// handleDepartmentAddUser handles adding user to an organization department
func (a *api) handleDepartmentAddUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		keyVal := a.getJSONRequestBody(r, w)

		vars := mux.Vars(r)
		DepartmentId := vars["departmentId"]
		UserEmail := strings.ToLower(keyVal["email"].(string))
		Role := keyVal["role"].(string)

		User, UserErr := a.db.GetUserByEmail(UserEmail)
		if UserErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, err := a.db.DepartmentAddUser(DepartmentId, User.UserID, Role)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		return
	}
}

// handleDepartmentRemoveUser handles removing user from a department (and department teams)
func (a *api) handleDepartmentRemoveUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		keyVal := a.getJSONRequestBody(r, w)

		vars := mux.Vars(r)
		DepartmentID := vars["departmentId"]
		UserID := keyVal["id"].(string)

		err := a.db.DepartmentRemoveUser(DepartmentID, UserID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		return
	}
}

// handleDepartmentTeamAddUser handles adding user to a team so long as they are in the department
func (a *api) handleDepartmentTeamAddUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		keyVal := a.getJSONRequestBody(r, w)

		vars := mux.Vars(r)
		OrgID := vars["orgId"]
		DepartmentID := vars["departmentId"]
		TeamID := vars["teamId"]
		UserEmail := strings.ToLower(keyVal["email"].(string))
		Role := keyVal["role"].(string)

		User, UserErr := a.db.GetUserByEmail(UserEmail)
		if UserErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, DepartmentRole, roleErr := a.db.DepartmentUserRole(User.UserID, OrgID, DepartmentID)
		if DepartmentRole == "" || roleErr != nil {
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

// handleDepartmentTeamByUser gets a team with users roles
func (a *api) handleDepartmentTeamByUser() http.HandlerFunc {
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

		Organization, err := a.db.OrganizationGet(OrgID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		Department, err := a.db.DepartmentGet(DepartmentID)
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
			Department:       Department,
			Team:             Team,
			OrganizationRole: OrgRole,
			DepartmentRole:   DepartmentRole,
			TeamRole:         TeamRole,
		})
	}
}
