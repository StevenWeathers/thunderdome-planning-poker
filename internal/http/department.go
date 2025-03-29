package http

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"go.uber.org/zap"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
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
//
//	@Summary		Get Departments
//	@Description	get list of organizations departments
//	@Tags			organization
//	@Produce		json
//	@Param			orgId	path	string	true	"the organization ID to get departments for"
//	@Success		200		object	standardJsonResponse{data=[]thunderdome.Department}
//	@Security		ApiKeyAuth
//	@Router			/organizations/{orgId}/departments [get]
func (s *Service) handleGetOrganizationDepartments() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !s.Config.OrganizationsEnabled {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, "ORGANIZATIONS_DISABLED"))
			return
		}
		ctx := r.Context()
		orgID := r.PathValue("orgId")
		oidErr := validate.Var(orgID, "required,uuid")
		if oidErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, oidErr.Error()))
			return
		}
		limit, offset := getLimitOffsetFromRequest(r)

		departments := s.OrganizationDataSvc.OrganizationDepartmentList(ctx, orgID, limit, offset)

		s.Success(w, r, http.StatusOK, departments, nil)
	}
}

// handleGetDepartmentByUser gets a department with user role
//
//	@Summary		Get Department
//	@Description	Gets an organization department with users role
//	@Tags			organization
//	@Produce		json
//	@Param			orgId			path	string	true	"the organization ID"
//	@Param			departmentId	path	string	true	"the department ID to get"
//	@Success		200				object	standardJsonResponse{data=departmentResponse}
//	@Failure		500				object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/organizations/{orgId}/departments/{departmentId} [get]
func (s *Service) handleGetDepartmentByUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !s.Config.OrganizationsEnabled {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, "ORGANIZATIONS_DISABLED"))
			return
		}
		ctx := r.Context()
		orgRole := ctx.Value(contextKeyOrgRole).(string)
		departmentRole := ctx.Value(contextKeyDepartmentRole).(string)
		sessionUserID := ctx.Value(contextKeyUserID).(string)
		orgID := r.PathValue("orgId")
		oidErr := validate.Var(orgID, "required,uuid")
		if oidErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, oidErr.Error()))
			return
		}
		departmentID := r.PathValue("departmentId")
		didErr := validate.Var(departmentID, "required,uuid")
		if didErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, didErr.Error()))
			return
		}

		organization, err := s.OrganizationDataSvc.OrganizationGetByID(ctx, orgID)
		if err != nil {
			s.Logger.Ctx(ctx).Error(
				"handleGetDepartmentByUser error", zap.Error(err), zap.String("session_user_id", sessionUserID),
				zap.String("organization_id", orgID), zap.String("department_id", departmentID),
				zap.String("organization_role", orgRole), zap.String("department_role", departmentRole))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		department, err := s.OrganizationDataSvc.DepartmentGetByID(ctx, departmentID)
		if err != nil {
			s.Logger.Ctx(ctx).Error(
				"handleGetDepartmentByUser error", zap.Error(err), zap.String("session_user_id", sessionUserID),
				zap.String("organization_id", orgID), zap.String("department_id", departmentID),
				zap.String("organization_role", orgRole), zap.String("department_role", departmentRole))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		result := &departmentResponse{
			Organization:     organization,
			Department:       department,
			OrganizationRole: orgRole,
			DepartmentRole:   departmentRole,
		}

		s.Success(w, r, http.StatusOK, result, nil)
	}
}

