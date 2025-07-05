package apperr

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	ogenhttp "github.com/ogen-go/ogen/http"
	"github.com/ogen-go/ogen/ogenerrors"
)

type Error struct {
	err             error
	stacktrace      []string
	id              string
	statusCode      int
	message         string
	messageJapanese string
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

func (e Error) StatusCode() int {
	return e.statusCode
}

func (e Error) MessageJapanese() string {
	return e.messageJapanese
}

func ToError(err error) Error {
	var appErr Error
	switch {
	case errors.As(err, &appErr):
	case errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded):
		appErr = ErrDeadlineExceeded(err)
	case errors.Is(err, ogenerrors.ErrSecurityRequirementIsNotSatisfied):
		appErr = ErrAuthorization(err)
	case errors.Is(err, ogenhttp.ErrNotImplemented):
		appErr = ErrNotImplemented()
	default:
		appErr = ErrUnknown(err)
	}
	return appErr
}

func ErrUnknown(err error) Error {
	return Error{
		err:             err,
		id:              "unknown",
		statusCode:      http.StatusInternalServerError,
		message:         "some error has occurred on the server side. please wait a few minutes and try again",
		messageJapanese: "サーバ側で何らかのエラーが発生しました。時間を置いてから再度お試しください",
	}
}

func ErrPanic(err error, stacktrace []string) Error {
	return Error{
		err:             err,
		stacktrace:      stacktrace,
		id:              "panic",
		statusCode:      http.StatusInternalServerError,
		message:         "some error has occurred on the server side. please wait a few minutes and try again",
		messageJapanese: "サーバ側で何らかのエラーが発生しました。時間を置いてから再度お試しください",
	}
}

func ErrNotImplemented() Error {
	return Error{
		id:              "not-implemented",
		statusCode:      http.StatusNotImplemented,
		message:         "this feature has not been implemented yet",
		messageJapanese: "この機能はまだ実装されていません",
	}
}

func ErrDeadlineExceeded(err error) Error {
	return Error{
		err:             err,
		id:              "deadline-exceeded",
		statusCode:      http.StatusGatewayTimeout,
		message:         "request was not processed within the specified time",
		messageJapanese: "リクエストは規定時間内に処理されませんでした",
	}
}

func ErrAuthorization(err error) Error {
	return Error{
		err:             err,
		id:              "authorization",
		statusCode:      http.StatusUnauthorized,
		message:         "user authentication failed",
		messageJapanese: "ユーザの認証に失敗しました",
	}
}

func ErrDuplicateUserEmail(err error) Error {
	return Error{
		err:             err,
		id:              "duplicate-user-email",
		statusCode:      http.StatusConflict,
		message:         "the mail address is already in use",
		messageJapanese: "そのメールアドレスは既に使用されています",
	}
}

func ErrProjectNotFound(err error) Error {
	return Error{
		err:             err,
		id:              "project-not-found",
		statusCode:      http.StatusNotFound,
		message:         "the specified project is not found",
		messageJapanese: "指定したプロジェクトは見つかりません",
	}
}
