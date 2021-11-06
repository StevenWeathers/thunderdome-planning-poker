package api

import (
	"context"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// adminOnly middleware checks if the user is an admin, otherwise reject their request
func (a *api) adminOnly(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get(apiKeyHeaderName)
		apiKey = strings.TrimSpace(apiKey)
		var UserID string

		if apiKey != "" {
			var apiKeyErr error
			UserID, apiKeyErr = a.db.ValidateAPIKey(apiKey)
			if apiKeyErr != nil {
				Failure(w, r, http.StatusUnauthorized, Errorf(EINVALID, "INVALID_APIKEY"))
				return
			}
		} else {
			var cookieErr error
			UserID, cookieErr = a.validateUserCookie(w, r)
			if cookieErr != nil {
				Failure(w, r, http.StatusUnauthorized, Errorf(EINVALID, "INVALID_USER"))
				return
			}
		}

		adminErr := a.db.ConfirmAdmin(UserID)
		if adminErr != nil {
			Failure(w, r, http.StatusForbidden, Errorf(EUNAUTHORIZED, "REQUIRES_ADMIN"))
			return
		}

		ctx := context.WithValue(r.Context(), contextKeyUserID, UserID)
		ctx = context.WithValue(ctx, contextKeyUserType, adminUserType)

		h(w, r.WithContext(ctx))
	}
}

// userOnly validates that the request was made by a valid user
func (a *api) userOnly(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get(apiKeyHeaderName)
		apiKey = strings.TrimSpace(apiKey)
		var UserID string

		if apiKey != "" && a.config.ExternalAPIEnabled == true {
			var apiKeyErr error
			UserID, apiKeyErr = a.db.ValidateAPIKey(apiKey)
			if apiKeyErr != nil {
				Failure(w, r, http.StatusUnauthorized, Errorf(EINVALID, "INVALID_APIKEY"))
				return
			}
		} else {
			var cookieErr error
			UserID, cookieErr = a.validateUserCookie(w, r)
			if cookieErr != nil {
				Failure(w, r, http.StatusUnauthorized, Errorf(EINVALID, "INVALID_USER"))
				return
			}
		}

		User, UserErr := a.db.GetUser(UserID)
		if UserErr != nil {
			a.clearUserCookies(w)
			Failure(w, r, http.StatusUnauthorized, Errorf(EINVALID, "INVALID_USER"))
			return
		}

		ctx := context.WithValue(r.Context(), contextKeyUserID, UserID)
		ctx = context.WithValue(ctx, contextKeyUserType, User.UserType)

		h(w, r.WithContext(ctx))
	}
}

// entityUserOnly validates that the request was made by the session user matching the {userId} of the entity (or ADMIN)
func (a *api) entityUserOnly(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		apiKey := r.Header.Get(apiKeyHeaderName)
		apiKey = strings.TrimSpace(apiKey)
		var UserID string
		EntityUserID := vars["userId"]

		if apiKey != "" && a.config.ExternalAPIEnabled == true {
			var apiKeyErr error
			UserID, apiKeyErr = a.db.ValidateAPIKey(apiKey)
			if apiKeyErr != nil {
				Failure(w, r, http.StatusUnauthorized, Errorf(EINVALID, "INVALID_APIKEY"))
				return
			}
		} else {
			var cookieErr error
			UserID, cookieErr = a.validateUserCookie(w, r)
			if cookieErr != nil {
				Failure(w, r, http.StatusUnauthorized, Errorf(EINVALID, "INVALID_USER"))
				return
			}
		}

		User, UserErr := a.db.GetUser(UserID)
		if UserErr != nil {
			a.clearUserCookies(w)
			Failure(w, r, http.StatusUnauthorized, Errorf(EINVALID, "INVALID_USER"))
			return
		}

		if User.UserType != adminUserType && EntityUserID != UserID {
			Failure(w, r, http.StatusForbidden, Errorf(EINVALID, "INVALID_USER"))
			return
		}

		ctx := context.WithValue(r.Context(), contextKeyUserID, UserID)
		ctx = context.WithValue(ctx, contextKeyUserType, User.UserType)

		h(w, r.WithContext(ctx))
	}
}