// handleCreateDepartment handles creating an organization department
//
//	@Summary		Create Department
//	@Description	Create an organization department
//	@Tags			organization
//	@Produce		json
//	@Param			orgId		path	string					true	"the organization ID to create department for"
//	@Param			department	body	teamCreateRequestBody	true	"new department object"
//	@Success		200			object	standardJsonResponse{data=thunderdome.Department}
//	@Failure		500			object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/organizations/{orgId}/departments [post]
func (s *Service) handleCreateDepartment() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !s.Config.OrganizationsEnabled {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, "ORGANIZATIONS_DISABLED"))
			return
		}
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)
		orgID := r.PathValue("orgId")
		oidErr := validate.Var(orgID, "required,uuid")
		if oidErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, oidErr.Error()))
			return
		}

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

		newDepartment, err := s.OrganizationDataSvc.DepartmentCreate(r.Context(), orgID, team.Name)
		if err != nil {
			s.Logger.Ctx(ctx).Error(
				"handleCreateDepartment error", zap.Error(err), zap.String("session_user_id", sessionUserID),
				zap.String("organization_id", orgID), zap.String("department_name", team.Name))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, newDepartment, nil)
	}
}

// handleDepartmentUpdate handles updating an organization department
//
//	@Summary		Update Department
//	@Description	Update an organization department
//	@Tags			organization
//	@Produce		json
//	@Param			orgId		path	string					true	"the organization ID"
//	@Param			deptId		path	string					true	"the department ID"
//	@Param			department	body	teamCreateRequestBody	true	"updated department object"
//	@Success		200			object	standardJsonResponse{data=thunderdome.Department}
//	@Failure		500			object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/organizations/{orgId}/departments/{departmentId} [put]
func (s *Service) handleDepartmentUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !s.Config.OrganizationsEnabled {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, "ORGANIZATIONS_DISABLED"))
			return
		}
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)
		orgID := r.PathValue("orgId")
		oidErr := validate.Var(orgID, "required,uuid")
		if oidErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, oidErr.Error()))
			return
		}
		deptID := r.PathValue("departmentId")
		didErr := validate.Var(deptID, "required,uuid")
		if didErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, didErr.Error()))
			return
		}

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

		newDepartment, err := s.OrganizationDataSvc.DepartmentUpdate(r.Context(), deptID, team.Name)
		if err != nil {
			s.Logger.Ctx(ctx).Error(
				"handleDepartmentUpdate error", zap.Error(err),
				zap.String("session_user_id", sessionUserID),
				zap.String("organization_id", orgID),
				zap.String("department_id", deptID),
				zap.String("department_name", team.Name))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, newDepartment, nil)
	}
}

// handleGetDepartmentTeams gets a list of teams associated to the department
//
//	@Summary		Get Department Teams
//	@Description	Gets a list of organization department teams
//	@Tags			organization
//	@Produce		json
//	@Param			orgId			path	string	true	"the organization ID"
//	@Param			departmentId	path	string	true	"the department ID to get teams for"
//	@Success		200				object	standardJsonResponse{data=[]thunderdome.Team}
//	@Security		ApiKeyAuth
//	@Router			/organizations/{orgId}/departments/{departmentId}/teams [get]
func (s *Service) handleGetDepartmentTeams() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !s.Config.OrganizationsEnabled {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, "ORGANIZATIONS_DISABLED"))
			return
		}
		ctx := r.Context()
		departmentID := r.PathValue("departmentId")
		didErr := validate.Var(departmentID, "required,uuid")
		if didErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, didErr.Error()))
			return
		}
		limit, offset := getLimitOffsetFromRequest(r)

		teams := s.OrganizationDataSvc.DepartmentTeamList(ctx, departmentID, limit, offset)

		s.Success(w, r, http.StatusOK, teams, nil)
	}
}

// handleGetDepartmentUsers gets a list of users associated to the department
//
//	@Summary		Get Department Users
//	@Description	Get a list of organization department users
//	@Tags			organization
//	@Produce		json
//	@Param			orgId			path	string	true	"the organization ID"
//	@Param			departmentId	path	string	true	"the department ID"
//	@Success		200				object	standardJsonResponse{data=[]thunderdome.DepartmentUser}
//	@Security		ApiKeyAuth
//	@Router			/organizations/{orgId}/departments/{departmentId}/users [get]
func (s *Service) handleGetDepartmentUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !s.Config.OrganizationsEnabled {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, "ORGANIZATIONS_DISABLED"))
			return
		}
		ctx := r.Context()
		departmentID := r.PathValue("departmentId")
		didErr := validate.Var(departmentID, "required,uuid")
		if didErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, didErr.Error()))
			return
		}
		limit, offset := getLimitOffsetFromRequest(r)

		users := s.OrganizationDataSvc.DepartmentUserList(ctx, departmentID, limit, offset)

		s.Success(w, r, http.StatusOK, users, nil)
	}
}

