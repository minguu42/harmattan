package database

import (
	"context"
	"time"

	"github.com/minguu42/harmattan/internal/domain"
)

type Project struct {
	ID         domain.ProjectID
	UserID     domain.UserID
	Name       string
	Color      string
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
