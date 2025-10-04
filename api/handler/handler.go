package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/minguu42/harmattan/api/factory"
	"github.com/minguu42/harmattan/api/handler/middleware"
	"github.com/minguu42/harmattan/api/handler/openapi"
	"github.com/minguu42/harmattan/api/usecase"
	"github.com/minguu42/harmattan/internal/alog"
	"github.com/ogen-go/ogen/ogenerrors"
)

type handler struct {
	openapi.UnimplementedHandler
	authentication usecase.Authentication
	monitoring     usecase.Monitoring
	project        usecase.Project
	step           usecase.Step
	tag            usecase.Tag
	task           usecase.Task
}

func New(f *factory.Factory, l *alog.Logger) (http.Handler, error) {
	h := handler{
		UnimplementedHandler: openapi.UnimplementedHandler{},
		authentication:       usecase.Authentication{Auth: f.Auth, DB: f.DB},
		monitoring:           usecase.Monitoring{},
		project:              usecase.Project{DB: f.DB},
		step:                 usecase.Step{DB: f.DB},
		tag:                  usecase.Tag{DB: f.DB},
		task:                 usecase.Task{DB: f.DB},
	}
	sh := securityHandler{auth: f.Auth, db: f.DB}
	middlewares := []openapi.Middleware{
		middleware.AttachRequestIDToLogger(l),
		middleware.AccessLog(l),
		middleware.Recover(),
	}
	return openapi.NewServer(&h, &sh,
		openapi.WithNotFound(notFound),
		openapi.WithMethodNotAllowed(methodNotAllowed),
		openapi.WithErrorHandler(errorHandler(l)),
		openapi.WithMiddleware(middlewares...),
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

		appErr := usecase.ToError(err)
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
