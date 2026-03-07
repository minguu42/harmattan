package usecase_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/minguu42/harmattan/internal/api/handler"
	"github.com/minguu42/harmattan/internal/api/openapi"
	"github.com/minguu42/harmattan/internal/database"
	"github.com/stretchr/testify/require"
)

func TestTask_CreateTask(t *testing.T) {
	require.NoError(t, tdb.TruncateAndInsert(t.Context(), []any{
		database.Projects{
			{
				ID:     "PROJECT-000000000000000001",
				UserID: testUserID,
				Name:   "プロジェクト1",
				Color:  "blue",
			},
			{
				ID:     "PROJECT-000000000000000002",
				UserID: "USER-000000000000000000002",
				Name:   "プロジェクト2",
				Color:  "gray",
			},
		},
		database.Tasks{},
		database.Steps{},
		database.Tags{},
	}))

	tests := map[string]test{
		"ok": {
			Method:     "POST",
			Path:       "/projects/PROJECT-000000000000000001/tasks",
			Headers:    http.Header{"Authorization": []string{"Bearer " + token}},
			Body:       `{"name": "タスク", "priority": 1}`,
			WantStatus: 200,
			WantJSON: &openapi.Task{
				ID:        fixedID,
				ProjectID: "PROJECT-000000000000000001",
				Name:      "タスク",
				Priority:  1,
				CreatedAt: fixedNow,
				UpdatedAt: fixedNow,
			},
			WantTables: []any{database.Tasks{
				{ID: fixedID, UserID: testUserID, ProjectID: "PROJECT-000000000000000001", Name: "タスク", Priority: 1, CreatedAt: fixedNow, UpdatedAt: fixedNow},
			}},
		},
		"not_found": {
			Method:     "POST",
			Path:       "/projects/PROJECT-000000000000000099/tasks",
			Headers:    http.Header{"Authorization": []string{"Bearer " + token}},
			Body:       `{"name": "タスク", "priority": 1}`,
			WantStatus: 404,
			WantJSON:   handler.ErrorResponse{Code: 404, Message: "指定したプロジェクトは見つかりません"},
		},
		"access_denied": {
			Method:     "POST",
			Path:       "/projects/PROJECT-000000000000000002/tasks",
			Headers:    http.Header{"Authorization": []string{"Bearer " + token}},
			Body:       `{"name": "タスク", "priority": 1}`,
			WantStatus: 404,
			WantJSON:   handler.ErrorResponse{Code: 404, Message: "指定したプロジェクトは見つかりません"},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			runTest(t, tt)
		})
	}
}

