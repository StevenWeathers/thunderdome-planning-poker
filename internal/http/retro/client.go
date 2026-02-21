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
func (s *Service) ServeWs() http.HandlerFunc {
	return s.hub.WebSocketHandler("retroId", func(w http.ResponseWriter, r *http.Request, c *wshub.Connection, roomID string) *wshub.AuthError {
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

		// make sure retro is legit
		retro, retroErr := s.RetroService.RetroGetByID(roomID, user.ID)
		if retroErr != nil {
			authErr := wshub.AuthError{
				Code:    4004,
				Message: "retro not found",
			}
			return &authErr
		}

		// check users retro active status
		userErr := s.RetroService.GetRetroUserActiveStatus(roomID, user.ID)
		userAllowed := userErr == nil || (errors.Is(userErr, sql.ErrNoRows) && retro.JoinCode == "")
		if userErr != nil && !errors.Is(userErr, sql.ErrNoRows) {
			usrErrMsg := userErr.Error()
			var authErr wshub.AuthError

			if usrErrMsg == "DUPLICATE_RETRO_USER" {
				authErr = wshub.AuthError{
					Code:    4003,
					Message: "duplicate session",
				}
			} else {
				s.logger.Ctx(ctx).Error("error finding user", zap.Error(userErr),
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
						s.logger.Ctx(ctx).Error("unexpected close error", zap.Error(err),
							zap.String("retro_id", roomID), zap.String("session_user_id", user.ID))
					}
					break
				}

				keyVal := make(map[string]string)
				err = json.Unmarshal(msg, &keyVal)
				if err != nil {
					s.logger.Error("unexpected message error", zap.Error(err),
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
			sub := s.hub.NewSubscriber(c.Ws, user.ID, roomID)

			users, _ := s.RetroService.RetroAddUser(roomID, user.ID)
			updatedUsers, _ := json.Marshal(users)

			Retro, _ := json.Marshal(retro)
			initEvent := wshub.CreateSocketEvent("init", string(Retro), user.ID)
			_ = sub.Conn.Write(websocket.TextMessage, initEvent)

			userJoinedEvent := wshub.CreateSocketEvent("user_joined", string(updatedUsers), user.ID)
			s.hub.Broadcast(wshub.Message{Data: userJoinedEvent, Room: roomID})

			go sub.WritePump()
			go sub.ReadPump(ctx, s.hub)
		}

		return nil
	})
}

func (s *Service) RetreatUser(roomID string, userID string) string {
	users := s.RetroService.RetroRetreatUser(roomID, userID)
	updatedUsers, _ := json.Marshal(users)

	return string(updatedUsers)
}

// APIEvent handles api driven events into the retro (if active)
func (s *Service) APIEvent(ctx context.Context, retroID string, userID, eventType string, eventValue string) (any, error) {
	return s.hub.ProcessAPIEventHandler(ctx, userID, retroID, eventType, eventValue)
}
