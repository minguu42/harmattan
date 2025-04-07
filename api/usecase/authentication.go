package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/minguu42/harmattan/api/apperr"
	"github.com/minguu42/harmattan/internal/auth"
	"github.com/minguu42/harmattan/internal/database"
	"github.com/minguu42/harmattan/internal/domain"
	"github.com/minguu42/harmattan/lib/idgen"
	"golang.org/x/crypto/bcrypt"
)

type Authentication struct {
	auth *auth.Authenticator
	db   *database.Client
}

func NewAuthentication(auth *auth.Authenticator, db *database.Client) Authentication {
	return Authentication{
		auth: auth,
		db:   db,
	}
}

type SignUpInput struct {
	Email    string
	Password string
}

func (in *SignUpInput) User(ctx context.Context) (*domain.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to generate hashed password: %w", err)
	}
	return &domain.User{
		ID:             domain.UserID(idgen.ULID(ctx)),
		Email:          in.Email,
		HashedPassword: string(hashedPassword),
	}, nil
}

type SignUpOutput struct {
	AccessToken  string
	RefreshToken string
}

func (uc *Authentication) SignUp(ctx context.Context, in *SignUpInput) (*SignUpOutput, error) {
	if _, err := uc.db.GetUserByEmail(ctx, in.Email); !errors.Is(err, database.ErrModelNotFound) {
		return nil, apperr.ErrDuplicateUserEmail(err)
	}

	user, err := in.User(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to generate user: %w", err)
	}

	accessToken, err := uc.auth.CreateAccessToken(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create access token: %w", err)
	}
	refreshToken, err := uc.auth.CreateRefreshToken(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create refresh token: %w", err)
	}

	if err := uc.db.CreateUser(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return &SignUpOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

type SignInInput struct {
	Email    string
	Password string
}

type SignInOutput struct {
	AccessToken  string
	RefreshToken string
}

func (uc *Authentication) SignIn(ctx context.Context, in *SignInInput) (*SignInOutput, error) {
	user, err := uc.db.GetUserByEmail(ctx, in.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(in.Password)) != nil {
		return nil, errors.New("password is not valid")
	}

	accessToken, err := uc.auth.CreateAccessToken(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create access token: %w", err)
	}
	refreshToken, err := uc.auth.CreateRefreshToken(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create refresh token: %w", err)
	}
	return &SignInOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

type RefreshTokenInput struct {
	RefreshToken string
}

type RefreshTokenOutput struct {
	AccessToken  string
	RefreshToken string
}

func (uc *Authentication) RefreshToken(ctx context.Context, in *RefreshTokenInput) (*RefreshTokenOutput, error) {
	id, err := uc.auth.ExtractIDFromRefreshToken(in.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to extract id from refresh token: %w", err)
	}
	user, err := uc.db.GetUserByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	accessToken, err := uc.auth.CreateAccessToken(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create access token: %w", err)
	}
	refreshToken, err := uc.auth.CreateRefreshToken(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create refresh token: %w", err)
	}
	return &RefreshTokenOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
