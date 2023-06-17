package entity

import "time"

type User struct {
	ID        int64
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// HasProject はユーザ u がプロジェクト p を所有しているかを返す
func (u *User) HasProject(p *Project) bool {
	return u.ID == p.UserID
}
