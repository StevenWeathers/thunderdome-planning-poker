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
		var User *thunderdome.User

		SessionId, cookieErr := b.validateSessionCookie(w, r)
		if cookieErr != nil && cookieErr.Error() != "COOKIE_NOT_FOUND" {
			authErr := wshub.AuthError{
				Code:    4001,
				Message: "unauthorized",
			}
			return &authErr
		}

		if SessionId != "" {
			var userErr error
			User, userErr = b.AuthService.GetSessionUser(ctx, SessionId)
			if userErr != nil {
				authErr := wshub.AuthError{
					Code:    4001,
					Message: "unauthorized",
				}
				return &authErr
			}
		} else {
			UserID, err := b.validateUserCookie(w, r)
			if err != nil {
				authErr := wshub.AuthError{
					Code:    4001,
					Message: "unauthorized",
				}
				return &authErr
			}

			var userErr error
			User, userErr = b.UserService.GetGuestUser(ctx, UserID)
			if userErr != nil {
				authErr := wshub.AuthError{
					Code:    4001,
					Message: "unauthorized",
				}
				return &authErr
			}
		}

		// make sure retro is legit
		retro, retroErr := b.RetroService.RetroGet(roomID, User.Id)
		if retroErr != nil {
			authErr := wshub.AuthError{
				Code:    4004,
				Message: "retro not found",
			}
			return &authErr
		}

		// check users retro active status
		UserErr := b.RetroService.GetRetroUserActiveStatus(roomID, User.Id)
		if UserErr != nil && !errors.Is(UserErr, sql.ErrNoRows) {
			usrErrMsg := UserErr.Error()
			var authErr wshub.AuthError

			if usrErrMsg == "DUPLICATE_RETRO_USER" {
				authErr = wshub.AuthError{
					Code:    4003,
					Message: "duplicate session",
				}
			} else {
				b.logger.Ctx(ctx).Error("error finding user", zap.Error(UserErr),
					zap.String("retro_id", roomID), zap.String("session_user_id", User.Id))
				authErr = wshub.AuthError{
					Code:    4005,
					Message: "internal error",
				}
			}
			return &authErr
		} else if retro.JoinCode != "" && (UserErr != nil && errors.Is(UserErr, sql.ErrNoRows)) {
			jcrEvent := wshub.CreateSocketEvent("join_code_required", "", User.Id)
			_ = c.Write(websocket.TextMessage, jcrEvent)

			for {
				_, msg, err := c.Ws.ReadMessage()
				if err != nil {
					if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
						b.logger.Ctx(ctx).Error("unexpected close error", zap.Error(err),
							zap.String("retro_id", roomID), zap.String("session_user_id", User.Id))
					}
					break
				}

				keyVal := make(map[string]string)
				err = json.Unmarshal(msg, &keyVal)
				if err != nil {
					b.logger.Error("unexpected message error", zap.Error(err),
						zap.String("retro_id", roomID), zap.String("session_user_id", User.Id))
				}

				if keyVal["type"] == "auth_retro" && keyVal["value"] == retro.JoinCode {
					// join code is valid, continue to room
					break
				} else if keyVal["type"] == "auth_retro" {
					authIncorrect := wshub.CreateSocketEvent("join_code_incorrect", "", User.Id)
					_ = c.Write(websocket.TextMessage, authIncorrect)
				}
			}
		}

		sub := b.hub.NewSubscriber(c.Ws, User.Id, roomID)

		Users, _ := b.RetroService.RetroAddUser(roomID, User.Id)
		UpdatedUsers, _ := json.Marshal(Users)

		Retro, _ := json.Marshal(retro)
		initEvent := wshub.CreateSocketEvent("init", string(Retro), User.Id)
		_ = sub.Conn.Write(websocket.TextMessage, initEvent)

		userJoinedEvent := wshub.CreateSocketEvent("user_joined", string(UpdatedUsers), User.Id)
		b.hub.Broadcast(wshub.Message{Data: userJoinedEvent, Room: roomID})

		go sub.WritePump()
		go sub.ReadPump(ctx, b.hub)

		return nil
	})
}

func (b *Service) RetreatUser(roomID string, userID string) string {
	Users := b.RetroService.RetroRetreatUser(roomID, userID)
	UpdatedUsers, _ := json.Marshal(Users)

	return string(UpdatedUsers)
}

// APIEvent handles api driven events into the retro (if active)
func (b *Service) APIEvent(ctx context.Context, retroID string, UserID, eventType string, eventValue string) error {
	return b.hub.ProcessAPIEventHandler(ctx, UserID, retroID, eventType, eventValue)
}
