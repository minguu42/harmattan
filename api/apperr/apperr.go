package apperr

import (
	"context"
	"errors"
	"fmt"

	ogenhttp "github.com/ogen-go/ogen/http"
	"github.com/ogen-go/ogen/ogenerrors"
)

type Error struct {
	err        error
	stacktrace []string
	status     int
	message    string
}

func (e Error) Error() string {
	if e.err != nil {
		return fmt.Sprintf("%s: %s", e.message, e.err)
	}
	return e.message
}

func (e Error) Stacktrace() []string {
	return e.stacktrace
}

func (e Error) Status() int {
	return e.status
}

func (e Error) Message() string {
	return e.message
}

func ToError(err error) Error {
	var requestErr *ogenerrors.DecodeRequestError
	var paramsErr *ogenerrors.DecodeParamsError

	var appErr Error
	switch {
	case errors.As(err, &appErr):
	case errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded):
		appErr = DeadlineExceededError(err)
	case errors.As(err, &paramsErr) || errors.As(err, &requestErr):
		appErr = ValidationError(err)
	case errors.Is(err, ogenerrors.ErrSecurityRequirementIsNotSatisfied):
		appErr = AuthorizationError(err)
	case errors.Is(err, ogenhttp.ErrNotImplemented):
		appErr = NotImplementedError()
	default:
		appErr = UnknownError(err)
	}
	return appErr
}

func ValidationError(err error) Error {
	return Error{err: err, status: 400, message: "リクエストに何らかの間違いがあります"}
}

func UnknownError(err error) Error {
	return Error{err: err, status: 500, message: "サーバ側で何らかのエラーが発生しました。時間を置いてから再度お試しください"}
}

func PanicError(err error, stacktrace []string) Error {
	return Error{
		err:        err,
		stacktrace: stacktrace,
		status:     500,
		message:    "サーバ側で何らかのエラーが発生しました。時間を置いてから再度お試しください",
	}
}

func NotImplementedError() Error {
	return Error{status: 501, message: "この機能はまだ実装されていません"}
}

func DeadlineExceededError(err error) Error {
	return Error{err: err, status: 504, message: "リクエストは規定時間内に処理されませんでした"}
}

func AuthorizationError(err error) Error {
	return Error{err: err, status: 401, message: "ユーザの認証に失敗しました"}
}
