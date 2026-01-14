package service

import (
	"context"
	"errors"

	"github.com/soydoradesu/product_discovery/internal/domain"
	"github.com/soydoradesu/product_discovery/internal/repository"
)

type ProductService struct {
	Products repository.ProductRepository
}

func NewProductService(products repository.ProductRepository) *ProductService {
	return &ProductService{Products: products}
}

func (s *ProductService) GetByID(ctx context.Context, id int64) (domain.Product, error) {
	p, err := s.Products.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return domain.Product{}, ErrProductNotFound
		}
		return domain.Product{}, err
	}
	return p, nil
}
