package api

import (
	"log"
	"net/http"

	"github.com/spf13/viper"
)

// handleCleanBattles handles cleaning up old battles (ADMIN Manaually Triggered)
// @Summary Clean Old Battles
// @Description Deletes battles older than {config.cleanup_battles_days_old} based on last activity date
// @Tags maintenance
// @Produce  json
// @Success 200 object standardJsonResponse{}
// @Failure 500 object standardJsonResponse{}
// @Router /maintenance/clean-battles [delete]
func (a *api) handleCleanBattles() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		DaysOld := viper.GetInt("config.cleanup_battles_days_old")

		err := a.db.CleanBattles(DaysOld)
		if err != nil {
			errors := make([]string, 0)
			errors = append(errors, err.Error())
			a.respondWithStandardJSON(w, http.StatusInternalServerError, false, errors, nil, nil)
		}

		a.respondWithStandardJSON(w, http.StatusOK, true, nil, nil, nil)
	}
}

// handleCleanGuests handles cleaning up old guests (ADMIN Manaually Triggered)
// @Summary Clean Old Guests
// @Description Deletes guest users older than {config.cleanup_guests_days_old} based on last activity date
// @Tags maintenance
// @Produce  json
// @Success 200 object standardJsonResponse{}
// @Failure 500 object standardJsonResponse{}
// @Router /maintenance/clean-guests [delete]
func (a *api) handleCleanGuests() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		DaysOld := viper.GetInt("config.cleanup_guests_days_old")

		err := a.db.CleanGuests(DaysOld)
		if err != nil {
			errors := make([]string, 0)
			errors = append(errors, err.Error())
			a.respondWithStandardJSON(w, http.StatusInternalServerError, false, errors, nil, nil)
		}

		a.respondWithStandardJSON(w, http.StatusOK, true, nil, nil, nil)
	}
}

// handleLowercaseUserEmails handles lowercasing any user emails that have any uppercase letters
// @Summary Lowercase User Emails
// @Description Lowercases any user emails that have uppercase letters to prevent duplicate email registration
// @Tags maintenance
// @Produce  json
// @Success 200 object standardJsonResponse{}
// @Failure 500 object standardJsonResponse{}
// @Router /maintenance/lowercase-emails [patch]
func (a *api) handleLowercaseUserEmails() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		lowercasedUsers, err := a.db.LowercaseUserEmails()
		if err != nil {
			errors := make([]string, 0)
			errors = append(errors, err.Error())
			a.respondWithStandardJSON(w, http.StatusInternalServerError, false, errors, nil, nil)
		}

		log.Println("Lowercased", len(lowercasedUsers), "user emails")
		for _, u := range lowercasedUsers {
			a.email.SendEmailUpdate(u.UserName, u.UserEmail)
		}

		mergedUsers, err := a.db.MergeDuplicateAccounts()
		if err != nil {
			errors := make([]string, 0)
			errors = append(errors, err.Error())
			a.respondWithStandardJSON(w, http.StatusInternalServerError, false, errors, nil, nil)
		}

		log.Println("Merged", len(mergedUsers), "user accounts")
		for _, u := range mergedUsers {
			a.email.SendMergedUpdate(u.UserName, u.UserEmail)
		}

		a.respondWithStandardJSON(w, http.StatusOK, true, nil, nil, nil)
	}
}
