package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/minguu42/opepe/gen/sqlc"
	"github.com/minguu42/opepe/pkg/entity"
	"github.com/minguu42/opepe/pkg/ttime"
)

func (db *DB) SaveTask(ctx context.Context, t *entity.Task) error {
	var dueOn sql.NullTime
	if t.DueOn != nil {
		dueOn = sql.NullTime{Time: *t.DueOn, Valid: true}
	}
	if t.ID != "" {
		if err := sqlc.New(db._db).UpdateTask(ctx, sqlc.UpdateTaskParams{
			Title:    t.Title,
			Content:  t.Content,
			Priority: uint32(t.Priority),
			DueOn:    dueOn,
			ID:       t.ID,
		}); err != nil {
			return fmt.Errorf("q.UpdateTask failed: %w", err)
		}
	}

	now := ttime.Now(ctx)
	if err := sqlc.New(db._db).CreateTask(ctx, sqlc.CreateTaskParams{
		ID:          db.idGenerator.Generate(),
		ProjectID:   t.ProjectID,
		Title:       t.Title,
		Content:     t.Content,
		Priority:    uint32(t.Priority),
		DueOn:       dueOn,
		CompletedAt: sql.NullTime{},
		CreatedAt:   now,
		UpdatedAt:   now,
	}); err != nil {
		return fmt.Errorf("q.CreateTask failed: %w", err)
	}
	return nil
}

func (db *DB) GetTasksByProjectID(ctx context.Context, projectID string, limit, offset int) ([]*entity.Task, error) {
	tasks, err := sqlc.New(db._db).GetTasksByProjectID(ctx, sqlc.GetTasksByProjectIDParams{
		ProjectID: projectID,
		Limit:     int32(limit),
		Offset:    int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("q.GetTasksByProjectID failed: %w", err)
	}

	ts := make([]*entity.Task, 0, len(tasks))
	for _, t := range tasks {
		var dueOn *time.Time
		if t.DueOn.Valid {
			dueOn = &t.DueOn.Time
		}
		var completedAt *time.Time
		if t.CompletedAt.Valid {
			completedAt = &t.CompletedAt.Time
		}
		ts = append(ts, &entity.Task{
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

func (db *DB) GetTaskByID(ctx context.Context, id string) (*entity.Task, error) {
	t, err := sqlc.New(db._db).GetTaskByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("q.GetTaskByID failed: %w", err)
		}
	}

	var dueOn *time.Time
	if t.DueOn.Valid {
		dueOn = &t.DueOn.Time
	}
	var completedAt *time.Time
	if t.CompletedAt.Valid {
		completedAt = &t.CompletedAt.Time
	}
	return &entity.Task{
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

func (db *DB) DeleteTask(ctx context.Context, id string) error {
	if err := sqlc.New(db._db).DeleteTask(ctx, id); err != nil {
		return fmt.Errorf("q.DeleteTask failed: %w", err)
	}
	return nil
}
