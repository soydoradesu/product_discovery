package main

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/soydoradesu/product_discovery/internal/config"
	"github.com/soydoradesu/product_discovery/internal/db"
	"github.com/soydoradesu/product_discovery/internal/seed"
)

func main() {
	log.SetFlags(log.LstdFlags | log.LUTC)

	cfg := config.Load()

	users := getenvInt("SEED_USERS", 1000)
	products := getenvInt("SEED_PRODUCTS", 1000)
	randomSeed := getenvInt64("SEED_RANDOM_SEED", 42)

	ctx := context.Background()

	pool, err := db.Connect(ctx, cfg.PostgresDSN())
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}
	defer pool.Close()

	if err := db.ApplyMigrations(ctx, pool, "./migrations"); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	opts := seed.Options{
		Users:      users,
		Products:   products,
		RandomSeed: randomSeed,
	}

	ctxSeed, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()

	if err := seed.Run(ctxSeed, pool, opts); err != nil {
		log.Fatalf("seed: %v", err)
	}

	log.Printf("seed done (users=%d, products=%d)", users, products)
}

func getenvInt(k string, def int) int {
	v := os.Getenv(k)
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return n
}

func getenvInt64(k string, def int64) int64 {
	v := os.Getenv(k)
	if v == "" {
		return def
	}
	n, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return def
	}
	return n
}
