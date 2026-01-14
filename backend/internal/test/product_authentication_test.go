package internal_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/soydoradesu/product_discovery/internal/auth"
	"github.com/soydoradesu/product_discovery/internal/config"
	"github.com/soydoradesu/product_discovery/internal/http/middleware"
)

func TestRequireAuth_Unauthorized_WhenMissingCookie(t *testing.T) {
	cfg := config.Config{JWTSecret: "test-secret"}

	r := chi.NewRouter()
	r.Use(middleware.RequireAuth(cfg))
	r.Get("/api/products/1", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/api/products/1", nil)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 got %d body=%s", rr.Code, rr.Body.String())
	}
}

func TestRequireAuth_Unauthorized_WhenCookieEmpty(t *testing.T) {
	cfg := config.Config{JWTSecret: "test-secret"}

	r := chi.NewRouter()
	r.Use(middleware.RequireAuth(cfg))
	r.Get("/api/products/1", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/api/products/1", nil)
	req.AddCookie(&http.Cookie{Name: "session", Value: ""})
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 got %d body=%s", rr.Code, rr.Body.String())
	}
}

func TestRequireAuth_Unauthorized_WhenTokenInvalid(t *testing.T) {
	cfg := config.Config{JWTSecret: "test-secret"}

	r := chi.NewRouter()
	r.Use(middleware.RequireAuth(cfg))
	r.Get("/api/products/1", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/api/products/1", nil)
	req.AddCookie(&http.Cookie{Name: "session", Value: "not-a-jwt"})
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 got %d body=%s", rr.Code, rr.Body.String())
	}
}

func TestRequireAuth_OK_WhenTokenValid_AndUserIDInContext(t *testing.T) {
	cfg := config.Config{JWTSecret: "test-secret"}

	token, err := auth.SignJWT(cfg.JWTSecret, 123, 10*time.Minute)
	if err != nil {
		t.Fatalf("sign jwt: %v", err)
	}

	r := chi.NewRouter()
	r.Use(middleware.RequireAuth(cfg))
	r.Get("/api/products/1", func(w http.ResponseWriter, r *http.Request) {
		uid, ok := middleware.UserIDFromContext(r.Context())
		if !ok || uid != 123 {
			t.Fatalf("expected uid=123 in context, got ok=%v uid=%d", ok, uid)
		}
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/api/products/1", nil)
	req.AddCookie(&http.Cookie{Name: "session", Value: token})
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d body=%s", rr.Code, rr.Body.String())
	}
}

func TestRequireAuth_Unauthorized_WhenTokenExpired(t *testing.T) {
	cfg := config.Config{JWTSecret: "test-secret"}

	token, err := auth.SignJWT(cfg.JWTSecret, 123, -1*time.Minute) // already expired
	if err != nil {
		t.Fatalf("sign jwt: %v", err)
	}

	r := chi.NewRouter()
	r.Use(middleware.RequireAuth(cfg))
	r.Get("/api/products/1", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/api/products/1", nil)
	req.AddCookie(&http.Cookie{Name: "session", Value: token})
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 got %d body=%s", rr.Code, rr.Body.String())
	}
}
