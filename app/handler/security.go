package handler

import (
	"context"

	"github.com/minguu42/mtasks/app/ogen"
	"github.com/minguu42/mtasks/app/repository"
)

type userKey struct{}

// Security -
type Security struct {
	Repository repository.Repository
}

// HandleIsAuthorized -
func (s *Security) HandleIsAuthorized(ctx context.Context, _ string, t ogen.IsAuthorized) (context.Context, error) {
	u, err := s.Repository.GetUserByAPIKey(ctx, t.APIKey)
	if err != nil {
		return nil, errUnauthorized
	}

	return context.WithValue(ctx, userKey{}, u), nil
}
