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
