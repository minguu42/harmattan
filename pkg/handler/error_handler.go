package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

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
		Message: "指定したエンドポイントが見つかりません。",
	})
}

func MethodNotAllowed(w http.ResponseWriter, _ *http.Request, allowed string) {
	w.Header().Set("Allow", allowed)
	w.Header().Set("Content-Type", "application/json")
	code := http.StatusMethodNotAllowed
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(ogen.Error{
		Code:    code,
		Message: fmt.Sprintf("指定したエンドポイントが対応しているメソッドは%sのみです。", strings.ReplaceAll(allowed, ",", "、")),
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
	case http.StatusUnauthorized:
		message = "ユーザの認証に失敗しました。もしくはユーザが認証されていません。"
	case http.StatusBadRequest:
		message = "入力に誤りがあります。"
	case http.StatusUnsupportedMediaType:
		message = "指定されたContent-Typeに対応していません。"
	case http.StatusNotImplemented:
		message = "このオペレーションはまだ実装されていません。"
	}
	_ = json.NewEncoder(w).Encode(ogen.Error{
		Code:    code,
		Message: message,
	})
}
