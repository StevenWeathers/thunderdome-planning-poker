package checkin

import (
	"context"
	"net/http"

	"github.com/StevenWeathers/thunderdome-planning-poker/internal/wshub"
	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"

	"go.uber.org/zap"

	"github.com/gorilla/websocket"
)

// ServeWs handles websocket requests from the peer.
func (s *Service) ServeWs() http.HandlerFunc {
	return s.hub.WebSocketHandler("teamId", func(w http.ResponseWriter, r *http.Request, c *wshub.Connection, roomID string) *wshub.AuthError {
		ctx := r.Context()
		var user *thunderdome.User

		sessionID, cookieErr := s.validateSessionCookie(w, r)
		if cookieErr != nil && cookieErr.Error() != "COOKIE_NOT_FOUND" {
			authErr := wshub.AuthError{
				Code:    4001,
				Message: "unauthorized",
			}
			return &authErr
		}

		if sessionID != "" {
			var userErr error
			user, userErr = s.AuthService.GetSessionUserByID(ctx, sessionID)
			if userErr != nil {
				authErr := wshub.AuthError{
					Code:    4001,
					Message: "unauthorized",
				}
				return &authErr
			}
		} else {
			userID, err := s.validateUserCookie(w, r)
			if err != nil {
				authErr := wshub.AuthError{
					Code:    4001,
					Message: "unauthorized",
				}
				return &authErr
			}

			var userErr error
			user, userErr = s.UserService.GetGuestUserByID(ctx, userID)
			if userErr != nil {
				authErr := wshub.AuthError{
					Code:    4001,
					Message: "unauthorized",
				}
				return &authErr
			}
		}

		// make sure team is legit
		_, retroErr := s.TeamService.TeamGetByID(context.Background(), roomID)
		if retroErr != nil {
			authErr := wshub.AuthError{
				Code:    4004,
				Message: "team not found",
			}
			return &authErr
		}

		// make sure user is a team user
		_, UserErr := s.TeamService.TeamUserRoleByUserID(ctx, user.ID, roomID)
		if UserErr != nil {
			s.logger.Ctx(ctx).Error("REQUIRES_TEAM_USER", zap.Error(UserErr),
				zap.String("team_id", roomID), zap.String("session_user_id", user.ID))

			authErr := wshub.AuthError{
				Code:    4005,
				Message: "REQUIRES_TEAM_USER",
			}
			return &authErr
		}

		sub := s.hub.NewSubscriber(c.Ws, user.ID, roomID)

		initEvent := wshub.CreateSocketEvent("init", "", user.ID)
		_ = sub.Conn.Write(websocket.TextMessage, initEvent)

		go sub.WritePump()
		go sub.ReadPump(ctx, s.hub)

		return nil
	})
}

// APIEvent handles api driven events into the team checkin (if active)
func (s *Service) APIEvent(ctx context.Context, teamID string, userID, eventType string, eventValue string) (any, error) {
	return s.hub.ProcessAPIEventHandler(ctx, userID, teamID, eventType, eventValue)
}
