package auth

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/minguu42/harmattan/internal/domain"
	"github.com/minguu42/harmattan/internal/lib/clock"
	"github.com/minguu42/harmattan/internal/lib/errtrace"
)

func NewAuthenticator(idTokenSecret string, idTokenExpiration time.Duration) (*Authenticator, error) {
	if idTokenSecret == "" {
		return nil, errtrace.Wrap(errors.New("id token secret is required"))
	}
	if idTokenExpiration == 0 {
		return nil, errtrace.Wrap(errors.New("id token expiration is required"))
	}
	return &Authenticator{
		idTokenSecret:     idTokenSecret,
		idTokenExpiration: idTokenExpiration,
	}, nil
}

type Authenticator struct {
	idTokenSecret     string
	idTokenExpiration time.Duration
}

func (a *Authenticator) CreateIDToken(ctx context.Context, id domain.UserID) (string, error) {
	now := clock.Now(ctx)
	claims := jwt.RegisteredClaims{
		Subject:   string(id),
		ExpiresAt: jwt.NewNumericDate(now.Add(a.idTokenExpiration)),
		IssuedAt:  jwt.NewNumericDate(now),
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(a.idTokenSecret))
	if err != nil {
		return "", errtrace.Wrap(err)
	}
	return token, nil
}

func (a *Authenticator) ParseIDToken(ctx context.Context, tokenString string) (domain.UserID, error) {
	parser := jwt.NewParser(
		jwt.WithValidMethods([]string{"HS256"}),
		jwt.WithExpirationRequired(),
		jwt.WithIssuedAt(),
		jwt.WithTimeFunc(func() time.Time { return clock.Now(ctx) }),
	)

	token, err := parser.Parse(tokenString, func(t *jwt.Token) (any, error) {
		return []byte(a.idTokenSecret), nil
	})
	if err != nil {
		return "", errtrace.Wrap(err)
	}

	claims := token.Claims.(jwt.MapClaims)
	if id, ok := claims["sub"].(string); ok && id != "" {
		return domain.UserID(id), nil
	}
	return "", errtrace.Wrap(errors.New("missing or empty sub claim"))
}

type userKey struct{}

func ContextWithUser(ctx context.Context, u *domain.User) context.Context {
	return context.WithValue(ctx, userKey{}, u)
}

func UserFromContext(ctx context.Context) (*domain.User, bool) {
	u, ok := ctx.Value(userKey{}).(*domain.User)
	return u, ok
}

func MustUserFromContext(ctx context.Context) *domain.User {
	if u, ok := ctx.Value(userKey{}).(*domain.User); ok {
		return u
	}
	panic("user not found in context")
}
