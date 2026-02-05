package main

import (
	"log"
	"net/http"
	"time"

	"tis-euprava/mup-gradjani/internal/api"
	"tis-euprava/mup-gradjani/internal/config"
)

func main() {
	log.Println("MUP service started")

	// 1) Load config (.env)
	cfg := config.LoadConfig()

	// 2) Open DB connection
	db, err := config.OpenDB(cfg)
	if err != nil {
		log.Fatalf("Greška pri otvaranju DB konekcije: %v", err)
	}
	defer db.Close()

	log.Println("Povezano na PostgreSQL bazu uspešno!")
	log.Printf("Servis '%s' sluša na portu %s\n", cfg.ServiceName, cfg.HTTPPort)

	// 3) Router (ServeMux)
	mux := http.NewServeMux()
	api.RegisterRoutes(mux, db)

	// 4) HTTP server )
	srv := &http.Server{
		Addr:              ":" + cfg.HTTPPort,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
