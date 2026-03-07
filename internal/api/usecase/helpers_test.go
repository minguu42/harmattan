package usecase_test

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type test struct {
	Method  string
	Path    string
	Headers http.Header
	Body    string

	WantStatus  int
	WantHeaders http.Header
	WantJSON    any
	WantTables  []any
}

func runTest(t *testing.T, tt test) {
	t.Helper()

	req, err := http.NewRequest(tt.Method, ts.URL+tt.Path, strings.NewReader(tt.Body))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	for k, v := range tt.Headers {
		req.Header[k] = v
	}
	if tt.Body != "" {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := ts.Client().Do(req)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			t.Fatalf("Failed to close response body: %v", err)
		}
	}()

	if tt.WantStatus != 0 && resp.StatusCode != tt.WantStatus {
		t.Errorf("Status code: got %d, want %d", resp.StatusCode, tt.WantStatus)
	}
	for k := range tt.WantHeaders {
		if got, want := resp.Header.Get(k), tt.WantHeaders.Get(k); got != want {
			t.Errorf("Response Header %q: got %q, want %q", k, got, want)
		}
	}
	if tt.WantJSON != nil {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Failed to read response body: %v", err)
		}
		var got, want any
		if err := json.Unmarshal(body, &got); err != nil {
			t.Fatalf("Failed to unmarshal response body: %v", err)
		}
		wantBytes, err := json.Marshal(tt.WantJSON)
		if err != nil {
			t.Fatalf("Failed to marshal want: %v", err)
		}
		if err := json.Unmarshal(wantBytes, &want); err != nil {
			t.Fatalf("Failed to unmarshal want: %v", err)
		}
		if diff := cmp.Diff(got, want); diff != "" {
			t.Errorf("Response body mismatch (-got +want):\n%s", diff)
		}
	}

	if len(tt.WantTables) > 0 {
		tdb.Assert(t, tt.WantTables)
	}
}

func doRequest(t *testing.T, method, path, body string) *http.Response {
	t.Helper()

	req, err := http.NewRequest(method, ts.URL+path, strings.NewReader(body))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := ts.Client().Do(req)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	return resp
}

func assertJSONEqual(t *testing.T, resp *http.Response, want any) {
	t.Helper()

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var got any
	require.NoError(t, json.Unmarshal(body, &got))

	wantBytes, err := json.Marshal(want)
	require.NoError(t, err)

	var wantParsed any
	require.NoError(t, json.Unmarshal(wantBytes, &wantParsed))
	assert.Equal(t, wantParsed, got)
}
