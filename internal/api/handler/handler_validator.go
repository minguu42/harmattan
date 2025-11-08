package handler

import (
	"net/mail"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/minguu42/harmattan/internal/lib/errors"
)

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
	ErrProjectNameLength        = errors.New("プロジェクト名は1文字以上80文字以下で指定できます")
	ErrStepNameLength           = errors.New("ステップ名は1文字以上100文字以下で指定できます")
	ErrTagNameLength            = errors.New("タグ名は1文字以上20文字以下で指定できます")
	ErrTaskNameLength           = errors.New("タスク名は1文字以上100文字以下で指定できます")
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

func validateProjectName(name string) []error {
	var errs []error
	if utf8.RuneCountInString(name) < 1 || 80 < utf8.RuneCountInString(name) {
		errs = append(errs, ErrProjectNameLength)
	}
	return errs
}

func validateStepName(name string) []error {
	var errs []error
	if utf8.RuneCountInString(name) < 1 || 100 < utf8.RuneCountInString(name) {
		errs = append(errs, ErrStepNameLength)
	}
	return errs
}

func validateTagName(name string) []error {
	var errs []error
	if utf8.RuneCountInString(name) < 1 || 20 < utf8.RuneCountInString(name) {
		errs = append(errs, ErrTagNameLength)
	}
	return errs
}

func validateTaskName(name string) []error {
	var errs []error
	if utf8.RuneCountInString(name) < 1 || 100 < utf8.RuneCountInString(name) {
		errs = append(errs, ErrTaskNameLength)
	}
	return errs
}
