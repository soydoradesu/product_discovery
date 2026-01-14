package internal_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/soydoradesu/product_discovery/internal/domain"
	"github.com/soydoradesu/product_discovery/internal/repository"
	"github.com/soydoradesu/product_discovery/internal/service"
)

type fakeProducts struct {
	byID map[int64]domain.Product

	searchItems []domain.ProductSummary
	searchTotal int64
	searchErr   error
}

func (f *fakeProducts) GetByID(ctx context.Context, id int64) (domain.Product, error) {
	p, ok := f.byID[id]
	if !ok {
		return domain.Product{}, repository.ErrNotFound
	}
	return p, nil
}

func (f *fakeProducts) Search(ctx context.Context, params domain.SearchParams) ([]domain.ProductSummary, int64, error) {
	if f.searchErr != nil {
		return nil, 0, f.searchErr
	}
	return f.searchItems, f.searchTotal, nil
}

func TestProductService_GetByID_OK(t *testing.T) {
	fp := &fakeProducts{
		byID: map[int64]domain.Product{
			1: {ID: 1, Name: "X", Price: 10, Description: "D", Rating: 4.5, InStock: true, CreatedAt: time.Now()},
		},
	}
	svc := service.NewProductService(fp)

	p, err := svc.GetByID(context.Background(), 1)
	if err != nil {
		t.Fatalf("expected nil err, got %v", err)
	}
	if p.ID != 1 {
		t.Fatalf("expected id=1 got %d", p.ID)
	}
}

func TestProductService_GetByID_NotFound(t *testing.T) {
	fp := &fakeProducts{byID: map[int64]domain.Product{}}
	svc := service.NewProductService(fp)

	_, err := svc.GetByID(context.Background(), 123)
	if err != service.ErrProductNotFound {
		t.Fatalf("expected ErrProductNotFound got %v", err)
	}
}

func TestProductService_Search_OK(t *testing.T) {
	fp := &fakeProducts{
		searchItems: []domain.ProductSummary{
			{ID: 1, Name: "A", Price: 10.5, Rating: 4.2, InStock: true, CreatedAt: time.Now()},
			{ID: 2, Name: "B", Price: 20.0, Rating: 4.8, InStock: false, CreatedAt: time.Now()},
		},
		searchTotal: 2,
	}
	svc := service.NewProductService(fp)

	params := domain.SearchParams{
		Q:        "abc",
		Page:     1,
		PageSize: 10,
		Sort:     "relevance",
		Method:   "desc",
	}

	items, total, normalized, err := svc.Search(context.Background(), params)
	if err != nil {
		t.Fatalf("expected nil err, got %v", err)
	}
	if total != 2 {
		t.Fatalf("expected total=2 got %d", total)
	}
	if len(items) != 2 {
		t.Fatalf("expected 2 items got %d", len(items))
	}
	if normalized.Page <= 0 || normalized.PageSize <= 0 {
		t.Fatalf("expected normalized page/pagesize > 0, got page=%d size=%d", normalized.Page, normalized.PageSize)
	}
}

func TestProductService_Search_PropagatesError(t *testing.T) {
	fp := &fakeProducts{searchErr: errors.New("db down")}
	svc := service.NewProductService(fp)

	_, _, _, err := svc.Search(context.Background(), domain.SearchParams{Page: 1, PageSize: 10})
	if err == nil {
		t.Fatalf("expected error")
	}
}
