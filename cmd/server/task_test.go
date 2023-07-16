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

func TestCreateTask(t *testing.T) {
	tests := []test{
		{
			name: "タスクを作成する",
			request: request{
				method: http.MethodPost,
				path:   "/projects/01DXF6DT000000000000000000/tasks",
				body:   strings.NewReader(`{"title": "新タスク", "content": "Hello, 世界!", "priority": 3, "dueOn": "2020-01-02"}`),
			},
			response: response{
				statusCode: http.StatusCreated,
				body: ogen.Task{
					ID:          "01DXF6DT000000000000000002",
					ProjectID:   "01DXF6DT000000000000000000",
					Title:       "新タスク",
					Content:     "Hello, 世界!",
					Priority:    3,
					DueOn:       ogen.OptDate{Value: time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC), Set: true},
					CompletedAt: ogen.OptDateTime{Set: false},
					CreatedAt:   time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
					UpdatedAt:   time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
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

			var got ogen.Task
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

func TestListTask(t *testing.T) {
	tests := []test{
		{
			name: "タスク一覧を取得する",
			request: request{
				method: http.MethodGet,
				path:   "/projects/01DXF6DT000000000000000000/tasks?sort=createdAt",
			},
			response: response{
				statusCode: http.StatusOK,
				body: ogen.Tasks{
					Tasks: []ogen.Task{
						{
							ID:          "01DXF6DT000000000000000000",
							ProjectID:   "01DXF6DT000000000000000000",
							Title:       "タスク1",
							Content:     "Hello, 世界!",
							Priority:    0,
							DueOn:       ogen.OptDate{Set: false},
							CompletedAt: ogen.OptDateTime{Set: false},
							CreatedAt:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
							UpdatedAt:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
						},
						{
							ID:          "01DXF6DT000000000000000001",
							ProjectID:   "01DXF6DT000000000000000000",
							Title:       "タスク2",
							Content:     "",
							Priority:    3,
							DueOn:       ogen.OptDate{Value: time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC), Set: true},
							CompletedAt: ogen.OptDateTime{Value: time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC), Set: true},
							CreatedAt:   time.Date(2020, 1, 1, 0, 0, 1, 0, time.UTC),
							UpdatedAt:   time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
						},
					},
					HasNext: false,
				},
			},
		},
	}
	for _, tt := range tests {
		var got ogen.Tasks
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
	}
}

func TestUpdateTask(t *testing.T) {
	tests := []test{
		{
			name: "タスク1を更新する",
			request: request{
				method: http.MethodPatch,
				path:   "/projects/01DXF6DT000000000000000000/tasks/01DXF6DT000000000000000000",
				body:   strings.NewReader(`{"isCompleted": true}`),
			},
			response: response{
				statusCode: http.StatusOK,
				body: ogen.Task{
					ID:          "01DXF6DT000000000000000000",
					ProjectID:   "01DXF6DT000000000000000000",
					Title:       "タスク1",
					Content:     "Hello, 世界!",
					CompletedAt: ogen.OptDateTime{Value: time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC), Set: true},
					CreatedAt:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt:   time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
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

			var got ogen.Task
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

func TestDeleteTask(t *testing.T) {
	tests := []test{
		{
			name: "タスク1を削除する",
			request: request{
				method: http.MethodDelete,
				path:   "/projects/01DXF6DT000000000000000000/tasks/01DXF6DT000000000000000000",
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
