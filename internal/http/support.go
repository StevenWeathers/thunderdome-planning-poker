package http

import (
	"encoding/json"
	"io"
	"net/http"

	"go.uber.org/zap"
)

// handleGetAlerts gets a list of support tickets
// @Summary      Get Support Ticket
// @Description  get list of support tickets
// @Tags         support
// @Produce      json
// @Param        limit   query   int  false  "Max number of results to return"
// @Param        offset  query   int  false  "Starting point to return rows from, should be multiplied by limit or 0"
// @Success      200     object  standardJsonResponse{data=[]thunderdome.Support}
// @Failure      500     object  standardJsonResponse{}
// @Security     ApiKeyAuth
// @Router       /support-tickets [get]
func (s *Service) handleGetSupportTickets() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		SessionUserID, _ := ctx.Value(contextKeyUserID).(*string)
		Limit, Offset := getLimitOffsetFromRequest(r)
		tickets, Count, err := s.SupportDataSvc.GetSupportTickets(ctx, Limit, Offset)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleGetSupportTickets error", zap.Error(err),
				zap.Int("limit", Limit), zap.Int("offset", Offset),
				zap.Stringp("session_user_id", SessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		Meta := &pagination{
			Count:  Count,
			Offset: Offset,
			Limit:  Limit,
		}

		s.Success(w, r, http.StatusOK, tickets, Meta)
	}
}

type supportRequestBody struct {
	UserName     string `json:"user_name" db:"user_name"`
	UserEmail    string `json:"user_email" db:"user_email"`
	UserQuestion string `json:"user_question" db:"user_question"`
}

// handleAlertCreate creates a new support ticket
// @Summary      Create Support Ticket
// @Description  Creates a support ticket
// @Tags         support
// @Produce      json
// @Param        support  body    supportRequestBody true  "new support object"
// @Success      200    object  standardJsonResponse{data=[]thunderdome.Support}
// @Failure      500    object  standardJsonResponse{}
// @Security     ApiKeyAuth
// @Router       /support-tickets [post]
func (s *Service) handleSupportCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		SessionUserID := ctx.Value(contextKeyUserID).(string)

		var sup = supportRequestBody{}
		body, bodyErr := io.ReadAll(r.Body)
		if bodyErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
			return
		}

		jsonErr := json.Unmarshal(body, &sup)
		if jsonErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		inputErr := validate.Struct(sup)
		if inputErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, inputErr.Error()))
			return
		}

		ticket, err := s.SupportDataSvc.CreateSupportTicket(ctx, SessionUserID, sup.UserName, sup.UserEmail, sup.UserQuestion)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleSupportCreate error", zap.Error(err),
				zap.String("user_email", sup.UserEmail), zap.String("question", sup.UserQuestion),
				zap.String("session_user_id", SessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, ticket, nil)
	}
}
