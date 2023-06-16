package app

import (
	"context"

	"github.com/minguu42/mtasks/app/ogen"
)

type SecurityHandler struct {
	Repository repository
}

type userKey struct{}

// HandleIsAuthorized -
func (s *SecurityHandler) HandleIsAuthorized(ctx context.Context, _ string, t ogen.IsAuthorized) (context.Context, error) {
	u, err := s.Repository.GetUserByAPIKey(ctx, t.APIKey)
	if err != nil {
		return nil, errUnauthorized
	}

	return context.WithValue(ctx, userKey{}, u), nil
}
