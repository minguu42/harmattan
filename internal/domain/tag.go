package domain

type TagID string

type Tag struct {
	ID     TagID
	UserID UserID
	Name   string
}

type Tags []Tag
