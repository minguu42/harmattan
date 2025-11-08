package errors

import (
	"errors"
	"runtime"
)

const maxStackDepth = 16

func New(text string) error {
	return errors.New(text)
}

func Wrap(err error) error {
	if err == nil {
		return nil
	}

	var serr *stackError
	if errors.As(err, &serr) {
		return err
	}

	pc := make([]uintptr, maxStackDepth)
	n := runtime.Callers(2, pc)

	return &stackError{
		err:   err,
		stack: pc[:n:n],
	}
}

func Join(errs ...error) error {
	return errors.Join(errs...)
}

func Is(err, target error) bool {
	return errors.Is(err, target)
}

func As(err error, target any) bool {
	return errors.As(err, target)
}

func AsType[E error](err error) (E, bool) {
	if err == nil {
		var zero E
		return zero, false
	}

	var e E
	ok := errors.As(err, &e)
	return e, ok
}
