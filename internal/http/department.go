package http

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"go.uber.org/zap"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"

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
// @Summary      Get Departments
// @Description  get list of organizations departments
// @Tags         organization
// @Produce      json
// @Param        orgId  path    string  true  "the organization ID to get departments for"
// @Success      200    object  standardJsonResponse{data=[]thunderdome.Department}
// @Security     ApiKeyAuth
// @Router       /organizations/{orgId}/departments [get]
func (s *Service) handleGetOrganizationDepartments() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !s.Config.OrganizationsEnabled {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, "ORGANIZATIONS_DISABLED"))
			return
		}
		ctx := r.Context()
		vars := mux.Vars(r)
		OrgID := vars["orgId"]
		Limit, Offset := getLimitOffsetFromRequest(r)

		Departments := s.OrganizationDataSvc.OrganizationDepartmentList(ctx, OrgID, Limit, Offset)

		s.Success(w, r, http.StatusOK, Departments, nil)
	}
}

// handleGetDepartmentByUser gets a department with user role
// @Summary      Get Department
// @Description  Gets an organization department with users role
// @Tags         organization
// @Produce      json
// @Param        orgId         path    string  true  "the organization ID"
// @Param        departmentId  path    string  true  "the department ID to get"
// @Success      200           object  standardJsonResponse{data=departmentResponse}
// @Failure      500           object  standardJsonResponse{}
// @Security     ApiKeyAuth
// @Router       /organizations/{orgId}/departments/{departmentId} [get]
func (s *Service) handleGetDepartmentByUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !s.Config.OrganizationsEnabled {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, "ORGANIZATIONS_DISABLED"))
			return
		}
		ctx := r.Context()
		OrgRole := ctx.Value(contextKeyOrgRole).(string)
		DepartmentRole := ctx.Value(contextKeyDepartmentRole).(string)
		SessionUserID := ctx.Value(contextKeyUserID).(string)
		vars := mux.Vars(r)
		OrgID := vars["orgId"]
		DepartmentID := vars["departmentId"]

		Organization, err := s.OrganizationDataSvc.OrganizationGet(ctx, OrgID)
		if err != nil {
			s.Logger.Ctx(ctx).Error(
				"handleGetDepartmentByUser error", zap.Error(err), zap.String("session_user_id", SessionUserID),
				zap.String("organization_id", OrgID), zap.String("department_id", DepartmentID),
				zap.String("organization_role", OrgRole), zap.String("department_role", DepartmentRole))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		Department, err := s.OrganizationDataSvc.DepartmentGet(ctx, DepartmentID)
		if err != nil {
			s.Logger.Ctx(ctx).Error(
				"handleGetDepartmentByUser error", zap.Error(err), zap.String("session_user_id", SessionUserID),
				zap.String("organization_id", OrgID), zap.String("department_id", DepartmentID),
				zap.String("organization_role", OrgRole), zap.String("department_role", DepartmentRole))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		result := &departmentResponse{
			Organization:     Organization,
			Department:       Department,
			OrganizationRole: OrgRole,
			DepartmentRole:   DepartmentRole,
		}

		s.Success(w, r, http.StatusOK, result, nil)
	}
}

// handleCreateDepartment handles creating an organization department
// @Summary      Create Department
// @Description  Create an organization department
// @Tags         organization
// @Produce      json
// @Param        orgId       path    string                 true  "the organization ID to create department for"
// @Param        department  body    teamCreateRequestBody  true  "new department object"
// @Success      200         object  standardJsonResponse{data=thunderdome.Department}
// @Failure      500         object  standardJsonResponse{}
// @Security     ApiKeyAuth
// @Router       /organizations/{orgId}/departments [post]
func (s *Service) handleCreateDepartment() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !s.Config.OrganizationsEnabled {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, "ORGANIZATIONS_DISABLED"))
			return
		}
		ctx := r.Context()
		SessionUserID := ctx.Value(contextKeyUserID).(string)
		vars := mux.Vars(r)
		OrgID := vars["orgId"]

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

		NewDepartment, err := s.OrganizationDataSvc.DepartmentCreate(r.Context(), OrgID, team.Name)
		if err != nil {
			s.Logger.Ctx(ctx).Error(
				"handleCreateDepartment error", zap.Error(err), zap.String("session_user_id", SessionUserID),
				zap.String("organization_id", OrgID), zap.String("department_name", team.Name))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, NewDepartment, nil)
	}
}

