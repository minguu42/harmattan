package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/minguu42/harmattan/api/apperr"
	"github.com/minguu42/harmattan/api/factory"
	"github.com/minguu42/harmattan/api/usecase"
	"github.com/minguu42/harmattan/internal/oapi"
	"github.com/minguu42/harmattan/lib/applog"
)

type handler struct {
	authentication usecase.Authentication
	monitoring     usecase.Monitoring
}

func New(f *factory.Factory, _ *applog.Logger) (http.Handler, error) {
	return oapi.NewServer(&handler{
		authentication: usecase.NewAuthentication(f.Auth, f.DB),
		monitoring:     usecase.Monitoring{},
	})
}

func (h *handler) NewError(_ context.Context, err error) *oapi.ErrorStatusCode {
	var appErr apperr.Error
	switch {
	case errors.As(err, &appErr):
	case errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded):
		appErr = apperr.ErrDeadlineExceeded(err)
	default:
		appErr = apperr.ErrUnknown(err)
	}
	return appErr.APIError()
}
