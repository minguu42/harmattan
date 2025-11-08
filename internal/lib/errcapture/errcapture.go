package errcapture

import (
	"errors"
	"log"
	"os"
	"runtime"
)

// Do はfを実行し、エラーがあればerrが示すエラーに追加する。
// errがnilの場合はfを実行して終了する。
func Do(err *error, f func() error) {
	ferr := f()
	if err == nil || ferr == nil {
		return
	}

	// os.ErrClosedは2重にクローズ処理を行った場合に返され、処理的には問題がないため無視する。
	if errors.Is(ferr, os.ErrClosed) {
		return
	}

	*err = errors.Join(*err, ferr)
}

// Fatal はfを実行し、エラーがあればログに出力し、os.Exit(1)で終了する。
func Fatal(f func() error) {
	err := f()
	if err == nil {
		return
	}

	// os.ErrClosedは2重にクローズ処理を行った場合に返され、処理的には問題がないため無視する。
	if errors.Is(err, os.ErrClosed) {
		return
	}

	_, file, line, _ := runtime.Caller(1)
	log.Fatalf("%s:%d: %v", file, line, err)
}
