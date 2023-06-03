package app

import (
	"context"

	"github.com/minguu42/mtasks/app/ogen"
)

// GetHealth は GET /health に対応するハンドラ関数
func (h *handler) GetHealth(_ context.Context) (ogen.GetHealthRes, error) {
	return &ogen.GetHealthOK{}, nil
}
