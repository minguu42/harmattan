package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/minguu42/harmattan/internal/alog"
	middleware2 "github.com/minguu42/harmattan/internal/api/handler/middleware"
	openapi2 "github.com/minguu42/harmattan/internal/api/handler/openapi"
	usecase2 "github.com/minguu42/harmattan/internal/api/usecase"
	"github.com/minguu42/harmattan/internal/factory"
	"github.com/ogen-go/ogen/ogenerrors"
)

type handler struct {
	openapi2.UnimplementedHandler
	authentication usecase2.Authentication
	monitoring     usecase2.Monitoring
	project        usecase2.Project
	step           usecase2.Step
	tag            usecase2.Tag
	task           usecase2.Task
}

func New(f *factory.Factory, l *alog.Logger) (http.Handler, error) {
	h := handler{
		UnimplementedHandler: openapi2.UnimplementedHandler{},
		authentication:       usecase2.Authentication{Auth: f.Auth, DB: f.DB},
		monitoring:           usecase2.Monitoring{},
		project:              usecase2.Project{DB: f.DB},
		step:                 usecase2.Step{DB: f.DB},
		tag:                  usecase2.Tag{DB: f.DB},
		task:                 usecase2.Task{DB: f.DB},
	}
	sh := securityHandler{auth: f.Auth, db: f.DB}
	middlewares := []openapi2.Middleware{
		middleware2.AttachRequestIDToLogger(l),
		middleware2.AccessLog(l),
		middleware2.Recover(),
	}
	return openapi2.NewServer(&h, &sh,
		openapi2.WithNotFound(notFound),
		openapi2.WithMethodNotAllowed(methodNotAllowed),
		openapi2.WithErrorHandler(errorHandler(l)),
		openapi2.WithMiddleware(middlewares...),
	)
}

func notFound(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	_, _ = w.Write([]byte(`{"code":404,"message":"指定したパスは見つかりません"}`))
}

func methodNotAllowed(w http.ResponseWriter, r *http.Request, allowed string) {
	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Methods", allowed)
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Set("Allow", allowed)
	w.WriteHeader(http.StatusMethodNotAllowed)
	_, _ = w.Write([]byte(`{"code":405,"message":"指定したメソッドは許可されていません"}`))
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func errorHandler(l *alog.Logger) func(context.Context, http.ResponseWriter, *http.Request, error) {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request, err error) {
		// パラメータとリクエストの解析に失敗した場合にミドルウェアは実行されないので、ここでアクセスログを出力する
		var operationID string
		if requestErr := new(ogenerrors.DecodeRequestError); errors.As(err, &requestErr) {
			operationID = requestErr.OperationID()
		} else if paramsErr := new(ogenerrors.DecodeParamsError); errors.As(err, &paramsErr) {
			operationID = paramsErr.OperationID()
		}

		appErr := usecase2.ToError(err)
		l.Access(ctx, &alog.AccessFields{
			Status: appErr.Status(),
			ErrorInfo: &alog.ErrorInfo{
				ErrorMessage: appErr.Error(),
				StackTrace:   appErr.Stacktrace(),
			},
			OperationID: operationID,
			Method:      r.Method,
			URL:         r.URL.String(),
			IPAddress:   r.RemoteAddr,
		})

		w.WriteHeader(appErr.Status())
		bs, _ := json.Marshal(ErrorResponse{
			Code:    appErr.Status(),
			Message: appErr.Message(),
		})
		_, _ = w.Write(bs)
	}
}
