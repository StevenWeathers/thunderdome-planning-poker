package storyboard

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
	return b.hub.WebSocketHandler("storyboardId", func(w http.ResponseWriter, r *http.Request, c *wshub.Connection, roomID string) *wshub.AuthError {
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

		// make sure storyboard is legit
		storyboard, storyboardErr := b.StoryboardService.GetStoryboardByID(roomID, user.ID)
		if storyboardErr != nil {
			authErr := wshub.AuthError{
				Code:    4004,
				Message: "storyboard not found",
			}
			return &authErr
		}

		// check users storyboard active status
		userErr := b.StoryboardService.GetStoryboardUserActiveStatus(roomID, user.ID)
		userAllowed := userErr == nil
		if userErr != nil && !errors.Is(userErr, sql.ErrNoRows) {
			usrErrMsg := userErr.Error()
			var authErr wshub.AuthError

			if usrErrMsg == "DUPLICATE_STORYBOARD_USER" {
				authErr = wshub.AuthError{
					Code:    4003,
					Message: "duplicate session",
				}
			} else {
				b.logger.Ctx(ctx).Error("error finding user", zap.Error(userErr),
					zap.String("storyboard_id", roomID), zap.String("session_user_id", user.ID))
				authErr = wshub.AuthError{
					Code:    4005,
					Message: "internal error",
				}
			}
			return &authErr
		} else if storyboard.JoinCode != "" && (userErr != nil && errors.Is(userErr, sql.ErrNoRows)) {
			jcrEvent := wshub.CreateSocketEvent("join_code_required", "", user.ID)
			_ = c.Write(websocket.TextMessage, jcrEvent)

			for {
				_, msg, err := c.Ws.ReadMessage()
				if err != nil {
					if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
						b.logger.Ctx(ctx).Error("unexpected close error", zap.Error(err),
							zap.String("storyboard_id", roomID), zap.String("session_user_id", user.ID))
					}
					break
				}

				keyVal := make(map[string]string)
				err = json.Unmarshal(msg, &keyVal)
				if err != nil {
					b.logger.Error("unexpected message error", zap.Error(err),
						zap.String("retro_id", roomID), zap.String("session_user_id", user.ID))
				}

				if keyVal["type"] == "auth_storyboard" && keyVal["value"] == storyboard.JoinCode {
					// join code is valid, continue to room
					userAllowed = true
					break
				} else if keyVal["type"] == "auth_storyboard" {
					authIncorrect := wshub.CreateSocketEvent("join_code_incorrect", "", user.ID)
					_ = c.Write(websocket.TextMessage, authIncorrect)
				}
			}
		} else {
			userAllowed = true
		}

		if userAllowed {
			sub := b.hub.NewSubscriber(c.Ws, user.ID, roomID)

			users, _ := b.StoryboardService.AddUserToStoryboard(roomID, user.ID)
			updatedUsers, _ := json.Marshal(users)

			Storyboard, _ := json.Marshal(storyboard)
			initEvent := wshub.CreateSocketEvent("init", string(Storyboard), user.ID)
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
	Users := b.StoryboardService.RetreatStoryboardUser(roomID, userID)
	UpdatedUsers, _ := json.Marshal(Users)

	return string(UpdatedUsers)
}

// APIEvent handles api driven events into the storyboard (if active)
func (b *Service) APIEvent(ctx context.Context, storyboardID string, userID, eventType string, eventValue string) (any, error) {
	return b.hub.ProcessAPIEventHandler(ctx, userID, storyboardID, eventType, eventValue)
}
