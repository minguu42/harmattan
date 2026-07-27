package apierror

import (
	"errors"
	"fmt"
	"strings"

	"github.com/minguu42/harmattan/internal/domain"
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

func InvalidEmailOrPasswordError() Error {
	return Error{status: 400, message: "メールアドレスかパスワードに誤りがあります"}
}

func DuplicateUserEmailError() Error {
	return Error{status: 409, message: "そのメールアドレスは既に使用されています"}
}

func ProjectNotFoundError() Error {
	return Error{status: 404, message: "指定したプロジェクトは見つかりません"}
}

func TooManyProjectsError() Error {
	return Error{status: 409, message: fmt.Sprintf("作成できるプロジェクトは%d件までです。不要なプロジェクトを削除してから再度お試しください", domain.MaxProjectsPerUser)}
}

func TooManyTasksError() Error {
	return Error{status: 409, message: fmt.Sprintf("1つのプロジェクトに作成できるタスクは%d件までです。不要なタスクを削除してから再度お試しください", domain.MaxTasksPerProject)}
}

func TooManyStepsError() Error {
	return Error{status: 409, message: fmt.Sprintf("1つのタスクに作成できるステップは%d件までです。不要なステップを削除してから再度お試しください", domain.MaxStepsPerTask)}
}

func TooManyTagsError() Error {
	return Error{status: 409, message: fmt.Sprintf("作成できるタグは%d件までです。不要なタグを削除してから再度お試しください", domain.MaxTagsPerUser)}
}

func TaskNotFoundError() Error {
	return Error{status: 404, message: "指定したタスクは見つかりません"}
}

func StepNotFoundError() Error {
	return Error{status: 404, message: "指定したステップは見つかりません"}
}

func TagNotFoundError() Error {
	return Error{status: 404, message: "指定したタグは見つかりません"}
}
