package handler_test

import (
	"testing"

	"github.com/ikawaha/httpcheck"
	"github.com/minguu42/harmattan/api/handler"
)

func TestHandler_NotFound(t *testing.T) {
	httpcheck.New(th).Test(t, "GET", "/non-existent-path").
		Check().
		HasStatus(404).
		HasJSON(handler.ErrorResponse{Code: 404, Message: "指定したパスは見つかりません"})
}

func TestHandler_MethodNotFound(t *testing.T) {
	t.Run("method not allowed", func(t *testing.T) {
		httpcheck.New(th).Test(t, "POST", "/health").
			Check().
			HasStatus(405).
			HasJSON(handler.ErrorResponse{Code: 405, Message: "指定したメソッドは許可されていません"})
	})
	t.Run("options", func(t *testing.T) {
		httpcheck.New(th).Test(t, "OPTIONS", "/health").
			Check().
			HasStatus(204).
			HasHeader("Access-Control-Allow-Methods", "GET").
			HasHeader("Access-Control-Allow-Headers", "Content-Type")
	})
}
