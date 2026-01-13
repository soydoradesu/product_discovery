package repository

import (
	"context"

	"github.com/soydoradesu/product_discovery/internal/domain"
)

type UserRepository interface {
	GetByEmail(ctx context.Context, email string) (domain.User, error)
	GetByID(ctx context.Context, id int64) (domain.User, error)
}