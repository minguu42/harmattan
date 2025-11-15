package errtrace

import (
	"errors"
	"runtime"
)

const MaxStackDepth = 8

// FromStack は err をスタックトレース stack を持つエラーでラップする
// err が nil の場合は nil を返す
// err が既にスタックトレース付きエラーの場合はラップせずそのまま返す
func FromStack(err error, stack []uintptr) error {
	if err == nil {
		return nil
	}

	if serr := new(StackError); errors.As(err, &serr) {
		return err
	}

	return &StackError{err: err, stack: stack}
}

// Wrap は err をスタックトレース付きエラーでラップする
// err が nil の場合は nil を返す
// err が既にスタックトレース付きエラーの場合はラップせずそのまま返す
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
