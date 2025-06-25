package domain

import "time"

type ProjectID string

type Project struct {
	ID         ProjectID
	UserID     UserID
	Name       string
	Color      string
	IsArchived bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type Projects []Project
