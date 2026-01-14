package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
	"errors"
	"crypto/rand"
	"encoding/base64"

	"github.com/soydoradesu/product_discovery/internal/auth"
	"github.com/soydoradesu/product_discovery/internal/config"
	"github.com/soydoradesu/product_discovery/internal/http/middleware"
	"github.com/soydoradesu/product_discovery/internal/http/respond"
	"github.com/soydoradesu/product_discovery/internal/service"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
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

const oauthStateCookie = "oauth_state"

func (h *AuthHandlers) oauthConfig() (*oauth2.Config, error) {
	if strings.TrimSpace(h.Cfg.GoogleClientID) == "" || strings.TrimSpace(h.Cfg.GoogleClientSecret) == "" {
		return nil, errors.New("google oauth not configured")
	}
	return &oauth2.Config{
		ClientID: h.Cfg.GoogleClientID,
		ClientSecret: h.Cfg.GoogleClientSecret,
		RedirectURL: h.Cfg.GoogleRedirectURL,
		Endpoint: google.Endpoint,
		Scopes: []string{"email", "profile"},
	}, nil
}

func (h *AuthHandlers) setSessionCookie(w http.ResponseWriter, userID int64) error {
	token, err := auth.SignJWT(h.Cfg.JWTSecret, userID, 7*24*time.Hour)
	if err != nil {
		return err
	}
	http.SetCookie(w, &http.Cookie{
		Name: "session",
		Value: token,
		Path: "/",
		HttpOnly: true,
		Secure: h.Cfg.CookieSecure,
		SameSite: http.SameSiteLaxMode,
		MaxAge: int((7 * 24 * time.Hour).Seconds()),
	})
	return nil
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
		Name: "session",
		Value: token,
		Path: "/",
		HttpOnly: true,
		Secure: h.Cfg.CookieSecure,
		SameSite: http.SameSiteLaxMode,
		MaxAge: int((7 * 24 * time.Hour).Seconds()),
	})

	respond.JSON(w, http.StatusOK, okResp{OK: true})
}

func (h *AuthHandlers) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name: "session",
		Value: "",
		Path: "/",
		HttpOnly: true,
		Secure: h.Cfg.CookieSecure,
		SameSite: http.SameSiteLaxMode,
		MaxAge: -1,
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

func (h *AuthHandlers) GoogleStart(w http.ResponseWriter, r *http.Request) {
	ocfg, err := h.oauthConfig()
	if err != nil {
		respond.Fail(w, http.StatusInternalServerError, "CONFIG_ERROR", "google oauth is not configured")
		return
	}

	state := make([]byte, 32)
	if _, err := rand.Read(state); err != nil {
		respond.Fail(w, http.StatusInternalServerError, "INTERNAL", "failed to start oauth")
		return
	}
	stateStr := base64.RawURLEncoding.EncodeToString(state)

	http.SetCookie(w, &http.Cookie{
		Name: oauthStateCookie,
		Value: stateStr,
		Path: "/api/auth/google",
		HttpOnly: true,
		Secure: h.Cfg.CookieSecure,
		SameSite: http.SameSiteLaxMode,
		MaxAge: int((10 * time.Minute).Seconds()),
	})

	url := ocfg.AuthCodeURL(stateStr, oauth2.AccessTypeOnline)
	http.Redirect(w, r, url, http.StatusFound)
}

type googleUserInfo struct {
	ID string `json:"id"`
	Email string `json:"email"`
	VerifiedEmail bool `json:"verified_email"`
}

// GET /api/auth/google/callback
func (h *AuthHandlers) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	ocfg, err := h.oauthConfig()
	if err != nil {
		respond.Fail(w, http.StatusInternalServerError, "CONFIG_ERROR", "google oauth is not configured")
		return
	}

	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")
	if code == "" || state == "" {
		respond.Fail(w, http.StatusBadRequest, "BAD_REQUEST", "missing code/state")
		return
	}

	c, err := r.Cookie(oauthStateCookie)
	if err != nil || c.Value == "" || c.Value != state {
		respond.Fail(w, http.StatusUnauthorized, "UNAUTHORIZED", "invalid oauth state")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name: oauthStateCookie,
		Value: "",
		Path: "/api/auth/google",
		HttpOnly: true,
		Secure: h.Cfg.CookieSecure,
		SameSite: http.SameSiteLaxMode,
		MaxAge: -1,
	})

	tok, err := ocfg.Exchange(r.Context(), code)
	if err != nil {
		respond.Fail(w, http.StatusUnauthorized, "UNAUTHORIZED", "oauth exchange failed")
		return
	}

	client := ocfg.Client(r.Context(), tok)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		respond.Fail(w, http.StatusUnauthorized, "UNAUTHORIZED", "failed to fetch userinfo")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respond.Fail(w, http.StatusUnauthorized, "UNAUTHORIZED", "userinfo request failed")
		return
	}

	var ui googleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&ui); err != nil {
		respond.Fail(w, http.StatusInternalServerError, "INTERNAL", "failed to parse userinfo")
		return
	}

	if ui.Email == "" || ui.ID == "" {
		respond.Fail(w, http.StatusUnauthorized, "UNAUTHORIZED", "invalid userinfo")
		return
	}
	if !ui.VerifiedEmail {
		respond.Fail(w, http.StatusUnauthorized, "UNAUTHORIZED", "email not verified")
		return
	}

	userID, err := h.Auth.OAuthLogin(r.Context(), ui.Email, ui.ID)
	if err != nil {
		if err == service.ErrOAuthAccountConflict {
			respond.Fail(w, http.StatusConflict, "CONFLICT", "account already linked to different google id")
			return
		}
		respond.Fail(w, http.StatusInternalServerError, "INTERNAL", "something went wrong")
		return
	}

	if err := h.setSessionCookie(w, userID); err != nil {
		respond.Fail(w, http.StatusInternalServerError, "INTERNAL", "failed to create session")
		return
	}

	// redirect to frontend
	http.Redirect(w, r, h.Cfg.FrontendURL+"/?oauth=success", http.StatusFound)
}