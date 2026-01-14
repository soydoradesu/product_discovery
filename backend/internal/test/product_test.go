package internal_test

import (
	"context"
	"testing"
	"time"

	"github.com/soydoradesu/product_discovery/internal/domain"
	"github.com/soydoradesu/product_discovery/internal/repository"
	"github.com/soydoradesu/product_discovery/internal/service"
)

type fakeProducts struct {
	byID map[int64]domain.Product
}

func (f *fakeProducts) GetByID(ctx context.Context, id int64) (domain.Product, error) {
	p, ok := f.byID[id]
	if !ok {
		return domain.Product{}, repository.ErrNotFound
	}
	return p, nil
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
