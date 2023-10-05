package http

import (
	"net/http"

	"go.uber.org/zap"
)

// handleCleanBattles handles cleaning up old battles (ADMIN Manually Triggered)
// @Summary      Clean Old Battles
// @Description  Deletes battles older than {config.cleanup_battles_days_old} based on last activity date
// @Tags         maintenance
// @Produce      json
// @Success      200  object  standardJsonResponse{}
// @Failure      500  object  standardJsonResponse{}
// @Security     ApiKeyAuth
// @Router       /maintenance/clean-battles [delete]
func (s *Service) handleCleanBattles() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		SessionUserID := ctx.Value(contextKeyUserID).(string)
		DaysOld := s.Config.CleanupBattlesDaysOld

		err := s.PokerDataSvc.PurgeOldGames(ctx, DaysOld)
		if err != nil {
			s.Logger.Ctx(ctx).Error(
				"handleCleanBattles error", zap.Error(err), zap.String("session_user_id", SessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, nil, nil)
	}
}

// handleCleanRetros handles cleaning up old retros (ADMIN Manually Triggered)
// @Summary      Clean Old Retros
// @Description  Deletes retros older than {config.cleanup_retros_days_old} based on last activity date
// @Tags         maintenance
// @Produce      json
// @Success      200  object  standardJsonResponse{}
// @Failure      500  object  standardJsonResponse{}
// @Security     ApiKeyAuth
// @Router       /maintenance/clean-retros [delete]
func (s *Service) handleCleanRetros() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		SessionUserID := ctx.Value(contextKeyUserID).(string)
		DaysOld := s.Config.CleanupRetrosDaysOld

		err := s.RetroDataSvc.CleanRetros(r.Context(), DaysOld)
		if err != nil {
			s.Logger.Ctx(ctx).Error(
				"handleCleanRetros error", zap.Error(err), zap.String("session_user_id", SessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, nil, nil)
	}
}

// handleCleanStoryboards handles cleaning up old storyboards (ADMIN Manually Triggered)
// @Summary      Clean Old Storyboards
// @Description  Deletes storyboards older than {config.cleanup_storyboards_days_old} based on last activity date
// @Tags         maintenance
// @Produce      json
// @Success      200  object  standardJsonResponse{}
// @Failure      500  object  standardJsonResponse{}
// @Security     ApiKeyAuth
// @Router       /maintenance/clean-storyboards [delete]
func (s *Service) handleCleanStoryboards() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		SessionUserID := ctx.Value(contextKeyUserID).(string)
		DaysOld := s.Config.CleanupStoryboardsDaysOld

		err := s.StoryboardDataSvc.CleanStoryboards(r.Context(), DaysOld)
		if err != nil {
			s.Logger.Ctx(ctx).Error(
				"handleCleanStoryboards error", zap.Error(err), zap.String("session_user_id", SessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, nil, nil)
	}
}

// handleCleanGuests handles cleaning up old guests (ADMIN Manaually Triggered)
// @Summary      Clean Old Guests
// @Description  Deletes guest users older than {config.cleanup_guests_days_old} based on last activity date
// @Tags         maintenance
// @Produce      json
// @Success      200  object  standardJsonResponse{}
// @Failure      500  object  standardJsonResponse{}
// @Security     ApiKeyAuth
// @Router       /maintenance/clean-guests [delete]
func (s *Service) handleCleanGuests() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		SessionUserID := ctx.Value(contextKeyUserID).(string)
		DaysOld := s.Config.CleanupGuestsDaysOld

		err := s.UserDataSvc.CleanGuests(r.Context(), DaysOld)
		if err != nil {
			s.Logger.Ctx(ctx).Error(
				"handleCleanGuests error", zap.Error(err), zap.String("session_user_id", SessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, nil, nil)
	}
}
