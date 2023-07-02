package entity

import "time"

type Project struct {
	ID        int64
	UserID    int64
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// ContainsTask はプロジェクト p がタスク t を含むかを返す
func (p *Project) ContainsTask(t *Task) bool {
	return p.ID == t.ProjectID
}
