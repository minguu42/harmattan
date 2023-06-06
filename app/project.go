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

// newProjectResponse はモデル Project からレスポンスモデルの ogen.Project を生成する
func newProjectResponse(p *Project) ogen.Project {
	return ogen.Project{
		ID:        p.ID,
		Name:      p.Name,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}
