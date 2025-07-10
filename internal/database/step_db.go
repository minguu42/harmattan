package database

import (
	"context"
	"errors"
	"time"

	"github.com/minguu42/harmattan/internal/domain"
	"gorm.io/gorm"
)

type Step struct {
	ID          domain.StepID
	UserID      domain.UserID
	TaskID      domain.TaskID
	Name        string
	CompletedAt *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (s *Step) ToDomain() *domain.Step {
	return &domain.Step{
		ID:          s.ID,
		UserID:      s.UserID,
		TaskID:      s.TaskID,
		Name:        s.Name,
		CompletedAt: s.CompletedAt,
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
		CompletedAt: s.CompletedAt,
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
	return c.db(ctx).Model(Step{}).Where("id = ?", s.ID).Updates(Step{
		Name:        s.Name,
		CompletedAt: s.CompletedAt,
		UpdatedAt:   s.UpdatedAt,
	}).Error
}

func (c *Client) DeleteStepByID(ctx context.Context, id domain.StepID) error {
	return c.db(ctx).Where("id = ?", id).Delete(Step{}).Error
}