// handleCreateDepartmentTeam handles creating an department team
//
//	@Summary		Create Department Team
//	@Description	Create a department team
//	@Tags			organization
//	@Produce		json
//	@Param			orgId			path	string					true	"the organization ID"
//	@Param			departmentId	path	string					true	"the department ID"
//	@Param			team			body	teamCreateRequestBody	true	"new team object"
//	@Success		200				object	standardJsonResponse{data=thunderdome.Team}
//	@Failure		500				object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/organizations/{orgId}/departments/{departmentId}/teams [post]
func (s *Service) handleCreateDepartmentTeam() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !s.Config.OrganizationsEnabled {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, "ORGANIZATIONS_DISABLED"))
			return
		}
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)
		orgID := r.PathValue("orgId")
		oidErr := validate.Var(orgID, "required,uuid")
		if oidErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, oidErr.Error()))
			return
		}
		departmentID := r.PathValue("departmentId")
		didErr := validate.Var(departmentID, "required,uuid")
		if didErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, didErr.Error()))
			return
		}

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

		newTeam, err := s.OrganizationDataSvc.DepartmentTeamCreate(ctx, departmentID, team.Name)
		if err != nil {
			s.Logger.Ctx(ctx).Error(
				"handleCreateDepartmentTeam error", zap.Error(err), zap.String("session_user_id", sessionUserID),
				zap.String("organization_id", orgID), zap.String("department_id", departmentID),
				zap.String("team_name", team.Name))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, newTeam, nil)
	}
}

// handleDepartmentAddUser handles adding user to an organization department
//
//	@Summary		Add Department User
//	@Description	Add a department User
//	@Tags			organization
//	@Produce		json
//	@Param			orgId			path	string				true	"the organization ID"
//	@Param			departmentId	path	string				true	"the department ID"
//	@Param			user			body	addUserRequestBody	true	"new department user object"
//	@Success		200				object	standardJsonResponse{}
//	@Failure		500				object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/organizations/{orgId}/departments/{departmentId}/users [post]
func (s *Service) handleDepartmentAddUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !s.Config.OrganizationsEnabled {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, "ORGANIZATIONS_DISABLED"))
			return
		}
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)
		orgID := r.PathValue("orgId")
		oidErr := validate.Var(orgID, "required,uuid")
		if oidErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, oidErr.Error()))
			return
		}
		departmentID := r.PathValue("departmentId")
		didErr := validate.Var(departmentID, "required,uuid")
		if didErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, didErr.Error()))
			return
		}

		var u = addUserRequestBody{}
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

		_, err := s.OrganizationDataSvc.DepartmentAddUser(ctx, departmentID, u.UserID, u.Role)
		if err != nil {
			s.Logger.Ctx(ctx).Error(
				"handleDepartmentAddUser error", zap.Error(err), zap.String("session_user_id", sessionUserID),
				zap.String("organization_id", orgID), zap.String("department_id", departmentID),
				zap.String("user_id", u.UserID), zap.String("user_role", u.Role))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, nil, nil)
	}
}

