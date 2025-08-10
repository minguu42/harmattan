package handler_test

import (
	"testing"
	"time"

	"github.com/ikawaha/httpcheck"
	"github.com/minguu42/harmattan/api/handler"
	"github.com/minguu42/harmattan/internal/database"
	"github.com/minguu42/harmattan/internal/openapi"
	"github.com/stretchr/testify/require"
)

func TestHandler_CreateTask(t *testing.T) {
	require.NoError(t, tdb.Reset(t.Context(), []any{database.Project{}, database.Task{}, database.Step{}, database.Tag{}}))
	require.NoError(t, tdb.Insert(t.Context(), []any{
		database.Projects{
			{
				ID:        "PROJECT-000000000000000001",
				UserID:    testUserID,
				Name:      "テストプロジェクト",
				Color:     "blue",
				CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
				UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
			},
			{
				ID:        "PROJECT-000000000000000002",
				UserID:    "USER-000000000000000000002",
				Name:      "テストプロジェクト2",
				Color:     "gray",
				CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
				UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
			},
		},
	}))

	t.Run("project not found", func(t *testing.T) {
		want := handler.ErrorResponse{Code: 404, Message: "指定したプロジェクトは見つかりません"}
		httpcheck.New(th).Test(t, "POST", "/projects/PROJECT-000000000000000099/tasks").
			WithHeader("Authorization", "Bearer "+token).
			WithHeader("Content-Type", "application/json").
			WithBody([]byte(`{"name": "テストタスク", "priority": 1}`)).
			Check().HasStatus(404).HasJSON(want)
	})
	t.Run("user does not own the project", func(t *testing.T) {
		want := handler.ErrorResponse{Code: 404, Message: "指定したプロジェクトは見つかりません"}
		httpcheck.New(th).Test(t, "POST", "/projects/PROJECT-000000000000000002/tasks").
			WithHeader("Authorization", "Bearer "+token).
			WithHeader("Content-Type", "application/json").
			WithBody([]byte(`{"name": "テストタスク", "priority": 1}`)).
			Check().HasStatus(404).HasJSON(want)
	})
	t.Run("ok", func(t *testing.T) {
		want := &openapi.Task{
			ID:        fixedID,
			ProjectID: "PROJECT-000000000000000001",
			Name:      "テストタスク",
			Priority:  1,
			CreatedAt: fixedNow,
			UpdatedAt: fixedNow,
		}
		httpcheck.New(th).Test(t, "POST", "/projects/PROJECT-000000000000000001/tasks").
			WithHeader("Authorization", "Bearer "+token).
			WithHeader("Content-Type", "application/json").
			WithBody([]byte(`{"name": "テストタスク", "priority": 1}`)).
			Check().HasStatus(200).HasJSON(want)
	})
}

