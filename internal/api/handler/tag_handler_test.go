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

func TestHandler_CreateTag(t *testing.T) {
	require.NoError(t, tdb.Reset(t.Context(), []any{database.Tag{}}))

	want := &openapi.Tag{
		ID:        fixedID,
		Name:      "タグ",
		CreatedAt: fixedNow,
		UpdatedAt: fixedNow,
	}
	httpcheck.New(th).Test(t, "POST", "/tags").
		WithHeader("Authorization", "Bearer "+token).
		WithHeader("Content-Type", "application/json").
		WithBody([]byte(`{"name": "タグ"}`)).
		Check().HasStatus(200).HasJSON(want)
}

func TestHandler_ListTags(t *testing.T) {
	require.NoError(t, tdb.Reset(t.Context(), []any{database.Tag{}}))
	require.NoError(t, tdb.Insert(t.Context(), []any{
		database.Tags{
			{
				ID:        "TAG-0000000000000000000001",
				UserID:    testUserID,
				Name:      "タグ1",
				CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst),
				UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst),
			},
			{
				ID:        "TAG-0000000000000000000002",
				UserID:    testUserID,
				Name:      "タグ2",
				CreatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst),
				UpdatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst),
			},
		},
	}))

	t.Run("no limit and offset", func(t *testing.T) {
		want := &openapi.ListTagsOK{
			Tags: []openapi.Tag{
				{ID: "TAG-0000000000000000000001", Name: "タグ1", CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst)},
				{ID: "TAG-0000000000000000000002", Name: "タグ2", CreatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst)},
			},
			HasNext: false,
		}
		httpcheck.New(th).Test(t, "GET", "/tags").
			WithHeader("Authorization", "Bearer "+token).
			Check().HasStatus(200).HasJSON(want)
	})
	t.Run("limit=1&offset=0", func(t *testing.T) {
		want := &openapi.ListTagsOK{
			Tags:    []openapi.Tag{{ID: "TAG-0000000000000000000001", Name: "タグ1", CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst)}},
			HasNext: true,
		}
		httpcheck.New(th).Test(t, "GET", "/tags?limit=1&offset=0").
			WithHeader("Authorization", "Bearer "+token).
			Check().HasStatus(200).HasJSON(want)
	})
	t.Run("limit=1&offset=1", func(t *testing.T) {
		want := &openapi.ListTagsOK{
			Tags:    []openapi.Tag{{ID: "TAG-0000000000000000000002", Name: "タグ2", CreatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst)}},
			HasNext: false,
		}
		httpcheck.New(th).Test(t, "GET", "/tags?limit=1&offset=1").
			WithHeader("Authorization", "Bearer "+token).
			Check().HasStatus(200).HasJSON(want)
	})
}

func TestHandler_GetTag(t *testing.T) {
	require.NoError(t, tdb.Reset(t.Context(), []any{database.Tag{}}))
	require.NoError(t, tdb.Insert(t.Context(), []any{
		database.Tags{
			{
				ID:        "TAG-0000000000000000000001",
				UserID:    testUserID,
				Name:      "タグ1",
				CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
				UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
			},
			{
				ID:        "TAG-0000000000000000000002",
				UserID:    "USER-000000000000000000002",
				Name:      "タグ2",
				CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
				UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
			},
		},
	}))

	t.Run("ok", func(t *testing.T) {
		want := &openapi.Tag{
			ID:        "TAG-0000000000000000000001",
			Name:      "タグ1",
			CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
			UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
		}
		httpcheck.New(th).Test(t, "GET", "/tags/TAG-0000000000000000000001").
			WithHeader("Authorization", "Bearer "+token).
			Check().HasStatus(200).HasJSON(want)
	})
	t.Run("tag not found", func(t *testing.T) {
		want := handler.ErrorResponse{Code: 404, Message: "指定したタグは見つかりません"}
		httpcheck.New(th).Test(t, "GET", "/tags/TAG-0000000000000000000099").
			WithHeader("Authorization", "Bearer "+token).
			Check().HasStatus(404).HasJSON(want)
	})
	t.Run("tag access denied", func(t *testing.T) {
		want := handler.ErrorResponse{Code: 404, Message: "指定したタグは見つかりません"}
		httpcheck.New(th).Test(t, "GET", "/tags/TAG-0000000000000000000002").
			WithHeader("Authorization", "Bearer "+token).
			Check().HasStatus(404).HasJSON(want)
	})
}

