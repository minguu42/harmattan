package errtrace

import (
	"errors"
	"runtime"
)

const MaxStackDepth = 8

func New(err error, stack []uintptr) error {
	return &StackError{err: err, stack: stack}
}

func Wrap(err error) error {
	if err == nil {
		return nil
	}

	var serr *StackError
	if errors.As(err, &serr) {
		return err
	}

	pc := make([]uintptr, MaxStackDepth)
	n := runtime.Callers(2, pc)

	return &StackError{err: err, stack: pc[:n:n]}
}
