package domain

import (
	"slices"
	"time"
)

type TaskID string

type Task struct {
	ID          TaskID
	UserID      UserID
	ProjectID   ProjectID
	Name        string
	TagIDs      []TagID
	Content     string
	Priority    int
	DueOn       *time.Time
	CompletedAt *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Steps       Steps
}

type Tasks []Task

func (ts Tasks) TagIDs() []TagID {
	var tagIDs []TagID
	for _, t := range ts {
		tagIDs = append(tagIDs, t.TagIDs...)
	}
	slices.Sort(tagIDs)
	return slices.Compact(tagIDs)
}
