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

func TestHandler_CreateStep(t *testing.T) {
	require.NoError(t, tdb.Reset(t.Context(), []any{database.Project{}, database.Task{}, database.Step{}}))
	require.NoError(t, tdb.Insert(t.Context(), []any{
		database.Project{
			ID:        "PROJECT-000000000000000001",
			UserID:    testUserID,
			Name:      "テストプロジェクト",
			Color:     "#1E3A8A",
			CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
			UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
		},
		database.Task{
			ID:        "TASK-000000000000000000001",
			UserID:    testUserID,
			ProjectID: "PROJECT-000000000000000001",
			Name:      "テストタスク",
			Content:   "テストタスクの内容",
			Priority:  1,
			CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
			UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
		},
	}))

	want := &openapi.Step{
		ID:        fixedID,
		TaskID:    "TASK-000000000000000000001",
		Name:      "テストステップ",
		CreatedAt: fixedNow,
		UpdatedAt: fixedNow,
	}
	httpcheck.New(th).Test(t, "POST", "/projects/PROJECT-000000000000000001/tasks/TASK-000000000000000000001/steps").
		WithHeader("Authorization", "Bearer "+token).
		WithHeader("Content-Type", "application/json").
		WithBody([]byte(`{"name": "テストステップ"}`)).
		Check().HasStatus(200).HasJSON(want)
}

func TestHandler_UpdateStep(t *testing.T) {
	require.NoError(t, tdb.Reset(t.Context(), []any{database.Project{}, database.Task{}, database.Step{}}))
	require.NoError(t, tdb.Insert(t.Context(), []any{
		database.Project{
			ID:        "PROJECT-000000000000000001",
			UserID:    testUserID,
			Name:      "テストプロジェクト",
			Color:     "#1E3A8A",
			CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
			UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
		},
		database.Task{
			ID:        "TASK-000000000000000000001",
			UserID:    testUserID,
			ProjectID: "PROJECT-000000000000000001",
			Name:      "テストタスク",
			Content:   "テストタスクの内容",
			Priority:  1,
			CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
			UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
		},
		database.Step{
			ID:        "STEP-000000000000000000001",
			UserID:    testUserID,
			TaskID:    "TASK-000000000000000000001",
			Name:      "テストステップ",
			CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
			UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
		},
	}))

	want := &openapi.Step{
		ID:        "STEP-000000000000000000001",
		TaskID:    "TASK-000000000000000000001",
		Name:      "更新されたステップ",
		CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
		UpdatedAt: fixedNow,
	}
	httpcheck.New(th).Test(t, "PATCH", "/projects/PROJECT-000000000000000001/tasks/TASK-000000000000000000001/steps/STEP-000000000000000000001").
		WithHeader("Authorization", "Bearer "+token).
		WithHeader("Content-Type", "application/json").
		WithBody([]byte(`{"name": "更新されたステップ"}`)).
		Check().HasStatus(200).HasJSON(want)
}

