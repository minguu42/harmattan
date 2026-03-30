package api

import (
	"context"

	"github.com/minguu42/harmattan/internal/api/apierror"
	"github.com/minguu42/harmattan/internal/api/openapi"
	"github.com/minguu42/harmattan/internal/atel"
	"github.com/minguu42/harmattan/internal/auth"
	"github.com/minguu42/harmattan/internal/database"
	"github.com/minguu42/harmattan/internal/domain"
	"github.com/minguu42/harmattan/internal/lib/errtrace"
)

type securityHandler struct {
	auth *auth.Authenticator
	db   *database.Client
}

func (h *securityHandler) HandleBearerAuth(ctx context.Context, _ openapi.OperationName, t openapi.BearerAuth) (context.Context, error) {
	// セキュリティハンドラはogenミドルウェアより先に実行されるため、ここでトレースIDをロガーに付与する
	ctx = atel.ContextWithTracedLogger(ctx)

	userID, err := h.auth.ParseIDToken(ctx, t.Token)
	if err != nil {
		return nil, errtrace.Wrap(apierror.AuthorizationError())
	}

	user, err := h.db.GetUserByID(ctx, userID)
	if err != nil {
		return nil, errtrace.Wrap(err)
	}
	return domain.ContextWithUser(ctx, user), nil
}
