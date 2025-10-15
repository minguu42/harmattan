package auth_test

import (
	"testing"
	"time"

	"github.com/minguu42/harmattan/internal/auth"
	"github.com/minguu42/harmattan/internal/domain"
	"github.com/minguu42/harmattan/internal/lib/clock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthenticator_CreateIDToken(t *testing.T) {
	a, err := auth.NewAuthenticator(auth.Config{
		IDTokenExpiration: 1 * time.Hour,
		IDTokenSecret:     "cIZ15duBB4CjZNxD6CH8jBgc5sP5Ch7G",
	})
	require.NoError(t, err)
	ctx := clock.WithFixedNow(t.Context(), time.Date(2025, 10, 1, 15, 40, 50, 0, time.UTC))

	got, err := a.CreateIDToken(ctx, "u1")
	require.NoError(t, err)

	want := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ1MSIsImV4cCI6MTc1OTMzNjg1MCwiaWF0IjoxNzU5MzMzMjUwfQ.vV6tn3H29xdhA67JvtvJSJ-YnNKFJoC2GTYP28ibFDQ"
	assert.Equal(t, want, got)
}

func TestAuthenticator_ParseIDToken(t *testing.T) {
	a, err := auth.NewAuthenticator(auth.Config{
		IDTokenExpiration: 1 * time.Hour,
		IDTokenSecret:     "cIZ15duBB4CjZNxD6CH8jBgc5sP5Ch7G",
	})
	require.NoError(t, err)
	ctx := clock.WithFixedNow(t.Context(), time.Date(2025, 10, 1, 16, 20, 50, 0, time.UTC))

	got, err := a.ParseIDToken(ctx, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ1MSIsImV4cCI6MTc1OTMzNjg1MCwiaWF0IjoxNzU5MzMzMjUwfQ.vV6tn3H29xdhA67JvtvJSJ-YnNKFJoC2GTYP28ibFDQ")
	require.NoError(t, err)

	want := domain.UserID("u1")
	assert.Equal(t, want, got)
}