func TestHandler_UpdateStep_Validation(t *testing.T) {
	require.NoError(t, tdb.Reset(t.Context(), []any{database.Project{}, database.Task{}, database.Step{}}))
	require.NoError(t, tdb.Insert(t.Context(), []any{
		database.Projects{
			{
				ID:        "PROJECT-000000000000000001",
				UserID:    testUserID,
				Name:      "テストプロジェクト",
				Color:     "#1E3A8A",
				CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
				UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
			},
			{
				ID:        "PROJECT-000000000000000002",
				UserID:    testUserID,
				Name:      "テストプロジェクト2",
				Color:     "#1E3A8A",
				CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
				UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
			},
			{
				ID:        "PROJECT-000000000000000003",
				UserID:    "USER-000000000000000000002",
				Name:      "他のユーザーのプロジェクト",
				Color:     "#1E3A8A",
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
				Content:   "テストタスクの内容",
				Priority:  1,
				CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
				UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
			},
			{
				ID:        "TASK-000000000000000000002",
				UserID:    testUserID,
				ProjectID: "PROJECT-000000000000000002",
				Name:      "別プロジェクトのタスク",
				Content:   "",
				Priority:  1,
				CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
				UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
			},
		},
		database.Steps{
			{
				ID:        "STEP-000000000000000000001",
				UserID:    testUserID,
				TaskID:    "TASK-000000000000000000001",
				Name:      "テストステップ",
				CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
				UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
			},
		},
	}))

	t.Run("project not found", func(t *testing.T) {
		want := handler.ErrorResponse{Code: 404, Message: "指定したプロジェクトは見つかりません"}
		httpcheck.New(th).Test(t, "PATCH", "/projects/PROJECT-000000000000000099/tasks/TASK-000000000000000000001/steps/STEP-000000000000000000001").
			WithHeader("Authorization", "Bearer "+token).
			WithHeader("Content-Type", "application/json").
			WithBody([]byte(`{"name": "更新されたステップ"}`)).
			Check().HasStatus(404).HasJSON(want)
	})

	t.Run("user does not own the project", func(t *testing.T) {
		want := handler.ErrorResponse{Code: 404, Message: "指定したプロジェクトは見つかりません"}
		httpcheck.New(th).Test(t, "PATCH", "/projects/PROJECT-000000000000000003/tasks/TASK-000000000000000000001/steps/STEP-000000000000000000001").
			WithHeader("Authorization", "Bearer "+token).
			WithHeader("Content-Type", "application/json").
			WithBody([]byte(`{"name": "更新されたステップ"}`)).
			Check().HasStatus(404).HasJSON(want)
	})

	t.Run("task not found", func(t *testing.T) {
		want := handler.ErrorResponse{Code: 404, Message: "指定したタスクは見つかりません"}
		httpcheck.New(th).Test(t, "PATCH", "/projects/PROJECT-000000000000000001/tasks/TASK-000000000000000000099/steps/STEP-000000000000000000001").
			WithHeader("Authorization", "Bearer "+token).
			WithHeader("Content-Type", "application/json").
			WithBody([]byte(`{"name": "更新されたステップ"}`)).
			Check().HasStatus(404).HasJSON(want)
	})

	t.Run("task does not belong to the project", func(t *testing.T) {
		want := handler.ErrorResponse{Code: 404, Message: "指定したタスクは見つかりません"}
		httpcheck.New(th).Test(t, "PATCH", "/projects/PROJECT-000000000000000001/tasks/TASK-000000000000000000002/steps/STEP-000000000000000000001").
			WithHeader("Authorization", "Bearer "+token).
			WithHeader("Content-Type", "application/json").
			WithBody([]byte(`{"name": "更新されたステップ"}`)).
			Check().HasStatus(404).HasJSON(want)
	})

	t.Run("step not found", func(t *testing.T) {
		want := handler.ErrorResponse{Code: 404, Message: "指定したステップは見つかりません"}
		httpcheck.New(th).Test(t, "PATCH", "/projects/PROJECT-000000000000000001/tasks/TASK-000000000000000000001/steps/STEP-000000000000000000099").
			WithHeader("Authorization", "Bearer "+token).
			WithHeader("Content-Type", "application/json").
			WithBody([]byte(`{"name": "更新されたステップ"}`)).
			Check().HasStatus(404).HasJSON(want)
	})
}

func TestHandler_DeleteStep(t *testing.T) {
	require.NoError(t, tdb.Reset(t.Context(), []any{database.Project{}, database.Task{}, database.Step{}}))
	require.NoError(t, tdb.Insert(t.Context(), []any{
		database.Project{
			ID:        "PROJECT-000000000000000001",
			UserID:    testUserID,
			Name:      "テストプロジェクト",
			Color:     "#1E3A8A",
			CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
			UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
		},
		database.Task{
			ID:        "TASK-000000000000000000001",
			UserID:    testUserID,
			ProjectID: "PROJECT-000000000000000001",
			Name:      "テストタスク",
			Content:   "テストタスクの内容",
			Priority:  1,
			CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
			UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
		},
		database.Step{
			ID:        "STEP-000000000000000000001",
			UserID:    testUserID,
			TaskID:    "TASK-000000000000000000001",
			Name:      "テストステップ",
			CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
			UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
		},
	}))

	httpcheck.New(th).Test(t, "DELETE", "/projects/PROJECT-000000000000000001/tasks/TASK-000000000000000000001/steps/STEP-000000000000000000001").
		WithHeader("Authorization", "Bearer "+token).
		Check().HasStatus(200)
}