// handleDepartmentUpdateUser handles updating an organization department user
//
//	@Summary		Update Department User
//	@Description	Update a department User
//	@Tags			organization
//	@Produce		json
//	@Param			orgId			path	string						true	"the organization ID"
//	@Param			departmentId	path	string						true	"the department ID"
//	@Param			userId			path	string						true	"the user ID"
//	@Param			user			body	teamUpdateUserRequestBody	true	"department user object"
//	@Success		200				object	standardJsonResponse{}
//	@Failure		500				object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/organizations/{orgId}/departments/{departmentId}/users/{userId} [put]
func (s *Service) handleDepartmentUpdateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !s.Config.OrganizationsEnabled {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, "ORGANIZATIONS_DISABLED"))
			return
		}
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)
		orgID := r.PathValue("orgId")
		oidErr := validate.Var(orgID, "required,uuid")
		if oidErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, oidErr.Error()))
			return
		}
		departmentID := r.PathValue("departmentId")
		didErr := validate.Var(departmentID, "required,uuid")
		if didErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, didErr.Error()))
			return
		}
		userID := r.PathValue("userId")
		uidErr := validate.Var(userID, "required,uuid")
		if uidErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, uidErr.Error()))
			return
		}

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

		_, err := s.OrganizationDataSvc.DepartmentUpdateUser(ctx, departmentID, userID, u.Role)
		if err != nil {
			s.Logger.Ctx(ctx).Error(
				"handleDepartmentUpdateUser error", zap.Error(err), zap.String("session_user_id", sessionUserID),
				zap.String("organization_id", orgID), zap.String("department_id", departmentID),
				zap.String("user_id", userID), zap.String("user_role", u.Role))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, nil, nil)
	}
}

// handleDepartmentRemoveUser handles removing user from a department (and department teams)
//
//	@Summary		Remove Department User
//	@Description	Remove a department User
//	@Tags			organization
//	@Produce		json
//	@Param			orgId			path	string	true	"the organization ID"
//	@Param			departmentId	path	string	true	"the department ID"
//	@Param			userId			path	string	true	"the user ID"
//	@Success		200				object	standardJsonResponse{}
//	@Failure		500				object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/organizations/{orgId}/departments/{departmentId}/users/{userId} [delete]
func (s *Service) handleDepartmentRemoveUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !s.Config.OrganizationsEnabled {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, "ORGANIZATIONS_DISABLED"))
			return
		}
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)
		orgID := r.PathValue("orgId")
		oidErr := validate.Var(orgID, "required,uuid")
		if oidErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, oidErr.Error()))
			return
		}
		departmentID := r.PathValue("departmentId")
		didErr := validate.Var(departmentID, "required,uuid")
		if didErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, didErr.Error()))
			return
		}
		userID := r.PathValue("userId")
		idErr := validate.Var(userID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		err := s.OrganizationDataSvc.DepartmentRemoveUser(ctx, departmentID, userID)
		if err != nil {
			s.Logger.Ctx(ctx).Error(
				"handleDepartmentRemoveUser error", zap.Error(err), zap.String("session_user_id", sessionUserID),
				zap.String("organization_id", orgID), zap.String("department_id", departmentID),
				zap.String("user_id", userID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, nil, nil)
	}
}

// handleDepartmentTeamAddUser handles adding user to a team so long as they are in the department
//
//	@Summary		Add Department Team User
//	@Description	Add a User to Department Team
//	@Tags			organization
//	@Produce		json
//	@Param			orgId			path	string				true	"the organization ID"
//	@Param			departmentId	path	string				true	"the department ID"
//	@Param			teamId			path	string				true	"the team ID"
//	@Param			user			body	addUserRequestBody	true	"new team user object"
//	@Success		200				object	standardJsonResponse{}
//	@Failure		500				object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/organizations/{orgId}/departments/{departmentId}/teams/{teamId}/users [post]
func (s *Service) handleDepartmentTeamAddUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !s.Config.OrganizationsEnabled {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, "ORGANIZATIONS_DISABLED"))
			return
		}
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)

		orgID := r.PathValue("orgId")
		oidErr := validate.Var(orgID, "required,uuid")
		if oidErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, oidErr.Error()))
			return
		}
		departmentID := r.PathValue("departmentId")
		didErr := validate.Var(departmentID, "required,uuid")
		if didErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, didErr.Error()))
			return
		}
		teamID := r.PathValue("teamId")
		tidErr := validate.Var(teamID, "required,uuid")
		if tidErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, tidErr.Error()))
			return
		}

		var u = addUserRequestBody{}
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

		_, departmentRole, roleErr := s.OrganizationDataSvc.DepartmentUserRole(ctx, u.UserID, orgID, departmentID)
		if departmentRole == "" || roleErr != nil {
			s.Logger.Ctx(ctx).Error(
				"handleDepartmentTeamAddUser error", zap.Error(roleErr), zap.String("session_user_id", sessionUserID),
				zap.String("organization_id", orgID), zap.String("department_id", departmentID),
				zap.String("team_id", teamID), zap.String("user_id", u.UserID))
			s.Failure(w, r, http.StatusInternalServerError, Errorf(EUNAUTHORIZED, "DEPARTMENT_USER_REQUIRED"))
			return
		}

		_, err := s.TeamDataSvc.TeamAddUser(ctx, teamID, u.UserID, u.Role)
		if err != nil {
			s.Logger.Ctx(ctx).Error(
				"handleDepartmentTeamAddUser error", zap.Error(err), zap.String("session_user_id", sessionUserID),
				zap.String("organization_id", orgID), zap.String("department_id", departmentID),
				zap.String("team_id", teamID), zap.String("user_id", u.UserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, nil, nil)
	}
}

