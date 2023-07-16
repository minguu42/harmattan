package main

import (
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/minguu42/mtasks/gen/mock"
	"github.com/minguu42/mtasks/gen/ogen"
)

func TestCreateProject(t *testing.T) {
	tests := []test{
		{
			name: "プロジェクトを作成する",
			request: request{
				method: http.MethodPost,
				path:   "/projects",
				body:   strings.NewReader(`{"name": "新プロジェクト", "color": "#1A2B3C"}`),
			},
			response: response{
				statusCode: http.StatusCreated,
				body: ogen.Project{
					ID:         "01DXF6DT000000000000000002",
					Name:       "新プロジェクト",
					Color:      "#1A2B3C",
					IsArchived: false,
					CreatedAt:  time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
					UpdatedAt:  time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := mock.NewMockIDGenerator(gomock.NewController(t))
			g.EXPECT().Generate().Return("01DXF6DT000000000000000002")
			tdb.SetIDGenerator(g)

			if err := tdb.Begin(); err != nil {
				t.Fatalf("tdb.Begin failed: %s", err)
			}
			defer tdb.Rollback()

			var got ogen.Project
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

func TestListProject(t *testing.T) {
	tests := []test{
		{
			name: "プロジェクト一覧を取得する",
			request: request{
				method: http.MethodGet,
				path:   "/projects?sort=createdAt",
			},
			response: response{
				statusCode: http.StatusOK,
				body: ogen.Projects{
					Projects: []ogen.Project{
						{
							ID:         "01DXF6DT000000000000000000",
							Name:       "プロジェクト1",
							Color:      "#1A2B3C",
							IsArchived: false,
							CreatedAt:  time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
							UpdatedAt:  time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
						},
						{
							ID:         "01DXF6DT000000000000000001",
							Name:       "プロジェクト2",
							Color:      "#1A2B3C",
							IsArchived: false,
							CreatedAt:  time.Date(2020, 1, 1, 0, 0, 1, 0, time.UTC),
							UpdatedAt:  time.Date(2020, 1, 1, 0, 0, 1, 0, time.UTC),
						},
					},
					HasNext: false,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got ogen.Projects
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

func TestUpdateProject(t *testing.T) {
	tests := []test{
		{
			name: "プロジェクト1を更新する",
			request: request{
				method: http.MethodPatch,
				path:   "/projects/01DXF6DT000000000000000000",
				body:   strings.NewReader(`{"name": "新プロジェクト1"}`),
			},
			response: response{
				statusCode: http.StatusOK,
				body: ogen.Project{
					ID:         "01DXF6DT000000000000000000",
					Name:       "新プロジェクト1",
					Color:      "#1A2B3C",
					IsArchived: false,
					CreatedAt:  time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt:  time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tdb.Begin(); err != nil {
				t.Fatalf("tdb.Begin failed: %s", err)
			}
			defer tdb.Rollback()

			var got ogen.Project
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

func TestDeleteProject(t *testing.T) {
	tests := []test{
		{
			name: "プロジェクト1を削除する",
			request: request{
				method: http.MethodDelete,
				path:   "/projects/01DXF6DT000000000000000000",
			},
			response: response{
				statusCode: http.StatusNoContent,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tdb.Begin(); err != nil {
				t.Fatalf("tdb.Begin failed: %s", err)
			}
			defer tdb.Rollback()

			resp, err := doTestRequest(tt.request, nil)
			if err != nil {
				t.Fatalf("doTestRequest failed: %s", err)
			}

			if tt.response.statusCode != resp.StatusCode {
				t.Errorf("status code want %d, but %d", tt.response.statusCode, resp.StatusCode)
			}
		})
	}
}
