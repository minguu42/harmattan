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

func TestClient_ListTasks(t *testing.T) {
	ctx := context.Background()

	completedAt := time.Date(2025, 1, 10, 0, 0, 0, 0, jst)
	dueOn := time.Date(2025, 2, 1, 0, 0, 0, 0, jst)

	require.NoError(t, tdb.TruncateAndInsert(ctx, []any{
		database.Users{
			{ID: "user1", Email: "user1@dummy.invalid", HashedPassword: "pass", CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst)},
		},
		database.Projects{
			{ID: "project1", UserID: "user1", Name: "プロジェクト1", Color: "blue", IsArchived: false, CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst)},
			{ID: "project2", UserID: "user1", Name: "プロジェクト2", Color: "red", IsArchived: false, CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst)},
		},
		database.Tags{
			{ID: "tag1", UserID: "user1", Name: "タグ1", CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst)},
			{ID: "tag2", UserID: "user1", Name: "タグ2", CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst)},
			{ID: "tag3", UserID: "user1", Name: "タグ3", CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst)},
		},
		database.Tasks{
			{ID: "task1", UserID: "user1", ProjectID: "project1", Name: "タスク1", Content: "Content 1", Priority: 1, DueOn: nil, CompletedAt: nil, CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst)},
			{ID: "task2", UserID: "user1", ProjectID: "project1", Name: "タスク2", Content: "Content 2", Priority: 2, DueOn: &dueOn, CompletedAt: nil, CreatedAt: time.Date(2025, 1, 2, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 2, 0, 0, 0, 0, jst)},
			{ID: "task3", UserID: "user1", ProjectID: "project1", Name: "タスク3", Content: "", Priority: 0, DueOn: nil, CompletedAt: &completedAt, CreatedAt: time.Date(2025, 1, 3, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 3, 0, 0, 0, 0, jst)},
			{ID: "task4", UserID: "user1", ProjectID: "project2", Name: "タスク4", Content: "", Priority: 0, DueOn: nil, CompletedAt: nil, CreatedAt: time.Date(2025, 1, 4, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 4, 0, 0, 0, 0, jst)},
			{ID: "task5", UserID: "user1", ProjectID: "project1", Name: "タスク5", Content: "", Priority: 0, DueOn: nil, CompletedAt: nil, CreatedAt: time.Date(2025, 1, 5, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 5, 0, 0, 0, 0, jst)},
		},
		database.Steps{
			{ID: "step1", UserID: "user1", TaskID: "task1", Name: "ステップ1-1", CompletedAt: nil, CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst)},
			{ID: "step2", UserID: "user1", TaskID: "task1", Name: "ステップ1-2", CompletedAt: nil, CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst)},
			{ID: "step3", UserID: "user1", TaskID: "task2", Name: "ステップ2-1", CompletedAt: nil, CreatedAt: time.Date(2025, 1, 2, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 2, 0, 0, 0, 0, jst)},
		},
		database.TaskTags{
			{TaskID: "task2", TagID: "tag1", CreatedAt: time.Date(2025, 1, 2, 0, 0, 0, 0, jst)},
			{TaskID: "task2", TagID: "tag2", CreatedAt: time.Date(2025, 1, 2, 0, 0, 0, 0, jst)},
		},
	}))

	tests := []struct {
		name          string
		projectID     domain.ProjectID
		limit         int
		offset        int
		showCompleted bool
		want          domain.Tasks
	}{
		{
			name:          "preloads_tasks_with_steps",
			projectID:     "project1",
			limit:         10,
			offset:        0,
			showCompleted: false,
			want: domain.Tasks{
				{
					ID:          "task1",
					UserID:      "user1",
					ProjectID:   "project1",
					Name:        "タスク1",
					TagIDs:      []domain.TagID{},
					Content:     "Content 1",
					Priority:    1,
					DueOn:       nil,
					CompletedAt: nil,
					CreatedAt:   time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
					UpdatedAt:   time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
					Steps: domain.Steps{
						{ID: "step1", UserID: "user1", TaskID: "task1", Name: "ステップ1-1", CompletedAt: nil, CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst)},
						{ID: "step2", UserID: "user1", TaskID: "task1", Name: "ステップ1-2", CompletedAt: nil, CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst)},
					},
				},
				{
					ID:          "task2",
					UserID:      "user1",
					ProjectID:   "project1",
					Name:        "タスク2",
					TagIDs:      []domain.TagID{"tag1", "tag2"},
					Content:     "Content 2",
					Priority:    2,
					DueOn:       &dueOn,
					CompletedAt: nil,
					CreatedAt:   time.Date(2025, 1, 2, 0, 0, 0, 0, jst),
					UpdatedAt:   time.Date(2025, 1, 2, 0, 0, 0, 0, jst),
					Steps: domain.Steps{
						{ID: "step3", UserID: "user1", TaskID: "task2", Name: "ステップ2-1", CompletedAt: nil, CreatedAt: time.Date(2025, 1, 2, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 2, 0, 0, 0, 0, jst)},
					},
				},
				{
					ID:          "task5",
					UserID:      "user1",
					ProjectID:   "project1",
					Name:        "タスク5",
					TagIDs:      []domain.TagID{},
					Content:     "",
					Priority:    0,
					DueOn:       nil,
					CompletedAt: nil,
					CreatedAt:   time.Date(2025, 1, 5, 0, 0, 0, 0, jst),
					UpdatedAt:   time.Date(2025, 1, 5, 0, 0, 0, 0, jst),
					Steps:       domain.Steps{},
				},
			},
		},
		{
			name:          "includes_tag_ids_from_task_tags",
			projectID:     "project1",
			limit:         10,
			offset:        0,
			showCompleted: false,
			want: domain.Tasks{
				{
					ID:          "task1",
					UserID:      "user1",
					ProjectID:   "project1",
					Name:        "タスク1",
					TagIDs:      []domain.TagID{},
					Content:     "Content 1",
					Priority:    1,
					DueOn:       nil,
					CompletedAt: nil,
					CreatedAt:   time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
					UpdatedAt:   time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
					Steps: domain.Steps{
						{ID: "step1", UserID: "user1", TaskID: "task1", Name: "ステップ1-1", CompletedAt: nil, CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst)},
						{ID: "step2", UserID: "user1", TaskID: "task1", Name: "ステップ1-2", CompletedAt: nil, CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst)},
					},
				},
				{
					ID:          "task2",
					UserID:      "user1",
					ProjectID:   "project1",
					Name:        "タスク2",
					TagIDs:      []domain.TagID{"tag1", "tag2"},
					Content:     "Content 2",
					Priority:    2,
					DueOn:       &dueOn,
					CompletedAt: nil,
					CreatedAt:   time.Date(2025, 1, 2, 0, 0, 0, 0, jst),
					UpdatedAt:   time.Date(2025, 1, 2, 0, 0, 0, 0, jst),
					Steps: domain.Steps{
						{ID: "step3", UserID: "user1", TaskID: "task2", Name: "ステップ2-1", CompletedAt: nil, CreatedAt: time.Date(2025, 1, 2, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 2, 0, 0, 0, 0, jst)},
					},
				},
				{
					ID:          "task5",
					UserID:      "user1",
					ProjectID:   "project1",
					Name:        "タスク5",
					TagIDs:      []domain.TagID{},
					Content:     "",
					Priority:    0,
					DueOn:       nil,
					CompletedAt: nil,
					CreatedAt:   time.Date(2025, 1, 5, 0, 0, 0, 0, jst),
					UpdatedAt:   time.Date(2025, 1, 5, 0, 0, 0, 0, jst),
					Steps:       domain.Steps{},
				},
			},
		},
		{
			name:          "excludes_completed_when_show_completed_is_false",
			projectID:     "project1",
			limit:         10,
			offset:        0,
			showCompleted: false,
			want: domain.Tasks{
				{
					ID:          "task1",
					UserID:      "user1",
					ProjectID:   "project1",
					Name:        "タスク1",
					TagIDs:      []domain.TagID{},
					Content:     "Content 1",
					Priority:    1,
					DueOn:       nil,
					CompletedAt: nil,
					CreatedAt:   time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
					UpdatedAt:   time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
					Steps: domain.Steps{
						{ID: "step1", UserID: "user1", TaskID: "task1", Name: "ステップ1-1", CompletedAt: nil, CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst)},
						{ID: "step2", UserID: "user1", TaskID: "task1", Name: "ステップ1-2", CompletedAt: nil, CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst)},
					},
				},
				{
					ID:          "task2",
					UserID:      "user1",
					ProjectID:   "project1",
					Name:        "タスク2",
					TagIDs:      []domain.TagID{"tag1", "tag2"},
					Content:     "Content 2",
					Priority:    2,
					DueOn:       &dueOn,
					CompletedAt: nil,
					CreatedAt:   time.Date(2025, 1, 2, 0, 0, 0, 0, jst),
					UpdatedAt:   time.Date(2025, 1, 2, 0, 0, 0, 0, jst),
					Steps: domain.Steps{
						{ID: "step3", UserID: "user1", TaskID: "task2", Name: "ステップ2-1", CompletedAt: nil, CreatedAt: time.Date(2025, 1, 2, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 2, 0, 0, 0, 0, jst)},
					},
				},
				{
					ID:          "task5",
					UserID:      "user1",
					ProjectID:   "project1",
					Name:        "タスク5",
					TagIDs:      []domain.TagID{},
					Content:     "",
					Priority:    0,
					DueOn:       nil,
					CompletedAt: nil,
					CreatedAt:   time.Date(2025, 1, 5, 0, 0, 0, 0, jst),
					UpdatedAt:   time.Date(2025, 1, 5, 0, 0, 0, 0, jst),
					Steps:       domain.Steps{},
				},
			},
		},
		{
			name:          "includes_completed_when_show_completed_is_true",
			projectID:     "project1",
			limit:         10,
			offset:        0,
			showCompleted: true,
			want: domain.Tasks{
				{
					ID:          "task1",
					UserID:      "user1",
					ProjectID:   "project1",
					Name:        "タスク1",
					TagIDs:      []domain.TagID{},
					Content:     "Content 1",
					Priority:    1,
					DueOn:       nil,
					CompletedAt: nil,
					CreatedAt:   time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
					UpdatedAt:   time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
					Steps: domain.Steps{
						{ID: "step1", UserID: "user1", TaskID: "task1", Name: "ステップ1-1", CompletedAt: nil, CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst)},
						{ID: "step2", UserID: "user1", TaskID: "task1", Name: "ステップ1-2", CompletedAt: nil, CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst)},
					},
				},
				{
					ID:          "task2",
					UserID:      "user1",
					ProjectID:   "project1",
					Name:        "タスク2",
					TagIDs:      []domain.TagID{"tag1", "tag2"},
					Content:     "Content 2",
					Priority:    2,
					DueOn:       &dueOn,
					CompletedAt: nil,
					CreatedAt:   time.Date(2025, 1, 2, 0, 0, 0, 0, jst),
					UpdatedAt:   time.Date(2025, 1, 2, 0, 0, 0, 0, jst),
					Steps: domain.Steps{
						{ID: "step3", UserID: "user1", TaskID: "task2", Name: "ステップ2-1", CompletedAt: nil, CreatedAt: time.Date(2025, 1, 2, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 2, 0, 0, 0, 0, jst)},
					},
				},
				{
					ID:          "task3",
					UserID:      "user1",
					ProjectID:   "project1",
					Name:        "タスク3",
					TagIDs:      []domain.TagID{},
					Content:     "",
					Priority:    0,
					DueOn:       nil,
					CompletedAt: &completedAt,
					CreatedAt:   time.Date(2025, 1, 3, 0, 0, 0, 0, jst),
					UpdatedAt:   time.Date(2025, 1, 3, 0, 0, 0, 0, jst),
					Steps:       domain.Steps{},
				},
				{
					ID:          "task5",
					UserID:      "user1",
					ProjectID:   "project1",
					Name:        "タスク5",
					TagIDs:      []domain.TagID{},
					Content:     "",
					Priority:    0,
					DueOn:       nil,
					CompletedAt: nil,
					CreatedAt:   time.Date(2025, 1, 5, 0, 0, 0, 0, jst),
					UpdatedAt:   time.Date(2025, 1, 5, 0, 0, 0, 0, jst),
					Steps:       domain.Steps{},
				},
			},
		},
		{
			name:          "returns_task_with_empty_steps",
			projectID:     "project1",
			limit:         1,
			offset:        2,
			showCompleted: false,
			want: domain.Tasks{
				{
					ID:          "task5",
					UserID:      "user1",
					ProjectID:   "project1",
					Name:        "タスク5",
					TagIDs:      []domain.TagID{},
					Content:     "",
					Priority:    0,
					DueOn:       nil,
					CompletedAt: nil,
					CreatedAt:   time.Date(2025, 1, 5, 0, 0, 0, 0, jst),
					UpdatedAt:   time.Date(2025, 1, 5, 0, 0, 0, 0, jst),
					Steps:       domain.Steps{},
				},
			},
		},
		{
			name:          "returns_task_with_empty_tag_ids",
			projectID:     "project1",
			limit:         1,
			offset:        0,
			showCompleted: false,
			want: domain.Tasks{
				{
					ID:          "task1",
					UserID:      "user1",
					ProjectID:   "project1",
					Name:        "タスク1",
					TagIDs:      []domain.TagID{},
					Content:     "Content 1",
					Priority:    1,
					DueOn:       nil,
					CompletedAt: nil,
					CreatedAt:   time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
					UpdatedAt:   time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
					Steps: domain.Steps{
						{ID: "step1", UserID: "user1", TaskID: "task1", Name: "ステップ1-1", CompletedAt: nil, CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst)},
						{ID: "step2", UserID: "user1", TaskID: "task1", Name: "ステップ1-2", CompletedAt: nil, CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst)},
					},
				},
			},
		},
		{
			name:          "pagination_with_limit_and_offset",
			projectID:     "project1",
			limit:         2,
			offset:        1,
			showCompleted: false,
			want: domain.Tasks{
				{
					ID:          "task2",
					UserID:      "user1",
					ProjectID:   "project1",
					Name:        "タスク2",
					TagIDs:      []domain.TagID{"tag1", "tag2"},
					Content:     "Content 2",
					Priority:    2,
					DueOn:       &dueOn,
					CompletedAt: nil,
					CreatedAt:   time.Date(2025, 1, 2, 0, 0, 0, 0, jst),
					UpdatedAt:   time.Date(2025, 1, 2, 0, 0, 0, 0, jst),
					Steps: domain.Steps{
						{ID: "step3", UserID: "user1", TaskID: "task2", Name: "ステップ2-1", CompletedAt: nil, CreatedAt: time.Date(2025, 1, 2, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 2, 0, 0, 0, 0, jst)},
					},
				},
				{
					ID:          "task5",
					UserID:      "user1",
					ProjectID:   "project1",
					Name:        "タスク5",
					TagIDs:      []domain.TagID{},
					Content:     "",
					Priority:    0,
					DueOn:       nil,
					CompletedAt: nil,
					CreatedAt:   time.Date(2025, 1, 5, 0, 0, 0, 0, jst),
					UpdatedAt:   time.Date(2025, 1, 5, 0, 0, 0, 0, jst),
					Steps:       domain.Steps{},
				},
			},
		},
		{
			name:          "excludes_other_projects_tasks",
			projectID:     "project1",
			limit:         10,
			offset:        0,
			showCompleted: false,
			want: domain.Tasks{
				{
					ID:          "task1",
					UserID:      "user1",
					ProjectID:   "project1",
					Name:        "タスク1",
					TagIDs:      []domain.TagID{},
					Content:     "Content 1",
					Priority:    1,
					DueOn:       nil,
					CompletedAt: nil,
					CreatedAt:   time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
					UpdatedAt:   time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
					Steps: domain.Steps{
						{ID: "step1", UserID: "user1", TaskID: "task1", Name: "ステップ1-1", CompletedAt: nil, CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst)},
						{ID: "step2", UserID: "user1", TaskID: "task1", Name: "ステップ1-2", CompletedAt: nil, CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst)},
					},
				},
				{
					ID:          "task2",
					UserID:      "user1",
					ProjectID:   "project1",
					Name:        "タスク2",
					TagIDs:      []domain.TagID{"tag1", "tag2"},
					Content:     "Content 2",
					Priority:    2,
					DueOn:       &dueOn,
					CompletedAt: nil,
					CreatedAt:   time.Date(2025, 1, 2, 0, 0, 0, 0, jst),
					UpdatedAt:   time.Date(2025, 1, 2, 0, 0, 0, 0, jst),
					Steps: domain.Steps{
						{ID: "step3", UserID: "user1", TaskID: "task2", Name: "ステップ2-1", CompletedAt: nil, CreatedAt: time.Date(2025, 1, 2, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 2, 0, 0, 0, 0, jst)},
					},
				},
				{
					ID:          "task5",
					UserID:      "user1",
					ProjectID:   "project1",
					Name:        "タスク5",
					TagIDs:      []domain.TagID{},
					Content:     "",
					Priority:    0,
					DueOn:       nil,
					CompletedAt: nil,
					CreatedAt:   time.Date(2025, 1, 5, 0, 0, 0, 0, jst),
					UpdatedAt:   time.Date(2025, 1, 5, 0, 0, 0, 0, jst),
					Steps:       domain.Steps{},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.ListTasks(ctx, tt.projectID, tt.limit, tt.offset, tt.showCompleted)

			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestClient_GetTaskByID(t *testing.T) {
	ctx := context.Background()

	completedAt := time.Date(2025, 1, 10, 0, 0, 0, 0, jst)
	dueOn := time.Date(2025, 2, 1, 0, 0, 0, 0, jst)

	require.NoError(t, tdb.TruncateAndInsert(ctx, []any{
		database.Users{
			{ID: "user1", Email: "user1@dummy.invalid", HashedPassword: "pass", CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst)},
		},
		database.Projects{
			{ID: "project1", UserID: "user1", Name: "プロジェクト1", Color: "blue", IsArchived: false, CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst)},
			{ID: "project2", UserID: "user1", Name: "プロジェクト2", Color: "red", IsArchived: false, CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst)},
		},
		database.Tags{
			{ID: "tag1", UserID: "user1", Name: "タグ1", CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst)},
			{ID: "tag2", UserID: "user1", Name: "タグ2", CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst)},
			{ID: "tag3", UserID: "user1", Name: "タグ3", CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst)},
		},
		database.Tasks{
			{ID: "task1", UserID: "user1", ProjectID: "project1", Name: "タスク1", Content: "Content 1", Priority: 1, DueOn: nil, CompletedAt: nil, CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst)},
			{ID: "task2", UserID: "user1", ProjectID: "project1", Name: "タスク2", Content: "Content 2", Priority: 2, DueOn: &dueOn, CompletedAt: nil, CreatedAt: time.Date(2025, 1, 2, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 2, 0, 0, 0, 0, jst)},
			{ID: "task3", UserID: "user1", ProjectID: "project1", Name: "タスク3", Content: "", Priority: 0, DueOn: nil, CompletedAt: &completedAt, CreatedAt: time.Date(2025, 1, 3, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 3, 0, 0, 0, 0, jst)},
			{ID: "task4", UserID: "user1", ProjectID: "project2", Name: "タスク4", Content: "", Priority: 0, DueOn: nil, CompletedAt: nil, CreatedAt: time.Date(2025, 1, 4, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 4, 0, 0, 0, 0, jst)},
			{ID: "task5", UserID: "user1", ProjectID: "project1", Name: "タスク5", Content: "", Priority: 0, DueOn: nil, CompletedAt: nil, CreatedAt: time.Date(2025, 1, 5, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 5, 0, 0, 0, 0, jst)},
		},
		database.Steps{
			{ID: "step1", UserID: "user1", TaskID: "task1", Name: "ステップ1", CompletedAt: nil, CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst)},
			{ID: "step2", UserID: "user1", TaskID: "task2", Name: "ステップ2-1", CompletedAt: nil, CreatedAt: time.Date(2025, 1, 2, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 2, 0, 0, 0, 0, jst)},
			{ID: "step3", UserID: "user1", TaskID: "task2", Name: "ステップ2-2", CompletedAt: nil, CreatedAt: time.Date(2025, 1, 2, 0, 0, 1, 0, jst), UpdatedAt: time.Date(2025, 1, 2, 0, 0, 1, 0, jst)},
			{ID: "step4", UserID: "user1", TaskID: "task2", Name: "ステップ2-3", CompletedAt: nil, CreatedAt: time.Date(2025, 1, 2, 0, 0, 2, 0, jst), UpdatedAt: time.Date(2025, 1, 2, 0, 0, 2, 0, jst)},
		},
		database.TaskTags{
			{TaskID: "task1", TagID: "tag1", CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst)},
			{TaskID: "task3", TagID: "tag1", CreatedAt: time.Date(2025, 1, 3, 0, 0, 0, 0, jst)},
			{TaskID: "task3", TagID: "tag2", CreatedAt: time.Date(2025, 1, 3, 0, 0, 0, 0, jst)},
			{TaskID: "task3", TagID: "tag3", CreatedAt: time.Date(2025, 1, 3, 0, 0, 0, 0, jst)},
		},
	}))

	tests := []struct {
		name    string
		input   domain.TaskID
		want    *domain.Task
		wantErr error
	}{
		{
			name:  "returns_task_with_steps_and_tag_ids",
			input: "task1",
			want: &domain.Task{
				ID:          "task1",
				UserID:      "user1",
				ProjectID:   "project1",
				Name:        "タスク1",
				TagIDs:      []domain.TagID{"tag1"},
				Content:     "Content 1",
				Priority:    1,
				DueOn:       nil,
				CompletedAt: nil,
				CreatedAt:   time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
				UpdatedAt:   time.Date(2025, 1, 1, 0, 0, 0, 0, jst),
				Steps: domain.Steps{
					{ID: "step1", UserID: "user1", TaskID: "task1", Name: "ステップ1", CompletedAt: nil, CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst)},
				},
			},
		},
		{
			name:  "returns_task_with_multiple_steps",
			input: "task2",
			want: &domain.Task{
				ID:          "task2",
				UserID:      "user1",
				ProjectID:   "project1",
				Name:        "タスク2",
				TagIDs:      []domain.TagID{},
				Content:     "Content 2",
				Priority:    2,
				DueOn:       &dueOn,
				CompletedAt: nil,
				CreatedAt:   time.Date(2025, 1, 2, 0, 0, 0, 0, jst),
				UpdatedAt:   time.Date(2025, 1, 2, 0, 0, 0, 0, jst),
				Steps: domain.Steps{
					{ID: "step2", UserID: "user1", TaskID: "task2", Name: "ステップ2-1", CompletedAt: nil, CreatedAt: time.Date(2025, 1, 2, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 2, 0, 0, 0, 0, jst)},
					{ID: "step3", UserID: "user1", TaskID: "task2", Name: "ステップ2-2", CompletedAt: nil, CreatedAt: time.Date(2025, 1, 2, 0, 0, 1, 0, jst), UpdatedAt: time.Date(2025, 1, 2, 0, 0, 1, 0, jst)},
					{ID: "step4", UserID: "user1", TaskID: "task2", Name: "ステップ2-3", CompletedAt: nil, CreatedAt: time.Date(2025, 1, 2, 0, 0, 2, 0, jst), UpdatedAt: time.Date(2025, 1, 2, 0, 0, 2, 0, jst)},
				},
			},
		},
		{
			name:  "returns_task_with_multiple_tag_ids",
			input: "task3",
			want: &domain.Task{
				ID:          "task3",
				UserID:      "user1",
				ProjectID:   "project1",
				Name:        "タスク3",
				TagIDs:      []domain.TagID{"tag1", "tag2", "tag3"},
				Content:     "",
				Priority:    0,
				DueOn:       nil,
				CompletedAt: &completedAt,
				CreatedAt:   time.Date(2025, 1, 3, 0, 0, 0, 0, jst),
				UpdatedAt:   time.Date(2025, 1, 3, 0, 0, 0, 0, jst),
				Steps:       domain.Steps{},
			},
		},
		{
			name:  "returns_task_with_empty_steps_and_tag_ids",
			input: "task5",
			want: &domain.Task{
				ID:          "task5",
				UserID:      "user1",
				ProjectID:   "project1",
				Name:        "タスク5",
				TagIDs:      []domain.TagID{},
				Content:     "",
				Priority:    0,
				DueOn:       nil,
				CompletedAt: nil,
				CreatedAt:   time.Date(2025, 1, 5, 0, 0, 0, 0, jst),
				UpdatedAt:   time.Date(2025, 1, 5, 0, 0, 0, 0, jst),
				Steps:       domain.Steps{},
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
			got, err := c.GetTaskByID(ctx, tt.input)

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
