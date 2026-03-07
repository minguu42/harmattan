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

func TestProject_CreateProject(t *testing.T) {
	require.NoError(t, tdb.TruncateAndInsert(t.Context(), []any{database.Projects{}}))

	runTest(t, test{
		Method:     "POST",
		Path:       "/projects",
		Headers:    http.Header{"Authorization": []string{"Bearer " + token}},
		Body:       `{"name": "プロジェクト1", "color": "blue"}`,
		WantStatus: 200,
		WantJSON: &openapi.Project{
			ID:        fixedID,
			Name:      "プロジェクト1",
			Color:     "blue",
			CreatedAt: fixedNow,
			UpdatedAt: fixedNow,
		},
		WantTables: []any{database.Projects{
			{ID: fixedID, UserID: testUserID, Name: "プロジェクト1", Color: "blue", CreatedAt: fixedNow, UpdatedAt: fixedNow},
		}},
	})
}

func TestProject_ListProjects(t *testing.T) {
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
				UserID:    testUserID,
				Name:      "プロジェクト2",
				Color:     "gray",
				CreatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst),
				UpdatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst),
			},
		},
	}))

	tests := map[string]test{
		"no_limit_and_offset": {
			Method:     "GET",
			Path:       "/projects",
			Headers:    http.Header{"Authorization": []string{"Bearer " + token}},
			WantStatus: 200,
			WantJSON: &openapi.ListProjectsOK{
				Projects: []openapi.Project{
					{
						ID:        "PROJECT-000000000000000001",
						Name:      "プロジェクト1",
						Color:     "blue",
						CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst),
						UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst),
					},
					{
						ID:        "PROJECT-000000000000000002",
						Name:      "プロジェクト2",
						Color:     "gray",
						CreatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst),
						UpdatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst),
					},
				},
				HasNext: false,
			},
		},
		"limit_1_offset_0": {
			Method:     "GET",
			Path:       "/projects?limit=1&offset=0",
			Headers:    http.Header{"Authorization": []string{"Bearer " + token}},
			WantStatus: 200,
			WantJSON: &openapi.ListProjectsOK{
				Projects: []openapi.Project{
					{
						ID:        "PROJECT-000000000000000001",
						Name:      "プロジェクト1",
						Color:     "blue",
						CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst),
						UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst),
					},
				},
				HasNext: true,
			},
		},
		"limit_1_offset_1": {
			Method:     "GET",
			Path:       "/projects?limit=1&offset=1",
			Headers:    http.Header{"Authorization": []string{"Bearer " + token}},
			WantStatus: 200,
			WantJSON: &openapi.ListProjectsOK{
				Projects: []openapi.Project{
					{
						ID:        "PROJECT-000000000000000002",
						Name:      "プロジェクト2",
						Color:     "gray",
						CreatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst),
						UpdatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst),
					},
				},
				HasNext: false,
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			runTest(t, tt)
		})
	}
}

func TestProject_GetProject(t *testing.T) {
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
	}))

	tests := map[string]test{
		"ok": {
			Method:     "GET",
			Path:       "/projects/PROJECT-000000000000000001",
			Headers:    http.Header{"Authorization": []string{"Bearer " + token}},
			WantStatus: 200,
			WantJSON: &openapi.Project{
				ID:        "PROJECT-000000000000000001",
				Name:      "プロジェクト1",
				Color:     "blue",
				CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst),
				UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst),
			},
		},
		"not_found": {
			Method:     "GET",
			Path:       "/projects/PROJECT-000000000000000099",
			Headers:    http.Header{"Authorization": []string{"Bearer " + token}},
			WantStatus: 404,
			WantJSON:   handler.ErrorResponse{Code: 404, Message: "指定したプロジェクトは見つかりません"},
		},
		"access_denied": {
			Method:     "GET",
			Path:       "/projects/PROJECT-000000000000000002",
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

func TestProject_UpdateProject(t *testing.T) {
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
	}))

	tests := map[string]test{
		"ok": {
			Method:     "PATCH",
			Path:       "/projects/PROJECT-000000000000000001",
			Headers:    http.Header{"Authorization": []string{"Bearer " + token}},
			Body:       `{"name": "更新後プロジェクト", "color": "gray", "is_archived": true}`,
			WantStatus: 200,
			WantJSON: &openapi.Project{
				ID:         "PROJECT-000000000000000001",
				Name:       "更新後プロジェクト",
				Color:      "gray",
				IsArchived: true,
				CreatedAt:  time.Date(2025, 1, 1, 0, 0, 1, 0, jst),
				UpdatedAt:  fixedNow,
			},
			WantTables: []any{database.Projects{
				{ID: "PROJECT-000000000000000001", UserID: testUserID, Name: "更新後プロジェクト", Color: "gray", IsArchived: true, CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst), UpdatedAt: fixedNow},
				{ID: "PROJECT-000000000000000002", UserID: "USER-000000000000000000002", Name: "プロジェクト2", Color: "gray", CreatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst)},
			}},
		},
		"not_found": {
			Method:     "PATCH",
			Path:       "/projects/PROJECT-000000000000000099",
			Headers:    http.Header{"Authorization": []string{"Bearer " + token}},
			Body:       `{"name": "更新後プロジェクト", "color": "gray", "is_archived": true}`,
			WantStatus: 404,
			WantJSON:   handler.ErrorResponse{Code: 404, Message: "指定したプロジェクトは見つかりません"},
		},
		"access_denied": {
			Method:     "PATCH",
			Path:       "/projects/PROJECT-000000000000000002",
			Headers:    http.Header{"Authorization": []string{"Bearer " + token}},
			Body:       `{"name": "更新後プロジェクト", "color": "gray", "is_archived": true}`,
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

func TestProject_DeleteProject(t *testing.T) {
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
	}))

	tests := map[string]test{
		"ok": {
			Method:     "DELETE",
			Path:       "/projects/PROJECT-000000000000000001",
			Headers:    http.Header{"Authorization": []string{"Bearer " + token}},
			WantStatus: 200,
			WantTables: []any{database.Projects{
				{ID: "PROJECT-000000000000000002", UserID: "USER-000000000000000000002", Name: "プロジェクト2", Color: "gray", CreatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst)},
			}},
		},
		"not_found": {
			Method:     "DELETE",
			Path:       "/projects/PROJECT-000000000000000099",
			Headers:    http.Header{"Authorization": []string{"Bearer " + token}},
			WantStatus: 404,
			WantJSON:   handler.ErrorResponse{Code: 404, Message: "指定したプロジェクトは見つかりません"},
		},
		"access_denied": {
			Method:     "DELETE",
			Path:       "/projects/PROJECT-000000000000000002",
			Headers:    http.Header{"Authorization": []string{"Bearer " + token}},
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
