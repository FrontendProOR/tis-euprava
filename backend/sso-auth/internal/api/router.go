package api

import (
	"net/http"

	"tis-euprava/sso-auth/internal/api/handlers"
)

func NewRouter() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", handlers.Health)

	// Auth
	mux.HandleFunc("/auth/register", handlers.Register)
	mux.HandleFunc("/auth/login", handlers.Login)
	mux.HandleFunc("/auth/logout", handlers.Logout)
	mux.HandleFunc("/auth/refresh", handlers.Refresh)

	// Token + profile
	mux.HandleFunc("/auth/token/validate", handlers.ValidateToken)
	mux.HandleFunc("/auth/me", handlers.Me)

	return handlers.WithCORS(mux)
}