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

func TestHandler_CreateProject(t *testing.T) {
	require.NoError(t, tdb.TruncateAndInsert(t.Context(), []any{database.Projects{}}))

	want := &openapi.Project{
		ID:        fixedID,
		Name:      "プロジェクト",
		Color:     "blue",
		CreatedAt: fixedNow,
		UpdatedAt: fixedNow,
	}
	httpcheck.New(th).Test(t, "POST", "/projects").
		WithHeader("Authorization", "Bearer "+token).
		WithHeader("Content-Type", "application/json").
		WithBody([]byte(`{"name": "プロジェクト", "color": "blue"}`)).
		Check().HasStatus(200).HasJSON(want)
}

func TestHandler_ListProjects(t *testing.T) {
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

	t.Run("no limit and offset", func(t *testing.T) {
		want := &openapi.ListProjectsOK{
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
		}
		httpcheck.New(th).Test(t, "GET", "/projects").
			WithHeader("Authorization", "Bearer "+token).
			Check().HasStatus(200).HasJSON(want)
	})
	t.Run("limit=1&offset=0", func(t *testing.T) {
		want := &openapi.ListProjectsOK{
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
		}
		httpcheck.New(th).Test(t, "GET", "/projects?limit=1&offset=0").
			WithHeader("Authorization", "Bearer "+token).
			Check().HasStatus(200).HasJSON(want)
	})
	t.Run("limit=1&offset=1", func(t *testing.T) {
		want := &openapi.ListProjectsOK{
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
		}
		httpcheck.New(th).Test(t, "GET", "/projects?limit=1&offset=1").
			WithHeader("Authorization", "Bearer "+token).
			Check().HasStatus(200).HasJSON(want)
	})
}

func TestHandler_GetProject(t *testing.T) {
	require.NoError(t, tdb.TruncateAndInsert(t.Context(), []any{
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
	}))

	t.Run("ok", func(t *testing.T) {
		want := &openapi.Project{
			ID:        "PROJECT-000000000000000001",
			Name:      "プロジェクト1",
			Color:     "blue",
			CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
			UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
		}
		httpcheck.New(th).Test(t, "GET", "/projects/PROJECT-000000000000000001").
			WithHeader("Authorization", "Bearer "+token).
			Check().HasStatus(200).HasJSON(want)
	})
	t.Run("project not found", func(t *testing.T) {
		want := handler.ErrorResponse{Code: 404, Message: "指定したプロジェクトは見つかりません"}
		httpcheck.New(th).Test(t, "GET", "/projects/PROJECT-000000000000000099").
			WithHeader("Authorization", "Bearer "+token).
			Check().HasStatus(404).HasJSON(want)
	})
	t.Run("project access denied", func(t *testing.T) {
		want := handler.ErrorResponse{Code: 404, Message: "指定したプロジェクトは見つかりません"}
		httpcheck.New(th).Test(t, "GET", "/projects/PROJECT-000000000000000002").
			WithHeader("Authorization", "Bearer "+token).
			Check().HasStatus(404).HasJSON(want)
	})
}

func TestHandler_UpdateProject(t *testing.T) {
	require.NoError(t, tdb.TruncateAndInsert(t.Context(), []any{
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
	}))

	t.Run("project not found", func(t *testing.T) {
		want := handler.ErrorResponse{Code: 404, Message: "指定したプロジェクトは見つかりません"}
		httpcheck.New(th).Test(t, "PATCH", "/projects/PROJECT-000000000000000099").
			WithHeader("Authorization", "Bearer "+token).
			WithHeader("Content-Type", "application/json").
			WithBody([]byte(`{"name": "更新後プロジェクト", "color": "gray", "is_archived": true}`)).
			Check().HasStatus(404).HasJSON(want)
	})
	t.Run("project access denied", func(t *testing.T) {
		want := handler.ErrorResponse{Code: 404, Message: "指定したプロジェクトは見つかりません"}
		httpcheck.New(th).Test(t, "PATCH", "/projects/PROJECT-000000000000000002").
			WithHeader("Authorization", "Bearer "+token).
			WithHeader("Content-Type", "application/json").
			WithBody([]byte(`{"name": "更新後プロジェクト", "color": "gray", "is_archived": true}`)).
			Check().HasStatus(404).HasJSON(want)
	})
	t.Run("ok", func(t *testing.T) {
		want := &openapi.Project{
			ID:         "PROJECT-000000000000000001",
			Name:       "更新後プロジェクト",
			Color:      "gray",
			IsArchived: true,
			CreatedAt:  time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
			UpdatedAt:  fixedNow,
		}
		httpcheck.New(th).Test(t, "PATCH", "/projects/PROJECT-000000000000000001").
			WithHeader("Authorization", "Bearer "+token).
			WithHeader("Content-Type", "application/json").
			WithBody([]byte(`{"name": "更新後プロジェクト", "color": "gray", "is_archived": true}`)).
			Check().HasStatus(200).HasJSON(want)
	})
}

func TestHandler_DeleteProject(t *testing.T) {
	require.NoError(t, tdb.TruncateAndInsert(t.Context(), []any{
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
	}))

	t.Run("project not found", func(t *testing.T) {
		want := handler.ErrorResponse{Code: 404, Message: "指定したプロジェクトは見つかりません"}
		httpcheck.New(th).Test(t, "DELETE", "/projects/PROJECT-000000000000000099").
			WithHeader("Authorization", "Bearer "+token).
			Check().HasStatus(404).HasJSON(want)
	})
	t.Run("project access denied", func(t *testing.T) {
		want := handler.ErrorResponse{Code: 404, Message: "指定したプロジェクトは見つかりません"}
		httpcheck.New(th).Test(t, "DELETE", "/projects/PROJECT-000000000000000002").
			WithHeader("Authorization", "Bearer "+token).
			Check().HasStatus(404).HasJSON(want)
	})
	t.Run("ok", func(t *testing.T) {
		httpcheck.New(th).Test(t, "DELETE", "/projects/PROJECT-000000000000000001").
			WithHeader("Authorization", "Bearer "+token).
			Check().HasStatus(200).HasString("")
	})
}