// handleGetDepartmentTeams gets a list of teams associated to the department
// @Summary      Get Department Teams
// @Description  Gets a list of organization department teams
// @Tags         organization
// @Produce      json
// @Param        orgId         path    string  true  "the organization ID"
// @Param        departmentId  path    string  true  "the department ID to get teams for"
// @Success      200           object  standardJsonResponse{data=[]thunderdome.Team}
// @Security     ApiKeyAuth
// @Router       /organizations/{orgId}/departments/{departmentId}/teams [get]
func (s *Service) handleGetDepartmentTeams() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !s.Config.OrganizationsEnabled {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, "ORGANIZATIONS_DISABLED"))
			return
		}
		ctx := r.Context()
		vars := mux.Vars(r)
		DepartmentID := vars["departmentId"]
		Limit, Offset := getLimitOffsetFromRequest(r)

		Teams := s.OrganizationDataSvc.DepartmentTeamList(ctx, DepartmentID, Limit, Offset)

		s.Success(w, r, http.StatusOK, Teams, nil)
	}
}

// handleGetDepartmentUsers gets a list of users associated to the department
// @Summary      Get Department Users
// @Description  Get a list of organization department users
// @Tags         organization
// @Produce      json
// @Param        orgId         path    string  true  "the organization ID"
// @Param        departmentId  path    string  true  "the department ID"
// @Success      200           object  standardJsonResponse{data=[]thunderdome.DepartmentUser}
// @Security     ApiKeyAuth
// @Router       /organizations/{orgId}/departments/{departmentId}/users [get]
func (s *Service) handleGetDepartmentUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !s.Config.OrganizationsEnabled {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, "ORGANIZATIONS_DISABLED"))
			return
		}
		ctx := r.Context()
		vars := mux.Vars(r)
		DepartmentID := vars["departmentId"]
		Limit, Offset := getLimitOffsetFromRequest(r)

		Users := s.OrganizationDataSvc.DepartmentUserList(ctx, DepartmentID, Limit, Offset)

		s.Success(w, r, http.StatusOK, Users, nil)
	}
}

// handleCreateDepartmentTeam handles creating an department team
// @Summary      Create Department Team
// @Description  Create a department team
// @Tags         organization
// @Produce      json
// @Param        orgId         path    string                 true  "the organization ID"
// @Param        departmentId  path    string                 true  "the department ID"
// @Param        team          body    teamCreateRequestBody  true  "new team object"
// @Success      200           object  standardJsonResponse{data=thunderdome.Team}
// @Failure      500           object  standardJsonResponse{}
// @Security     ApiKeyAuth
// @Router       /organizations/{orgId}/departments/{departmentId}/teams [post]
func (s *Service) handleCreateDepartmentTeam() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !s.Config.OrganizationsEnabled {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, "ORGANIZATIONS_DISABLED"))
			return
		}
		ctx := r.Context()
		SessionUserID := ctx.Value(contextKeyUserID).(string)
		vars := mux.Vars(r)
		OrgID := vars["orgId"]
		DepartmentID := vars["departmentId"]

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

		NewTeam, err := s.OrganizationDataSvc.DepartmentTeamCreate(ctx, DepartmentID, team.Name)
		if err != nil {
			s.Logger.Ctx(ctx).Error(
				"handleCreateDepartmentTeam error", zap.Error(err), zap.String("session_user_id", SessionUserID),
				zap.String("organization_id", OrgID), zap.String("department_id", DepartmentID),
				zap.String("team_name", team.Name))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, NewTeam, nil)
	}
}

