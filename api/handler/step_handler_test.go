package handler_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/ikawaha/httpcheck"
	"github.com/minguu42/harmattan/internal/database"
	"github.com/minguu42/harmattan/internal/domain"
	"github.com/minguu42/harmattan/internal/openapi"
	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/require"
)

func TestHandler_CreateStep(t *testing.T) {
	projectID := domain.ProjectID(ulid.Make().String())
	taskID := domain.TaskID(ulid.Make().String())
	require.NoError(t, tdb.Reset(t.Context(), []any{database.User{}, database.Project{}, database.Task{}, database.Step{}}))
	require.NoError(t, tdb.Insert(t.Context(), []any{
		database.User{
			ID:             "user_01",
			Email:          "user1@example.com",
			HashedPassword: "password",
			CreatedAt:      time.Date(2025, 1, 1, 0, 0, 0, 0, time.Local),
			UpdatedAt:      time.Date(2025, 1, 1, 0, 0, 0, 0, time.Local),
		},
		database.Project{
			ID:        projectID,
			UserID:    "user_01",
			Name:      "テストプロジェクト",
			Color:     "#1E3A8A",
			CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, time.Local),
			UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, time.Local),
		},
		database.Task{
			ID:        taskID,
			UserID:    "user_01",
			ProjectID: projectID,
			Name:      "テストタスク",
			Content:   "テストタスクの内容",
			Priority:  1,
			CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, time.Local),
			UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, time.Local),
		},
	}))

	wantResponse := &openapi.Step{
		ID:     "01JGFJJZ000000000000000000",
		TaskID: string(taskID),
		Name:   "テストステップ",
	}
	checker := httpcheck.New(fixTimeMiddleware(h, now))
	checker.Test(t, "POST", fmt.Sprintf("/projects/%s/tasks/%s/steps", projectID, taskID)).
		WithHeader("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ1c2VyXzAxIiwiZXhwIjoxNzM1NjYwODAwLCJpYXQiOjE3MzU2NTcyMDB9.bT7pyLGRAxG784_cg1DoZ9GD3GbGbNFichSlETzYfPc").
		WithHeader("Content-Type", "application/json").
		WithBody([]byte(`{"name": "テストステップ"}`)).
		Check().
		HasStatus(200).
		HasJSON(wantResponse)
}
