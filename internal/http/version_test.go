package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
)

func TestHandleGetAppVersion(t *testing.T) {
	s := &Service{UIConfig: thunderdome.UIConfig{AppConfig: thunderdome.AppConfig{AppVersion: "v1.2.3"}}}
	rr := httptest.NewRecorder()

	s.handleGetAppVersion().ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/api/version", nil))

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "no-store", rr.Header().Get("Cache-Control"))
	assert.JSONEq(t, `{"success":true,"error":"","data":{"version":"v1.2.3"},"meta":{}}`, rr.Body.String())
}
