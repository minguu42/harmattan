package database

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/minguu42/harmattan/internal/domain"
	"gorm.io/gorm"
)

type Task struct {
	ID          domain.TaskID
	UserID      domain.UserID
	ProjectID   domain.ProjectID
	Name        string
	Content     string
	Priority    int
	DueOn       *time.Time
	CompletedAt *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time

	Steps Steps
}

func (t *Task) ToDomain(taskTags TaskTags) *domain.Task {
	return &domain.Task{
		ID:          t.ID,
		UserID:      t.UserID,
		ProjectID:   t.ProjectID,
		Name:        t.Name,
		TagIDs:      taskTags.TagIDs(),
		Content:     t.Content,
		Priority:    t.Priority,
		DueOn:       t.DueOn,
		CompletedAt: t.CompletedAt,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
		Steps:       t.Steps.ToDomain(),
	}
}

type Tasks []Task

func (ts Tasks) ToDomain(taskTags TaskTags) domain.Tasks {
	tasks := make(domain.Tasks, 0, len(ts))
	taskTagsByTaskID := taskTags.ByTaskID()
	for _, t := range ts {
		tasks = append(tasks, *t.ToDomain(taskTagsByTaskID[t.ID]))
	}
	return tasks
}

func (ts Tasks) IDs() []domain.TaskID {
	ids := make([]domain.TaskID, 0, len(ts))
	for _, t := range ts {
		ids = append(ids, t.ID)
	}
	return ids
}

func (c *Client) CreateTask(ctx context.Context, t *domain.Task) error {
	task := Task{
		ID:          t.ID,
		UserID:      t.UserID,
		ProjectID:   t.ProjectID,
		Name:        t.Name,
		Content:     t.Content,
		Priority:    t.Priority,
		DueOn:       t.DueOn,
		CompletedAt: t.CompletedAt,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
	return c.db(ctx).Create(&task).Error
}

func (c *Client) ListTasks(ctx context.Context, projectID domain.ProjectID, limit, offset int) (domain.Tasks, error) {
	var ts Tasks
	if err := c.db(ctx).Preload("Steps").Where("project_id = ?", projectID).Limit(limit).Offset(offset).Find(&ts).Error; err != nil {
		return nil, fmt.Errorf("failed to query tasks: %w", err)
	}

	var tts TaskTags
	if err := c.db(ctx).Where("task_id in ?", ts.IDs()).Find(&tts).Error; err != nil {
		return nil, fmt.Errorf("failed to query task tags: %w", err)
	}
	return ts.ToDomain(tts), nil
}

func (c *Client) GetTaskByID(ctx context.Context, id domain.TaskID) (*domain.Task, error) {
	var t Task
	if err := c.db(ctx).Preload("Steps").Where("id = ?", id).Take(&t).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrModelNotFound
		}
		return nil, fmt.Errorf("failed to query task: %w", err)
	}

	var tts TaskTags
	if err := c.db(ctx).Where("task_id = ?", t.ID).Find(&tts).Error; err != nil {
		return nil, fmt.Errorf("failed to query task tags: %w", err)
	}
	return t.ToDomain(tts), nil
}

func (c *Client) UpdateTask(ctx context.Context, t *domain.Task) error {
	err := c.db(ctx).Model(Task{}).Where("id = ?", t.ID).Updates(map[string]any{
		"name":         t.Name,
		"content":      t.Content,
		"priority":     t.Priority,
		"due_on":       t.DueOn,
		"completed_at": t.CompletedAt,
		"updated_at":   t.UpdatedAt,
	}).Error
	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	if err := c.db(ctx).Where("task_id = ?", t.ID).Delete(TaskTag{}).Error; err != nil {
		return fmt.Errorf("failed to delete task tag: %w", err)
	}
	if len(t.TagIDs) == 0 {
		return nil
	}
	taskTags := make(TaskTags, 0, len(t.TagIDs))
	for _, tagID := range t.TagIDs {
		taskTags = append(taskTags, TaskTag{TaskID: t.ID, TagID: tagID, CreatedAt: t.UpdatedAt})
	}
	if err := c.db(ctx).Create(&taskTags).Error; err != nil {
		return fmt.Errorf("failed to create task tags: %w", err)
	}
	return nil
}

func (c *Client) DeleteTaskByID(ctx context.Context, id domain.TaskID) error {
	return c.db(ctx).Where("id = ?", id).Delete(Task{}).Error
}