// handleDepartmentTeamByUser gets a team with users roles
//
//	@Summary		Get Department Team
//	@Description	Get a department team with users role
//	@Tags			organization
//	@Produce		json
//	@Param			orgId			path	string	true	"the organization ID"
//	@Param			departmentId	path	string	true	"the department ID"
//	@Param			teamId			path	string	true	"the team ID"
//	@Success		200				object	standardJsonResponse{data=departmentTeamResponse}
//	@Failure		500				object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/organizations/{orgId}/departments/{departmentId}/teams/{teamId} [get]
func (s *Service) handleDepartmentTeamByUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !s.Config.OrganizationsEnabled {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, "ORGANIZATIONS_DISABLED"))
			return
		}
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)
		teamUserRoles := ctx.Value(contextKeyUserTeamRoles).(*thunderdome.UserTeamRoleInfo)
		var emptyRole = ""
		orgRole := teamUserRoles.OrganizationRole
		if orgRole == nil {
			orgRole = &emptyRole
		}
		departmentRole := teamUserRoles.DepartmentRole
		if departmentRole == nil {
			departmentRole = &emptyRole
		}
		teamRole := teamUserRoles.TeamRole
		if teamRole == nil {
			teamRole = &emptyRole
		}

		orgID := r.PathValue("orgId")
		oidErr := validate.Var(orgID, "required,uuid")
		if oidErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, oidErr.Error()))
			return
		}
		departmentID := r.PathValue("departmentId")
		didErr := validate.Var(departmentID, "required,uuid")
		if didErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, didErr.Error()))
			return
		}
		teamID := r.PathValue("teamId")
		tidErr := validate.Var(teamID, "required,uuid")
		if tidErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, tidErr.Error()))
			return
		}

		organization, err := s.OrganizationDataSvc.OrganizationGetByID(ctx, orgID)
		if err != nil {
			s.Logger.Ctx(ctx).Error(
				"handleDepartmentTeamByUser error", zap.Error(err), zap.String("session_user_id", sessionUserID),
				zap.String("organization_id", orgID), zap.String("department_id", departmentID),
				zap.String("team_id", teamID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		department, err := s.OrganizationDataSvc.DepartmentGetByID(ctx, departmentID)
		if err != nil {
			s.Logger.Ctx(ctx).Error(
				"handleDepartmentTeamByUser error", zap.Error(err), zap.String("session_user_id", sessionUserID),
				zap.String("organization_id", orgID), zap.String("department_id", departmentID),
				zap.String("team_id", teamID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		team, err := s.TeamDataSvc.TeamGetByID(ctx, teamID)
		if err != nil {
			s.Logger.Ctx(ctx).Error(
				"handleDepartmentTeamByUser error", zap.Error(err), zap.String("session_user_id", sessionUserID),
				zap.String("organization_id", orgID), zap.String("department_id", departmentID),
				zap.String("team_id", teamID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		result := &departmentTeamResponse{
			Organization:     organization,
			Department:       department,
			Team:             team,
			OrganizationRole: *orgRole,
			DepartmentRole:   *departmentRole,
			TeamRole:         *teamRole,
		}

		s.Success(w, r, http.StatusOK, result, nil)
	}
}

// handleDeleteDepartment handles deleting a department
//
//	@Summary		Delete Department
//	@Description	Delete a Department
//	@Tags			organization
//	@Produce		json
//	@Param			departmentId	path	string	true	"the department ID"
//	@Success		200				object	standardJsonResponse{}
//	@Success		403				object	standardJsonResponse{}
//	@Success		500				object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/organizations/{orgId}/departments/{departmentId} [delete]
func (s *Service) handleDeleteDepartment() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)

		orgID := r.PathValue("orgId")
		oidErr := validate.Var(orgID, "required,uuid")
		if oidErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, oidErr.Error()))
			return
		}
		departmentID := r.PathValue("departmentId")
		idErr := validate.Var(departmentID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		err := s.OrganizationDataSvc.DepartmentDelete(ctx, departmentID)
		if err != nil {
			s.Logger.Ctx(ctx).Error(
				"handleDeleteDepartment error", zap.Error(err), zap.String("session_user_id", sessionUserID),
				zap.String("organization_id", orgID), zap.String("department_id", departmentID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, nil, nil)
	}
}

// handleGetDepartmentUserInvites gets a list of user invites associated to the department
//
//	@Summary		Get Department User Invites
//	@Description	Get a list of user invites associated to the department
//	@Tags			organization
//	@Produce		json
//	@Param			organizationId	path	string	true	"the org ID"
//	@Param			departmentId	path	string	true	"the dept ID"
//	@Success		200				object	standardJsonResponse{data=[]thunderdome.DepartmentUserInvite}
//	@Security		ApiKeyAuth
//	@Router			/organizations/{orgId}/departments/{departmentId}/invites [get]
func (s *Service) handleGetDepartmentUserInvites() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)

		deptID := r.PathValue("departmentId")
		didErr := validate.Var(deptID, "required,uuid")
		if didErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, didErr.Error()))
			return
		}

		invites, err := s.OrganizationDataSvc.DepartmentGetUserInvites(ctx, deptID)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleGetDepartmentUserInvites error", zap.Error(err), zap.String("department_id", deptID),
				zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, invites, nil)
	}
}

