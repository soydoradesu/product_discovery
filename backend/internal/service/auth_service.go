package service

import (
	"context"
	"errors"
	"strings"

	"github.com/soydoradesu/product_discovery/internal/repository"
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
		if errors.Is(err, repository.ErrNotFound) {
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

func (s *AuthService) OAuthLogin(ctx context.Context, email, googleID string) (int64, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	googleID = strings.TrimSpace(googleID)

	if email == "" || googleID == "" {
		return 0, ErrInvalidCredentials
	}

	// existing google account
	u, err := s.Users.GetByGoogleID(ctx, googleID)
	if err == nil {
		return u.ID, nil
	}
	if !errors.Is(err, repository.ErrNotFound) {
		return 0, err
	}

	// existing email account then link with google ID
	u2, err := s.Users.GetByEmail(ctx, email)
	if err == nil {
		if u2.GoogleID != nil && *u2.GoogleID != "" && *u2.GoogleID != googleID {
			return 0, ErrOAuthAccountConflict
		}
		if u2.GoogleID == nil || *u2.GoogleID == "" {
			if err := s.Users.SetGoogleID(ctx, u2.ID, googleID); err != nil {
				return 0, err
			}
		}
		return u2.ID, nil
	}
	if !errors.Is(err, repository.ErrNotFound) {
		return 0, err
	}

	// new user
	id, err := s.Users.CreateOAuthUser(ctx, email, googleID)
	if err != nil {
		return 0, err
	}
	return id, nil
}