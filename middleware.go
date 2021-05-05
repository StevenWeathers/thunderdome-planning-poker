package main

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// adminOnly middleware checks if the user is an admin, otherwise reject their request
func (s *server) adminOnly(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get(apiKeyHeaderName)
		apiKey = strings.TrimSpace(apiKey)
		var UserID string

		if apiKey != "" {
			var apiKeyErr error
			UserID, apiKeyErr = s.database.ValidateAPIKey(apiKey)
			if apiKeyErr != nil {
				log.Println("error validating api key : " + apiKeyErr.Error() + "\n")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		} else {
			var cookieErr error
			UserID, cookieErr = s.validateUserCookie(w, r)
			if cookieErr != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		}

		adminErr := s.database.ConfirmAdmin(UserID)
		if adminErr != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), contextKeyUserID, UserID)

		h(w, r.WithContext(ctx))
	}
}

// userOnly validates that the request was made by a valid user
func (s *server) userOnly(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get(apiKeyHeaderName)
		apiKey = strings.TrimSpace(apiKey)
		var UserID string

		if apiKey != "" {
			var apiKeyErr error
			UserID, apiKeyErr = s.database.ValidateAPIKey(apiKey)
			if apiKeyErr != nil {
				log.Println("error validating api key : " + apiKeyErr.Error() + "\n")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		} else {
			var cookieErr error
			UserID, cookieErr = s.validateUserCookie(w, r)
			if cookieErr != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		}

		_, UserErr := s.database.GetUser(UserID)
		if UserErr != nil {
			log.Println("error finding user : " + UserErr.Error() + "\n")
			s.clearUserCookies(w)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), contextKeyUserID, UserID)

		h(w, r.WithContext(ctx))
	}
}

// orgUserOnly validates that the request was made by a valid user of the organization
func (s *server) orgUserOnly(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		UserID := r.Context().Value(contextKeyUserID).(string)
		OrgID := vars["orgId"]

		Role, UserErr := s.database.OrganizationUserRole(UserID, OrgID)
		if UserErr != nil {
			log.Println("error finding user in organization : " + UserErr.Error() + "\n")
			w.WriteHeader(http.StatusForbidden)
			return
		}

		ctx := context.WithValue(r.Context(), contextKeyOrgRole, Role)

		h(w, r.WithContext(ctx))
	}
}

// orgAdminOnly validates that the request was made by an ADMIN of the organization
func (s *server) orgAdminOnly(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		UserID := r.Context().Value(contextKeyUserID).(string)
		OrgID := vars["orgId"]

		Role, UserErr := s.database.OrganizationUserRole(UserID, OrgID)
		if UserErr != nil {
			log.Println("error finding user in organization : " + UserErr.Error() + "\n")
			w.WriteHeader(http.StatusForbidden)
			return
		}
		if Role != "ADMIN" {
			log.Println("user is not an ADMIN of organization")
			w.WriteHeader(http.StatusForbidden)
			return
		}

		ctx := context.WithValue(r.Context(), contextKeyOrgRole, Role)

		h(w, r.WithContext(ctx))
	}
}

// orgTeamOnly validates that the request was made by an user of the organization team (or organization)
func (s *server) orgTeamOnly(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		UserID := r.Context().Value(contextKeyUserID).(string)
		OrgID := vars["orgId"]
		TeamID := vars["teamId"]

		OrgRole, TeamRole, UserErr := s.database.OrganizationTeamUserRole(UserID, OrgID, TeamID)
		if UserErr != nil {
			log.Println("error finding user in organization : " + UserErr.Error() + "\n")
			w.WriteHeader(http.StatusForbidden)
			return
		}

		ctx := context.WithValue(r.Context(), contextKeyOrgRole, OrgRole)
		ctx = context.WithValue(ctx, contextKeyTeamRole, TeamRole)

		h(w, r.WithContext(ctx))
	}
}

// orgTeamAdminOnly validates that the request was made by an ADMIN of the organization team (or organization)
func (s *server) orgTeamAdminOnly(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		UserID := r.Context().Value(contextKeyUserID).(string)
		OrgID := vars["orgId"]
		TeamID := vars["teamId"]

		OrgRole, TeamRole, UserErr := s.database.OrganizationTeamUserRole(UserID, OrgID, TeamID)
		if UserErr != nil {
			log.Println("error finding user in organization : " + UserErr.Error() + "\n")
			w.WriteHeader(http.StatusForbidden)
			return
		}
		if TeamRole != "ADMIN" && OrgRole != "ADMIN" {
			log.Println("user is not an ADMIN of organization")
			w.WriteHeader(http.StatusForbidden)
			return
		}

		ctx := context.WithValue(r.Context(), contextKeyOrgRole, OrgRole)
		ctx = context.WithValue(ctx, contextKeyTeamRole, TeamRole)

		h(w, r.WithContext(ctx))
	}
}

