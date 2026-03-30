package database_test

import (
	"testing"
	"time"

	"github.com/minguu42/harmattan/internal/database"
	"github.com/minguu42/harmattan/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_CreateTask(t *testing.T) {
	t.Parallel()
	c, testDB := setupTest(t, []any{
		database.Users{
			{ID: "user01", Email: "user01@dummy.invalid", HashedPassword: "pass", CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst)},
		},
		database.Projects{
			{ID: "project01", UserID: "user01", Name: "プロジェクト1", Color: "blue", CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst)},
		},
		database.Tasks{},
	})

	err := c.CreateTask(t.Context(), &domain.Task{
		ID:        "task01",
		UserID:    "user01",
		ProjectID: "project01",
		Name:      "タスク1",
		Content:   "Content 1",
		Priority:  1,
		CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst),
		UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst),
	})
	require.NoError(t, err)

	testDB.Assert(t, []any{
		database.Tasks{
			{ID: "task01", UserID: "user01", ProjectID: "project01", Name: "タスク1", Content: "Content 1", Priority: 1, CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst)},
		},
	})
}

func TestClient_ListTasks(t *testing.T) {
	t.Parallel()
	completedAt := time.Date(2025, 1, 10, 0, 0, 0, 0, jst)
	dueOn := time.Date(2025, 2, 1, 0, 0, 0, 0, jst)

	c, _ := setupTest(t, []any{
		database.Users{
			{ID: "user01", Email: "user01@dummy.invalid", HashedPassword: "pass", CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst)},
		},
		database.Projects{
			{ID: "project01", UserID: "user01", Name: "プロジェクト1", Color: "blue", CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst)},
		},
		database.Tags{
			{ID: "tag01", UserID: "user01", Name: "タグ1", CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst)},
			{ID: "tag02", UserID: "user01", Name: "タグ2", CreatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst)},
		},
		database.Tasks{
			{ID: "task01", UserID: "user01", ProjectID: "project01", Name: "タスク1", Content: "Content 1", Priority: 1, CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst)},
			{ID: "task02", UserID: "user01", ProjectID: "project01", Name: "タスク2", Content: "Content 2", Priority: 2, DueOn: &dueOn, CreatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst)},
			{ID: "task03", UserID: "user01", ProjectID: "project01", Name: "タスク3", CompletedAt: &completedAt, CreatedAt: time.Date(2025, 1, 1, 0, 0, 3, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 3, 0, jst)},
			{ID: "task04", UserID: "user01", ProjectID: "project01", Name: "タスク4", CreatedAt: time.Date(2025, 1, 1, 0, 0, 4, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 4, 0, jst)},
		},
		database.Steps{
			{ID: "step01", UserID: "user01", TaskID: "task01", Name: "ステップ1-1", CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst)},
			{ID: "step02", UserID: "user01", TaskID: "task01", Name: "ステップ1-2", CreatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst)},
			{ID: "step03", UserID: "user01", TaskID: "task02", Name: "ステップ2-1", CreatedAt: time.Date(2025, 1, 1, 0, 0, 3, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 3, 0, jst)},
		},
		database.TaskTags{
			{TaskID: "task02", TagID: "tag01", CreatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst)},
			{TaskID: "task02", TagID: "tag02", CreatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst)},
		},
	})

	tests := []struct {
		name          string
		projectID     domain.ProjectID
		limit         int
		offset        int
		showCompleted bool
		want          domain.Tasks
	}{
		{
			name:      "multiple",
			projectID: "project01",
			limit:     10,
			offset:    0,
			want: domain.Tasks{
				{
					ID: "task01", UserID: "user01", ProjectID: "project01",
					Name: "タスク1", TagIDs: []domain.TagID{}, Content: "Content 1",
					Priority:  1,
					CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst),
					Steps: domain.Steps{
						{ID: "step01", UserID: "user01", TaskID: "task01", Name: "ステップ1-1", CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst)},
						{ID: "step02", UserID: "user01", TaskID: "task01", Name: "ステップ1-2", CreatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst)},
					},
				},
				{
					ID: "task02", UserID: "user01", ProjectID: "project01",
					Name: "タスク2", TagIDs: []domain.TagID{"tag01", "tag02"}, Content: "Content 2",
					Priority: 2, DueOn: &dueOn,
					CreatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst),
					Steps: domain.Steps{
						{ID: "step03", UserID: "user01", TaskID: "task02", Name: "ステップ2-1", CreatedAt: time.Date(2025, 1, 1, 0, 0, 3, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 3, 0, jst)},
					},
				},
				{
					ID: "task04", UserID: "user01", ProjectID: "project01",
					Name: "タスク4", TagIDs: []domain.TagID{},
					CreatedAt: time.Date(2025, 1, 1, 0, 0, 4, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 4, 0, jst),
					Steps: domain.Steps{},
				},
			},
		},
		{
			name:          "show_completed",
			projectID:     "project01",
			limit:         10,
			offset:        0,
			showCompleted: true,
			want: domain.Tasks{
				{
					ID:        "task01",
					UserID:    "user01",
					ProjectID: "project01",
					Name:      "タスク1",
					TagIDs:    []domain.TagID{},
					Content:   "Content 1",
					Priority:  1,
					CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst),
					UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst),
					Steps: domain.Steps{
						{ID: "step01", UserID: "user01", TaskID: "task01", Name: "ステップ1-1", CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst)},
						{ID: "step02", UserID: "user01", TaskID: "task01", Name: "ステップ1-2", CreatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst)},
					},
				},
				{
					ID:        "task02",
					UserID:    "user01",
					ProjectID: "project01",
					Name:      "タスク2",
					TagIDs:    []domain.TagID{"tag01", "tag02"}, Content: "Content 2",
					Priority:  2,
					DueOn:     &dueOn,
					CreatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst),
					UpdatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst),
					Steps: domain.Steps{
						{ID: "step03", UserID: "user01", TaskID: "task02", Name: "ステップ2-1", CreatedAt: time.Date(2025, 1, 1, 0, 0, 3, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 3, 0, jst)},
					},
				},
				{
					ID:          "task03",
					UserID:      "user01",
					ProjectID:   "project01",
					Name:        "タスク3",
					TagIDs:      []domain.TagID{},
					CompletedAt: &completedAt,
					CreatedAt:   time.Date(2025, 1, 1, 0, 0, 3, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 3, 0, jst),
					Steps: domain.Steps{},
				},
				{
					ID:        "task04",
					UserID:    "user01",
					ProjectID: "project01",
					Name:      "タスク4",
					TagIDs:    []domain.TagID{},
					CreatedAt: time.Date(2025, 1, 1, 0, 0, 4, 0, jst),
					UpdatedAt: time.Date(2025, 1, 1, 0, 0, 4, 0, jst),
					Steps:     domain.Steps{},
				},
			},
		},
		{
			name:      "pagination",
			projectID: "project01",
			limit:     2,
			offset:    1,
			want: domain.Tasks{
				{
					ID:        "task02",
					UserID:    "user01",
					ProjectID: "project01",
					Name:      "タスク2",
					TagIDs:    []domain.TagID{"tag01", "tag02"},
					Content:   "Content 2",
					Priority:  2,
					DueOn:     &dueOn,
					CreatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst),
					UpdatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst),
					Steps: domain.Steps{
						{ID: "step03", UserID: "user01", TaskID: "task02", Name: "ステップ2-1", CreatedAt: time.Date(2025, 1, 1, 0, 0, 3, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 3, 0, jst)},
					},
				},
				{
					ID:        "task04",
					UserID:    "user01",
					ProjectID: "project01",
					Name:      "タスク4",
					TagIDs:    []domain.TagID{},
					CreatedAt: time.Date(2025, 1, 1, 0, 0, 4, 0, jst),
					UpdatedAt: time.Date(2025, 1, 1, 0, 0, 4, 0, jst),
					Steps:     domain.Steps{},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := c.ListTasks(t.Context(), tt.projectID, tt.limit, tt.offset, tt.showCompleted)
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestClient_GetTaskByID(t *testing.T) {
	t.Parallel()
	c, _ := setupTest(t, []any{
		database.Users{
			{ID: "user01", Email: "user01@dummy.invalid", HashedPassword: "pass", CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst)},
		},
		database.Projects{
			{ID: "project01", UserID: "user01", Name: "プロジェクト1", Color: "blue", CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst)},
		},
		database.Tags{
			{ID: "tag01", UserID: "user01", Name: "タグ1", CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst)},
		},
		database.Tasks{
			{ID: "task01", UserID: "user01", ProjectID: "project01", Name: "タスク1", Content: "Content 1", Priority: 1, CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst)},
		},
		database.Steps{
			{ID: "step01", UserID: "user01", TaskID: "task01", Name: "ステップ1", CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst)},
		},
		database.TaskTags{
			{TaskID: "task01", TagID: "tag01", CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst)},
		},
	})

	tests := []struct {
		name    string
		id      domain.TaskID
		want    *domain.Task
		wantErr error
	}{
		{
			name: "found",
			id:   "task01",
			want: &domain.Task{
				ID:        "task01",
				UserID:    "user01",
				ProjectID: "project01",
				Name:      "タスク1",
				TagIDs:    []domain.TagID{"tag01"},
				Content:   "Content 1",
				Priority:  1,
				CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst),
				UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst),
				Steps: domain.Steps{
					{ID: "step01", UserID: "user01", TaskID: "task01", Name: "ステップ1", CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst)},
				},
			},
		},
		{
			name:    "not_found",
			id:      "task99",
			wantErr: database.ErrNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := c.GetTaskByID(t.Context(), tt.id)
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

func TestClient_UpdateTask(t *testing.T) {
	t.Parallel()
	c, testDB := setupTest(t, []any{
		database.Users{
			{ID: "user01", Email: "user01@dummy.invalid", HashedPassword: "pass", CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst)},
		},
		database.Projects{
			{ID: "project01", UserID: "user01", Name: "プロジェクト1", Color: "blue", CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst)},
		},
		database.Tags{
			{ID: "tag01", UserID: "user01", Name: "タグ1", CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst)},
			{ID: "tag02", UserID: "user01", Name: "タグ2", CreatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst)},
		},
		database.Tasks{
			{ID: "task01", UserID: "user01", ProjectID: "project01", Name: "タスク1", CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst)},
		},
		database.TaskTags{
			{TaskID: "task01", TagID: "tag01", CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst)},
		},
	})

	err := c.UpdateTask(t.Context(), &domain.Task{
		ID:        "task01",
		Name:      "更新後タスク",
		Content:   "Updated content",
		Priority:  2,
		TagIDs:    []domain.TagID{"tag02"},
		UpdatedAt: time.Date(2025, 2, 1, 0, 0, 0, 0, jst),
	})
	require.NoError(t, err)

	testDB.Assert(t, []any{
		database.Tasks{
			{ID: "task01", UserID: "user01", ProjectID: "project01", Name: "更新後タスク", Content: "Updated content", Priority: 2, CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst), UpdatedAt: time.Date(2025, 2, 1, 0, 0, 0, 0, jst)},
		},
		database.TaskTags{
			{TaskID: "task01", TagID: "tag02", CreatedAt: time.Date(2025, 2, 1, 0, 0, 0, 0, jst)},
		},
	})
}

func TestClient_DeleteTaskByID(t *testing.T) {
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
			{ID: "task02", UserID: "user01", ProjectID: "project01", Name: "タスク2", CreatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst)},
		},
	})

	err := c.DeleteTaskByID(t.Context(), "task01")
	require.NoError(t, err)

	testDB.Assert(t, []any{
		database.Tasks{
			{ID: "task02", UserID: "user01", ProjectID: "project01", Name: "タスク2", CreatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 2, 0, jst)},
		},
	})
}
