package database_test

import (
	"testing"
	"time"

	"github.com/minguu42/harmattan/internal/database"
	"github.com/minguu42/harmattan/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_CreateProject(t *testing.T) {
	require.NoError(t, tdb.TruncateAndInsert(t.Context(), []any{
		database.Users{
			{ID: "user01", Email: "user01@dummy.invalid", HashedPassword: "pass", CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst)},
		},
		database.Projects{},
	}))

	err := c.CreateProject(t.Context(), &domain.Project{
		ID:         "project01",
		UserID:     "user01",
		Name:       "プロジェクト1",
		Color:      "blue",
		IsArchived: false,
		CreatedAt:  time.Date(2025, 1, 1, 0, 0, 1, 0, jst),
		UpdatedAt:  time.Date(2025, 1, 1, 0, 0, 1, 0, jst),
	})
	require.NoError(t, err)

	tdb.Assert(t, []any{
		database.Projects{
			{ID: "project01", UserID: "user01", Name: "プロジェクト1", Color: "blue", IsArchived: false, CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst)},
		},
	})
}

func TestClient_ListProjects(t *testing.T) {
	require.NoError(t, tdb.TruncateAndInsert(t.Context(), []any{
		database.Users{
			{ID: "user01", Email: "user01@dummy.invalid", HashedPassword: "pass", CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst)},
		},
		database.Projects{
			{ID: "project01", UserID: "user01", Name: "プロジェクト1", Color: "blue", IsArchived: false, CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst)},
			{ID: "project02", UserID: "user01", Name: "プロジェクト2", Color: "red", IsArchived: true, CreatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst)},
			{ID: "project03", UserID: "user01", Name: "プロジェクト3", Color: "green", IsArchived: false, CreatedAt: time.Date(2025, 1, 1, 0, 0, 3, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 3, 0, jst)},
		},
	}))

	tests := []struct {
		name   string
		userID domain.UserID
		limit  int
		offset int
		want   domain.Projects
	}{
		{
			name:   "multiple",
			userID: "user01",
			limit:  10,
			offset: 0,
			want: domain.Projects{
				{ID: "project01", UserID: "user01", Name: "プロジェクト1", Color: "blue", IsArchived: false, CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst)},
				{ID: "project02", UserID: "user01", Name: "プロジェクト2", Color: "red", IsArchived: true, CreatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst)},
				{ID: "project03", UserID: "user01", Name: "プロジェクト3", Color: "green", IsArchived: false, CreatedAt: time.Date(2025, 1, 1, 0, 0, 3, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 3, 0, jst)},
			},
		},
		{
			name:   "no_match",
			userID: "user99",
			limit:  10,
			offset: 0,
			want:   domain.Projects{},
		},
		{
			name:   "pagination",
			userID: "user01",
			limit:  2,
			offset: 1,
			want: domain.Projects{
				{ID: "project02", UserID: "user01", Name: "プロジェクト2", Color: "red", IsArchived: true, CreatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst)},
				{ID: "project03", UserID: "user01", Name: "プロジェクト3", Color: "green", IsArchived: false, CreatedAt: time.Date(2025, 1, 1, 0, 0, 3, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 3, 0, jst)},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := c.ListProjects(t.Context(), tt.userID, tt.limit, tt.offset)
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestClient_GetProjectByID(t *testing.T) {
	require.NoError(t, tdb.TruncateAndInsert(t.Context(), []any{
		database.Users{
			{ID: "user01", Email: "user01@dummy.invalid", HashedPassword: "pass", CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst)},
		},
		database.Projects{
			{ID: "project01", UserID: "user01", Name: "プロジェクト1", Color: "blue", IsArchived: false, CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst)},
		},
	}))

	tests := []struct {
		name    string
		id      domain.ProjectID
		want    *domain.Project
		wantErr error
	}{
		{
			name: "found",
			id:   "project01",
			want: &domain.Project{
				ID:         "project01",
				UserID:     "user01",
				Name:       "プロジェクト1",
				Color:      "blue",
				IsArchived: false,
				CreatedAt:  time.Date(2025, 1, 1, 0, 0, 1, 0, jst),
				UpdatedAt:  time.Date(2025, 1, 1, 0, 0, 1, 0, jst),
			},
		},
		{
			name:    "not_found",
			id:      "project99",
			wantErr: database.ErrNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := c.GetProjectByID(t.Context(), tt.id)
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

func TestClient_UpdateProject(t *testing.T) {
	require.NoError(t, tdb.TruncateAndInsert(t.Context(), []any{
		database.Users{
			{ID: "user01", Email: "user01@dummy.invalid", HashedPassword: "pass", CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst)},
		},
		database.Projects{
			{ID: "project01", UserID: "user01", Name: "プロジェクト1", Color: "blue", IsArchived: false, CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst)},
		},
	}))

	err := c.UpdateProject(t.Context(), &domain.Project{
		ID:         "project01",
		Name:       "更新後プロジェクト",
		Color:      "red",
		IsArchived: true,
		UpdatedAt:  time.Date(2025, 2, 1, 0, 0, 0, 0, jst),
	})
	require.NoError(t, err)

	tdb.Assert(t, []any{
		database.Projects{
			{ID: "project01", UserID: "user01", Name: "更新後プロジェクト", Color: "red", IsArchived: true, CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst), UpdatedAt: time.Date(2025, 2, 1, 0, 0, 0, 0, jst)},
		},
	})
}

func TestClient_DeleteProjectByID(t *testing.T) {
	require.NoError(t, tdb.TruncateAndInsert(t.Context(), []any{
		database.Users{
			{ID: "user01", Email: "user01@dummy.invalid", HashedPassword: "pass", CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst)},
		},
		database.Projects{
			{ID: "project01", UserID: "user01", Name: "プロジェクト1", Color: "blue", IsArchived: false, CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst)},
			{ID: "project02", UserID: "user01", Name: "プロジェクト2", Color: "red", IsArchived: false, CreatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst)},
		},
	}))

	err := c.DeleteProjectByID(t.Context(), "project01")
	require.NoError(t, err)

	tdb.Assert(t, []any{
		database.Projects{
			{ID: "project02", UserID: "user01", Name: "プロジェクト2", Color: "red", IsArchived: false, CreatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst)},
		},
	})
}
