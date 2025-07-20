package handler

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"regexp"
	"strings"
	"unicode"

	"github.com/minguu42/harmattan/api/apperr"
	"github.com/minguu42/harmattan/api/usecase"
	"github.com/minguu42/harmattan/internal/openapi"
)

var (
	ErrInvalidEmail             = errors.New("メールアドレスの形式が正しくありません。")
	ErrInvalidEmailLength       = errors.New("メールアドレスは3文字以上254文字以下で指定できます。")
	ErrInvalidPasswordLength    = errors.New("パスワードは12文字以上64文字以下で指定できます。")
	ErrInvalidPassword          = errors.New("パスワードは英数字と一部の記号のみ使用できます。")
	ErrPasswordMissingUppercase = errors.New("パスワードには少なくとも1つの大文字が必要です。")
	ErrPasswordMissingLowercase = errors.New("パスワードには少なくとも1つの小文字が必要です。")
	ErrPasswordMissingNumber    = errors.New("パスワードには少なくとも1つの数字が必要です。")
	ErrPasswordMissingSymbol    = errors.New("パスワードには少なくとも1つの記号が必要です。")
)

func (h *handler) SignUp(ctx context.Context, req *openapi.SignUpReq) (*openapi.SignUpOK, error) {
	var validationErrors []error
	if errs := validateEmail(req.Email); len(errs) > 0 {
		validationErrors = append(validationErrors, errs...)
	}
	if errs := validatePassword(req.Password); len(errs) > 0 {
		validationErrors = append(validationErrors, errs...)
	}
	if len(validationErrors) > 0 {
		return nil, apperr.DomainValidationError(validationErrors)
	}

	out, err := h.authentication.SignUp(ctx, &usecase.SignUpInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute SignUp usecase: %w", err)
	}
	return &openapi.SignUpOK{IDToken: out.IDToken}, nil
}

func (h *handler) SignIn(ctx context.Context, req *openapi.SignInReq) (*openapi.SignInOK, error) {
	var validationErrors []error
	if errs := validateEmail(req.Email); len(errs) > 0 {
		validationErrors = append(validationErrors, errs...)
	}
	if errs := validatePassword(req.Password); len(errs) > 0 {
		validationErrors = append(validationErrors, errs...)
	}
	if len(validationErrors) > 0 {
		return nil, apperr.DomainValidationError(validationErrors)
	}

	out, err := h.authentication.SignIn(ctx, &usecase.SignInInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute SignIn usecase: %w", err)
	}
	return &openapi.SignInOK{IDToken: out.IDToken}, nil
}

func validateEmail(email string) []error {
	var errs []error
	if addr, err := mail.ParseAddress(email); err != nil || addr.Address != email {
		errs = append(errs, ErrInvalidEmail)
	}
	if len(email) < 3 || 254 < len(email) {
		errs = append(errs, ErrInvalidEmailLength)
	}
	return errs
}

const allowedSymbol = `!@#$%^&*()-_+=[]{}|\:;"'<>,.?/`

var allowedCharsPattern = regexp.MustCompile(`^[a-zA-Z0-9` + regexp.QuoteMeta(allowedSymbol) + `]+$`)

func validatePassword(password string) []error {
	var errs []error
	if len(password) < 12 || 64 < len(password) {
		errs = append(errs, ErrInvalidPasswordLength)
	}
	if !allowedCharsPattern.MatchString(password) {
		errs = append(errs, ErrInvalidPassword)
	}
	if !strings.ContainsFunc(password, func(r rune) bool { return unicode.IsUpper(r) }) {
		errs = append(errs, ErrPasswordMissingUppercase)
	}
	if !strings.ContainsFunc(password, func(r rune) bool { return unicode.IsLower(r) }) {
		errs = append(errs, ErrPasswordMissingLowercase)
	}
	if !strings.ContainsFunc(password, func(r rune) bool { return unicode.IsNumber(r) }) {
		errs = append(errs, ErrPasswordMissingNumber)
	}
	if !strings.ContainsAny(password, allowedSymbol) {
		errs = append(errs, ErrPasswordMissingSymbol)
	}
	return errs
}