// handleDeleteDepartmentUserInvite handles deleting user invite from an department
//
//	@Summary		Delete Department User Invite
//	@Description	Delete user invite from department
//	@Tags			organization
//	@Produce		json
//	@Param			orgId			path	string	true	"organization id"
//	@Param			departmentId	path	string	true	"the dept ID"
//	@Param			inviteId		path	string	true	"invite id"
//	@Success		200				object	standardJsonResponse{}
//	@Failure		403				object	standardJsonResponse{}
//	@Failure		500				object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/organizations/{orgId}/departments/{departmentId}/invites/{inviteId} [delete]
func (s *Service) handleDeleteDepartmentUserInvite() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !s.Config.OrganizationsEnabled {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, "ORGANIZATIONS_DISABLED"))
			return
		}
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)

		orgID := r.PathValue("orgId")
		oidErr := validate.Var(orgID, "required,uuid")
		if oidErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, oidErr.Error()))
			return
		}
		deptID := r.PathValue("departmentId")
		didErr := validate.Var(deptID, "required,uuid")
		if didErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, didErr.Error()))
			return
		}
		inviteID := r.PathValue("inviteId")
		idErr := validate.Var(inviteID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		err := s.OrganizationDataSvc.DepartmentDeleteUserInvite(ctx, inviteID)
		if err != nil {
			s.Logger.Ctx(ctx).Error(
				"handleDeleteDepartmentUserInvite error", zap.Error(err),
				zap.String("invite_id", inviteID),
				zap.String("session_user_id", sessionUserID),
				zap.String("organization_id", orgID),
				zap.String("department_id", deptID),
			)
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, nil, nil)
	}
}

