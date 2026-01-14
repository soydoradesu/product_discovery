package postgres

import (
	"context"
	"errors"

	"github.com/soydoradesu/product_discovery/internal/domain"
	"github.com/soydoradesu/product_discovery/internal/repository"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductRepo struct {
	pool *pgxpool.Pool
}

func NewProductRepo(pool *pgxpool.Pool) repository.ProductRepository {
	return &ProductRepo{pool: pool}
}

func (r *ProductRepo) GetByID(ctx context.Context, id int64) (domain.Product, error) {
	var p domain.Product
	var price int64

	err := r.pool.QueryRow(ctx, `
		SELECT id, name, price, description, rating, in_stock, created_at
		FROM products
		WHERE id = $1
	`, id).Scan(&p.ID, &p.Name, &price, &p.Description, &p.Rating, &p.InStock, &p.CreatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Product{}, repository.ErrNotFound
	}
	if err != nil {
		return domain.Product{}, err
	}
	p.Price = float64(price)

	// images
	rows, err := r.pool.Query(ctx, `
		SELECT url, position
		FROM product_images
		WHERE product_id = $1
		ORDER BY position ASC
	`, id)
	if err != nil {
		return domain.Product{}, err
	}
	defer rows.Close()

	var imgs []domain.ProductImage
	for rows.Next() {
		var img domain.ProductImage
		if err := rows.Scan(&img.URL, &img.Position); err != nil {
			return domain.Product{}, err
		}
		imgs = append(imgs, img)
	}
	if rows.Err() != nil {
		return domain.Product{}, rows.Err()
	}
	p.Images = imgs

	// categories
	crows, err := r.pool.Query(ctx, `
		SELECT c.id, c.name
		FROM categories c
		JOIN product_categories pc ON pc.category_id = c.id
		WHERE pc.product_id = $1
		ORDER BY c.id ASC
	`, id)
	if err != nil {
		return domain.Product{}, err
	}
	defer crows.Close()

	var cats []domain.Category
	for crows.Next() {
		var c domain.Category
		if err := crows.Scan(&c.ID, &c.Name); err != nil {
			return domain.Product{}, err
		}
		cats = append(cats, c)
	}
	if crows.Err() != nil {
		return domain.Product{}, crows.Err()
	}
	p.Categories = cats

	return p, nil
}
