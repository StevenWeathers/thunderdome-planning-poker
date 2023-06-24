package api

import (
	"encoding/json"
	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

type departmentResponse struct {
	Organization     *thunderdome.Organization `json:"organization"`
	Department       *thunderdome.Department   `json:"department"`
	OrganizationRole string                    `json:"organizationRole"`
	DepartmentRole   string                    `json:"departmentRole"`
}

type departmentTeamResponse struct {
	Organization     *thunderdome.Organization `json:"organization"`
	Department       *thunderdome.Department   `json:"department"`
	Team             *thunderdome.Team         `json:"team"`
	OrganizationRole string                    `json:"organizationRole"`
	DepartmentRole   string                    `json:"departmentRole"`
	TeamRole         string                    `json:"teamRole"`
}

// handleGetOrganizationDepartments gets a list of departments associated to the organization
// @Summary Get Departments
// @Description get list of organizations departments
// @Tags organization
// @Produce  json
// @Param orgId path string true "the organization ID to get departments for"
// @Success 200 object standardJsonResponse{data=[]thunderdome.Department}
// @Security ApiKeyAuth
// @Router /organizations/{orgId}/departments [get]
func (a *APIService) handleGetOrganizationDepartments() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !a.Config.OrganizationsEnabled {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, "ORGANIZATIONS_DISABLED"))
			return
		}
		vars := mux.Vars(r)
		OrgID := vars["orgId"]
		Limit, Offset := getLimitOffsetFromRequest(r)

		Departments := a.DB.OrganizationDepartmentList(r.Context(), OrgID, Limit, Offset)

		a.Success(w, r, http.StatusOK, Departments, nil)
	}
}

// handleGetDepartmentByUser gets a department with user role
// @Summary Get Department
// @Description Gets an organization department with users role
// @Tags organization
// @Produce  json
// @Param orgId path string true "the organization ID"
// @Param departmentId path string true "the department ID to get"
// @Success 200 object standardJsonResponse{data=departmentResponse}
// @Failure 500 object standardJsonResponse{}
// @Security ApiKeyAuth
// @Router /organizations/{orgId}/departments/{departmentId} [get]
func (a *APIService) handleGetDepartmentByUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !a.Config.OrganizationsEnabled {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, "ORGANIZATIONS_DISABLED"))
			return
		}
		ctx := r.Context()
		OrgRole := ctx.Value(contextKeyOrgRole).(string)
		DepartmentRole := ctx.Value(contextKeyDepartmentRole).(string)
		vars := mux.Vars(r)
		OrgID := vars["orgId"]
		DepartmentID := vars["departmentId"]

		Organization, err := a.DB.OrganizationGet(ctx, OrgID)
		if err != nil {
			a.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		Department, err := a.DB.DepartmentGet(ctx, DepartmentID)
		if err != nil {
			a.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		result := &departmentResponse{
			Organization:     Organization,
			Department:       Department,
			OrganizationRole: OrgRole,
			DepartmentRole:   DepartmentRole,
		}

		a.Success(w, r, http.StatusOK, result, nil)
	}
}

// handleCreateDepartment handles creating an organization department
// @Summary Create Department
// @Description Create an organization department
// @Tags organization
// @Produce  json
// @Param orgId path string true "the organization ID to create department for"
// @Param department body teamCreateRequestBody true "new department object"
// @Success 200 object standardJsonResponse{data=thunderdome.Department}
// @Failure 500 object standardJsonResponse{}
// @Security ApiKeyAuth
// @Router /organizations/{orgId}/departments [post]
func (a *APIService) handleCreateDepartment() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !a.Config.OrganizationsEnabled {
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

		NewDepartment, err := a.DB.DepartmentCreate(r.Context(), OrgID, team.Name)
		if err != nil {
			a.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		a.Success(w, r, http.StatusOK, NewDepartment, nil)
	}
}

// handleGetDepartmentTeams gets a list of teams associated to the department
// @Summary Get Department Teams
// @Description Gets a list of organization department teams
// @Tags organization
// @Produce  json
// @Param orgId path string true "the organization ID"
// @Param departmentId path string true "the department ID to get teams for"
// @Success 200 object standardJsonResponse{data=[]thunderdome.Team}
// @Security ApiKeyAuth
// @Router /organizations/{orgId}/departments/{departmentId}/teams [get]
func (a *APIService) handleGetDepartmentTeams() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !a.Config.OrganizationsEnabled {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, "ORGANIZATIONS_DISABLED"))
			return
		}
		vars := mux.Vars(r)
		DepartmentID := vars["departmentId"]
		Limit, Offset := getLimitOffsetFromRequest(r)

		Teams := a.DB.DepartmentTeamList(r.Context(), DepartmentID, Limit, Offset)

		a.Success(w, r, http.StatusOK, Teams, nil)
	}
}

