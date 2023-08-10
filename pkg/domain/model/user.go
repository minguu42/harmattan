package model

import "time"

type User struct {
	ID        string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// HasProject はユーザがプロジェクトを所有しているかを返す
func (u *User) HasProject(p *Project) bool {
	return u.ID == p.UserID
}