// handleDepartmentInviteUser handles inviting user to an organization department
//
//	@Summary		Invite Department User
//	@Description	Invite a department User
//	@Tags			organization
//	@Produce		json
//	@Param			orgId			path	string						true	"the organization ID"
//	@Param			departmentId	path	string						true	"the department ID"
//	@Param			user			body	teamInviteUserRequestBody	true	"new department user object"
//	@Success		200				object	standardJsonResponse{}
//	@Failure		500				object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/organizations/{orgId}/departments/{departmentId}/invites [post]
func (s *Service) handleDepartmentInviteUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !s.Config.OrganizationsEnabled {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, "ORGANIZATIONS_DISABLED"))
			return
		}
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)

		orgID := r.PathValue("orgId")
		oidErr := validate.Var(orgID, "required,uuid")
		if oidErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, oidErr.Error()))
			return
		}
		departmentID := r.PathValue("departmentId")
		didErr := validate.Var(departmentID, "required,uuid")
		if didErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, didErr.Error()))
			return
		}

		var u = teamInviteUserRequestBody{}
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

		userEmail := strings.ToLower(u.Email)

		inviteID, inviteErr := s.OrganizationDataSvc.DepartmentInviteUser(ctx, departmentID, userEmail, u.Role)
		if inviteErr != nil {
			s.Logger.Ctx(ctx).Error("handleDepartmentInviteUser error", zap.Error(inviteErr),
				zap.String("organization_id", orgID),
				zap.String("department_id", departmentID),
				zap.String("session_user_id", sessionUserID),
				zap.String("user_email", userEmail),
			)
			s.Failure(w, r, http.StatusInternalServerError, inviteErr)
			return
		}

		org, orgErr := s.OrganizationDataSvc.OrganizationGetByID(ctx, orgID)
		if orgErr != nil {
			s.Logger.Ctx(ctx).Error("handleDepartmentInviteUser error", zap.Error(orgErr),
				zap.String("organization_id", orgID),
				zap.String("department_id", departmentID),
				zap.String("session_user_id", sessionUserID),
				zap.String("user_email", userEmail),
			)
			s.Failure(w, r, http.StatusInternalServerError, orgErr)
			return
		}
		dept, deptErr := s.OrganizationDataSvc.DepartmentGetByID(ctx, departmentID)
		if deptErr != nil {
			s.Logger.Ctx(ctx).Error("handleDepartmentInviteUser error", zap.Error(orgErr),
				zap.String("organization_id", orgID),
				zap.String("department_id", departmentID),
				zap.String("session_user_id", sessionUserID),
				zap.String("user_email", userEmail),
			)
			s.Failure(w, r, http.StatusInternalServerError, orgErr)
			return
		}
		emailErr := s.Email.SendDepartmentInvite(org.Name, dept.Name, userEmail, inviteID)
		if emailErr != nil {
			s.Logger.Ctx(ctx).Error("handleDepartmentInviteUser error", zap.Error(emailErr),
				zap.String("organization_id", orgID),
				zap.String("department_id", departmentID),
				zap.String("session_user_id", sessionUserID),
				zap.String("user_email", userEmail),
			)
			s.Failure(w, r, http.StatusInternalServerError, emailErr)
			return
		}
		s.Success(w, r, http.StatusOK, nil, nil)
	}
}
