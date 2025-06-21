package handler

import (
	"context"
	"fmt"

	"github.com/minguu42/harmattan/api/usecase"
	"github.com/minguu42/harmattan/internal/oapi"
)

func (h *handler) SignUp(ctx context.Context, req *oapi.SignUpReq) (*oapi.SignUpOK, error) {
	out, err := h.authentication.SignUp(ctx, &usecase.SignUpInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute SignUp usecase: %w", err)
	}
	return &oapi.SignUpOK{IDToken: out.IDToken}, nil
}

func (h *handler) SignIn(ctx context.Context, req *oapi.SignInReq) (*oapi.SignInOK, error) {
	out, err := h.authentication.SignIn(ctx, &usecase.SignInInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute SignIn usecase: %w", err)
	}
	return &oapi.SignInOK{IDToken: out.IDToken}, nil
}
