package domain

import "time"

type TagID string

type Tag struct {
	ID        TagID
	UserID    UserID
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Tags []Tag

func (ts Tags) TagByID() map[TagID]Tag {
	m := make(map[TagID]Tag, len(ts))
	for _, t := range ts {
		m[t.ID] = t
	}
	return m
}
