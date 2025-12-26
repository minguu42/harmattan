package middleware

import (
	"errors"
	"fmt"
	"runtime"

	"github.com/minguu42/harmattan/internal/api/usecase"
	"github.com/minguu42/harmattan/internal/lib/errtrace"
	"github.com/ogen-go/ogen/middleware"
)

func Recover() middleware.Middleware {
	return func(req middleware.Request, next middleware.Next) (resp middleware.Response, err error) {
		defer func() {
			if r := recover(); r != nil {
				message := "panic: "
				switch v := r.(type) {
				case string:
					message += v
				case fmt.Stringer:
					message += v.String()
				case error:
					message += v.Error()
				default:
					message += fmt.Sprintf("%v", v)
				}

				pc := make([]uintptr, errtrace.MaxStackDepth)
				n := runtime.Callers(2, pc)

				err = errtrace.FromStack(usecase.UnknownError(errors.New(message)), pc[:n:n])
			}
		}()

		resp, err = next(req)
		return
	}
}
