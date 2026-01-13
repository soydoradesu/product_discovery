package service

import (
	"context"
	"errors"

	"github.com/soydoradesu/product_discovery/internal/repository"
	"github.com/soydoradesu/product_discovery/internal/repository/postgres"
)

type AuthService struct {
	Users repository.UserRepository
}

func NewAuthService(users repository.UserRepository) *AuthService {
	return &AuthService{Users: users}
}

func (s *AuthService) Login(ctx context.Context, email, password string) (int64, error) {
	u, err := s.Users.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, postgres.ErrNotFound) {
			return 0, ErrUserNotFound
		}
		return 0, err
	}

	if u.PasswordHash == nil || *u.PasswordHash == "" {
		return 0, ErrInvalidCredentials
	}
	if !CheckPassword(*u.PasswordHash, password) {
		return 0, ErrInvalidCredentials
	}

	return u.ID, nil
}