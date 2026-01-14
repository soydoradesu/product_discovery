package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"regexp"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/soydoradesu/product_discovery/internal/domain"
	"github.com/soydoradesu/product_discovery/internal/repository"
)

type ProductRepo struct {
	pool *pgxpool.Pool
}

func NewProductRepo(pool *pgxpool.Pool) repository.ProductRepository {
	return &ProductRepo{pool: pool}
}

func (r *ProductRepo) GetByID(ctx context.Context, id int64) (domain.Product, error) {
	var p domain.Product
	var price float64

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
	p.Price = price

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

func (r *ProductRepo) Search(ctx context.Context, params domain.SearchParams) ([]domain.ProductSummary, int64, error) {
	q := strings.TrimSpace(params.Q)
	tsq := buildPrefixTSQuery(q)

	// $1 = q, $2 = category array
	args := []any{tsq, params.CategoryID}
	where := "WHERE ($1 = '' OR p.search_vector @@ to_tsquery('simple', $1))"

	// category multi-value: match ANY selected category
	where += `
	AND (
		COALESCE(array_length($2::bigint[], 1), 0) = 0
		OR EXISTS (
			SELECT 1
			FROM product_categories pc2
			WHERE pc2.product_id = p.id
			AND pc2.category_id = ANY($2::bigint[])
		)
	)`

	idx := 3
	if params.MinPrice != nil {
		where += fmt.Sprintf(" AND p.price >= $%d", idx)
		args = append(args, *params.MinPrice)
		idx++
	}
	if params.MaxPrice != nil {
		where += fmt.Sprintf(" AND p.price <= $%d", idx)
		args = append(args, *params.MaxPrice)
		idx++
	}
	if params.InStock != nil {
		where += fmt.Sprintf(" AND p.in_stock = $%d", idx)
		args = append(args, *params.InStock)
		idx++
	}

	// count for distinct products
	countSQL := `
		SELECT COUNT(DISTINCT p.id)
		FROM products p
	` + where

	var total int64
	if err := r.pool.QueryRow(ctx, countSQL, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// sorting
	orderBy := "ORDER BY p.id ASC"
	method := "DESC"
	if params.Method == "asc" {
		method = "ASC"
	}

	switch params.Sort {
	case "relevance":
		orderBy = fmt.Sprintf("ORDER BY rank %s, p.id ASC", method)
	case "price":
		orderBy = fmt.Sprintf("ORDER BY p.price %s, p.id ASC", method)
	case "created_at":
		orderBy = fmt.Sprintf("ORDER BY p.created_at %s, p.id ASC", method)
	case "rating":
		orderBy = fmt.Sprintf("ORDER BY p.rating %s, p.id ASC", method)
	}

	limit := params.PageSize
	offset := (params.Page - 1) * params.PageSize

	itemsSQL := `
	SELECT
		p.id,
		p.name,
		p.price,
		p.rating,
		p.in_stock,
		p.created_at,
		(SELECT url FROM product_images pi WHERE pi.product_id = p.id ORDER BY pi.position ASC LIMIT 1) AS thumbnail,
		COALESCE(
			jsonb_agg(DISTINCT jsonb_build_object('id', c.id, 'name', c.name))
			FILTER (WHERE c.id IS NOT NULL),
			'[]'::jsonb
		) AS categories_json,
		CASE
			WHEN $1 <> '' THEN ts_rank_cd(p.search_vector, to_tsquery('simple', $1))
			ELSE 0
		END AS rank
	FROM products p
	LEFT JOIN product_categories pc ON pc.product_id = p.id
	LEFT JOIN categories c ON c.id = pc.category_id
	` + where + `
	GROUP BY p.id
	` + orderBy + `
	LIMIT ` + fmt.Sprintf("%d", limit) + ` OFFSET ` + fmt.Sprintf("%d", offset)

	rows, err := r.pool.Query(ctx, itemsSQL, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var out []domain.ProductSummary
	for rows.Next() {
		var (
			ps domain.ProductSummary
			price float64
			thumb *string
			catsJSON []byte
			rank float64
		)

		if err := rows.Scan(
			&ps.ID,
			&ps.Name,
			&price,
			&ps.Rating,
			&ps.InStock,
			&ps.CreatedAt,
			&thumb,
			&catsJSON,
			&rank,
		); err != nil {
			return nil, 0, err
		}

		ps.Price = price
		ps.Thumbnail = thumb

		var cats []domain.Category
		if err := json.Unmarshal(catsJSON, &cats); err != nil {
			return nil, 0, err
		}
		ps.Categories = cats

		out = append(out, ps)
	}
	if rows.Err() != nil {
		return nil, 0, rows.Err()
	}

	return out, total, nil
}

var tsTokenRe = regexp.MustCompile(`[A-Za-z0-9]+`)

func buildPrefixTSQuery(input string) string {
	tokens := tsTokenRe.FindAllString(strings.ToLower(input), -1)
	if len(tokens) == 0 {
		return ""
	}

	parts := make([]string, 0, len(tokens))
	for _, t := range tokens {
		if t == "" {
			continue
		}
		parts = append(parts, t+":*")
	}

	return strings.Join(parts, " & ")
}