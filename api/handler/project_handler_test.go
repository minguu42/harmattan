package handler_test

import (
	"testing"
	"time"

	"github.com/ikawaha/httpcheck"
	"github.com/minguu42/harmattan/internal/database"
	"github.com/minguu42/harmattan/internal/oapi"
	"github.com/stretchr/testify/require"
)

func TestHandler_CreateProject(t *testing.T) {
	wantResponse := &oapi.Project{
		ID:        "01JGFJJZ000000000000000000",
		Name:      "テストプロジェクト",
		Color:     "#1E3A8A",
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := tdb.Insert(t.Context(), []any{
		database.User{
			ID:             "user_01",
			Email:          "user1@example.com",
			HashedPassword: "password",
			CreatedAt:      time.Date(2025, 1, 1, 0, 0, 0, 0, time.Local),
			UpdatedAt:      time.Date(2025, 1, 1, 0, 0, 0, 0, time.Local),
		},
	}); err != nil {
		t.Fatal(err)
	}
	checker := httpcheck.New(fixTimeMiddleware(h, now))
	checker.Test(t, "POST", "/projects").
		WithHeader("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ1c2VyXzAxIiwiZXhwIjoxNzM1NjYwODAwLCJpYXQiOjE3MzU2NTcyMDB9.bT7pyLGRAxG784_cg1DoZ9GD3GbGbNFichSlETzYfPc").
		WithHeader("Content-Type", "application/json").
		WithBody([]byte(`{"name": "テストプロジェクト", "color": "#1E3A8A"}`)).
		Check().
		HasStatus(200).
		HasJSON(wantResponse)
}

func TestHandler_ListProjects(t *testing.T) {
	require.NoError(t, tdb.Reset(t.Context(), []any{database.User{}, database.Project{}}))

	users := []database.User{{
		ID:             "user_01",
		Email:          "user1@example.com",
		HashedPassword: "password",
		CreatedAt:      time.Date(2025, 1, 1, 0, 0, 0, 0, time.Local),
		UpdatedAt:      time.Date(2025, 1, 1, 0, 0, 0, 0, time.Local),
	}}
	projects := []database.Project{
		{ID: "project_01", UserID: "user_01", Name: "プロジェクト1", Color: "#1E3A8A", CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, time.Local), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, time.Local)},
		{ID: "project_02", UserID: "user_01", Name: "プロジェクト2", Color: "#FFFFFF", CreatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, time.Local), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, time.Local)},
	}
	err := tdb.Insert(t.Context(), []any{users, projects})
	require.NoError(t, err)

	t.Run("no limit and offset", func(t *testing.T) {
		checker := httpcheck.New(fixTimeMiddleware(h, now))
		want := &oapi.Projects{
			Projects: []oapi.Project{
				{ID: "project_01", Name: "プロジェクト1", Color: "#1E3A8A", CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, time.Local), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, time.Local)},
				{ID: "project_02", Name: "プロジェクト2", Color: "#FFFFFF", CreatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, time.Local), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, time.Local)},
			},
			HasNext: false,
		}
		checker.Test(t, "GET", "/projects").
			WithHeader("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ1c2VyXzAxIiwiZXhwIjoxNzM1NjYwODAwLCJpYXQiOjE3MzU2NTcyMDB9.bT7pyLGRAxG784_cg1DoZ9GD3GbGbNFichSlETzYfPc").
			Check().
			HasStatus(200).
			HasJSON(want)
	})
	t.Run("limit=1&offset=0", func(t *testing.T) {
		checker := httpcheck.New(fixTimeMiddleware(h, now))
		want := &oapi.Projects{
			Projects: []oapi.Project{
				{ID: "project_01", Name: "プロジェクト1", Color: "#1E3A8A", CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, time.Local), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, time.Local)},
			},
			HasNext: true,
		}
		checker.Test(t, "GET", "/projects?limit=1&offset=0").
			WithHeader("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ1c2VyXzAxIiwiZXhwIjoxNzM1NjYwODAwLCJpYXQiOjE3MzU2NTcyMDB9.bT7pyLGRAxG784_cg1DoZ9GD3GbGbNFichSlETzYfPc").
			Check().
			HasStatus(200).
			HasJSON(want)
	})
	t.Run("limit=1&offset=1", func(t *testing.T) {
		checker := httpcheck.New(fixTimeMiddleware(h, now))
		want := &oapi.Projects{
			Projects: []oapi.Project{
				{ID: "project_02", Name: "プロジェクト2", Color: "#FFFFFF", CreatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, time.Local), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, time.Local)},
			},
			HasNext: false,
		}
		checker.Test(t, "GET", "/projects?limit=1&offset=1").
			WithHeader("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ1c2VyXzAxIiwiZXhwIjoxNzM1NjYwODAwLCJpYXQiOjE3MzU2NTcyMDB9.bT7pyLGRAxG784_cg1DoZ9GD3GbGbNFichSlETzYfPc").
			Check().
			HasStatus(200).
			HasJSON(want)
	})
}
