package domain

type UserID string

type User struct {
	ID             UserID
	Email          string
	HashedPassword string
}

func (u *User) HasProject(p *Project) bool {
	return u.ID == p.UserID
}
