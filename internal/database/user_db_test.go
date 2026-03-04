package database_test

import (
	"testing"
	"time"

	"github.com/minguu42/harmattan/internal/database"
	"github.com/minguu42/harmattan/internal/domain"
	"github.com/minguu42/harmattan/internal/lib/clock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_CreateUser(t *testing.T) {
	require.NoError(t, tdb.TruncateAndInsert(t.Context(), []any{
		database.Users{},
	}))

	now := time.Date(2025, 1, 1, 0, 0, 1, 0, jst)
	err := c.CreateUser(clock.WithFixedNow(t.Context(), now), &domain.User{
		ID:             "user01",
		Email:          "user01@dummy.invalid",
		HashedPassword: "hashedpassword1",
	})
	require.NoError(t, err)

	tdb.Assert(t, []any{
		database.Users{
			{ID: "user01", Email: "user01@dummy.invalid", HashedPassword: "hashedpassword1", CreatedAt: now, UpdatedAt: now},
		},
	})
}

func TestClient_GetUserByID(t *testing.T) {
	require.NoError(t, tdb.TruncateAndInsert(t.Context(), []any{
		database.Users{
			{ID: "user01", Email: "user01@dummy.invalid", HashedPassword: "hashedpassword1", CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst)},
		},
	}))

	tests := []struct {
		name    string
		id      domain.UserID
		want    *domain.User
		wantErr error
	}{
		{
			name: "found",
			id:   "user01",
			want: &domain.User{
				ID:             "user01",
				Email:          "user01@dummy.invalid",
				HashedPassword: "hashedpassword1",
			},
		},
		{
			name:    "not_found",
			id:      "user99",
			wantErr: database.ErrNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.GetUserByID(t.Context(), tt.id)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, got)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestClient_GetUserByEmail(t *testing.T) {
	require.NoError(t, tdb.TruncateAndInsert(t.Context(), []any{
		database.Users{
			{ID: "user01", Email: "user01@dummy.invalid", HashedPassword: "hashedpassword1", CreatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 1, 0, jst)},
		},
	}))

	tests := []struct {
		name    string
		email   string
		want    *domain.User
		wantErr error
	}{
		{
			name:  "found",
			email: "user01@dummy.invalid",
			want: &domain.User{
				ID:             "user01",
				Email:          "user01@dummy.invalid",
				HashedPassword: "hashedpassword1",
			},
		},
		{
			name:    "not_found",
			email:   "user99@dummy.invalid",
			wantErr: database.ErrNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.GetUserByEmail(t.Context(), tt.email)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, got)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
