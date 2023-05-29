package app

import "time"

type user struct {
	id        int64
	name      string
	createdAt time.Time
	updatedAt time.Time
}
