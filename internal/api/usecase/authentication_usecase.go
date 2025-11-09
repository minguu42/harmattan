package usecase

import (
	"context"
	"errors"

	"github.com/minguu42/harmattan/internal/auth"
	"github.com/minguu42/harmattan/internal/database"
	"github.com/minguu42/harmattan/internal/domain"
	"github.com/minguu42/harmattan/internal/lib/errtrace"
	"github.com/minguu42/harmattan/internal/lib/idgen"
	"golang.org/x/crypto/bcrypt"
)

type Authentication struct {
	Auth *auth.Authenticator
	DB   *database.Client
}

type SignUpInput struct {
	Email    string
	Password string
}

func (in *SignUpInput) User(ctx context.Context) (*domain.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errtrace.Wrap(err)
	}
	return &domain.User{
		ID:             domain.UserID(idgen.ULID(ctx)),
		Email:          in.Email,
		HashedPassword: string(hashedPassword),
	}, nil
}

type SignUpOutput struct {
	IDToken string
}

func (uc *Authentication) SignUp(ctx context.Context, in *SignUpInput) (*SignUpOutput, error) {
	u, err := uc.DB.GetUserByEmail(ctx, in.Email)
	if err != nil && !errors.Is(err, database.ErrModelNotFound) {
		return nil, errtrace.Wrap(err)
	}
	if u != nil {
		return nil, errtrace.Wrap(DuplicateUserEmailError())
	}

	user, err := in.User(ctx)
	if err != nil {
		return nil, errtrace.Wrap(err)
	}

	token, err := uc.Auth.CreateIDToken(ctx, user.ID)
	if err != nil {
		return nil, errtrace.Wrap(err)
	}

	if err := uc.DB.CreateUser(ctx, user); err != nil {
		return nil, errtrace.Wrap(err)
	}
	return &SignUpOutput{IDToken: token}, nil
}

type SignInInput struct {
	Email    string
	Password string
}

type SignInOutput struct {
	IDToken string
}

func (uc *Authentication) SignIn(ctx context.Context, in *SignInInput) (*SignInOutput, error) {
	user, err := uc.DB.GetUserByEmail(ctx, in.Email)
	if err != nil {
		if errors.Is(err, database.ErrModelNotFound) {
			return nil, errtrace.Wrap(InvalidEmailOrPasswordError())
		}
		return nil, errtrace.Wrap(err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(in.Password)); err != nil {
		return nil, errtrace.Wrap(InvalidEmailOrPasswordError())
	}

	token, err := uc.Auth.CreateIDToken(ctx, user.ID)
	if err != nil {
		return nil, errtrace.Wrap(err)
	}
	return &SignInOutput{IDToken: token}, nil
}
