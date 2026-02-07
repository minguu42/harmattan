package handler

import (
	"context"

	"github.com/minguu42/harmattan/internal/api/openapi"
	"github.com/minguu42/harmattan/internal/lib/errtrace"
)

func (h *Handler) CheckHealth(ctx context.Context) (*openapi.CheckHealthOK, error) {
	out, err := h.Monitoring.CheckHealth(ctx)
	if err != nil {
		return nil, errtrace.Wrap(err)
	}
	return &openapi.CheckHealthOK{Revision: out.Revision}, nil
}
