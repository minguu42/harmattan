package database_test

import (
	"testing"
	"time"

	"github.com/minguu42/harmattan/internal/database"
	"github.com/minguu42/harmattan/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_CreateTag(t *testing.T) {
	require.NoError(t, tdb.TruncateAndInsert(t.Context(), []any{
		database.Users{
			{ID: "user01", Email: "user01@dummy.invalid", HashedPassword: "pass", CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst)},
		},
		database.Tags{},
	}))

	err := c.CreateTag(t.Context(), &domain.Tag{
		ID:        "tag01",
		UserID:    "user01",
		Name:      "タグ1",
		CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst),
		UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst),
	})
	require.NoError(t, err)

	tdb.Assert(t, []any{
		database.Tags{
			{ID: "tag01", UserID: "user01", Name: "タグ1", CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst)},
		},
	})
}

func TestClient_ListTags(t *testing.T) {
	require.NoError(t, tdb.TruncateAndInsert(t.Context(), []any{
		database.Users{
			{ID: "user01", Email: "user01@dummy.invalid", HashedPassword: "pass", CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst)},
		},
		database.Tags{
			{ID: "tag01", UserID: "user01", Name: "タグ1", CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst)},
			{ID: "tag02", UserID: "user01", Name: "タグ2", CreatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst)},
			{ID: "tag03", UserID: "user01", Name: "タグ3", CreatedAt: time.Date(2025, 1, 1, 0, 0, 3, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 3, 0, jst)},
		},
	}))

	tests := []struct {
		name   string
		userID domain.UserID
		limit  int
		offset int
		want   domain.Tags
	}{
		{
			name:   "multiple",
			userID: "user01",
			limit:  10,
			offset: 0,
			want: domain.Tags{
				{ID: "tag01", UserID: "user01", Name: "タグ1", CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst)},
				{ID: "tag02", UserID: "user01", Name: "タグ2", CreatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst)},
				{ID: "tag03", UserID: "user01", Name: "タグ3", CreatedAt: time.Date(2025, 1, 1, 0, 0, 3, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 3, 0, jst)},
			},
		},
		{
			name:   "no_match",
			userID: "user99",
			limit:  10,
			offset: 0,
			want:   domain.Tags{},
		},
		{
			name:   "pagination",
			userID: "user01",
			limit:  2,
			offset: 1,
			want: domain.Tags{
				{ID: "tag02", UserID: "user01", Name: "タグ2", CreatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst)},
				{ID: "tag03", UserID: "user01", Name: "タグ3", CreatedAt: time.Date(2025, 1, 1, 0, 0, 3, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 3, 0, jst)},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.ListTags(t.Context(), tt.userID, tt.limit, tt.offset)
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestClient_GetTagByID(t *testing.T) {
	require.NoError(t, tdb.TruncateAndInsert(t.Context(), []any{
		database.Users{
			{ID: "user01", Email: "user01@dummy.invalid", HashedPassword: "pass", CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst)},
		},
		database.Tags{
			{ID: "tag01", UserID: "user01", Name: "タグ1", CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst)},
		},
	}))

	tests := []struct {
		name    string
		id      domain.TagID
		want    *domain.Tag
		wantErr error
	}{
		{
			name: "found",
			id:   "tag01",
			want: &domain.Tag{
				ID:        "tag01",
				UserID:    "user01",
				Name:      "タグ1",
				CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst),
				UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst),
			},
		},
		{
			name:    "not_found",
			id:      "tag99",
			wantErr: database.ErrNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.GetTagByID(t.Context(), tt.id)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, got)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestClient_GetTagsByIDs(t *testing.T) {
	require.NoError(t, tdb.TruncateAndInsert(t.Context(), []any{
		database.Users{
			{ID: "user01", Email: "user01@dummy.invalid", HashedPassword: "pass", CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst)},
		},
		database.Tags{
			{ID: "tag01", UserID: "user01", Name: "タグ1", CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst)},
			{ID: "tag02", UserID: "user01", Name: "タグ2", CreatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst)},
			{ID: "tag03", UserID: "user01", Name: "タグ3", CreatedAt: time.Date(2025, 1, 1, 0, 0, 3, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 3, 0, jst)},
		},
	}))

	tests := []struct {
		name string
		ids  []domain.TagID
		want domain.Tags
	}{
		{
			name: "multiple",
			ids:  []domain.TagID{"tag01", "tag03"},
			want: domain.Tags{
				{ID: "tag01", UserID: "user01", Name: "タグ1", CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst)},
				{ID: "tag03", UserID: "user01", Name: "タグ3", CreatedAt: time.Date(2025, 1, 1, 0, 0, 3, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 3, 0, jst)},
			},
		},
		{
			name: "empty",
			ids:  []domain.TagID{},
			want: domain.Tags{},
		},
		{
			name: "partial_match",
			ids:  []domain.TagID{"tag01", "tag99", "tag02"},
			want: domain.Tags{
				{ID: "tag01", UserID: "user01", Name: "タグ1", CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst)},
				{ID: "tag02", UserID: "user01", Name: "タグ2", CreatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst)},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.GetTagsByIDs(t.Context(), tt.ids)
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestClient_UpdateTag(t *testing.T) {
	require.NoError(t, tdb.TruncateAndInsert(t.Context(), []any{
		database.Users{
			{ID: "user01", Email: "user01@dummy.invalid", HashedPassword: "pass", CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst)},
		},
		database.Tags{
			{ID: "tag01", UserID: "user01", Name: "タグ1", CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst)},
		},
	}))

	err := c.UpdateTag(t.Context(), &domain.Tag{
		ID:        "tag01",
		Name:      "更新後タグ",
		UpdatedAt: time.Date(2025, 2, 1, 0, 0, 0, 0, jst),
	})
	require.NoError(t, err)

	tdb.Assert(t, []any{
		database.Tags{
			{ID: "tag01", UserID: "user01", Name: "更新後タグ", CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst), UpdatedAt: time.Date(2025, 2, 1, 0, 0, 0, 0, jst)},
		},
	})
}

func TestClient_DeleteTagByID(t *testing.T) {
	require.NoError(t, tdb.TruncateAndInsert(t.Context(), []any{
		database.Users{
			{ID: "user01", Email: "user01@dummy.invalid", HashedPassword: "pass", CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst)},
		},
		database.Tags{
			{ID: "tag01", UserID: "user01", Name: "タグ1", CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst)},
			{ID: "tag02", UserID: "user01", Name: "タグ2", CreatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst)},
		},
	}))

	err := c.DeleteTagByID(t.Context(), "tag01")
	require.NoError(t, err)

	tdb.Assert(t, []any{
		database.Tags{
			{ID: "tag02", UserID: "user01", Name: "タグ2", CreatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst)},
		},
	})
}
