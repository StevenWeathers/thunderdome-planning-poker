package http

import (
	"bytes"
	"encoding/json"
	"image"
	"image/png"
	"io"
	"net/http"
	"strconv"

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
	Country              string `json:"country" validate:"len=2"`
	Locale               string `json:"locale" validate:"len=2"`
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

		if SessionUserType == adminUserType {
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
		SessionUserID := ctx.Value(contextKeyUserID).(*string)
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
		SessionUserID := ctx.Value(contextKeyUserID).(*string)
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
		SessionUserID := ctx.Value(contextKeyUserID).(*string)

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
