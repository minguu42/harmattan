package database

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/minguu42/harmattan/internal/domain"
	"github.com/minguu42/harmattan/internal/lib/pointers"
	"gorm.io/gorm"
)

type Step struct {
	ID          domain.StepID
	UserID      domain.UserID
	TaskID      domain.TaskID
	Name        string
	CompletedAt sql.Null[time.Time]
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (s *Step) ToDomain() *domain.Step {
	return &domain.Step{
		ID:          s.ID,
		UserID:      s.UserID,
		TaskID:      s.TaskID,
		Name:        s.Name,
		CompletedAt: pointers.RefOrNil(!s.CompletedAt.Valid, s.CompletedAt.V),
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
	}
}

type Steps []Step

func (ss Steps) ToDomain() domain.Steps {
	steps := make(domain.Steps, 0, len(ss))
	for _, s := range ss {
		steps = append(steps, *s.ToDomain())
	}
	return steps
}

func (c *Client) CreateStep(ctx context.Context, s *domain.Step) error {
	step := Step{
		ID:          s.ID,
		UserID:      s.UserID,
		TaskID:      s.TaskID,
		Name:        s.Name,
		CompletedAt: sql.Null[time.Time]{V: pointers.OrZero(s.CompletedAt), Valid: s.CompletedAt != nil},
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
	}
	return c.db(ctx).Create(&step).Error
}

func (c *Client) GetStepByID(ctx context.Context, id domain.StepID) (*domain.Step, error) {
	var s Step
	if err := c.db(ctx).Where("id = ?", id).Take(&s).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrModelNotFound
		}
		return nil, err
	}
	return s.ToDomain(), nil
}

func (c *Client) UpdateStep(ctx context.Context, s *domain.Step) error {
	return c.db(ctx).Model(Step{}).Where("id = ?", s.ID).Updates(map[string]any{
		"name":         s.Name,
		"completed_at": s.CompletedAt,
		"updated_at":   s.UpdatedAt,
	}).Error
}

func (c *Client) DeleteStepByID(ctx context.Context, id domain.StepID) error {
	return c.db(ctx).Where("id = ?", id).Delete(Step{}).Error
}
