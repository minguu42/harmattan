package api_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/minguu42/harmattan/internal/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNotFound(t *testing.T) {
	h, err := api.NewHandler(&api.Factory{}, "unknown", []string{"*"})
	require.NoError(t, err)

	req := httptest.NewRequest("GET", "/non-existent-path", nil)
	resp := httptest.NewRecorder()
	h.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code)
	assert.Equal(t, `{"code":404,"message":"指定したパスは見つかりません"}`, resp.Body.String())
}

func TestMethodNotAllowed(t *testing.T) {
	h, err := api.NewHandler(&api.Factory{}, "unknown", []string{"*"})
	require.NoError(t, err)

	t.Run("options", func(t *testing.T) {
		req := httptest.NewRequest("OPTIONS", "/health", nil)
		resp := httptest.NewRecorder()
		h.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusNoContent, resp.Code)
		assert.Equal(t, "GET", resp.Header().Get("Access-Control-Allow-Methods"))
		assert.Equal(t, "Content-Type", resp.Header().Get("Access-Control-Allow-Headers"))
	})
	t.Run("non_options", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/health", nil)
		resp := httptest.NewRecorder()
		h.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusMethodNotAllowed, resp.Code)
		assert.Equal(t, "GET", resp.Header().Get("Allow"))
		assert.Equal(t, `{"code":405,"message":"指定したメソッドは許可されていません"}`, resp.Body.String())
	})
}

func TestErrorHandler(t *testing.T) {
	h, err := api.NewHandler(&api.Factory{}, "unknown", []string{"*"})
	require.NoError(t, err)

	t.Run("authorization_error", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/projects", nil)
		resp := httptest.NewRecorder()
		h.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusUnauthorized, resp.Code)
		assert.Equal(t, `{"code":401,"message":"ユーザの認証に失敗しました"}`, resp.Body.String())
	})
	t.Run("decode_request_error", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/sign-in", strings.NewReader("invalid"))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		h.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.Equal(t, `{"code":400,"message":"リクエストに何らかの間違いがあります"}`, resp.Body.String())
	})
}

func TestCORS(t *testing.T) {
	t.Run("allowed_origin", func(t *testing.T) {
		h, err := api.NewHandler(&api.Factory{}, "unknown", []string{"http://localhost:5173"})
		require.NoError(t, err)

		req := httptest.NewRequest("GET", "/non-existent-path", nil)
		req.Header.Set("Origin", "http://localhost:5173")
		resp := httptest.NewRecorder()
		h.ServeHTTP(resp, req)

		assert.Equal(t, "http://localhost:5173", resp.Header().Get("Access-Control-Allow-Origin"))
	})
	t.Run("disallowed_origin", func(t *testing.T) {
		h, err := api.NewHandler(&api.Factory{}, "unknown", []string{"http://localhost:5173"})
		require.NoError(t, err)

		req := httptest.NewRequest("GET", "/non-existent-path", nil)
		req.Header.Set("Origin", "http://evil.example.invalid")
		resp := httptest.NewRecorder()
		h.ServeHTTP(resp, req)

		assert.Empty(t, resp.Header().Get("Access-Control-Allow-Origin"))
	})
}