// handleGetDepartmentUsers gets a list of users associated to the department
// @Summary Get Department Users
// @Description Get a list of organization department users
// @Tags organization
// @Produce  json
// @Param orgId path string true "the organization ID"
// @Param departmentId path string true "the department ID"
// @Success 200 object standardJsonResponse{data=[]thunderdome.DepartmentUser}
// @Security ApiKeyAuth
// @Router /organizations/{orgId}/departments/{departmentId}/users [get]
func (a *APIService) handleGetDepartmentUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !a.Config.OrganizationsEnabled {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, "ORGANIZATIONS_DISABLED"))
			return
		}
		vars := mux.Vars(r)
		DepartmentID := vars["departmentId"]
		Limit, Offset := getLimitOffsetFromRequest(r)

		Users := a.DB.DepartmentUserList(r.Context(), DepartmentID, Limit, Offset)

		a.Success(w, r, http.StatusOK, Users, nil)
	}
}

// handleCreateDepartmentTeam handles creating an department team
// @Summary Create Department Team
// @Description Create a department team
// @Tags organization
// @Produce  json
// @Param orgId path string true "the organization ID"
// @Param departmentId path string true "the department ID"
// @Param team body teamCreateRequestBody true "new team object"
// @Success 200 object standardJsonResponse{data=thunderdome.Team}
// @Failure 500 object standardJsonResponse{}
// @Security ApiKeyAuth
// @Router /organizations/{orgId}/departments/{departmentId}/teams [post]
func (a *APIService) handleCreateDepartmentTeam() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !a.Config.OrganizationsEnabled {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, "ORGANIZATIONS_DISABLED"))
			return
		}
		vars := mux.Vars(r)
		DepartmentID := vars["departmentId"]

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

		NewTeam, err := a.DB.DepartmentTeamCreate(r.Context(), DepartmentID, team.Name)
		if err != nil {
			a.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		a.Success(w, r, http.StatusOK, NewTeam, nil)
	}
}

// handleDepartmentAddUser handles adding user to an organization department
// @Summary Add Department User
// @Description Add a department User
// @Tags organization
// @Produce  json
// @Param orgId path string true "the organization ID"
// @Param departmentId path string true "the department ID"
// @Param user body teamAddUserRequestBody true "new department user object"
// @Success 200 object standardJsonResponse{}
// @Failure 500 object standardJsonResponse{}
// @Security ApiKeyAuth
// @Router /organizations/{orgId}/departments/{departmentId}/users [post]
func (a *APIService) handleDepartmentAddUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !a.Config.OrganizationsEnabled {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, "ORGANIZATIONS_DISABLED"))
			return
		}

		vars := mux.Vars(r)
		DepartmentId := vars["departmentId"]

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

		User, UserErr := a.UserService.GetUserByEmail(r.Context(), UserEmail)
		if UserErr != nil {
			a.Failure(w, r, http.StatusInternalServerError, Errorf(ENOTFOUND, "USER_NOT_FOUND"))
			return
		}

		_, err := a.DB.DepartmentAddUser(r.Context(), DepartmentId, User.Id, u.Role)
		if err != nil {
			a.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		a.Success(w, r, http.StatusOK, nil, nil)
	}
}

// handleDepartmentRemoveUser handles removing user from a department (and department teams)
// @Summary Remove Department User
// @Description Remove a department User
// @Tags organization
// @Produce  json
// @Param orgId path string true "the organization ID"
// @Param departmentId path string true "the department ID"
// @Param userId path string true "the user ID"
// @Success 200 object standardJsonResponse{}
// @Failure 500 object standardJsonResponse{}
// @Security ApiKeyAuth
// @Router /organizations/{orgId}/departments/{departmentId}/users/{userId} [delete]
func (a *APIService) handleDepartmentRemoveUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !a.Config.OrganizationsEnabled {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, "ORGANIZATIONS_DISABLED"))
			return
		}
		vars := mux.Vars(r)
		DepartmentID := vars["departmentId"]
		UserID := vars["userId"]
		idErr := validate.Var(UserID, "required,uuid")
		if idErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		err := a.DB.DepartmentRemoveUser(r.Context(), DepartmentID, UserID)
		if err != nil {
			a.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		a.Success(w, r, http.StatusOK, nil, nil)
	}
}

