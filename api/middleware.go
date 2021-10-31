package api

import (
	"context"
	"log"
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
				log.Println("error validating api key : " + apiKeyErr.Error() + "\n")
				a.respondWithStandardJSON(w, http.StatusUnauthorized, false, nil, nil, nil)
				return
			}
		} else {
			var cookieErr error
			UserID, cookieErr = a.validateUserCookie(w, r)
			if cookieErr != nil {
				a.respondWithStandardJSON(w, http.StatusUnauthorized, false, nil, nil, nil)
				return
			}
		}

		adminErr := a.db.ConfirmAdmin(UserID)
		if adminErr != nil {
			a.respondWithStandardJSON(w, http.StatusForbidden, false, nil, nil, nil)
			return
		}

		ctx := context.WithValue(r.Context(), contextKeyUserID, UserID)

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
				log.Println("error validating api key : " + apiKeyErr.Error() + "\n")
				a.respondWithStandardJSON(w, http.StatusUnauthorized, false, nil, nil, nil)
				return
			}
		} else {
			var cookieErr error
			UserID, cookieErr = a.validateUserCookie(w, r)
			if cookieErr != nil {
				a.respondWithStandardJSON(w, http.StatusUnauthorized, false, nil, nil, nil)
				return
			}
		}

		_, UserErr := a.db.GetUser(UserID)
		if UserErr != nil {
			log.Println("error finding user : " + UserErr.Error() + "\n")
			a.clearUserCookies(w)
			a.respondWithStandardJSON(w, http.StatusUnauthorized, false, nil, nil, nil)
			return
		}

		ctx := context.WithValue(r.Context(), contextKeyUserID, UserID)

		h(w, r.WithContext(ctx))
	}
}

// verifiedUserOnly validates that the request was made by a verified registered user
func (a *api) verifiedUserOnly(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get(apiKeyHeaderName)
		apiKey = strings.TrimSpace(apiKey)
		var UserID string

		if apiKey != "" && a.config.ExternalAPIEnabled == true {
			var apiKeyErr error
			UserID, apiKeyErr = a.db.ValidateAPIKey(apiKey)
			if apiKeyErr != nil {
				log.Println("error validating api key : " + apiKeyErr.Error() + "\n")
				a.respondWithStandardJSON(w, http.StatusUnauthorized, false, nil, nil, nil)
				return
			}
		} else {
			var cookieErr error
			UserID, cookieErr = a.validateUserCookie(w, r)
			if cookieErr != nil {
				a.respondWithStandardJSON(w, http.StatusUnauthorized, false, nil, nil, nil)
				return
			}
		}

		User, UserErr := a.db.GetUser(UserID)
		if UserErr != nil {
			log.Println("error finding user : " + UserErr.Error() + "\n")
			a.clearUserCookies(w)
			a.respondWithStandardJSON(w, http.StatusUnauthorized, false, nil, nil, nil)
			return
		}

		if User.Verified == false {
			errors := make([]string, 0)
			errors = append(errors, "USER_NOT_VERIFIED")
			a.respondWithStandardJSON(w, http.StatusForbidden, false, errors, nil, nil)
			return
		}

		ctx := context.WithValue(r.Context(), contextKeyUserID, UserID)

		h(w, r.WithContext(ctx))
	}
}

