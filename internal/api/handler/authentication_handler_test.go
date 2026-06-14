package handler_test

import (
	"strings"
	"testing"

	"github.com/minguu42/harmattan/internal/api/handler"
	"github.com/stretchr/testify/assert"
)

func TestValidateEmail(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		email string
		want  []error
	}{
		{name: "valid", email: "user@example.com"},
		{name: "min_length_boundary", email: "a@b"},
		{name: "below_min_length", email: "a@", want: []error{handler.ErrEmailFormat, handler.ErrEmailLength}},
		{name: "max_length_boundary", email: strings.Repeat("a", 242) + "@example.com"},
		{name: "above_max_length", email: strings.Repeat("a", 243) + "@example.com", want: []error{handler.ErrEmailLength}},
		{name: "invalid_format", email: "userexample.com", want: []error{handler.ErrEmailFormat}},
		{name: "non_ascii", email: "üser@example.com", want: []error{handler.ErrEmailCharacter}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.ElementsMatch(t, tt.want, handler.ValidateEmail(tt.email))
		})
	}
}

func TestValidatePassword(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		password string
		want     []error
	}{
		{name: "valid_min_length_boundary", password: "Aa1!aaaaaaaa"},
		{name: "below_min_length", password: "Aa1!aaaaaaa", want: []error{handler.ErrPasswordLength}},
		{name: "valid_max_length_boundary", password: "Aa1!" + strings.Repeat("a", 60)},
		{name: "above_max_length", password: "Aa1!" + strings.Repeat("a", 61), want: []error{handler.ErrPasswordLength}},
		{name: "missing_uppercase", password: "aa1!aaaaaaaa", want: []error{handler.ErrPasswordMissingUppercase}},
		{name: "missing_lowercase", password: "AA1!AAAAAAAA", want: []error{handler.ErrPasswordMissingLowercase}},
		{name: "missing_digit", password: "Aa!!aaaaaaaa", want: []error{handler.ErrPasswordMissingDigit}},
		{name: "missing_symbol", password: "Aa1aaaaaaaaa", want: []error{handler.ErrPasswordMissingSymbol}},
		{name: "invalid_character", password: "Aa1! aaaaaaa", want: []error{handler.ErrPasswordCharacter}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.ElementsMatch(t, tt.want, handler.ValidatePassword(tt.password))
		})
	}
}
