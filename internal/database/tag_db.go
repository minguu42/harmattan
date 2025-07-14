package database

import (
	"time"

	"github.com/minguu42/harmattan/internal/domain"
)

type Tag struct {
	ID        domain.TagID
	UserID    domain.UserID
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
