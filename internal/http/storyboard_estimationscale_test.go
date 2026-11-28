package http

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
)

type stubStoryboardCreateSvc struct {
	StoryboardDataSvc
	confirmErr error
	createdIn  *[]string
}

func (s stubStoryboardCreateSvc) ConfirmStoryboardEstimationScale(_ context.Context, _ string, _ string, _ string, _ bool) error {
	return s.confirmErr
}

func (s stubStoryboardCreateSvc) CreateStoryboard(_ context.Context, _ string, _ string, _ string, _ string, estimationScaleID string, _ []*thunderdome.Color) (*thunderdome.Storyboard, error) {
	*s.createdIn = append(*s.createdIn, estimationScaleID)
	return &thunderdome.Storyboard{ID: "board"}, nil
}

func TestHandleStoryboardCreateEstimationScale(t *testing.T) {
	validate = validator.New()
	userID := "11111111-1111-1111-1111-111111111111"
	scaleID := "22222222-2222-2222-2222-222222222222"

	tests := []struct {
		name        string
		scaleID     string
		confirmErr  error
		wantStatus  int
		wantCreated []string
	}{
		{name: "without a scale", wantStatus: http.StatusOK, wantCreated: []string{""}},
		{name: "with an available scale", scaleID: scaleID, wantStatus: http.StatusOK, wantCreated: []string{scaleID}},
		{name: "unavailable scale", scaleID: scaleID, confirmErr: thunderdome.ErrStoryboardEstimationScaleNotFound, wantStatus: http.StatusBadRequest},
		{name: "scale lookup failure", scaleID: scaleID, confirmErr: errors.New("boom"), wantStatus: http.StatusInternalServerError},
		{name: "invalid scale id", scaleID: "nope", wantStatus: http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			created := make([]string, 0)
			service := &Service{
				Config:            &Config{},
				StoryboardDataSvc: stubStoryboardCreateSvc{confirmErr: tt.confirmErr, createdIn: &created},
				Logger:            otelzap.New(zap.NewNop()),
			}

			body := `{"storyboardName":"Board","estimationScaleId":"` + tt.scaleID + `"}`
			req := httptest.NewRequest(http.MethodPost, "/api/users/x/storyboards", strings.NewReader(body))
			req = req.WithContext(context.WithValue(req.Context(), contextKeyUserID, userID))
			req.SetPathValue("userId", userID)
			rr := httptest.NewRecorder()

			service.handleStoryboardCreate()(rr, req)

			assert.Equal(t, tt.wantStatus, rr.Code)
			assert.ElementsMatch(t, tt.wantCreated, created)
		})
	}
}