func TestTask_ListTasks(t *testing.T) {
	completedAt := time.Date(2025, 1, 1, 12, 0, 0, 0, jst)
	require.NoError(t, tdb.TruncateAndInsert(t.Context(), []any{
		database.Projects{
			{
				ID:        "PROJECT-000000000000000001",
				UserID:    testUserID,
				Name:      "プロジェクト1",
				Color:     "blue",
				CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst),
				UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst),
			},
			{
				ID:        "PROJECT-000000000000000002",
				UserID:    "USER-000000000000000000002",
				Name:      "プロジェクト2",
				Color:     "gray",
				CreatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst),
				UpdatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst),
			},
		},
		database.Tasks{
			{
				ID:        "TASK-000000000000000000001",
				UserID:    testUserID,
				ProjectID: "PROJECT-000000000000000001",
				Name:      "タスク1",
				Content:   "タスク1の内容",
				Priority:  1,
				CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst),
				UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst),
			},
			{
				ID:        "TASK-000000000000000000002",
				UserID:    testUserID,
				ProjectID: "PROJECT-000000000000000001",
				Name:      "タスク2",
				Content:   "タスク2の内容",
				Priority:  2,
				CreatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst),
				UpdatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst),
			},
			{
				ID:          "TASK-000000000000000000003",
				UserID:      testUserID,
				ProjectID:   "PROJECT-000000000000000001",
				Name:        "タスク3（完了済み）",
				Content:     "タスク3の内容",
				Priority:    3,
				CompletedAt: &completedAt,
				CreatedAt:   time.Date(2025, 1, 1, 0, 0, 3, 0, jst),
				UpdatedAt:   time.Date(2025, 1, 1, 0, 0, 3, 0, jst),
			},
		},
		database.Steps{},
		database.TaskTags{},
	}))

	tests := map[string]test{
		"show_completed_false_default": {
			Method:     "GET",
			Path:       "/projects/PROJECT-000000000000000001/tasks",
			Headers:    http.Header{"Authorization": []string{"Bearer " + token}},
			WantStatus: 200,
			WantJSON: &openapi.ListTasksOK{
				Tasks: []openapi.Task{
					{
						ID:        "TASK-000000000000000000001",
						ProjectID: "PROJECT-000000000000000001",
						Name:      "タスク1",
						Content:   "タスク1の内容",
						Priority:  1,
						CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst),
						UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst),
					},
					{
						ID:        "TASK-000000000000000000002",
						ProjectID: "PROJECT-000000000000000001",
						Name:      "タスク2",
						Content:   "タスク2の内容",
						Priority:  2,
						CreatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst),
						UpdatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst),
					},
				},
				HasNext: false,
			},
		},
		"show_completed_true": {
			Method:     "GET",
			Path:       "/projects/PROJECT-000000000000000001/tasks?showCompleted=true",
			Headers:    http.Header{"Authorization": []string{"Bearer " + token}},
			WantStatus: 200,
			WantJSON: &openapi.ListTasksOK{
				Tasks: []openapi.Task{
					{
						ID:        "TASK-000000000000000000001",
						ProjectID: "PROJECT-000000000000000001",
						Name:      "タスク1",
						Content:   "タスク1の内容",
						Priority:  1,
						CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst),
						UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst),
					},
					{
						ID:        "TASK-000000000000000000002",
						ProjectID: "PROJECT-000000000000000001",
						Name:      "タスク2",
						Content:   "タスク2の内容",
						Priority:  2,
						CreatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst),
						UpdatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst),
					},
					{
						ID:          "TASK-000000000000000000003",
						ProjectID:   "PROJECT-000000000000000001",
						Name:        "タスク3（完了済み）",
						Content:     "タスク3の内容",
						Priority:    3,
						CompletedAt: openapi.OptDateTime{Value: completedAt, Set: true},
						CreatedAt:   time.Date(2025, 1, 1, 0, 0, 3, 0, jst),
						UpdatedAt:   time.Date(2025, 1, 1, 0, 0, 3, 0, jst),
					},
				},
				HasNext: false,
			},
		},
		"limit_1_offset_0": {
			Method:     "GET",
			Path:       "/projects/PROJECT-000000000000000001/tasks?limit=1&offset=0",
			Headers:    http.Header{"Authorization": []string{"Bearer " + token}},
			WantStatus: 200,
			WantJSON: &openapi.ListTasksOK{
				Tasks: []openapi.Task{
					{
						ID:        "TASK-000000000000000000001",
						ProjectID: "PROJECT-000000000000000001",
						Name:      "タスク1",
						Content:   "タスク1の内容",
						Priority:  1,
						CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst),
						UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst),
					},
				},
				HasNext: true,
			},
		},
		"limit_1_offset_1": {
			Method:     "GET",
			Path:       "/projects/PROJECT-000000000000000001/tasks?limit=1&offset=1",
			Headers:    http.Header{"Authorization": []string{"Bearer " + token}},
			WantStatus: 200,
			WantJSON: &openapi.ListTasksOK{
				Tasks: []openapi.Task{
					{
						ID:        "TASK-000000000000000000002",
						ProjectID: "PROJECT-000000000000000001",
						Name:      "タスク2",
						Content:   "タスク2の内容",
						Priority:  2,
						CreatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst),
						UpdatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst),
					},
				},
				HasNext: false,
			},
		},
		"not_found": {
			Method:     "GET",
			Path:       "/projects/PROJECT-000000000000000099/tasks",
			Headers:    http.Header{"Authorization": []string{"Bearer " + token}},
			WantStatus: 404,
			WantJSON:   handler.ErrorResponse{Code: 404, Message: "指定したプロジェクトは見つかりません"},
		},
		"access_denied": {
			Method:     "GET",
			Path:       "/projects/PROJECT-000000000000000002/tasks",
			Headers:    http.Header{"Authorization": []string{"Bearer " + token}},
			WantStatus: 404,
			WantJSON:   handler.ErrorResponse{Code: 404, Message: "指定したプロジェクトは見つかりません"},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			runTest(t, tt)
		})
	}
}

