package http

import (
	"encoding/json"
	"io"
	"net/http"

	"go.uber.org/zap"

	"github.com/gorilla/mux"
)

// handleAppStats gets the applications stats
// @Summary      Get Application Stats
// @Description  Get application stats such as count of registered users
// @Tags         admin
// @Produce      json
// @Success      200  object  standardJsonResponse{data=[]thunderdome.ApplicationStats}
// @Failure      500  object  standardJsonResponse{}
// @Security     ApiKeyAuth
// @Router       /admin/stats [get]
func (s *Service) handleAppStats() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		SessionUserID := ctx.Value(contextKeyUserID).(string)
		AppStats, err := s.AdminDataSvc.GetAppStats(r.Context())
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleAppStats error", zap.Error(err), zap.String("session_user_id", SessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, AppStats, nil)
	}
}

// handleGetRegisteredUsers gets a list of registered users
// @Summary      Get Registered Users
// @Description  Get list of registered users
// @Tags         admin
// @Produce      json
// @Param        limit   query   int  false  "Max number of results to return"
// @Param        offset  query   int  false  "Starting point to return rows from, should be multiplied by limit or 0"
// @Success      200     object  standardJsonResponse{data=[]thunderdome.User}
// @Failure      500     object  standardJsonResponse{}
// @Security     ApiKeyAuth
// @Router       /admin/users [get]
func (s *Service) handleGetRegisteredUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		SessionUserID := ctx.Value(contextKeyUserID).(string)
		Limit, Offset := getLimitOffsetFromRequest(r)

		Users, Count, err := s.UserDataSvc.GetRegisteredUsers(ctx, Limit, Offset)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleGetRegisteredUsers error", zap.Error(err),
				zap.Int("limit", Limit), zap.Int("offset", Offset), zap.String("session_user_id", SessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		Meta := &pagination{
			Count:  Count,
			Offset: Offset,
			Limit:  Limit,
		}

		s.Success(w, r, http.StatusOK, Users, Meta)
	}
}

