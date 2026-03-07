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

func TestStep_CreateStep(t *testing.T) {
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

	tests := map[string]test{
		"ok": {
			Method:     "POST",
			Path:       "/tasks/TASK-000000000000000000001/steps",
			Headers:    http.Header{"Authorization": []string{"Bearer " + token}},
			Body:       `{"name": "ステップ"}`,
			WantStatus: 200,
			WantJSON: &openapi.Step{
				ID:        fixedID,
				TaskID:    "TASK-000000000000000000001",
				Name:      "ステップ",
				CreatedAt: fixedNow,
				UpdatedAt: fixedNow,
			},
			WantTables: []any{database.Steps{
				{ID: fixedID, UserID: testUserID, TaskID: "TASK-000000000000000000001", Name: "ステップ", CreatedAt: fixedNow, UpdatedAt: fixedNow},
			}},
		},
		"not_found": {
			Method:     "POST",
			Path:       "/tasks/TASK-000000000000000000099/steps",
			Headers:    http.Header{"Authorization": []string{"Bearer " + token}},
			Body:       `{"name": "ステップ"}`,
			WantStatus: 404,
			WantJSON:   handler.ErrorResponse{Code: 404, Message: "指定したタスクは見つかりません"},
		},
		"access_denied": {
			Method:     "POST",
			Path:       "/tasks/TASK-000000000000000000002/steps",
			Headers:    http.Header{"Authorization": []string{"Bearer " + token}},
			Body:       `{"name": "ステップ"}`,
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

func TestStep_UpdateStep(t *testing.T) {
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
				CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst),
			},
			{
				ID:     "STEP-000000000000000000002",
				UserID: "USER-000000000000000000002",
				TaskID: "TASK-000000000000000000002",
				Name:   "ステップ2",
			},
		},
	}))

	tests := map[string]test{
		"ok": {
			Method:     "PATCH",
			Path:       "/steps/STEP-000000000000000000001",
			Headers:    http.Header{"Authorization": []string{"Bearer " + token}},
			Body:       `{"name": "更新後ステップ"}`,
			WantStatus: 200,
			WantJSON: &openapi.Step{
				ID:        "STEP-000000000000000000001",
				TaskID:    "TASK-000000000000000000001",
				Name:      "更新後ステップ",
				CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst),
				UpdatedAt: fixedNow,
			},
		},
		"not_found": {
			Method:     "PATCH",
			Path:       "/steps/STEP-000000000000000000099",
			Headers:    http.Header{"Authorization": []string{"Bearer " + token}},
			Body:       `{"name": "更新後ステップ"}`,
			WantStatus: 404,
			WantJSON:   handler.ErrorResponse{Code: 404, Message: "指定したステップは見つかりません"},
		},
		"access_denied": {
			Method:     "PATCH",
			Path:       "/steps/STEP-000000000000000000002",
			Headers:    http.Header{"Authorization": []string{"Bearer " + token}},
			Body:       `{"name": "更新後ステップ"}`,
			WantStatus: 404,
			WantJSON:   handler.ErrorResponse{Code: 404, Message: "指定したステップは見つかりません"},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			runTest(t, tt)
		})
	}
}

func TestStep_DeleteStep(t *testing.T) {
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
				CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst),
				UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst),
			},
			{
				ID:        "STEP-000000000000000000002",
				UserID:    "USER-000000000000000000002",
				TaskID:    "TASK-000000000000000000002",
				Name:      "ステップ2",
				CreatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst),
				UpdatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst),
			},
		},
	}))

	tests := map[string]test{
		"ok": {
			Method:     "DELETE",
			Path:       "/steps/STEP-000000000000000000001",
			Headers:    http.Header{"Authorization": []string{"Bearer " + token}},
			WantStatus: 200,
			WantTables: []any{database.Steps{
				{ID: "STEP-000000000000000000002", UserID: "USER-000000000000000000002", TaskID: "TASK-000000000000000000002", Name: "ステップ2", CreatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst)},
			}},
		},
		"not_found": {
			Method:     "DELETE",
			Path:       "/steps/STEP-000000000000000000099",
			Headers:    http.Header{"Authorization": []string{"Bearer " + token}},
			WantStatus: 404,
			WantJSON:   handler.ErrorResponse{Code: 404, Message: "指定したステップは見つかりません"},
		},
		"access_denied": {
			Method:     "DELETE",
			Path:       "/steps/STEP-000000000000000000002",
			Headers:    http.Header{"Authorization": []string{"Bearer " + token}},
			WantStatus: 404,
			WantJSON:   handler.ErrorResponse{Code: 404, Message: "指定したステップは見つかりません"},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			runTest(t, tt)
		})
	}
}
