package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/soydoradesu/product_discovery/internal/auth"
	"github.com/soydoradesu/product_discovery/internal/config"
	"github.com/soydoradesu/product_discovery/internal/http/middleware"
	"github.com/soydoradesu/product_discovery/internal/http/respond"
	"github.com/soydoradesu/product_discovery/internal/service"
)

type AuthHandlers struct {
	Cfg  config.Config
	Auth *service.AuthService
}

type loginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type okResp struct {
	OK bool `json:"ok"`
}

type meResp struct {
	UserID int64 `json:"userId"`
	Email  string `json:"email"`
}

func (h *AuthHandlers) Login(w http.ResponseWriter, r *http.Request) {
	var req loginReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respond.Fail(w, http.StatusBadRequest, "BAD_REQUEST", "invalid json")
		return
	}

	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	req.Password = strings.TrimSpace(req.Password)

	if req.Email == "" || req.Password == "" {
		respond.Fail(w, http.StatusBadRequest, "VALIDATION_ERROR", "email and password are required")
		return
	}

	userID, err := h.Auth.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		switch err {
		case service.ErrInvalidCredentials, service.ErrUserNotFound:
			respond.Fail(w, http.StatusUnauthorized, "INVALID_CREDENTIALS", "email or password is incorrect")
			return
		default:
			respond.Fail(w, http.StatusInternalServerError, "INTERNAL", "something went wrong")
			return
		}
	}

	token, err := auth.SignJWT(h.Cfg.JWTSecret, userID, 7*24*time.Hour)
	if err != nil {
		respond.Fail(w, http.StatusInternalServerError, "INTERNAL", "failed to create session")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   h.Cfg.CookieSecure,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int((7 * 24 * time.Hour).Seconds()),
	})

	respond.JSON(w, http.StatusOK, okResp{OK: true})
}

func (h *AuthHandlers) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   h.Cfg.CookieSecure,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})
	respond.JSON(w, http.StatusOK, okResp{OK: true})
}

func (h *AuthHandlers) Me(w http.ResponseWriter, r *http.Request) {
	uid, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		respond.Fail(w, http.StatusUnauthorized, "UNAUTHORIZED", "missing session")
		return
	}

	u, err := h.Auth.Users.GetByID(r.Context(), uid)
	if err != nil {
		respond.Fail(w, http.StatusUnauthorized, "UNAUTHORIZED", "invalid session")
		return
	}

	respond.JSON(w, http.StatusOK, meResp{UserID: uid, Email: u.Email})
}