package app

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/minguu42/mtasks/app/ogen"
)

type authMiddleware struct {
	next       *ogen.Server
	repository repository
}

type userKey struct{}

func (a *authMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	route, ok := a.next.FindRoute(r.Method, r.URL.Path)
	if !ok || route.OperationID() == "getHealth" {
		a.next.ServeHTTP(w, r)
		return
	}

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")

	apiKey := r.Header.Get("X-Api-Key")
	if apiKey == "" {
		w.WriteHeader(http.StatusUnauthorized)
		_ = encoder.Encode(ogen.Error{
			Message: "ユーザの認証に失敗しました。",
			Debug:   "X-Api-Key is required",
		})
		return
	}

	ctx := r.Context()
	u, err := a.repository.GetUserByAPIKey(ctx, apiKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = encoder.Encode(ogen.Error{
			Message: "サーバで何らかのエラーが発生しました。もう一度お試しください。",
			Debug:   fmt.Sprintf("repository.GetUserByAPIKey failed: %v", err),
		})
		return
	}

	ctx = context.WithValue(ctx, userKey{}, u)
	a.next.ServeHTTP(w, r.WithContext(ctx))
}