func TestHandler_DeleteStep_Validation(t *testing.T) {
	require.NoError(t, tdb.Reset(t.Context(), []any{database.Project{}, database.Task{}, database.Step{}}))
	require.NoError(t, tdb.Insert(t.Context(), []any{
		database.Projects{
			{
				ID:        "PROJECT-000000000000000001",
				UserID:    testUserID,
				Name:      "テストプロジェクト",
				Color:     "#1E3A8A",
				CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
				UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
			},
			{
				ID:        "PROJECT-000000000000000002",
				UserID:    testUserID,
				Name:      "テストプロジェクト2",
				Color:     "#1E3A8A",
				CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
				UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
			},
			{
				ID:        "PROJECT-000000000000000003",
				UserID:    "USER-000000000000000000002",
				Name:      "他のユーザーのプロジェクト",
				Color:     "#1E3A8A",
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
				Content:   "テストタスクの内容",
				Priority:  1,
				CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
				UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
			},
			{
				ID:        "TASK-000000000000000000002",
				UserID:    testUserID,
				ProjectID: "PROJECT-000000000000000002",
				Name:      "別プロジェクトのタスク",
				Content:   "",
				Priority:  1,
				CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
				UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
			},
		},
		database.Steps{
			{
				ID:        "STEP-000000000000000000001",
				UserID:    testUserID,
				TaskID:    "TASK-000000000000000000001",
				Name:      "テストステップ",
				CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
				UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
			},
		},
	}))

	t.Run("project not found", func(t *testing.T) {
		want := handler.ErrorResponse{Code: 404, Message: "指定したプロジェクトは見つかりません"}
		httpcheck.New(th).Test(t, "DELETE", "/projects/PROJECT-000000000000000099/tasks/TASK-000000000000000000001/steps/STEP-000000000000000000001").
			WithHeader("Authorization", "Bearer "+token).
			Check().HasStatus(404).HasJSON(want)
	})
	t.Run("user does not own the project", func(t *testing.T) {
		want := handler.ErrorResponse{Code: 404, Message: "指定したプロジェクトは見つかりません"}
		httpcheck.New(th).Test(t, "DELETE", "/projects/PROJECT-000000000000000003/tasks/TASK-000000000000000000001/steps/STEP-000000000000000000001").
			WithHeader("Authorization", "Bearer "+token).
			Check().HasStatus(404).HasJSON(want)
	})
	t.Run("task not found", func(t *testing.T) {
		want := handler.ErrorResponse{Code: 404, Message: "指定したタスクは見つかりません"}
		httpcheck.New(th).Test(t, "DELETE", "/projects/PROJECT-000000000000000001/tasks/TASK-000000000000000000099/steps/STEP-000000000000000000001").
			WithHeader("Authorization", "Bearer "+token).
			Check().HasStatus(404).HasJSON(want)
	})
	t.Run("task does not belong to the project", func(t *testing.T) {
		want := handler.ErrorResponse{Code: 404, Message: "指定したタスクは見つかりません"}
		httpcheck.New(th).Test(t, "DELETE", "/projects/PROJECT-000000000000000001/tasks/TASK-000000000000000000002/steps/STEP-000000000000000000001").
			WithHeader("Authorization", "Bearer "+token).
			Check().HasStatus(404).HasJSON(want)
	})
	t.Run("step not found", func(t *testing.T) {
		want := handler.ErrorResponse{Code: 404, Message: "指定したステップは見つかりません"}
		httpcheck.New(th).Test(t, "DELETE", "/projects/PROJECT-000000000000000001/tasks/TASK-000000000000000000001/steps/STEP-000000000000000000099").
			WithHeader("Authorization", "Bearer "+token).
			Check().HasStatus(404).HasJSON(want)
	})
}
