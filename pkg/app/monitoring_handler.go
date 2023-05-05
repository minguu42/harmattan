package app

import "net/http"

func GetHealth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}
