package errtrace

import (
	"errors"
	"runtime"
)

const MaxStackDepth = 8

func New(err error, stack []uintptr) error {
	if err == nil {
		return nil
	}

	if serr := new(StackError); errors.As(err, &serr) {
		return err
	}

	return &StackError{err: err, stack: stack}
}

func Wrap(err error) error {
	if err == nil {
		return nil
	}

	if serr := new(StackError); errors.As(err, &serr) {
		return err
	}

	pc := make([]uintptr, MaxStackDepth)
	n := runtime.Callers(2, pc)

	return &StackError{err: err, stack: pc[:n:n]}
}