// verifiedUserOnly validates that the request was made by a verified registered user
func (a *api) verifiedUserOnly(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		apiKey := r.Header.Get(apiKeyHeaderName)
		apiKey = strings.TrimSpace(apiKey)
		var UserID string
		EntityUserID := vars["userId"]

		if apiKey != "" && a.config.ExternalAPIEnabled == true {
			var apiKeyErr error
			UserID, apiKeyErr = a.db.ValidateAPIKey(apiKey)
			if apiKeyErr != nil {
				Failure(w, r, http.StatusUnauthorized, Errorf(EINVALID, "INVALID_APIKEY"))
				return
			}
		} else {
			var cookieErr error
			UserID, cookieErr = a.validateUserCookie(w, r)
			if cookieErr != nil {
				Failure(w, r, http.StatusUnauthorized, Errorf(EINVALID, "INVALID_USER"))
				return
			}
		}

		User, UserErr := a.db.GetUser(UserID)
		if UserErr != nil {
			a.clearUserCookies(w)
			Failure(w, r, http.StatusUnauthorized, Errorf(EINVALID, "INVALID_USER"))
			return
		}

		if User.UserType != adminUserType && EntityUserID != UserID && User.Verified == false {
			Failure(w, r, http.StatusForbidden, Errorf(EUNAUTHORIZED, "REQUIRES_VERIFIED_USER"))
			return
		}

		ctx := context.WithValue(r.Context(), contextKeyUserID, UserID)
		ctx = context.WithValue(ctx, contextKeyUserType, User.UserType)

		h(w, r.WithContext(ctx))
	}
}

// orgUserOnly validates that the request was made by a valid user of the organization
func (a *api) orgUserOnly(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		UserID := r.Context().Value(contextKeyUserID).(string)
		UserType := r.Context().Value(contextKeyUserType).(string)
		OrgID := vars["orgId"]

		Role, UserErr := a.db.OrganizationUserRole(UserID, OrgID)
		if UserType != adminUserType && UserErr != nil {
			Failure(w, r, http.StatusForbidden, Errorf(EUNAUTHORIZED, "ORGANIZATION_USER_REQUIRED"))
			return
		}

		ctx := context.WithValue(r.Context(), contextKeyOrgRole, Role)

		h(w, r.WithContext(ctx))
	}
}

// orgAdminOnly validates that the request was made by an ADMIN of the organization
func (a *api) orgAdminOnly(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		UserID := r.Context().Value(contextKeyUserID).(string)
		UserType := r.Context().Value(contextKeyUserType).(string)
		OrgID := vars["orgId"]

		Role, UserErr := a.db.OrganizationUserRole(UserID, OrgID)
		if UserType != adminUserType && UserErr != nil {
			Failure(w, r, http.StatusForbidden, Errorf(EUNAUTHORIZED, "ORGANIZATION_USER_REQUIRED"))
			return
		}
		if UserType != adminUserType && Role != "ADMIN" {
			Failure(w, r, http.StatusForbidden, Errorf(EUNAUTHORIZED, "REQUIRES_ORG_ADMIN"))
			return
		}

		ctx := context.WithValue(r.Context(), contextKeyOrgRole, Role)

		h(w, r.WithContext(ctx))
	}
}

// orgTeamOnly validates that the request was made by an user of the organization team (or organization)
func (a *api) orgTeamOnly(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		UserID := r.Context().Value(contextKeyUserID).(string)
		UserType := r.Context().Value(contextKeyUserType).(string)
		OrgID := vars["orgId"]
		TeamID := vars["teamId"]

		OrgRole, TeamRole, UserErr := a.db.OrganizationTeamUserRole(UserID, OrgID, TeamID)
		if UserType != adminUserType && UserErr != nil {
			Failure(w, r, http.StatusForbidden, Errorf(EUNAUTHORIZED, "REQUIRES_TEAM_USER"))
			return
		}

		ctx := context.WithValue(r.Context(), contextKeyOrgRole, OrgRole)
		ctx = context.WithValue(ctx, contextKeyTeamRole, TeamRole)

		h(w, r.WithContext(ctx))
	}
}

