package http

import (
	"bytes"
	"encoding/json"
	"image"
	"image/png"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"

	"go.uber.org/zap"

	"github.com/anthonynsimon/bild/transform"
	"github.com/ipsn/go-adorable"
	"github.com/o1egl/govatar"
)

// handleSessionUserProfile returns the users profile by session user ID
//
//	@Summary		Get Session User Profile
//	@Description	Gets a users profile by session user ID
//	@Tags			auth, user
//	@Produce		json
//	@Success		200	object	standardJsonResponse{data=thunderdome.User}
//	@Failure		403	object	standardJsonResponse{}
//	@Failure		500	object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/auth/user [get]
func (s *Service) handleSessionUserProfile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)

		user, userErr := s.UserDataSvc.GetUserByID(ctx, sessionUserID)
		if userErr != nil {
			s.Logger.Ctx(ctx).Error("handleSessionUserProfile error", zap.Error(userErr),
				zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, userErr)
			return
		}

		s.Success(w, r, http.StatusOK, user, nil)
	}
}

// handleUserProfile returns the users profile if it matches their session
//
//	@Summary		Get User Profile
//	@Description	Gets a users profile
//	@Tags			user
//	@Produce		json
//	@Param			userId	path	string	true	"the user ID"
//	@Success		200		object	standardJsonResponse{data=thunderdome.User}
//	@Failure		403		object	standardJsonResponse{}
//	@Failure		500		object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/users/{userId} [get]
func (s *Service) handleUserProfile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)

		userID := r.PathValue("userId")
		idErr := validate.Var(userID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		user, userErr := s.UserDataSvc.GetUserByID(ctx, userID)
		if userErr != nil {
			s.Logger.Ctx(ctx).Error("handleUserProfile error", zap.Error(userErr),
				zap.String("entity_user_id", userID), zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, userErr)
			return
		}

		s.Success(w, r, http.StatusOK, user, nil)
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
//
//	@Summary		Update User Profile
//	@Description	Update a users profile
//	@Tags			user
//	@Produce		json
//	@Param			userId	path	string							true	"the user ID"
//	@Param			user	body	userprofileUpdateRequestBody	true	"the user profile object to update"
//	@Success		200		object	standardJsonResponse{data=thunderdome.User}
//	@Failure		403		object	standardJsonResponse{}
//	@Failure		500		object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/users/{userId} [put]
func (s *Service) handleUserProfileUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)
		sessionUserType := ctx.Value(contextKeyUserType).(string)

		userID := r.PathValue("userId")
		idErr := validate.Var(userID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

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

		invalidUsername := containsLink(profile.Name)
		if invalidUsername {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, "INVALID_USERNAME"))
			return
		}

		if sessionUserType == thunderdome.AdminUserType {
			_, _, vErr := validateUserAccount(profile.Name, profile.Email)
			if vErr != nil {
				s.Failure(w, r, http.StatusBadRequest, vErr)
				return
			}
			updateErr := s.UserDataSvc.UpdateUserAccount(ctx, userID, profile.Name, profile.Email, profile.Avatar, profile.NotificationsEnabled, profile.Country, profile.Locale, profile.Company, profile.JobTitle, profile.Theme)
			if updateErr != nil {
				s.Logger.Ctx(ctx).Error("handleUserProfileUpdate error", zap.Error(updateErr),
					zap.String("entity_user_id", userID), zap.String("session_user_id", sessionUserID))
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
				updateErr = s.UserDataSvc.UpdateUserProfile(ctx, userID, profile.Name, profile.Avatar, profile.NotificationsEnabled, profile.Country, profile.Locale, profile.Company, profile.JobTitle, profile.Theme)
			} else {
				updateErr = s.UserDataSvc.UpdateUserProfileLdap(ctx, userID, profile.Avatar, profile.NotificationsEnabled, profile.Country, profile.Locale, profile.Company, profile.JobTitle, profile.Theme)
			}
			if updateErr != nil {
				s.Logger.Ctx(ctx).Error("handleUserProfileUpdate error", zap.Error(updateErr),
					zap.String("entity_user_id", userID), zap.String("session_user_id", sessionUserID))
				s.Failure(w, r, http.StatusInternalServerError, updateErr)
				return
			}
		}

		user, userErr := s.UserDataSvc.GetUserByID(ctx, userID)
		if userErr != nil {
			s.Logger.Ctx(ctx).Error("handleUserProfileUpdate error", zap.Error(userErr),
				zap.String("entity_user_id", userID), zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, userErr)
			return
		}

		s.Success(w, r, http.StatusOK, user, nil)
	}
}

