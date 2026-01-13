package httpapi

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"

	"github.com/soydoradesu/product_discovery/internal/http/respond"
)

func NewRouter() http.Handler {
	r := chi.NewRouter()
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(chimw.Recoverer)
	r.Use(chimw.Compress(5))
	r.Use(chimw.Timeout(15 * time.Second))

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) { 
		w.WriteHeader(http.StatusOK) 
	})
	r.Get("/debug/fail", func(w http.ResponseWriter, r *http.Request) {
    	respond.Fail(w, 400, "bad_request", "invalid input")
	})


	// todo: api routes

	return r
}
