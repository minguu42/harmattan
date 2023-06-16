package app

import (
	"context"

	"github.com/minguu42/mtasks/app/ogen"
)

type securityHandler struct {
	repository repository
}

type userKey struct{}

// HandleIsAuthorized -
func (s *securityHandler) HandleIsAuthorized(ctx context.Context, _ string, t ogen.IsAuthorized) (context.Context, error) {
	u, err := s.repository.GetUserByAPIKey(ctx, t.APIKey)
	if err != nil {
		return nil, errUnauthorized
	}

	return context.WithValue(ctx, userKey{}, u), nil
}
