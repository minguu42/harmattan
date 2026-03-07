package api_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSecurityHandler(t *testing.T) {
	t.Run("valid_token", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/projects", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp := httptest.NewRecorder()
		th.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
	})
	t.Run("missing_token", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/projects", nil)
		resp := httptest.NewRecorder()
		th.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusUnauthorized, resp.Code)
		assert.Equal(t, `{"code":401,"message":"ユーザの認証に失敗しました"}`, resp.Body.String())
	})
	t.Run("invalid_token", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/projects", nil)
		req.Header.Set("Authorization", "Bearer invalid-token")
		resp := httptest.NewRecorder()
		th.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusUnauthorized, resp.Code)
		assert.Equal(t, `{"code":401,"message":"ユーザの認証に失敗しました"}`, resp.Body.String())
	})
}
