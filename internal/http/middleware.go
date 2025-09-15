package http

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"go.uber.org/zap"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
)

func (s *Service) panicRecovery(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				s.Logger.Error(fmt.Sprintf("http handler recovering from panic error: %v", err))
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()

		h.ServeHTTP(w, r)
	})
}

// userOnly validates that the request was made by a valid user
func (s *Service) userOnly(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey := strings.TrimSpace(r.Header.Get(apiKeyHeaderName))
		ctx := r.Context()
		var user *thunderdome.User

		if apiKey != "" && s.Config.ExternalAPIEnabled {
			var apiKeyErr error
			user, apiKeyErr = s.ApiKeyDataSvc.GetAPIKeyUser(ctx, apiKey)
			if apiKeyErr != nil {
				s.Failure(w, r, http.StatusUnauthorized, Errorf(EINVALID, "INVALID_APIKEY"))
				return
			}
		} else {
			sessionID, cookieErr := s.Cookie.ValidateSessionCookie(w, r)
			if cookieErr != nil && cookieErr.Error() != "COOKIE_NOT_FOUND" {
				s.Failure(w, r, http.StatusUnauthorized, Errorf(EINVALID, "INVALID_USER"))
				return
			}

			if sessionID != "" {
				var userErr error
				user, userErr = s.AuthDataSvc.GetSessionUserByID(ctx, sessionID)
				if userErr != nil {
					s.Failure(w, r, http.StatusUnauthorized, Errorf(EINVALID, "INVALID_USER"))
					return
				}
			} else {
				userID, err := s.Cookie.ValidateUserCookie(w, r)
				if err != nil {
					s.Failure(w, r, http.StatusUnauthorized, Errorf(EINVALID, "INVALID_USER"))
					return
				}

				var userErr error
				user, userErr = s.UserDataSvc.GetGuestUserByID(ctx, userID)
				if userErr != nil {
					s.Failure(w, r, http.StatusUnauthorized, Errorf(EINVALID, "INVALID_USER"))
					return
				}
			}
		}

		ctx = context.WithValue(ctx, contextKeyUserID, user.ID)
		ctx = context.WithValue(ctx, contextKeyUserType, user.Type)

		h(w, r.WithContext(ctx))
	}
}

// entityUserOnly validates that the request was made by the session user matching the {userId} of the entity (or ADMIN)
func (s *Service) entityUserOnly(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		userID := ctx.Value(contextKeyUserID).(string)
		userType := ctx.Value(contextKeyUserType).(string)
		entityUserID := r.PathValue("userId")
		idErr := validate.Var(entityUserID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		if userType != thunderdome.AdminUserType && entityUserID != userID {
			s.Failure(w, r, http.StatusForbidden, Errorf(EINVALID, "INVALID_USER"))
			return
		}

		h(w, r)
	}
}

// registeredUserOnly validates that the request was made by a registered user
func (s *Service) registeredUserOnly(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userType := r.Context().Value(contextKeyUserType).(string)

		if userType == thunderdome.GuestUserType {
			s.Failure(w, r, http.StatusForbidden, Errorf(EINVALID, "REGISTERED_USER_ONLY"))
			return
		}

		h(w, r)
	}
}

// adminOnly middleware checks if the user is an admin, otherwise reject their request
func (s *Service) adminOnly(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userType := r.Context().Value(contextKeyUserType).(string)

		if userType != thunderdome.AdminUserType {
			s.Failure(w, r, http.StatusForbidden, Errorf(EUNAUTHORIZED, "REQUIRES_ADMIN"))
			return
		}

		h(w, r)
	}
}

// verifiedUserOnly validates that the request was made by a verified registered user
func (s *Service) verifiedUserOnly(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)
		userType := ctx.Value(contextKeyUserType).(string)
		entityUserID := r.PathValue("userId")
		idErr := validate.Var(entityUserID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		if userType != thunderdome.AdminUserType && (entityUserID != sessionUserID) {
			s.Failure(w, r, http.StatusForbidden, Errorf(EINVALID, "INVALID_USER"))
			return
		}

		entityUser, entityUserErr := s.UserDataSvc.GetUserByID(ctx, entityUserID)
		if entityUserErr != nil {
			s.Logger.Ctx(ctx).Error(
				"verifiedUserOnly error", zap.Error(entityUserErr), zap.String("entity_user_id", entityUserID),
				zap.String("session_user_id", sessionUserID), zap.String("session_user_type", userType))
			s.Failure(w, r, http.StatusInternalServerError, entityUserErr)
			return
		}

		if s.Config.ExternalAPIVerifyRequired && !entityUser.Verified {
			s.Failure(w, r, http.StatusForbidden, Errorf(EUNAUTHORIZED, "REQUIRES_VERIFIED_USER"))
			return
		}

		h(w, r)
	}
}

