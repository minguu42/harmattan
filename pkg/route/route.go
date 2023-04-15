// Package route はルーティングに関する設定を記述するパッケージである
package route

import (
	"github.com/go-chi/chi/v5"
	"github.com/minguu42/mtasks/pkg"
)

func Route(r chi.Router) {
	r.Get("/health", pkg.GetHealth())
}
