package config

import (
	"os"
	"strconv"
)

type Config struct {
	BackendAddr string

	PostgresHost string
	PostgresPort int
	PostgresUser string
	PostgresPass string
	PostgresDB   string

	JWTSecret    string
	CookieSecure bool
	FrontendURL  string

	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string
}

func Load() Config {
	return Config{
		BackendAddr: getenv("BACKEND_ADDR", ":8080"),

		PostgresHost: getenv("POSTGRES_HOST", "localhost"),
		PostgresPort: getenvInt("POSTGRES_PORT", 5432),
		PostgresUser: getenv("POSTGRES_USER", "app"),
		PostgresPass: getenv("POSTGRES_PASSWORD", "app"),
		PostgresDB:   getenv("POSTGRES_DB", "productdb"),

		JWTSecret:    getenv("JWT_SECRET", "change-me"),
		CookieSecure: getenvBool("COOKIE_SECURE", false),
		FrontendURL:  getenv("FRONTEND_URL", "http://localhost:5173"),

		GoogleClientID:     getenv("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret: getenv("GOOGLE_CLIENT_SECRET", ""),
		GoogleRedirectURL:  getenv("GOOGLE_REDIRECT_URL", "http://localhost:8080/api/auth/google/callback"),
	}
}

func (c Config) PostgresDSN() string {
	return "postgres://" + c.PostgresUser + ":" + c.PostgresPass + "@" + c.PostgresHost +
		":" + strconv.Itoa(c.PostgresPort) + "/" + c.PostgresDB + "?sslmode=disable"
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
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

func getenvBool(k string, def bool) bool {
	v := os.Getenv(k)
	if v == "" {
		return def
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return def
	}
	return b
}