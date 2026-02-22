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

func TestClient_ListProjects(t *testing.T) {
	ctx := context.Background()

	setup := []any{
		database.Users{
			{ID: "user1", Email: "user1@dummy.invalid", HashedPassword: "pass", CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst)},
			{ID: "user2", Email: "user2@dummy.invalid", HashedPassword: "pass", CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst)},
		},
		database.Projects{
			{ID: "project1", UserID: "user1", Name: "プロジェクト1", Color: "blue", IsArchived: false, CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst)},
			{ID: "project2", UserID: "user1", Name: "プロジェクト2", Color: "red", IsArchived: true, CreatedAt: time.Date(2025, 1, 2, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 2, 0, 0, 0, 0, jst)},
			{ID: "project3", UserID: "user1", Name: "プロジェクト3", Color: "green", IsArchived: false, CreatedAt: time.Date(2025, 1, 3, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 3, 0, 0, 0, 0, jst)},
			{ID: "project4", UserID: "user1", Name: "プロジェクト4", Color: "yellow", IsArchived: false, CreatedAt: time.Date(2025, 1, 4, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 4, 0, 0, 0, 0, jst)},
			{ID: "project5", UserID: "user1", Name: "プロジェクト5", Color: "purple", IsArchived: false, CreatedAt: time.Date(2025, 1, 5, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 5, 0, 0, 0, 0, jst)},
			{ID: "project6", UserID: "user2", Name: "プロジェクト6", Color: "orange", IsArchived: false, CreatedAt: time.Date(2025, 1, 6, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 6, 0, 0, 0, 0, jst)},
		},
	}
	require.NoError(t, tdb.TruncateAndInsert(ctx, setup))

	tests := []struct {
		name   string
		userID domain.UserID
		limit  int
		offset int
		want   domain.Projects
	}{
		{
			name:   "returns_multiple_projects",
			userID: "user1",
			limit:  10,
			offset: 0,
			want: domain.Projects{
				{ID: "project1", UserID: "user1", Name: "プロジェクト1", Color: "blue", IsArchived: false, CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst)},
				{ID: "project2", UserID: "user1", Name: "プロジェクト2", Color: "red", IsArchived: true, CreatedAt: time.Date(2025, 1, 2, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 2, 0, 0, 0, 0, jst)},
				{ID: "project3", UserID: "user1", Name: "プロジェクト3", Color: "green", IsArchived: false, CreatedAt: time.Date(2025, 1, 3, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 3, 0, 0, 0, 0, jst)},
				{ID: "project4", UserID: "user1", Name: "プロジェクト4", Color: "yellow", IsArchived: false, CreatedAt: time.Date(2025, 1, 4, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 4, 0, 0, 0, 0, jst)},
				{ID: "project5", UserID: "user1", Name: "プロジェクト5", Color: "purple", IsArchived: false, CreatedAt: time.Date(2025, 1, 5, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 5, 0, 0, 0, 0, jst)},
			},
		},
		{
			name:   "returns_empty_array_when_no_results",
			userID: "nonexistent",
			limit:  10,
			offset: 0,
			want:   domain.Projects{},
		},
		{
			name:   "pagination_with_limit_and_offset",
			userID: "user1",
			limit:  2,
			offset: 1,
			want: domain.Projects{
				{ID: "project2", UserID: "user1", Name: "プロジェクト2", Color: "red", IsArchived: true, CreatedAt: time.Date(2025, 1, 2, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 2, 0, 0, 0, 0, jst)},
				{ID: "project3", UserID: "user1", Name: "プロジェクト3", Color: "green", IsArchived: false, CreatedAt: time.Date(2025, 1, 3, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 3, 0, 0, 0, 0, jst)},
			},
		},
		{
			name:   "excludes_other_users_projects",
			userID: "user1",
			limit:  10,
			offset: 0,
			want: domain.Projects{
				{ID: "project1", UserID: "user1", Name: "プロジェクト1", Color: "blue", IsArchived: false, CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst)},
				{ID: "project2", UserID: "user1", Name: "プロジェクト2", Color: "red", IsArchived: true, CreatedAt: time.Date(2025, 1, 2, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 2, 0, 0, 0, 0, jst)},
				{ID: "project3", UserID: "user1", Name: "プロジェクト3", Color: "green", IsArchived: false, CreatedAt: time.Date(2025, 1, 3, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 3, 0, 0, 0, 0, jst)},
				{ID: "project4", UserID: "user1", Name: "プロジェクト4", Color: "yellow", IsArchived: false, CreatedAt: time.Date(2025, 1, 4, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 4, 0, 0, 0, 0, jst)},
				{ID: "project5", UserID: "user1", Name: "プロジェクト5", Color: "purple", IsArchived: false, CreatedAt: time.Date(2025, 1, 5, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 5, 0, 0, 0, 0, jst)},
			},
		},
		{
			name:   "includes_archived_projects",
			userID: "user1",
			limit:  10,
			offset: 0,
			want: domain.Projects{
				{ID: "project1", UserID: "user1", Name: "プロジェクト1", Color: "blue", IsArchived: false, CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst)},
				{ID: "project2", UserID: "user1", Name: "プロジェクト2", Color: "red", IsArchived: true, CreatedAt: time.Date(2025, 1, 2, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 2, 0, 0, 0, 0, jst)},
				{ID: "project3", UserID: "user1", Name: "プロジェクト3", Color: "green", IsArchived: false, CreatedAt: time.Date(2025, 1, 3, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 3, 0, 0, 0, 0, jst)},
				{ID: "project4", UserID: "user1", Name: "プロジェクト4", Color: "yellow", IsArchived: false, CreatedAt: time.Date(2025, 1, 4, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 4, 0, 0, 0, 0, jst)},
				{ID: "project5", UserID: "user1", Name: "プロジェクト5", Color: "purple", IsArchived: false, CreatedAt: time.Date(2025, 1, 5, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 5, 0, 0, 0, 0, jst)},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.ListProjects(ctx, tt.userID, tt.limit, tt.offset)

			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestClient_GetProjectByID(t *testing.T) {
	ctx := context.Background()

	setup := []any{
		database.Users{
			{ID: "user1", Email: "user1@dummy.invalid", HashedPassword: "pass", CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst)},
		},
		database.Projects{
			{ID: "project1", UserID: "user1", Name: "プロジェクト1", Color: "blue", IsArchived: false, CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst)},
		},
	}
	require.NoError(t, tdb.TruncateAndInsert(ctx, setup))

	tests := []struct {
		name    string
		input   domain.ProjectID
		want    *domain.Project
		wantErr error
	}{
		{
			name:  "returns_project_when_exists",
			input: "project1",
			want: &domain.Project{
				ID:         "project1",
				UserID:     "user1",
				Name:       "プロジェクト1",
				Color:      "blue",
				IsArchived: false,
				CreatedAt:  time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
				UpdatedAt:  time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
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
			got, err := c.GetProjectByID(ctx, tt.input)

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
