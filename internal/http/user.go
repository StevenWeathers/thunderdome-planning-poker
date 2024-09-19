package http

import (
	"bytes"
	"encoding/json"
	"image"
	"image/png"
	"io"
	"net/http"
	"strconv"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"

	"go.uber.org/zap"

	"github.com/anthonynsimon/bild/transform"
	"github.com/ipsn/go-adorable"
	"github.com/o1egl/govatar"

	"github.com/gorilla/mux"
)

// handleSessionUserProfile returns the users profile by session user ID
// @Summary      Get Session User Profile
// @Description  Gets a users profile by session user ID
// @Tags         auth, user
// @Produce      json
// @Success      200  object  standardJsonResponse{data=thunderdome.User}
// @Failure      403  object  standardJsonResponse{}
// @Failure      500  object  standardJsonResponse{}
// @Security     ApiKeyAuth
// @Router       /auth/user [get]
func (s *Service) handleSessionUserProfile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		SessionUserID := ctx.Value(contextKeyUserID).(string)

		User, UserErr := s.UserDataSvc.GetUser(ctx, SessionUserID)
		if UserErr != nil {
			s.Logger.Ctx(ctx).Error("handleSessionUserProfile error", zap.Error(UserErr),
				zap.String("session_user_id", SessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, UserErr)
			return
		}

		s.Success(w, r, http.StatusOK, User, nil)
	}
}

// handleUserProfile returns the users profile if it matches their session
// @Summary      Get User Profile
// @Description  Gets a users profile
// @Tags         user
// @Produce      json
// @Param        userId  path    string  true  "the user ID"
// @Success      200     object  standardJsonResponse{data=thunderdome.User}
// @Failure      403     object  standardJsonResponse{}
// @Failure      500     object  standardJsonResponse{}
// @Security     ApiKeyAuth
// @Router       /users/{userId} [get]
func (s *Service) handleUserProfile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		SessionUserID := ctx.Value(contextKeyUserID).(string)
		vars := mux.Vars(r)
		UserID := vars["userId"]

		User, UserErr := s.UserDataSvc.GetUser(ctx, UserID)
		if UserErr != nil {
			s.Logger.Ctx(ctx).Error("handleUserProfile error", zap.Error(UserErr),
				zap.String("entity_user_id", UserID), zap.String("session_user_id", SessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, UserErr)
			return
		}

		s.Success(w, r, http.StatusOK, User, nil)
	}
}

type userprofileUpdateRequestBody struct {
	Name                 string `json:"name" validate:"required,max=64"`
	Avatar               string `json:"avatar" validate:"max=128"`
	NotificationsEnabled bool   `json:"notificationsEnabled"`
	Country              string `json:"country" validate:"omitempty,len=2"`
	Locale               string `json:"locale" validate:"omitempty,len=2"`
	Company              string `json:"company" validate:"max=256"`
	JobTitle             string `json:"jobTitle" validate:"max=128"`
	Email                string `json:"email" validate:"omitempty,email"`
	Theme                string `json:"theme" validate:"max=5"`
}

