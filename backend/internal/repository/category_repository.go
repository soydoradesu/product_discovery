package repository

import (
	"context"

	"github.com/soydoradesu/product_discovery/internal/domain"
)

type CategoryRepository interface {
	List(ctx context.Context) ([]domain.Category, error)
}