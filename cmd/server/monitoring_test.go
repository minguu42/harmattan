package main

import (
	"net/http"
	"testing"
)

func TestMonitoring(t *testing.T) {
	run(t, []test{
		{
			id:         "getHealth",
			method:     http.MethodGet,
			path:       "/health",
			statusCode: http.StatusOK,
		},
	})
}
