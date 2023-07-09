package main

import (
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/minguu42/mtasks/gen/ogen"
)

func TestGetHealth(t *testing.T) {
	tests := []test{
		{
			name: "サーバの状態を取得する",
			request: request{
				method: http.MethodGet,
				path:   "/health",
			},
			response: response{
				statusCode: http.StatusOK,
				body: ogen.GetHealthOK{
					Version:  "",
					Revision: "",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got ogen.GetHealthOK
			resp, err := doTestRequest(tt.request, &got)
			if err != nil {
				t.Fatalf("doTestRequest failed: %s", err)
			}

			if tt.response.statusCode != resp.StatusCode {
				t.Errorf("status code want %d, but %d", tt.response.statusCode, resp.StatusCode)
			}
			if diff := cmp.Diff(tt.response.body, got); diff != "" {
				t.Errorf("response body mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
