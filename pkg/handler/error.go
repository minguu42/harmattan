package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/minguu42/opepe/gen/ogen"
	ht "github.com/ogen-go/ogen/http"
	"github.com/ogen-go/ogen/ogenerrors"
	"github.com/ogen-go/ogen/validate"
)

func ErrorCode(err error) int {
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
	code := ErrorCode(err)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	var message string
	switch code {
	case http.StatusBadRequest:
		message = "入力に誤りがあります。入力を確認し、もう一度お試しください。"
	case http.StatusNotImplemented:
		message = "この機能はもうすぐ使用できます。お楽しみに♪"
	case http.StatusUnsupportedMediaType:
		message = "指定されたContent-Typeに対応していません。"
	}
	_ = json.NewEncoder(w).Encode(ogen.Error{
		Message: message,
		Debug:   err.Error(),
	})
}
