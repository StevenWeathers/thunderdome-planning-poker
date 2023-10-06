package http

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"go.uber.org/zap"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"

	"github.com/gorilla/mux"
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
		apiKey := r.Header.Get(apiKeyHeaderName)
		apiKey = strings.TrimSpace(apiKey)
		ctx := r.Context()
		var User *thunderdome.User

		if apiKey != "" && s.Config.ExternalAPIEnabled {
			var apiKeyErr error
			User, apiKeyErr = s.ApiKeyDataSvc.GetApiKeyUser(ctx, apiKey)
			if apiKeyErr != nil {
				s.Failure(w, r, http.StatusUnauthorized, Errorf(EINVALID, "INVALID_APIKEY"))
				return
			}
		} else {
			SessionId, cookieErr := s.Cookie.ValidateSessionCookie(w, r)
			if cookieErr != nil && cookieErr.Error() != "NO_SESSION_COOKIE" {
				s.Failure(w, r, http.StatusUnauthorized, Errorf(EINVALID, "INVALID_USER"))
				return
			}

			if SessionId != "" {
				var userErr error
				User, userErr = s.AuthDataSvc.GetSessionUser(ctx, SessionId)
				if userErr != nil {
					s.Failure(w, r, http.StatusUnauthorized, Errorf(EINVALID, "INVALID_USER"))
					return
				}
			} else {
				UserID, err := s.Cookie.ValidateUserCookie(w, r)
				if err != nil {
					s.Failure(w, r, http.StatusUnauthorized, Errorf(EINVALID, "INVALID_USER"))
					return
				}

				var userErr error
				User, userErr = s.UserDataSvc.GetGuestUser(ctx, UserID)
				if userErr != nil {
					s.Failure(w, r, http.StatusUnauthorized, Errorf(EINVALID, "INVALID_USER"))
					return
				}
			}
		}

		ctx = context.WithValue(ctx, contextKeyUserID, User.Id)
		ctx = context.WithValue(ctx, contextKeyUserType, User.Type)

		h(w, r.WithContext(ctx))
	}
}

// entityUserOnly validates that the request was made by the session user matching the {userId} of the entity (or ADMIN)
func (s *Service) entityUserOnly(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		ctx := r.Context()
		UserID := ctx.Value(contextKeyUserID).(string)
		UserType := ctx.Value(contextKeyUserType).(string)
		EntityUserID := vars["userId"]
		idErr := validate.Var(EntityUserID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		if UserType != adminUserType && EntityUserID != UserID {
			s.Failure(w, r, http.StatusForbidden, Errorf(EINVALID, "INVALID_USER"))
			return
		}

		h(w, r)
	}
}

// registeredUserOnly validates that the request was made by a registered user
func (s *Service) registeredUserOnly(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		UserType := r.Context().Value(contextKeyUserType).(string)

		if UserType == guestUserType {
			s.Failure(w, r, http.StatusForbidden, Errorf(EINVALID, "REGISTERED_USER_ONLY"))
			return
		}

		h(w, r)
	}
}

// adminOnly middleware checks if the user is an admin, otherwise reject their request
func (s *Service) adminOnly(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		UserType := r.Context().Value(contextKeyUserType).(string)

		if UserType != adminUserType {
			s.Failure(w, r, http.StatusForbidden, Errorf(EUNAUTHORIZED, "REQUIRES_ADMIN"))
			return
		}

		h(w, r)
	}
}

// verifiedUserOnly validates that the request was made by a verified registered user
func (s *Service) verifiedUserOnly(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		ctx := r.Context()
		SessionUserID := ctx.Value(contextKeyUserID).(string)
		UserType := ctx.Value(contextKeyUserType).(string)
		EntityUserID := vars["userId"]
		idErr := validate.Var(EntityUserID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		if UserType != adminUserType && (EntityUserID != SessionUserID) {
			s.Failure(w, r, http.StatusForbidden, Errorf(EINVALID, "INVALID_USER"))
			return
		}

		EntityUser, EntityUserErr := s.UserDataSvc.GetUser(ctx, EntityUserID)
		if EntityUserErr != nil {
			s.Logger.Ctx(ctx).Error(
				"verifiedUserOnly error", zap.Error(EntityUserErr), zap.String("entity_user_id", EntityUserID),
				zap.String("session_user_id", SessionUserID), zap.String("session_user_type", UserType))
			s.Failure(w, r, http.StatusInternalServerError, EntityUserErr)
			return
		}

		if s.Config.ExternalAPIVerifyRequired && !EntityUser.Verified {
			s.Failure(w, r, http.StatusForbidden, Errorf(EUNAUTHORIZED, "REQUIRES_VERIFIED_USER"))
			return
		}

		h(w, r)
	}
}

