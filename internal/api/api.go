package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/minguu42/harmattan/internal/api/handler"
	"github.com/minguu42/harmattan/internal/api/middleware"
	"github.com/minguu42/harmattan/internal/api/openapi"
	"github.com/minguu42/harmattan/internal/api/usecase"
	"github.com/minguu42/harmattan/internal/atel"
	"github.com/minguu42/harmattan/internal/lib/errtrace"
	"github.com/ogen-go/ogen/ogenerrors"
	"github.com/rs/cors"
)

//go:generate go tool ogen -clean -config ../../.ogen.yaml -package openapi -target ./openapi ../../doc/openapi.yaml

func NewHandler(f *Factory, revision string, allowedOrigins []string) (http.Handler, error) {
	h := &handler.Handler{
		UnimplementedHandler: openapi.UnimplementedHandler{},
		Authentication:       usecase.Authentication{Auth: f.Auth, DB: f.DB},
		Monitoring:           usecase.Monitoring{Revision: revision, DB: f.DB},
		Project:              usecase.Project{DB: f.DB},
		Step:                 usecase.Step{DB: f.DB},
		Tag:                  usecase.Tag{DB: f.DB},
		Task:                 usecase.Task{DB: f.DB},
	}

	sh := securityHandler{auth: f.Auth, db: f.DB}
	middlewares := []openapi.Middleware{
		middleware.AttachTraceIDToLogger(),
		middleware.AccessLog(),
		middleware.Recover(),
	}
	ogenServer, err := openapi.NewServer(h, &sh,
		openapi.WithNotFound(notFound),
		openapi.WithMethodNotAllowed(methodNotAllowed),
		openapi.WithErrorHandler(errorHandler),
		openapi.WithMiddleware(middlewares...),
	)
	if err != nil {
		return nil, errtrace.Wrap(err)
	}

	corsSetting := cors.New(cors.Options{
		AllowedOrigins: allowedOrigins,
		AllowedMethods: []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Authorization", "Content-Type"},
	})
	return corsSetting.Handler(ogenServer), nil
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

func errorHandler(ctx context.Context, w http.ResponseWriter, r *http.Request, err error) {
	appErr := usecase.ToError(err)

	// パラメータとリクエストの解析に失敗した場合にミドルウェアは実行されないので、ここでアクセスログを出力する
	// 上記以外の場合のアクセスログは middleware.AccessLog で出力される
	var operationID string
	if requestErr, ok := errors.AsType[*ogenerrors.DecodeRequestError](err); ok {
		operationID = requestErr.OperationID()
	} else if paramsErr, ok := errors.AsType[*ogenerrors.DecodeParamsError](err); ok {
		operationID = paramsErr.OperationID()
	}
	if operationID != "" {
		atel.AccessLog(ctx, &atel.AccessFields{
			Status:      appErr.Status(),
			OperationID: operationID,
			Method:      r.Method,
			URL:         r.URL.String(),
			IPAddress:   r.RemoteAddr,
			UserAgent:   r.UserAgent(),
		})
	}

	w.WriteHeader(appErr.Status())
	bs, _ := json.Marshal(ErrorResponse{
		Code:    appErr.Status(),
		Message: appErr.Message(),
	})
	_, _ = w.Write(bs)
}