func TestHandler_ListTasks(t *testing.T) {
	require.NoError(t, tdb.Reset(t.Context(), []any{database.Project{}, database.Task{}, database.Step{}, database.Tag{}}))
	require.NoError(t, tdb.Insert(t.Context(), []any{
		database.Projects{
			{
				ID:        "PROJECT-000000000000000001",
				UserID:    testUserID,
				Name:      "テストプロジェクト",
				Color:     "blue",
				CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst),
				UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst),
			},
			{
				ID:        "PROJECT-000000000000000002",
				UserID:    "USER-000000000000000000002",
				Name:      "テストプロジェクト2",
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
		},
	}))

	t.Run("no limit and offset", func(t *testing.T) {
		want := &openapi.ListTasksOK{
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
		}
		httpcheck.New(th).Test(t, "GET", "/projects/PROJECT-000000000000000001/tasks").
			WithHeader("Authorization", "Bearer "+token).
			Check().HasStatus(200).HasJSON(want)
	})
	t.Run("limit=1&offset=0", func(t *testing.T) {
		want := &openapi.ListTasksOK{
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
		}
		httpcheck.New(th).Test(t, "GET", "/projects/PROJECT-000000000000000001/tasks?limit=1&offset=0").
			WithHeader("Authorization", "Bearer "+token).
			Check().HasStatus(200).HasJSON(want)
	})
	t.Run("limit=1&offset=1", func(t *testing.T) {
		want := &openapi.ListTasksOK{
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
		}
		httpcheck.New(th).Test(t, "GET", "/projects/PROJECT-000000000000000001/tasks?limit=1&offset=1").
			WithHeader("Authorization", "Bearer "+token).
			Check().HasStatus(200).HasJSON(want)
	})
	t.Run("project not found", func(t *testing.T) {
		want := handler.ErrorResponse{Code: 404, Message: "指定したプロジェクトは見つかりません"}
		httpcheck.New(th).Test(t, "GET", "/projects/PROJECT-000000000000000099/tasks").
			WithHeader("Authorization", "Bearer "+token).
			Check().HasStatus(404).HasJSON(want)
	})
	t.Run("user does not own the project", func(t *testing.T) {
		want := handler.ErrorResponse{Code: 404, Message: "指定したプロジェクトは見つかりません"}
		httpcheck.New(th).Test(t, "GET", "/projects/PROJECT-000000000000000002/tasks").
			WithHeader("Authorization", "Bearer "+token).
			Check().HasStatus(404).HasJSON(want)
	})
}

