package handler

import (
	"context"
	"runtime/debug"
	"slices"

	"github.com/minguu42/opepe/gen/ogen"
)

// version の値はビルド時に埋め込まれる
var version = "v0.0.0+unknown"

// GetHealth は GET /health に対応するハンドラ関数
func (h *Handler) GetHealth(_ context.Context) (*ogen.GetHealthOK, error) {
	revision := "xxxxxxx"
	if info, ok := debug.ReadBuildInfo(); ok {
		if i := slices.IndexFunc(info.Settings, func(s debug.BuildSetting) bool {
			return s.Key == "vcs.revision"
		}); i != -1 {
			revision = info.Settings[i].Value[:len(revision)]
		}
	}

	return &ogen.GetHealthOK{
		Version:  version,
		Revision: revision,
	}, nil
}
