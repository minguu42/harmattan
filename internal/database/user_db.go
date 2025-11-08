package database

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/minguu42/harmattan/internal/domain"
	"github.com/minguu42/harmattan/internal/lib/clock"
	"gorm.io/gorm"
)

type User struct {
	ID             string
	Email          string
	HashedPassword string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (u *User) ToDomain() *domain.User {
	return &domain.User{
		ID:             domain.UserID(u.ID),
		Email:          u.Email,
		HashedPassword: u.HashedPassword,
	}
}

type Users []User

func (c *Client) CreateUser(ctx context.Context, u *domain.User) error {
	now := clock.Now(ctx)
	user := User{
		ID:             string(u.ID),
		Email:          u.Email,
		HashedPassword: u.HashedPassword,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	if err := c.db(ctx).Create(&user).Error; err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (c *Client) GetUserByID(ctx context.Context, id domain.UserID) (*domain.User, error) {
	var u User
	if err := c.db(ctx).Where("id = ?", string(id)).Take(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrModelNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return u.ToDomain(), nil
}

func (c *Client) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	var u User
	if err := c.db(ctx).Where("email = ?", email).Take(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrModelNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return u.ToDomain(), nil
}
