package api

import "net/http"

func getHealth() http.HandlerFunc {
	return func(_ http.ResponseWriter, _ *http.Request) {}
}
