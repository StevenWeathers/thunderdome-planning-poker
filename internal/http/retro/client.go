package retro

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/StevenWeathers/thunderdome-planning-poker/internal/wshub"
	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"

	"go.uber.org/zap"

	"github.com/gorilla/websocket"
)

// ServeWs handles websocket requests from the peer.
func (b *Service) ServeWs() http.HandlerFunc {
	return b.hub.WebSocketHandler("retroId", func(w http.ResponseWriter, r *http.Request, c *wshub.Connection, roomID string) *wshub.AuthError {
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
			user, userErr = b.AuthService.GetSessionUserByID(ctx, sessionID)
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
			user, userErr = b.UserService.GetGuestUserByID(ctx, userID)
			if userErr != nil {
				authErr := wshub.AuthError{
					Code:    4001,
					Message: "unauthorized",
				}
				return &authErr
			}
		}

		// make sure retro is legit
		retro, retroErr := b.RetroService.RetroGetByID(roomID, user.ID)
		if retroErr != nil {
			authErr := wshub.AuthError{
				Code:    4004,
				Message: "retro not found",
			}
			return &authErr
		}

		// check users retro active status
		userErr := b.RetroService.GetRetroUserActiveStatus(roomID, user.ID)
		userAllowed := userErr == nil
		if userErr != nil && !errors.Is(userErr, sql.ErrNoRows) {
			usrErrMsg := userErr.Error()
			var authErr wshub.AuthError

			if usrErrMsg == "DUPLICATE_RETRO_USER" {
				authErr = wshub.AuthError{
					Code:    4003,
					Message: "duplicate session",
				}
			} else {
				b.logger.Ctx(ctx).Error("error finding user", zap.Error(userErr),
					zap.String("retro_id", roomID), zap.String("session_user_id", user.ID))
				authErr = wshub.AuthError{
					Code:    4005,
					Message: "internal error",
				}
			}
			return &authErr
		} else if retro.JoinCode != "" && (userErr != nil && errors.Is(userErr, sql.ErrNoRows)) {
			jcrEvent := wshub.CreateSocketEvent("join_code_required", "", user.ID)
			_ = c.Write(websocket.TextMessage, jcrEvent)

			for {
				_, msg, err := c.Ws.ReadMessage()
				if err != nil {
					if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
						b.logger.Ctx(ctx).Error("unexpected close error", zap.Error(err),
							zap.String("retro_id", roomID), zap.String("session_user_id", user.ID))
					}
					break
				}

				keyVal := make(map[string]string)
				err = json.Unmarshal(msg, &keyVal)
				if err != nil {
					b.logger.Error("unexpected message error", zap.Error(err),
						zap.String("retro_id", roomID), zap.String("session_user_id", user.ID))
				}

				if keyVal["type"] == "auth_retro" && keyVal["value"] == retro.JoinCode {
					// join code is valid, continue to room
					userAllowed = true
					break
				} else if keyVal["type"] == "auth_retro" {
					authIncorrect := wshub.CreateSocketEvent("join_code_incorrect", "", user.ID)
					_ = c.Write(websocket.TextMessage, authIncorrect)
				}
			}
		}

		if userAllowed {
			sub := b.hub.NewSubscriber(c.Ws, user.ID, roomID)

			users, _ := b.RetroService.RetroAddUser(roomID, user.ID)
			updatedUsers, _ := json.Marshal(users)

			Retro, _ := json.Marshal(retro)
			initEvent := wshub.CreateSocketEvent("init", string(Retro), user.ID)
			_ = sub.Conn.Write(websocket.TextMessage, initEvent)

			userJoinedEvent := wshub.CreateSocketEvent("user_joined", string(updatedUsers), user.ID)
			b.hub.Broadcast(wshub.Message{Data: userJoinedEvent, Room: roomID})

			go sub.WritePump()
			go sub.ReadPump(ctx, b.hub)
		}

		return nil
	})
}

func (b *Service) RetreatUser(roomID string, userID string) string {
	users := b.RetroService.RetroRetreatUser(roomID, userID)
	updatedUsers, _ := json.Marshal(users)

	return string(updatedUsers)
}

// APIEvent handles api driven events into the retro (if active)
func (b *Service) APIEvent(ctx context.Context, retroID string, userID, eventType string, eventValue string) (any, error) {
	return b.hub.ProcessAPIEventHandler(ctx, userID, retroID, eventType, eventValue)
}