// handleDepartmentTeamAddUser handles adding user to a team so long as they are in the department
// @Summary Add Department Team User
// @Description Add a User to Department Team
// @Tags organization
// @Produce  json
// @Param orgId path string true "the organization ID"
// @Param departmentId path string true "the department ID"
// @Param teamId path string true "the team ID"
// @Param user body teamAddUserRequestBody true "new team user object"
// @Success 200 object standardJsonResponse{}
// @Failure 500 object standardJsonResponse{}
// @Security ApiKeyAuth
// @Router /organizations/{orgId}/departments/{departmentId}/teams/{teamId}/users [post]
func (a *APIService) handleDepartmentTeamAddUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !a.Config.OrganizationsEnabled {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, "ORGANIZATIONS_DISABLED"))
			return
		}

		vars := mux.Vars(r)
		OrgID := vars["orgId"]
		DepartmentID := vars["departmentId"]
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

		User, UserErr := a.UserService.GetUserByEmail(r.Context(), UserEmail)
		if UserErr != nil {
			a.Failure(w, r, http.StatusInternalServerError, Errorf(ENOTFOUND, "USER_NOT_FOUND"))
			return
		}

		_, DepartmentRole, roleErr := a.DB.DepartmentUserRole(r.Context(), User.Id, OrgID, DepartmentID)
		if DepartmentRole == "" || roleErr != nil {
			a.Failure(w, r, http.StatusInternalServerError, Errorf(EUNAUTHORIZED, "DEPARTMENT_USER_REQUIRED"))
			return
		}

		_, err := a.DB.TeamAddUser(r.Context(), TeamID, User.Id, u.Role)
		if err != nil {
			a.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		a.Success(w, r, http.StatusOK, nil, nil)
	}
}

// handleDepartmentTeamByUser gets a team with users roles
// @Summary Get Department Team
// @Description Get a department team with users role
// @Tags organization
// @Produce  json
// @Param orgId path string true "the organization ID"
// @Param departmentId path string true "the department ID"
// @Param teamId path string true "the team ID"
// @Success 200 object standardJsonResponse{data=departmentTeamResponse}
// @Failure 500 object standardJsonResponse{}
// @Security ApiKeyAuth
// @Router /organizations/{orgId}/departments/{departmentId}/teams/{teamId} [get]
func (a *APIService) handleDepartmentTeamByUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !a.Config.OrganizationsEnabled {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, "ORGANIZATIONS_DISABLED"))
			return
		}
		ctx := r.Context()
		OrgRole := ctx.Value(contextKeyOrgRole).(string)
		DepartmentRole := ctx.Value(contextKeyDepartmentRole).(string)
		TeamRole := ctx.Value(contextKeyTeamRole).(string)
		vars := mux.Vars(r)
		OrgID := vars["orgId"]
		DepartmentID := vars["departmentId"]
		TeamID := vars["teamId"]

		Organization, err := a.DB.OrganizationGet(r.Context(), OrgID)
		if err != nil {
			a.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		Department, err := a.DB.DepartmentGet(r.Context(), DepartmentID)
		if err != nil {
			a.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		Team, err := a.DB.TeamGet(ctx, TeamID)
		if err != nil {
			a.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		result := &departmentTeamResponse{
			Organization:     Organization,
			Department:       Department,
			Team:             Team,
			OrganizationRole: OrgRole,
			DepartmentRole:   DepartmentRole,
			TeamRole:         TeamRole,
		}

		a.Success(w, r, http.StatusOK, result, nil)
	}
}

// handleDeleteDepartment handles deleting a department
// @Summary Delete Department
// @Description Delete a Department
// @Tags organization
// @Produce  json
// @Param departmentId path string true "the department ID"
// @Success 200 object standardJsonResponse{}
// @Success 403 object standardJsonResponse{}
// @Success 500 object standardJsonResponse{}
// @Security ApiKeyAuth
// @Router /organizations/{orgId}/departments/{departmentId} [delete]
func (a *APIService) handleDeleteDepartment() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		DepartmentID := vars["departmentId"]
		idErr := validate.Var(DepartmentID, "required,uuid")
		if idErr != nil {
			a.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		err := a.DB.DepartmentDelete(r.Context(), DepartmentID)
		if err != nil {
			a.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		a.Success(w, r, http.StatusOK, nil, nil)
	}
}
