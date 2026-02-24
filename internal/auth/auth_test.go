package auth_test

import (
	"context"
	"testing"
	"time"

	"github.com/minguu42/harmattan/internal/auth"
	"github.com/minguu42/harmattan/internal/domain"
	"github.com/minguu42/harmattan/internal/lib/clock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthenticator_CreateIDToken(t *testing.T) {
	t.Parallel()

	authn, err := auth.NewAuthenticator("cIZ15duBB4CjZNxD6CH8jBgc5sP5Ch7G", 1*time.Hour)
	require.NoError(t, err)
	ctx := clock.WithFixedNow(t.Context(), time.Date(2025, 10, 1, 15, 40, 50, 0, time.UTC))

	got, err := authn.CreateIDToken(ctx, "u1")
	require.NoError(t, err)

	want := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ1MSIsImV4cCI6MTc1OTMzNjg1MCwiaWF0IjoxNzU5MzMzMjUwfQ.vV6tn3H29xdhA67JvtvJSJ-YnNKFJoC2GTYP28ibFDQ"
	assert.Equal(t, want, got)
}

func TestAuthenticator_ParseIDToken(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		ctx     context.Context
		token   string
		want    domain.UserID
		wantErr bool
	}{
		{
			name:  "valid token",
			ctx:   clock.WithFixedNow(context.Background(), time.Date(2025, 10, 1, 16, 20, 50, 0, time.UTC)),
			token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ1MSIsImV4cCI6MTc1OTMzNjg1MCwiaWF0IjoxNzU5MzMzMjUwfQ.vV6tn3H29xdhA67JvtvJSJ-YnNKFJoC2GTYP28ibFDQ",
			want:  domain.UserID("u1"),
		},
		{
			name:    "invalid token format",
			ctx:     clock.WithFixedNow(context.Background(), time.Date(2025, 10, 1, 16, 20, 50, 0, time.UTC)),
			token:   "invalid-token",
			wantErr: true,
		},
		{
			name:    "expired token",
			ctx:     clock.WithFixedNow(context.Background(), time.Date(2025, 10, 1, 17, 0, 0, 0, time.UTC)),
			token:   "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ1MSIsImV4cCI6MTc1OTMzNjg1MCwiaWF0IjoxNzU5MzMzMjUwfQ.vV6tn3H29xdhA67JvtvJSJ-YnNKFJoC2GTYP28ibFDQ",
			wantErr: true,
		},
		{
			name:    "invalid signature",
			ctx:     clock.WithFixedNow(context.Background(), time.Date(2025, 10, 1, 16, 20, 50, 0, time.UTC)),
			token:   "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTkzMzY4NTAsImlhdCI6MTc1OTMzMzI1MCwic3ViIjoidTEifQ._yo4Qd7YkZygDRrcEq6g400isF7rBNp2b0Rj9_KeWwY",
			wantErr: true,
		},
		{
			name:    "empty sub claim",
			ctx:     clock.WithFixedNow(context.Background(), time.Date(2025, 10, 1, 16, 20, 50, 0, time.UTC)),
			token:   "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIiLCJleHAiOjE3NTkzMzY4NTAsImlhdCI6MTc1OTMzMzI1MH0.KpZrhJmdPOBwK2xOTj-0quFWKPuqjDqc1rcQ66t3Y-M",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			authn, err := auth.NewAuthenticator("cIZ15duBB4CjZNxD6CH8jBgc5sP5Ch7G", 1*time.Hour)
			require.NoError(t, err)

			got, err := authn.ParseIDToken(tt.ctx, tt.token)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
