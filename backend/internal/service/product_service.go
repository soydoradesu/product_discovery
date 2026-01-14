package service

import (
	"context"
	"errors"
	"strings"

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

func (s *ProductService) Search(ctx context.Context, params domain.SearchParams) ([]domain.ProductSummary, int64, domain.SearchParams, error) {
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 20
	}
	if params.PageSize > 100 {
		params.PageSize = 100
	}

	params.Sort = strings.TrimSpace(strings.ToLower(params.Sort))
	params.Method = strings.TrimSpace(strings.ToLower(params.Method))

	// default sorting
	if params.Sort == "" {
		if strings.TrimSpace(params.Q) != "" {
			params.Sort = "relevance"
		} else {
			params.Sort = "created_at"
		}
	}

	// validation for sort field
	switch params.Sort {
	case "relevance", "price", "created_at", "rating":
	default:
		params.Sort = "created_at"
	}

	// relevance only makes sense when q present
	if params.Sort == "relevance" && strings.TrimSpace(params.Q) == "" {
		params.Sort = "created_at"
	}

	// default method per sort
	if params.Method == "" {
		switch params.Sort {
		case "price":
			params.Method = "asc"
		default:
			params.Method = "desc"
		}
	}

	if params.Method != "asc" && params.Method != "desc" {
		params.Method = "desc"
	}

	items, total, err := s.Products.Search(ctx, params)
	if err != nil {
		return nil, 0, params, err
	}
	return items, total, params, nil
}