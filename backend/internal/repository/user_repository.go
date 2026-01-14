package repository

import (
	"context"

	"github.com/soydoradesu/product_discovery/internal/domain"
)

type UserRepository interface {
	GetByEmail(ctx context.Context, email string) (domain.User, error)
	GetByID(ctx context.Context, id int64) (domain.User, error)

	// OAuth
	GetByGoogleID(ctx context.Context, googleID string) (domain.User, error)
	SetGoogleID(ctx context.Context, userID int64, googleID string) error
	CreateOAuthUser(ctx context.Context, email, googleID string) (int64, error)
}