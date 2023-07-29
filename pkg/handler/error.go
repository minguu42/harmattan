package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/minguu42/opepe/gen/ogen"
	ht "github.com/ogen-go/ogen/http"
	"github.com/ogen-go/ogen/ogenerrors"
	"github.com/ogen-go/ogen/validate"
)

func NotFound(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	code := http.StatusNotFound
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(ogen.Error{
		Code:    code,
		Message: "指定したパスに対応するエンドポイントが見つかりません。もう一度ご確認ください。",
	})
}

func MethodNotAllowed(w http.ResponseWriter, _ *http.Request, allowed string) {
	w.Header().Set("Allow", allowed)
	w.Header().Set("Content-Type", "text/plain")
	code := http.StatusMethodNotAllowed
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(ogen.Error{
		Code:    code,
		Message: fmt.Sprintf("このパスに対応しているメソッドは%sのみである", allowed),
	})
}

func StatusCode(err error) int {
	var (
		code    = http.StatusInternalServerError
		ctError *validate.InvalidContentTypeError
		ogenErr ogenerrors.Error
	)
	switch {
	case errors.Is(err, ht.ErrNotImplemented):
		code = http.StatusNotImplemented
	case errors.As(err, &ctError):
		code = http.StatusUnsupportedMediaType
	case errors.As(err, &ogenErr):
		code = ogenErr.Code()
	}
	return code
}

func ErrorHandler(_ context.Context, w http.ResponseWriter, _ *http.Request, err error) {
	code := StatusCode(err)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	var message string
	switch code {
	case http.StatusBadRequest:
		message = "入力に誤りがあります。入力を確認し、もう一度お試しください。"
	case http.StatusUnsupportedMediaType:
		message = "指定されたContent-Typeに対応していません。"
	case http.StatusNotImplemented:
		message = "この機能はもうすぐ使用できます。お楽しみに♪"
	}
	_ = json.NewEncoder(w).Encode(ogen.Error{
		Code:    code,
		Message: message,
	})
}
