package handler

import (
	"context"
	"fmt"

	"github.com/minguu42/harmattan/api/apperr"
	"github.com/minguu42/harmattan/internal/auth"
	"github.com/minguu42/harmattan/internal/database"
	"github.com/minguu42/harmattan/internal/openapi"
)

type securityHandler struct {
	auth *auth.Authenticator
	db   *database.Client
}

func (h *securityHandler) HandleBearerAuth(ctx context.Context, _ openapi.OperationName, t openapi.BearerAuth) (context.Context, error) {
	userID, err := h.auth.ParseIDToken(ctx, t.Token)
	if err != nil {
		return nil, apperr.AuthorizationError(fmt.Errorf("failed to parse id token: %w", err))
	}

	u, err := h.db.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}
	return auth.ContextWithUser(ctx, u), nil
}
