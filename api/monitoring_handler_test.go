package api

import (
	"io"
	"net/http/httptest"
	"testing"
)

func TestGetHealth(t *testing.T) {
	r := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	getHealth()(w, r)

	resp := w.Result()
	if resp.StatusCode != 200 {
		t.Errorf("got = %d, want = 200", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("io.ReadAll failed: %v", err)
	}
	if string(body) != "" {
		t.Errorf("got = %s, want = ", body)
	}
}
