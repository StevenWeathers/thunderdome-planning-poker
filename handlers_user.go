package main

import (
	"bytes"
	"image"
	"image/png"
	"log"
	"net/http"
	"strconv"

	"github.com/anthonynsimon/bild/transform"
	"github.com/gorilla/mux"
	"github.com/ipsn/go-adorable"
	"github.com/o1egl/govatar"
)

// handleUpdatePassword attempts to update a users password
func (s *server) handleUpdatePassword() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		keyVal := s.getJSONRequestBody(r, w)
		UserID := r.Context().Value(contextKeyUserID).(string)

		UserPassword, passwordErr := ValidateUserPassword(
			keyVal["warriorPassword1"].(string),
			keyVal["warriorPassword2"].(string),
		)

		if passwordErr != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		UserName, UserEmail, updateErr := s.database.UserUpdatePassword(UserID, UserPassword)
		if updateErr != nil {
			log.Println("error attempting to update user password : " + updateErr.Error() + "\n")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		s.email.SendPasswordUpdate(UserName, UserEmail)

		return
	}
}

// handleUserProfile returns the users profile if it matches their session
func (s *server) handleUserProfile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		UserID := vars["id"]
		UserCookieID := r.Context().Value(contextKeyUserID).(string)

		if UserID != UserCookieID {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		User, UserErr := s.database.GetUser(UserID)
		if UserErr != nil {
			log.Println("error finding user : " + UserErr.Error() + "\n")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		s.respondWithJSON(w, http.StatusOK, User)
	}
}

// handleUserProfileUpdate attempts to update users profile
func (s *server) handleUserProfileUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		keyVal := s.getJSONRequestBody(r, w)
		UserName := keyVal["warriorName"].(string)
		UserAvatar := keyVal["warriorAvatar"].(string)
		NotificationsEnabled, _ := keyVal["notificationsEnabled"].(bool)
		Country := keyVal["country"].(string)
		Locale := keyVal["locale"].(string)
		Company := keyVal["company"].(string)
		JobTitle := keyVal["jobTitle"].(string)

		UserID := vars["id"]
		UserCookieID := r.Context().Value(contextKeyUserID).(string)
		if UserID != UserCookieID {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		updateErr := s.database.UpdateUserProfile(UserID, UserName, UserAvatar, NotificationsEnabled, Country, Locale, Company, JobTitle)
		if updateErr != nil {
			log.Println("error attempting to update user profile : " + updateErr.Error() + "\n")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		user, UserErr := s.database.GetUser(UserID)
		if UserErr != nil {
			log.Println("error reloading user after update : " + UserErr.Error() + "\n")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		s.respondWithJSON(w, http.StatusOK, user)
	}
}

// handleAccountVerification attempts to verify a users account
func (s *server) handleAccountVerification() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		keyVal := s.getJSONRequestBody(r, w)
		VerifyID := keyVal["verifyId"].(string)

		verifyErr := s.database.VerifyUserAccount(VerifyID)
		if verifyErr != nil {
			log.Println("error attempting to verify user account : " + verifyErr.Error() + "\n")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		return
	}
}

// handleUserDelete attempts to delete a users account
func (s *server) handleUserDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		UserID := vars["id"]
		UserCookieID := r.Context().Value(contextKeyUserID).(string)
		if UserID != UserCookieID {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		User, UserErr := s.database.GetUser(UserID)
		if UserErr != nil {
			log.Println("error finding user : " + UserErr.Error() + "\n")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		updateErr := s.database.DeleteUser(UserID)
		if updateErr != nil {
			log.Println("error attempting to delete user : " + updateErr.Error() + "\n")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		s.email.SendDeleteConfirmation(User.UserName, User.UserEmail)

		s.clearUserCookies(w)

		return
	}
}

// handleUserAvatar creates an avatar for the given user by ID
func (s *server) handleUserAvatar() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

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
		if s.config.AvatarService == "govatar" {
			avatar, _ = govatar.GenerateForUsername(AvatarGender, UserID)
		} else { // must be goadorable
			var err error
			avatar, _, err = image.Decode(bytes.NewReader(adorable.PseudoRandom([]byte(UserID))))
			if err != nil {
				log.Fatalln(err)
			}
		}

		img := transform.Resize(avatar, Width, Width, transform.Linear)
		buffer := new(bytes.Buffer)

		if err := png.Encode(buffer, img); err != nil {
			log.Println("unable to encode image.")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "image/png")
		w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))

		if _, err := w.Write(buffer.Bytes()); err != nil {
			log.Println("unable to write image.")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

// handleGetActiveCountries gets a list of registered users countries
func (s *server) handleGetActiveCountries() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		countries, err := s.database.GetActiveCountries()

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Cache-Control", "max-age=3600") // cache for 1 hour just to decrease load
		s.respondWithJSON(w, http.StatusOK, countries)
	}
}
