package database_test

import (
	"context"
	"testing"
	"time"

	"github.com/minguu42/harmattan/internal/database"
	"github.com/minguu42/harmattan/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_ListTags(t *testing.T) {
	ctx := context.Background()

	setup := []any{
		database.Users{
			{ID: "user1", Email: "user1@dummy.invalid", HashedPassword: "pass", CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst)},
			{ID: "user2", Email: "user2@dummy.invalid", HashedPassword: "pass", CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst)},
		},
		database.Tags{
			{ID: "tag1", UserID: "user1", Name: "タグ1", CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst)},
			{ID: "tag2", UserID: "user1", Name: "タグ2", CreatedAt: time.Date(2025, 1, 2, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 2, 0, 0, 0, 0, jst)},
			{ID: "tag3", UserID: "user1", Name: "タグ3", CreatedAt: time.Date(2025, 1, 3, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 3, 0, 0, 0, 0, jst)},
			{ID: "tag4", UserID: "user1", Name: "タグ4", CreatedAt: time.Date(2025, 1, 4, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 4, 0, 0, 0, 0, jst)},
			{ID: "tag5", UserID: "user1", Name: "タグ5", CreatedAt: time.Date(2025, 1, 5, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 5, 0, 0, 0, 0, jst)},
			{ID: "tag6", UserID: "user2", Name: "タグ6", CreatedAt: time.Date(2025, 1, 6, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 6, 0, 0, 0, 0, jst)},
		},
	}
	require.NoError(t, tdb.TruncateAndInsert(ctx, setup))

	tests := []struct {
		name   string
		userID domain.UserID
		limit  int
		offset int
		want   domain.Tags
	}{
		{
			name:   "returns_multiple_tags",
			userID: "user1",
			limit:  10,
			offset: 0,
			want: domain.Tags{
				{ID: "tag1", UserID: "user1", Name: "タグ1", CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst)},
				{ID: "tag2", UserID: "user1", Name: "タグ2", CreatedAt: time.Date(2025, 1, 2, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 2, 0, 0, 0, 0, jst)},
				{ID: "tag3", UserID: "user1", Name: "タグ3", CreatedAt: time.Date(2025, 1, 3, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 3, 0, 0, 0, 0, jst)},
				{ID: "tag4", UserID: "user1", Name: "タグ4", CreatedAt: time.Date(2025, 1, 4, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 4, 0, 0, 0, 0, jst)},
				{ID: "tag5", UserID: "user1", Name: "タグ5", CreatedAt: time.Date(2025, 1, 5, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 5, 0, 0, 0, 0, jst)},
			},
		},
		{
			name:   "returns_empty_array_when_no_results",
			userID: "nonexistent",
			limit:  10,
			offset: 0,
			want:   domain.Tags{},
		},
		{
			name:   "pagination_with_limit_and_offset",
			userID: "user1",
			limit:  2,
			offset: 1,
			want: domain.Tags{
				{ID: "tag2", UserID: "user1", Name: "タグ2", CreatedAt: time.Date(2025, 1, 2, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 2, 0, 0, 0, 0, jst)},
				{ID: "tag3", UserID: "user1", Name: "タグ3", CreatedAt: time.Date(2025, 1, 3, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 3, 0, 0, 0, 0, jst)},
			},
		},
		{
			name:   "excludes_other_users_tags",
			userID: "user1",
			limit:  10,
			offset: 0,
			want: domain.Tags{
				{ID: "tag1", UserID: "user1", Name: "タグ1", CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst)},
				{ID: "tag2", UserID: "user1", Name: "タグ2", CreatedAt: time.Date(2025, 1, 2, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 2, 0, 0, 0, 0, jst)},
				{ID: "tag3", UserID: "user1", Name: "タグ3", CreatedAt: time.Date(2025, 1, 3, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 3, 0, 0, 0, 0, jst)},
				{ID: "tag4", UserID: "user1", Name: "タグ4", CreatedAt: time.Date(2025, 1, 4, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 4, 0, 0, 0, 0, jst)},
				{ID: "tag5", UserID: "user1", Name: "タグ5", CreatedAt: time.Date(2025, 1, 5, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 5, 0, 0, 0, 0, jst)},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.ListTags(ctx, tt.userID, tt.limit, tt.offset)

			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestClient_GetTagByID(t *testing.T) {
	ctx := context.Background()

	setup := []any{
		database.Users{
			{ID: "user1", Email: "user1@dummy.invalid", HashedPassword: "pass", CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst)},
		},
		database.Tags{
			{ID: "tag1", UserID: "user1", Name: "タグ1", CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst)},
		},
	}
	require.NoError(t, tdb.TruncateAndInsert(ctx, setup))

	tests := []struct {
		name    string
		input   domain.TagID
		want    *domain.Tag
		wantErr error
	}{
		{
			name:  "returns_tag_when_exists",
			input: "tag1",
			want: &domain.Tag{
				ID:        "tag1",
				UserID:    "user1",
				Name:      "タグ1",
				CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
				UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
			},
		},
		{
			name:    "returns_error_when_not_found",
			input:   "nonexistent",
			wantErr: database.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.GetTagByID(ctx, tt.input)

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
	ctx := context.Background()

	setup := []any{
		database.Users{
			{ID: "user1", Email: "user1@dummy.invalid", HashedPassword: "pass", CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst)},
		},
		database.Tags{
			{ID: "tag1", UserID: "user1", Name: "タグ1", CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst)},
			{ID: "tag2", UserID: "user1", Name: "タグ2", CreatedAt: time.Date(2025, 1, 2, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 2, 0, 0, 0, 0, jst)},
			{ID: "tag3", UserID: "user1", Name: "タグ3", CreatedAt: time.Date(2025, 1, 3, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 3, 0, 0, 0, 0, jst)},
		},
	}
	require.NoError(t, tdb.TruncateAndInsert(ctx, setup))

	tests := []struct {
		name  string
		input []domain.TagID
		want  domain.Tags
	}{
		{
			name:  "returns_multiple_tags_by_ids",
			input: []domain.TagID{"tag1", "tag3"},
			want: domain.Tags{
				{ID: "tag1", UserID: "user1", Name: "タグ1", CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst)},
				{ID: "tag3", UserID: "user1", Name: "タグ3", CreatedAt: time.Date(2025, 1, 3, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 3, 0, 0, 0, 0, jst)},
			},
		},
		{
			name:  "returns_empty_array_when_given_empty_ids",
			input: []domain.TagID{},
			want:  domain.Tags{},
		},
		{
			name:  "filters_out_nonexistent_ids",
			input: []domain.TagID{"tag1", "nonexistent", "tag2"},
			want: domain.Tags{
				{ID: "tag1", UserID: "user1", Name: "タグ1", CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst)},
				{ID: "tag2", UserID: "user1", Name: "タグ2", CreatedAt: time.Date(2025, 1, 2, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 2, 0, 0, 0, 0, jst)},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.GetTagsByIDs(ctx, tt.input)

			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
