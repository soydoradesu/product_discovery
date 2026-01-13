package middleware

import (
	"context"
	"net/http"

	"github.com/soydoradesu/product_discovery/internal/auth"
	"github.com/soydoradesu/product_discovery/internal/config"
	"github.com/soydoradesu/product_discovery/internal/http/respond"
)

type ctxKey string

const userIDKey ctxKey = "userID"

func UserIDFromContext(ctx context.Context) (int64, bool) {
	v := ctx.Value(userIDKey)
	id, ok := v.(int64)
	return id, ok
}

func RequireAuth(cfg config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c, err := r.Cookie("session")
			if err != nil || c.Value == "" {
				respond.Fail(w, http.StatusUnauthorized, "UNAUTHORIZED", "missing session")
				return
			}
			
			claims, err := auth.VerifyJWT(cfg.JWTSecret, c.Value)
			if err != nil {
				respond.Fail(w, http.StatusUnauthorized, "UNAUTHORIZED", "invalid session")
				return
			}
			ctx := context.WithValue(r.Context(), userIDKey, claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}