func TestHandler_UpdateTask(t *testing.T) {
	require.NoError(t, tdb.Reset(t.Context(), []any{database.Project{}, database.Task{}, database.Step{}, database.Tag{}}))
	require.NoError(t, tdb.Insert(t.Context(), []any{
		database.Projects{
			{
				ID:        "PROJECT-000000000000000001",
				UserID:    testUserID,
				Name:      "テストプロジェクト",
				Color:     "blue",
				CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
				UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
			},
			{
				ID:        "PROJECT-000000000000000002",
				UserID:    "USER-000000000000000000002",
				Name:      "他のユーザーのプロジェクト",
				Color:     "gray",
				CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
				UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
			},
			{
				ID:        "PROJECT-000000000000000003",
				UserID:    testUserID,
				Name:      "テストプロジェクト3",
				Color:     "gray",
				CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
				UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
			},
		},
		database.Tasks{
			{
				ID:        "TASK-000000000000000000001",
				UserID:    testUserID,
				ProjectID: "PROJECT-000000000000000001",
				Name:      "更新前タスク",
				Content:   "更新前内容",
				Priority:  1,
				CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
				UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
			},
			{
				ID:        "TASK-000000000000000000002",
				UserID:    "USER-000000000000000000002",
				ProjectID: "PROJECT-000000000000000002",
				Name:      "他のユーザーのタスク",
				Content:   "他のユーザーの内容",
				Priority:  2,
				CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
				UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
			},
			{
				ID:        "TASK-000000000000000000003",
				UserID:    testUserID,
				ProjectID: "PROJECT-000000000000000003",
				Name:      "更新前タスク",
				Content:   "更新前内容",
				Priority:  1,
				CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
				UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
			},
		},
	}))

	t.Run("project not found", func(t *testing.T) {
		want := handler.ErrorResponse{Code: 404, Message: "指定したプロジェクトは見つかりません"}
		httpcheck.New(th).Test(t, "PATCH", "/projects/PROJECT-000000000000000099/tasks/TASK-000000000000000000001").
			WithHeader("Authorization", "Bearer "+token).
			WithHeader("Content-Type", "application/json").
			WithBody([]byte(`{"name": "更新後タスク", "content": "更新後内容", "priority": 3}`)).
			Check().HasStatus(404).HasJSON(want)
	})
	t.Run("user does not own the project", func(t *testing.T) {
		want := handler.ErrorResponse{Code: 404, Message: "指定したプロジェクトは見つかりません"}
		httpcheck.New(th).Test(t, "PATCH", "/projects/PROJECT-000000000000000002/tasks/TASK-000000000000000000001").
			WithHeader("Authorization", "Bearer "+token).
			WithHeader("Content-Type", "application/json").
			WithBody([]byte(`{"name": "更新後タスク", "content": "更新後内容", "priority": 3}`)).
			Check().HasStatus(404).HasJSON(want)
	})
	t.Run("task not found", func(t *testing.T) {
		want := handler.ErrorResponse{Code: 404, Message: "指定したタスクは見つかりません"}
		httpcheck.New(th).Test(t, "PATCH", "/projects/PROJECT-000000000000000001/tasks/TASK-000000000000000000099").
			WithHeader("Authorization", "Bearer "+token).
			WithHeader("Content-Type", "application/json").
			WithBody([]byte(`{"name": "更新後タスク", "content": "更新後内容", "priority": 3}`)).
			Check().HasStatus(404).HasJSON(want)
	})
	t.Run("user does not own the task", func(t *testing.T) {
		want := handler.ErrorResponse{Code: 404, Message: "指定したタスクは見つかりません"}
		httpcheck.New(th).Test(t, "PATCH", "/projects/PROJECT-000000000000000001/tasks/TASK-000000000000000000002").
			WithHeader("Authorization", "Bearer "+token).
			WithHeader("Content-Type", "application/json").
			WithBody([]byte(`{"name": "更新後タスク", "content": "更新後内容", "priority": 3}`)).
			Check().HasStatus(404).HasJSON(want)
	})
	t.Run("task does not belong to the project", func(t *testing.T) {
		want := handler.ErrorResponse{Code: 404, Message: "指定したタスクは見つかりません"}
		httpcheck.New(th).Test(t, "PATCH", "/projects/PROJECT-000000000000000001/tasks/TASK-000000000000000000003").
			WithHeader("Authorization", "Bearer "+token).
			WithHeader("Content-Type", "application/json").
			WithBody([]byte(`{"name": "更新後タスク", "content": "更新後内容", "priority": 3}`)).
			Check().HasStatus(404).HasJSON(want)
	})
	t.Run("ok", func(t *testing.T) {
		want := &openapi.Task{
			ID:          "TASK-000000000000000000001",
			ProjectID:   "PROJECT-000000000000000001",
			Name:        "更新後タスク",
			Content:     "更新後内容",
			Priority:    3,
			DueOn:       openapi.OptDate{Value: time.Date(2025, 1, 2, 0, 0, 0, 0, jst), Set: true},
			CompletedAt: openapi.OptDateTime{Value: time.Date(2025, 1, 1, 12, 0, 0, 0, jst), Set: true},
			CreatedAt:   time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
			UpdatedAt:   fixedNow,
		}
		httpcheck.New(th).Test(t, "PATCH", "/projects/PROJECT-000000000000000001/tasks/TASK-000000000000000000001").
			WithHeader("Authorization", "Bearer "+token).
			WithHeader("Content-Type", "application/json").
			WithBody([]byte(`{"name": "更新後タスク", "content": "更新後内容", "priority": 3, "due_on": "2025-01-02T00:00:00+09:00", "completed_at": "2025-01-01T12:00:00+09:00"}`)).
			Check().HasStatus(200).HasJSON(want)
	})
}

