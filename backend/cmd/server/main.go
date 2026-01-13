package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/soydoradesu/product_discovery/internal/config"
	"github.com/soydoradesu/product_discovery/internal/db"
	httpapi "github.com/soydoradesu/product_discovery/internal/http"

)

func main() {
	log.SetFlags(log.LstdFlags | log.LUTC)

	cfg := config.Load()

	ctx := context.Background()
	pool, err := db.Connect(ctx, cfg.PostgresDSN())
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}
	defer pool.Close()

	if err := db.ApplyMigrations(ctx, pool, "./migrations"); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	r := httpapi.NewRouter(cfg, pool)

	srv := &http.Server{
		Addr: cfg.BackendAddr,
		Handler: r,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout: 60 * time.Second,
	}

	log.Printf("listening on %s", cfg.BackendAddr)
	log.Fatal(srv.ListenAndServe())
}
