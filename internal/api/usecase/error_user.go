package usecase

import (
	"errors"
	"fmt"
	"strings"
)

func DomainValidationError(errs []error) Error {
	var message string
	switch len(errs) {
	case 0:
		message = "リクエストに何らかの間違いがあります"
	case 1:
		message = errs[0].Error()
	default:
		messages := make([]string, 0, len(errs))
		for _, err := range errs {
			messages = append(messages, fmt.Sprintf("・%s", err.Error()))
		}
		message = "リクエストに以下の問題があります。\n"
		message += strings.Join(messages, "\n")
	}
	return Error{err: errors.Join(errs...), status: 400, message: message}
}

func InvalidEmailOrPasswordError(err error) Error {
	return Error{err: err, status: 400, message: "メールアドレスかパスワードに誤りがあります"}
}

func DuplicateUserEmailError(err error) Error {
	return Error{err: err, status: 409, message: "そのメールアドレスは既に使用されています"}
}

func ProjectNotFoundError(err error) Error {
	return Error{err: err, status: 404, message: "指定したプロジェクトは見つかりません"}
}

func TaskNotFoundError(err error) Error {
	return Error{err: err, status: 404, message: "指定したタスクは見つかりません"}
}

func ProjectAccessDeniedError() Error {
	return Error{status: 404, message: "指定したプロジェクトは見つかりません"}
}

func TaskAccessDeniedError() Error {
	return Error{status: 404, message: "指定したタスクは見つかりません"}
}

func StepAccessDeniedError() Error {
	return Error{status: 404, message: "指定したステップは見つかりません"}
}

func TagAccessDeniedError() Error {
	return Error{status: 404, message: "指定したタグは見つかりません"}
}

func StepNotFoundError(err error) Error {
	return Error{err: err, status: 404, message: "指定したステップは見つかりません"}
}

func TagNotFoundError(err error) Error {
	return Error{err: err, status: 404, message: "指定したタグは見つかりません"}
}
