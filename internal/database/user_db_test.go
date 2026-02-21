package database_test

import (
	"context"
	"testing"
	"time"

	"github.com/minguu42/harmattan/internal/database"
	"github.com/minguu42/harmattan/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_GetUserByID(t *testing.T) {
	ctx := context.Background()

	require.NoError(t, tdb.TruncateAndInsert(ctx, []any{
		database.Users{
			{ID: "user1", Email: "user1@dummy.invalid", HashedPassword: "hashedpassword1", CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst)},
			{ID: "user2", Email: "user2@dummy.invalid", HashedPassword: "hashedpassword2", CreatedAt: time.Date(2025, 1, 2, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 2, 0, 0, 0, 0, jst)},
			{ID: "user3", Email: "user3@dummy.invalid", HashedPassword: "hashedpassword3", CreatedAt: time.Date(2025, 1, 3, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 3, 0, 0, 0, 0, jst)},
		},
	}))

	tests := []struct {
		name    string
		input   domain.UserID
		want    *domain.User
		wantErr error
	}{
		{
			name:  "returns_user_when_exists",
			input: "user1",
			want: &domain.User{
				ID:             "user1",
				Email:          "user1@dummy.invalid",
				HashedPassword: "hashedpassword1",
			},
		},
		{
			name:    "returns_error_when_not_found",
			input:   "nonexistent",
			wantErr: database.ErrNotFound,
		},
		{
			name:  "returns_correct_user_among_multiple_users",
			input: "user2",
			want: &domain.User{
				ID:             "user2",
				Email:          "user2@dummy.invalid",
				HashedPassword: "hashedpassword2",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.GetUserByID(ctx, tt.input)

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
	ctx := context.Background()

	require.NoError(t, tdb.TruncateAndInsert(ctx, []any{
		database.Users{
			{ID: "user1", Email: "user1@dummy.invalid", HashedPassword: "hashedpassword1", CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst), UpdatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, jst)},
		},
	}))

	tests := []struct {
		name    string
		input   string
		want    *domain.User
		wantErr error
	}{
		{
			name:  "returns_user_when_email_exists",
			input: "user1@dummy.invalid",
			want: &domain.User{
				ID:             "user1",
				Email:          "user1@dummy.invalid",
				HashedPassword: "hashedpassword1",
			},
		},
		{
			name:    "returns_error_when_email_not_found",
			input:   "nonexistent@dummy.invalid",
			wantErr: database.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.GetUserByEmail(ctx, tt.input)

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