// subscribedEntityUserOnly validates that the request was made by the subscribed entity user
func (s *Service) subscribedEntityUserOnly(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		userID := ctx.Value(contextKeyUserID).(string)
		userType := ctx.Value(contextKeyUserType).(string)
		entityUserID := r.PathValue("userId")
		idErr := validate.Var(entityUserID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		if userType != thunderdome.AdminUserType && (entityUserID != userID) {
			s.Failure(w, r, http.StatusForbidden, Errorf(EINVALID, "INVALID_USER"))
			return
		}

		if !s.Config.SubscriptionsEnabled {
			h(w, r)
			return
		}

		// admins can bypass active subscriber functions
		if userType != thunderdome.AdminUserType {
			subscriberErr := s.SubscriptionDataSvc.CheckActiveSubscriber(ctx, entityUserID)
			if subscriberErr != nil {
				s.Failure(w, r, http.StatusForbidden, Errorf(EUNAUTHORIZED, "REQUIRES_SUBSCRIBED_USER"))
				return
			}
		}

		h(w, r)
	}
}

// subscribedUserOnly validates that the request was made by a subscribed user
func (s *Service) subscribedUserOnly(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID := ctx.Value(contextKeyUserID).(string)
		userType := ctx.Value(contextKeyUserType).(string)

		if !s.Config.SubscriptionsEnabled {
			h(w, r)
			return
		}

		// admins can bypass active subscriber functions
		if userType != thunderdome.AdminUserType {
			subscriberErr := s.SubscriptionDataSvc.CheckActiveSubscriber(ctx, userID)
			if subscriberErr != nil {
				s.Failure(w, r, http.StatusForbidden, Errorf(EUNAUTHORIZED, "REQUIRES_SUBSCRIBED_USER"))
				return
			}
		}

		h(w, r)
	}
}

// orgUserOnly validates that the request was made by a valid user of the organization
func (s *Service) orgUserOnly(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		userID := ctx.Value(contextKeyUserID).(string)
		userType := ctx.Value(contextKeyUserType).(string)
		orgID := r.PathValue("orgId")
		idErr := validate.Var(orgID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		_, err := s.OrganizationDataSvc.OrganizationGetByID(ctx, orgID)
		if err != nil && err.Error() == "ORGANIZATION_NOT_FOUND" {
			s.Failure(w, r, http.StatusNotFound, Errorf(ENOTFOUND, "ORGANIZATION_NOT_FOUND"))
			return
		} else if err != nil {
			s.Failure(w, r, http.StatusInternalServerError, Errorf(EINTERNAL, err.Error()))
			return
		}

		var role string
		if userType != thunderdome.AdminUserType {
			var userErr error
			role, userErr = s.OrganizationDataSvc.OrganizationUserRole(ctx, userID, orgID)
			if userErr != nil {
				s.Failure(w, r, http.StatusForbidden, Errorf(EUNAUTHORIZED, "ORGANIZATION_USER_REQUIRED"))
				return
			}
		} else {
			role = thunderdome.AdminUserType
		}

		ctx = context.WithValue(ctx, contextKeyOrgRole, role)

		h(w, r.WithContext(ctx))
	}
}

// orgAdminOnly validates that the request was made by an ADMIN of the organization
func (s *Service) orgAdminOnly(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		userID := ctx.Value(contextKeyUserID).(string)
		userType := ctx.Value(contextKeyUserType).(string)
		orgID := r.PathValue("orgId")
		idErr := validate.Var(orgID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		var role string
		if userType != thunderdome.AdminUserType {
			var userErr error
			role, userErr = s.OrganizationDataSvc.OrganizationUserRole(ctx, userID, orgID)
			if userErr != nil {
				s.Failure(w, r, http.StatusForbidden, Errorf(EUNAUTHORIZED, "ORGANIZATION_USER_REQUIRED"))
				return
			}
			if role != thunderdome.AdminUserType {
				s.Failure(w, r, http.StatusForbidden, Errorf(EUNAUTHORIZED, "REQUIRES_ORG_ADMIN"))
				return
			}
		} else {
			role = thunderdome.AdminUserType
		}

		ctx = context.WithValue(ctx, contextKeyOrgRole, role)

		h(w, r.WithContext(ctx))
	}
}

