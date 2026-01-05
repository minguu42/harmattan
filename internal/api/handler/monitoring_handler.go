package handler

import (
	"context"

	"github.com/minguu42/harmattan/internal/api/openapi"
)

func (h *Handler) CheckHealth(_ context.Context) (*openapi.CheckHealthOK, error) {
	out := h.Monitoring.CheckHealth()
	return &openapi.CheckHealthOK{Revision: out.Revision}, nil
}
