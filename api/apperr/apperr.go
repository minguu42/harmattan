package apperr

import (
	"fmt"
	"net/http"

	"github.com/minguu42/harmattan/internal/oapi"
)

type Error struct {
	err             error
	id              string
	code            int
	message         string
	messageJapanese string
}

func (e Error) Error() string {
	if e.err != nil {
		return fmt.Sprintf("%s: %s", e.message, e.err)
	}
	return e.message
}

func (e Error) APIError() *oapi.ErrorStatusCode {
	return &oapi.ErrorStatusCode{
		StatusCode: e.code,
		Response: oapi.Error{
			Code:    e.code,
			Message: e.messageJapanese,
		},
	}
}

func ErrDeadlineExceeded(err error) Error {
	return Error{
		err:             err,
		id:              "deadline-exceeded",
		code:            http.StatusGatewayTimeout,
		message:         "request was not processed within the specified time",
		messageJapanese: "リクエストは規定時間内に処理されませんでした",
	}
}

func ErrDuplicateUserEmail(err error) Error {
	return Error{
		err:             err,
		id:              "duplicate-user-email",
		code:            http.StatusConflict,
		message:         "the mail address is already in use",
		messageJapanese: "そのメールアドレスは既に使用されています",
	}
}

func ErrUnknown(err error) Error {
	return Error{
		err:             err,
		id:              "unknown",
		code:            http.StatusInternalServerError,
		message:         "some error has occurred on the server side. please wait a few minutes and try again",
		messageJapanese: "サーバ側で何らかのエラーが発生しました。時間を置いてから再度お試しください。",
	}
}
