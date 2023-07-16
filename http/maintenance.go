package http

import (
	"net/http"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// handleCleanBattles handles cleaning up old battles (ADMIN Manually Triggered)
// @Summary Clean Old Battles
// @Description Deletes battles older than {config.cleanup_battles_days_old} based on last activity date
// @Tags maintenance
// @Produce  json
// @Success 200 object standardJsonResponse{}
// @Failure 500 object standardJsonResponse{}
// @Security ApiKeyAuth
// @Router /maintenance/clean-battles [delete]
func (s *Service) handleCleanBattles() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		DaysOld := viper.GetInt("config.cleanup_battles_days_old")

		err := s.PokerDataSvc.PurgeOldGames(r.Context(), DaysOld)
		if err != nil {
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, nil, nil)
	}
}

// handleCleanRetros handles cleaning up old retros (ADMIN Manually Triggered)
// @Summary Clean Old Retros
// @Description Deletes retros older than {config.cleanup_retros_days_old} based on last activity date
// @Tags maintenance
// @Produce  json
// @Success 200 object standardJsonResponse{}
// @Failure 500 object standardJsonResponse{}
// @Security ApiKeyAuth
// @Router /maintenance/clean-retros [delete]
func (s *Service) handleCleanRetros() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		DaysOld := viper.GetInt("config.cleanup_retros_days_old")

		err := s.RetroDataSvc.CleanRetros(r.Context(), DaysOld)
		if err != nil {
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, nil, nil)
	}
}

// handleCleanStoryboards handles cleaning up old storyboards (ADMIN Manually Triggered)
// @Summary Clean Old Storyboards
// @Description Deletes storyboards older than {config.cleanup_storyboards_days_old} based on last activity date
// @Tags maintenance
// @Produce  json
// @Success 200 object standardJsonResponse{}
// @Failure 500 object standardJsonResponse{}
// @Security ApiKeyAuth
// @Router /maintenance/clean-storyboards [delete]
func (s *Service) handleCleanStoryboards() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		DaysOld := viper.GetInt("config.cleanup_storyboards_days_old")

		err := s.StoryboardDataSvc.CleanStoryboards(r.Context(), DaysOld)
		if err != nil {
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, nil, nil)
	}
}

// handleCleanGuests handles cleaning up old guests (ADMIN Manaually Triggered)
// @Summary Clean Old Guests
// @Description Deletes guest users older than {config.cleanup_guests_days_old} based on last activity date
// @Tags maintenance
// @Produce  json
// @Success 200 object standardJsonResponse{}
// @Failure 500 object standardJsonResponse{}
// @Security ApiKeyAuth
// @Router /maintenance/clean-guests [delete]
func (s *Service) handleCleanGuests() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		DaysOld := viper.GetInt("config.cleanup_guests_days_old")

		err := s.UserDataSvc.CleanGuests(r.Context(), DaysOld)
		if err != nil {
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, nil, nil)
	}
}

// handleLowercaseUserEmails handles lowercasing any user emails that have any uppercase letters
// @Summary Lowercase User Emails
// @Description Lowercases any user emails that have uppercase letters to prevent duplicate Email registration
// @Tags maintenance
// @Produce  json
// @Success 200 object standardJsonResponse{}
// @Failure 500 object standardJsonResponse{}
// @Security ApiKeyAuth
// @Router /maintenance/lowercase-emails [patch]
func (s *Service) handleLowercaseUserEmails() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		lowercasedUsers, err := s.UserDataSvc.LowercaseUserEmails(r.Context())
		if err != nil {
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Logger.Info("Lowercased user emails", zap.Int("count", len(lowercasedUsers)))
		for _, u := range lowercasedUsers {
			_ = s.Email.SendEmailUpdate(u.Name, u.Email)
		}

		mergedUsers, err := s.UserDataSvc.MergeDuplicateAccounts(r.Context())
		if err != nil {
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Logger.Info("Merged user accounts", zap.Int("count", len(mergedUsers)))
		for _, u := range mergedUsers {
			_ = s.Email.SendMergedUpdate(u.Name, u.Email)
		}

		s.Success(w, r, http.StatusOK, nil, nil)
	}
}
