package api

import (
	"database/sql"
	"net/http"

	"tis-euprava/mup-gradjani/internal/api/handlers"
)

func RegisterRoutes(mux *http.ServeMux, db *sql.DB) {
	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	// Requests collection: GET list, POST create
	mux.HandleFunc("/api/requests", handlers.Requests(db))
}
