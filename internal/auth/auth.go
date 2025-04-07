package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/minguu42/harmattan/internal/domain"
	"github.com/minguu42/harmattan/lib/clock"
)

type Authenticator struct {
	accessTokenExpiry  time.Duration
	refreshTokenExpiry time.Duration
	accessTokenSecret  string
	refreshTokenSecret string
}

type Config struct {
	AccessTokenExpiry  time.Duration `env:"ACCESS_TOKEN_EXPIRY" default:"2h"`
	RefreshTokenExpiry time.Duration `env:"REFRESH_TOKEN_EXPIRY" default:"168h"`
	AccessTokenSecret  string        `env:"ACCESS_TOKEN_SECRET,required"`
	RefreshTokenSecret string        `env:"REFRESH_TOKEN_SECRET,required"`
}

func NewAuthenticator(conf Config) (*Authenticator, error) {
	if conf.AccessTokenSecret == "" {
		return nil, errors.New("access token secret is required")
	}
	if conf.RefreshTokenSecret == "" {
		return nil, errors.New("refresh token secret is required")
	}
	return &Authenticator{
		accessTokenExpiry:  conf.AccessTokenExpiry,
		refreshTokenExpiry: conf.RefreshTokenExpiry,
		accessTokenSecret:  conf.RefreshTokenSecret,
		refreshTokenSecret: conf.AccessTokenSecret,
	}, nil
}

type accessTokenClaims struct {
	jwt.RegisteredClaims
	ID domain.UserID `json:"id"`
}

type refreshTokenClaims struct {
	jwt.RegisteredClaims
	ID domain.UserID `json:"id"`
}

func (a Authenticator) CreateAccessToken(ctx context.Context, user *domain.User) (string, error) {
	claims := &accessTokenClaims{
		ID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(clock.Now(ctx).Add(a.accessTokenExpiry)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(a.accessTokenSecret))
	if err != nil {
		return "", fmt.Errorf("failed to create signed JWT: %w", err)
	}
	return t, nil
}

func (a Authenticator) CreateRefreshToken(ctx context.Context, user *domain.User) (string, error) {
	claimsRefresh := &refreshTokenClaims{
		ID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(clock.Now(ctx).Add(a.refreshTokenExpiry)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsRefresh)
	rt, err := token.SignedString([]byte(a.refreshTokenSecret))
	if err != nil {
		return "", fmt.Errorf("failed to create signed JWT: %w", err)
	}
	return rt, nil
}

func (a Authenticator) ExtractIDFromAccessToken(token string) (domain.UserID, error) {
	jwtToken, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %s", t.Header["alg"])
		}
		return []byte(a.accessTokenSecret), nil
	})
	if err != nil {
		return "", fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok && !jwtToken.Valid {
		return "", errors.New("invalid token")
	}
	return domain.UserID(claims["id"].(string)), nil
}

func (a Authenticator) ExtractIDFromRefreshToken(tokenString string) (domain.UserID, error) {
	jwtToken, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %s", t.Header["alg"])
		}
		return []byte(a.refreshTokenSecret), nil
	})
	if err != nil {
		return "", fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok && !jwtToken.Valid {
		return "", errors.New("invalid token")
	}
	return domain.UserID(claims["id"].(string)), nil
}
