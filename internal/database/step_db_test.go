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

func TestClient_GetStepByID(t *testing.T) {
	ctx := context.Background()

	completedAt := time.Date(2025, 1, 10, 0, 0, 0, 0, jst)

	require.NoError(t, tdb.TruncateAndInsert(ctx, []any{
		database.Users{
			{ID: "user1", Email: "user1@dummy.invalid", HashedPassword: "pass", CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst)},
		},
		database.Projects{
			{ID: "project1", UserID: "user1", Name: "プロジェクト1", Color: "blue", IsArchived: false, CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst)},
		},
		database.Tasks{
			{ID: "task1", UserID: "user1", ProjectID: "project1", Name: "タスク1", Content: "", Priority: 0, DueOn: nil, CompletedAt: nil, CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst)},
		},
		database.Steps{
			{ID: "step1", UserID: "user1", TaskID: "task1", Name: "ステップ1", CompletedAt: nil, CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst)},
			{ID: "step2", UserID: "user1", TaskID: "task1", Name: "ステップ2", CompletedAt: nil, CreatedAt: time.Date(2025, 1, 2, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 2, 0, 0, 0, 0, jst)},
			{ID: "step3", UserID: "user1", TaskID: "task1", Name: "ステップ3", CompletedAt: &completedAt, CreatedAt: time.Date(2025, 1, 3, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 3, 0, 0, 0, 0, jst)},
		},
	}))

	tests := []struct {
		name    string
		input   domain.StepID
		want    *domain.Step
		wantErr error
	}{
		{
			name:  "returns_step_when_exists",
			input: "step1",
			want: &domain.Step{
				ID:          "step1",
				UserID:      "user1",
				TaskID:      "task1",
				Name:        "ステップ1",
				CompletedAt: nil,
				CreatedAt:   time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
				UpdatedAt:   time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
			},
		},
		{
			name:  "returns_step_with_null_completed_at",
			input: "step2",
			want: &domain.Step{
				ID:          "step2",
				UserID:      "user1",
				TaskID:      "task1",
				Name:        "ステップ2",
				CompletedAt: nil,
				CreatedAt:   time.Date(2025, 1, 2, 0, 0, 0, 0, jst),
				UpdatedAt:   time.Date(2025, 1, 2, 0, 0, 0, 0, jst),
			},
		},
		{
			name:  "returns_step_with_non_null_completed_at",
			input: "step3",
			want: &domain.Step{
				ID:          "step3",
				UserID:      "user1",
				TaskID:      "task1",
				Name:        "ステップ3",
				CompletedAt: &completedAt,
				CreatedAt:   time.Date(2025, 1, 3, 0, 0, 0, 0, jst),
				UpdatedAt:   time.Date(2025, 1, 3, 0, 0, 0, 0, jst),
			},
		},
		{
			name:    "returns_error_when_not_found",
			input:   "nonexistent",
			wantErr: database.ErrModelNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.GetStepByID(ctx, tt.input)

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
