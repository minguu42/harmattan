package errtrace

import (
	"errors"
	"runtime"
)

const maxStackDepth = 8

func Wrap(err error) error {
	if err == nil {
		return nil
	}

	var serr *StackError
	if errors.As(err, &serr) {
		return err
	}

	pc := make([]uintptr, maxStackDepth)
	n := runtime.Callers(2, pc)

	return &StackError{err: err, stack: pc[:n:n]}
}
