package app

import (
	"io"
	"net/http/httptest"
	"testing"
)

func TestGetHealth(t *testing.T) {
	r := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	GetHealth()(w, r)

	resp := w.Result()

	if resp.StatusCode != 200 {
		t.Errorf("got = %d, want = 200", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("cannot read test response: %v", err)
	}
	if string(body) != "" {
		t.Errorf("got = %s, want = ", body)
	}
}