// subscribedUserOnly validates that the request was made by a subscribed user
func (s *Service) subscribedUserOnly(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		ctx := r.Context()
		UserID := ctx.Value(contextKeyUserID).(string)
		UserType := ctx.Value(contextKeyUserType).(string)
		EntityUserID := vars["userId"]
		idErr := validate.Var(EntityUserID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		if UserType != adminUserType && (EntityUserID != UserID) {
			s.Failure(w, r, http.StatusForbidden, Errorf(EINVALID, "INVALID_USER"))
			return
		}

		subscriberErr := s.SubscriptionDataSvc.CheckActiveSubscriber(ctx, EntityUserID)
		if subscriberErr != nil {
			s.Failure(w, r, http.StatusForbidden, Errorf(EUNAUTHORIZED, "REQUIRES_SUBSCRIBED_USER"))
			return
		}

		h(w, r)
	}
}

// orgUserOnly validates that the request was made by a valid user of the organization
func (s *Service) orgUserOnly(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		ctx := r.Context()
		UserID := ctx.Value(contextKeyUserID).(string)
		UserType := ctx.Value(contextKeyUserType).(string)
		OrgID := vars["orgId"]
		idErr := validate.Var(OrgID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		var Role string
		if UserType != adminUserType {
			var UserErr error
			Role, UserErr = s.OrganizationDataSvc.OrganizationUserRole(ctx, UserID, OrgID)
			if UserErr != nil {
				s.Failure(w, r, http.StatusForbidden, Errorf(EUNAUTHORIZED, "ORGANIZATION_USER_REQUIRED"))
				return
			}
		} else {
			Role = adminUserType
		}

		ctx = context.WithValue(ctx, contextKeyOrgRole, Role)

		h(w, r.WithContext(ctx))
	}
}

// orgAdminOnly validates that the request was made by an ADMIN of the organization
func (s *Service) orgAdminOnly(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		ctx := r.Context()
		UserID := ctx.Value(contextKeyUserID).(string)
		UserType := ctx.Value(contextKeyUserType).(string)
		OrgID := vars["orgId"]
		idErr := validate.Var(OrgID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		var Role string
		if UserType != adminUserType {
			var UserErr error
			Role, UserErr := s.OrganizationDataSvc.OrganizationUserRole(ctx, UserID, OrgID)
			if UserErr != nil {
				s.Failure(w, r, http.StatusForbidden, Errorf(EUNAUTHORIZED, "ORGANIZATION_USER_REQUIRED"))
				return
			}
			if Role != adminUserType {
				s.Failure(w, r, http.StatusForbidden, Errorf(EUNAUTHORIZED, "REQUIRES_ORG_ADMIN"))
				return
			}
		} else {
			Role = adminUserType
		}

		ctx = context.WithValue(ctx, contextKeyOrgRole, Role)

		h(w, r.WithContext(ctx))
	}
}

