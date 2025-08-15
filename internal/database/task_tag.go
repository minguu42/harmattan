package database

import (
	"time"

	"github.com/minguu42/harmattan/internal/domain"
)

type TaskTag struct {
	TaskID    domain.TaskID
	TagID     domain.TagID
	CreatedAt time.Time
}

type TaskTags []TaskTag

func (tts TaskTags) ByTaskID() map[domain.TaskID]TaskTags {
	m := map[domain.TaskID]TaskTags{}
	for _, tt := range tts {
		m[tt.TaskID] = append(m[tt.TaskID], tt)
	}
	return m
}

func (tts TaskTags) TagIDs() []domain.TagID {
	ids := make([]domain.TagID, 0, len(tts))
	for _, tt := range tts {
		ids = append(ids, tt.TagID)
	}
	return ids
}