// orgTeamAdminOnly validates that the request was made by an ADMIN of the organization team (or organization)
func (a *api) orgTeamAdminOnly(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		UserID := r.Context().Value(contextKeyUserID).(string)
		UserType := r.Context().Value(contextKeyUserType).(string)
		OrgID := vars["orgId"]
		TeamID := vars["teamId"]

		OrgRole, TeamRole, UserErr := a.db.OrganizationTeamUserRole(UserID, OrgID, TeamID)
		if UserType != adminUserType && UserErr != nil {
			Failure(w, r, http.StatusForbidden, Errorf(EUNAUTHORIZED, "REQUIRES_TEAM_USER"))
			return
		}
		if UserType != adminUserType && TeamRole != "ADMIN" && OrgRole != "ADMIN" {
			Failure(w, r, http.StatusForbidden, Errorf(EUNAUTHORIZED, "REQUIRES_TEAM_OR_ORGANIZATION_ADMIN"))
			return
		}

		ctx := context.WithValue(r.Context(), contextKeyOrgRole, OrgRole)
		ctx = context.WithValue(ctx, contextKeyTeamRole, TeamRole)

		h(w, r.WithContext(ctx))
	}
}

// departmentUserOnly validates that the request was made by a valid user of the organization (with department role)
func (a *api) departmentUserOnly(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		UserID := r.Context().Value(contextKeyUserID).(string)
		UserType := r.Context().Value(contextKeyUserType).(string)
		OrgID := vars["orgId"]
		DepartmentID := vars["departmentId"]

		OrgRole, DepartmentRole, UserErr := a.db.DepartmentUserRole(UserID, OrgID, DepartmentID)
		if UserType != adminUserType && UserErr != nil {
			Failure(w, r, http.StatusForbidden, Errorf(EUNAUTHORIZED, "REQUIRES_DEPARTMENT_USER"))
			return
		}

		ctx := context.WithValue(r.Context(), contextKeyOrgRole, OrgRole)
		ctx = context.WithValue(ctx, contextKeyDepartmentRole, DepartmentRole)

		h(w, r.WithContext(ctx))
	}
}

// departmentAdminOnly validates that the request was made by an ADMIN of the organization (with department role)
func (a *api) departmentAdminOnly(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		UserID := r.Context().Value(contextKeyUserID).(string)
		UserType := r.Context().Value(contextKeyUserType).(string)
		OrgID := vars["orgId"]
		DepartmentID := vars["departmentId"]

		OrgRole, DepartmentRole, UserErr := a.db.DepartmentUserRole(UserID, OrgID, DepartmentID)
		if UserType != adminUserType && UserErr != nil {
			Failure(w, r, http.StatusForbidden, Errorf(EUNAUTHORIZED, "REQUIRES_DEPARTMENT_USER"))
			return
		}
		if UserType != adminUserType && DepartmentRole != "ADMIN" && OrgRole != "ADMIN" {
			Failure(w, r, http.StatusForbidden, Errorf(EUNAUTHORIZED, "REQUIRES_DEPARTMENT_OR_ORGANIZATION_ADMIN"))
			return
		}

		ctx := context.WithValue(r.Context(), contextKeyOrgRole, OrgRole)
		ctx = context.WithValue(ctx, contextKeyDepartmentRole, DepartmentRole)

		h(w, r.WithContext(ctx))
	}
}

