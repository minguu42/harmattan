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

func (u *User) HasTask(t *Task) bool {
	return u.ID == t.UserID
}

func (u *User) HasStep(s *Step) bool {
	return u.ID == s.UserID
}

func (u *User) HasTag(t *Tag) bool {
	return u.ID == t.UserID
}