func TestHandler_GetTask(t *testing.T) {
	require.NoError(t, tdb.Reset(t.Context(), []any{database.Project{}, database.Task{}, database.Step{}, database.Tag{}}))
	require.NoError(t, tdb.Insert(t.Context(), []any{
		database.Projects{
			{
				ID:        "PROJECT-000000000000000001",
				UserID:    testUserID,
				Name:      "テストプロジェクト",
				Color:     "blue",
				CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
				UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
			},
			{
				ID:        "PROJECT-000000000000000002",
				UserID:    "USER-000000000000000000002",
				Name:      "他のユーザーのプロジェクト",
				Color:     "gray",
				CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
				UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
			},
			{
				ID:        "PROJECT-000000000000000003",
				UserID:    testUserID,
				Name:      "テストプロジェクト3",
				Color:     "gray",
				CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
				UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
			},
		},
		database.Tasks{
			{
				ID:        "TASK-000000000000000000001",
				UserID:    testUserID,
				ProjectID: "PROJECT-000000000000000001",
				Name:      "テストタスク",
				Content:   "テスト内容",
				Priority:  1,
				CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
				UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
			},
			{
				ID:        "TASK-000000000000000000002",
				UserID:    "USER-000000000000000000002",
				ProjectID: "PROJECT-000000000000000002",
				Name:      "他のユーザーのタスク",
				Content:   "他のユーザーの内容",
				Priority:  2,
				CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
				UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
			},
			{
				ID:        "TASK-000000000000000000003",
				UserID:    testUserID,
				ProjectID: "PROJECT-000000000000000003",
				Name:      "タスク3",
				Content:   "内容3",
				Priority:  1,
				CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
				UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
			},
		},
	}))

	t.Run("project not found", func(t *testing.T) {
		want := handler.ErrorResponse{Code: 404, Message: "指定したプロジェクトは見つかりません"}
		httpcheck.New(th).Test(t, "GET", "/projects/PROJECT-000000000000000099/tasks/TASK-000000000000000000001").
			WithHeader("Authorization", "Bearer "+token).
			Check().HasStatus(404).HasJSON(want)
	})
	t.Run("user does not own the project", func(t *testing.T) {
		want := handler.ErrorResponse{Code: 404, Message: "指定したプロジェクトは見つかりません"}
		httpcheck.New(th).Test(t, "GET", "/projects/PROJECT-000000000000000002/tasks/TASK-000000000000000000001").
			WithHeader("Authorization", "Bearer "+token).
			Check().HasStatus(404).HasJSON(want)
	})
	t.Run("task not found", func(t *testing.T) {
		want := handler.ErrorResponse{Code: 404, Message: "指定したタスクは見つかりません"}
		httpcheck.New(th).Test(t, "GET", "/projects/PROJECT-000000000000000001/tasks/TASK-000000000000000000099").
			WithHeader("Authorization", "Bearer "+token).
			Check().HasStatus(404).HasJSON(want)
	})
	t.Run("user does not own the task", func(t *testing.T) {
		want := handler.ErrorResponse{Code: 404, Message: "指定したタスクは見つかりません"}
		httpcheck.New(th).Test(t, "GET", "/projects/PROJECT-000000000000000001/tasks/TASK-000000000000000000002").
			WithHeader("Authorization", "Bearer "+token).
			Check().HasStatus(404).HasJSON(want)
	})
	t.Run("task does not belong to the project", func(t *testing.T) {
		want := handler.ErrorResponse{Code: 404, Message: "指定したタスクは見つかりません"}
		httpcheck.New(th).Test(t, "GET", "/projects/PROJECT-000000000000000001/tasks/TASK-000000000000000000003").
			WithHeader("Authorization", "Bearer "+token).
			Check().HasStatus(404).HasJSON(want)
	})
	t.Run("ok", func(t *testing.T) {
		want := &openapi.Task{
			ID:        "TASK-000000000000000000001",
			ProjectID: "PROJECT-000000000000000001",
			Name:      "テストタスク",
			Content:   "テスト内容",
			Priority:  1,
			CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
			UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
		}
		httpcheck.New(th).Test(t, "GET", "/projects/PROJECT-000000000000000001/tasks/TASK-000000000000000000001").
			WithHeader("Authorization", "Bearer "+token).
			Check().HasStatus(200).HasJSON(want)
	})
}

