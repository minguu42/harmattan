package database

import (
	"context"
	"errors"
	"time"

	"github.com/minguu42/harmattan/internal/domain"
	"gorm.io/gorm"
)

type Project struct {
	ID         domain.ProjectID
	UserID     domain.UserID
	Name       string
	Color      domain.ProjectColor
	IsArchived bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (p *Project) ToDomain() *domain.Project {
	return &domain.Project{
		ID:         p.ID,
		UserID:     p.UserID,
		Name:       p.Name,
		Color:      p.Color,
		IsArchived: p.IsArchived,
		CreatedAt:  p.CreatedAt,
		UpdatedAt:  p.UpdatedAt,
	}
}

type Projects []Project

func (ps Projects) ToDomain() domain.Projects {
	projects := make(domain.Projects, 0, len(ps))
	for _, p := range ps {
		projects = append(projects, *p.ToDomain())
	}
	return projects
}

func (c *Client) CreateProject(ctx context.Context, p *domain.Project) error {
	project := Project{
		ID:         p.ID,
		UserID:     p.UserID,
		Name:       p.Name,
		Color:      p.Color,
		IsArchived: p.IsArchived,
		CreatedAt:  p.CreatedAt,
		UpdatedAt:  p.UpdatedAt,
	}
	return c.db(ctx).Create(&project).Error
}

func (c *Client) ListProjects(ctx context.Context, id domain.UserID, limit, offset int) (domain.Projects, error) {
	var ps Projects
	if err := c.db(ctx).Where("user_id = ?", id).Limit(limit).Offset(offset).Find(&ps).Error; err != nil {
		return nil, err
	}
	return ps.ToDomain(), nil
}

func (c *Client) GetProjectByID(ctx context.Context, id domain.ProjectID) (*domain.Project, error) {
	var p Project
	if err := c.db(ctx).Where("id = ?", id).Take(&p).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrModelNotFound
		}
		return nil, err
	}
	return p.ToDomain(), nil
}

func (c *Client) UpdateProject(ctx context.Context, p *domain.Project) error {
	return c.db(ctx).Model(Project{}).Where("id = ?", p.ID).Updates(map[string]any{
		"name":        p.Name,
		"color":       p.Color,
		"is_archived": p.IsArchived,
		"updated_at":  p.UpdatedAt,
	}).Error
}

func (c *Client) DeleteProjectByID(ctx context.Context, id domain.ProjectID) error {
	return c.db(ctx).Where("id = ?", id).Delete(Project{}).Error
}
