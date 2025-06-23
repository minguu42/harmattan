package handler_test

import (
	"testing"
	"time"

	"github.com/ikawaha/httpcheck"
	"github.com/minguu42/harmattan/internal/database"
	"github.com/minguu42/harmattan/internal/oapi"
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