func TestTask_GetTask(t *testing.T) {
	require.NoError(t, tdb.TruncateAndInsert(t.Context(), []any{
		database.Projects{
			{
				ID:        "PROJECT-000000000000000001",
				UserID:    testUserID,
				Name:      "プロジェクト1",
				Color:     "blue",
				CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst),
				UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst),
			},
			{
				ID:        "PROJECT-000000000000000002",
				UserID:    "USER-000000000000000000002",
				Name:      "プロジェクト2",
				Color:     "gray",
				CreatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst),
				UpdatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst),
			},
		},
		database.Tasks{
			{
				ID:        "TASK-000000000000000000001",
				UserID:    testUserID,
				ProjectID: "PROJECT-000000000000000001",
				Name:      "タスク1",
				Content:   "内容",
				Priority:  1,
				CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst),
				UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst),
			},
			{
				ID:        "TASK-000000000000000000002",
				UserID:    "USER-000000000000000000002",
				ProjectID: "PROJECT-000000000000000002",
				Name:      "タスク2",
				Content:   "内容",
				Priority:  2,
				CreatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst),
				UpdatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst),
			},
		},
		database.Steps{},
		database.TaskTags{},
	}))

	tests := map[string]test{
		"ok": {
			Method:     "GET",
			Path:       "/tasks/TASK-000000000000000000001",
			Headers:    http.Header{"Authorization": []string{"Bearer " + token}},
			WantStatus: 200,
			WantJSON: &openapi.Task{
				ID:        "TASK-000000000000000000001",
				ProjectID: "PROJECT-000000000000000001",
				Name:      "タスク1",
				Content:   "内容",
				Priority:  1,
				CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst),
				UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst),
			},
		},
		"not_found": {
			Method:     "GET",
			Path:       "/tasks/TASK-000000000000000000099",
			Headers:    http.Header{"Authorization": []string{"Bearer " + token}},
			WantStatus: 404,
			WantJSON:   handler.ErrorResponse{Code: 404, Message: "指定したタスクは見つかりません"},
		},
		"access_denied": {
			Method:     "GET",
			Path:       "/tasks/TASK-000000000000000000002",
			Headers:    http.Header{"Authorization": []string{"Bearer " + token}},
			WantStatus: 404,
			WantJSON:   handler.ErrorResponse{Code: 404, Message: "指定したタスクは見つかりません"},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			runTest(t, tt)
		})
	}
}

