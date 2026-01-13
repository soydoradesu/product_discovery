package httpapi

import (
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"

	"github.com/soydoradesu/product_discovery/internal/http/respond"
	"github.com/soydoradesu/product_discovery/internal/config"
	"github.com/soydoradesu/product_discovery/internal/http/handlers"
	"github.com/soydoradesu/product_discovery/internal/http/middleware"
	"github.com/soydoradesu/product_discovery/internal/repository/postgres"
	"github.com/soydoradesu/product_discovery/internal/service"
)

func NewRouter(cfg config.Config, pool *pgxpool.Pool) http.Handler {
	userRepo := postgres.NewUserRepo(pool)
	authSvc := service.NewAuthService(userRepo)
	authH := &handlers.AuthHandlers{Cfg: cfg, Auth: authSvc}


	r := chi.NewRouter()
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(chimw.Recoverer)
	r.Use(chimw.Compress(5))
	r.Use(chimw.Timeout(15 * time.Second))
	r.Use(middleware.CORS(cfg))

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) { 
		w.WriteHeader(http.StatusOK) 
	})
	r.Get("/debug/fail", func(w http.ResponseWriter, r *http.Request) {
    	respond.Fail(w, 400, "bad_request", "invalid input")
	})

	r.Route("/api", func(api chi.Router) {
		api.Route("/auth", func(ar chi.Router) {
			ar.Post("/login", authH.Login)
			ar.Post("/logout", authH.Logout)
		})

		api.With(middleware.RequireAuth(cfg)).Get("/me", authH.Me)
	})
	return r
}