// handleDepartmentAddUser handles adding user to an organization department
// @Summary      Add Department User
// @Description  Add a department User
// @Tags         organization
// @Produce      json
// @Param        orgId         path    string                  true  "the organization ID"
// @Param        departmentId  path    string                  true  "the department ID"
// @Param        user          body    teamAddUserRequestBody  true  "new department user object"
// @Success      200           object  standardJsonResponse{}
// @Failure      500           object  standardJsonResponse{}
// @Security     ApiKeyAuth
// @Router       /organizations/{orgId}/departments/{departmentId}/users [post]
func (s *Service) handleDepartmentAddUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !s.Config.OrganizationsEnabled {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, "ORGANIZATIONS_DISABLED"))
			return
		}
		ctx := r.Context()
		SessionUserID := ctx.Value(contextKeyUserID).(string)
		vars := mux.Vars(r)
		OrgID := vars["orgId"]
		DepartmentId := vars["departmentId"]

		var u = teamAddUserRequestBody{}
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

		UserEmail := strings.ToLower(u.Email)

		User, UserErr := s.UserDataSvc.GetUserByEmail(ctx, UserEmail)
		if UserErr != nil {
			s.Logger.Ctx(ctx).Error(
				"handleDepartmentAddUser error", zap.Error(UserErr), zap.String("session_user_id", SessionUserID),
				zap.String("organization_id", OrgID), zap.String("department_id", DepartmentId),
				zap.String("user_email", UserEmail))
			s.Failure(w, r, http.StatusInternalServerError, Errorf(ENOTFOUND, "USER_NOT_FOUND"))
			return
		}

		_, err := s.OrganizationDataSvc.DepartmentAddUser(ctx, DepartmentId, User.Id, u.Role)
		if err != nil {
			s.Logger.Ctx(ctx).Error(
				"handleDepartmentAddUser error", zap.Error(err), zap.String("session_user_id", SessionUserID),
				zap.String("organization_id", OrgID), zap.String("department_id", DepartmentId),
				zap.String("user_id", User.Id), zap.String("user_role", u.Role))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, nil, nil)
	}
}

// handleDepartmentUpdateUser handles updating an organization department user
// @Summary      Update Department User
// @Description  Update a department User
// @Tags         organization
// @Produce      json
// @Param        orgId         path    string                  true  "the organization ID"
// @Param        departmentId  path    string                  true  "the department ID"
// @Param        userId  path    string                  true  "the user ID"
// @Param        user          body    teamUpdateUserRequestBody  true  "department user object"
// @Success      200           object  standardJsonResponse{}
// @Failure      500           object  standardJsonResponse{}
// @Security     ApiKeyAuth
// @Router       /organizations/{orgId}/departments/{departmentId}/users/{userId} [put]
func (s *Service) handleDepartmentUpdateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !s.Config.OrganizationsEnabled {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, "ORGANIZATIONS_DISABLED"))
			return
		}
		ctx := r.Context()
		SessionUserID := ctx.Value(contextKeyUserID).(string)
		vars := mux.Vars(r)
		OrgID := vars["orgId"]
		DepartmentId := vars["departmentId"]
		UserId := vars["userId"]

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

		_, err := s.OrganizationDataSvc.DepartmentUpdateUser(ctx, DepartmentId, UserId, u.Role)
		if err != nil {
			s.Logger.Ctx(ctx).Error(
				"handleDepartmentUpdateUser error", zap.Error(err), zap.String("session_user_id", SessionUserID),
				zap.String("organization_id", OrgID), zap.String("department_id", DepartmentId),
				zap.String("user_id", UserId), zap.String("user_role", u.Role))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, nil, nil)
	}
}