// handleUserDelete attempts to delete a users account
//
//	@Summary		Delete User
//	@Description	Deletes a user
//	@Tags			user
//	@Produce		json
//	@Param			userId	path	string	true	"the user ID"
//	@Success		200		object	standardJsonResponse{}
//	@Failure		403		object	standardJsonResponse{}
//	@Failure		500		object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/users/{userId} [delete]
func (s *Service) handleUserDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		userID := r.PathValue("userId")
		idErr := validate.Var(userID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)

		user, userErr := s.UserDataSvc.GetUserByID(ctx, userID)
		if userErr != nil {
			s.Logger.Ctx(ctx).Error("handleUserDelete error", zap.Error(userErr),
				zap.String("entity_user_id", userID), zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, userErr)
			return
		}

		updateErr := s.UserDataSvc.DeleteUser(ctx, userID)
		if updateErr != nil {
			s.Logger.Ctx(ctx).Error("handleUserDelete error", zap.Error(updateErr),
				zap.String("user_id", userID), zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, updateErr)
			return
		}

		// don't attempt to send email to guest users
		if user.Email != "" {
			_ = s.Email.SendDeleteConfirmation(user.Name, user.Email)
		}

		// don't clear admins user cookies when deleting other users
		if userID == sessionUserID {
			s.Cookie.ClearUserCookies(w)
		}

		s.Success(w, r, http.StatusOK, nil, nil)
	}
}

// handleVerifyRequest sends verification Email
//
//	@Summary		Request Verification Email
//	@Description	Sends verification Email
//	@Tags			user
//	@Param			userId	path	string	true	"the user ID"
//	@Success		200		object	standardJsonResponse{}
//	@Success		400		object	standardJsonResponse{}
//	@Success		500		object	standardJsonResponse{}
//	@Router			/users/{userId}/request-verify [post]
func (s *Service) handleVerifyRequest() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID, _ := ctx.Value(contextKeyUserID).(*string)

		userID := r.PathValue("userId")
		idErr := validate.Var(userID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		user, verifyID, err := s.AuthDataSvc.UserVerifyRequest(ctx, userID)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleVerifyRequest error", zap.Error(err),
				zap.String("entity_user_id", userID), zap.Stringp("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		_ = s.Email.SendEmailVerification(user.Name, user.Email, verifyID)

		s.Success(w, r, http.StatusOK, nil, nil)
	}
}