// handleUserProfileUpdate attempts to update users profile
// @Summary      Update User Profile
// @Description  Update a users profile
// @Tags         user
// @Produce      json
// @Param        userId  path    string                        true  "the user ID"
// @Param        user    body    userprofileUpdateRequestBody  true  "the user profile object to update"
// @Success      200     object  standardJsonResponse{data=thunderdome.User}
// @Failure      403     object  standardJsonResponse{}
// @Failure      500     object  standardJsonResponse{}
// @Security     ApiKeyAuth
// @Router       /users/{userId} [put]
func (s *Service) handleUserProfileUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		SessionUserID := ctx.Value(contextKeyUserID).(string)
		SessionUserType := ctx.Value(contextKeyUserType).(string)
		vars := mux.Vars(r)
		UserID := vars["userId"]

		var profile = userprofileUpdateRequestBody{}
		body, bodyErr := io.ReadAll(r.Body)
		if bodyErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
			return
		}

		jsonErr := json.Unmarshal(body, &profile)
		if jsonErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		inputErr := validate.Struct(profile)
		if inputErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, inputErr.Error()))
			return
		}

		if SessionUserType == thunderdome.AdminUserType {
			_, _, vErr := validateUserAccount(profile.Name, profile.Email)
			if vErr != nil {
				s.Failure(w, r, http.StatusBadRequest, vErr)
				return
			}
			updateErr := s.UserDataSvc.UpdateUserAccount(ctx, UserID, profile.Name, profile.Email, profile.Avatar, profile.NotificationsEnabled, profile.Country, profile.Locale, profile.Company, profile.JobTitle, profile.Theme)
			if updateErr != nil {
				s.Logger.Ctx(ctx).Error("handleUserProfileUpdate error", zap.Error(updateErr),
					zap.String("entity_user_id", UserID), zap.String("session_user_id", SessionUserID))
				s.Failure(w, r, http.StatusInternalServerError, updateErr)
				return
			}
		} else {
			var updateErr error
			if !s.Config.LdapEnabled {
				if profile.Name == "" {
					s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, "INVALID_USERNAME"))
					return
				}
				updateErr = s.UserDataSvc.UpdateUserProfile(ctx, UserID, profile.Name, profile.Avatar, profile.NotificationsEnabled, profile.Country, profile.Locale, profile.Company, profile.JobTitle, profile.Theme)
			} else {
				updateErr = s.UserDataSvc.UpdateUserProfileLdap(ctx, UserID, profile.Avatar, profile.NotificationsEnabled, profile.Country, profile.Locale, profile.Company, profile.JobTitle, profile.Theme)
			}
			if updateErr != nil {
				s.Logger.Ctx(ctx).Error("handleUserProfileUpdate error", zap.Error(updateErr),
					zap.String("entity_user_id", UserID), zap.String("session_user_id", SessionUserID))
				s.Failure(w, r, http.StatusInternalServerError, updateErr)
				return
			}
		}

		user, UserErr := s.UserDataSvc.GetUser(ctx, UserID)
		if UserErr != nil {
			s.Logger.Ctx(ctx).Error("handleUserProfileUpdate error", zap.Error(UserErr),
				zap.String("entity_user_id", UserID), zap.String("session_user_id", SessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, UserErr)
			return
		}

		s.Success(w, r, http.StatusOK, user, nil)
	}
}

// handleUserDelete attempts to delete a users account
// @Summary      Delete User
// @Description  Deletes a user
// @Tags         user
// @Produce      json
// @Param        userId  path    string  true  "the user ID"
// @Success      200     object  standardJsonResponse{}
// @Failure      403     object  standardJsonResponse{}
// @Failure      500     object  standardJsonResponse{}
// @Security     ApiKeyAuth
// @Router       /users/{userId} [delete]
func (s *Service) handleUserDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		UserID := vars["userId"]
		ctx := r.Context()
		SessionUserID := ctx.Value(contextKeyUserID).(string)

		User, UserErr := s.UserDataSvc.GetUser(ctx, UserID)
		if UserErr != nil {
			s.Logger.Ctx(ctx).Error("handleUserDelete error", zap.Error(UserErr),
				zap.String("entity_user_id", UserID), zap.String("session_user_id", SessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, UserErr)
			return
		}

		updateErr := s.UserDataSvc.DeleteUser(ctx, UserID)
		if updateErr != nil {
			s.Logger.Ctx(ctx).Error("handleUserDelete error", zap.Error(updateErr),
				zap.String("user_id", UserID), zap.String("session_user_id", SessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, updateErr)
			return
		}

		// don't attempt to send email to guest users
		if User.Email != "" {
			_ = s.Email.SendDeleteConfirmation(User.Name, User.Email)
		}

		// don't clear admins user cookies when deleting other users
		if UserID == SessionUserID {
			s.Cookie.ClearUserCookies(w)
		}

		s.Success(w, r, http.StatusOK, nil, nil)
	}
}