// departmentTeamUserOnly validates that the request was made by an user of the department team (or organization)
func (a *api) departmentTeamUserOnly(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		UserID := r.Context().Value(contextKeyUserID).(string)
		UserType := r.Context().Value(contextKeyUserType).(string)
		OrgID := vars["orgId"]
		DepartmentID := vars["departmentId"]
		TeamID := vars["teamId"]

		OrgRole, DepartmentRole, TeamRole, UserErr := a.db.DepartmentTeamUserRole(UserID, OrgID, DepartmentID, TeamID)
		if UserType != adminUserType && UserErr != nil {
			Failure(w, r, http.StatusForbidden, Errorf(EUNAUTHORIZED, "REQUIRES_TEAM_USER"))
			return
		}

		ctx := context.WithValue(r.Context(), contextKeyOrgRole, OrgRole)
		ctx = context.WithValue(ctx, contextKeyDepartmentRole, DepartmentRole)
		ctx = context.WithValue(ctx, contextKeyTeamRole, TeamRole)

		h(w, r.WithContext(ctx))
	}
}

// departmentTeamAdminOnly validates that the request was made by an ADMIN of the department team (or organization)
func (a *api) departmentTeamAdminOnly(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		UserID := r.Context().Value(contextKeyUserID).(string)
		UserType := r.Context().Value(contextKeyUserType).(string)
		OrgID := vars["orgId"]
		DepartmentID := vars["departmentId"]
		TeamID := vars["teamId"]

		OrgRole, DepartmentRole, TeamRole, UserErr := a.db.DepartmentTeamUserRole(UserID, OrgID, DepartmentID, TeamID)
		if UserType != adminUserType && UserErr != nil {
			Failure(w, r, http.StatusForbidden, Errorf(EUNAUTHORIZED, "REQUIRES_TEAM_USER"))
			return
		}

		if UserType != adminUserType && TeamRole != "ADMIN" && DepartmentRole != "ADMIN" && OrgRole != "ADMIN" {
			Failure(w, r, http.StatusForbidden, Errorf(EUNAUTHORIZED, "REQUIRES_TEAM_OR_DEPARTMENT_OR_ORGANIZATION_ADMIN"))
			return
		}

		ctx := context.WithValue(r.Context(), contextKeyOrgRole, OrgRole)
		ctx = context.WithValue(ctx, contextKeyDepartmentRole, DepartmentRole)
		ctx = context.WithValue(ctx, contextKeyTeamRole, TeamRole)

		h(w, r.WithContext(ctx))
	}
}

// teamUserOnly validates that the request was made by a valid user of the team
func (a *api) teamUserOnly(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		UserID := r.Context().Value(contextKeyUserID).(string)
		UserType := r.Context().Value(contextKeyUserType).(string)
		TeamID := vars["teamId"]

		Role, UserErr := a.db.TeamUserRole(UserID, TeamID)
		if UserType != adminUserType && UserErr != nil {
			Failure(w, r, http.StatusForbidden, Errorf(EUNAUTHORIZED, "REQUIRES_TEAM_USER"))
			return
		}

		ctx := context.WithValue(r.Context(), contextKeyTeamRole, Role)

		h(w, r.WithContext(ctx))
	}
}

// teamAdminOnly validates that the request was made by an ADMIN of the team
func (a *api) teamAdminOnly(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		UserID := r.Context().Value(contextKeyUserID).(string)
		UserType := r.Context().Value(contextKeyUserType).(string)
		TeamID := vars["teamId"]

		Role, UserErr := a.db.TeamUserRole(UserID, TeamID)
		if UserType != adminUserType && UserErr != nil {
			Failure(w, r, http.StatusForbidden, Errorf(EUNAUTHORIZED, "REQUIRES_TEAM_USER"))
			return
		}
		if UserType != adminUserType && Role != "ADMIN" {
			Failure(w, r, http.StatusForbidden, Errorf(EUNAUTHORIZED, "REQUIRES_TEAM_ADMIN"))
			return
		}

		ctx := context.WithValue(r.Context(), contextKeyTeamRole, Role)

		h(w, r.WithContext(ctx))
	}
}
