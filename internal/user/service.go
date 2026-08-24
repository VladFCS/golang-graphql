package user

import (
	"context"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repository Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repository: repo,
	}
}

func (s *Service) GetUsers(ctx context.Context) ([]User, error) {
	return s.repository.List(ctx)
}

func (s *Service) GetUserByID(ctx context.Context, id string) (*User, error) {
	return s.repository.FindByID(ctx, id)
}

func (s *Service) CreateUser(ctx context.Context, input CreateUserInput) (*User, error) {
	email := strings.TrimSpace(input.Email)
	if email == "" {
		return nil, ErrEmailMustPresent
	}

	username := strings.TrimSpace(input.Username)
	if username == "" {
		return nil, ErrUsernameMustPresent
	}

	if strings.TrimSpace(input.Password) == "" {
		return nil, ErrPasswordMustPresent
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	user := User{
		Email:        email,
		Username:     username,
		PasswordHash: string(passwordHash),
	}
	createdUser, err := s.repository.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return createdUser, nil
}
