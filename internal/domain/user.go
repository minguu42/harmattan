package domain

import (
	"context"
	"errors"

	"github.com/minguu42/harmattan/internal/lib/errtrace"
)

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

type userKey struct{}

func ContextWithUser(ctx context.Context, u *User) context.Context {
	return context.WithValue(ctx, userKey{}, u)
}

func UserFromContext(ctx context.Context) (*User, error) {
	user, ok := ctx.Value(userKey{}).(*User)
	if !ok {
		return nil, errtrace.Wrap(errors.New("user not found in context"))
	}
	return user, nil
}
