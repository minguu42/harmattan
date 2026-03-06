package handler

import (
	"context"
	"errors"
	"net/mail"
	"regexp"
	"strings"
	"unicode"

	"github.com/minguu42/harmattan/internal/api/apierror"
	"github.com/minguu42/harmattan/internal/api/openapi"
	"github.com/minguu42/harmattan/internal/api/usecase"
	"github.com/minguu42/harmattan/internal/lib/errtrace"
)

func (h *Handler) SignUp(ctx context.Context, req *openapi.SignUpReq) (*openapi.SignUpOK, error) {
	var errs []error
	errs = append(errs, validateEmail(req.Email)...)
	errs = append(errs, validatePassword(req.Password)...)
	if len(errs) > 0 {
		return nil, errtrace.Wrap(apierror.DomainValidationError(errs))
	}

	out, err := h.Authentication.SignUp(ctx, &usecase.SignUpInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return nil, errtrace.Wrap(err)
	}
	return &openapi.SignUpOK{IDToken: out.IDToken}, nil
}

func (h *Handler) SignIn(ctx context.Context, req *openapi.SignInReq) (*openapi.SignInOK, error) {
	var errs []error
	errs = append(errs, validateEmail(req.Email)...)
	errs = append(errs, validatePassword(req.Password)...)
	if len(errs) > 0 {
		return nil, errtrace.Wrap(apierror.DomainValidationError(errs))
	}

	out, err := h.Authentication.SignIn(ctx, &usecase.SignInInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return nil, errtrace.Wrap(err)
	}
	return &openapi.SignInOK{IDToken: out.IDToken}, nil
}

var (
	ErrEmailCharacter           = errors.New("メールアドレスにはASCII文字のみ使用できます")
	ErrEmailFormat              = errors.New("メールアドレスの形式が正しくありません")
	ErrEmailLength              = errors.New("メールアドレスは3文字以上254文字以下で指定できます")
	ErrPasswordCharacter        = errors.New("パスワードは英数字と一部の記号のみ使用できます")
	ErrPasswordLength           = errors.New("パスワードは12文字以上64文字以下で指定できます")
	ErrPasswordMissingUppercase = errors.New("パスワードには少なくとも1つの大文字が必要です")
	ErrPasswordMissingLowercase = errors.New("パスワードには少なくとも1つの小文字が必要です")
	ErrPasswordMissingDigit     = errors.New("パスワードには少なくとも1つの数字が必要です")
	ErrPasswordMissingSymbol    = errors.New("パスワードには少なくとも1つの記号が必要です")
)

func validateEmail(email string) []error {
	var errs []error
	if strings.ContainsFunc(email, func(r rune) bool { return r > unicode.MaxASCII }) {
		errs = append(errs, ErrEmailCharacter)
	}
	if addr, err := mail.ParseAddress(email); err != nil || addr.Address != email {
		errs = append(errs, ErrEmailFormat)
	}
	if len(email) < 3 || 254 < len(email) {
		errs = append(errs, ErrEmailLength)
	}
	return errs
}

const allowedPasswordSymbol = `!@#$%^&*()-_+=[]{}|\:;"'<>,.?/`

var allowedPasswordCharsPattern = regexp.MustCompile(`^[a-zA-Z0-9` + regexp.QuoteMeta(allowedPasswordSymbol) + `]+$`)

func validatePassword(password string) []error {
	var errs []error
	if !allowedPasswordCharsPattern.MatchString(password) {
		errs = append(errs, ErrPasswordCharacter)
	}
	if len(password) < 12 || 64 < len(password) {
		errs = append(errs, ErrPasswordLength)
	}
	if !strings.ContainsFunc(password, unicode.IsUpper) {
		errs = append(errs, ErrPasswordMissingUppercase)
	}
	if !strings.ContainsFunc(password, unicode.IsLetter) {
		errs = append(errs, ErrPasswordMissingLowercase)
	}
	if !strings.ContainsFunc(password, func(r rune) bool { return r >= '0' && r <= '9' }) {
		errs = append(errs, ErrPasswordMissingDigit)
	}
	if !strings.ContainsAny(password, allowedPasswordSymbol) {
		errs = append(errs, ErrPasswordMissingSymbol)
	}
	return errs
}
