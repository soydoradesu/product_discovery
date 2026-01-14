package repository

import (
	"context"

	"github.com/soydoradesu/product_discovery/internal/domain"
)

type ProductRepository interface {
	GetByID(ctx context.Context, id int64) (domain.Product, error)
}