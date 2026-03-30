package database_test

import (
	"testing"
	"time"

	"github.com/minguu42/harmattan/internal/database"
	"github.com/minguu42/harmattan/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_CreateStep(t *testing.T) {
	t.Parallel()
	c, testDB := setupTest(t, []any{
		database.Users{
			{ID: "user01", Email: "user01@dummy.invalid", HashedPassword: "pass", CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst)},
		},
		database.Projects{
			{ID: "project01", UserID: "user01", Name: "プロジェクト1", Color: "blue", CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst)},
		},
		database.Tasks{
			{ID: "task01", UserID: "user01", ProjectID: "project01", Name: "タスク1", CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst)},
		},
		database.Steps{},
	})

	err := c.CreateStep(t.Context(), &domain.Step{
		ID:        "step01",
		UserID:    "user01",
		TaskID:    "task01",
		Name:      "ステップ1",
		CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst),
		UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst),
	})
	require.NoError(t, err)

	testDB.Assert(t, []any{
		database.Steps{
			{ID: "step01", UserID: "user01", TaskID: "task01", Name: "ステップ1", CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst)},
		},
	})
}

func TestClient_GetStepByID(t *testing.T) {
	t.Parallel()
	c, _ := setupTest(t, []any{
		database.Users{
			{ID: "user01", Email: "user01@dummy.invalid", HashedPassword: "pass", CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst)},
		},
		database.Projects{
			{ID: "project01", UserID: "user01", Name: "プロジェクト1", Color: "blue", CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst)},
		},
		database.Tasks{
			{ID: "task01", UserID: "user01", ProjectID: "project01", Name: "タスク1", CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst)},
		},
		database.Steps{
			{ID: "step01", UserID: "user01", TaskID: "task01", Name: "ステップ1", CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst)},
		},
	})

	tests := []struct {
		name    string
		id      domain.StepID
		want    *domain.Step
		wantErr error
	}{
		{
			name: "found",
			id:   "step01",
			want: &domain.Step{
				ID:        "step01",
				UserID:    "user01",
				TaskID:    "task01",
				Name:      "ステップ1",
				CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst),
				UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst),
			},
		},
		{
			name:    "not_found",
			id:      "step99",
			wantErr: database.ErrNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := c.GetStepByID(t.Context(), tt.id)
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

func TestClient_UpdateStep(t *testing.T) {
	t.Parallel()
	c, testDB := setupTest(t, []any{
		database.Users{
			{ID: "user01", Email: "user01@dummy.invalid", HashedPassword: "pass", CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst)},
		},
		database.Projects{
			{ID: "project01", UserID: "user01", Name: "プロジェクト1", Color: "blue", CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst)},
		},
		database.Tasks{
			{ID: "task01", UserID: "user01", ProjectID: "project01", Name: "タスク1", CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst)},
		},
		database.Steps{
			{ID: "step01", UserID: "user01", TaskID: "task01", Name: "ステップ1", CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst)},
		},
	})

	completedAt := time.Date(2025, 1, 10, 0, 0, 0, 0, jst)
	err := c.UpdateStep(t.Context(), &domain.Step{
		ID:          "step01",
		Name:        "更新後ステップ",
		CompletedAt: &completedAt,
		UpdatedAt:   time.Date(2025, 2, 1, 0, 0, 0, 0, jst),
	})
	require.NoError(t, err)

	testDB.Assert(t, []any{
		database.Steps{
			{ID: "step01", UserID: "user01", TaskID: "task01", Name: "更新後ステップ", CompletedAt: &completedAt, CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst), UpdatedAt: time.Date(2025, 2, 1, 0, 0, 0, 0, jst)},
		},
	})
}

func TestClient_DeleteStepByID(t *testing.T) {
	t.Parallel()
	c, testDB := setupTest(t, []any{
		database.Users{
			{ID: "user01", Email: "user01@dummy.invalid", HashedPassword: "pass", CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst)},
		},
		database.Projects{
			{ID: "project01", UserID: "user01", Name: "プロジェクト1", Color: "blue", CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst)},
		},
		database.Tasks{
			{ID: "task01", UserID: "user01", ProjectID: "project01", Name: "タスク1", CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst)},
		},
		database.Steps{
			{ID: "step01", UserID: "user01", TaskID: "task01", Name: "ステップ1", CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst)},
			{ID: "step02", UserID: "user01", TaskID: "task01", Name: "ステップ2", CreatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst)},
		},
	})

	err := c.DeleteStepByID(t.Context(), "step01")
	require.NoError(t, err)

	testDB.Assert(t, []any{
		database.Steps{
			{ID: "step02", UserID: "user01", TaskID: "task01", Name: "ステップ2", CreatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst)},
		},
	})
}