// handleGetActiveCountries gets a list of registered users countries
//
//	@Summary		Get Active Countries
//	@Description	Gets a list of users countries
//	@Produce		json
//	@Success		200	object	standardJsonResponse{[]string}
//	@Failure		500	object	standardJsonResponse{}
//	@Router			/active-countries [get]
func (s *Service) handleGetActiveCountries() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID, _ := ctx.Value(contextKeyUserID).(*string)
		countries, err := s.UserDataSvc.GetActiveCountries(ctx)

		if err != nil {
			s.Logger.Ctx(ctx).Error("handleGetActiveCountries error", zap.Error(err),
				zap.Stringp("session_user_id", sessionUserID))
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

		ctx := r.Context()
		sessionUserID, _ := ctx.Value(contextKeyUserID).(*string)

		width, _ := strconv.Atoi(r.PathValue("width"))
		userID := r.PathValue("id")
		idErr := validate.Var(userID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		avatarGender := govatar.MALE
		userGender := r.PathValue("avatar")
		if userGender == "female" {
			avatarGender = govatar.FEMALE
		}

		var avatar image.Image
		if s.Config.AvatarService == "govatar" {
			avatar, _ = govatar.GenerateForUsername(avatarGender, userID)
		} else { // must be goadorable
			var err error
			avatar, _, err = image.Decode(bytes.NewReader(adorable.PseudoRandom([]byte(userID))))
			if err != nil {
				s.Logger.Ctx(ctx).Error(err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		img := transform.Resize(avatar, width, width, transform.Linear)
		buffer := new(bytes.Buffer)

		if err := png.Encode(buffer, img); err != nil {
			s.Logger.Ctx(ctx).Error("handleUserAvatar error", zap.Error(err), zap.String("entity_user_id", userID),
				zap.Stringp("session_user_id", sessionUserID))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "image/png")
		w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))

		if _, err := w.Write(buffer.Bytes()); err != nil {
			s.Logger.Ctx(ctx).Error("handleUserAvatar error", zap.Error(err), zap.String("entity_user_id", userID),
				zap.Stringp("session_user_id", sessionUserID))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

// handleUserOrganizationInvite processes an organization invite for the user
//
//	@Summary		User Organization Invite
//	@Description	Processes an organization invite for the user
//	@Tags			user
//	@Param			userId		path	string	true	"the user ID"
//	@Param			inviteId	path	string	true	"the invite ID"
//	@Success		200			object	standardJsonResponse{}
//	@Success		400			object	standardJsonResponse{}
//	@Success		500			object	standardJsonResponse{}
//	@Router			/users/{userId}/invite/organization/{inviteId} [post]
func (s *Service) handleUserOrganizationInvite() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID, _ := ctx.Value(contextKeyUserID).(*string)

		userID := r.PathValue("userId")
		idErr := validate.Var(userID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		inviteID := r.PathValue("inviteId")
		idErr = validate.Var(inviteID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		user, err := s.UserDataSvc.GetUserByID(ctx, userID)
		if err != nil {
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		orgInvite, err := s.OrganizationDataSvc.OrganizationUserGetInviteByID(ctx, inviteID)
		if err != nil {
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}
		if user.Email != orgInvite.Email {
			s.Failure(w, r, http.StatusInternalServerError, Errorf(EINVALID, err.Error()))
			return
		}

		orgID, inviteErr := s.OrganizationDataSvc.OrganizationAddUser(ctx, orgInvite.OrganizationID, userID, orgInvite.Role)
		if inviteErr != nil {
			s.Logger.Ctx(ctx).Error("handleUserOrganizationInvite error adding invited user to organization",
				zap.Error(inviteErr),
				zap.String("session_user_id", *sessionUserID),
				zap.String("invite_id", orgInvite.InviteID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		delInviteErr := s.OrganizationDataSvc.OrganizationDeleteUserInvite(ctx, orgInvite.InviteID)
		if delInviteErr != nil {
			s.Logger.Ctx(ctx).Error("handleUserOrganizationInvite error deleting user invite to organization",
				zap.Error(delInviteErr),
				zap.String("session_user_id", *sessionUserID),
				zap.String("invite_id", orgInvite.InviteID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		organization, orgErr := s.OrganizationDataSvc.OrganizationGetByID(ctx, orgID)
		if orgErr != nil {
			s.Logger.Ctx(ctx).Error("handleUserOrganizationInvite error getting organization",
				zap.Error(delInviteErr),
				zap.String("session_user_id", *sessionUserID),
				zap.String("invite_id", orgInvite.InviteID))
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
//
//	@Summary		User Department Invite
//	@Description	Processes an department invite for the user
//	@Tags			user
//	@Param			userId		path	string	true	"the user ID"
//	@Param			inviteId	path	string	true	"the invite ID"
//	@Success		200			object	standardJsonResponse{}
//	@Success		400			object	standardJsonResponse{}
//	@Success		500			object	standardJsonResponse{}
//	@Router			/users/{userId}/invite/department/{inviteId} [post]
func (s *Service) handleUserDepartmentInvite() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID, _ := ctx.Value(contextKeyUserID).(*string)

		userID := r.PathValue("userId")
		idErr := validate.Var(userID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		inviteID := r.PathValue("inviteId")
		idErr = validate.Var(inviteID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		user, err := s.UserDataSvc.GetUserByID(ctx, userID)
		if err != nil {
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		deptInvite, err := s.OrganizationDataSvc.DepartmentUserGetInviteByID(ctx, inviteID)
		if err != nil {
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}
		if user.Email != deptInvite.Email {
			s.Failure(w, r, http.StatusInternalServerError, Errorf(EINVALID, err.Error()))
			return
		}

		dept, deptErr := s.OrganizationDataSvc.DepartmentGetByID(ctx, deptInvite.DepartmentID)
		if deptErr != nil {
			s.Logger.Ctx(ctx).Error("handleUserDepartmentInvite get department error", zap.Error(deptErr),
				zap.String("department_id", dept.ID),
				zap.String("session_user_id", *sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, deptErr)
			return
		}

		org, orgErr := s.OrganizationDataSvc.OrganizationGetByID(ctx, dept.OrganizationID)
		if orgErr != nil {
			s.Logger.Ctx(ctx).Error("handleUserDepartmentInvite get organization error", zap.Error(orgErr),
				zap.String("organization_id", org.ID),
				zap.String("department_id", dept.ID),
				zap.String("session_user_id", *sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, orgErr)
			return
		}

		_, orgAddErr := s.OrganizationDataSvc.OrganizationUpsertUser(ctx, org.ID, userID, thunderdome.EntityMemberUserType)
		if orgAddErr != nil {
			s.Logger.Ctx(ctx).Error(
				"handleUserDepartmentInvite upsert organization user error", zap.Error(orgAddErr),
				zap.String("session_user_id", *sessionUserID),
				zap.String("organization_id", org.ID),
				zap.String("department_id", dept.ID),
				zap.String("user_id", userID), zap.String("user_role", thunderdome.EntityMemberUserType))
			s.Failure(w, r, http.StatusInternalServerError, orgAddErr)
			return
		}

		_, inviteErr := s.OrganizationDataSvc.DepartmentAddUser(ctx, deptInvite.DepartmentID, userID, deptInvite.Role)
		if inviteErr != nil {
			s.Logger.Ctx(ctx).Error("handleUserDepartmentInvite error adding invited user to organization",
				zap.Error(inviteErr),
				zap.String("session_user_id", *sessionUserID),
				zap.String("invite_id", deptInvite.InviteID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		delInviteErr := s.OrganizationDataSvc.OrganizationDeleteUserInvite(ctx, deptInvite.InviteID)
		if delInviteErr != nil {
			s.Logger.Ctx(ctx).Error("handleUserDepartmentInvite error deleting user invite to organization",
				zap.Error(delInviteErr),
				zap.String("session_user_id", *sessionUserID),
				zap.String("invite_id", deptInvite.InviteID))
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
//
//	@Summary		User Team Invite
//	@Description	Processes a team invite for the user
//	@Tags			user
//	@Param			userId		path	string	true	"the user ID"
//	@Param			inviteId	path	string	true	"the invite ID"
//	@Success		200			object	standardJsonResponse{}
//	@Success		400			object	standardJsonResponse{}
//	@Success		500			object	standardJsonResponse{}
//	@Router			/users/{userId}/invite/team/{inviteId} [post]
func (s *Service) handleUserTeamInvite() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID, _ := ctx.Value(contextKeyUserID).(*string)

		userID := r.PathValue("userId")
		idErr := validate.Var(userID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		inviteID := r.PathValue("inviteId")
		idErr = validate.Var(inviteID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		user, err := s.UserDataSvc.GetUserByID(ctx, userID)
		if err != nil {
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		teamInvite, err := s.TeamDataSvc.TeamUserGetInviteByID(ctx, inviteID)
		if err != nil {
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}
		if user.Email != teamInvite.Email {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, err.Error()))
			return
		}

		team, teamErr := s.TeamDataSvc.TeamGetByID(ctx, teamInvite.TeamID)
		if teamErr != nil {
			s.Logger.Ctx(ctx).Error("handleTeamInviteUser error", zap.Error(teamErr),
				zap.String("team_id", teamInvite.TeamID), zap.String("session_user_id", *sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, teamErr)
			return
		}

		// team is associated to organization, upsert user to organization
		if team.OrganizationID != "" {
			_, orgAddErr := s.OrganizationDataSvc.OrganizationUpsertUser(ctx, team.OrganizationID, userID, thunderdome.EntityMemberUserType)
			if orgAddErr != nil {
				s.Logger.Ctx(ctx).Error(
					"handleTeamInviteUser upsert organization user error", zap.Error(orgAddErr),
					zap.String("session_user_id", *sessionUserID),
					zap.String("organization_id", team.OrganizationID),
					zap.String("user_id", userID), zap.String("user_role", thunderdome.EntityMemberUserType))
				s.Failure(w, r, http.StatusInternalServerError, orgAddErr)
				return
			}
		}

		// team is associated to department, upsert user to department and organization
		if team.DepartmentID != "" {
			dept, deptErr := s.OrganizationDataSvc.DepartmentGetByID(ctx, team.DepartmentID)
			if deptErr != nil {
				s.Logger.Ctx(ctx).Error("handleTeamInviteUser get department error", zap.Error(deptErr),
					zap.String("department_id", dept.ID),
					zap.String("session_user_id", *sessionUserID))
				s.Failure(w, r, http.StatusInternalServerError, deptErr)
				return
			}

			org, orgErr := s.OrganizationDataSvc.OrganizationGetByID(ctx, dept.OrganizationID)
			if orgErr != nil {
				s.Logger.Ctx(ctx).Error("handleTeamInviteUser get organization error", zap.Error(orgErr),
					zap.String("organization_id", org.ID),
					zap.String("department_id", dept.ID),
					zap.String("session_user_id", *sessionUserID))
				s.Failure(w, r, http.StatusInternalServerError, orgErr)
				return
			}

			team.OrganizationID = org.ID

			_, orgAddErr := s.OrganizationDataSvc.OrganizationUpsertUser(ctx, org.ID, userID, thunderdome.EntityMemberUserType)
			if orgAddErr != nil {
				s.Logger.Ctx(ctx).Error(
					"handleTeamInviteUser upsert organization user error", zap.Error(orgAddErr),
					zap.String("session_user_id", *sessionUserID),
					zap.String("organization_id", org.ID),
					zap.String("department_id", dept.ID),
					zap.String("user_id", userID), zap.String("user_role", thunderdome.EntityMemberUserType))
				s.Failure(w, r, http.StatusInternalServerError, orgAddErr)
				return
			}

			_, deptAddErr := s.OrganizationDataSvc.DepartmentUpsertUser(ctx, team.DepartmentID, userID, thunderdome.EntityMemberUserType)
			if deptAddErr != nil {
				s.Logger.Ctx(ctx).Error(
					"handleTeamInviteUser upsert department user error", zap.Error(deptAddErr),
					zap.String("session_user_id", *sessionUserID),
					zap.String("department_id", team.DepartmentID),
					zap.String("user_id", userID), zap.String("user_role", thunderdome.EntityMemberUserType))
				s.Failure(w, r, http.StatusInternalServerError, deptAddErr)
				return
			}
		}

		_, inviteErr := s.TeamDataSvc.TeamAddUser(ctx, teamInvite.TeamID, userID, teamInvite.Role)
		if inviteErr != nil {
			s.Logger.Ctx(ctx).Error("handleUserRegistration error adding invited user to team", zap.Error(inviteErr),
				zap.String("session_user_id", *sessionUserID),
				zap.String("invite_id", teamInvite.InviteID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		delInviteErr := s.TeamDataSvc.TeamDeleteUserInvite(ctx, teamInvite.InviteID)
		if delInviteErr != nil {
			s.Logger.Ctx(ctx).Error("handleUserRegistration error deleting user invite to team", zap.Error(delInviteErr),
				zap.String("session_user_id", *sessionUserID),
				zap.String("invite_id", teamInvite.InviteID))
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
//
//	@Summary		Get User Credential
//	@Description	Gets a users credential
//	@Tags			user
//	@Produce		json
//	@Param			userId	path	string	true	"the user ID"
//	@Success		200		object	standardJsonResponse{data=thunderdome.Credential}
//	@Failure		403		object	standardJsonResponse{}
//	@Failure		500		object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/users/{userId}/credential [get]
func (s *Service) handleUserCredential() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)

		userID := r.PathValue("userId")
		idErr := validate.Var(userID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		credential, userErr := s.UserDataSvc.GetUserCredentialByUserID(ctx, userID)
		if userErr != nil {
			s.Logger.Ctx(ctx).Error("handleUserCredential error", zap.Error(userErr),
				zap.String("entity_user_id", userID),
				zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, userErr)
			return
		}

		s.Success(w, r, http.StatusOK, credential, nil)
	}
}

// handleChangeEmailRequest attempts to send a change email Email
//
//	@Summary		Change Email Request
//	@Description	Sends a change email request Email to user
//	@Tags			auth
//	@Produce		json
//	@Param			userId	path	string	true	"the user ID"
//	@Success		200		object	standardJsonResponse{}
//	@Failure		403		object	standardJsonResponse{}
//	@Failure		500		object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/user/{userId}/email-change [get]
func (s *Service) handleChangeEmailRequest() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID, _ := ctx.Value(contextKeyUserID).(string)

		userID := r.PathValue("userId")
		idErr := validate.Var(userID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		u, userErr := s.UserDataSvc.GetUserByID(ctx, userID)
		if userErr != nil {
			s.Logger.Ctx(ctx).Error("handleChangeEmailRequest error", zap.Error(userErr),
				zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, userErr)
			return
		}

		changeId, changeErr := s.UserDataSvc.RequestEmailChange(ctx, userID)
		if changeErr == nil {
			_ = s.Email.SendEmailChangeRequest(u.Name, u.Email, changeId)
		} else {
			s.Logger.Ctx(ctx).Error("handleChangeEmailRequest error", zap.Error(changeErr),
				zap.String("user_email", sanitizeUserInputForLogs(u.Email)))
			s.Failure(w, r, http.StatusInternalServerError, userErr)
			return
		}

		s.Success(w, r, http.StatusOK, nil, nil)
	}
}

type changeEmailRequestBody struct {
	Email string `json:"email" validate:"required,email"`
}

// handleChangeEmailAction attempts to change the users email
//
//	@Summary		Change Email Action
//	@Description	Attempts to change the users email
//	@Description	Requires a valid change ID
//	@Tags			auth
//	@Produce		json
//	@Param			userId		path	string					true	"the user ID"
//	@Param			changeId	path	string					true	"the change ID"
//	@Param			user		body	changeEmailRequestBody	true	"the user object to update"
//	@Success		200			object	standardJsonResponse{data=thunderdome.User}
//	@Failure		403			object	standardJsonResponse{}
//	@Failure		500			object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/user/{userId}/email-change/{changeId} [post]
func (s *Service) handleChangeEmailAction() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID, _ := ctx.Value(contextKeyUserID).(string)

		userID := r.PathValue("userId")
		idErr := validate.Var(userID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		changeId := r.PathValue("changeId")
		idErr = validate.Var(changeId, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		u, userErr := s.UserDataSvc.GetUserByID(ctx, userID)
		if userErr != nil {
			s.Logger.Ctx(ctx).Error("handleChangeEmailAction error", zap.Error(userErr),
				zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, userErr)
			return
		}

		body, bodyErr := io.ReadAll(r.Body)
		if bodyErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
			return
		}

		var cr = changeEmailRequestBody{}
		jsonErr := json.Unmarshal(body, &cr)
		if jsonErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		inputErr := validate.Struct(cr)
		if inputErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, inputErr.Error()))
			return
		}

		newEmail := strings.ToLower(cr.Email)

		changeErr := s.UserDataSvc.ConfirmEmailChange(ctx, userID, changeId, newEmail)
		if changeErr == nil {
			_ = s.Email.SendEmailChangeConfirmation(u.Name, u.Email, newEmail)
		} else {
			s.Logger.Ctx(ctx).Error("handleChangeEmailAction error", zap.Error(changeErr),
				zap.String("user_email", sanitizeUserInputForLogs(u.Email)))
			s.Failure(w, r, http.StatusInternalServerError, userErr)
			return
		}

		updatedUser, userErr := s.UserDataSvc.GetUserByID(ctx, userID)
		if userErr != nil {
			s.Logger.Ctx(ctx).Error("handleChangeEmailAction error", zap.Error(userErr),
				zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, userErr)
			return
		}

		s.Success(w, r, http.StatusOK, updatedUser, nil)
	}
}
