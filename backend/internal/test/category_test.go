package internal_test

import (
	"context"
	"testing"

	"github.com/soydoradesu/product_discovery/internal/domain"
	"github.com/soydoradesu/product_discovery/internal/service"
)

type fakeCategories struct {
	items []domain.Category
	err   error
}

func (f *fakeCategories) List(ctx context.Context) ([]domain.Category, error) {
	return f.items, f.err
}

func TestCategoryService_List_OK(t *testing.T) {
	svc := service.NewCategoryService(&fakeCategories{
		items: []domain.Category{{ID: 1, Name: "Laptop"}, {ID: 2, Name: "Phone"}},
	})

	got, err := svc.List(context.Background())
	if err != nil {
		t.Fatalf("expected nil err, got %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}