func TestHandler_UpdateTag(t *testing.T) {
	require.NoError(t, tdb.Reset(t.Context(), []any{database.Tag{}}))
	require.NoError(t, tdb.Insert(t.Context(), []any{
		database.Tags{
			{
				ID:        "TAG-0000000000000000000001",
				UserID:    testUserID,
				Name:      "タグ1",
				CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
				UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
			},
			{
				ID:        "TAG-0000000000000000000002",
				UserID:    "USER-000000000000000000002",
				Name:      "タグ2",
				CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
				UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
			},
		},
	}))

	t.Run("tag not found", func(t *testing.T) {
		want := handler.ErrorResponse{Code: 404, Message: "指定したタグは見つかりません"}
		httpcheck.New(th).Test(t, "PATCH", "/tags/TAG-0000000000000000000099").
			WithHeader("Authorization", "Bearer "+token).
			WithHeader("Content-Type", "application/json").
			WithBody([]byte(`{"name": "更新後タグ"}`)).
			Check().HasStatus(404).HasJSON(want)
	})
	t.Run("tag access denied", func(t *testing.T) {
		want := handler.ErrorResponse{Code: 404, Message: "指定したタグは見つかりません"}
		httpcheck.New(th).Test(t, "PATCH", "/tags/TAG-0000000000000000000002").
			WithHeader("Authorization", "Bearer "+token).
			WithHeader("Content-Type", "application/json").
			WithBody([]byte(`{"name": "更新後タグ"}`)).
			Check().HasStatus(404).HasJSON(want)
	})
	t.Run("ok", func(t *testing.T) {
		want := &openapi.Tag{
			ID:        "TAG-0000000000000000000001",
			Name:      "更新後タグ",
			CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
			UpdatedAt: fixedNow,
		}
		httpcheck.New(th).Test(t, "PATCH", "/tags/TAG-0000000000000000000001").
			WithHeader("Authorization", "Bearer "+token).
			WithHeader("Content-Type", "application/json").
			WithBody([]byte(`{"name": "更新後タグ"}`)).
			Check().HasStatus(200).HasJSON(want)
	})
}

func TestHandler_DeleteTag(t *testing.T) {
	require.NoError(t, tdb.Reset(t.Context(), []any{database.Tag{}}))
	require.NoError(t, tdb.Insert(t.Context(), []any{
		database.Tags{
			{
				ID:        "TAG-0000000000000000000001",
				UserID:    testUserID,
				Name:      "タグ1",
				CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
				UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
			},
			{
				ID:        "TAG-0000000000000000000002",
				UserID:    "USER-000000000000000000002",
				Name:      "タグ2",
				CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
				UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
			},
		},
	}))

	t.Run("tag not found", func(t *testing.T) {
		want := handler.ErrorResponse{Code: 404, Message: "指定したタグは見つかりません"}
		httpcheck.New(th).Test(t, "DELETE", "/tags/TAG-0000000000000000000099").
			WithHeader("Authorization", "Bearer "+token).
			Check().HasStatus(404).HasJSON(want)
	})
	t.Run("tag access denied", func(t *testing.T) {
		want := handler.ErrorResponse{Code: 404, Message: "指定したタグは見つかりません"}
		httpcheck.New(th).Test(t, "DELETE", "/tags/TAG-0000000000000000000002").
			WithHeader("Authorization", "Bearer "+token).
			Check().HasStatus(404).HasJSON(want)
	})
	t.Run("ok", func(t *testing.T) {
		httpcheck.New(th).Test(t, "DELETE", "/tags/TAG-0000000000000000000001").
			WithHeader("Authorization", "Bearer "+token).
			Check().HasStatus(200).HasString("")
	})
}
