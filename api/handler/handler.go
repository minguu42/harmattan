package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/minguu42/harmattan/api/apperr"
	"github.com/minguu42/harmattan/api/factory"
	"github.com/minguu42/harmattan/api/handler/middleware"
	"github.com/minguu42/harmattan/api/usecase"
	"github.com/minguu42/harmattan/internal/oapi"
	"github.com/minguu42/harmattan/lib/applog"
)

type handler struct {
	oapi.UnimplementedHandler
	authentication usecase.Authentication
	monitoring     usecase.Monitoring
	project        usecase.Project
}

func New(f *factory.Factory, l *applog.Logger) (http.Handler, error) {
	h := handler{
		UnimplementedHandler: oapi.UnimplementedHandler{},
		authentication:       usecase.Authentication{Auth: f.Auth, DB: f.DB},
		monitoring:           usecase.Monitoring{},
		project:              usecase.Project{DB: f.DB},
	}
	sh := securityHandler{auth: f.Auth, db: f.DB}
	middlewares := []oapi.Middleware{
		middleware.AttachRequestIDToLogger(l),
		middleware.AccessLog(l),
		middleware.Recover(),
	}
	return oapi.NewServer(&h, &sh,
		oapi.WithNotFound(notFound),
		oapi.WithMethodNotAllowed(methodNotAllowed),
		oapi.WithErrorHandler(errorHandler),
		oapi.WithMiddleware(middlewares...),
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

func errorHandler(_ context.Context, w http.ResponseWriter, _ *http.Request, err error) {
	appErr := apperr.ToError(err)
	w.WriteHeader(appErr.StatusCode())
	bs, _ := json.Marshal(ErrorResponse{
		Code:    appErr.StatusCode(),
		Message: appErr.MessageJapanese(),
	})
	_, _ = w.Write(bs)
}
