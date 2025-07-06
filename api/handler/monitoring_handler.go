package handler

import (
	"context"

	"github.com/minguu42/harmattan/internal/openapi"
)

func (h *handler) CheckHealth(_ context.Context) (*openapi.CheckHealthOK, error) {
	out := h.monitoring.CheckHealth()
	return &openapi.CheckHealthOK{Revision: out.Revision}, nil
}