// handleDepartmentRemoveUser handles removing user from a department (and department teams)
// @Summary      Remove Department User
// @Description  Remove a department User
// @Tags         organization
// @Produce      json
// @Param        orgId         path    string  true  "the organization ID"
// @Param        departmentId  path    string  true  "the department ID"
// @Param        userId        path    string  true  "the user ID"
// @Success      200           object  standardJsonResponse{}
// @Failure      500           object  standardJsonResponse{}
// @Security     ApiKeyAuth
// @Router       /organizations/{orgId}/departments/{departmentId}/users/{userId} [delete]
func (s *Service) handleDepartmentRemoveUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !s.Config.OrganizationsEnabled {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, "ORGANIZATIONS_DISABLED"))
			return
		}
		ctx := r.Context()
		SessionUserID := ctx.Value(contextKeyUserID).(string)
		vars := mux.Vars(r)
		OrgID := vars["orgId"]
		DepartmentID := vars["departmentId"]
		UserID := vars["userId"]
		idErr := validate.Var(UserID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		err := s.OrganizationDataSvc.DepartmentRemoveUser(ctx, DepartmentID, UserID)
		if err != nil {
			s.Logger.Ctx(ctx).Error(
				"handleDepartmentRemoveUser error", zap.Error(err), zap.String("session_user_id", SessionUserID),
				zap.String("organization_id", OrgID), zap.String("department_id", DepartmentID),
				zap.String("user_id", UserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, nil, nil)
	}
}

// handleDepartmentTeamAddUser handles adding user to a team so long as they are in the department
// @Summary      Add Department Team User
// @Description  Add a User to Department Team
// @Tags         organization
// @Produce      json
// @Param        orgId         path    string                  true  "the organization ID"
// @Param        departmentId  path    string                  true  "the department ID"
// @Param        teamId        path    string                  true  "the team ID"
// @Param        user          body    teamAddUserRequestBody  true  "new team user object"
// @Success      200           object  standardJsonResponse{}
// @Failure      500           object  standardJsonResponse{}
// @Security     ApiKeyAuth
// @Router       /organizations/{orgId}/departments/{departmentId}/teams/{teamId}/users [post]
func (s *Service) handleDepartmentTeamAddUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !s.Config.OrganizationsEnabled {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, "ORGANIZATIONS_DISABLED"))
			return
		}
		ctx := r.Context()
		SessionUserID := ctx.Value(contextKeyUserID).(string)
		vars := mux.Vars(r)
		OrgID := vars["orgId"]
		DepartmentID := vars["departmentId"]
		TeamID := vars["teamId"]

		var u = teamAddUserRequestBody{}
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

		UserEmail := u.Email

		User, UserErr := s.UserDataSvc.GetUserByEmail(ctx, UserEmail)
		if UserErr != nil {
			s.Logger.Ctx(ctx).Error(
				"handleDepartmentTeamAddUser error", zap.Error(UserErr), zap.String("session_user_id", SessionUserID),
				zap.String("organization_id", OrgID), zap.String("department_id", DepartmentID),
				zap.String("team_id", TeamID), zap.String("user_email", UserEmail))
			s.Failure(w, r, http.StatusInternalServerError, Errorf(ENOTFOUND, "USER_NOT_FOUND"))
			return
		}

		_, DepartmentRole, roleErr := s.OrganizationDataSvc.DepartmentUserRole(ctx, User.Id, OrgID, DepartmentID)
		if DepartmentRole == "" || roleErr != nil {
			s.Logger.Ctx(ctx).Error(
				"handleDepartmentTeamAddUser error", zap.Error(roleErr), zap.String("session_user_id", SessionUserID),
				zap.String("organization_id", OrgID), zap.String("department_id", DepartmentID),
				zap.String("team_id", TeamID), zap.String("user_id", User.Id))
			s.Failure(w, r, http.StatusInternalServerError, Errorf(EUNAUTHORIZED, "DEPARTMENT_USER_REQUIRED"))
			return
		}

		_, err := s.TeamDataSvc.TeamAddUser(ctx, TeamID, User.Id, u.Role)
		if err != nil {
			s.Logger.Ctx(ctx).Error(
				"handleDepartmentTeamAddUser error", zap.Error(err), zap.String("session_user_id", SessionUserID),
				zap.String("organization_id", OrgID), zap.String("department_id", DepartmentID),
				zap.String("team_id", TeamID), zap.String("user_id", User.Id))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, nil, nil)
	}
}

