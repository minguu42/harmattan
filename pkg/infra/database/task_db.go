package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/minguu42/opepe/gen/sqlc"
	"github.com/minguu42/opepe/pkg/domain/model"
	"github.com/minguu42/opepe/pkg/domain/repository"
)

func (db *DB) CreateTask(ctx context.Context, t *model.Task) error {
	var dueOn sql.NullTime
	if t.DueOn != nil {
		dueOn = sql.NullTime{Time: *t.DueOn, Valid: true}
	}
	var completedAt sql.NullTime
	if t.CompletedAt != nil {
		completedAt = sql.NullTime{Time: *t.CompletedAt, Valid: true}
	}
	if err := sqlc.New(db.sqlDB).CreateTask(ctx, sqlc.CreateTaskParams{
		ID:          t.ID,
		ProjectID:   t.ProjectID,
		Title:       t.Title,
		Content:     t.Content,
		Priority:    uint32(t.Priority),
		DueOn:       dueOn,
		CompletedAt: completedAt,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}); err != nil {
		return fmt.Errorf("q.CreateTask failed: %w", err)
	}
	return nil
}

func (db *DB) GetTasksByProjectID(ctx context.Context, projectID string, limit, offset int) ([]model.Task, error) {
	tasks, err := sqlc.New(db.sqlDB).GetTasksByProjectID(ctx, sqlc.GetTasksByProjectIDParams{
		ProjectID: projectID,
		Limit:     int32(limit),
		Offset:    int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("q.GetTasksByProjectID failed: %w", err)
	}

	ts := make([]model.Task, 0, len(tasks))
	for _, t := range tasks {
		var dueOn *time.Time
		if t.DueOn.Valid {
			dueOn = &t.DueOn.Time
		}
		var completedAt *time.Time
		if t.CompletedAt.Valid {
			completedAt = &t.CompletedAt.Time
		}
		ts = append(ts, model.Task{
			ID:          t.ID,
			ProjectID:   t.ProjectID,
			Title:       t.Title,
			Content:     t.Content,
			Priority:    int(t.Priority),
			DueOn:       dueOn,
			CompletedAt: completedAt,
			CreatedAt:   t.CreatedAt,
			UpdatedAt:   t.UpdatedAt,
		})
	}
	return ts, nil
}

func (db *DB) GetTaskByID(ctx context.Context, id string) (*model.Task, error) {
	t, err := sqlc.New(db.sqlDB).GetTaskByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrModelNotFound
		}
		return nil, fmt.Errorf("q.GetTaskByID failed: %w", err)
	}

	var dueOn *time.Time
	if t.DueOn.Valid {
		dueOn = &t.DueOn.Time
	}
	var completedAt *time.Time
	if t.CompletedAt.Valid {
		completedAt = &t.CompletedAt.Time
	}
	return &model.Task{
		ID:          t.ID,
		ProjectID:   t.ProjectID,
		Title:       t.Title,
		Content:     t.Content,
		Priority:    int(t.Priority),
		DueOn:       dueOn,
		CompletedAt: completedAt,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}, nil
}

func (db *DB) UpdateTask(ctx context.Context, t *model.Task) error {
	var dueOn sql.NullTime
	if t.DueOn != nil {
		dueOn = sql.NullTime{Time: *t.DueOn, Valid: true}
	}
	if err := sqlc.New(db.sqlDB).UpdateTask(ctx, sqlc.UpdateTaskParams{
		Title:    t.Title,
		Content:  t.Content,
		Priority: uint32(t.Priority),
		DueOn:    dueOn,
		ID:       t.ID,
	}); err != nil {
		return fmt.Errorf("q.UpdateTask failed: %w", err)
	}
	return nil
}

func (db *DB) DeleteTask(ctx context.Context, id string) error {
	if err := sqlc.New(db.sqlDB).DeleteTask(ctx, id); err != nil {
		return fmt.Errorf("q.DeleteTask failed: %w", err)
	}
	return nil
}
