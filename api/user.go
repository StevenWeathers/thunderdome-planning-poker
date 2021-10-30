package api

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// handleUserProfile returns the users profile if it matches their session
func (a *api) handleUserProfile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		UserID := vars["id"]
		UserCookieID := r.Context().Value(contextKeyUserID).(string)

		if UserID != UserCookieID {
			a.respondWithStandardJSON(w, http.StatusForbidden, false, nil, nil, nil)
		}

		User, UserErr := a.db.GetUser(UserID)
		if UserErr != nil {
			log.Println("error finding user : " + UserErr.Error() + "\n")
			errors := make([]string, 0)
			errors = append(errors, UserErr.Error())
			a.respondWithStandardJSON(w, http.StatusInternalServerError, false, errors, nil, nil)
		}

		a.respondWithStandardJSON(w, http.StatusOK, true, nil, User, nil)
	}
}

// handleUserProfileUpdate attempts to update users profile
func (a *api) handleUserProfileUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		keyVal := a.getJSONRequestBody(r, w)
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
			a.respondWithStandardJSON(w, http.StatusForbidden, false, nil, nil, nil)
		}

		updateErr := a.db.UpdateUserProfile(UserID, UserName, UserAvatar, NotificationsEnabled, Country, Locale, Company, JobTitle)
		if updateErr != nil {
			log.Println("error attempting to update user profile : " + updateErr.Error() + "\n")
			errors := make([]string, 0)
			errors = append(errors, updateErr.Error())
			a.respondWithStandardJSON(w, http.StatusInternalServerError, false, errors, nil, nil)
		}

		user, UserErr := a.db.GetUser(UserID)
		if UserErr != nil {
			log.Println("error reloading user after update : " + UserErr.Error() + "\n")
			errors := make([]string, 0)
			errors = append(errors, UserErr.Error())
			a.respondWithStandardJSON(w, http.StatusInternalServerError, false, errors, nil, nil)
		}

		a.respondWithStandardJSON(w, http.StatusOK, true, nil, user, nil)
	}
}

// handleUserDelete attempts to delete a users account
func (a *api) handleUserDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		UserID := vars["id"]
		UserCookieID := r.Context().Value(contextKeyUserID).(string)
		if UserID != UserCookieID {
			a.respondWithStandardJSON(w, http.StatusForbidden, false, nil, nil, nil)
		}

		User, UserErr := a.db.GetUser(UserID)
		if UserErr != nil {
			log.Println("error finding user : " + UserErr.Error() + "\n")
			errors := make([]string, 0)
			errors = append(errors, UserErr.Error())
			a.respondWithStandardJSON(w, http.StatusInternalServerError, false, errors, nil, nil)
		}

		updateErr := a.db.DeleteUser(UserID)
		if updateErr != nil {
			log.Println("error attempting to delete user : " + updateErr.Error() + "\n")
			errors := make([]string, 0)
			errors = append(errors, updateErr.Error())
			a.respondWithStandardJSON(w, http.StatusInternalServerError, false, errors, nil, nil)
		}

		a.email.SendDeleteConfirmation(User.UserName, User.UserEmail)

		a.clearUserCookies(w)

		a.respondWithStandardJSON(w, http.StatusOK, true, nil, nil, nil)
	}
}

// handleGetActiveCountries gets a list of registered users countries
func (a *api) handleGetActiveCountries() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		countries, err := a.db.GetActiveCountries()

		if err != nil {
			errors := make([]string, 0)
			errors = append(errors, err.Error())
			a.respondWithStandardJSON(w, http.StatusInternalServerError, false, errors, nil, nil)
		}

		w.Header().Set("Cache-Control", "max-age=3600") // cache for 1 hour just to decrease load
		a.respondWithStandardJSON(w, http.StatusOK, true, nil, countries, nil)
	}
}
