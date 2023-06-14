package app

import (
	"time"

	"github.com/minguu42/mtasks/app/ogen"
)

type Project struct {
	ID        int64
	UserID    int64
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// newProjectResponse は Project から ogen.Project を生成する
func newProjectResponse(p *Project) *ogen.Project {
	return &ogen.Project{
		ID:        p.ID,
		Name:      p.Name,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}

// newProjectsResponse は Project のスライスから ogen.Project のスライスを生成する
func newProjectsResponse(ps []*Project) []ogen.Project {
	projects := make([]ogen.Project, 0, len(ps))
	for _, p := range ps {
		projects = append(projects, *newProjectResponse(p))
	}
	return projects
}
