package handler

import (
	"context"

	"github.com/minguu42/harmattan/internal/api/openapi"
	"github.com/minguu42/harmattan/internal/api/usecase"
	"github.com/minguu42/harmattan/internal/auth"
	"github.com/minguu42/harmattan/internal/database"
	"github.com/minguu42/harmattan/internal/lib/errtrace"
)

type securityHandler struct {
	auth *auth.Authenticator
	db   *database.Client
}

func (h *securityHandler) HandleBearerAuth(ctx context.Context, _ openapi.OperationName, t openapi.BearerAuth) (context.Context, error) {
	userID, err := h.auth.ParseIDToken(ctx, t.Token)
	if err != nil {
		return nil, usecase.AuthorizationError()
	}

	u, err := h.db.GetUserByID(ctx, userID)
	if err != nil {
		return nil, errtrace.Wrap(err)
	}
	return auth.ContextWithUser(ctx, u), nil
}
