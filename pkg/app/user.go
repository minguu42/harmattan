package app

import "time"

type user struct {
	id        uint64
	name      string
	createdAt time.Time
	updatedAt time.Time
}
