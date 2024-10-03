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
func (b *Service) ServeWs() http.HandlerFunc {
	return b.hub.WebSocketHandler("teamId", func(w http.ResponseWriter, r *http.Request, c *wshub.Connection, roomID string) *wshub.AuthError {
		ctx := r.Context()
		var user *thunderdome.User

		sessionID, cookieErr := b.validateSessionCookie(w, r)
		if cookieErr != nil && cookieErr.Error() != "COOKIE_NOT_FOUND" {
			authErr := wshub.AuthError{
				Code:    4001,
				Message: "unauthorized",
			}
			return &authErr
		}

		if sessionID != "" {
			var userErr error
			user, userErr = b.AuthService.GetSessionUser(ctx, sessionID)
			if userErr != nil {
				authErr := wshub.AuthError{
					Code:    4001,
					Message: "unauthorized",
				}
				return &authErr
			}
		} else {
			userID, err := b.validateUserCookie(w, r)
			if err != nil {
				authErr := wshub.AuthError{
					Code:    4001,
					Message: "unauthorized",
				}
				return &authErr
			}

			var userErr error
			user, userErr = b.UserService.GetGuestUser(ctx, userID)
			if userErr != nil {
				authErr := wshub.AuthError{
					Code:    4001,
					Message: "unauthorized",
				}
				return &authErr
			}
		}

		// make sure team is legit
		_, retroErr := b.TeamService.TeamGet(context.Background(), roomID)
		if retroErr != nil {
			authErr := wshub.AuthError{
				Code:    4004,
				Message: "team not found",
			}
			return &authErr
		}

		// make sure user is a team user
		_, UserErr := b.TeamService.TeamUserRole(ctx, user.ID, roomID)
		if UserErr != nil {
			b.logger.Ctx(ctx).Error("REQUIRES_TEAM_USER", zap.Error(UserErr),
				zap.String("team_id", roomID), zap.String("session_user_id", user.ID))

			authErr := wshub.AuthError{
				Code:    4005,
				Message: "REQUIRES_TEAM_USER",
			}
			return &authErr
		}

		sub := b.hub.NewSubscriber(c.Ws, user.ID, roomID)

		initEvent := wshub.CreateSocketEvent("init", "", user.ID)
		_ = sub.Conn.Write(websocket.TextMessage, initEvent)

		go sub.WritePump()
		go sub.ReadPump(ctx, b.hub)

		return nil
	})
}

// APIEvent handles api driven events into the team checkin (if active)
func (b *Service) APIEvent(ctx context.Context, teamID string, userID, eventType string, eventValue string) error {
	return b.hub.ProcessAPIEventHandler(ctx, userID, teamID, eventType, eventValue)
}
