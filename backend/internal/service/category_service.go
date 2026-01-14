package service

import (
	"context"

	"github.com/soydoradesu/product_discovery/internal/domain"
	"github.com/soydoradesu/product_discovery/internal/repository"
)

type CategoryService struct {
	Categories repository.CategoryRepository
}

func NewCategoryService(categories repository.CategoryRepository) *CategoryService {
	return &CategoryService{Categories: categories}
}

func (s *CategoryService) List(ctx context.Context) ([]domain.Category, error) {
	return s.Categories.List(ctx)
}