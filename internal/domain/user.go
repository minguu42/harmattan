package domain

type UserID string

type User struct {
	ID             UserID
	Email          string
	HashedPassword string
}
