package handler_test

import (
	"testing"
	"time"

	"github.com/minguu42/harmattan/internal/api/handler"
	"github.com/minguu42/harmattan/internal/api/openapi"
	"github.com/minguu42/harmattan/internal/database"
	"github.com/minguu42/harmattan/internal/lib/httpcheck"
	"github.com/stretchr/testify/require"
)

func TestHandler_CreateStep(t *testing.T) {
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
			},
			{
				ID:        "TASK-000000000000000000002",
				UserID:    "USER-000000000000000000002",
				ProjectID: "PROJECT-000000000000000002",
				Name:      "タスク2",
			},
		},
		database.Steps{},
	}))

	t.Run("task not found", func(t *testing.T) {
		want := handler.ErrorResponse{Code: 404, Message: "指定したタスクは見つかりません"}
		httpcheck.New(th).Test(t, "POST", "/tasks/TASK-000000000000000000099/steps").
			WithHeader("Authorization", "Bearer "+token).
			WithHeader("Content-Type", "application/json").
			WithBody([]byte(`{"name": "ステップ"}`)).
			Check().HasStatus(404).HasJSON(want)
	})
	t.Run("task access denied", func(t *testing.T) {
		want := handler.ErrorResponse{Code: 404, Message: "指定したタスクは見つかりません"}
		httpcheck.New(th).Test(t, "POST", "/tasks/TASK-000000000000000000002/steps").
			WithHeader("Authorization", "Bearer "+token).
			WithHeader("Content-Type", "application/json").
			WithBody([]byte(`{"name": "ステップ"}`)).
			Check().HasStatus(404).HasJSON(want)
	})
	t.Run("ok", func(t *testing.T) {
		want := &openapi.Step{
			ID:        fixedID,
			TaskID:    "TASK-000000000000000000001",
			Name:      "ステップ",
			CreatedAt: fixedNow,
			UpdatedAt: fixedNow,
		}
		httpcheck.New(th).Test(t, "POST", "/tasks/TASK-000000000000000000001/steps").
			WithHeader("Authorization", "Bearer "+token).
			WithHeader("Content-Type", "application/json").
			WithBody([]byte(`{"name": "ステップ"}`)).
			Check().HasStatus(200).HasJSON(want)
	})
}

func TestHandler_UpdateStep(t *testing.T) {
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
			},
			{
				ID:        "TASK-000000000000000000002",
				UserID:    "USER-000000000000000000002",
				ProjectID: "PROJECT-000000000000000002",
				Name:      "タスク2",
			},
		},
		database.Steps{
			{
				ID:        "STEP-000000000000000000001",
				UserID:    testUserID,
				TaskID:    "TASK-000000000000000000001",
				Name:      "ステップ1",
				CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
			},
			{
				ID:     "STEP-000000000000000000002",
				UserID: "USER-000000000000000000002",
				TaskID: "TASK-000000000000000000002",
				Name:   "ステップ2",
			},
		},
	}))

	t.Run("step not found", func(t *testing.T) {
		want := handler.ErrorResponse{Code: 404, Message: "指定したステップは見つかりません"}
		httpcheck.New(th).Test(t, "PATCH", "/steps/STEP-000000000000000000099").
			WithHeader("Authorization", "Bearer "+token).
			WithHeader("Content-Type", "application/json").
			WithBody([]byte(`{"name": "更新後ステップ"}`)).
			Check().HasStatus(404).HasJSON(want)
	})
	t.Run("step access denied", func(t *testing.T) {
		want := handler.ErrorResponse{Code: 404, Message: "指定したステップは見つかりません"}
		httpcheck.New(th).Test(t, "PATCH", "/steps/STEP-000000000000000000002").
			WithHeader("Authorization", "Bearer "+token).
			WithHeader("Content-Type", "application/json").
			WithBody([]byte(`{"name": "更新後ステップ"}`)).
			Check().HasStatus(404).HasJSON(want)
	})
	t.Run("ok", func(t *testing.T) {
		want := &openapi.Step{
			ID:        "STEP-000000000000000000001",
			TaskID:    "TASK-000000000000000000001",
			Name:      "更新後ステップ",
			CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
			UpdatedAt: fixedNow,
		}
		httpcheck.New(th).Test(t, "PATCH", "/steps/STEP-000000000000000000001").
			WithHeader("Authorization", "Bearer "+token).
			WithHeader("Content-Type", "application/json").
			WithBody([]byte(`{"name": "更新後ステップ"}`)).
			Check().HasStatus(200).HasJSON(want)
	})
}

func TestHandler_DeleteStep(t *testing.T) {
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
			},
			{
				ID:        "TASK-000000000000000000002",
				UserID:    "USER-000000000000000000002",
				ProjectID: "PROJECT-000000000000000002",
				Name:      "タスク2",
			},
		},
		database.Steps{
			{
				ID:     "STEP-000000000000000000001",
				UserID: testUserID,
				TaskID: "TASK-000000000000000000001",
				Name:   "ステップ1",
			},
			{
				ID:     "STEP-000000000000000000002",
				UserID: "USER-000000000000000000002",
				TaskID: "TASK-000000000000000000002",
				Name:   "ステップ2",
			},
		},
	}))

	t.Run("step not found", func(t *testing.T) {
		want := handler.ErrorResponse{Code: 404, Message: "指定したステップは見つかりません"}
		httpcheck.New(th).Test(t, "DELETE", "/steps/STEP-000000000000000000099").
			WithHeader("Authorization", "Bearer "+token).
			Check().HasStatus(404).HasJSON(want)
	})
	t.Run("step access denied", func(t *testing.T) {
		want := handler.ErrorResponse{Code: 404, Message: "指定したステップは見つかりません"}
		httpcheck.New(th).Test(t, "DELETE", "/steps/STEP-000000000000000000002").
			WithHeader("Authorization", "Bearer "+token).
			Check().HasStatus(404).HasJSON(want)
	})
	t.Run("ok", func(t *testing.T) {
		httpcheck.New(th).Test(t, "DELETE", "/steps/STEP-000000000000000000001").
			WithHeader("Authorization", "Bearer "+token).
			Check().HasStatus(200).HasString("")
	})
}