// departmentUserOnly validates that the request was made by a valid user of the organization (with department role)
func (s *Service) departmentUserOnly(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		userID := ctx.Value(contextKeyUserID).(string)
		userType := ctx.Value(contextKeyUserType).(string)
		orgID := r.PathValue("orgId")
		idErr := validate.Var(orgID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		departmentID := r.PathValue("departmentId")
		idErr = validate.Var(departmentID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		var orgRole string
		var departmentRole string
		if userType != thunderdome.AdminUserType {
			var userErr error
			orgRole, departmentRole, userErr = s.OrganizationDataSvc.DepartmentUserRole(ctx, userID, orgID, departmentID)
			if userErr != nil || (departmentRole == "" && orgRole != thunderdome.AdminUserType) {
				s.Logger.Ctx(ctx).Warn("middleware departmentUserOnly REQUIRES_DEPARTMENT_USER",
					zap.Error(userErr),
					zap.String("user_id", userID),
					zap.String("org_id", orgID),
					zap.String("department_id", departmentID))
				s.Failure(w, r, http.StatusForbidden, Errorf(EUNAUTHORIZED, "REQUIRES_DEPARTMENT_USER"))
				return
			}
		} else {
			orgRole = thunderdome.AdminUserType
			departmentRole = thunderdome.AdminUserType
		}

		ctx = context.WithValue(ctx, contextKeyOrgRole, orgRole)
		ctx = context.WithValue(ctx, contextKeyDepartmentRole, departmentRole)

		h(w, r.WithContext(ctx))
	}
}

// departmentAdminOnly validates that the request was made by an ADMIN of the organization (with department role)
func (s *Service) departmentAdminOnly(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		userID := ctx.Value(contextKeyUserID).(string)
		userType := ctx.Value(contextKeyUserType).(string)
		orgID := r.PathValue("orgId")
		idErr := validate.Var(orgID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		departmentID := r.PathValue("departmentId")
		idErr = validate.Var(departmentID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		var orgRole string
		var departmentRole string
		if userType != thunderdome.AdminUserType {
			var userErr error
			orgRole, departmentRole, userErr = s.OrganizationDataSvc.DepartmentUserRole(ctx, userID, orgID, departmentID)
			if userErr != nil {
				s.Failure(w, r, http.StatusForbidden, Errorf(EUNAUTHORIZED, "REQUIRES_DEPARTMENT_USER"))
				return
			}
			if departmentRole != thunderdome.AdminUserType && orgRole != thunderdome.AdminUserType {
				s.Failure(w, r, http.StatusForbidden, Errorf(EUNAUTHORIZED, "REQUIRES_DEPARTMENT_OR_ORGANIZATION_ADMIN"))
				return
			}
		} else {
			orgRole = thunderdome.AdminUserType
			departmentRole = thunderdome.AdminUserType
		}

		ctx = context.WithValue(ctx, contextKeyOrgRole, orgRole)
		ctx = context.WithValue(ctx, contextKeyDepartmentRole, departmentRole)

		h(w, r.WithContext(ctx))
	}
}

// teamUserOnly validates that the request was made by a valid user of the team
// with bypass for global admins, and if associated to team department and/or organization admins
func (s *Service) teamUserOnly(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		userID := ctx.Value(contextKeyUserID).(string)
		userType := ctx.Value(contextKeyUserType).(string)
		teamID := r.PathValue("teamId")
		idErr := validate.Var(teamID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		roles, err := s.TeamDataSvc.TeamUserRolesByUserID(ctx, userID, teamID)
		if err != nil && err.Error() == "TEAM_NOT_FOUND" {
			s.Logger.Ctx(ctx).Warn("middleware teamUserOnly TEAM_NOT_FOUND",
				zap.Any("team_user_roles", roles),
				zap.String("user_type", userType))
			s.Failure(w, r, http.StatusNotFound, Errorf(ENOTFOUND, "TEAM_NOT_FOUND"))
			return
		} else if err != nil || (userType != thunderdome.AdminUserType &&
			roles.AssociationLevel != "TEAM" &&
			(roles.DepartmentRole == nil || (roles.DepartmentRole != nil && *roles.DepartmentRole != thunderdome.AdminUserType)) &&
			(roles.OrganizationRole == nil || (roles.OrganizationRole != nil && *roles.OrganizationRole != thunderdome.AdminUserType))) {
			s.Logger.Ctx(ctx).Warn("middleware teamUserOnly REQUIRES_TEAM_USER",
				zap.Any("team_user_roles", roles),
				zap.String("user_type", userType))
			s.Failure(w, r, http.StatusForbidden, Errorf(EUNAUTHORIZED, "REQUIRES_TEAM_USER"))
			return
		}

		ctx = context.WithValue(ctx, contextKeyUserTeamRoles, roles)

		h(w, r.WithContext(ctx))
	}
}

// teamAdminOnly validates that the request was made by an ADMIN of the team
// or an ADMIN of the team's parent entities if associated (department or organization)
func (s *Service) teamAdminOnly(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userType := ctx.Value(contextKeyUserType).(string)
		teamUserRoles := ctx.Value(contextKeyUserTeamRoles).(*thunderdome.UserTeamRoleInfo)

		if userType != thunderdome.AdminUserType &&
			(teamUserRoles.TeamRole == nil || *teamUserRoles.TeamRole != thunderdome.AdminUserType) &&
			(teamUserRoles.DepartmentRole == nil || *teamUserRoles.DepartmentRole != thunderdome.AdminUserType) &&
			(teamUserRoles.OrganizationRole == nil || *teamUserRoles.OrganizationRole != thunderdome.AdminUserType) {
			s.Logger.Ctx(ctx).Warn("middleware teamAdminOnly REQUIRES_TEAM_ADMIN",
				zap.Any("team_user_roles", teamUserRoles),
				zap.String("user_type", userType))
			s.Failure(w, r, http.StatusForbidden, Errorf(EUNAUTHORIZED, "REQUIRES_TEAM_ADMIN"))
			return
		}

		h(w, r.WithContext(ctx))
	}
}

// subscribedOrgOnly validates that the request was made by a subscribed organization only
func (s *Service) subscribedOrgOnly(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		userType := ctx.Value(contextKeyUserType).(string)
		orgID := r.PathValue("orgId")
		idErr := validate.Var(orgID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		if !s.Config.SubscriptionsEnabled {
			h(w, r)
			return
		}

		if userType != thunderdome.AdminUserType {
			subscribed, err := s.OrganizationDataSvc.OrganizationIsSubscribed(ctx, orgID)
			if err != nil || !subscribed {
				s.Failure(w, r, http.StatusForbidden, Errorf(EUNAUTHORIZED, "ORGANIZATION_SUBSCRIPTION_REQUIRED"))
				return
			}
		}

		h(w, r.WithContext(ctx))
	}
}

// subscribedTeamOnly validates that the request was made by a subscribed team only
func (s *Service) subscribedTeamOnly(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		userType := ctx.Value(contextKeyUserType).(string)
		teamID := r.PathValue("teamId")
		idErr := validate.Var(teamID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		if !s.Config.SubscriptionsEnabled {
			h(w, r)
			return
		}

		if userType != thunderdome.AdminUserType {
			subscribed, err := s.TeamDataSvc.TeamIsSubscribed(ctx, teamID)
			if err != nil || !subscribed {
				s.Failure(w, r, http.StatusForbidden, Errorf(EUNAUTHORIZED, "TEAM_SUBSCRIPTION_REQUIRED"))
				return
			}
		}

		h(w, r.WithContext(ctx))
	}
}

func (s *Service) subscribedProjectOnly(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		userType := ctx.Value(contextKeyUserType).(string)
		projectID := r.PathValue("projectId")
		idErr := validate.Var(projectID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		if !s.Config.SubscriptionsEnabled {
			h(w, r)
			return
		}

		if userType != thunderdome.AdminUserType {
			subscribed, err := s.SubscriptionDataSvc.ProjectIsSubscribed(ctx, projectID)
			if err != nil || !subscribed {
				s.Failure(w, r, http.StatusForbidden, Errorf(EUNAUTHORIZED, "ORGANIZATION_OR_TEAM_SUBSCRIPTION_REQUIRED"))
				return
			}
		}

		h(w, r.WithContext(ctx))
	}
}

func (s *Service) projectUserOnly(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		userType := ctx.Value(contextKeyUserType).(string)
		projectID := r.PathValue("projectId")
		idErr := validate.Var(projectID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		userID := ctx.Value(contextKeyUserID).(string)

		var role string
		if userType != thunderdome.AdminUserType {
			var isMember bool
			var err error
			isMember, role, err = s.ProjectDataSvc.IsUserProjectMember(ctx, userID, projectID)
			if err != nil || !isMember {
				s.Failure(w, r, http.StatusForbidden, Errorf(EUNAUTHORIZED, "PROJECT_MEMBERSHIP_REQUIRED"))
				return
			}
		}

		ctx = context.WithValue(ctx, contextKeyUserProjectRole, role)

		h(w, r.WithContext(ctx))
	}
}

func (s *Service) projectAdminOnly(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		userType := ctx.Value(contextKeyUserType).(string)
		projectID := r.PathValue("projectId")
		idErr := validate.Var(projectID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		userID := ctx.Value(contextKeyUserID).(string)

		var role string
		if userType != thunderdome.AdminUserType {
			var isMember bool
			var err error

			isMember, role, err := s.ProjectDataSvc.IsUserProjectMember(ctx, userID, projectID)
			if err != nil || !isMember {
				s.Failure(w, r, http.StatusForbidden, Errorf(EUNAUTHORIZED, "PROJECT_MEMBERSHIP_REQUIRED"))
				return
			}
			if role != thunderdome.AdminUserType {
				s.Failure(w, r, http.StatusForbidden, Errorf(EUNAUTHORIZED, "PROJECT_ADMIN_REQUIRED"))
				return
			}
		}

		ctx = context.WithValue(ctx, contextKeyUserProjectRole, role)

		h(w, r.WithContext(ctx))
	}
}