func TestHandler_DeleteTask(t *testing.T) {
	require.NoError(t, tdb.Reset(t.Context(), []any{database.Project{}, database.Task{}, database.Step{}, database.Tag{}}))
	require.NoError(t, tdb.Insert(t.Context(), []any{
		database.Projects{
			{
				ID:        "PROJECT-000000000000000001",
				UserID:    testUserID,
				Name:      "テストプロジェクト",
				Color:     "blue",
				CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
				UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
			},
			{
				ID:        "PROJECT-000000000000000002",
				UserID:    "USER-000000000000000000002",
				Name:      "他のユーザーのプロジェクト",
				Color:     "gray",
				CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
				UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
			},
			{
				ID:        "PROJECT-000000000000000003",
				UserID:    testUserID,
				Name:      "テストプロジェクト3",
				Color:     "gray",
				CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
				UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
			},
		},
		database.Tasks{
			{
				ID:        "TASK-000000000000000000001",
				UserID:    testUserID,
				ProjectID: "PROJECT-000000000000000001",
				Name:      "削除対象タスク",
				Content:   "削除対象内容",
				Priority:  1,
				CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
				UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
			},
			{
				ID:        "TASK-000000000000000000002",
				UserID:    "USER-000000000000000000002",
				ProjectID: "PROJECT-000000000000000002",
				Name:      "他のユーザーのタスク",
				Content:   "他のユーザーの内容",
				Priority:  2,
				CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
				UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
			},

			{
				ID:        "TASK-000000000000000000003",
				UserID:    testUserID,
				ProjectID: "PROJECT-000000000000000003",
				Name:      "タスク3",
				Content:   "更新前内容",
				Priority:  1,
				CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
				UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
			},
		},
	}))

	t.Run("project not found", func(t *testing.T) {
		want := handler.ErrorResponse{Code: 404, Message: "指定したプロジェクトは見つかりません"}
		httpcheck.New(th).Test(t, "DELETE", "/projects/PROJECT-000000000000000099/tasks/TASK-000000000000000000001").
			WithHeader("Authorization", "Bearer "+token).
			Check().HasStatus(404).HasJSON(want)
	})
	t.Run("user does not own the project", func(t *testing.T) {
		want := handler.ErrorResponse{Code: 404, Message: "指定したプロジェクトは見つかりません"}
		httpcheck.New(th).Test(t, "DELETE", "/projects/PROJECT-000000000000000002/tasks/TASK-000000000000000000001").
			WithHeader("Authorization", "Bearer "+token).
			Check().HasStatus(404).HasJSON(want)
	})
	t.Run("task not found", func(t *testing.T) {
		want := handler.ErrorResponse{Code: 404, Message: "指定したタスクは見つかりません"}
		httpcheck.New(th).Test(t, "DELETE", "/projects/PROJECT-000000000000000001/tasks/TASK-000000000000000000099").
			WithHeader("Authorization", "Bearer "+token).
			Check().HasStatus(404).HasJSON(want)
	})
	t.Run("user does not own the task", func(t *testing.T) {
		want := handler.ErrorResponse{Code: 404, Message: "指定したタスクは見つかりません"}
		httpcheck.New(th).Test(t, "DELETE", "/projects/PROJECT-000000000000000001/tasks/TASK-000000000000000000002").
			WithHeader("Authorization", "Bearer "+token).
			Check().HasStatus(404).HasJSON(want)
	})
	t.Run("task does not belong to the project", func(t *testing.T) {
		want := handler.ErrorResponse{Code: 404, Message: "指定したタスクは見つかりません"}
		httpcheck.New(th).Test(t, "DELETE", "/projects/PROJECT-000000000000000001/tasks/TASK-000000000000000000003").
			WithHeader("Authorization", "Bearer "+token).
			Check().HasStatus(404).HasJSON(want)
	})
	t.Run("ok", func(t *testing.T) {
		httpcheck.New(th).Test(t, "DELETE", "/projects/PROJECT-000000000000000001/tasks/TASK-000000000000000000001").
			WithHeader("Authorization", "Bearer "+token).
			Check().HasStatus(200)
	})
}