type userCreateRequestBody struct {
	Name      string `json:"name" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password1 string `json:"password1" validate:"required,min=6,max=72"`
	Password2 string `json:"password2" validate:"required,min=6,max=72,eqfield=Password1"`
}

// handleUserCreate registers a new authenticated user
// @Summary      Create Registered User
// @Description  Create a registered user
// @Tags         admin
// @Produce      json
// @param        newUser  body    userCreateRequestBody  true  "new user object"
// @Success      200      object  standardJsonResponse{data=thunderdome.User}
// @Failure      400      object  standardJsonResponse{}
// @Failure      500      object  standardJsonResponse{}
// @Security     ApiKeyAuth
// @Router       /admin/users [post]
func (s *Service) handleUserCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		SessionUserID := ctx.Value(contextKeyUserID).(string)
		var user = userCreateRequestBody{}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, err.Error()))
			return
		}

		jsonErr := json.Unmarshal(body, &user)
		if jsonErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		accountErr := validate.Struct(user)

		if accountErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, accountErr.Error()))
			return
		}

		newUser, VerifyID, err := s.UserDataSvc.CreateUser(ctx, user.Name, user.Email, user.Password1)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleUserCreate error", zap.Error(err),
				zap.String("user_email", user.Email), zap.String("session_user_id", SessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		_ = s.Email.SendWelcome(user.Name, user.Email, VerifyID)

		s.Success(w, r, http.StatusOK, newUser, nil)
	}
}

// handleUserPromote handles promoting a user to admin
// @Summary      Promotes User
// @Description  Promotes a user to admin
// @Description  Grants read and write access to administrative information
// @Tags         admin
// @Produce      json
// @Param        userId  path    string  true  "the user ID to promote"
// @Success      200     object  standardJsonResponse{}
// @Failure      500     object  standardJsonResponse{}
// @Security     ApiKeyAuth
// @Router       /admin/users/{userId}/promote/ [patch]
func (s *Service) handleUserPromote() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		vars := mux.Vars(r)
		UserID := vars["userId"]
		SessionUserID := ctx.Value(contextKeyUserID).(string)
		idErr := validate.Var(UserID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		err := s.UserDataSvc.PromoteUser(ctx, UserID)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleUserPromote error", zap.Error(err),
				zap.String("entity_user_id", UserID), zap.String("session_user_id", SessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, nil, nil)
	}
}

// handleUserDemote handles demoting a user to registered
// @Summary      Demote User
// @Description  Demotes a user from admin to registered
// @Tags         admin
// @Produce      json
// @Param        userId  path    string  true  "the user ID to demote"
// @Success      200     object  standardJsonResponse{}
// @Failure      500     object  standardJsonResponse{}
// @Security     ApiKeyAuth
// @Router       /admin/users/{userId}/demote [patch]
func (s *Service) handleUserDemote() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		vars := mux.Vars(r)
		UserID := vars["userId"]
		SessionUserID := ctx.Value(contextKeyUserID).(string)
		idErr := validate.Var(UserID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		err := s.UserDataSvc.DemoteUser(ctx, UserID)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleUserDemote error", zap.Error(err),
				zap.String("entity_user_id", UserID), zap.String("session_user_id", SessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, nil, nil)
	}
}

// handleUserDisable handles disabling a user
// @Summary      Disable User
// @Description  Disable a user from logging in
// @Tags         admin
// @Produce      json
// @Param        userId  path    string  true  "the user ID to disable"
// @Success      200     object  standardJsonResponse{}
// @Failure      500     object  standardJsonResponse{}
// @Security     ApiKeyAuth
// @Router       /admin/users/{userId}/disable [patch]
func (s *Service) handleUserDisable() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		vars := mux.Vars(r)
		UserID := vars["userId"]
		SessionUserID := ctx.Value(contextKeyUserID).(string)
		idErr := validate.Var(UserID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		err := s.UserDataSvc.DisableUser(ctx, UserID)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleUserDisable error", zap.Error(err),
				zap.String("entity_user_id", UserID), zap.String("session_user_id", SessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, nil, nil)
	}
}

// handleUserEnable handles enabling a user
// @Summary      Enable User
// @Description  Enable a user to allow login
// @Tags         admin
// @Produce      json
// @Param        userId  path    string  true  "the user ID to enable"
// @Success      200     object  standardJsonResponse{}
// @Failure      500     object  standardJsonResponse{}
// @Security     ApiKeyAuth
// @Router       /admin/users/{userId}/enable [patch]
func (s *Service) handleUserEnable() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		vars := mux.Vars(r)
		UserID := vars["userId"]
		SessionUserID := ctx.Value(contextKeyUserID).(string)
		idErr := validate.Var(UserID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		err := s.UserDataSvc.EnableUser(r.Context(), UserID)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleUserEnable error", zap.Error(err),
				zap.String("entity_user_id", UserID), zap.String("session_user_id", SessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, nil, nil)
	}
}

// handleAdminUpdateUserPassword attempts to update a user's password
// @Summary      Update Password
// @Description  Updates the user's password
// @Tags         admin
// @Param        userId     path    string                     true   "the user ID to update password for"
// @Param        passwords  body    updatePasswordRequestBody  false  "update password object"
// @Success      200        object  standardJsonResponse{}
// @Success      400        object  standardJsonResponse{}
// @Success      500        object  standardJsonResponse{}
// @Security     ApiKeyAuth
// @Router       /admin/users/{userId}/password [patch]
func (s *Service) handleAdminUpdateUserPassword() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		vars := mux.Vars(r)
		UserID := vars["userId"]
		SessionUserID := ctx.Value(contextKeyUserID).(string)
		idErr := validate.Var(UserID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		body, bodyErr := io.ReadAll(r.Body)
		if bodyErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
			return
		}

		var u = updatePasswordRequestBody{}
		jsonErr := json.Unmarshal(body, &u)
		if jsonErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		inputErr := validate.Struct(u)
		if inputErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, inputErr.Error()))
			return
		}

		UserName, UserEmail, updateErr := s.AuthDataSvc.UserUpdatePassword(ctx, UserID, u.Password1)
		if updateErr != nil {
			s.Logger.Ctx(ctx).Error("handleAdminUpdateUserPassword error", zap.Error(updateErr),
				zap.String("entity_user_id", UserID), zap.String("session_user_id", SessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, updateErr)
			return
		}

		_ = s.Email.SendPasswordUpdate(UserName, UserEmail)

		s.Success(w, r, http.StatusOK, nil, nil)
	}
}

// handleGetOrganizations gets a list of organizations
// @Summary      Get Organizations
// @Description  Get a list of organizations
// @Tags         admin
// @Produce      json
// @Param        limit   query   int  false  "Max number of results to return"
// @Param        offset  query   int  false  "Starting point to return rows from, should be multiplied by limit or 0"
// @Success      200     object  standardJsonResponse{data=[]thunderdome.Organization}
// @Failure      500     object  standardJsonResponse{}
// @Security     ApiKeyAuth
// @Router       /admin/organizations [get]
func (s *Service) handleGetOrganizations() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if !s.Config.OrganizationsEnabled {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, "ORGANIZATIONS_DISABLED"))
			return
		}
		Limit, Offset := getLimitOffsetFromRequest(r)

		Organizations := s.OrganizationDataSvc.OrganizationList(ctx, Limit, Offset)

		s.Success(w, r, http.StatusOK, Organizations, nil)
	}
}

// handleGetTeams gets a list of teams
// @Summary      Get Teams
// @Description  Get a list of teams
// @Tags         admin
// @Produce      json
// @Param        limit   query   int  false  "Max number of results to return"
// @Param        offset  query   int  false  "Starting point to return rows from, should be multiplied by limit or 0"
// @Success      200     object  standardJsonResponse{data=[]thunderdome.Team}
// @Failure      500     object  standardJsonResponse{}
// @Security     ApiKeyAuth
// @Router       /admin/teams [get]
func (s *Service) handleGetTeams() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		Limit, Offset := getLimitOffsetFromRequest(r)

		Teams, Count := s.TeamDataSvc.TeamList(ctx, Limit, Offset)

		Meta := &pagination{
			Count:  Count,
			Offset: Offset,
			Limit:  Limit,
		}

		s.Success(w, r, http.StatusOK, Teams, Meta)
	}
}

// handleGetAPIKeys gets a list of APIKeys
// @Summary      Get API Keys
// @Description  Get a list of users API Keys
// @Tags         admin
// @Produce      json
// @Param        limit   query   int  false  "Max number of results to return"
// @Param        offset  query   int  false  "Starting point to return rows from, should be multiplied by limit or 0"
// @Success      200     object  standardJsonResponse{data=[]thunderdome.UserAPIKey}
// @Failure      500     object  standardJsonResponse{}
// @Security     ApiKeyAuth
// @Router       /admin/apikeys [get]
func (s *Service) handleGetAPIKeys() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		Limit, Offset := getLimitOffsetFromRequest(r)

		Teams := s.ApiKeyDataSvc.GetAPIKeys(ctx, Limit, Offset)

		s.Success(w, r, http.StatusOK, Teams, nil)
	}
}

// handleSearchRegisteredUsersByEmail gets a list of registered users filtered by Email likeness
// @Summary      Search Registered Users by Email
// @Description  Get list of registered users filtered by Email likeness
// @Tags         admin
// @Produce      json
// @Param        search  query   string  true   "The user Email to search for"
// @Param        limit   query   int     false  "Max number of results to return"
// @Param        offset  query   int     false  "Starting point to return rows from, should be multiplied by limit or 0"
// @Success      200     object  standardJsonResponse{data=[]thunderdome.User}
// @Failure      400     object  standardJsonResponse{}
// @Failure      500     object  standardJsonResponse{}
// @Security     ApiKeyAuth
// @Router       /admin/search/users/email [get]
func (s *Service) handleSearchRegisteredUsersByEmail() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		SessionUserID := ctx.Value(contextKeyUserID).(string)
		Limit, Offset := getLimitOffsetFromRequest(r)
		Search, err := getSearchFromRequest(r)
		if err != nil {
			s.Failure(w, r, http.StatusBadRequest, err)
			return
		}

		Users, Count, err := s.UserDataSvc.SearchRegisteredUsersByEmail(ctx, Search, Limit, Offset)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleSearchRegisteredUsersByEmail error", zap.Error(err),
				zap.Int("limit", Limit), zap.Int("offset", Offset), zap.String("user_email", Search),
				zap.String("session_user_id", SessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		Meta := &pagination{
			Count:  Count,
			Offset: Offset,
			Limit:  Limit,
		}

		s.Success(w, r, http.StatusOK, Users, Meta)
	}
}
