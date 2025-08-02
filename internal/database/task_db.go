package database

import (
	"context"
	"errors"
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
	Steps       Steps `gorm:"foreignKey:TaskID"`
	Tags        Tags  `gorm:"many2many:tasks_tags;"`
}

func (t *Task) ToDomain() *domain.Task {
	return &domain.Task{
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
		Steps:       t.Steps.ToDomain(),
		Tags:        t.Tags.ToDomain(),
	}
}

type Tasks []Task

func (ts Tasks) ToDomain() domain.Tasks {
	tasks := make(domain.Tasks, 0, len(ts))
	for _, t := range ts {
		tasks = append(tasks, *t.ToDomain())
	}
	return tasks
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
	if err := c.db(ctx).Preload("Steps").Preload("Tags").Where("project_id = ?", projectID).Limit(limit).Offset(offset).Find(&ts).Error; err != nil {
		return nil, err
	}
	return ts.ToDomain(), nil
}

func (c *Client) GetTaskByID(ctx context.Context, id domain.TaskID) (*domain.Task, error) {
	var t Task
	if err := c.db(ctx).Preload("Steps").Preload("Tags").Where("id = ?", id).Take(&t).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrModelNotFound
		}
		return nil, err
	}
	return t.ToDomain(), nil
}

func (c *Client) UpdateTask(ctx context.Context, t *domain.Task) error {
	return c.db(ctx).Model(Task{}).Where("id = ?", t.ID).Updates(map[string]any{
		"name":         t.Name,
		"content":      t.Content,
		"priority":     t.Priority,
		"due_on":       t.DueOn,
		"completed_at": t.CompletedAt,
		"updated_at":   t.UpdatedAt,
	}).Error
}

func (c *Client) DeleteTaskByID(ctx context.Context, id domain.TaskID) error {
	return c.db(ctx).Where("id = ?", id).Delete(Task{}).Error
}
