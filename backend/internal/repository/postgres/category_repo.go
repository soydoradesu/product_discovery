package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/soydoradesu/product_discovery/internal/domain"
	"github.com/soydoradesu/product_discovery/internal/repository"
)

type CategoryRepo struct {
	pool *pgxpool.Pool
}

func NewCategoryRepo(pool *pgxpool.Pool) repository.CategoryRepository {
	return &CategoryRepo{pool: pool}
}

func (r *CategoryRepo) List(ctx context.Context) ([]domain.Category, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, name
		FROM categories
		ORDER BY name ASC, id ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.Category
	for rows.Next() {
		var c domain.Category
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return out, nil
}