// handleDepartmentTeamByUser gets a team with users roles
// @Summary      Get Department Team
// @Description  Get a department team with users role
// @Tags         organization
// @Produce      json
// @Param        orgId         path    string  true  "the organization ID"
// @Param        departmentId  path    string  true  "the department ID"
// @Param        teamId        path    string  true  "the team ID"
// @Success      200           object  standardJsonResponse{data=departmentTeamResponse}
// @Failure      500           object  standardJsonResponse{}
// @Security     ApiKeyAuth
// @Router       /organizations/{orgId}/departments/{departmentId}/teams/{teamId} [get]
func (s *Service) handleDepartmentTeamByUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !s.Config.OrganizationsEnabled {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, "ORGANIZATIONS_DISABLED"))
			return
		}
		ctx := r.Context()
		SessionUserID := ctx.Value(contextKeyUserID).(string)
		OrgRole := ctx.Value(contextKeyOrgRole).(string)
		DepartmentRole := ctx.Value(contextKeyDepartmentRole).(string)
		TeamRole := ctx.Value(contextKeyTeamRole).(string)
		vars := mux.Vars(r)
		OrgID := vars["orgId"]
		DepartmentID := vars["departmentId"]
		TeamID := vars["teamId"]

		Organization, err := s.OrganizationDataSvc.OrganizationGet(ctx, OrgID)
		if err != nil {
			s.Logger.Ctx(ctx).Error(
				"handleDepartmentTeamByUser error", zap.Error(err), zap.String("session_user_id", SessionUserID),
				zap.String("organization_id", OrgID), zap.String("department_id", DepartmentID),
				zap.String("team_id", TeamID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		Department, err := s.OrganizationDataSvc.DepartmentGet(ctx, DepartmentID)
		if err != nil {
			s.Logger.Ctx(ctx).Error(
				"handleDepartmentTeamByUser error", zap.Error(err), zap.String("session_user_id", SessionUserID),
				zap.String("organization_id", OrgID), zap.String("department_id", DepartmentID),
				zap.String("team_id", TeamID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		Team, err := s.TeamDataSvc.TeamGet(ctx, TeamID)
		if err != nil {
			s.Logger.Ctx(ctx).Error(
				"handleDepartmentTeamByUser error", zap.Error(err), zap.String("session_user_id", SessionUserID),
				zap.String("organization_id", OrgID), zap.String("department_id", DepartmentID),
				zap.String("team_id", TeamID))
			s.Failure(w, r, http.StatusInternalServerError, err)
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

		s.Success(w, r, http.StatusOK, result, nil)
	}
}

// handleDeleteDepartment handles deleting a department
// @Summary      Delete Department
// @Description  Delete a Department
// @Tags         organization
// @Produce      json
// @Param        departmentId  path    string  true  "the department ID"
// @Success      200           object  standardJsonResponse{}
// @Success      403           object  standardJsonResponse{}
// @Success      500           object  standardJsonResponse{}
// @Security     ApiKeyAuth
// @Router       /organizations/{orgId}/departments/{departmentId} [delete]
func (s *Service) handleDeleteDepartment() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		SessionUserID := ctx.Value(contextKeyUserID).(string)
		vars := mux.Vars(r)
		OrgID := vars["orgId"]
		DepartmentID := vars["departmentId"]
		idErr := validate.Var(DepartmentID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		err := s.OrganizationDataSvc.DepartmentDelete(ctx, DepartmentID)
		if err != nil {
			s.Logger.Ctx(ctx).Error(
				"handleDeleteDepartment error", zap.Error(err), zap.String("session_user_id", SessionUserID),
				zap.String("organization_id", OrgID), zap.String("department_id", DepartmentID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, nil, nil)
	}
}
