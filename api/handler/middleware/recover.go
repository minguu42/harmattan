package middleware

import (
	"errors"
	"fmt"
	"runtime"
	"strings"

	"github.com/minguu42/harmattan/api/apperr"
	"github.com/ogen-go/ogen/middleware"
)

func Recover() middleware.Middleware {
	return func(req middleware.Request, next middleware.Next) (resp middleware.Response, err error) {
		defer func() {
			if r := recover(); r != nil {
				message := "panic: "
				switch v := r.(type) {
				case string:
					message = message + v
				case fmt.Stringer:
					message = message + v.String()
				case error:
					message = message + v.Error()
				}

				var stacktrace []string
				for depth := 1; ; depth++ {
					pc, f, line, ok := runtime.Caller(depth)
					if !ok {
						break
					}
					// 出力するスタックトレースの量を減らすために基盤部分のスタックトレースは出力しない
					if name := f[strings.LastIndex(f, "/")+1:]; name == "oas_handlers_gen.go" {
						break
					}

					fullFuncName := runtime.FuncForPC(pc).Name()
					funcName := fullFuncName[strings.LastIndex(fullFuncName, "/")+1:]
					stacktrace = append(stacktrace, fmt.Sprintf("%s:%d %s", f, line, funcName))
				}
				err = apperr.PanicError(errors.New(message), stacktrace)
			}
		}()

		resp, err = next(req)
		return
	}
}
