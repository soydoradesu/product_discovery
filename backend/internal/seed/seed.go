package seed

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/soydoradesu/product_discovery/internal/service"
)

type Options struct {
	Users int
	Products int
	RandomSeed int64
}

func Run(ctx context.Context, pool *pgxpool.Pool, opt Options) error {
	if opt.Users <= 0 {
		opt.Users = 1000
	}
	if opt.Products <= 0 {
		opt.Products = 1000
	}
	if opt.RandomSeed == 0 {
		opt.RandomSeed = 42
	}

	rng := rand.New(rand.NewSource(opt.RandomSeed))

	var userCount int64
	if err := pool.QueryRow(ctx, `SELECT COUNT(*) FROM users`).Scan(&userCount); err != nil {
		return err
	}
	var productCount int64
	if err := pool.QueryRow(ctx, `SELECT COUNT(*) FROM products`).Scan(&productCount); err != nil {
		return err
	}

	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	categoryIDs, err := ensureCategories(ctx, tx)
	if err != nil {
		return err
	}

	if userCount < int64(opt.Users) {
		if err := seedUsers(ctx, tx, rng, opt.Users-int(userCount)); err != nil {
			return err
		}
	}

	if productCount < int64(opt.Products) {
		if err := seedProducts(ctx, tx, rng, opt.Products-int(productCount), categoryIDs); err != nil {
			return err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}
	return nil
}

func ensureCategories(ctx context.Context, tx pgx.Tx) ([]int64, error) {
	base := []string{
		"Laptop", "Phone", "Audio", "Wearables", "Gaming",
		"Accessories", "Camera", "Networking", "Storage", "Home",
	}
	for _, name := range base {
		_, err := tx.Exec(ctx, `INSERT INTO categories(name) VALUES($1) ON CONFLICT (name) DO NOTHING`, name)
		if err != nil {
			return nil, err
		}
	}

	rows, err := tx.Query(ctx, `SELECT id FROM categories ORDER BY id ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return ids, nil
}

func seedUsers(ctx context.Context, tx pgx.Tx, rng *rand.Rand, n int) error {
	// demo user always present
	demoPass := "Password123!"
	hash, err := service.HashPassword(demoPass)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO users(email, password_hash)
		VALUES ($1, $2)
		ON CONFLICT (email) DO NOTHING
	`, "demo@example.com", hash)
	if err != nil {
		return err
	}

	reuseHash := hash

	rows := make([][]any, 0, n)
	for i := 0; i < n; i++ {
		email := fmt.Sprintf("user%04d_%d@example.com", i+1, rng.Intn(1_000_000))
		rows = append(rows, []any{email, reuseHash})
	}

	_, err = tx.CopyFrom(
		ctx,
		pgx.Identifier{"users"},
		[]string{"email", "password_hash"},
		pgx.CopyFromRows(rows),
	)
	if err != nil {
		for _, r := range rows {
			_, e := tx.Exec(ctx, `
				INSERT INTO users(email, password_hash)
				VALUES ($1, $2)
				ON CONFLICT (email) DO NOTHING
			`, r[0], r[1])
			if e != nil {
				return e
			}
		}
	}

	return nil
}

func seedProducts(ctx context.Context, tx pgx.Tx, rng *rand.Rand, n int, categoryIDs []int64) error {
	adjs := []string{"Ultra", "Pro", "Air", "Max", "Mini", "Prime", "Edge", "Nova", "Zen", "Core"}
	nouns := []string{"Speaker", "Headphones", "Laptop", "Phone", "Mouse", "Keyboard", "Router", "SSD", "Camera", "Monitor"}

	now := time.Now().UTC()

	for i := 0; i < n; i++ {
		name := fmt.Sprintf("%s %s %04d", adjs[rng.Intn(len(adjs))], nouns[rng.Intn(len(nouns))], i+1)
		desc := fmt.Sprintf("Seeded product %04d — %s for everyday use.", i+1, strings.ToLower(nouns[rng.Intn(len(nouns))]))

		rating := 2.5 + rng.Float64()*2.5
		inStock := rng.Intn(100) < 70 
		createdAt := now.Add(-time.Duration(rng.Intn(180*24)) * time.Hour)

		price := 10.0 + rng.Float64()*990.0 

		var productID int64
		err := tx.QueryRow(ctx, `
			INSERT INTO products(name, price, description, rating, in_stock, created_at)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING id
		`, name, price, desc, rating, inStock, createdAt).Scan(&productID)
		if err != nil {
			return err
		}

		// images: 1..4
		imgN := 1 + rng.Intn(4)
		for pos := 1; pos <= imgN; pos++ {
			url := fmt.Sprintf("https://picsum.photos/seed/p%d_%d/800/600", productID, pos)
			_, err := tx.Exec(ctx, `
				INSERT INTO product_images(product_id, url, position)
				VALUES ($1, $2, $3)
				ON CONFLICT (product_id, position) DO NOTHING
			`, productID, url, pos)
			if err != nil {
				return err
			}
		}

		cn := 1 + rng.Intn(3)
		seen := map[int64]bool{}
		for len(seen) < cn {
			cid := categoryIDs[rng.Intn(len(categoryIDs))]
			if seen[cid] {
				continue
			}
			seen[cid] = true
			_, err := tx.Exec(ctx, `
				INSERT INTO product_categories(product_id, category_id)
				VALUES ($1, $2)
				ON CONFLICT (product_id, category_id) DO NOTHING
			`, productID, cid)
			if err != nil {
				return err
			}
		}
	}

	return nil
}