// orgUserOnly validates that the request was made by a valid user of the organization
func (a *api) orgUserOnly(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		UserID := r.Context().Value(contextKeyUserID).(string)
		OrgID := vars["orgId"]

		Role, UserErr := a.db.OrganizationUserRole(UserID, OrgID)
		if UserErr != nil {
			log.Println("error finding user in organization : " + UserErr.Error() + "\n")
			a.respondWithStandardJSON(w, http.StatusForbidden, false, nil, nil, nil)
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
		OrgID := vars["orgId"]

		Role, UserErr := a.db.OrganizationUserRole(UserID, OrgID)
		if UserErr != nil {
			log.Println("error finding user in organization : " + UserErr.Error() + "\n")
			a.respondWithStandardJSON(w, http.StatusForbidden, false, nil, nil, nil)
			return
		}
		if Role != "ADMIN" {
			log.Println("user is not an ADMIN of organization")
			a.respondWithStandardJSON(w, http.StatusForbidden, false, nil, nil, nil)
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
		OrgID := vars["orgId"]
		TeamID := vars["teamId"]

		OrgRole, TeamRole, UserErr := a.db.OrganizationTeamUserRole(UserID, OrgID, TeamID)
		if UserErr != nil {
			log.Println("error finding user in organization : " + UserErr.Error() + "\n")
			a.respondWithStandardJSON(w, http.StatusForbidden, false, nil, nil, nil)
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
		OrgID := vars["orgId"]
		TeamID := vars["teamId"]

		OrgRole, TeamRole, UserErr := a.db.OrganizationTeamUserRole(UserID, OrgID, TeamID)
		if UserErr != nil {
			log.Println("error finding user in organization : " + UserErr.Error() + "\n")
			a.respondWithStandardJSON(w, http.StatusForbidden, false, nil, nil, nil)
			return
		}
		if TeamRole != "ADMIN" && OrgRole != "ADMIN" {
			log.Println("user is not an ADMIN of organization")
			a.respondWithStandardJSON(w, http.StatusForbidden, false, nil, nil, nil)
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
		OrgID := vars["orgId"]
		DepartmentID := vars["departmentId"]

		OrgRole, DepartmentRole, UserErr := a.db.DepartmentUserRole(UserID, OrgID, DepartmentID)
		if UserErr != nil {
			log.Println("error finding user in organization : " + UserErr.Error() + "\n")
			a.respondWithStandardJSON(w, http.StatusForbidden, false, nil, nil, nil)
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
		OrgID := vars["orgId"]
		DepartmentID := vars["departmentId"]

		OrgRole, DepartmentRole, UserErr := a.db.DepartmentUserRole(UserID, OrgID, DepartmentID)
		if UserErr != nil {
			log.Println("error finding user in organization : " + UserErr.Error() + "\n")
			a.respondWithStandardJSON(w, http.StatusForbidden, false, nil, nil, nil)
			return
		}
		if DepartmentRole != "ADMIN" && OrgRole != "ADMIN" {
			log.Println("user is not an ADMIN of department or organization")
			a.respondWithStandardJSON(w, http.StatusForbidden, false, nil, nil, nil)
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
		OrgID := vars["orgId"]
		DepartmentID := vars["departmentId"]
		TeamID := vars["teamId"]

		OrgRole, DepartmentRole, TeamRole, UserErr := a.db.DepartmentTeamUserRole(UserID, OrgID, DepartmentID, TeamID)
		if UserErr != nil {
			log.Println("error finding user in department team : " + UserErr.Error() + "\n")
			a.respondWithStandardJSON(w, http.StatusForbidden, false, nil, nil, nil)
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
		OrgID := vars["orgId"]
		DepartmentID := vars["departmentId"]
		TeamID := vars["teamId"]

		OrgRole, DepartmentRole, TeamRole, UserErr := a.db.DepartmentTeamUserRole(UserID, OrgID, DepartmentID, TeamID)
		if UserErr != nil {
			log.Println("error finding user in department team : " + UserErr.Error() + "\n")
			a.respondWithStandardJSON(w, http.StatusForbidden, false, nil, nil, nil)
			return
		}

		if TeamRole != "ADMIN" && DepartmentRole != "ADMIN" && OrgRole != "ADMIN" {
			log.Println("user is not an ADMIN of organization")
			a.respondWithStandardJSON(w, http.StatusForbidden, false, nil, nil, nil)
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
		TeamID := vars["teamId"]

		Role, UserErr := a.db.TeamUserRole(UserID, TeamID)
		if UserErr != nil {
			log.Println("error finding user in team : " + UserErr.Error() + "\n")
			a.respondWithStandardJSON(w, http.StatusForbidden, false, nil, nil, nil)
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
		TeamID := vars["teamId"]

		Role, UserErr := a.db.TeamUserRole(UserID, TeamID)
		if UserErr != nil {
			log.Println("error finding user in team : " + UserErr.Error() + "\n")
			a.respondWithStandardJSON(w, http.StatusForbidden, false, nil, nil, nil)
			return
		}
		if Role != "ADMIN" {
			log.Println("user is not an ADMIN of team")
			a.respondWithStandardJSON(w, http.StatusForbidden, false, nil, nil, nil)
			return
		}

		ctx := context.WithValue(r.Context(), contextKeyTeamRole, Role)

		h(w, r.WithContext(ctx))
	}
}