// orgTeamOnly validates that the request was made by an user of the organization team (or organization)
func (s *Service) orgTeamOnly(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		ctx := r.Context()
		UserID := ctx.Value(contextKeyUserID).(string)
		UserType := ctx.Value(contextKeyUserType).(string)
		OrgID := vars["orgId"]
		idErr := validate.Var(OrgID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		TeamID := vars["teamId"]
		idErr = validate.Var(TeamID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		var OrgRole string
		var TeamRole string
		if UserType != adminUserType {
			var UserErr error
			OrgRole, TeamRole, UserErr = s.OrganizationDataSvc.OrganizationTeamUserRole(ctx, UserID, OrgID, TeamID)
			if UserErr != nil {
				s.Failure(w, r, http.StatusForbidden, Errorf(EUNAUTHORIZED, "REQUIRES_TEAM_USER"))
				return
			}
		} else {
			OrgRole = adminUserType
			TeamRole = adminUserType
		}

		ctx = context.WithValue(ctx, contextKeyOrgRole, OrgRole)
		ctx = context.WithValue(ctx, contextKeyTeamRole, TeamRole)

		h(w, r.WithContext(ctx))
	}
}

// orgTeamAdminOnly validates that the request was made by an ADMIN of the organization team (or organization)
func (s *Service) orgTeamAdminOnly(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		ctx := r.Context()
		UserID := ctx.Value(contextKeyUserID).(string)
		UserType := ctx.Value(contextKeyUserType).(string)
		OrgID := vars["orgId"]
		idErr := validate.Var(OrgID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		TeamID := vars["teamId"]
		idErr = validate.Var(TeamID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		var OrgRole string
		var TeamRole string
		if UserType != adminUserType {
			var UserErr error
			OrgRole, TeamRole, UserErr := s.OrganizationDataSvc.OrganizationTeamUserRole(ctx, UserID, OrgID, TeamID)
			if UserErr != nil {
				s.Failure(w, r, http.StatusForbidden, Errorf(EUNAUTHORIZED, "REQUIRES_TEAM_USER"))
				return
			}
			if TeamRole != adminUserType && OrgRole != adminUserType {
				s.Failure(w, r, http.StatusForbidden, Errorf(EUNAUTHORIZED, "REQUIRES_TEAM_OR_ORGANIZATION_ADMIN"))
				return
			}
		} else {
			OrgRole = adminUserType
			TeamRole = adminUserType
		}

		ctx = context.WithValue(ctx, contextKeyOrgRole, OrgRole)
		ctx = context.WithValue(ctx, contextKeyTeamRole, TeamRole)

		h(w, r.WithContext(ctx))
	}
}

// departmentUserOnly validates that the request was made by a valid user of the organization (with department role)
func (s *Service) departmentUserOnly(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		ctx := r.Context()
		UserID := ctx.Value(contextKeyUserID).(string)
		UserType := ctx.Value(contextKeyUserType).(string)
		OrgID := vars["orgId"]
		idErr := validate.Var(OrgID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		DepartmentID := vars["departmentId"]
		idErr = validate.Var(DepartmentID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		var OrgRole string
		var DepartmentRole string
		if UserType != adminUserType {
			var UserErr error
			OrgRole, DepartmentRole, UserErr = s.OrganizationDataSvc.DepartmentUserRole(ctx, UserID, OrgID, DepartmentID)
			if UserErr != nil {
				s.Failure(w, r, http.StatusForbidden, Errorf(EUNAUTHORIZED, "REQUIRES_DEPARTMENT_USER"))
				return
			}
		} else {
			OrgRole = adminUserType
			DepartmentRole = adminUserType
		}

		ctx = context.WithValue(ctx, contextKeyOrgRole, OrgRole)
		ctx = context.WithValue(ctx, contextKeyDepartmentRole, DepartmentRole)

		h(w, r.WithContext(ctx))
	}
}

// departmentAdminOnly validates that the request was made by an ADMIN of the organization (with department role)
func (s *Service) departmentAdminOnly(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		ctx := r.Context()
		UserID := ctx.Value(contextKeyUserID).(string)
		UserType := ctx.Value(contextKeyUserType).(string)
		OrgID := vars["orgId"]
		idErr := validate.Var(OrgID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		DepartmentID := vars["departmentId"]
		idErr = validate.Var(DepartmentID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		var OrgRole string
		var DepartmentRole string
		if UserType != adminUserType {
			var UserErr error
			OrgRole, DepartmentRole, UserErr := s.OrganizationDataSvc.DepartmentUserRole(ctx, UserID, OrgID, DepartmentID)
			if UserErr != nil {
				s.Failure(w, r, http.StatusForbidden, Errorf(EUNAUTHORIZED, "REQUIRES_DEPARTMENT_USER"))
				return
			}
			if DepartmentRole != adminUserType && OrgRole != adminUserType {
				s.Failure(w, r, http.StatusForbidden, Errorf(EUNAUTHORIZED, "REQUIRES_DEPARTMENT_OR_ORGANIZATION_ADMIN"))
				return
			}
		} else {
			OrgRole = adminUserType
			DepartmentRole = adminUserType
		}

		ctx = context.WithValue(ctx, contextKeyOrgRole, OrgRole)
		ctx = context.WithValue(ctx, contextKeyDepartmentRole, DepartmentRole)

		h(w, r.WithContext(ctx))
	}
}

// departmentTeamUserOnly validates that the request was made by an user of the department team (or organization)
func (s *Service) departmentTeamUserOnly(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		ctx := r.Context()
		UserID := ctx.Value(contextKeyUserID).(string)
		UserType := ctx.Value(contextKeyUserType).(string)
		OrgID := vars["orgId"]
		idErr := validate.Var(OrgID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		DepartmentID := vars["departmentId"]
		idErr = validate.Var(DepartmentID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		TeamID := vars["teamId"]
		idErr = validate.Var(TeamID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		var OrgRole string
		var DepartmentRole string
		var TeamRole string
		if UserType != adminUserType {
			var UserErr error
			OrgRole, DepartmentRole, TeamRole, UserErr = s.OrganizationDataSvc.DepartmentTeamUserRole(ctx, UserID, OrgID, DepartmentID, TeamID)
			if UserErr != nil {
				s.Failure(w, r, http.StatusForbidden, Errorf(EUNAUTHORIZED, "REQUIRES_TEAM_USER"))
				return
			}
		} else {
			OrgRole = adminUserType
			DepartmentRole = adminUserType
			TeamRole = adminUserType
		}

		ctx = context.WithValue(ctx, contextKeyOrgRole, OrgRole)
		ctx = context.WithValue(ctx, contextKeyDepartmentRole, DepartmentRole)
		ctx = context.WithValue(ctx, contextKeyTeamRole, TeamRole)

		h(w, r.WithContext(ctx))
	}
}

// departmentTeamAdminOnly validates that the request was made by an ADMIN of the department team (or organization)
func (s *Service) departmentTeamAdminOnly(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		ctx := r.Context()
		UserID := ctx.Value(contextKeyUserID).(string)
		UserType := ctx.Value(contextKeyUserType).(string)
		OrgID := vars["orgId"]
		idErr := validate.Var(OrgID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		DepartmentID := vars["departmentId"]
		idErr = validate.Var(DepartmentID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		TeamID := vars["teamId"]
		idErr = validate.Var(TeamID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		var OrgRole string
		var DepartmentRole string
		var TeamRole string
		if UserType != adminUserType {
			var UserErr error
			OrgRole, DepartmentRole, TeamRole, UserErr = s.OrganizationDataSvc.DepartmentTeamUserRole(ctx, UserID, OrgID, DepartmentID, TeamID)
			if UserErr != nil {
				s.Failure(w, r, http.StatusForbidden, Errorf(EUNAUTHORIZED, "REQUIRES_TEAM_USER"))
				return
			}

			if TeamRole != adminUserType && DepartmentRole != adminUserType && OrgRole != adminUserType {
				s.Failure(w, r, http.StatusForbidden, Errorf(EUNAUTHORIZED, "REQUIRES_TEAM_OR_DEPARTMENT_OR_ORGANIZATION_ADMIN"))
				return
			}
		} else {
			OrgRole = adminUserType
			DepartmentRole = adminUserType
			TeamRole = adminUserType
		}

		ctx = context.WithValue(ctx, contextKeyOrgRole, OrgRole)
		ctx = context.WithValue(ctx, contextKeyDepartmentRole, DepartmentRole)
		ctx = context.WithValue(ctx, contextKeyTeamRole, TeamRole)

		h(w, r.WithContext(ctx))
	}
}

// teamUserOnly validates that the request was made by a valid user of the team
func (s *Service) teamUserOnly(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		ctx := r.Context()
		UserID := ctx.Value(contextKeyUserID).(string)
		UserType := ctx.Value(contextKeyUserType).(string)
		TeamID := vars["teamId"]
		idErr := validate.Var(TeamID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		var Role string
		if UserType != adminUserType {
			var UserErr error
			Role, UserErr = s.TeamDataSvc.TeamUserRole(ctx, UserID, TeamID)
			if UserType != adminUserType && UserErr != nil {
				s.Failure(w, r, http.StatusForbidden, Errorf(EUNAUTHORIZED, "REQUIRES_TEAM_USER"))
				return
			}
		} else {
			Role = adminUserType
		}

		ctx = context.WithValue(ctx, contextKeyTeamRole, Role)

		h(w, r.WithContext(ctx))
	}
}

// teamAdminOnly validates that the request was made by an ADMIN of the team
func (s *Service) teamAdminOnly(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		ctx := r.Context()
		UserID := ctx.Value(contextKeyUserID).(string)
		UserType := ctx.Value(contextKeyUserType).(string)
		TeamID := vars["teamId"]
		idErr := validate.Var(TeamID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		var Role string
		if UserType != adminUserType {
			var UserErr error
			Role, UserErr = s.TeamDataSvc.TeamUserRole(ctx, UserID, TeamID)
			if UserErr != nil {
				s.Failure(w, r, http.StatusForbidden, Errorf(EUNAUTHORIZED, "REQUIRES_TEAM_USER"))
				return
			}
			if Role != adminUserType {
				s.Failure(w, r, http.StatusForbidden, Errorf(EUNAUTHORIZED, "REQUIRES_TEAM_ADMIN"))
				return
			}
		} else {
			Role = adminUserType
		}

		ctx = context.WithValue(ctx, contextKeyTeamRole, Role)

		h(w, r.WithContext(ctx))
	}
}
