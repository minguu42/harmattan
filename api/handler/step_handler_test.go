package handler_test

import (
	"testing"
	"time"

	"github.com/ikawaha/httpcheck"
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
		ID:     fixedID,
		TaskID: "TASK-000000000000000000001",
		Name:   "テストステップ",
	}
	httpcheck.New(th).Test(t, "POST", "/projects/PROJECT-000000000000000001/tasks/TASK-000000000000000000001/steps").
		WithHeader("Authorization", "Bearer "+token).
		WithHeader("Content-Type", "application/json").
		WithBody([]byte(`{"name": "テストステップ"}`)).
		Check().HasStatus(200).HasJSON(want)
}
