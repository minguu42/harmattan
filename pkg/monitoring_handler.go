package pkg

import "net/http"

func GetHealth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte{'a', 'b', 'c', '\n'})
	}
}