func TestTask_UpdateTask(t *testing.T) {
	require.NoError(t, tdb.TruncateAndInsert(t.Context(), []any{
		database.Projects{
			{
				ID:     "PROJECT-000000000000000001",
				UserID: testUserID,
				Name:   "プロジェクト1",
				Color:  "blue",
			},
			{
				ID:     "PROJECT-000000000000000002",
				UserID: "USER-000000000000000000002",
				Name:   "プロジェクト2",
				Color:  "gray",
			},
		},
		database.Tasks{
			{
				ID:        "TASK-000000000000000000001",
				UserID:    testUserID,
				ProjectID: "PROJECT-000000000000000001",
				Name:      "タスク1",
				CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst),
				UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst),
			},
			{
				ID:        "TASK-000000000000000000002",
				UserID:    "USER-000000000000000000002",
				ProjectID: "PROJECT-000000000000000002",
				Name:      "タスク2",
			},
		},
		database.Tags{
			{
				ID:        "TAG-0000000000000000000001",
				UserID:    testUserID,
				Name:      "タグ1",
				CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst),
				UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst),
			},
			{
				ID:     "TAG-0000000000000000000002",
				UserID: "USER-000000000000000000002",
				Name:   "タグ2",
			},
		},
		database.Steps{},
		database.TaskTags{},
	}))

	tests := map[string]test{
		"ok": {
			Method:     "PATCH",
			Path:       "/tasks/TASK-000000000000000000001",
			Headers:    http.Header{"Authorization": []string{"Bearer " + token}},
			Body:       `{"name": "更新後タスク", "content": "更新後内容", "priority": 3, "due_on": "2025-01-02", "completed_at": "2025-01-01T12:00:00+09:00", "tag_ids": ["TAG-0000000000000000000001", "TAG-0000000000000000000002", "TAG-0000000000000000000099"]}`,
			WantStatus: 200,
			WantJSON: &openapi.Task{
				ID:          "TASK-000000000000000000001",
				ProjectID:   "PROJECT-000000000000000001",
				Name:        "更新後タスク",
				Content:     "更新後内容",
				Priority:    3,
				DueOn:       openapi.OptDate{Value: time.Date(2025, 1, 2, 0, 0, 0, 0, jst), Set: true},
				CompletedAt: openapi.OptDateTime{Value: time.Date(2025, 1, 1, 12, 0, 0, 0, jst), Set: true},
				CreatedAt:   time.Date(2025, 1, 1, 0, 0, 1, 0, jst),
				UpdatedAt:   fixedNow,
				Tags: []openapi.Tag{
					{
						ID:        "TAG-0000000000000000000001",
						Name:      "タグ1",
						CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst),
						UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst),
					},
				},
			},
		},
		"not_found": {
			Method:     "PATCH",
			Path:       "/tasks/TASK-000000000000000000099",
			Headers:    http.Header{"Authorization": []string{"Bearer " + token}},
			Body:       `{"name": "更新後タスク", "content": "更新後内容", "priority": 3}`,
			WantStatus: 404,
			WantJSON:   handler.ErrorResponse{Code: 404, Message: "指定したタスクは見つかりません"},
		},
		"access_denied": {
			Method:     "PATCH",
			Path:       "/tasks/TASK-000000000000000000002",
			Headers:    http.Header{"Authorization": []string{"Bearer " + token}},
			Body:       `{"name": "更新後タスク", "content": "更新後内容", "priority": 3}`,
			WantStatus: 404,
			WantJSON:   handler.ErrorResponse{Code: 404, Message: "指定したタスクは見つかりません"},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			runTest(t, tt)
		})
	}
}

func TestTask_DeleteTask(t *testing.T) {
	require.NoError(t, tdb.TruncateAndInsert(t.Context(), []any{
		database.Projects{
			{
				ID:     "PROJECT-000000000000000001",
				UserID: testUserID,
				Name:   "プロジェクト1",
				Color:  "blue",
			},
			{
				ID:     "PROJECT-000000000000000002",
				UserID: "USER-000000000000000000002",
				Name:   "プロジェクト2",
				Color:  "gray",
			},
		},
		database.Tasks{
			{
				ID:        "TASK-000000000000000000001",
				UserID:    testUserID,
				ProjectID: "PROJECT-000000000000000001",
				Name:      "タスク1",
				CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst),
				UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst),
			},
			{
				ID:        "TASK-000000000000000000002",
				UserID:    "USER-000000000000000000002",
				ProjectID: "PROJECT-000000000000000002",
				Name:      "タスク2",
				CreatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst),
				UpdatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst),
			},
		},
	}))

	tests := map[string]test{
		"ok": {
			Method:     "DELETE",
			Path:       "/tasks/TASK-000000000000000000001",
			Headers:    http.Header{"Authorization": []string{"Bearer " + token}},
			WantStatus: 200,
			WantTables: []any{database.Tasks{
				{ID: "TASK-000000000000000000002", UserID: "USER-000000000000000000002", ProjectID: "PROJECT-000000000000000002", Name: "タスク2", CreatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst)},
			}},
		},
		"not_found": {
			Method:     "DELETE",
			Path:       "/tasks/TASK-000000000000000000099",
			Headers:    http.Header{"Authorization": []string{"Bearer " + token}},
			WantStatus: 404,
			WantJSON:   handler.ErrorResponse{Code: 404, Message: "指定したタスクは見つかりません"},
		},
		"access_denied": {
			Method:     "DELETE",
			Path:       "/tasks/TASK-000000000000000000002",
			Headers:    http.Header{"Authorization": []string{"Bearer " + token}},
			WantStatus: 404,
			WantJSON:   handler.ErrorResponse{Code: 404, Message: "指定したタスクは見つかりません"},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			runTest(t, tt)
		})
	}
}