// departmentUserOnly validates that the request was made by a valid user of the organization (with department role)
func (s *server) departmentUserOnly(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		UserID := r.Context().Value(contextKeyUserID).(string)
		OrgID := vars["orgId"]
		DepartmentID := vars["departmentId"]

		OrgRole, DepartmentRole, UserErr := s.database.DepartmentUserRole(UserID, OrgID, DepartmentID)
		if UserErr != nil {
			log.Println("error finding user in organization : " + UserErr.Error() + "\n")
			w.WriteHeader(http.StatusForbidden)
			return
		}

		ctx := context.WithValue(r.Context(), contextKeyOrgRole, OrgRole)
		ctx = context.WithValue(ctx, contextKeyDepartmentRole, DepartmentRole)

		h(w, r.WithContext(ctx))
	}
}

// departmentAdminOnly validates that the request was made by an ADMIN of the organization (with department role)
func (s *server) departmentAdminOnly(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		UserID := r.Context().Value(contextKeyUserID).(string)
		OrgID := vars["orgId"]
		DepartmentID := vars["departmentId"]

		OrgRole, DepartmentRole, UserErr := s.database.DepartmentUserRole(UserID, OrgID, DepartmentID)
		if UserErr != nil {
			log.Println("error finding user in organization : " + UserErr.Error() + "\n")
			w.WriteHeader(http.StatusForbidden)
			return
		}
		if DepartmentRole != "ADMIN" && OrgRole != "ADMIN" {
			log.Println("user is not an ADMIN of department or organization")
			w.WriteHeader(http.StatusForbidden)
			return
		}

		ctx := context.WithValue(r.Context(), contextKeyOrgRole, OrgRole)
		ctx = context.WithValue(ctx, contextKeyDepartmentRole, DepartmentRole)

		h(w, r.WithContext(ctx))
	}
}

// departmentTeamUserOnly validates that the request was made by an user of the department team (or organization)
func (s *server) departmentTeamUserOnly(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		UserID := r.Context().Value(contextKeyUserID).(string)
		OrgID := vars["orgId"]
		DepartmentID := vars["departmentId"]
		TeamID := vars["teamId"]

		OrgRole, DepartmentRole, TeamRole, UserErr := s.database.DepartmentTeamUserRole(UserID, OrgID, DepartmentID, TeamID)
		if UserErr != nil {
			log.Println("error finding user in department team : " + UserErr.Error() + "\n")
			w.WriteHeader(http.StatusForbidden)
			return
		}

		ctx := context.WithValue(r.Context(), contextKeyOrgRole, OrgRole)
		ctx = context.WithValue(ctx, contextKeyDepartmentRole, DepartmentRole)
		ctx = context.WithValue(ctx, contextKeyTeamRole, TeamRole)

		h(w, r.WithContext(ctx))
	}
}

// departmentTeamAdminOnly validates that the request was made by an ADMIN of the department team (or organization)
func (s *server) departmentTeamAdminOnly(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		UserID := r.Context().Value(contextKeyUserID).(string)
		OrgID := vars["orgId"]
		DepartmentID := vars["departmentId"]
		TeamID := vars["teamId"]

		OrgRole, DepartmentRole, TeamRole, UserErr := s.database.DepartmentTeamUserRole(UserID, OrgID, DepartmentID, TeamID)
		if UserErr != nil {
			log.Println("error finding user in department team : " + UserErr.Error() + "\n")
			w.WriteHeader(http.StatusForbidden)
			return
		}

		if TeamRole != "ADMIN" && DepartmentRole != "ADMIN" && OrgRole != "ADMIN" {
			log.Println("user is not an ADMIN of organization")
			w.WriteHeader(http.StatusForbidden)
			return
		}

		ctx := context.WithValue(r.Context(), contextKeyOrgRole, OrgRole)
		ctx = context.WithValue(ctx, contextKeyDepartmentRole, DepartmentRole)
		ctx = context.WithValue(ctx, contextKeyTeamRole, TeamRole)

		h(w, r.WithContext(ctx))
	}
}

// teamUserOnly validates that the request was made by a valid user of the team
func (s *server) teamUserOnly(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		UserID := r.Context().Value(contextKeyUserID).(string)
		TeamID := vars["teamId"]

		Role, UserErr := s.database.TeamUserRole(UserID, TeamID)
		if UserErr != nil {
			log.Println("error finding user in team : " + UserErr.Error() + "\n")
			w.WriteHeader(http.StatusForbidden)
			return
		}

		ctx := context.WithValue(r.Context(), contextKeyTeamRole, Role)

		h(w, r.WithContext(ctx))
	}
}

// teamAdminOnly validates that the request was made by an ADMIN of the team
func (s *server) teamAdminOnly(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		UserID := r.Context().Value(contextKeyUserID).(string)
		TeamID := vars["teamId"]

		Role, UserErr := s.database.TeamUserRole(UserID, TeamID)
		if UserErr != nil {
			log.Println("error finding user in team : " + UserErr.Error() + "\n")
			w.WriteHeader(http.StatusForbidden)
			return
		}
		if Role != "ADMIN" {
			log.Println("user is not an ADMIN of team")
			w.WriteHeader(http.StatusForbidden)
			return
		}

		ctx := context.WithValue(r.Context(), contextKeyTeamRole, Role)

		h(w, r.WithContext(ctx))
	}
}
