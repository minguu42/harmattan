package handler_test

import (
	"testing"
	"time"

	"github.com/ikawaha/httpcheck"
	"github.com/minguu42/harmattan/internal/api/handler"
	"github.com/minguu42/harmattan/internal/api/openapi"
	"github.com/minguu42/harmattan/internal/database"
	"github.com/stretchr/testify/require"
)

func TestHandler_CreateTask(t *testing.T) {
	require.NoError(t, tdb.Reset(t.Context(), []any{database.Project{}, database.Task{}, database.Step{}, database.Tag{}}))
	require.NoError(t, tdb.Insert(t.Context(), []any{
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
	}))

	t.Run("project not found", func(t *testing.T) {
		want := handler.ErrorResponse{Code: 404, Message: "指定したプロジェクトは見つかりません"}
		httpcheck.New(th).Test(t, "POST", "/projects/PROJECT-000000000000000099/tasks").
			WithHeader("Authorization", "Bearer "+token).
			WithHeader("Content-Type", "application/json").
			WithBody([]byte(`{"name": "タスク", "priority": 1}`)).
			Check().HasStatus(404).HasJSON(want)
	})
	t.Run("project access denied", func(t *testing.T) {
		want := handler.ErrorResponse{Code: 404, Message: "指定したプロジェクトは見つかりません"}
		httpcheck.New(th).Test(t, "POST", "/projects/PROJECT-000000000000000002/tasks").
			WithHeader("Authorization", "Bearer "+token).
			WithHeader("Content-Type", "application/json").
			WithBody([]byte(`{"name": "タスク", "priority": 1}`)).
			Check().HasStatus(404).HasJSON(want)
	})
	t.Run("ok", func(t *testing.T) {
		want := &openapi.Task{
			ID:        fixedID,
			ProjectID: "PROJECT-000000000000000001",
			Name:      "タスク",
			Priority:  1,
			CreatedAt: fixedNow,
			UpdatedAt: fixedNow,
		}
		httpcheck.New(th).Test(t, "POST", "/projects/PROJECT-000000000000000001/tasks").
			WithHeader("Authorization", "Bearer "+token).
			WithHeader("Content-Type", "application/json").
			WithBody([]byte(`{"name": "タスク", "priority": 1}`)).
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
	t.Run("project access denied", func(t *testing.T) {
		want := handler.ErrorResponse{Code: 404, Message: "指定したプロジェクトは見つかりません"}
		httpcheck.New(th).Test(t, "GET", "/projects/PROJECT-000000000000000002/tasks").
			WithHeader("Authorization", "Bearer "+token).
			Check().HasStatus(404).HasJSON(want)
	})
}

func TestHandler_GetTask(t *testing.T) {
	require.NoError(t, tdb.Reset(t.Context(), []any{database.Project{}, database.Task{}, database.Step{}, database.Tag{}}))
	require.NoError(t, tdb.Insert(t.Context(), []any{
		database.Projects{
			{
				ID:        "PROJECT-000000000000000001",
				UserID:    testUserID,
				Name:      "プロジェクト1",
				Color:     "blue",
				CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
				UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
			},
			{
				ID:        "PROJECT-000000000000000002",
				UserID:    "USER-000000000000000000002",
				Name:      "プロジェクト2",
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
				Name:      "タスク1",
				Content:   "内容",
				Priority:  1,
				CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
				UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
			},
			{
				ID:        "TASK-000000000000000000002",
				UserID:    "USER-000000000000000000002",
				ProjectID: "PROJECT-000000000000000002",
				Name:      "タスク2",
				Content:   "内容",
				Priority:  2,
				CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
				UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
			},
		},
	}))

	t.Run("task not found", func(t *testing.T) {
		want := handler.ErrorResponse{Code: 404, Message: "指定したタスクは見つかりません"}
		httpcheck.New(th).Test(t, "GET", "/tasks/TASK-000000000000000000099").
			WithHeader("Authorization", "Bearer "+token).
			Check().HasStatus(404).HasJSON(want)
	})
	t.Run("task access denied", func(t *testing.T) {
		want := handler.ErrorResponse{Code: 404, Message: "指定したタスクは見つかりません"}
		httpcheck.New(th).Test(t, "GET", "/tasks/TASK-000000000000000000002").
			WithHeader("Authorization", "Bearer "+token).
			Check().HasStatus(404).HasJSON(want)
	})
	t.Run("ok", func(t *testing.T) {
		want := &openapi.Task{
			ID:        "TASK-000000000000000000001",
			ProjectID: "PROJECT-000000000000000001",
			Name:      "タスク1",
			Content:   "内容",
			Priority:  1,
			CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
			UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
		}
		httpcheck.New(th).Test(t, "GET", "/tasks/TASK-000000000000000000001").
			WithHeader("Authorization", "Bearer "+token).
			Check().HasStatus(200).HasJSON(want)
	})
}

func TestHandler_UpdateTask(t *testing.T) {
	require.NoError(t, tdb.Reset(t.Context(), []any{database.Project{}, database.Task{}, database.Step{}, database.Tag{}, database.TaskTag{}}))
	require.NoError(t, tdb.Insert(t.Context(), []any{
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
				CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
				UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
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
				CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
				UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
			},
			{
				ID:     "TAG-0000000000000000000002",
				UserID: "USER-000000000000000000002",
				Name:   "タグ2",
			},
		},
	}))

	t.Run("task not found", func(t *testing.T) {
		want := handler.ErrorResponse{Code: 404, Message: "指定したタスクは見つかりません"}
		httpcheck.New(th).Test(t, "PATCH", "/tasks/TASK-000000000000000000099").
			WithHeader("Authorization", "Bearer "+token).
			WithHeader("Content-Type", "application/json").
			WithBody([]byte(`{"name": "更新後タスク", "content": "更新後内容", "priority": 3}`)).
			Check().HasStatus(404).HasJSON(want)
	})
	t.Run("task access denied", func(t *testing.T) {
		want := handler.ErrorResponse{Code: 404, Message: "指定したタスクは見つかりません"}
		httpcheck.New(th).Test(t, "PATCH", "/tasks/TASK-000000000000000000002").
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
			Tags: []openapi.Tag{
				{
					ID:        "TAG-0000000000000000000001",
					Name:      "タグ1",
					CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
					UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
				},
			},
		}
		httpcheck.New(th).Test(t, "PATCH", "/tasks/TASK-000000000000000000001").
			WithHeader("Authorization", "Bearer "+token).
			WithHeader("Content-Type", "application/json").
			WithBody([]byte(`{"name": "更新後タスク", "content": "更新後内容", "priority": 3, "due_on": "2025-01-02", "completed_at": "2025-01-01T12:00:00+09:00", "tag_ids": ["TAG-0000000000000000000001", "TAG-0000000000000000000002", "TAG-0000000000000000000099"]}`)).
			Check().HasStatus(200).HasJSON(want)
	})
}

