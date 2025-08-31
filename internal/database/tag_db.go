package database

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/minguu42/harmattan/internal/domain"
	"gorm.io/gorm"
)

type Tag struct {
	ID        domain.TagID
	UserID    domain.UserID
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (t *Tag) ToDomain() *domain.Tag {
	return &domain.Tag{
		ID:        t.ID,
		UserID:    t.UserID,
		Name:      t.Name,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
}

type Tags []Tag

func (ts Tags) ToDomain() domain.Tags {
	tags := make(domain.Tags, 0, len(ts))
	for _, t := range ts {
		tags = append(tags, *t.ToDomain())
	}
	return tags
}

func (c *Client) CreateTag(ctx context.Context, t *domain.Tag) error {
	tag := Tag{
		ID:        t.ID,
		UserID:    t.UserID,
		Name:      t.Name,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
	if err := c.db(ctx).Create(&tag).Error; err != nil {
		return fmt.Errorf("failed to create tag: %w", err)
	}
	return nil
}

func (c *Client) ListTags(ctx context.Context, id domain.UserID, limit, offset int) (domain.Tags, error) {
	var ts Tags
	if err := c.db(ctx).Where("user_id = ?", id).Limit(limit).Offset(offset).Find(&ts).Error; err != nil {
		return nil, fmt.Errorf("failed to get tags: %w", err)
	}
	return ts.ToDomain(), nil
}

func (c *Client) GetTagByID(ctx context.Context, id domain.TagID) (*domain.Tag, error) {
	var t Tag
	if err := c.db(ctx).Where("id = ?", id).Take(&t).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrModelNotFound
		}
		return nil, fmt.Errorf("failed to get tag: %w", err)
	}
	return t.ToDomain(), nil
}

func (c *Client) GetTagsByIDs(ctx context.Context, ids []domain.TagID) (domain.Tags, error) {
	if len(ids) == 0 {
		return domain.Tags{}, nil
	}

	var ts Tags
	if err := c.db(ctx).Where("id in ?", ids).Find(&ts).Error; err != nil {
		return nil, fmt.Errorf("failed to query tags: %w", err)
	}
	return ts.ToDomain(), nil
}

func (c *Client) UpdateTag(ctx context.Context, t *domain.Tag) error {
	err := c.db(ctx).Model(Tag{}).Where("id = ?", t.ID).Updates(map[string]any{
		"name":       t.Name,
		"updated_at": t.UpdatedAt,
	}).Error
	if err != nil {
		return fmt.Errorf("failed to update tag: %w", err)
	}
	return nil
}

func (c *Client) DeleteTagByID(ctx context.Context, id domain.TagID) error {
	if err := c.db(ctx).Where("id = ?", id).Delete(Tag{}).Error; err != nil {
		return fmt.Errorf("failed to delete tag: %w", err)
	}
	return nil
}
