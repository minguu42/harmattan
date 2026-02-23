package errtrace

import (
	"errors"
	"log/slog"
	"runtime"
)

const MaxStackDepth = 8

// FromStack は err をスタックトレース stack を持つエラーでラップする
// err が nil の場合は nil を返す
// err が既にスタックトレース付きエラーの場合はラップせずそのまま返す
// ただし attrs が指定されている場合は既存のエラーの属性に追加する
func FromStack(err error, stack []uintptr, attrs ...slog.Attr) error {
	if err == nil {
		return nil
	}

	if serr, ok := errors.AsType[*StackError](err); ok {
		if len(attrs) > 0 {
			serr.attrs = append(serr.attrs, attrs...)
		}
		return err
	}

	return &StackError{err: err, stack: stack, attrs: attrs}
}

// Wrap は err をスタックトレース付きエラーでラップする
// err が nil の場合は nil を返す
// err が既にスタックトレース付きエラーの場合はラップせずそのまま返す
// ただし attrs が指定されている場合は既存のエラーの属性に追加する
func Wrap(err error, attrs ...slog.Attr) error {
	if err == nil {
		return nil
	}

	if serr, ok := errors.AsType[*StackError](err); ok {
		if len(attrs) > 0 {
			serr.attrs = append(serr.attrs, attrs...)
		}
		return err
	}

	pc := make([]uintptr, MaxStackDepth)
	n := runtime.Callers(2, pc)

	return &StackError{err: err, stack: pc[:n:n], attrs: attrs}
}
