package handler

import (
	"context"
	"fmt"

	openapi2 "github.com/minguu42/harmattan/internal/api/handler/openapi"
	"github.com/minguu42/harmattan/internal/api/usecase"
	"github.com/minguu42/harmattan/internal/auth"
	"github.com/minguu42/harmattan/internal/database"
)

type securityHandler struct {
	auth *auth.Authenticator
	db   *database.Client
}

func (h *securityHandler) HandleBearerAuth(ctx context.Context, _ openapi2.OperationName, t openapi2.BearerAuth) (context.Context, error) {
	userID, err := h.auth.ParseIDToken(ctx, t.Token)
	if err != nil {
		return nil, usecase.AuthorizationError(fmt.Errorf("failed to parse id token: %w", err))
	}

	u, err := h.db.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}
	return auth.ContextWithUser(ctx, u), nil
}
