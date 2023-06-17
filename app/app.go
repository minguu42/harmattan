// Package app はアプリケーションの中心に関するパッケージ
package app

import (
	"context"
	"errors"
	"github.com/minguu42/mtasks/app/ogen"
	"github.com/minguu42/mtasks/app/repository"
)

type Handler struct {
	Repository repository.Repository
}

var (
	errBadRequest          = errors.New("there is an input error")
	errUnauthorized        = errors.New("user is not authenticated")
	errTaskNotFound        = errors.New("the specified task is not found")
	errProjectNotFound     = errors.New("the specified project is not found")
	errInternalServerError = errors.New("some error occurred on the server")
	errNotImplemented      = errors.New("this API operation is not yet implemented")
	errServerUnavailable   = errors.New("server is temporarily unavailable")
)

// NewError -
func (h *Handler) NewError(_ context.Context, err error) *ogen.ErrorStatusCode {
	var (
		statusCode int
		message    string
	)
	switch err {
	case errBadRequest:
		statusCode = 400
		message = "入力に誤りがあります。入力をご確認ください。"
	case errUnauthorized:
		statusCode = 401
		message = "ユーザが認証されていません。ユーザの認証後にもう一度お試しください。"
	case errProjectNotFound:
		statusCode = 404
		message = "指定されたプロジェクトが見つかりません。もう一度ご確認ください。"
	case errTaskNotFound:
		statusCode = 404
		message = "指定されたタスクが見つかりません。もう一度ご確認ください。"
	case errInternalServerError:
		statusCode = 500
		message = "不明なエラーが発生しました。もう一度お試しください。"
	case errNotImplemented:
		statusCode = 501
		message = "この機能はもうすぐ使用できます。お楽しみに♪"
	case errServerUnavailable:
		statusCode = 503
		message = "サーバが一時的に利用できない状態です。時間を空けてから、もう一度お試しください。"
	default:
		statusCode = 500
		message = "不明なエラーが発生しました。もう一度お試しください。"
		err = errInternalServerError
	}

	return &ogen.ErrorStatusCode{
		StatusCode: statusCode,
		Response: ogen.Error{
			Message: message,
			Debug:   err.Error(),
		},
	}
}
