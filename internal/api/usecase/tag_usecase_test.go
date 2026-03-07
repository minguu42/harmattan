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

func TestTag_CreateTag(t *testing.T) {
	require.NoError(t, tdb.TruncateAndInsert(t.Context(), []any{database.Tags{}}))

	runTest(t, test{
		Method:     "POST",
		Path:       "/tags",
		Headers:    http.Header{"Authorization": []string{"Bearer " + token}},
		Body:       `{"name": "タグ"}`,
		WantStatus: 200,
		WantJSON: &openapi.Tag{
			ID:        fixedID,
			Name:      "タグ",
			CreatedAt: fixedNow,
			UpdatedAt: fixedNow,
		},
		WantTables: []any{database.Tags{
			{ID: fixedID, UserID: testUserID, Name: "タグ", CreatedAt: fixedNow, UpdatedAt: fixedNow},
		}},
	})
}

func TestTag_ListTags(t *testing.T) {
	require.NoError(t, tdb.TruncateAndInsert(t.Context(), []any{
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

	tests := map[string]test{
		"no_limit_and_offset": {
			Method:     "GET",
			Path:       "/tags",
			Headers:    http.Header{"Authorization": []string{"Bearer " + token}},
			WantStatus: 200,
			WantJSON: &openapi.ListTagsOK{
				Tags: []openapi.Tag{
					{ID: "TAG-0000000000000000000001", Name: "タグ1", CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst)},
					{ID: "TAG-0000000000000000000002", Name: "タグ2", CreatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst)},
				},
				HasNext: false,
			},
		},
		"limit_1_offset_0": {
			Method:     "GET",
			Path:       "/tags?limit=1&offset=0",
			Headers:    http.Header{"Authorization": []string{"Bearer " + token}},
			WantStatus: 200,
			WantJSON: &openapi.ListTagsOK{
				Tags:    []openapi.Tag{{ID: "TAG-0000000000000000000001", Name: "タグ1", CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst)}},
				HasNext: true,
			},
		},
		"limit_1_offset_1": {
			Method:     "GET",
			Path:       "/tags?limit=1&offset=1",
			Headers:    http.Header{"Authorization": []string{"Bearer " + token}},
			WantStatus: 200,
			WantJSON: &openapi.ListTagsOK{
				Tags:    []openapi.Tag{{ID: "TAG-0000000000000000000002", Name: "タグ2", CreatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst)}},
				HasNext: false,
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			runTest(t, tt)
		})
	}
}

func TestTag_GetTag(t *testing.T) {
	require.NoError(t, tdb.TruncateAndInsert(t.Context(), []any{
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
				UserID:    "USER-000000000000000000002",
				Name:      "タグ2",
				CreatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst),
				UpdatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst),
			},
		},
	}))

	tests := map[string]test{
		"ok": {
			Method:     "GET",
			Path:       "/tags/TAG-0000000000000000000001",
			Headers:    http.Header{"Authorization": []string{"Bearer " + token}},
			WantStatus: 200,
			WantJSON: &openapi.Tag{
				ID:        "TAG-0000000000000000000001",
				Name:      "タグ1",
				CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst),
				UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst),
			},
		},
		"not_found": {
			Method:     "GET",
			Path:       "/tags/TAG-0000000000000000000099",
			Headers:    http.Header{"Authorization": []string{"Bearer " + token}},
			WantStatus: 404,
			WantJSON:   handler.ErrorResponse{Code: 404, Message: "指定したタグは見つかりません"},
		},
		"access_denied": {
			Method:     "GET",
			Path:       "/tags/TAG-0000000000000000000002",
			Headers:    http.Header{"Authorization": []string{"Bearer " + token}},
			WantStatus: 404,
			WantJSON:   handler.ErrorResponse{Code: 404, Message: "指定したタグは見つかりません"},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			runTest(t, tt)
		})
	}
}

func TestTag_UpdateTag(t *testing.T) {
	require.NoError(t, tdb.TruncateAndInsert(t.Context(), []any{
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
				UserID:    "USER-000000000000000000002",
				Name:      "タグ2",
				CreatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst),
				UpdatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst),
			},
		},
	}))

	tests := map[string]test{
		"ok": {
			Method:     "PATCH",
			Path:       "/tags/TAG-0000000000000000000001",
			Headers:    http.Header{"Authorization": []string{"Bearer " + token}},
			Body:       `{"name": "更新後タグ"}`,
			WantStatus: 200,
			WantJSON: &openapi.Tag{
				ID:        "TAG-0000000000000000000001",
				Name:      "更新後タグ",
				CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst),
				UpdatedAt: fixedNow,
			},
			WantTables: []any{database.Tags{
				{ID: "TAG-0000000000000000000001", UserID: testUserID, Name: "更新後タグ", CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst), UpdatedAt: fixedNow},
				{ID: "TAG-0000000000000000000002", UserID: "USER-000000000000000000002", Name: "タグ2", CreatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst)},
			}},
		},
		"not_found": {
			Method:     "PATCH",
			Path:       "/tags/TAG-0000000000000000000099",
			Headers:    http.Header{"Authorization": []string{"Bearer " + token}},
			Body:       `{"name": "更新後タグ"}`,
			WantStatus: 404,
			WantJSON:   handler.ErrorResponse{Code: 404, Message: "指定したタグは見つかりません"},
		},
		"access_denied": {
			Method:     "PATCH",
			Path:       "/tags/TAG-0000000000000000000002",
			Headers:    http.Header{"Authorization": []string{"Bearer " + token}},
			Body:       `{"name": "更新後タグ"}`,
			WantStatus: 404,
			WantJSON:   handler.ErrorResponse{Code: 404, Message: "指定したタグは見つかりません"},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			runTest(t, tt)
		})
	}
}

func TestTag_DeleteTag(t *testing.T) {
	require.NoError(t, tdb.TruncateAndInsert(t.Context(), []any{
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
				UserID:    "USER-000000000000000000002",
				Name:      "タグ2",
				CreatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst),
				UpdatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst),
			},
		},
	}))

	tests := map[string]test{
		"ok": {
			Method:     "DELETE",
			Path:       "/tags/TAG-0000000000000000000001",
			Headers:    http.Header{"Authorization": []string{"Bearer " + token}},
			WantStatus: 200,
			WantTables: []any{database.Tags{
				{ID: "TAG-0000000000000000000002", UserID: "USER-000000000000000000002", Name: "タグ2", CreatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst)},
			}},
		},
		"not_found": {
			Method:     "DELETE",
			Path:       "/tags/TAG-0000000000000000000099",
			Headers:    http.Header{"Authorization": []string{"Bearer " + token}},
			WantStatus: 404,
			WantJSON:   handler.ErrorResponse{Code: 404, Message: "指定したタグは見つかりません"},
		},
		"access_denied": {
			Method:     "DELETE",
			Path:       "/tags/TAG-0000000000000000000002",
			Headers:    http.Header{"Authorization": []string{"Bearer " + token}},
			WantStatus: 404,
			WantJSON:   handler.ErrorResponse{Code: 404, Message: "指定したタグは見つかりません"},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			runTest(t, tt)
		})
	}
}
