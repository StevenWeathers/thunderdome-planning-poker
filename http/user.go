package http

import (
	"bytes"
	"encoding/json"
	"image"
	"image/png"
	"io"
	"net/http"
	"strconv"

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
		UserID := ctx.Value(contextKeyUserID).(string)

		User, UserErr := s.UserDataSvc.GetUser(ctx, UserID)
		if UserErr != nil {
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
		vars := mux.Vars(r)
		UserID := vars["userId"]

		User, UserErr := s.UserDataSvc.GetUser(r.Context(), UserID)
		if UserErr != nil {
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
			updateErr := s.UserDataSvc.UpdateUserAccount(ctx, UserID, profile.Name, profile.Email, profile.Avatar, profile.NotificationsEnabled, profile.Country, profile.Locale, profile.Company, profile.JobTitle)
			if updateErr != nil {
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
				updateErr = s.UserDataSvc.UpdateUserProfile(ctx, UserID, profile.Name, profile.Avatar, profile.NotificationsEnabled, profile.Country, profile.Locale, profile.Company, profile.JobTitle)
			} else {
				updateErr = s.UserDataSvc.UpdateUserProfileLdap(ctx, UserID, profile.Avatar, profile.NotificationsEnabled, profile.Country, profile.Locale, profile.Company, profile.JobTitle)
			}
			if updateErr != nil {
				s.Failure(w, r, http.StatusInternalServerError, updateErr)
				return
			}
		}

		user, UserErr := s.UserDataSvc.GetUser(ctx, UserID)
		if UserErr != nil {
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
		UserCookieID := ctx.Value(contextKeyUserID).(string)

		User, UserErr := s.UserDataSvc.GetUser(ctx, UserID)
		if UserErr != nil {
			s.Failure(w, r, http.StatusInternalServerError, UserErr)
			return
		}

		updateErr := s.UserDataSvc.DeleteUser(ctx, UserID)
		if updateErr != nil {
			s.Failure(w, r, http.StatusInternalServerError, updateErr)
			return
		}

		_ = s.Email.SendDeleteConfirmation(User.Name, User.Email)

		// don't clear admins user cookies when deleting other users
		if UserID == UserCookieID {
			s.clearUserCookies(w)
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
		vars := mux.Vars(r)
		UserID := vars["userId"]

		User, VerifyId, err := s.AuthDataSvc.UserVerifyRequest(r.Context(), UserID)
		if err != nil {
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
		countries, err := s.UserDataSvc.GetActiveCountries(r.Context())

		if err != nil {
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
				s.Logger.Ctx(ctx).Fatal(err.Error())
			}
		}

		img := transform.Resize(avatar, Width, Width, transform.Linear)
		buffer := new(bytes.Buffer)

		if err := png.Encode(buffer, img); err != nil {
			s.Logger.Ctx(ctx).Error("unable to encode image.")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "image/png")
		w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))

		if _, err := w.Write(buffer.Bytes()); err != nil {
			s.Logger.Ctx(ctx).Error("unable to write image.")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