func TestHandler_DeleteTask(t *testing.T) {
	require.NoError(t, tdb.Reset(t.Context(), []any{database.Project{}, database.Task{}, database.Step{}, database.Tag{}}))
	require.NoError(t, tdb.Insert(t.Context(), []any{
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
				CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
				UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
			},
			{
				ID:        "TASK-000000000000000000002",
				UserID:    "USER-000000000000000000002",
				ProjectID: "PROJECT-000000000000000002",
				Name:      "タスク2",
			},
		},
	}))

	t.Run("task not found", func(t *testing.T) {
		want := handler.ErrorResponse{Code: 404, Message: "指定したタスクは見つかりません"}
		httpcheck.New(th).Test(t, "DELETE", "/tasks/TASK-000000000000000000099").
			WithHeader("Authorization", "Bearer "+token).
			Check().HasStatus(404).HasJSON(want)
	})
	t.Run("task access denied", func(t *testing.T) {
		want := handler.ErrorResponse{Code: 404, Message: "指定したタスクは見つかりません"}
		httpcheck.New(th).Test(t, "DELETE", "/tasks/TASK-000000000000000000002").
			WithHeader("Authorization", "Bearer "+token).
			Check().HasStatus(404).HasJSON(want)
	})
	t.Run("ok", func(t *testing.T) {
		httpcheck.New(th).Test(t, "DELETE", "/tasks/TASK-000000000000000000001").
			WithHeader("Authorization", "Bearer "+token).
			Check().HasStatus(200).HasString("")
	})
}
