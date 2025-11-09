package handler

import (
	"context"

	"github.com/minguu42/harmattan/internal/api/openapi"
	"github.com/minguu42/harmattan/internal/api/usecase"
	"github.com/minguu42/harmattan/internal/lib/errtrace"
)

func (h *handler) SignUp(ctx context.Context, req *openapi.SignUpReq) (*openapi.SignUpOK, error) {
	var errs []error
	errs = append(errs, validateEmail(req.Email)...)
	errs = append(errs, validatePassword(req.Password)...)
	if len(errs) > 0 {
		return nil, usecase.DomainValidationError(errs)
	}

	out, err := h.authentication.SignUp(ctx, &usecase.SignUpInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return nil, errtrace.Wrap(err)
	}
	return &openapi.SignUpOK{IDToken: out.IDToken}, nil
}

func (h *handler) SignIn(ctx context.Context, req *openapi.SignInReq) (*openapi.SignInOK, error) {
	var errs []error
	errs = append(errs, validateEmail(req.Email)...)
	errs = append(errs, validatePassword(req.Password)...)
	if len(errs) > 0 {
		return nil, usecase.DomainValidationError(errs)
	}

	out, err := h.authentication.SignIn(ctx, &usecase.SignInInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return nil, errtrace.Wrap(err)
	}
	return &openapi.SignInOK{IDToken: out.IDToken}, nil
}
