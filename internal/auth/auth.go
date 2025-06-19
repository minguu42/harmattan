package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/minguu42/harmattan/internal/domain"
	"github.com/minguu42/harmattan/lib/clock"
)

type Config struct {
	IDTokenExpiration time.Duration `env:"ID_TOKEN_EXPIRATION" default:"1h"`
	IDTokenSecret     string        `env:"ID_TOKEN_SECRET,required"`
}

func NewAuthenticator(conf Config) (*Authenticator, error) {
	if conf.IDTokenSecret == "" {
		return nil, errors.New("id token secret is required")
	}
	return &Authenticator{
		idTokenExpiration: conf.IDTokenExpiration,
		idTokenSecret:     conf.IDTokenSecret,
	}, nil
}

type Authenticator struct {
	idTokenExpiration time.Duration
	idTokenSecret     string
}

func (a *Authenticator) CreateIDToken(ctx context.Context, id domain.UserID) (string, error) {
	now := clock.Now(ctx)
	claims := jwt.RegisteredClaims{
		Subject:   string(id),
		ExpiresAt: jwt.NewNumericDate(now.Add(a.idTokenExpiration)),
		IssuedAt:  jwt.NewNumericDate(now),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(a.idTokenSecret))
}

func (a *Authenticator) ParseIDToken(tokenString string) (domain.UserID, error) {
	parser := jwt.NewParser(
		jwt.WithValidMethods([]string{"HS256"}),
		jwt.WithExpirationRequired(),
		jwt.WithIssuedAt(),
	)

	token, err := parser.Parse(tokenString, func(t *jwt.Token) (any, error) {
		return []byte(a.idTokenSecret), nil
	})
	if err != nil {
		return "", fmt.Errorf("failed to parse token: %w", err)
	}

	claims := token.Claims.(jwt.MapClaims)
	if id, ok := claims["sub"].(string); ok && id != "" {
		return domain.UserID(id), nil
	}
	return "", errors.New("missing or empty sub claim")
}

type userKey struct{}

func ContextWithUser(ctx context.Context, u *domain.User) context.Context {
	return context.WithValue(ctx, userKey{}, u)
}

func UserFromContext(ctx context.Context) *domain.User {
	v, _ := ctx.Value(userKey{}).(*domain.User)
	return v
}