// handleVerifyRequest sends verification Email
// @Summary      Request Verification Email
// @Description  Sends verification Email
// @Tags         user
// @Param        userId  path    string  true  "the user ID"
// @Success      200     object  standardJsonResponse{}
// @Success      400     object  standardJsonResponse{}
// @Success      500     object  standardJsonResponse{}
// @Router       /users/{userId}/request-verify [post]
func (s *Service) handleVerifyRequest() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		SessionUserID, _ := ctx.Value(contextKeyUserID).(*string)
		vars := mux.Vars(r)
		UserID := vars["userId"]

		User, VerifyId, err := s.AuthDataSvc.UserVerifyRequest(ctx, UserID)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleVerifyRequest error", zap.Error(err),
				zap.String("entity_user_id", UserID), zap.Stringp("session_user_id", SessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		_ = s.Email.SendEmailVerification(User.Name, User.Email, VerifyId)

		s.Success(w, r, http.StatusOK, nil, nil)
	}
}

// handleGetActiveCountries gets a list of registered users countries
// @Summary      Get Active Countries
// @Description  Gets a list of users countries
// @Produce      json
// @Success      200  object  standardJsonResponse{[]string}
// @Failure      500  object  standardJsonResponse{}
// @Router       /active-countries [get]
func (s *Service) handleGetActiveCountries() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		SessionUserID, _ := ctx.Value(contextKeyUserID).(*string)
		countries, err := s.UserDataSvc.GetActiveCountries(ctx)

		if err != nil {
			s.Logger.Ctx(ctx).Error("handleGetActiveCountries error", zap.Error(err),
				zap.Stringp("session_user_id", SessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		w.Header().Set("Cache-Control", "max-age=3600") // cache for 1 hour just to decrease load
		s.Success(w, r, http.StatusOK, countries, nil)
	}
}

// handleUserAvatar creates an avatar for the given user by ID
func (s *Service) handleUserAvatar() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		ctx := r.Context()
		SessionUserID, _ := ctx.Value(contextKeyUserID).(*string)

		Width, _ := strconv.Atoi(vars["width"])
		UserID := vars["id"]
		AvatarGender := govatar.MALE
		userGender, ok := vars["avatar"]
		if ok {
			if userGender == "female" {
				AvatarGender = govatar.FEMALE
			}
		}

		var avatar image.Image
		if s.Config.AvatarService == "govatar" {
			avatar, _ = govatar.GenerateForUsername(AvatarGender, UserID)
		} else { // must be goadorable
			var err error
			avatar, _, err = image.Decode(bytes.NewReader(adorable.PseudoRandom([]byte(UserID))))
			if err != nil {
				s.Logger.Ctx(ctx).Error(err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		img := transform.Resize(avatar, Width, Width, transform.Linear)
		buffer := new(bytes.Buffer)

		if err := png.Encode(buffer, img); err != nil {
			s.Logger.Ctx(ctx).Error("handleUserAvatar error", zap.Error(err), zap.String("entity_user_id", UserID),
				zap.Stringp("session_user_id", SessionUserID))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "image/png")
		w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))

		if _, err := w.Write(buffer.Bytes()); err != nil {
			s.Logger.Ctx(ctx).Error("handleUserAvatar error", zap.Error(err), zap.String("entity_user_id", UserID),
				zap.Stringp("session_user_id", SessionUserID))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

// handleUserOrganizationInvite processes an organization invite for the user
// @Summary      User Organization Invite
// @Description  Processes an organization invite for the user
// @Tags         user
// @Param        userId  path    string  true  "the user ID"
// @Param        inviteId  path    string  true  "the invite ID"
// @Success      200     object  standardJsonResponse{}
// @Success      400     object  standardJsonResponse{}
// @Success      500     object  standardJsonResponse{}
// @Router       /users/{userId}/invite/organization/{inviteId} [post]
func (s *Service) handleUserOrganizationInvite() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		SessionUserID, _ := ctx.Value(contextKeyUserID).(*string)
		vars := mux.Vars(r)
		UserID := vars["userId"]
		InviteID := vars["inviteId"]

		user, err := s.UserDataSvc.GetUser(ctx, UserID)
		if err != nil {
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		orgInvite, err := s.OrganizationDataSvc.OrganizationUserGetInviteByID(ctx, InviteID)
		if err != nil {
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}
		if user.Email != orgInvite.Email {
			s.Failure(w, r, http.StatusInternalServerError, Errorf(EINVALID, err.Error()))
			return
		}

		orgId, inviteErr := s.OrganizationDataSvc.OrganizationAddUser(ctx, orgInvite.OrganizationId, UserID, orgInvite.Role)
		if inviteErr != nil {
			s.Logger.Ctx(ctx).Error("handleUserOrganizationInvite error adding invited user to organization",
				zap.Error(inviteErr),
				zap.String("session_user_id", *SessionUserID),
				zap.String("invite_id", orgInvite.InviteId))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		delInviteErr := s.OrganizationDataSvc.OrganizationDeleteUserInvite(ctx, orgInvite.InviteId)
		if delInviteErr != nil {
			s.Logger.Ctx(ctx).Error("handleUserOrganizationInvite error deleting user invite to organization",
				zap.Error(delInviteErr),
				zap.String("session_user_id", *SessionUserID),
				zap.String("invite_id", orgInvite.InviteId))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		organization, orgErr := s.OrganizationDataSvc.OrganizationGet(ctx, orgId)
		if orgErr != nil {
			s.Logger.Ctx(ctx).Error("handleUserOrganizationInvite error getting organization",
				zap.Error(delInviteErr),
				zap.String("session_user_id", *SessionUserID),
				zap.String("invite_id", orgInvite.InviteId))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		result := thunderdome.UserOrganization{
			Organization: *organization,
			Role:         orgInvite.Role,
		}

		s.Success(w, r, http.StatusOK, result, nil)
	}
}

// handleUserDepartmentInvite processes an department invite for the user
// @Summary      User Department Invite
// @Description  Processes an department invite for the user
// @Tags         user
// @Param        userId  path    string  true  "the user ID"
// @Param        inviteId  path    string  true  "the invite ID"
// @Success      200     object  standardJsonResponse{}
// @Success      400     object  standardJsonResponse{}
// @Success      500     object  standardJsonResponse{}
// @Router       /users/{userId}/invite/department/{inviteId} [post]
func (s *Service) handleUserDepartmentInvite() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		SessionUserID, _ := ctx.Value(contextKeyUserID).(*string)
		vars := mux.Vars(r)
		UserID := vars["userId"]
		InviteID := vars["inviteId"]

		user, err := s.UserDataSvc.GetUser(ctx, UserID)
		if err != nil {
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		deptInvite, err := s.OrganizationDataSvc.DepartmentUserGetInviteByID(ctx, InviteID)
		if err != nil {
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}
		if user.Email != deptInvite.Email {
			s.Failure(w, r, http.StatusInternalServerError, Errorf(EINVALID, err.Error()))
			return
		}

		dept, deptErr := s.OrganizationDataSvc.DepartmentGet(ctx, deptInvite.DepartmentId)
		if deptErr != nil {
			s.Logger.Ctx(ctx).Error("handleUserDepartmentInvite get department error", zap.Error(deptErr),
				zap.String("department_id", dept.Id),
				zap.String("session_user_id", *SessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, deptErr)
			return
		}

		org, orgErr := s.OrganizationDataSvc.OrganizationGet(ctx, dept.OrganizationId)
		if orgErr != nil {
			s.Logger.Ctx(ctx).Error("handleUserDepartmentInvite get organization error", zap.Error(orgErr),
				zap.String("organization_id", org.Id),
				zap.String("department_id", dept.Id),
				zap.String("session_user_id", *SessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, orgErr)
			return
		}

		_, orgAddErr := s.OrganizationDataSvc.OrganizationUpsertUser(ctx, org.Id, UserID, "MEMBER")
		if orgAddErr != nil {
			s.Logger.Ctx(ctx).Error(
				"handleUserDepartmentInvite upsert organization user error", zap.Error(orgAddErr),
				zap.String("session_user_id", *SessionUserID),
				zap.String("organization_id", org.Id),
				zap.String("department_id", dept.Id),
				zap.String("user_id", UserID), zap.String("user_role", "MEMBER"))
			s.Failure(w, r, http.StatusInternalServerError, orgAddErr)
			return
		}

		_, inviteErr := s.OrganizationDataSvc.DepartmentAddUser(ctx, deptInvite.DepartmentId, UserID, deptInvite.Role)
		if inviteErr != nil {
			s.Logger.Ctx(ctx).Error("handleUserDepartmentInvite error adding invited user to organization",
				zap.Error(inviteErr),
				zap.String("session_user_id", *SessionUserID),
				zap.String("invite_id", deptInvite.InviteId))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		delInviteErr := s.OrganizationDataSvc.OrganizationDeleteUserInvite(ctx, deptInvite.InviteId)
		if delInviteErr != nil {
			s.Logger.Ctx(ctx).Error("handleUserDepartmentInvite error deleting user invite to organization",
				zap.Error(delInviteErr),
				zap.String("session_user_id", *SessionUserID),
				zap.String("invite_id", deptInvite.InviteId))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		result := thunderdome.UserDepartment{
			Department: *dept,
			Role:       deptInvite.Role,
		}

		s.Success(w, r, http.StatusOK, result, nil)
	}
}

// handleUserTeamInvite processes a team invite for the user
// @Summary      User Team Invite
// @Description  Processes a team invite for the user
// @Tags         user
// @Param        userId  path    string  true  "the user ID"
// @Param        inviteId  path    string  true  "the invite ID"
// @Success      200     object  standardJsonResponse{}
// @Success      400     object  standardJsonResponse{}
// @Success      500     object  standardJsonResponse{}
// @Router       /users/{userId}/invite/team/{inviteId} [post]
func (s *Service) handleUserTeamInvite() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		SessionUserID, _ := ctx.Value(contextKeyUserID).(*string)
		vars := mux.Vars(r)
		UserID := vars["userId"]
		InviteID := vars["inviteId"]

		user, err := s.UserDataSvc.GetUser(ctx, UserID)
		if err != nil {
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		teamInvite, err := s.TeamDataSvc.TeamUserGetInviteByID(ctx, InviteID)
		if err != nil {
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}
		if user.Email != teamInvite.Email {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, err.Error()))
			return
		}

		team, teamErr := s.TeamDataSvc.TeamGet(ctx, teamInvite.TeamId)
		if teamErr != nil {
			s.Logger.Ctx(ctx).Error("handleTeamInviteUser error", zap.Error(teamErr),
				zap.String("team_id", teamInvite.TeamId), zap.String("session_user_id", *SessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, teamErr)
			return
		}

		// team is associated to organization, upsert user to organization
		if team.OrganizationId != "" {
			_, orgAddErr := s.OrganizationDataSvc.OrganizationUpsertUser(ctx, team.OrganizationId, UserID, "MEMBER")
			if orgAddErr != nil {
				s.Logger.Ctx(ctx).Error(
					"handleTeamInviteUser upsert organization user error", zap.Error(orgAddErr),
					zap.String("session_user_id", *SessionUserID),
					zap.String("organization_id", team.OrganizationId),
					zap.String("user_id", UserID), zap.String("user_role", "MEMBER"))
				s.Failure(w, r, http.StatusInternalServerError, orgAddErr)
				return
			}
		}

		// team is associated to department, upsert user to department and organization
		if team.DepartmentId != "" {
			dept, deptErr := s.OrganizationDataSvc.DepartmentGet(ctx, team.DepartmentId)
			if deptErr != nil {
				s.Logger.Ctx(ctx).Error("handleTeamInviteUser get department error", zap.Error(deptErr),
					zap.String("department_id", dept.Id),
					zap.String("session_user_id", *SessionUserID))
				s.Failure(w, r, http.StatusInternalServerError, deptErr)
				return
			}

			org, orgErr := s.OrganizationDataSvc.OrganizationGet(ctx, dept.OrganizationId)
			if orgErr != nil {
				s.Logger.Ctx(ctx).Error("handleTeamInviteUser get organization error", zap.Error(orgErr),
					zap.String("organization_id", org.Id),
					zap.String("department_id", dept.Id),
					zap.String("session_user_id", *SessionUserID))
				s.Failure(w, r, http.StatusInternalServerError, orgErr)
				return
			}

			team.OrganizationId = org.Id

			_, orgAddErr := s.OrganizationDataSvc.OrganizationUpsertUser(ctx, org.Id, UserID, "MEMBER")
			if orgAddErr != nil {
				s.Logger.Ctx(ctx).Error(
					"handleTeamInviteUser upsert organization user error", zap.Error(orgAddErr),
					zap.String("session_user_id", *SessionUserID),
					zap.String("organization_id", org.Id),
					zap.String("department_id", dept.Id),
					zap.String("user_id", UserID), zap.String("user_role", "MEMBER"))
				s.Failure(w, r, http.StatusInternalServerError, orgAddErr)
				return
			}

			_, deptAddErr := s.OrganizationDataSvc.DepartmentUpsertUser(ctx, team.DepartmentId, UserID, "MEMBER")
			if deptAddErr != nil {
				s.Logger.Ctx(ctx).Error(
					"handleTeamInviteUser upsert department user error", zap.Error(deptAddErr),
					zap.String("session_user_id", *SessionUserID),
					zap.String("department_id", team.DepartmentId),
					zap.String("user_id", UserID), zap.String("user_role", "MEMBER"))
				s.Failure(w, r, http.StatusInternalServerError, deptAddErr)
				return
			}
		}

		_, inviteErr := s.TeamDataSvc.TeamAddUser(ctx, teamInvite.TeamId, UserID, teamInvite.Role)
		if inviteErr != nil {
			s.Logger.Ctx(ctx).Error("handleUserRegistration error adding invited user to team", zap.Error(inviteErr),
				zap.String("session_user_id", *SessionUserID),
				zap.String("invite_id", teamInvite.InviteId))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		delInviteErr := s.TeamDataSvc.TeamDeleteUserInvite(ctx, teamInvite.InviteId)
		if delInviteErr != nil {
			s.Logger.Ctx(ctx).Error("handleUserRegistration error deleting user invite to team", zap.Error(delInviteErr),
				zap.String("session_user_id", *SessionUserID),
				zap.String("invite_id", teamInvite.InviteId))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		result := thunderdome.UserTeam{
			Team: *team,
			Role: teamInvite.Role,
		}

		s.Success(w, r, http.StatusOK, result, nil)
	}
}

// handleUserCredential returns the users credential if they have one
// @Summary      Get User Credential
// @Description  Gets a users credential
// @Tags         user
// @Produce      json
// @Param        userId  path    string  true  "the user ID"
// @Success      200     object  standardJsonResponse{data=thunderdome.Credential}
// @Failure      403     object  standardJsonResponse{}
// @Failure      500     object  standardJsonResponse{}
// @Security     ApiKeyAuth
// @Router       /users/{userId}/credential [get]
func (s *Service) handleUserCredential() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		SessionUserID := ctx.Value(contextKeyUserID).(string)
		vars := mux.Vars(r)
		UserID := vars["userId"]

		credential, UserErr := s.UserDataSvc.GetUserCredential(ctx, UserID)
		if UserErr != nil {
			s.Logger.Ctx(ctx).Error("handleUserCredential error", zap.Error(UserErr),
				zap.String("entity_user_id", UserID),
				zap.String("session_user_id", SessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, UserErr)
			return
		}

		s.Success(w, r, http.StatusOK, credential, nil)
	}
}
