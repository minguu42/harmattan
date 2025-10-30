package errorsx

import (
	"errors"
	"os"
)

func Capture(err *error, f func() error) {
	ferr := f()
	if err == nil {
		return
	}

	// 2重にクローズ処理を行なっても処理的には問題がないため無視する
	if errors.Is(ferr, os.ErrClosed) {
		return
	}
	*err = errors.Join(*err, ferr)
}
