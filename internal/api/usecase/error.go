package usecase

import (
	"context"
	"fmt"

	"github.com/minguu42/harmattan/internal/lib/errors"
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
	var appErr Error
	if errors.As(err, &appErr) {
		return appErr
	}

	var requestErr *ogenerrors.DecodeRequestError
	var paramsErr *ogenerrors.DecodeParamsError
	switch {
	case errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded):
		return DeadlineExceededError(err)
	case errors.As(err, &paramsErr) || errors.As(err, &requestErr):
		return ValidationError(err)
	case errors.Is(err, ogenerrors.ErrSecurityRequirementIsNotSatisfied):
		return AuthorizationError(err)
	case errors.Is(err, ogenhttp.ErrNotImplemented):
		return NotImplementedError()
	default:
		return UnknownError(err)
	}
}
