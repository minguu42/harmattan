// Package render はリクエスト、レスポンスの処理を手伝うヘルパー関数を含む
package render

import (
	"encoding/json"
	"net/http"
)

// Response はレスポンスを生成する
func Response(w http.ResponseWriter, statusCode int, v any) error {
	w.WriteHeader(statusCode)
	if statusCode == http.StatusNoContent {
		return nil
	}

	w.Header().Set("Content-Type", "application/json")

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(v)
}

// Error はエラーレスポンスを生成する
func Error(w http.ResponseWriter, statusCode int, err error) {
	w.WriteHeader(statusCode)

	w.Header().Set("Content-Type", "application/json")

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	switch statusCode {
	case http.StatusBadRequest:
		_ = encoder.Encode(newBadRequest(err))
	case http.StatusUnauthorized:
		_ = encoder.Encode(newUnauthorized(err))
	case http.StatusNotFound:
		_ = encoder.Encode(newNotFound(err))
	case http.StatusInternalServerError:
		_ = encoder.Encode(newInternalServerError(err))
	}
}
