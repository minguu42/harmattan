package apierror

import (
	"context"
	"errors"

	ogenhttp "github.com/ogen-go/ogen/http"
	"github.com/ogen-go/ogen/ogenerrors"
)

type Error struct {
	err     error
	status  int
	message string
}

func (e Error) Error() string {
	if e.err != nil {
		return e.err.Error()
	}
	return e.message
}

func (e Error) Status() int {
	return e.status
}

func (e Error) Message() string {
	return e.message
}

func ToError(err error) Error {
	if appErr, ok := errors.AsType[Error](err); ok {
		return appErr
	}

	var requestErr *ogenerrors.DecodeRequestError
	var paramsErr *ogenerrors.DecodeParamsError
	switch {
	case errors.Is(err, context.Canceled):
		return ClientDisconnectedError()
	case errors.Is(err, context.DeadlineExceeded):
		return DeadlineExceededError()
	case errors.As(err, &paramsErr) || errors.As(err, &requestErr):
		return ValidationError(err)
	case errors.Is(err, ogenerrors.ErrSecurityRequirementIsNotSatisfied):
		return AuthorizationError()
	case errors.Is(err, ogenhttp.ErrNotImplemented):
		return NotImplementedError()
	default:
		return UnknownError(err)
	}
}
